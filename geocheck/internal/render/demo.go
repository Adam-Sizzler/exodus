package render

import (
	"net/netip"
	"time"

	"exodus-node/geocheck/internal/access"
	"exodus-node/geocheck/internal/ai"
	"exodus-node/geocheck/internal/asn"
	"exodus-node/geocheck/internal/geo"
	"exodus-node/geocheck/internal/mtr"
	"exodus-node/geocheck/internal/netx"
	"exodus-node/geocheck/internal/portal"
	"exodus-node/geocheck/internal/reputation"
)

// DemoReport builds a complete report from invented measurements, so the
// documentation can show what a run looks like without publishing anyone's
// address, network or route. Every address here is documentation space —
// RFC 5737 for IPv4, RFC 3849 for IPv6, RFC 5398 for the AS number — which is
// reserved precisely so that examples cannot collide with a real host.
//
// It is assembled from the real catalogues rather than a hardcoded list of
// names, and the path verdicts come from the real classifier, so adding a
// check or changing the analysis shows up here instead of quietly drifting
// out of date.
func DemoReport(version string) Report {
	families := []netx.Family{netx.V4, netx.V6}

	report := Report{
		Version:  version,
		Identity: demoIdentity(),
		Resolver: "https://cloudflare-dns.com/dns-query",
		Families: families,
		Geo:      demoGeo(),
		Portal:   demoPortal(),
		Access:   demoAccess(),
		AI:       demoAI(),
		Trace:    demoTrace(),
		TraceCap: mtr.Capability{ICMP: true, Raw: true, PathVisible: true},
		Duration: 6*time.Second + 200*time.Millisecond,
	}
	report.Reputation = demoReputation(report.Identity.IPv4)
	return report
}

func demoIdentity() geo.Identity {
	return geo.Identity{
		IPv4: netip.MustParseAddr("198.51.100.34"),
		IPv6: netip.MustParseAddr("2001:db8:4f2a::34"),
		ASN: asn.Info{
			Number:  64496, // RFC 5398 documentation ASN
			Name:    "EXAMPLE-AS",
			Country: "NL",
			Prefix:  "198.51.100.0/24",
		},
		Org: "Example Networks B.V.",
	}
}

func demoReputation(ip netip.Addr) *reputation.Info {
	return &reputation.Info{
		IP:           ip,
		Type:         "Hosting",
		ASN:          "AS64496",
		Range:        "198.51.100.0/24",
		Hostname:     "static.198-51-100-34.example.net",
		Provider:     "Example Networks",
		Organisation: "Example Networks B.V.",
		City:         "Amsterdam",
		Region:       "North Holland",
		Country:      "Netherlands",
		Code:         "NL",
		Risk:         13,
		Confidence:   91,
		Hosting:      true,
		FirstSeen:    "2019-03-11",
		LastSeen:     "2026-08-17",
	}
}

// demoGeo answers every real check. Most agree on the Netherlands; a few
// disagree, err or are skipped, because a screenshot in which everything
// succeeds teaches nothing about how disagreement is displayed.
func demoGeo() []geo.Result {
	checks := append(geo.ServiceChecks(), geo.DatabaseChecks()...)
	checks = append(checks, geo.CDNChecks()...)

	out := make([]geo.Result, 0, len(checks))
	for i, c := range checks {
		out = append(out, geo.Result{
			Check: c,
			V4:    demoOutcome(c, i, false),
			V6:    demoOutcome(c, i, true),
		})
	}
	return out
}

func demoOutcome(c geo.Check, i int, v6 bool) geo.Outcome {
	switch c.Kind {
	case geo.KindAvailability:
		// One service in the set refuses hosting ranges, which is the single
		// most common real finding on a server and worth showing.
		if i%11 == 4 {
			return geo.Outcome{Value: "no"}
		}
		return geo.Outcome{Value: "yes"}
	case geo.KindBlocked:
		return geo.Outcome{Value: "no"}
	}

	// IPv6 is not offered by every provider, so some cells are legitimately
	// blank rather than wrong.
	if v6 && i%7 == 3 {
		return geo.Outcome{Skipped: true}
	}
	switch {
	case i%13 == 6:
		return geo.Outcome{Value: "DE"}
	case i%17 == 9:
		return geo.Outcome{Value: "GB"}
	}
	return geo.Outcome{Value: "NL"}
}

func demoPortal() []portal.Result {
	eps := portal.DefaultEndpoints()
	out := make([]portal.Result, 0, len(eps))
	for i, ep := range eps {
		res := portal.Result{
			Endpoint: ep,
			Verdict:  portal.VerdictOK,
			Status:   ep.WantStatus,
			Body:     ep.WantBody,
			RTT:      time.Duration(8+3*i) * time.Millisecond,
		}
		out = append(out, res)
	}
	return out
}

