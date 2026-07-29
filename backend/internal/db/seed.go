package db

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"exodus/internal/config"
	"exodus/internal/jobqueue"
	"exodus/internal/security"

	"github.com/google/uuid"
)

const defaultResponseRules = `{"rules":[{"name":"Browser Subscription","enabled":true,"operator":"AND","conditions":[{"value":"text/html","operator":"CONTAINS","headerName":"accept","caseSensitive":true}],"description":"System critical: do not delete or disable this rule.","responseType":"BROWSER"},{"name":"Mihomo Clients","enabled":true,"operator":"AND","conditions":[{"value":"^(?:FlClash|FlClashX|Flowvy|[Cc]lash-[Vv]erge|[Kk]oala-[Cc]lash|[Cc]lash-?[Mm]eta|[Mm]urge|[Cc]lashX [Mm]eta|[Mm]ihomo|[Cc]lash-nyanpasu|clash.meta|prizrak-box)","operator":"REGEX","headerName":"user-agent","caseSensitive":false}],"description":"Response with generated YAML config (Mihomo Template)","responseType":"MIHOMO"},{"name":"Stash (iOS, macOS)","enabled":true,"operator":"AND","conditions":[{"value":"^stash","operator":"REGEX","headerName":"user-agent","caseSensitive":false}],"description":"Response with generated YAML config (Stash Template)","responseType":"STASH"},{"name":"Sing-box clients","enabled":true,"operator":"AND","conditions":[{"value":"^sfa|sfi|sfm|sft|karing|singbox|rabbithole","operator":"REGEX","headerName":"user-agent","caseSensitive":false}],"description":"Resonse with generated JSON config (Singbox template)","responseType":"SINGBOX"},{"name":"Clash Core Clients","enabled":true,"operator":"AND","conditions":[{"value":"^clash","operator":"REGEX","headerName":"user-agent","caseSensitive":false}],"description":"Response with generated YAML config (Clash Template)","responseType":"CLASH"},{"name":"Fallback Base64","enabled":true,"operator":"AND","conditions":[],"description":"System critical: do not delete or disable this rule.","responseType":"XRAY_BASE64"}],"version":"1"}`

// PrevResponseRulesHash stores the canonical SHA-256 hash of the PREVIOUS defaultResponseRules.
// IMPORTANT: When updating defaultResponseRules, calculate its current canonical hash and set it as PrevResponseRulesHash FIRST before updating defaultResponseRules.
const PrevResponseRulesHash = "4e61e34171f1ad37b23b2ef9dc22a441dca4beef75bc98fe2a54fe3685c52ad8"

const defaultHWIDSettings = `{"enabled":false,"maxDevicesAnnounce":null,"fallbackDeviceLimit":999}`
const defaultCustomRemarks = `{"emptyHosts":["→ exodus","→ No hosts found","→ Check Hosts tab","→ Check Internal Squads tab"],"expiredUsers":["⌛ Subscription expired","Contact support"],"limitedUsers":["🚧 Subscription limited","Contact support"],"disabledUsers":["🚫 Subscription disabled","Contact support"],"HWIDNotSupported":["App not supported"],"HWIDMaxDevicesExceeded":["Limit of devices reached"]}`
const defaultSubpageConfigUUID = "00000000-0000-0000-0000-000000000000"
const defaultSingboxConfig = `{"log":{"level":"info"},"dns":{"servers":[{"tag":"dns-remote","type":"udp","server":"1.1.1.1","detour":"direct"}]},"inbounds":[{"type":"shadowsocks","tag":"ss-in","listen":"127.0.0.1","listen_port":2080,"method":"chacha20-ietf-poly1305"},{"type":"trojan","tag":"trojan-in","listen":"127.0.0.1","listen_port":2443,"users":[]}],"outbounds":[{"type":"direct","tag":"direct"},{"type":"block","tag":"block"}],"route":{"final":"direct"}}`

