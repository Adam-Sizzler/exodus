package geo

import (
	"context"
	"net/netip"
	"strings"
	"sync"

	"exodus-node/geocheck/internal/asn"
	"exodus-node/geocheck/internal/jsonx"
	"exodus-node/geocheck/internal/netx"
)

// identityServices are echo endpoints that return the caller's address. They
// are queried in parallel and the answer is taken by majority, so a single
// misbehaving or hijacked endpoint cannot skew the result.
var identityServices = []string{
	"https://www.cloudflare.com/cdn-cgi/trace",
	"https://api64.ipify.org",
	"https://ident.me",
	"https://ifconfig.me/ip",
	"https://icanhazip.com",
}

// PublicIP discovers the external address for one family.
func PublicIP(ctx context.Context, stack *netx.Stack, f netx.Family) (netip.Addr, bool) {
	type vote struct {
		addr netip.Addr
		ok   bool
	}

	votes := make([]vote, len(identityServices))
	var wg sync.WaitGroup
	for i, url := range identityServices {
		wg.Add(1)
		go func() {
			defer wg.Done()
			resp, err := stack.Do(ctx, f, netx.Request{URL: url})
			if err != nil {
				return
			}
			raw := strings.TrimSpace(resp.Text())
			if strings.Contains(raw, "ip=") {
				for _, line := range strings.Split(raw, "\n") {
					if k, v, ok := strings.Cut(line, "="); ok && k == "ip" {
						raw = v
						break
					}
				}
			}
			addr, err := netip.ParseAddr(strings.TrimSpace(raw))
			if err != nil {
				return
			}
			addr = addr.Unmap()
			if (f == netx.V4) != addr.Is4() {
				return
			}
			votes[i] = vote{addr, true}
		}()
	}
	wg.Wait()

	counts := map[netip.Addr]int{}
	var firstSeen netip.Addr
	for _, v := range votes {
		if !v.ok {
			continue
		}
		if !firstSeen.IsValid() {
			firstSeen = v.addr
		}
		counts[v.addr]++
		if counts[v.addr] >= 2 {
			return v.addr, true
		}
	}
	return firstSeen, firstSeen.IsValid()
}

// Identity is the summary shown above the tables.
type Identity struct {
	IPv4 netip.Addr
	IPv6 netip.Addr
	ASN  asn.Info
	// Org is a human readable network operator name, richer than the AS name
	// when a lookup service provides one.
	Org string
}

// Describe resolves the AS and operator behind the detected addresses.
func Describe(ctx context.Context, stack *netx.Stack, ar *asn.Resolver, v4, v6 netip.Addr) Identity {
	id := Identity{IPv4: v4, IPv6: v6}

	ip, family := v4, netx.V4
	if !ip.IsValid() {
		ip, family = v6, netx.V6
	}
	if !ip.IsValid() {
		return id
	}

	id.ASN = ar.Lookup(ctx, ip)
	id.Org = id.ASN.Name

	// Cymru names are terse registry strings; prefer a descriptive org name
	// when one of the enrichment endpoints answers.
	for _, src := range []struct{ url, numPath, namePath string }{
		{"https://ipinfo.check.place/" + ip.String(), "ASN.AutonomousSystemNumber", "ASN.AutonomousSystemOrganization"},
		{"https://geoip.oxl.app/api/ip/" + ip.String(), "asn", "organization.name"},
	} {
		resp, err := stack.Do(ctx, family, netx.Request{URL: src.url})
		if err != nil || !resp.OK() {
			continue
		}
		org := strings.TrimSpace(jsonx.String(resp.Body, src.namePath))
		if org == "" {
			continue
		}
		id.Org = org
		if id.ASN.Number == 0 {
			if n := jsonx.String(resp.Body, src.numPath); n != "" {
				id.ASN.Number = atoiSafe(n)
			}
		}
		break
	}
	return id
}

func atoiSafe(s string) int {
	n := 0
	for _, c := range s {
		if c < '0' || c > '9' {
			return n
		}
		n = n*10 + int(c-'0')
	}
	return n
}
