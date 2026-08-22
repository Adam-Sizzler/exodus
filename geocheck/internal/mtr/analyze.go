package mtr

import (
	"fmt"
	"time"
)

// Class is the verdict on how a target is reached.
type Class int

const (
	// ClassUnknown means the path carried too little information to judge.
	ClassUnknown Class = iota
	// ClassUnreachable means nothing answered at all.
	ClassUnreachable
	// ClassIntercepted means something answered on the destination's behalf.
	ClassIntercepted
	// ClassDirect means the destination network is entered without crossing a
	// transit provider, at the best latency this connection achieves anywhere.
	ClassDirect
	// ClassPeered means the destination network is entered without transit,
	// but from further away: a regional exchange rather than an on-net cache.
	ClassPeered
	// ClassTransit means one or more transit carriers sit between you and the
	// destination network.
	ClassTransit
	// ClassDetour means the traffic travels far further than the destination
	// warrants, typically a tunnel exiting in another region.
	ClassDetour
)

func (c Class) String() string {
	switch c {
	case ClassUnreachable:
		return "unreachable"
	case ClassIntercepted:
		return "intercepted"
	case ClassDirect:
		return "direct"
	case ClassPeered:
		return "peered"
	case ClassTransit:
		return "transit"
	case ClassDetour:
		return "detour"
	default:
		return "unknown"
	}
}

// Label is a short human phrase for the verdict.
func (c Class) Label() string {
	switch c {
	case ClassUnreachable:
		return "no answer"
	case ClassIntercepted:
		return "intercepted"
	case ClassDirect:
		return "direct / on-net"
	case ClassPeered:
		return "regional"
	case ClassTransit:
		return "via transit"
	case ClassDetour:
		return "detour"
	default:
		return "inconclusive"
	}
}

// Verdict summarises what the measured path means.
type Verdict struct {
	Class Class
	// Score is 0..100, higher being a shorter and more direct path.
	Score int

	RTT time.Duration // best RTT to the destination
	// Excess is RTT above the best latency this connection achieved to any
	// target. It is the figure the classification actually turns on.
	Excess    time.Duration
	Loss      float64 // loss at the destination, 0..1
	Jitter    time.Duration
	HopCount  int          // responding hops up to the destination
	Transits  []TransitHop // transit carriers crossed, in path order
	Networks  int          // distinct autonomous systems on the path
	OnNet     bool         // the destination AS was entered without transit
	MaxJumpMs float64      // largest RTT increase between adjacent hops
	Notes     []string
}

// TransitHop names one transit carrier seen on the path.
type TransitHop struct {
	TTL  int
	ASN  int
	Name string
}

// Latency is interpreted relative to the floor: the lowest round trip this
// connection achieved to any target. That floor is the cost of the access
// network itself (fibre idles at 7-14 ms, cable at 12-24 ms, DSL higher), and
// no amount of good peering gets below it. Judging targets against a fixed
// millisecond threshold instead would label an entire well-connected metro as
// "not direct" purely because the last mile is DSL.
//
// The excess above the floor is what carries information, and these are the
// bands it falls into.
const (
	excessOnNetMs  = 5  // indistinguishable from your best-connected destination
	excessRegionMs = 20 // same country or a neighbouring metro
	excessContMs   = 60 // elsewhere on the continent
)

// interceptedMs is the round trip below which a reply deserves a second look.
// Light covers only about 100 km of fibre in 1 ms of round trip, and consumer
// access technologies cost several milliseconds before a packet leaves the
// building at all.
//
// It is NOT on its own a verdict. From a datacenter or a VPS, sub-millisecond
// replies are routine and are the best possible outcome: the content network
// has a cache in the same facility. Treating that as interception would
// mislabel exactly the hosts this tool is most often run on. So a low figure
// only counts as interception when the path also shows that nothing was
// traversed to earn it — see locallyAnswered.
const interceptedMs = 1.0

// jumpDetourMs is the inter-hop RTT increase that implies a long-haul link
// where the destination should not have needed one.
const jumpDetourMs = 40

