package constants

import "errors"

// ErrExcludedSecret is the error returned when an event for a Secret that has been excluded is triggered.
var ErrExcludedSecret = errors.New("secret was marked as excluded")

// ErrNotIncludedSecret is the error returned when an event for a Secret that is not included is triggered.
var ErrNotIncludedSecret = errors.New("secret was not marked as included")
