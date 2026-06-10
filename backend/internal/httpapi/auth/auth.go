package auth

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"exodus/internal/config"
	dbmanager "exodus/internal/db/manager"
	"exodus/internal/httpapi/middleware"
	"exodus/internal/httpapi/shared"
	"exodus/internal/notifications"
	"exodus/internal/security"

	"github.com/google/uuid"
)

const sessionCookieName = "exodus_session"

type authContextKey string

const authPrincipalContextKey authContextKey = "auth_principal"

type AuthPrincipal struct {
	AdminUUID string `json:"admin_uuid,omitempty"`
	Username  string `json:"username,omitempty"`
	Role      string `json:"role"`
	TokenType string `json:"token_type"`
	ExpiresAt int64  `json:"expires_at,omitempty"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type SetupRequest struct {
	Username          string `json:"username"`
	Password          string `json:"password"`
	SessionTTLMinutes int    `json:"session_ttl_minutes,omitempty"`
}

type AuthAdminInfo struct {
	UUID              string `json:"uuid"`
	Username          string `json:"username"`
	Role              string `json:"role"`
	SessionTTLMinutes int    `json:"session_ttl_minutes"`
}

type LoginResponse struct {
	Admin            AuthAdminInfo     `json:"admin"`
	ExpiresAt        int64             `json:"expires_at"`
	BrandingSettings map[string]any    `json:"branding_settings"`
	PasswordSettings map[string]any    `json:"password_settings"`
	Message          string            `json:"message,omitempty"`
	Principal        *AuthPrincipal    `json:"principal,omitempty"`
	Meta             map[string]string `json:"meta,omitempty"`
}

type BootstrapResponse struct {
	BrandingSettings   map[string]any `json:"branding_settings"`
	PasswordSettings   map[string]any `json:"password_settings"`
	DefaultUsername    string         `json:"default_username"`
	HasAdminConfigured bool           `json:"has_admin_configured"`
}

var errAdminAlreadyConfigured = errors.New("admin already configured")

func CurrentAuthPrincipal(ctx context.Context) (*AuthPrincipal, bool) {
	value := ctx.Value(authPrincipalContextKey)
	if value == nil {
		return nil, false
	}
	principal, ok := value.(*AuthPrincipal)
	return principal, ok && principal != nil
}

func emitLoginNotification(ctx context.Context, cfg *config.BackendConfig, event, username, adminUUID, password, reason string, r *http.Request) {
	loginAttempt := map[string]any{
		"username":    username,
		"password":    password,
		"ip":          notificationClientIP(r),
		"userAgent":   "",
		"description": reason,
	}
	data := map[string]any{
		"username":     username,
		"password":     password,
		"ip":           loginAttempt["ip"],
		"description":  reason,
		"loginAttempt": loginAttempt,
	}
	if adminUUID != "" {
		data["adminUuid"] = adminUUID
		loginAttempt["adminUuid"] = adminUUID
	}
	if r != nil {
		data["remoteAddr"] = r.RemoteAddr
		data["userAgent"] = r.UserAgent()
		data["path"] = r.URL.Path
		loginAttempt["remoteAddr"] = r.RemoteAddr
		loginAttempt["userAgent"] = r.UserAgent()
		loginAttempt["path"] = r.URL.Path
	}
	notifications.Emit(ctx, cfg, notifications.Event{
		Scope: notifications.ScopeService,
		Event: event,
		Data:  data,
	})
}

func notificationClientIP(r *http.Request) string {
	if r == nil {
		return ""
	}
	for _, header := range []string{"CF-Connecting-IP", "X-Real-IP", "X-Forwarded-For"} {
		value := strings.TrimSpace(r.Header.Get(header))
		if value == "" {
			continue
		}
		if comma := strings.Index(value, ","); comma >= 0 {
			value = strings.TrimSpace(value[:comma])
		}
		if value != "" {
			return value
		}
	}
	remoteAddr := strings.TrimSpace(r.RemoteAddr)
	if colon := strings.LastIndex(remoteAddr, ":"); colon > 0 && !strings.Contains(remoteAddr[colon+1:], "]") {
		return strings.Trim(remoteAddr[:colon], "[]")
	}
	return strings.Trim(remoteAddr, "[]")
}

func WithPanelAuth(manager *dbmanager.DatabaseManager, cfg *config.BackendConfig, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

		ctx := context.WithValue(r.Context(), authPrincipalContextKey, principal)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func isPublicAPIPath(path string) bool {
	switch path {
	case "/api/health",
		"/api/auth/status", "/api/auth/login", "/api/auth/register",
		"/api/auth/oauth2/authorize", "/api/auth/oauth2/callback", "/api/auth/oauth2/tg/callback",
		"/api/auth/passkey/authentication/options", "/api/auth/passkey/authentication/verify",
		"/api/system/metadata":
		return true
	default:
		if strings.HasPrefix(path, "/api/sub/") || path == "/api/sub" {
			return true
		}
		return false
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

func readBearerToken(authHeader string) string {
	authHeader = strings.TrimSpace(authHeader)
	if authHeader == "" {
		return ""
	}
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return ""
	}
	return strings.TrimSpace(parts[1])
}

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
			SELECT uuid, token_name
			FROM api_tokens
			WHERE uuid = ?
			LIMIT 1
		`, payload.UUID)

		var tokenUUID, tokenName string
		if scanErr := row.Scan(&tokenUUID, &tokenName); scanErr != nil {
			return scanErr
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

func AuthTTLMinutes() int {
	return int(security.AuthTokenLifetime / time.Minute)
}

func AuthBootstrapHandler(manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			shared.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}

		brandingSettings, passwordSettings, defaultUsername, hasAdmin, err := getBootstrapData(manager)
		if err != nil {
			cfg.Logger.Error("Failed to load auth bootstrap data", "error", err)
			shared.WriteJSONError(w, http.StatusInternalServerError, "failed to load auth bootstrap data")
			return
		}

		shared.WriteJSON(w, http.StatusOK, BootstrapResponse{
			BrandingSettings:   brandingSettings,
			PasswordSettings:   passwordSettings,
			DefaultUsername:    defaultUsername,
			HasAdminConfigured: hasAdmin,
		})
	}
}

