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

	signalChan := initSignalChannel()

	for {
		select {
		case secretEvent := <-secretWatcher.ResultChan():
			secretEventHandler(ctx, secretEvent)
		case namespaceEvent := <-namespaceWatcher.ResultChan():
			namespaceEventHandler(ctx, namespaceEvent)
		case secretsyncruleEvent := <-secretSyncRuleWatcher.ResultChan():
			secretSyncRuleEventHandler(ctx, secretsyncruleEvent)
		case s := <-signalChan:
			log.Infof("Shutting down from signal: %s", s)
			secretWatcher.Stop()
			namespaceWatcher.Stop()
			return nil
		}
	}
}
