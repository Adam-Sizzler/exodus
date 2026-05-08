package config

import "testing"

func TestLoadNodeConfigTokenMode(t *testing.T) {
	t.Setenv("SECRET_KEY", "")
	t.Setenv("NODE_GRPC_TOKEN", "1234567890abcdef")
	t.Setenv("NODE_GRPC_PATH", "/node/")
	t.Setenv("LOG_LEVEL", "warn")
	t.Setenv("LOG_MODE", "inclusive")

	cfg, err := LoadNodeConfig()
	if err != nil {
		t.Fatalf("LoadNodeConfig() error = %v", err)
	}
	if cfg.Exodus.MTLSConfig != nil {
		t.Fatalf("MTLSConfig = %#v, want nil in token mode", cfg.Exodus.MTLSConfig)
	}
	if !cfg.Exodus.RequireGRPCToken {
		t.Fatalf("RequireGRPCToken = false, want true")
	}
	if cfg.Exodus.GRPCToken != "1234567890abcdef" {
		t.Fatalf("GRPCToken = %q", cfg.Exodus.GRPCToken)
	}
	if cfg.Exodus.GrpcPath != "node" {
		t.Fatalf("GrpcPath = %q, want node", cfg.Exodus.GrpcPath)
	}
}