func AuthStatusHandler(manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			shared.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}

		brandingSettings, passwordSettings, _, hasAdmin, err := getBootstrapData(manager)
		if err != nil {
			cfg.Logger.Error("Failed to load auth status", "error", err)
			shared.WriteJSONError(w, http.StatusInternalServerError, "failed to load auth status")
			return
		}
		cfg.Logger.Trace("Auth status requested", "has_admin", hasAdmin, "remote_addr", r.RemoteAddr)

		title, _ := brandingSettings["title"].(string)
		logoURL, _ := brandingSettings["logoUrl"].(string)
		if hasAdmin {
			passkeyEnabled, tgEnabled, tgBotID, oauth2Providers := getAuthMethodsStatus(manager)
			passwordEnabled := resolvePasswordAuthEnabled(
				passwordSettings,
				passkeyEnabled,
				tgEnabled,
				oauth2Providers,
				cfg,
			)

			shared.WriteJSON(w, http.StatusOK, map[string]any{
				"response": map[string]any{
					"isLoginAllowed":    true,
					"isRegisterAllowed": false,
					"authentication": map[string]any{
						"passkey": map[string]any{
							"enabled": passkeyEnabled,
						},
						"tgAuth": map[string]any{
							"enabled": tgEnabled,
							"botId":   tgBotID,
						},
						"oauth2": map[string]any{
							"providers": oauth2Providers,
						},
						"password": map[string]any{
							"enabled": passwordEnabled,
						},
					},
					"branding": map[string]any{
						"title":   nullableString(title),
						"logoUrl": nullableString(logoURL),
					},
				},
			})
			return
		}

		shared.WriteJSON(w, http.StatusOK, map[string]any{
			"response": map[string]any{
				"isLoginAllowed":    false,
				"isRegisterAllowed": true,
				"authentication":    nil,
				"branding": map[string]any{
					"title":   nullableString(title),
					"logoUrl": nullableString(logoURL),
				},
			},
		})
	}
}

