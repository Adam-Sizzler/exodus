package subscription

import (
	"encoding/json"
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
