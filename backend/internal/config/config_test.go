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

func TestLoadConfigFailFastValidations(t *testing.T) {
	t.Setenv("METRICS_USER", "")
	t.Setenv("METRICS_PASS", "")
	t.Setenv("DATABASE_URL", "")

	_, err := LoadConfig()
	if err == nil {
		t.Fatalf("expected error when METRICS_USER is empty")
	}

	t.Setenv("METRICS_USER", "admin")
	_, err = LoadConfig()
	if err == nil {
		t.Fatalf("expected error when METRICS_PASS is empty")
	}

	t.Setenv("METRICS_PASS", "pass")
	_, err = LoadConfig()
	if err == nil || err.Error() != "DATABASE_URL is not set" {
		t.Fatalf("expected 'DATABASE_URL is not set' error, got %v", err)
	}

	t.Setenv("DATABASE_URL", "postgresql://postgres:postgres@exodus-db:5432/postgres")
	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("unexpected error loading valid config: %v", err)
	}
	if cfg.Metrics.User != "admin" || cfg.Metrics.Pass != "pass" {
		t.Fatalf("unexpected metrics config: user=%q pass=%q", cfg.Metrics.User, cfg.Metrics.Pass)
	}
}

func TestNormalizePanelConfigKeepsDocsPathsWithoutTrailingSlash(t *testing.T) {
	cfg := defaultConfig
	cfg.Docs.SwaggerPath = "/docs/"
	cfg.Docs.ScalarPath = "scalar/"

	normalizePanelConfig(&cfg)

	if cfg.Docs.SwaggerPath != "/docs" {
		t.Fatalf("swagger path got %q, want /docs", cfg.Docs.SwaggerPath)
	}
	if cfg.Docs.ScalarPath != "/scalar" {
		t.Fatalf("scalar path got %q, want /scalar", cfg.Docs.ScalarPath)
	}
}
