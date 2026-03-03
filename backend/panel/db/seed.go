package db

import (
	"context"
	"database/sql"
	"fmt"

	"v2ray-stat/backend/panel/config"
	"v2ray-stat/backend/panel/dbutil"

	"github.com/google/uuid"
)

const defaultResponseRules = `{"rules":[{"name":"Browser Subscription","enabled":true,"operator":"AND","conditions":[{"value":"text/html","operator":"CONTAINS","headerName":"accept","caseSensitive":true}],"description":"System critical: do not delete or disable this rule.","responseType":"BROWSER"},{"name":"Mihomo Clients","enabled":true,"operator":"AND","conditions":[{"value":"^(?:FlClash|FlClashX|Flowvy|[Cc]lash-[Vv]erge|[Kk]oala-[Cc]lash|[Cc]lash-?[Mm]eta|[Mm]urge|[Cc]lashX [Mm]eta|[Mm]ihomo|[Cc]lash-nyanpasu|clash.meta|prizrak-box)","operator":"REGEX","headerName":"user-agent","caseSensitive":false}],"description":"Response with generated YAML config (Mihomo Template)","responseType":"MIHOMO"},{"name":"Stash (iOS, macOS)","enabled":true,"operator":"AND","conditions":[{"value":"^stash","operator":"REGEX","headerName":"user-agent","caseSensitive":false}],"description":"Response with generated YAML config (Stash Template)","responseType":"STASH"},{"name":"Sing-box clients","enabled":true,"operator":"AND","conditions":[{"value":"^sfa|sfi|sfm|sft|karing|singbox|rabbithole","operator":"REGEX","headerName":"user-agent","caseSensitive":false}],"description":"Resonse with generated JSON config (Singbox template)","responseType":"SINGBOX"},{"name":"Clash Core Clients","enabled":true,"operator":"AND","conditions":[{"value":"^clash","operator":"REGEX","headerName":"user-agent","caseSensitive":false}],"description":"Response with generated YAML config (Clash Template)","responseType":"CLASH"},{"name":"Fallback Base64","enabled":true,"operator":"AND","conditions":[],"description":"System critical: do not delete or disable this rule.","responseType":"XRAY_BASE64"}],"version":"1"}`
const defaultHWIDSettings = `{"enabled":false,"maxDevicesAnnounce":null,"fallbackDeviceLimit":999}`
const defaultCustomRemarks = `{"emptyHosts":["→ v2rs","→ No hosts found","→ Check Hosts tab","→ Check Internal Squads tab"],"expiredUsers":["⌛ Subscription expired","Contact support"],"limitedUsers":["🚧 Subscription limited","Contact support"],"disabledUsers":["🚫 Subscription disabled","Contact support"],"HWIDNotSupported":["App not supported"],"HWIDMaxDevicesExceeded":["Limit of devices reached"]}`
const defaultSubpageConfigUUID = "00000000-0000-0000-0000-000000000000"

// SeedDefaults inserts base settings and templates if they do not exist.
func SeedDefaults(ctx context.Context, dbConn *sql.DB, cfg *config.BackendConfig) error {
	if dbConn == nil {
		return fmt.Errorf("database connection is nil")
	}

	tx, err := dbConn.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin defaults transaction: %w", err)
	}

	if err := ensureV2rsSettings(tx); err != nil {
		_ = tx.Rollback()
		return err
	}
	if err := ensureDefaultSubscriptionSettings(tx); err != nil {
		_ = tx.Rollback()
		return err
	}
	if err := ensureDefaultTemplates(tx); err != nil {
		_ = tx.Rollback()
		return err
	}
	if err := ensureDefaultSubscriptionPageConfig(tx); err != nil {
		_ = tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit defaults transaction: %w", err)
	}

	cfg.Logger.Info("Default configuration data ensured")
	return nil
}

