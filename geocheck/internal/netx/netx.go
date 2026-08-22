// Package netx builds the network stack shared by every check: HTTP clients
// pinned to an address family, optional source-interface binding, SOCKS5 and
// DNS-over-HTTPS resolution.
package netx

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/netip"
	"strings"
	"time"

	"golang.org/x/net/proxy"
)

// Family selects the IP version a request must travel over.
type Family int

const (
	// V4 selects IPv4. The numeric values match the IP version so they read
	// naturally in flags and logs.
	V4 Family = 4
	// V6 selects IPv6.
	V6 Family = 6
)

func (f Family) String() string {
	if f == V6 {
		return "IPv6"
	}
	return "IPv4"
}

// Network appends the family suffix Go's dialers expect ("tcp" -> "tcp4").
func (f Family) Network(base string) string {
	if f == V6 {
		return base + "6"
	}
	return base + "4"
}

// Options configures the stack. The zero value is usable.
type Options struct {
	// Interface pins every outgoing socket to a source. It accepts either an
	// interface name ("eth0") or a literal address already assigned to this
	// host ("203.0.113.10"), which is the only way to choose between many
	// aliases on one NIC.
	Interface string
	// Proxy is a SOCKS5 endpoint as host:port. Names are still resolved
	// locally so that the -4/-6 selection keeps its meaning.
	Proxy string
	// Timeout bounds a single request end to end.
	Timeout time.Duration
	// DoH selects the DNS-over-HTTPS resolver: "" for the system resolver,
	// "auto" to probe the built-in list, or an explicit https:// URL.
	DoH string
	// UserAgent is sent with every HTTP request that does not override it.
	UserAgent string
}

const defaultUserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 " +
	"(KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36"

// Stack owns the dialers, resolver and HTTP clients for one run.
type Stack struct {
	opts     Options
	resolver *Resolver

	localV4 netip.Addr
	localV6 netip.Addr

	// device is the interface name to pin sockets to. It is derived, not taken
	// verbatim from Options.Interface, because that may hold an address.
	device string
	// pinned is set when a literal source address was given, which fixes the
	// address family too.
	pinned  Family
	hasPin  bool
	clients map[Family]*http.Client
}

// New builds a stack. It fails only on unusable configuration (missing
// interface, malformed proxy); resolver probing is best effort.
func New(ctx context.Context, opts Options) (*Stack, error) {
	if opts.Timeout <= 0 {
		opts.Timeout = 6 * time.Second
	}
	if opts.UserAgent == "" {
		opts.UserAgent = defaultUserAgent
	}

	s := &Stack{opts: opts, clients: make(map[Family]*http.Client, 2)}

	if opts.Interface != "" {
		if err := s.resolveBinding(opts.Interface); err != nil {
			return nil, err
		}
	}

	if opts.Proxy != "" {
		if _, _, err := net.SplitHostPort(opts.Proxy); err != nil {
			return nil, fmt.Errorf("invalid proxy address %q: expected host:port", opts.Proxy)
		}
	}

	s.resolver = newResolver(s, opts.DoH)
	s.resolver.warmup(ctx)

	for _, f := range []Family{V4, V6} {
		s.clients[f] = s.newHTTPClient(f)
	}
	return s, nil
}

// UserAgent returns the default UA string for this run.
func (s *Stack) UserAgent() string { return s.opts.UserAgent }

// Timeout returns the per-request budget.
func (s *Stack) Timeout() time.Duration { return s.opts.Timeout }

// Resolver exposes the DoH-or-system resolver for non-HTTP consumers (MTR).
func (s *Stack) Resolver() *Resolver { return s.resolver }

// LocalAddr returns the pinned source address for a family, if any.
func (s *Stack) LocalAddr(f Family) netip.Addr {
	if f == V6 {
		return s.localV6
	}
	return s.localV4
}

// HTTP returns the client pinned to the given address family.
func (s *Stack) HTTP(f Family) *http.Client { return s.clients[f] }

