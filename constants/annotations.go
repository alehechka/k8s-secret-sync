package constants

// ManagedByAnnotationKey is an annotation key appended to kube-secret-sync managed secrets.
const ManagedByAnnotationKey = "app.kubernetes.io/managed-by"

// ManagedByAnnotationValue is an annotation value appended to kube-secret-sync managed secrets.
const ManagedByAnnotationValue = "kube-secret-sync"

// LastAppliedConfigurationAnnotationKey is an annotation created by Kubernetes to keep track of last config
const LastAppliedConfigurationAnnotationKey = "kubectl.kubernetes.io/last-applied-configuration"
