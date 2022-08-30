package client

// StringSlice is a slice of strings
type StringSlice []string

// IsExcluded determines whether or not a given string is excluded in terms of a blacklist StringSlice
// If the StringSlice is empty, then all provided strings are considered not excluded.
// If the provided string exists within the StringSlice, then it is considered excluded.
func (slice StringSlice) IsExcluded(str string) bool {
	if len(slice) == 0 {
		return false
	}

	for _, obj := range slice {
		if obj == str {
			return true
		}
	}

	return false
}

// IsIncluded determines whether or not a given string is included in terms of a whitelist StringSlice
// If the StringSlice is empty, then all provided strings are considered included.
// If the StringSlice is not empty, then only strings that exist in the StringSlice will be considered included.
func (slice StringSlice) IsIncluded(str string) bool {
	if len(slice) == 0 {
		return true
	}

	for _, obj := range slice {
		if obj == str {
			return true
		}
	}

	return false
}

// SyncConfig contains the configuration options for the SyncSecrets operation.
type SyncConfig struct {
	ExcludeSecrets StringSlice
	IncludeSecrets StringSlice

	ExcludeNamespaces StringSlice
	IncludeNamespaces StringSlice

	SecretsNamespace string

	OutOfCluster bool
	KubeConfig   string
}