// dialer returns a net.Dialer bound to the configured interface and family.
func (s *Stack) dialer(f Family) *net.Dialer {
	d := &net.Dialer{
		Timeout:   s.opts.Timeout,
		KeepAlive: 30 * time.Second,
		Control:   bindToDevice(s.device),
	}
	if local := s.LocalAddr(f); local.IsValid() {
		d.LocalAddr = &net.TCPAddr{IP: local.AsSlice(), Zone: local.Zone()}
	}
	return d
}

// dialContext resolves the host itself (so DoH and the family pin are honoured)
// and then connects, optionally through SOCKS5.
func (s *Stack) dialContext(f Family) func(context.Context, string, string) (net.Conn, error) {
	// The requested network is ignored: the family pin decides it, so that
	// Happy Eyeballs cannot quietly fall back to the other one.
	return func(ctx context.Context, _, addr string) (net.Conn, error) {
		host, port, err := net.SplitHostPort(addr)
		if err != nil {
			return nil, err
		}

		var ips []netip.Addr
		if ip, err := netip.ParseAddr(host); err == nil {
			if (f == V4) != ip.Is4() && !ip.Is4In6() {
				return nil, fmt.Errorf("literal %s is not %s", ip, f)
			}
			ips = []netip.Addr{ip}
		} else {
			ips, err = s.resolver.Lookup(ctx, host, f)
			if err != nil {
				return nil, err
			}
		}

		var errs []error
		for _, ip := range ips {
			target := net.JoinHostPort(ip.String(), port)
			conn, err := s.rawDial(ctx, f, target)
			if err == nil {
				return conn, nil
			}
			errs = append(errs, err)
			if ctx.Err() != nil {
				break
			}
		}
		if len(errs) == 0 {
			return nil, fmt.Errorf("no %s address for %s", f, host)
		}
		return nil, errors.Join(errs...)
	}
}

// rawDial performs the actual connect, through SOCKS5 when configured.
func (s *Stack) rawDial(ctx context.Context, f Family, target string) (net.Conn, error) {
	d := s.dialer(f)
	if s.opts.Proxy == "" {
		return d.DialContext(ctx, f.Network("tcp"), target)
	}
	// The proxy itself is reached over whichever family works; only the
	// tunnelled destination is pinned.
	pd, err := proxy.SOCKS5("tcp", s.opts.Proxy, nil, &net.Dialer{
		Timeout: s.opts.Timeout,
		Control: bindToDevice(s.device),
	})
	if err != nil {
		return nil, err
	}
	if cd, ok := pd.(proxy.ContextDialer); ok {
		return cd.DialContext(ctx, "tcp", target)
	}
	return pd.Dial("tcp", target)
}

func (s *Stack) newHTTPClient(f Family) *http.Client {
	tr := &http.Transport{
		DialContext:           s.dialContext(f),
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          64,
		MaxIdleConnsPerHost:   4,
		IdleConnTimeout:       30 * time.Second,
		TLSHandshakeTimeout:   s.opts.Timeout,
		ExpectContinueTimeout: time.Second,
	}
	return &http.Client{
		Transport: tr,
		Timeout:   s.opts.Timeout,
		CheckRedirect: func(_ *http.Request, via []*http.Request) error {
			if len(via) >= 10 {
				return errors.New("stopped after 10 redirects")
			}
			return nil
		},
	}
}

// Close releases idle connections held by the HTTP clients.
func (s *Stack) Close() {
	for _, c := range s.clients {
		if tr, ok := c.Transport.(*http.Transport); ok {
			tr.CloseIdleConnections()
		}
	}
}

