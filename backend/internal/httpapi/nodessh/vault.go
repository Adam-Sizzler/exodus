package nodessh

import (
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"fmt"
	"io"
	"math/big"
	"sync"

	"github.com/gtank/ristretto255"
	"golang.org/x/crypto/hkdf"
)

var (
	vaultScalarLock sync.RWMutex
	cachedSecret    string
	cachedScalar    *ristretto255.Scalar
)

// expandMessageXMD_SHA512 implements hash_to_field (RFC 9380 §5.3) with SHA-512.
//
// NOTE: lenInBytes is hardcoded to 64 in the only call site (getOrDeriveVaultScalar).
// RFC 9497 §3.2 for ristretto255-SHA512 specifies L = ceil((255+128)/8) = 48.
// Using L=64 (one full SHA-512 block) is a deliberate deviation: the resulting
// OPRF scalar is cryptographically sound and self-consistent, but it is NOT the
// same scalar that the upstream Remnawave NestJS backend would derive from the
// same APP_SECRET. Vault data encrypted against one scalar cannot be decrypted
// with the other. This is acceptable as long as Exodus vaults are never migrated
// from upstream; if migration support is ever required, L must be changed to 48
// and all existing vault entries re-evaluated against the new scalar.
func expandMessageXMD_SHA512(msg, dst []byte, lenInBytes int) []byte {
	bLen := 64
	rLen := 128
	ell := (lenInBytes + bLen - 1) / bLen
	zPad := make([]byte, rLen)
	dstPrime := append(dst, byte(len(dst)))

	h := sha512.New()
	h.Write(zPad)
	h.Write(msg)
	h.Write([]byte{byte(lenInBytes >> 8), byte(lenInBytes)})
	h.Write([]byte{0})
	h.Write(dstPrime)
	b0 := h.Sum(nil)

	h = sha512.New()
	h.Write(b0)
	h.Write([]byte{1})
	h.Write(dstPrime)
	b1 := h.Sum(nil)

	res := make([]byte, 0, ell*bLen)
	res = append(res, b1...)
	return res[:lenInBytes]
}

func getOrDeriveVaultScalar(appSecret string) (*ristretto255.Scalar, error) {
	vaultScalarLock.RLock()
	if cachedScalar != nil && cachedSecret == appSecret {
		s := cachedScalar
		vaultScalarLock.RUnlock()
		return s, nil
	}
	vaultScalarLock.RUnlock()

	vaultScalarLock.Lock()
	defer vaultScalarLock.Unlock()

	if cachedScalar != nil && cachedSecret == appSecret {
		return cachedScalar, nil
	}

	kdf := hkdf.New(sha256.New, []byte(appSecret), nil, []byte("rw-vault-oprf-v1"))
	seed := make([]byte, 32)
	if _, err := io.ReadFull(kdf, seed); err != nil {
		return nil, err
	}

	dst := append([]byte("DeriveKeyPairOPRFV1-"), 0)
	dst = append(dst, []byte("-ristretto255-SHA512")...)

	info := []byte("rw-vault")
	infoEncoded := append([]byte{byte(len(info) >> 8), byte(len(info))}, info...)
	msg := append(seed, infoEncoded...)
	msg = append(msg, 0)

	order, _ := new(big.Int).SetString("1000000000000000000000000000000014def9dea2f79cd65812631a5cf5d3ed", 16)

	for counter := 0; counter <= 255; counter++ {
		msg[len(msg)-1] = byte(counter)
		xmd := expandMessageXMD_SHA512(msg, dst, 64)

		num := new(big.Int)
		reversed := make([]byte, 64)
		for i := 0; i < 64; i++ {
			reversed[i] = xmd[63-i]
		}
		num.SetBytes(reversed)
		num.Mod(num, order)

		if num.Sign() == 0 {
			continue
		}

		numBytes := num.Bytes()
		skBytes := make([]byte, 32)
		for i, b := range numBytes {
			skBytes[len(numBytes)-1-i] = b
		}

		// Use SetCanonicalBytes (replaces deprecated Scalar.Decode).
		scalar, err := new(ristretto255.Scalar).SetCanonicalBytes(skBytes)
		if err != nil {
			return nil, err
		}

		cachedSecret = appSecret
		cachedScalar = scalar
		return cachedScalar, nil
	}

	return nil, fmt.Errorf("failed to derive OPRF scalar")
}

// EvaluateBlindedElement performs OPRF scalar multiplication (RFC 9497 blindEvaluate).
func EvaluateBlindedElement(appSecret string, blindedBase64 string) (string, error) {
	blindedBytes, err := base64.StdEncoding.DecodeString(blindedBase64)
	if err != nil {
		return "", fmt.Errorf("invalid base64 blinded element: %w", err)
	}

	// Use SetCanonicalBytes (replaces deprecated Element.Decode).
	point, err := new(ristretto255.Element).SetCanonicalBytes(blindedBytes)
	if err != nil {
		return "", fmt.Errorf("invalid ristretto255 point: %w", err)
	}

	// #7: RFC 9497 §3.3.1 — reject the identity element.
	// Use NewIdentityElement (replaces deprecated new(Element).Zero()).
	identity := ristretto255.NewIdentityElement()
	if point.Equal(identity) == 1 {
		return "", fmt.Errorf("invalid blinded element: identity element is not allowed (RFC 9497 §3.3.1)")
	}

	scalar, err := getOrDeriveVaultScalar(appSecret)
	if err != nil {
		return "", fmt.Errorf("failed to get vault scalar: %w", err)
	}

	evaluated := new(ristretto255.Element).ScalarMult(scalar, point)

	return base64.StdEncoding.EncodeToString(evaluated.Bytes()), nil
}
