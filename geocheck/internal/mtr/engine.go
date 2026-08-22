package mtr

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"net/netip"
	"os"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"exodus-node/geocheck/internal/netx"
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
	"golang.org/x/net/ipv6"
)

// magic tags our probes so stray ICMP traffic is never mistaken for a reply.
var magic = [4]byte{'g', 'c', 'h', 'k'}

// probeReply is what the read loop hands back to a waiting probe.
type probeReply struct {
	from  netip.Addr
	at    time.Time
	final bool // an echo reply: the destination itself answered

	// quotedDst is the destination a router echoed back inside its error.
	// It is checked against the probe's target to reject stray ICMP traffic.
	quotedDst netip.Addr
}

type pendingProbe struct {
	sentAt time.Time
	dst    netip.Addr
	ch     chan probeReply
}

// engine owns one ICMP socket and multiplexes every in-flight probe over it,
// so tracing twenty targets costs one file descriptor rather than twenty.
type engine struct {
	family netx.Family
	conn   *icmp.PacketConn
	v4     *ipv4.PacketConn
	v6     *ipv6.PacketConn
	raw    bool
	id     uint16

	sendMu sync.Mutex // guards the socket-wide TTL and the write
	curTTL int

	mu      sync.Mutex
	pending map[uint16]*pendingProbe
	seq     atomic.Uint32

	closeOnce sync.Once
	done      chan struct{}
}

// newEngine opens the best ICMP socket available: a raw socket when we hold
// the privilege, otherwise an unprivileged datagram ("ping") socket.
func newEngine(family netx.Family, src netip.Addr, iface string) (*engine, error) {
	rawNet, dgramNet, listen := "ip4:icmp", "udp4", "0.0.0.0"
	if family == netx.V6 {
		rawNet, dgramNet, listen = "ip6:ipv6-icmp", "udp6", "::"
	}
	if src.IsValid() {
		listen = src.String()
	}

	e := &engine{
		family:  family,
		pending: make(map[uint16]*pendingProbe),
		done:    make(chan struct{}),
		id:      uint16(os.Getpid() & 0xffff),
		curTTL:  -1,
	}

	conn, raw, err := openICMPSocket(rawNet, dgramNet, listen)
	if err != nil {
		return nil, err
	}
	e.conn, e.raw = conn, raw

	if family == netx.V6 {
		e.v6 = conn.IPv6PacketConn()
		_ = e.v6.SetControlMessage(ipv6.FlagHopLimit|ipv6.FlagSrc, true)
		// Ask the kernel to drop everything a traceroute cannot use, which
		// also blunts the BSDs' habit of copying all host ICMP to every socket.
		var f ipv6.ICMPFilter
		f.SetAll(true)
		f.Accept(ipv6.ICMPTypeEchoReply)
		f.Accept(ipv6.ICMPTypeTimeExceeded)
		f.Accept(ipv6.ICMPTypeDestinationUnreachable)
		_ = e.v6.SetICMPFilter(&f)
	} else {
		e.v4 = conn.IPv4PacketConn()
		// Only FlagTTL is accepted on a Darwin datagram ICMP socket, and the
		// call fails as a unit, so ask for nothing else. It is advisory here
		// in any case: replies are matched on the payload, not the header.
		_ = e.v4.SetControlMessage(ipv4.FlagTTL, true)
	}

	if iface != "" {
		e.bindInterface(iface)
	}

	go e.readLoop()
	return e, nil
}

// openICMPSocket opens the best ICMP socket available: a raw one when the
// privilege can be obtained, otherwise an unprivileged datagram socket.
//
// The goroutine is pinned to its OS thread for the whole operation because
// Linux capabilities are per-thread, and raisePrivilege() raises the effective
// bit on the calling thread only. Without pinning, the runtime is free to move
// this goroutine to a different thread between raising the capability and
// calling socket(), and the raw socket then fails with EPERM — intermittently,
// which looks like a flaky network rather than a scheduling race.
func openICMPSocket(rawNet, dgramNet, listen string) (*icmp.PacketConn, bool, error) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	raisePrivilege()

	if conn, err := icmp.ListenPacket(rawNet, listen); err == nil {
		return conn, true, nil
	}
	conn, err := icmp.ListenPacket(dgramNet, listen)
	if err != nil {
		return nil, false, fmt.Errorf("open icmp socket: %w", err)
	}
	return conn, false, nil
}

// Raw reports whether the engine holds a raw socket, which is what makes
// TTL manipulation reliable across platforms.
func (e *engine) Raw() bool { return e.raw }

