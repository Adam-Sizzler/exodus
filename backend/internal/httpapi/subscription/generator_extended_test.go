package subscription

import (
	"encoding/base64"
	"encoding/json"
	"strings"
	"testing"
)

func TestExtendedClientsLinks(t *testing.T) {
	serverDesc := "Premium Node Tokyo"
	inboundType := "vless"
	inboundNet := "tcp"
	inboundSec := "none"

	host := SubscriptionHost{
		UUID:              "host-uuid-1",
		Address:           "1.2.3.4",
		Port:              443,
		Remark:            "Tokyo-01",
		InboundType:       &inboundType,
		InboundNetwork:    &inboundNet,
		InboundSecurity:   &inboundSec,
		ServerDescription: &serverDesc,
	}
	user := SubscriptionUser{
		UUID:      "user-uuid-1",
		ShortUUID: "u123",
		VlessUUID: "11111111-1111-1111-1111-111111111111",
	}

	t.Run("standard client does not receive serverDescription query param", func(t *testing.T) {
		links, _ := buildSubscriptionLinksExt([]SubscriptionHost{host}, user, false)
		if len(links) != 1 {
			t.Fatalf("expected 1 link, got %d", len(links))
		}
		if strings.Contains(links[0], "serverDescription") {
			t.Errorf("link contains serverDescription for non-extended client: %s", links[0])
		}
	})

	t.Run("extended client receives base64 serverDescription query param", func(t *testing.T) {
		links, _ := buildSubscriptionLinksExt([]SubscriptionHost{host}, user, true)
		if len(links) != 1 {
			t.Fatalf("expected 1 link, got %d", len(links))
		}
		expectedB64 := base64.StdEncoding.EncodeToString([]byte(serverDesc))
		expectedParam := "?serverDescription=" + expectedB64
		if !strings.Contains(links[0], expectedParam) {
			t.Errorf("link missing expected serverDescription param: %s; want %s", links[0], expectedParam)
		}
	})
}

func TestExtendedClientsMihomo(t *testing.T) {
	serverDesc := "Frankfurt Node 2"
	inboundType := "vless"
	host := SubscriptionHost{
		UUID:              "host-uuid-2",
		Address:           "5.6.7.8",
		Port:              8443,
		Remark:            "FRA-02",
		InboundType:       &inboundType,
		ServerDescription: &serverDesc,
	}
	user := SubscriptionUser{
		UUID:      "user-uuid-2",
		ShortUUID: "u456",
		VlessUUID: "22222222-2222-2222-2222-222222222222",
	}

	templateYAML := []byte("proxies: []\n")

	t.Run("standard client YAML omits serverDescription", func(t *testing.T) {
		out, err := generateYAMLConfigExt(templateYAML, []SubscriptionHost{host}, user, false, responseTypeMihomo)
		if err != nil {
			t.Fatalf("generateYAMLConfigExt failed: %v", err)
		}
		if strings.Contains(out, "serverDescription") {
			t.Errorf("YAML should not contain serverDescription for standard client: %s", out)
		}
	})

	t.Run("extended client YAML includes serverDescription", func(t *testing.T) {
		out, err := generateYAMLConfigExt(templateYAML, []SubscriptionHost{host}, user, true, responseTypeMihomo)
		if err != nil {
			t.Fatalf("generateYAMLConfigExt failed: %v", err)
		}
		if !strings.Contains(out, "serverDescription: Frankfurt Node 2") {
			t.Errorf("YAML should contain serverDescription: %s", out)
		}
	})
}

func TestExtendedClientsAndTemplatesXrayJSON(t *testing.T) {
	serverDesc := "Amsterdam Ultra"
	inboundType := "vless"
	customTplUUID := "custom-tpl-uuid-1"

	host := SubscriptionHost{
		UUID:                 "host-uuid-3",
		Address:              "9.10.11.12",
		Port:                 443,
		Remark:               "AMS-01",
		InboundType:          &inboundType,
		ServerDescription:    &serverDesc,
		XrayJSONTemplateUUID: &customTplUUID,
	}
	user := SubscriptionUser{
		UUID:      "user-uuid-3",
		ShortUUID: "u789",
		VlessUUID: "33333333-3333-3333-3333-333333333333",
	}

	defaultTemplate := []byte(`{"tag":"default-tpl","outbounds":[{"tag":"direct","protocol":"freedom"}]}`)
	customTemplate := []byte(`{"tag":"custom-host-tpl","outbounds":[{"tag":"custom-out","protocol":"freedom"}]}`)

	loader := func(uuid string) ([]byte, error) {
		if uuid == customTplUUID {
			return customTemplate, nil
		}
		return nil, nil
	}

	t.Run("XRAY_JSON outputs array of configs with per-host template and meta serverDescription", func(t *testing.T) {
		out, err := generateXrayJSONConfigExt(defaultTemplate, []SubscriptionHost{host}, user, true, false, loader)
		if err != nil {
			t.Fatalf("generateXrayJSONConfigExt failed: %v", err)
		}

		var configs []map[string]interface{}
		if err := json.Unmarshal([]byte(out), &configs); err != nil {
			t.Fatalf("unmarshal output as JSON array failed: %v; output: %s", err, out)
		}
		if len(configs) != 1 {
			t.Fatalf("expected 1 config in array, got %d", len(configs))
		}

		cfg := configs[0]
		if cfg["tag"] != "custom-host-tpl" {
			t.Errorf("expected host to use custom template, got tag: %v", cfg["tag"])
		}
		if cfg["remarks"] != "AMS-01" {
			t.Errorf("expected remarks AMS-01, got %v", cfg["remarks"])
		}

		meta, ok := cfg["meta"].(map[string]interface{})
		if !ok || meta["serverDescription"] != serverDesc {
			t.Errorf("expected meta.serverDescription = %q, got %v", serverDesc, cfg["meta"])
		}
	})

	t.Run("XRAY_JSON ignores host template when ignoreHostTemplate is true", func(t *testing.T) {
		out, err := generateXrayJSONConfigExt(defaultTemplate, []SubscriptionHost{host}, user, false, true, loader)
		if err != nil {
			t.Fatalf("generateXrayJSONConfigExt failed: %v", err)
		}

		var configs []map[string]interface{}
		if err := json.Unmarshal([]byte(out), &configs); err != nil {
			t.Fatalf("unmarshal output: %v", err)
		}

		cfg := configs[0]
		if cfg["tag"] != "default-tpl" {
			t.Errorf("expected host to use default template when ignoreHostTemplate=true, got: %v", cfg["tag"])
		}
		if cfg["meta"] != nil {
			t.Errorf("expected meta to be nil for standard client, got %v", cfg["meta"])
		}
	})
}
