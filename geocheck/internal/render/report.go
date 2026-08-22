package render

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/netip"
	"strconv"
	"strings"
	"time"

	lg "charm.land/lipgloss/v2"
	"charm.land/lipgloss/v2/table"
	"exodus-node/geocheck/internal/access"
	"exodus-node/geocheck/internal/ai"
	"exodus-node/geocheck/internal/countries"
	"exodus-node/geocheck/internal/detect"
	"exodus-node/geocheck/internal/geo"
	"exodus-node/geocheck/internal/mtr"
	"exodus-node/geocheck/internal/netx"
	"exodus-node/geocheck/internal/portal"
	"exodus-node/geocheck/internal/reputation"
)

// Report is everything one run produced.
type Report struct {
	Version   string
	Identity  geo.Identity
	Resolver  string
	Interface string
	Proxy     string
	Families  []netx.Family
	Geo       []geo.Result
	// Reputation is what proxycheck.io says about the exit address.
	Reputation    *reputation.Info
	ReputationErr error
	Portal        []portal.Result
	Access        []access.Result
	AI            []ai.Result
	Trace         []*mtr.Report
	TraceCap      mtr.Capability
	Duration      time.Duration

	// MaskIP hides the last octets of the public address in the header.
	MaskIP bool
	// TraceDetail prints the full hop table for every target.
	TraceDetail bool
	// EmbedSVG asks the JSON document to carry the rendered picture as well,
	// so a consumer that wants both gets one parseable stream instead of two.
	EmbedSVG bool
}

// PrintFindings reports conditions that undermine the measurement itself.
// They come first because they change how everything below should be read.
func (o *Output) PrintFindings(findings []detect.Finding) {
	if len(findings) == 0 {
		return
	}
	t := o.Theme
	o.line("")
	o.section("Measurement integrity")
	for _, f := range findings {
		style, mark := t.Info, "i"
		switch f.Severity {
		case detect.Alert:
			style, mark = t.Bad, "!"
		case detect.Warn:
			style, mark = t.Warn, "!"
		}
		o.line("  %s %s", style.Render(mark), t.Value.Render(f.Title))
		o.line("    %s", t.Muted.Render(wrap(f.Detail, o.Width-6)))
	}
}

// wrap reflows text to a column width, indenting continuation lines.
func wrap(s string, width int) string {
	if width < 20 {
		width = 20
	}
	var out strings.Builder
	col := 0
	for i, word := range strings.Fields(s) {
		if i > 0 {
			if col+1+len(word) > width {
				out.WriteString("\n    ")
				col = 0
			} else {
				out.WriteString(" ")
				col++
			}
		}
		out.WriteString(word)
		col += len(word)
	}
	return out.String()
}

// Print writes the whole human-readable report.
//
// The three geolocation tables are narrow, so when the terminal is wide enough
// they are laid out in two columns rather than stacked. The report already
// needs 115 columns for its widest table, and pairing needs 122, so this asks
// for very little more than it asked for already — and only takes it when the
// terminal actually has it.
//
// The check is against the measured width of those tables rather than a fixed
// number, so a table growing a column cannot quietly push the layout past what
// fits.
func (o *Output) Print(r Report) {
	if text, ok := o.pairedReport(r); ok {
		_, _ = io.WriteString(o.W, text)
		return
	}
	o.printStacked(r)
}

// pairedReport renders the report, pairs the geolocation tables, and reports
// whether the result fits. It renders into a buffer first because the decision
// needs the tables' widths, which are only known once they are drawn.
func (o *Output) pairedReport(r Report) (string, bool) {
	var buf bytes.Buffer
	inner := &Output{W: &buf, Theme: o.Theme, Color: o.Color, Width: o.Width}
	inner.printStacked(r)

	lines := ansiLines(buf.String())
	paired, ok := pairGeoTablesWithin(lines, o.Width)
	if !ok {
		return "", false
	}
	return ansiText(paired), true
}

