package auth

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"exodus/internal/config"
	dbmanager "exodus/internal/db/manager"
	"exodus/internal/httpapi/shared"
)

func WithPanelAuth(manager *dbmanager.DatabaseManager, cfg *config.BackendConfig, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = cleanPath(r.URL.Path)

		// Protect API routes only.
		if !strings.HasPrefix(r.URL.Path, "/api/") {
			next.ServeHTTP(w, r)
			return
		}

		if isPublicAPIPath(r.URL.Path) {
			next.ServeHTTP(w, r)
			return
		}

		principal, err := authenticateRequest(r, manager, cfg)
		if err != nil {
			cfg.Logger.Warn("Auth rejected request",
				"path", r.URL.Path,
				"method", r.Method,
				"remote_addr", r.RemoteAddr,
				"error", err,
			)
			shared.WriteJSONError(w, http.StatusUnauthorized, "unauthorized")
			return
		}
		if !requireAPITokenScope(principal, r) {
			shared.WriteJSONError(w, http.StatusForbidden, "api token scope is not allowed for this endpoint")
			return
		}

		ctx := context.WithValue(r.Context(), authPrincipalContextKey, principal)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func RequireAdminRole(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		principal, ok := CurrentAuthPrincipal(r.Context())
		if !ok {
			shared.WriteJSONError(w, http.StatusUnauthorized, "unauthorized")
			return
		}
		if !strings.EqualFold(principal.Role, "ADMIN") {
			shared.WriteJSONError(w, http.StatusForbidden, "forbidden: admin role required")
			return
		}
		next(w, r)
	}
}

func authenticateRequest(r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) (*AuthPrincipal, error) {
	if token := readBearerToken(r.Header.Get("Authorization")); token != "" {
		return resolveToken(token, manager, cfg)
	}

	cookie, err := r.Cookie(sessionCookieName)
	if err == nil && strings.TrimSpace(cookie.Value) != "" {
		return resolveToken(strings.TrimSpace(cookie.Value), manager, cfg)
	}

	return nil, errors.New("missing auth token")
}

func isPublicAPIPath(p string) bool {
	p = cleanPath(p)
	if !strings.HasPrefix(p, "/api/") {
		return true
	}
	trimmed := p[len("/api/"):]

	publicSuffixes := []string{
		"auth/bootstrap", "auth/status", "auth/setup", "auth/login", "auth/register",
		"auth/passkey/authentication/options", "auth/passkey/authentication/verify",
		"auth/oauth2/authorize", "auth/oauth2/callback",
	}

	if strings.HasPrefix(trimmed, "queues/static") {
		return true
	}

	for _, s := range publicSuffixes {
		if trimmed == s || trimmed == s+"/" {
			return true
		}
	}

	// Legacy backward compatibility paths
	if trimmed == "login" || trimmed == "login/" || trimmed == "register" || trimmed == "register/" {
		return true
	}

	return false
}
