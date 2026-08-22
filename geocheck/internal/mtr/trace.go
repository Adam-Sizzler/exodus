package mtr

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/netip"
	"sync"
	"time"

	"exodus-node/geocheck/internal/asn"
	"exodus-node/geocheck/internal/netx"
)

// Config tunes the tracer.
type Config struct {
	Family    netx.Family
	SourceIP  netip.Addr
	Interface string

	MaxTTL  int           // highest TTL to probe
	Rounds  int           // probes per hop
	Timeout time.Duration // per-probe reply budget
	Pace    time.Duration // delay between consecutive sends

	// Targets is how many destinations are traced at once.
	Targets int

	Resolver   *netx.Resolver
	ASN        *asn.Resolver
	ReverseDNS bool
}

func (c *Config) applyDefaults() {
	if c.MaxTTL <= 0 {
		c.MaxTTL = 30
	}
	if c.Rounds <= 0 {
		c.Rounds = 5
	}
	if c.Timeout <= 0 {
		c.Timeout = 1500 * time.Millisecond
	}
	if c.Pace <= 0 {
		c.Pace = 3 * time.Millisecond
	}
	if c.Targets <= 0 {
		c.Targets = 6
	}
}

// Tracer traces many targets over one shared ICMP socket.
type Tracer struct {
	cfg    Config
	engine *engine

	// gate paces every send across all targets so routers do not rate-limit
	// our probes into looking like packet loss.
	gate chan struct{}
}

// NewTracer opens the ICMP socket. The returned tracer must be closed.
// A nil error with ICMPAvailable() false means only TCP probing is possible.
func NewTracer(cfg Config) (*Tracer, error) {
	cfg.applyDefaults()
	t := &Tracer{cfg: cfg, gate: make(chan struct{}, 1)}

	eng, err := newEngine(cfg.Family, cfg.SourceIP, cfg.Interface)
	if err != nil {
		return t, nil // TCP-only mode; not an error the caller must abort on
	}
	t.engine = eng
	return t, nil
}

// ICMPAvailable reports whether hop-by-hop tracing is possible.
func (t *Tracer) ICMPAvailable() bool { return t.engine != nil }

// Privileged reports whether the ICMP socket is a raw one.
func (t *Tracer) Privileged() bool { return t.engine != nil && t.engine.Raw() }

// Capability describes what this process can actually measure, so the CLI can
// explain a thin result instead of silently printing an empty path.
type Capability struct {
	// ICMP is true when an ICMP socket was opened at all.
	ICMP bool
	// Raw is true when it is a raw socket, which sees every ICMP error.
	Raw bool
	// PathVisible is true when intermediate hops can be observed. When false,
	// only destination latency is measurable and targets fall back to TCP.
	PathVisible bool
	// Hint tells the user how to unlock the full path, when something is missing.
	Hint string
}

// Capability reports the tracer's measurement ability on this host.
func (t *Tracer) Capability() Capability {
	switch {
	case t.engine == nil:
		return Capability{Hint: "could not open an ICMP socket; " + privilegeHint}
	case t.engine.Raw():
		return Capability{ICMP: true, Raw: true, PathVisible: true}
	case datagramSeesTimeExceeded:
		return Capability{ICMP: true, PathVisible: true}
	default:
		return Capability{
			ICMP: true,
			Hint: "unprivileged ICMP sockets on this platform never receive the " +
				"TimeExceeded replies traceroute needs, so per-hop detail is unavailable; " +
				privilegeHint,
		}
	}
}

// BinaryPathPlaceholder is substituted by the CLI with the running executable.
const BinaryPathPlaceholder = "{binary}"

// Close releases the socket.
func (t *Tracer) Close() error {
	if t.engine == nil {
		return nil
	}
	return t.engine.Close()
}