func (o *Output) printStacked(r Report) {
	o.banner(r)
	o.identity(r)

	if r.Reputation != nil || r.ReputationErr != nil {
		o.reputation(r)
	}

	if len(r.Geo) > 0 {
		o.consensus(r)
		o.geoTable(r, "Popular services", geo.GroupServices)
		o.geoTable(r, "GeoIP databases", geo.GroupGeoIP)
		o.geoTable(r, "CDN edge location", geo.GroupCDN)
	}

	if len(r.Portal) > 0 {
		o.portalTable(r)
	}

	if len(r.Trace) > 0 {
		o.traceSummary(r)
		if r.TraceDetail {
			o.traceDetail(r)
		}
	}

	if len(r.Access) > 0 {
		o.accessTable(r)
	}

	if len(r.AI) > 0 {
		o.aiTable(r)
	}

	o.footer(r)
}

func (o *Output) line(format string, args ...any) {
	fmt.Fprintf(o.W, format+"\n", args...)
}

const (
	projectURL      = "github.com/adam-sizzler/exodus"
	projectShortURL = "github.com/adam-sizzler/exodus"
)

func (o *Output) banner(r Report) {
	t := o.Theme
	o.line("")
	o.line("  %s %s", t.Title.Render("geocheck"), t.Muted.Render(r.Version))
	o.line("  %s", t.Subtitle.Render("where the internet thinks you are, and how directly you reach it"))
	o.line("  %s %s %s %s",
		t.Muted.Render("↗"),
		t.Accent.Render(projectURL),
		t.Muted.Render("·"),
		t.Accent.Render(projectShortURL))
	o.line("")
}

func (o *Output) identity(r Report) {
	t := o.Theme
	id := r.Identity

	row := func(label, value string) {
		if value == "" {
			return
		}
		o.line("  %s  %s", t.Label.Render(pad(label, 10)), value)
	}

	if id.IPv4.IsValid() {
		row("IPv4", t.Value.Render(maskAddr(id.IPv4, r.MaskIP)))
	}
	if id.IPv6.IsValid() {
		row("IPv6", t.Value.Render(maskAddr(id.IPv6, r.MaskIP)))
	}
	if !id.ASN.Empty() {
		org := id.Org
		if org == "" || strings.EqualFold(org, id.ASN.Name) {
			org = id.ASN.Name
		}
		row("Network", t.Value.Render(fmt.Sprintf("AS%d", id.ASN.Number))+"  "+t.Muted.Render(org))
	}

	var env []string
	if r.Interface != "" {
		env = append(env, "interface "+r.Interface)
	}
	if r.Proxy != "" {
		env = append(env, "socks5 "+r.Proxy)
	}
	if r.Resolver != "" {
		env = append(env, "DoH "+r.Resolver)
	}
	if len(env) > 0 {
		row("Via", t.Muted.Render(strings.Join(env, " · ")))
	}
	o.line("")
}

