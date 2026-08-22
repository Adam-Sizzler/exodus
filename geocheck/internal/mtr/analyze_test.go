package mtr

import (
	"net/netip"
	"strings"
	"testing"
	"time"

	"exodus-node/geocheck/internal/asn"
)

func ms(v float64) time.Duration {
	return time.Duration(v * float64(time.Millisecond))
}

// report builds a single-hop report answering at the given RTT.
func report(id string, rtt time.Duration, opts ...func(*Report)) *Report {
	dst := netip.MustParseAddr("203.0.113.10")
	r := &Report{
		Target:   Target{ID: id, Name: id, ASN: 64500},
		Resolved: dst,
		DestASN:  asn.Info{Number: 64500, Name: "EXAMPLE"},
		Method:   MethodICMP,
		Hops: []Hop{{
			TTL: 1, Addr: dst, Sent: 5, Recv: 5,
			RTTs: []time.Duration{rtt, rtt, rtt, rtt, rtt},
		}},
	}
	for _, o := range opts {
		o(r)
	}
	return r
}

func TestFloorIgnoresSingleOutlier(t *testing.T) {
	// One endpoint answers unusually fast. With four or more samples the floor
	// must come from the second fastest so that one outlier cannot define it.
	reports := []*Report{
		report("fast", ms(1.6)),
		report("a", ms(6.4)),
		report("b", ms(12)),
		report("c", ms(35)),
	}
	if got, want := Floor(reports), ms(6.4); got != want {
		t.Fatalf("Floor = %v, want %v", got, want)
	}
}

func TestFloorWithFewTargetsUsesMinimum(t *testing.T) {
	reports := []*Report{report("a", ms(9)), report("b", ms(30))}
	if got, want := Floor(reports), ms(9); got != want {
		t.Fatalf("Floor = %v, want %v", got, want)
	}
}

func TestFloorSkipsSubMillisecondReplies(t *testing.T) {
	// A sub-millisecond answer is local interception; it must not become the
	// baseline every other target is measured against.
	reports := []*Report{
		report("proxy", ms(0.2)),
		report("a", ms(8)),
		report("b", ms(9)),
	}
	if got, want := Floor(reports), ms(8); got != want {
		t.Fatalf("Floor = %v, want %v", got, want)
	}
}

// withPath prepends intermediate hops so the report looks like a real trace.
func withPath(rtts ...time.Duration) func(*Report) {
	return func(r *Report) {
		hops := make([]Hop, 0, len(rtts)+1)
		for i, rtt := range rtts {
			hops = append(hops, Hop{
				TTL:  i + 1,
				Addr: netip.MustParseAddr("198.51.100." + string(rune('1'+i))),
				Sent: 5, Recv: 5, RTTs: []time.Duration{rtt},
			})
		}
		final := r.Hops[len(r.Hops)-1]
		final.TTL = len(rtts) + 1
		r.Hops = append(hops, final)
	}
}

func TestClassifySubMillisecondWithNoPathIsIntercepted(t *testing.T) {
	// One hop, which is the destination itself: nothing was traversed, so the
	// reply cannot have come from far away.
	v := Classify(report("x", ms(0.3)), ms(8))
	if v.Class != ClassIntercepted {
		t.Fatalf("Class = %v, want intercepted", v.Class)
	}
	if v.Score != 0 {
		t.Fatalf("Score = %d, want 0 for an intercepted target", v.Score)
	}
}

// A datacenter host reaching an on-net cache is the case this tool is most
// often run on, and it must not be mistaken for interception: sub-millisecond
// latency behind a real multi-hop path is the best possible result.
func TestClassifySubMillisecondBehindRealPathIsDirect(t *testing.T) {
	r := report("x", ms(0.45), withPath(ms(0.1), ms(0.2), ms(0.3), ms(0.4), ms(0.42)))

	v := Classify(r, ms(0.4))
	if v.Class == ClassIntercepted {
		t.Fatalf("Class = intercepted; a six-hop path proves the packets travelled")
	}
	if v.Class != ClassDirect {
		t.Fatalf("Class = %v, want direct", v.Class)
	}
	if v.Score == 0 {
		t.Error("Score = 0; an on-net cache is the best outcome, not the worst")
	}
	if v.HopCount != 6 {
		t.Errorf("HopCount = %d, want 6", v.HopCount)
	}
}

// A TCP handshake never walks the path, so having no intermediate hops is the
// method rather than a finding. Datacenter hosts commonly lack a raw socket and
// reach every CDN in well under a millisecond over TCP; reading that as
// interception condemns almost every target at once.
func TestClassifySubMillisecondOverTCPIsNotIntercepted(t *testing.T) {
	r := report("x", ms(0.63))
	r.Method = MethodTCP

	v := Classify(r, ms(0.87))
	if v.Class == ClassIntercepted {
		t.Fatal("Class = intercepted; a TCP measurement cannot establish interception")
	}
	if v.Score == 0 {
		t.Error("Score = 0; nothing was actually wrong with this target")
	}

	// The ambiguity still has to be stated, just not as a verdict.
	found := false
	for _, n := range v.Notes {
		if strings.Contains(n, "TCP handshake") {
			found = true
		}
	}
	if !found {
		t.Error("the TCP caveat should still appear in the notes")
	}
}

