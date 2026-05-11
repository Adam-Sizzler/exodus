package subscription

import (
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"gopkg.in/yaml.v3"
)

func TestGetSubscriptionRefillDateAt(t *testing.T) {
	loc := time.FixedZone("test", 7*60*60)
	now := time.Date(2026, 5, 10, 13, 0, 0, 0, loc) // Sunday.

	cases := []struct {
		strategy string
		want     time.Time
	}{
		{
			strategy: "DAY",
			want:     time.Date(2026, 5, 11, 0, 0, 0, 0, loc),
		},
		{
			strategy: "WEEK",
			want:     time.Date(2026, 5, 11, 0, 0, 0, 0, loc),
		},
		{
			strategy: "MONTH",
			want:     time.Date(2026, 6, 1, 0, 0, 0, 0, loc),
		},
	}

	for _, tc := range cases {
		t.Run(tc.strategy, func(t *testing.T) {
			got := getSubscriptionRefillDateAt(tc.strategy, now)
			want := fmt.Sprintf("%d", tc.want.Unix())
			if got != want {
				t.Fatalf("got %s, want %s", got, want)
			}
		})
	}

	if got := getSubscriptionRefillDateAt("NO_RESET", now); got != "" {
		t.Fatalf("NO_RESET got %q, want empty", got)
	}
}

func TestGenerateYAMLConfigPreservesTemplateOrderAndSelectorNodePlacement(t *testing.T) {
	template := []byte(`allow-lan: true
mode: rule

proxy-groups:
  - name: auto
    type: url-test
    proxies:
    url: https://cp.cloudflare.com/generate_204
    interval: 300
    tolerance: 50

  - name: select
    type: select
    proxies:
      - auto

proxies:

rules:
  - MATCH,select
`)

	network := "tcp"
	security := "tls"
	protocol := "trojan"
	hosts := []SubscriptionHost{
		{
			Remark:             "seiko",
			Address:            "one.example.com",
			Port:               443,
			InboundType:        &protocol,
			InboundNetwork:     &network,
			InboundSecurity:    &security,
			SelectorNodesFirst: true,
		},
		{
			Remark:             "aeza",
			Address:            "two.example.com",
			Port:               443,
			InboundType:        &protocol,
			InboundNetwork:     &network,
			InboundSecurity:    &security,
			SelectorNodesFirst: true,
		},
		{
			Remark:          "timeweb",
			Address:         "three.example.com",
			Port:            443,
			InboundType:     &protocol,
			InboundNetwork:  &network,
			InboundSecurity: &security,
		},
	}

	rendered, err := generateYAMLConfig(template, hosts, SubscriptionUser{TrojanPassword: "secret"})
	if err != nil {
		t.Fatalf("generateYAMLConfig returned error: %v", err)
	}

	if !strings.Contains(rendered, "mode: rule\n\nproxy-groups:") {
		t.Fatalf("expected blank line before proxy-groups, got:\n%s", rendered)
	}

	if !strings.Contains(rendered, `proxy-groups:
  - name: auto
    type: url-test
    proxies:
      - seiko
      - aeza
      - timeweb
    url: https://cp.cloudflare.com/generate_204
    interval: 300
    tolerance: 50
`) {
		t.Fatalf("expected auto group to contain all hosts in host order, got:\n%s", rendered)
	}

	if !strings.Contains(rendered, `  - name: select
    type: select
    proxies:
      - seiko
      - aeza
      - auto
      - timeweb
`) {
		t.Fatalf("expected select group to place flagged nodes before auto and others after, got:\n%s", rendered)
	}

	proxyGroupsIndex := strings.Index(rendered, "proxy-groups:")
	proxiesIndex := strings.Index(rendered, "\nproxies:")
	rulesIndex := strings.Index(rendered, "\nrules:")
	if proxyGroupsIndex == -1 || proxiesIndex == -1 || rulesIndex == -1 {
		t.Fatalf("expected proxy-groups, proxies, and rules sections, got:\n%s", rendered)
	}
	if !(proxyGroupsIndex < proxiesIndex && proxiesIndex < rulesIndex) {
		t.Fatalf("expected template section order to be preserved, got:\n%s", rendered)
	}

	var parsed yaml.Node
	if err := yaml.Unmarshal([]byte(rendered), &parsed); err != nil {
		t.Fatalf("expected rendered YAML to stay valid, got error: %v\n%s", err, rendered)
	}
}

