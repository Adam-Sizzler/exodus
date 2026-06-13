package nodes

import (
	"math"
	"strings"
)

func normalizeCountryCode(value *string) string {
	if value == nil || strings.TrimSpace(*value) == "" {
		return "XX"
	}
	return strings.ToUpper(strings.TrimSpace(*value))
}

func normalizeTags(tags []string) []string {
	if tags == nil {
		return []string{}
	}
	normalized := make([]string, 0, len(tags))
	for _, tag := range tags {
		tag = strings.ToUpper(strings.TrimSpace(tag))
		if tag == "" {
			continue
		}
		normalized = append(normalized, tag)
	}
	return dedupeStrings(normalized)
}

func normalizeNullableString(value *string) any {
	if value == nil || strings.TrimSpace(*value) == "" {
		return nil
	}
	return strings.TrimSpace(*value)
}

func coalesceBool(value *bool, fallback bool) bool {
	if value == nil {
		return fallback
	}
	return *value
}

func coalesceInt(value *int, fallback int) int {
	if value == nil {
		return fallback
	}
	return *value
}

func coalesceInt64(value *int64, fallback int64) int64 {
	if value == nil {
		return fallback
	}
	return *value
}

func coalesceFloat(value *float64, fallback float64) float64 {
	if value == nil {
		return fallback
	}
	return *value
}

func coalesceInt64Ptr(value *int64) int64 {
	if value == nil {
		return 0
	}
	return *value
}

func toNanoMultiplier(value float64) int64 {
	return int64(math.Round(value * 1_000_000_000))
}

func normalizeAPISchema(value *string) string {
	if value == nil {
		return "mtls"
	}
	normalized := strings.ToLower(strings.TrimSpace(*value))
	switch normalized {
	case "tls":
		return "tls"
	case "mtls", "grpc", "grpcs", "https", "":
		return "mtls"
	default:
		return "mtls"
	}
}

func normalizeAPIPath(value *string) string {
	if value == nil {
		return "/"
	}
	trimmed := strings.TrimSpace(*value)
	if trimmed == "" || trimmed == "/" {
		return "/"
	}
	return "/" + strings.Trim(trimmed, "/")
}

func isAllowedNodeAPISchema(value string) bool {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "mtls", "tls", "grpc", "grpcs", "https":
		return true
	default:
		return false
	}
}

func fromNanoMultiplier(value int64) float64 {
	return float64(value) / 1_000_000_000
}

func ensureStringSlice(values []string) []string {
	if values == nil {
		return []string{}
	}
	return values
}

func ensureInboundSlice(values []configProfileInboundResponse) []configProfileInboundResponse {
	if values == nil {
		return []configProfileInboundResponse{}
	}
	return values
}

func dedupeStrings(values []string) []string {
	if len(values) < 2 {
		return values
	}
	seen := make(map[string]struct{}, len(values))
	result := make([]string, 0, len(values))
	for _, value := range values {
		if _, ok := seen[value]; ok {
			continue
		}
		seen[value] = struct{}{}
		result = append(result, value)
	}
	return result
}
