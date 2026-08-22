// Package geo determines what the internet thinks your location is, by asking
// GeoIP APIs and by reading the region that popular consumer services serve you.
package geo

import (
	"context"
	"net/netip"
	"regexp"
	"sort"
	"strings"
	"sync"

	"exodus-node/geocheck/internal/countries"
	"exodus-node/geocheck/internal/netx"
	"golang.org/x/sync/errgroup"
)

// Group buckets checks for display.
type Group string

const (
	// GroupServices covers consumer products whose served region matters.
	GroupServices Group = "services"
	// GroupGeoIP covers databases and APIs that map an IP to a country.
	GroupGeoIP Group = "geoip"
	// GroupCDN covers which CDN edge node actually serves you.
	GroupCDN Group = "cdn"
)

// Kind describes how a result should be interpreted and rendered.
type Kind int

const (
	// KindCountry yields an ISO 3166-1 alpha-2 code.
	KindCountry Kind = iota
	// KindAvailability yields yes/no: is the service usable from here.
	KindAvailability
	// KindBlocked yields yes/no where "yes" is the bad outcome (captcha, ban).
	KindBlocked
)

// Env is the shared context a check runs against.
type Env struct {
	Stack *netx.Stack
	IPv4  netip.Addr
	IPv6  netip.Addr

	mu     sync.RWMutex
	shared map[string]string // cross-check fallbacks, e.g. YouTube reusing Google
}

// PublicIP returns the detected external address for a family.
func (e *Env) PublicIP(f netx.Family) netip.Addr {
	if f == netx.V6 {
		return e.IPv6
	}
	return e.IPv4
}

func (e *Env) share(key string, f netx.Family, val string) {
	if val == "" {
		return
	}
	e.mu.Lock()
	defer e.mu.Unlock()
	if e.shared == nil {
		e.shared = map[string]string{}
	}
	e.shared[key+"/"+f.String()] = val
}

func (e *Env) lookupShared(key string, f netx.Family) string {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.shared[key+"/"+f.String()]
}

// RunFunc performs one probe and returns a raw value ("" when unavailable).
type RunFunc func(ctx context.Context, env *Env, f netx.Family) (string, error)

// Check is a single named probe.
type Check struct {
	ID    string
	Name  string
	Group Group
	Kind  Kind
	Run   RunFunc

	// DependsOn names a check that must complete first (fallback chains).
	DependsOn string
}

// Result holds one check's outcome per address family.
type Result struct {
	Check Check
	V4    Outcome
	V6    Outcome
}

// Outcome is the value a check produced for one address family.
type Outcome struct {
	Value   string // country code, or "yes"/"no"
	Err     error
	Skipped bool
}

// Country returns the uppercased ISO code, or "" when there is none.
func (o Outcome) Country() string {
	if o.Skipped || o.Err != nil {
		return ""
	}
	return o.Value
}

// Runner executes a set of checks concurrently.
type Runner struct {
	Env         *Env
	Families    []netx.Family
	Concurrency int
	// OnProgress, when set, is called as each check finishes.
	OnProgress func(done, total int, name string)
}

// Run executes checks and returns results in the order given.
func (r *Runner) Run(ctx context.Context, checks []Check) []Result {
	if r.Concurrency <= 0 {
		r.Concurrency = 12
	}

	results := make([]Result, len(checks))
	index := make(map[string]int, len(checks))
	for i, c := range checks {
		results[i] = Result{Check: c}
		index[c.ID] = i
	}

	// Checks that others depend on run first, in their own wave, so the
	// fallback data is present by the time the dependents execute.
	var first, second []int
	needed := map[string]bool{}
	for _, c := range checks {
		if c.DependsOn != "" {
			needed[c.DependsOn] = true
		}
	}
	for i, c := range checks {
		if needed[c.ID] {
			first = append(first, i)
		} else {
			second = append(second, i)
		}
	}

	var done int
	var progressMu sync.Mutex
	report := func(name string) {
		if r.OnProgress == nil {
			return
		}
		progressMu.Lock()
		done++
		n := done
		progressMu.Unlock()
		r.OnProgress(n, len(checks), name)
	}

	runWave := func(idxs []int) {
		g, gctx := errgroup.WithContext(ctx)
		g.SetLimit(r.Concurrency)
		for _, i := range idxs {
			g.Go(func() error {
				c := checks[i]
				var res Result
				res.Check = c
				for _, f := range r.Families {
					out := r.runOne(gctx, c, f)
					if f == netx.V6 {
						res.V6 = out
					} else {
						res.V4 = out
					}
				}
				results[i] = res
				report(c.Name)
				return nil
			})
		}
		_ = g.Wait()
	}

	runWave(first)
	runWave(second)
	return results
}

func (r *Runner) runOne(ctx context.Context, c Check, f netx.Family) Outcome {
	if !r.Env.PublicIP(f).IsValid() {
		return Outcome{Skipped: true}
	}
	val, err := c.Run(ctx, r.Env, f)
	if err != nil {
		return Outcome{Err: err}
	}
	val = normalize(val, c.Kind)
	if val != "" {
		r.Env.share(c.ID, f, val)
	}
	return Outcome{Value: val}
}

var nonAlpha = regexp.MustCompile(`[^A-Za-z]`)

// normalize cleans a raw provider value into a canonical form.
func normalize(v string, kind Kind) string {
	v = strings.TrimSpace(v)
	if v == "" || strings.EqualFold(v, "null") || strings.EqualFold(v, "undefined") {
		return ""
	}
	switch kind {
	case KindAvailability, KindBlocked:
		return strings.ToLower(v)
	default:
		v = nonAlpha.ReplaceAllString(v, "")
		if len(v) != 2 {
			return ""
		}
		return strings.ToUpper(v)
	}
}

// Consensus summarises how many providers agreed on each country code.
type Consensus struct {
	Code    string
	Name    string
	Count   int
	Total   int
	Percent float64
}

// Summarize tallies country answers across results for one family.
func Summarize(results []Result, f netx.Family) []Consensus {
	counts := map[string]int{}
	total := 0
	for _, r := range results {
		if r.Check.Kind != KindCountry || r.Check.Group == GroupCDN {
			continue
		}
		out := r.V4
		if f == netx.V6 {
			out = r.V6
		}
		code := out.Country()
		if code == "" {
			continue
		}
		counts[code]++
		total++
	}
	if total == 0 {
		return nil
	}

	out := make([]Consensus, 0, len(counts))
	for code, n := range counts {
		out = append(out, Consensus{
			Code:    code,
			Name:    countries.Name(code),
			Count:   n,
			Total:   total,
			Percent: float64(n) / float64(total) * 100,
		})
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Count != out[j].Count {
			return out[i].Count > out[j].Count
		}
		return out[i].Code < out[j].Code
	})
	return out
}
