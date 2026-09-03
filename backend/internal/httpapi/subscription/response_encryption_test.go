package subscription

import (
	"bytes"
	"io"
	"strings"
	"testing"

	"filippo.io/age"
	"filippo.io/age/armor"
)

func TestEncryptResponseBody(t *testing.T) {
	identity, err := age.GenerateX25519Identity()
	if err != nil {
		t.Fatalf("failed to generate identity: %v", err)
	}
	pubKey := identity.Recipient().String()

	plaintext := []byte("vless://test-uuid@1.2.3.4:443?type=ws#TestHost\n")

	t.Run("successful age1 encryption and decryption", func(t *testing.T) {
		ciphertext, err := encryptResponseBody(plaintext, "age1", pubKey)
		if err != nil {
			t.Fatalf("encryptResponseBody failed: %v", err)
		}

		if !strings.HasPrefix(ciphertext, "-----BEGIN AGE ENCRYPTED FILE-----") {
			t.Errorf("ciphertext missing armor prefix: %s", ciphertext)
		}

		// Decrypt and verify
		armorReader := armor.NewReader(strings.NewReader(ciphertext))
		r, err := age.Decrypt(armorReader, identity)
		if err != nil {
			t.Fatalf("age.Decrypt failed: %v", err)
		}

		var decrypted bytes.Buffer
		if _, err := io.Copy(&decrypted, r); err != nil {
			t.Fatalf("read decrypted content: %v", err)
		}

		if !bytes.Equal(decrypted.Bytes(), plaintext) {
			t.Errorf("decrypted content = %q; want %q", decrypted.String(), string(plaintext))
		}
	})

	t.Run("invalid method returns error", func(t *testing.T) {
		_, err := encryptResponseBody(plaintext, "rsa", pubKey)
		if err == nil {
			t.Error("expected error for unsupported method, got nil")
		}
	})

	t.Run("empty key returns error", func(t *testing.T) {
		_, err := encryptResponseBody(plaintext, "age1", "")
		if err == nil {
			t.Error("expected error for empty key, got nil")
		}
	})

	t.Run("malformed key returns error", func(t *testing.T) {
		_, err := encryptResponseBody(plaintext, "age1", "not-a-valid-key")
		if err == nil {
			t.Error("expected error for malformed key, got nil")
		}
	})

	t.Run("age1pq1 method rejects standard age1 key", func(t *testing.T) {
		_, err := encryptResponseBody(plaintext, "age1pq1", pubKey)
		if err == nil {
			t.Error("expected error when method is age1pq1 but key is age1, got nil")
		}
	})
}
