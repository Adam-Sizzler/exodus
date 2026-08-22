package netx

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/netip"
	"strings"
	"sync"
	"time"

	"github.com/miekg/dns"
)

// DefaultDoHResolvers are probed in order when Options.DoH is "auto".
// All of them are IP literals so resolving the resolver is never a problem.
var DefaultDoHResolvers = []struct {
	Name string
	URL  string
}{
	{"Cloudflare", "https://1.1.1.1/dns-query"},
	{"Quad9", "https://9.9.9.9/dns-query"},
	{"Google", "https://8.8.8.8/dns-query"},
	{"AdGuard", "https://94.140.14.140/dns-query"},
}

// Resolver answers A/AAAA queries over DNS-over-HTTPS when one is configured
// and falls back to the system resolver otherwise. Results are memoised for
// the lifetime of the run so a burst of checks against the same host costs one
// lookup.
type Resolver struct {
	stack *Stack

	mu       sync.RWMutex
	dohURL   string
	dohName  string
	dohHTTP  *http.Client
	requestC string // configured value: "", "auto" or an explicit URL

	cache sync.Map // cacheKey -> *cacheEntry
}

type cacheKey struct {
	host   string
	family Family
}

type cacheEntry struct {
	once  sync.Once
	addrs []netip.Addr
	err   error
}

func newResolver(s *Stack, doh string) *Resolver {
	r := &Resolver{stack: s, requestC: strings.TrimSpace(doh)}
	// A dedicated client: it must never route through Stack.dialContext or a
	// DoH lookup would need a DoH lookup.
	r.dohHTTP = &http.Client{
		Timeout: s.opts.Timeout,
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout: s.opts.Timeout,
				Control: bindToDevice(s.device),
			}).DialContext,
			TLSHandshakeTimeout: s.opts.Timeout,
			ForceAttemptHTTP2:   true,
		},
	}
	return r
}

// warmup selects a working DoH endpoint. It is best effort: on failure the
// resolver silently keeps using the system resolver.
func (r *Resolver) warmup(ctx context.Context) {
	switch r.requestC {
	case "", "off", "none", "system":
		return
	case "auto":
		for _, c := range DefaultDoHResolvers {
			probe, cancel := context.WithTimeout(ctx, 2500*time.Millisecond)
			addrs, err := r.queryDoH(probe, c.URL, "www.google.com", V4)
			cancel()
			if err == nil && len(addrs) > 0 {
				r.mu.Lock()
				r.dohURL, r.dohName = c.URL, c.Name
				r.mu.Unlock()
				return
			}
		}
	default:
		r.mu.Lock()
		r.dohURL, r.dohName = r.requestC, r.requestC
		r.mu.Unlock()
	}
}

// Active reports the DoH endpoint in use, or "" when resolving via the system.
func (r *Resolver) Active() string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.dohName
}

// Lookup returns the addresses of host for the given family.
func (r *Resolver) Lookup(ctx context.Context, host string, f Family) ([]netip.Addr, error) {
	if ip, err := netip.ParseAddr(host); err == nil {
		return []netip.Addr{ip}, nil
	}
	host = strings.TrimSuffix(strings.ToLower(host), ".")

	v, _ := r.cache.LoadOrStore(cacheKey{host, f}, &cacheEntry{})
	e := v.(*cacheEntry)
	e.once.Do(func() {
		e.addrs, e.err = r.lookupUncached(ctx, host, f)
	})
	return e.addrs, e.err
}

func (r *Resolver) lookupUncached(ctx context.Context, host string, f Family) ([]netip.Addr, error) {
	r.mu.RLock()
	url := r.dohURL
	r.mu.RUnlock()

	if url != "" {
		if addrs, err := r.queryDoH(ctx, url, host, f); err == nil && len(addrs) > 0 {
			return addrs, nil
		}
		// Fall through to the system resolver rather than failing the check.
	}
	return r.querySystem(ctx, host, f)
}

