package client

import (
	"context"

	typesv1 "github.com/alehechka/kube-secret-sync/api/types/v1"
	"github.com/alehechka/kube-secret-sync/api/types/v1/clientset"
	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
)

func secretSyncRuleEventHandler(ctx context.Context, event watch.Event) error {
	rule, ok := event.Object.(*typesv1.SecretSyncRule)
	if !ok {
		log.Error("failed to cast SecretSyncRule")
		return nil
	}

	switch event.Type {
	case watch.Added:
		return addedSecretSyncRuleHandler(ctx, rule)
	case watch.Modified:
		return modifiedSecretSyncRuleHandler(ctx, rule)
	case watch.Deleted:
		return deletedSecretSyncRuleHandler(ctx, rule)
	}

	return nil
}

func addedSecretSyncRuleHandler(ctx context.Context, rule *typesv1.SecretSyncRule) error {
	ruleLogger(rule).Infof("added")

	secret, err := getSecret(ctx, rule.Spec.Namespace, rule.Spec.Secret)
	if err != nil {
		return err
	}

	for _, namespace := range rule.Namespaces(ctx) {
		createUpdateSecret(ctx, rule.Spec.Rules, &namespace, secret)
	}

	return nil
}

func modifiedSecretSyncRuleHandler(ctx context.Context, rule *typesv1.SecretSyncRule) error {
	ruleLogger(rule).Infof("modified")
	return nil
}

func deletedSecretSyncRuleHandler(ctx context.Context, rule *typesv1.SecretSyncRule) error {
	ruleLogger(rule).Infof("deleted")
	return nil
}

func listSecretSyncRules(ctx context.Context) (rules *typesv1.SecretSyncRuleList, err error) {
	rules, err = clientset.KubeSecretSync.SecretSyncRules().List(ctx, metav1.ListOptions{})
	if err != nil {
		log.Errorf("failed to list SecretSyncRules: %s", err.Error())
	}
	return
}
