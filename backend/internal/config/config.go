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

	AppVersion string
}

type MTLSConfig struct {
	Cert   string
	Key    string
	CACert string
}

func Load() (Config, error) {
	pathPrefix := normalizePathPrefix(firstEnv("SUB_GRPC_PATH", "NODE_GRPC_PATH", "GRPC_PATH"))
	grpcToken := strings.TrimSpace(firstEnv("SUB_GRPC_TOKEN", "NODE_GRPC_TOKEN", "GRPC_TOKEN"))

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
	}

	grpcPort, err := parsePort(firstEnv("SUB_GRPC_PORT", "NODE_GRPC_PORT", "NODE_PORT"), DefaultGRPCPort)
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
			return cfg, fmt.Errorf("SUB_GRPC_TOKEN must be at least 16 characters")
		}
		if len(cfg.GRPCToken) > 512 {
			return cfg, fmt.Errorf("SUB_GRPC_TOKEN must be less than 512 characters")
		}
	}
	if cfg.RequireGRPCToken && cfg.GRPCToken == "" {
		return cfg, fmt.Errorf("SUB_GRPC_TOKEN is required when SUB_SECRET_KEY is not provided")
	}

	seed := strings.TrimSpace(os.Getenv("SUB_SECRET_KEY"))
	if seed == "" {
		seed = strings.TrimSpace(os.Getenv("SECRET_KEY"))
	}
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
		return 0, fmt.Errorf("invalid SUB_GRPC_PORT value: %q", raw)
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
