package client

import (
	"context"

	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// SyncSecrets syncs Secrets across all selected Namespaces
func SyncSecrets(config *SyncConfig) (err error) {
	ctx := context.Background()

	log.Debugf("Starting with following configuration: %#v", *config)

	clientset, err := clientset(config)
	if err != nil {
		return err
	}

	secretWatcher, err := clientset.CoreV1().Secrets(config.SecretsNamespace).Watch(ctx, metav1.ListOptions{})
	if err != nil {
		return err
	}

	namespaceWatcher, err := clientset.CoreV1().Namespaces().Watch(ctx, metav1.ListOptions{})
	if err != nil {
		return err
	}

	for {
		select {
		case secretEvent := <-secretWatcher.ResultChan():
			secretEventHandler(ctx, clientset, config, secretEvent)
		case namespaceEvent := <-namespaceWatcher.ResultChan():
			namespaceEventHandler(ctx, clientset, config, namespaceEvent)
		}
	}
}
