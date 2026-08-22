package access

import "testing"

// challengePage reproduces the shape of Cloudflare's interstitial, including
// the detail that motivates the ordering in classifyClaude: the form action
// echoes the path that was requested. Asking for the region-refusal page
// therefore yields a Cloudflare page containing "app-unavailable-in-region"
// even though the origin never answered.
func challengePage(path string) string {
	return `<!DOCTYPE html><html><head><title>Just a moment...</title></head>
<body><div class="main-wrapper"><form id="challenge-form" action="` + path +
		`?__cf_chl_tk=PpajGHGRUMFSgUWDnOn" method="POST">
</form><script>window._cf_chl_opt={cvId: '3', cType: 'managed'};</script>
<div>Enable JavaScript and cookies to continue</div></div></body></html>`
}

const regionPage = `<!DOCTYPE html><html><head><title>Claude</title></head>
<body><h1>Unfortunately, Claude isn't available here.</h1>
<p>Claude is not yet available in your region.</p></body></html>`

const claudeTestRegion = "XA"

func TestClassifyClaude(t *testing.T) {
	cases := []struct {
		name   string
		status int
		body   string
		want   State
		detail string
	}{
		{
			name:   "a challenge is not an answer about availability",
			status: 403,
			body:   challengePage("/"),
			want:   StateError,
		},
		{
			// The regression this whole split exists for. Before the challenge
			// was tested first, the path echoed inside Cloudflare's own markup
			// matched the region marker and produced a confident "blocked" for
			// a request the service never saw.
			name:   "a challenge echoing the region path is still not a region block",
			status: 403,
			body:   challengePage("/app-unavailable-in-region"),
			want:   StateError,
		},
		{
			name:   "the real refusal page is a region block",
			status: 200,
			body:   regionPage,
			want:   StateBlocked,
		},
		{
			name:   "a 403 refusal page is still a region block",
			status: 403,
			body:   regionPage,
			want:   StateBlocked,
		},
		{
			name:   "an ordinary page is available",
			status: 200,
			body:   `<!DOCTYPE html><html><body><div id="root"></div></body></html>`,
			want:   StateAvailable,
		},
		{
			name:   "a bare 403 proves nothing either way",
			status: 403,
			body:   `Forbidden`,
			want:   StateError,
		},
		{
			name:   "a server error is not a block",
			status: 502,
			body:   `Bad Gateway`,
			want:   StateError,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := classifyClaude(c.status, c.body, claudeTestRegion)
			if got.State != c.want {
				t.Errorf("state = %v, want %v (detail %q)", got.State, c.want, got.Detail)
			}
			if got.Region != claudeTestRegion {
				t.Errorf("region = %q, want %s", got.Region, claudeTestRegion)
			}
		})
	}
}

// TestChallengeNeedsARefusalStatus keeps the detector from firing on a page
// that merely happens to contain the wording — a blog post about Cloudflare,
// say — when the service answered normally.
func TestChallengeNeedsARefusalStatus(t *testing.T) {
	body := challengePage("/")
	if isChallenge(200, body) {
		t.Error("a 200 response should never be read as a challenge")
	}
	for _, status := range []int{403, 429, 503} {
		if !isChallenge(status, body) {
			t.Errorf("status %d with challenge markup was not detected", status)
		}
	}
}

// TestClassifyClaudeNeverGuessesBlocked is the property that matters: the only
// route to "blocked" is a page that says so in words.
func TestClassifyClaudeNeverGuessesBlocked(t *testing.T) {
	inconclusive := []struct {
		status int
		body   string
	}{
		{403, challengePage("/")},
		{403, challengePage("/app-unavailable-in-region")},
		{403, ""},
		{429, challengePage("/")},
		{503, "origin unreachable"},
		{500, "internal error"},
	}
	for _, c := range inconclusive {
		if got := classifyClaude(c.status, c.body, claudeTestRegion); got.State == StateBlocked {
			t.Errorf("status %d claimed a region block from an inconclusive response: %q",
				c.status, got.Detail)
		}
	}
}
