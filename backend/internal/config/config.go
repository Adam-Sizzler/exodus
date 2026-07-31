package config

import (
	"encoding/json"
	"fmt"
	"net"
	"net/url"
	"os"
	"strconv"
	"strings"

	"exodus/internal/logger"

	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

type BackendConfig struct {
	Logger        *logger.Logger
	Log           LogConfig
	EXODUS        EXODUSConfig
	Panel         PanelConfig
	Docs          DocsConfig
	Metrics       MetricsConfig
	Notifications NotificationsConfig
	Scheduler     SchedulerConfig
	Database      DatabaseConfig
	Redis         RedisConfig
	JWT           JWTConfig
	CORS          CORSConfig
}

type PanelConfig struct {
	StaticDir         string
	BasePath          string
	AllowInsecureHTTP bool
	TrustedProxies    []string
	AppPort           int
	ShortUUIDLength   int
	trustedProxyNets  []*net.IPNet
}

type DocsConfig struct {
	IsEnabled   bool
	ScalarPath  string
	SwaggerPath string
}

type CORSConfig struct {
	AllowedOrigins []string
}

type JWTConfig struct {
	AuthSecret        string
	APITokensSecret   string
	AuthLifetimeHours int
}

type MetricsConfig struct {
	Address         string
	Port            int
	User            string
	Pass            string
	CacheTTLSeconds int
}

type SchedulerConfig struct {
	ServiceCleanUsageHistory                 bool
	NotificationsEnabled                     bool
	BandwidthUsageNotificationsEnabled       bool
	BandwidthUsageNotificationsThreshold     []int
	NotConnectedUsersNotificationsEnabled    bool
	NotConnectedUsersNotificationsAfterHours []int
	ExpirationNotificationsEnabled           bool
	ExpirationNotifications                  []int
}

type NotificationsConfig struct {
	TelegramEnabled         bool
	TelegramBotToken        string
	TelegramBotProxy        string
	TelegramUsersChatID     string
	TelegramUsersThreadID   string
	TelegramNodesChatID     string
	TelegramNodesThreadID   string
	TelegramCRMChatID       string
	TelegramCRMThreadID     string
	TelegramServiceChatID   string
	TelegramServiceThreadID string

	WebhookEnabled bool
	WebhookURLs    []string
	WebhookSecret  string

	ConfigPath    string
	EventChannels map[string]NotificationEventChannelConfig
}

type NotificationEventChannelConfig struct {
	Telegram *bool `yaml:"telegram" json:"telegram"`
	Webhook  *bool `yaml:"webhook" json:"webhook"`
}

type LogConfig struct {
	LogLevel             string
	LogFormat            string
	EnableDebugLogs      bool
	IsHTTPLoggingEnabled bool
	NodeEnv              string
}

type EXODUSConfig struct {
	Address string
}

type DatabaseConfig struct {
	URL                string
	Socket             string
}

type RedisConfig struct {
	Host                         string
	Port                         int
	Password                     string
	DB                           int
	Socket                       string
	UserUsageHistoryTTLSeconds   int
	UserUsageHistoryDelaySeconds int
	DisableUserUsageRecords      bool
	UserUsageIgnoreBelowBytes    int64
	JobQueueVisibilitySeconds    int
	SubscriptionQueueConcurrency int
	PushToDBQueueConcurrency     int
}

var defaultConfig = BackendConfig{
	Log: LogConfig{
		LogLevel:             "info",
		LogFormat:            "console",
		EnableDebugLogs:      false,
		IsHTTPLoggingEnabled: false,
		NodeEnv:              "production",
	},
	EXODUS: EXODUSConfig{
		Address: "0.0.0.0",
	},
	Database: DatabaseConfig{
		URL: "",
	},
	Redis: RedisConfig{
		Host:                         "",
		Port:                         6379,
		Password:                     "",
		DB:                           0,
		Socket:                       "",
		UserUsageHistoryTTLSeconds:   10800,
		UserUsageHistoryDelaySeconds: 120,
		DisableUserUsageRecords:      false,
		UserUsageIgnoreBelowBytes:    0,
		JobQueueVisibilitySeconds:    300,
		SubscriptionQueueConcurrency: 50,
		PushToDBQueueConcurrency:     3,
	},
	Panel: PanelConfig{
		StaticDir:         "/opt/app/ui",
		BasePath:          "/",
		AllowInsecureHTTP: false,
		TrustedProxies:    []string{},
		AppPort:           3000,
		ShortUUIDLength:   16,
		trustedProxyNets:  nil,
	},
	JWT: JWTConfig{
		AuthLifetimeHours: 12,
	},
	Docs: DocsConfig{
		IsEnabled:   false,
		ScalarPath:  "/scalar",
		SwaggerPath: "/docs",
	},
	Metrics: MetricsConfig{
		Address:         "0.0.0.0",
		Port:            3001,
		User:            "",
		Pass:            "",
		CacheTTLSeconds: 10,
	},
	Notifications: NotificationsConfig{},
	Scheduler: SchedulerConfig{
		ServiceCleanUsageHistory:                 false,
		NotificationsEnabled:                     false,
		BandwidthUsageNotificationsEnabled:       false,
		BandwidthUsageNotificationsThreshold:     nil,
		NotConnectedUsersNotificationsEnabled:    false,
		NotConnectedUsersNotificationsAfterHours: nil,
		ExpirationNotificationsEnabled:           false,
		ExpirationNotifications:                  nil,
	},
	CORS: CORSConfig{
		AllowedOrigins: []string{},
	},
}

// LoadConfig loads backend configuration from .env and environment variables only.
// .env values are loaded first, then real environment variables win.
func LoadConfig() (BackendConfig, error) {
	cfg := defaultConfig

	fallbackLogger, _ := logger.NewLoggerWithValidation(defaultConfig.Log.LogLevel, "UTC", os.Stderr)
	cfg.Logger = fallbackLogger

	if err := loadDotEnv(); err != nil {
		cfg.Logger.Warn("Failed to load .env file", "error", err)
	}

	applyEnvOverrides(&cfg)
	loadNotificationsYAMLConfig(&cfg)
	normalizePanelConfig(&cfg)

	if cfg.Panel.AppPort < 1 || cfg.Panel.AppPort > 65535 {
		cfg.Logger.Warn("Invalid app port, using default", "port", cfg.Panel.AppPort, "default", defaultConfig.Panel.AppPort)
		cfg.Panel.AppPort = defaultConfig.Panel.AppPort
	}
	if cfg.Metrics.Port < 1 || cfg.Metrics.Port > 65535 {
		cfg.Logger.Warn("Invalid metrics port, using default", "port", cfg.Metrics.Port, "default", defaultConfig.Metrics.Port)
		cfg.Metrics.Port = defaultConfig.Metrics.Port
	}

	realLogger, err := logger.NewLoggerFromEnv(cfg.Log.LogLevel, cfg.Log.LogFormat, "UTC", os.Stderr)
	if err != nil {
		return cfg, fmt.Errorf("failed to initialize logger: %w", err)
	}
	cfg.Logger = realLogger

	if strings.TrimSpace(cfg.Metrics.User) == "" {
		return cfg, fmt.Errorf("METRICS_USER cannot be empty")
	}
	if strings.TrimSpace(cfg.Metrics.Pass) == "" {
		return cfg, fmt.Errorf("METRICS_PASS cannot be empty")
	}
	if strings.TrimSpace(cfg.Database.URL) == "" {
		return cfg, fmt.Errorf("DATABASE_URL is not set")
	}

	cfg.Logger.RoleService(logger.RoleAPI, logger.ServiceConfig).Debug("Configuration loaded from environment", "node_env", cfg.Log.NodeEnv, "log_format", cfg.Log.LogFormat, "log_level", cfg.Log.LogLevel)

	return cfg, nil
}

func resolveConfiguredLogLevel(value string) string {
	level := strings.ToLower(strings.TrimSpace(value))
	switch level {
	case "trace", "verbose":
		return "trace"
	case "debug":
		return "debug"
	case "warn", "warning":
		return "warn"
	case "error":
		return "error"
	case "none", "silent":
		return "none"
	case "", "info", "log":
		return "info"
	default:
		return "info"
	}
}

func isDevelopmentEnv(value string) bool {
	env := strings.ToLower(strings.TrimSpace(value))
	return env == "development" || env == "dev"
}

func isDebugOrTraceLogLevel(value string) bool {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "debug", "trace", "verbose":
		return true
	default:
		return false
	}
}

