package render

import (
	"bytes"
	"net/netip"
	"testing"
)

// documentationRanges are the prefixes reserved for use in documentation:
// RFC 5737 for IPv4 and RFC 3849 for IPv6. Nothing in them is routable, so an
// example built from these cannot name a real host.
var documentationRanges = []netip.Prefix{
	netip.MustParsePrefix("192.0.2.0/24"),    // TEST-NET-1
	netip.MustParsePrefix("198.51.100.0/24"), // TEST-NET-2
	netip.MustParsePrefix("203.0.113.0/24"),  // TEST-NET-3
	netip.MustParsePrefix("2001:db8::/32"),
}

// TestDemoUsesOnlyDocumentationAddresses is the guard that matters. The whole
// point of the demo report is that it can be published — recorded into a GIF,
// pasted into the README — without exposing anyone's address or route. If a
// real address ever reaches it, that protection is silently gone, so every
// address the report carries is checked rather than a sampled few.
func TestDemoUsesOnlyDocumentationAddresses(t *testing.T) {
	r := DemoReport("test")

	seen := 0
	check := func(what string, a netip.Addr) {
		t.Helper()
		if !a.IsValid() {
			return
		}
		seen++
		for _, p := range documentationRanges {
			if p.Contains(a) {
				return
			}
		}
		t.Errorf("%s is %s, which is outside documentation space", what, a)
	}

	check("identity IPv4", r.Identity.IPv4)
	check("identity IPv6", r.Identity.IPv6)
	if r.Reputation != nil {
		check("reputation address", r.Reputation.IP)
	}
	for _, rep := range r.Trace {
		check("resolved "+rep.Target.ID, rep.Resolved)
		for _, h := range rep.Hops {
			check("hop in "+rep.Target.ID, h.Addr)
			for _, a := range h.Addrs {
				check("hop address in "+rep.Target.ID, a)
			}
		}
	}

	if seen < 20 {
		t.Fatalf("only %d addresses examined; the report looks empty, so this test proves nothing", seen)
	}
}

// TestDemoASNIsReserved keeps the invented operator from being attributed to a
// real company. RFC 5398 reserves 64496-64511 for documentation.
func TestDemoASNIsReserved(t *testing.T) {
	got := DemoReport("test").Identity.ASN.Number
	if got < 64496 || got > 64511 {
		t.Errorf("identity ASN is AS%d, which belongs to a real operator; use RFC 5398 space", got)
	}
}

// TestDemoIsDeterministic means regenerating the recorded demo produces no diff
// unless the report genuinely changed.
func TestDemoIsDeterministic(t *testing.T) {
	render := func() []byte {
		var buf bytes.Buffer
		// A zero Theme renders every style as plain text, which is all this
		// test needs and keeps it independent of colour detection.
		o := &Output{W: &buf, Theme: Theme{}, Width: 100}
		o.Print(DemoReport("test"))
		return buf.Bytes()
	}
	if a, b := render(), render(); !bytes.Equal(a, b) {
		t.Error("two renders of the demo report differ")
	}
}

// TestDemoFillsEverySection stops the demo from quietly losing a section: an
// empty slice renders as nothing at all, which would look intentional.
func TestDemoFillsEverySection(t *testing.T) {
	r := DemoReport("test")
	for _, c := range []struct {
		name string
		n    int
	}{
		{"geo", len(r.Geo)},
		{"portal", len(r.Portal)},
		{"access", len(r.Access)},
		{"trace", len(r.Trace)},
	} {
		if c.n == 0 {
			t.Errorf("the %s section is empty", c.name)
		}
	}
	if r.Reputation == nil {
		t.Error("the reputation section is missing")
	}
}
