package config

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"time"

	"v2ray-stat/logger"

	"gopkg.in/yaml.v3"
)

// NodeConfig holds the configuration settings for the node.
type NodeConfig struct {
	Log            LogConfig       `yaml:"log_config"`
	V2RS           V2RSConfig      `yaml:"node"`
	TZ             string          `yaml:"TZ"`
	Features       map[string]bool `yaml:"features"`
	Core           CoreConfig      `yaml:"core"`
	Paths          PathsConfig     `yaml:"paths"`
	ServiceManager string          `yaml:"service_manager"`
	Logger         *logger.Logger
}

type LogConfig struct {
	LogLevel string `yaml:"level"`
	LogMode  string `yaml:"mode"`
}

type V2RSConfig struct {
	GrpcAddress string      `yaml:"listen_grpc_address"`
	GrpcPort    int         `yaml:"listen_grpc_port"`
	GrpcPath    string      `yaml:"path"`
	MTLSConfig  *MTLSConfig `yaml:"mtls"`
}

type CoreConfig struct {
	Type           string      `yaml:"type"`
	ApiGrpcAddress string      `yaml:"api_grpc_address"`
	ApiGrpcPort    int         `yaml:"api_grpc_port"`
	Dir            string      `yaml:"dir"`
	Config         string      `yaml:"config"`
	AccessLog      string      `yaml:"access_log"`
	AccessLogRegex string      `yaml:"access_log_regex"`
	LogSource      string      `yaml:"log_source"`
	LogServiceName string      `yaml:"log_service_name"`
	LogMap         LogFieldMap `yaml:"access_log_map"`
}

// PathsConfig holds paths and logging settings.
type PathsConfig struct {
	F2BLog       string `yaml:"f2b_log"`
	F2BBannedLog string `yaml:"f2b_banned_log"`
	HAProxyAuth  string `yaml:"haproxy_auth"`
}

type MTLSConfig struct {
	Cert   string `yaml:"cert"`
	Key    string `yaml:"key"`
	CACert string `yaml:"ca_cert"`
}

type LogFieldMap struct {
	User   int `yaml:"user"`
	IP     int `yaml:"ip"`
	Domain int `yaml:"domain"`
}

var defaultConfig = NodeConfig{
	Log: LogConfig{
		LogLevel: "none",
		LogMode:  "inclusive",
	},
	TZ: "",
	V2RS: V2RSConfig{
		GrpcAddress: "127.0.0.1",
		GrpcPort:    9253,
		GrpcPath:    "",
	},
	Features: make(map[string]bool),
	Core: CoreConfig{
		Type:           "xray",
		ApiGrpcAddress: "127.0.0.1",
		ApiGrpcPort:    10812,
		Dir:            "/app/xray/",
		Config:         "/app/xray/config.json",
		LogSource:      "file",
		AccessLog:      "/app/logs/access.log",
		AccessLogRegex: `from (?:[\w]+:)?([\d\.]+):\d+ accepted (?:tcp|udp):\[?([\w\.\-:]+)\]?:\d+ \[[^\]]+\] email: (\S+)`,
		LogMap: LogFieldMap{
			IP:     1,
			Domain: 2,
			User:   3,
		},
	},
	Paths: PathsConfig{
		F2BLog:       "/var/log/v2ray-stat.log",
		F2BBannedLog: "/var/log/v2ray-stat-banned.log",
		HAProxyAuth:  "/etc/haproxy/data/users.csv",
	},
}

