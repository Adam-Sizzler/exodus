package config

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"v2ray-stat/common"
	"v2ray-stat/logger"

	"github.com/joho/godotenv"
)

type BackendConfig struct {
	Logger   *logger.Logger
	Log      LogConfig
	TZ       string
	V2RS     V2RSConfig
	Panel    PanelConfig
	Database DatabaseConfig
	Redis    RedisConfig
	CORS     CORSConfig
}

type PanelConfig struct {
	StaticDir         string
	AllowInsecureHTTP bool
	TrustedProxies    []string
	WebPort           int
	trustedProxyNets  []*net.IPNet
}

type CORSConfig struct {
	AllowedOrigins []string
}

type LogConfig struct {
	LogLevel string
	LogMode  string
}

type V2RSConfig struct {
	Address string
	Port    int
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
	TZ:       "UTC",
	V2RS: V2RSConfig{
		Address: "0.0.0.0",
		Port:    9243,
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
		AllowInsecureHTTP: false,
		TrustedProxies:    []string{},
		WebPort:           9242,
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

	fallbackLogger, _ := logger.NewLoggerWithValidation(defaultConfig.Log.LogLevel, defaultConfig.Log.LogMode, defaultConfig.TZ, os.Stderr)
	cfg.Logger = fallbackLogger

	if err := loadDotEnv(); err != nil {
		cfg.Logger.Warn("Failed to load .env file", "error", err)
	}

	applyEnvOverrides(&cfg)
	normalizePanelConfig(&cfg)

	if cfg.V2RS.Port < 1 || cfg.V2RS.Port > 65535 {
		cfg.Logger.Warn("Invalid backend port, using default", "port", cfg.V2RS.Port, "default", defaultConfig.V2RS.Port)
		cfg.V2RS.Port = defaultConfig.V2RS.Port
	}
	if cfg.Panel.WebPort < 1 || cfg.Panel.WebPort > 65535 {
		cfg.Logger.Warn("Invalid web port, using default", "port", cfg.Panel.WebPort, "default", defaultConfig.Panel.WebPort)
		cfg.Panel.WebPort = defaultConfig.Panel.WebPort
	}

	if strings.TrimSpace(cfg.TZ) == "" {
		cfg.TZ = defaultConfig.TZ
	}
	if _, err := time.LoadLocation(cfg.TZ); err != nil {
		cfg.Logger.Warn("Invalid timezone value, using default", "timezone", cfg.TZ, "default", defaultConfig.TZ)
		cfg.TZ = defaultConfig.TZ
	}

	realLogger, err := logger.NewLoggerWithValidation(cfg.Log.LogLevel, cfg.Log.LogMode, cfg.TZ, os.Stderr)
	if err != nil {
		return cfg, fmt.Errorf("failed to initialize logger: %w", err)
	}
	cfg.Logger = realLogger

	common.InitTimezone(cfg.TZ, cfg.Logger)
	cfg.Logger.Info("v2rs listen", "address", cfg.V2RS.Address, "port", cfg.V2RS.Port)
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

	if value := envFirst("LOG_LEVEL", "V2RS_LOG_LEVEL"); value != "" {
		cfg.Log.LogLevel = value
	}
	if value := envFirst("LOG_MODE", "V2RS_LOG_MODE"); value != "" {
		cfg.Log.LogMode = value
	}

	if value := envFirst("TZ"); value != "" {
		cfg.TZ = value
	}

	if value := envFirst("BACKEND_ADDRESS", "V2RS_LISTEN_ADDRESS"); value != "" {
		cfg.V2RS.Address = value
	}

	if value := envFirst("BACKEND_PORT"); value != "" {
		if parsed, err := strconv.Atoi(value); err == nil && parsed > 0 && parsed <= 65535 {
			cfg.V2RS.Port = parsed
		} else if cfg.Logger != nil {
			cfg.Logger.Warn("Invalid backend port value, ignoring", "value", value)
		}
	}

	if value := envFirst("APP_PORT"); value != "" {
		if parsed, err := strconv.Atoi(value); err == nil && parsed > 0 && parsed <= 65535 {
			cfg.Panel.WebPort = parsed
		} else if cfg.Logger != nil {
			cfg.Logger.Warn("Invalid APP_PORT value, ignoring", "value", value)
		}
	}

	if value := envFirst("V2RS_ALLOW_INSECURE_HTTP"); value != "" {
		if parsed, err := strconv.ParseBool(value); err == nil {
			cfg.Panel.AllowInsecureHTTP = parsed
		} else if cfg.Logger != nil {
			cfg.Logger.Warn("Invalid V2RS_ALLOW_INSECURE_HTTP value, ignoring", "value", value)
		}
	}

	if value := envFirst("V2RS_TRUSTED_PROXIES"); value != "" {
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
	if cfg.Panel.WebPort < 1 || cfg.Panel.WebPort > 65535 {
		cfg.Panel.WebPort = defaultConfig.Panel.WebPort
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
