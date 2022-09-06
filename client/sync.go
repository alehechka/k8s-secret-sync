package client

import (
	"context"
	"time"

	"github.com/alehechka/kube-secret-sync/clientset"
	log "github.com/sirupsen/logrus"
)

var startTime time.Time

// SyncSecrets syncs Secrets across all selected Namespaces
func SyncSecrets(config *clientset.SyncConfig) (err error) {
	ctx := context.Background()

	startTime = time.Now()

	log.Debugf("Starting with following configuration: %#v", *config)

	if err = initClientsets(config); err != nil {
		return err
	}

	secretWatcher, namespaceWatcher, secretSyncRuleWatcher, err := initWatchers(ctx)
	if err != nil {
		return err
	}

	defer secretSyncRuleWatcher.Stop()
	defer secretWatcher.Stop()
	defer namespaceWatcher.Stop()

	signalChan := initSignalChannel()

	for {
		select {
		case secretEvent, ok := <-secretWatcher.ResultChan():
			if !ok {
				log.Debug("Secret watcher timed out, restarting now.")
				if secretWatcher, err = SecretWatcher(ctx); err != nil {
					return err
				}
				continue
			}
			secretEventHandler(ctx, secretEvent)
		case namespaceEvent, ok := <-namespaceWatcher.ResultChan():
			if !ok {
				log.Debug("Namespace watcher timed out, restarting now.")
				if namespaceWatcher, err = NamespaceWatcher(ctx); err != nil {
					return err
				}
				continue
			}
			namespaceEventHandler(ctx, namespaceEvent)
		case secretSyncRuleEvent, ok := <-secretSyncRuleWatcher.ResultChan():
			if !ok {
				log.Debug("SecretSyncRule watcher timed out, restarting now.")
				if secretSyncRuleWatcher, err = SecretSyncRuleWatcher(ctx); err != nil {
					return err
				}
				continue
			}
			secretSyncRuleEventHandler(ctx, secretSyncRuleEvent)
		case s := <-signalChan:
			log.Infof("Shutting down from signal: %s", s)
			return nil
		}
	}
}
