package client

import (
	"context"
	"reflect"

	typesv1 "github.com/alehechka/kube-secret-sync/api/types/v1"
	"github.com/alehechka/kube-secret-sync/clientset"
	"github.com/alehechka/kube-secret-sync/constants"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
)

func secretEventHandler(ctx context.Context, event watch.Event) {
	secret := event.Object.(*v1.Secret)

	if isManagedBy(secret) {
		return
	}

	switch event.Type {
	case watch.Added:
		addedSecretHandler(ctx, secret)
	case watch.Modified:
		modifiedSecretHandler(ctx, secret)
	case watch.Deleted:
		deletedSecretHandler(ctx, secret)
	}
}

func addedSecretHandler(ctx context.Context, secret *v1.Secret) error {
	logger := secretLogger(secret)
	logger.Infof("added")

	if secret.CreationTimestamp.Time.Before(startTime) {
		logger.Debugf("secret will be synced on startup by SecretSyncRule watcher")
		return nil
	}

	return syncAddedModifiedSecret(ctx, secret)
}

func modifiedSecretHandler(ctx context.Context, secret *v1.Secret) error {
	if secret.DeletionTimestamp != nil {
		return nil
	}

	secretLogger(secret).Infof("modified")

	return syncAddedModifiedSecret(ctx, secret)
}

func syncAddedModifiedSecret(ctx context.Context, secret *v1.Secret) error {
	rules, err := listSecretSyncRules(ctx)
	if err != nil {
		return err
	}

	for _, rule := range rules.Items {
		if rule.ShouldSyncSecret(secret) {
			for _, namespace := range rule.Namespaces(ctx) {
				createUpdateSecret(ctx, rule.Spec.Rules, &namespace, secret)
			}
		}
	}

	return nil
}

func createUpdateSecret(ctx context.Context, rules typesv1.Rules, namespace *v1.Namespace, secret *v1.Secret) error {
	logger := secretLogger(prepareSecret(namespace, secret))

	if namespaceSecret, err := getSecret(ctx, namespace.Name, secret.Name); err == nil {
		logger.Debugf("already exists")

		if !rules.Force && !isManagedBy(namespaceSecret) {
			logger.Debugf("existing secret is not managed and will not be force updated")
			return nil
		}

		if isManagedBy(namespaceSecret) && secretsAreEqual(secret, namespaceSecret) {
			logger.Debugf("existing secret contains same data")
			return nil
		}

		return updateSecret(ctx, namespace, secret)
	}

	return createSecret(ctx, namespace, secret)
}

// TODO - rebuild this function
func deletedSecretHandler(ctx context.Context, secret *v1.Secret) error {
	secretLogger(secret).Infof("deleted")

	return nil
}

func syncDeletedSecret(ctx context.Context, rules typesv1.Rules, namespace *v1.Namespace, secret *v1.Secret) error {
	if namespaceSecret, err := getSecret(ctx, namespace.Name, secret.Name); err == nil {
		if rules.Force || isManagedBy(namespaceSecret) {
			return deleteSecret(ctx, namespace, secret)
		}
	}

	secretLogger(secret).Debugf("not found for deletion")
	return nil
}

func createSecret(ctx context.Context, namespace *v1.Namespace, secret *v1.Secret) error {
	newSecret := prepareSecret(namespace, secret)

	logger := secretLogger(newSecret)
	logger.Infof("creating secret")

	_, err := clientset.Default.CoreV1().Secrets(namespace.Name).Create(ctx, newSecret, metav1.CreateOptions{})

	if err != nil {
		logger.Errorf("failed to create secret - %s", err.Error())
	}

	return err
}

func updateSecret(ctx context.Context, namespace *v1.Namespace, secret *v1.Secret) (err error) {
	updateSecret := prepareSecret(namespace, secret)

	logger := secretLogger(updateSecret)
	logger.Infof("updating secret")

	_, err = clientset.Default.CoreV1().Secrets(namespace.Name).Update(ctx, updateSecret, metav1.UpdateOptions{})
	if err != nil {
		logger.Errorf("failed to update secret - %s", err.Error())
	}

	return
}

func deleteSecret(ctx context.Context, namespace *v1.Namespace, secret *v1.Secret) (err error) {
	logger := secretLogger(prepareSecret(namespace, secret))

	logger.Infof("deleting secret")

	err = clientset.Default.CoreV1().Secrets(namespace.Name).Delete(ctx, secret.Name, metav1.DeleteOptions{})
	if err != nil {
		logger.Errorf("failed to delete secret - %s", err.Error())
	}

	return
}

func getSecret(ctx context.Context, namespace, name string) (secret *v1.Secret, err error) {
	secret, err = clientset.Default.CoreV1().Secrets(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		secretLogger(&v1.Secret{ObjectMeta: metav1.ObjectMeta{Namespace: namespace, Name: name}}).
			Errorf("does not exist to sync: %s", err.Error())
	}
	return
}

func listSecrets(ctx context.Context, namespace string) (list *v1.SecretList, err error) {
	list, err = clientset.Default.CoreV1().Secrets(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		namespaceLogger(&v1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: namespace}}).
			Errorf("failed to list secrets: %s", err.Error())
	}
	return
}

func secretsAreEqual(a, b *v1.Secret) bool {
	return (a.Type == b.Type &&
		reflect.DeepEqual(a.Data, b.Data) &&
		reflect.DeepEqual(a.StringData, b.StringData) &&
		annotationsAreEqual(a.Annotations, b.Annotations))
}

func annotationsAreEqual(a, b map[string]string) bool {
	if a == nil {
		a = make(map[string]string)
	}
	delete(a, constants.ManagedByAnnotationKey)
	delete(a, constants.LastAppliedConfigurationAnnotationKey)

	if b == nil {
		b = make(map[string]string)
	}
	delete(b, constants.ManagedByAnnotationKey)
	delete(b, constants.LastAppliedConfigurationAnnotationKey)

	return reflect.DeepEqual(a, b)
}

func prepareSecret(namespace *v1.Namespace, secret *v1.Secret) *v1.Secret {
	annotations := secret.Annotations
	if annotations == nil {
		annotations = make(map[string]string)
	}
	annotations[constants.ManagedByAnnotationKey] = constants.ManagedByAnnotationValue
	delete(annotations, constants.LastAppliedConfigurationAnnotationKey)

	return &v1.Secret{
		TypeMeta: secret.TypeMeta,
		ObjectMeta: metav1.ObjectMeta{
			Name:        secret.Name,
			Namespace:   namespace.Name,
			Labels:      secret.Labels,
			Annotations: annotations,
		},
		Immutable:  secret.Immutable,
		Data:       secret.Data,
		StringData: secret.StringData,
		Type:       secret.Type,
	}

}

func isManagedBy(secret *v1.Secret) bool {
	managedBy, ok := secret.Annotations[constants.ManagedByAnnotationKey]

	return ok && managedBy == constants.ManagedByAnnotationValue
}
