package constants

import "errors"

// ErrExcludedSecret is the error returned when an event for a Secret that has been excluded is triggered.
var ErrExcludedSecret = errors.New("secret was marked as excluded")

// ErrNotIncludedSecret is the error returned when an event for a Secret that is not included is triggered.
var ErrNotIncludedSecret = errors.New("secret was not marked as included")

// ErrExcludedNamespace is the error returned when an event for a Namespace that has been excluded is triggered.
var ErrExcludedNamespace = errors.New("namespaces was not marked as included")

// ErrNotIncludedNamespace is the error returned when an event for a Namespace that is not included is triggered.
var ErrNotIncludedNamespace = errors.New("namespaces was not marked as included")

// ErrSecretsNamespace is the error returned when an event for the Namespace that secrets are synced from is triggered.
var ErrSecretsNamespace = errors.New("namespace is used to sync secrets from")
