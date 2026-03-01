package db

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"v2ray-stat/backend/config"
	"v2ray-stat/backend/db/manager"
	_ "github.com/mattn/go-sqlite3"
)

// OpenAndInitDB opens and initializes a SQLite database.
func OpenAndInitDB(dbPath string, dbType string, cfg *config.BackendConfig) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		cfg.Logger.Error("Failed to open database", "dbType", dbType, "path", dbPath, "error", err)
		return nil, fmt.Errorf("failed to open %s database: %v", dbType, err)
	}
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)

	var tableCount int
	err = db.QueryRow("SELECT count(*) FROM sqlite_master WHERE type='table' AND name='v2raystat_settings'").Scan(&tableCount)
	if err != nil {
		cfg.Logger.Error("Failed to check table existence", "dbType", dbType, "error", err)
		db.Close()
		return nil, fmt.Errorf("failed to check table existence for %s database: %v", dbType, err)
	}
	if tableCount > 0 {
		cfg.Logger.Debug("Tables already exist", "dbType", dbType)
		if err := applySchemaMigrations(db, cfg); err != nil {
			cfg.Logger.Error("Failed to apply schema migrations", "dbType", dbType, "error", err)
			db.Close()
			return nil, fmt.Errorf("failed to apply schema migrations for %s database: %v", dbType, err)
		}
		if err := ensurePanelAuthDefaults(db, cfg); err != nil {
			cfg.Logger.Error("Failed to ensure panel auth defaults", "dbType", dbType, "error", err)
			db.Close()
			return nil, fmt.Errorf("failed to ensure panel auth defaults for %s database: %v", dbType, err)
		}
		return db, nil
	}

	sqlStmt := `
		PRAGMA foreign_keys = ON;
		PRAGMA cache_size = 2000;
		PRAGMA journal_mode = WAL;
		PRAGMA synchronous = NORMAL;
		PRAGMA temp_store = MEMORY;
		PRAGMA busy_timeout = 5000;

		-- ==========================================
		-- GENERAL SETTINGS & ADMIN
		-- ==========================================

		CREATE TABLE IF NOT EXISTS remnawave_settings (
			id INTEGER PRIMARY KEY DEFAULT 1,
			passkey_settings TEXT,
			oauth2_settings TEXT,
			tg_auth_settings TEXT,
			password_settings TEXT,
			branding_settings TEXT
		);

		CREATE TABLE IF NOT EXISTS admin (
			uuid TEXT PRIMARY KEY,
			username TEXT UNIQUE NOT NULL,
			password_hash TEXT NOT NULL,
			role TEXT NOT NULL,
			session_ttl_minutes INTEGER NOT NULL DEFAULT 60,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);

		CREATE TABLE IF NOT EXISTS admin_sessions (
			session_token TEXT PRIMARY KEY,
			admin_uuid TEXT NOT NULL,
			expires_at INTEGER NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (admin_uuid) REFERENCES admin(uuid) ON DELETE CASCADE
		);
		CREATE INDEX IF NOT EXISTS idx_admin_sessions_admin_uuid ON admin_sessions(admin_uuid);
		CREATE INDEX IF NOT EXISTS idx_admin_sessions_expires_at ON admin_sessions(expires_at);

		CREATE TABLE IF NOT EXISTS passkeys (
			id TEXT PRIMARY KEY,
			admin_uuid TEXT NOT NULL,
			public_key BLOB NOT NULL,
			counter INTEGER NOT NULL,
			device_type TEXT NOT NULL,
			backed_up BOOLEAN NOT NULL,
			transports TEXT,
			passkey_provider TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (admin_uuid) REFERENCES admin(uuid) ON DELETE CASCADE
		);
		CREATE INDEX IF NOT EXISTS idx_passkeys_id ON passkeys(id);
		CREATE INDEX IF NOT EXISTS idx_passkeys_admin_uuid ON passkeys(admin_uuid);

		CREATE TABLE IF NOT EXISTS api_tokens (
			uuid TEXT PRIMARY KEY,
			token TEXT UNIQUE NOT NULL,
			token_name TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);

		CREATE TABLE IF NOT EXISTS keygen (
			uuid TEXT PRIMARY KEY,
			priv_key TEXT NOT NULL,
			pub_key TEXT NOT NULL,
			ca_cert TEXT,
			ca_key TEXT,
			client_cert TEXT,
			client_key TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);

		-- ==========================================
		-- INFRASTRUCTURE & NODES
		-- ==========================================

		CREATE TABLE IF NOT EXISTS infra_providers (
			uuid TEXT PRIMARY KEY,
			name TEXT UNIQUE NOT NULL,
			favicon_link TEXT,
			login_url TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);

		CREATE TABLE IF NOT EXISTS subscription_page_config (
			uuid TEXT PRIMARY KEY,
			view_position INTEGER DEFAULT 0,
			name TEXT UNIQUE NOT NULL,
			config TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);

		CREATE TABLE IF NOT EXISTS config_profiles (
			uuid TEXT PRIMARY KEY,
			view_position INTEGER DEFAULT 0,
			name TEXT UNIQUE NOT NULL,
			config TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);

		CREATE TABLE IF NOT EXISTS nodes (
			uuid TEXT PRIMARY KEY,
			id INTEGER UNIQUE,
			name TEXT NOT NULL,
			address TEXT NOT NULL,
			port INTEGER,
			api_schema TEXT DEFAULT 'grpc',
    		api_path TEXT DEFAULT '',
			active_config_profile_uuid TEXT,
			is_connected BOOLEAN DEFAULT 0,
			is_connecting BOOLEAN DEFAULT 0,
			is_disabled BOOLEAN DEFAULT 0,
			last_status_change DATETIME,
			last_status_message TEXT,
			xray_version TEXT,
			node_version TEXT,
			xray_uptime TEXT DEFAULT '0',
			users_online INTEGER DEFAULT 0,
			consumption_multiplier INTEGER DEFAULT 1000000000,
			is_traffic_tracking_active BOOLEAN DEFAULT 0,
			traffic_reset_day INTEGER DEFAULT 1,
			traffic_limit_bytes INTEGER DEFAULT 0,
			traffic_used_bytes INTEGER DEFAULT 0,
			notify_percent INTEGER DEFAULT 0,
			provider_uuid TEXT,
			view_position INTEGER DEFAULT 0,
			country_code TEXT DEFAULT 'XX',
			tags TEXT DEFAULT '[]',
			cpu_count INTEGER,
			cpu_model TEXT,
			total_ram TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (active_config_profile_uuid) REFERENCES config_profiles(uuid) ON DELETE SET NULL,
			FOREIGN KEY (provider_uuid) REFERENCES infra_providers(uuid) ON DELETE SET NULL
		);
		CREATE INDEX IF NOT EXISTS idx_nodes_name ON nodes(name);
		CREATE INDEX IF NOT EXISTS idx_nodes_address ON nodes(address);

		-- ==========================================
		-- TASKS & TASK NODES
		-- ==========================================

		CREATE TABLE IF NOT EXISTS tasks (
			id TEXT PRIMARY KEY,
			operation TEXT NOT NULL,
			payload TEXT,
			status TEXT NOT NULL DEFAULT 'pending',
			created_at INTEGER NOT NULL,
			completed_at INTEGER,
			timeout_at INTEGER NOT NULL
		);
		CREATE INDEX IF NOT EXISTS idx_tasks_status ON tasks(status);
		CREATE INDEX IF NOT EXISTS idx_tasks_created ON tasks(created_at);

		CREATE TABLE IF NOT EXISTS task_nodes (
			id TEXT PRIMARY KEY,
			task_id TEXT NOT NULL,
			node_name TEXT NOT NULL,
			status TEXT NOT NULL DEFAULT 'pending',
			error_message TEXT,
			sent_at INTEGER,
			completed_at INTEGER,
			FOREIGN KEY (task_id) REFERENCES tasks(id) ON DELETE CASCADE
		);
		CREATE INDEX IF NOT EXISTS idx_task_nodes_task_id ON task_nodes(task_id);
		CREATE INDEX IF NOT EXISTS idx_task_nodes_status ON task_nodes(status);

		CREATE TABLE IF NOT EXISTS infra_billing_nodes (
			uuid TEXT PRIMARY KEY,
			node_uuid TEXT NOT NULL,
			provider_uuid TEXT NOT NULL,
			next_billing_at DATETIME NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			UNIQUE(node_uuid, provider_uuid),
			FOREIGN KEY (node_uuid) REFERENCES nodes(uuid) ON DELETE CASCADE,
			FOREIGN KEY (provider_uuid) REFERENCES infra_providers(uuid) ON DELETE CASCADE
		);
		CREATE INDEX IF NOT EXISTS idx_infra_billing_next_at ON infra_billing_nodes(next_billing_at);

		CREATE TABLE IF NOT EXISTS infra_billing_history (
			uuid TEXT PRIMARY KEY,
			provider_uuid TEXT NOT NULL,
			amount REAL NOT NULL,
			billed_at DATETIME NOT NULL,
			FOREIGN KEY (provider_uuid) REFERENCES infra_providers(uuid) ON DELETE CASCADE
		);

		CREATE TABLE IF NOT EXISTS nodes_usage_history (
			node_uuid TEXT NOT NULL,
			download_bytes INTEGER NOT NULL,
			upload_bytes INTEGER NOT NULL,
			total_bytes INTEGER NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			PRIMARY KEY (node_uuid, created_at),
			FOREIGN KEY (node_uuid) REFERENCES nodes(uuid) ON DELETE CASCADE
		);
		CREATE INDEX IF NOT EXISTS idx_nodes_usage_hist ON nodes_usage_history(node_uuid, created_at DESC);

		CREATE TABLE IF NOT EXISTS nodes_traffic_usage_history (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			node_uuid TEXT NOT NULL,
			traffic_bytes INTEGER NOT NULL,
			reset_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (node_uuid) REFERENCES nodes(uuid) ON DELETE CASCADE
		);

		-- ==========================================
		-- HOSTS & CONFIGURATIONS
		-- ==========================================

		CREATE TABLE IF NOT EXISTS subscription_templates (
			uuid TEXT PRIMARY KEY,
			view_position INTEGER DEFAULT 0,
			name TEXT DEFAULT 'Default',
			template_type TEXT NOT NULL,
			template_yaml TEXT,
			template_json TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			UNIQUE(template_type, name)
		);

		CREATE TABLE IF NOT EXISTS config_profile_inbounds (
			uuid TEXT PRIMARY KEY,
			profile_uuid TEXT NOT NULL,
			tag TEXT UNIQUE NOT NULL,
			type TEXT NOT NULL,
			network TEXT,
			security TEXT,
			port INTEGER,
			raw_inbound TEXT,
			FOREIGN KEY (profile_uuid) REFERENCES config_profiles(uuid) ON DELETE CASCADE
		);
		CREATE INDEX IF NOT EXISTS idx_cp_inbounds_profile ON config_profile_inbounds(profile_uuid, uuid);

		CREATE TABLE IF NOT EXISTS hosts (
			uuid TEXT PRIMARY KEY,
			view_position INTEGER DEFAULT 0,
			remark TEXT NOT NULL,
			address TEXT NOT NULL,
			port INTEGER NOT NULL,
			path TEXT,
			sni TEXT,
			host TEXT,
			alpn TEXT,
				fingerprint TEXT,
				security_layer TEXT DEFAULT 'DEFAULT',
			xhttp_extra_params TEXT,
			mux_params TEXT,
			sockopt_params TEXT,
			is_disabled BOOLEAN DEFAULT 0,
			server_description TEXT,
			vless_route_id INTEGER,
			allow_insecure BOOLEAN DEFAULT 0,
			shuffle_host BOOLEAN DEFAULT 0,
			mihomo_x25519 BOOLEAN DEFAULT 0,
			xray_json_template_uuid TEXT,
			keep_sni_blank BOOLEAN DEFAULT 0,
			tag TEXT,
			is_hidden BOOLEAN DEFAULT 0,
			override_sni_from_address BOOLEAN DEFAULT 0,
			config_profile_uuid TEXT,
			config_profile_inbound_uuid TEXT,
			FOREIGN KEY (config_profile_inbound_uuid) REFERENCES config_profile_inbounds(uuid) ON DELETE SET NULL,
			FOREIGN KEY (config_profile_uuid) REFERENCES config_profiles(uuid) ON DELETE SET NULL,
			FOREIGN KEY (xray_json_template_uuid) REFERENCES subscription_templates(uuid) ON DELETE SET NULL
		);

		CREATE TABLE IF NOT EXISTS hosts_to_nodes (
			host_uuid TEXT NOT NULL,
			node_uuid TEXT NOT NULL,
			PRIMARY KEY (host_uuid, node_uuid),
			FOREIGN KEY (host_uuid) REFERENCES hosts(uuid) ON DELETE CASCADE,
			FOREIGN KEY (node_uuid) REFERENCES nodes(uuid) ON DELETE CASCADE
		);

		CREATE TABLE IF NOT EXISTS config_profile_inbounds_to_nodes (
			config_profile_inbound_uuid TEXT NOT NULL,
			node_uuid TEXT NOT NULL,
			PRIMARY KEY (config_profile_inbound_uuid, node_uuid),
			FOREIGN KEY (config_profile_inbound_uuid) REFERENCES config_profile_inbounds(uuid) ON DELETE CASCADE,
			FOREIGN KEY (node_uuid) REFERENCES nodes(uuid) ON DELETE CASCADE
		);

		CREATE TABLE IF NOT EXISTS config_profile_snippets (
			name TEXT PRIMARY KEY,
			snippet TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);

		-- ==========================================
		-- SQUADS (EXTERNAL & INTERNAL)
		-- ==========================================

		CREATE TABLE IF NOT EXISTS external_squads (
			uuid TEXT PRIMARY KEY,
			view_position INTEGER DEFAULT 0,
			name TEXT UNIQUE NOT NULL,
			subscription_settings TEXT,
			host_overrides TEXT,
			response_headers TEXT,
			hwid_settings TEXT,
			custom_remarks TEXT,
			subpage_config_uuid TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (subpage_config_uuid) REFERENCES subscription_page_config(uuid) ON DELETE SET NULL
		);

		CREATE TABLE IF NOT EXISTS external_squads_templates (
			external_squad_uuid TEXT NOT NULL,
			template_uuid TEXT NOT NULL,
			template_type TEXT NOT NULL,
			PRIMARY KEY (external_squad_uuid, template_type),
			FOREIGN KEY (external_squad_uuid) REFERENCES external_squads(uuid) ON DELETE CASCADE,
			FOREIGN KEY (template_uuid) REFERENCES subscription_templates(uuid) ON DELETE CASCADE
		);

		CREATE TABLE IF NOT EXISTS internal_squads (
			uuid TEXT PRIMARY KEY,
			view_position INTEGER DEFAULT 0,
			name TEXT UNIQUE NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);

		CREATE TABLE IF NOT EXISTS internal_squad_inbounds (
			internal_squad_uuid TEXT NOT NULL,
			inbound_uuid TEXT NOT NULL,
			PRIMARY KEY (internal_squad_uuid, inbound_uuid),
			FOREIGN KEY (internal_squad_uuid) REFERENCES internal_squads(uuid) ON DELETE CASCADE,
			FOREIGN KEY (inbound_uuid) REFERENCES config_profile_inbounds(uuid) ON DELETE CASCADE
		);

		CREATE TABLE IF NOT EXISTS internal_squad_host_exclusions (
			host_uuid TEXT NOT NULL,
			squad_uuid TEXT NOT NULL,
			PRIMARY KEY (host_uuid, squad_uuid),
			FOREIGN KEY (host_uuid) REFERENCES hosts(uuid) ON DELETE CASCADE,
			FOREIGN KEY (squad_uuid) REFERENCES internal_squads(uuid) ON DELETE CASCADE
		);

		-- ==========================================
		-- USERS & TRAFFIC
		-- ==========================================

		CREATE TABLE IF NOT EXISTS users (
			t_id INTEGER PRIMARY KEY AUTOINCREMENT,
			uuid TEXT UNIQUE NOT NULL,
			short_uuid TEXT UNIQUE NOT NULL,
				username TEXT UNIQUE NOT NULL,
			status TEXT DEFAULT 'ACTIVE',
			traffic_limit_bytes INTEGER DEFAULT 0,
			traffic_limit_strategy TEXT DEFAULT 'NO_RESET',
			expire_at DATETIME NOT NULL,
			sub_last_user_agent TEXT,
			sub_last_opened_at DATETIME,
			last_traffic_reset_at DATETIME,
			sub_revoked_at DATETIME,
			trojan_password TEXT NOT NULL,
			vless_uuid TEXT NOT NULL,
			ss_password TEXT NOT NULL,
			description TEXT,
			tag TEXT,
			telegram_id INTEGER,
			email TEXT,
			hwid_device_limit INTEGER,
			external_squad_uuid TEXT,
			last_triggered_threshold INTEGER DEFAULT 0,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (external_squad_uuid) REFERENCES external_squads(uuid) ON DELETE SET NULL
		);

		CREATE TABLE IF NOT EXISTS user_traffic (
			t_id INTEGER PRIMARY KEY,
			used_traffic_bytes INTEGER DEFAULT 0,
			lifetime_used_traffic_bytes INTEGER DEFAULT 0,
			online_at DATETIME,
			last_connected_node_uuid TEXT,
			first_connected_at DATETIME,
			FOREIGN KEY (t_id) REFERENCES users(t_id) ON DELETE CASCADE,
			FOREIGN KEY (last_connected_node_uuid) REFERENCES nodes(uuid) ON DELETE SET NULL
		);

		CREATE TABLE IF NOT EXISTS nodes_user_usage_history (
			node_id INTEGER NOT NULL,
			user_id INTEGER NOT NULL,
			total_bytes INTEGER NOT NULL,
			created_at DATETIME DEFAULT CURRENT_DATE,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			PRIMARY KEY (node_id, created_at, user_id),
			FOREIGN KEY (node_id) REFERENCES nodes(id) ON DELETE CASCADE,
			FOREIGN KEY (user_id) REFERENCES users(t_id) ON DELETE CASCADE
		);

		CREATE TABLE IF NOT EXISTS hwid_user_devices (
			hwid TEXT NOT NULL,
			user_uuid TEXT NOT NULL,
			platform TEXT,
			os_version TEXT,
			device_model TEXT,
			user_agent TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			PRIMARY KEY (hwid, user_uuid),
			FOREIGN KEY (user_uuid) REFERENCES users(uuid) ON DELETE CASCADE
		);

		CREATE TABLE IF NOT EXISTS internal_squad_members (
			internal_squad_uuid TEXT NOT NULL,
			user_id INTEGER NOT NULL,
			PRIMARY KEY (internal_squad_uuid, user_id),
			FOREIGN KEY (internal_squad_uuid) REFERENCES internal_squads(uuid) ON DELETE CASCADE,
			FOREIGN KEY (user_id) REFERENCES users(t_id) ON DELETE CASCADE
		);
		CREATE INDEX IF NOT EXISTS idx_ism_squad ON internal_squad_members(internal_squad_uuid);
		CREATE INDEX IF NOT EXISTS idx_ism_user ON internal_squad_members(user_id);

		CREATE TABLE IF NOT EXISTS user_subscription_request_history (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_uuid TEXT NOT NULL,
			request_ip TEXT,
			user_agent TEXT,
			request_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_uuid) REFERENCES users(uuid) ON DELETE CASCADE
		);
		CREATE INDEX IF NOT EXISTS idx_usr_hist_user ON user_subscription_request_history(user_uuid);
		CREATE INDEX IF NOT EXISTS idx_usr_hist_time ON user_subscription_request_history(request_at);

		-- ==========================================
		-- SUBSCRIPTION SETTINGS
		-- ==========================================

		CREATE TABLE IF NOT EXISTS subscription_settings (
			uuid TEXT PRIMARY KEY,
			profile_title TEXT NOT NULL,
			support_link TEXT NOT NULL,
			profile_update_interval INTEGER NOT NULL,
			address TEXT NOT NULL DEFAULT '',
			port INTEGER NOT NULL DEFAULT 9263,
			api_schema TEXT NOT NULL DEFAULT 'grpc',
			api_path TEXT NOT NULL DEFAULT '',
			is_profile_webpage_url_enabled BOOLEAN DEFAULT 1,
			serve_json_at_base_subscription BOOLEAN DEFAULT 0,
			happ_announce TEXT,
			happ_routing TEXT,
			is_show_custom_remarks BOOLEAN DEFAULT 1,
			custom_remarks TEXT NOT NULL,
			custom_response_headers TEXT,
			randomize_hosts BOOLEAN DEFAULT 0,
			response_rules TEXT,
			hwid_settings TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);

		-- ==========================================
		-- TRIGGERS FOR UPDATED_AT (Optional, highly recommended for SQLite)
		-- ==========================================
		-- Пример триггера, обновляющего updated_at для таблицы users.
		-- Вы можете продублировать его для остальных таблиц при необходимости:
		CREATE TRIGGER IF NOT EXISTS update_users_updated_at
		AFTER UPDATE ON users
		BEGIN
			UPDATE users SET updated_at = CURRENT_TIMESTAMP WHERE t_id = NEW.t_id;
		END;
	`
	if _, err = db.Exec(sqlStmt); err != nil {
		cfg.Logger.Error("Failed to execute SQL script", "dbType", dbType, "error", err)
		db.Close()
		return nil, fmt.Errorf("failed to execute SQL script for %s database: %v", dbType, err)
	}

	if err = applySchemaMigrations(db, cfg); err != nil {
		cfg.Logger.Error("Failed to apply schema migrations", "dbType", dbType, "error", err)
		db.Close()
		return nil, fmt.Errorf("failed to apply schema migrations for %s database: %v", dbType, err)
	}
	if err = ensurePanelAuthDefaults(db, cfg); err != nil {
		cfg.Logger.Error("Failed to ensure panel auth defaults", "dbType", dbType, "error", err)
		db.Close()
		return nil, fmt.Errorf("failed to ensure panel auth defaults for %s database: %v", dbType, err)
	}

	cfg.Logger.Debug("Database initialized", "dbType", dbType)
	return db, nil
}