func ensureV2rsSettings(tx *sql.Tx) error {
	var exists int
	if err := tx.QueryRowContext(context.Background(), `SELECT COUNT(*) FROM v2rs_settings WHERE id = 1`).Scan(&exists); err != nil {
		return fmt.Errorf("check v2rs_settings row: %w", err)
	}
	if exists > 0 {
		return nil
	}

	passkeySettings := `{"rpId":null,"origin":null,"enabled":false}`
	oauth2Settings := `{"github":{"enabled":false,"clientId":null,"clientSecret":null,"allowedEmails":[]},"yandex":{"enabled":false,"clientId":null,"clientSecret":null,"allowedEmails":[]},"generic":{"enabled":false,"clientId":null,"tokenUrl":null,"withPkce":false,"clientSecret":null,"allowedEmails":[],"frontendDomain":null,"authorizationUrl":null},"keycloak":{"realm":null,"enabled":false,"clientId":null,"clientSecret":null,"allowedEmails":[],"frontendDomain":null,"keycloakDomain":null},"pocketid":{"enabled":false,"clientId":null,"plainDomain":null,"clientSecret":null,"allowedEmails":[]}}`
	tgAuthSettings := `{"enabled":false,"adminIds":[],"botToken":null}`
	passwordSettings := `{"enabled":true}`
	brandingSettings := `{"title":"V2RS","logoUrl":null}`

	query := dbutil.Rebind(`
		INSERT INTO v2rs_settings (
			id, passkey_settings, oauth2_settings, tg_auth_settings, password_settings, branding_settings
		) VALUES (1, ?, ?, ?, ?, ?)
	`)
	if _, err := tx.ExecContext(context.Background(), query, passkeySettings, oauth2Settings, tgAuthSettings, passwordSettings, brandingSettings); err != nil {
		return fmt.Errorf("insert default v2rs_settings row: %w", err)
	}

	return nil
}

func ensureDefaultSubscriptionSettings(tx *sql.Tx) error {
	var count int
	if err := tx.QueryRowContext(context.Background(), `SELECT COUNT(*) FROM subscription_settings`).Scan(&count); err != nil {
		return fmt.Errorf("count subscription_settings: %w", err)
	}
	if count > 0 {
		return nil
	}

	query := dbutil.Rebind(`
		INSERT INTO subscription_settings (
			uuid, profile_title, support_link, profile_update_interval,
			address, port, api_schema, api_path,
			is_profile_webpage_url_enabled, serve_json_at_base_subscription,
			happ_announce, happ_routing, is_show_custom_remarks,
			custom_remarks, custom_response_headers, randomize_hosts,
			response_rules, hwid_settings
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`)
	_, err := tx.ExecContext(context.Background(), query,
		uuid.NewString(),
		"v2rs",
		"https://github.com",
		12,
		"",
		9263,
		"grpc",
		"",
		true,
		false,
		"",
		"",
		true,
		defaultCustomRemarks,
		`{}`,
		false,
		defaultResponseRules,
		defaultHWIDSettings,
	)
	if err != nil {
		return fmt.Errorf("insert default subscription_settings: %w", err)
	}

	return nil
}

func ensureDefaultTemplates(tx *sql.Tx) error {
	for _, tmpl := range DefaultSubscriptionTemplates() {
		var count int
		if err := tx.QueryRowContext(context.Background(), dbutil.Rebind(`SELECT COUNT(*) FROM subscription_templates WHERE template_type = ?`), tmpl.TemplateType).Scan(&count); err != nil {
			return fmt.Errorf("count template %s: %w", tmpl.TemplateType, err)
		}
		if count > 0 {
			continue
		}

		query := dbutil.Rebind(`
			INSERT INTO subscription_templates (
				uuid, view_position, name, template_type, template_yaml, template_json
			) VALUES (?, ?, ?, ?, ?, ?)
		`)
		if _, err := tx.ExecContext(context.Background(), query,
			uuid.NewString(),
			tmpl.ViewPosition,
			tmpl.Name,
			tmpl.TemplateType,
			tmpl.TemplateYAML,
			tmpl.TemplateJSON,
		); err != nil {
			return fmt.Errorf("insert template %s: %w", tmpl.TemplateType, err)
		}
	}

	return nil
}

func ensureDefaultSubscriptionPageConfig(tx *sql.Tx) error {
	var count int
	if err := tx.QueryRowContext(context.Background(), `SELECT COUNT(*) FROM subscription_page_config`).Scan(&count); err != nil {
		return fmt.Errorf("count subscription_page_config: %w", err)
	}
	if count > 0 {
		return nil
	}

	query := dbutil.Rebind(`
		INSERT INTO subscription_page_config (
			uuid, view_position, name, config
		) VALUES (?, ?, ?, ?)
	`)
	if _, err := tx.ExecContext(context.Background(), query, defaultSubpageConfigUUID, 1, "Default", DefaultSubscriptionPageConfig); err != nil {
		return fmt.Errorf("insert default subscription_page_config: %w", err)
	}

	return nil
}

func strPtr(v string) *string {
	if v == "" {
		return nil
	}
	return &v
}
