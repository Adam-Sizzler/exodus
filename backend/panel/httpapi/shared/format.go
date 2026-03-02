package shared

import "fmt"

// FormatBytes converts bytes to a human-readable string.
func FormatBytes(bytes int64) string {
	if bytes == 0 {
		return "0 B"
	}
	const unit = 1024
	sizes := []string{"B", "KB", "MB", "GB", "TB", "PB"}
	i := int64(0)
	for bytes >= unit && i < int64(len(sizes)-1) {
		bytes /= unit
		i++
	}
	return fmt.Sprintf("%.2f %s", float64(bytes)/float64(1024), sizes[i])
}