// idPreserved reports whether the ICMP identifier we set survives to the wire.
// Raw sockets always keep it; Linux overwrites it on datagram sockets with the
// socket's own port so it can demultiplex replies itself.
func (e *engine) idPreserved() bool { return e.raw || datagramKeepsICMPID }

func (e *engine) bindInterface(name string) {
	ifi, err := net.InterfaceByName(name)
	if err != nil {
		return
	}
	if e.v6 != nil {
		_ = e.v6.SetMulticastInterface(ifi)
	} else if e.v4 != nil {
		_ = e.v4.SetMulticastInterface(ifi)
	}
}

func (e *engine) Close() error {
	var err error
	e.closeOnce.Do(func() {
		close(e.done)
		err = e.conn.Close()
	})
	return err
}

// send transmits one echo request towards dst with the given TTL and returns
// the sequence number identifying it.
func (e *engine) send(dst netip.Addr, ttl int) (uint16, *pendingProbe, error) {
	seq := uint16(e.seq.Add(1))

	payload := make([]byte, 16)
	copy(payload, magic[:])
	binary.BigEndian.PutUint16(payload[4:], seq)
	binary.BigEndian.PutUint16(payload[6:], uint16(ttl))

	var msgType icmp.Type = ipv4.ICMPTypeEcho
	if e.family == netx.V6 {
		msgType = ipv6.ICMPTypeEchoRequest
	}
	wb, err := (&icmp.Message{
		Type: msgType,
		Code: 0,
		Body: &icmp.Echo{ID: int(e.id), Seq: int(seq), Data: payload},
	}).Marshal(nil)
	if err != nil {
		return 0, nil, err
	}

	p := &pendingProbe{dst: dst, ch: make(chan probeReply, 1)}

	var target net.Addr
	if e.raw {
		target = &net.IPAddr{IP: dst.AsSlice()}
	} else {
		target = &net.UDPAddr{IP: dst.AsSlice()}
	}

	e.sendMu.Lock()
	if e.curTTL != ttl {
		var terr error
		if e.family == netx.V6 {
			terr = e.v6.SetHopLimit(ttl)
		} else {
			terr = e.v4.SetTTL(ttl)
		}
		if terr != nil {
			e.sendMu.Unlock()
			return 0, nil, fmt.Errorf("set ttl %d: %w", ttl, terr)
		}
		e.curTTL = ttl
	}
	// Register before writing so a fast reply cannot beat the bookkeeping.
	e.mu.Lock()
	p.sentAt = time.Now()
	e.pending[seq] = p
	e.mu.Unlock()

	_, err = e.conn.WriteTo(wb, target)
	e.sendMu.Unlock()

	if err != nil {
		e.mu.Lock()
		delete(e.pending, seq)
		e.mu.Unlock()
		return 0, nil, err
	}
	return seq, p, nil
}

func (e *engine) forget(seq uint16) {
	e.mu.Lock()
	delete(e.pending, seq)
	e.mu.Unlock()
}

// probe sends one probe, waits for its reply and returns the round-trip time.
func (e *engine) probe(ctx context.Context, dst netip.Addr, ttl int, timeout time.Duration) (probeReply, time.Duration, error) {
	seq, p, err := e.send(dst, ttl)
	if err != nil {
		return probeReply{}, 0, err
	}
	defer e.forget(seq)

	timer := time.NewTimer(timeout)
	defer timer.Stop()

	select {
	case r := <-p.ch:
		return r, r.at.Sub(p.sentAt), nil
	case <-timer.C:
		return probeReply{}, 0, errTimeout
	case <-ctx.Done():
		return probeReply{}, 0, ctx.Err()
	case <-e.done:
		return probeReply{}, 0, net.ErrClosed
	}
}

var errTimeout = errors.New("probe timed out")

// readLoop demultiplexes every inbound ICMP packet to the probe that is
// waiting for it.
func (e *engine) readLoop() {
	buf := make([]byte, 1500)
	proto := 1 // ICMPv4
	if e.family == netx.V6 {
		proto = 58
	}

	for {
		select {
		case <-e.done:
			return
		default:
		}

		n, peer, err := e.conn.ReadFrom(buf)
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				return
			}
			continue
		}
		at := time.Now()

		msg, err := icmp.ParseMessage(proto, buf[:n])
		if err != nil {
			continue
		}
		seq, quotedDst, final, ok := e.match(msg)
		if !ok {
			continue
		}

		from, _ := netip.AddrFromSlice(addrBytes(peer))
		e.deliver(seq, probeReply{from: from.Unmap(), at: at, final: final, quotedDst: quotedDst})
	}
}

