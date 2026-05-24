package server

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/iancoleman/orderedmap"
)

func TestBuildSingboxConfigWithV2RayAPIPreservesTopLevelOrder(t *testing.T) {
	raw := []byte(`{
  "log": {"level":"debug","output":"access.log","timestamp":true},
  "dns": {"servers":[{"type":"local","tag":"dns-main"}],"final":"dns-main"},
  "inbounds": [{"tag":"trojan-in","type":"trojan","listen":"0.0.0.0","listen_port":10443,"users":[{"name":"u1","password":"p1"}]}],
  "outbounds": [{"tag":"direct","type":"direct"}],
  "route": {"rules":[{"action":"sniff"}]}
}`)

	out, _, err := BuildSingboxConfigWithV2RayAPI(raw, BuildOptions{
		Listen:  "127.0.0.1:10085",
		Enabled: true,
	})
	if err != nil {
		t.Fatalf("BuildSingboxConfigWithV2RayAPI failed: %v", err)
	}

	var cfg orderedmap.OrderedMap
	if err := json.Unmarshal(out, &cfg); err != nil {
		t.Fatalf("unmarshal built config: %v", err)
	}

	wantTop := []string{"log", "dns", "inbounds", "outbounds", "route", "experimental"}
	if !reflect.DeepEqual(cfg.Keys(), wantTop) {
		t.Fatalf("top-level keys mismatch, want=%v got=%v", wantTop, cfg.Keys())
	}

	experimentalRaw, ok := cfg.Get("experimental")
	if !ok {
		t.Fatalf("missing experimental section")
	}
	experimental := mustOrderedMap(t, experimentalRaw)
	wantExperimental := []string{"cache_file", "v2ray_api"}
	if !reflect.DeepEqual(experimental.Keys(), wantExperimental) {
		t.Fatalf("experimental keys mismatch, want=%v got=%v", wantExperimental, experimental.Keys())
	}
}

func TestBuildSingboxConfigWithV2RayAPIPreservesNestedObjectOrder(t *testing.T) {
	raw := []byte(`{
  "log": {"level":"debug","output":"access.log","timestamp":true},
  "dns": {
    "servers": [{"type":"local","tag":"dns-main"}],
    "rules": [{"rule_set":["category-ads-all"],"action":"predefined","rcode":"NOERROR"}],
    "final": "dns-main"
  },
  "inbounds": [{"tag":"trojan-in","type":"trojan","listen":"0.0.0.0","listen_port":10443,"users":[]}],
  "outbounds": [{"tag":"direct","type":"direct"}],
  "route": {"rules":[{"action":"sniff"}]}
}`)

	out, _, err := BuildSingboxConfigWithV2RayAPI(raw, BuildOptions{
		Listen:  "127.0.0.1:10085",
		Enabled: true,
	})
	if err != nil {
		t.Fatalf("BuildSingboxConfigWithV2RayAPI failed: %v", err)
	}

	var cfg orderedmap.OrderedMap
	if err := json.Unmarshal(out, &cfg); err != nil {
		t.Fatalf("unmarshal built config: %v", err)
	}

	dns := mustOrderedMap(t, mustGet(t, cfg, "dns"))
	dnsRules := mustArray(t, mustGet(t, dns, "rules"))
	firstRule := mustOrderedMap(t, dnsRules[0])
	wantRuleOrder := []string{"rule_set", "action", "rcode"}
	if !reflect.DeepEqual(firstRule.Keys(), wantRuleOrder) {
		t.Fatalf("dns.rules[0] keys mismatch, want=%v got=%v", wantRuleOrder, firstRule.Keys())
	}
}

func TestShouldReloadCoreHonorsForceRestart(t *testing.T) {
	restart := true
	noRestart := false
	force := true
	soft := false

	cases := []struct {
		name          string
		task          DeployConfigTaskPayload
		configChanged bool
		want          bool
	}{
		{
			name:          "soft restart changed config",
			task:          DeployConfigTaskPayload{Restart: &restart, ForceRestart: &soft},
			configChanged: true,
			want:          true,
		},
		{
			name:          "soft restart unchanged config",
			task:          DeployConfigTaskPayload{Restart: &restart, ForceRestart: &soft},
			configChanged: false,
			want:          false,
		},
		{
			name:          "force restart unchanged config",
			task:          DeployConfigTaskPayload{Restart: &restart, ForceRestart: &force},
			configChanged: false,
			want:          true,
		},
		{
			name:          "camel force restart unchanged config",
			task:          DeployConfigTaskPayload{Restart: &restart, ForceRestartCamel: &force},
			configChanged: false,
			want:          true,
		},
		{
			name:          "restart disabled",
			task:          DeployConfigTaskPayload{Restart: &noRestart},
			configChanged: true,
			want:          false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := shouldReloadCore(tc.task, tc.configChanged); got != tc.want {
				t.Fatalf("unexpected reload decision: got=%v want=%v", got, tc.want)
			}
		})
	}
}

func mustGet(t *testing.T, o orderedmap.OrderedMap, key string) any {
	t.Helper()
	v, ok := o.Get(key)
	if !ok {
		t.Fatalf("missing key %q", key)
	}
	return v
}

func mustOrderedMap(t *testing.T, v any) orderedmap.OrderedMap {
	t.Helper()
	m, ok := v.(orderedmap.OrderedMap)
	if !ok {
		t.Fatalf("expected ordered map, got %T", v)
	}
	return m
}

func mustArray(t *testing.T, v any) []any {
	t.Helper()
	arr, ok := v.([]any)
	if !ok {
		t.Fatalf("expected array, got %T", v)
	}
	return arr
}