func canonicalHash(rawJSON string) (string, error) {
	var v any
	if err := json.Unmarshal([]byte(rawJSON), &v); err != nil {
		return "", err
	}
	canonical, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(canonical)
	return hex.EncodeToString(sum[:]), nil
}

func clearRedis(ctx context.Context, cfg *config.BackendConfig) error {
	client, err := jobqueue.NewRedisClient(cfg)
	if err != nil || client == nil {
		return err
	}
	defer client.Close()
	return client.FlushDB(ctx).Err()
}

// SeedDefaults inserts base settings and templates if they do not exist.
func SeedDefaults(ctx context.Context, dbConn *sql.DB, cfg *config.BackendConfig) error {
	if dbConn == nil {
		return fmt.Errorf("database connection is nil")
	}

	if err := clearRedis(ctx, cfg); err != nil {
		cfg.Logger.Warn("Failed to clear Redis on startup", "error", err)
	}

	tx, err := dbConn.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin defaults transaction: %w", err)
	}

	// 2. Checkup external squads
	if n, err := checkupExternalSquads(ctx, tx); err != nil {
		_ = tx.Rollback()
		return err
	} else if n > 0 {
		cfg.Logger.Info("External squads checkup completed", "fixedRows", n)
	}

	// 3. Ensure Exodus settings
	if err := ensureExodusSettings(ctx, tx, cfg); err != nil {
		_ = tx.Rollback()
		return err
	}

	// 4. Ensure default templates
	if err := ensureDefaultTemplates(ctx, tx, cfg); err != nil {
		_ = tx.Rollback()
		return err
	}

	// 5. Ensure default config profile
	if err := ensureDefaultConfigProfile(ctx, tx, cfg); err != nil {
		_ = tx.Rollback()
		return err
	}

	// 6. Resync config profile inbounds
	if n, err := resyncConfigProfileInbounds(ctx, tx); err != nil {
		_ = tx.Rollback()
		return err
	} else if n > 0 {
		cfg.Logger.Info("Config profile inbounds resynced", "profiles", n)
	}

	// 7. Ensure default internal squad
	if err := ensureDefaultInternalSquad(ctx, tx, cfg); err != nil {
		_ = tx.Rollback()
		return err
	}

	// 8. Ensure default subscription settings
	if err := ensureDefaultSubscriptionSettings(ctx, tx, cfg); err != nil {
		_ = tx.Rollback()
		return err
	}

	// 11. Ensure default subscription page config
	if err := ensureDefaultSubscriptionPageConfig(ctx, tx); err != nil {
		_ = tx.Rollback()
		return err
	}

	// 9. Ensure keygen
	if err := ensureKeygen(ctx, tx); err != nil {
		_ = tx.Rollback()
		return err
	}

	// 12. Ensure single admin
	if err := ensureSingleAdmin(ctx, tx, cfg); err != nil {
		_ = tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit defaults transaction: %w", err)
	}

	cfg.Logger.Info("Default configuration data ensured")
	return nil
}

