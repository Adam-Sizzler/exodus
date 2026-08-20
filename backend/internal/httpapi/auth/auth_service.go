package auth

import (
	"database/sql"
	"encoding/json"
	"errors"
	"strings"

	"exodus/internal/config"
	"exodus/internal/panelsettings"
)

func getBootstrapData(db *sql.DB) (brandingSettings map[string]any, passwordSettings map[string]any, defaultUsername string, hasAdmin bool, err error) {
	brandingSettings = panelsettings.DefaultBrandingSettings()
	passwordSettings = panelsettings.DefaultPasswordSettings()
	defaultUsername = "admin"
	hasAdmin = false

	row := db.QueryRow(`
		SELECT branding_settings, password_settings
		FROM exodus_settings
		WHERE id = 1
		LIMIT 1
	`)

	var brandingRaw, passwordRaw sql.NullString
	if scanErr := row.Scan(&brandingRaw, &passwordRaw); scanErr != nil && !errors.Is(scanErr, sql.ErrNoRows) {
		return brandingSettings, passwordSettings, defaultUsername, hasAdmin, scanErr
	}

	if brandingRaw.Valid && strings.TrimSpace(brandingRaw.String) != "" {
		var tmp map[string]any
		if json.Unmarshal([]byte(brandingRaw.String), &tmp) == nil && len(tmp) > 0 {
			brandingSettings = tmp
		}
	}
	if passwordRaw.Valid && strings.TrimSpace(passwordRaw.String) != "" {
		var tmp map[string]any
		if json.Unmarshal([]byte(passwordRaw.String), &tmp) == nil && len(tmp) > 0 {
			passwordSettings = tmp
		}
	}

	var adminCount int
	if countErr := db.QueryRow("SELECT COUNT(*) FROM admin").Scan(&adminCount); countErr != nil {
		return brandingSettings, passwordSettings, defaultUsername, hasAdmin, countErr
	}
	hasAdmin = adminCount > 0

	if hasAdmin {
		var firstUsername sql.NullString
		if firstErr := db.QueryRow("SELECT username FROM admin ORDER BY created_at ASC LIMIT 1").Scan(&firstUsername); firstErr != nil && !errors.Is(firstErr, sql.ErrNoRows) {
			return brandingSettings, passwordSettings, defaultUsername, hasAdmin, firstErr
		}
		if firstUsername.Valid && strings.TrimSpace(firstUsername.String) != "" {
			defaultUsername = firstUsername.String
		}
	}
	return brandingSettings, passwordSettings, defaultUsername, hasAdmin, nil
}

func getAuthMethodsStatus(db *sql.DB) (passkeyEnabled bool, oauth2Providers map[string]bool) {
	passkeyEnabled = false
	oauth2Providers = map[string]bool{
		"github":   false,
		"yandex":   false,
		"generic":  false,
		"keycloak": false,
		"pocketid": false,
		"telegram": false,
	}

	row := db.QueryRow(`
		SELECT passkey_settings, oauth2_settings
		FROM exodus_settings
		WHERE id = 1
		LIMIT 1
	`)

	var passkeyRaw, oauth2Raw sql.NullString
	if err := row.Scan(&passkeyRaw, &oauth2Raw); err != nil {
		return passkeyEnabled, oauth2Providers
	}

	if passkeyRaw.Valid && strings.TrimSpace(passkeyRaw.String) != "" {
		var passkeyObj map[string]any
		if json.Unmarshal([]byte(passkeyRaw.String), &passkeyObj) == nil {
			if enabled, ok := passkeyObj["enabled"].(bool); ok {
				passkeyEnabled = enabled
			}
		}
	}

	if oauth2Raw.Valid && strings.TrimSpace(oauth2Raw.String) != "" {
		var oauthObj map[string]map[string]any
		if json.Unmarshal([]byte(oauth2Raw.String), &oauthObj) == nil {
			for provider := range oauth2Providers {
				if providerCfg, ok := oauthObj[provider]; ok {
					if enabled, ok := providerCfg["enabled"].(bool); ok {
						oauth2Providers[provider] = enabled
					}
				}
			}
		}
	}

	return passkeyEnabled, oauth2Providers
}

func resolvePasswordAuthEnabled(
	passwordSettings map[string]any,
	passkeyEnabled bool,
	oauth2Providers map[string]bool,
	cfg *config.BackendConfig,
) bool {
	passwordEnabled := true
	if raw, ok := passwordSettings["enabled"]; ok {
		if value, ok := raw.(bool); ok {
			passwordEnabled = value
		}
	}

	if passwordEnabled {
		return true
	}

	hasOAuth2Enabled := false
	for _, enabled := range oauth2Providers {
		if enabled {
			hasOAuth2Enabled = true
			break
		}
	}

	if !passkeyEnabled && !hasOAuth2Enabled {
		if cfg != nil && cfg.Logger != nil {
			cfg.Logger.Warn("All authentication methods are disabled. Falling back to password auth enabled=true to prevent lockout.")
		}
		return true
	}

	return false
}
