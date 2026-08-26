package util

import "strings"

// StringPtrIfNotEmpty trims value and returns a pointer to it, or nil if
// the trimmed result is empty. Was duplicated identically in
// internal/httpapi/nodes and internal/httpapi/subscription.
func StringPtrIfNotEmpty(value string) *string {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil
	}
	return &value
}

// LowerStringPtr trims and lowercases *value, returning a pointer to the
// result, or nil if value is nil or the trimmed result is empty. Was
// duplicated identically in internal/jobqueue and internal/httpapi/subscription.
func LowerStringPtr(value *string) *string {
	if value == nil {
		return nil
	}
	lowered := strings.ToLower(strings.TrimSpace(*value))
	if lowered == "" {
		return nil
	}
	return &lowered
}