// reputation shows how the exit address is classified. It sits at the top
// because it usually explains the rest: a service that geolocates you correctly
// and still refuses is almost always refusing the address type, not the place.
func (o *Output) reputation(r Report) {
	t := o.Theme
	o.section("Address reputation")

	if r.ReputationErr != nil {
		msg := "unavailable: " + r.ReputationErr.Error()
		if errors.Is(r.ReputationErr, reputation.ErrQuotaExceeded) {
			msg += ". The unregistered allowance is 100 lookups a day per source " +
				"address; a free key raises it to 1000 — pass it with --proxycheck-key " +
				"or set PROXYCHECK_API_KEY."
		}
		o.line("  %s %s", t.Muted.Render("·"), t.Muted.Render(wrap(msg, o.Width-6)))
		o.line("")
		return
	}

	info := r.Reputation
	row := func(label, value string) {
		if value == "" {
			return
		}
		o.line("  %s  %s", t.Label.Render(pad(label, 10)), value)
	}

	// Type first: it is the field that most often explains a refusal.
	if info.Type != "" {
		style, gloss := t.Good, "end-user address space"
		if !info.Residential() {
			style, gloss = t.Warn, "datacenter space; many consumer services refuse it"
		}
		row("Type", style.Render(info.Type)+"  "+t.Muted.Render(gloss))
	}

	riskStyle := t.Good
	switch {
	case info.Risk >= 66:
		riskStyle = t.Bad
	case info.Risk >= 33:
		riskStyle = t.Warn
	}
	risk := riskStyle.Render(fmt.Sprintf("%d/100", info.Risk)) + "  " +
		t.Muted.Render(meter(float64(info.Risk), 16))
	if info.Confidence > 0 {
		risk += t.Muted.Render(fmt.Sprintf("  confidence %d%%", info.Confidence))
	}
	row("Risk", risk)

	if flags := info.Flags(); len(flags) > 0 {
		styled := make([]string, 0, len(flags))
		for _, f := range flags {
			switch f {
			case "hosting":
				styled = append(styled, t.Warn.Render(f))
			case "anonymous":
				styled = append(styled, t.Warn.Render(f))
			default:
				styled = append(styled, t.Bad.Render(f))
			}
		}
		row("Flagged", strings.Join(styled, t.Muted.Render(" · ")))
	} else {
		row("Flagged", t.Good.Render("nothing"))
	}

	if info.Operator != "" {
		op := t.Value.Render(info.Operator)
		var extra []string
		if info.Anonymity != "" {
			extra = append(extra, info.Anonymity+" anonymity")
		}
		if info.HasPolicies {
			if info.NoLogging {
				extra = append(extra, "no logging")
			} else {
				extra = append(extra, "keeps logs")
			}
		}
		if len(extra) > 0 {
			op += "  " + t.Muted.Render("("+strings.Join(extra, ", ")+")")
		}
		row("Operator", op)
	}

	if place := joinNonEmpty([]string{info.City, info.Region, info.Country}, ", "); place != "" {
		row("Placed at", t.Value.Render(place))
	}

	// A high count on one address is CGNAT; a high count across the subnet with
	// a low one here is ordinary shared hosting. A zero means "not estimated",
	// so it is left out rather than printed as a misleading "~0".
	switch {
	case info.DevicesAddr > 0:
		row("Devices", t.Muted.Render(fmt.Sprintf(
			"~%d behind this address, ~%d across %s",
			info.DevicesAddr, info.DevicesSubnt, orDash(info.Range))))
	case info.DevicesSubnt > 0:
		row("Devices", t.Muted.Render(fmt.Sprintf(
			"~%d across %s", info.DevicesSubnt, orDash(info.Range))))
	}

	if info.LastSeen != "" {
		seen := "last flagged " + info.LastSeen
		if info.FirstSeen != "" {
			seen += ", first " + info.FirstSeen
		}
		row("History", t.Muted.Render(seen))
	}

	if info.Warning != "" {
		o.line("  %s %s", t.Warn.Render("!"), t.Muted.Render(info.Warning))
	}
	o.line("")
}

func joinNonEmpty(parts []string, sep string) string {
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if strings.TrimSpace(p) != "" {
			out = append(out, p)
		}
	}
	return strings.Join(out, sep)
}

func orDash(s string) string {
	if s == "" {
		return "the subnet"
	}
	return s
}

// consensus shows how the answers distribute across countries, which is the
// fastest way to spot a leak: one dissenting provider stands out immediately.
func (o *Output) consensus(r Report) {
	t := o.Theme
	for _, f := range r.Families {
		rows := geo.Summarize(r.Geo, f)
		if len(rows) == 0 {
			continue
		}
		o.section(fmt.Sprintf("Consensus (%s)", f))
		for _, c := range rows {
			bar := meter(c.Percent, 18)
			style := t.Good
			if c.Percent < 50 {
				style = t.Warn
			}
			if c.Percent < 20 {
				style = t.Muted
			}
			o.line("  %s %s  %s %s",
				t.Muted.Render(pad(c.Code, 4)),
				t.Value.Render(pad(c.Name, 26)),
				style.Render(bar),
				t.Muted.Render(fmt.Sprintf("%3.0f%%  %d/%d", c.Percent, c.Count, c.Total)))
		}
		o.line("")
	}
}

