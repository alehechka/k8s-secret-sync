package client

import (
	typesv1 "github.com/alehechka/kube-secret-sync/api/types/v1"
	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
)

func (client *Client) SecretSyncRuleEventHandler(event watch.Event) error {
	rule, ok := event.Object.(*typesv1.SecretSyncRule)
	if !ok {
		log.Error("failed to cast SecretSyncRule")
		return nil
	}

	switch event.Type {
	case watch.Added:
		return client.AddedSecretSyncRuleHandler(rule)
	case watch.Modified:
		return client.ModifiedSecretSyncRuleHandler(rule)
	case watch.Deleted:
		return client.DeletedSecretSyncRuleHandler(rule)
	}

	return nil
}

func (client *Client) AddedSecretSyncRuleHandler(rule *typesv1.SecretSyncRule) error {
	ruleLogger(rule).Infof("added")

	secret, err := client.GetSecret(rule.Spec.Secret.Namespace, rule.Spec.Secret.Name)
	if err != nil {
		return err
	}

	for _, namespace := range rule.Namespaces(client.Context, client.DefaultClientset) {
		client.CreateUpdateSecret(rule.Spec.Rules, &namespace, secret)
	}

	return nil
}

// modifiedSecretSyncRuleHandler handles syncing secrets after a SecretSyncRule has been modified
//
// Due to the event watcher only providing the new state of the modified resource, it is impossible to know the previous state.
// (The exception to this is potentially "applied" changes and parsing the last-applied-configuration annotation)
// In coping with this limitation, a modified SecretSyncRule will simply attempt to resync the rule across all applicable namespaces.
func (client *Client) ModifiedSecretSyncRuleHandler(rule *typesv1.SecretSyncRule) error {
	if rule.DeletionTimestamp != nil {
		return nil
	}

	ruleLogger(rule).Infof("modified")

	secret, err := client.GetSecret(rule.Spec.Secret.Namespace, rule.Spec.Secret.Name)
	if err != nil {
		return err
	}

	for _, namespace := range rule.Namespaces(client.Context, client.DefaultClientset) {
		client.CreateUpdateSecret(rule.Spec.Rules, &namespace, secret)
	}

	return nil
}

func (client *Client) DeletedSecretSyncRuleHandler(rule *typesv1.SecretSyncRule) error {
	ruleLogger(rule).Infof("deleted")

	secret, err := client.GetSecret(rule.Spec.Secret.Namespace, rule.Spec.Secret.Name)
	if err != nil {
		return err
	}

	for _, namespace := range rule.Namespaces(client.Context, client.DefaultClientset) {
		client.SyncDeletedSecret(rule.Spec.Rules, &namespace, secret)
	}

	return nil
}

func (client *Client) ListSecretSyncRules() (rules *typesv1.SecretSyncRuleList, err error) {
	rules, err = client.KubeSecretSyncClientset.SecretSyncRules().List(client.Context, metav1.ListOptions{})
	if err != nil {
		log.Errorf("failed to list SecretSyncRules: %s", err.Error())
	}
	return
}