func applySchemaMigrations(db *sql.DB, cfg *config.BackendConfig) error {
	if err := ensureSubscriptionSettingsConnectionColumns(db); err != nil {
		return fmt.Errorf("ensure subscription_settings connection columns: %w", err)
	}
	if err := ensureAdminSessionTTLColumn(db); err != nil {
		return fmt.Errorf("ensure admin.session_ttl_minutes column: %w", err)
	}
	if err := ensureAdminSessionsTable(db); err != nil {
		return fmt.Errorf("ensure admin_sessions table: %w", err)
	}
	if err := ensureNodesWithoutAPIMetadata(db, cfg); err != nil {
		return fmt.Errorf("remove nodes.api_metadata: %w", err)
	}
	return nil
}

func ensureSubscriptionSettingsConnectionColumns(db *sql.DB) error {
	columns := []struct {
		name string
		def  string
	}{
		{name: "address", def: "address TEXT NOT NULL DEFAULT ''"},
		{name: "port", def: "port INTEGER NOT NULL DEFAULT 9263"},
		{name: "api_schema", def: "api_schema TEXT NOT NULL DEFAULT 'grpc'"},
		{name: "api_path", def: "api_path TEXT NOT NULL DEFAULT ''"},
	}

	for _, column := range columns {
		exists, err := columnExists(db, "subscription_settings", column.name)
		if err != nil {
			return err
		}
		if exists {
			continue
		}

		if _, err := db.Exec(fmt.Sprintf("ALTER TABLE subscription_settings ADD COLUMN %s", column.def)); err != nil {
			return fmt.Errorf("add column %s to subscription_settings: %w", column.name, err)
		}
	}

	return nil
}

