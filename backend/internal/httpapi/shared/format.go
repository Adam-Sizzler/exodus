package shared

import (
	"strconv"
	"strings"
)

// ParseUptimeSeconds parses an uptime value stored as a decimal-seconds
// string (as kept in the database, since it originates from sing-box's
// own uptime counter) into an int64 number of seconds for API responses.
// Invalid or empty values default to 0 rather than erroring, since uptime
// is best-effort telemetry and should never block a response.
func ParseUptimeSeconds(raw string) int64 {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return 0
	}

	seconds, err := strconv.ParseInt(raw, 10, 64)
	if err != nil {
		return 0
	}

	return seconds
}
