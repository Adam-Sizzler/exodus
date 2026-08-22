package render

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io"
	"time"

	"exodus-node/geocheck/internal/access"

	"exodus-node/geocheck/internal/countries"
	"exodus-node/geocheck/internal/detect"
	"exodus-node/geocheck/internal/geo"
	"exodus-node/geocheck/internal/mtr"
	"exodus-node/geocheck/internal/netx"
	"exodus-node/geocheck/internal/portal"
	"exodus-node/geocheck/internal/reputation"
)

// schemaVersion is bumped whenever the JSON shape changes incompatibly.
const schemaVersion = 1

type jsonReport struct {
	Schema     int             `json:"schema"`
	Tool       string          `json:"tool"`
	Timestamp  string          `json:"timestamp"`
	DurationMS int64           `json:"duration_ms"`
	Identity   jsonIdentity    `json:"identity"`
	Transport  jsonTransport   `json:"transport"`
	Findings   []jsonFinding   `json:"findings"`
	Reputation *jsonReputation `json:"reputation,omitempty"`
	Consensus  map[string]any  `json:"consensus,omitempty"`
	Geo        map[string]any  `json:"geo,omitempty"`
	Portal     *jsonPortal     `json:"connectivity_checks,omitempty"`
	Access     []jsonAccess    `json:"stash_checks,omitempty"`
	AI         []jsonAI        `json:"ai_endpoints,omitempty"`
	Path       *jsonPathReport `json:"connectivity,omitempty"`
	Image      *jsonImage      `json:"image,omitempty"`
}

// jsonImage carries the rendered report as a picture. MediaType and Encoding
// are spelled out rather than implied so a consumer can assemble a data: URI
// straight from the fields:
//
//	data:<media_type>;<encoding>,<data>
type jsonImage struct {
	Format    string `json:"format"`
	MediaType string `json:"media_type"`
	Encoding  string `json:"encoding"`
	Data      string `json:"data"`
}

type jsonIdentity struct {
	IPv4    string `json:"ipv4,omitempty"`
	IPv6    string `json:"ipv6,omitempty"`
	ASN     int    `json:"asn,omitempty"`
	ASName  string `json:"as_name,omitempty"`
	Org     string `json:"org,omitempty"`
	Country string `json:"as_country,omitempty"`
}

type jsonTransport struct {
	Interface string `json:"interface,omitempty"`
	Proxy     string `json:"proxy,omitempty"`
	Resolver  string `json:"resolver,omitempty"`
}

type jsonFinding struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Severity string `json:"severity"`
	Detail   string `json:"detail"`
}

type jsonReputation struct {
	Type         string   `json:"type,omitempty"`
	Residential  bool     `json:"residential"`
	Risk         int      `json:"risk"`
	Confidence   int      `json:"confidence,omitempty"`
	Flags        []string `json:"flags,omitempty"`
	Proxy        bool     `json:"proxy"`
	VPN          bool     `json:"vpn"`
	Tor          bool     `json:"tor"`
	Hosting      bool     `json:"hosting"`
	Scraper      bool     `json:"scraper"`
	Compromised  bool     `json:"compromised"`
	Anonymous    bool     `json:"anonymous"`
	Operator     string   `json:"operator,omitempty"`
	OperatorURL  string   `json:"operator_url,omitempty"`
	Anonymity    string   `json:"operator_anonymity,omitempty"`
	Provider     string   `json:"provider,omitempty"`
	Range        string   `json:"range,omitempty"`
	Hostname     string   `json:"hostname,omitempty"`
	City         string   `json:"city,omitempty"`
	Region       string   `json:"region,omitempty"`
	Country      string   `json:"country,omitempty"`
	CountryCode  string   `json:"country_code,omitempty"`
	DevicesAddr  int      `json:"devices_on_address,omitempty"`
	DevicesSubnt int      `json:"devices_in_subnet,omitempty"`
	FirstSeen    string   `json:"first_seen,omitempty"`
	LastSeen     string   `json:"last_seen,omitempty"`
	Error        string   `json:"error,omitempty"`
}

type jsonConsensus struct {
	Code    string  `json:"code"`
	Country string  `json:"country"`
	Count   int     `json:"count"`
	Total   int     `json:"total"`
	Percent float64 `json:"percent"`
}

type jsonCheck struct {
	ID   string          `json:"id"`
	Name string          `json:"name"`
	Kind string          `json:"kind"`
	IPv4 *jsonCheckValue `json:"ipv4,omitempty"`
	IPv6 *jsonCheckValue `json:"ipv6,omitempty"`
}

type jsonCheckValue struct {
	Value   string `json:"value,omitempty"`
	Country string `json:"country,omitempty"`
	Error   string `json:"error,omitempty"`
}

