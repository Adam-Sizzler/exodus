package passkeys

import (
	"encoding/json"
	"time"

	gowebauthn "github.com/go-webauthn/webauthn/webauthn"
)

type passkeyRecord struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	CreatedAt  time.Time `json:"createdAt"`
	LastUsedAt time.Time `json:"lastUsedAt"`
}

type passkeyWriteRequest struct {
	ID   string `json:"id"`
	Name string `json:"name,omitempty"`
}

type passkeySettings struct {
	Enabled bool    `json:"enabled"`
	RPID    *string `json:"rpId"`
	Origin  *string `json:"origin"`
}

type resolvedPasskeySettings struct {
	RPID    string
	Origins []string
}

type webAuthnAdmin struct {
	uuid        string
	username    string
	credentials []gowebauthn.Credential
}

func (a *webAuthnAdmin) WebAuthnID() []byte {
	return []byte(a.uuid)
}

func (a *webAuthnAdmin) WebAuthnName() string {
	return a.username
}

func (a *webAuthnAdmin) WebAuthnDisplayName() string {
	return "Exodus Administrator"
}

func (a *webAuthnAdmin) WebAuthnCredentials() []gowebauthn.Credential {
	return a.credentials
}

type storedPasskeySession struct {
	session gowebauthn.SessionData
	expires time.Time
}

type verifyRegistrationRequest struct {
	Response json.RawMessage `json:"response"`
}

type verifyAuthenticationRequest struct {
	Response json.RawMessage `json:"response"`
}
