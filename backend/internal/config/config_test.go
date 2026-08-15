package config

import "testing"

func TestConfigSubPathMethods(t *testing.T) {
	tests := []struct {
		name          string
		subPath       string
		wantTrimmed   string
		wantIsCustom  bool
		wantWithSlash string
	}{
		{
			name:          "root slash",
			subPath:       "/",
			wantTrimmed:   "",
			wantIsCustom:  false,
			wantWithSlash: "/",
		},
		{
			name:          "empty string",
			subPath:       "",
			wantTrimmed:   "",
			wantIsCustom:  false,
			wantWithSlash: "/",
		},
		{
			name:          "custom prefix with trailing slash",
			subPath:       "/subscription/",
			wantTrimmed:   "/subscription",
			wantIsCustom:  true,
			wantWithSlash: "/subscription/",
		},
		{
			name:          "custom prefix without trailing slash",
			subPath:       "/subscription",
			wantTrimmed:   "/subscription",
			wantIsCustom:  true,
			wantWithSlash: "/subscription/",
		},
		{
			name:          "custom prefix without leading slash",
			subPath:       "custom-sub",
			wantTrimmed:   "/custom-sub",
			wantIsCustom:  true,
			wantWithSlash: "/custom-sub/",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := Config{SubPath: tt.subPath}
			if got := cfg.SubPathTrimmed(); got != tt.wantTrimmed {
				t.Errorf("SubPathTrimmed() = %q, want %q", got, tt.wantTrimmed)
			}
			if got := cfg.IsCustomSubPath(); got != tt.wantIsCustom {
				t.Errorf("IsCustomSubPath() = %v, want %v", got, tt.wantIsCustom)
			}
			if got := cfg.SubPathWithSlash(); got != tt.wantWithSlash {
				t.Errorf("SubPathWithSlash() = %q, want %q", got, tt.wantWithSlash)
			}
		})
	}
}

func TestValidateSubPath(t *testing.T) {
	valid := []string{
		"/",
		"",
		"/subscription",
		"/subscription/",
		"/custom-sub_123/v1",
		"sub",
	}

	for _, path := range valid {
		if err := validateSubPath(path); err != nil {
			t.Errorf("expected valid for %q, got error: %v", path, err)
		}
	}

	invalid := []string{
		"/../escape",
		"/sub?query=1",
		"/sub#hash",
		"/sub with space",
		"/sub\\backslash",
		"/sub;injection",
	}

	for _, path := range invalid {
		if err := validateSubPath(path); err == nil {
			t.Errorf("expected invalid for %q, got nil error", path)
		}
	}
}