type jsonPortal struct {
	Clean        bool             `json:"clean"`
	PlainBlocked bool             `json:"plain_http_blocked"`
	OK           int              `json:"ok"`
	Portal       int              `json:"captive_portal"`
	Altered      int              `json:"altered"`
	Unreachable  int              `json:"unreachable"`
	Endpoints    []jsonPortalItem `json:"endpoints"`
}

type jsonPortalItem struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	Vendor   string  `json:"vendor"`
	URL      string  `json:"url"`
	Verdict  string  `json:"verdict"`
	Status   int     `json:"status,omitempty"`
	Expected int     `json:"expected_status"`
	Redirect string  `json:"redirect,omitempty"`
	Body     string  `json:"body,omitempty"`
	RTTMS    float64 `json:"rtt_ms,omitempty"`
	Detail   string  `json:"detail,omitempty"`
	Error    string  `json:"error,omitempty"`
}

type jsonAI struct {
	ID     string  `json:"id"`
	Name   string  `json:"name"`
	Vendor string  `json:"vendor,omitempty"`
	State  string  `json:"state"`
	Status int     `json:"http_status,omitempty"`
	Detail string  `json:"detail,omitempty"`
	RTTMS  float64 `json:"rtt_ms,omitempty"`
	Error  string  `json:"error,omitempty"`
}

type jsonAccess struct {
	ID     string  `json:"id"`
	Name   string  `json:"name"`
	State  string  `json:"state"`
	Region string  `json:"region,omitempty"`
	Detail string  `json:"detail,omitempty"`
	RTTMS  float64 `json:"rtt_ms,omitempty"`
	Error  string  `json:"error,omitempty"`
}

type jsonPathReport struct {
	Available bool         `json:"icmp_available"`
	Raw       bool         `json:"privileged"`
	Hint      string       `json:"hint,omitempty"`
	Score     int          `json:"score"`
	FloorMS   float64      `json:"latency_floor_ms"`
	Breakdown jsonCounts   `json:"breakdown"`
	Targets   []jsonTarget `json:"targets"`
}

type jsonCounts struct {
	Direct      int `json:"direct"`
	Peered      int `json:"peered"`
	Transit     int `json:"transit"`
	Detour      int `json:"detour"`
	Intercepted int `json:"intercepted"`
	Failed      int `json:"failed"`
}

type jsonTarget struct {
	ID       string     `json:"id"`
	Name     string     `json:"name"`
	Host     string     `json:"host"`
	Resolved string     `json:"resolved,omitempty"`
	Method   string     `json:"method"`
	Anycast  bool       `json:"anycast"`
	DestASN  int        `json:"dest_asn,omitempty"`
	DestName string     `json:"dest_as_name,omitempty"`
	Verdict  string     `json:"verdict"`
	Score    int        `json:"score"`
	RTTMS    float64    `json:"rtt_ms"`
	ExcessMS float64    `json:"excess_ms"`
	JitterMS float64    `json:"jitter_ms"`
	Loss     float64    `json:"loss"`
	Hops     []jsonHop  `json:"hops,omitempty"`
	Transits []jsonTran `json:"transits,omitempty"`
	Notes    []string   `json:"notes,omitempty"`
	Error    string     `json:"error,omitempty"`
}

type jsonTran struct {
	ASN  int    `json:"asn"`
	Name string `json:"name"`
	TTL  int    `json:"ttl"`
}

type jsonHop struct {
	TTL     int      `json:"ttl"`
	Addr    string   `json:"addr,omitempty"`
	Host    string   `json:"host,omitempty"`
	ASN     int      `json:"asn,omitempty"`
	ASName  string   `json:"as_name,omitempty"`
	Sent    int      `json:"sent"`
	Recv    int      `json:"recv"`
	Loss    float64  `json:"loss"`
	BestMS  float64  `json:"best_ms"`
	AvgMS   float64  `json:"avg_ms"`
	WorstMS float64  `json:"worst_ms"`
	Addrs   []string `json:"addrs,omitempty"`
}

