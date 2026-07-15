package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// NodeConfig holds the configuration settings for the node.
type NodeConfig struct {
	Log             LogConfig
	Exodus          ExodusConfig
	CoreAPIGRPCPort int
	Logger          *Logger
}

type LogConfig struct {
	LogLevel string
}

type ExodusConfig struct {
	GrpcAddress      string
	GrpcPort         int
	GrpcPath         string
	GRPCToken        string
	RequireGRPCToken bool
	MTLSConfig       *MTLSConfig
}

const (
	FixedSingboxDir          = "/app/singbox/"
	FixedSingboxConfigPath   = "/app/singbox/config.json"
	FixedCoreType            = "singbox"
	FixedCoreAPIAddress      = "127.0.0.1"
	FixedCoreAPIGRPCPort     = 10085
	DefaultNodeGRPCAddress   = "0.0.0.0"
	DefaultNodeGRPCPort      = 2222
	DefaultNodeGRPCPath      = ""
	DefaultNodeLogLevel      = ""
	DefaultNodeTLSServerName = "internal.exodus.local"
)

type MTLSConfig struct {
	Cert   string
	Key    string
	CACert string
}

var defaultConfig = NodeConfig{
	Log: LogConfig{
		LogLevel: DefaultNodeLogLevel,
	},
	Exodus: ExodusConfig{
		GrpcAddress: DefaultNodeGRPCAddress,
		GrpcPort:    DefaultNodeGRPCPort,
		GrpcPath:    DefaultNodeGRPCPath,
	},
	CoreAPIGRPCPort: FixedCoreAPIGRPCPort,
}

// LoadNodeConfig loads node configuration from environment variables only.
func LoadNodeConfig() (NodeConfig, error) {
	cfg := defaultConfig

	cfg.Log.LogLevel = ResolveExodusLogLevel(cfg.Log.LogLevel)

	if value := os.Getenv("SINGBOX_API_PORT"); value != "" {
		port, err := strconv.Atoi(strings.TrimSpace(value))
		if err == nil && port >= 1 && port <= 65535 {
			cfg.CoreAPIGRPCPort = port
		}
	}

	if value := firstEnv("NODE_ADDRESS"); value != "" {
		cfg.Exodus.GrpcAddress = value
	}

	if value := firstEnv("NODE_PORT"); value != "" {
		port, err := strconv.Atoi(value)
		if err != nil || port < 1 || port > 65535 {
			return cfg, fmt.Errorf("invalid NODE_PORT value: %q", value)
		}
		cfg.Exodus.GrpcPort = port
	}

	if value := firstEnv("NODE_PATH"); value != "" {
		cfg.Exodus.GrpcPath = strings.Trim(value, "/")
	}
	if value := firstEnv("NODE_TOKEN"); value != "" {
		cfg.Exodus.GRPCToken = value
	}

	nodePayload, err := ParseNodePayloadFromSecret()
	if err == nil {
		cfg.Exodus.MTLSConfig = &MTLSConfig{
			Cert:   nodePayload.NodeCertPem,
			Key:    nodePayload.NodeKeyPem,
			CACert: nodePayload.CaCertPem,
		}
	} else if !errors.Is(err, ErrSecretKeyNotSet) {
		return cfg, err
	}

	cfg.Exodus.RequireGRPCToken = cfg.Exodus.MTLSConfig == nil
	if cfg.Exodus.GRPCToken != "" {
		if len(cfg.Exodus.GRPCToken) < 16 {
			return cfg, fmt.Errorf("NODE_TOKEN must be at least 16 characters")
		}
		if len(cfg.Exodus.GRPCToken) > 512 {
			return cfg, fmt.Errorf("NODE_TOKEN must be less than 512 characters")
		}
	}
	if cfg.Exodus.RequireGRPCToken && cfg.Exodus.GRPCToken == "" {
		return cfg, fmt.Errorf("NODE_TOKEN is required when SECRET_KEY is not provided")
	}

	cfg.Logger = NewExodusLogger(os.Stderr, cfg.Log.LogLevel)
	cfg.Logger.WithContext("ConfigService").Debug(
		"Node configuration validated",
		"address", cfg.Exodus.GrpcAddress,
		"port", cfg.Exodus.GrpcPort,
		"path", cfg.Exodus.GrpcPath,
		"mtls_enabled", cfg.Exodus.MTLSConfig != nil,
		"token_required", cfg.Exodus.RequireGRPCToken,
	)
	return cfg, nil
}

func (cfg *NodeConfig) LoggerFor(context string) *Logger {
	if cfg == nil || cfg.Logger == nil {
		return nil
	}
	return cfg.Logger.WithContext(context)
}

func firstEnv(keys ...string) string {
	for _, key := range keys {
		value, ok := lookupEnvTrimmed(key)
		if ok {
			return value
		}
	}
	return ""
}

func lookupEnvTrimmed(key string) (string, bool) {
	value, ok := os.LookupEnv(key)
	if !ok {
		return "", false
	}
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return "", false
	}
	return trimmed, true
}