func checkupExternalSquads(ctx context.Context, tx *sql.Tx) (int64, error) {
	res, err := tx.ExecContext(ctx, `
		UPDATE external_squads SET
			subscription_settings = CASE 
				WHEN subscription_settings::text IN ('{}', 'null', '[]') THEN NULL 
				ELSE subscription_settings 
			END,
			host_overrides = CASE 
				WHEN host_overrides::text IN ('{}', 'null', '[]') THEN NULL 
				ELSE host_overrides 
			END,
			response_headers = CASE 
				WHEN response_headers::text IN ('{}', 'null', '[]') THEN NULL 
				ELSE response_headers 
			END,
			hwid_settings = CASE 
				WHEN hwid_settings::text IN ('{}', 'null', '[]') OR (hwid_settings IS NOT NULL AND jsonb_typeof(hwid_settings) != 'object') THEN NULL 
				ELSE hwid_settings 
			END,
			custom_remarks = CASE 
				WHEN custom_remarks::text IN ('{}', 'null', '[]') OR (custom_remarks IS NOT NULL AND jsonb_typeof(custom_remarks) != 'object') THEN NULL 
				ELSE custom_remarks 
			END
		WHERE 
			subscription_settings::text IN ('{}', 'null', '[]')
			OR host_overrides::text IN ('{}', 'null', '[]')
			OR response_headers::text IN ('{}', 'null', '[]')
			OR hwid_settings::text IN ('{}', 'null', '[]') OR (hwid_settings IS NOT NULL AND jsonb_typeof(hwid_settings) != 'object')
			OR custom_remarks::text IN ('{}', 'null', '[]') OR (custom_remarks IS NOT NULL AND jsonb_typeof(custom_remarks) != 'object')
	`)
	if err != nil {
		return 0, fmt.Errorf("checkup external squads: %w", err)
	}
	return res.RowsAffected()
}

type passkeySettingsStruct struct {
	Enabled bool    `json:"enabled"`
	RpID    *string `json:"rpId"`
	Origin  *string `json:"origin"`
}

type passwordSettingsStruct struct {
	Enabled bool `json:"enabled"`
}

