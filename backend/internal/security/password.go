package security

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"fmt"
	"strings"

	"golang.org/x/crypto/scrypt"
)

// This mirrors upstream Exodus's own password scheme (see
// src/modules/auth/auth.service.ts: applySecretHmac + scrypt), ported
// idiomatically to Go:
//
//  1. Pepper: HMAC-SHA256(password, secret) where secret is the panel's
//     JWT auth secret (cfg.JWT.AuthSecret) - a value that lives only in
//     server config/env, never in the database. A full DB dump (salts +
//     hashes) alone is not enough to attempt an offline dictionary attack;
//     the secret has to be known too.
//  2. scrypt (memory-hard, resists GPU/ASIC-accelerated offline
//     cracking far better than a CPU-only KDF like PBKDF2) over the
//     peppered value, with a random per-password salt.
//
// There is deliberately no support for any other stored hash format and
// no migration path from the project's previous PBKDF2-based scheme -
// this project only ever has one admin password in practice, reset by
// hand via CLI after this change ships (see release notes/commit
// message), rather than carrying legacy-format verification code
// indefinitely for a single account.

const (
	passwordSaltBytes = 16
	passwordKeyBytes  = 64

	// scrypt cost parameters. N=16384, r=8, p=1 match Node's
	// crypto.scrypt() defaults, which is what upstream Exodus uses
	// (it never overrides them), so this preserves the same effective
	// work factor as upstream rather than an arbitrary Go-side choice.
	scryptN = 16384
	scryptR = 8
	scryptP = 1
)

// pepperPassword applies the HMAC-SHA256 pepper step, matching upstream's
// applySecretHmac(password, jwtSecret). Upstream then hex-encodes the HMAC
// digest before feeding it to scrypt (hmacResult.toString('hex')) rather
// than using the raw bytes directly - preserved here for parity, though it
// has no security effect either way.
func pepperPassword(password, secret string) []byte {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(password))
	return []byte(hex.EncodeToString(mac.Sum(nil)))
}

// HashPassword returns a "salt_hex:hash_hex" value. secret must be the
// panel's JWT auth secret (cfg.JWT.AuthSecret) - the exact same secret
// must be supplied to VerifyPassword later, or the password will never
// verify again (this is the whole point of a pepper: it is not stored
// anywhere alongside the hash).
func HashPassword(password, secret string) (string, error) {
	if password == "" {
		return "", fmt.Errorf("password cannot be empty")
	}
	if strings.TrimSpace(secret) == "" {
		return "", fmt.Errorf("secret cannot be empty")
	}

	saltBytes := make([]byte, passwordSaltBytes)
	if _, err := rand.Read(saltBytes); err != nil {
		return "", fmt.Errorf("generate salt: %w", err)
	}
	// Upstream passes the salt to Node's scrypt() as a *string*
	// (randomBytes(16).toString('hex')) with no explicit encoding, so Node
	// treats it as UTF-8 - meaning the actual bytes scrypt uses as salt
	// are the ASCII bytes of the hex string itself (32 bytes), not the 16
	// raw random bytes it represents. Matched here deliberately, not a
	// mistake: hex-encode first, then use *those* bytes as the scrypt salt.
	saltHex := hex.EncodeToString(saltBytes)

	derived, err := scrypt.Key(pepperPassword(password, secret), []byte(saltHex), scryptN, scryptR, scryptP, passwordKeyBytes)
	if err != nil {
		return "", fmt.Errorf("derive key: %w", err)
	}

	return fmt.Sprintf("%s:%s", saltHex, hex.EncodeToString(derived)), nil
}

// VerifyPassword checks password (peppered with secret, the panel's JWT
// auth secret) against a hash produced by HashPassword.
func VerifyPassword(password, secret, stored string) bool {
	stored = strings.TrimSpace(stored)
	if password == "" || strings.TrimSpace(secret) == "" || stored == "" {
		return false
	}

	parts := strings.SplitN(stored, ":", 2)
	if len(parts) != 2 {
		return false
	}

	saltHex := strings.TrimSpace(parts[0])
	expected, err := hex.DecodeString(strings.TrimSpace(parts[1]))
	if err != nil {
		return false
	}

	// Upstream uses the UTF-8 bytes of the salt's hex representation as the actual scrypt salt
	derived, err := scrypt.Key(pepperPassword(password, secret), []byte(saltHex), scryptN, scryptR, scryptP, len(expected))
	if err != nil {
		return false
	}

	return subtle.ConstantTimeCompare(derived, expected) == 1
}
