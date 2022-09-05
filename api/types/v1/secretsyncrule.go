package v1

import (
	"context"

	"github.com/alehechka/kube-secret-sync/api/types"
	"github.com/alehechka/kube-secret-sync/clientset"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +kubebuilder:object:root=true

// SecretSyncRule is the definition for the SecretSyncRule CRD
type SecretSyncRule struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec SecretSyncRuleSpec `json:"spec"`
}

// +kubebuilder:object:root=true

// SecretSyncRuleList is the definition for the SecretSyncRule CRD list
type SecretSyncRuleList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []SecretSyncRule `json:"items"`
}

// +kubebuilder:object:generate=true

// SecretSyncRuleSpec is the spec attribute of the SecretSyncRule CRD
type SecretSyncRuleSpec struct {
	Secret    string `json:"secret"`
	Namespace string `json:"namespace"`
	Rules     Rules  `json:"rules"`
}

// +kubebuilder:object:generate=true

// Rules contains all rules for the secret to follow
type Rules struct {
	Namespaces NamespaceRules `json:"namespaces"`
	Force      bool           `json:"force"`
}

// +kubebuilder:object:generate=true

// NamespaceRules include all rules for namepsaces to sync to.
type NamespaceRules struct {
	Exclude      types.StringSlice `json:"exclude"`
	ExcludeRegex types.StringSlice `json:"excludeRegex"`
	Include      types.StringSlice `json:"include"`
	IncludeRegex types.StringSlice `json:"includeRegex"`
}

// ShouldSyncSecret determines whether or not the given Secret should be synced
func (rule *SecretSyncRule) ShouldSyncSecret(secret *v1.Secret) bool {
	return rule.Spec.Secret == secret.Name && rule.Spec.Namespace == secret.Namespace
}

// ShouldSyncNamespace determines whether or not the given Namespace should be synced
func (rule *SecretSyncRule) ShouldSyncNamespace(namespace *v1.Namespace) bool {
	rules := rule.Spec.Rules

	if rule.Spec.Namespace == namespace.Name {
		return false
	}

	if rules.Namespaces.Exclude.IsExcluded(namespace.Name) || rules.Namespaces.ExcludeRegex.IsRegexExcluded(namespace.Name) {
		return false
	}

	if rules.Namespaces.Include.IsEmpty() && rules.Namespaces.IncludeRegex.IsEmpty() {
		return true
	}

	if rules.Namespaces.Include.IsIncluded(namespace.Name) || rules.Namespaces.IncludeRegex.IsRegexIncluded(namespace.Name) {
		return true
	}

	return false
}

// Namespaces returns a list of all namespaces that the given Rule allows for syncing
func (rule *SecretSyncRule) Namespaces(ctx context.Context) (namespaces []v1.Namespace) {
	list, err := clientset.Default.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
	if err != nil {
		return
	}

	for _, namespace := range list.Items {
		if rule.ShouldSyncNamespace(&namespace) {
			namespaces = append(namespaces, namespace)
		}
	}

	return
}

// ShouldSyncSecret iterates over the list to determine whether or not the given Secret should be synced
func (list *SecretSyncRuleList) ShouldSyncSecret(secret *v1.Secret) bool {
	for _, rule := range list.Items {
		if rule.ShouldSyncSecret(secret) {
			return true
		}
	}

	return false
}

// ShouldSyncNamespace iterates over the list to determine whether or not the given Namespace should be synced
func (list *SecretSyncRuleList) ShouldSyncNamespace(namespace *v1.Namespace) bool {
	for _, rule := range list.Items {
		if rule.ShouldSyncNamespace(namespace) {
			return true
		}
	}

	return false
}