// transitASNs are the carriers whose presence means your traffic is being
// carried by a third party rather than handed straight to the destination
// network. Ordered by CAIDA AS Rank customer-cone size.
var transitASNs = map[int]string{
	3356:  "Lumen (Level3)",
	1299:  "Arelion (Telia)",
	3257:  "GTT",
	174:   "Cogent",
	2914:  "NTT",
	6939:  "Hurricane Electric",
	6453:  "Tata",
	6461:  "Zayo",
	6762:  "Telecom Italia Sparkle",
	3491:  "PCCW",
	9002:  "RETN",
	5511:  "Orange",
	4637:  "Telstra Global",
	1273:  "Vodafone",
	12956: "Telxius",
	3320:  "Deutsche Telekom",
	12389: "Rostelecom",
	701:   "Verizon UUNET",
	7018:  "AT&T",
	1239:  "Sprint",
	4134:  "China Telecom",
	4809:  "China Telecom CN2",
	4837:  "China Unicom",
	58453: "China Mobile International",
	7473:  "Singtel",
	9498:  "Bharti Airtel",
	20485: "TransTelecom",
	31133: "MegaFon",
	6830:  "Liberty Global",
}

// TransitName returns the carrier name for an AS number, if it is a known
// transit provider.
func TransitName(num int) (string, bool) {
	n, ok := transitASNs[num]
	return n, ok
}

// Floor estimates what the access network costs, by taking the fastest
// destination this connection reached. It deliberately does not use the single
// minimum: one unusually close endpoint — a local cache, or a proxy answering
// early — would otherwise define the floor and push every honest target into
// looking distant. With enough samples the second-fastest is used instead, so
// the floor has to be corroborated by two independent destinations.
func Floor(reports []*Report) time.Duration {
	var rtts []time.Duration
	for _, r := range reports {
		if r == nil {
			continue
		}
		h, ok := r.FinalHop()
		if !ok {
			continue
		}
		best := h.Best()
		// A reply that something local produced is not a floor; letting one in
		// would drag every honest target's excess upwards. A genuinely fast
		// on-net cache, by contrast, belongs in the floor and defines it.
		if best <= 0 || locallyAnswered(r) {
			continue
		}
		rtts = append(rtts, best)
	}
	if len(rtts) == 0 {
		return 0
	}
	sortDurations(rtts)
	if len(rtts) >= 4 {
		return rtts[1]
	}
	return rtts[0]
}

func sortDurations(d []time.Duration) {
	for i := 1; i < len(d); i++ {
		for j := i; j > 0 && d[j] < d[j-1]; j-- {
			d[j], d[j-1] = d[j-1], d[j]
		}
	}
}

// locallyAnswered reports whether a reply was almost certainly produced nearby
// rather than by the destination.
//
// A sub-millisecond figure alone does not establish this. What does is a
// sub-millisecond figure with no evidence of travel behind it: the reply
// arrived at the first TTL, so nothing was traversed to reach it. That is the
// signature of a userspace tunnel or transparent proxy answering on the
// destination's behalf.
//
// When several hops did answer, the packets demonstrably crossed routers and
// came back, and sub-millisecond simply means the content sits in this
// building — normal, and excellent, on a datacenter host.
func locallyAnswered(r *Report) bool {
	// Only a traceroute can establish this. A TCP handshake never walks the
	// path at all, so the absence of intermediate hops is the method rather
	// than a finding — concluding "intercepted" from it would assert something
	// the measurement cannot know. On a datacenter host, where the raw socket
	// is often unavailable and every CDN answers a handshake in well under a
	// millisecond, that mistake condemns almost every target at once.
	//
	// Interception is still caught, just where the evidence actually lives: by
	// the tunnel and resolver checks in the detect package, and by ICMP traces
	// when a raw socket is available.
	if r.Method != MethodICMP {
		return false
	}

	h, ok := r.FinalHop()
	if !ok || h.Best() <= 0 {
		return false
	}
	if h.Best().Seconds()*1000 >= interceptedMs {
		return false
	}
	return respondingHops(r) <= 1
}

// respondingHops counts the hops that answered at all.
func respondingHops(r *Report) int {
	n := 0
	for _, h := range r.Hops {
		if h.Responded() {
			n++
		}
	}
	return n
}

// ClassifyAll judges every report against the connection's own latency floor.
func ClassifyAll(reports []*Report) {
	floor := Floor(reports)
	for _, r := range reports {
		if r != nil {
			r.Verdict = Classify(r, floor)
		}
	}
}

