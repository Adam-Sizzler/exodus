package config

import (
	"crypto"
	"crypto/sha256"
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

	"golang.org/x/crypto/hkdf"
)

var ErrSecretKeyNotSet = errors.New("SECRET_KEY is not set")

type NodePayload struct {
	CaCertPem    string `json:"caCertPem"`
	JWTPublicKey string `json:"jwtPublicKey"`
	NodeCertPem  string `json:"nodeCertPem"`
	NodeKeyPem   string `json:"nodeKeyPem"`
}

// SecretKeyCheck represents a single validation check result.
type SecretKeyCheck struct {
	Name   string
	OK     bool
	Detail string
}

// SecretKeyReport contains all validation results for SECRET_KEY.
type SecretKeyReport struct {
	Checks []SecretKeyCheck
	AllOK  bool
	SNI    string
}

// ParseNodePayloadFromSecret parses and strictly validates the SECRET_KEY payload,
// returning both the payload and the structured validation report.
func ParseNodePayloadFromSecret() (NodePayload, *SecretKeyReport, error) {
	secret := strings.TrimSpace(os.Getenv("SECRET_KEY"))
	if secret == "" {
		return NodePayload{}, nil, ErrSecretKeyNotSet
	}
	secret = strings.Trim(secret, "\"'")

	raw, err := base64.StdEncoding.DecodeString(secret)
	if err != nil {
		report := &SecretKeyReport{
			Checks: []SecretKeyCheck{{Name: "SECRET_KEY base64", OK: false, Detail: err.Error()}},
			AllOK:  false,
		}
		return NodePayload{}, report, fmt.Errorf("SECRET_KEY is not valid base64 payload: %w", err)
	}

	var payload NodePayload
	if err := json.Unmarshal(raw, &payload); err != nil {
		report := &SecretKeyReport{
			Checks: []SecretKeyCheck{{Name: "SECRET_KEY json", OK: false, Detail: err.Error()}},
			AllOK:  false,
		}
		return NodePayload{}, report, fmt.Errorf("SECRET_KEY contains invalid JSON payload: %w", err)
	}

	payload.CaCertPem = normalizePEM(payload.CaCertPem)
	payload.JWTPublicKey = normalizePEM(payload.JWTPublicKey)
	payload.NodeCertPem = normalizePEM(payload.NodeCertPem)
	payload.NodeKeyPem = normalizePEM(payload.NodeKeyPem)

	if payload.CaCertPem == "" || payload.NodeCertPem == "" || payload.NodeKeyPem == "" {
		report := &SecretKeyReport{
			Checks: []SecretKeyCheck{{Name: "SECRET_KEY fields", OK: false, Detail: "missing required mTLS fields"}},
			AllOK:  false,
		}
		return NodePayload{}, report, fmt.Errorf("SECRET_KEY payload is missing required mTLS fields (caCertPem, nodeCertPem, nodeKeyPem)")
	}

	checks := make([]SecretKeyCheck, 0, 7)
	var caCert *x509.Certificate
	var nodeCert *x509.Certificate
	var firstError error

	// 1. CA parses
	func() {
		caBlock, _ := pem.Decode([]byte(payload.CaCertPem))
		if caBlock == nil {
			checks = append(checks, SecretKeyCheck{"CA parses", false, "invalid PEM"})
			if firstError == nil {
				firstError = fmt.Errorf("SECRET_KEY CA certificate is not a valid PEM block")
			}
			return
		}
		var err error
		caCert, err = x509.ParseCertificate(caBlock.Bytes)
		if err != nil {
			checks = append(checks, SecretKeyCheck{"CA parses", false, err.Error()})
			if firstError == nil {
				firstError = fmt.Errorf("SECRET_KEY CA certificate parsing failed: %w", err)
			}
			return
		}
		fp := fmt.Sprintf("%X", sha256.Sum256(caCert.Raw))
		checks = append(checks, SecretKeyCheck{"CA parses", true, formatFingerprint(fp)})
	}()

	// 2. CA not expired
	func() {
		if caCert == nil {
			checks = append(checks, SecretKeyCheck{"CA not expired", false, "CA unavailable"})
			return
		}
		now := time.Now()
		if now.Before(caCert.NotBefore) {
			checks = append(checks, SecretKeyCheck{"CA not expired", false, "not yet valid"})
			if firstError == nil {
				firstError = fmt.Errorf("SECRET_KEY CA certificate is not yet valid")
			}
			return
		}
		if now.After(caCert.NotAfter) {
			checks = append(checks, SecretKeyCheck{"CA not expired", false, "expired"})
			if firstError == nil {
				firstError = fmt.Errorf("SECRET_KEY CA certificate expired on %s", caCert.NotAfter.UTC().Format(time.RFC3339))
			}
			return
		}
		checks = append(checks, SecretKeyCheck{"CA not expired", true, "until " + caCert.NotAfter.UTC().Format("Jan 02 15:04:05 2006 GMT")})
	}()

	// 3. CA self-signature
	func() {
		if caCert == nil {
			checks = append(checks, SecretKeyCheck{"CA self-signature", false, "CA unavailable"})
			return
		}
		if err := caCert.CheckSignatureFrom(caCert); err != nil {
			checks = append(checks, SecretKeyCheck{"CA self-signature", false, "mismatch – corrupted"})
			if firstError == nil {
				firstError = fmt.Errorf("SECRET_KEY CA self-signature check failed: %w", err)
			}
			return
		}
		checks = append(checks, SecretKeyCheck{"CA self-signature", true, "valid"})
	}()

	// 4. node cert parses
	func() {
		nodeBlock, _ := pem.Decode([]byte(payload.NodeCertPem))
		if nodeBlock == nil {
			checks = append(checks, SecretKeyCheck{"node cert parses", false, "invalid PEM"})
			if firstError == nil {
				firstError = fmt.Errorf("SECRET_KEY node certificate is not a valid PEM block")
			}
			return
		}
		var err error
		nodeCert, err = x509.ParseCertificate(nodeBlock.Bytes)
		if err != nil {
			checks = append(checks, SecretKeyCheck{"node cert parses", false, err.Error()})
			if firstError == nil {
				firstError = fmt.Errorf("SECRET_KEY node certificate parsing failed: %w", err)
			}
			return
		}
		fp := fmt.Sprintf("%X", sha256.Sum256(nodeCert.Raw))
		checks = append(checks, SecretKeyCheck{"node cert parses", true, formatFingerprint(fp)})
	}()

	// 5. node signed by CA
	func() {
		if caCert == nil || nodeCert == nil {
			checks = append(checks, SecretKeyCheck{"node signed by CA", false, "cert unavailable"})
			return
		}
		if err := nodeCert.CheckSignatureFrom(caCert); err != nil {
			checks = append(checks, SecretKeyCheck{"node signed by CA", false, "not signed by this CA"})
			if firstError == nil {
				firstError = fmt.Errorf("SECRET_KEY node certificate verification against CA failed: %w", err)
			}
			return
		}
		checks = append(checks, SecretKeyCheck{"node signed by CA", true, "valid"})
	}()

	// 6. node key matches cert
	func() {
		if nodeCert == nil {
			checks = append(checks, SecretKeyCheck{"node key matches cert", false, "cert unavailable"})
			return
		}
		_, err := tls.X509KeyPair([]byte(payload.NodeCertPem), []byte(payload.NodeKeyPem))
		if err != nil {
			checks = append(checks, SecretKeyCheck{"node key matches cert", false, "key does not match cert"})
			if firstError == nil {
				firstError = fmt.Errorf("SECRET_KEY node certificate does not match private key: %w", err)
			}
			return
		}
		checks = append(checks, SecretKeyCheck{"node key matches cert", true, "valid"})
	}()

	// 7. jwt public key
	func() {
		if payload.JWTPublicKey == "" {
			checks = append(checks, SecretKeyCheck{"jwt public key", false, "missing"})
			return
		}
		jwtBlock, _ := pem.Decode([]byte(payload.JWTPublicKey))
		if jwtBlock == nil {
			checks = append(checks, SecretKeyCheck{"jwt public key", false, "invalid PEM"})
			if firstError == nil {
				firstError = fmt.Errorf("SECRET_KEY jwtPublicKey is not a valid PEM block")
			}
			return
		}
		if _, err := x509.ParsePKIXPublicKey(jwtBlock.Bytes); err != nil {
			if _, err2 := x509.ParsePKCS1PublicKey(jwtBlock.Bytes); err2 != nil {
				checks = append(checks, SecretKeyCheck{"jwt public key", false, err.Error()})
				if firstError == nil {
					firstError = fmt.Errorf("SECRET_KEY jwtPublicKey parsing failed: %w", err)
				}
				return
			}
		}
		checks = append(checks, SecretKeyCheck{"jwt public key", true, "ok"})
	}()

	allOK := firstError == nil
	sni := ""
	if allOK {
		sni = deriveSNI(payload.CaCertPem, payload.JWTPublicKey)
	}

	report := &SecretKeyReport{
		Checks: checks,
		AllOK:  allOK,
		SNI:    sni,
	}

	return payload, report, firstError
}

