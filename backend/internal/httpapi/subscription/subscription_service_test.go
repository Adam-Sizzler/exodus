package subscription

import (
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"gopkg.in/yaml.v3"
)

func getSubscriptionRefillDateAt(strategy string, now time.Time) string {
	var next time.Time
	local := now.Local()
	switch strategy {
	case "DAY":
		next = time.Date(local.Year(), local.Month(), local.Day()+1, 0, 0, 0, 0, local.Location())
	case "WEEK":
		daysUntilMonday := (8 - int(local.Weekday())) % 7
		if daysUntilMonday == 0 {
			daysUntilMonday = 7
		}
		next = time.Date(local.Year(), local.Month(), local.Day()+daysUntilMonday, 0, 0, 0, 0, local.Location())
	case "MONTH":
		next = time.Date(local.Year(), local.Month()+1, 1, 0, 0, 0, 0, local.Location())
	default:
		return ""
	}
	return fmt.Sprintf("%d", next.Unix())
}

func firstHostTagFromSlice(tags []string) *string {
	for _, t := range tags {
		trimmed := strings.TrimSpace(t)
		if trimmed != "" {
			return &trimmed
		}
	}
	return nil
}

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

func TestFirstHostTag(t *testing.T) {
	tag := firstHostTagFromSlice([]string{"", " edge ", "backup"})
	if tag == nil || *tag != "edge" {
		t.Fatalf("got %#v, want edge", tag)
	}

	if tag := firstHostTagFromSlice([]string{"", "   "}); tag != nil {
		t.Fatalf("got %#v, want nil", tag)
	}
}

