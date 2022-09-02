package client

// SyncConfig contains the configuration options for the SyncSecrets operation.
type SyncConfig struct {
	ExcludeSecrets      StringSlice
	ExcludeRegexSecrets RegexSlice
	IncludeSecrets      StringSlice
	IncludeRegexSecrets RegexSlice

	ExcludeNamespaces      StringSlice
	ExcludeRegexNamespaces RegexSlice
	IncludeNamespaces      StringSlice
	IncludeRegexNamespaces RegexSlice

	SecretsNamespace string

	ForceSync bool

	OutOfCluster bool
	KubeConfig   string
}
