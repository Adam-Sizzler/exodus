package configprofiles

import (
	"encoding/json"
	"testing"
)

func TestValidateConfigStructure(t *testing.T) {
	tests := []struct {
		name        string
		configJSON  string
		wantErr     bool
		errContains string
	}{
		{
			name: "valid config with unique inbounds and outbounds",
			configJSON: `{
				"inbounds": [
					{"type": "shadowsocks", "tag": "ss-in-1", "listen_port": 2080},
					{"type": "trojan", "tag": "trojan-1", "listen_port": 2443}
				],
				"outbounds": [
					{"type": "direct", "tag": "direct"}
				]
			}`,
			wantErr: false,
		},
		{
			name: "invalid config with duplicate inbound tags in same profile",
			configJSON: `{
				"inbounds": [
					{"type": "shadowsocks", "tag": "ss-in-1", "listen_port": 2080},
					{"type": "trojan", "tag": "ss-in-1", "listen_port": 2443}
				],
				"outbounds": [
					{"type": "direct", "tag": "direct"}
				]
			}`,
			wantErr:     true,
			errContains: "duplicate inbound tag \"ss-in-1\" found. All inbound tags must be unique",
		},
		{
			name: "invalid config with empty inbound tag",
			configJSON: `{
				"inbounds": [
					{"type": "shadowsocks", "tag": "   ", "listen_port": 2080}
				],
				"outbounds": [
					{"type": "direct", "tag": "direct"}
				]
			}`,
			wantErr:     true,
			errContains: "all inbounds must have a non-empty tag",
		},
		{
			name: "invalid config with comma in inbound tag",
			configJSON: `{
				"inbounds": [
					{"type": "shadowsocks", "tag": "ss,in,1", "listen_port": 2080}
				],
				"outbounds": [
					{"type": "direct", "tag": "direct"}
				]
			}`,
			wantErr:     true,
			errContains: "character ',' is not allowed in inbound tag",
		},
		{
			name: "valid config without outbounds array",
			configJSON: `{
				"log": {"level": "info"}
			}`,
			wantErr: false,
		},
		{
			name: "valid config with empty outbounds array",
			configJSON: `{
				"outbounds": []
			}`,
			wantErr: false,
		},
		{
			name: "valid config with valid SS 2022 32-byte password",
			configJSON: `{
				"inbounds": [
					{
						"type": "shadowsocks",
						"tag": "ss-2022",
						"method": "2022-blake3-aes-256-gcm",
						"password": "MDEyMzQ1Njc4OTAxMjM0NTY3ODkwMTIzNDU2Nzg5MDE="
					}
				],
				"outbounds": [{"type": "direct", "tag": "direct"}]
			}`,
			wantErr: false,
		},
		{
			name: "invalid config with SS 2022 missing password",
			configJSON: `{
				"inbounds": [
					{
						"type": "shadowsocks",
						"tag": "ss-2022",
						"method": "2022-blake3-aes-256-gcm"
					}
				],
				"outbounds": [{"type": "direct", "tag": "direct"}]
			}`,
			wantErr:     true,
			errContains: "Shadowsocks password is required for 2022-blake3-* methods",
		},
		{
			name: "invalid config with SS 2022 invalid key length",
			configJSON: `{
				"inbounds": [
					{
						"type": "shadowsocks",
						"tag": "ss-2022",
						"method": "2022-blake3-aes-256-gcm",
						"password": "c2hvcnQ="
					}
				],
				"outbounds": [{"type": "direct", "tag": "direct"}]
			}`,
			wantErr:     true,
			errContains: "must be a base64 string that decodes to exactly 32 bytes",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateConfigStructure(json.RawMessage(tt.configJSON))
			if (err != nil) != tt.wantErr {
				t.Fatalf("validateConfigStructure() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr && tt.errContains != "" {
				if err == nil || !testingStringContains(err.Error(), tt.errContains) {
					t.Fatalf("expected error containing %q, got %v", tt.errContains, err)
				}
			}
		})
	}
}

func testingStringContains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 || (len(s) > 0 && len(substr) > 0 && testingFindSubstr(s, substr)))
}

func testingFindSubstr(s, substr string) bool {
	for i := 0; i+len(substr) <= len(s); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