func TestFloorUsesTCPReadingsToo(t *testing.T) {
	// With no raw socket every target is measured over TCP, so excluding those
	// readings would leave no baseline at all.
	mk := func(id string, rtt time.Duration) *Report {
		r := report(id, rtt)
		r.Method = MethodTCP
		return r
	}
	got := Floor([]*Report{mk("a", ms(0.63)), mk("b", ms(0.87)), mk("c", ms(0.9)), mk("d", ms(6.35))})
	if got != ms(0.87) {
		t.Fatalf("Floor = %v, want %v", got, ms(0.87))
	}
}

func TestFloorAcceptsGenuineSubMillisecondButRejectsLocalReplies(t *testing.T) {
	// A fast on-net cache with a real path defines the floor.
	onNet := report("cache", ms(0.45), withPath(ms(0.1), ms(0.2)))
	// A local responder with no path behind it must not.
	local := report("proxy", ms(0.2))

	if got, want := Floor([]*Report{onNet, local}), ms(0.45); got != want {
		t.Fatalf("Floor = %v, want %v (the corroborated reading)", got, want)
	}
	if got := Floor([]*Report{local}); got != 0 {
		t.Errorf("Floor = %v, want 0 when the only reading is locally answered", got)
	}
}

func TestClassifyRelativeToFloor(t *testing.T) {
	floor := ms(10)
	cases := []struct {
		name string
		rtt  time.Duration
		want Class
	}{
		{"at the floor", ms(11), ClassDirect},
		{"just inside the on-net band", ms(15), ClassDirect},
		{"same country", ms(25), ClassPeered},
		{"same continent", ms(60), ClassPeered},
		{"far beyond the continent", ms(200), ClassDetour},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := Classify(report("x", tc.rtt), floor).Class; got != tc.want {
				t.Fatalf("Class = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestClassifyTransitCarrierWins(t *testing.T) {
	// A low RTT does not make a path direct if a transit carrier is on it.
	r := report("x", ms(11), func(r *Report) {
		r.Hops = append([]Hop{{
			TTL: 1, Addr: netip.MustParseAddr("198.51.100.1"), Sent: 5, Recv: 5,
			RTTs: []time.Duration{ms(9)},
			ASN:  asn.Info{Number: 174, Name: "COGENT"},
		}}, r.Hops...)
		r.Hops[1].TTL = 2
	})

	v := Classify(r, ms(10))
	if v.Class != ClassTransit {
		t.Fatalf("Class = %v, want transit", v.Class)
	}
	if len(v.Transits) != 1 || v.Transits[0].ASN != 174 {
		t.Fatalf("Transits = %+v, want a single AS174", v.Transits)
	}
	if v.OnNet {
		t.Error("OnNet should be false when a transit carrier is on the path")
	}
}

func TestClassifyUnreachable(t *testing.T) {
	r := report("x", ms(10))
	r.Hops = []Hop{{TTL: 1, Sent: 5, Recv: 0}}

	v := Classify(r, ms(10))
	if v.Class != ClassUnreachable {
		t.Fatalf("Class = %v, want unreachable", v.Class)
	}
	if v.Loss != 1 {
		t.Fatalf("Loss = %v, want 1", v.Loss)
	}
}

func TestHopStatistics(t *testing.T) {
	h := Hop{
		Sent: 4, Recv: 3,
		RTTs: []time.Duration{ms(10), ms(30), ms(20)},
	}
	if got, want := h.Best(), ms(10); got != want {
		t.Errorf("Best = %v, want %v", got, want)
	}
	if got, want := h.Worst(), ms(30); got != want {
		t.Errorf("Worst = %v, want %v", got, want)
	}
	if got, want := h.Avg(), ms(20); got != want {
		t.Errorf("Avg = %v, want %v", got, want)
	}
	if got, want := h.Loss(), 0.25; got != want {
		t.Errorf("Loss = %v, want %v", got, want)
	}
}

func TestSelectTargets(t *testing.T) {
	if got := Select("all"); len(got) != len(Catalog) {
		t.Errorf("Select(all) returned %d targets, want %d", len(got), len(Catalog))
	}
	if got := Select("nonexistent-tag"); len(got) != 0 {
		t.Errorf("Select of an unknown tag returned %d targets, want 0", len(got))
	}
	got := Select("google")
	if len(got) == 0 {
		t.Fatal("Select(google) returned nothing")
	}
	for _, tg := range got {
		if !tg.HasTag("google") && tg.ID != "google" {
			t.Errorf("target %q does not carry the google tag", tg.ID)
		}
	}
	// An id selects exactly one target.
	if got := Select("cloudflare_dns"); len(got) != 1 {
		t.Errorf("Select by id returned %d targets, want 1", len(got))
	}
}

func TestCatalogIntegrity(t *testing.T) {
	seen := map[string]bool{}
	for _, tg := range Catalog {
		if tg.ID == "" || tg.Name == "" || tg.Host == "" {
			t.Errorf("target %+v has an empty required field", tg)
		}
		if seen[tg.ID] {
			t.Errorf("duplicate target id %q", tg.ID)
		}
		seen[tg.ID] = true
		if tg.Port == 0 {
			t.Errorf("target %q has no TCP port, so it cannot fall back from ICMP", tg.ID)
		}
	}
}
