package client

import (
	log "github.com/sirupsen/logrus"
)

// SyncSecrets syncs Secrets across all selected Namespaces
func SyncSecrets(config *SyncConfig) (err error) {
	log.Debugf("Starting with following configuration: %#v", *config)

	client := new(Client)

	if err = client.Initialize(config); err != nil {
		return err
	}

	defer client.SecretWatcher.Stop()
	defer client.NamespaceWatcher.Stop()
	defer client.SecretSyncRuleWatcher.Stop()

	for {
		select {
		case secretEvent, ok := <-client.SecretWatcher.ResultChan():
			if !ok {
				log.Debug("Secret watcher timed out, restarting now.")
				if err := client.StartSecretWatcher(); err != nil {
					return err
				}
				defer client.SecretWatcher.Stop()
				continue
			}
			client.SecretEventHandler(secretEvent)
		case namespaceEvent, ok := <-client.NamespaceWatcher.ResultChan():
			if !ok {
				log.Debug("Namespace watcher timed out, restarting now.")
				if err := client.StartNamespaceWatcher(); err != nil {
					return err
				}
				defer client.NamespaceWatcher.Stop()
				continue
			}
			client.NamespaceEventHandler(namespaceEvent)
		case secretSyncRuleEvent, ok := <-client.SecretSyncRuleWatcher.ResultChan():
			if !ok {
				log.Debug("SecretSyncRule watcher timed out, restarting now.")
				if err := client.StartSecretSyncRuleWatcher(); err != nil {
					return err
				}
				defer client.SecretSyncRuleWatcher.Stop()
				continue
			}
			client.SecretSyncRuleEventHandler(secretSyncRuleEvent)
		case s := <-client.SignalChannel:
			log.Infof("Shutting down from signal: %s", s)
			return nil
		}
	}
}
