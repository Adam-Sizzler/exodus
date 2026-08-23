package nodeplugins

import (
	"encoding/json"
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
	inboundTags, ok := haproxyAuth["inboundTags"].([]any)
	if !ok {
		t.Fatalf("expected haproxyAuth.inboundTags array, got %#v", haproxyAuth["inboundTags"])
	}
	if len(inboundTags) != 0 {
		t.Fatalf("expected haproxyAuth.inboundTags=[], got %#v", inboundTags)
	}
}

func TestNormalizePluginConfigPreservesPartialConfig(t *testing.T) {
	raw := json.RawMessage(`{"ingressFilter":{"enabled":true,"blockedIps":["203.0.113.1"]},"sharedLists":[]}`)

	normalized, err := normalizePluginConfig(raw)
	if err != nil {
		t.Fatalf("normalizePluginConfig returned error: %v", err)
	}

	var config map[string]any
	if err := json.Unmarshal(normalized, &config); err != nil {
		t.Fatalf("normalized config is not valid JSON: %v", err)
	}
	haproxyAuth, ok := config["haproxyAuth"].(map[string]any)
	if !ok {
		t.Fatal("normalized config is missing haproxyAuth object")
	}
	if inboundTags, ok := haproxyAuth["inboundTags"].([]any); !ok || len(inboundTags) != 0 {
		t.Fatalf("expected haproxyAuth.inboundTags=[], got %#v", haproxyAuth["inboundTags"])
	}
}

func TestNormalizePluginConfigNormalizesHaproxyConfig(t *testing.T) {
	raw := json.RawMessage(`{"haproxyAuth":{"enabled":false}}`)

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
	haproxyAuth := config["haproxyAuth"].(map[string]any)
	if inboundTags, ok := haproxyAuth["inboundTags"].([]any); !ok || len(inboundTags) != 0 {
		t.Fatalf("expected legacy enabled=false to normalize to inboundTags=[], got %#v", haproxyAuth["inboundTags"])
	}
	if _, ok := haproxyAuth["enabled"]; ok {
		t.Fatal("normalized haproxyAuth should not preserve legacy enabled flag")
	}
	if _, ok := config["ingressFilter"]; ok {
		t.Fatal("normalized config should not add ingressFilter")
	}
	if _, ok := config["egressFilter"]; ok {
		t.Fatal("normalized config should not add egressFilter")
	}
}

func TestNormalizePluginConfigConvertsLegacyHaproxyEnabled(t *testing.T) {
	raw := json.RawMessage(`{"haproxyAuth":{"enabled":true}}`)

	normalized, err := normalizePluginConfig(raw)
	if err != nil {
		t.Fatalf("normalizePluginConfig returned error: %v", err)
	}

	var config map[string]any
	if err := json.Unmarshal(normalized, &config); err != nil {
		t.Fatalf("normalized config is not valid JSON: %v", err)
	}
	haproxyAuth := config["haproxyAuth"].(map[string]any)
	inboundTags, ok := haproxyAuth["inboundTags"].([]any)
	if !ok || len(inboundTags) != 1 || inboundTags[0] != "*" {
		t.Fatalf("expected legacy enabled=true to normalize to inboundTags=[*], got %#v", haproxyAuth["inboundTags"])
	}
	if _, ok := haproxyAuth["enabled"]; ok {
		t.Fatal("normalized haproxyAuth should not preserve legacy enabled flag")
	}
}
