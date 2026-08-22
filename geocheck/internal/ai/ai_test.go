package ai

import "testing"

func find(t *testing.T, id string) endpoint {
	t.Helper()
	for _, e := range catalog {
		if e.id == id {
			return e
		}
	}
	t.Fatalf("no endpoint %q in the catalogue", id)
	return endpoint{}
}

// TestAuthErrorProvesReachability is the idea the whole package rests on: an
// API that asks for credentials has already agreed to serve the caller's
// network, so 401 is a positive result rather than a failure.
func TestAuthErrorProvesReachability(t *testing.T) {
	e := find(t, "openai")
	got := e.classify(401, `{"error":{"message":"Missing bearer authentication in header","type":"invalid_request_error"}}`)
	if got.State != StateReachable {
		t.Errorf("state = %v, want reachable", got.State)
	}
}

// TestRegionRefusalBeatsTheStatus covers the case that makes a status-only
// reading wrong: the refusal and a missing key can share a status code, and
// only the body separates them.
func TestRegionRefusalBeatsTheStatus(t *testing.T) {
	cases := []struct {
		id     string
		status int
		body   string
	}{
		{"openai", 403, `{"error":{"code":"unsupported_country_region_territory","message":"Country, region, or territory not supported"}}`},
		{"anthropic", 403, `{"error":{"type":"forbidden","message":"Request not allowed"}}`},
	}
	for _, c := range cases {
		e := find(t, c.id)
		if got := e.classify(c.status, c.body); got.State != StateBlocked {
			t.Errorf("%s: state = %v, want blocked", c.id, got.State)
		}
	}
}

// TestChallengeIsNotAnAnswer keeps an interstitial from being read as either
// verdict — the same discipline the website checks follow.
func TestChallengeIsNotAnAnswer(t *testing.T) {
	e := find(t, "anthropic")
	body := `<html><head><title>Just a moment...</title></head><body>` +
		`<script>window._cf_chl_opt={};</script></body></html>`
	got := e.classify(403, body)
	if got.State != StateError {
		t.Errorf("state = %v, want error", got.State)
	}
}

// TestGeminiClaimsOnlyReachability records a measured limitation: Google
// answers 403 "unregistered caller" whether or not the region is served, so
// this endpoint must never produce a blocked verdict.
func TestGeminiClaimsOnlyReachability(t *testing.T) {
	e := find(t, "gemini")
	if len(e.regionMarkers) != 0 {
		t.Error("the Gemini endpoint has no region signal; it must not carry markers")
	}
	body := `{"error":{"code":403,"message":"Method doesn't allow unregistered callers"}}`
	if got := e.classify(403, body); got.State != StateReachable {
		t.Errorf("state = %v, want reachable", got.State)
	}
}

// TestUnexpectedStatusIsInconclusive stops a strange answer from being turned
// into a verdict.
func TestUnexpectedStatusIsInconclusive(t *testing.T) {
	e := find(t, "deepseek")
	for _, status := range []int{500, 502, 418} {
		if got := e.classify(status, "server exploded"); got.State != StateError {
			t.Errorf("status %d gave %v, want error", status, got.State)
		}
	}
}

func TestSummarize(t *testing.T) {
	got := Summarize([]Result{
		{State: StateReachable}, {State: StateReachable},
		{State: StateBlocked}, {State: StateError},
	})
	want := Summary{Total: 4, Reachable: 2, Blocked: 1, Failed: 1}
	if got != want {
		t.Errorf("got %+v, want %+v", got, want)
	}
}

// TestEveryEndpointIsHTTPS keeps an unencrypted probe from leaking which
// services are being asked about.
func TestEveryEndpointIsHTTPS(t *testing.T) {
	for _, e := range catalog {
		if len(e.url) < 8 || e.url[:8] != "https://" {
			t.Errorf("%s is not https: %s", e.id, e.url)
		}
		if len(e.authStatus) == 0 {
			t.Errorf("%s declares no status that proves reachability", e.id)
		}
	}
}