// LookupTXT returns the TXT records of host. It is used for the Team Cymru
// IP-to-ASN service, which is published purely over DNS.
func (r *Resolver) LookupTXT(ctx context.Context, host string) ([]string, error) {
	r.mu.RLock()
	url := r.dohURL
	r.mu.RUnlock()

	if url != "" {
		if txt, err := r.txtDoH(ctx, url, host); err == nil && len(txt) > 0 {
			return txt, nil
		}
	}
	return net.DefaultResolver.LookupTXT(ctx, host)
}

// LookupPTR returns the reverse DNS name of an address, without the trailing
// dot, or "" when there is none.
func (r *Resolver) LookupPTR(ctx context.Context, ip netip.Addr) string {
	arpa, err := dns.ReverseAddr(ip.String())
	if err != nil {
		return ""
	}

	r.mu.RLock()
	url := r.dohURL
	r.mu.RUnlock()

	if url != "" {
		if reply, err := r.exchange(ctx, url, arpa, dns.TypePTR); err == nil {
			for _, rr := range reply.Answer {
				if p, ok := rr.(*dns.PTR); ok {
					return strings.TrimSuffix(p.Ptr, ".")
				}
			}
		}
	}
	names, err := net.DefaultResolver.LookupAddr(ctx, ip.String())
	if err != nil || len(names) == 0 {
		return ""
	}
	return strings.TrimSuffix(names[0], ".")
}

func (r *Resolver) txtDoH(ctx context.Context, url, host string) ([]string, error) {
	reply, err := r.exchange(ctx, url, host, dns.TypeTXT)
	if err != nil {
		return nil, err
	}
	var out []string
	for _, rr := range reply.Answer {
		if t, ok := rr.(*dns.TXT); ok {
			out = append(out, strings.Join(t.Txt, ""))
		}
	}
	return out, nil
}

func (r *Resolver) querySystem(ctx context.Context, host string, f Family) ([]netip.Addr, error) {
	network := "ip4"
	if f == V6 {
		network = "ip6"
	}
	addrs, err := net.DefaultResolver.LookupNetIP(ctx, network, host)
	if err != nil {
		return nil, err
	}
	out := make([]netip.Addr, 0, len(addrs))
	for _, a := range addrs {
		out = append(out, a.Unmap())
	}
	if len(out) == 0 {
		return nil, fmt.Errorf("no %s address for %s", f, host)
	}
	return out, nil
}

// exchange performs one RFC 8484 wire-format query over HTTPS POST.
func (r *Resolver) exchange(ctx context.Context, url, host string, qtype uint16) (*dns.Msg, error) {
	m := new(dns.Msg)
	m.SetQuestion(dns.Fqdn(host), qtype)
	m.RecursionDesired = true
	wire, err := m.Pack()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(wire))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/dns-message")
	req.Header.Set("Accept", "application/dns-message")

	resp, err := r.dohHTTP.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("doh %s: http %d", url, resp.StatusCode)
	}
	body, err := io.ReadAll(io.LimitReader(resp.Body, 64<<10))
	if err != nil {
		return nil, err
	}

	reply := new(dns.Msg)
	if err := reply.Unpack(body); err != nil {
		return nil, fmt.Errorf("doh %s: %w", url, err)
	}
	return reply, nil
}

// queryDoH resolves A/AAAA records for host over DoH.
func (r *Resolver) queryDoH(ctx context.Context, url, host string, f Family) ([]netip.Addr, error) {
	qtype := dns.TypeA
	if f == V6 {
		qtype = dns.TypeAAAA
	}
	reply, err := r.exchange(ctx, url, host, qtype)
	if err != nil {
		return nil, err
	}

	var out []netip.Addr
	for _, rr := range reply.Answer {
		switch v := rr.(type) {
		case *dns.A:
			if ip, ok := netip.AddrFromSlice(v.A.To4()); ok {
				out = append(out, ip)
			}
		case *dns.AAAA:
			if ip, ok := netip.AddrFromSlice(v.AAAA.To16()); ok {
				out = append(out, ip)
			}
		}
	}
	if len(out) == 0 {
		return nil, fmt.Errorf("doh %s: no %s records for %s", url, f, host)
	}
	return out, nil
}
