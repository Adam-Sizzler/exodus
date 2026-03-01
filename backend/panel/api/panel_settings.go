package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"v2ray-stat/backend/panel/config"
	dbmanager "v2ray-stat/backend/panel/db/manager"
	"v2ray-stat/backend/panel/security"

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
				writeJSONError(w, http.StatusInternalServerError, "failed to load panel settings")
				return
			}
			writeJSON(w, http.StatusOK, PanelSettingsResponse{Settings: settings})
		case http.MethodPatch, http.MethodPut:
			var payload map[string]json.RawMessage
			if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
				writeJSONError(w, http.StatusBadRequest, "invalid JSON body")
				return
			}

			validKeys := map[string]string{
				"passkey_settings":  "passkey_settings",
				"oauth2_settings":   "oauth2_settings",
				"tg_auth_settings":  "tg_auth_settings",
				"password_settings": "password_settings",
				"branding_settings": "branding_settings",
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
					writeJSONError(w, http.StatusBadRequest, "invalid JSON value for "+key)
					return
				}
				setClauses = append(setClauses, column+" = ?")
				args = append(args, string(raw))
			}

			if len(setClauses) == 0 {
				writeJSONError(w, http.StatusBadRequest, "no valid fields provided")
				return
			}

			err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
				if _, execErr := db.Exec(`
					INSERT INTO v2rs_settings (
						id, passkey_settings, oauth2_settings, tg_auth_settings, password_settings, branding_settings
					) VALUES (
						1, ?, ?, ?, ?, ?
					)
					ON CONFLICT (id) DO NOTHING
				`,
					`{"rpId":null,"origin":null,"enabled":false}`,
					`{"github":{"enabled":false,"clientId":null,"clientSecret":null,"allowedEmails":[]},"yandex":{"enabled":false,"clientId":null,"clientSecret":null,"allowedEmails":[]},"generic":{"enabled":false,"clientId":null,"tokenUrl":null,"withPkce":false,"clientSecret":null,"allowedEmails":[],"frontendDomain":null,"authorizationUrl":null},"keycloak":{"realm":null,"enabled":false,"clientId":null,"clientSecret":null,"allowedEmails":[],"frontendDomain":null,"keycloakDomain":null},"pocketid":{"enabled":false,"clientId":null,"plainDomain":null,"clientSecret":null,"allowedEmails":[]}}`,
					`{"enabled":false,"adminIds":[],"botToken":null}`,
					`{"enabled":true}`,
					`{"title":"V2RS","logoUrl":null}`,
				); execErr != nil {
					return execErr
				}

				query := "UPDATE v2rs_settings SET " + strings.Join(setClauses, ", ") + " WHERE id = 1"
				_, execErr := db.Exec(query, args...)
				return execErr
			})
			if err != nil {
				cfg.Logger.Error("Failed to update panel settings", "error", err)
				writeJSONError(w, http.StatusInternalServerError, "failed to update panel settings")
				return
			}

			settings, err := loadPanelSettings(manager)
			if err != nil {
				cfg.Logger.Error("Failed to load panel settings after update", "error", err)
				writeJSONError(w, http.StatusInternalServerError, "failed to load updated panel settings")
				return
			}
			writeJSON(w, http.StatusOK, PanelSettingsResponse{Settings: settings})
		default:
			writeJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
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
				writeJSONError(w, http.StatusInternalServerError, "failed to load api tokens")
				return
			}
			writeJSON(w, http.StatusOK, map[string]any{
				"tokens": tokens,
				"count":  len(tokens),
			})
		case http.MethodPost:
			var payload struct {
				TokenName string `json:"token_name"`
			}
			if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
				writeJSONError(w, http.StatusBadRequest, "invalid JSON body")
				return
			}
			tokenName := strings.TrimSpace(payload.TokenName)
			if tokenName == "" {
				writeJSONError(w, http.StatusBadRequest, "token_name is required")
				return
			}

			tokenValue, err := security.GenerateRandomToken(64)
			if err != nil {
				cfg.Logger.Error("Failed to generate api token", "error", err)
				writeJSONError(w, http.StatusInternalServerError, "failed to generate api token")
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
				writeJSONError(w, http.StatusInternalServerError, "failed to save api token")
				return
			}

			// Return full token only in create response.
			writeJSON(w, http.StatusCreated, map[string]any{
				"token":   record,
				"message": "api token created",
			})
		default:
			writeJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
		}
	}
}

func PanelAPITokenByUUIDHandler(manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			writeJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}

		tokenUUID := strings.TrimSpace(strings.TrimPrefix(r.URL.Path, "/api/v1/panel/api-tokens/"))
		if _, err := uuid.Parse(tokenUUID); err != nil {
			writeJSONError(w, http.StatusBadRequest, "invalid token UUID")
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
			writeJSONError(w, http.StatusInternalServerError, "failed to delete api token")
			return
		}
		if rowsAffected == 0 {
			writeJSONError(w, http.StatusNotFound, "api token not found")
			return
		}

		writeJSON(w, http.StatusOK, map[string]any{
			"message": "api token deleted",
			"uuid":    tokenUUID,
		})
	}
}

func loadPanelSettings(manager *dbmanager.DatabaseManager) (map[string]any, error) {
	settings := map[string]any{
		"passkey_settings":  map[string]any{"rpId": nil, "origin": nil, "enabled": false},
		"oauth2_settings":   map[string]any{},
		"tg_auth_settings":  map[string]any{"enabled": false, "adminIds": []any{}, "botToken": nil},
		"password_settings": map[string]any{"enabled": true},
		"branding_settings": map[string]any{"title": "V2RS", "logoUrl": nil},
	}

	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		row := db.QueryRow(`
			SELECT passkey_settings, oauth2_settings, tg_auth_settings, password_settings, branding_settings
			FROM v2rs_settings
			WHERE id = 1
			LIMIT 1
		`)

		var passkeyRaw, oauth2Raw, tgAuthRaw, passwordRaw, brandingRaw sql.NullString
		if scanErr := row.Scan(&passkeyRaw, &oauth2Raw, &tgAuthRaw, &passwordRaw, &brandingRaw); scanErr != nil {
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