func loadDotEnv() error {
	err := godotenv.Load(".env")
	if err == nil || os.IsNotExist(err) {
		return nil
	}
	return err
}

func applyEnvOverrides(cfg *BackendConfig) {
	if cfg == nil {
		return
	}

	rawLogLevel := envFirst("LOG_LEVEL", "EXODUS_LOG_LEVEL")
	if rawLogLevel != "" {
		cfg.Log.LogLevel = resolveConfiguredLogLevel(rawLogLevel)
	}
	if value := envFirst("LOG_FORMAT", "EXODUS_LOG_FORMAT"); value != "" {
		cfg.Log.LogFormat = value
	}
	if value := envFirst("NODE_ENV"); value != "" {
		cfg.Log.NodeEnv = value
	}
	if value := envFirst("ENABLE_DEBUG_LOGS", "EXODUS_ENABLE_DEBUG_LOGS"); value != "" {
		cfg.Log.EnableDebugLogs = parseBoolEnv(value)
	}
	if value := envFirst("IS_HTTP_LOGGING_ENABLED"); value != "" {
		cfg.Log.IsHTTPLoggingEnabled = parseBoolEnv(value)
	}
	// Fallback to debug log level if ENABLE_DEBUG_LOGS=true or NODE_ENV=development and LOG_LEVEL was not explicitly specified
	if rawLogLevel == "" && (cfg.Log.EnableDebugLogs || isDevelopmentEnv(cfg.Log.NodeEnv)) {
		cfg.Log.LogLevel = "debug"
	}

	if value := envFirst("APP_ADDRESS"); value != "" {
		cfg.EXODUS.Address = value
	}

	if value := envFirst("APP_PORT"); value != "" {
		if parsed, err := strconv.Atoi(value); err == nil && parsed > 0 && parsed <= 65535 {
			cfg.Panel.AppPort = parsed
		} else if cfg.Logger != nil {
			cfg.Logger.Warn("Invalid APP_PORT value, ignoring", "value", value)
		}
	}
	if value := envFirst("SHORT_UUID_LENGTH"); value != "" {
		if parsed, err := strconv.Atoi(value); err == nil && parsed >= 16 && parsed <= 64 {
			cfg.Panel.ShortUUIDLength = parsed
		} else if cfg.Logger != nil {
			cfg.Logger.Warn("Invalid SHORT_UUID_LENGTH value (must be 16-64), ignoring", "value", value)
		}
	}
	if value := envFirst("APP_PATH"); value != "" {
		cfg.Panel.BasePath = value
	}
	if value := envFirst("PANEL_STATIC_DIR", "EXODUS_STATIC_DIR"); value != "" {
		cfg.Panel.StaticDir = value
	}

	if value := envFirst("IS_DOCS_ENABLED"); value != "" {
		cfg.Docs.IsEnabled = strings.EqualFold(strings.TrimSpace(value), "true")
	}
	if value := envFirst("SCALAR_PATH"); value != "" {
		cfg.Docs.ScalarPath = strings.TrimSpace(value)
	}
	if value := envFirst("SWAGGER_PATH"); value != "" {
		cfg.Docs.SwaggerPath = strings.TrimSpace(value)
	}

	if value := envFirst("METRICS_ADDRESS"); value != "" {
		cfg.Metrics.Address = value
	}
	if value := envFirst("METRICS_PORT"); value != "" {
		if parsed, err := strconv.Atoi(value); err == nil && parsed > 0 && parsed <= 65535 {
			cfg.Metrics.Port = parsed
		} else if cfg.Logger != nil {
			cfg.Logger.Warn("Invalid METRICS_PORT value, ignoring", "value", value)
		}
	}
	if value := envFirst("METRICS_USER"); value != "" {
		cfg.Metrics.User = value
	}
	if value := envFirst("METRICS_PASS"); value != "" {
		cfg.Metrics.Pass = value
	}
	if value := envFirst("METRICS_CACHE_TTL_SECONDS"); value != "" {
		if parsed, err := strconv.Atoi(value); err == nil && parsed >= 0 {
			cfg.Metrics.CacheTTLSeconds = parsed
		} else if cfg.Logger != nil {
			cfg.Logger.Warn("Invalid METRICS_CACHE_TTL_SECONDS value, ignoring", "value", value)
		}
	}

	if value := envFirst("EXODUS_ALLOW_INSECURE_HTTP"); value != "" {
		if parsed, err := strconv.ParseBool(value); err == nil {
			cfg.Panel.AllowInsecureHTTP = parsed
		} else if cfg.Logger != nil {
			cfg.Logger.Warn("Invalid EXODUS_ALLOW_INSECURE_HTTP value, ignoring", "value", value)
		}
	}

	if value := envFirst("EXODUS_TRUSTED_PROXIES"); value != "" {
		cfg.Panel.TrustedProxies = splitCSV(value)
	}

	if value := envFirst("DATABASE_URL"); value != "" {
		cfg.Database.URL = value
	}
	if value := envFirst("DATABASE_SOCKET", "POSTGRES_SOCKET"); value != "" {
		cfg.Database.Socket = value
		cfg.Database.URL = postgresSocketDatabaseURL(value)
	}

	if value := envFirst("REDIS_HOST"); value != "" {
		cfg.Redis.Host = value
	}
	if value := envFirst("REDIS_PORT"); value != "" {
		if parsed, err := strconv.Atoi(value); err == nil && parsed > 0 {
			cfg.Redis.Port = parsed
		} else if cfg.Logger != nil {
			cfg.Logger.Warn("Invalid REDIS_PORT value, ignoring", "value", value)
		}
	}
	if value := envFirst("REDIS_PASSWORD"); value != "" {
		cfg.Redis.Password = value
	}
	if value := envFirst("REDIS_DB"); value != "" {
		if parsed, err := strconv.Atoi(value); err == nil && parsed >= 0 {
			cfg.Redis.DB = parsed
		} else if cfg.Logger != nil {
			cfg.Logger.Warn("Invalid REDIS_DB value, ignoring", "value", value)
		}
	}
	if value := envFirst("REDIS_SOCKET"); value != "" {
		cfg.Redis.Socket = value
	}
	if value := envFirst("NODE_USER_USAGE_REDIS_TTL_SECONDS", "EXODUS_NODE_USER_USAGE_REDIS_TTL_SECONDS"); value != "" {
		if parsed, err := strconv.Atoi(value); err == nil && parsed > 0 {
			cfg.Redis.UserUsageHistoryTTLSeconds = parsed
		} else if cfg.Logger != nil {
			cfg.Logger.Warn("Invalid NODE_USER_USAGE_REDIS_TTL_SECONDS value, ignoring", "value", value)
		}
	}
	if value := envFirst("NODE_USER_USAGE_FLUSH_DELAY_SECONDS", "EXODUS_NODE_USER_USAGE_FLUSH_DELAY_SECONDS"); value != "" {
		if parsed, err := strconv.Atoi(value); err == nil && parsed >= 0 {
			cfg.Redis.UserUsageHistoryDelaySeconds = parsed
		} else if cfg.Logger != nil {
			cfg.Logger.Warn("Invalid NODE_USER_USAGE_FLUSH_DELAY_SECONDS value, ignoring", "value", value)
		}
	}
	if value := envFirst("SERVICE_DISABLE_USER_USAGE_RECORDS"); value != "" {
		cfg.Redis.DisableUserUsageRecords = parseBoolEnv(value)
	}
	if value := envFirst("USER_USAGE_IGNORE_BELOW_BYTES"); value != "" {
		if parsed, err := strconv.ParseInt(value, 10, 64); err == nil && parsed >= 0 {
			cfg.Redis.UserUsageIgnoreBelowBytes = parsed
		} else if cfg.Logger != nil {
			cfg.Logger.Warn("Invalid USER_USAGE_IGNORE_BELOW_BYTES value, ignoring", "value", value)
		}
	}
	if value := envFirst("JOB_QUEUE_VISIBILITY_SECONDS", "EXODUS_JOB_QUEUE_VISIBILITY_SECONDS"); value != "" {
		if parsed, err := strconv.Atoi(value); err == nil && parsed > 0 {
			cfg.Redis.JobQueueVisibilitySeconds = parsed
		} else if cfg.Logger != nil {
			cfg.Logger.Warn("Invalid JOB_QUEUE_VISIBILITY_SECONDS value, ignoring", "value", value)
		}
	}
	if value := envFirst("SUBSCRIPTION_QUEUE_CONCURRENCY", "EXODUS_SUBSCRIPTION_QUEUE_CONCURRENCY"); value != "" {
		if parsed, err := strconv.Atoi(value); err == nil && parsed > 0 {
			cfg.Redis.SubscriptionQueueConcurrency = parsed
		} else if cfg.Logger != nil {
			cfg.Logger.Warn("Invalid SUBSCRIPTION_QUEUE_CONCURRENCY value, ignoring", "value", value)
		}
	}
	if value := envFirst("PUSH_TO_DB_QUEUE_CONCURRENCY", "EXODUS_PUSH_TO_DB_QUEUE_CONCURRENCY"); value != "" {
		if parsed, err := strconv.Atoi(value); err == nil && parsed > 0 {
			cfg.Redis.PushToDBQueueConcurrency = parsed
		} else if cfg.Logger != nil {
			cfg.Logger.Warn("Invalid PUSH_TO_DB_QUEUE_CONCURRENCY value, ignoring", "value", value)
		}
	}

	if value := envFirst("FRONT_END_DOMAIN"); value != "" {
		cfg.CORS.AllowedOrigins = splitCSV(value)
	}

	cfg.JWT.AuthSecret = strings.TrimSpace(envFirst("JWT_AUTH_SECRET"))
	cfg.JWT.APITokensSecret = strings.TrimSpace(envFirst("JWT_API_TOKENS_SECRET"))
	if value := envFirst("JWT_AUTH_LIFETIME"); value != "" {
		if parsed, err := strconv.Atoi(value); err == nil && parsed >= 12 && parsed <= 168 {
			cfg.JWT.AuthLifetimeHours = parsed
		} else if cfg.Logger != nil {
			cfg.Logger.Warn("Invalid JWT_AUTH_LIFETIME value (must be between 12 and 168 hours), ignoring", "value", value)
		}
	}

	if value := envFirst("SERVICE_CLEAN_USAGE_HISTORY"); value != "" {
		cfg.Scheduler.ServiceCleanUsageHistory = parseBoolEnv(value)
	}
	cfg.Notifications.TelegramEnabled = parseBoolEnv(envFirst("IS_TELEGRAM_NOTIFICATIONS_ENABLED"))
	cfg.Notifications.TelegramBotToken = envFirst("TELEGRAM_BOT_TOKEN")
	cfg.Notifications.TelegramBotProxy = strings.TrimSpace(envFirst("TELEGRAM_BOT_PROXY"))
	cfg.Notifications.TelegramUsersChatID, cfg.Notifications.TelegramUsersThreadID = parseTelegramTarget(envFirst("TELEGRAM_NOTIFY_USERS"))
	cfg.Notifications.TelegramNodesChatID, cfg.Notifications.TelegramNodesThreadID = parseTelegramTarget(envFirst("TELEGRAM_NOTIFY_NODES"))
	cfg.Notifications.TelegramCRMChatID, cfg.Notifications.TelegramCRMThreadID = parseTelegramTarget(envFirst("TELEGRAM_NOTIFY_CRM"))
	cfg.Notifications.TelegramServiceChatID, cfg.Notifications.TelegramServiceThreadID = parseTelegramTarget(envFirst("TELEGRAM_NOTIFY_SERVICE"))
	cfg.Notifications.WebhookEnabled = parseBoolEnv(envFirst("WEBHOOK_ENABLED"))
	cfg.Notifications.WebhookURLs = splitCSV(envFirst("WEBHOOK_URL"))
	cfg.Notifications.WebhookSecret = envFirst("WEBHOOK_SECRET_HEADER")
	cfg.Scheduler.NotificationsEnabled = cfg.Notifications.TelegramEnabled || cfg.Notifications.WebhookEnabled
	if value := envFirst("BANDWIDTH_USAGE_NOTIFICATIONS_ENABLED"); value != "" {
		cfg.Scheduler.BandwidthUsageNotificationsEnabled = parseBoolEnv(value)
	}
	if value := envFirst("BANDWIDTH_USAGE_NOTIFICATIONS_THRESHOLD"); value != "" {
		parsed, err := parseIntJSONArray(value)
		if err != nil {
			if cfg.Logger != nil {
				cfg.Logger.Warn("Invalid BANDWIDTH_USAGE_NOTIFICATIONS_THRESHOLD value, ignoring", "error", err)
			}
		} else {
			cfg.Scheduler.BandwidthUsageNotificationsThreshold = parsed
		}
	}
	if value := envFirst("NOT_CONNECTED_USERS_NOTIFICATIONS_ENABLED"); value != "" {
		cfg.Scheduler.NotConnectedUsersNotificationsEnabled = parseBoolEnv(value)
	}
	if value := envFirst("NOT_CONNECTED_USERS_NOTIFICATIONS_AFTER_HOURS"); value != "" {
		parsed, err := parseIntJSONArray(value)
		if err != nil {
			if cfg.Logger != nil {
				cfg.Logger.Warn("Invalid NOT_CONNECTED_USERS_NOTIFICATIONS_AFTER_HOURS value, ignoring", "error", err)
			}
		} else {
			cfg.Scheduler.NotConnectedUsersNotificationsAfterHours = parsed
		}
	}
	if value := envFirst("EXPIRATION_NOTIFICATIONS_ENABLED"); value != "" {
		cfg.Scheduler.ExpirationNotificationsEnabled = parseBoolEnv(value)
	}
	if value := envFirst("EXPIRATION_NOTIFICATIONS"); value != "" {
		parsed, err := parseExpirationNotificationIntervals(value)
		if err != nil {
			if cfg.Logger != nil {
				cfg.Logger.Warn("Invalid EXPIRATION_NOTIFICATIONS value, ignoring", "error", err)
			}
		} else {
			cfg.Scheduler.ExpirationNotifications = parsed
		}
	}
}

