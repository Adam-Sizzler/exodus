package panelsettings

import "testing"

func TestValidateAuthenticationSettingsRequiresAtLeastOneMethod(t *testing.T) {
	settings := map[string]any{
		"passkey_settings":  map[string]any{"enabled": false},
		"password_settings": map[string]any{"enabled": false},
		"oauth2_settings":   defaultOAuth2Settings(),
	}

	if err := validateAuthenticationSettings(settings); err == nil {
		t.Fatal("expected disabled authentication methods to be rejected")
	}
}

func TestValidateAuthenticationSettingsAcceptsPasswordOnly(t *testing.T) {
	settings := map[string]any{
		"passkey_settings":  map[string]any{"enabled": false},
		"password_settings": map[string]any{"enabled": true},
		"oauth2_settings":   defaultOAuth2Settings(),
	}

	if err := validateAuthenticationSettings(settings); err != nil {
		t.Fatalf("expected password-only auth to be valid, got %v", err)
	}
}

func TestValidateAuthenticationSettingsRejectsIncompletePasskey(t *testing.T) {
	settings := map[string]any{
		"passkey_settings":  map[string]any{"enabled": true, "rpId": "example.com", "origin": ""},
		"password_settings": map[string]any{"enabled": false},
		"oauth2_settings":   defaultOAuth2Settings(),
	}

	if err := validateAuthenticationSettings(settings); err == nil {
		t.Fatal("expected incomplete passkey settings to be rejected")
	}
}

func TestValidateAuthenticationSettingsValidatesOAuthProviders(t *testing.T) {
	oauth := defaultOAuth2Settings()
	oauth["github"].(map[string]any)["enabled"] = true
	oauth["github"].(map[string]any)["clientId"] = "client"
	oauth["github"].(map[string]any)["clientSecret"] = "secret"
	oauth["github"].(map[string]any)["allowedEmails"] = []any{"admin@example.com"}

	settings := map[string]any{
		"passkey_settings":  map[string]any{"enabled": false},
		"password_settings": map[string]any{"enabled": false},
		"oauth2_settings":   oauth,
	}

	if err := validateAuthenticationSettings(settings); err != nil {
		t.Fatalf("expected complete github settings to be valid, got %v", err)
	}

	oauth["github"].(map[string]any)["allowedEmails"] = []any{}
	if err := validateAuthenticationSettings(settings); err == nil {
		t.Fatal("expected github without allowed emails to be rejected")
	}
}

func TestValidateAuthenticationSettingsValidatesTelegram(t *testing.T) {
	oauth := defaultOAuth2Settings()
	oauth["telegram"].(map[string]any)["enabled"] = true
	oauth["telegram"].(map[string]any)["clientId"] = "client"
	oauth["telegram"].(map[string]any)["clientSecret"] = "secret"
	oauth["telegram"].(map[string]any)["frontendDomain"] = "panel.example.com"
	oauth["telegram"].(map[string]any)["allowedIds"] = []any{"12345"}

	settings := map[string]any{
		"passkey_settings":  map[string]any{"enabled": false},
		"password_settings": map[string]any{"enabled": false},
		"oauth2_settings":   oauth,
	}

	if err := validateAuthenticationSettings(settings); err != nil {
		t.Fatalf("expected complete telegram settings to be valid, got %v", err)
	}

	oauth["telegram"].(map[string]any)["allowedIds"] = []any{}
	if err := validateAuthenticationSettings(settings); err == nil {
		t.Fatal("expected telegram without allowed ids to be rejected")
	}
}
