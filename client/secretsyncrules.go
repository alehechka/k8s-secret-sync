package client

import (
	"context"

	typesv1 "github.com/alehechka/kube-secret-sync/api/types/v1"
	log "github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/watch"
)

func secretSyncRuleEventHandler(ctx context.Context, event watch.Event) {
	rule, ok := event.Object.(*typesv1.SecretSyncRule)
	if !ok {
		log.Error("failed to cast SecretSyncRule")
	}

	ruleLogger := ruleLogger(rule)

	switch event.Type {
	case watch.Added:
		addSecretSyncRule(ctx, ruleLogger, rule)
	case watch.Modified:
		modifySecretSyncRule(ctx, ruleLogger, rule)
	case watch.Deleted:
		deleteSecretSyncRule(ctx, ruleLogger, rule)
	}
}

func addSecretSyncRule(ctx context.Context, ruleLogger *log.Entry, rule *typesv1.SecretSyncRule) error {
	ruleLogger.Infof("added")

	secret, err := getSecret(ctx, rule.Spec.Namespace, rule.Spec.Secret)
	if err != nil {
		secretLogger(secret).Errorf("does not exist to sync: %s", err.Error())
		return err
	}

	namespaces, err := listNamespaces(ctx)
	if err != nil {
		return err
	}

	for _, namespace := range namespaces.Items {
		if rule.ShouldSyncNamespace(&namespace) {
			syncAddedModifiedSecret(ctx, rule.Spec.Rules, namespace, secret)
		}
	}

	return nil
}

func modifySecretSyncRule(ctx context.Context, ruleLogger *log.Entry, rule *typesv1.SecretSyncRule) {
	ruleLogger.Infof("modified")
}

func deleteSecretSyncRule(ctx context.Context, ruleLogger *log.Entry, rule *typesv1.SecretSyncRule) {
	ruleLogger.Infof("deleted")
}
