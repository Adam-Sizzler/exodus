package auth

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"exodus/internal/config"
	dbmanager "exodus/internal/db/manager"
	"exodus/internal/httpapi/middleware"
	"exodus/internal/security"
)

func resolveToken(token string, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) (*AuthPrincipal, error) {
	if token == "" {
		return nil, errors.New("empty token")
	}

	if principal, err := resolveAdminJWT(token, manager, cfg); err == nil {
		return principal, nil
	}

	if principal, err := resolveAPIJWT(token, manager, cfg); err == nil {
		return principal, nil
	}

	return nil, errors.New("invalid auth token")
}

func resolveAdminJWT(token string, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) (*AuthPrincipal, error) {
	payload, err := security.ParseJWT(cfg.JWT.AuthSecret, token)
	if err != nil {
		return nil, err
	}
	if !strings.EqualFold(payload.Role, "ADMIN") {
		return nil, errors.New("jwt role is not admin")
	}
	if payload.Username == nil || strings.TrimSpace(*payload.Username) == "" {
		return nil, errors.New("jwt username is empty")
	}

	username := strings.TrimSpace(*payload.Username)
	var principal *AuthPrincipal
	err = manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		row := db.QueryRow(`
			SELECT uuid, username, role
			FROM admin
			WHERE username = ? AND UPPER(role) = 'ADMIN'
			LIMIT 1
		`, username)

		var adminUUID, dbUsername, role string
		if scanErr := row.Scan(&adminUUID, &dbUsername, &role); scanErr != nil {
			return scanErr
		}
		if adminUUID != payload.UUID {
			return errors.New("jwt uuid does not match admin")
		}

		expiresAt := int64(0)
		if payload.ExpiresAt != nil {
			expiresAt = payload.ExpiresAt.Time.Unix()
		}
		principal = &AuthPrincipal{
			AdminUUID: adminUUID,
			Username:  dbUsername,
			Role:      strings.ToUpper(role),
			TokenType: "jwt_auth",
			ExpiresAt: expiresAt,
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	if principal == nil {
		return nil, errors.New("admin not found")
	}
	return principal, nil
}

func resolveAPIJWT(token string, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) (*AuthPrincipal, error) {
	payload, err := security.ParseJWT(cfg.JWT.APITokensSecret, token)
	if err != nil {
		return nil, err
	}
	if !strings.EqualFold(payload.Role, "API") {
		return nil, errors.New("jwt role is not api")
	}

	var principal *AuthPrincipal
	err = manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		row := db.QueryRow(`
			SELECT uuid, name, expire_at, array_to_json(COALESCE(scopes, ARRAY['*']::text[]))::text AS scopes
			FROM api_tokens
			WHERE uuid = ?
			LIMIT 1
		`, payload.UUID)

		var tokenUUID, tokenName string
		var tokenExpireAt time.Time
		var scopesRaw string
		if scanErr := row.Scan(&tokenUUID, &tokenName, &tokenExpireAt, &scopesRaw); scanErr != nil {
			return scanErr
		}
		if time.Now().After(tokenExpireAt) {
			return errors.New("api token expired")
		}

		expiresAt := int64(0)
		if payload.ExpiresAt != nil {
			expiresAt = payload.ExpiresAt.Time.Unix()
		}
		principal = &AuthPrincipal{
			AdminUUID: tokenUUID,
			Username:  tokenName,
			Role:      "API",
			TokenType: "jwt_api_token",
			ExpiresAt: expiresAt,
			Scopes:    parseAPITokenPrincipalScopes(scopesRaw),
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	if principal == nil {
		return nil, errors.New("api token not found")
	}
	return principal, nil
}

func parseAPITokenPrincipalScopes(raw string) []string {
	var scopes []string
	if err := json.Unmarshal([]byte(strings.TrimSpace(raw)), &scopes); err != nil {
		return []string{"*"}
	}
	return normalizeAPITokenPrincipalScopes(scopes)
}

func createAdminAccessToken(cfg *config.BackendConfig, username, adminUUID, role string) (string, int64, error) {
	if strings.TrimSpace(role) == "" {
		role = "ADMIN"
	}
	return security.SignAuthJWT(cfg.JWT.AuthSecret, username, adminUUID, role)
}

func setAuthCookie(w http.ResponseWriter, r *http.Request, cfg *config.BackendConfig, token string, expiresAt int64) {
	secureCookie := middleware.IsSecureRequest(r, cfg)
	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieName,
		Value:    token,
		Path:     "/",
		Expires:  time.Unix(expiresAt, 0).UTC(),
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   secureCookie,
	})
}
