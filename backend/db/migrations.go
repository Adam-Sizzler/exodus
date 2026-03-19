package db

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"embed"
	"encoding/hex"
	"errors"
	"fmt"
	"io/fs"
	"sort"
	"strings"

	"cerberus/backend/config"
)

//go:embed prisma/migrations/*/migration.sql
var migrationsFS embed.FS

const migrationsAdvisoryLockKey int64 = 2203092601

func ApplyMigrations(ctx context.Context, dbConn *sql.DB, cfg *config.BackendConfig) error {
	if dbConn == nil {
		return fmt.Errorf("database connection is nil")
	}

	conn, err := dbConn.Conn(ctx)
	if err != nil {
		return fmt.Errorf("acquire migration connection: %w", err)
	}
	defer conn.Close()

	if _, err := conn.ExecContext(ctx, `SELECT pg_advisory_lock($1)`, migrationsAdvisoryLockKey); err != nil {
		return fmt.Errorf("acquire migration advisory lock: %w", err)
	}
	defer func() {
		if _, unlockErr := conn.ExecContext(context.Background(), `SELECT pg_advisory_unlock($1)`, migrationsAdvisoryLockKey); unlockErr != nil {
			cfg.Logger.Warn("Failed to release migration advisory lock", "error", unlockErr)
		}
	}()

	if _, err := conn.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			migration_name TEXT PRIMARY KEY,
			checksum TEXT NOT NULL,
			applied_at TIMESTAMPTZ NOT NULL DEFAULT now()
		)
	`); err != nil {
		return fmt.Errorf("create schema_migrations table: %w", err)
	}

	dirs, err := fs.ReadDir(migrationsFS, "prisma/migrations")
	if err != nil {
		return fmt.Errorf("read migrations directory: %w", err)
	}

	var names []string
	for _, dir := range dirs {
		if !dir.IsDir() {
			continue
		}
		names = append(names, dir.Name())
	}
	sort.Strings(names)

	for _, name := range names {
		sqlPath := fmt.Sprintf("prisma/migrations/%s/migration.sql", name)
		sqlBytes, err := migrationsFS.ReadFile(sqlPath)
		if err != nil {
			return fmt.Errorf("read migration %s: %w", name, err)
		}
		sqlText := strings.TrimSpace(string(sqlBytes))
		if sqlText == "" {
			continue
		}

		checksum := sha256.Sum256(sqlBytes)
		checksumHex := hex.EncodeToString(checksum[:])

		var existing string
		err = conn.QueryRowContext(ctx, `SELECT checksum FROM schema_migrations WHERE migration_name = $1`, name).Scan(&existing)
		switch {
		case err == nil:
			if existing != checksumHex {
				return fmt.Errorf("migration %s checksum mismatch: stored=%s current=%s", name, existing, checksumHex)
			}
			continue
		case errors.Is(err, sql.ErrNoRows):
			// apply
		default:
			if err != nil {
				return fmt.Errorf("check migration %s: %w", name, err)
			}
		}

		cfg.Logger.Info("Applying migration", "name", name)
		tx, err := conn.BeginTx(ctx, nil)
		if err != nil {
			return fmt.Errorf("begin migration %s: %w", name, err)
		}
		if _, err := tx.ExecContext(ctx, sqlText); err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("apply migration %s: %w", name, err)
		}
		if _, err := tx.ExecContext(ctx, `INSERT INTO schema_migrations (migration_name, checksum) VALUES ($1, $2)`, name, checksumHex); err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("record migration %s: %w", name, err)
		}
		if err := tx.Commit(); err != nil {
			return fmt.Errorf("commit migration %s: %w", name, err)
		}
	}

	return nil
}
