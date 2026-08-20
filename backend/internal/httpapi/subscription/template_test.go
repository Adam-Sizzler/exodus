package subscription

import (
	"testing"
	"time"
)

func TestFormatTemplateValueStatusAndResetStrategy(t *testing.T) {
	user := SubscriptionUser{
		Username:             "test_user",
		ShortUUID:            "abc123xyz",
		Status:               "ACTIVE",
		TrafficLimitStrategy: "MONTH",
		TrafficLimitBytes:    10 * 1024 * 1024 * 1024,
		UsedTrafficBytes:     2 * 1024 * 1024 * 1024,
		ExpireAt:             time.Now().Add(24 * 5 * time.Hour),
	}
	settings := SubscriptionSettingsParsed{}
	subURL := "https://sub.example.com/abc123xyz"

	// 1. Default STATUS
	val1 := formatTemplateValue("Ваша подписка: {{STATUS}}", user, settings, subURL)
	if val1 != "Ваша подписка: Active" {
		t.Errorf("expected 'Ваша подписка: Active', got %q", val1)
	}

	// 2. Custom STATUS override
	val2 := formatTemplateValue("Ваша подписка: {{STATUS:ACTIVE=✅ Активна|EXPIRED=😓 Истекла}}", user, settings, subURL)
	if val2 != "Ваша подписка: ✅ Активна" {
		t.Errorf("expected 'Ваша подписка: ✅ Активна', got %q", val2)
	}

	// 3. Default RESET_STRATEGY
	val3 := formatTemplateValue("Сброс: {{RESET_STRATEGY}}", user, settings, subURL)
	if val3 != "Сброс: MONTH" {
		t.Errorf("expected 'Сброс: MONTH', got %q", val3)
	}

	// 4. Custom RESET_STRATEGY override
	val4 := formatTemplateValue("Сброс: {{RESET_STRATEGY:NO_RESET=не сбрасывается|MONTH=раз в месяц}}", user, settings, subURL)
	if val4 != "Сброс: раз в месяц" {
		t.Errorf("expected 'Сброс: раз в месяц', got %q", val4)
	}
}

func TestAllTemplateVariables(t *testing.T) {
	email := "test@example.com"
	tag := "VIP"
	tgID := int64(123456789)
	desc := "Premium user"
	hwidLimit := 3
	now := time.Now()
	expireAt := now.Add(24 * 10 * time.Hour)
	createdAt := now.Add(-24 * 5 * time.Hour)
	lastReset := now.Add(-24 * 2 * time.Hour)

	user := SubscriptionUser{
		TID:                  42,
		Username:             "johndoe",
		ShortUUID:            "xyz789",
		Status:               "ACTIVE",
		TrafficLimitStrategy: "MONTH",
		TrafficLimitBytes:    100 * 1024 * 1024 * 1024,
		UsedTrafficBytes:     25 * 1024 * 1024 * 1024,
		LifetimeUsedBytes:    500 * 1024 * 1024 * 1024,
		ExpireAt:             expireAt,
		CreatedAt:            createdAt,
		LastTrafficResetAt:   &lastReset,
		Email:                &email,
		Tag:                  &tag,
		TelegramID:           &tgID,
		Description:          &desc,
		HwidDeviceLimit:      &hwidLimit,
	}

	settings := SubscriptionSettingsParsed{}
	subURL := "https://sub.domain.com/xyz789"

	tests := []struct {
		input    string
		expected string
	}{
		{"{{DAYS_LEFT}}", "10"},
		{"{{TRAFFIC_USED}}", "25 GB"},
		{"{{TRAFFIC_LEFT}}", "75 GB"},
		{"{{TOTAL_TRAFFIC}}", "100 GB"},
		{"{{STATUS}}", "Active"},
		{"{{USERNAME}}", "johndoe"},
		{"{{EMAIL}}", "test@example.com"},
		{"{{TELEGRAM_ID}}", "123456789"},
		{"{{SUBSCRIPTION_URL}}", "https://sub.domain.com/xyz789"},
		{"{{TAG}}", "VIP"},
		{"{{SHORT_UUID}}", "xyz789"},
		{"{{ID}}", "42"},
		{"{{TRAFFIC_USED_BYTES}}", "26843545600"},
		{"{{TRAFFIC_LEFT_BYTES}}", "80530636800"},
		{"{{TOTAL_TRAFFIC_BYTES}}", "107374182400"},
		{"{{RESET_STRATEGY}}", "MONTH"},
		{"{{LIFETIME_USED_BYTES}}", "536870912000"},
		{"{{SS_HWID_LIMIT}}", "3"},
		{"{{DESCRIPTION}}", "Premium user"},
	}

	for _, tc := range tests {
		got := resolveTemplateVariables(tc.input, user, settings, subURL)
		if got != tc.expected {
			t.Errorf("template %q: expected %q, got %q", tc.input, tc.expected, got)
		}
	}
}

