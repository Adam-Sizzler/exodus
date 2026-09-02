package cmd

import (
	"strings"
	"testing"

	"filippo.io/age"
)

func TestAgeKeypairGeneration(t *testing.T) {
	// Test standard X25519 (age1)
	id, err := age.GenerateX25519Identity()
	if err != nil {
		t.Fatalf("GenerateX25519Identity failed: %v", err)
	}
	secKey := id.String()
	pubKey := id.Recipient().String()

	if !strings.HasPrefix(secKey, "AGE-SECRET-KEY-1") {
		t.Errorf("Expected secret key to start with AGE-SECRET-KEY-1, got %s", secKey)
	}
	if !strings.HasPrefix(pubKey, "age1") {
		t.Errorf("Expected recipient to start with age1, got %s", pubKey)
	}

	// Test hybrid post-quantum (age1pq1)
	pqId, err := age.GenerateHybridIdentity()
	if err != nil {
		t.Fatalf("GenerateHybridIdentity failed: %v", err)
	}
	pqSecKey := pqId.String()
	pqPubKey := pqId.Recipient().String()

	if !strings.HasPrefix(pqSecKey, "AGE-SECRET-KEY-PQ-1") {
		t.Errorf("Expected secret key to start with AGE-SECRET-KEY-PQ-1, got %s", pqSecKey)
	}
	if !strings.HasPrefix(pqPubKey, "age1pq1") {
		t.Errorf("Expected recipient to start with age1pq1, got %s", pqPubKey)
	}
}
