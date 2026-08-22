// Package portal runs the connectivity-check endpoints operating systems use
// to decide whether they are really online: Google's generate_204, Apple's
// hotspot-detect, Microsoft's NCSI and their equivalents.
//
// Each endpoint has a precisely specified answer — a 204 with an empty body, or
// a fixed string — so any deviation is evidence rather than noise. A redirect
// means a captive portal; a 200 where 204 was promised, or a body that does not
// match, means something is rewriting HTTP in flight.
package portal

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"exodus-node/geocheck/internal/netx"
)

// Verdict is the outcome of one endpoint check.
type Verdict int

const (
	// VerdictOK means the endpoint answered exactly as specified.
	VerdictOK Verdict = iota
	// VerdictPortal means the request was redirected, which is how a captive
	// portal announces itself.
	VerdictPortal
	// VerdictAltered means the status or body was not what the vendor
	// documents, so something on the path rewrote the response.
	VerdictAltered
	// VerdictUnreachable means the request failed outright.
	VerdictUnreachable
)

func (v Verdict) String() string {
	switch v {
	case VerdictPortal:
		return "portal"
	case VerdictAltered:
		return "altered"
	case VerdictUnreachable:
		return "unreachable"
	default:
		return "ok"
	}
}

// Label is a short human phrase for the verdict.
func (v Verdict) Label() string {
	switch v {
	case VerdictPortal:
		return "captive portal"
	case VerdictAltered:
		return "altered"
	case VerdictUnreachable:
		return "no answer"
	default:
		return "clean"
	}
}

// Endpoint describes one connectivity-check URL and the answer it must give.
type Endpoint struct {
	ID     string
	Name   string // display label, e.g. "Google (Android)"
	Vendor string
	URL    string

	// WantStatus is the documented status code, normally 204 or 200.
	WantStatus int
	// WantBody is the exact body expected. An empty string means the body must
	// be empty, which is the whole point of the 204 endpoints.
	WantBody string
	// BodyContains relaxes the check to a substring, for endpoints that return
	// a full HTML page whose wording is not contractual.
	BodyContains string

	// Tags group endpoints for selection.
	Tags []string
}

// Result is one endpoint's outcome.
type Result struct {
	Endpoint Endpoint
	Verdict  Verdict

	Status   int
	Redirect string // Location header, when redirected
	Body     string // trimmed, for reporting a mismatch
	RTT      time.Duration
	Err      error
	Detail   string // why the verdict came out this way
}

// Run checks every endpoint concurrently over the given address family.
func Run(ctx context.Context, stack *netx.Stack, f netx.Family, endpoints []Endpoint, concurrency int) []Result {
	if concurrency <= 0 {
		concurrency = 8
	}

	results := make([]Result, len(endpoints))
	sem := make(chan struct{}, concurrency)
	var wg sync.WaitGroup

	for i, ep := range endpoints {
		wg.Add(1)
		go func() {
			defer wg.Done()
			select {
			case sem <- struct{}{}:
			case <-ctx.Done():
				results[i] = Result{Endpoint: ep, Verdict: VerdictUnreachable, Err: ctx.Err()}
				return
			}
			defer func() { <-sem }()
			results[i] = check(ctx, stack, f, ep)
		}()
	}
	wg.Wait()
	return results
}

func check(ctx context.Context, stack *netx.Stack, f netx.Family, ep Endpoint) Result {
	res := Result{Endpoint: ep}

	start := time.Now()
	resp, err := stack.Do(ctx, f, netx.Request{
		URL: ep.URL,
		// The redirect is the finding, so it must not be followed.
		NoRedirect: true,
		// A neutral agent: some portals only intercept browser-looking traffic,
		// and the real OS probes do not pretend to be browsers either.
		UserAgent: "geocheck/1 (connectivity-check)",
	})
	res.RTT = time.Since(start)

	if err != nil {
		res.Verdict = VerdictUnreachable
		res.Err = err
		res.Detail = failureReason(err)
		return res
	}

	res.Status = resp.Status
	body := resp.Text()
	res.Body = truncate(strings.TrimSpace(body), 60)

	if loc := resp.Header.Get("Location"); resp.Status >= 300 && resp.Status < 400 {
		res.Verdict = VerdictPortal
		res.Redirect = loc
		res.Detail = "redirected to " + truncate(loc, 60)
		return res
	}

	// A 5xx is a gateway speaking for the endpoint, not the endpoint speaking
	// differently. The request never arrived, so there is no basis for saying
	// the response was rewritten — but something did answer in its place, and
	// that is worth naming.
	if resp.Status >= 500 {
		res.Verdict = VerdictUnreachable
		res.Detail = fmt.Sprintf(
			"HTTP %d from a gateway; the request never reached the endpoint, so a proxy "+
				"answered on its behalf", resp.Status)
		return res
	}

	if resp.Status != ep.WantStatus {
		res.Verdict = VerdictAltered
		res.Detail = statusMismatch(ep.WantStatus, resp.Status)
		return res
	}

	switch {
	case ep.BodyContains != "":
		if !strings.Contains(body, ep.BodyContains) {
			res.Verdict = VerdictAltered
			res.Detail = "body does not contain " + quote(ep.BodyContains)
			return res
		}
	case ep.WantBody == "":
		if strings.TrimSpace(body) != "" {
			res.Verdict = VerdictAltered
			res.Detail = "expected an empty body, got " + quote(res.Body)
			return res
		}
	default:
		if body != ep.WantBody && strings.TrimSpace(body) != strings.TrimSpace(ep.WantBody) {
			res.Verdict = VerdictAltered
			res.Detail = "expected " + quote(strings.TrimSpace(ep.WantBody)) +
				", got " + quote(res.Body)
			return res
		}
	}

	res.Verdict = VerdictOK
	return res
}

