package server

import (
	"strings"
	"testing"
)

func TestGetStartMessageMatchesRemnawaveStyle(t *testing.T) {
	message := GetStartMessage("v1.2.3")
	if !strings.Contains(message, "Exodus Subscription Page v1.2.3") {
		t.Fatalf("missing title: %s", message)
	}
	if strings.Contains(message, "│ Docs") {
		t.Fatalf("startup table must be single-column, got: %s", message)
	}
	if !strings.Contains(message, "Docs → https://docs.exodus.dev") {
		t.Fatalf("missing docs row: %s", message)
	}
	if !strings.Contains(message, "Community → https://t.me/exodus") {
		t.Fatalf("missing community row: %s", message)
	}
}