// JSON writes the machine-readable form of a report.
func JSON(w io.Writer, r Report, findings []detect.Finding, now time.Time) error {
	out := jsonReport{
		Schema:     schemaVersion,
		Tool:       r.Version,
		Timestamp:  now.UTC().Format(time.RFC3339),
		DurationMS: r.Duration.Milliseconds(),
		Identity: jsonIdentity{
			ASN:     r.Identity.ASN.Number,
			ASName:  r.Identity.ASN.Name,
			Org:     r.Identity.Org,
			Country: r.Identity.ASN.Country,
		},
		Transport: jsonTransport{
			Interface: r.Interface,
			Proxy:     r.Proxy,
			Resolver:  r.Resolver,
		},
		Findings: make([]jsonFinding, 0, len(findings)),
	}
	if r.Identity.IPv4.IsValid() {
		out.Identity.IPv4 = maskAddr(r.Identity.IPv4, r.MaskIP)
	}
	if r.Identity.IPv6.IsValid() {
		out.Identity.IPv6 = maskAddr(r.Identity.IPv6, r.MaskIP)
	}

	for _, f := range findings {
		out.Findings = append(out.Findings, jsonFinding{
			ID: f.ID, Title: f.Title, Severity: f.Severity.String(), Detail: f.Detail,
		})
	}

	out.Reputation = jsonReputationReport(r.Reputation, r.ReputationErr)

	if len(r.Geo) > 0 {
		out.Consensus = map[string]any{}
		out.Geo = map[string]any{}
		for _, f := range r.Families {
			rows := geo.Summarize(r.Geo, f)
			cons := make([]jsonConsensus, 0, len(rows))
			for _, c := range rows {
				cons = append(cons, jsonConsensus{
					Code: c.Code, Country: c.Name, Count: c.Count,
					Total: c.Total, Percent: round2(c.Percent),
				})
			}
			out.Consensus[familyKey(f)] = cons
		}
		for _, g := range []geo.Group{geo.GroupServices, geo.GroupGeoIP, geo.GroupCDN} {
			out.Geo[string(g)] = jsonChecks(r, g)
		}
	}

	if len(r.Portal) > 0 {
		out.Portal = jsonPortalReport(r.Portal)
	}

	if len(r.Access) > 0 {
		out.Access = jsonAccessReport(r.Access)
	}

	if len(r.AI) > 0 {
		out.AI = make([]jsonAI, 0, len(r.AI))
		for _, x := range r.AI {
			item := jsonAI{
				ID: x.Check.ID, Name: x.Check.Name, Vendor: x.Check.Vendor,
				State: x.State.String(), Status: x.Status, Detail: x.Detail,
				RTTMS: round2(ms(x.RTT)),
			}
			if x.Err != nil {
				item.Error = x.Err.Error()
			}
			out.AI = append(out.AI, item)
		}
	}

	if len(r.Trace) > 0 {
		out.Path = jsonPath(r)
	}

	if r.EmbedSVG {
		var svg bytes.Buffer
		if err := SVG(&svg, r, findings); err != nil {
			return err
		}
		out.Image = &jsonImage{
			Format:    "svg",
			MediaType: "image/svg+xml",
			Encoding:  "base64",
			Data:      base64.StdEncoding.EncodeToString(svg.Bytes()),
		}
	}

	return json.NewEncoder(w).Encode(out)
}

func jsonReputationReport(info *reputation.Info, err error) *jsonReputation {
	if err != nil {
		return &jsonReputation{Error: err.Error()}
	}
	if info == nil {
		return nil
	}
	return &jsonReputation{
		Type: info.Type, Residential: info.Residential(),
		Risk: info.Risk, Confidence: info.Confidence, Flags: info.Flags(),
		Proxy: info.Proxy, VPN: info.VPN, Tor: info.Tor, Hosting: info.Hosting,
		Scraper: info.Scraper, Compromised: info.Compromised, Anonymous: info.Anonymous,
		Operator: info.Operator, OperatorURL: info.OperatorURL, Anonymity: info.Anonymity,
		Provider: info.Provider, Range: info.Range, Hostname: info.Hostname,
		City: info.City, Region: info.Region, Country: info.Country, CountryCode: info.Code,
		DevicesAddr: info.DevicesAddr, DevicesSubnt: info.DevicesSubnt,
		FirstSeen: info.FirstSeen, LastSeen: info.LastSeen,
	}
}

func jsonChecks(r Report, group geo.Group) []jsonCheck {
	var out []jsonCheck
	for _, res := range r.Geo {
		if res.Check.Group != group {
			continue
		}
		c := jsonCheck{ID: res.Check.ID, Name: res.Check.Name, Kind: kindName(res.Check.Kind)}
		for _, f := range r.Families {
			o := res.V4
			if f == netx.V6 {
				o = res.V6
			}
			v := jsonValue(o, res.Check.Kind)
			if v == nil {
				continue
			}
			if f == netx.V6 {
				c.IPv6 = v
			} else {
				c.IPv4 = v
			}
		}
		out = append(out, c)
	}
	return out
}

func jsonValue(o geo.Outcome, kind geo.Kind) *jsonCheckValue {
	if o.Skipped {
		return nil
	}
	v := &jsonCheckValue{Value: o.Value}
	if o.Err != nil {
		v.Error = o.Err.Error()
	}
	if kind == geo.KindCountry && o.Value != "" {
		v.Country = countries.Name(o.Value)
	}
	return v
}

