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

	"exodus/backend/config"
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
		CREATE EXTENSION IF NOT EXISTS pgcrypto;

		CREATE TABLE IF NOT EXISTS schema_migrations (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			checksum TEXT NOT NULL,
			finished_at TIMESTAMPTZ,
			migration_name TEXT NOT NULL UNIQUE,
			logs TEXT,
			rolled_back_at TIMESTAMPTZ,
			started_at TIMESTAMPTZ NOT NULL DEFAULT now(),
			applied_steps_count INTEGER NOT NULL DEFAULT 0,
			applied_at TIMESTAMPTZ NOT NULL DEFAULT now()
		);

		ALTER TABLE schema_migrations ADD COLUMN IF NOT EXISTS id UUID;
		ALTER TABLE schema_migrations ADD COLUMN IF NOT EXISTS checksum TEXT;
		ALTER TABLE schema_migrations ADD COLUMN IF NOT EXISTS finished_at TIMESTAMPTZ;
		ALTER TABLE schema_migrations ADD COLUMN IF NOT EXISTS migration_name TEXT;
		ALTER TABLE schema_migrations ADD COLUMN IF NOT EXISTS logs TEXT;
		ALTER TABLE schema_migrations ADD COLUMN IF NOT EXISTS rolled_back_at TIMESTAMPTZ;
		ALTER TABLE schema_migrations ADD COLUMN IF NOT EXISTS started_at TIMESTAMPTZ;
		ALTER TABLE schema_migrations ADD COLUMN IF NOT EXISTS applied_steps_count INTEGER;
		ALTER TABLE schema_migrations ADD COLUMN IF NOT EXISTS applied_at TIMESTAMPTZ;

		UPDATE schema_migrations
		SET id = gen_random_uuid()
		WHERE id IS NULL;

		UPDATE schema_migrations
		SET started_at = COALESCE(started_at, applied_at, now())
		WHERE started_at IS NULL;

		UPDATE schema_migrations
		SET finished_at = COALESCE(finished_at, applied_at, started_at, now())
		WHERE finished_at IS NULL;

		UPDATE schema_migrations
		SET applied_steps_count = COALESCE(applied_steps_count, 1);

		UPDATE schema_migrations
		SET applied_at = COALESCE(applied_at, finished_at, started_at, now())
		WHERE applied_at IS NULL;

		ALTER TABLE schema_migrations ALTER COLUMN id SET NOT NULL;
		ALTER TABLE schema_migrations ALTER COLUMN checksum SET NOT NULL;
		ALTER TABLE schema_migrations ALTER COLUMN migration_name SET NOT NULL;
		ALTER TABLE schema_migrations ALTER COLUMN started_at SET NOT NULL;
		ALTER TABLE schema_migrations ALTER COLUMN applied_steps_count SET NOT NULL;
		ALTER TABLE schema_migrations ALTER COLUMN applied_at SET NOT NULL;

		ALTER TABLE schema_migrations ALTER COLUMN id SET DEFAULT gen_random_uuid();
		ALTER TABLE schema_migrations ALTER COLUMN started_at SET DEFAULT now();
		ALTER TABLE schema_migrations ALTER COLUMN applied_steps_count SET DEFAULT 0;
		ALTER TABLE schema_migrations ALTER COLUMN applied_at SET DEFAULT now();

		DO $$
		BEGIN
			IF EXISTS (
				SELECT 1
				FROM pg_constraint
				WHERE conrelid = 'schema_migrations'::regclass
				  AND conname = 'schema_migrations_pkey'
			) THEN
				ALTER TABLE schema_migrations DROP CONSTRAINT schema_migrations_pkey;
			END IF;
		END
		$$;

		ALTER TABLE schema_migrations ADD CONSTRAINT schema_migrations_pkey PRIMARY KEY (id);
		CREATE UNIQUE INDEX IF NOT EXISTS schema_migrations_migration_name_key ON schema_migrations(migration_name);
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
		if _, err := tx.ExecContext(ctx, `
			INSERT INTO schema_migrations (
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