func AuthLoginCompatHandler(manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			shared.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}

		var req LoginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			cfg.Logger.Warn("Auth login decode failed", "error", err, "remote_addr", r.RemoteAddr)
			shared.WriteJSONError(w, http.StatusBadRequest, "invalid JSON body")
			return
		}

		username := strings.TrimSpace(req.Username)
		password := strings.TrimSpace(req.Password)
		if username == "" || password == "" {
			shared.WriteJSONError(w, http.StatusBadRequest, "username and password are required")
			return
		}
		cfg.Logger.Trace("Auth login attempt", "username", username, "remote_addr", r.RemoteAddr)

		_, passwordSettings, _, hasAdmin, err := getBootstrapData(manager)
		if err != nil {
			cfg.Logger.Error("Failed to read auth bootstrap for login", "error", err)
			shared.WriteJSONError(w, http.StatusInternalServerError, "failed to validate credentials")
			return
		}

		if !hasAdmin {
			cfg.Logger.Warn("Auth login blocked: no admin configured", "username", username, "remote_addr", r.RemoteAddr)
			emitLoginNotification(r.Context(), cfg, notifications.EventLoginAttemptFailed, username, "", password, "no_admin_configured", r)
			shared.WriteJSONError(w, http.StatusForbidden, "login is not allowed")
			return
		}

		passkeyEnabled, tgEnabled, _, oauth2Providers := getAuthMethodsStatus(manager)
		passwordEnabled := resolvePasswordAuthEnabled(
			passwordSettings,
			passkeyEnabled,
			tgEnabled,
			oauth2Providers,
			cfg,
		)
		if !passwordEnabled {
			cfg.Logger.Warn("Auth login blocked: password auth disabled", "username", username, "remote_addr", r.RemoteAddr)
			emitLoginNotification(r.Context(), cfg, notifications.EventLoginAttemptFailed, username, "", password, "password_auth_disabled", r)
			shared.WriteJSONError(w, http.StatusForbidden, "login is not allowed")
			return
		}

		var (
			adminUUID          string
			storedPasswordHash string
			role               string
		)

		err = manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
			row := db.QueryRow(`
				SELECT uuid, password_hash, role
				FROM admin
				WHERE username = ?
				LIMIT 1
			`, username)

			if scanErr := row.Scan(&adminUUID, &storedPasswordHash, &role); scanErr != nil {
				if errors.Is(scanErr, sql.ErrNoRows) {
					return nil
				}
				return scanErr
			}
			return nil
		})
		if err != nil {
			cfg.Logger.Error("Failed to read admin credentials", "error", err)
			shared.WriteJSONError(w, http.StatusInternalServerError, "failed to validate credentials")
			return
		}
		if adminUUID == "" || !security.VerifyPassword(password, storedPasswordHash) {
			cfg.Logger.Warn("Auth login failed: invalid credentials", "username", username, "remote_addr", r.RemoteAddr)
			emitLoginNotification(r.Context(), cfg, notifications.EventLoginAttemptFailed, username, adminUUID, password, "invalid_credentials", r)
			shared.WriteJSONError(w, http.StatusForbidden, "invalid username or password")
			return
		}

		accessToken, expiresAt, err := createAdminAccessToken(cfg, username, adminUUID, role)
		if err != nil {
			cfg.Logger.Error("Failed to create JWT auth token", "error", err)
			shared.WriteJSONError(w, http.StatusInternalServerError, "failed to create auth token")
			return
		}
		cfg.Logger.Info("Auth login success", "username", username, "remote_addr", r.RemoteAddr)
		emitLoginNotification(r.Context(), cfg, notifications.EventLoginAttemptSuccess, username, adminUUID, "", "", r)

		setAuthCookie(w, r, cfg, accessToken, expiresAt)

		shared.WriteJSON(w, http.StatusOK, map[string]any{
			"response": map[string]any{
				"accessToken": accessToken,
			},
		})
	}
}

