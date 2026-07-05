package users

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

func validateUUIDList(values []string) error {
	if len(values) == 0 {
		return fmt.Errorf("uuids cannot be empty")
	}
	return validateUUIDListAllowEmpty(values)
}

func validateUUIDListAllowEmpty(values []string) error {
	for _, value := range values {
		if _, err := uuid.Parse(strings.TrimSpace(value)); err != nil {
			return fmt.Errorf("invalid uuid value")
		}
	}
	return nil
}

func dedupeStrings(values []string) []string {
	result := make([]string, 0, len(values))
	seen := make(map[string]struct{}, len(values))
	for _, value := range values {
		trimmed := strings.TrimSpace(value)
		if trimmed == "" {
			continue
		}
		if _, ok := seen[trimmed]; ok {
			continue
		}
		seen[trimmed] = struct{}{}
		result = append(result, trimmed)
	}
	return result
}

func normalizeUserStatus(value *string) string {
	if value == nil || strings.TrimSpace(*value) == "" {
		return "ACTIVE"
	}
	return strings.ToUpper(strings.TrimSpace(*value))
}

func normalizeTrafficStrategy(value *string) string {
	if value == nil || strings.TrimSpace(*value) == "" {
		return "NO_RESET"
	}
	return strings.ToUpper(strings.TrimSpace(*value))
}

func isValidUserStatus(value string) bool {
	switch strings.ToUpper(strings.TrimSpace(value)) {
	case "ACTIVE", "DISABLED", "LIMITED", "EXPIRED":
		return true
	default:
		return false
	}
}

func isValidTrafficStrategy(value string) bool {
	switch strings.ToUpper(strings.TrimSpace(value)) {
	case "NO_RESET", "DAY", "WEEK", "MONTH", "MONTH_ROLLING":
		return true
	default:
		return false
	}
}

func normalizeNullableString(value *string) any {
	if value == nil || strings.TrimSpace(*value) == "" {
		return nil
	}
	return strings.TrimSpace(*value)
}

func protocolCredentialString(value *string, fallback string) string {
	if value != nil && strings.TrimSpace(*value) != "" {
		return strings.TrimSpace(*value)
	}
	return strings.TrimSpace(fallback)
}

func normalizeUserTag(value *string) any {
	if value == nil || strings.TrimSpace(*value) == "" {
		return nil
	}
	return strings.ToUpper(strings.TrimSpace(*value))
}

func coalesceInt64(value *int64, fallback int64) int64 {
	if value == nil {
		return fallback
	}
	return *value
}

func coalesceUUID(value *string) string {
	if value != nil && strings.TrimSpace(*value) != "" {
		return strings.TrimSpace(*value)
	}
	return uuid.NewString()
}

func coalesceShortUUID(value *string) string {
	if value != nil && strings.TrimSpace(*value) != "" {
		return strings.TrimSpace(*value)
	}
	return generateSubscriptionShortUUID()
}

func coalesceRandomString(value *string, length int) string {
	if value != nil && strings.TrimSpace(*value) != "" {
		return strings.TrimSpace(*value)
	}
	return generateRandomString(length)
}

func generateRandomString(length int) string {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return ""
	}
	return hex.EncodeToString(bytes)[:length]
}

func generateSubscriptionShortUUID() string {
	bytes := make([]byte, 12)
	if _, err := rand.Read(bytes); err != nil {
		return ""
	}
	return base64.RawURLEncoding.EncodeToString(bytes)
}
