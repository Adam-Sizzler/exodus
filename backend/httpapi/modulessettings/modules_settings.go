package modulessettings

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"exodus/backend/config"
	dbmanager "exodus/backend/db/manager"
	"exodus/backend/httpapi/shared"
	monitor "exodus/backend/nodes"
)

type modulesSettingsResponse struct {
	Response modulesSettingsPayload `json:"response"`
}

type modulesSettingsPayload struct {
	Haproxy modulesHaproxySettings `json:"haproxy"`
}

type modulesHaproxySettings struct {
	Enabled bool `json:"enabled"`
}

func ModulesSettingsHandler(manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			settings, err := loadModulesSettings(manager)
			if err != nil {
				cfg.Logger.Error("Failed to load modules settings", "error", err)
				shared.WriteJSONError(w, http.StatusInternalServerError, "failed to load modules settings")
				return
			}
			shared.WriteJSON(w, http.StatusOK, modulesSettingsResponse{Response: settings})
		case http.MethodPatch, http.MethodPut:
			var payload modulesSettingsPayload
			if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
				shared.WriteJSONError(w, http.StatusBadRequest, "invalid JSON body")
				return
			}

			err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
				if _, execErr := db.Exec(`
					INSERT INTO modules_settings (id, haproxy_enabled)
					VALUES (1, false)
					ON CONFLICT (id) DO NOTHING
				`); execErr != nil {
					return execErr
				}

				_, execErr := db.Exec(`
					UPDATE modules_settings
					SET haproxy_enabled = ?, updated_at = CURRENT_TIMESTAMP
					WHERE id = 1
				`, payload.Haproxy.Enabled)
				return execErr
			})
			if err != nil {
				cfg.Logger.Error("Failed to update modules settings", "error", err)
				shared.WriteJSONError(w, http.StatusInternalServerError, "failed to update modules settings")
				return
			}

			settings, err := loadModulesSettings(manager)
			if err != nil {
				cfg.Logger.Error("Failed to load modules settings after update", "error", err)
				shared.WriteJSONError(w, http.StatusInternalServerError, "failed to load updated modules settings")
				return
			}

			// Apply module setting changes to connected nodes immediately.
			monitor.RequestNodeDeploy(true)
			shared.WriteJSON(w, http.StatusOK, modulesSettingsResponse{Response: settings})
		default:
			shared.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
		}
	}
}

func loadModulesSettings(manager *dbmanager.DatabaseManager) (modulesSettingsPayload, error) {
	settings := modulesSettingsPayload{
		Haproxy: modulesHaproxySettings{Enabled: false},
	}

	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		if _, execErr := db.Exec(`
			INSERT INTO modules_settings (id, haproxy_enabled)
			VALUES (1, false)
			ON CONFLICT (id) DO NOTHING
		`); execErr != nil {
			return execErr
		}

		row := db.QueryRow(`SELECT haproxy_enabled FROM modules_settings WHERE id = 1 LIMIT 1`)
		if scanErr := row.Scan(&settings.Haproxy.Enabled); scanErr != nil {
			if errors.Is(scanErr, sql.ErrNoRows) {
				return nil
			}
			return scanErr
		}

		return nil
	})

	return settings, err
}
