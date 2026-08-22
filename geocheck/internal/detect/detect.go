// Package detect looks for the things that quietly invalidate a network
// measurement: a local tunnel carrying the default route, or a resolver that
// is answering on someone else's behalf. These are structural checks — they
// hold regardless of how plausible the latency numbers look.
package detect

import (
	"context"
	"net"
	"net/netip"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/miekg/dns"
)

// Severity ranks how much a finding should change your reading of the report.
type Severity int

const (
	// Info is context worth knowing about.
	Info Severity = iota
	// Warn means some results may be misleading.
	Warn
	// Alert means the measurement is being intercepted and the routing
	// conclusions cannot be trusted.
	Alert
)

func (s Severity) String() string {
	switch s {
	case Alert:
		return "alert"
	case Warn:
		return "warn"
	default:
		return "info"
	}
}

// Finding is one detected condition.
type Finding struct {
	ID       string
	Title    string
	Severity Severity
	Detail   string
}

// Options configures the checks.
type Options struct {
	Timeout time.Duration
	// Skip disables the DNS probes, which send plain port-53 queries and so
	// bypass any configured proxy.
	SkipDNS bool
}

// Run executes every check concurrently and returns what it found.
func Run(ctx context.Context, opts Options) []Finding {
	if opts.Timeout <= 0 {
		opts.Timeout = 3 * time.Second
	}

	checks := []func(context.Context, Options) *Finding{
		checkTunnelRoute,
		checkResolverAddress,
	}
	if !opts.SkipDNS {
		checks = append(checks, checkOpenDNSIdentity, checkGoogleIDServer)
	}

	out := make([]*Finding, len(checks))
	var wg sync.WaitGroup
	for i, fn := range checks {
		wg.Add(1)
		go func() {
			defer wg.Done()
			out[i] = fn(ctx, opts)
		}()
	}
	wg.Wait()

	findings := make([]Finding, 0, len(out))
	for _, f := range out {
		if f != nil {
			findings = append(findings, *f)
		}
	}
	return findings
}

// checkTunnelRoute reports whether the address the kernel would use to reach
// the internet belongs to a point-to-point tunnel device. A userspace VPN or
// tun2socks client carrying the default route rewrites TTLs and answers ICMP
// itself, which makes every traceroute result fiction.
func checkTunnelRoute(_ context.Context, _ Options) *Finding {
	src, ok := defaultSourceAddr()
	if !ok {
		return nil
	}
	iface, ok := interfaceFor(src)
	if !ok {
		return nil
	}

	name := strings.ToLower(iface.Name)
	tunnelish := iface.Flags&net.FlagPointToPoint != 0
	for _, prefix := range []string{"tun", "utun", "tap", "ppp", "wg", "ipsec", "gpd", "nordlynx"} {
		if strings.HasPrefix(name, prefix) {
			tunnelish = true
			break
		}
	}
	if !tunnelish {
		return nil
	}

	return &Finding{
		ID:       "tunnel_default_route",
		Title:    "Default route runs through a tunnel",
		Severity: Alert,
		Detail: "traffic leaves via " + iface.Name + " (" + src.String() + "), a point-to-point " +
			"tunnel device. A userspace VPN or proxy client on this path answers ICMP itself " +
			"and rewrites TTLs, so hop-by-hop results describe the tunnel, not the internet.",
	}
}

// checkResolverAddress flags a system resolver inside 198.18.0.0/15. That range
// is reserved for benchmarking and never appears on a real network; it is the
// fake-IP range used by tun2socks-style clients (Clash, sing-box), so finding
// it there means DNS is being answered locally.
func checkResolverAddress(_ context.Context, _ Options) *Finding {
	benchmark := netip.MustParsePrefix("198.18.0.0/15")

	for _, addr := range systemResolvers() {
		if benchmark.Contains(addr) {
			return &Finding{
				ID:       "resolver_fake_ip",
				Title:    "System resolver is a local proxy",
				Severity: Alert,
				Detail: "the configured nameserver " + addr.String() + " is inside 198.18.0.0/15, " +
					"a benchmarking range that never appears on a real network. A local proxy is " +
					"synthesising DNS answers.",
			}
		}
	}
	return nil
}