// Run traces every target, honouring the configured parallelism, and reports
// each result through onDone as soon as it is ready.
func (t *Tracer) Run(ctx context.Context, targets []Target, onDone func(*Report)) []*Report {
	reports := make([]*Report, len(targets))
	sem := make(chan struct{}, t.cfg.Targets)

	var wg sync.WaitGroup
	for i, tgt := range targets {
		wg.Add(1)
		go func() {
			defer wg.Done()
			select {
			case sem <- struct{}{}:
			case <-ctx.Done():
				return
			}
			defer func() { <-sem }()

			r := t.Trace(ctx, tgt)
			reports[i] = r
			if onDone != nil {
				onDone(r)
			}
		}()
	}
	wg.Wait()

	out := reports[:0]
	for _, r := range reports {
		if r != nil {
			out = append(out, r)
		}
	}
	// Verdicts are only meaningful relative to the whole run: each target is
	// judged against the lowest latency this connection achieved anywhere.
	ClassifyAll(out)
	return out
}

// Trace measures the path to one target.
func (t *Tracer) Trace(ctx context.Context, tgt Target) *Report {
	rep := &Report{Target: tgt}

	dst, err := t.resolve(ctx, tgt)
	if err != nil {
		rep.Err = err
		rep.Verdict = Verdict{Class: ClassUnreachable, Notes: []string{err.Error()}}
		return rep
	}
	rep.Resolved = dst
	if t.cfg.ASN != nil {
		rep.DestASN = t.cfg.ASN.Lookup(ctx, dst)
	}

	if t.engine == nil || tgt.ICMPSilent || !t.Capability().PathVisible {
		t.traceTCP(ctx, rep)
	} else {
		t.traceICMP(ctx, rep)
		// A path where nothing at all answered tells us nothing; fall back to
		// a TCP handshake so the target still yields a latency figure.
		if !anyResponse(rep.Hops) {
			t.traceTCP(ctx, rep)
		}
	}

	t.annotate(ctx, rep)
	// A provisional verdict, so a single Trace call is usable on its own.
	// Run re-classifies the whole set once the latency floor is known.
	rep.Verdict = Classify(rep, 0)
	return rep
}

func (t *Tracer) resolve(ctx context.Context, tgt Target) (netip.Addr, error) {
	if ip, err := netip.ParseAddr(tgt.Host); err == nil {
		ip = ip.Unmap()
		if (t.cfg.Family == netx.V4) != ip.Is4() {
			return netip.Addr{}, fmt.Errorf("%s is not reachable over %s", tgt.Host, t.cfg.Family)
		}
		return ip, nil
	}
	if t.cfg.Resolver == nil {
		return netip.Addr{}, errors.New("no resolver configured")
	}
	addrs, err := t.cfg.Resolver.Lookup(ctx, tgt.Host, t.cfg.Family)
	if err != nil {
		return netip.Addr{}, fmt.Errorf("resolve %s: %w", tgt.Host, err)
	}
	if len(addrs) == 0 {
		return netip.Addr{}, fmt.Errorf("resolve %s: no %s address", tgt.Host, t.cfg.Family)
	}
	return addrs[0].Unmap(), nil
}

// traceICMP walks the path with TTL-limited echo requests. Every TTL of a
// round is in flight at once, so a round costs one timeout rather than thirty.
func (t *Tracer) traceICMP(ctx context.Context, rep *Report) {
	rep.Method = MethodICMP

	type slot struct {
		mu    sync.Mutex
		sent  int
		recv  int
		rtts  []time.Duration
		addrs []netip.Addr
	}
	slots := make([]*slot, t.cfg.MaxTTL+1)
	for i := range slots {
		slots[i] = &slot{}
	}

	// limit shrinks as soon as we learn which TTL reaches the destination.
	var limitMu sync.Mutex
	limit := t.cfg.MaxTTL

	for round := 0; round < t.cfg.Rounds; round++ {
		if ctx.Err() != nil {
			break
		}
		limitMu.Lock()
		hi := limit
		limitMu.Unlock()

		var wg sync.WaitGroup
		for ttl := 1; ttl <= hi; ttl++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				t.throttle(ctx)

				s := slots[ttl]
				s.mu.Lock()
				s.sent++
				s.mu.Unlock()

				reply, rtt, err := t.engine.probe(ctx, rep.Resolved, ttl, t.cfg.Timeout)
				if err != nil {
					return
				}

				s.mu.Lock()
				s.recv++
				s.rtts = append(s.rtts, rtt)
				if !containsAddr(s.addrs, reply.from) && reply.from.IsValid() {
					s.addrs = append(s.addrs, reply.from)
				}
				s.mu.Unlock()

				if reply.final || reply.from == rep.Resolved {
					limitMu.Lock()
					if ttl < limit {
						limit = ttl
					}
					limitMu.Unlock()
				}
			}()
		}
		wg.Wait()
	}

	limitMu.Lock()
	hi := limit
	limitMu.Unlock()

	hops := make([]Hop, 0, hi)
	for ttl := 1; ttl <= hi; ttl++ {
		s := slots[ttl]
		h := Hop{TTL: ttl, Sent: s.sent, Recv: s.recv, RTTs: s.rtts, Addrs: s.addrs}
		if len(s.addrs) > 0 {
			h.Addr = s.addrs[0]
		}
		hops = append(hops, h)
	}
	rep.Hops = trimTrailingSilence(hops)
}

