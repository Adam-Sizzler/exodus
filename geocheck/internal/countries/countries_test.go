package countries

import "testing"

func TestName(t *testing.T) {
	if got := Name("nl"); got != "Netherlands" {
		t.Errorf("Name(nl) = %q, want Netherlands", got)
	}
	// Unknown codes come back unchanged so there is always something to print.
	if got := Name("ZZ"); got != "ZZ" {
		t.Errorf("Name(ZZ) = %q, want ZZ", got)
	}
}

func TestCode(t *testing.T) {
	cases := []struct{ in, want string }{
		{"Netherlands", "NL"},
		{"united states", "US"},
		{"United States of America", "US"},
		{"UK", "GB"},
		{"Türkiye", "TR"},
		{"Nowhere", ""},
	}
	for _, c := range cases {
		if got := Code(c.in); got != c.want {
			t.Errorf("Code(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestKnown(t *testing.T) {
	if !Known("de") {
		t.Error("Known(de) = false, want true")
	}
	if Known("ZZ") {
		t.Error("Known(ZZ) = true, want false")
	}
}
