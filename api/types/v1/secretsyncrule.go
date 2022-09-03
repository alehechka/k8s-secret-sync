package v1

import (
	"github.com/alehechka/kube-secret-sync/api/types"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +kubebuilder:object:root=true

// SecretSyncRule is the definition for the SecretSyncRule CRD
type SecretSyncRule struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec SecretSyncRuleSpec
}

// +kubebuilder:object:root=true

// SecretSyncRuleList is the definition for the SecretSyncRule CRD list
type SecretSyncRuleList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []SecretSyncRule
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