func demoAccess() []access.Result {
	checks := access.Checks()
	out := make([]access.Result, 0, len(checks))
	for i, c := range checks {
		res := access.Result{
			Check:  c,
			State:  access.StateAvailable,
			Region: "NL",
			RTT:    time.Duration(48+29*i) * time.Millisecond,
		}
		// A hosting address typically clears some services and not others.
		switch i % 5 {
		case 2:
			res.State, res.Detail, res.Region = access.StateRestricted, "Originals only", ""
		case 4:
			res.State, res.Detail, res.Region = access.StateBlocked, "Disallowed ISP", ""
		}
		out = append(out, res)
	}
	return out
}

// demoAI answers every real endpoint. One is refused so the blocked state is
// visible; the rest answer with the authentication error that proves the
// request arrived.
func demoAI() []ai.Result {
	checks := ai.Checks()
	out := make([]ai.Result, 0, len(checks))
	for i, c := range checks {
		res := ai.Result{
			Check: c, State: ai.StateReachable, Status: 401,
			Detail: "answered, credentials not supplied",
			RTT:    time.Duration(60+21*i) * time.Millisecond,
		}
		if i%7 == 3 {
			res.State, res.Status = ai.StateBlocked, 403
			res.Detail = "the API refused this region"
		}
		out = append(out, res)
	}
	return out
}

// demoTrace invents a path per target and then hands the lot to the real
// classifier, so the verdict column is genuinely computed rather than typed in.
func demoTrace() []*mtr.Report {
	targets := mtr.DefaultTargets()
	reports := make([]*mtr.Report, 0, len(targets))

	for i, t := range targets {
		reports = append(reports, demoPath(t, i))
	}
	mtr.ClassifyAll(reports)
	return reports
}

// demoProfile is one invented route: how far away the destination is, and
// whether a transit carrier sits in the middle.
type demoProfile struct {
	baseMs  float64 // RTT at the destination
	transit bool
}

// The floor of this connection is the fastest of these, so the profiles are
// chosen to span the classifier's bands: on-net, regional, transit-carried and
// a long detour.
var demoProfiles = []demoProfile{
	{baseMs: 3.4},
	{baseMs: 4.1},
	{baseMs: 11.6},
	{baseMs: 27.4, transit: true},
	{baseMs: 38.2, transit: true},
	{baseMs: 121.5, transit: true},
}

func demoPath(t mtr.Target, i int) *mtr.Report {
	p := demoProfiles[i%len(demoProfiles)]
	dest := demoDest(i)

	destAS := asn.Info{Number: t.ASN, Name: demoASName(t), Country: "NL"}

	hops := []mtr.Hop{
		demoHop(1, "192.0.2.1", "gw.example.net", asn.Info{Number: 64496, Name: "EXAMPLE-AS"}, 0.4),
		demoHop(2, "192.0.2.9", "core1.ams.example.net", asn.Info{Number: 64496, Name: "EXAMPLE-AS"}, 0.9),
	}

	ttl := 3
	if p.transit {
		hops = append(hops,
			demoHop(ttl, "203.0.113.17", "ae0.ams.transit.example.com",
				asn.Info{Number: 64500, Name: "EXAMPLE-TRANSIT"}, p.baseMs*0.45))
		ttl++
		hops = append(hops,
			demoHop(ttl, "203.0.113.42", "ae3.fra.transit.example.com",
				asn.Info{Number: 64500, Name: "EXAMPLE-TRANSIT"}, p.baseMs*0.8))
		ttl++
	} else {
		hops = append(hops,
			demoHop(ttl, "198.51.100.129", "ix-ams.example.net",
				asn.Info{Number: 64496, Name: "EXAMPLE-AS"}, p.baseMs*0.6))
		ttl++
	}

	hops = append(hops, demoHop(ttl, dest.String(), "", destAS, p.baseMs))

	return &mtr.Report{
		Target:   t,
		Resolved: dest,
		DestASN:  destAS,
		Method:   mtr.MethodICMP,
		Hops:     hops,
	}
}

// demoDest hands each target a distinct address out of the TEST-NET-1 range.
func demoDest(i int) netip.Addr {
	return netip.AddrFrom4([4]byte{192, 0, 2, byte(200 + i)})
}

func demoASName(t mtr.Target) string {
	if t.Net != "" {
		return t.Net
	}
	return "EXAMPLE-DEST"
}

// demoHop fills a hop with five probes clustered around a mean, so the jitter
// and loss columns have something realistic to show.
func demoHop(ttl int, addr, host string, info asn.Info, ms float64) mtr.Hop {
	a := netip.MustParseAddr(addr)
	spread := []float64{-0.06, 0.02, -0.01, 0.09, 0.03}

	rtts := make([]time.Duration, 0, len(spread))
	for _, s := range spread {
		rtts = append(rtts, time.Duration((ms+ms*s)*float64(time.Millisecond)))
	}

	return mtr.Hop{
		TTL:   ttl,
		Addr:  a,
		Host:  host,
		ASN:   info,
		Addrs: []netip.Addr{a},
		Sent:  len(spread),
		Recv:  len(spread),
		RTTs:  rtts,
	}
}
