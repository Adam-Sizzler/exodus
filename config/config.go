package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"exodus-node/logger"
)

// NodeConfig holds the configuration settings for the node.
type NodeConfig struct {
	Log      LogConfig
	Exodus   ExodusConfig
	Logger   *logger.Logger
}

type LogConfig struct {
	LogLevel string
	LogMode  string
}

type ExodusConfig struct {
	GrpcAddress string
	GrpcPort    int
	GrpcPath    string
	MTLSConfig  *MTLSConfig
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
	DefaultNodeLogLevel      = "warn"
	DefaultNodeLogMode       = "inclusive"
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
		LogMode:  DefaultNodeLogMode,
	},
	Exodus: ExodusConfig{
		GrpcAddress: DefaultNodeGRPCAddress,
		GrpcPort:    DefaultNodeGRPCPort,
		GrpcPath:    DefaultNodeGRPCPath,
	},
}

// LoadNodeConfig loads node configuration from environment variables only.
func LoadNodeConfig() (NodeConfig, error) {
	cfg := defaultConfig

	if value := firstEnv("LOG_LEVEL", "EXODUS_LOG_LEVEL"); value != "" {
		cfg.Log.LogLevel = value
	}
	if value := firstEnv("LOG_MODE", "EXODUS_LOG_MODE"); value != "" {
		cfg.Log.LogMode = value
	}

	if value := firstEnv("NODE_GRPC_ADDRESS", "LISTEN_GRPC_ADDRESS"); value != "" {
		cfg.Exodus.GrpcAddress = value
	}

	if value := firstEnv("NODE_GRPC_PORT", "LISTEN_GRPC_PORT", "NODE_PORT"); value != "" {
		port, err := strconv.Atoi(value)
		if err != nil || port < 1 || port > 65535 {
			return cfg, fmt.Errorf("invalid NODE_GRPC_PORT/NODE_PORT value: %q", value)
		}
		cfg.Exodus.GrpcPort = port
	}

	if value := firstEnv("NODE_GRPC_PATH", "GRPC_PATH"); value != "" {
		cfg.Exodus.GrpcPath = strings.Trim(value, "/")
	}

	nodePayload, err := ParseNodePayloadFromSecret()
	if err != nil {
		return cfg, err
	}
	cfg.Exodus.MTLSConfig = &MTLSConfig{
		Cert:   nodePayload.NodeCertPem,
		Key:    nodePayload.NodeKeyPem,
		CACert: nodePayload.CaCertPem,
	}

	cfg.Logger, err = logger.NewLoggerWithValidation(cfg.Log.LogLevel, cfg.Log.LogMode, "UTC", os.Stderr)
	if err != nil {
		return cfg, fmt.Errorf("failed to initialize logger: %w", err)
	}
	cfg.Logger.Info("Node configuration validated", "address", cfg.Exodus.GrpcAddress, "port", cfg.Exodus.GrpcPort, "mtls_enabled", cfg.Exodus.MTLSConfig != nil)
	return cfg, nil
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
