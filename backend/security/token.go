package security

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

// GenerateRandomToken returns a URL-safe random token.
func GenerateRandomToken(byteLen int) (string, error) {
	if byteLen < 16 {
		byteLen = 16
	}
	raw := make([]byte, byteLen)
	if _, err := rand.Read(raw); err != nil {
		return "", fmt.Errorf("generate random token: %w", err)
	}
	return base64.RawURLEncoding.EncodeToString(raw), nil
}