func parseTelegramTarget(rawTarget string) (string, string) {
	target := strings.TrimSpace(rawTarget)
	if target == "" {
		return "", ""
	}
	if !strings.Contains(target, ":") {
		return target, ""
	}

	parts := strings.SplitN(target, ":", 2)
	return strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])
}

func loadNotificationsYAMLConfig(cfg *BackendConfig) {
	if cfg == nil {
		return
	}

	explicitPath := envFirst("NOTIFICATIONS_CONFIG_PATH", "EXODUS_NOTIFICATIONS_CONFIG_PATH")
	candidates := []string{}
	if explicitPath != "" {
		candidates = append(candidates, explicitPath)
	} else {
		candidates = append(candidates,
			"/var/lib/exodus/configs/notifications/notifications-config.yml",
			"configs/notifications/notifications-config.yml",
		)
	}

	var selected string
	for _, candidate := range candidates {
		if strings.TrimSpace(candidate) == "" {
			continue
		}
		if _, err := os.Stat(candidate); err == nil {
			selected = candidate
			break
		}
	}
	if selected == "" {
		if explicitPath != "" && cfg.Logger != nil {
			cfg.Logger.Warn("Notifications config file not found", "path", explicitPath)
		}
		cfg.Notifications.EventChannels = map[string]NotificationEventChannelConfig{}
		return
	}

	content, err := os.ReadFile(selected)
	if err != nil {
		if cfg.Logger != nil {
			cfg.Logger.Warn("Failed to read notifications config", "path", selected, "error", err)
		}
		cfg.Notifications.EventChannels = map[string]NotificationEventChannelConfig{}
		return
	}

	var parsed struct {
		Events map[string]NotificationEventChannelConfig `yaml:"events"`
	}
	if err := yaml.Unmarshal(content, &parsed); err != nil {
		if cfg.Logger != nil {
			cfg.Logger.Warn("Failed to parse notifications config", "path", selected, "error", err)
		}
		cfg.Notifications.EventChannels = map[string]NotificationEventChannelConfig{}
		return
	}
	if parsed.Events == nil {
		parsed.Events = map[string]NotificationEventChannelConfig{}
	}
	cfg.Notifications.ConfigPath = selected
	cfg.Notifications.EventChannels = parsed.Events
	if cfg.Logger != nil {
		cfg.Logger.Info("Notifications event config loaded", "path", selected, "events", len(parsed.Events))
	}
}