// Classify turns one measured path into a verdict, given the floor computed
// across the whole run.
func Classify(r *Report, floor time.Duration) Verdict {
	v := Verdict{Class: ClassUnknown}

	final, ok := r.FinalHop()
	if !ok {
		v.Class = ClassUnreachable
		v.Loss = 1
		v.Notes = append(v.Notes, "no hop answered; the path is fully filtered")
		return v
	}

	v.RTT = final.Best()
	v.Loss = final.Loss()
	v.Jitter = final.StdDev()
	rttMs := v.RTT.Seconds() * 1000

	if floor > 0 && v.RTT > floor {
		v.Excess = v.RTT - floor
	}
	excessMs := v.Excess.Seconds() * 1000

	reachedDest := final.Addr == r.Resolved
	if !reachedDest {
		v.Notes = append(v.Notes, "the destination never answered; the last responding hop is shown")
	}

	v.analysePath(r)

	// TCP probing gives one honest destination RTT but no path, so the transit
	// analysis simply does not apply.
	pathVisible := r.Method == MethodICMP && v.HopCount > 1

	switch {
	case locallyAnswered(r):
		v.Class = ClassIntercepted
		v.Notes = append(v.Notes, fmt.Sprintf(
			"%.2f ms with nothing answering in between: the reply came back before any "+
				"router was crossed, so something local produced it", rttMs))

	case !reachedDest && v.HopCount <= 1:
		v.Class = ClassUnknown
		v.Notes = append(v.Notes, "not enough responding hops to judge the route")

	case pathVisible && v.MaxJumpMs > jumpDetourMs:
		v.Class = ClassDetour
		v.Notes = append(v.Notes, fmt.Sprintf(
			"one hop adds %.0f ms, so a single long-haul link carries the whole route",
			v.MaxJumpMs))

	case len(v.Transits) > 0:
		v.Class = ClassTransit
		names := make([]string, 0, len(v.Transits))
		for _, t := range v.Transits {
			names = append(names, t.Name)
		}
		v.Notes = append(v.Notes, "carried by "+joinNames(names)+
			" rather than handed over directly")

	case floor == 0:
		// Nothing to compare against; fall back to raw latency.
		v.Class = absoluteClass(rttMs)

	case excessMs <= excessOnNetMs:
		v.Class = ClassDirect
		if rttMs < interceptedMs {
			// Worth spelling out, because sub-millisecond looks alarming until
			// you know the path behind it was real.
			v.Notes = append(v.Notes, fmt.Sprintf(
				"%.2f ms across %d hops: the content is cached inside this facility, "+
					"which is the best result available",
				rttMs, v.HopCount))
		} else {
			v.Notes = append(v.Notes,
				"as fast as this connection reaches anything, so the content sits in your "+
					"own network or metro")
		}

	case excessMs <= excessRegionMs:
		v.Class = ClassPeered
		v.Notes = append(v.Notes, fmt.Sprintf(
			"%.0f ms further than your best target: a national or neighbouring-metro handover",
			excessMs))

	case excessMs <= excessContMs:
		v.Class = ClassPeered
		v.Notes = append(v.Notes, fmt.Sprintf(
			"%.0f ms further than your best target, which puts the handover elsewhere on the continent",
			excessMs))

	default:
		v.Class = ClassDetour
		v.Notes = append(v.Notes, fmt.Sprintf(
			"%.0f ms further than your best target, roughly %.0f km of extra cable each way",
			excessMs, excessMs*100/2))
	}

	v.addContext(r, pathVisible)
	v.Score = score(v, rttMs, excessMs, floor, pathVisible)
	return v
}

// analysePath walks the hops, counting networks and noting transit carriers
// crossed before the destination network is entered.
func (v *Verdict) analysePath(r *Report) {
	seenASN := map[int]bool{}
	destASN := r.DestASN.Number
	if destASN == 0 {
		destASN = r.Target.ASN
	}
	enteredDest := false
	var prevRTT time.Duration

	for _, h := range r.Hops {
		if !h.Responded() {
			continue
		}
		v.HopCount++

		if best := h.Best(); best > 0 {
			if prevRTT > 0 {
				if jump := (best - prevRTT).Seconds() * 1000; jump > v.MaxJumpMs {
					v.MaxJumpMs = jump
				}
			}
			prevRTT = best
		}

		num := h.ASN.Number
		if num == 0 {
			continue
		}
		if !seenASN[num] {
			seenASN[num] = true
			v.Networks++
		}
		if destASN != 0 && num == destASN {
			enteredDest = true
			continue
		}
		if enteredDest {
			continue
		}
		if name, isTransit := transitASNs[num]; isTransit && !hasTransit(v.Transits, num) {
			v.Transits = append(v.Transits, TransitHop{TTL: h.TTL, ASN: num, Name: name})
		}
	}
	v.OnNet = len(v.Transits) == 0
}

