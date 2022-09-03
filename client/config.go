package client

import (
	"github.com/alehechka/kube-secret-sync/api/types"
)

// SyncConfig contains the configuration options for the SyncSecrets operation.
type SyncConfig struct {
	ExcludeSecrets      types.StringSlice
	ExcludeRegexSecrets types.RegexSlice
	IncludeSecrets      types.StringSlice
	IncludeRegexSecrets types.RegexSlice

	ExcludeNamespaces      types.StringSlice
	ExcludeRegexNamespaces types.RegexSlice
	IncludeNamespaces      types.StringSlice
	IncludeRegexNamespaces types.RegexSlice

	SecretsNamespace string

	ForceSync bool

	OutOfCluster bool
	KubeConfig   string
}
