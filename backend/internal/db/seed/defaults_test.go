package seed

import (
	"testing"
)

func TestCanonicalHashMatch(t *testing.T) {
	h, err := canonicalHash(defaultResponseRules)
	if err != nil {
		t.Fatalf("failed to calculate canonical hash: %v", err)
	}
	if h != PrevResponseRulesHash {
		t.Errorf("Canonical hash mismatch!\nGot:      %s\nExpected: %s", h, PrevResponseRulesHash)
	}
}