func (c NotificationsConfig) EventChannelEnabled(eventName, channel string) bool {
	if strings.TrimSpace(eventName) == "" || strings.TrimSpace(channel) == "" {
		return true
	}
	eventConfig, exists := c.EventChannels[eventName]
	if !exists {
		return true
	}
	switch strings.ToLower(strings.TrimSpace(channel)) {
	case "telegram":
		if eventConfig.Telegram == nil {
			return true
		}
		return *eventConfig.Telegram
	case "webhook":
		if eventConfig.Webhook == nil {
			return true
		}
		return *eventConfig.Webhook
	default:
		return true
	}
}

func normalizePanelConfig(cfg *BackendConfig) {
	if cfg == nil {
		return
	}

	if strings.TrimSpace(cfg.Panel.StaticDir) == "" {
		cfg.Panel.StaticDir = defaultConfig.Panel.StaticDir
	}
	if strings.TrimSpace(cfg.Panel.BasePath) == "" {
		cfg.Panel.BasePath = defaultConfig.Panel.BasePath
	}
	cfg.Panel.BasePath = normalizeBasePath(cfg.Panel.BasePath)
	// Docs paths are endpoint mount paths, not application base paths. They must
	// stay without a trailing slash; otherwise net/http treats e.g. "/docs/" as a
	// subtree route and emits a root-relative redirect to "/docs/" for requests
	// to "/docs", which drops APP_PATH behind a reverse proxy.
	cfg.Docs.ScalarPath = normalizeRoutePath(cfg.Docs.ScalarPath)
	cfg.Docs.SwaggerPath = normalizeRoutePath(cfg.Docs.SwaggerPath)
	if cfg.Panel.AppPort < 1 || cfg.Panel.AppPort > 65535 {
		cfg.Panel.AppPort = defaultConfig.Panel.AppPort
	}
	if strings.TrimSpace(cfg.Metrics.Address) == "" {
		cfg.Metrics.Address = defaultConfig.Metrics.Address
	}
	if cfg.Metrics.Port < 1 || cfg.Metrics.Port > 65535 {
		cfg.Metrics.Port = defaultConfig.Metrics.Port
	}
	if cfg.Metrics.CacheTTLSeconds < 0 {
		cfg.Metrics.CacheTTLSeconds = defaultConfig.Metrics.CacheTTLSeconds
	}

	proxyNets, invalid := parseTrustedProxies(cfg.Panel.TrustedProxies)
	if len(invalid) > 0 && cfg.Logger != nil {
		cfg.Logger.Warn("Invalid trusted proxy CIDR entries ignored", "entries", invalid)
	}
	cfg.Panel.trustedProxyNets = proxyNets
}

