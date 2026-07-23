package config

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	DefaultPort        = "3010"
	DefaultGRPCAddress = "0.0.0.0"
	DefaultGRPCPort    = 2222
)

type Config struct {
	AppPort                         string
	SubPath                         string
	CaddyAuthAPIToken               string
	CloudflareZeroTrustClientID     string
	CloudflareZeroTrustClientSecret string
	SessionSecret                   string
	NodeEnv                         string

	GRPCAddress      string
	GRPCPort         int
	GRPCPath         string
	GRPCToken        string
	RequireGRPCToken bool
	MTLSConfig       *MTLSConfig

	AppVersion                      string
	MarzbanLegacyDropRevokedSubscriptions bool
}

type MTLSConfig struct {
	Cert   string
	Key    string
	CACert string
}

// EnvError is a schema-like validation item rendered in the same style as
// Exodus subscription-page Zod validation errors.
type EnvError struct {
	Key     string
	Message string
}

// EnvErrors groups all .env validation failures so the logger can print a
// single grouped configuration error block instead of one fatal Go error.
type EnvErrors []EnvError

func (e EnvErrors) Error() string {
	if len(e) == 0 {
		return ".env configuration validation error"
	}
	parts := make([]string, 0, len(e))
	for _, item := range e {
		if strings.TrimSpace(item.Key) == "" {
			parts = append(parts, strings.TrimSpace(item.Message))
			continue
		}
		parts = append(parts, strings.TrimSpace(item.Key)+": "+strings.TrimSpace(item.Message))
	}
	return strings.Join(parts, "; ")
}

func NewEnvError(key, message string) EnvErrors {
	return EnvErrors{{Key: strings.TrimSpace(key), Message: strings.TrimSpace(message)}}
}

// FormatEnvironmentErrors renders configuration failures like Exodus
// subscription-page: a readable block with exact env names and setup hints.
func FormatEnvironmentErrors(err error) string {
	if err == nil {
		return ""
	}

	var envErrs EnvErrors
	if errors.As(err, &envErrs) {
		lines := make([]string, 0, len(envErrs))
		for _, item := range envErrs {
			key := strings.TrimSpace(item.Key)
			msg := strings.TrimSpace(item.Message)
			if key == "" {
				lines = append(lines, "❌ "+msg)
				continue
			}
			lines = append(lines, fmt.Sprintf("❌ %s: %s", key, msg))
		}
		return fmt.Sprintf(`🔧 Environment Configuration Errors:
%s

Please fix your .env file and restart the application.`, strings.Join(lines, "\n"))
	}

	return fmt.Sprintf(`🔧 Environment Configuration Errors:
❌ %s

Please fix your .env file and restart the application.`, err.Error())
}

func Load() (Config, error) {
	pathPrefix := normalizePathPrefix(firstEnv("SUB_GRPC_PATH", "SUB_PATH"))
	grpcToken := strings.TrimSpace(os.Getenv("SUB_GRPC_TOKEN"))

	cfg := Config{
		AppPort:                         getEnvOrDefault("APP_PORT_SUB", DefaultPort),
		SubPath:                         pathPrefix,
		CaddyAuthAPIToken:               os.Getenv("CADDY_AUTH_API_TOKEN"),
		CloudflareZeroTrustClientID:     os.Getenv("CLOUDFLARE_ZERO_TRUST_CLIENT_ID"),
		CloudflareZeroTrustClientSecret: os.Getenv("CLOUDFLARE_ZERO_TRUST_CLIENT_SECRET"),
		NodeEnv:                         strings.TrimSpace(os.Getenv("NODE_ENV")),
		GRPCAddress:                     getEnvOrDefault("SUB_GRPC_ADDRESS", DefaultGRPCAddress),
		GRPCPath:                        pathPrefix,
		GRPCToken:                       grpcToken,
		AppVersion:                      strings.TrimSpace(firstEnv("SUB_APP_VERSION", "APP_VERSION")),
		MarzbanLegacyDropRevokedSubscriptions: parseBool(os.Getenv("MARZBAN_LEGACY_DROP_REVOKED_SUBSCRIPTIONS"), false),
	}

	grpcPort, err := parsePort(os.Getenv("SUB_GRPC_PORT"), DefaultGRPCPort)
	if err != nil {
		return cfg, err
	}
	cfg.GRPCPort = grpcPort

	payload, err := ParseNodePayloadFromSecret()
	if err == nil {
		cfg.MTLSConfig = &MTLSConfig{
			Cert:   payload.NodeCertPem,
			Key:    payload.NodeKeyPem,
			CACert: payload.CaCertPem,
		}
	} else if !errors.Is(err, ErrSecretKeyNotSet) {
		return cfg, err
	}

	cfg.RequireGRPCToken = cfg.MTLSConfig == nil
	if cfg.GRPCToken != "" {
		if len(cfg.GRPCToken) < 16 {
			return cfg, NewEnvError("SUB_GRPC_TOKEN", "Must be at least 16 characters. Dashboard → Subscription → Nodes → Current node → gRPC Token (SUB_GRPC_TOKEN) or Secret Key (SUB_SECRET_KEY).")
		}
		if len(cfg.GRPCToken) > 512 {
			return cfg, NewEnvError("SUB_GRPC_TOKEN", "Must be less than 512 characters.")
		}
	}
	if cfg.RequireGRPCToken && cfg.GRPCToken == "" {
		return cfg, NewEnvError("SUB_GRPC_TOKEN", "Required when SUB_SECRET_KEY is not provided. Dashboard → Subscription → Nodes → Current node → gRPC Token (SUB_GRPC_TOKEN) or Secret Key (SUB_SECRET_KEY).")
	}

	seed := strings.TrimSpace(os.Getenv("SUB_SECRET_KEY"))
	if seed == "" {
		seed = cfg.GRPCToken
	}
	cfg.SessionSecret = deriveSessionSecret(seed)
	return cfg, nil
}

func (c Config) IsDevelopment() bool {
	return strings.EqualFold(c.NodeEnv, "development")
}

func DisplayPrefix(prefix string) string {
	if prefix == "" {
		return "/"
	}
	return prefix
}

func getEnvOrDefault(key, fallback string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	return value
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

func parsePort(raw string, fallback int) (int, error) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return fallback, nil
	}
	port, err := strconv.Atoi(trimmed)
	if err != nil || port < 1 || port > 65535 {
		return 0, NewEnvError("SUB_GRPC_PORT", fmt.Sprintf("Invalid port value %q. Use a number from 1 to 65535.", raw))
	}
	return port, nil
}

func normalizePathPrefix(raw string) string {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" || trimmed == "/" {
		return ""
	}
	return "/" + strings.Trim(trimmed, "/")
}

func deriveSessionSecret(apiToken string) string {
	hash := sha256.Sum256([]byte(strings.TrimSpace(apiToken)))
	return hex.EncodeToString(hash[:])
}

func parseBool(raw string, fallback bool) bool {
	val := strings.TrimSpace(raw)
	if val == "" {
		return fallback
	}
	b, err := strconv.ParseBool(val)
	if err != nil {
		return fallback
	}
	return b
}
