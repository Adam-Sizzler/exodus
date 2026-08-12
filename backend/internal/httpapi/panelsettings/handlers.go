package panelsettings

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"exodus/internal/config"
	"exodus/internal/httpapi/shared"
	"exodus/internal/security"

	"github.com/google/uuid"
)

func PanelSettingsHandler(db *sql.DB, cfg *config.BackendConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			settings, err := loadPanelSettings(r.Context(), db)
			if err != nil {
				cfg.Logger.Error("Failed to load panel settings", "error", err)
				shared.WriteJSONError(w, http.StatusInternalServerError, "failed to load panel settings")
				return
			}
			shared.WriteJSON(w, http.StatusOK, PanelSettingsResponse{Settings: settings})
		case http.MethodPatch, http.MethodPut:
			var payload map[string]json.RawMessage
			if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
				shared.WriteJSONError(w, http.StatusBadRequest, "invalid JSON body")
				return
			}
			candidateSettings, err := loadPanelSettings(r.Context(), db)
			if err != nil {
				cfg.Logger.Error("Failed to load panel settings before update", "error", err)
				shared.WriteJSONError(w, http.StatusInternalServerError, "failed to load panel settings")
				return
			}

			validKeys := map[string]string{
				"passkey_settings":  "passkey_settings",
				"oauth2_settings":   "oauth2_settings",
				"password_settings": "password_settings",
				"branding_settings": "branding_settings",
			}

			setClauses := make([]string, 0, len(payload))
			args := make([]any, 0, len(payload))
			authSettingsTouched := false
			idx := 1
			for key, raw := range payload {
				column, ok := validKeys[key]
				if !ok {
					continue
				}
				raw = json.RawMessage(strings.TrimSpace(string(raw)))
				if len(raw) == 0 || string(raw) == "null" {
					continue
				}
				if !json.Valid(raw) {
					shared.WriteJSONError(w, http.StatusBadRequest, "invalid JSON value for "+key)
					return
				}
				var decoded any
				if err := json.Unmarshal(raw, &decoded); err != nil {
					shared.WriteJSONError(w, http.StatusBadRequest, "invalid JSON value for "+key)
					return
				}
				candidateSettings[column] = decoded
				if column == "passkey_settings" || column == "oauth2_settings" || column == "password_settings" {
					authSettingsTouched = true
				}
				setClauses = append(setClauses, fmt.Sprintf("%s = $%d", column, idx))
				args = append(args, string(raw))
				idx++
			}

			if len(setClauses) == 0 {
				shared.WriteJSONError(w, http.StatusBadRequest, "no valid fields provided")
				return
			}
			if authSettingsTouched {
				if err := validateAuthenticationSettings(candidateSettings); err != nil {
					shared.WriteJSONError(w, http.StatusBadRequest, err.Error())
					return
				}
			}

			if _, execErr := db.ExecContext(r.Context(), `
				INSERT INTO exodus_settings (
					id, passkey_settings, oauth2_settings, password_settings, branding_settings
				) VALUES (
					1, $1, $2, $3, $4
				)
				ON CONFLICT (id) DO NOTHING
			`,
				`{"rpId":null,"origin":null,"enabled":false}`,
				`{"github":{"enabled":false,"clientId":null,"clientSecret":null,"allowedEmails":[]},"yandex":{"enabled":false,"clientId":null,"clientSecret":null,"allowedEmails":[]},"generic":{"enabled":false,"clientId":null,"tokenUrl":null,"withPkce":false,"clientSecret":null,"allowedEmails":[],"frontendDomain":null,"authorizationUrl":null},"keycloak":{"realm":null,"enabled":false,"clientId":null,"clientSecret":null,"allowedEmails":[],"frontendDomain":null,"keycloakDomain":null},"pocketid":{"enabled":false,"clientId":null,"plainDomain":null,"frontendDomain":null,"clientSecret":null,"allowedEmails":[]},"telegram":{"enabled":false,"clientId":null,"clientSecret":null,"allowedIds":[],"frontendDomain":null}}`,
				`{"enabled":true}`,
				`{"title":"EXODUS","logoUrl":null}`,
			); execErr != nil {
				cfg.Logger.Error("Failed to seed panel settings", "error", execErr)
				shared.WriteJSONError(w, http.StatusInternalServerError, "failed to update panel settings")
				return
			}

			query := "UPDATE exodus_settings SET " + strings.Join(setClauses, ", ") + " WHERE id = 1"
			if _, execErr := db.ExecContext(r.Context(), query, args...); execErr != nil {
				cfg.Logger.Error("Failed to update panel settings", "error", execErr)
				shared.WriteJSONError(w, http.StatusInternalServerError, "failed to update panel settings")
				return
			}

			settings, err := loadPanelSettings(r.Context(), db)
			if err != nil {
				cfg.Logger.Error("Failed to load panel settings after update", "error", err)
				shared.WriteJSONError(w, http.StatusInternalServerError, "failed to load updated panel settings")
				return
			}
			shared.WriteJSON(w, http.StatusOK, PanelSettingsResponse{Settings: settings})
		default:
			shared.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
		}
	}
}

