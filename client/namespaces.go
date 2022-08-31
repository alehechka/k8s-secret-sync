package client

import (
	"context"

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
	log.Infof("[%s]: Namespace created", namespace.Name)
}

func listNamespaces(ctx context.Context, clientset *kubernetes.Clientset) (*v1.NamespaceList, error) {
	return clientset.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
}
