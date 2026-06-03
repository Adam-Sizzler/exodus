package users

import "testing"

func TestBuildInboundUsersUsesSingboxProtocolCredentials(t *testing.T) {
	users := []inboundUserCredentials{
		{
			Username:       "alice",
			VLESSUUID:      "9f76f8d8-daf1-4db6-b045-0987cd5e09a2",
			TrojanPassword: "trojan-secret",
			Hysteria2Pass:  "hy-secret",
		},
	}

	tests := []struct {
		name      string
		protocol  string
		wantKey   string
		wantValue any
	}{
		{name: "vmess uses uuid", protocol: "vmess", wantKey: "uuid", wantValue: users[0].VLESSUUID},
		{name: "hysteria uses auth_str", protocol: "hysteria", wantKey: "auth_str", wantValue: "hy-secret"},
		{name: "hysteria2 uses password", protocol: "hysteria2", wantKey: "password", wantValue: "hy-secret"},
		{name: "tuic uses password", protocol: "tuic", wantKey: "password", wantValue: "trojan-secret"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			items := buildInboundUsers(tc.protocol, users)
			if len(items) != 1 {
				t.Fatalf("got %d users, want 1", len(items))
			}
			item, ok := items[0].(map[string]any)
			if !ok {
				t.Fatalf("unexpected item type %T", items[0])
			}
			if got := item[tc.wantKey]; got != tc.wantValue {
				t.Fatalf("field %s got %#v, want %#v", tc.wantKey, got, tc.wantValue)
			}
		})
	}

	tuic := buildInboundUsers("tuic", users)[0].(map[string]any)
	if got := tuic["uuid"]; got != users[0].VLESSUUID {
		t.Fatalf("tuic uuid got %#v, want %#v", got, users[0].VLESSUUID)
	}
}