func AuthRegisterCompatHandler(manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			shared.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}

		var req LoginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			cfg.Logger.Warn("Auth register decode failed", "error", err, "remote_addr", r.RemoteAddr)
			shared.WriteJSONError(w, http.StatusBadRequest, "invalid JSON body")
			return
		}

		username := strings.TrimSpace(req.Username)
		password := strings.TrimSpace(req.Password)
		if username == "" || password == "" {
			shared.WriteJSONError(w, http.StatusBadRequest, "username and password are required")
			return
		}
		cfg.Logger.Trace("Auth register attempt", "username", username, "remote_addr", r.RemoteAddr)

		passwordHash, err := security.HashPassword(password)
		if err != nil {
			shared.WriteJSONError(w, http.StatusInternalServerError, "failed to create admin account")
			return
		}

		adminUUID := uuid.NewString()
		err = manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
			var adminCount int
			if countErr := db.QueryRow("SELECT COUNT(*) FROM admin").Scan(&adminCount); countErr != nil {
				return countErr
			}
			if adminCount > 0 {
				return errAdminAlreadyConfigured
			}

			_, execErr := db.Exec(`
				INSERT INTO admin (uuid, username, password_hash, role)
				VALUES (?, ?, ?, 'ADMIN')
			`, adminUUID, username, passwordHash)
			return execErr
		})
		if errors.Is(err, errAdminAlreadyConfigured) {
			cfg.Logger.Warn("Auth register blocked: admin already configured", "username", username, "remote_addr", r.RemoteAddr)
			shared.WriteJSONError(w, http.StatusForbidden, "registration is not allowed")
			return
		}
		if err != nil {
			cfg.Logger.Error("Auth register failed", "username", username, "remote_addr", r.RemoteAddr, "error", err)
			shared.WriteJSONError(w, http.StatusInternalServerError, "failed to create initial admin account")
			return
		}

		accessToken, expiresAt, err := createAdminAccessToken(cfg, username, adminUUID, "ADMIN")
		if err != nil {
			cfg.Logger.Error("Failed to create JWT auth token", "error", err)
			shared.WriteJSONError(w, http.StatusInternalServerError, "failed to create auth token")
			return
		}
		cfg.Logger.Info("Auth register success", "username", username, "remote_addr", r.RemoteAddr)
		setAuthCookie(w, r, cfg, accessToken, expiresAt)

		shared.WriteJSON(w, http.StatusCreated, map[string]any{
			"response": map[string]any{
				"accessToken": accessToken,
			},
		})
	}
}

