package config

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"
)

var ErrSecretKeyNotSet = errors.New("SECRET_KEY is not set")

type NodePayload struct {
	CaCertPem    string `json:"caCertPem"`
	JWTPublicKey string `json:"jwtPublicKey"`
	NodeCertPem  string `json:"nodeCertPem"`
	NodeKeyPem   string `json:"nodeKeyPem"`
}

// ParseNodePayloadFromSecret parses and strictly validates the SECRET_KEY payload.
func ParseNodePayloadFromSecret() (NodePayload, error) {
	secret := strings.TrimSpace(os.Getenv("SECRET_KEY"))
	if secret == "" {
		return NodePayload{}, ErrSecretKeyNotSet
	}
	secret = strings.Trim(secret, "\"'")

	raw, err := base64.StdEncoding.DecodeString(secret)
	if err != nil {
		return NodePayload{}, fmt.Errorf("SECRET_KEY is not valid base64 payload: %w", err)
	}

	var payload NodePayload
	if err := json.Unmarshal(raw, &payload); err != nil {
		return NodePayload{}, fmt.Errorf("SECRET_KEY contains invalid JSON payload: %w", err)
	}

	payload.CaCertPem = normalizePEM(payload.CaCertPem)
	payload.JWTPublicKey = normalizePEM(payload.JWTPublicKey)
	payload.NodeCertPem = normalizePEM(payload.NodeCertPem)
	payload.NodeKeyPem = normalizePEM(payload.NodeKeyPem)

	if payload.CaCertPem == "" || payload.NodeCertPem == "" || payload.NodeKeyPem == "" {
		return NodePayload{}, fmt.Errorf("SECRET_KEY payload is missing required mTLS fields (caCertPem, nodeCertPem, nodeKeyPem)")
	}

	// 1. Strict keypair match validation
	if _, err := tls.X509KeyPair([]byte(payload.NodeCertPem), []byte(payload.NodeKeyPem)); err != nil {
		return NodePayload{}, fmt.Errorf("SECRET_KEY node certificate does not match private key: %w", err)
	}

	// 2. Strict CA certificate parsing
	caBlock, _ := pem.Decode([]byte(payload.CaCertPem))
	if caBlock == nil {
		return NodePayload{}, fmt.Errorf("SECRET_KEY CA certificate is not a valid PEM block")
	}
	caCert, err := x509.ParseCertificate(caBlock.Bytes)
	if err != nil {
		return NodePayload{}, fmt.Errorf("SECRET_KEY CA certificate parsing failed: %w", err)
	}

	// 3. Strict Node certificate parsing and expiration check
	nodeBlock, _ := pem.Decode([]byte(payload.NodeCertPem))
	if nodeBlock == nil {
		return NodePayload{}, fmt.Errorf("SECRET_KEY node certificate is not a valid PEM block")
	}
	nodeCert, err := x509.ParseCertificate(nodeBlock.Bytes)
	if err != nil {
		return NodePayload{}, fmt.Errorf("SECRET_KEY node certificate parsing failed: %w", err)
	}

	now := time.Now()
	if now.After(nodeCert.NotAfter) {
		return NodePayload{}, fmt.Errorf("SECRET_KEY node certificate has expired on %s", nodeCert.NotAfter.UTC().Format(time.RFC3339))
	}
	if now.Before(nodeCert.NotBefore) {
		return NodePayload{}, fmt.Errorf("SECRET_KEY node certificate is not yet valid (valid from %s)", nodeCert.NotBefore.UTC().Format(time.RFC3339))
	}

	// 4. Verify certificate chain (node cert signed by CA cert)
	caPool := x509.NewCertPool()
	caPool.AddCert(caCert)
	verifyOpts := x509.VerifyOptions{
		Roots:       caPool,
		CurrentTime: now,
		KeyUsages:   []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageAny},
	}
	if _, err := nodeCert.Verify(verifyOpts); err != nil {
		// Also allow without strict KeyUsages constraint if custom cert
		basicOpts := x509.VerifyOptions{
			Roots:       caPool,
			CurrentTime: now,
		}
		if _, bErr := nodeCert.Verify(basicOpts); bErr != nil {
			return NodePayload{}, fmt.Errorf("SECRET_KEY node certificate verification against CA failed: %w", bErr)
		}
	}

	// 5. Strict JWT Public Key validation (if present)
	if payload.JWTPublicKey != "" {
		jwtBlock, _ := pem.Decode([]byte(payload.JWTPublicKey))
		if jwtBlock == nil {
			return NodePayload{}, fmt.Errorf("SECRET_KEY jwtPublicKey is not a valid PEM block")
		}
		if _, err := x509.ParsePKIXPublicKey(jwtBlock.Bytes); err != nil {
			if _, err2 := x509.ParsePKCS1PublicKey(jwtBlock.Bytes); err2 != nil {
				return NodePayload{}, fmt.Errorf("SECRET_KEY jwtPublicKey parsing failed: %w", err)
			}
		}
	}

	return payload, nil
}

func normalizePEM(pem string) string {
	normalized := strings.ReplaceAll(pem, "\\n", "\n")
	normalized = strings.ReplaceAll(normalized, "\r\n", "\n")
	normalized = strings.TrimSpace(normalized)
	return normalized
}
