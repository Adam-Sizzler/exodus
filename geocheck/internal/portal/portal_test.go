package portal

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"exodus-node/geocheck/internal/netx"
)

// newStack builds a stack pointed at the local test server. The endpoints are
// plain HTTP on loopback, which the stack dials directly.
func newStack(t *testing.T) *netx.Stack {
	t.Helper()
	s, err := netx.New(context.Background(), netx.Options{
		Timeout: 5 * time.Second,
		DoH:     "off",
	})
	if err != nil {
		t.Fatalf("netx.New: %v", err)
	}
	t.Cleanup(s.Close)
	return s
}

func runOne(t *testing.T, ep Endpoint) Result {
	t.Helper()
	res := Run(context.Background(), newStack(t), netx.V4, []Endpoint{ep}, 1)
	if len(res) != 1 {
		t.Fatalf("got %d results, want 1", len(res))
	}
	return res[0]
}

func TestVerdictOKOn204(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	got := runOne(t, Endpoint{ID: "t", Name: "t", URL: srv.URL, WantStatus: 204})
	if got.Verdict != VerdictOK {
		t.Fatalf("Verdict = %v (%s), want ok", got.Verdict, got.Detail)
	}
}

func TestVerdictPortalOnRedirect(t *testing.T) {
	// A captive portal answers the probe with a redirect to its login page.
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		http.Redirect(w, &http.Request{}, "http://portal.example/login", http.StatusFound)
	}))
	defer srv.Close()

	got := runOne(t, Endpoint{ID: "t", Name: "t", URL: srv.URL, WantStatus: 204})
	if got.Verdict != VerdictPortal {
		t.Fatalf("Verdict = %v, want portal", got.Verdict)
	}
	if !strings.Contains(got.Redirect, "portal.example") {
		t.Errorf("Redirect = %q, want the Location header", got.Redirect)
	}
	// The redirect must not be followed, or the portal's page would be
	// reported as though it were the endpoint's own answer.
	if got.Status != http.StatusFound {
		t.Errorf("Status = %d, want 302; the redirect was followed", got.Status)
	}
}

func TestVerdictAlteredOn200InsteadOf204(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	got := runOne(t, Endpoint{ID: "t", Name: "t", URL: srv.URL, WantStatus: 204})
	if got.Verdict != VerdictAltered {
		t.Fatalf("Verdict = %v, want altered", got.Verdict)
	}
	if !strings.Contains(got.Detail, "transparent-proxy") {
		t.Errorf("Detail = %q, want the 200-instead-of-204 explanation", got.Detail)
	}
}

func TestVerdictAlteredOnInjectedBody(t *testing.T) {
	// A proxy that answers with the right status but injects content where the
	// endpoint promises nothing.
	//
	// Note this uses a 200 rather than a 204: net/http refuses to write a body
	// alongside a 204, exactly as RFC 9110 requires, so a 204-with-content
	// cannot be produced through the standard server. The injection that is
	// actually observable in the wild arrives as a 200, which is covered here
	// and by TestVerdictAlteredOn200InsteadOf204.
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte("<html>ads</html>"))
	}))
	defer srv.Close()

	got := runOne(t, Endpoint{ID: "t", Name: "t", URL: srv.URL, WantStatus: 200})
	if got.Verdict != VerdictAltered {
		t.Fatalf("Verdict = %v, want altered", got.Verdict)
	}
	if !strings.Contains(got.Detail, "empty body") {
		t.Errorf("Detail = %q, want it to name the unexpected body", got.Detail)
	}
}

func TestBodyMatching(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte("success\n"))
	}))
	defer srv.Close()

	t.Run("exact body matches", func(t *testing.T) {
		got := runOne(t, Endpoint{ID: "t", Name: "t", URL: srv.URL, WantStatus: 200, WantBody: "success\n"})
		if got.Verdict != VerdictOK {
			t.Fatalf("Verdict = %v (%s), want ok", got.Verdict, got.Detail)
		}
	})

	t.Run("trailing whitespace is tolerated", func(t *testing.T) {
		got := runOne(t, Endpoint{ID: "t", Name: "t", URL: srv.URL, WantStatus: 200, WantBody: "success"})
		if got.Verdict != VerdictOK {
			t.Fatalf("Verdict = %v (%s), want ok", got.Verdict, got.Detail)
		}
	})

	t.Run("wrong body is altered", func(t *testing.T) {
		got := runOne(t, Endpoint{ID: "t", Name: "t", URL: srv.URL, WantStatus: 200, WantBody: "different"})
		if got.Verdict != VerdictAltered {
			t.Fatalf("Verdict = %v, want altered", got.Verdict)
		}
	})

	t.Run("substring match", func(t *testing.T) {
		got := runOne(t, Endpoint{ID: "t", Name: "t", URL: srv.URL, WantStatus: 200, BodyContains: "succ"})
		if got.Verdict != VerdictOK {
			t.Fatalf("Verdict = %v (%s), want ok", got.Verdict, got.Detail)
		}
	})
}