func ensureAdminSessionTTLColumn(db *sql.DB) error {
	exists, err := columnExists(db, "admin", "session_ttl_minutes")
	if err != nil {
		return err
	}
	if !exists {
		if _, err := db.Exec("ALTER TABLE admin ADD COLUMN session_ttl_minutes INTEGER NOT NULL DEFAULT 60"); err != nil {
			return fmt.Errorf("add column session_ttl_minutes to admin: %w", err)
		}
	}

	if _, err := db.Exec("UPDATE admin SET session_ttl_minutes = 60 WHERE session_ttl_minutes IS NULL OR session_ttl_minutes <= 0"); err != nil {
		return fmt.Errorf("normalize admin.session_ttl_minutes defaults: %w", err)
	}
	return nil
}

func ensureAdminSessionsTable(db *sql.DB) error {
	if _, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS admin_sessions (
			session_token TEXT PRIMARY KEY,
			admin_uuid TEXT NOT NULL,
			expires_at INTEGER NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (admin_uuid) REFERENCES admin(uuid) ON DELETE CASCADE
		)`); err != nil {
		return fmt.Errorf("create admin_sessions table: %w", err)
	}
	if _, err := db.Exec("CREATE INDEX IF NOT EXISTS idx_admin_sessions_admin_uuid ON admin_sessions(admin_uuid)"); err != nil {
		return fmt.Errorf("create idx_admin_sessions_admin_uuid: %w", err)
	}
	if _, err := db.Exec("CREATE INDEX IF NOT EXISTS idx_admin_sessions_expires_at ON admin_sessions(expires_at)"); err != nil {
		return fmt.Errorf("create idx_admin_sessions_expires_at: %w", err)
	}
	return nil
}

func ensurePanelAuthDefaults(db *sql.DB, _ *config.BackendConfig) error {
	if err := ensureRemnawaveSettingsDefaultRow(db); err != nil {
		return fmt.Errorf("ensure remnawave_settings default row: %w", err)
	}
	if err := ensureAdminSessionTTLDefaults(db); err != nil {
		return fmt.Errorf("ensure admin session_ttl_minutes defaults: %w", err)
	}
	return nil
}

func ensureRemnawaveSettingsDefaultRow(db *sql.DB) error {
	var exists int
	if err := db.QueryRow("SELECT COUNT(*) FROM remnawave_settings WHERE id = 1").Scan(&exists); err != nil {
		return fmt.Errorf("check remnawave_settings row: %w", err)
	}
	if exists > 0 {
		return nil
	}

	passkeySettings := `{"rpId":null,"origin":null,"enabled":false}`
	oauth2Settings := `{"github":{"enabled":false,"clientId":null,"clientSecret":null,"allowedEmails":[]},"yandex":{"enabled":false,"clientId":null,"clientSecret":null,"allowedEmails":[]},"generic":{"enabled":false,"clientId":null,"tokenUrl":null,"withPkce":false,"clientSecret":null,"allowedEmails":[],"frontendDomain":null,"authorizationUrl":null},"keycloak":{"realm":null,"enabled":false,"clientId":null,"clientSecret":null,"allowedEmails":[],"frontendDomain":null,"keycloakDomain":null},"pocketid":{"enabled":false,"clientId":null,"plainDomain":null,"clientSecret":null,"allowedEmails":[]}}`
	tgAuthSettings := `{"enabled":false,"adminIds":[],"botToken":null}`
	passwordSettings := `{"enabled":true}`
	brandingSettings := `{"title":"V2RayStat","logoUrl":null}`

	if _, err := db.Exec(`
		INSERT INTO remnawave_settings (
			id, passkey_settings, oauth2_settings, tg_auth_settings, password_settings, branding_settings
		) VALUES (1, ?, ?, ?, ?, ?)
	`, passkeySettings, oauth2Settings, tgAuthSettings, passwordSettings, brandingSettings); err != nil {
		return fmt.Errorf("insert default remnawave_settings row: %w", err)
	}

	return nil
}

func ensureAdminSessionTTLDefaults(db *sql.DB) error {
	_, err := db.Exec("UPDATE admin SET session_ttl_minutes = 60 WHERE session_ttl_minutes IS NULL OR session_ttl_minutes <= 0")
	return err
}

func ensureNodesWithoutAPIMetadata(db *sql.DB, cfg *config.BackendConfig) error {
	exists, err := columnExists(db, "nodes", "api_metadata")
	if err != nil {
		return err
	}
	if !exists {
		return nil
	}

	if _, err := db.Exec("ALTER TABLE nodes DROP COLUMN api_metadata"); err == nil {
		cfg.Logger.Info("Dropped nodes.api_metadata column using ALTER TABLE")
		return nil
	}

	cfg.Logger.Warn("ALTER TABLE DROP COLUMN failed, rebuilding nodes table without api_metadata")
	return rebuildNodesTableWithoutAPIMetadata(db)
}

func rebuildNodesTableWithoutAPIMetadata(db *sql.DB) (err error) {
	if _, err = db.Exec("PRAGMA foreign_keys = OFF"); err != nil {
		return fmt.Errorf("disable foreign keys: %w", err)
	}
	defer func() {
		_, _ = db.Exec("PRAGMA foreign_keys = ON")
	}()

	tx, err := db.BeginTx(context.Background(), nil)
	if err != nil {
		return fmt.Errorf("begin nodes rebuild transaction: %w", err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	if _, err = tx.Exec("DROP TABLE IF EXISTS nodes_new"); err != nil {
		return fmt.Errorf("drop stale nodes_new table: %w", err)
	}

	if _, err = tx.Exec(`
		CREATE TABLE nodes_new (
			uuid TEXT PRIMARY KEY,
			id INTEGER UNIQUE,
			name TEXT NOT NULL,
			address TEXT NOT NULL,
			port INTEGER,
			api_schema TEXT DEFAULT 'grpc',
			api_path TEXT DEFAULT '',
			active_config_profile_uuid TEXT,
			is_connected BOOLEAN DEFAULT 0,
			is_connecting BOOLEAN DEFAULT 0,
			is_disabled BOOLEAN DEFAULT 0,
			last_status_change DATETIME,
			last_status_message TEXT,
			xray_version TEXT,
			node_version TEXT,
			xray_uptime TEXT DEFAULT '0',
			users_online INTEGER DEFAULT 0,
			consumption_multiplier INTEGER DEFAULT 1000000000,
			is_traffic_tracking_active BOOLEAN DEFAULT 0,
			traffic_reset_day INTEGER DEFAULT 1,
			traffic_limit_bytes INTEGER DEFAULT 0,
			traffic_used_bytes INTEGER DEFAULT 0,
			notify_percent INTEGER DEFAULT 0,
			provider_uuid TEXT,
			view_position INTEGER DEFAULT 0,
			country_code TEXT DEFAULT 'XX',
			tags TEXT DEFAULT '[]',
			cpu_count INTEGER,
			cpu_model TEXT,
			total_ram TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (active_config_profile_uuid) REFERENCES config_profiles(uuid) ON DELETE SET NULL,
			FOREIGN KEY (provider_uuid) REFERENCES infra_providers(uuid) ON DELETE SET NULL
		)`); err != nil {
		return fmt.Errorf("create nodes_new: %w", err)
	}

	if _, err = tx.Exec(`
		INSERT INTO nodes_new (
			uuid, id, name, address, port, api_schema, api_path, active_config_profile_uuid,
			is_connected, is_connecting, is_disabled, last_status_change, last_status_message,
			xray_version, node_version, xray_uptime, users_online, consumption_multiplier,
			is_traffic_tracking_active, traffic_reset_day, traffic_limit_bytes, traffic_used_bytes,
			notify_percent, provider_uuid, view_position, country_code, tags,
			cpu_count, cpu_model, total_ram, created_at, updated_at
		)
		SELECT
			uuid, id, name, address, port, api_schema, api_path, active_config_profile_uuid,
			is_connected, is_connecting, is_disabled, last_status_change, last_status_message,
			xray_version, node_version, xray_uptime, users_online, consumption_multiplier,
			is_traffic_tracking_active, traffic_reset_day, traffic_limit_bytes, traffic_used_bytes,
			notify_percent, provider_uuid, view_position, country_code, tags,
			cpu_count, cpu_model, total_ram, created_at, updated_at
		FROM nodes`); err != nil {
		return fmt.Errorf("copy data to nodes_new: %w", err)
	}

	if _, err = tx.Exec("DROP TABLE nodes"); err != nil {
		return fmt.Errorf("drop old nodes table: %w", err)
	}

	if _, err = tx.Exec("ALTER TABLE nodes_new RENAME TO nodes"); err != nil {
		return fmt.Errorf("rename nodes_new to nodes: %w", err)
	}

	if _, err = tx.Exec("CREATE INDEX IF NOT EXISTS idx_nodes_name ON nodes(name)"); err != nil {
		return fmt.Errorf("create idx_nodes_name: %w", err)
	}
	if _, err = tx.Exec("CREATE INDEX IF NOT EXISTS idx_nodes_address ON nodes(address)"); err != nil {
		return fmt.Errorf("create idx_nodes_address: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("commit nodes rebuild transaction: %w", err)
	}

	return nil
}

func columnExists(db *sql.DB, tableName, columnName string) (bool, error) {
	rows, err := db.Query(fmt.Sprintf("PRAGMA table_info(%s)", tableName))
	if err != nil {
		return false, fmt.Errorf("query table info for %s: %w", tableName, err)
	}
	defer rows.Close()

	for rows.Next() {
		var cid int
		var name, ctype string
		var notNull, pk int
		var defaultValue sql.NullString
		if err := rows.Scan(&cid, &name, &ctype, &notNull, &defaultValue, &pk); err != nil {
			return false, fmt.Errorf("scan table info for %s: %w", tableName, err)
		}
		if strings.EqualFold(name, columnName) {
			return true, nil
		}
	}

	if err := rows.Err(); err != nil {
		return false, fmt.Errorf("iterate table info for %s: %w", tableName, err)
	}

	return false, nil
}

// InitDatabase initializes in-memory and file databases.
func InitDatabase(cfg *config.BackendConfig) (memDB, fileDB *sql.DB, err error) {
	memDB, err = OpenAndInitDB(":memory:", "in-memory", cfg)
	if err != nil {
		return nil, nil, err
	}
	defer func() {
		if err != nil {
			memDB.Close()
		}
	}()

	if _, err := os.Stat(cfg.Paths.Database); os.IsNotExist(err) {
		cfg.Logger.Warn("File database does not exist, will create new", "path", cfg.Paths.Database)
	} else if err != nil {
		cfg.Logger.Error("Failed to check file database", "path", cfg.Paths.Database, "error", err)
		return nil, nil, fmt.Errorf("error checking file database: %v", err)
	}

	fileDB, err = OpenAndInitDB(cfg.Paths.Database, "file", cfg)
	if err != nil {
		return nil, nil, err
	}
	defer func() {
		if err != nil {
			fileDB.Close()
		}
	}()

	if err = ensureDefaultSubscriptionData(fileDB, cfg); err != nil {
		cfg.Logger.Error("Failed to ensure default subscription data", "error", err)
		return nil, nil, fmt.Errorf("failed to ensure default subscription data: %w", err)
	}

	var tableCount int
	err = fileDB.QueryRow("SELECT count(*) FROM sqlite_master WHERE type='table' AND name='user_traffic'").Scan(&tableCount)
	if err == nil && tableCount > 0 {
		tempManager, err := manager.NewDatabaseManager(fileDB, context.Background(), 1, 300, 500, cfg)
		if err != nil {
			cfg.Logger.Error("Failed to create temporary DatabaseManager", "error", err)
			return nil, nil, fmt.Errorf("failed to create temporary DatabaseManager: %v", err)
		}
		syncCtx, syncCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer syncCancel()
		if err = tempManager.SyncDBWithContext(syncCtx, memDB, "file to memory"); err != nil {
			cfg.Logger.Error("Failed to synchronize database (file to memory)", "error", err)
		}
		tempManager.Close()
	}

	cfg.Logger.Debug("Database initialization completed", "in-memory", true, "file", true)
	return memDB, fileDB, nil
}

func MonitorSubscriptionsAndSync(ctx context.Context, manager *manager.DatabaseManager, fileDB *sql.DB, cfg *config.BackendConfig, wg *sync.WaitGroup) {
	go func() {
		defer wg.Done()

		ticker := time.NewTicker(1 * time.Hour)
		defer ticker.Stop()

		cfg.Logger.Debug("Starting subscription and sync monitoring")

		for {
			select {
			case <-ticker.C:
				fileDB = performSyncCycle(ctx, manager, fileDB, cfg)

			case <-ctx.Done():
				cfg.Logger.Debug("Stopped subscription and sync monitoring")
				return
			}
		}
	}()
}

func performSyncCycle(ctx context.Context, mgr *manager.DatabaseManager, currentFileDB *sql.DB, cfg *config.BackendConfig) *sql.DB {
	if _, err := os.Stat(cfg.Paths.Database); os.IsNotExist(err) {
		cfg.Logger.Warn("File database does not exist, recreating", "path", cfg.Paths.Database)

		newDB, err := OpenAndInitDB(cfg.Paths.Database, "file", cfg)
		if err != nil {
			cfg.Logger.Error("Failed to recreate file database", "error", err)
			return currentFileDB
		}

		if currentFileDB != nil {
			currentFileDB.Close()
		}
		currentFileDB = newDB
	}

	syncCtx, syncCancel := context.WithTimeout(ctx, 15*time.Second)
	defer syncCancel()

	if err := mgr.SyncDBWithContext(syncCtx, currentFileDB, "memory to file"); err != nil {
		cfg.Logger.Error("Failed to synchronize database", "error", err)
	} else {
		cfg.Logger.Info("Database synchronized successfully")
	}

	return currentFileDB
}
