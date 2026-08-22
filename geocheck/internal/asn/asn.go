// Package asn resolves IP addresses to their originating autonomous system
// using the Team Cymru IP-to-ASN service, which is published over plain DNS
// and therefore works through the same DoH transport as everything else.
package asn

import (
	"context"
	"fmt"
	"net/netip"
	"strconv"
	"strings"
	"sync"

	"exodus-node/geocheck/internal/netx"
)

// Info describes the autonomous system an address belongs to.
type Info struct {
	Number  int    // 0 when unknown
	Name    string // e.g. "CLOUDFLARENET"
	Country string // registry country of the AS
	Prefix  string // announced BGP prefix covering the address
}

// String renders the AS compactly, e.g. "AS13335 CLOUDFLARENET".
func (i Info) String() string {
	switch {
	case i.Number == 0 && i.Name == "":
		return ""
	case i.Name == "":
		return "AS" + strconv.Itoa(i.Number)
	case i.Number == 0:
		return i.Name
	default:
		return fmt.Sprintf("AS%d %s", i.Number, i.Name)
	}
}

// Empty reports whether nothing was resolved.
func (i Info) Empty() bool { return i.Number == 0 && i.Name == "" }

// Resolver looks up and caches AS information.
type Resolver struct {
	dns *netx.Resolver

	ipCache   sync.Map // netip.Addr -> *ipEntry
	nameCache sync.Map // int -> *nameEntry
}

type ipEntry struct {
	once sync.Once
	info Info
}

type nameEntry struct {
	once    sync.Once
	name    string
	country string
}

// New builds a resolver on top of the shared DNS resolver.
func New(dns *netx.Resolver) *Resolver { return &Resolver{dns: dns} }

// Lookup returns the AS that originates the prefix covering ip. Addresses that
// cannot have a public origin (private, loopback, link-local) return an empty
// Info without touching the network.
func (r *Resolver) Lookup(ctx context.Context, ip netip.Addr) Info {
	ip = ip.Unmap()
	if !routable(ip) {
		return Info{}
	}
	v, _ := r.ipCache.LoadOrStore(ip, &ipEntry{})
	e := v.(*ipEntry)
	e.once.Do(func() { e.info = r.lookupUncached(ctx, ip) })
	return e.info
}

func (r *Resolver) lookupUncached(ctx context.Context, ip netip.Addr) Info {
	name := cymruOriginName(ip)
	if name == "" {
		return Info{}
	}
	txt, err := r.dns.LookupTXT(ctx, name)
	if err != nil || len(txt) == 0 {
		return Info{}
	}

	// "15169 | 8.8.8.0/24 | US | arin | 1992-12-01"
	// The first field can carry several ASNs for a multi-origin prefix.
	fields := splitCymru(txt[0])
	if len(fields) < 1 {
		return Info{}
	}
	info := Info{}
	if first, _, _ := strings.Cut(fields[0], " "); first != "" {
		if n, err := strconv.Atoi(first); err == nil {
			info.Number = n
		}
	}
	if len(fields) > 1 {
		info.Prefix = fields[1]
	}
	if len(fields) > 2 {
		info.Country = strings.ToUpper(fields[2])
	}
	if info.Number != 0 {
		info.Name, _ = r.describe(ctx, info.Number)
	}
	return info
}

// describe resolves an AS number to its registered name.
func (r *Resolver) describe(ctx context.Context, num int) (name, country string) {
	v, _ := r.nameCache.LoadOrStore(num, &nameEntry{})
	e := v.(*nameEntry)
	e.once.Do(func() {
		txt, err := r.dns.LookupTXT(ctx, fmt.Sprintf("AS%d.asn.cymru.com", num))
		if err != nil || len(txt) == 0 {
			return
		}
		// "15169 | US | arin | 2000-03-30 | GOOGLE, US"
		fields := splitCymru(txt[0])
		if len(fields) >= 2 {
			e.country = strings.ToUpper(fields[1])
		}
		if len(fields) >= 5 {
			e.name = trimASSuffix(fields[4])
		}
	})
	return e.name, e.country
}

// Describe returns the registered name of an AS number.
func (r *Resolver) Describe(ctx context.Context, num int) Info {
	if num == 0 {
		return Info{}
	}
	name, country := r.describe(ctx, num)
	return Info{Number: num, Name: name, Country: country}
}

func splitCymru(s string) []string {
	parts := strings.Split(s, "|")
	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}
	return parts
}

// trimASSuffix drops the trailing country Cymru appends to AS names,
// turning "GOOGLE, US" into "GOOGLE".
func trimASSuffix(s string) string {
	if i := strings.LastIndex(s, ","); i > 0 && len(s)-i <= 4 {
		return strings.TrimSpace(s[:i])
	}
	return strings.TrimSpace(s)
}

// cymruOriginName builds the reverse-nibble query name for an address.
func cymruOriginName(ip netip.Addr) string {
	if ip.Is4() {
		b := ip.As4()
		return fmt.Sprintf("%d.%d.%d.%d.origin.asn.cymru.com", b[3], b[2], b[1], b[0])
	}
	if !ip.Is6() {
		return ""
	}
	b := ip.As16()
	var sb strings.Builder
	sb.Grow(80)
	const hex = "0123456789abcdef"
	for i := 15; i >= 0; i-- {
		sb.WriteByte(hex[b[i]&0x0f])
		sb.WriteByte('.')
		sb.WriteByte(hex[b[i]>>4])
		sb.WriteByte('.')
	}
	sb.WriteString("origin6.asn.cymru.com")
	return sb.String()
}

// routable reports whether an address can plausibly have a public BGP origin.
func routable(ip netip.Addr) bool {
	return ip.IsValid() &&
		ip.IsGlobalUnicast() &&
		!ip.IsPrivate() &&
		!ip.IsLoopback() &&
		!ip.IsLinkLocalUnicast() &&
		!isCGNAT(ip) &&
		!isULA(ip)
}

func isCGNAT(ip netip.Addr) bool {
	if !ip.Is4() {
		return false
	}
	b := ip.As4()
	return b[0] == 100 && b[1] >= 64 && b[1] <= 127
}

func isULA(ip netip.Addr) bool {
	return ip.Is6() && !ip.Is4In6() && ip.As16()[0]&0xfe == 0xfc
}
