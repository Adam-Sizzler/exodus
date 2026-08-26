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

	db.SetMaxOpenConns(50)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(15 * time.Minute)
	db.SetConnMaxIdleTime(5 * time.Minute)

	for {
		pingCtx, cancelPing := context.WithTimeout(context.Background(), 5*time.Second)
		err := db.PingContext(pingCtx)
		cancelPing()
		if err == nil {
			break
		}
		cfg.Logger.Warn("Database not ready, retrying in 5s", "error", err)
		time.Sleep(5 * time.Second)
	}

	execCtx, cancelExec := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelExec()
	if _, err := db.ExecContext(execCtx, `CREATE EXTENSION IF NOT EXISTS "pgcrypto"`); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("ensure pgcrypto extension: %w", err)
	}

	initCtx, cancelInit := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancelInit()

	fmt.Println("Migrating database...")
	if err := ApplyMigrations(initCtx, db, cfg); err != nil {
		_ = db.Close()
		return nil, err
	}
	fmt.Println("Migrations deployed successfully!")
	fmt.Println("Seeding database...")

	return db, nil
}

// InitDatabase initializes the database and returns a ready connection.
func InitDatabase(cfg *config.BackendConfig) (*sql.DB, error) {
	return OpenAndInitDB(cfg)
}