func ensureExodusSettings(ctx context.Context, tx *sql.Tx, cfg *config.BackendConfig) error {
	defaultPasskeySettings := `{"enabled":false,"origin":null,"rpId":null}`
	defaultOAuth2Settings := `{"generic":{"allowedEmails":[],"authorizationUrl":null,"clientId":null,"clientSecret":null,"enabled":false,"frontendDomain":null,"tokenUrl":null,"withPkce":false},"github":{"allowedEmails":[],"clientId":null,"clientSecret":null,"enabled":false},"keycloak":{"allowedEmails":[],"clientId":null,"clientSecret":null,"enabled":false,"frontendDomain":null,"keycloakDomain":null,"realm":null},"pocketid":{"allowedEmails":[],"clientId":null,"clientSecret":null,"enabled":false,"plainDomain":null},"telegram":{"allowedIds":[],"clientId":null,"clientSecret":null,"enabled":false,"frontendDomain":null},"yandex":{"allowedEmails":[],"clientId":null,"clientSecret":null,"enabled":false}}`
	defaultPasswordSettings := `{"enabled":true}`
	defaultBrandingSettings := `{"logoUrl":null,"title":"EXODUS"}`

	var (
		passkeyRaw  sql.NullString
		oauth2Raw   sql.NullString
		passwordRaw sql.NullString
		brandingRaw sql.NullString
	)

	err := tx.QueryRowContext(ctx, `SELECT passkey_settings::text, oauth2_settings::text, password_settings::text, branding_settings::text FROM exodus_settings WHERE id = 1`).Scan(&passkeyRaw, &oauth2Raw, &passwordRaw, &brandingRaw)
	if errors.Is(err, sql.ErrNoRows) {
		query := `
			INSERT INTO exodus_settings (
				id, passkey_settings, oauth2_settings, password_settings, branding_settings
			) VALUES (1, $1, $2, $3, $4)
		`
		if _, err := tx.ExecContext(ctx, query, defaultPasskeySettings, defaultOAuth2Settings, defaultPasswordSettings, defaultBrandingSettings); err != nil {
			return fmt.Errorf("insert default exodus_settings row: %w", err)
		}
		cfg.Logger.Info("Seeded exodus_settings default row")
		return nil
	} else if err != nil {
		return fmt.Errorf("check exodus_settings row: %w", err)
	}

	var resetFields []string
	updates := make(map[string]string)

	// 1. Validate passkey_settings
	passkeyValid := false
	if passkeyRaw.Valid && strings.TrimSpace(passkeyRaw.String) != "" && passkeyRaw.String != "null" {
		var pk passkeySettingsStruct
		if err := json.Unmarshal([]byte(passkeyRaw.String), &pk); err == nil {
			passkeyValid = true
		}
	}
	if !passkeyValid {
		updates["passkey_settings"] = defaultPasskeySettings
		resetFields = append(resetFields, "passkey_settings")
		// Special case: if passkey_settings is invalid, force password_settings.enabled = true
		updates["password_settings"] = defaultPasswordSettings
		resetFields = append(resetFields, "password_settings (forced enabled)")
	}

	// 2. Validate oauth2_settings
	oauth2Valid := false
	if oauth2Raw.Valid && strings.TrimSpace(oauth2Raw.String) != "" && oauth2Raw.String != "null" {
		var m map[string]json.RawMessage
		if err := json.Unmarshal([]byte(oauth2Raw.String), &m); err == nil {
			oauth2Valid = true
		}
	}
	if !oauth2Valid {
		updates["oauth2_settings"] = defaultOAuth2Settings
		resetFields = append(resetFields, "oauth2_settings")
	}

	// 3. Validate password_settings (if not already queued for reset)
	if _, alreadyReset := updates["password_settings"]; !alreadyReset {
		passwordValid := false
		if passwordRaw.Valid && strings.TrimSpace(passwordRaw.String) != "" && passwordRaw.String != "null" {
			var pw passwordSettingsStruct
			if err := json.Unmarshal([]byte(passwordRaw.String), &pw); err == nil {
				passwordValid = true
			}
		}
		if !passwordValid {
			updates["password_settings"] = defaultPasswordSettings
			resetFields = append(resetFields, "password_settings")
		}
	}

	// 4. Validate branding_settings
	brandingValid := false
	if brandingRaw.Valid && strings.TrimSpace(brandingRaw.String) != "" && brandingRaw.String != "null" {
		var m map[string]json.RawMessage
		if err := json.Unmarshal([]byte(brandingRaw.String), &m); err == nil {
			brandingValid = true
		}
	}
	if !brandingValid {
		updates["branding_settings"] = defaultBrandingSettings
		resetFields = append(resetFields, "branding_settings")
	}

	if len(updates) > 0 {
		setParts := make([]string, 0, len(updates))
		args := make([]any, 0, len(updates))
		idx := 1
		for col, val := range updates {
			setParts = append(setParts, fmt.Sprintf("%s = $%d::jsonb", col, idx))
			args = append(args, val)
			idx++
		}
		query := fmt.Sprintf("UPDATE exodus_settings SET %s WHERE id = 1", strings.Join(setParts, ", "))
		if _, err := tx.ExecContext(ctx, query, args...); err != nil {
			return fmt.Errorf("reset invalid exodus_settings fields: %w", err)
		}
		cfg.Logger.Info("ExodusSettings reset invalid fields", "reset", strings.Join(resetFields, ", "))
	} else {
		cfg.Logger.Debug("ExodusSettings already valid")
	}

	return nil
}

func ensureDefaultConfigProfile(ctx context.Context, tx *sql.Tx, cfg *config.BackendConfig) error {
	var count int
	if err := tx.QueryRowContext(ctx, `SELECT COUNT(*) FROM config_profiles`).Scan(&count); err != nil {
		return fmt.Errorf("count config_profiles: %w", err)
	}
	if count > 0 {
		return nil
	}

	query := `
		INSERT INTO config_profiles (
			uuid, view_position, name, config, created_at, updated_at
		) VALUES ($1, $2, $3, $4::jsonb, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	`
	if _, err := tx.ExecContext(ctx, query, "00000000-0000-0000-0000-000000000000", 1, "Default-Profile", defaultSingboxConfig); err != nil {
		return fmt.Errorf("insert default config_profile: %w", err)
	}

	cfg.Logger.Info("Default config profile seeded")
	return nil
}

