package geo

import (
	"testing"

	"exodus-node/geocheck/internal/netx"
)

func TestNormalizeCountry(t *testing.T) {
	cases := []struct{ in, want string }{
		{"nl", "NL"},
		{" DE ", "DE"},
		{"\"US\"", "US"},
		{"null", ""},
		{"", ""},
		{"USA", ""},   // three letters is not an alpha-2 code
		{"X", ""},     // too short
		{"12", ""},    // digits are stripped, leaving nothing
		{"n/a", "NA"}, // punctuation is stripped; a real code shape survives
	}
	for _, c := range cases {
		if got := normalize(c.in, KindCountry); got != c.want {
			t.Errorf("normalize(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestNormalizeAvailability(t *testing.T) {
	if got := normalize("Yes", KindAvailability); got != "yes" {
		t.Errorf("normalize(Yes) = %q, want yes", got)
	}
	if got := normalize("null", KindAvailability); got != "" {
		t.Errorf("normalize(null) = %q, want empty", got)
	}
}

func TestSummarizeCountsCountriesOnly(t *testing.T) {
	results := []Result{
		{Check: Check{Group: GroupServices, Kind: KindCountry}, V4: Outcome{Value: "NL"}},
		{Check: Check{Group: GroupServices, Kind: KindCountry}, V4: Outcome{Value: "NL"}},
		{Check: Check{Group: GroupGeoIP, Kind: KindCountry}, V4: Outcome{Value: "DE"}},
		// Availability answers are not countries and must not be tallied.
		{Check: Check{Group: GroupServices, Kind: KindAvailability}, V4: Outcome{Value: "yes"}},
		// CDN answers describe the edge, not the exit, so they are excluded.
		{Check: Check{Group: GroupCDN, Kind: KindCountry}, V4: Outcome{Value: "FR"}},
		// Failures contribute nothing.
		{Check: Check{Group: GroupGeoIP, Kind: KindCountry}, V4: Outcome{Value: ""}},
	}

	got := Summarize(results, netx.V4)
	if len(got) != 2 {
		t.Fatalf("got %d countries, want 2: %+v", len(got), got)
	}
	if got[0].Code != "NL" || got[0].Count != 2 {
		t.Errorf("top entry = %+v, want NL x2", got[0])
	}
	if got[0].Total != 3 {
		t.Errorf("Total = %d, want 3", got[0].Total)
	}
	if got[0].Percent < 66 || got[0].Percent > 67 {
		t.Errorf("Percent = %v, want ~66.7", got[0].Percent)
	}
}

func TestSummarizeEmpty(t *testing.T) {
	if got := Summarize(nil, netx.V4); got != nil {
		t.Errorf("Summarize(nil) = %+v, want nil", got)
	}
}

func TestParseGGCCluster(t *testing.T) {
	cases := []struct{ in, want string }{
		{"192.0.2.0/24 => fra16s52", "FRA"},
		{"198.51.100.7 => exampleisp-arn2", "ARN"},
		{"203.0.113.0/24 => sto03s07\n", "STO"},
		{"garbage", ""},
		{"", ""},
	}
	for _, c := range cases {
		if got := parseGGCCluster(c.in); got != c.want {
			t.Errorf("parseGGCCluster(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestCheckSetsAreWellFormed(t *testing.T) {
	all := append(ServiceChecks(), DatabaseChecks()...)
	all = append(all, CDNChecks()...)

	seen := map[string]bool{}
	for _, c := range all {
		if c.ID == "" || c.Name == "" || c.Run == nil {
			t.Errorf("check %+v is incomplete", c)
		}
		if seen[c.ID] {
			t.Errorf("duplicate check id %q", c.ID)
		}
		seen[c.ID] = true
	}
	// Dependencies must name a check that actually exists, or the fallback
	// silently never fires.
	for _, c := range all {
		if c.DependsOn != "" && !seen[c.DependsOn] {
			t.Errorf("check %q depends on unknown check %q", c.ID, c.DependsOn)
		}
	}
}

func TestIATATableCoversCommonPoPs(t *testing.T) {
	for code, want := range map[string]string{
		"ARN": "SE", "FRA": "DE", "AMS": "NL", "LHR": "GB",
		"IAD": "US", "NRT": "JP", "GRU": "BR", "SIN": "SG",
	} {
		if got := iataToCountry[code]; got != want {
			t.Errorf("iataToCountry[%q] = %q, want %q", code, got, want)
		}
	}
}
