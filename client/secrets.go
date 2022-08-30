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
	log.Infof("[%s] Secret added: %s", secret.ObjectMeta.Namespace, secret.ObjectMeta.Name)

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

		// Secret already exists
		if namespaceSecret, err := getSecret(ctx, clientset, namespace.Name, secret.Name); err == nil {
			log.Debugf("Secret already exists: %s/%s", namespace.Name, secret.Name)

			// Existing secret contains same data
			if reflect.DeepEqual(namespaceSecret.Data, secret.Data) {
				log.Debugf("Existing secret contains same data: %s/%s", namespace.Name, secret.Name)
				continue
			}

			log.Debugf("TODO: update secret to match new secret: %s/%s", namespace.Name, secret.Name)
		}

		createSecret(ctx, clientset, namespace, secret)
	}

	return nil
}

func createSecret(ctx context.Context, clientset *kubernetes.Clientset, namespace v1.Namespace, secret *v1.Secret) error {
	log.Infof("Creating secret: %s/%s", namespace.Name, secret.Name)

	newSecret := &v1.Secret{
		TypeMeta: secret.TypeMeta,
		ObjectMeta: metav1.ObjectMeta{
			Name:        secret.Name,
			Namespace:   namespace.Name,
			Labels:      secret.Labels,
			Annotations: secret.Annotations,
		},
		Immutable:  secret.Immutable,
		Data:       secret.Data,
		StringData: secret.StringData,
		Type:       secret.Type,
	}

	_, err := clientset.CoreV1().Secrets(namespace.Name).Create(ctx, newSecret, metav1.CreateOptions{})

	if err != nil {
		log.Errorf("Failed to create secret %s/%s: %s", namespace.Name, secret.Name, err.Error())
	}

	return err
}

func modifySecrets(ctx context.Context, clientset *kubernetes.Clientset, config *SyncConfig, secret *v1.Secret) {
	log.Infof("[%s] Secret modified: %s", secret.ObjectMeta.Namespace, secret.ObjectMeta.Name)
}

func deleteSecrets(ctx context.Context, clientset *kubernetes.Clientset, config *SyncConfig, secret *v1.Secret) {
	log.Infof("[%s] Secret deleted: %s", secret.ObjectMeta.Namespace, secret.ObjectMeta.Name)
}

func getSecret(ctx context.Context, clientset *kubernetes.Clientset, namespace, name string) (*v1.Secret, error) {
	return clientset.CoreV1().Secrets(namespace).Get(ctx, name, metav1.GetOptions{})
}
