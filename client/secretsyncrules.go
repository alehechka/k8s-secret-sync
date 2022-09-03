package client

import (
	"context"

	v1 "github.com/alehechka/kube-secret-sync/api/types/v1"
	log "github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/watch"
)

func secretSyncRuleEventHandler(ctx context.Context, config *SyncConfig, event watch.Event) {
	rule, ok := event.Object.(*v1.SecretSyncRule)
	if !ok {
		log.Error("failed to cast SecretSyncRule")
	}

	switch event.Type {
	case watch.Added:
		addSecretSyncRule(ctx, config, rule)
	case watch.Modified:
		modifySecretSyncRule(ctx, config, rule)
	case watch.Deleted:
		deleteSecretSyncRule(ctx, config, rule)
	}
}

func addSecretSyncRule(ctx context.Context, config *SyncConfig, rule *v1.SecretSyncRule) {
	log.Infof("[%s/%s]: SecretSyncRule added", rule.Namespace, rule.Name)

	if rule.CreationTimestamp.Time.Before(startTime) {
		log.Debugf("[%s/%s]: SecretSyncRule will be synced on startup by Secrets watcher", rule.Namespace, rule.Name)
		return
	}

	log.Infof("[%s/%s]: SecretSyncRule syncing", rule.Namespace, rule.Name)
}

func modifySecretSyncRule(ctx context.Context, config *SyncConfig, rule *v1.SecretSyncRule) {
	log.Infof("[%s/%s]: SecretSyncRule modified", rule.Namespace, rule.Name)
}

func deleteSecretSyncRule(ctx context.Context, config *SyncConfig, rule *v1.SecretSyncRule) {
	log.Infof("[%s/%s]: SecretSyncRule deleted", rule.Namespace, rule.Name)
}