// resolveBinding interprets the --interface value, which may be either an
// interface name or a literal source address.
//
// Accepting an address matters on hosts that carry many of them: a server with
// a dozen aliases on one NIC cannot be steered by interface name alone, and the
// address the kernel picks by default is rarely the one you wanted to test.
func (s *Stack) resolveBinding(spec string) error {
	addr, err := netip.ParseAddr(spec)
	if err != nil {
		// Not an address, so treat it as an interface name.
		v4, v6, err := interfaceAddrs(spec)
		if err != nil {
			return err
		}
		if !v4.IsValid() && !v6.IsValid() {
			return fmt.Errorf("interface %q has no usable global address", spec)
		}
		s.device, s.localV4, s.localV6 = spec, v4, v6
		return nil
	}

	addr = addr.Unmap()
	owner, ok := interfaceOwning(addr)
	if !ok {
		return fmt.Errorf("address %s is not assigned to any interface on this host%s",
			addr, localAddressHint())
	}

	// Pinning a specific address also pins the family: asking for a v4 source
	// and then running v6 checks over some other address would be a lie.
	s.device = owner
	if addr.Is4() {
		s.localV4, s.pinned = addr, V4
	} else {
		s.localV6, s.pinned = addr, V6
	}
	s.hasPin = true
	return nil
}

// PinnedFamily reports the address family fixed by a literal --interface
// address, if one was given.
func (s *Stack) PinnedFamily() (Family, bool) { return s.pinned, s.hasPin }

// BindDevice returns the interface name sockets should be pinned to, which is
// resolved from an address when one was supplied.
func (s *Stack) BindDevice() string { return s.device }

// interfaceOwning finds the interface that holds an address.
func interfaceOwning(addr netip.Addr) (string, bool) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", false
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
				return iface.Name, true
			}
		}
	}
	return "", false
}

// localAddressHint lists the addresses that could have been meant, so the error
// tells the user what to type instead of only what was wrong.
func localAddressHint() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	var usable []string
	for _, a := range addrs {
		pfx, ok := a.(*net.IPNet)
		if !ok {
			continue
		}
		ip, ok := netip.AddrFromSlice(pfx.IP)
		if !ok {
			continue
		}
		ip = ip.Unmap()
		if !ip.IsGlobalUnicast() || ip.IsLinkLocalUnicast() {
			continue
		}
		usable = append(usable, ip.String())
		if len(usable) == 12 {
			usable = append(usable, "...")
			break
		}
	}
	if len(usable) == 0 {
		return ""
	}
	return "; available: " + strings.Join(usable, ", ")
}

// interfaceAddrs picks the best global source address of each family on iface.
func interfaceAddrs(name string) (v4, v6 netip.Addr, err error) {
	iface, err := net.InterfaceByName(name)
	if err != nil {
		return netip.Addr{}, netip.Addr{}, fmt.Errorf("interface %q: %w", name, err)
	}
	addrs, err := iface.Addrs()
	if err != nil {
		return netip.Addr{}, netip.Addr{}, fmt.Errorf("interface %q: %w", name, err)
	}
	for _, a := range addrs {
		pfx, ok := a.(*net.IPNet)
		if !ok {
			continue
		}
		ip, ok := netip.AddrFromSlice(pfx.IP)
		if !ok {
			continue
		}
		ip = ip.Unmap()
		if ip.IsLoopback() || ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() {
			continue
		}
		if ip.Is4() {
			if !v4.IsValid() {
				v4 = ip
			}
			continue
		}
		// Prefer a global unicast v6 over a ULA.
		if !v6.IsValid() || (isULA(v6) && !isULA(ip)) {
			v6 = ip
		}
	}
	return v4, v6, nil
}

func isULA(ip netip.Addr) bool {
	return ip.Is6() && ip.As16()[0]&0xfe == 0xfc
}

// HasFamily reports whether the host holds an address that could plausibly
// reach the internet over f. IPv4 accepts RFC1918 (NAT is the norm); IPv6
// requires a real global address because ULA never leaves the site.
func HasFamily(f Family) bool {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return false
	}
	for _, a := range addrs {
		pfx, ok := a.(*net.IPNet)
		if !ok {
			continue
		}
		ip, ok := netip.AddrFromSlice(pfx.IP)
		if !ok {
			continue
		}
		ip = ip.Unmap()
		if !ip.IsGlobalUnicast() {
			continue
		}
		if f == V4 && ip.Is4() {
			return true
		}
		if f == V6 && ip.Is6() && !isULA(ip) {
			return true
		}
	}
	return false
}
