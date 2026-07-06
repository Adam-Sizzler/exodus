package hosts

import (
	"encoding/json"
	"fmt"
	"strings"
)

func normalizeOptionalStringAllowEmpty(value *string) interface{} {
	if value == nil {
		return nil
	}
	return strings.TrimSpace(*value)
}

func normalizeNullableString(value *string) interface{} {
	if value == nil {
		return nil
	}
	trimmed := strings.TrimSpace(*value)
	if trimmed == "" {
		return nil
	}
	return trimmed
}

func normalizeNullableInt(value *int) interface{} {
	if value == nil {
		return nil
	}
	return *value
}

func normalizeProtocolCredentialForCreate(override *bool, value *string) *string {
	if !coalesceBool(override, false) {
		return nil
	}
	return normalizeProtocolCredentialPointer(value)
}

func normalizeProtocolCredentialPointer(value *string) *string {
	if value == nil {
		return nil
	}
	trimmed := strings.TrimSpace(*value)
	if trimmed == "" {
		return nil
	}
	return &trimmed
}

func normalizeSecurityLayer(value *string) string {
	if value == nil || strings.TrimSpace(*value) == "" {
		return "DEFAULT"
	}
	upper := strings.ToUpper(strings.TrimSpace(*value))
	if _, ok := allowedSecurityLayers[upper]; ok {
		return upper
	}
	return "DEFAULT"
}

func normalizeMihomoIPVersion(value *string) *string {
	if value == nil {
		return nil
	}
	normalized := strings.ToLower(strings.TrimSpace(*value))
	if normalized == "" {
		return nil
	}
	return &normalized
}

func normalizeJSONField(raw *json.RawMessage, emptyObjectAsNull bool) (bool, []byte, error) {
	if raw == nil {
		return false, nil, nil
	}
	trimmed := strings.TrimSpace(string(*raw))
	if trimmed == "" || trimmed == "null" {
		return true, nil, nil
	}
	if !json.Valid(*raw) {
		return true, nil, fmt.Errorf("invalid JSON payload")
	}
	if emptyObjectAsNull {
		var obj map[string]any
		if err := json.Unmarshal(*raw, &obj); err == nil {
			if len(obj) == 0 {
				return true, nil, nil
			}
		}
	}
	return true, []byte(*raw), nil
}

func normalizeJSONValue(raw *json.RawMessage, emptyObjectAsNull bool) ([]byte, error) {
	_, val, err := normalizeJSONField(raw, emptyObjectAsNull)
	return val, err
}

func normalizeOptionalJSONField(raw OptionalJSON, emptyObjectAsNull bool) (bool, []byte, error) {
	if !raw.Set {
		return false, nil, nil
	}
	value := json.RawMessage(raw.Raw)
	return normalizeJSONField(&value, emptyObjectAsNull)
}

func coalesceBool(value *bool, fallback bool) bool {
	if value == nil {
		return fallback
	}
	return *value
}

func ensureStringSlice(values []string) []string {
	if values == nil {
		return []string{}
	}
	return values
}

func normalizeTags(tags []string) []string {
	if tags == nil {
		return []string{}
	}
	normalized := make([]string, 0, len(tags))
	for _, tag := range tags {
		tag = strings.TrimSpace(tag)
		if tag == "" {
			continue
		}
		normalized = append(normalized, tag)
	}
	return dedupeStrings(normalized)
}

func parseJSONAny(raw json.RawMessage) interface{} {
	if len(raw) == 0 {
		return nil
	}
	var value interface{}
	if err := json.Unmarshal(raw, &value); err != nil {
		return nil
	}
	return value
}

func dedupeStrings(values []string) []string {
	if len(values) < 2 {
		return values
	}
	seen := make(map[string]struct{}, len(values))
	out := make([]string, 0, len(values))
	for _, value := range values {
		if _, exists := seen[value]; exists {
			continue
		}
		seen[value] = struct{}{}
		out = append(out, value)
	}
	return out
}