func resyncConfigProfileInbounds(ctx context.Context, tx *sql.Tx) (int, error) {
	rows, err := tx.QueryContext(ctx, `SELECT uuid, config FROM config_profiles`)
	if err != nil {
		return 0, fmt.Errorf("list config profiles: %w", err)
	}
	defer rows.Close()

	type profile struct {
		uuid   string
		config json.RawMessage
	}
	var profiles []profile
	for rows.Next() {
		var p profile
		if err := rows.Scan(&p.uuid, &p.config); err != nil {
			return 0, err
		}
		profiles = append(profiles, p)
	}
	if err := rows.Err(); err != nil {
		return 0, err
	}

	for _, p := range profiles {
		if _, err := SyncConfigProfileInboundsTx(ctx, tx, p.uuid, p.config); err != nil {
			return 0, fmt.Errorf("sync inbounds for profile %s: %w", p.uuid, err)
		}
	}
	return len(profiles), nil
}

func ensureDefaultInternalSquad(ctx context.Context, tx *sql.Tx, cfg *config.BackendConfig) error {
	var count int
	if err := tx.QueryRowContext(ctx, `SELECT COUNT(*) FROM internal_squads`).Scan(&count); err != nil {
		return fmt.Errorf("count internal_squads: %w", err)
	}
	if count > 0 {
		return nil
	}

	var profileUUID string
	if err := tx.QueryRowContext(ctx, `SELECT uuid FROM config_profiles ORDER BY view_position ASC LIMIT 1`).Scan(&profileUUID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}
		return fmt.Errorf("read config profile for internal squad: %w", err)
	}

	rows, err := tx.QueryContext(ctx, `SELECT uuid FROM config_profile_inbounds WHERE profile_uuid = $1`, profileUUID)
	if err != nil {
		return fmt.Errorf("read config profile inbounds: %w", err)
	}
	defer rows.Close()

	var inboundUUIDs []string
	for rows.Next() {
		var inbUUID string
		if err := rows.Scan(&inbUUID); err != nil {
			return err
		}
		inboundUUIDs = append(inboundUUIDs, inbUUID)
	}
	if err := rows.Err(); err != nil {
		return err
	}

	squadUUID := "00000000-0000-0000-0000-000000000000"
	if _, err := tx.ExecContext(ctx, `
		INSERT INTO internal_squads (uuid, view_position, name, created_at, updated_at)
		VALUES ($1, $2, $3, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	`, squadUUID, 1, "Default-Squad"); err != nil {
		return fmt.Errorf("insert default internal squad: %w", err)
	}

	for _, inbUUID := range inboundUUIDs {
		if _, err := tx.ExecContext(ctx, `
			INSERT INTO internal_squad_inbounds (internal_squad_uuid, inbound_uuid)
			VALUES ($1, $2) ON CONFLICT DO NOTHING
		`, squadUUID, inbUUID); err != nil {
			return fmt.Errorf("link inbound %s to default squad: %w", inbUUID, err)
		}
	}

	cfg.Logger.Info("Default internal squad seeded", "inbounds", len(inboundUUIDs))
	return nil
}

