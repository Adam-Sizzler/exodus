package subscription

import (
	"net"
	"strings"
)

// isDomainAddress reports whether addr looks like a domain name rather than an IP address.
func isDomainAddress(addr string) bool {
	addr = strings.TrimSpace(addr)
	if addr == "" {
		return false
	}
	// Check IPv4
	if net.ParseIP(addr) != nil {
		return false
	}
	// Check bracketed IPv6 [::1]
	if strings.HasPrefix(addr, "[") && strings.HasSuffix(addr, "]") {
		if net.ParseIP(addr[1:len(addr)-1]) != nil {
			return false
		}
	}
	return strings.Contains(addr, ".")
}

// resolveFinalServerName resolves the effective SNI / ServerName for a host,
// matching the exact precedence rules of upstream:
// 1. If host.KeepSNIBlank -> empty string ""
// 2. If host.OverrideSNIFromAddress -> host.Address (highest priority)
// 3. If host.SNI is set and non-empty -> host.SNI
// 4. If fallbackSNI (from inbound / defaults) is non-empty -> fallbackSNI
// 5. If host.Address is a domain -> host.Address
// 6. Otherwise -> empty string ""
func resolveFinalServerName(host SubscriptionHost, fallbackSNI string) string {
	if host.KeepSNIBlank {
		return ""
	}
	if host.OverrideSNIFromAddress {
		return host.Address
	}
	if host.SNI != nil && strings.TrimSpace(*host.SNI) != "" {
		return strings.TrimSpace(*host.SNI)
	}
	fallbackSNI = strings.TrimSpace(fallbackSNI)
	if fallbackSNI != "" {
		return fallbackSNI
	}
	if isDomainAddress(host.Address) {
		return host.Address
	}
	return ""
}
