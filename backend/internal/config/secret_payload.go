package config

import (
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"errors"
	"os"
	"strings"
)

var ErrSecretKeyNotSet = errors.New("SUB_SECRET_KEY is not set")

type NodePayload struct {
	CaCertPem    string `json:"caCertPem"`
	JWTPublicKey string `json:"jwtPublicKey"`
	NodeCertPem  string `json:"nodeCertPem"`
	NodeKeyPem   string `json:"nodeKeyPem"`
}

func ParseNodePayloadFromSecret() (NodePayload, error) {
	secret := strings.TrimSpace(os.Getenv("SUB_SECRET_KEY"))
	if secret == "" {
		return NodePayload{}, ErrSecretKeyNotSet
	}

	secret = strings.Trim(secret, "\"'")
	raw, err := base64.StdEncoding.DecodeString(secret)
	if err != nil {
		return NodePayload{}, NewEnvError("SUB_SECRET_KEY", "Invalid Secret Key payload. Dashboard → Subscription → Nodes → Current node → Secret Key (SUB_SECRET_KEY).")
	}

	var payload NodePayload
	if err := json.Unmarshal(raw, &payload); err != nil {
		return NodePayload{}, NewEnvError("SUB_SECRET_KEY", "Contains invalid JSON payload. Copy Secret Key from Dashboard → Subscription → Nodes → Current node.")
	}

	payload.CaCertPem = normalizePEM(payload.CaCertPem)
	payload.JWTPublicKey = normalizePEM(payload.JWTPublicKey)
	payload.NodeCertPem = normalizePEM(payload.NodeCertPem)
	payload.NodeKeyPem = normalizePEM(payload.NodeKeyPem)

	if payload.CaCertPem == "" || payload.NodeCertPem == "" || payload.NodeKeyPem == "" {
		return NodePayload{}, NewEnvError("SUB_SECRET_KEY", "Payload is missing mTLS fields. Copy Secret Key from Dashboard → Subscription → Nodes → Current node.")
	}

	if _, err := tls.X509KeyPair([]byte(payload.NodeCertPem), []byte(payload.NodeKeyPem)); err != nil {
		return NodePayload{}, NewEnvError("SUB_SECRET_KEY", "Payload has invalid node certificate/key pair. Regenerate Secret Key in Dashboard → Subscription → Nodes → Current node.")
	}

	return payload, nil
}

func normalizePEM(pem string) string {
	normalized := strings.ReplaceAll(pem, "\\n", "\n")
	normalized = strings.ReplaceAll(normalized, "\r\n", "\n")
	normalized = strings.TrimSpace(normalized)
	return normalized
}