// throttle spaces sends globally.
func (t *Tracer) throttle(ctx context.Context) {
	select {
	case t.gate <- struct{}{}:
	case <-ctx.Done():
		return
	}
	timer := time.NewTimer(t.cfg.Pace)
	go func() {
		<-timer.C
		<-t.gate
	}()
}

// traceTCP measures the TCP handshake to the target's service port. There is
// no per-hop view, but the destination RTT is measured honestly and, unlike
// ICMP, it is never deprioritised or filtered.
func (t *Tracer) traceTCP(ctx context.Context, rep *Report) {
	rep.Method = MethodTCP

	port := rep.Target.Port
	if port == 0 {
		port = 443
	}
	addr := net.JoinHostPort(rep.Resolved.String(), fmt.Sprint(port))

	d := &net.Dialer{Timeout: t.cfg.Timeout}
	if t.cfg.SourceIP.IsValid() {
		d.LocalAddr = &net.TCPAddr{IP: t.cfg.SourceIP.AsSlice()}
	}

	hop := Hop{TTL: 0, Addr: rep.Resolved}
	for i := 0; i < t.cfg.Rounds; i++ {
		if ctx.Err() != nil {
			break
		}
		hop.Sent++
		start := time.Now()
		conn, err := d.DialContext(ctx, t.cfg.Family.Network("tcp"), addr)
		if err != nil {
			continue
		}
		rtt := time.Since(start)
		_ = conn.Close()
		hop.Recv++
		hop.RTTs = append(hop.RTTs, rtt)
		// Handshakes are heavier than ICMP; do not hammer the endpoint.
		select {
		case <-time.After(50 * time.Millisecond):
		case <-ctx.Done():
		}
	}
	hop.Addrs = []netip.Addr{rep.Resolved}
	rep.Hops = []Hop{hop}
}

// annotate resolves the AS and reverse name of every responding hop.
func (t *Tracer) annotate(ctx context.Context, rep *Report) {
	var wg sync.WaitGroup
	for i := range rep.Hops {
		h := &rep.Hops[i]
		if !h.Addr.IsValid() {
			continue
		}
		wg.Add(1)
		go func() {
			defer wg.Done()
			if t.cfg.ASN != nil {
				h.ASN = t.cfg.ASN.Lookup(ctx, h.Addr)
			}
			if t.cfg.ReverseDNS && t.cfg.Resolver != nil {
				rctx, cancel := context.WithTimeout(ctx, 1200*time.Millisecond)
				h.Host = t.cfg.Resolver.LookupPTR(rctx, h.Addr)
				cancel()
			}
		}()
	}
	wg.Wait()

	if rep.DestASN.Empty() && t.cfg.ASN != nil && rep.Resolved.IsValid() {
		rep.DestASN = t.cfg.ASN.Lookup(ctx, rep.Resolved)
	}
}

func containsAddr(list []netip.Addr, a netip.Addr) bool {
	for _, x := range list {
		if x == a {
			return true
		}
	}
	return false
}

func anyResponse(hops []Hop) bool {
	for _, h := range hops {
		if h.Responded() {
			return true
		}
	}
	return false
}

// trimTrailingSilence drops unanswered hops past the last responding one, so a
// filtered destination does not print twenty rows of stars.
func trimTrailingSilence(hops []Hop) []Hop {
	last := -1
	for i, h := range hops {
		if h.Responded() {
			last = i
		}
	}
	if last < 0 {
		return hops[:0]
	}
	return hops[:last+1]
}
