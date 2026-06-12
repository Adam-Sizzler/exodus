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

	"exodus/internal/config"
)

//go:embed prisma/migrations/*/migration.sql
var migrationsFS embed.FS

const migrationsAdvisoryLockKey int64 = 2203092601

var retiredMigrations = map[string]struct{}{
	"20260518013000_drop_config_profile_snippets": {},
}

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
		CREATE EXTENSION IF NOT EXISTS pgcrypto WITH SCHEMA public;

		CREATE TABLE IF NOT EXISTS public.schema_migrations (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			migration_name TEXT NOT NULL UNIQUE,
			checksum TEXT NOT NULL,
			started_at TIMESTAMPTZ NOT NULL DEFAULT now(),
			finished_at TIMESTAMPTZ,
			applied_steps_count INTEGER NOT NULL DEFAULT 1,
			applied_at TIMESTAMPTZ NOT NULL DEFAULT now()
		);
	`); err != nil {
		return fmt.Errorf("initialize schema_migrations table: %w", err)
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

	knownMigrations := make(map[string]struct{}, len(names))
	for _, name := range names {
		knownMigrations[name] = struct{}{}
	}

	appliedRows, err := conn.QueryContext(ctx, `SELECT migration_name FROM public.schema_migrations ORDER BY migration_name`)
	if err != nil {
		return fmt.Errorf("read applied migrations: %w", err)
	}
	var legacyMigrations []string
	for appliedRows.Next() {
		var appliedName string
		if err := appliedRows.Scan(&appliedName); err != nil {
			_ = appliedRows.Close()
			return fmt.Errorf("scan applied migration: %w", err)
		}
		if _, ok := knownMigrations[appliedName]; ok {
			continue
		}
		if _, ok := retiredMigrations[appliedName]; !ok {
			legacyMigrations = append(legacyMigrations, appliedName)
		}
	}
	if err := appliedRows.Close(); err != nil {
		return fmt.Errorf("close applied migrations cursor: %w", err)
	}
	if len(legacyMigrations) > 0 {
		return fmt.Errorf(
			"unsupported legacy migration history detected: %s; this release supports only the %s baseline, use a clean database or prepare the database for the new baseline",
			strings.Join(legacyMigrations, ", "),
			strings.Join(names, ", "),
		)
	}

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
		err = conn.QueryRowContext(ctx, `SELECT checksum FROM public.schema_migrations WHERE migration_name = $1`, name).Scan(&existing)
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
		if _, err := tx.ExecContext(ctx, `
			INSERT INTO public.schema_migrations (
				migration_name, checksum, started_at, finished_at, applied_steps_count, applied_at
			) VALUES ($1, $2, now(), now(), 1, now())
		`, name, checksumHex); err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("record migration %s: %w", name, err)
		}
		if err := tx.Commit(); err != nil {
			return fmt.Errorf("commit migration %s: %w", name, err)
		}
	}

	return nil
}
