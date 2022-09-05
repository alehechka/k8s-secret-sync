package client

import (
	"context"

	log "github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
)

func namespaceEventHandler(ctx context.Context, config *SyncConfig, event watch.Event) {
	namespace := event.Object.(*v1.Namespace)

	switch event.Type {
	case watch.Added:
		addNamespace(ctx, config, namespace)
	}
}

func addNamespace(ctx context.Context, config *SyncConfig, namespace *v1.Namespace) {
	logger := namespaceLogger(namespace)
	logger.Infof("added")

	if namespace.CreationTimestamp.Time.Before(startTime) {
		logger.Debugf("namespace will be synced on startup by Secrets watcher")
		return
	}

	syncNamespace(ctx, config, namespace)
}

// TODO - rebuild this function
func syncNamespace(ctx context.Context, config *SyncConfig, namespace *v1.Namespace) error {
	namespaceLogger(namespace).Debugf("syncing new namespace")

	return nil
}

func listNamespaces(ctx context.Context) (namespaces *v1.NamespaceList, err error) {
	namespaces, err = DefaultClientset.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
	if err != nil {
		log.Errorf("failed to list namespaces: %s", err.Error())
	}
	return
}
