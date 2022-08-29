package client

type Slice []string

func (slice Slice) IsExcluded(str string) bool {
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

func (slice Slice) IsIncluded(str string) bool {
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

type Config struct {
	ExcludeSecrets Slice
	IncludeSecrets Slice

	ExcludeNamespaces Slice
	IncludeNamespaces Slice

	SecretsNamespace string

	OutOfCluster bool
	KubeConfig   string
}