func (o *Output) geoTable(r Report, title string, group geo.Group) {
	rows := make([][]string, 0, len(r.Geo))
	for _, res := range r.Geo {
		if res.Check.Group != group {
			continue
		}
		row := make([]string, 0, 1+len(r.Families))
		row = append(row, res.Check.Name)
		for _, f := range r.Families {
			out := res.V4
			if f == netx.V6 {
				out = res.V6
			}
			row = append(row, o.outcomeCell(out, res.Check.Kind))
		}
		rows = append(rows, row)
	}
	if len(rows) == 0 {
		return
	}

	headers := []string{"Service"}
	for _, f := range r.Families {
		headers = append(headers, f.String())
	}

	o.section(title)
	o.line("%s", o.table(headers, rows))
	o.line("")
}

// outcomeCell formats one provider answer, colouring it by what it means.
func (o *Output) outcomeCell(out geo.Outcome, kind geo.Kind) string {
	t := o.Theme
	switch {
	case out.Skipped:
		return t.Muted.Render("—")
	case out.Err != nil:
		return t.Muted.Render("err")
	case out.Value == "":
		return t.Muted.Render("n/a")
	}

	switch kind {
	case geo.KindAvailability:
		if out.Value == "yes" {
			return t.Good.Render("available")
		}
		return t.Bad.Render("blocked")
	case geo.KindBlocked:
		if out.Value == "yes" {
			return t.Bad.Render("captcha")
		}
		return t.Good.Render("clean")
	default:
		return t.Muted.Render(out.Value) + "  " + t.Value.Render(countries.Name(out.Value))
	}
}

// portalTable shows the operating-system connectivity checks. Each of these
// endpoints has a contractual answer, so a deviation is proof that something
// rewrote the response rather than a matter of interpretation.
func (o *Output) portalTable(r Report) {
	t := o.Theme
	o.section("Connectivity checks")

	rows := make([][]string, 0, len(r.Portal))
	for _, res := range r.Portal {
		rows = append(rows, []string{
			res.Endpoint.Name,
			o.portalVerdictCell(res.Verdict),
			o.portalExpectedCell(res.Endpoint),
			o.portalGotCell(res),
			o.portalTimeCell(res),
		})
	}
	o.line("%s", o.table([]string{"Endpoint", "Result", "Expected", "Got", "Time"}, rows))

	s := portal.Summarize(r.Portal)
	switch {
	case s.Portal > 0:
		o.line("  %s %s", t.Bad.Render("!"), t.Value.Render(
			"A captive portal is intercepting HTTP. Nothing below reflects the real internet."))
	case s.Altered > 0:
		o.line("  %s %s", t.Warn.Render("!"), t.Value.Render(
			"Responses were rewritten in flight, so a transparent HTTP proxy sits on the path."))
	case s.OK == 0:
		o.line("  %s %s", t.Warn.Render("!"), t.Muted.Render(
			"No endpoint answered at all, so nothing can be concluded about interception."))
	case s.PlainHTTPBlocked():
		o.line("  %s %s", t.Warn.Render("!"), t.Value.Render(
			"Plain HTTP is being dropped: every port-80 endpoint timed out while HTTPS answered."))
		o.line("    %s", t.Muted.Render(
			"Traffic is not being rewritten, it is being discarded — a tunnel that forwards "+
				"only TLS, or a firewall blocking port 80, behaves this way."))
	default:
		msg := fmt.Sprintf("%d of %d endpoints answered exactly as specified", s.OK, s.Total)
		if s.PlainOK > 0 {
			msg += "; plain HTTP arrives unmodified, so nothing is rewriting it"
		}
		o.line("  %s %s", t.Good.Render("✓"), t.Muted.Render(msg+"."))
		if s.PlainFail > 0 {
			o.line("    %s", t.Muted.Render(fmt.Sprintf(
				"%d port-80 endpoint(s) did not answer; that is reachability, not interception.",
				s.PlainFail)))
		}
	}

	// Each deviation is worth spelling out; a bare "altered" is not actionable.
	for _, res := range r.Portal {
		if res.Verdict == portal.VerdictOK || res.Detail == "" {
			continue
		}
		o.line("    %s %s %s",
			t.Muted.Render("·"), t.Muted.Render(res.Endpoint.Name+":"), t.Muted.Render(res.Detail))
	}
	o.line("")
}

