package api

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
)

func parseUUIDCSV(raw string) ([]string, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil, fmt.Errorf("uuids query parameter is required")
	}

	parts := strings.Split(raw, ",")
	out := make([]string, 0, len(parts))
	seen := make(map[string]struct{}, len(parts))

	for _, part := range parts {
		id := strings.TrimSpace(part)
		if id == "" {
			continue
		}
		if _, err := uuid.Parse(id); err != nil {
			return nil, fmt.Errorf("invalid uuid in list")
		}
		if _, exists := seen[id]; exists {
			continue
		}
		seen[id] = struct{}{}
		out = append(out, id)
	}

	if len(out) == 0 {
		return nil, fmt.Errorf("no valid uuids provided")
	}

	return out, nil
}