// notAnOpenDNSResolver is the authoritative TXT record for which.opendns.com,
// and the entire signal this check rests on.
//
// The name is answered two different ways. A genuine OpenDNS resolver
// intercepts it internally and replies with its own identifier — "r3001.lon",
// "r3007.ams". Every other resolver has to fetch the record from the
// authoritative servers, where it reads this sentence. Getting it back from an
// address that is supposed to *be* an OpenDNS resolver therefore proves the
// query was answered by something else.
const notAnOpenDNSResolver = "not an opendns resolver"

// checkOpenDNSIdentity asks OpenDNS which of its servers answered.
func checkOpenDNSIdentity(ctx context.Context, opts Options) *Finding {
	answers, err := queryTXT(ctx, "208.67.222.222:53", "which.opendns.com.", dns.ClassINET, opts.Timeout)
	if err != nil || len(answers) == 0 {
		return nil // unreachable is not the same as intercepted
	}
	return openDNSFinding(answers)
}

// openDNSFinding is the decision, split from the query so it can be tested
// against captured answers instead of the network.
func openDNSFinding(answers []string) *Finding {
	for _, a := range answers {
		if !strings.Contains(strings.ToLower(a), notAnOpenDNSResolver) {
			continue
		}
		return &Finding{
			ID:       "dns_hijack_opendns",
			Title:    "Port 53 DNS is being intercepted",
			Severity: Alert,
			Detail: "208.67.222.222 was asked which OpenDNS server answered and replied " +
				quote(a) + ". That is the record the public authoritative servers hold, " +
				"not what an OpenDNS resolver returns, so the query was answered by " +
				"something else on the path.",
		}
	}
	return nil
}

// checkGoogleIDServer uses a query Google Public DNS deliberately does not
// implement. A real answer can only have come from an impostor.
func checkGoogleIDServer(ctx context.Context, opts Options) *Finding {
	answers, err := queryTXT(ctx, "8.8.8.8:53", "id.server.", dns.ClassCHAOS, opts.Timeout)
	if err != nil || len(answers) == 0 {
		return nil
	}
	return &Finding{
		ID:       "dns_hijack_google",
		Title:    "Port 53 DNS is being intercepted",
		Severity: Alert,
		Detail: "8.8.8.8 answered a CHAOS id.server query with " + quote(answers[0]) + ". " +
			"Google Public DNS does not implement that query, so the reply came from " +
			"something else on the path.",
	}
}

func queryTXT(ctx context.Context, server, name string, class uint16, timeout time.Duration) ([]string, error) {
	m := new(dns.Msg)
	m.Question = []dns.Question{{Name: name, Qtype: dns.TypeTXT, Qclass: class}}
	m.RecursionDesired = true
	m.Id = dns.Id()

	c := &dns.Client{Timeout: timeout}
	reply, _, err := c.ExchangeContext(ctx, m, server)
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

// defaultSourceAddr returns the local address the kernel would use to reach the
// public internet. Dialling a UDP socket sends no packets but performs the
// route lookup, which works identically on every platform.
func defaultSourceAddr() (netip.Addr, bool) {
	conn, err := net.Dial("udp4", "8.8.8.8:53")
	if err != nil {
		return netip.Addr{}, false
	}
	// No packets are sent; the socket exists only to make the kernel perform
	// the route lookup, so closing it cannot fail in a way that matters.
	defer func() { _ = conn.Close() }()
	ap, ok := conn.LocalAddr().(*net.UDPAddr)
	if !ok {
		return netip.Addr{}, false
	}
	addr, ok := netip.AddrFromSlice(ap.IP)
	return addr.Unmap(), ok
}

func interfaceFor(addr netip.Addr) (net.Interface, bool) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return net.Interface{}, false
	}
	for _, iface := range ifaces {
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}
		for _, a := range addrs {
			pfx, ok := a.(*net.IPNet)
			if !ok {
				continue
			}
			if ip, ok := netip.AddrFromSlice(pfx.IP); ok && ip.Unmap() == addr {
				return iface, true
			}
		}
	}
	return net.Interface{}, false
}

// systemResolvers reads the configured nameservers. resolv.conf is absent on
// some platforms, in which case there is simply nothing to check.
func systemResolvers() []netip.Addr {
	data, err := os.ReadFile("/etc/resolv.conf")
	if err != nil {
		return nil
	}
	var out []netip.Addr
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, "nameserver") {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}
		if addr, err := netip.ParseAddr(fields[1]); err == nil {
			out = append(out, addr.Unmap())
		}
	}
	return out
}

func quote(s string) string {
	s = strings.TrimSpace(s)
	if len(s) > 60 {
		s = s[:57] + "..."
	}
	return "\"" + s + "\""
}