func (o *Output) portalVerdictCell(v portal.Verdict) string {
	t := o.Theme
	switch v {
	case portal.VerdictOK:
		return t.Good.Render("● " + v.Label())
	case portal.VerdictPortal:
		return t.Bad.Render("● " + v.Label())
	case portal.VerdictAltered:
		return t.Warn.Render("● " + v.Label())
	default:
		return t.Muted.Render("○ " + v.Label())
	}
}

func (o *Output) portalExpectedCell(ep portal.Endpoint) string {
	t := o.Theme
	switch {
	case ep.BodyContains != "":
		return t.Muted.Render(fmt.Sprintf("%d + %q", ep.WantStatus, clip(ep.BodyContains, 28)))
	case ep.WantBody == "":
		return t.Muted.Render(fmt.Sprintf("%d, empty", ep.WantStatus))
	default:
		return t.Muted.Render(fmt.Sprintf("%d + %q",
			ep.WantStatus, clip(strings.TrimSpace(ep.WantBody), 28)))
	}
}

// clip shortens a string for a table cell. Apple's expected body is 69 bytes of
// HTML, which would otherwise set the column width for the whole table.
func clip(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n-1] + "…"
}

// portalTimeCell hides the elapsed time for a failed request, which is just the
// timeout being reported back and would read as if something were measured.
func (o *Output) portalTimeCell(res portal.Result) string {
	if res.Err != nil {
		return o.Theme.Muted.Render("—")
	}
	return o.rttCell(res.RTT)
}

func (o *Output) portalGotCell(res portal.Result) string {
	t := o.Theme
	if res.Err != nil {
		return t.Muted.Render("—")
	}
	got := fmt.Sprint(res.Status)
	if res.Body != "" {
		got += " + " + strconv.Quote(clip(res.Body, 28))
	}
	if res.Verdict == portal.VerdictOK {
		return t.Muted.Render(got)
	}
	return t.Warn.Render(got)
}

// traceSummary is the headline answer to "is my connection direct".
func (o *Output) traceSummary(r Report) {
	t := o.Theme
	o.section("Direct connectivity")

	if r.TraceCap.Hint != "" {
		o.line("  %s %s", t.Warn.Render("!"), t.Muted.Render(r.TraceCap.Hint))
		o.line("")
	}

	rows := make([][]string, 0, len(r.Trace))
	for _, tr := range r.Trace {
		if tr == nil {
			continue
		}
		v := tr.Verdict
		hops := "—"
		if tr.Method == mtr.MethodICMP && v.HopCount > 0 {
			hops = fmt.Sprint(v.HopCount)
		}
		rows = append(rows, []string{
			tr.Target.Name,
			o.verdictCell(v.Class),
			o.rttCell(v.RTT),
			o.lossCell(v.Loss),
			hops,
			o.pathCell(tr),
		})
	}
	if len(rows) == 0 {
		return
	}

	o.line("%s", o.table(
		[]string{"Target", "Verdict", "RTT", "Loss", "Hops", "Route"}, rows))
	o.line("")

	s := mtr.Summarize(r.Trace)
	o.line("  %s  %s %s",
		t.Label.Render(pad("Directness", 10)),
		o.scoreStyle(s.Score).Render(fmt.Sprintf("%d/100", s.Score)),
		t.Muted.Render(meter(float64(s.Score), 20)))

	var parts []string
	add := func(n int, label string, style lg.Style) {
		if n > 0 {
			parts = append(parts, style.Render(fmt.Sprintf("%d %s", n, label)))
		}
	}
	add(s.Direct, "direct", t.Good)
	add(s.Peered, "regional", t.Info)
	add(s.Transit, "via transit", t.Warn)
	add(s.Detour, "detour", t.Bad)
	add(s.Intercepted, "intercepted", t.Bad)
	add(s.Failed, "no answer", t.Muted)
	if len(parts) > 0 {
		o.line("  %s  %s", t.Label.Render(pad("Breakdown", 10)), strings.Join(parts, t.Muted.Render(" · ")))
	}
	if s.MedianRTT > 0 {
		o.line("  %s  %s", t.Label.Render(pad("Median RTT", 10)), t.Value.Render(formatRTT(s.MedianRTT)))
	}
	if s.Floor > 0 {
		o.line("  %s  %s %s",
			t.Label.Render(pad("Baseline", 10)),
			t.Value.Render(formatRTT(s.Floor)),
			t.Muted.Render("— the fastest this connection reached anything; "+
				"verdicts measure distance from here"))
	}
	o.line("")
}

