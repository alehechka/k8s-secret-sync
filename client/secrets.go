package client

import (
	"context"
	"reflect"

	"github.com/alehechka/kube-secret-sync/constants"
	log "github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"

	"k8s.io/client-go/kubernetes"
)

func secretEventHandler(ctx context.Context, clientset *kubernetes.Clientset, config *SyncConfig, event watch.Event) {
	secret := event.Object.(*v1.Secret)

	switch event.Type {
	case watch.Added:
		addSecrets(ctx, clientset, config, secret)
	case watch.Modified:
		modifySecrets(ctx, clientset, config, secret)
	case watch.Deleted:
		deleteSecrets(ctx, clientset, config, secret)
	}
}

func addSecrets(ctx context.Context, clientset *kubernetes.Clientset, config *SyncConfig, secret *v1.Secret) error {
	log.Infof("[%s/%s]: Secret created", secret.ObjectMeta.Namespace, secret.ObjectMeta.Name)

	return syncNamespaceSecret(ctx, clientset, config, secret, syncAddedModifiedSecret)
}

func modifySecrets(ctx context.Context, clientset *kubernetes.Clientset, config *SyncConfig, secret *v1.Secret) error {
	if secret.DeletionTimestamp != nil {
		return nil
	}

	log.Infof("[%s/%s]: Secret modified", secret.ObjectMeta.Namespace, secret.ObjectMeta.Name)

	return syncNamespaceSecret(ctx, clientset, config, secret, syncAddedModifiedSecret)
}

func deleteSecrets(ctx context.Context, clientset *kubernetes.Clientset, config *SyncConfig, secret *v1.Secret) error {
	log.Infof("[%s/%s]: Secret deleted", secret.ObjectMeta.Namespace, secret.ObjectMeta.Name)

	return syncNamespaceSecret(ctx, clientset, config, secret, syncDeletedSecret)
}

type SecretSyncFunc func(context.Context, *kubernetes.Clientset, *SyncConfig, v1.Namespace, *v1.Secret) error

func syncNamespaceSecret(ctx context.Context, clientset *kubernetes.Clientset, config *SyncConfig, secret *v1.Secret, sync SecretSyncFunc) error {
	if config.ExcludeSecrets.IsExcluded(secret.Name) {
		log.Debugf("[%s/%s]: Secret is excluded from sync", secret.Namespace, secret.Name)
		return constants.ErrExcludedSecret
	}

	if !config.IncludeSecrets.IsIncluded(secret.Name) {
		log.Debugf("[%s/%s]: Secret is not included for sync", secret.Namespace, secret.Name)
		return constants.ErrNotIncludedSecret
	}

	namespaces, err := listNamespaces(ctx, clientset)
	if err != nil {
		log.Errorf("Failed to list namespaces: %s", err.Error())
		return err
	}

	for _, namespace := range namespaces.Items {
		if namespace.Name == config.SecretsNamespace {
			log.Debugf("[%s]: Skipping secrets namespace", namespace.Name)
			continue
		}

		if config.ExcludeNamespaces.IsExcluded(namespace.Name) {
			log.Debugf("[%s]: Namespace has been excluded from sync", namespace.Name)
			continue
		}

		if !config.IncludeNamespaces.IsIncluded(namespace.Name) {
			log.Debugf("[%s]: Namespace is not included for sync", namespace.Name)
			continue
		}

		sync(ctx, clientset, config, namespace, secret)
	}

	return nil
}

func syncAddedModifiedSecret(ctx context.Context, clientset *kubernetes.Clientset, config *SyncConfig, namespace v1.Namespace, secret *v1.Secret) error {
	if namespaceSecret, err := getSecret(ctx, clientset, namespace.Name, secret.Name); err == nil {
		log.Debugf("[%s/%s]: Secret already exists", namespace.Name, secret.Name)

		if !config.ForceSync && !isManagedBy(namespaceSecret) {
			log.Debugf("[%s/%s]: Existing secret is not managed and will not be force updated", namespace.Name, secret.Name)
			return nil
		}

		if isManagedBy(namespaceSecret) && secretsAreEqual(secret, namespaceSecret) {
			log.Debugf("[%s/%s]: Existing secret contains same data", namespace.Name, secret.Name)
			return nil
		}

		_, err = updateSecret(ctx, clientset, namespace, secret)
		return err
	}

	return createSecret(ctx, clientset, namespace, secret)
}

func syncDeletedSecret(ctx context.Context, clientset *kubernetes.Clientset, config *SyncConfig, namespace v1.Namespace, secret *v1.Secret) error {
	if namespaceSecret, err := getSecret(ctx, clientset, namespace.Name, secret.Name); err == nil {
		if config.ForceSync || isManagedBy(namespaceSecret) {
			return deleteSecret(ctx, clientset, namespace, secret)
		}
	}

	log.Debugf("[%s/%s]: Secret not found for deletion", namespace.Name, secret.Name)
	return nil
}

func createSecret(ctx context.Context, clientset *kubernetes.Clientset, namespace v1.Namespace, secret *v1.Secret) error {
	log.Infof("[%s/%s]: Creating secret", namespace.Name, secret.Name)

	newSecret := prepareSecret(namespace, secret)

	_, err := clientset.CoreV1().Secrets(namespace.Name).Create(ctx, newSecret, metav1.CreateOptions{})

	if err != nil {
		log.Errorf("[%s/%s]: Failed to create secret - %s", namespace.Name, secret.Name, err.Error())
	}

	return err
}

func deleteSecret(ctx context.Context, clientset *kubernetes.Clientset, namespace v1.Namespace, secret *v1.Secret) (err error) {
	log.Infof("[%s/%s]: Deleting secret", namespace.Name, secret.Name)

	err = clientset.CoreV1().Secrets(namespace.Name).Delete(ctx, secret.Name, metav1.DeleteOptions{})
	if err != nil {
		log.Errorf("[%s/%s]: Failed to delete secret - %s", namespace.Name, secret.Name, err.Error())
	}

	return
}

func updateSecret(ctx context.Context, clientset *kubernetes.Clientset, namespace v1.Namespace, secret *v1.Secret) (updated *v1.Secret, err error) {
	log.Infof("[%s/%s]: Updating secret", namespace.Name, secret.Name)

	updateSecret := prepareSecret(namespace, secret)

	updated, err = clientset.CoreV1().Secrets(namespace.Name).Update(ctx, updateSecret, metav1.UpdateOptions{})
	if err != nil {
		log.Errorf("[%s/%s]: Failed to update secret - %s", namespace.Name, secret.Name, err.Error())
	}

	return
}

func getSecret(ctx context.Context, clientset *kubernetes.Clientset, namespace, name string) (*v1.Secret, error) {
	return clientset.CoreV1().Secrets(namespace).Get(ctx, name, metav1.GetOptions{})
}

func secretsAreEqual(a, b *v1.Secret) bool {
	aAnns := a.Annotations
	bAnns := b.Annotations
	delete(aAnns, constants.ManagedByAnnotationKey)
	delete(bAnns, constants.ManagedByAnnotationKey)

	return (a.Type == b.Type &&
		reflect.DeepEqual(a.Data, b.Data) &&
		reflect.DeepEqual(a.StringData, b.StringData) &&
		reflect.DeepEqual(aAnns, bAnns))
}

func prepareSecret(namespace v1.Namespace, secret *v1.Secret) *v1.Secret {
	annotations := secret.Annotations
	if annotations == nil {
		annotations = make(map[string]string)
	}
	annotations[constants.ManagedByAnnotationKey] = constants.ManagedByAnnotationValue

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
