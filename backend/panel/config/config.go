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

	"gopkg.in/yaml.v3"
)

type BackendConfig struct {
	Logger   *logger.Logger
	Log      LogConfig      `yaml:"log_config"`
	APIToken string         `yaml:"api_token"`
	TZ       string         `yaml:"TZ"`
	V2RS     V2RSConfig     `yaml:"backend"`
	Panel    PanelConfig    `yaml:"panel"`
	Database DatabaseConfig `yaml:"database"`
	Redis    RedisConfig    `yaml:"redis"`
	CORS     CORSConfig     `yaml:"cors"`
}

type PanelConfig struct {
	StaticDir         string   `yaml:"static_dir"`
	AllowInsecureHTTP bool     `yaml:"allow_insecure_http"`
	TrustedProxies    []string `yaml:"trusted_proxies"`
	WebPort           int      `yaml:"web_port"`
	trustedProxyNets  []*net.IPNet
}

type CORSConfig struct {
	AllowedOrigins []string `yaml:"allowed_origins"`
	AllowedMethods []string `yaml:"allowed_methods"`
	AllowedHeaders []string `yaml:"allowed_headers"`
}

type LogConfig struct {
	LogLevel string `yaml:"level"`
	LogMode  string `yaml:"mode"`
}

type V2RSConfig struct {
	Address string `yaml:"listen_address"`
	Port    int    `yaml:"listen_port"`
}

type DatabaseConfig struct {
	URL string `yaml:"url"`
}

type RedisConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

var defaultConfig = BackendConfig{
	Log: LogConfig{
		LogLevel: "none",
		LogMode:  "inclusive",
	},
	APIToken: "",
	TZ:       "UTC",
	V2RS: V2RSConfig{
		Address: "127.0.0.1",
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
		StaticDir:         "./ui",
		AllowInsecureHTTP: false,
		TrustedProxies:    []string{},
		WebPort:           9242,
		trustedProxyNets:  nil,
	},
	CORS: CORSConfig{
		AllowedOrigins: []string{},
		AllowedMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Content-Type", "Authorization", "X-API-Token", "Cookie"},
	},
}

func LoadConfig(configFile string) (BackendConfig, error) {
	cfg := defaultConfig

	_, err := os.Stat(configFile)
	if err != nil {
		if os.IsNotExist(err) {
			cfg.Logger, _ = logger.NewLoggerWithValidation("warn", "inclusive", cfg.TZ, os.Stderr)
			cfg.Logger.Warn("Configuration file not found, using default values", "file", configFile)
			return cfg, nil
		}
		return cfg, fmt.Errorf("error accessing configuration file %s: %v", configFile, err)
	}

	data, err := os.ReadFile(configFile)
	if err != nil {
		return cfg, fmt.Errorf("error reading configuration file %s: %v", configFile, err)
	}
	cfg.Logger, _ = logger.NewLoggerWithValidation("debug", "inclusive", cfg.TZ, os.Stderr)

	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return cfg, fmt.Errorf("error parsing YAML configuration from %s: %v", configFile, err)
	}

	cfg.Logger, err = logger.NewLoggerWithValidation(cfg.Log.LogLevel, cfg.Log.LogMode, cfg.TZ, os.Stderr)
	if err != nil {
		return cfg, fmt.Errorf("failed to initialize logger: %v", err)
	}

	applyEnvOverrides(&cfg)
	normalizePanelConfig(&cfg)

	// Validate configuration
	if cfg.V2RS.Port != 0 {
		if cfg.V2RS.Port < 1 || cfg.V2RS.Port > 65535 {
			cfg.Logger.Warn("Invalid v2rs.port, using default", "port", cfg.V2RS.Port, "default", defaultConfig.V2RS.Port)
			cfg.V2RS.Port = defaultConfig.V2RS.Port
		}
		cfg.Logger.Info("v2rs listen", "address", cfg.V2RS.Address, "port", cfg.V2RS.Port)
	}

	if cfg.TZ != "" {
		if _, err := time.LoadLocation(cfg.TZ); err != nil {
			cfg.Logger.Warn("Invalid timezone value, using default", "timezone", cfg.TZ)
			cfg.TZ = defaultConfig.TZ
		}
	}

	if cfg.TZ != "" {
		common.InitTimezone(cfg.TZ, cfg.Logger)
	} else {
		cfg.TZ = "UTC"
		common.InitTimezone("UTC", cfg.Logger)
	}

	cfg.Logger.Debug("Configuration validated")
	return cfg, nil
}

func applyEnvOverrides(cfg *BackendConfig) {
	if cfg == nil {
		return
	}

	if value := strings.TrimSpace(os.Getenv("BACKEND_PORT")); value != "" {
		if parsed, err := strconv.Atoi(value); err == nil && parsed > 0 && parsed <= 65535 {
			cfg.V2RS.Port = parsed
		} else if cfg.Logger != nil {
			cfg.Logger.Warn("Invalid BACKEND_PORT value, ignoring", "value", value)
		}
	}

	if value := strings.TrimSpace(os.Getenv("WEB_PORT")); value != "" {
		if parsed, err := strconv.Atoi(value); err == nil && parsed > 0 && parsed <= 65535 {
			cfg.Panel.WebPort = parsed
		} else if cfg.Logger != nil {
			cfg.Logger.Warn("Invalid WEB_PORT value, ignoring", "value", value)
		}
	}

	if value := strings.TrimSpace(os.Getenv("V2RS_ALLOW_INSECURE_HTTP")); value != "" {
		if parsed, err := strconv.ParseBool(value); err == nil {
			cfg.Panel.AllowInsecureHTTP = parsed
		} else if cfg.Logger != nil {
			cfg.Logger.Warn("Invalid V2RS_ALLOW_INSECURE_HTTP value, ignoring", "value", value)
		}
	}

	if value := strings.TrimSpace(os.Getenv("V2RS_PANEL_STATIC_DIR")); value != "" {
		cfg.Panel.StaticDir = value
	}

	if value := strings.TrimSpace(os.Getenv("V2RS_TRUSTED_PROXIES")); value != "" {
		cfg.Panel.TrustedProxies = splitCSV(value)
	}

	if value := strings.TrimSpace(os.Getenv("DATABASE_URL")); value != "" {
		cfg.Database.URL = value
	}

	if value := strings.TrimSpace(os.Getenv("REDIS_HOST")); value != "" {
		cfg.Redis.Host = value
	}
	if value := strings.TrimSpace(os.Getenv("REDIS_PORT")); value != "" {
		if parsed, err := strconv.Atoi(value); err == nil && parsed > 0 {
			cfg.Redis.Port = parsed
		}
	}

	if value := strings.TrimSpace(os.Getenv("FRONT_END_DOMAIN")); value != "" {
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
