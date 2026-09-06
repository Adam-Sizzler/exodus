package subscription

import (
	"encoding/json"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestApplyShuffle(t *testing.T) {
	t.Run("no shuffleHost retains original order", func(t *testing.T) {
		hosts := []SubscriptionHost{
			{UUID: "h1", Remark: "Host 1", ShuffleHost: false},
			{UUID: "h2", Remark: "Host 2", ShuffleHost: false},
			{UUID: "h3", Remark: "Host 3", ShuffleHost: false},
		}
		result := applyShuffle(hosts)
		if len(result) != 3 {
			t.Fatalf("expected 3 hosts, got %d", len(result))
		}
		if result[0].UUID != "h1" || result[1].UUID != "h2" || result[2].UUID != "h3" {
			t.Fatalf("order changed when no shuffleHost was set: %+v", result)
		}
	})

	t.Run("shuffleHost puts shuffled hosts at the beginning", func(t *testing.T) {
		hosts := []SubscriptionHost{
			{UUID: "h1", Remark: "Host 1", ShuffleHost: false},
			{UUID: "h2", Remark: "Host 2", ShuffleHost: true},
			{UUID: "h3", Remark: "Host 3", ShuffleHost: false},
			{UUID: "h4", Remark: "Host 4", ShuffleHost: true},
			{UUID: "h5", Remark: "Host 5", ShuffleHost: false},
		}
		result := applyShuffle(hosts)
		if len(result) != 5 {
			t.Fatalf("expected 5 hosts, got %d", len(result))
		}

		// First 2 must be h2 and h4 in some order
		shuffledSet := map[string]bool{result[0].UUID: true, result[1].UUID: true}
		if !shuffledSet["h2"] || !shuffledSet["h4"] {
			t.Fatalf("expected first two to be h2 and h4, got %s and %s", result[0].UUID, result[1].UUID)
		}

		// Remaining 3 must preserve relative order: h1, h3, h5
		if result[2].UUID != "h1" || result[3].UUID != "h3" || result[4].UUID != "h5" {
			t.Fatalf("expected remaining hosts to retain relative order (h1, h3, h5), got [%s, %s, %s]",
				result[2].UUID, result[3].UUID, result[4].UUID)
		}
	})
}

func TestClashGeneratorParity(t *testing.T) {
	vlessProto := "vless"
	trojanProto := "trojan"
	ssProto := "shadowsocks"
	hy2Proto := "hysteria2"
	tlsSec := "tls"
	tcpNet := "tcp"
	xhttpNet := "xhttp"
	pinnedCert := "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"

	user := SubscriptionUser{
		ID:             1,
		UUID:           "u1",
		ShortUUID:      "usr123",
		Username:       "testuser",
		VlessUUID:      "11111111-1111-1111-1111-111111111111",
		TrojanPassword: "trojan-password-123",
		SSPassword:     "ss-password-123",
	}

	hosts := []SubscriptionHost{
		{
			UUID:            "h-vless",
			Remark:          "VLESS Host",
			Address:         "vless.example.com",
			Port:            443,
			InboundType:     &vlessProto,
			InboundSecurity: &tlsSec,
			InboundNetwork:  &tcpNet,
		},
		{
			UUID:            "h-hy2",
			Remark:          "Hysteria2 Host",
			Address:         "hy2.example.com",
			Port:            443,
			InboundType:     &hy2Proto,
			InboundSecurity: &tlsSec,
		},
		{
			UUID:                 "h-trojan",
			Remark:               "Trojan Host",
			Address:              "trojan.example.com",
			Port:                 443,
			InboundType:          &trojanProto,
			InboundSecurity:      &tlsSec,
			InboundNetwork:       &tcpNet,
			PinnedPeerCertSha256: &pinnedCert,
		},
		{
			UUID:                         "h-trojan-excluded",
			Remark:                       "Trojan Excluded",
			Address:                      "trojan-ex.example.com",
			Port:                         443,
			InboundType:                  &trojanProto,
			InboundSecurity:              &tlsSec,
			InboundNetwork:               &tcpNet,
			ExcludeFromSubscriptionTypes: []string{"CLASH"},
		},
		{
			UUID:            "h-ss",
			Remark:          "SS Host",
			Address:         "ss.example.com",
			Port:            8388,
			InboundType:     &ssProto,
			InboundNetwork:  &tcpNet,
			InboundRaw:      json.RawMessage(`{"settings":{"method":"aes-256-gcm"}}`),
		},
		{
			UUID:            "h-trojan-xhttp",
			Remark:          "Trojan xHTTP",
			Address:         "xhttp.example.com",
			Port:            443,
			InboundType:     &trojanProto,
			InboundSecurity: &tlsSec,
			InboundNetwork:  &xhttpNet,
		},
	}

	template := []byte(`
proxies: []
proxy-groups:
  - name: PROXY
    type: select
    proxies: []
rules: []
`)

	clashGen := NewClashGenerator(nil)
	out, err := clashGen.Generate(template, user, hosts, SubscriptionSettingsParsed{})
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	var parsed map[string]interface{}
	if err := yaml.Unmarshal([]byte(out), &parsed); err != nil {
		t.Fatalf("Unmarshal YAML failed: %v", err)
	}

	proxies, ok := parsed["proxies"].([]interface{})
	if !ok {
		t.Fatalf("expected proxies array in clash config")
	}

	// Should only contain "Trojan Host" and "SS Host"
	// VLESS -> unsupported
	// Hysteria2 -> unsupported
	// Trojan Excluded -> excludeFromSubscriptionTypes CLASH
	// Trojan xHTTP -> unsupported transport
	if len(proxies) != 2 {
		t.Fatalf("expected exactly 2 proxies in Clash, got %d: %+v", len(proxies), proxies)
	}

	proxy0 := proxies[0].(map[string]interface{})
	proxy1 := proxies[1].(map[string]interface{})

	names := []string{proxy0["name"].(string), proxy1["name"].(string)}
	if names[0] != "Trojan Host" || names[1] != "SS Host" {
		t.Fatalf("unexpected proxy names: %+v", names)
	}

	// Verify Trojan host has skip-cert-verify: true due to PinnedPeerCertSha256
	if skip, ok := proxy0["skip-cert-verify"].(bool); !ok || !skip {
		t.Fatalf("expected skip-cert-verify: true for Trojan Host with pinned cert, got %+v", proxy0["skip-cert-verify"])
	}
}

func TestStashAndMihomoIsolation(t *testing.T) {
	vlessProto := "vless"
	tlsSec := "tls"
	tcpNet := "tcp"
	xhttpNet := "xhttp"
	desc := "My Server Description"

	user := SubscriptionUser{
		ID:        1,
		UUID:      "u1",
		ShortUUID: "usr123",
		Username:  "testuser",
		VlessUUID: "11111111-1111-1111-1111-111111111111",
	}

	hosts := []SubscriptionHost{
		{
			UUID:              "h-visible",
			Remark:            "Visible Host",
			Address:           "vis.example.com",
			Port:              443,
			InboundType:       &vlessProto,
			InboundSecurity:   &tlsSec,
			InboundNetwork:    &tcpNet,
			ServerDescription: &desc,
			IsHidden:          false,
		},
		{
			UUID:              "h-hidden",
			Remark:            "Hidden Host",
			Address:           "hid.example.com",
			Port:              443,
			InboundType:       &vlessProto,
			InboundSecurity:   &tlsSec,
			InboundNetwork:    &tcpNet,
			ServerDescription: &desc,
			IsHidden:          true,
		},
		{
			UUID:            "h-xhttp",
			Remark:          "xHTTP Host",
			Address:         "xhttp.example.com",
			Port:            443,
			InboundType:     &vlessProto,
			InboundSecurity: &tlsSec,
			InboundNetwork:  &xhttpNet,
			IsHidden:        false,
		},
	}

	template := []byte(`
proxies: []
proxy-groups:
  - name: PROXY
    type: select
    proxies: []
rules: []
`)

	mihomoGen := NewMihomoGenerator(nil)

	t.Run("STASH includes hidden hosts, skips xhttp, suppresses serverDescription", func(t *testing.T) {
		settings := SubscriptionSettingsParsed{IsExtendedClient: true}
		out, err := mihomoGen.Generate(template, user, hosts, settings, "STASH")
		if err != nil {
			t.Fatalf("Generate failed: %v", err)
		}

		var parsed map[string]interface{}
		if err := yaml.Unmarshal([]byte(out), &parsed); err != nil {
			t.Fatalf("Unmarshal YAML failed: %v", err)
		}

		proxies := parsed["proxies"].([]interface{})
		// Stash must include "Visible Host" and "Hidden Host", but skip "xHTTP Host"
		if len(proxies) != 2 {
			t.Fatalf("expected 2 proxies in Stash, got %d", len(proxies))
		}

		p0 := proxies[0].(map[string]interface{})
		p1 := proxies[1].(map[string]interface{})
		if p0["name"] != "Visible Host" || p1["name"] != "Hidden Host" {
			t.Fatalf("unexpected proxy names: %v, %v", p0["name"], p1["name"])
		}

		// serverDescription must NOT be present in Stash
		if _, exists := p0["serverDescription"]; exists {
			t.Fatalf("serverDescription must not be present in Stash output")
		}
	})

	t.Run("MIHOMO excludes hidden hosts by default, supports xhttp and serverDescription", func(t *testing.T) {
		settings := SubscriptionSettingsParsed{IsExtendedClient: true}
		out, err := mihomoGen.Generate(template, user, hosts, settings, "MIHOMO")
		if err != nil {
			t.Fatalf("Generate failed: %v", err)
		}

		var parsed map[string]interface{}
		if err := yaml.Unmarshal([]byte(out), &parsed); err != nil {
			t.Fatalf("Unmarshal YAML failed: %v", err)
		}

		proxies := parsed["proxies"].([]interface{})
		// Mihomo must include "Visible Host" and "xHTTP Host", but skip "Hidden Host"
		if len(proxies) != 2 {
			t.Fatalf("expected 2 proxies in Mihomo, got %d", len(proxies))
		}

		p0 := proxies[0].(map[string]interface{})
		if p0["name"] != "Visible Host" {
			t.Fatalf("expected Visible Host, got %v", p0["name"])
		}

		// serverDescription MUST be present for extended client in Mihomo
		if p0["serverDescription"] != "My Server Description" {
			t.Fatalf("expected serverDescription 'My Server Description', got %v", p0["serverDescription"])
		}
	})
}

func TestSingboxInsecurePinnedCert(t *testing.T) {
	vlessProto := "vless"
	tlsSec := "tls"
	tcpNet := "tcp"
	pinnedCert := "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"

	user := SubscriptionUser{
		ID:        1,
		UUID:      "u1",
		ShortUUID: "usr123",
		Username:  "testuser",
		VlessUUID: "11111111-1111-1111-1111-111111111111",
	}

	hosts := []SubscriptionHost{
		{
			UUID:                 "h-vless",
			Remark:               "VLESS Pinned",
			Address:              "vless.example.com",
			Port:                 443,
			InboundType:          &vlessProto,
			InboundSecurity:      &tlsSec,
			InboundNetwork:       &tcpNet,
			PinnedPeerCertSha256: &pinnedCert,
		},
	}

	template := []byte(`{"outbounds": []}`)
	singboxGen := NewSingboxGenerator(nil)
	out, err := singboxGen.Generate(template, user, hosts, SubscriptionSettingsParsed{})
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	var parsed map[string]interface{}
	if err := json.Unmarshal([]byte(out), &parsed); err != nil {
		t.Fatalf("Unmarshal JSON failed: %v", err)
	}

	outbounds := parsed["outbounds"].([]interface{})
	if len(outbounds) == 0 {
		t.Fatalf("expected outbounds")
	}

	ob := outbounds[0].(map[string]interface{})
	tlsCfg, ok := ob["tls"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected tls config in outbound")
	}

	if insecure, ok := tlsCfg["insecure"].(bool); !ok || !insecure {
		t.Fatalf("expected tls.insecure: true when PinnedPeerCertSha256 is present, got %v", tlsCfg["insecure"])
	}
}

func TestXrayLinksPinnedCertAndVerifyPeer(t *testing.T) {
	vlessProto := "vless"
	tlsSec := "tls"
	tcpNet := "tcp"
	pinnedCert := "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
	verifyPeer := "custom-peer.com"

	user := SubscriptionUser{
		ID:        1,
		UUID:      "u1",
		ShortUUID: "usr123",
		Username:  "testuser",
		VlessUUID: "11111111-1111-1111-1111-111111111111",
	}

	hosts := []SubscriptionHost{
		{
			UUID:                 "h-vless",
			Remark:               "VLESS Pinned",
			Address:              "vless.example.com",
			Port:                 443,
			InboundType:          &vlessProto,
			InboundSecurity:      &tlsSec,
			InboundNetwork:       &tcpNet,
			PinnedPeerCertSha256: &pinnedCert,
			VerifyPeerCertByName: &verifyPeer,
		},
	}

	xrayGen := NewXrayGenerator(nil)
	links, err := xrayGen.GenerateLinks(user, hosts, SubscriptionSettingsParsed{})
	if err != nil {
		t.Fatalf("GenerateLinks failed: %v", err)
	}

	if len(links) != 1 {
		t.Fatalf("expected 1 link, got %d", len(links))
	}

	link := links[0]
	if !strings.Contains(link, "pinSHA256="+pinnedCert) && !strings.Contains(link, "pcs="+pinnedCert) {
		t.Fatalf("link missing pinSHA256 or pcs query param: %s", link)
	}
	if !strings.Contains(link, "vcn="+verifyPeer) {
		t.Fatalf("link missing vcn query param: %s", link)
	}
}