func TestBuildXrayOutboundIncludesXHTTPExtraAndSockopt(t *testing.T) {
	protocol := "vless"
	network := "xhttp"
	path := "/host-path"
	hostHeader := "cdn.example.com"
	xhttpExtra := `{"xPaddingBytes":"100-1000","xmux":{"maxConcurrency":"16-32"}}`
	sockopt := `{"tcpNoDelay":true}`

	outbound := buildXrayOutbound(SubscriptionHost{
		Remark:           "xhttp",
		Address:          "edge.example.com",
		Port:             443,
		Path:             &path,
		Host:             &hostHeader,
		InboundType:      &protocol,
		InboundNetwork:   &network,
		XHTTPExtraParams: &xhttpExtra,
		SockoptParams:    &sockopt,
		InboundRaw: json.RawMessage(`{
			"streamSettings": {
				"xhttpSettings": {
					"mode": "auto",
					"path": "/inbound-path"
				}
			}
		}`),
	}, SubscriptionUser{VlessUUID: "9f76f8d8-daf1-4db6-b045-0987cd5e09a2"})

	streamSettings, ok := outbound["streamSettings"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected streamSettings map, got %#v", outbound["streamSettings"])
	}
	xhttpSettings, ok := streamSettings["xhttpSettings"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected xhttpSettings map, got %#v", streamSettings["xhttpSettings"])
	}
	if got := xhttpSettings["mode"]; got != "auto" {
		t.Fatalf("expected inbound xhttp mode to be preserved, got %#v", got)
	}
	if got := xhttpSettings["path"]; got != path {
		t.Fatalf("expected host path override %q, got %#v", path, got)
	}
	if got := xhttpSettings["host"]; got != hostHeader {
		t.Fatalf("expected host header override %q, got %#v", hostHeader, got)
	}
	extra, ok := xhttpSettings["extra"].(map[string]interface{})
	if !ok || extra["xPaddingBytes"] != "100-1000" {
		t.Fatalf("expected XHTTP extra params, got %#v", xhttpSettings["extra"])
	}
	gotSockopt, ok := streamSettings["sockopt"].(map[string]interface{})
	if !ok || gotSockopt["tcpNoDelay"] != true {
		t.Fatalf("expected SockOpt params, got %#v", streamSettings["sockopt"])
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
			Remark:          "seiko",
			Address:         "one.example.com",
			Port:            443,
			InboundType:     &protocol,
			InboundNetwork:  &network,
			InboundSecurity: &security,
		},
		{
			Remark:          "aeza",
			Address:         "two.example.com",
			Port:            443,
			InboundType:     &protocol,
			InboundNetwork:  &network,
			InboundSecurity: &security,
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
      - auto
      - seiko
      - aeza
      - timeweb
`) {
		t.Fatalf("expected select group to keep template entries before host nodes, got:\n%s", rendered)
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
			Remark:          "seiko",
			Address:         "one.example.com",
			Port:            443,
			InboundType:     &protocol,
			InboundNetwork:  &network,
			InboundSecurity: &security,
		},
		{
			Remark:          "aeza",
			Address:         "two.example.com",
			Port:            443,
			InboundType:     &protocol,
			InboundNetwork:  &network,
			InboundSecurity: &security,
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
        "auto",
        "seiko",
        "aeza",
        "timeweb"
      ]
    }`) {
		t.Fatalf("expected selector order to keep template entries before host nodes, got:\n%s", rendered)
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
        "enabled": true,
        "server_name": "one.example.com"
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
        "enabled": true,
        "server_name": "two.example.com"
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
        "enabled": true,
        "server_name": "three.example.com"
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
			Remark:          "aeza",
			Address:         "two.example.com",
			Port:            443,
			InboundType:     &protocol,
			InboundNetwork:  &network,
			InboundSecurity: &security,
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
        "auto",
        "direct",
        "aeza",
        "timeweb"
      ]
    }`) {
		t.Fatalf("expected selector to preserve template entries before host nodes, got:\n%s", rendered)
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
			Remark:          "aeza",
			Address:         "two.example.com",
			Port:            443,
			InboundType:     &protocol,
			InboundNetwork:  &network,
			InboundSecurity: &security,
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
        "auto",
        "seiko",
        "aeza",
        "timeweb"
      ]
    }`) {
		t.Fatalf("expected selector to keep template entries before host nodes, got:\n%s", rendered)
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

func TestGenerateSingboxConfigUsesSingboxProtocolCredentials(t *testing.T) {
	template := []byte(`{"outbounds":[]}`)
	security := "tls"
	vlessUUID := "9f76f8d8-daf1-4db6-b045-0987cd5e09a2"
	user := SubscriptionUser{
		Username:          "alice",
		VlessUUID:         vlessUUID,
		TrojanPassword:    "trojan-secret",
		Hysteria2Password: "hy-secret",
	}

	tests := []struct {
		name     string
		protocol string
		network  string
		want     []string
		notWant  []string
	}{
		{
			name:     "hysteria uses auth_str password",
			protocol: "hysteria",
			network:  "udp",
			want: []string{
				`"type": "hysteria"`,
				`"auth_str": "hy-secret"`,
			},
			notWant: []string{`"password": "` + vlessUUID + `"`},
		},
		{
			name:     "hysteria2 uses password not uuid",
			protocol: "hysteria2",
			network:  "udp",
			want: []string{
				`"type": "hysteria2"`,
				`"password": "hy-secret"`,
			},
			notWant: []string{`"password": "` + vlessUUID + `"`},
		},
		{
			name:     "vmess uses uuid",
			protocol: "vmess",
			network:  "tcp",
			want: []string{
				`"type": "vmess"`,
				`"uuid": "` + vlessUUID + `"`,
			},
		},
		{
			name:     "tuic uses uuid and password",
			protocol: "tuic",
			network:  "udp",
			want: []string{
				`"type": "tuic"`,
				`"uuid": "` + vlessUUID + `"`,
				`"password": "trojan-secret"`,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			rendered, err := generateSingboxConfig(template, []SubscriptionHost{
				{
					Remark:          tc.protocol,
					Address:         "example.com",
					Port:            443,
					InboundType:     &tc.protocol,
					InboundNetwork:  &tc.network,
					InboundSecurity: &security,
				},
			}, user)
			if err != nil {
				t.Fatalf("generateSingboxConfig returned error: %v", err)
			}
			for _, want := range tc.want {
				if !strings.Contains(rendered, want) {
					t.Fatalf("expected rendered config to contain %q, got:\n%s", want, rendered)
				}
			}
			for _, notWant := range tc.notWant {
				if strings.Contains(rendered, notWant) {
					t.Fatalf("expected rendered config not to contain %q, got:\n%s", notWant, rendered)
				}
			}
		})
	}
}

func TestHysteria2CredentialDoesNotFallbackToOtherProtocols(t *testing.T) {
	protocol := "hysteria2"
	host := SubscriptionHost{InboundType: &protocol}
	user := SubscriptionUser{
		VlessUUID:      "9f76f8d8-daf1-4db6-b045-0987cd5e09a2",
		TrojanPassword: "trojan-secret",
	}

	if got := effectiveProtocolCredential(host, user); got != "" {
		t.Fatalf("got %q, want empty dedicated hysteria2 password", got)
	}

	user.Hysteria2Password = "hysteria-secret"
	if got := effectiveProtocolCredential(host, user); got != "hysteria-secret" {
		t.Fatalf("got %q, want dedicated hysteria2 password", got)
	}
}

func TestExtractHwidHeadersExodusCompatibleHeaders(t *testing.T) {
	req := httptest.NewRequest("GET", "/api/sub/short", nil)
	req.Header.Set("X-HWID", "device-123")
	req.Header.Set("X-Device-OS", "Android")
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

func TestExtractHwidHeadersInfersSingBoxPlatformFromUserAgent(t *testing.T) {
	req := httptest.NewRequest("GET", "/api/sub/short", nil)
	req.Header.Set("X-HWID", "device-789")
	req.Header.Set("User-Agent", "SFA/1.13.11 (662; sing-box 1.13.11; language ru_RU_u_fw_mon_ms_metric_mu_celsius)")

	got := extractHwidHeaders(req)
	if got == nil {
		t.Fatal("expected hwid headers, got nil")
	}
	assertStringPtr(t, "platform", got.Platform, "android")
	assertStringPtr(t, "device model", got.DeviceModel, "unknown")
	assertStringPtr(t, "user agent", got.UserAgent, "SFA/1.13.11 (662; sing-box 1.13.11; language ru_RU_u_fw_mon_ms_metric_mu_celsius)")
}

func TestExtractHwidHeadersRequiresHwid(t *testing.T) {
	req := httptest.NewRequest("GET", "/api/sub/short", nil)
	req.Header.Set("X-Device-OS", "android")

	if got := extractHwidHeaders(req); got != nil {
		t.Fatalf("expected nil without X-HWID, got %#v", got)
	}
}

func TestExtractSyntheticHwidHeadersFromV2rayNUserAgent(t *testing.T) {
	req := httptest.NewRequest("GET", "/api/sub/short", nil)
	req.Header.Set("User-Agent", "v2rayN/7.16.4")

	got := extractSyntheticHwidHeaders(req, "1af12108-83fa-4f23-bacd-2fe7b4df3e71", "203.0.113.10")
	if got == nil {
		t.Fatal("expected synthetic hwid headers, got nil")
	}
	if !got.Synthetic {
		t.Fatal("expected synthetic flag")
	}
	if _, err := uuid.Parse(got.Hwid); err != nil {
		t.Fatalf("unexpected synthetic hwid: %q", got.Hwid)
	}
	assertStringPtr(t, "platform", got.Platform, "windows")
	assertStringPtr(t, "device model", got.DeviceModel, "unknown")
	assertStringPtr(t, "user agent", got.UserAgent, "v2rayN/7.16.4")
}

func TestInferKnownSinglePlatformClientsFromUserAgent(t *testing.T) {
	cases := map[string]string{
		"v2rayN/7.16.4":      "windows",
		"v2rayNG/1.10.31":    "android",
		"Streisand/1.7.8":    "ios",
		"V2Box/1.5.5":        "ios",
		"rabbithole/1.0":     "ios",
		"Exclave/1.0":        "android",
		"SFW/26.4.12":        "windows",
		"SFA/1.13.11 (662)":  "android",
		"SFI/1.13.11 (662)":  "ios",
		"SFM/1.13.11 (662)":  "macos",
		"SFL/1.13.11 (662)":  "linux",
		"SFATV/1.13.11 (12)": "android",
	}

	for userAgent, want := range cases {
		t.Run(userAgent, func(t *testing.T) {
			if got := inferPlatformFromUserAgent(userAgent); got != want {
				t.Fatalf("unexpected platform: got %q, want %q", got, want)
			}
		})
	}
}

func TestExtractSyntheticHwidHeadersIsStableAcrossRequestIP(t *testing.T) {
	req := httptest.NewRequest("GET", "/api/sub/short", nil)
	req.Header.Set("User-Agent", "SFA/1.13.11 (662; sing-box 1.13.11; language ru_RU_u_fw_mon_ms_metric_mu_celsius)")

	first := extractSyntheticHwidHeaders(req, "1af12108-83fa-4f23-bacd-2fe7b4df3e71", "203.0.113.10")
	second := extractSyntheticHwidHeaders(req, "1af12108-83fa-4f23-bacd-2fe7b4df3e71", "198.51.100.7")
	if first == nil || second == nil {
		t.Fatalf("expected synthetic hwid headers, got first=%#v second=%#v", first, second)
	}
	if first.Hwid != second.Hwid {
		t.Fatalf("expected stable hwid across IPs, got %q and %q", first.Hwid, second.Hwid)
	}
	if _, err := uuid.Parse(first.Hwid); err != nil {
		t.Fatalf("expected UUID synthetic hwid, got %q", first.Hwid)
	}
	assertStringPtr(t, "platform", first.Platform, "android")
	assertStringPtr(t, "device model", first.DeviceModel, "unknown")
}

func TestExtractSyntheticHwidHeadersChangesWhenAdditionalMetadataAppears(t *testing.T) {
	userUUID := "1af12108-83fa-4f23-bacd-2fe7b4df3e71"
	req := httptest.NewRequest("GET", "/api/sub/short", nil)
	req.Header.Set("User-Agent", "SFA/1.13.11 (662; sing-box 1.13.11; language ru_RU_u_fw_mon_ms_metric_mu_celsius)")

	withModel := httptest.NewRequest("GET", "/api/sub/short", nil)
	withModel.Header.Set("User-Agent", "SFA/1.13.11 (662; sing-box 1.13.11; language ru_RU_u_fw_mon_ms_metric_mu_celsius)")
	withModel.Header.Set("X-Device-Model", "Pixel 8")

	first := extractSyntheticHwidHeaders(req, userUUID, "203.0.113.10")
	second := extractSyntheticHwidHeaders(withModel, userUUID, "203.0.113.10")
	if first == nil || second == nil {
		t.Fatalf("expected synthetic hwid headers, got first=%#v second=%#v", first, second)
	}
	if first.Hwid == second.Hwid {
		t.Fatalf("expected different hwid when device metadata appears, got %q", first.Hwid)
	}
}

func TestExtractSyntheticHwidHeadersFromPlatformUserAgent(t *testing.T) {
	req := httptest.NewRequest("GET", "/api/sub/short", nil)
	req.Header.Set("User-Agent", "FlClash/v0.8.92 clash-verge Platform/windows")

	got := extractSyntheticHwidHeaders(req, "user-uuid", "203.0.113.10")
	if got == nil {
		t.Fatal("expected synthetic hwid headers, got nil")
	}
	assertStringPtr(t, "platform", got.Platform, "windows")
	assertStringPtr(t, "device model", got.DeviceModel, "unknown")
}

func TestExtractSyntheticHwidHeadersPrefersRealDeviceHeaders(t *testing.T) {
	req := httptest.NewRequest("GET", "/api/sub/short", nil)
	req.Header.Set("User-Agent", "UnknownClient/1.0")
	req.Header.Set("X-Device-OS", "Android")
	req.Header.Set("X-Ver-OS", "15")
	req.Header.Set("X-Device-Model", "Pixel 8")

	got := extractSyntheticHwidHeaders(req, "user-uuid", "203.0.113.10")
	if got == nil {
		t.Fatal("expected synthetic hwid headers, got nil")
	}
	assertStringPtr(t, "platform", got.Platform, "android")
	assertStringPtr(t, "os version", got.OsVersion, "15")
	assertStringPtr(t, "device model", got.DeviceModel, "Pixel 8")
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

func TestBuildSubscriptionLinksAndResolvedProxies(t *testing.T) {
	user := SubscriptionUser{
		TID:               1,
		UUID:              "11111111-1111-1111-1111-111111111111",
		ShortUUID:         "testuser",
		Username:          "testuser",
		Status:            "ACTIVE",
		TrafficLimitBytes: 107374182400,
		VlessUUID:         "22222222-2222-2222-2222-222222222222",
		TrojanPassword:    "trojan-pass",
		SSPassword:        "ss-pass",
		Hysteria2Password: "hy2-pass",
		AnytlsPassword:    "anytls-pass",
		UsedTrafficBytes:  1024,
		LifetimeUsedBytes: 2048,
	}

	inboundTypeVless := "vless"
	inboundTypeTrojan := "trojan"
	inboundTypeSS := "shadowsocks"
	inboundTypeHy2 := "hysteria2"
	sni := "example.com"
	networkWs := "ws"
	pathWs := "/ws"

	hosts := []SubscriptionHost{
		{
			UUID:           "host-1",
			Remark:         "🇸🇪 Sweden VLESS",
			Address:        "se.example.com",
			Port:           443,
			InboundType:    &inboundTypeVless,
			InboundNetwork: &networkWs,
			Path:           &pathWs,
			SNI:            &sni,
			SecurityLayer:  "TLS",
		},
		{
			UUID:          "host-2",
			Remark:        "🇩🇪 Germany Trojan",
			Address:       "de.example.com",
			Port:          443,
			InboundType:   &inboundTypeTrojan,
			SNI:           &sni,
			SecurityLayer: "TLS",
		},
		{
			UUID:        "host-3",
			Remark:      "🇺🇸 US Shadowsocks",
			Address:     "us.example.com",
			Port:        8388,
			InboundType: &inboundTypeSS,
		},
		{
			UUID:        "host-4",
			Remark:      "🇫🇷 France Hy2",
			Address:     "fr.example.com",
			Port:        443,
			InboundType: &inboundTypeHy2,
			SNI:         &sni,
		},
	}

	links, _ := buildSubscriptionLinks(hosts, user)
	if len(links) != 4 {
		t.Fatalf("expected 4 links, got %d", len(links))
	}
	if !strings.HasPrefix(links[0], "vless://") {
		t.Fatalf("expected vless link, got %s", links[0])
	}
	if !strings.HasPrefix(links[1], "trojan://") {
		t.Fatalf("expected trojan link, got %s", links[1])
	}
	if !strings.HasPrefix(links[2], "ss://") {
		t.Fatalf("expected ss link, got %s", links[2])
	}
	if !strings.HasPrefix(links[3], "hysteria2://") {
		t.Fatalf("expected hysteria2 link, got %s", links[3])
	}

	resolved := buildResolvedProxyConfigs(hosts, user)
	if len(resolved) != 4 {
		t.Fatalf("expected 4 resolved proxy configs, got %d", len(resolved))
	}
	if resolved[0].Protocol != "vless" || resolved[0].Transport != "ws" {
		t.Fatalf("unexpected resolved proxy: %+v", resolved[0])
	}
	if resolved[1].Protocol != "trojan" {
		t.Fatalf("unexpected resolved proxy: %+v", resolved[1])
	}
	if resolved[2].Protocol != "shadowsocks" {
		t.Fatalf("unexpected resolved proxy: %+v", resolved[2])
	}
	if resolved[3].Protocol != "hysteria" {
		t.Fatalf("unexpected resolved proxy: %+v", resolved[3])
	}
}