func ensureDefaultSubscriptionSettings(ctx context.Context, tx *sql.Tx, cfg *config.BackendConfig) error {
	var (
		subUUID          string
		hwidRaw          sql.NullString
		customRemarksRaw sql.NullString
		responseRulesRaw sql.NullString
	)

	err := tx.QueryRowContext(ctx, `SELECT uuid, hwid_settings, custom_remarks, response_rules FROM subscription_settings LIMIT 1`).Scan(&subUUID, &hwidRaw, &customRemarksRaw, &responseRulesRaw)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("read subscription_settings: %w", err)
	}

	if err == nil {
		var (
			updates     []string
			args        []any
			resetFields []string
			argIdx      = 1
		)

		needHwidUpdate := !hwidRaw.Valid || strings.TrimSpace(hwidRaw.String) == "" || hwidRaw.String == "null"
		if needHwidUpdate {
			updates = append(updates, fmt.Sprintf("hwid_settings = $%d::jsonb", argIdx))
			args = append(args, defaultHWIDSettings)
			argIdx++
			resetFields = append(resetFields, "hwid_settings")
		}

		needRemarksUpdate := false
		var currentRemarks map[string]any
		if customRemarksRaw.Valid && strings.TrimSpace(customRemarksRaw.String) != "" && customRemarksRaw.String != "null" {
			if err := json.Unmarshal([]byte(customRemarksRaw.String), &currentRemarks); err != nil {
				needRemarksUpdate = true
			} else {
				if _, ok1 := currentRemarks["HWIDMaxDevicesExceeded"]; !ok1 {
					needRemarksUpdate = true
				}
				if _, ok2 := currentRemarks["HWIDNotSupported"]; !ok2 {
					needRemarksUpdate = true
				}
			}
		} else {
			needRemarksUpdate = true
		}

		if needRemarksUpdate {
			var defaultRemarks map[string]any
			_ = json.Unmarshal([]byte(defaultCustomRemarks), &defaultRemarks)
			if currentRemarks == nil {
				currentRemarks = defaultRemarks
			} else {
				for k, v := range defaultRemarks {
					if _, exists := currentRemarks[k]; !exists {
						currentRemarks[k] = v
					}
				}
			}
			mergedBytes, _ := json.Marshal(currentRemarks)
			updates = append(updates, fmt.Sprintf("custom_remarks = $%d::jsonb", argIdx))
			args = append(args, string(mergedBytes))
			argIdx++
			resetFields = append(resetFields, "custom_remarks")
		}

		var currentRules string
		if responseRulesRaw.Valid && strings.TrimSpace(responseRulesRaw.String) != "" && responseRulesRaw.String != "null" {
			currentRules = responseRulesRaw.String
		}

		if currentRules == "" {
			updates = append(updates, fmt.Sprintf("response_rules = $%d::jsonb", argIdx))
			args = append(args, defaultResponseRules)
			argIdx++
			resetFields = append(resetFields, "response_rules (seeded missing)")
		} else {
			existingHash, _ := canonicalHash(currentRules)
			defaultHash, _ := canonicalHash(defaultResponseRules)
			if existingHash == PrevResponseRulesHash && defaultHash != PrevResponseRulesHash {
				updates = append(updates, fmt.Sprintf("response_rules = $%d::jsonb", argIdx))
				args = append(args, defaultResponseRules)
				argIdx++
				resetFields = append(resetFields, "response_rules (upgraded to new default)")
			}
		}

		if len(updates) > 0 {
			args = append(args, subUUID)
			query := fmt.Sprintf("UPDATE subscription_settings SET %s WHERE uuid = $%d", strings.Join(updates, ", "), argIdx)
			if _, err := tx.ExecContext(ctx, query, args...); err != nil {
				return fmt.Errorf("update subscription_settings backfill: %w", err)
			}
			cfg.Logger.Info("SubscriptionSettings field backfill applied", "updated", strings.Join(resetFields, ", "))
		}
		return nil
	}

	query := `
		INSERT INTO subscription_settings (
			uuid, profile_title, support_link, profile_update_interval,
			address, port, api_schema, api_path,
			is_profile_webpage_url_enabled, serve_json_at_base_subscription,
			happ_announce, happ_routing, is_show_custom_remarks,
			custom_remarks, custom_response_headers, randomize_hosts,
			response_rules, hwid_settings
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18)
	`
	_, err = tx.ExecContext(ctx, query,
		"00000000-0000-0000-0000-000000000000",
		"exodus",
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

	cfg.Logger.Info("Default subscription settings seeded")
	return nil
}