func parseTrustedProxies(entries []string) ([]*net.IPNet, []string) {
	var (
		nets    []*net.IPNet
		invalid []string
	)

	for _, entry := range entries {
		value := strings.TrimSpace(entry)
		if value == "" {
			continue
		}
		if !strings.Contains(value, "/") {
			if ip := net.ParseIP(value); ip != nil {
				bits := 32
				if ip.To4() == nil {
					bits = 128
				}
				_, netBlock, err := net.ParseCIDR(fmt.Sprintf("%s/%d", ip.String(), bits))
				if err == nil {
					nets = append(nets, netBlock)
					continue
				}
			}
			invalid = append(invalid, value)
			continue
		}
		_, netBlock, err := net.ParseCIDR(value)
		if err != nil {
			invalid = append(invalid, value)
			continue
		}
		nets = append(nets, netBlock)
	}

	return nets, invalid
}

func splitCSV(input string) []string {
	parts := strings.Split(input, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		value := strings.TrimSpace(part)
		if value != "" {
			result = append(result, value)
		}
	}
	return result
}

func parseBoolEnv(value string) bool {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "1", "true", "yes", "y", "on":
		return true
	default:
		return false
	}
}

func parseIntJSONArray(value string) ([]int, error) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return nil, nil
	}
	var result []int
	if err := json.Unmarshal([]byte(trimmed), &result); err != nil {
		return nil, err
	}
	return result, nil
}

