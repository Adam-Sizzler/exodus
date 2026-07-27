package auth

import (
	"net/http/httptest"
	"testing"
)

func TestExternalLoginNotificationDataUsesResolvedClientIP(t *testing.T) {
	req := httptest.NewRequest("POST", "/api/auth/oauth2/callback", nil)
	req.RemoteAddr = "172.18.0.8:47242"
	req.Header.Set("X-Forwarded-For", "74.2.54.120, 172.18.0.1")
	req.Header.Set("X-Real-IP", "172.18.0.1")
	req.Header.Set("User-Agent", "Mozilla/5.0")

	data := externalLoginNotificationData(nil, "oauth2", "telegram", "406150372", "admin-uuid", "", req)

	if got := data["ip"]; got != "74.2.54.120" {
		t.Fatalf("data ip got %v, want %q", got, "74.2.54.120")
	}
	loginAttempt, ok := data["loginAttempt"].(map[string]any)
	if !ok {
		t.Fatalf("loginAttempt missing or has invalid type: %T", data["loginAttempt"])
	}
	if got := loginAttempt["ip"]; got != "74.2.54.120" {
		t.Fatalf("loginAttempt ip got %v, want %q", got, "74.2.54.120")
	}
	if _, ok := loginAttempt["remoteAddr"]; ok {
		t.Fatalf("loginAttempt must not expose remoteAddr: %v", loginAttempt["remoteAddr"])
	}
	if got := loginAttempt["username"]; got != "406150372" {
		t.Fatalf("loginAttempt username got %v, want %q", got, "406150372")
	}
}

func TestExternalLoginNotificationUsernameFallback(t *testing.T) {
	if got := externalLoginNotificationUsername("oauth2", "telegram", ""); got != "oauth2:telegram" {
		t.Fatalf("username fallback got %q, want %q", got, "oauth2:telegram")
	}
	if got := externalLoginNotificationUsername("oauth2", "telegram", " 406150372 "); got != "406150372" {
		t.Fatalf("username identifier got %q, want %q", got, "406150372")
	}
}
