package config

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"

	"cerberus/logger"

	"github.com/joho/godotenv"
)

type BackendConfig struct {
	Logger   *logger.Logger
	Log      LogConfig
	CERBERUS CERBERUSConfig
	Panel    PanelConfig
	Database DatabaseConfig
	Redis    RedisConfig
	CORS     CORSConfig
}

type PanelConfig struct {
	StaticDir         string
	BasePath          string
	AllowInsecureHTTP bool
	TrustedProxies    []string
	AppPort           int
	trustedProxyNets  []*net.IPNet
}

type CORSConfig struct {
	AllowedOrigins []string
}

type LogConfig struct {
	LogLevel string
	LogMode  string
}

type CERBERUSConfig struct {
	Address string
}

type DatabaseConfig struct {
	URL string
}

type RedisConfig struct {
	Host string
	Port int
}

var defaultConfig = BackendConfig{
	Log: LogConfig{
		LogLevel: "warn",
		LogMode:  "inclusive",
	},
	CERBERUS: CERBERUSConfig{
		Address: "0.0.0.0",
	},
	Database: DatabaseConfig{
		URL: "",
	},
	Redis: RedisConfig{
		Host: "",
		Port: 6379,
	},
	Panel: PanelConfig{
		StaticDir:         "/app/ui",
		BasePath:          "/",
		AllowInsecureHTTP: false,
		TrustedProxies:    []string{},
		AppPort:           3000,
		trustedProxyNets:  nil,
	},
	CORS: CORSConfig{
		AllowedOrigins: []string{},
	},
}

// LoadConfig loads backend configuration from .env and environment variables only.
// .env values are loaded first, then real environment variables win.
func LoadConfig() (BackendConfig, error) {
	cfg := defaultConfig

	fallbackLogger, _ := logger.NewLoggerWithValidation(defaultConfig.Log.LogLevel, defaultConfig.Log.LogMode, "UTC", os.Stderr)
	cfg.Logger = fallbackLogger

	if err := loadDotEnv(); err != nil {
		cfg.Logger.Warn("Failed to load .env file", "error", err)
	}

	applyEnvOverrides(&cfg)
	normalizePanelConfig(&cfg)

	if cfg.Panel.AppPort < 1 || cfg.Panel.AppPort > 65535 {
		cfg.Logger.Warn("Invalid app port, using default", "port", cfg.Panel.AppPort, "default", defaultConfig.Panel.AppPort)
		cfg.Panel.AppPort = defaultConfig.Panel.AppPort
	}

	realLogger, err := logger.NewLoggerWithValidation(cfg.Log.LogLevel, cfg.Log.LogMode, "UTC", os.Stderr)
	if err != nil {
		return cfg, fmt.Errorf("failed to initialize logger: %w", err)
	}
	cfg.Logger = realLogger

	cfg.Logger.Info("cerberus listen", "address", cfg.CERBERUS.Address, "port", cfg.Panel.AppPort)
	cfg.Logger.Debug("Configuration loaded from environment")

	return cfg, nil
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

	if value := envFirst("LOG_LEVEL", "CERBERUS_LOG_LEVEL"); value != "" {
		cfg.Log.LogLevel = value
	}
	if value := envFirst("LOG_MODE", "CERBERUS_LOG_MODE"); value != "" {
		cfg.Log.LogMode = value
	}

	if value := envFirst("APP_ADDRESS"); value != "" {
		cfg.CERBERUS.Address = value
	}

	if value := envFirst("APP_PORT"); value != "" {
		if parsed, err := strconv.Atoi(value); err == nil && parsed > 0 && parsed <= 65535 {
			cfg.Panel.AppPort = parsed
		} else if cfg.Logger != nil {
			cfg.Logger.Warn("Invalid APP_PORT value, ignoring", "value", value)
		}
	}
	if value := envFirst("APP_PATH"); value != "" {
		cfg.Panel.BasePath = value
	}

	if value := envFirst("CERBERUS_ALLOW_INSECURE_HTTP"); value != "" {
		if parsed, err := strconv.ParseBool(value); err == nil {
			cfg.Panel.AllowInsecureHTTP = parsed
		} else if cfg.Logger != nil {
			cfg.Logger.Warn("Invalid CERBERUS_ALLOW_INSECURE_HTTP value, ignoring", "value", value)
		}
	}

	if value := envFirst("CERBERUS_TRUSTED_PROXIES"); value != "" {
		cfg.Panel.TrustedProxies = splitCSV(value)
	}

	if value := envFirst("DATABASE_URL"); value != "" {
		cfg.Database.URL = value
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

	if value := envFirst("FRONT_END_DOMAIN"); value != "" {
		cfg.CORS.AllowedOrigins = splitCSV(value)
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
	if cfg.Panel.AppPort < 1 || cfg.Panel.AppPort > 65535 {
		cfg.Panel.AppPort = defaultConfig.Panel.AppPort
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
