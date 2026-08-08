package security

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	// Exodus-compatible admin access token lifetime.
	AuthTokenLifetime = 12 * time.Hour
	APITokenLifetime  = 99999 * 24 * time.Hour
)

type JWTPayload struct {
	Username *string `json:"username"`
	UUID     string  `json:"uuid"`
	Role     string  `json:"role"`
	jwt.RegisteredClaims
}

func SignAuthJWT(secret, username, uuid, role string) (string, int64, error) {
	return SignAuthJWTWithLifetime(secret, username, uuid, role, AuthTokenLifetime)
}

func SignAuthJWTWithLifetime(secret, username, uuid, role string, lifetime time.Duration) (string, int64, error) {
	if lifetime <= 0 {
		lifetime = AuthTokenLifetime
	}
	username = strings.TrimSpace(username)
	if username == "" {
		return "", 0, errors.New("username is required")
	}
	usernamePtr := username
	return signJWT(secret, &usernamePtr, uuid, role, lifetime)
}

func SignAPITokenJWT(secret, uuid string) (string, int64, error) {
	return signJWT(secret, nil, uuid, "API", APITokenLifetime)
}

func SignAPITokenJWTWithLifetime(secret, uuid string, lifetime time.Duration) (string, int64, error) {
	if lifetime <= 0 {
		lifetime = APITokenLifetime
	}
	return signJWT(secret, nil, uuid, "API", lifetime)
}

type OttJWTPayload struct {
	Scope string `json:"scope"`
	jwt.RegisteredClaims
}

func SignOttJWT(secret string) (string, error) {
	if err := validateJWTSecret(secret); err != nil {
		return "", err
	}
	now := time.Now().UTC()
	claims := OttJWTPayload{
		Scope: "ott",
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "exodus",
			ExpiresAt: jwt.NewNumericDate(now.Add(30 * time.Second)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func ParseJWT(secret, rawToken string) (*JWTPayload, error) {
	if err := validateJWTSecret(secret); err != nil {
		return nil, err
	}
	rawToken = strings.TrimSpace(rawToken)
	if rawToken == "" {
		return nil, errors.New("token is required")
	}

	claims := &JWTPayload{}
	token, err := jwt.ParseWithClaims(rawToken, claims, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected jwt signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	if token == nil || !token.Valid {
		return nil, errors.New("invalid jwt")
	}

	claims.UUID = strings.TrimSpace(claims.UUID)
	claims.Role = strings.ToUpper(strings.TrimSpace(claims.Role))
	if claims.UUID == "" || claims.Role == "" {
		return nil, errors.New("jwt payload is missing uuid or role")
	}
	return claims, nil
}

func signJWT(secret string, username *string, uuid, role string, lifetime time.Duration) (string, int64, error) {
	if err := validateJWTSecret(secret); err != nil {
		return "", 0, err
	}
	uuid = strings.TrimSpace(uuid)
	role = strings.ToUpper(strings.TrimSpace(role))
	if uuid == "" || role == "" {
		return "", 0, errors.New("uuid and role are required")
	}

	now := time.Now().UTC()
	expiresAt := now.Add(lifetime)
	claims := JWTPayload{
		Username: username,
		UUID:     uuid,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}

	signed, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secret))
	if err != nil {
		return "", 0, err
	}
	return signed, expiresAt.Unix(), nil
}

func validateJWTSecret(secret string) error {
	secret = strings.TrimSpace(secret)
	if secret == "" {
		return errors.New("jwt secret is not configured")
	}
	if secret == "change_me" {
		return errors.New("jwt secret cannot be change_me")
	}
	return nil
}
