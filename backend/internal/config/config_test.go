package config

import (
	"net/url"
	"testing"
)

func TestPostgresSocketDatabaseURL(t *testing.T) {
	t.Setenv("POSTGRES_USER", "postgres")
	t.Setenv("POSTGRES_PASSWORD", "p@ss/word")
	t.Setenv("POSTGRES_DB", "exodus")

	dsn := postgresSocketDatabaseURL("/var/run/postgresql")
	parsed, err := url.Parse(dsn)
	if err != nil {
		t.Fatalf("parse dsn: %v", err)
	}

	if parsed.Scheme != "postgresql" {
		t.Fatalf("unexpected scheme: %s", parsed.Scheme)
	}
	if parsed.User.Username() != "postgres" {
		t.Fatalf("unexpected user: %s", parsed.User.Username())
	}
	if password, _ := parsed.User.Password(); password != "p@ss/word" {
		t.Fatalf("unexpected password: %s", password)
	}
	if parsed.Path != "/exodus" {
		t.Fatalf("unexpected database path: %s", parsed.Path)
	}
	if parsed.Query().Get("host") != "/var/run/postgresql" {
		t.Fatalf("unexpected socket host: %s", parsed.Query().Get("host"))
	}
}

func TestDatabaseSocketOverridesDatabaseURL(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgresql://postgres:postgres@exodus-db:5432/postgres")
	t.Setenv("DATABASE_SOCKET", "/var/run/postgresql")
	t.Setenv("POSTGRES_SOCKET", "")
	t.Setenv("POSTGRES_USER", "postgres")
	t.Setenv("POSTGRES_PASSWORD", "postgres")
	t.Setenv("POSTGRES_DB", "postgres")

	cfg := defaultConfig
	applyEnvOverrides(&cfg)

	if cfg.Database.Socket != "/var/run/postgresql" {
		t.Fatalf("unexpected socket: %s", cfg.Database.Socket)
	}
	parsed, err := url.Parse(cfg.Database.URL)
	if err != nil {
		t.Fatalf("parse dsn: %v", err)
	}
	if parsed.Query().Get("host") != "/var/run/postgresql" {
		t.Fatalf("DATABASE_SOCKET did not override DATABASE_URL: %s", cfg.Database.URL)
	}
}

func TestDatabaseURLFallsBackToPostgresEnvWhenSocketDisabled(t *testing.T) {
	t.Setenv("DATABASE_URL", "")
	t.Setenv("DATABASE_SOCKET", "")
	t.Setenv("POSTGRES_SOCKET", "")
	t.Setenv("POSTGRES_USER", "exodus")
	t.Setenv("POSTGRES_PASSWORD", "s3cret")
	t.Setenv("POSTGRES_DB", "panel")
	t.Setenv("DATABASE_HOST", "postgres.local")
	t.Setenv("DATABASE_PORT", "6543")

	cfg := defaultConfig
	applyEnvOverrides(&cfg)

	if cfg.Database.Socket != "" {
		t.Fatalf("unexpected socket: %s", cfg.Database.Socket)
	}
	parsed, err := url.Parse(cfg.Database.URL)
	if err != nil {
		t.Fatalf("parse dsn: %v", err)
	}
	if parsed.Host != "postgres.local:6543" {
		t.Fatalf("unexpected host: %s", parsed.Host)
	}
	if parsed.User.Username() != "exodus" {
		t.Fatalf("unexpected user: %s", parsed.User.Username())
	}
	if password, _ := parsed.User.Password(); password != "s3cret" {
		t.Fatalf("unexpected password: %s", password)
	}
	if parsed.Path != "/panel" {
		t.Fatalf("unexpected database path: %s", parsed.Path)
	}
}