func jsonPortalReport(results []portal.Result) *jsonPortal {
	s := portal.Summarize(results)
	out := &jsonPortal{
		Clean:        s.Clean(),
		PlainBlocked: s.PlainHTTPBlocked(),
		OK:           s.OK,
		Portal:       s.Portal,
		Altered:      s.Altered,
		Unreachable:  s.Unreachable,
		Endpoints:    make([]jsonPortalItem, 0, len(results)),
	}
	for _, r := range results {
		item := jsonPortalItem{
			ID: r.Endpoint.ID, Name: r.Endpoint.Name, Vendor: r.Endpoint.Vendor,
			URL: r.Endpoint.URL, Verdict: r.Verdict.String(),
			Status: r.Status, Expected: r.Endpoint.WantStatus,
			Redirect: r.Redirect, Body: r.Body, Detail: r.Detail,
		}
		if r.Err == nil {
			item.RTTMS = round2(ms(r.RTT))
		} else {
			item.Error = r.Err.Error()
		}
		out.Endpoints = append(out.Endpoints, item)
	}
	return out
}

func jsonAccessReport(results []access.Result) []jsonAccess {
	out := make([]jsonAccess, 0, len(results))
	for _, r := range results {
		item := jsonAccess{
			ID: r.Check.ID, Name: r.Check.Name,
			State: r.State.String(), Region: r.Region, Detail: r.Detail,
		}
		if r.Err != nil {
			item.Error = r.Err.Error()
		} else {
			item.RTTMS = round2(ms(r.RTT))
		}
		out = append(out, item)
	}
	return out
}

func jsonPath(r Report) *jsonPathReport {
	s := mtr.Summarize(r.Trace)
	out := &jsonPathReport{
		Available: r.TraceCap.ICMP,
		Raw:       r.TraceCap.Raw,
		Hint:      r.TraceCap.Hint,
		Score:     s.Score,
		FloorMS:   round2(ms(s.Floor)),
		Breakdown: jsonCounts{
			Direct: s.Direct, Peered: s.Peered, Transit: s.Transit,
			Detour: s.Detour, Intercepted: s.Intercepted, Failed: s.Failed,
		},
		Targets: make([]jsonTarget, 0, len(r.Trace)),
	}

	for _, tr := range r.Trace {
		if tr == nil {
			continue
		}
		t := jsonTarget{
			ID: tr.Target.ID, Name: tr.Target.Name, Host: tr.Target.Host,
			Method:   string(tr.Method),
			Anycast:  tr.Target.Anycast,
			DestASN:  tr.DestASN.Number,
			DestName: tr.DestASN.Name,
			Verdict:  tr.Verdict.Class.String(),
			Score:    tr.Verdict.Score,
			RTTMS:    round2(ms(tr.Verdict.RTT)),
			ExcessMS: round2(ms(tr.Verdict.Excess)),
			JitterMS: round2(ms(tr.Verdict.Jitter)),
			Loss:     round2(tr.Verdict.Loss),
			Notes:    tr.Verdict.Notes,
		}
		if tr.Resolved.IsValid() {
			t.Resolved = tr.Resolved.String()
		}
		if tr.Err != nil {
			t.Error = tr.Err.Error()
		}
		for _, x := range tr.Verdict.Transits {
			t.Transits = append(t.Transits, jsonTran{ASN: x.ASN, Name: x.Name, TTL: x.TTL})
		}
		for _, h := range tr.Hops {
			jh := jsonHop{
				TTL: h.TTL, Host: h.Host, Sent: h.Sent, Recv: h.Recv,
				ASN: h.ASN.Number, ASName: h.ASN.Name,
				Loss:    round2(h.Loss()),
				BestMS:  round2(ms(h.Best())),
				AvgMS:   round2(ms(h.Avg())),
				WorstMS: round2(ms(h.Worst())),
			}
			if h.Addr.IsValid() {
				jh.Addr = h.Addr.String()
			}
			for _, a := range h.Addrs {
				jh.Addrs = append(jh.Addrs, a.String())
			}
			t.Hops = append(t.Hops, jh)
		}
		out.Targets = append(out.Targets, t)
	}
	return out
}

func kindName(k geo.Kind) string {
	switch k {
	case geo.KindAvailability:
		return "availability"
	case geo.KindBlocked:
		return "blocked"
	default:
		return "country"
	}
}

func familyKey(f netx.Family) string {
	if f == netx.V6 {
		return "ipv6"
	}
	return "ipv4"
}

func ms(d time.Duration) float64 { return d.Seconds() * 1000 }

func round2(f float64) float64 {
	return float64(int64(f*100+0.5)) / 100
}
