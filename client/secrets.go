package client

import (
	"reflect"

	typesv1 "github.com/alehechka/kube-secret-sync/api/types/v1"
	"github.com/alehechka/kube-secret-sync/constants"
	log "github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
)

func (client *Client) SecretEventHandler(event watch.Event) error {
	secret, ok := event.Object.(*v1.Secret)
	if !ok {
		log.Error("failed to cast Secret")
		return nil
	}

	if IsManagedBy(secret) {
		return nil
	}

	switch event.Type {
	case watch.Added:
		return client.AddedSecretHandler(secret)
	case watch.Modified:
		return client.ModifiedSecretHandler(secret)
	case watch.Deleted:
		return client.DeletedSecretHandler(secret)
	}

	return nil
}

func (client *Client) AddedSecretHandler(secret *v1.Secret) error {
	logger := secretLogger(secret)

	if secret.CreationTimestamp.Time.Before(client.StartTime) {
		logger.Debugf("secret will be synced on startup by SecretSyncRule watcher")
		return nil
	}

	logger.Infof("added")
	return client.SyncAddedModifiedSecret(secret)
}

func (client *Client) ModifiedSecretHandler(secret *v1.Secret) error {
	if secret.DeletionTimestamp != nil {
		return nil
	}

	secretLogger(secret).Infof("modified")
	return client.SyncAddedModifiedSecret(secret)
}

func (client *Client) SyncAddedModifiedSecret(secret *v1.Secret) error {
	rules, err := client.ListSecretSyncRules()
	if err != nil {
		return err
	}

	for _, rule := range rules.Items {
		if rule.ShouldSyncSecret(secret) {
			for _, namespace := range rule.Namespaces(client.Context, client.DefaultClientset) {
				client.CreateUpdateSecret(&rule, &namespace, secret)
			}
		}
	}

	return nil
}

func (client *Client) CreateUpdateSecret(rule *typesv1.SecretSyncRule, namespace *v1.Namespace, secret *v1.Secret) error {
	logger := secretLogger(secret, namespace)

	if namespaceSecret, err := client.GetSecret(namespace.Name, secret.Name); err == nil {
		logger.Debugf("already exists")

		if !rule.Spec.Rules.Force && !IsManagedBy(namespaceSecret) {
			logger.Debugf("existing secret is not managed and will not be force updated")
			return nil
		}

		if IsManagedBy(namespaceSecret) && SecretsAreEqual(secret, namespaceSecret) {
			logger.Debugf("existing secret contains same data")
			return nil
		}

		return client.UpdateSecret(rule, namespace, secret)
	}

	return client.CreateSecret(rule, namespace, secret)
}

func (client *Client) DeletedSecretHandler(secret *v1.Secret) error {
	secretLogger(secret).Infof("deleted")

	rules, err := client.ListSecretSyncRules()
	if err != nil {
		return err
	}

	for _, rule := range rules.Items {
		if rule.ShouldSyncSecret(secret) {
			for _, namespace := range rule.Namespaces(client.Context, client.DefaultClientset) {
				client.SyncDeletedSecret(rule.Spec.Rules, &namespace, secret)
			}
		}
	}

	return nil
}

func (client *Client) SyncDeletedSecret(rules typesv1.Rules, namespace *v1.Namespace, secret *v1.Secret) error {
	logger := secretLogger(secret, namespace)

	if namespaceSecret, err := client.GetSecret(namespace.Name, secret.Name); err == nil {
		if rules.Force || IsManagedBy(namespaceSecret) {
			return client.DeleteSecret(namespace, secret)
		}

		logger.Debugf("existing secret is not managed and will not be force deleted")
	}

	return nil
}

func (client *Client) CreateSecret(rule *typesv1.SecretSyncRule, namespace *v1.Namespace, secret *v1.Secret) error {
	newSecret := PrepareSecret(rule, namespace, secret)

	logger := secretLogger(newSecret)
	logger.Infof("creating secret")

	_, err := client.DefaultClientset.CoreV1().Secrets(namespace.Name).Create(client.Context, newSecret, metav1.CreateOptions{})

	if err != nil {
		logger.Errorf("failed to create secret - %s", err.Error())
	}

	return err
}

func (client *Client) UpdateSecret(rule *typesv1.SecretSyncRule, namespace *v1.Namespace, secret *v1.Secret) (err error) {
	updateSecret := PrepareSecret(rule, namespace, secret)

	logger := secretLogger(updateSecret)
	logger.Infof("updating secret")

	_, err = client.DefaultClientset.CoreV1().Secrets(namespace.Name).Update(client.Context, updateSecret, metav1.UpdateOptions{})
	if err != nil {
		logger.Errorf("failed to update secret - %s", err.Error())
	}

	return
}

func (client *Client) DeleteSecret(namespace *v1.Namespace, secret *v1.Secret) (err error) {
	logger := secretLogger(secret, namespace)

	logger.Infof("deleting secret")

	err = client.DefaultClientset.CoreV1().Secrets(namespace.Name).Delete(client.Context, secret.Name, metav1.DeleteOptions{})
	if err != nil {
		logger.Errorf("failed to delete secret - %s", err.Error())
	}

	return
}

func (client *Client) GetSecret(namespace, name string) (secret *v1.Secret, err error) {
	secret, err = client.DefaultClientset.CoreV1().Secrets(namespace).Get(client.Context, name, metav1.GetOptions{})
	if err != nil {
		secretNameLogger(namespace, name).
			Errorf("failed to get secret: %s", err.Error())
	}
	return
}

func (client *Client) ListSecrets(namespace string) (list *v1.SecretList, err error) {
	list, err = client.DefaultClientset.CoreV1().Secrets(namespace).List(client.Context, metav1.ListOptions{})
	if err != nil {
		namespaceNameLogger(namespace).
			Errorf("failed to list secrets: %s", err.Error())
	}
	return
}

func SecretsAreEqual(a, b *v1.Secret) bool {
	return (a.Type == b.Type &&
		reflect.DeepEqual(a.Data, b.Data) &&
		reflect.DeepEqual(a.StringData, b.StringData) &&
		AnnotationsAreEqual(a.Annotations, b.Annotations))
}

func PrepareSecret(rule *typesv1.SecretSyncRule, namespace *v1.Namespace, secret *v1.Secret) *v1.Secret {
	annotations := CopyAnnotations(secret.Annotations)
	annotations[constants.ManagedByAnnotationKey] = constants.ManagedByAnnotationValue

	return &v1.Secret{
		TypeMeta: secret.TypeMeta,
		ObjectMeta: metav1.ObjectMeta{
			Name:            secret.Name,
			Namespace:       namespace.Name,
			Labels:          secret.Labels,
			Annotations:     annotations,
			OwnerReferences: []metav1.OwnerReference{OwnerReference(rule)},
		},
		Immutable:  secret.Immutable,
		Data:       secret.Data,
		StringData: secret.StringData,
		Type:       secret.Type,
	}

}

func IsManagedBy(secret *v1.Secret) bool {
	managedBy, ok := secret.Annotations[constants.ManagedByAnnotationKey]

	return ok && managedBy == constants.ManagedByAnnotationValue
}