func AuthSetupHandler(manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			shared.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}

		var req SetupRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			shared.WriteJSONError(w, http.StatusBadRequest, "invalid JSON body")
			return
		}

		username := strings.TrimSpace(req.Username)
		password := strings.TrimSpace(req.Password)
		if username == "" || password == "" {
			shared.WriteJSONError(w, http.StatusBadRequest, "username and password are required")
			return
		}

		passwordHash, err := security.HashPassword(password)
		if err != nil {
			cfg.Logger.Error("Failed to hash setup password", "error", err)
			shared.WriteJSONError(w, http.StatusInternalServerError, "failed to create admin account")
			return
		}

		adminUUID := uuid.NewString()
		err = manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
			var adminCount int
			if countErr := db.QueryRow("SELECT COUNT(*) FROM admin").Scan(&adminCount); countErr != nil {
				return countErr
			}
			if adminCount > 0 {
				return errAdminAlreadyConfigured
			}

			_, execErr := db.Exec(`
				INSERT INTO admin (uuid, username, password_hash, role)
				VALUES (?, ?, ?, 'ADMIN')
			`, adminUUID, username, passwordHash)
			return execErr
		})
		if errors.Is(err, errAdminAlreadyConfigured) {
			shared.WriteJSONError(w, http.StatusConflict, "admin account already exists")
			return
		}
		if err != nil {
			cfg.Logger.Error("Failed to create initial admin account", "error", err)
			shared.WriteJSONError(w, http.StatusInternalServerError, "failed to create initial admin account")
			return
		}

		accessToken, expiresAt, err := createAdminAccessToken(cfg, username, adminUUID, "ADMIN")
		if err != nil {
			cfg.Logger.Error("Failed to create JWT auth token", "error", err)
			shared.WriteJSONError(w, http.StatusInternalServerError, "failed to create auth token")
			return
		}
		setAuthCookie(w, r, cfg, accessToken, expiresAt)

		brandingSettings, passwordSettings, _, _, bootstrapErr := getBootstrapData(manager)
		if bootstrapErr != nil {
			cfg.Logger.Warn("Failed to include bootstrap settings in setup response", "error", bootstrapErr)
			brandingSettings = defaultBrandingSettings()
			passwordSettings = defaultPasswordSettings()
		}

		shared.WriteJSON(w, http.StatusCreated, LoginResponse{
			Admin: AuthAdminInfo{
				UUID:              adminUUID,
				Username:          username,
				Role:              "ADMIN",
				SessionTTLMinutes: AuthTTLMinutes(),
			},
			ExpiresAt:        expiresAt,
			BrandingSettings: brandingSettings,
			PasswordSettings: passwordSettings,
			Message:          "admin account created",
		})
	}
}

func AuthLoginHandler(manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			shared.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}

		var req LoginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			shared.WriteJSONError(w, http.StatusBadRequest, "invalid JSON body")
			return
		}

		username := strings.TrimSpace(req.Username)
		password := strings.TrimSpace(req.Password)
		if username == "" || password == "" {
			shared.WriteJSONError(w, http.StatusBadRequest, "username and password are required")
			return
		}

		var (
			adminUUID          string
			storedPasswordHash string
			role               string
		)

		err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
			row := db.QueryRow(`
				SELECT uuid, password_hash, role
				FROM admin
				WHERE username = ?
				LIMIT 1
			`, username)

			if scanErr := row.Scan(&adminUUID, &storedPasswordHash, &role); scanErr != nil {
				if errors.Is(scanErr, sql.ErrNoRows) {
					return nil
				}
				return scanErr
			}
			return nil
		})
		if err != nil {
			cfg.Logger.Error("Failed to read admin credentials", "error", err)
			shared.WriteJSONError(w, http.StatusInternalServerError, "failed to validate credentials")
			return
		}
		if adminUUID == "" || !security.VerifyPassword(password, storedPasswordHash) {
			emitLoginNotification(r.Context(), cfg, notifications.EventLoginAttemptFailed, username, adminUUID, password, "invalid_credentials", r)
			shared.WriteJSONError(w, http.StatusUnauthorized, "invalid username or password")
			return
		}

		accessToken, expiresAt, err := createAdminAccessToken(cfg, username, adminUUID, role)
		if err != nil {
			cfg.Logger.Error("Failed to create JWT auth token", "error", err)
			shared.WriteJSONError(w, http.StatusInternalServerError, "failed to create auth token")
			return
		}

		emitLoginNotification(r.Context(), cfg, notifications.EventLoginAttemptSuccess, username, adminUUID, "", "", r)
		setAuthCookie(w, r, cfg, accessToken, expiresAt)

		brandingSettings, passwordSettings, _, _, err := getBootstrapData(manager)
		if err != nil {
			cfg.Logger.Warn("Failed to include bootstrap settings in login response", "error", err)
			brandingSettings = defaultBrandingSettings()
			passwordSettings = defaultPasswordSettings()
		}

		shared.WriteJSON(w, http.StatusOK, LoginResponse{
			Admin: AuthAdminInfo{
				UUID:              adminUUID,
				Username:          username,
				Role:              strings.ToUpper(role),
				SessionTTLMinutes: AuthTTLMinutes(),
			},
			ExpiresAt:        expiresAt,
			BrandingSettings: brandingSettings,
			PasswordSettings: passwordSettings,
			Message:          "login successful",
		})
	}
}