func TestUnreachable(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}))
	url := srv.URL
	srv.Close() // nothing is listening now

	got := runOne(t, Endpoint{ID: "t", Name: "t", URL: url, WantStatus: 204})
	if got.Verdict != VerdictUnreachable {
		t.Fatalf("Verdict = %v, want unreachable", got.Verdict)
	}
	if got.Err == nil {
		t.Error("Err should be set for an unreachable endpoint")
	}
}

func TestSummarizeSplitsBySchemeAndDetectsBlockedHTTP(t *testing.T) {
	results := []Result{
		{Endpoint: Endpoint{URL: "http://a/"}, Verdict: VerdictUnreachable},
		{Endpoint: Endpoint{URL: "http://b/"}, Verdict: VerdictUnreachable},
		{Endpoint: Endpoint{URL: "https://c/"}, Verdict: VerdictOK},
	}
	s := Summarize(results)

	if s.Total != 3 || s.OK != 1 || s.Unreachable != 2 {
		t.Fatalf("Summary = %+v", s)
	}
	if s.PlainTotal != 2 || s.PlainOK != 0 || s.TLSOK != 1 {
		t.Fatalf("scheme split wrong: %+v", s)
	}
	if !s.PlainHTTPBlocked() {
		t.Error("PlainHTTPBlocked should be true when only HTTPS answers")
	}
	if s.Clean() {
		t.Error("Clean should be false when plain HTTP never answered")
	}
}

func TestSummarizeCleanRun(t *testing.T) {
	results := []Result{
		{Endpoint: Endpoint{URL: "http://a/"}, Verdict: VerdictOK},
		{Endpoint: Endpoint{URL: "https://b/"}, Verdict: VerdictOK},
	}
	s := Summarize(results)
	if !s.Clean() {
		t.Error("Clean should be true when everything answered correctly")
	}
	if s.PlainHTTPBlocked() {
		t.Error("PlainHTTPBlocked should be false when plain HTTP answered")
	}
}

func TestCatalogIntegrity(t *testing.T) {
	seen := map[string]bool{}
	for _, ep := range Catalog {
		if ep.ID == "" || ep.Name == "" || ep.URL == "" || ep.Vendor == "" {
			t.Errorf("endpoint %+v has an empty required field", ep)
		}
		if seen[ep.ID] {
			t.Errorf("duplicate endpoint id %q", ep.ID)
		}
		seen[ep.ID] = true

		if ep.WantStatus == 0 {
			t.Errorf("endpoint %q has no expected status", ep.ID)
		}
		// A 204 carries no body by definition, so specifying one is a mistake
		// that would make the check always fail.
		if ep.WantStatus == 204 && (ep.WantBody != "" || ep.BodyContains != "") {
			t.Errorf("endpoint %q expects 204 yet also expects a body", ep.ID)
		}
		if !strings.HasPrefix(ep.URL, "http://") && !strings.HasPrefix(ep.URL, "https://") {
			t.Errorf("endpoint %q has a malformed URL %q", ep.ID, ep.URL)
		}
	}

	if len(DefaultEndpoints()) == 0 {
		t.Error("the default endpoint set is empty")
	}

	// Endpoints with a known-legitimate TLS quirk must never be fetched over
	// HTTPS, or the run reports a shared CDN certificate as interception.
	for _, ep := range Catalog {
		if reason, unsafe := HTTPSUnsafe(ep.ID); unsafe {
			if strings.HasPrefix(ep.URL, "https://") {
				t.Errorf("endpoint %q is fetched over HTTPS but %s", ep.ID, reason)
			}
		}
	}
	// The default set must span vendors, or agreement between endpoints would
	// say nothing beyond "one provider is reachable".
	vendors := map[string]bool{}
	for _, ep := range DefaultEndpoints() {
		vendors[ep.Vendor] = true
	}
	if len(vendors) < 4 {
		t.Errorf("default set covers only %d vendors, want at least 4", len(vendors))
	}
}

func TestSelect(t *testing.T) {
	if got := Select("all"); len(got) != len(Catalog) {
		t.Errorf("Select(all) = %d, want %d", len(got), len(Catalog))
	}
	if got := Select("nope"); len(got) != 0 {
		t.Errorf("Select of an unknown tag = %d, want 0", len(got))
	}
	if got := Select("apple"); len(got) != 1 || got[0].ID != "apple" {
		t.Errorf("Select by id returned %+v", got)
	}
}
