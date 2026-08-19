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
