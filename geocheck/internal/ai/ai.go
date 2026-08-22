// Package ai answers whether the model providers' APIs will talk to this
// address, which is a different question from whether their consumer websites
// will.
//
// The websites sit behind bot protection: claude.ai and gemini.google.com hand
// out Cloudflare challenges to anything that is not a browser, so a check
// against them is inconclusive on exactly the hosts people most want to test.
// The APIs have no such defence, and they answer in their own vocabulary.
//
// That vocabulary is what these checks read, and it makes an unauthenticated
// probe useful:
//
//   - an authentication error means the request reached the service. It looked
//     at the credentials and found them missing, which it can only do after
//     deciding it is willing to serve the caller at all.
//   - a refusal naming the region means the opposite, and says so outright:
//     OpenAI answers "unsupported_country_region_territory", Anthropic answers
//     "Request not allowed".
//
// So no API key is needed to learn the thing worth knowing. A key would only
// upgrade "reachable" to "usable", which is a question about the account
// rather than about the network.
//
// The verdict is deliberately "reachable" rather than "available": without
// credentials the check proves the endpoint answers and does not refuse this
// region, and stops short of claiming the service would actually serve a
// request.
package ai

import (
	"context"
	"sync"
	"time"

	"exodus-node/geocheck/internal/netx"
)

// State is what the endpoint said about this address.
type State int

const (
	// StateReachable means the API answered and did not refuse the region.
	StateReachable State = iota
	// StateBlocked means it refused, naming the region or country.
	StateBlocked
	// StateError means the probe could not reach a conclusion.
	StateError
)

func (s State) String() string {
	switch s {
	case StateReachable:
		return "reachable"
	case StateBlocked:
		return "blocked"
	default:
		return "error"
	}
}

// Check is one endpoint probe.
type Check struct {
	ID   string
	Name string
	// Vendor groups endpoints that belong to the same company.
	Vendor string
	Run    func(ctx context.Context, env Env) Result
}

// Env is what a check needs to make its request.
type Env struct {
	Stack  *netx.Stack
	Family netx.Family
}

// Result is one endpoint's answer.
type Result struct {
	Check  Check
	State  State
	Status int
	// Detail is the reason, in the service's own words where it gave one.
	Detail string
	RTT    time.Duration
	Err    error
}

// Run probes every endpoint concurrently.
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

// Summary counts the outcomes for the footer line.
type Summary struct {
	Total     int
	Reachable int
	Blocked   int
	Failed    int
}

// Summarize tallies results.
func Summarize(results []Result) Summary {
	s := Summary{Total: len(results)}
	for _, r := range results {
		switch r.State {
		case StateReachable:
			s.Reachable++
		case StateBlocked:
			s.Blocked++
		default:
			s.Failed++
		}
	}
	return s
}
