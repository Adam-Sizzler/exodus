package asn

import (
	"net/netip"
	"strings"
	"testing"
)

func TestCymruOriginName(t *testing.T) {
	cases := []struct{ ip, want string }{
		{"8.8.8.8", "8.8.8.8.origin.asn.cymru.com"},
		{"1.0.0.1", "1.0.0.1.origin.asn.cymru.com"},
		{"203.0.113.42", "42.113.0.203.origin.asn.cymru.com"},
	}
	for _, c := range cases {
		if got := cymruOriginName(netip.MustParseAddr(c.ip)); got != c.want {
			t.Errorf("cymruOriginName(%s) = %q, want %q", c.ip, got, c.want)
		}
	}
}

func TestCymruOriginNameIPv6(t *testing.T) {
	// 2001:4860:4860::8888 expands to 32 nibbles, reversed, one label each.
	got := cymruOriginName(netip.MustParseAddr("2001:4860:4860::8888"))
	want := "8.8.8.8." + strings.Repeat("0.", 16) +
		"0.6.8.4.0.6.8.4.1.0.0.2.origin6.asn.cymru.com"
	if got != want {
		t.Errorf("cymruOriginName(v6) = %q, want %q", got, want)
	}
	if n := strings.Count(strings.TrimSuffix(got, ".origin6.asn.cymru.com"), "."); n != 31 {
		t.Errorf("got %d separators, want 31 for 32 nibbles", n)
	}
}

func TestRoutable(t *testing.T) {
	cases := []struct {
		ip   string
		want bool
	}{
		{"8.8.8.8", true},
		{"203.0.113.5", true},
		{"2001:4860:4860::8888", true},
		{"10.0.0.1", false},        // RFC1918
		{"192.168.1.1", false},     // RFC1918
		{"172.16.5.4", false},      // RFC1918
		{"127.0.0.1", false},       // loopback
		{"169.254.1.1", false},     // link-local
		{"100.64.0.1", false},      // CGNAT, 100.64.0.0/10
		{"100.127.255.255", false}, // last CGNAT address
		{"100.128.0.1", true},      // just above CGNAT, publicly routable
		{"100.63.255.255", true},   // just below CGNAT
		{"fd00::1", false},         // ULA
		{"fe80::1", false},         // link-local v6
	}
	for _, c := range cases {
		if got := routable(netip.MustParseAddr(c.ip)); got != c.want {
			t.Errorf("routable(%s) = %v, want %v", c.ip, got, c.want)
		}
	}
}

func TestSplitAndTrim(t *testing.T) {
	fields := splitCymru("15169 | 8.8.8.0/24 | US | arin | 1992-12-01")
	if len(fields) != 5 {
		t.Fatalf("got %d fields, want 5", len(fields))
	}
	if fields[0] != "15169" || fields[1] != "8.8.8.0/24" || fields[2] != "US" {
		t.Errorf("unexpected fields: %q", fields)
	}

	if got := trimASSuffix("GOOGLE - Google LLC, US"); got != "GOOGLE - Google LLC" {
		t.Errorf("trimASSuffix = %q", got)
	}
	if got := trimASSuffix("CLOUDFLARENET"); got != "CLOUDFLARENET" {
		t.Errorf("trimASSuffix = %q", got)
	}
}

func TestInfoString(t *testing.T) {
	cases := []struct {
		info Info
		want string
	}{
		{Info{Number: 13335, Name: "CLOUDFLARENET"}, "AS13335 CLOUDFLARENET"},
		{Info{Number: 13335}, "AS13335"},
		{Info{Name: "CLOUDFLARENET"}, "CLOUDFLARENET"},
		{Info{}, ""},
	}
	for _, c := range cases {
		if got := c.info.String(); got != c.want {
			t.Errorf("Info%+v.String() = %q, want %q", c.info, got, c.want)
		}
	}
	if !(Info{}).Empty() {
		t.Error("zero Info should be Empty")
	}
}
