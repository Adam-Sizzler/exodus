package security

import (
	"strings"
	"testing"
)

func TestResolveGRPCAuthToken(t *testing.T) {
	first, err := ResolveGRPCAuthToken("")
	if err != nil {
		t.Fatalf("generate first token: %v", err)
	}
	second, err := ResolveGRPCAuthToken("")
	if err != nil {
		t.Fatalf("generate second token: %v", err)
	}

	if len(first) != 64 {
		t.Fatalf("unexpected generated token length: %d", len(first))
	}
	if first == second {
		t.Fatal("generated tokens must be unique")
	}

	normalized, err := ResolveGRPCAuthToken(strings.ToUpper(first))
	if err != nil {
		t.Fatalf("normalize token: %v", err)
	}
	if normalized != first {
		t.Fatalf("unexpected normalized token: %s", normalized)
	}

	if _, err := ResolveGRPCAuthToken("not-a-valid-token"); err == nil {
		t.Fatal("expected invalid token error")
	}
}