func TestGenerateSingboxConfigPreservesNestedTemplateOrderAndSelectorNodePlacement(t *testing.T) {
	template := []byte(`{
  "outbounds": [
    {
      "type": "selector",
      "tag": "select",
      "outbounds": null
    },
    {
      "type": "urltest",
      "tag": "auto",
      "url": "https://cp.cloudflare.com/generate_204",
      "outbounds": null
    },
    {
      "type": "direct",
      "tag": "direct"
    }
  ]
}`)

	protocol := "trojan"
	network := "tcp"
	security := "tls"
	rendered, err := generateSingboxConfig(template, []SubscriptionHost{
		{
			Remark:             "seiko",
			Address:            "one.example.com",
			Port:               443,
			InboundType:        &protocol,
			InboundNetwork:     &network,
			InboundSecurity:    &security,
			SelectorNodesFirst: true,
		},
		{
			Remark:             "aeza",
			Address:            "two.example.com",
			Port:               443,
			InboundType:        &protocol,
			InboundNetwork:     &network,
			InboundSecurity:    &security,
			SelectorNodesFirst: true,
		},
		{
			Remark:          "timeweb",
			Address:         "three.example.com",
			Port:            443,
			InboundType:     &protocol,
			InboundNetwork:  &network,
			InboundSecurity: &security,
		},
	}, SubscriptionUser{TrojanPassword: "secret"})
	if err != nil {
		t.Fatalf("generateSingboxConfig returned error: %v", err)
	}

	if !strings.Contains(rendered, `{
      "type": "selector",
      "tag": "select",
      "outbounds": [
        "seiko",
        "aeza",
        "auto",
        "timeweb"
      ]
    }`) {
		t.Fatalf("expected selector order to be preferred nodes, template entries, then remaining nodes, got:\n%s", rendered)
	}

	if !strings.Contains(rendered, `{
      "type": "urltest",
      "tag": "auto",
      "url": "https://cp.cloudflare.com/generate_204",
      "outbounds": [
        "seiko",
        "aeza",
        "timeweb"
      ]
    }`) {
		t.Fatalf("expected urltest to keep node order by view position, got:\n%s", rendered)
	}

	if !strings.Contains(rendered, `{
      "type": "trojan",
      "tag": "seiko",
      "server": "one.example.com",
      "server_port": 443,
      "password": "secret",
      "tls": {
        "enabled": true
      }
    }`) {
		t.Fatalf("expected first generated outbound order to be stable, got:\n%s", rendered)
	}

	if !strings.Contains(rendered, `{
      "type": "trojan",
      "tag": "aeza",
      "server": "two.example.com",
      "server_port": 443,
      "password": "secret",
      "tls": {
        "enabled": true
      }
    }`) {
		t.Fatalf("expected second generated outbound order to be stable, got:\n%s", rendered)
	}

	if !strings.Contains(rendered, `{
      "type": "trojan",
      "tag": "timeweb",
      "server": "three.example.com",
      "server_port": 443,
      "password": "secret",
      "tls": {
        "enabled": true
      }
    }`) {
		t.Fatalf("expected third generated outbound order to be stable, got:\n%s", rendered)
	}

	var parsed map[string]interface{}
	if err := json.Unmarshal([]byte(rendered), &parsed); err != nil {
		t.Fatalf("expected rendered JSON to stay valid, got error: %v\n%s", err, rendered)
	}
}

func TestGenerateSingboxConfigKeepsTemplateSelectorEntriesBetweenLeadingAndTrailingNodes(t *testing.T) {
	template := []byte(`{
  "outbounds": [
    {
      "type": "selector",
      "tag": "select",
      "outbounds": ["auto", "direct"]
    },
    {
      "type": "urltest",
      "tag": "auto",
      "url": "https://cp.cloudflare.com/generate_204",
      "outbounds": null
    },
    {
      "type": "direct",
      "tag": "direct"
    }
  ]
}`)

	protocol := "trojan"
	network := "tcp"
	security := "tls"
	rendered, err := generateSingboxConfig(template, []SubscriptionHost{
		{
			Remark:             "aeza",
			Address:            "two.example.com",
			Port:               443,
			InboundType:        &protocol,
			InboundNetwork:     &network,
			InboundSecurity:    &security,
			SelectorNodesFirst: true,
		},
		{
			Remark:          "timeweb",
			Address:         "three.example.com",
			Port:            443,
			InboundType:     &protocol,
			InboundNetwork:  &network,
			InboundSecurity: &security,
		},
	}, SubscriptionUser{TrojanPassword: "secret"})
	if err != nil {
		t.Fatalf("generateSingboxConfig returned error: %v", err)
	}

	if !strings.Contains(rendered, `{
      "type": "selector",
      "tag": "select",
      "outbounds": [
        "aeza",
        "auto",
        "direct",
        "timeweb"
      ]
    }`) {
		t.Fatalf("expected selector to preserve template entries between leading and trailing nodes, got:\n%s", rendered)
	}
}

