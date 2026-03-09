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

	"v2ray-stat/backend/config"
	dbmanager "v2ray-stat/backend/db/manager"
	"v2ray-stat/backend/httpapi/middleware"
	"v2ray-stat/backend/httpapi/shared"
	"v2ray-stat/backend/security"

	"github.com/google/uuid"
)

const sessionCookieName = "v2rs_session"

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
	case "/api/health", "/api/v1/health",
		"/api/v1/auth/login", "/api/v1/auth/bootstrap", "/api/v1/auth/setup",
		"/api/auth/status", "/api/auth/login", "/api/auth/register",
		"/api/auth/oauth2/authorize", "/api/auth/oauth2/callback", "/api/auth/oauth2/tg/callback",
		"/api/auth/passkey/authentication/options", "/api/auth/passkey/authentication/verify",
		"/api/system/metadata", "/api/v1/system/metadata":
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

	nowUnix := time.Now().UTC().Unix()

	var sessionPrincipal *AuthPrincipal
	sessionErr := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		row := db.QueryRow(`
			SELECT s.admin_uuid, a.username, a.role, s.expires_at
			FROM admin_sessions s
			JOIN admin a ON a.uuid = s.admin_uuid
			WHERE s.session_token = ?
			LIMIT 1
		`, token)

		var adminUUID, username, role string
		var expiresAt int64
		if err := row.Scan(&adminUUID, &username, &role, &expiresAt); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil
			}
			return err
		}

		if expiresAt <= nowUnix {
			_, _ = db.Exec("DELETE FROM admin_sessions WHERE session_token = ?", token)
			return nil
		}

		sessionPrincipal = &AuthPrincipal{
			AdminUUID: adminUUID,
			Username:  username,
			Role:      strings.ToUpper(role),
			TokenType: "session",
			ExpiresAt: expiresAt,
		}
		return nil
	})
	if sessionErr != nil {
		return nil, sessionErr
	}
	if sessionPrincipal != nil {
		return sessionPrincipal, nil
	}

	var apiPrincipal *AuthPrincipal
	apiErr := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		row := db.QueryRow(`
			SELECT uuid, token_name
			FROM api_tokens
			WHERE token = ?
			LIMIT 1
		`, token)

		var tokenUUID, tokenName string
		if err := row.Scan(&tokenUUID, &tokenName); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil
			}
			return err
		}

		apiPrincipal = &AuthPrincipal{
			AdminUUID: tokenUUID,
			Username:  tokenName,
			Role:      "API",
			TokenType: "db_api_token",
		}
		return nil
	})
	if apiErr != nil {
		return nil, apiErr
	}
	if apiPrincipal != nil {
		return apiPrincipal, nil
	}

	return nil, errors.New("invalid auth token")
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
			shared.WriteJSONError(w, http.StatusForbidden, "login is not allowed")
			return
		}

		var (
			adminUUID          string
			storedPasswordHash string
			sessionTTLMinutes  int
		)

		err = manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
			row := db.QueryRow(`
				SELECT uuid, password_hash, session_ttl_minutes
				FROM admin
				WHERE username = ?
				LIMIT 1
			`, username)

			var ttl sql.NullInt64
			if scanErr := row.Scan(&adminUUID, &storedPasswordHash, &ttl); scanErr != nil {
				if errors.Is(scanErr, sql.ErrNoRows) {
					return nil
				}
				return scanErr
			}

			if ttl.Valid && ttl.Int64 > 0 {
				sessionTTLMinutes = int(ttl.Int64)
			} else {
				sessionTTLMinutes = 60
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
			shared.WriteJSONError(w, http.StatusForbidden, "invalid username or password")
			return
		}

		sessionToken, err := security.GenerateRandomToken(48)
		if err != nil {
			cfg.Logger.Error("Failed to generate session token", "error", err)
			shared.WriteJSONError(w, http.StatusInternalServerError, "failed to create session")
			return
		}

		expiresAt := time.Now().UTC().Add(time.Duration(sessionTTLMinutes) * time.Minute).Unix()
		if err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
			if _, execErr := db.Exec("DELETE FROM admin_sessions WHERE admin_uuid = ?", adminUUID); execErr != nil {
				return execErr
			}
			_, execErr := db.Exec(`
				INSERT INTO admin_sessions (session_token, admin_uuid, expires_at)
				VALUES (?, ?, ?)
			`, sessionToken, adminUUID, expiresAt)
			return execErr
		}); err != nil {
			cfg.Logger.Error("Failed to persist admin session", "error", err)
			shared.WriteJSONError(w, http.StatusInternalServerError, "failed to persist session")
			return
		}
		cfg.Logger.Info("Auth login success", "username", username, "remote_addr", r.RemoteAddr)

		secureCookie := middleware.IsSecureRequest(r, cfg)
		http.SetCookie(w, &http.Cookie{
			Name:     sessionCookieName,
			Value:    sessionToken,
			Path:     "/",
			Expires:  time.Unix(expiresAt, 0).UTC(),
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode,
			Secure:   secureCookie,
		})

		shared.WriteJSON(w, http.StatusOK, map[string]any{
			"response": map[string]any{
				"accessToken": sessionToken,
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
		sessionToken, err := security.GenerateRandomToken(48)
		if err != nil {
			shared.WriteJSONError(w, http.StatusInternalServerError, "failed to create admin session")
			return
		}
		sessionTTLMinutes := 60
		expiresAt := time.Now().UTC().Add(time.Duration(sessionTTLMinutes) * time.Minute).Unix()

		err = manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
			var adminCount int
			if countErr := db.QueryRow("SELECT COUNT(*) FROM admin").Scan(&adminCount); countErr != nil {
				return countErr
			}
			if adminCount > 0 {
				return errAdminAlreadyConfigured
			}

			if _, execErr := db.Exec(`
				INSERT INTO admin (uuid, username, password_hash, role, session_ttl_minutes)
				VALUES (?, ?, ?, 'ADMIN', ?)
			`, adminUUID, username, passwordHash, sessionTTLMinutes); execErr != nil {
				return execErr
			}
			_, execErr := db.Exec(`
				INSERT INTO admin_sessions (session_token, admin_uuid, expires_at)
				VALUES (?, ?, ?)
			`, sessionToken, adminUUID, expiresAt)
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
		cfg.Logger.Info("Auth register success", "username", username, "remote_addr", r.RemoteAddr)

		secureCookie := middleware.IsSecureRequest(r, cfg)
		http.SetCookie(w, &http.Cookie{
			Name:     sessionCookieName,
			Value:    sessionToken,
			Path:     "/",
			Expires:  time.Unix(expiresAt, 0).UTC(),
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode,
			Secure:   secureCookie,
		})

		shared.WriteJSON(w, http.StatusCreated, map[string]any{
			"response": map[string]any{
				"accessToken": sessionToken,
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

		sessionTTLMinutes := req.SessionTTLMinutes
		if sessionTTLMinutes <= 0 {
			sessionTTLMinutes = 60
		}

		passwordHash, err := security.HashPassword(password)
		if err != nil {
			cfg.Logger.Error("Failed to hash setup password", "error", err)
			shared.WriteJSONError(w, http.StatusInternalServerError, "failed to create admin account")
			return
		}

		adminUUID := uuid.NewString()
		sessionToken, err := security.GenerateRandomToken(48)
		if err != nil {
			cfg.Logger.Error("Failed to generate setup session token", "error", err)
			shared.WriteJSONError(w, http.StatusInternalServerError, "failed to create admin session")
			return
		}
		expiresAt := time.Now().UTC().Add(time.Duration(sessionTTLMinutes) * time.Minute).Unix()

		err = manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
			var adminCount int
			if countErr := db.QueryRow("SELECT COUNT(*) FROM admin").Scan(&adminCount); countErr != nil {
				return countErr
			}
			if adminCount > 0 {
				return errAdminAlreadyConfigured
			}

			if _, execErr := db.Exec(`
				INSERT INTO admin (uuid, username, password_hash, role, session_ttl_minutes)
				VALUES (?, ?, ?, 'ADMIN', ?)
			`, adminUUID, username, passwordHash, sessionTTLMinutes); execErr != nil {
				return execErr
			}

			_, execErr := db.Exec(`
				INSERT INTO admin_sessions (session_token, admin_uuid, expires_at)
				VALUES (?, ?, ?)
			`, sessionToken, adminUUID, expiresAt)
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

		secureCookie := middleware.IsSecureRequest(r, cfg)
		http.SetCookie(w, &http.Cookie{
			Name:     sessionCookieName,
			Value:    sessionToken,
			Path:     "/",
			Expires:  time.Unix(expiresAt, 0).UTC(),
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode,
			Secure:   secureCookie,
		})

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
				SessionTTLMinutes: sessionTTLMinutes,
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
			sessionTTLMinutes  int
		)

		err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
			row := db.QueryRow(`
				SELECT uuid, password_hash, role, session_ttl_minutes
				FROM admin
				WHERE username = ?
				LIMIT 1
			`, username)

			var ttl sql.NullInt64
			if scanErr := row.Scan(&adminUUID, &storedPasswordHash, &role, &ttl); scanErr != nil {
				if errors.Is(scanErr, sql.ErrNoRows) {
					return nil
				}
				return scanErr
			}

			if ttl.Valid && ttl.Int64 > 0 {
				sessionTTLMinutes = int(ttl.Int64)
			} else {
				sessionTTLMinutes = 60
			}
			return nil
		})
		if err != nil {
			cfg.Logger.Error("Failed to read admin credentials", "error", err)
			shared.WriteJSONError(w, http.StatusInternalServerError, "failed to validate credentials")
			return
		}
		if adminUUID == "" || !security.VerifyPassword(password, storedPasswordHash) {
			shared.WriteJSONError(w, http.StatusUnauthorized, "invalid username or password")
			return
		}

		sessionToken, err := security.GenerateRandomToken(48)
		if err != nil {
			cfg.Logger.Error("Failed to generate session token", "error", err)
			shared.WriteJSONError(w, http.StatusInternalServerError, "failed to create session")
			return
		}

		expiresAt := time.Now().UTC().Add(time.Duration(sessionTTLMinutes) * time.Minute).Unix()
		if err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
			// Single active session per admin account simplifies emergency password reset flows.
			if _, execErr := db.Exec("DELETE FROM admin_sessions WHERE admin_uuid = ?", adminUUID); execErr != nil {
				return execErr
			}
			_, execErr := db.Exec(`
				INSERT INTO admin_sessions (session_token, admin_uuid, expires_at)
				VALUES (?, ?, ?)
			`, sessionToken, adminUUID, expiresAt)
			return execErr
		}); err != nil {
			cfg.Logger.Error("Failed to persist admin session", "error", err)
			shared.WriteJSONError(w, http.StatusInternalServerError, "failed to persist session")
			return
		}

		secureCookie := middleware.IsSecureRequest(r, cfg)
		http.SetCookie(w, &http.Cookie{
			Name:     sessionCookieName,
			Value:    sessionToken,
			Path:     "/",
			Expires:  time.Unix(expiresAt, 0).UTC(),
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode,
			Secure:   secureCookie,
		})

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
				SessionTTLMinutes: sessionTTLMinutes,
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
		if !ok || principal == nil || principal.TokenType != "session" {
			shared.WriteJSONError(w, http.StatusUnauthorized, "unauthorized")
			return
		}

		var (
			adminInfo AuthAdminInfo
			expiresAt = principal.ExpiresAt
		)
		err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
			row := db.QueryRow(`
				SELECT uuid, username, role, session_ttl_minutes
				FROM admin
				WHERE uuid = ?
				LIMIT 1
			`, principal.AdminUUID)

			var ttl sql.NullInt64
			if scanErr := row.Scan(&adminInfo.UUID, &adminInfo.Username, &adminInfo.Role, &ttl); scanErr != nil {
				return scanErr
			}
			adminInfo.Role = strings.ToUpper(adminInfo.Role)
			if ttl.Valid && ttl.Int64 > 0 {
				adminInfo.SessionTTLMinutes = int(ttl.Int64)
			} else {
				adminInfo.SessionTTLMinutes = 60
			}
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

		token := readBearerToken(r.Header.Get("Authorization"))
		if token == "" {
			if cookie, err := r.Cookie(sessionCookieName); err == nil {
				token = strings.TrimSpace(cookie.Value)
			}
		}

		if token != "" {
			if err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
				_, execErr := db.Exec("DELETE FROM admin_sessions WHERE session_token = ?", token)
				return execErr
			}); err != nil {
				cfg.Logger.Warn("Failed to delete session during logout", "error", err)
			}
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
			FROM v2rs_settings
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
		"title":   "V2RS",
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
			FROM v2rs_settings
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
