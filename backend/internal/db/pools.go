package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"exodus/internal/config"
)

// Pools holds separate connection pools for interactive HTTP requests and background workers.
type Pools struct {
	Interactive *sql.DB
	Background  *sql.DB
}

// NewPools configures the interactive pool and opens an isolated background pool.
func NewPools(ctx context.Context, interactive *sql.DB, cfg *config.BackendConfig) (*Pools, error) {
	interactive.SetMaxOpenConns(32)
	interactive.SetMaxIdleConns(16)
	interactive.SetConnMaxLifetime(30 * time.Minute)
	interactive.SetConnMaxIdleTime(5 * time.Minute)

	bg, err := sql.Open("pgx", cfg.Database.URL)
	if err != nil {
		return nil, fmt.Errorf("open background pool: %w", err)
	}

	pingCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	if err := bg.PingContext(pingCtx); err != nil {
		_ = bg.Close()
		return nil, fmt.Errorf("ping background pool: %w", err)
	}

	bg.SetMaxOpenConns(8)
	bg.SetMaxIdleConns(4)
	bg.SetConnMaxLifetime(30 * time.Minute)
	bg.SetConnMaxIdleTime(5 * time.Minute)

	return &Pools{
		Interactive: interactive,
		Background:  bg,
	}, nil
}

// Close closes both database connection pools.
func (p *Pools) Close() error {
	var errs []error
	if p.Interactive != nil {
		if err := p.Interactive.Close(); err != nil {
			errs = append(errs, fmt.Errorf("close interactive pool: %w", err))
		}
	}
	if p.Background != nil {
		if err := p.Background.Close(); err != nil {
			errs = append(errs, fmt.Errorf("close background pool: %w", err))
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf("errors closing pools: %v", errs)
	}
	return nil
}
