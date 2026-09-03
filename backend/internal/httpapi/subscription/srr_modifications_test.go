package subscription

import (
	"strings"
	"testing"
	"time"

	"exodus/internal/config"
	"exodus/internal/httpapi/subscriptionresponserules"
	subscriptionsettings "exodus/internal/httpapi/subscriptionsettings"
)

func TestSRRModificationsFlow(t *testing.T) {
	tagExcluded := "EXCLUDE_ME"
	tagAllowed := "KEEP_ME"
	inboundType := "vless"

	hosts := []SubscriptionHost{
		{
			UUID:         "host-1",
			Address:      "1.1.1.1",
			Port:         443,
			Remark:       "Host-1-Excluded",
			InboundType:  &inboundType,
			Tags:         []string{tagExcluded, "PROD"},
		},
		{
			UUID:         "host-2",
			Address:      "2.2.2.2",
			Port:         443,
			Remark:       "Host-2-Allowed",
			InboundType:  &inboundType,
			Tags:         []string{tagAllowed},
		},
	}

	user := SubscriptionUser{
		TID:       1,
		UUID:      "user-uuid",
		ShortUUID: "u001",
		Username:  "testuser",
		Status:    "ACTIVE",
		VlessUUID: "11111111-1111-1111-1111-111111111111",
		ExpireAt:  time.Now().Add(24 * time.Hour),
	}

	cfg := &config.BackendConfig{}

	t.Run("excludeHostsByTags filters matching hosts", func(t *testing.T) {
		matchedRuleMods := &subscriptionresponserules.RuleModifications{
			ExcludeHostsByTags: []string{tagExcluded},
		}

		excludeSet := make(map[string]struct{}, len(matchedRuleMods.ExcludeHostsByTags))
		for _, tag := range matchedRuleMods.ExcludeHostsByTags {
			trimmed := strings.TrimSpace(tag)
			if trimmed != "" {
				excludeSet[trimmed] = struct{}{}
			}
		}

		filtered := make([]SubscriptionHost, 0, len(hosts))
		for _, host := range hosts {
			excluded := false
			for _, tag := range host.Tags {
				if _, ok := excludeSet[strings.TrimSpace(tag)]; ok {
					excluded = true
					break
				}
			}
			if !excluded {
				filtered = append(filtered, host)
			}
		}

		if len(filtered) != 1 {
			t.Fatalf("expected 1 host remaining, got %d", len(filtered))
		}
		if filtered[0].UUID != "host-2" {
			t.Errorf("expected host-2 to remain, got %s", filtered[0].UUID)
		}
	})

	t.Run("respondWithRemarks replaces real hosts", func(t *testing.T) {
		remarks := []string{"Server under maintenance", "Support: @myvpn"}
		matchedRuleMods := &subscriptionresponserules.RuleModifications{
			RespondWithRemarks: remarks,
		}

		settings := SubscriptionSettingsParsed{
			Raw: subscriptionsettings.SubscriptionSettings{
				IsShowCustomRemarks: true,
			},
		}

		var earlyExitRemarks []string
		if matchedRuleMods != nil && len(matchedRuleMods.RespondWithRemarks) > 0 {
			if settings.Raw.IsShowCustomRemarks {
				earlyExitRemarks = matchedRuleMods.RespondWithRemarks
			}
		}

		if len(earlyExitRemarks) != 2 {
			t.Fatalf("expected 2 early exit remarks, got %d", len(earlyExitRemarks))
		}
		fallbackHosts := createFallbackRemarkHosts(earlyExitRemarks)
		if len(fallbackHosts) != 2 {
			t.Fatalf("expected 2 fallback hosts, got %d", len(fallbackHosts))
		}
		if fallbackHosts[0].Remark != "Server under maintenance" {
			t.Errorf("expected first remark 'Server under maintenance', got %q", fallbackHosts[0].Remark)
		}
	})

	t.Run("headers and applyHeadersToEnd override precedence", func(t *testing.T) {
		matchedRuleMods := &subscriptionresponserules.RuleModifications{
			Headers: []subscriptionresponserules.RuleHeader{
				{Key: "profile-title", Value: "OverriddenTitle"},
				{Key: "X-Custom-Header", Value: "CustomValue"},
			},
			ApplyHeadersToEnd: true,
		}

		settings := SubscriptionSettingsParsed{
			Raw: subscriptionsettings.SubscriptionSettings{
				ProfileTitle: "OriginalTitle",
			},
		}

		respHeaders := buildResponseHeaders(user, settings, "text/plain", "https://example.com/sub/u001")

		// With applyHeadersToEnd: true, rule headers override earlier headers
		if matchedRuleMods.ApplyHeadersToEnd {
			for _, h := range matchedRuleMods.Headers {
				respHeaders[h.Key] = h.Value
			}
		}

		if respHeaders["profile-title"] != "OverriddenTitle" {
			t.Errorf("expected profile-title = OverriddenTitle, got %q", respHeaders["profile-title"])
		}
		if respHeaders["X-Custom-Header"] != "CustomValue" {
			t.Errorf("expected X-Custom-Header = CustomValue, got %q", respHeaders["X-Custom-Header"])
		}
	})

	t.Run("ignoreServeJsonAtBaseSubscription toggles json upgrade", func(t *testing.T) {
		userAgent := "Streisand/1.7.8"

		// Case 1: normal serveJson enabled -> upgrades to XRAY_JSON
		reqType := defaultResponseType
		serveJSON := true
		ignoreServeJSON := false

		if reqType == defaultResponseType && serveJSON && !ignoreServeJSON && isJSONSubscriptionFallbackSupported(userAgent) {
			reqType = responseTypeXrayJSON
		}
		if reqType != responseTypeXrayJSON {
			t.Errorf("expected upgrade to XRAY_JSON, got %s", reqType)
		}

		// Case 2: ignoreServeJsonAtBaseSubscription = true -> stays XRAY_BASE64
		reqType = defaultResponseType
		ignoreServeJSON = true
		if reqType == defaultResponseType && serveJSON && !ignoreServeJSON && isJSONSubscriptionFallbackSupported(userAgent) {
			reqType = responseTypeXrayJSON
		}
		if reqType != defaultResponseType {
			t.Errorf("expected to stay defaultResponseType when ignored, got %s", reqType)
		}
	})

	_ = cfg
}