func AuthMeHandler(manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			shared.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}

		principal, ok := CurrentAuthPrincipal(r.Context())
		if !ok || principal == nil || principal.TokenType != "jwt_auth" {
			shared.WriteJSONError(w, http.StatusUnauthorized, "unauthorized")
			return
		}

		var (
			adminInfo AuthAdminInfo
			expiresAt = principal.ExpiresAt
		)
		err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
			row := db.QueryRow(`
				SELECT uuid, username, role
				FROM admin
				WHERE uuid = ?
				LIMIT 1
			`, principal.AdminUUID)

			if scanErr := row.Scan(&adminInfo.UUID, &adminInfo.Username, &adminInfo.Role); scanErr != nil {
				return scanErr
			}
			adminInfo.Role = strings.ToUpper(adminInfo.Role)
			adminInfo.SessionTTLMinutes = AuthTTLMinutes()
			return nil
		})
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				shared.WriteJSONError(w, http.StatusUnauthorized, "unauthorized")
				return
			}
			cfg.Logger.Error("Failed to load auth me payload", "error", err)
			shared.WriteJSONError(w, http.StatusInternalServerError, "failed to load current session")
			return
		}

		brandingSettings, passwordSettings, _, _, err := getBootstrapData(manager)
		if err != nil {
			cfg.Logger.Warn("Failed to load branding for auth/me response", "error", err)
			brandingSettings = defaultBrandingSettings()
			passwordSettings = defaultPasswordSettings()
		}

		shared.WriteJSON(w, http.StatusOK, LoginResponse{
			Admin:            adminInfo,
			ExpiresAt:        expiresAt,
			BrandingSettings: brandingSettings,
			PasswordSettings: passwordSettings,
			Principal:        principal,
		})
	}
}

func AuthLogoutHandler(manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			shared.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}

		secureCookie := middleware.IsSecureRequest(r, cfg)
		http.SetCookie(w, &http.Cookie{
			Name:     sessionCookieName,
			Value:    "",
			Path:     "/",
			Expires:  time.Unix(0, 0).UTC(),
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode,
			Secure:   secureCookie,
			MaxAge:   -1,
		})

		shared.WriteJSON(w, http.StatusOK, map[string]string{"message": "logged out"})
	}
}

func getBootstrapData(manager *dbmanager.DatabaseManager) (brandingSettings map[string]any, passwordSettings map[string]any, defaultUsername string, hasAdmin bool, err error) {
	brandingSettings = defaultBrandingSettings()
	passwordSettings = defaultPasswordSettings()
	defaultUsername = "admin"
	hasAdmin = false

	err = manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		row := db.QueryRow(`
			SELECT branding_settings, password_settings
			FROM exodus_settings
			WHERE id = 1
			LIMIT 1
		`)

		var brandingRaw, passwordRaw sql.NullString
		if scanErr := row.Scan(&brandingRaw, &passwordRaw); scanErr != nil {
			if errors.Is(scanErr, sql.ErrNoRows) {
				return nil
			}
			return scanErr
		}

		if brandingRaw.Valid && strings.TrimSpace(brandingRaw.String) != "" {
			var tmp map[string]any
			if json.Unmarshal([]byte(brandingRaw.String), &tmp) == nil && len(tmp) > 0 {
				brandingSettings = tmp
			}
		}
		if passwordRaw.Valid && strings.TrimSpace(passwordRaw.String) != "" {
			var tmp map[string]any
			if json.Unmarshal([]byte(passwordRaw.String), &tmp) == nil && len(tmp) > 0 {
				passwordSettings = tmp
			}
		}

		var adminCount int
		if countErr := db.QueryRow("SELECT COUNT(*) FROM admin").Scan(&adminCount); countErr != nil {
			return countErr
		}
		hasAdmin = adminCount > 0

		if hasAdmin {
			var firstUsername sql.NullString
			if firstErr := db.QueryRow("SELECT username FROM admin ORDER BY created_at ASC LIMIT 1").Scan(&firstUsername); firstErr != nil && !errors.Is(firstErr, sql.ErrNoRows) {
				return firstErr
			}
			if firstUsername.Valid && strings.TrimSpace(firstUsername.String) != "" {
				defaultUsername = firstUsername.String
			}
		}
		return nil
	})
	return
}

