package client

import (
	"context"

	log "github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
)

// SyncSecrets syncs Secrets across all selected Namespaces
func SyncSecrets(config *SyncConfig) (err error) {
	ctx := context.Background()

	log.Debugf("Starting with following configuration: %#v", *config)

	clientset, err := clientset(config)
	if err != nil {
		return err
	}

	secretWatcher, err := clientset.CoreV1().Secrets(config.SecretsNamespace).Watch(ctx, metav1.ListOptions{})
	if err != nil {
		return err
	}

	namespaceWatcher, err := clientset.CoreV1().Namespaces().Watch(ctx, metav1.ListOptions{})
	if err != nil {
		return err
	}

	for {
		select {
		case secretEvent := <-secretWatcher.ResultChan():
			secretEventHandler(ctx, clientset, config, secretEvent)
		case namespaceEvent := <-namespaceWatcher.ResultChan():
			namespaceEventHandler(ctx, clientset, config, namespaceEvent)
		}
	}
}

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

func namespaceEventHandler(ctx context.Context, clientset *kubernetes.Clientset, config *SyncConfig, event watch.Event) {
	namespace := event.Object.(*v1.Namespace)

	switch event.Type {
	case watch.Added:
		addNamespace(ctx, clientset, config, namespace)
	}
}

func addNamespace(ctx context.Context, clientset *kubernetes.Clientset, config *SyncConfig, namespace *v1.Namespace) {
	log.Infof("[%s]: Namespace created", namespace.Name)
}
