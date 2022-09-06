package client

import (
	"context"

	typesv1 "github.com/alehechka/kube-secret-sync/api/types/v1"
	"github.com/alehechka/kube-secret-sync/clientset"
	log "github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
)

func namespaceEventHandler(ctx context.Context, event watch.Event) error {
	namespace, ok := event.Object.(*v1.Namespace)
	if !ok {
		log.Error("failed to cast Namespace")
		return nil
	}

	switch event.Type {
	case watch.Added:
		return addedNamespaceHandler(ctx, namespace)
	}

	return nil
}

func addedNamespaceHandler(ctx context.Context, namespace *v1.Namespace) error {
	logger := namespaceLogger(namespace)

	if namespace.CreationTimestamp.Time.Before(startTime) {
		logger.Debugf("namespace will be synced on startup by SecretSyncRule watcher")
		return nil
	}

	logger.Infof("added")
	return syncNamespace(ctx, namespace)
}

func syncNamespace(ctx context.Context, namespace *v1.Namespace) error {
	namespaceLogger(namespace).Debugf("syncing new namespace")

	rules, err := listSecretSyncRules(ctx)
	if err != nil {
		return err
	}

	for _, rule := range rules.Items {
		if rule.ShouldSyncNamespace(namespace) {
			syncSecretToNamespace(ctx, namespace, &rule)
		}
	}

	return nil
}

func syncSecretToNamespace(ctx context.Context, namespace *v1.Namespace, rule *typesv1.SecretSyncRule) error {
	secret, err := getSecret(ctx, rule.Spec.Namespace, rule.Spec.Secret)
	if err != nil {
		return err
	}

	return createUpdateSecret(ctx, rule.Spec.Rules, namespace, secret)
}

func listNamespaces(ctx context.Context) (namespaces *v1.NamespaceList, err error) {
	namespaces, err = clientset.Default.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
	if err != nil {
		log.Errorf("failed to list namespaces: %s", err.Error())
	}
	return
}