// traceDetail prints the per-hop table, the part that shows *why* a verdict
// came out the way it did.
func (o *Output) traceDetail(r Report) {
	t := o.Theme
	for _, tr := range r.Trace {
		if tr == nil || len(tr.Hops) == 0 {
			continue
		}
		o.section("Path to " + tr.Target.Name)
		o.line("  %s %s   %s",
			t.Muted.Render("resolved"),
			t.Value.Render(tr.Resolved.String()),
			t.Muted.Render(tr.DestASN.String()))

		if tr.Method == mtr.MethodTCP {
			o.line("  %s", t.Muted.Render("TCP handshake only; this host does not answer ICMP"))
		}

		rows := make([][]string, 0, len(tr.Hops))
		for _, h := range tr.Hops {
			if !h.Responded() {
				rows = append(rows, []string{
					fmt.Sprint(h.TTL), t.Muted.Render("* * *"), "", "", "", "",
				})
				continue
			}
			host := h.Host
			if host == "" {
				host = h.Addr.String()
			} else {
				host = fmt.Sprintf("%s (%s)", host, h.Addr)
			}
			if len(h.Addrs) > 1 {
				host += t.Muted.Render(fmt.Sprintf(" +%d", len(h.Addrs)-1))
			}
			ttl := fmt.Sprint(h.TTL)
			if tr.Method == mtr.MethodTCP {
				ttl = "→"
			}
			rows = append(rows, []string{
				ttl,
				host,
				o.asnCell(h.ASN.Number, h.ASN.String()),
				o.lossCell(h.Loss()),
				formatRTT(h.Best()),
				formatRTT(h.Avg()),
			})
		}
		o.line("%s", o.table(
			[]string{"#", "Host", "Network", "Loss", "Best", "Avg"}, rows))

		for _, n := range tr.Verdict.Notes {
			o.line("  %s %s", t.Muted.Render("·"), t.Muted.Render(n))
		}
		o.line("")
	}
}

// accessTable answers "will this service actually serve me", which is a
// different question from where it thinks I am — a service can geolocate you
// correctly and still refuse the address.
func (o *Output) accessTable(r Report) {
	t := o.Theme
	o.section("Stash checks")

	rows := make([][]string, 0, len(r.Access))
	for _, res := range r.Access {
		region := res.Region
		if region == "" {
			region = "—"
		}
		rows = append(rows, []string{
			res.Check.Name,
			o.accessStateCell(res.State),
			t.Muted.Render(region),
			t.Muted.Render(clip(res.Detail, 52)),
			o.accessTimeCell(res),
		})
	}
	o.line("%s", o.table([]string{"Service", "Status", "Region", "Detail", "Time"}, rows))

	s := access.Summarize(r.Access)
	var parts []string
	add := func(n int, label string, style lg.Style) {
		if n > 0 {
			parts = append(parts, style.Render(fmt.Sprintf("%d %s", n, label)))
		}
	}
	add(s.Available, "available", t.Good)
	add(s.Restricted, "restricted", t.Warn)
	add(s.Blocked, "blocked", t.Bad)
	add(s.Errors, "inconclusive", t.Muted)
	if len(parts) > 0 {
		o.line("  %s", strings.Join(parts, t.Muted.Render(" · ")))
	}
	o.line("")
}

