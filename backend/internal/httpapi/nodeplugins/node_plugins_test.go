package nodeplugins

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestNormalizePluginConfigDefaults(t *testing.T) {
	normalized, err := normalizePluginConfig(nil)
	if err != nil {
		t.Fatalf("normalizePluginConfig(nil) returned error: %v", err)
	}

	var config map[string]any
	if err := json.Unmarshal(normalized, &config); err != nil {
		t.Fatalf("normalized default config is not valid JSON: %v", err)
	}

	if _, ok := config["ingressFilter"]; !ok {
		t.Fatal("default config is missing ingressFilter")
	}
	if _, ok := config["egressFilter"]; !ok {
		t.Fatal("default config is missing egressFilter")
	}

	haproxyAuth, ok := config["haproxyAuth"].(map[string]any)
	if !ok {
		t.Fatal("default config is missing haproxyAuth object")
	}
	if enabled, ok := haproxyAuth["enabled"].(bool); !ok || enabled {
		t.Fatalf("expected haproxyAuth.enabled=false, got %#v", haproxyAuth["enabled"])
	}
}

func TestNormalizePluginConfigAddsHaproxyAuth(t *testing.T) {
	raw := json.RawMessage(`{"ingressFilter":{"enabled":true,"blockedIps":["203.0.113.1"]},"sharedLists":[]}`)

	normalized, err := normalizePluginConfig(raw)
	if err != nil {
		t.Fatalf("normalizePluginConfig returned error: %v", err)
	}

	var config map[string]any
	if err := json.Unmarshal(normalized, &config); err != nil {
		t.Fatalf("normalized config is not valid JSON: %v", err)
	}
	if _, ok := config["haproxyAuth"]; !ok {
		t.Fatal("normalized config is missing haproxyAuth")
	}
}

func TestNormalizePluginConfigRejectsUnsupportedPlugins(t *testing.T) {
	cases := []string{
		`{"torrentBlocker":{"enabled":true}}`,
		`{"connectionDrop":{"enabled":true}}`,
	}

	for _, tc := range cases {
		_, err := normalizePluginConfig(json.RawMessage(tc))
		if err == nil {
			t.Fatalf("expected unsupported plugin config to fail: %s", tc)
		}
		if !strings.Contains(err.Error(), "not supported") {
			t.Fatalf("expected unsupported error, got: %v", err)
		}
	}
}
