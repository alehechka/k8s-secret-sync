package client

import (
	"context"

	typesv1 "github.com/alehechka/kube-secret-sync/api/types/v1"
	"github.com/alehechka/kube-secret-sync/constants"
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
	log.Infof("[%s]: Namespace added", namespace.Name)

	if namespace.CreationTimestamp.Time.Before(startTime) {
		log.Debugf("[%s]: Namespace will be synced on startup by Secrets watcher", namespace.Name)
		return
	}

	syncNamespace(ctx, config, namespace)
}

func syncNamespace(ctx context.Context, config *SyncConfig, namespace *v1.Namespace) error {
	log.Debugf("[%s]: Syncing new namespace", namespace.Name)

	if err := verifyNamespace(config, *namespace); err != nil {
		return err
	}

	secrets, err := listSecrets(ctx, config.SecretsNamespace)
	if err != nil {
		log.Errorf("Failed to list secrets: %s", err.Error())
		return err
	}

	for _, secret := range secrets.Items {
		if isInvalidSecret(config, &secret) {
			continue
		}

		syncAddedModifiedSecret(ctx, config.ForceSync, *namespace, &secret)
	}

	return nil
}

func listNamespaces(ctx context.Context) (namespaces *v1.NamespaceList, err error) {
	namespaces, err = DefaultClientset.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
	if err != nil {
		log.Errorf("Failed to list namespaces: %s", err.Error())
	}
	return
}

func verifyNamespace(config *SyncConfig, namespace v1.Namespace) error {
	if namespace.Name == config.SecretsNamespace {
		log.Debugf("[%s]: Skipping secrets namespace", namespace.Name)
		return constants.ErrSecretsNamespace
	}

	if config.ExcludeNamespaces.IsExcluded(namespace.Name) || config.ExcludeRegexNamespaces.IsExcluded(namespace.Name) {
		log.Debugf("[%s]: Namespace has been excluded from sync", namespace.Name)
		return constants.ErrExcludedNamespace
	}

	if (config.IncludeNamespaces.IsEmpty() && config.IncludeRegexNamespaces.IsEmpty()) ||
		config.IncludeNamespaces.IsIncluded(namespace.Name) ||
		config.IncludeRegexNamespaces.IsIncluded(namespace.Name) {
		return nil
	}

	log.Debugf("[%s]: Namespace is not included for sync", namespace.Name)
	return constants.ErrNotIncludedNamespace

}

func isInvalidNamespace(config *SyncConfig, namespace v1.Namespace) bool {
	err := verifyNamespace(config, namespace)

	return err != nil
}

func verifyNamespaceRules(rules typesv1.Rules, namespace v1.Namespace) error {
	if rules.Namespaces.Exclude.IsExcluded(namespace.Name) || rules.Namespaces.ExcludeRegex.IsRegexExcluded(namespace.Name) {
		log.Debugf("[%s]: Skipping secrets namespace", namespace.Name)
		return constants.ErrSecretsNamespace
	}

	if (rules.Namespaces.Include.IsEmpty() && rules.Namespaces.IncludeRegex.IsEmpty()) ||
		rules.Namespaces.Include.IsIncluded(namespace.Name) ||
		rules.Namespaces.IncludeRegex.IsIncluded(namespace.Name) {
		return nil
	}

	log.Debugf("[%s]: Namespace is not included for sync", namespace.Name)
	return constants.ErrNotIncludedNamespace
}

func isInvalidNamespaceRules(rules typesv1.Rules, namespace v1.Namespace) bool {
	err := verifyNamespaceRules(rules, namespace)

	return err != nil
}