// aiTable reports what the model providers' APIs said. The verdict column
// reads "reachable" rather than "available" on purpose: no credentials are
// sent, so the check proves the endpoint answered and did not refuse the
// region, and claims nothing about whether a real request would succeed.
func (o *Output) aiTable(r Report) {
	t := o.Theme
	o.section("AI endpoints")

	rows := make([][]string, 0, len(r.AI))
	for _, res := range r.AI {
		status := "—"
		if res.Status > 0 {
			status = strconv.Itoa(res.Status)
		}
		rows = append(rows, []string{
			res.Check.Name,
			o.aiStateCell(res.State),
			t.Muted.Render(status),
			t.Muted.Render(clip(res.Detail, 46)),
			o.aiTimeCell(res),
		})
	}
	o.line("%s", o.table([]string{"Endpoint", "Status", "HTTP", "Detail", "Time"}, rows))

	s := ai.Summarize(r.AI)
	var parts []string
	add := func(n int, label string, style lg.Style) {
		if n > 0 {
			parts = append(parts, style.Render(fmt.Sprintf("%d %s", n, label)))
		}
	}
	add(s.Reachable, "reachable", t.Good)
	add(s.Blocked, "blocked", t.Bad)
	add(s.Failed, "inconclusive", t.Muted)
	if len(parts) > 0 {
		o.line("  %s", strings.Join(parts, t.Muted.Render(" · ")))
	}
	o.line("")
}

func (o *Output) aiStateCell(s ai.State) string {
	t := o.Theme
	switch s {
	case ai.StateReachable:
		return t.Good.Render("● reachable")
	case ai.StateBlocked:
		return t.Bad.Render("● blocked")
	default:
		return t.Muted.Render("○ error")
	}
}

func (o *Output) aiTimeCell(res ai.Result) string {
	if res.RTT <= 0 {
		return o.Theme.Muted.Render("—")
	}
	return o.rttCell(res.RTT)
}

func (o *Output) accessStateCell(s access.State) string {
	t := o.Theme
	switch s {
	case access.StateAvailable:
		return t.Good.Render("● available")
	case access.StateRestricted:
		return t.Warn.Render("● restricted")
	case access.StateBlocked:
		return t.Bad.Render("● blocked")
	default:
		return t.Muted.Render("○ error")
	}
}

func (o *Output) accessTimeCell(res access.Result) string {
	if res.Err != nil {
		return o.Theme.Muted.Render("—")
	}
	return o.rttCell(res.RTT)
}

func (o *Output) footer(r Report) {
	t := o.Theme
	if r.Duration > 0 {
		o.line("  %s", t.Muted.Render(fmt.Sprintf("completed in %.1fs", r.Duration.Seconds())))
	}
	o.line("")
}

func (o *Output) section(title string) {
	t := o.Theme
	width := o.Width - 4
	if width < 20 {
		width = 20
	}
	rule := strings.Repeat("─", maxInt(0, width-lg.Width(title)-3))
	o.line("  %s %s", t.Section.Render(title), t.Border.Render(rule))
}

// table builds a bordered table styled for the current theme.
func (o *Output) table(headers []string, rows [][]string) string {
	t := o.Theme
	return table.New().
		Border(lg.RoundedBorder()).
		BorderStyle(t.Border).
		BorderColumn(false).
		BorderRow(false).
		Headers(headers...).
		Rows(rows...).
		StyleFunc(func(row, _ int) lg.Style {
			if row == table.HeaderRow {
				return t.TableHead
			}
			return t.TableCell
		}).
		Render()
}