func ensureDefaultTemplates(ctx context.Context, tx *sql.Tx, cfg *config.BackendConfig) error {
	defaults := DefaultSubscriptionTemplates()
	validTypes := make([]string, 0, len(defaults))
	for _, tmpl := range defaults {
		validTypes = append(validTypes, tmpl.TemplateType)
	}

	var deletedCount int64
	if len(validTypes) > 0 {
		placeholders := make([]string, len(validTypes))
		args := make([]any, len(validTypes))
		for i, t := range validTypes {
			placeholders[i] = fmt.Sprintf("$%d", i+1)
			args[i] = t
		}
		query := fmt.Sprintf("DELETE FROM subscription_templates WHERE template_type NOT IN (%s)", strings.Join(placeholders, ", "))
		res, err := tx.ExecContext(ctx, query, args...)
		if err != nil {
			return fmt.Errorf("delete obsolete subscription templates: %w", err)
		}
		deletedCount, _ = res.RowsAffected()
	}

	var seededCount int
	for _, tmpl := range defaults {
		var count int
		if err := tx.QueryRowContext(ctx, `SELECT COUNT(*) FROM subscription_templates WHERE template_type = $1`, tmpl.TemplateType).Scan(&count); err != nil {
			return fmt.Errorf("count template %s: %w", tmpl.TemplateType, err)
		}
		if count > 0 {
			continue
		}

		query := `
			INSERT INTO subscription_templates (
				uuid, view_position, name, template_type, template_yaml, template_json
			) VALUES ($1, $2, $3, $4, $5, $6)
		`
		if _, err := tx.ExecContext(ctx, query,
			uuid.NewString(),
			tmpl.ViewPosition,
			tmpl.Name,
			tmpl.TemplateType,
			tmpl.TemplateYAML,
			tmpl.TemplateJSON,
		); err != nil {
			return fmt.Errorf("insert template %s: %w", tmpl.TemplateType, err)
		}
		seededCount++
	}

	if deletedCount > 0 || seededCount > 0 {
		cfg.Logger.Info("Subscription templates updated", "deletedObsolete", deletedCount, "seededMissing", seededCount)
	} else {
		cfg.Logger.Debug("Subscription templates up to date")
	}

	return nil
}

