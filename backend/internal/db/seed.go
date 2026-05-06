package db

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"fmt"
	"strings"

	"exodus/internal/config"
	"exodus/internal/dbutil"
	"exodus/internal/security"

	"github.com/google/uuid"
)

const defaultResponseRules = `{"rules":[{"name":"Browser Subscription","enabled":true,"operator":"AND","conditions":[{"value":"text/html","operator":"CONTAINS","headerName":"accept","caseSensitive":true}],"description":"System critical: do not delete or disable this rule.","responseType":"BROWSER"},{"name":"Mihomo Clients","enabled":true,"operator":"AND","conditions":[{"value":"^(?:FlClash|FlClashX|Flowvy|[Cc]lash-[Vv]erge|[Kk]oala-[Cc]lash|[Cc]lash-?[Mm]eta|[Mm]urge|[Cc]lashX [Mm]eta|[Mm]ihomo|[Cc]lash-nyanpasu|clash.meta|prizrak-box)","operator":"REGEX","headerName":"user-agent","caseSensitive":false}],"description":"Response with generated YAML config (Mihomo Template)","responseType":"MIHOMO"},{"name":"Stash (iOS, macOS)","enabled":true,"operator":"AND","conditions":[{"value":"^stash","operator":"REGEX","headerName":"user-agent","caseSensitive":false}],"description":"Response with generated YAML config (Stash Template)","responseType":"STASH"},{"name":"Sing-box clients","enabled":true,"operator":"AND","conditions":[{"value":"^sfa|sfi|sfm|sft|karing|singbox|rabbithole","operator":"REGEX","headerName":"user-agent","caseSensitive":false}],"description":"Resonse with generated JSON config (Singbox template)","responseType":"SINGBOX"},{"name":"Clash Core Clients","enabled":true,"operator":"AND","conditions":[{"value":"^clash","operator":"REGEX","headerName":"user-agent","caseSensitive":false}],"description":"Response with generated YAML config (Clash Template)","responseType":"CLASH"},{"name":"Fallback Base64","enabled":true,"operator":"AND","conditions":[],"description":"System critical: do not delete or disable this rule.","responseType":"XRAY_BASE64"}],"version":"1"}`
const defaultHWIDSettings = `{"enabled":false,"maxDevicesAnnounce":null,"fallbackDeviceLimit":999}`
const defaultCustomRemarks = `{"emptyHosts":["→ exodus","→ No hosts found","→ Check Hosts tab","→ Check Internal Squads tab"],"expiredUsers":["⌛ Subscription expired","Contact support"],"limitedUsers":["🚧 Subscription limited","Contact support"],"disabledUsers":["🚫 Subscription disabled","Contact support"],"HWIDNotSupported":["App not supported"],"HWIDMaxDevicesExceeded":["Limit of devices reached"]}`
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

	if err := ensureV2rsSettings(ctx, tx); err != nil {
		_ = tx.Rollback()
		return err
	}
	if err := ensureModulesSettings(ctx, tx); err != nil {
		_ = tx.Rollback()
		return err
	}
	if err := ensureDefaultSubscriptionSettings(ctx, tx); err != nil {
		_ = tx.Rollback()
		return err
	}
	if err := ensureDefaultTemplates(ctx, tx); err != nil {
		_ = tx.Rollback()
		return err
	}
	if err := ensureDefaultSubscriptionPageConfig(ctx, tx); err != nil {
		_ = tx.Rollback()
		return err
	}
	if err := ensureKeygen(ctx, tx); err != nil {
		_ = tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit defaults transaction: %w", err)
	}

	cfg.Logger.Info("Default configuration data ensured")
	return nil
}

func ensureV2rsSettings(ctx context.Context, tx *sql.Tx) error {
	var exists int
	if err := tx.QueryRowContext(ctx, `SELECT COUNT(*) FROM exodus_settings WHERE id = 1`).Scan(&exists); err != nil {
		return fmt.Errorf("check exodus_settings row: %w", err)
	}
	if exists > 0 {
		return nil
	}

	passkeySettings := `{"rpId":null,"origin":null,"enabled":false}`
	oauth2Settings := `{"github":{"enabled":false,"clientId":null,"clientSecret":null,"allowedEmails":[]},"yandex":{"enabled":false,"clientId":null,"clientSecret":null,"allowedEmails":[]},"generic":{"enabled":false,"clientId":null,"tokenUrl":null,"withPkce":false,"clientSecret":null,"allowedEmails":[],"frontendDomain":null,"authorizationUrl":null},"keycloak":{"realm":null,"enabled":false,"clientId":null,"clientSecret":null,"allowedEmails":[],"frontendDomain":null,"keycloakDomain":null},"pocketid":{"enabled":false,"clientId":null,"plainDomain":null,"clientSecret":null,"allowedEmails":[]}}`
	tgAuthSettings := `{"enabled":false,"adminIds":[],"botToken":null}`
	passwordSettings := `{"enabled":true}`
	brandingSettings := `{"title":"EXODUS","logoUrl":null}`
	modulesSettings := `{"haproxy":{"enabled":false}}`

	query := dbutil.Rebind(`
		INSERT INTO exodus_settings (
			id, passkey_settings, oauth2_settings, tg_auth_settings, password_settings, branding_settings, modules_settings
		) VALUES (1, ?, ?, ?, ?, ?, ?)
	`)
	if _, err := tx.ExecContext(ctx, query, passkeySettings, oauth2Settings, tgAuthSettings, passwordSettings, brandingSettings, modulesSettings); err != nil {
		return fmt.Errorf("insert default exodus_settings row: %w", err)
	}

	return nil
}

