package client

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var startTime time.Time

// SyncSecrets syncs Secrets across all selected Namespaces
func SyncSecrets(config *SyncConfig) (err error) {
	ctx := context.Background()

	startTime = time.Now()

	log.Debugf("Starting with following configuration: %#v", *config)

	if err := InitializeDefaultClientset(config); err != nil {
		return err
	}

	if err := InitializeKubeSecretSyncClient(config); err != nil {
		return err
	}

	secretWatcher, err := DefaultClientset.CoreV1().Secrets(config.SecretsNamespace).Watch(ctx, metav1.ListOptions{})
	if err != nil {
		return err
	}

	namespaceWatcher, err := DefaultClientset.CoreV1().Namespaces().Watch(ctx, metav1.ListOptions{})
	if err != nil {
		return err
	}

	secretsyncruleWatcher, err := KubeSecretSyncClientset.SecretSyncRules(v1.NamespaceAll).Watch(ctx, metav1.ListOptions{})
	if err != nil {
		log.Error("Failed to start watching secretsyncrules")
		return err
	}

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGKILL)

	for {
		select {
		case secretEvent := <-secretWatcher.ResultChan():
			secretEventHandler(ctx, config, secretEvent)
		case namespaceEvent := <-namespaceWatcher.ResultChan():
			namespaceEventHandler(ctx, config, namespaceEvent)
		case secretsyncruleEvent := <-secretsyncruleWatcher.ResultChan():
			secretSyncRuleEventHandler(ctx, config, secretsyncruleEvent)
		case s := <-sigc:
			log.Infof("Shutting down from signal: %s", s)
			secretWatcher.Stop()
			namespaceWatcher.Stop()
			return nil
		}
	}
}
