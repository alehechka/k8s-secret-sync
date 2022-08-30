package client

import (
	"context"

	log "github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
)

// SyncSecrets syncs Secrets across all selected Namespaces
func SyncSecrets(config *SyncConfig) (err error) {
	ctx := context.Background()

	log.Debugf("Starting with following configuration: %#v", *config)

	clientset, err := clientset(config)
	if err != nil {
		return err
	}

	watcher, err := clientset.CoreV1().Secrets(config.SecretsNamespace).Watch(ctx, metav1.ListOptions{})
	if err != nil {
		return err
	}

	for event := range watcher.ResultChan() {
		secret := event.Object.(*v1.Secret)

		switch event.Type {
		case watch.Added:
			addSecrets(ctx, clientset, config, secret)
		case watch.Modified:
			modifySecrets(ctx, clientset, config, secret)
		case watch.Deleted:
			deleteSecrets(ctx, clientset, config, secret)
		}

	}

	return nil
}