func defaultBrandingSettings() map[string]any {
	return map[string]any{
		"title":   "EXODUS",
		"logoUrl": nil,
	}
}

func defaultPasswordSettings() map[string]any {
	return map[string]any{
		"enabled": true,
	}
}

func getAuthMethodsStatus(manager *dbmanager.DatabaseManager) (passkeyEnabled bool, tgEnabled bool, tgBotID *int, oauth2Providers map[string]bool) {
	passkeyEnabled = false
	tgEnabled = false
	oauth2Providers = map[string]bool{
		"github":   false,
		"yandex":   false,
		"generic":  false,
		"keycloak": false,
		"pocketid": false,
	}

	_ = manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		row := db.QueryRow(`
			SELECT passkey_settings, oauth2_settings, tg_auth_settings
			FROM exodus_settings
			WHERE id = 1
			LIMIT 1
		`)

		var passkeyRaw, oauth2Raw, tgRaw sql.NullString
		if err := row.Scan(&passkeyRaw, &oauth2Raw, &tgRaw); err != nil {
			return nil
		}

		if passkeyRaw.Valid && strings.TrimSpace(passkeyRaw.String) != "" {
			var passkeyObj map[string]any
			if json.Unmarshal([]byte(passkeyRaw.String), &passkeyObj) == nil {
				if enabled, ok := passkeyObj["enabled"].(bool); ok {
					passkeyEnabled = enabled
				}
			}
		}

		if tgRaw.Valid && strings.TrimSpace(tgRaw.String) != "" {
			var tgObj map[string]any
			if json.Unmarshal([]byte(tgRaw.String), &tgObj) == nil {
				if enabled, ok := tgObj["enabled"].(bool); ok {
					tgEnabled = enabled
				}
				if token, ok := tgObj["botToken"].(string); ok && strings.TrimSpace(token) != "" {
					prefix := strings.SplitN(token, ":", 2)[0]
					if parsed, err := strconv.Atoi(prefix); err == nil {
						botID := parsed
						tgBotID = &botID
					}
				}
			}
		}

		if oauth2Raw.Valid && strings.TrimSpace(oauth2Raw.String) != "" {
			var oauthObj map[string]map[string]any
			if json.Unmarshal([]byte(oauth2Raw.String), &oauthObj) == nil {
				for provider := range oauth2Providers {
					if providerCfg, ok := oauthObj[provider]; ok {
						if enabled, ok := providerCfg["enabled"].(bool); ok {
							oauth2Providers[provider] = enabled
						}
					}
				}
			}
		}

		return nil
	})

	return passkeyEnabled, tgEnabled, tgBotID, oauth2Providers
}

func resolvePasswordAuthEnabled(
	passwordSettings map[string]any,
	passkeyEnabled bool,
	tgEnabled bool,
	oauth2Providers map[string]bool,
	cfg *config.BackendConfig,
) bool {
	passwordEnabled := true
	if raw, ok := passwordSettings["enabled"]; ok {
		if value, ok := raw.(bool); ok {
			passwordEnabled = value
		}
	}

	if passwordEnabled {
		return true
	}

	hasOAuth2Enabled := false
	for _, enabled := range oauth2Providers {
		if enabled {
			hasOAuth2Enabled = true
			break
		}
	}

	// Avoid hard lockout when all auth methods are disabled by misconfiguration.
	if !passkeyEnabled && !tgEnabled && !hasOAuth2Enabled {
		if cfg != nil && cfg.Logger != nil {
			cfg.Logger.Warn("All authentication methods are disabled. Falling back to password auth enabled=true to prevent lockout.")
		}
		return true
	}

	return false
}

func nullableString(value string) *string {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil
	}
	return &value
}
