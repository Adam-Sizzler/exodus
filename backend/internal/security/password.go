package security

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/sha512"
	"crypto/subtle"
	"encoding/hex"
	"fmt"
	"strings"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/pbkdf2"
)

const (
	passwordSaltBytes = 16
	passwordKeyBytes  = 64
	// Keep a strong default and preserve the legacy "salt:hash" layout.
	passwordPBKDF2Iterations = 210_000
)

var legacyPBKDF2Iterations = []int{
	passwordPBKDF2Iterations,
	1_000, 5_000, 10_000, 20_000, 50_000, 100_000,
	150_000, 200_000, 250_000, 300_000, 310_000, 400_000, 500_000, 600_000, 750_000, 1_000_000,
}

// HashPassword returns a legacy-compatible "salt_hex:hash_hex" value.
func HashPassword(password string) (string, error) {
	password = strings.TrimSpace(password)
	if password == "" {
		return "", fmt.Errorf("password cannot be empty")
	}

	salt := make([]byte, passwordSaltBytes)
	if _, err := rand.Read(salt); err != nil {
		return "", fmt.Errorf("generate salt: %w", err)
	}

	derived := pbkdf2.Key([]byte(password), salt, passwordPBKDF2Iterations, passwordKeyBytes, sha512.New)
	return fmt.Sprintf("%s:%s", hex.EncodeToString(salt), hex.EncodeToString(derived)), nil
}

// HashPasswordBcrypt exists for future migrations and compatibility.
func HashPasswordBcrypt(password string) (string, error) {
	password = strings.TrimSpace(password)
	if password == "" {
		return "", fmt.Errorf("password cannot be empty")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("bcrypt hash failed: %w", err)
	}
	return string(hash), nil
}

// VerifyPassword validates password against supported stored formats.
// Supported formats:
// 1) "salt_hex:hash_hex" (legacy)
// 2) "pbkdf2$<iterations>$<salt_hex>$<hash_hex>"
// 3) bcrypt hash ($2a/$2b/$2y)
func VerifyPassword(password, stored string) bool {
	password = strings.TrimSpace(password)
	stored = strings.TrimSpace(stored)
	if password == "" || stored == "" {
		return false
	}

	if strings.HasPrefix(stored, "$2a$") || strings.HasPrefix(stored, "$2b$") || strings.HasPrefix(stored, "$2y$") {
		return bcrypt.CompareHashAndPassword([]byte(stored), []byte(password)) == nil
	}

	if strings.HasPrefix(stored, "pbkdf2$") {
		return verifyPBKDF2Tagged(password, stored)
	}

	if strings.Contains(stored, ":") {
		return verifyLegacySaltHash(password, stored)
	}

	return false
}

func verifyPBKDF2Tagged(password, stored string) bool {
	parts := strings.Split(stored, "$")
	if len(parts) != 4 {
		return false
	}
	if parts[0] != "pbkdf2" {
		return false
	}

	iterations := 0
	if _, err := fmt.Sscanf(parts[1], "%d", &iterations); err != nil || iterations <= 0 {
		return false
	}

	salt, err := hex.DecodeString(parts[2])
	if err != nil {
		return false
	}
	expected, err := hex.DecodeString(parts[3])
	if err != nil {
		return false
	}

	derived := pbkdf2.Key([]byte(password), salt, iterations, len(expected), sha512.New)
	return subtle.ConstantTimeCompare(derived, expected) == 1
}

func verifyLegacySaltHash(password, stored string) bool {
	parts := strings.SplitN(stored, ":", 2)
	if len(parts) != 2 {
		return false
	}

	saltHex := strings.TrimSpace(parts[0])
	hashHex := strings.TrimSpace(parts[1])
	if saltHex == "" || hashHex == "" {
		return false
	}

	salt, err := hex.DecodeString(saltHex)
	if err != nil {
		return false
	}
	expected, err := hex.DecodeString(hashHex)
	if err != nil {
		return false
	}

	// 1) Legacy PBKDF2 candidates.
	for _, it := range legacyPBKDF2Iterations {
		derived := pbkdf2.Key([]byte(password), salt, it, len(expected), sha512.New)
		if subtle.ConstantTimeCompare(derived, expected) == 1 {
			return true
		}
		derived = pbkdf2.Key([]byte(password), salt, it, len(expected), sha256.New)
		if subtle.ConstantTimeCompare(derived, expected) == 1 {
			return true
		}
	}

	return false
}