// PanelAPITokensHandler godoc
// @Summary      Manage API tokens
// @Description  List or create (201) scoped panel API tokens
// @Tags         API Tokens Controller
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body      object  false  "API token creation parameters"
// @Success      200   {object}  map[string]any
// @Success      201   {object}  map[string]any
// @Failure      400   {object}  shared.ErrorResponse
// @Failure      500   {object}  shared.ErrorResponse
// @Router       /tokens [get]
// @Router       /tokens [post]
func PanelAPITokensHandler(db *sql.DB, cfg *config.BackendConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			tokens, err := loadAPITokens(r.Context(), db)
			if err != nil {
				cfg.Logger.Error("Failed to load api tokens", "error", err)
				shared.WriteJSONError(w, http.StatusInternalServerError, "failed to load api tokens")
				return
			}
			shared.WriteJSON(w, http.StatusOK, map[string]any{
				"response": map[string]any{
					"tokens": toAPITokensResponse(tokens),
				},
			})
		case http.MethodPost:
			var payload struct {
				TokenNameLegacy string   `json:"token_name"`
				TokenName       string   `json:"tokenName"`
				Name            string   `json:"name"`
				ExpiresInDays   int      `json:"expiresInDays"`
				Scopes          []string `json:"scopes"`
			}
			if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
				shared.WriteJSONError(w, http.StatusBadRequest, "invalid JSON body")
				return
			}
			tokenName := strings.TrimSpace(payload.Name)
			if tokenName == "" {
				tokenName = strings.TrimSpace(payload.TokenName)
			}
			if tokenName == "" {
				tokenName = strings.TrimSpace(payload.TokenNameLegacy)
			}
			if tokenName == "" {
				shared.WriteJSONError(w, http.StatusBadRequest, "name is required")
				return
			}
			expiresInDays := payload.ExpiresInDays
			if expiresInDays <= 0 {
				expiresInDays = int(security.APITokenLifetime / (24 * time.Hour))
			}
			scopes := normalizeAPITokenScopes(payload.Scopes)

			tokenUUID := uuid.NewString()
			lifetime := time.Duration(expiresInDays) * 24 * time.Hour
			tokenValue, expiresAtUnix, err := security.SignAPITokenJWTWithLifetime(cfg.JWT.AuthSecret, tokenUUID, lifetime)
			if err != nil {
				cfg.Logger.Error("Failed to generate api token", "error", err)
				shared.WriteJSONError(w, http.StatusInternalServerError, "failed to generate api token")
				return
			}
			expireAt := time.Unix(expiresAtUnix, 0).UTC()

			record := APITokenRecord{
				UUID:     tokenUUID,
				Token:    tokenValue,
				Name:     tokenName,
				ExpireAt: expireAt,
				Scopes:   scopes,
			}

			if _, execErr := db.ExecContext(r.Context(), `
				INSERT INTO api_tokens (uuid, name, expire_at, scopes)
				VALUES ($1, $2, $3, $4::text[])
			`, record.UUID, record.Name, record.ExpireAt, postgresTextArrayLiteral(record.Scopes)); execErr != nil {
				cfg.Logger.Error("Failed to insert api token", "error", execErr)
				shared.WriteJSONError(w, http.StatusInternalServerError, "failed to save api token")
				return
			}

			shared.WriteJSON(w, http.StatusCreated, map[string]any{
				"response": toAPITokenResponse(record, true),
			})
		default:
			shared.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
		}
	}
}

// PanelAPITokenScopesHandler godoc
// @Summary      Get available API token scopes
// @Description  Get list of all supported permission scopes for panel tokens
// @Tags         API Tokens Controller
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  map[string]any
// @Router       /tokens/scopes [get]
func PanelAPITokenScopesHandler(_ *sql.DB, cfg *config.BackendConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			shared.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}
		shared.WriteJSON(w, http.StatusOK, map[string]any{
			"response": map[string]any{
				"wildcard":  "*",
				"resources": buildAPITokenScopes(cfg),
			},
		})
	}
}

// PanelAPITokensOttHandler godoc
// @Summary      Issue One-Time Token (OTT)
// @Description  Issue short-lived token for internal utilities and redirects
// @Tags         API Tokens Controller
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  map[string]any
// @Failure      500  {object}  shared.ErrorResponse
// @Router       /tokens/ott [post]
func PanelAPITokensOttHandler(_ *sql.DB, cfg *config.BackendConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			shared.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}

		secret := cfg.JWT.AuthSecret
		ott, err := security.SignOttJWT(secret)
		if err != nil {
			cfg.Logger.Error("Failed to issue OTT token", "error", err)
			shared.WriteJSONError(w, http.StatusInternalServerError, "failed to issue ott token")
			return
		}

		shared.WriteJSON(w, http.StatusOK, map[string]any{
			"response": map[string]any{
				"ott": ott,
			},
		})
	}
}

