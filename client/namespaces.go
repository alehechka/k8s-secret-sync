package client

import (
	"context"

	log "github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
)

func namespaceEventHandler(ctx context.Context, event watch.Event) {
	namespace := event.Object.(*v1.Namespace)

	switch event.Type {
	case watch.Added:
		addNamespace(ctx, namespace)
	}
}

func addNamespace(ctx context.Context, namespace *v1.Namespace) {
	logger := namespaceLogger(namespace)
	logger.Infof("added")

	if namespace.CreationTimestamp.Time.Before(startTime) {
		logger.Debugf("namespace will be synced on startup by SecretSyncRule watcher")
		return
	}

	syncNamespace(ctx, namespace)
}

// TODO - rebuild this function
func syncNamespace(ctx context.Context, namespace *v1.Namespace) error {
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
