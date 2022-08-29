package client

import (
	"context"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func listNamespaces(ctx context.Context, clientset *kubernetes.Clientset) (*v1.NamespaceList, error) {
	return clientset.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
}