// PanelAPITokenByUUIDHandler godoc
// @Summary      Delete API token
// @Description  Revoke API token by token UUID
// @Tags         API Tokens Controller
// @Security     BearerAuth
// @Param        uuid  path  string  true  "Token UUID" format(uuid)
// @Success      204
// @Failure      400  {object}  shared.ErrorResponse
// @Failure      404  {object}  shared.ErrorResponse
// @Failure      500  {object}  shared.ErrorResponse
// @Router       /tokens/{uuid} [delete]
func PanelAPITokenByUUIDHandler(db *sql.DB, cfg *config.BackendConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenUUID := strings.TrimSpace(strings.TrimPrefix(r.URL.Path, "/api/tokens/"))
		if tokenUUID == "" {
			switch r.Method {
			case http.MethodGet, http.MethodPost:
				PanelAPITokensHandler(db, cfg)(w, r)
			default:
				shared.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
			}
			return
		}

		if _, err := uuid.Parse(tokenUUID); err != nil {
			shared.WriteJSONError(w, http.StatusBadRequest, "invalid token UUID")
			return
		}
		if r.Method != http.MethodDelete {
			shared.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}

		result, execErr := db.ExecContext(r.Context(), "DELETE FROM api_tokens WHERE uuid = $1", tokenUUID)
		if execErr != nil {
			cfg.Logger.Error("Failed to delete api token", "uuid", tokenUUID, "error", execErr)
			shared.WriteJSONError(w, http.StatusInternalServerError, "failed to delete api token")
			return
		}
		rowsAffected, execErr := result.RowsAffected()
		if execErr != nil {
			cfg.Logger.Error("Failed to read rows affected for api token deletion", "uuid", tokenUUID, "error", execErr)
			shared.WriteJSONError(w, http.StatusInternalServerError, "failed to delete api token")
			return
		}
		if rowsAffected == 0 {
			shared.WriteJSONError(w, http.StatusNotFound, "api token not found")
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

// ExodusSettingsHandler godoc
// @Summary      Manage Exodus panel settings
// @Description  Get or update Exodus panel branding, passkey, oauth2, and password policies
// @Tags         Exodus Settings Controller
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body      object  false  "Settings payload to update"
// @Success      200   {object}  map[string]any
// @Failure      400   {object}  shared.ErrorResponse
// @Failure      500   {object}  shared.ErrorResponse
// @Router       /exodus-settings [get]
// @Router       /exodus-settings [patch]
func ExodusSettingsHandler(db *sql.DB, cfg *config.BackendConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			settings, err := loadPanelSettings(r.Context(), db)
			if err != nil {
				cfg.Logger.Error("Failed to load exodus settings", "error", err)
				shared.WriteJSONError(w, http.StatusInternalServerError, "failed to load settings")
				return
			}
			shared.WriteJSON(w, http.StatusOK, map[string]any{"response": toExodusSettingsResponse(settings)})
		case http.MethodPatch:
			var payload map[string]json.RawMessage
			if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
				shared.WriteJSONError(w, http.StatusBadRequest, "invalid JSON body")
				return
			}
			adapted := map[string]json.RawMessage{}
			for key, value := range payload {
				switch key {
				case "passkeySettings":
					adapted["passkey_settings"] = value
				case "oauth2Settings":
					adapted["oauth2_settings"] = value
				case "passwordSettings":
					adapted["password_settings"] = value
				case "brandingSettings":
					adapted["branding_settings"] = value
				}
			}
			if len(adapted) == 0 {
				shared.WriteJSONError(w, http.StatusBadRequest, "no valid fields provided")
				return
			}

			body, _ := json.Marshal(adapted)
			r2 := r.Clone(r.Context())
			r2.Method = http.MethodPatch
			r2.Body = http.NoBody
			r2.Body = ioNopCloser{strings.NewReader(string(body))}

			rec := responseRecorder{header: http.Header{}}
			PanelSettingsHandler(db, cfg)(&rec, r2)
			if rec.status >= 400 {
				w.WriteHeader(rec.status)
				_, _ = w.Write(rec.body)
				return
			}

			settings, err := loadPanelSettings(r.Context(), db)
			if err != nil {
				cfg.Logger.Error("Failed to load exodus settings after update", "error", err)
				shared.WriteJSONError(w, http.StatusInternalServerError, "failed to load updated settings")
				return
			}
			shared.WriteJSON(w, http.StatusOK, map[string]any{"response": toExodusSettingsResponse(settings)})
		default:
			shared.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
		}
	}
}