// failureReason turns a transport error into something that says what to do
// about it, rather than echoing a Go error string.
func failureReason(err error) string {
	s := err.Error()
	switch {
	case strings.Contains(s, "certificate") || strings.Contains(s, "x509"):
		return "TLS certificate rejected, which is what interception looks like"
	case strings.Contains(s, "timeout") || strings.Contains(s, "deadline exceeded"):
		return "timed out; the traffic is being dropped rather than refused"
	case strings.Contains(s, "connection reset"):
		return "connection reset, so something on the path tore it down"
	case strings.HasSuffix(s, "EOF"):
		// The TCP handshake succeeded and then the peer hung up without
		// sending a byte. A real web server does not do that; a middlebox
		// that accepts the connection and declines to forward it does.
		return "the connection was accepted then closed with no reply, " +
			"which is how a middlebox drops traffic it will not carry"
	case strings.Contains(s, "connection refused"):
		return "connection refused"
	case strings.Contains(s, "no such host") || strings.Contains(s, "no IPv4 address"):
		return "the name did not resolve"
	default:
		return "request failed"
	}
}

func statusMismatch(want, got int) string {
	d := "expected HTTP " + itoa(want) + ", got " + itoa(got)
	if want == 204 && got == 200 {
		d += "; a 200 where 204 was promised is the classic transparent-proxy signature"
	}
	return d
}

// IsTLS reports whether the endpoint is fetched over HTTPS.
func (r Result) IsTLS() bool { return strings.HasPrefix(r.Endpoint.URL, "https://") }

// Summary counts outcomes across a run. The plain-HTTP and HTTPS tallies are
// kept apart because the interesting cases are asymmetric: HTTP failing while
// HTTPS succeeds means port 80 is filtered, and HTTPS failing while HTTP
// succeeds points at TLS interception.
type Summary struct {
	Total       int
	OK          int
	Portal      int
	Altered     int
	Unreachable int

	PlainTotal int
	PlainOK    int
	PlainFail  int
	TLSTotal   int
	TLSOK      int
}

// Clean reports that no interception was observed. It deliberately requires
// having seen unmodified plain HTTP: if every port-80 probe failed, nothing was
// altered but nothing was verified either, and calling that "clean" would
// overstate what the run established.
func (s Summary) Clean() bool {
	if s.Portal > 0 || s.Altered > 0 || s.OK == 0 {
		return false
	}
	return !s.PlainHTTPBlocked()
}

// PlainHTTPBlocked reports the case where no plain-HTTP endpoint answered while
// HTTPS did, which means port 80 is being dropped rather than rewritten.
func (s Summary) PlainHTTPBlocked() bool {
	return s.PlainTotal > 0 && s.PlainOK == 0 && s.TLSOK > 0
}

// Summarize tallies results.
func Summarize(results []Result) Summary {
	var s Summary
	for _, r := range results {
		s.Total++
		if r.IsTLS() {
			s.TLSTotal++
		} else {
			s.PlainTotal++
		}

		switch r.Verdict {
		case VerdictOK:
			s.OK++
			if r.IsTLS() {
				s.TLSOK++
			} else {
				s.PlainOK++
			}
		case VerdictPortal:
			s.Portal++
		case VerdictAltered:
			s.Altered++
		default:
			s.Unreachable++
			if !r.IsTLS() {
				s.PlainFail++
			}
		}
	}
	return s
}

func truncate(s string, n int) string {
	s = strings.ReplaceAll(s, "\n", " ")
	if len(s) <= n {
		return s
	}
	return s[:n-3] + "..."
}

func quote(s string) string { return "\"" + s + "\"" }

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	var b [8]byte
	i := len(b)
	for n > 0 {
		i--
		b[i] = byte('0' + n%10)
		n /= 10
	}
	return string(b[i:])
}
