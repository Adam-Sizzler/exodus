package util

// Coalesce returns *value if value is non-nil, otherwise fallback. This is
// the single generic implementation behind what used to be four identical,
// hand-copied functions (coalesceInt, coalesceInt64, coalesceBool,
// coalesceFloat) scattered across internal/httpapi/nodes, hosts, users, and
// externalsquads.
func Coalesce[T any](value *T, fallback T) T {
	if value == nil {
		return fallback
	}
	return *value
}
