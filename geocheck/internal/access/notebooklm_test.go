package access

import "testing"

// The destinations below are the real ones, taken from live runs: a served
// region is sent to a Google sign-in, an unserved one to a marketing page that
// states the reason in the query string. Neither names a country, so nothing
// here reveals where the captures were taken.
const (
	signInURL = "https://accounts.google.com/v3/signin/identifier?continue=" +
		"https://notebook.google.com/login&followup=https://notebook.google.com/login" +
		"&osid=1&passive=1209600&flowName=GlifWebSignIn&flowEntry=ServiceLogin"
	loginURL      = "https://notebook.google.com/login?continue=https://notebook.google.com/"
	unsupportedLM = "https://notebook.google/?location=unsupported"
)

func TestClassifyNotebookLM(t *testing.T) {
	cases := []struct {
		name   string
		status int
		final  string
		body   string
		want   State
	}{
		{
			// Reaching a sign-in is the evidence: the product needs an account,
			// so this is as far as an unauthenticated check can get in a region
			// that is served.
			name:   "sign-in means the region is served",
			status: 200, final: signInURL, want: StateAvailable,
		},
		{
			name:   "the product's own login path counts too",
			status: 200, final: loginURL, want: StateAvailable,
		},
		{
			name:   "the redirect states the refusal",
			status: 200, final: unsupportedLM, want: StateBlocked,
		},
		{
			// Nothing to read the answer from.
			name: "no destination", status: 200, final: "", want: StateError,
		},
		{
			name:   "an unrecognised destination is not a verdict",
			status: 200, final: "https://example.com/somewhere", want: StateError,
		},
		{
			name:   "a challenge is not an answer",
			status: 403,
			final:  "https://notebooklm.google.com/",
			body: `<html><head><title>Just a moment...</title></head><body>` +
				`<script>window._cf_chl_opt={};</script></body></html>`,
			want: StateError,
		},
		{
			name: "server error", status: 502,
			final: "https://notebooklm.google.com/", want: StateError,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := classifyNotebookLM(c.status, c.final, c.body)
			if got.State != c.want {
				t.Errorf("state = %v, want %v (detail %q)", got.State, c.want, got.Detail)
			}
		})
	}
}

// TestNotebookLMRefusalWinsOverTheHost guards the ordering. The refusal is
// served from a *.google host too, so testing the host first would read a
// stated refusal as a served region.
func TestNotebookLMRefusalWinsOverTheHost(t *testing.T) {
	for _, final := range []string{
		unsupportedLM,
		"https://notebook.google.com/login?location=unsupported",
		"https://accounts.google.com/signin?location=unsupported",
	} {
		if got := classifyNotebookLM(200, final, ""); got.State != StateBlocked {
			t.Errorf("%q came back %v, want blocked", final, got.State)
		}
	}
}

// TestNotebookLMNeedsNoCountryList is what makes this check age better than the
// Gemini one next to it: the service says why, so nothing has to be maintained
// as geography changes.
func TestNotebookLMNeedsNoCountryList(t *testing.T) {
	if got := classifyNotebookLM(200, signInURL, ""); got.State != StateAvailable {
		t.Fatalf("a served region should not depend on any list, got %v", got.State)
	}
	if got := classifyNotebookLM(200, unsupportedLM, ""); got.Region != "" {
		t.Errorf("the verdict invented a region (%q); the redirect names none", got.Region)
	}
}
