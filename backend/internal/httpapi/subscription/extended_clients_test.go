package subscription

import (
	"testing"
)

func TestIsExtendedClient(t *testing.T) {
	tests := []struct {
		name       string
		userAgent  string
		additional []string
		expected   bool
	}{
		{
			name:      "Happ user-agent builtin",
			userAgent: "Happ/1.2.3 (iPhone; iOS 16.0)",
			expected:  true,
		},
		{
			name:      "INCY user-agent builtin",
			userAgent: "INCY/2.0.0",
			expected:  true,
		},
		{
			name:      "FlClash X builtin",
			userAgent: "FlClash X/0.14.0",
			expected:  true,
		},
		{
			name:      "FlClashX without space builtin",
			userAgent: "FlClashX/0.14.0",
			expected:  true,
		},
		{
			name:      "Flowvy builtin",
			userAgent: "Flowvy/1.0",
			expected:  true,
		},
		{
			name:      "prizrak-box builtin",
			userAgent: "prizrak-box/1.0",
			expected:  true,
		},
		{
			name:      "koala-clash builtin",
			userAgent: "koala-clash/2.1",
			expected:  true,
		},
		{
			name:      "standard v2rayNG is not extended",
			userAgent: "v2rayNG/1.8.5",
			expected:  false,
		},
		{
			name:      "empty user agent",
			userAgent: "",
			expected:  false,
		},
		{
			name:       "custom regex match",
			userAgent:  "CustomVPNClient/3.0",
			additional: []string{"^CustomVPNClient/"},
			expected:   true,
		},
		{
			name:       "custom regex no match",
			userAgent:  "OtherClient/1.0",
			additional: []string{"^CustomVPNClient/"},
			expected:   false,
		},
		{
			name:       "invalid regex pattern does not panic",
			userAgent:  "AnyClient/1.0",
			additional: []string{"[invalid("},
			expected:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isExtendedClient(tt.userAgent, tt.additional)
			if got != tt.expected {
				t.Errorf("isExtendedClient(%q, %v) = %v; want %v", tt.userAgent, tt.additional, got, tt.expected)
			}
		})
	}
}

func TestIsJSONSubscriptionFallbackSupported(t *testing.T) {
	tests := []struct {
		userAgent string
		expected  bool
	}{
		{"Streisand/1.7.8", true},
		{"streisand/2.0", true},
		{"Happ/1.0", true},
		{"INCY/1.0", true},
		{"ktor-client/2.3.0", true},
		{"V2Box/1.5.5", true},
		{"io.github.saeeddev94.xray/1.0", true},
		{"v2rayNG/1.8.12", true},
		{"v2rayN/6.23.0", true},
		{"v2plus/1.2.3", true},
		{"curl/8.0.1", false},
		{"ClashMeta/1.16.0", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.userAgent, func(t *testing.T) {
			got := isJSONSubscriptionFallbackSupported(tt.userAgent)
			if got != tt.expected {
				t.Errorf("isJSONSubscriptionFallbackSupported(%q) = %v; want %v", tt.userAgent, got, tt.expected)
			}
		})
	}
}
