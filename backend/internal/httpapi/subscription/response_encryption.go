package subscription

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"filippo.io/age"
	"filippo.io/age/armor"
)

// encryptResponseBody encrypts the response body with age (X25519 age1 or post-quantum hybrid age1pq1)
// and returns the ASCII-armored ciphertext.
func encryptResponseBody(body []byte, method, key string) (string, error) {
	key = strings.TrimSpace(key)
	method = strings.ToLower(strings.TrimSpace(method))

	switch method {
	case "age1", "age1pq1":
		// supported
	default:
		return "", fmt.Errorf("unsupported age encryption method: %q", method)
	}

	if key == "" {
		return "", fmt.Errorf("age encryption key cannot be empty")
	}

	recipients, err := age.ParseRecipients(strings.NewReader(key))
	if err != nil {
		return "", fmt.Errorf("parse age recipient key: %w", err)
	}
	if len(recipients) == 0 {
		return "", fmt.Errorf("no valid age recipient found in key")
	}

	// Validate method matches key type
	_, isHybrid := recipients[0].(*age.HybridRecipient)
	if method == "age1pq1" && !isHybrid {
		return "", fmt.Errorf("method is age1pq1 but provided key is not a hybrid (post-quantum) recipient")
	}
	if method == "age1" && isHybrid {
		return "", fmt.Errorf("method is age1 but provided key is a hybrid (post-quantum) recipient")
	}

	var out bytes.Buffer
	armorWriter := armor.NewWriter(&out)
	w, err := age.Encrypt(armorWriter, recipients...)
	if err != nil {
		return "", fmt.Errorf("init age encryption: %w", err)
	}

	if _, err := io.Copy(w, bytes.NewReader(body)); err != nil {
		return "", fmt.Errorf("write encrypted body: %w", err)
	}
	if err := w.Close(); err != nil {
		return "", fmt.Errorf("close age encrypter: %w", err)
	}
	if err := armorWriter.Close(); err != nil {
		return "", fmt.Errorf("close age armor writer: %w", err)
	}

	return out.String(), nil
}