func parseExpirationNotificationIntervals(value string) ([]int, error) {
	result, err := parseIntJSONArray(value)
	if err != nil {
		return nil, err
	}
	if len(result) == 0 {
		return nil, fmt.Errorf("EXPIRATION_NOTIFICATIONS must not be empty")
	}

	negativeCount := 0
	positiveCount := 0
	seen := make(map[int]struct{}, len(result))
	for i, interval := range result {
		if interval == 0 || interval < -168 || interval > 168 {
			return nil, fmt.Errorf("all expiration values must be non-zero integers between -168 and 168")
		}
		if i > 0 && interval <= result[i-1] {
			return nil, fmt.Errorf("EXPIRATION_NOTIFICATIONS values must be in strictly ascending order")
		}
		if _, exists := seen[interval]; exists {
			return nil, fmt.Errorf("EXPIRATION_NOTIFICATIONS must not contain duplicate values")
		}
		seen[interval] = struct{}{}
		if interval < 0 {
			negativeCount++
		} else {
			positiveCount++
		}
	}
	if negativeCount > 5 {
		return nil, fmt.Errorf("EXPIRATION_NOTIFICATIONS must contain at most 5 negative values")
	}
	if positiveCount > 5 {
		return nil, fmt.Errorf("EXPIRATION_NOTIFICATIONS must contain at most 5 positive values")
	}

	return result, nil
}

