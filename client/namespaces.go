package client

import (
	"context"

	"github.com/alehechka/kube-secret-sync/constants"
	log "github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
)

func namespaceEventHandler(ctx context.Context, clientset *kubernetes.Clientset, config *SyncConfig, event watch.Event) {
	namespace := event.Object.(*v1.Namespace)

	switch event.Type {
	case watch.Added:
		addNamespace(ctx, clientset, config, namespace)
	}
}

func addNamespace(ctx context.Context, clientset *kubernetes.Clientset, config *SyncConfig, namespace *v1.Namespace) {
	log.Infof("[%s]: Namespace added", namespace.Name)

	if namespace.CreationTimestamp.Time.Before(startTime) {
		log.Debugf("[%s]: Namespace will be synced on startup by Secrets watcher", namespace.Name)
		return
	}

	syncNamespace(ctx, clientset, config, namespace)
}

func syncNamespace(ctx context.Context, clientset *kubernetes.Clientset, config *SyncConfig, namespace *v1.Namespace) error {
	log.Debugf("[%s]: Syncing new namespace", namespace.Name)

	if err := verifyNamespace(config, *namespace); err != nil {
		return err
	}

	secrets, err := listSecrets(ctx, clientset, config.SecretsNamespace)
	if err != nil {
		log.Errorf("Failed to list secrets: %s", err.Error())
		return err
	}

	for _, secret := range secrets.Items {
		if isInvalidSecret(config, &secret) {
			continue
		}

		syncAddedModifiedSecret(ctx, clientset, config, *namespace, &secret)
	}

	return nil
}

func listNamespaces(ctx context.Context, clientset *kubernetes.Clientset) (*v1.NamespaceList, error) {
	return clientset.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
}

func verifyNamespace(config *SyncConfig, namespace v1.Namespace) error {
	if namespace.Name == config.SecretsNamespace {
		log.Debugf("[%s]: Skipping secrets namespace", namespace.Name)
		return constants.ErrSecretsNamespace
	}

	if config.ExcludeNamespaces.IsExcluded(namespace.Name) {
		log.Debugf("[%s]: Namespace has been excluded from sync", namespace.Name)
		return constants.ErrExcludedNamespace
	}

	if !config.IncludeNamespaces.IsIncluded(namespace.Name) {
		log.Debugf("[%s]: Namespace is not included for sync", namespace.Name)
		return constants.ErrNotIncludedNamespace
	}

	return nil
}

func isInvalidNamespace(config *SyncConfig, namespace v1.Namespace) bool {
	err := verifyNamespace(config, namespace)

	return err != nil
}
