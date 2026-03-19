package config

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"strings"
)

const (
	SubpageDefaultConfigUUID = "00000000-0000-0000-0000-000000000000"
	DefaultPort              = "3010"
)

type Config struct {
	AppPort                         string
	CerberusPanelURL               string
	CerberusAPIToken               string
	SubpageConfigUUID               string
	CustomSubPrefix                 string
	CaddyAuthAPIToken               string
	CloudflareZeroTrustClientID     string
	CloudflareZeroTrustClientSecret string
	SessionSecret                   string
	NodeEnv                         string
}

func Load() (Config, error) {
	cfg := Config{
		AppPort:                         getEnvOrDefault("APP_PORT_SUB", DefaultPort),
		CerberusPanelURL:               strings.TrimRight(os.Getenv("CERBERUS_PANEL_URL"), "/"),
		CerberusAPIToken:               os.Getenv("CERBERUS_API_TOKEN"),
		SubpageConfigUUID:               getEnvOrDefault("SUBPAGE_CONFIG_UUID", SubpageDefaultConfigUUID),
		CustomSubPrefix:                 strings.Trim(os.Getenv("CUSTOM_SUB_PREFIX"), "/"),
		CaddyAuthAPIToken:               os.Getenv("CADDY_AUTH_API_TOKEN"),
		CloudflareZeroTrustClientID:     os.Getenv("CLOUDFLARE_ZERO_TRUST_CLIENT_ID"),
		CloudflareZeroTrustClientSecret: os.Getenv("CLOUDFLARE_ZERO_TRUST_CLIENT_SECRET"),
		NodeEnv:                         strings.TrimSpace(os.Getenv("NODE_ENV")),
	}

	if cfg.CerberusPanelURL == "" {
		return cfg, fmt.Errorf("CERBERUS_PANEL_URL is required")
	}

	if !strings.HasPrefix(cfg.CerberusPanelURL, "http://") &&
		!strings.HasPrefix(cfg.CerberusPanelURL, "https://") {
		return cfg, fmt.Errorf("CERBERUS_PANEL_URL must start with http:// or https://")
	}

	if strings.TrimSpace(cfg.CerberusAPIToken) == "" {
		return cfg, fmt.Errorf("CERBERUS_API_TOKEN is required")
	}

	cfg.SessionSecret = deriveSessionSecret(cfg.CerberusAPIToken)

	return cfg, nil
}

func (c Config) IsDevelopment() bool {
	return strings.EqualFold(c.NodeEnv, "development")
}

func DisplayPrefix(prefix string) string {
	if prefix == "" {
		return "not set"
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

func deriveSessionSecret(apiToken string) string {
	hash := sha256.Sum256([]byte(strings.TrimSpace(apiToken)))
	return hex.EncodeToString(hash[:])
}