// match extracts the sequence number of the probe an ICMP message refers to,
// plus the destination quoted back by a router when there is one.
func (e *engine) match(msg *icmp.Message) (seq uint16, quotedDst netip.Addr, final bool, ok bool) {
	switch body := msg.Body.(type) {
	case *icmp.Echo:
		// An echo reply: the destination answered, quoting our payload back
		// verbatim. Requiring our magic is what keeps this correct on the BSDs,
		// where every ICMP datagram socket receives a copy of all inbound ICMP
		// on the host, including other processes' replies.
		if s, ok := seqFromPayload(body.Data); ok {
			return s, netip.Addr{}, true, true
		}
		// Without the payload the only remaining evidence is the identifier,
		// which Linux rewrites on datagram sockets; trust it only when we know
		// it survived.
		if e.idPreserved() && uint16(body.ID) == e.id {
			return uint16(body.Seq), netip.Addr{}, true, true
		}
		return 0, netip.Addr{}, false, false

	case *icmp.TimeExceeded:
		return e.matchQuoted(body.Data)
	case *icmp.DstUnreach:
		return e.matchQuoted(body.Data)
	default:
		return 0, netip.Addr{}, false, false
	}
}

// matchQuoted digs the original echo request out of the packet an intermediate
// router quoted back to us.
func (e *engine) matchQuoted(quoted []byte) (seq uint16, quotedDst netip.Addr, final bool, ok bool) {
	quotedDst = quotedDestination(quoted, e.family)

	inner, err := stripIPHeader(quoted, e.family)
	if err != nil {
		return 0, quotedDst, false, false
	}
	proto := 1
	if e.family == netx.V6 {
		proto = 58
	}
	msg, err := icmp.ParseMessage(proto, inner)
	if err != nil {
		return 0, quotedDst, false, false
	}
	echo, isEcho := msg.Body.(*icmp.Echo)
	if !isEcho {
		return 0, quotedDst, false, false
	}
	if s, ok := seqFromPayload(echo.Data); ok {
		return s, quotedDst, false, true
	}
	// Routers must quote at least 8 bytes of the original datagram, which
	// covers the ICMP header (id and seq) but usually stops before our payload.
	return uint16(echo.Seq), quotedDst, false, true
}

// quotedDestination reads the destination address out of the IP header a
// router echoed back. Matching it against the probe's target is what stops us
// adopting another process's ICMP errors as our own.
func quotedDestination(b []byte, family netx.Family) netip.Addr {
	if family == netx.V6 {
		if len(b) < ipv6.HeaderLen {
			return netip.Addr{}
		}
		addr, _ := netip.AddrFromSlice(b[24:40])
		return addr
	}
	if len(b) < ipv4.HeaderLen {
		return netip.Addr{}
	}
	addr, _ := netip.AddrFromSlice(b[16:20])
	return addr
}

// seqFromPayload recovers the sequence number we stamped into the probe body.
func seqFromPayload(data []byte) (uint16, bool) {
	if len(data) < 6 || [4]byte(data[:4]) != magic {
		return 0, false
	}
	return binary.BigEndian.Uint16(data[4:6]), true
}

// stripIPHeader removes the quoted IP header, honouring IPv4's variable length.
func stripIPHeader(b []byte, family netx.Family) ([]byte, error) {
	if family == netx.V6 {
		if len(b) < ipv6.HeaderLen {
			return nil, errors.New("truncated quoted ipv6 header")
		}
		return b[ipv6.HeaderLen:], nil
	}
	if len(b) < ipv4.HeaderLen {
		return nil, errors.New("truncated quoted ipv4 header")
	}
	ihl := int(b[0]&0x0f) * 4
	if ihl < ipv4.HeaderLen || ihl > len(b) {
		ihl = ipv4.HeaderLen
	}
	return b[ihl:], nil
}

func (e *engine) deliver(seq uint16, r probeReply) {
	e.mu.Lock()
	p, ok := e.pending[seq]
	if ok && r.quotedDst.IsValid() && r.quotedDst != p.dst {
		// Another process's probe happened to reuse our sequence number.
		e.mu.Unlock()
		return
	}
	if ok {
		delete(e.pending, seq)
	}
	e.mu.Unlock()
	if !ok {
		return
	}
	select {
	case p.ch <- r:
	default:
	}
}

func addrBytes(a net.Addr) []byte {
	switch v := a.(type) {
	case *net.IPAddr:
		if ip := v.IP.To4(); ip != nil {
			return ip
		}
		return v.IP
	case *net.UDPAddr:
		if ip := v.IP.To4(); ip != nil {
			return ip
		}
		return v.IP
	default:
		return nil
	}
}
