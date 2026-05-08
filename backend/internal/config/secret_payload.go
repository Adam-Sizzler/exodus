package config

import (
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
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
	secret, ok := lookupEnvTrimmed("SUB_SECRET_KEY")
	if !ok {
		return NodePayload{}, ErrSecretKeyNotSet
	}

	secret = strings.Trim(secret, "\"'")
	raw, err := base64.StdEncoding.DecodeString(secret)
	if err != nil {
		return NodePayload{}, fmt.Errorf("SUB_SECRET_KEY is not valid base64 payload: %w", err)
	}

	var payload NodePayload
	if err := json.Unmarshal(raw, &payload); err != nil {
		return NodePayload{}, fmt.Errorf("SUB_SECRET_KEY contains invalid JSON payload: %w", err)
	}

	payload.CaCertPem = normalizePEM(payload.CaCertPem)
	payload.JWTPublicKey = normalizePEM(payload.JWTPublicKey)
	payload.NodeCertPem = normalizePEM(payload.NodeCertPem)
	payload.NodeKeyPem = normalizePEM(payload.NodeKeyPem)

	if payload.CaCertPem == "" || payload.NodeCertPem == "" || payload.NodeKeyPem == "" {
		return NodePayload{}, fmt.Errorf("SUB_SECRET_KEY payload is missing mTLS fields")
	}

	if _, err := tls.X509KeyPair([]byte(payload.NodeCertPem), []byte(payload.NodeKeyPem)); err != nil {
		return NodePayload{}, fmt.Errorf("SUB_SECRET_KEY payload has invalid node certificate/key pair: %w", err)
	}

	return payload, nil
}

func normalizePEM(pem string) string {
	normalized := strings.ReplaceAll(pem, "\\n", "\n")
	normalized = strings.ReplaceAll(normalized, "\r\n", "\n")
	normalized = strings.TrimSpace(normalized)
	return normalized
}
