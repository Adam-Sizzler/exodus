// Package mtr measures the network path to a set of well-known destinations
// and judges how directly your traffic reaches them.
package mtr

import (
	"math"
	"net/netip"
	"sort"
	"time"

	"exodus-node/geocheck/internal/asn"
)

// Hop is one TTL position on the path, aggregated over all probe rounds.
type Hop struct {
	TTL  int
	Addr netip.Addr // zero when no reply was ever received
	Host string     // reverse DNS name, when resolved
	ASN  asn.Info

	// Addrs lists every address that answered at this TTL. More than one
	// means the path is load balanced across parallel routers.
	Addrs []netip.Addr

	Sent int
	Recv int
	RTTs []time.Duration
}

// Loss is the fraction of probes that went unanswered, 0..1.
func (h Hop) Loss() float64 {
	if h.Sent == 0 {
		return 0
	}
	return float64(h.Sent-h.Recv) / float64(h.Sent)
}

// Responded reports whether anything ever answered at this TTL.
func (h Hop) Responded() bool { return h.Recv > 0 && h.Addr.IsValid() }

// Best returns the minimum RTT, which is the only figure that reflects the
// path itself; higher samples include router queueing and ICMP deprioritisation.
func (h Hop) Best() time.Duration { return pick(h.RTTs, 0) }

// Worst returns the maximum RTT.
func (h Hop) Worst() time.Duration { return pick(h.RTTs, len(h.RTTs)-1) }

// Avg returns the mean RTT.
func (h Hop) Avg() time.Duration {
	if len(h.RTTs) == 0 {
		return 0
	}
	var sum time.Duration
	for _, d := range h.RTTs {
		sum += d
	}
	return sum / time.Duration(len(h.RTTs))
}

// StdDev returns the jitter across samples.
func (h Hop) StdDev() time.Duration {
	if len(h.RTTs) < 2 {
		return 0
	}
	mean := float64(h.Avg())
	var acc float64
	for _, d := range h.RTTs {
		diff := float64(d) - mean
		acc += diff * diff
	}
	return time.Duration(math.Sqrt(acc / float64(len(h.RTTs)-1)))
}

func pick(rtts []time.Duration, i int) time.Duration {
	if len(rtts) == 0 {
		return 0
	}
	s := append([]time.Duration(nil), rtts...)
	sort.Slice(s, func(a, b int) bool { return s[a] < s[b] })
	return s[i]
}

// Report is the outcome of tracing one target.
type Report struct {
	Target   Target
	Resolved netip.Addr
	DestASN  asn.Info
	Method   Method
	Hops     []Hop
	Verdict  Verdict
	Err      error
}

// FinalHop returns the hop that answered as the destination, or the last hop
// that answered at all.
func (r *Report) FinalHop() (Hop, bool) {
	for i := len(r.Hops) - 1; i >= 0; i-- {
		if r.Hops[i].Responded() {
			return r.Hops[i], true
		}
	}
	return Hop{}, false
}

// Method is how a target was probed.
type Method string

const (
	// MethodICMP is a classic TTL-limited ICMP echo traceroute.
	MethodICMP Method = "icmp"
	// MethodTCP measures the TCP handshake to the service port. It carries no
	// per-hop detail but works without privileges and against ICMP-silent hosts.
	MethodTCP Method = "tcp"
)
