package types

import (
	"regexp"
)

// RegexSlice is a slice of compiled regular expressions
type RegexSlice []*regexp.Regexp

// CompileAll compiles a slice of strings into a RegexSlice
func CompileAll(patterns []string) (slice RegexSlice, err error) {
	for _, pattern := range patterns {
		reg, err := regexp.Compile(pattern)
		if err != nil {
			return nil, err
		}
		slice = append(slice, reg)
	}
	return
}

// IsExcluded determines whether or not a given string is excluded in terms of a blacklist RegexSlice of regular expressions.
// If the RegexSlice is empty, then all provided strings are considered not excluded.
// If the provided string exists within the RegexSlice regular expressions, then it is considered excluded.
func (slice RegexSlice) IsExcluded(str string) bool {
	if slice.IsEmpty() {
		return false
	}

	return slice.exists(str)
}

// IsIncluded determines whether or not a given string is included in terms of a whitelist RegexSlice of regular expressions.
// If the RegexSlice is empty, then all provided strings are considered included.
// If the RegexSlice is not empty, then only strings that exist in the RegexSlice regular expressions will be considered included.
func (slice RegexSlice) IsIncluded(str string) bool {
	return slice.exists(str)
}

func (slice RegexSlice) exists(str string) bool {
	for _, reg := range slice {
		if reg.MatchString(str) {
			return true
		}
	}

	return false
}

// IsEmpty checks if length of slice is greater than zero
func (slice RegexSlice) IsEmpty() bool {
	return len(slice) == 0
}
