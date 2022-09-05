package client

import (
	"context"

	typesv1 "github.com/alehechka/kube-secret-sync/api/types/v1"
	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
)

func secretSyncRuleEventHandler(ctx context.Context, event watch.Event) {
	rule, ok := event.Object.(*typesv1.SecretSyncRule)
	if !ok {
		log.Error("failed to cast SecretSyncRule")
	}

	switch event.Type {
	case watch.Added:
		addSecretSyncRule(ctx, rule)
	case watch.Modified:
		modifySecretSyncRule(ctx, rule)
	case watch.Deleted:
		deleteSecretSyncRule(ctx, rule)
	}
}

func addSecretSyncRule(ctx context.Context, rule *typesv1.SecretSyncRule) error {
	ruleLogger(rule).Infof("added")

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
			createUpdateSecret(ctx, rule.Spec.Rules, &namespace, secret)
		}
	}

	return nil
}

func modifySecretSyncRule(ctx context.Context, rule *typesv1.SecretSyncRule) {
	ruleLogger(rule).Infof("modified")
}

func deleteSecretSyncRule(ctx context.Context, rule *typesv1.SecretSyncRule) {
	ruleLogger(rule).Infof("deleted")
}

func listSecretSyncRules(ctx context.Context) (rules *typesv1.SecretSyncRuleList, err error) {
	rules, err = KubeSecretSyncClientset.SecretSyncRules().List(ctx, metav1.ListOptions{})
	if err != nil {
		log.Errorf("failed to list SecretSyncRules: %s", err.Error())
	}
	return
}
