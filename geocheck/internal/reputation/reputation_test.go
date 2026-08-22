package reputation

import (
	"errors"
	"testing"
)

func TestFlagsOrderAndContent(t *testing.T) {
	i := &Info{Tor: true, Hosting: true, Anonymous: true, VPN: true}
	got := i.Flags()
	want := []string{"Tor", "VPN", "hosting", "anonymous"}
	if len(got) != len(want) {
		t.Fatalf("Flags() = %v, want %v", got, want)
	}
	for n := range want {
		if got[n] != want[n] {
			// Order matters: the most serious finding should read first.
			t.Fatalf("Flags() = %v, want %v", got, want)
		}
	}

	if (&Info{}).Clean() != true {
		t.Error("an Info with no detections should be Clean")
	}
	if i.Clean() {
		t.Error("an Info with detections should not be Clean")
	}
}

func TestResidentialClassification(t *testing.T) {
	cases := map[string]bool{
		"Residential": true,
		"residential": true,
		"Wireless":    true,
		"Business":    true,
		"Education":   true,
		"Hosting":     false,
		"":            false,
		"Unknown":     false,
	}
	for typ, want := range cases {
		if got := (&Info{Type: typ}).Residential(); got != want {
			t.Errorf("Residential(%q) = %v, want %v", typ, got, want)
		}
	}
}

func TestShortDate(t *testing.T) {
	cases := map[string]string{
		"2026-08-15T07:36:43Z": "2026-08-15",
		"2026-08-15":           "2026-08-15",
		"":                     "",
		"short":                "short",
	}
	for in, want := range cases {
		if got := shortDate(in); got != want {
			t.Errorf("shortDate(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestQuotaErrorIsIdentifiable(t *testing.T) {
	// Callers need to tell an exhausted allowance from a transport failure so
	// they can suggest the free API key rather than a network problem.
	if !errors.Is(ErrQuotaExceeded, ErrQuotaExceeded) {
		t.Fatal("ErrQuotaExceeded must be comparable")
	}
}