func ensureDefaultSubscriptionPageConfig(ctx context.Context, tx *sql.Tx) error {
	var count int
	if err := tx.QueryRowContext(ctx, `SELECT COUNT(*) FROM subscription_page_config`).Scan(&count); err != nil {
		return fmt.Errorf("count subscription_page_config: %w", err)
	}
	if count > 0 {
		return nil
	}

	query := `
		INSERT INTO subscription_page_config (
			uuid, view_position, name, config
		) VALUES ($1, $2, $3, $4)
	`
	if _, err := tx.ExecContext(ctx, query, defaultSubpageConfigUUID, 1, "Default", DefaultSubscriptionPageConfig); err != nil {
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

func ensureKeygen(ctx context.Context, tx *sql.Tx) error {
	var count int
	if err := tx.QueryRowContext(ctx, `SELECT COUNT(*) FROM keygen`).Scan(&count); err != nil {
		return fmt.Errorf("count keygen rows: %w", err)
	}

	if count > 1 {
		if _, err := tx.ExecContext(ctx, `DELETE FROM keygen`); err != nil {
			return fmt.Errorf("delete old keygen rows: %w", err)
		}
		count = 0
	}

	if count == 0 {
		pubKey, privKey, err := security.GenerateJWTKeypair()
		if err != nil {
			return fmt.Errorf("generate jwt keypair: %w", err)
		}
		masterCerts, err := security.GenerateMasterCerts()
		if err != nil {
			return fmt.Errorf("generate master certs: %w", err)
		}

		query := `
			INSERT INTO keygen (uuid, priv_key, pub_key, ca_cert, ca_key, client_cert, client_key)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
		`
		if _, err := tx.ExecContext(
			ctx,
			query,
			uuid.NewString(),
			privKey,
			pubKey,
			masterCerts.CACertPEM,
			masterCerts.CAKeyPEM,
			masterCerts.ClientCertPEM,
			masterCerts.ClientKeyPEM,
		); err != nil {
			return fmt.Errorf("insert keygen row: %w", err)
		}
		return nil
	}

	var (
		id         string
		pubKey     sql.NullString
		privKey    sql.NullString
		caCert     sql.NullString
		caKey      sql.NullString
		clientCert sql.NullString
		clientKey  sql.NullString
	)
	if err := tx.QueryRowContext(
		ctx,
		`SELECT uuid, pub_key, priv_key, ca_cert, ca_key, client_cert, client_key FROM keygen ORDER BY created_at ASC LIMIT 1`,
	).Scan(&id, &pubKey, &privKey, &caCert, &caKey, &clientCert, &clientKey); err != nil {
		return fmt.Errorf("read keygen row: %w", err)
	}

	needJWT := !pubKey.Valid || pubKey.String == "" || !privKey.Valid || privKey.String == ""
	needMTLS := !caCert.Valid || caCert.String == "" || !caKey.Valid || caKey.String == "" || !clientCert.Valid || clientCert.String == "" || !clientKey.Valid || clientKey.String == ""
	if !needJWT && !needMTLS {
		return nil
	}

	updateParts := make([]string, 0, 6)
	args := make([]any, 0, 7)
	idx := 1
	if needJWT {
		newPubKey, newPrivKey, err := security.GenerateJWTKeypair()
		if err != nil {
			return fmt.Errorf("regenerate jwt keypair: %w", err)
		}
		updateParts = append(updateParts, fmt.Sprintf("pub_key = $%d", idx), fmt.Sprintf("priv_key = $%d", idx+1))
		args = append(args, newPubKey, newPrivKey)
		idx += 2
	}
	if needMTLS {
		masterCerts, err := security.GenerateMasterCerts()
		if err != nil {
			return fmt.Errorf("regenerate mTLS certificates: %w", err)
		}
		updateParts = append(updateParts, fmt.Sprintf("ca_cert = $%d", idx), fmt.Sprintf("ca_key = $%d", idx+1), fmt.Sprintf("client_cert = $%d", idx+2), fmt.Sprintf("client_key = $%d", idx+3))
		args = append(args, masterCerts.CACertPEM, masterCerts.CAKeyPEM, masterCerts.ClientCertPEM, masterCerts.ClientKeyPEM)
		idx += 4
	}

	query := fmt.Sprintf("UPDATE keygen SET %s, updated_at = CURRENT_TIMESTAMP WHERE uuid = $%d", joinWithComma(updateParts), idx)
	args = append(args, id)
	if _, err := tx.ExecContext(ctx, query, args...); err != nil {
		return fmt.Errorf("update keygen row: %w", err)
	}

	return nil
}

func joinWithComma(parts []string) string {
	if len(parts) == 0 {
		return ""
	}
	out := parts[0]
	for i := 1; i < len(parts); i++ {
		out += ", " + parts[i]
	}
	return out
}

func ensureSingleAdmin(ctx context.Context, tx *sql.Tx, cfg *config.BackendConfig) error {
	var count int
	if err := tx.QueryRowContext(ctx, `SELECT COUNT(*) FROM admin`).Scan(&count); err != nil {
		return fmt.Errorf("count admin rows: %w", err)
	}

	if count <= 1 {
		return nil
	}

	res, err := tx.ExecContext(ctx, `
		DELETE FROM admin
		WHERE uuid NOT IN (
			SELECT uuid FROM admin ORDER BY created_at ASC LIMIT 1
		)
	`)
	if err != nil {
		return fmt.Errorf("deduplicate admin records: %w", err)
	}

	deleted, _ := res.RowsAffected()
	cfg.Logger.Warn("Multiple admin records found, retained oldest admin and removed others", "deleted", deleted)
	return nil
}
