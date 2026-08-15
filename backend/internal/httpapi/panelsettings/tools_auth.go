package panelsettings

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"exodus/internal/config"
	"exodus/internal/httpapi/middleware"
	"exodus/internal/httpapi/shared"
	"exodus/internal/security"

	"github.com/golang-jwt/jwt/v5"
)

const (
	BackendToolsAuthCookieName = "ex-tools"
	BackendToolsJWTIssuer      = "exodus"
	BackendToolsJWTScopeAccess = "access"
	BackendToolsJWTScopeOtt    = "ott"
	BackendToolsJWTLifetime = 2 * time.Hour
)

type ToolsClaims struct {
	Scope string `json:"scope"`
	jwt.RegisteredClaims
}

func ToolsAuthMiddleware(cfg *config.BackendConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			secret := cfg.JWT.AuthSecret
			if err := security.ValidateJWTSecret(secret); err != nil {
				shared.WriteJSONError(w, http.StatusForbidden, "forbidden")
				return
			}

			ott := r.URL.Query().Get("ott")
			if ott != "" {
				if verifyToolsJWT(ott, secret, BackendToolsJWTScopeOtt) {
					accessToken, err := signToolsAccessJWT(secret)
					if err == nil {
						cookiePath := cfg.Panel.Trimmed() + "/api/backend-tools"

						http.SetCookie(w, &http.Cookie{
							Name:     BackendToolsAuthCookieName,
							Value:    accessToken,
							Path:     cookiePath,
							HttpOnly: true,
							Secure:   middleware.IsSecureRequest(r, cfg),
							SameSite: http.SameSiteLaxMode,
							MaxAge:   int(BackendToolsJWTLifetime.Seconds()),
						})
					}

					basePath := ""
					if cfg != nil {
						basePath = cfg.Panel.Trimmed()
					}

					redirectPath := r.URL.Path
					if basePath != "" && !strings.HasPrefix(redirectPath, basePath) {
						redirectPath = basePath + redirectPath
					}

					q := r.URL.Query()
					q.Del("ott")
					if len(q) > 0 {
						redirectPath += "?" + q.Encode()
					}

					http.Redirect(w, r, redirectPath, http.StatusFound)
					return
				}

				shared.WriteJSONError(w, http.StatusForbidden, "invalid OTT token")
				return
			}

			// 2. Check for backend tools cookie (ex-tools)
			if cookie, err := r.Cookie(BackendToolsAuthCookieName); err == nil && cookie.Value != "" {
				if verifyToolsJWT(cookie.Value, secret, BackendToolsJWTScopeAccess) {
					next.ServeHTTP(w, r)
					return
				}
			}

			shared.WriteJSONError(w, http.StatusForbidden, "forbidden")
		})
	}
}

func signToolsAccessJWT(secret string) (string, error) {
	now := time.Now().UTC()
	claims := ToolsClaims{
		Scope: BackendToolsJWTScopeAccess,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    BackendToolsJWTIssuer,
			ExpiresAt: jwt.NewNumericDate(now.Add(BackendToolsJWTLifetime)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func verifyToolsJWT(rawToken, secret, expectedScope string) bool {
	token, err := jwt.ParseWithClaims(rawToken, &ToolsClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected jwt signing method: %v", t.Header["alg"])
		}
		return []byte(secret), nil
	}, jwt.WithIssuer(BackendToolsJWTIssuer))
	if err != nil || !token.Valid {
		return false
	}
	claims, ok := token.Claims.(*ToolsClaims)
	if !ok {
		return false
	}
	return claims.Scope == expectedScope
}