func ensureModulesSettings(ctx context.Context, tx *sql.Tx) error {
	var exists int
	if err := tx.QueryRowContext(ctx, `SELECT COUNT(*) FROM modules_settings WHERE id = 1`).Scan(&exists); err != nil {
		return fmt.Errorf("check modules_settings row: %w", err)
	}
	if exists > 0 {
		return nil
	}

	query := dbutil.Rebind(`
		INSERT INTO modules_settings (
			id, haproxy_enabled
		) VALUES (1, ?)
	`)
	if _, err := tx.ExecContext(ctx, query, false); err != nil {
		return fmt.Errorf("insert default modules_settings row: %w", err)
	}
	return nil
}

func ensureDefaultSubscriptionSettings(ctx context.Context, tx *sql.Tx) error {
	var count int
	if err := tx.QueryRowContext(ctx, `SELECT COUNT(*) FROM subscription_settings`).Scan(&count); err != nil {
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
	_, err := tx.ExecContext(ctx, query,
		uuid.NewString(),
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

	return nil
}

func ensureDefaultTemplates(ctx context.Context, tx *sql.Tx) error {
	for _, tmpl := range DefaultSubscriptionTemplates() {
		var count int
		if err := tx.QueryRowContext(ctx, dbutil.Rebind(`SELECT COUNT(*) FROM subscription_templates WHERE template_type = ?`), tmpl.TemplateType).Scan(&count); err != nil {
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

	query := dbutil.Rebind(`
		INSERT INTO subscription_page_config (
			uuid, view_position, name, config
		) VALUES (?, ?, ?, ?)
	`)
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

		query := dbutil.Rebind(`
			INSERT INTO keygen (uuid, priv_key, pub_key, ca_cert, ca_key, client_cert, client_key, grpc_auth_token)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?)
		`)
		grpcAuthToken, err := generateGRPCAuthToken()
		if err != nil {
			return fmt.Errorf("generate grpc auth token: %w", err)
		}
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
			grpcAuthToken,
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
		grpcToken  sql.NullString
	)
	if err := tx.QueryRowContext(
		ctx,
		`SELECT uuid, pub_key, priv_key, ca_cert, ca_key, client_cert, client_key, grpc_auth_token FROM keygen ORDER BY created_at ASC LIMIT 1`,
	).Scan(&id, &pubKey, &privKey, &caCert, &caKey, &clientCert, &clientKey, &grpcToken); err != nil {
		return fmt.Errorf("read keygen row: %w", err)
	}

	needJWT := !pubKey.Valid || pubKey.String == "" || !privKey.Valid || privKey.String == ""
	needMTLS := !caCert.Valid || caCert.String == "" || !caKey.Valid || caKey.String == "" || !clientCert.Valid || clientCert.String == "" || !clientKey.Valid || clientKey.String == ""
	needGRPCToken := !grpcToken.Valid || strings.TrimSpace(grpcToken.String) == ""
	if !needJWT && !needMTLS && !needGRPCToken {
		return nil
	}

	updateParts := make([]string, 0, 6)
	args := make([]any, 0, 7)
	if needJWT {
		newPubKey, newPrivKey, err := security.GenerateJWTKeypair()
		if err != nil {
			return fmt.Errorf("regenerate jwt keypair: %w", err)
		}
		updateParts = append(updateParts, "pub_key = ?", "priv_key = ?")
		args = append(args, newPubKey, newPrivKey)
	}
	if needMTLS {
		masterCerts, err := security.GenerateMasterCerts()
		if err != nil {
			return fmt.Errorf("regenerate mTLS certificates: %w", err)
		}
		updateParts = append(updateParts, "ca_cert = ?", "ca_key = ?", "client_cert = ?", "client_key = ?")
		args = append(args, masterCerts.CACertPEM, masterCerts.CAKeyPEM, masterCerts.ClientCertPEM, masterCerts.ClientKeyPEM)
	}
	if needGRPCToken {
		token, err := generateGRPCAuthToken()
		if err != nil {
			return fmt.Errorf("regenerate grpc auth token: %w", err)
		}
		updateParts = append(updateParts, "grpc_auth_token = ?")
		args = append(args, token)
	}

	query := dbutil.Rebind(fmt.Sprintf("UPDATE keygen SET %s, updated_at = CURRENT_TIMESTAMP WHERE uuid = ?", joinWithComma(updateParts)))
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

func generateGRPCAuthToken() (string, error) {
	raw := make([]byte, 32)
	if _, err := rand.Read(raw); err != nil {
		return "", err
	}
	return hex.EncodeToString(raw), nil
}