// addContext appends the caveats that change how a verdict should be read.
func (v *Verdict) addContext(r *Report, pathVisible bool) {
	if !pathVisible && r.Method == MethodICMP && v.Class != ClassUnreachable {
		v.Notes = append(v.Notes,
			"intermediate hops are hidden, so only the destination latency is meaningful")
	}
	if r.Method == MethodTCP {
		v.Notes = append(v.Notes,
			"measured with a TCP handshake rather than ICMP; note that a local proxy or "+
				"tunnel completes the handshake itself, which makes this latency look "+
				"better than the real path")
	}
	if r.Target.Anycast {
		v.Notes = append(v.Notes,
			"anycast: this shows how near the closest announcement is, not where a server sits")
	}
	// A brand name resolving into a CDN is normal, not a fault, so this is
	// reported as context rather than counted against the score.
	if r.Target.ASN != 0 && r.DestASN.Number != 0 && r.DestASN.Number != r.Target.ASN {
		v.Notes = append(v.Notes, fmt.Sprintf(
			"served from %s rather than the target's own AS%d, which usually means a CDN",
			r.DestASN.String(), r.Target.ASN))
	}
	if v.Loss > 0.2 && v.Loss < 1 {
		v.Notes = append(v.Notes,
			fmt.Sprintf("%.0f%% of probes to the destination were lost", v.Loss*100))
	}
	if pathVisible {
		v.Notes = append(v.Notes,
			"loss shown at intermediate hops is usually router rate limiting; only loss "+
				"that persists to the destination is real")
	}
}

// absoluteClass is the fallback when no floor could be established.
func absoluteClass(rttMs float64) Class {
	switch {
	case rttMs <= 20:
		return ClassDirect
	case rttMs <= 60:
		return ClassPeered
	case rttMs <= 150:
		return ClassPeered
	default:
		return ClassDetour
	}
}

// score condenses the verdict into 0..100 so targets can be ranked and
// averaged into one headline number.
func score(v Verdict, rttMs, excessMs float64, floor time.Duration, pathVisible bool) int {
	if v.Class == ClassUnreachable {
		return 0
	}
	if v.Class == ClassIntercepted {
		return 0
	}

	s := 100

	if floor > 0 {
		switch {
		case excessMs <= excessOnNetMs:
		case excessMs <= excessRegionMs:
			s -= 10
		case excessMs <= excessContMs:
			s -= 30
		default:
			s -= 55
		}
	} else {
		switch {
		case rttMs <= 20:
		case rttMs <= 60:
			s -= 15
		case rttMs <= 150:
			s -= 35
		default:
			s -= 55
		}
	}

	s -= 12 * len(v.Transits)

	if pathVisible && v.MaxJumpMs > jumpDetourMs {
		s -= 15
	}
	if v.Loss > 0 && v.Loss < 1 {
		s -= int(v.Loss * 25)
	}

	if s < 0 {
		s = 0
	}
	if s > 100 {
		s = 100
	}
	return s
}

func hasTransit(list []TransitHop, num int) bool {
	for _, t := range list {
		if t.ASN == num {
			return true
		}
	}
	return false
}

func joinNames(names []string) string {
	switch len(names) {
	case 0:
		return ""
	case 1:
		return names[0]
	case 2:
		return names[0] + " and " + names[1]
	default:
		out := ""
		for i, n := range names[:len(names)-1] {
			if i > 0 {
				out += ", "
			}
			out += n
		}
		return out + " and " + names[len(names)-1]
	}
}

// Summary aggregates verdicts across all traced targets.
type Summary struct {
	Score       int
	Direct      int
	Peered      int
	Transit     int
	Detour      int
	Intercepted int
	Failed      int
	Total       int
	Floor       time.Duration
	MedianRTT   time.Duration
}

// Summarize folds a set of reports into one headline.
func Summarize(reports []*Report) Summary {
	var s Summary
	var total, scored int
	var rtts []time.Duration

	s.Floor = Floor(reports)

	for _, r := range reports {
		if r == nil {
			continue
		}
		s.Total++
		switch r.Verdict.Class {
		case ClassDirect:
			s.Direct++
		case ClassPeered:
			s.Peered++
		case ClassTransit:
			s.Transit++
		case ClassDetour:
			s.Detour++
		case ClassIntercepted:
			s.Intercepted++
		default:
			s.Failed++
		}
		if r.Verdict.Class != ClassUnreachable && r.Verdict.Class != ClassUnknown {
			total += r.Verdict.Score
			scored++
			if r.Verdict.RTT > 0 {
				rtts = append(rtts, r.Verdict.RTT)
			}
		}
	}

	if scored > 0 {
		s.Score = total / scored
	}
	if len(rtts) > 0 {
		for i := 1; i < len(rtts); i++ {
			for j := i; j > 0 && rtts[j] < rtts[j-1]; j-- {
				rtts[j], rtts[j-1] = rtts[j-1], rtts[j]
			}
		}
		s.MedianRTT = rtts[len(rtts)/2]
	}
	return s
}
