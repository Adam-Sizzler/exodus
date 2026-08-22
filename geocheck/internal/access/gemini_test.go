package access

import "testing"

// ISO 3166-1 reserves XAA-XZZ for private use, so these stand in for "some
// country Google serves" without naming one. The blocked cases below are taken
// from geminiUnsupported itself rather than written out, which keeps any real
// country code — and anyone's location — out of this file.
const (
	servedCode  = "XAA"
	servedCode2 = "XAB"
)

// geminiPage reproduces the shape of the configuration block Google embeds in
// its pages. The country code sits at a fixed offset inside that array, which
// is what the check reads.
func geminiPage(alpha3 string) string {
	return `<!DOCTYPE html><html><head><title>&#8206;Google Gemini</title></head><body>` +
		`<script>window.WIZ_global_data={};var _bar_={CONFIG:[[[0,"www.gstatic.com",` +
		`"og.qtm.en_US.VN3NDeUdSgI.2019.O","example","en","658",0,[4,2,"","","","961884838","0"],` +
		`null,"YdSDapniAbW3-d8P15ie4Ao",null,0,"og.qtm.OCU-UcflS4s.L.W.O",` +
		`"AA2YrTsruG48AvvyaPPTzRzHSmp1qDGl3A","AA2YrTsiQfTsY-xZhEUlZp1Kwq7Yf8c9jw","",` +
		`2,1,200,"` + alpha3 + `",null,null,"269","658",1,null,null,103135050,null,0,0,0,0]]]}` +
		`</script></body></html>`
}

func TestClassifyGemini(t *testing.T) {
	cases := []struct {
		name   string
		status int
		body   string
		want   State
		region string
	}{
		{"a served country", 200, geminiPage(servedCode), StateAvailable, servedCode},
		{"another served country", 200, geminiPage(servedCode2), StateAvailable, servedCode2},
		{
			// The page is the only place the region can be read from, so losing
			// it must be reported, not guessed around.
			name: "no region in the page", status: 200,
			body: `<!DOCTYPE html><html><body>nothing useful here</body></html>`,
			want: StateError,
		},
		{
			name: "a challenge is not an answer", status: 403,
			body: `<html><head><title>Just a moment...</title></head><body>` +
				`<script>window._cf_chl_opt={};</script></body></html>`,
			want: StateError,
		},
		{"server error", 503, "Service Unavailable", StateError, ""},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := classifyGemini(c.status, c.body)
			if got.State != c.want {
				t.Errorf("state = %v, want %v (detail %q)", got.State, c.want, got.Detail)
			}
			if got.Region != c.region {
				t.Errorf("region = %q, want %q", got.Region, c.region)
			}
		})
	}
}

// TestEveryUnsupportedCountryIsBlocked walks the list instead of naming a
// country, so the whole of it is covered and adding an entry needs no test
// change.
func TestEveryUnsupportedCountryIsBlocked(t *testing.T) {
	if len(geminiUnsupported) == 0 {
		t.Fatal("the unsupported list is empty; this test would prove nothing")
	}
	for code := range geminiUnsupported {
		got := classifyGemini(200, geminiPage(code))
		if got.State != StateBlocked {
			t.Errorf("a country on the unsupported list came back %v, want blocked", got.State)
		}
		if got.Region != code {
			t.Error("the verdict did not carry the region it was based on")
		}
	}
}

// TestGeminiNeverGuessesFromAnAbsentRegion is the discipline this check exists
// under: the only route to a verdict is a country code actually present in the
// response. The widely copied shell implementation keys on an opaque feature
// flag instead, and reports "blocked" everywhere now that Google retired it.
func TestGeminiNeverGuessesFromAnAbsentRegion(t *testing.T) {
	for _, body := range []string{
		"",
		"<html></html>",
		`{"error":"nope"}`,
		geminiPage("xaa"), // lower case is not the documented form
	} {
		if got := classifyGemini(200, body); got.State != StateError {
			t.Errorf("body %.20q produced %v, want an inconclusive result", body, got.State)
		}
	}
}

// TestGeminiRegionPatternToleratesFieldShifts keeps the extraction from being
// pinned to the exact leading numbers, which are positional and not ours.
func TestGeminiRegionPatternToleratesFieldShifts(t *testing.T) {
	body := `var _bar_={CONFIG:[[[0,"x","y","example","en","658",0,` +
		`"",3,7,200,"` + servedCode + `",null,null,"269"]]]}`
	got := classifyGemini(200, body)
	if got.State != StateAvailable || got.Region != servedCode {
		t.Errorf("got %v/%q, want available/%s", got.State, got.Region, servedCode)
	}
}
