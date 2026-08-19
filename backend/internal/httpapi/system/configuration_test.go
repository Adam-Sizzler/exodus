package system

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"exodus/internal/config"
)

func TestConfigurationHandler(t *testing.T) {
	cfg := &config.BackendConfig{
		Backend: config.BackendAppConfig{
			ShortUUIDLength: 12,
		},
		Notifications: config.NotificationsConfig{
			WebhookEnabled: true,
		},
		Scheduler: config.SchedulerConfig{
			ServiceCleanUsageHistory:                 true,
			BandwidthUsageNotificationsEnabled:       true,
			BandwidthUsageNotificationsThreshold:     []int{80, 90},
			NotConnectedUsersNotificationsEnabled:    true,
			NotConnectedUsersNotificationsAfterHours: []int{24, 48},
			ExpirationNotificationsEnabled:           true,
			ExpirationNotifications:                  []int{1, 3},
		},
		Redis: config.RedisConfig{
			ExportToStreamEnabled: true,
		},
	}

	handler := ConfigurationHandler(nil, cfg)

	req := httptest.NewRequest(http.MethodGet, "/api/system/configuration", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	var resp struct {
		Response SystemConfigurationResponse `json:"response"`
	}
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if !resp.Response.Notifications.Webhook {
		t.Errorf("expected Webhook to be true")
	}
	if !resp.Response.Notifications.BandwidthUsage {
		t.Errorf("expected BandwidthUsage to be true")
	}
	if !resp.Response.Notifications.NotConnectedAfter {
		t.Errorf("expected NotConnectedAfter to be true")
	}
	if !resp.Response.Notifications.ExpirationNotifications {
		t.Errorf("expected ExpirationNotifications to be true")
	}
	if !resp.Response.Service.CleanUsageHistory {
		t.Errorf("expected CleanUsageHistory to be true")
	}
	if !resp.Response.Service.ExportToRedisStream {
		t.Errorf("expected ExportToRedisStream to be true")
	}
	if resp.Response.Misc.ShortUUIDLength != 12 {
		t.Errorf("expected ShortUUIDLength 12, got %d", resp.Response.Misc.ShortUUIDLength)
	}
}
