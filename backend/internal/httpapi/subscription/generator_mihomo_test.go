package subscription

import (
	"encoding/json"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestMihomoHysteria2Basic(t *testing.T) {
	inboundType := "hysteria2"
	host := SubscriptionHost{
		Remark:      "Hy2-Node",
		Address:     "hy2.example.com",
		Port:        443,
		InboundType: &inboundType,
	}
	user := SubscriptionUser{
		Hysteria2Password: "my-hy2-pass",
	}

	proxy := buildMihomoProxy(host, user)
	if proxy == nil {
		t.Fatal("expected proxy, got nil")
	}

	if proxy["type"] != "hysteria2" {
		t.Errorf("expected type 'hysteria2', got %v", proxy["type"])
	}
	if proxy["password"] != "my-hy2-pass" {
		t.Errorf("expected password 'my-hy2-pass', got %v", proxy["password"])
	}
	if proxy["server"] != "hy2.example.com" {
		t.Errorf("expected server 'hy2.example.com', got %v", proxy["server"])
	}
	if proxy["port"] != 443 {
		t.Errorf("expected port 443, got %v", proxy["port"])
	}
	if proxy["udp"] != true {
		t.Errorf("expected udp true, got %v", proxy["udp"])
	}
	if alpn, ok := proxy["alpn"].([]string); !ok || len(alpn) != 1 || alpn[0] != "h3" {
		t.Errorf("expected alpn [h3], got %v", proxy["alpn"])
	}
	if _, ok := proxy["network"]; ok {
		t.Errorf("hysteria2 should not have network property, got %v", proxy["network"])
	}
}

func TestMihomoHysteria2SalamanderObfs(t *testing.T) {
	inboundType := "hysteria2"
	finalMaskJSON := `{
		"udp": [
			{
				"type": "salamander",
				"settings": {
					"password": "salamander-secret"
				}
			}
		]
	}`
	host := SubscriptionHost{
		Remark:      "Hy2-Salamander",
		Address:     "hy2.example.com",
		Port:        443,
		InboundType: &inboundType,
		FinalMask:   &finalMaskJSON,
	}
	user := SubscriptionUser{
		Hysteria2Password: "my-hy2-pass",
	}

	proxy := buildMihomoProxy(host, user)
	if proxy == nil {
		t.Fatal("expected proxy, got nil")
	}

	if proxy["obfs"] != "salamander" {
		t.Errorf("expected obfs 'salamander', got %v", proxy["obfs"])
	}
	if proxy["obfs-password"] != "salamander-secret" {
		t.Errorf("expected obfs-password 'salamander-secret', got %v", proxy["obfs-password"])
	}
	if _, ok := proxy["obfs-min-packet-size"]; ok {
		t.Errorf("salamander should not have obfs-min-packet-size, got %v", proxy["obfs-min-packet-size"])
	}
	if _, ok := proxy["obfs-max-packet-size"]; ok {
		t.Errorf("salamander should not have obfs-max-packet-size, got %v", proxy["obfs-max-packet-size"])
	}
}

func TestMihomoHysteria2GeckoObfsRange(t *testing.T) {
	inboundType := "hysteria2"
	finalMaskJSON := `{
		"udp": [
			{
				"type": "salamander",
				"settings": {
					"password": "gecko-secret",
					"packetSize": "100-200"
				}
			}
		]
	}`
	host := SubscriptionHost{
		Remark:      "Hy2-Gecko",
		Address:     "hy2.example.com",
		Port:        443,
		InboundType: &inboundType,
		FinalMask:   &finalMaskJSON,
	}
	user := SubscriptionUser{
		Hysteria2Password: "my-hy2-pass",
	}

	proxy := buildMihomoProxy(host, user)
	if proxy == nil {
		t.Fatal("expected proxy, got nil")
	}

	if proxy["obfs"] != "gecko" {
		t.Errorf("expected obfs 'gecko', got %v", proxy["obfs"])
	}
	if proxy["obfs-password"] != "gecko-secret" {
		t.Errorf("expected obfs-password 'gecko-secret', got %v", proxy["obfs-password"])
	}
	if proxy["obfs-min-packet-size"] != 100 {
		t.Errorf("expected obfs-min-packet-size 100, got %v", proxy["obfs-min-packet-size"])
	}
	if proxy["obfs-max-packet-size"] != 200 {
		t.Errorf("expected obfs-max-packet-size 200, got %v", proxy["obfs-max-packet-size"])
	}
}

func TestMihomoHysteria2GeckoObfsReversedRange(t *testing.T) {
	inboundType := "hysteria2"
	finalMaskJSON := `{
		"udp": [
			{
				"type": "salamander",
				"settings": {
					"password": "gecko-secret",
					"packetSize": "500-200"
				}
			}
		]
	}`
	host := SubscriptionHost{
		Remark:      "Hy2-Gecko-Rev",
		Address:     "hy2.example.com",
		Port:        443,
		InboundType: &inboundType,
		FinalMask:   &finalMaskJSON,
	}
	user := SubscriptionUser{
		Hysteria2Password: "my-hy2-pass",
	}

	proxy := buildMihomoProxy(host, user)
	if proxy == nil {
		t.Fatal("expected proxy, got nil")
	}

	if proxy["obfs"] != "gecko" {
		t.Errorf("expected obfs 'gecko', got %v", proxy["obfs"])
	}
	if proxy["obfs-min-packet-size"] != 200 {
		t.Errorf("expected obfs-min-packet-size 200, got %v", proxy["obfs-min-packet-size"])
	}
	if proxy["obfs-max-packet-size"] != 500 {
		t.Errorf("expected obfs-max-packet-size 500, got %v", proxy["obfs-max-packet-size"])
	}
}

func TestMihomoHysteria2GeckoObfsSingleInt(t *testing.T) {
	inboundType := "hysteria2"
	finalMaskJSON := `{
		"udp": [
			{
				"type": "salamander",
				"settings": {
					"password": "gecko-secret",
					"packetSize": 150
				}
			}
		]
	}`
	host := SubscriptionHost{
		Remark:      "Hy2-Gecko-Int",
		Address:     "hy2.example.com",
		Port:        443,
		InboundType: &inboundType,
		FinalMask:   &finalMaskJSON,
	}
	user := SubscriptionUser{
		Hysteria2Password: "my-hy2-pass",
	}

	proxy := buildMihomoProxy(host, user)
	if proxy == nil {
		t.Fatal("expected proxy, got nil")
	}

	if proxy["obfs"] != "gecko" {
		t.Errorf("expected obfs 'gecko', got %v", proxy["obfs"])
	}
	if proxy["obfs-min-packet-size"] != 150 {
		t.Errorf("expected obfs-min-packet-size 150, got %v", proxy["obfs-min-packet-size"])
	}
	if proxy["obfs-max-packet-size"] != 150 {
		t.Errorf("expected obfs-max-packet-size 150, got %v", proxy["obfs-max-packet-size"])
	}
}

func TestMihomoHysteria2QuicParams(t *testing.T) {
	inboundType := "hysteria2"
	finalMaskJSON := `{
		"quicParams": {
			"brutalUp": 100,
			"brutalDown": 50,
			"udpHop": {
				"ports": "10000-20000",
				"interval": 30
			},
			"bbrProfile": "bbr"
		}
	}`
	host := SubscriptionHost{
		Remark:      "Hy2-Quic",
		Address:     "hy2.example.com",
		Port:        443,
		InboundType: &inboundType,
		FinalMask:   &finalMaskJSON,
	}
	user := SubscriptionUser{
		Hysteria2Password: "my-hy2-pass",
	}

	proxy := buildMihomoProxy(host, user)
	if proxy == nil {
		t.Fatal("expected proxy, got nil")
	}

	if proxy["up"] != "100" {
		t.Errorf("expected up '100', got %v", proxy["up"])
	}
	if proxy["down"] != "50" {
		t.Errorf("expected down '50', got %v", proxy["down"])
	}
	if proxy["ports"] != "10000-20000" {
		t.Errorf("expected ports '10000-20000', got %v", proxy["ports"])
	}
	if proxy["hop-interval"] != "30" {
		t.Errorf("expected hop-interval '30', got %v", proxy["hop-interval"])
	}
	if proxy["bbr-profile"] != "bbr" {
		t.Errorf("expected bbr-profile 'bbr', got %v", proxy["bbr-profile"])
	}
}

func TestMihomoHysteria2TLSParams(t *testing.T) {
	inboundType := "hysteria2"
	sni := "custom-sni.com"
	pinned := "sha256/pinned-cert"
	fingerprint := "chrome"
	alpn := "h3,h2"
	host := SubscriptionHost{
		Remark:               "Hy2-TLS",
		Address:              "hy2.example.com",
		Port:                 443,
		InboundType:          &inboundType,
		SNI:                  &sni,
		PinnedPeerCertSha256: &pinned,
		Fingerprint:          &fingerprint,
		ALPN:                 &alpn,
	}
	user := SubscriptionUser{
		Hysteria2Password: "my-hy2-pass",
	}

	proxy := buildMihomoProxy(host, user)
	if proxy == nil {
		t.Fatal("expected proxy, got nil")
	}

	if proxy["sni"] != "custom-sni.com" {
		t.Errorf("expected sni 'custom-sni.com', got %v", proxy["sni"])
	}
	if proxy["skip-cert-verify"] != true {
		t.Errorf("expected skip-cert-verify true, got %v", proxy["skip-cert-verify"])
	}
	if proxy["client-fingerprint"] != "chrome" {
		t.Errorf("expected client-fingerprint 'chrome', got %v", proxy["client-fingerprint"])
	}
	alpnList, ok := proxy["alpn"].([]string)
	if !ok || len(alpnList) != 2 || alpnList[0] != "h3" || alpnList[1] != "h2" {
		t.Errorf("expected alpn [h3, h2], got %v", proxy["alpn"])
	}
}

func TestMihomoHysteria2FullGeneration(t *testing.T) {
	inboundType := "hysteria2"
	finalMaskJSON := `{
		"udp": [
			{
				"type": "salamander",
				"settings": {
					"password": "gecko-pass",
					"packetSize": "100-200"
				}
			}
		],
		"quicParams": {
			"brutalUp": 100,
			"brutalDown": 50
		}
	}`
	hosts := []SubscriptionHost{
		{
			Remark:      "Hy2-Production",
			Address:     "hy2.example.com",
			Port:        443,
			InboundType: &inboundType,
			FinalMask:   &finalMaskJSON,
		},
	}
	user := SubscriptionUser{
		Hysteria2Password: "my-hy2-pass",
	}

	gen := NewMihomoGenerator(nil)
	output, err := gen.Generate(nil, user, hosts, SubscriptionSettingsParsed{})
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	var parsed map[string]interface{}
	if err := yaml.Unmarshal([]byte(output), &parsed); err != nil {
		t.Fatalf("Unmarshal generated YAML failed: %v\nOutput:\n%s", err, output)
	}

	proxies, ok := parsed["proxies"].([]interface{})
	if !ok || len(proxies) != 1 {
		t.Fatalf("expected 1 proxy in output, got %v", proxies)
	}

	p := proxies[0].(map[string]interface{})
	if p["name"] != "Hy2-Production" {
		t.Errorf("expected name 'Hy2-Production', got %v", p["name"])
	}
	if p["type"] != "hysteria2" {
		t.Errorf("expected type 'hysteria2', got %v", p["type"])
	}
	if p["password"] != "my-hy2-pass" {
		t.Errorf("expected password 'my-hy2-pass', got %v", p["password"])
	}
	if p["obfs"] != "gecko" {
		t.Errorf("expected obfs 'gecko', got %v", p["obfs"])
	}
	if p["obfs-password"] != "gecko-pass" {
		t.Errorf("expected obfs-password 'gecko-pass', got %v", p["obfs-password"])
	}
	if p["obfs-min-packet-size"] != 100 {
		t.Errorf("expected obfs-min-packet-size 100, got %v", p["obfs-min-packet-size"])
	}
	if p["obfs-max-packet-size"] != 200 {
		t.Errorf("expected obfs-max-packet-size 200, got %v", p["obfs-max-packet-size"])
	}
	if p["up"] != "100" {
		t.Errorf("expected up '100', got %v", p["up"])
	}
	if p["down"] != "50" {
		t.Errorf("expected down '50', got %v", p["down"])
	}
}

func TestMihomoHysteria2FromInboundRaw(t *testing.T) {
	inboundType := "hysteria2"
	inboundRaw := json.RawMessage(`{
		"type": "hysteria2",
		"finalMask": {
			"udp": [
				{
					"type": "salamander",
					"settings": {
						"password": "raw-gecko-pass",
						"packetSize": "120-250"
					}
				}
			]
		}
	}`)
	hosts := []SubscriptionHost{
		{
			Remark:      "Hy2-From-Raw",
			Address:     "hy2-raw.example.com",
			Port:        443,
			InboundType: &inboundType,
			InboundRaw:  inboundRaw,
		},
	}
	user := SubscriptionUser{
		Hysteria2Password: "raw-hy2-pass",
	}

	proxy := buildMihomoProxy(hosts[0], user)
	if proxy == nil {
		t.Fatal("expected proxy, got nil")
	}

	if proxy["obfs"] != "gecko" {
		t.Errorf("expected obfs 'gecko', got %v", proxy["obfs"])
	}
	if proxy["obfs-password"] != "raw-gecko-pass" {
		t.Errorf("expected obfs-password 'raw-gecko-pass', got %v", proxy["obfs-password"])
	}
	if proxy["obfs-min-packet-size"] != 120 {
		t.Errorf("expected obfs-min-packet-size 120, got %v", proxy["obfs-min-packet-size"])
	}
	if proxy["obfs-max-packet-size"] != 250 {
		t.Errorf("expected obfs-max-packet-size 250, got %v", proxy["obfs-max-packet-size"])
	}
}

func TestMihomoVlessReality(t *testing.T) {
	inboundType := "vless"
	inboundRaw := json.RawMessage(`{
		"type": "vless",
		"streamSettings": {
			"network": "tcp",
			"security": "reality",
			"realitySettings": {
				"serverNames": ["reality.example.com"],
				"publicKey": "my-public-key-12345",
				"shortIds": ["0123456789abcdef"],
				"fingerprint": "chrome"
			}
		}
	}`)
	host := SubscriptionHost{
		Remark:      "VLESS-Reality-Node",
		Address:     "node.example.com",
		Port:        443,
		InboundType: &inboundType,
		InboundRaw:  inboundRaw,
	}
	user := SubscriptionUser{
		VlessUUID: "11111111-2222-3333-4444-555555555555",
	}

	proxy := buildMihomoProxy(host, user)
	if proxy == nil {
		t.Fatal("expected proxy, got nil")
	}

	if proxy["type"] != "vless" {
		t.Errorf("expected type 'vless', got %v", proxy["type"])
	}
	if proxy["uuid"] != "11111111-2222-3333-4444-555555555555" {
		t.Errorf("expected uuid, got %v", proxy["uuid"])
	}
	if proxy["tls"] != true {
		t.Errorf("expected tls true, got %v", proxy["tls"])
	}
	if proxy["servername"] != "reality.example.com" {
		t.Errorf("expected servername 'reality.example.com', got %v", proxy["servername"])
	}
	if proxy["flow"] != "xtls-rprx-vision" {
		t.Errorf("expected flow 'xtls-rprx-vision', got %v", proxy["flow"])
	}
	if proxy["client-fingerprint"] != "chrome" {
		t.Errorf("expected client-fingerprint 'chrome', got %v", proxy["client-fingerprint"])
	}
	realityOpts, ok := proxy["reality-opts"].(map[string]interface{})
	if !ok || realityOpts == nil {
		t.Fatalf("expected reality-opts map, got %v", proxy["reality-opts"])
	}
	if realityOpts["public-key"] != "my-public-key-12345" {
		t.Errorf("expected public-key, got %v", realityOpts["public-key"])
	}
	if realityOpts["short-id"] != "0123456789abcdef" {
		t.Errorf("expected short-id, got %v", realityOpts["short-id"])
	}
}

func TestMihomoExodusTemplateDirective(t *testing.T) {
	templateYAML := []byte(`
proxy-groups:
  - name: PROXY
    type: select
    proxies:
      - DIRECT
  - name: MANUAL
    type: select
    exodus:
      include-proxies: false
    proxies:
      - DIRECT
`)
	inboundType := "vless"
	hosts := []SubscriptionHost{
		{
			Remark:      "Node 1",
			Address:     "1.1.1.1",
			Port:        443,
			InboundType: &inboundType,
		},
	}
	user := SubscriptionUser{
		VlessUUID: "11111111-2222-3333-4444-555555555555",
	}

	result, err := generateYAMLConfig(templateYAML, hosts, user)
	if err != nil {
		t.Fatalf("generateYAMLConfig error: %v", err)
	}

	var parsed struct {
		ProxyGroups []map[string]interface{} `yaml:"proxy-groups"`
	}
	if err := yaml.Unmarshal([]byte(result), &parsed); err != nil {
		t.Fatalf("failed to unmarshal output: %v", err)
	}

	if len(parsed.ProxyGroups) != 2 {
		t.Fatalf("expected 2 proxy groups, got %d", len(parsed.ProxyGroups))
	}

	// First group should have Node 1
	p0 := parsed.ProxyGroups[0]["proxies"].([]interface{})
	if len(p0) != 2 || p0[1] != "Node 1" {
		t.Errorf("expected PROXY group to include Node 1, got %v", p0)
	}

	// Second group has exodus.include-proxies: false -> should NOT have Node 1
	p1 := parsed.ProxyGroups[1]["proxies"].([]interface{})
	if len(p1) != 1 || p1[0] != "DIRECT" {
		t.Errorf("expected MANUAL group to remain untouched, got %v", p1)
	}

	// The 'exodus' key must be deleted
	if _, ok := parsed.ProxyGroups[1]["exodus"]; ok {
		t.Errorf("expected exodus key to be removed from YAML")
	}
}

func TestXrayVlessRealityLinkAndOutbound(t *testing.T) {
	inboundType := "vless"
	inboundRaw := json.RawMessage(`{
		"type": "vless",
		"streamSettings": {
			"network": "tcp",
			"security": "reality",
			"realitySettings": {
				"serverNames": ["reality.example.com"],
				"publicKey": "my-public-key-12345",
				"shortIds": ["0123456789abcdef"],
				"fingerprint": "chrome"
			}
		}
	}`)
	muxParams := `{"smux": {"enabled": true, "protocol": "h2mux"}}`
	host := SubscriptionHost{
		Remark:      "VLESS-Reality-Link",
		Address:     "node.example.com",
		Port:        443,
		InboundType: &inboundType,
		InboundRaw:  inboundRaw,
		MuxParams:   &muxParams,
	}
	user := SubscriptionUser{
		VlessUUID: "11111111-2222-3333-4444-555555555555",
	}

	link, proto := buildHostLink(host, user)
	if proto != "vless" {
		t.Errorf("expected proto vless, got %s", proto)
	}
	if !strings.Contains(link, "security=reality") {
		t.Errorf("link missing security=reality: %s", link)
	}
	if !strings.Contains(link, "pbk=my-public-key-12345") {
		t.Errorf("link missing pbk: %s", link)
	}
	if !strings.Contains(link, "sid=0123456789abcdef") {
		t.Errorf("link missing sid: %s", link)
	}
	if !strings.Contains(link, "flow=xtls-rprx-vision") {
		t.Errorf("link missing flow: %s", link)
	}
	if !strings.Contains(link, "fp=chrome") {
		t.Errorf("link missing fp: %s", link)
	}

	// Xray Outbound with smux only -> mux must be omitted
	outbound := buildXrayOutbound(host, user)
	if outbound == nil {
		t.Fatal("expected outbound, got nil")
	}
	if _, hasMux := outbound["mux"]; hasMux {
		t.Errorf("expected outbound.mux to be omitted when only smux is present, got %v", outbound["mux"])
	}

	// Xray Outbound with native Xray mux fields
	nativeMux := `{"enabled": true, "concurrency": 8, "smux": {"enabled": true}}`
	host.MuxParams = &nativeMux
	outboundWithMux := buildXrayOutbound(host, user)
	if muxObj, ok := outboundWithMux["mux"].(map[string]interface{}); !ok {
		t.Errorf("expected outbound.mux to be present, got %v", outboundWithMux["mux"])
	} else {
		if _, hasSmux := muxObj["smux"]; hasSmux {
			t.Errorf("expected smux key to be stripped from outbound.mux")
		}
		if muxObj["enabled"] != true || muxObj["concurrency"] != float64(8) {
			t.Errorf("expected native mux fields, got %v", muxObj)
		}
	}
}
