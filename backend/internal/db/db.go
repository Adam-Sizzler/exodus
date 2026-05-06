package db

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"exodus/internal/config"

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

	pingCtx, cancelPing := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelPing()
	if err := db.PingContext(pingCtx); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("ping database: %w", err)
	}

	if _, err := db.ExecContext(pingCtx, `CREATE EXTENSION IF NOT EXISTS "pgcrypto"`); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("ensure pgcrypto extension: %w", err)
	}

	initCtx, cancelInit := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancelInit()

	if err := ApplyMigrations(initCtx, db, cfg); err != nil {
		_ = db.Close()
		return nil, err
	}

	if err := SeedDefaults(initCtx, db, cfg); err != nil {
		_ = db.Close()
		return nil, err
	}

	return db, nil
}

// InitDatabase initializes the database and returns a ready connection.
func InitDatabase(cfg *config.BackendConfig) (*sql.DB, error) {
	return OpenAndInitDB(cfg)
}