func postgresSocketDatabaseURL(socketDir string) string {
	socketDir = strings.TrimSpace(socketDir)
	if socketDir == "" {
		return ""
	}

	user := envFirst("POSTGRES_USER", "DATABASE_USER")
	password := envFirst("POSTGRES_PASSWORD", "DATABASE_PASSWORD")
	database := envFirst("POSTGRES_DB", "DATABASE_NAME")
	if user == "" || database == "" {
		return ""
	}

	dsn := url.URL{
		Scheme: "postgresql",
		Host:   "localhost",
		Path:   "/" + database,
	}
	if password != "" {
		dsn.User = url.UserPassword(user, password)
	} else {
		dsn.User = url.User(user)
	}

	query := dsn.Query()
	query.Set("host", socketDir)
	dsn.RawQuery = query.Encode()

	return dsn.String()
}

func normalizeBasePath(input string) string {
	trimmed := strings.TrimSpace(input)
	if trimmed == "" || trimmed == "/" {
		return "/"
	}

	cleaned := strings.Trim(trimmed, "/")
	if cleaned == "" {
		return "/"
	}

	return "/" + cleaned + "/"
}

// normalizeRoutePath normalizes a single endpoint mount path (e.g. docs paths)
// without a trailing slash, unlike normalizeBasePath.
func normalizeRoutePath(input string) string {
	trimmed := strings.TrimSpace(input)
	if trimmed == "" || trimmed == "/" {
		return "/"
	}

	cleaned := strings.Trim(trimmed, "/")
	if cleaned == "" {
		return "/"
	}

	return "/" + cleaned
}

func envFirst(keys ...string) string {
	for _, key := range keys {
		if value := strings.TrimSpace(os.Getenv(key)); value != "" {
			return value
		}
	}
	return ""
}

func (p PanelConfig) IsTrustedProxy(ip net.IP) bool {
	if ip == nil {
		return false
	}
	for _, netBlock := range p.trustedProxyNets {
		if netBlock != nil && netBlock.Contains(ip) {
			return true
		}
	}
	return false
}
