package auth

import (
	"net/http/httptest"
	"testing"
)

func TestExternalLoginNotificationDataUsesForwardedIP(t *testing.T) {
	req := httptest.NewRequest("POST", "/api/auth/oauth2/callback", nil)
	req.RemoteAddr = "172.18.0.8:47242"
	req.Header.Set("X-Real-IP", "144.31.119.150")
	req.Header.Set("User-Agent", "Mozilla/5.0")

	data := externalLoginNotificationData("oauth2", "telegram", "306100972", "admin-uuid", "", req)

	if got := data["ip"]; got != "144.31.119.150" {
		t.Fatalf("data ip got %v, want %q", got, "144.31.119.150")
	}
	loginAttempt, ok := data["loginAttempt"].(map[string]any)
	if !ok {
		t.Fatalf("loginAttempt missing or has invalid type: %T", data["loginAttempt"])
	}
	if got := loginAttempt["ip"]; got != "144.31.119.150" {
		t.Fatalf("loginAttempt ip got %v, want %q", got, "144.31.119.150")
	}
	if got := loginAttempt["remoteAddr"]; got != "172.18.0.8:47242" {
		t.Fatalf("loginAttempt remoteAddr got %v, want %q", got, "172.18.0.8:47242")
	}
	if got := loginAttempt["username"]; got != "306100972" {
		t.Fatalf("loginAttempt username got %v, want %q", got, "306100972")
	}
}

func TestExternalLoginNotificationUsernameFallback(t *testing.T) {
	if got := externalLoginNotificationUsername("oauth2", "telegram", ""); got != "oauth2:telegram" {
		t.Fatalf("username fallback got %q, want %q", got, "oauth2:telegram")
	}
	if got := externalLoginNotificationUsername("oauth2", "telegram", " 306100972 "); got != "306100972" {
		t.Fatalf("username identifier got %q, want %q", got, "306100972")
	}
}
