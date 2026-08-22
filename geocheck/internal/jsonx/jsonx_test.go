package jsonx

import "testing"

var sample = []byte(`{
  "country": "NL",
  "location": {"country": {"code": "DE"}},
  "list": [{"country": "FR"}, {"country": "ES"}],
  "num": 15169,
  "flag": false,
  "nothing": null
}`)

func TestString(t *testing.T) {
	cases := []struct{ path, want string }{
		{"country", "NL"},
		{"location.country.code", "DE"},
		{"list.0.country", "FR"},
		{"list.1.country", "ES"},
		{"num", "15169"},
		{"flag", "false"},
		{"nothing", ""},
		{"missing", ""},
		{"country.deeper", ""},
		{"list.9.country", ""},
		{"list.x", ""},
	}
	for _, c := range cases {
		if got := String(sample, c.path); got != c.want {
			t.Errorf("String(%q) = %q, want %q", c.path, got, c.want)
		}
	}
}

func TestStringOnInvalidJSON(t *testing.T) {
	if got := String([]byte("<html>not json</html>"), "country"); got != "" {
		t.Errorf("String on invalid JSON = %q, want empty", got)
	}
}

func TestTopLevelArray(t *testing.T) {
	if got := String([]byte(`[{"country":"IT"}]`), "0.country"); got != "IT" {
		t.Errorf("String = %q, want IT", got)
	}
}