func LoadNodeConfig(configFile string) (NodeConfig, error) {
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
	cfg.Logger, _ = logger.NewLoggerWithValidation("info", "inclusive", cfg.TZ, os.Stderr)
	cfg.Logger.Debug("Read configuration file", "file", configFile, "size", len(data))

	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return cfg, fmt.Errorf("error parsing YAML configuration from %s: %v", configFile, err)
	}

	cfg.Logger, err = logger.NewLoggerWithValidation(cfg.Log.LogLevel, cfg.Log.LogMode, cfg.TZ, os.Stderr)
	if err != nil {
		return cfg, fmt.Errorf("failed to initialize logger: %v", err)
	}

	if cfg.Core.Type != "xray" && cfg.Core.Type != "singbox" {
		cfg.Logger.Warn("Invalid v2ray-stat.type, using default", "type", cfg.Core.Type, "default", defaultConfig.Core.Type)
		cfg.Core.Type = defaultConfig.Core.Type
	}

	if cfg.V2RS.GrpcPort != 0 {
		if cfg.V2RS.GrpcPort < 1 || cfg.V2RS.GrpcPort > 65535 {
			cfg.Logger.Warn("Invalid v2ray-stat.grpc_port, using default", "port", cfg.V2RS.GrpcPort, "default", defaultConfig.V2RS.GrpcPort)
			cfg.V2RS.GrpcPort = defaultConfig.V2RS.GrpcPort
		}
	}

	if cfg.Core.AccessLogRegex != "" {
		if _, err := regexp.Compile(cfg.Core.AccessLogRegex); err != nil {
			cfg.Logger.Warn("Invalid core.access_log_regex, using default", "regex", cfg.Core.AccessLogRegex, "default", defaultConfig.Core.AccessLogRegex)
			cfg.Core.AccessLogRegex = defaultConfig.Core.AccessLogRegex
		}
	}

	if cfg.Core.Config == "" {
		cfg.Logger.Warn("Core config path is empty, using default", "default", defaultConfig.Core.Config)
		cfg.Core.Config = defaultConfig.Core.Config
	} else if _, err := os.Stat(cfg.Core.Config); os.IsNotExist(err) {
		cfg.Logger.Warn("Core config file not found, using default", "file", cfg.Core.Config, "default", defaultConfig.Core.Config)
		cfg.Core.Config = defaultConfig.Core.Config
	}

	isDefaultAddr := cfg.Core.ApiGrpcAddress == "" || cfg.Core.ApiGrpcAddress == defaultConfig.Core.ApiGrpcAddress
	isDefaultPort := cfg.Core.ApiGrpcPort == 0 || cfg.Core.ApiGrpcPort == defaultConfig.Core.ApiGrpcPort

	if isDefaultAddr && isDefaultPort {
		cfg.Logger.Debug("API address/port are default, attempting auto-detection from core config", "file", cfg.Core.Config)

		detAddr, detPortStr := detectApiPort(cfg.Core.Config)
		port, err := strconv.Atoi(detPortStr)
		if err != nil {
			cfg.Core.ApiGrpcPort = port
			if detAddr != "" && detAddr != "0.0.0.0" {
				cfg.Core.ApiGrpcAddress = detAddr
			}
			cfg.Logger.Info("Auto-detected API settings", "address", cfg.Core.ApiGrpcAddress, "port", cfg.Core.ApiGrpcPort)
		}
	} else {
		cfg.Logger.Info("Using API settings from config.yml (override active)", "address", cfg.Core.ApiGrpcAddress, "port", cfg.Core.ApiGrpcPort)
	}

	if cfg.Core.ApiGrpcPort != 0 {
		if cfg.Core.ApiGrpcPort < 1 || cfg.Core.ApiGrpcPort > 65535 {
			cfg.Logger.Warn("Invalid API port, falling back to default", "port", cfg.Core.ApiGrpcPort, "default", defaultConfig.Core.ApiGrpcPort)
			cfg.Core.ApiGrpcPort = defaultConfig.Core.ApiGrpcPort
		}
	}

	if cfg.Core.ApiGrpcPort < 1 || cfg.Core.ApiGrpcPort > 65535 {
		cfg.Logger.Warn("Invalid core.api_grpc_port, using fallback default", "port", cfg.Core.ApiGrpcPort, "default", defaultConfig.Core.ApiGrpcPort)
		cfg.Core.ApiGrpcPort = defaultConfig.Core.ApiGrpcPort
	}

	if cfg.Core.ApiGrpcAddress == "" {
		cfg.Core.ApiGrpcAddress = defaultConfig.Core.ApiGrpcAddress
	}

	if cfg.V2RS.MTLSConfig != nil {
		if cfg.V2RS.GrpcAddress != "127.0.0.1" && cfg.V2RS.GrpcAddress != "0.0.0.0" && cfg.V2RS.GrpcAddress != "localhost" {
			if cfg.V2RS.MTLSConfig.Cert == "" || cfg.V2RS.MTLSConfig.Key == "" || cfg.V2RS.MTLSConfig.CACert == "" {
				cfg.Logger.Error("Incomplete mTLS configuration for non-localhost address", "address", cfg.V2RS.GrpcAddress)
				return cfg, fmt.Errorf("incomplete mTLS configuration for non-localhost address")
			}
			for _, file := range []string{cfg.V2RS.MTLSConfig.Cert, cfg.V2RS.MTLSConfig.Key, cfg.V2RS.MTLSConfig.CACert} {
				if _, err := os.Stat(file); os.IsNotExist(err) {
					cfg.Logger.Error("mTLS certificate file not found for non-localhost address", "file", file, "address", cfg.V2RS.GrpcAddress)
					return cfg, fmt.Errorf("mTLS certificate file not found: %s", file)
				}
			}
		}
	}

	if cfg.TZ != "" {
		if _, err := time.LoadLocation(cfg.TZ); err != nil {
			cfg.Logger.Warn("Invalid timezone value, using default", "timezone", cfg.TZ)
			cfg.TZ = defaultConfig.TZ
		}
	}

	if cfg.Core.LogSource != "file" && cfg.Core.LogSource != "journal" {
		cfg.Logger.Warn("Invalid core.log_source, using default", "source", cfg.Core.LogSource, "default", "file")
		cfg.Core.LogSource = "file"
	}
	if cfg.Core.LogSource == "journal" && cfg.Core.LogServiceName == "" {
		cfg.Logger.Warn("LogServiceName empty for journal source, falling back to file")
		cfg.Core.LogSource = "file"
	}

	if cfg.Features == nil {
		cfg.Features = make(map[string]bool)
	}

	if cfg.ServiceManager == "" {
		cfg.ServiceManager = "systemd"
	}

	cfg.Logger.Info("Node configuration validated", "address", cfg.V2RS.GrpcAddress, "port", cfg.V2RS.GrpcPort, "mtls_enabled", cfg.V2RS.MTLSConfig != nil)
	return cfg, nil
}
