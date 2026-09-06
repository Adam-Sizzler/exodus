package subscription

import (
	"testing"
)

func TestResolveFinalServerName(t *testing.T) {
	customSNI := "custom.domain.com"
	inboundSNI := "inbound.domain.com"

	tests := []struct {
		name        string
		host        SubscriptionHost
		fallbackSNI string
		expected    string
	}{
		{
			name: "KeepSNIBlank forces empty string",
			host: SubscriptionHost{
				Address:                "example.com",
				KeepSNIBlank:           true,
				OverrideSNIFromAddress: true,
				SNI:                    &customSNI,
			},
			fallbackSNI: inboundSNI,
			expected:    "",
		},
		{
			name: "OverrideSNIFromAddress takes priority over host.SNI and fallbackSNI",
			host: SubscriptionHost{
				Address:                "host-address.com",
				OverrideSNIFromAddress: true,
				SNI:                    &customSNI,
			},
			fallbackSNI: inboundSNI,
			expected:    "host-address.com",
		},
		{
			name: "host.SNI takes priority over fallbackSNI",
			host: SubscriptionHost{
				Address: "host-address.com",
				SNI:     &customSNI,
			},
			fallbackSNI: inboundSNI,
			expected:    customSNI,
		},
		{
			name: "fallbackSNI is used when host.SNI is empty",
			host: SubscriptionHost{
				Address: "host-address.com",
			},
			fallbackSNI: inboundSNI,
			expected:    inboundSNI,
		},
		{
			name: "host.Address is used when it is a domain and no SNI specified",
			host: SubscriptionHost{
				Address: "my-node.example.org",
			},
			fallbackSNI: "",
			expected:    "my-node.example.org",
		},
		{
			name: "IP address does not become SNI when no SNI specified",
			host: SubscriptionHost{
				Address: "192.168.1.1",
			},
			fallbackSNI: "",
			expected:    "",
		},
		{
			name: "IPv6 address does not become SNI",
			host: SubscriptionHost{
				Address: "2001:db8::1",
			},
			fallbackSNI: "",
			expected:    "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual := resolveFinalServerName(tc.host, tc.fallbackSNI)
			if actual != tc.expected {
				t.Fatalf("expected %q, got %q", tc.expected, actual)
			}
		})
	}
}