func (o *Output) verdictCell(c mtr.Class) string {
	t := o.Theme
	switch c {
	case mtr.ClassDirect:
		return t.Good.Render("● " + c.Label())
	case mtr.ClassPeered:
		return t.Info.Render("● " + c.Label())
	case mtr.ClassTransit:
		return t.Warn.Render("● " + c.Label())
	case mtr.ClassDetour:
		return t.Bad.Render("● " + c.Label())
	case mtr.ClassUnreachable:
		return t.Muted.Render("○ " + c.Label())
	default:
		return t.Muted.Render("○ " + c.Label())
	}
}

func (o *Output) rttCell(d time.Duration) string {
	t := o.Theme
	if d == 0 {
		return t.Muted.Render("—")
	}
	ms := d.Seconds() * 1000
	style := t.Good
	switch {
	case ms > 160:
		style = t.Bad
	case ms > 80:
		style = t.Warn
	case ms > 30:
		style = t.Info
	}
	return style.Render(formatRTT(d))
}

func (o *Output) lossCell(loss float64) string {
	t := o.Theme
	if loss <= 0 {
		return t.Muted.Render("0%")
	}
	style := t.Warn
	if loss >= 0.5 {
		style = t.Bad
	}
	return style.Render(fmt.Sprintf("%.0f%%", loss*100))
}

// pathCell summarises the route: the destination network, or the carriers in
// the way when there are any.
func (o *Output) pathCell(tr *mtr.Report) string {
	t := o.Theme
	if len(tr.Verdict.Transits) > 0 {
		names := make([]string, 0, len(tr.Verdict.Transits))
		for _, x := range tr.Verdict.Transits {
			names = append(names, x.Name)
		}
		return t.Warn.Render(strings.Join(names, " → ")) +
			t.Muted.Render(" → "+shortNet(tr))
	}
	if !tr.DestASN.Empty() {
		return t.Muted.Render(tr.DestASN.String())
	}
	return t.Muted.Render(tr.Target.Net)
}

func (o *Output) asnCell(num int, s string) string {
	t := o.Theme
	if s == "" {
		return t.Muted.Render("—")
	}
	if _, isTransit := mtr.TransitName(num); isTransit {
		return t.Warn.Render(s)
	}
	return t.Muted.Render(s)
}

func (o *Output) scoreStyle(score int) lg.Style {
	t := o.Theme
	switch {
	case score >= 85:
		return t.Good
	case score >= 65:
		return t.Info
	case score >= 40:
		return t.Warn
	default:
		return t.Bad
	}
}

func shortNet(tr *mtr.Report) string {
	if !tr.DestASN.Empty() {
		return tr.DestASN.String()
	}
	return tr.Target.Net
}

// Country names are printed as a plain ISO code rather than a flag emoji.
// A flag is a pair of regional-indicator code points, and terminals that lack
// the glyph fall back to drawing two boxed letters, which looks like mojibake
// and confuses width measurement. The code is unambiguous everywhere.

func maskAddr(a netip.Addr, mask bool) string {
	if !mask {
		return a.String()
	}
	if a.Is4() {
		b := a.As4()
		return fmt.Sprintf("%d.%d.x.x", b[0], b[1])
	}
	pfx, err := a.Prefix(32)
	if err != nil {
		return a.String()
	}
	return strings.TrimSuffix(pfx.Addr().String(), "::") + "::/32"
}

// meter draws a proportional bar.
func meter(percent float64, width int) string {
	if percent < 0 {
		percent = 0
	}
	if percent > 100 {
		percent = 100
	}
	filled := int(percent/100*float64(width) + 0.5)
	return strings.Repeat("█", filled) + strings.Repeat("░", width-filled)
}

func formatRTT(d time.Duration) string {
	if d == 0 {
		return "—"
	}
	ms := d.Seconds() * 1000
	switch {
	case ms < 10:
		return fmt.Sprintf("%.2f ms", ms)
	case ms < 100:
		return fmt.Sprintf("%.1f ms", ms)
	default:
		return fmt.Sprintf("%.0f ms", ms)
	}
}

func pad(s string, n int) string {
	if w := lg.Width(s); w < n {
		return s + strings.Repeat(" ", n-w)
	}
	return s
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
