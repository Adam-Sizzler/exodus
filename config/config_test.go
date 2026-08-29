package config

import "testing"

func TestLoadNodeConfigTokenMode(t *testing.T) {
	t.Setenv("SECRET_KEY", "")
	t.Setenv("NODE_GRPC_TOKEN", "1234567890abcdef")
	t.Setenv("NODE_GRPC_PATH", "/node/")
	t.Setenv("LOG_LEVEL", "debug")

	cfg, err := LoadNodeConfig()
	if err != nil {
		t.Fatalf("LoadNodeConfig() error = %v", err)
	}
	if cfg.Backend.MTLSConfig != nil {
		t.Fatalf("MTLSConfig = %#v, want nil in token mode", cfg.Backend.MTLSConfig)
	}
	if !cfg.Backend.RequireGRPCToken {
		t.Fatalf("RequireGRPCToken = false, want true")
	}
	if cfg.Backend.GRPCToken != "1234567890abcdef" {
		t.Fatalf("GRPCToken = %q", cfg.Backend.GRPCToken)
	}
	if cfg.Backend.Trimmed() != "/node" {
		t.Fatalf("Trimmed() = %q, want /node", cfg.Backend.Trimmed())
	}
	if !cfg.Backend.IsCustom() {
		t.Fatalf("IsCustom() = false, want true")
	}
	if cfg.Backend.WithSlash() != "/node/" {
		t.Fatalf("WithSlash() = %q, want /node/", cfg.Backend.WithSlash())
	}
	if cfg.Log.LogLevel != "debug" {
		t.Fatalf("LogLevel = %q, want debug", cfg.Log.LogLevel)
	}
}

func TestConfigPathMethods(t *testing.T) {
	tests := []struct {
		name          string
		grpcPath      string
		wantTrimmed   string
		wantIsCustom  bool
		wantWithSlash string
	}{
		{
			name:          "root slash",
			grpcPath:      "/",
			wantTrimmed:   "",
			wantIsCustom:  false,
			wantWithSlash: "/",
		},
		{
			name:          "empty string",
			grpcPath:      "",
			wantTrimmed:   "",
			wantIsCustom:  false,
			wantWithSlash: "/",
		},
		{
			name:          "custom path with trailing slash",
			grpcPath:      "/node/",
			wantTrimmed:   "/node",
			wantIsCustom:  true,
			wantWithSlash: "/node/",
		},
		{
			name:          "custom path without trailing slash",
			grpcPath:      "/node",
			wantTrimmed:   "/node",
			wantIsCustom:  true,
			wantWithSlash: "/node/",
		},
		{
			name:          "custom path without leading slash",
			grpcPath:      "node",
			wantTrimmed:   "/node",
			wantIsCustom:  true,
			wantWithSlash: "/node/",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := BackendConfig{GrpcPath: tt.grpcPath}
			if got := cfg.Trimmed(); got != tt.wantTrimmed {
				t.Errorf("Trimmed() = %q, want %q", got, tt.wantTrimmed)
			}
			if got := cfg.IsCustom(); got != tt.wantIsCustom {
				t.Errorf("IsCustom() = %v, want %v", got, tt.wantIsCustom)
			}
			if got := cfg.WithSlash(); got != tt.wantWithSlash {
				t.Errorf("WithSlash() = %q, want %q", got, tt.wantWithSlash)
			}
		})
	}
}

func TestValidateBasePath(t *testing.T) {
	valid := []string{
		"/",
		"",
		"/node",
		"/node/",
		"/custom-node_123/v1",
		"node",
	}

	for _, path := range valid {
		if err := validateBasePath(path); err != nil {
			t.Errorf("expected valid for %q, got error: %v", path, err)
		}
	}

	invalid := []string{
		"/../escape",
		"/node?query=1",
		"/node#hash",
		"/node with space",
		"/node\\backslash",
		"/node;injection",
	}

	for _, path := range invalid {
		if err := validateBasePath(path); err == nil {
			t.Errorf("expected invalid for %q, got nil error", path)
		}
	}
}

func TestLoadNodeConfigEnvFlags(t *testing.T) {
	t.Setenv("SECRET_KEY", "")
	t.Setenv("NODE_GRPC_TOKEN", "1234567890abcdef")
	t.Setenv("SNI_VERIFICATION", "true")
	t.Setenv("NFTABLES_LOGGING", "false")
	t.Setenv("NFTABLES_ACCEPT_REPLY_TRAFFIC", "true")

	cfg, err := LoadNodeConfig()
	if err != nil {
		t.Fatalf("LoadNodeConfig() error = %v", err)
	}

	if !cfg.Backend.SNIVerification {
		t.Errorf("SNIVerification = false, want true")
	}
	if cfg.Backend.NftablesLogging {
		t.Errorf("NftablesLogging = true, want false")
	}
	if !cfg.Backend.NftablesAcceptReplyTraffic {
		t.Errorf("NftablesAcceptReplyTraffic = false, want true")
	}
}

func TestParseBoolEnvDefault(t *testing.T) {
	t.Setenv("TEST_BOOL_TRUE", "true")
	t.Setenv("TEST_BOOL_FALSE", "false")
	t.Setenv("TEST_BOOL_ONE", "1")
	t.Setenv("TEST_BOOL_ZERO", "0")
	t.Setenv("TEST_BOOL_EMPTY", "")

	if !parseBoolEnvDefault("TEST_BOOL_TRUE", false) {
		t.Errorf("expected true")
	}
	if parseBoolEnvDefault("TEST_BOOL_FALSE", true) {
		t.Errorf("expected false")
	}
	if !parseBoolEnvDefault("TEST_BOOL_ONE", false) {
		t.Errorf("expected true")
	}
	if parseBoolEnvDefault("TEST_BOOL_ZERO", true) {
		t.Errorf("expected false")
	}
	if !parseBoolEnvDefault("TEST_BOOL_EMPTY", true) {
		t.Errorf("expected default true")
	}
	if parseBoolEnvDefault("NON_EXISTING_ENV", false) {
		t.Errorf("expected default false")
	}
}

