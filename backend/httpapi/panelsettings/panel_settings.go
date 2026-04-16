package panelsettings

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"exodus/backend/config"
	dbmanager "exodus/backend/db/manager"
	"exodus/backend/httpapi/shared"
	"exodus/backend/security"

	"github.com/google/uuid"
)

type PanelSettingsResponse struct {
	Settings map[string]any `json:"settings"`
}

type APITokenRecord struct {
	UUID      string    `json:"uuid"`
	Token     string    `json:"token"`
	TokenName string    `json:"token_name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func PanelSettingsHandler(manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			settings, err := loadPanelSettings(manager)
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

			validKeys := map[string]string{
				"passkey_settings":  "passkey_settings",
				"oauth2_settings":   "oauth2_settings",
				"tg_auth_settings":  "tg_auth_settings",
				"password_settings": "password_settings",
				"branding_settings": "branding_settings",
				"modules_settings":  "modules_settings",
			}

			setClauses := make([]string, 0, len(payload))
			args := make([]any, 0, len(payload))
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
				setClauses = append(setClauses, column+" = ?")
				args = append(args, string(raw))
			}

			if len(setClauses) == 0 {
				shared.WriteJSONError(w, http.StatusBadRequest, "no valid fields provided")
				return
			}

			err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
				if _, execErr := db.Exec(`
					INSERT INTO exodus_settings (
						id, passkey_settings, oauth2_settings, tg_auth_settings, password_settings, branding_settings, modules_settings
					) VALUES (
						1, ?, ?, ?, ?, ?, ?
					)
					ON CONFLICT (id) DO NOTHING
				`,
					`{"rpId":null,"origin":null,"enabled":false}`,
					`{"github":{"enabled":false,"clientId":null,"clientSecret":null,"allowedEmails":[]},"yandex":{"enabled":false,"clientId":null,"clientSecret":null,"allowedEmails":[]},"generic":{"enabled":false,"clientId":null,"tokenUrl":null,"withPkce":false,"clientSecret":null,"allowedEmails":[],"frontendDomain":null,"authorizationUrl":null},"keycloak":{"realm":null,"enabled":false,"clientId":null,"clientSecret":null,"allowedEmails":[],"frontendDomain":null,"keycloakDomain":null},"pocketid":{"enabled":false,"clientId":null,"plainDomain":null,"clientSecret":null,"allowedEmails":[]}}`,
					`{"enabled":false,"adminIds":[],"botToken":null}`,
					`{"enabled":true}`,
					`{"title":"EXODUS","logoUrl":null}`,
					`{"haproxy":{"enabled":false}}`,
				); execErr != nil {
					return execErr
				}

				query := "UPDATE exodus_settings SET " + strings.Join(setClauses, ", ") + " WHERE id = 1"
				_, execErr := db.Exec(query, args...)
				return execErr
			})
			if err != nil {
				cfg.Logger.Error("Failed to update panel settings", "error", err)
				shared.WriteJSONError(w, http.StatusInternalServerError, "failed to update panel settings")
				return
			}

			settings, err := loadPanelSettings(manager)
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

func PanelAPITokensHandler(manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			tokens, err := loadAPITokens(manager)
			if err != nil {
				cfg.Logger.Error("Failed to load api tokens", "error", err)
				shared.WriteJSONError(w, http.StatusInternalServerError, "failed to load api tokens")
				return
			}
			shared.WriteJSON(w, http.StatusOK, map[string]any{
				"response": map[string]any{
					"apiKeys": toAPITokensResponse(tokens),
					"docs": map[string]any{
						"isDocsEnabled": false,
						"scalarPath":    nil,
						"swaggerPath":   nil,
					},
				},
			})
		case http.MethodPost:
			var payload struct {
				TokenNameLegacy string `json:"token_name"`
				TokenName       string `json:"tokenName"`
			}
			if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
				shared.WriteJSONError(w, http.StatusBadRequest, "invalid JSON body")
				return
			}
			tokenName := strings.TrimSpace(payload.TokenName)
			if tokenName == "" {
				tokenName = strings.TrimSpace(payload.TokenNameLegacy)
			}
			if tokenName == "" {
				shared.WriteJSONError(w, http.StatusBadRequest, "tokenName is required")
				return
			}

			tokenValue, err := security.GenerateRandomToken(64)
			if err != nil {
				cfg.Logger.Error("Failed to generate api token", "error", err)
				shared.WriteJSONError(w, http.StatusInternalServerError, "failed to generate api token")
				return
			}

			record := APITokenRecord{
				UUID:      uuid.NewString(),
				Token:     tokenValue,
				TokenName: tokenName,
			}

			err = manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
				_, execErr := db.Exec(`
					INSERT INTO api_tokens (uuid, token, token_name)
					VALUES (?, ?, ?)
				`, record.UUID, record.Token, record.TokenName)
				return execErr
			})
			if err != nil {
				cfg.Logger.Error("Failed to insert api token", "error", err)
				shared.WriteJSONError(w, http.StatusInternalServerError, "failed to save api token")
				return
			}

			// Return full token only in create response.
			shared.WriteJSON(w, http.StatusCreated, map[string]any{
				"response": map[string]any{
					"uuid":  record.UUID,
					"token": record.Token,
				},
			})
		default:
			shared.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
		}
	}
}

func PanelAPITokenByUUIDHandler(manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenUUID := strings.TrimSpace(strings.TrimPrefix(r.URL.Path, "/api/tokens/"))
		if tokenUUID == "" {
			switch r.Method {
			case http.MethodGet, http.MethodPost:
				PanelAPITokensHandler(manager, cfg)(w, r)
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

		var rowsAffected int64
		err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
			result, execErr := db.Exec("DELETE FROM api_tokens WHERE uuid = ?", tokenUUID)
			if execErr != nil {
				return execErr
			}
			ra, execErr := result.RowsAffected()
			if execErr != nil {
				return execErr
			}
			rowsAffected = ra
			return nil
		})
		if err != nil {
			cfg.Logger.Error("Failed to delete api token", "uuid", tokenUUID, "error", err)
			shared.WriteJSONError(w, http.StatusInternalServerError, "failed to delete api token")
			return
		}
		if rowsAffected == 0 {
			shared.WriteJSONError(w, http.StatusNotFound, "api token not found")
			return
		}

		shared.WriteJSON(w, http.StatusOK, map[string]any{"response": true})
	}
}

func ExodusSettingsHandler(manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			settings, err := loadPanelSettings(manager)
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
				case "tgAuthSettings":
					adapted["tg_auth_settings"] = value
				case "passwordSettings":
					adapted["password_settings"] = value
				case "brandingSettings":
					adapted["branding_settings"] = value
				case "modulesSettings":
					adapted["modules_settings"] = value
				}
			}
			if len(adapted) == 0 {
				shared.WriteJSONError(w, http.StatusBadRequest, "no valid fields provided")
				return
			}

			// Reuse existing update logic through PanelSettingsHandler payload format.
			body, _ := json.Marshal(adapted)
			r2 := r.Clone(r.Context())
			r2.Method = http.MethodPatch
			r2.Body = http.NoBody
			r2.Body = ioNopCloser{strings.NewReader(string(body))}

			rec := responseRecorder{header: http.Header{}}
			PanelSettingsHandler(manager, cfg)(&rec, r2)
			if rec.status >= 400 {
				w.WriteHeader(rec.status)
				_, _ = w.Write(rec.body)
				return
			}

			settings, err := loadPanelSettings(manager)
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

type apiTokenResponse struct {
	UUID      string    `json:"uuid"`
	Token     string    `json:"token"`
	TokenName string    `json:"tokenName"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func toAPITokensResponse(items []APITokenRecord) []apiTokenResponse {
	out := make([]apiTokenResponse, 0, len(items))
	for _, item := range items {
		out = append(out, apiTokenResponse{
			UUID:      item.UUID,
			Token:     item.Token,
			TokenName: item.TokenName,
			CreatedAt: item.CreatedAt,
			UpdatedAt: item.UpdatedAt,
		})
	}
	return out
}

func toExodusSettingsResponse(settings map[string]any) map[string]any {
	return map[string]any{
		"passkeySettings":  settings["passkey_settings"],
		"oauth2Settings":   settings["oauth2_settings"],
		"tgAuthSettings":   settings["tg_auth_settings"],
		"passwordSettings": settings["password_settings"],
		"brandingSettings": settings["branding_settings"],
		"modulesSettings":  settings["modules_settings"],
	}
}

// local lightweight helpers to internally call handlers without importing httptest.
type ioNopCloser struct {
	*strings.Reader
}

func (ioNopCloser) Close() error { return nil }

type responseRecorder struct {
	header http.Header
	status int
	body   []byte
}

func (r *responseRecorder) Header() http.Header { return r.header }
func (r *responseRecorder) Write(p []byte) (int, error) {
	r.body = append(r.body, p...)
	if r.status == 0 {
		r.status = http.StatusOK
	}
	return len(p), nil
}
func (r *responseRecorder) WriteHeader(statusCode int) { r.status = statusCode }

func loadPanelSettings(manager *dbmanager.DatabaseManager) (map[string]any, error) {
	settings := map[string]any{
		"passkey_settings":  map[string]any{"rpId": nil, "origin": nil, "enabled": false},
		"oauth2_settings":   map[string]any{},
		"tg_auth_settings":  map[string]any{"enabled": false, "adminIds": []any{}, "botToken": nil},
		"password_settings": map[string]any{"enabled": true},
		"branding_settings": map[string]any{"title": "EXODUS", "logoUrl": nil},
		"modules_settings":  map[string]any{"haproxy": map[string]any{"enabled": false}},
	}

	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		row := db.QueryRow(`
			SELECT passkey_settings, oauth2_settings, tg_auth_settings, password_settings, branding_settings, modules_settings
			FROM exodus_settings
			WHERE id = 1
			LIMIT 1
		`)

		var passkeyRaw, oauth2Raw, tgAuthRaw, passwordRaw, brandingRaw, modulesRaw sql.NullString
		if scanErr := row.Scan(&passkeyRaw, &oauth2Raw, &tgAuthRaw, &passwordRaw, &brandingRaw, &modulesRaw); scanErr != nil {
			if errors.Is(scanErr, sql.ErrNoRows) {
				return nil
			}
			return scanErr
		}

		mergeJSONObject(settings, "passkey_settings", passkeyRaw.String)
		mergeJSONObject(settings, "oauth2_settings", oauth2Raw.String)
		mergeJSONObject(settings, "tg_auth_settings", tgAuthRaw.String)
		mergeJSONObject(settings, "password_settings", passwordRaw.String)
		mergeJSONObject(settings, "branding_settings", brandingRaw.String)
		mergeJSONObject(settings, "modules_settings", modulesRaw.String)
		return nil
	})

	return settings, err
}

func loadAPITokens(manager *dbmanager.DatabaseManager) ([]APITokenRecord, error) {
	tokens := make([]APITokenRecord, 0)

	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		rows, queryErr := db.Query(`
			SELECT uuid, token, token_name, created_at, updated_at
			FROM api_tokens
			ORDER BY created_at DESC
		`)
		if queryErr != nil {
			return queryErr
		}
		defer rows.Close()

		for rows.Next() {
			var row APITokenRecord
			if scanErr := rows.Scan(&row.UUID, &row.Token, &row.TokenName, &row.CreatedAt, &row.UpdatedAt); scanErr != nil {
				return scanErr
			}
			tokens = append(tokens, row)
		}
		return rows.Err()
	})

	return tokens, err
}

func mergeJSONObject(target map[string]any, key, raw string) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return
	}

	var decoded map[string]any
	if err := json.Unmarshal([]byte(raw), &decoded); err != nil {
		return
	}

	if existing, ok := target[key].(map[string]any); ok {
		for k, v := range decoded {
			existing[k] = v
		}
		target[key] = existing
		return
	}

	target[key] = decoded
}
