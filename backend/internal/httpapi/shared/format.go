package shared

import "fmt"

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
