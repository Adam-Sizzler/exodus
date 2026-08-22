// Package access answers a narrower question than the geolocation checks: not
// "which country does this service think I am in", but "will it actually serve
// me". The two differ more often than expected — a service can geolocate you
// correctly and still refuse, because the exit address belongs to a hosting
// provider, or because the catalogue in that region is partial.
//
// Each check probes the endpoint the service's own clients use and reads the
// specific marker it returns, rather than inferring availability from a status
// code alone.
//
// The probes are ported from the Stash collapsed-tile panel scripts, which is
// why the report groups them under "Stash checks". They are kept apart from the
// geolocation checks on purpose: those ask which country a service assigns you,
// these ask whether it will serve you at all, and the two disagree exactly in
// the cases worth knowing about.
package access

import (
	"context"
	"sync"
	"time"

	"exodus-node/geocheck/internal/netx"
)

// State is how usable a service is from here.
type State int

const (
	// StateAvailable means the service serves this address normally.
	StateAvailable State = iota
	// StateRestricted means it serves a reduced offering: reachable, but not
	// the full catalogue or feature set.
	StateRestricted
	// StateBlocked means the service refuses this address or region.
	StateBlocked
	// StateError means the check could not reach a conclusion.
	StateError
)

func (s State) String() string {
	switch s {
	case StateAvailable:
		return "available"
	case StateRestricted:
		return "restricted"
	case StateBlocked:
		return "blocked"
	default:
		return "error"
	}
}

// Check is one service-availability probe.
type Check struct {
	ID   string
	Name string
	Run  func(ctx context.Context, env Env) Result
}

// Env is what a check needs to run.
type Env struct {
	Stack  *netx.Stack
	Family netx.Family
}

// Result is one check's outcome.
type Result struct {
	Check Check
	State State
	// Detail is the specific finding, e.g. "Originals only" or
	// "Disallowed ISP" — the part that says *why*.
	Detail string
	// Region is the country the service reports serving, when it says so.
	Region string
	RTT    time.Duration
	Err    error
}

// Run executes every check concurrently.
func Run(ctx context.Context, stack *netx.Stack, f netx.Family, checks []Check, concurrency int) []Result {
	if concurrency <= 0 {
		concurrency = 6
	}

	results := make([]Result, len(checks))
	sem := make(chan struct{}, concurrency)
	var wg sync.WaitGroup

	for i, c := range checks {
		wg.Add(1)
		go func() {
			defer wg.Done()
			select {
			case sem <- struct{}{}:
			case <-ctx.Done():
				results[i] = Result{Check: c, State: StateError, Err: ctx.Err()}
				return
			}
			defer func() { <-sem }()

			start := time.Now()
			res := c.Run(ctx, Env{Stack: stack, Family: f})
			res.Check = c
			if res.RTT == 0 {
				res.RTT = time.Since(start)
			}
			results[i] = res
		}()
	}
	wg.Wait()
	return results
}

// Summary counts outcomes.
type Summary struct {
	Total      int
	Available  int
	Restricted int
	Blocked    int
	Errors     int
}

// Summarize tallies results.
func Summarize(results []Result) Summary {
	var s Summary
	for _, r := range results {
		s.Total++
		switch r.State {
		case StateAvailable:
			s.Available++
		case StateRestricted:
			s.Restricted++
		case StateBlocked:
			s.Blocked++
		default:
			s.Errors++
		}
	}
	return s
}
