package config

import (
	"errors"
	"fmt"
	"os"
	"regexp"
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

	if value := strings.TrimSpace(os.Getenv("NODE_GRPC_ADDRESS")); value != "" {
		cfg.Exodus.GrpcAddress = value
	}

	if value := strings.TrimSpace(os.Getenv("NODE_GRPC_PORT")); value != "" {
		port, err := strconv.Atoi(value)
		if err != nil || port < 1 || port > 65535 {
			return cfg, fmt.Errorf("invalid NODE_GRPC_PORT value: %q", value)
		}
		cfg.Exodus.GrpcPort = port
	}

	if value := strings.TrimSpace(os.Getenv("NODE_GRPC_PATH")); value != "" {
		if err := validateBasePath(value); err != nil {
			return cfg, err
		}
		cfg.Exodus.GrpcPath = normalizeBasePath(value)
	}

	if value := strings.TrimSpace(os.Getenv("NODE_GRPC_TOKEN")); value != "" {
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
			return cfg, fmt.Errorf("NODE_GRPC_TOKEN must be at least 16 characters")
		}
		if len(cfg.Exodus.GRPCToken) > 512 {
			return cfg, fmt.Errorf("NODE_GRPC_TOKEN must be less than 512 characters")
		}
	}
	if cfg.Exodus.RequireGRPCToken && cfg.Exodus.GRPCToken == "" {
		return cfg, fmt.Errorf("NODE_GRPC_TOKEN is required when SECRET_KEY is not provided")
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

var validBasePathRegex = regexp.MustCompile(`^[a-zA-Z0-9_\-\/]*$`)

func validateBasePath(input string) error {
	trimmed := strings.TrimSpace(input)
	if trimmed == "" || trimmed == "/" {
		return nil
	}
	if strings.Contains(trimmed, "..") {
		return fmt.Errorf("invalid NODE_GRPC_PATH %q: directory traversal is not allowed", input)
	}
	if strings.ContainsAny(trimmed, "?#\\ \t\r\n") {
		return fmt.Errorf("invalid NODE_GRPC_PATH %q: contains illegal characters (spaces, query/fragment markers, or backslashes)", input)
	}
	if !validBasePathRegex.MatchString(trimmed) {
		return fmt.Errorf("invalid NODE_GRPC_PATH %q: must contain only alphanumeric characters, hyphens, underscores, and slashes", input)
	}
	return nil
}

// Trimmed returns GrpcPath without trailing slash (e.g. "/node" or "" for root "/").
func (e ExodusConfig) Trimmed() string {
	normalized := normalizeBasePath(e.GrpcPath)
	if normalized == "/" {
		return ""
	}
	return strings.TrimSuffix(normalized, "/")
}

// IsCustom reports whether a custom non-root GrpcPath is configured.
func (e ExodusConfig) IsCustom() bool {
	return e.Trimmed() != ""
}

// WithSlash returns GrpcPath with leading and trailing slashes (e.g. "/node/" or "/" for root).
func (e ExodusConfig) WithSlash() string {
	return normalizeBasePath(e.GrpcPath)
}

func normalizeBasePath(input string) string {
	cleaned := strings.Trim(strings.TrimSpace(input), "/")
	if cleaned == "" {
		return "/"
	}
	return "/" + cleaned + "/"
}