func normalizePEM(pem string) string {
	normalized := strings.ReplaceAll(pem, "\\n", "\n")
	normalized = strings.ReplaceAll(normalized, "\r\n", "\n")
	normalized = strings.TrimSpace(normalized)
	return normalized
}

func formatFingerprint(hex string) string {
	hex = strings.ReplaceAll(hex, ":", "")
	if len(hex) > 16 {
		hex = hex[:16]
	}
	parts := make([]string, 0, len(hex)/2)
	for i := 0; i+1 < len(hex); i += 2 {
		parts = append(parts, hex[i:i+2])
	}
	return strings.Join(parts, ":") + ":"
}

const hkdfInfo = "rw-v1"

var sniTLDs = []string{"com", "net", "org", "io", "dev", "app"}

func deriveSNI(caCertPem, jwtPublicKey string) string {
	canonicalize := func(pemStr string) string {
		lines := strings.Split(pemStr, "\n")
		var sb strings.Builder
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if strings.HasPrefix(line, "-----") {
				continue
			}
			sb.WriteString(line)
		}
		return sb.String()
	}

	ikm := []byte(canonicalize(jwtPublicKey) + canonicalize(caCertPem))

	reader := hkdf.New(crypto.SHA256.New, ikm, nil, []byte(hkdfInfo))
	okm := make([]byte, 22)
	if _, err := reader.Read(okm); err != nil {
		return "unknown"
	}

	host := fmt.Sprintf("%x", okm[:16])
	mid := fmt.Sprintf("%x", okm[16:21])
	tld := sniTLDs[okm[21]%byte(len(sniTLDs))]
	return host + "." + mid + "." + tld
}
