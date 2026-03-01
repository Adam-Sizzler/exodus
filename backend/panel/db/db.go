package db

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"v2ray-stat/backend/panel/config"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// OpenAndInitDB opens a PostgreSQL database, applies migrations, and seeds defaults.
func OpenAndInitDB(cfg *config.BackendConfig) (*sql.DB, error) {
	dsn := strings.TrimSpace(cfg.Database.URL)
	if dsn == "" {
		return nil, fmt.Errorf("DATABASE_URL is not set")
	}

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("ping database: %w", err)
	}

	if _, err := db.ExecContext(ctx, `CREATE EXTENSION IF NOT EXISTS "pgcrypto"`); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("ensure pgcrypto extension: %w", err)
	}

	if err := ApplyMigrations(ctx, db, cfg); err != nil {
		_ = db.Close()
		return nil, err
	}

	if err := SeedDefaults(ctx, db, cfg); err != nil {
		_ = db.Close()
		return nil, err
	}

	return db, nil
}

// InitDatabase initializes the database and returns a ready connection.
func InitDatabase(cfg *config.BackendConfig) (*sql.DB, error) {
	return OpenAndInitDB(cfg)
}