func TestGenerateSingboxConfigKeepsURLTestOrderByHostOrder(t *testing.T) {
	template := []byte(`{
  "outbounds": [
    {
      "type": "selector",
      "tag": "select",
      "outbounds": null
    },
    {
      "type": "urltest",
      "tag": "auto",
      "url": "https://cp.cloudflare.com/generate_204",
      "outbounds": null
    }
  ]
}`)

	protocol := "trojan"
	network := "tcp"
	security := "tls"
	rendered, err := generateSingboxConfig(template, []SubscriptionHost{
		{
			Remark:          "seiko",
			Address:         "one.example.com",
			Port:            443,
			InboundType:     &protocol,
			InboundNetwork:  &network,
			InboundSecurity: &security,
		},
		{
			Remark:             "aeza",
			Address:            "two.example.com",
			Port:               443,
			InboundType:        &protocol,
			InboundNetwork:     &network,
			InboundSecurity:    &security,
			SelectorNodesFirst: true,
		},
		{
			Remark:          "timeweb",
			Address:         "three.example.com",
			Port:            443,
			InboundType:     &protocol,
			InboundNetwork:  &network,
			InboundSecurity: &security,
		},
	}, SubscriptionUser{TrojanPassword: "secret"})
	if err != nil {
		t.Fatalf("generateSingboxConfig returned error: %v", err)
	}

	if !strings.Contains(rendered, `{
      "type": "selector",
      "tag": "select",
      "outbounds": [
        "aeza",
        "auto",
        "seiko",
        "timeweb"
      ]
    }`) {
		t.Fatalf("expected selector to place preferred nodes before template entries, got:\n%s", rendered)
	}

	if !strings.Contains(rendered, `{
      "type": "urltest",
      "tag": "auto",
      "url": "https://cp.cloudflare.com/generate_204",
      "outbounds": [
        "seiko",
        "aeza",
        "timeweb"
      ]
    }`) {
		t.Fatalf("expected urltest to keep original host order, got:\n%s", rendered)
	}
}

func TestExtractHwidHeadersRemnawaveCompatibleHeaders(t *testing.T) {
	req := httptest.NewRequest("GET", "/api/sub/short", nil)
	req.Header.Set("X-HWID", "device-123")
	req.Header.Set("X-Device-OS", "android")
	req.Header.Set("X-Ver-OS", "15")
	req.Header.Set("X-Device-Model", "Pixel 8 Pro")
	req.Header.Set("User-Agent", "FlClash/0.8.86")

	got := extractHwidHeaders(req)
	if got == nil {
		t.Fatal("expected hwid headers, got nil")
	}
	if got.Hwid != "device-123" {
		t.Fatalf("unexpected hwid: %q", got.Hwid)
	}
	assertStringPtr(t, "platform", got.Platform, "android")
	assertStringPtr(t, "os version", got.OsVersion, "15")
	assertStringPtr(t, "device model", got.DeviceModel, "Pixel 8 Pro")
	assertStringPtr(t, "user agent", got.UserAgent, "FlClash/0.8.86")
}

func TestExtractHwidHeadersSupportsHwidFallbackHeaders(t *testing.T) {
	req := httptest.NewRequest("GET", "/api/sub/short", nil)
	req.Header.Set("X-HWID", "device-456")
	req.Header.Set("X-HWID-Platform", "ios")
	req.Header.Set("X-HWID-OS-Version", "18.1")
	req.Header.Set("X-HWID-Device-Model", "iPhone")
	req.Header.Set("X-HWID-User-Agent", "Shadowrocket/2.2.59")

	got := extractHwidHeaders(req)
	if got == nil {
		t.Fatal("expected hwid headers, got nil")
	}
	assertStringPtr(t, "platform", got.Platform, "ios")
	assertStringPtr(t, "os version", got.OsVersion, "18.1")
	assertStringPtr(t, "device model", got.DeviceModel, "iPhone")
	assertStringPtr(t, "user agent", got.UserAgent, "Shadowrocket/2.2.59")
}

func TestExtractHwidHeadersRequiresHwid(t *testing.T) {
	req := httptest.NewRequest("GET", "/api/sub/short", nil)
	req.Header.Set("X-Device-OS", "android")

	if got := extractHwidHeaders(req); got != nil {
		t.Fatalf("expected nil without X-HWID, got %#v", got)
	}
}

func assertStringPtr(t *testing.T, field string, got *string, want string) {
	t.Helper()
	if got == nil {
		t.Fatalf("expected %s %q, got nil", field, want)
	}
	if *got != want {
		t.Fatalf("unexpected %s: got %q, want %q", field, *got, want)
	}
}
