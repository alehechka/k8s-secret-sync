package client

import (
	"context"
	"reflect"

	"github.com/alehechka/kube-secret-sync/constants"
	log "github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/client-go/kubernetes"
)

func addSecrets(ctx context.Context, clientset *kubernetes.Clientset, config *SyncConfig, secret *v1.Secret) error {
	log.Infof("Secret added: %s/%s", secret.ObjectMeta.Namespace, secret.ObjectMeta.Name)

	if config.ExcludeSecrets.IsExcluded(secret.Name) {
		log.Debugf("Secret is excluded from sync: %s", secret.Name)
		return constants.ErrExcludedSecret
	}

	if !config.IncludeSecrets.IsIncluded(secret.Name) {
		log.Debugf("Secret is not included for sync: %s", secret.Name)
		return constants.ErrNotIncludedSecret
	}

	namespaces, err := listNamespaces(ctx, clientset)
	if err != nil {
		log.Errorf("Failed to list namespaces during add: %s", err.Error())
		return err
	}

	for _, namespace := range namespaces.Items {
		if namespace.Name == config.SecretsNamespace {
			log.Debugf("Skipping secrets namespace: %s", namespace.Name)
			continue
		}

		if config.ExcludeNamespaces.IsExcluded(namespace.Name) {
			log.Debugf("Namespace has been excluded from sync: %s", namespace.Name)
			continue
		}

		if !config.IncludeNamespaces.IsIncluded(namespace.Name) {
			log.Debugf("Namespace is not included for sync: %s", namespace.Name)
			continue
		}

		if namespaceSecret, err := getSecret(ctx, clientset, namespace.Name, secret.Name); err == nil {
			log.Debugf("Secret already exists: %s/%s", namespace.Name, secret.Name)

			if !config.ForceSync && !isManagedBy(namespaceSecret) {
				log.Debugf("Existing secret is not managed and will not be force updated: %s/%s", namespace.Name, secret.Name)
				continue
			}

			if isManagedBy(namespaceSecret) && secretsAreEqual(secret, namespaceSecret) {
				log.Debugf("Existing secret contains same data: %s/%s", namespace.Name, secret.Name)
				continue
			}

			updateSecret(ctx, clientset, namespace, secret)
			continue
		}

		createSecret(ctx, clientset, namespace, secret)
	}

	return nil
}

func createSecret(ctx context.Context, clientset *kubernetes.Clientset, namespace v1.Namespace, secret *v1.Secret) error {
	log.Infof("Creating secret: %s/%s", namespace.Name, secret.Name)

	newSecret := prepareSecret(namespace, secret)

	_, err := clientset.CoreV1().Secrets(namespace.Name).Create(ctx, newSecret, metav1.CreateOptions{})

	if err != nil {
		log.Errorf("Failed to create secret %s/%s: %s", namespace.Name, secret.Name, err.Error())
	}

	return err
}

func modifySecrets(ctx context.Context, clientset *kubernetes.Clientset, config *SyncConfig, secret *v1.Secret) {
	log.Infof("Secret modified: %s/%s", secret.ObjectMeta.Namespace, secret.ObjectMeta.Name)
}

func deleteSecrets(ctx context.Context, clientset *kubernetes.Clientset, config *SyncConfig, secret *v1.Secret) {
	log.Infof("Secret deleted: %s/%s", secret.ObjectMeta.Namespace, secret.ObjectMeta.Name)
}

func deleteSecret(ctx context.Context, clientset *kubernetes.Clientset, namespace, name string) (err error) {
	log.Infof("Deleting secret: %s/%s", namespace, name)

	err = clientset.CoreV1().Secrets(namespace).Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil {
		log.Errorf("Failed to delete secret %s/%s: %s", namespace, name, err.Error())
	}

	return
}

func updateSecret(ctx context.Context, clientset *kubernetes.Clientset, namespace v1.Namespace, secret *v1.Secret) (updated *v1.Secret, err error) {
	log.Infof("Updating secret: %s/%s", namespace.Name, secret.Name)

	updateSecret := prepareSecret(namespace, secret)

	updated, err = clientset.CoreV1().Secrets(namespace.Name).Update(ctx, updateSecret, metav1.UpdateOptions{})
	if err != nil {
		log.Errorf("Failed to update secret %s/%s: %s", namespace.Name, secret.Name, err.Error())
	}

	return
}

func getSecret(ctx context.Context, clientset *kubernetes.Clientset, namespace, name string) (*v1.Secret, error) {
	return clientset.CoreV1().Secrets(namespace).Get(ctx, name, metav1.GetOptions{})
}

func secretsAreEqual(a, b *v1.Secret) bool {
	return (a.Type == b.Type &&
		reflect.DeepEqual(a.Data, b.Data) &&
		reflect.DeepEqual(a.StringData, b.StringData))
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
