package auth

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
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

var errAdminAlreadyConfigured = errors.New("admin already configured")

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
			passkeyEnabled, oauth2Providers := getAuthMethodsStatus(manager)
			passwordEnabled := resolvePasswordAuthEnabled(
				passwordSettings,
				passkeyEnabled,
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

		rateLimitKey := loginRateLimitKey(r)
		if !globalLoginRateLimiter.allow(rateLimitKey) {
			cfg.Logger.Warn("Auth login rate limited", "remote_addr", r.RemoteAddr)
			shared.WriteJSONError(w, http.StatusTooManyRequests, "too many login attempts, please try again later")
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

		passkeyEnabled, oauth2Providers := getAuthMethodsStatus(manager)
		passwordEnabled := resolvePasswordAuthEnabled(
			passwordSettings,
			passkeyEnabled,
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
		hashToVerify := storedPasswordHash
		if adminUUID == "" {
			hashToVerify = dummyPasswordHash
		}
		if adminUUID == "" || !security.VerifyPassword(password, hashToVerify) {
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
		globalLoginRateLimiter.reset(rateLimitKey)
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

func OAuth2AuthorizeHandler(manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			shared.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}
		var req oauthAuthorizeRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			shared.SendError(w, http.StatusBadRequest, "invalid JSON", err, cfg)
			return
		}
		provider := strings.ToLower(strings.TrimSpace(req.Provider))
		if !isSupportedOAuthProvider(provider) {
			shared.SendError(w, http.StatusBadRequest, "OAuth2 provider not found", nil, cfg)
			return
		}
		if !isLoginAllowed(manager) {
			shared.SendError(w, http.StatusForbidden, "login is not allowed", nil, cfg)
			return
		}
		settings, err := loadOAuthSettings(manager)
		if err != nil {
			shared.SendError(w, http.StatusInternalServerError, "failed to load OAuth2 settings", err, cfg)
			return
		}
		providerSettings := getOAuthProviderSettings(settings, provider)
		if !providerSettings.Enabled {
			shared.SendError(w, http.StatusForbidden, "OAuth2 provider is disabled", nil, cfg)
			return
		}

		state, err := security.GenerateRandomToken(32)
		if err != nil {
			shared.SendError(w, http.StatusInternalServerError, "failed to generate OAuth2 state", err, cfg)
			return
		}
		codeVerifier := ""
		authURL, err := buildAuthorizationURL(provider, providerSettings, cfg.Panel.BasePath, state, &codeVerifier)
		if err != nil {
			shared.SendError(w, http.StatusInternalServerError, "OAuth2 authorize error", err, cfg)
			return
		}
		storeOAuthState(provider, state, codeVerifier)
		shared.WriteJSON(w, http.StatusOK, map[string]any{
			"response": map[string]any{"authorizationUrl": authURL},
		})
	}
}

func OAuth2CallbackHandler(manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			shared.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}
		var req oauthCallbackRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			shared.SendError(w, http.StatusBadRequest, "invalid JSON", err, cfg)
			return
		}
		provider := strings.ToLower(strings.TrimSpace(req.Provider))
		if !isSupportedOAuthProvider(provider) {
			emitExternalLoginNotification(r.Context(), cfg, notifications.EventLoginAttemptFailed, "oauth2", provider, "", "", "unsupported_provider", r)
			shared.SendError(w, http.StatusBadRequest, "OAuth2 provider not found", nil, cfg)
			return
		}
		stateEntry, ok := takeOAuthState(provider)
		if !ok || stateEntry.State != strings.TrimSpace(req.State) {
			emitExternalLoginNotification(r.Context(), cfg, notifications.EventLoginAttemptFailed, "oauth2", provider, "", "", "state_mismatch", r)
			shared.SendError(w, http.StatusForbidden, "OAuth2 state mismatch", nil, cfg)
			return
		}
		if !isLoginAllowed(manager) {
			emitExternalLoginNotification(r.Context(), cfg, notifications.EventLoginAttemptFailed, "oauth2", provider, "", "", "login_not_allowed", r)
			shared.SendError(w, http.StatusForbidden, "login is not allowed", nil, cfg)
			return
		}
		settings, err := loadOAuthSettings(manager)
		if err != nil {
			shared.SendError(w, http.StatusInternalServerError, "failed to load OAuth2 settings", err, cfg)
			return
		}
		providerSettings := getOAuthProviderSettings(settings, provider)
		if !providerSettings.Enabled {
			emitExternalLoginNotification(r.Context(), cfg, notifications.EventLoginAttemptFailed, "oauth2", provider, "", "", "provider_disabled", r)
			shared.SendError(w, http.StatusForbidden, "OAuth2 provider is disabled", nil, cfg)
			return
		}
		email, hasCustomClaim, err := exchangeOAuthCode(r.Context(), provider, providerSettings, cfg.Panel.BasePath, strings.TrimSpace(req.Code), stateEntry.CodeVerifier)
		if err != nil {
			emitExternalLoginNotification(r.Context(), cfg, notifications.EventLoginAttemptFailed, "oauth2", provider, "", "", "callback_error", r)
			shared.SendError(w, http.StatusForbidden, "OAuth2 callback error", err, cfg)
			return
		}
		if email == "" || !isOAuthPrincipalAllowed(provider, email, hasCustomClaim, providerSettings) {
			emitExternalLoginNotification(r.Context(), cfg, notifications.EventLoginAttemptFailed, "oauth2", provider, email, "", "email_not_allowed", r)
			shared.SendError(w, http.StatusForbidden, "OAuth2 principal is not allowed", nil, cfg)
			return
		}
		token, adminUUID, err := createFirstAdminSession(w, r, manager, cfg)
		if err != nil {
			emitExternalLoginNotification(r.Context(), cfg, notifications.EventLoginAttemptFailed, "oauth2", provider, email, "", "session_create_failed", r)
			shared.SendError(w, http.StatusForbidden, "failed to create OAuth2 session", err, cfg)
			return
		}
		emitExternalLoginNotification(r.Context(), cfg, notifications.EventLoginAttemptSuccess, "oauth2", provider, email, adminUUID, "", r)
		shared.WriteJSON(w, http.StatusOK, map[string]any{
			"response": map[string]any{"accessToken": token},
		})
	}
}
