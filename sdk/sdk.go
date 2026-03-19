package sdk

import (
	"fmt"
	"time"
)

const (
	defaultDialTimeout    = 5 * time.Second
	defaultRequestTimeout = 5 * time.Second
)

// API is a unified core SDK entrypoint (similar in spirit to cerberus XtlsApi).
type API struct {
	Stats StatsService
}

// New creates a sing-box core SDK instance.
func New(cfg Config) (*API, error) {
	if cfg.Address == "" {
		return nil, fmt.Errorf("core API address is required")
	}
	if cfg.Port < 1 || cfg.Port > 65535 {
		return nil, fmt.Errorf("invalid core API port: %d", cfg.Port)
	}
	if cfg.CoreType != "" && cfg.CoreType != "singbox" {
		return nil, fmt.Errorf("unsupported core type: %s (only singbox is supported)", cfg.CoreType)
	}
	if cfg.DialTimeout <= 0 {
		cfg.DialTimeout = defaultDialTimeout
	}
	if cfg.RequestTimeout <= 0 {
		cfg.RequestTimeout = defaultRequestTimeout
	}

	stats, err := newSingboxStatsService(cfg)
	if err != nil {
		return nil, err
	}

	return &API{Stats: stats}, nil
}

// Close closes all API services.
func (a *API) Close() error {
	if a == nil || a.Stats == nil {
		return nil
	}
	return a.Stats.Close()
}
