package shared

import (
	"fmt"
	"strconv"
	"strings"
)

// FormatBytes converts bytes to a human-readable string.
func FormatBytes(bytes int64) string {
	if bytes == 0 {
		return "0 B"
	}

	sign := ""
	value := float64(bytes)
	if value < 0 {
		sign = "-"
		value = -value
	}

	const unit = 1024.0
	sizes := []string{"B", "KiB", "MiB", "GiB", "TiB", "PiB"}
	i := 0
	for value >= unit && i < len(sizes)-1 {
		value /= unit
		i++
	}

	if i == 0 {
		return fmt.Sprintf("%s%d %s", sign, int64(value), sizes[i])
	}

	return fmt.Sprintf("%s%.2f %s", sign, value, sizes[i])
}

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
