package access

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"exodus-node/geocheck/internal/netx"
)

func testEnv(t *testing.T) Env {
	t.Helper()
	s, err := netx.New(context.Background(), netx.Options{
		Timeout: 5 * time.Second,
		DoH:     "off",
	})
	if err != nil {
		t.Fatalf("netx.New: %v", err)
	}
	t.Cleanup(s.Close)
	return Env{Stack: s, Family: netx.V4}
}

// netflixServer mimics Netflix's behaviour: a redirect to a locale-prefixed
// path for titles it serves, and a 404 for those it does not.
func netflixServer(t *testing.T, locale string, serve map[string]bool) *httptest.Server {
	t.Helper()
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		// Already redirected: serve it.
		if locale != "" && strings.HasPrefix(path, "/"+locale+"/title/") {
			id := strings.TrimPrefix(path, "/"+locale+"/title/")
			if serve[id] {
				_, _ = w.Write([]byte("ok"))
				return
			}
			w.WriteHeader(http.StatusNotFound)
			return
		}
		id := strings.TrimPrefix(path, "/title/")
		if !serve[id] {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		if locale == "" {
			_, _ = w.Write([]byte("ok"))
			return
		}
		http.Redirect(w, r, "/"+locale+"/title/"+id, http.StatusFound)
	})
	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)
	return srv
}

func TestNetflixRegionFromRedirect(t *testing.T) {
	env := testEnv(t)

	t.Run("locale prefix names the catalogue", func(t *testing.T) {
		srv := netflixServer(t, "fr-en", map[string]bool{netflixLicensedTitle: true})
		region, ok, err := netflixTitleAt(context.Background(), env, srv.URL, netflixLicensedTitle)
		if err != nil || !ok {
			t.Fatalf("ok=%v err=%v", ok, err)
		}
		if region != "FR" {
			t.Errorf("region = %q, want FR", region)
		}
	})

	t.Run("bare /title means the US catalogue", func(t *testing.T) {
		srv := netflixServer(t, "", map[string]bool{netflixLicensedTitle: true})
		region, ok, err := netflixTitleAt(context.Background(), env, srv.URL, netflixLicensedTitle)
		if err != nil || !ok {
			t.Fatalf("ok=%v err=%v", ok, err)
		}
		if region != "US" {
			t.Errorf("region = %q, want US", region)
		}
	})

	t.Run("404 means the title is withheld", func(t *testing.T) {
		srv := netflixServer(t, "fr-en", map[string]bool{})
		_, ok, err := netflixTitleAt(context.Background(), env, srv.URL, netflixLicensedTitle)
		if err != nil {
			t.Fatalf("err = %v", err)
		}
		if ok {
			t.Error("a 404 must not count as served")
		}
	})
}

func TestStateStrings(t *testing.T) {
	for state, want := range map[State]string{
		StateAvailable:  "available",
		StateRestricted: "restricted",
		StateBlocked:    "blocked",
		StateError:      "error",
	} {
		if got := state.String(); got != want {
			t.Errorf("State(%d).String() = %q, want %q", state, got, want)
		}
	}
}

func TestSummarize(t *testing.T) {
	s := Summarize([]Result{
		{State: StateAvailable},
		{State: StateAvailable},
		{State: StateRestricted},
		{State: StateBlocked},
		{State: StateError},
	})
	if s.Total != 5 || s.Available != 2 || s.Restricted != 1 || s.Blocked != 1 || s.Errors != 1 {
		t.Fatalf("Summary = %+v", s)
	}
}

func TestChecksAreWellFormed(t *testing.T) {
	seen := map[string]bool{}
	for _, c := range Checks() {
		if c.ID == "" || c.Name == "" || c.Run == nil {
			t.Errorf("check %+v is incomplete", c)
		}
		if seen[c.ID] {
			t.Errorf("duplicate check id %q", c.ID)
		}
		seen[c.ID] = true
	}
	if len(Checks()) < 5 {
		t.Errorf("got %d checks, want at least 5", len(Checks()))
	}
}

func TestRunSetsCheckAndTiming(t *testing.T) {
	// The probe sleeps so that the elapsed time is unambiguously non-zero.
	// Returning instantly made this test flaky: both clock readings could land
	// in the same tick, Run would record a legitimate 0, and the assertion
	// below would fail perhaps one run in six for no reason worth chasing.
	probe := Check{ID: "x", Name: "X", Run: func(context.Context, Env) Result {
		time.Sleep(time.Millisecond)
		return Result{State: StateAvailable}
	}}
	got := Run(context.Background(), testEnv(t).Stack, netx.V4, []Check{probe}, 1)
	if len(got) != 1 {
		t.Fatalf("got %d results, want 1", len(got))
	}
	if got[0].Check.ID != "x" {
		t.Errorf("Check was not attached to the result")
	}
	// A check that does not time itself still gets a duration.
	if got[0].RTT <= 0 {
		t.Errorf("RTT = %v, want it filled in by Run", got[0].RTT)
	}
}

func TestItoa(t *testing.T) {
	for in, want := range map[int]string{0: "0", 7: "7", 403: "403", 1234: "1234", -5: "-5"} {
		if got := itoa(in); got != want {
			t.Errorf("itoa(%d) = %q, want %q", in, got, want)
		}
	}
}
