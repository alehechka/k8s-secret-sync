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

	if isManagedBy(secret) {
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
				client.CreateUpdateSecret(rule.Spec.Rules, &namespace, secret)
			}
		}
	}

	return nil
}

func (client *Client) CreateUpdateSecret(rules typesv1.Rules, namespace *v1.Namespace, secret *v1.Secret) error {
	logger := secretLogger(prepareSecret(namespace, secret))

	if namespaceSecret, err := client.GetSecret(namespace.Name, secret.Name); err == nil {
		logger.Debugf("already exists")

		if !rules.Force && !isManagedBy(namespaceSecret) {
			logger.Debugf("existing secret is not managed and will not be force updated")
			return nil
		}

		if isManagedBy(namespaceSecret) && secretsAreEqual(secret, namespaceSecret) {
			logger.Debugf("existing secret contains same data")
			return nil
		}

		return client.UpdateSecret(namespace, secret)
	}

	return client.CreateSecret(namespace, secret)
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
	logger := secretLogger(prepareSecret(namespace, secret))

	if namespaceSecret, err := client.GetSecret(namespace.Name, secret.Name); err == nil {
		if rules.Force || isManagedBy(namespaceSecret) {
			return client.DeleteSecret(namespace, secret)
		}

		logger.Debugf("existing secret is not managed and will not be force deleted")
	}

	return nil
}

func (client *Client) CreateSecret(namespace *v1.Namespace, secret *v1.Secret) error {
	newSecret := prepareSecret(namespace, secret)

	logger := secretLogger(newSecret)
	logger.Infof("creating secret")

	_, err := client.DefaultClientset.CoreV1().Secrets(namespace.Name).Create(client.Context, newSecret, metav1.CreateOptions{})

	if err != nil {
		logger.Errorf("failed to create secret - %s", err.Error())
	}

	return err
}

func (client *Client) UpdateSecret(namespace *v1.Namespace, secret *v1.Secret) (err error) {
	updateSecret := prepareSecret(namespace, secret)

	logger := secretLogger(updateSecret)
	logger.Infof("updating secret")

	_, err = client.DefaultClientset.CoreV1().Secrets(namespace.Name).Update(client.Context, updateSecret, metav1.UpdateOptions{})
	if err != nil {
		logger.Errorf("failed to update secret - %s", err.Error())
	}

	return
}

func (client *Client) DeleteSecret(namespace *v1.Namespace, secret *v1.Secret) (err error) {
	logger := secretLogger(prepareSecret(namespace, secret))

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
		secretLogger(&v1.Secret{ObjectMeta: metav1.ObjectMeta{Namespace: namespace, Name: name}}).
			Errorf("failed to get secret: %s", err.Error())
	}
	return
}

func (client *Client) ListSecrets(namespace string) (list *v1.SecretList, err error) {
	list, err = client.DefaultClientset.CoreV1().Secrets(namespace).List(client.Context, metav1.ListOptions{})
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
