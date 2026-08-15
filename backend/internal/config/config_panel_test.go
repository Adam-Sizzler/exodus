package config

import "testing"

func TestPanelConfigMethods(t *testing.T) {
	tests := []struct {
		name          string
		basePath      string
		wantTrimmed   string
		wantIsCustom  bool
		wantWithSlash string
	}{
		{
			name:          "root slash",
			basePath:      "/",
			wantTrimmed:   "",
			wantIsCustom:  false,
			wantWithSlash: "/",
		},
		{
			name:          "empty string",
			basePath:      "",
			wantTrimmed:   "",
			wantIsCustom:  false,
			wantWithSlash: "/",
		},
		{
			name:          "custom path with trailing slash",
			basePath:      "/exodus_path/",
			wantTrimmed:   "/exodus_path",
			wantIsCustom:  true,
			wantWithSlash: "/exodus_path/",
		},
		{
			name:          "custom path without trailing slash",
			basePath:      "/exodus_path",
			wantTrimmed:   "/exodus_path",
			wantIsCustom:  true,
			wantWithSlash: "/exodus_path/",
		},
		{
			name:          "custom path without leading slash",
			basePath:      "panel",
			wantTrimmed:   "/panel",
			wantIsCustom:  true,
			wantWithSlash: "/panel/",
		},
		{
			name:          "nested custom path",
			basePath:      "/admin/panel/",
			wantTrimmed:   "/admin/panel",
			wantIsCustom:  true,
			wantWithSlash: "/admin/panel/",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := PanelConfig{BasePath: tt.basePath}
			if got := p.Trimmed(); got != tt.wantTrimmed {
				t.Errorf("Trimmed() = %q, want %q", got, tt.wantTrimmed)
			}
			if got := p.IsCustom(); got != tt.wantIsCustom {
				t.Errorf("IsCustom() = %v, want %v", got, tt.wantIsCustom)
			}
			if got := p.WithSlash(); got != tt.wantWithSlash {
				t.Errorf("WithSlash() = %q, want %q", got, tt.wantWithSlash)
			}
		})
	}
}

func TestValidateBasePath(t *testing.T) {
	valid := []string{
		"/",
		"",
		"/panel",
		"/panel/",
		"/exodus_path",
		"/exodus-path_123/sub",
		"custom-panel",
	}

	for _, path := range valid {
		if err := validateBasePath(path); err != nil {
			t.Errorf("expected valid for %q, got error: %v", path, err)
		}
	}

	invalid := []string{
		"/../escape",
		"/panel?query=1",
		"/panel#hash",
		"/panel with space",
		"/panel\\backslash",
		"/panel;injection",
		"/panel<script>",
	}

	for _, path := range invalid {
		if err := validateBasePath(path); err == nil {
			t.Errorf("expected invalid for %q, got nil error", path)
		}
	}
}
