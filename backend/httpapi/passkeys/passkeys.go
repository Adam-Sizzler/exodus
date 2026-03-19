package passkeys

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"regexp"
	"strings"
	"time"

	"cerberus/backend/config"
	dbmanager "cerberus/backend/db/manager"
	"cerberus/backend/httpapi/auth"
	"cerberus/backend/httpapi/shared"
)

var passkeyNameRegexp = regexp.MustCompile(`^[A-Za-z0-9_\s-]+$`)

type passkeyRecord struct {
	ID         string
	Name       string
	CreatedAt  time.Time
	LastUsedAt time.Time
}

type passkeyWriteRequest struct {
	ID   string `json:"id"`
	Name string `json:"name,omitempty"`
}

func PasskeysHandler(manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handleGetPasskeys(w, r, manager, cfg)
		case http.MethodPatch:
			handlePatchPasskey(w, r, manager, cfg)
		case http.MethodDelete:
			handleDeletePasskey(w, r, manager, cfg)
		default:
			shared.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
		}
	}
}

func handleGetPasskeys(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) {
	adminUUID, ok := currentAdminUUID(r)
	if !ok {
		shared.WriteJSONError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	items, err := listPasskeysForAdmin(r.Context(), manager, adminUUID)
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to fetch passkeys", err, cfg)
		return
	}
	shared.WriteJSON(w, http.StatusOK, map[string]any{
		"response": map[string]any{
			"passkeys": items,
		},
	})
}

func handlePatchPasskey(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) {
	adminUUID, ok := currentAdminUUID(r)
	if !ok {
		shared.WriteJSONError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req passkeyWriteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		shared.SendError(w, http.StatusBadRequest, "invalid JSON", err, cfg)
		return
	}
	req.ID = strings.TrimSpace(req.ID)
	req.Name = strings.TrimSpace(req.Name)
	if req.ID == "" {
		shared.SendError(w, http.StatusBadRequest, "id is required", nil, cfg)
		return
	}
	if len(req.Name) < 2 || len(req.Name) > 30 || !passkeyNameRegexp.MatchString(req.Name) {
		shared.SendError(w, http.StatusBadRequest, "invalid passkey name", nil, cfg)
		return
	}

	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		result, execErr := db.ExecContext(r.Context(), `
			UPDATE passkeys
			SET passkey_provider = ?, updated_at = CURRENT_TIMESTAMP
			WHERE id = ? AND admin_uuid = ?
		`, req.Name, req.ID, adminUUID)
		if execErr != nil {
			return execErr
		}
		rows, rowsErr := result.RowsAffected()
		if rowsErr != nil {
			return rowsErr
		}
		if rows == 0 {
			return sql.ErrNoRows
		}
		return nil
	})
	if err != nil {
		if err == sql.ErrNoRows {
			shared.SendError(w, http.StatusNotFound, "passkey not found", nil, cfg)
			return
		}
		shared.SendError(w, http.StatusInternalServerError, "failed to update passkey", err, cfg)
		return
	}

	items, err := listPasskeysForAdmin(r.Context(), manager, adminUUID)
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to fetch passkeys", err, cfg)
		return
	}
	shared.WriteJSON(w, http.StatusOK, map[string]any{"response": map[string]any{"passkeys": items}})
}

func handleDeletePasskey(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) {
	adminUUID, ok := currentAdminUUID(r)
	if !ok {
		shared.WriteJSONError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req passkeyWriteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		shared.SendError(w, http.StatusBadRequest, "invalid JSON", err, cfg)
		return
	}
	req.ID = strings.TrimSpace(req.ID)
	if req.ID == "" {
		shared.SendError(w, http.StatusBadRequest, "id is required", nil, cfg)
		return
	}

	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		result, execErr := db.ExecContext(r.Context(), `DELETE FROM passkeys WHERE id = ? AND admin_uuid = ?`, req.ID, adminUUID)
		if execErr != nil {
			return execErr
		}
		rows, rowsErr := result.RowsAffected()
		if rowsErr != nil {
			return rowsErr
		}
		if rows == 0 {
			return sql.ErrNoRows
		}
		return nil
	})
	if err != nil {
		if err == sql.ErrNoRows {
			shared.SendError(w, http.StatusNotFound, "passkey not found", nil, cfg)
			return
		}
		shared.SendError(w, http.StatusInternalServerError, "failed to delete passkey", err, cfg)
		return
	}

	items, err := listPasskeysForAdmin(r.Context(), manager, adminUUID)
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to fetch passkeys", err, cfg)
		return
	}
	shared.WriteJSON(w, http.StatusOK, map[string]any{"response": map[string]any{"passkeys": items}})
}

func listPasskeysForAdmin(ctx context.Context, manager *dbmanager.DatabaseManager, adminUUID string) ([]passkeyRecord, error) {
	records := make([]passkeyRecord, 0)
	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		rows, err := db.QueryContext(ctx, `
			SELECT id,
			       COALESCE(NULLIF(passkey_provider, ''), id) AS name,
			       created_at,
			       updated_at
			FROM passkeys
			WHERE admin_uuid = ?
			ORDER BY created_at DESC
		`, adminUUID)
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			var item passkeyRecord
			if scanErr := rows.Scan(&item.ID, &item.Name, &item.CreatedAt, &item.LastUsedAt); scanErr != nil {
				return scanErr
			}
			records = append(records, item)
		}
		return rows.Err()
	})
	return records, err
}

func currentAdminUUID(r *http.Request) (string, bool) {
	principal, ok := auth.CurrentAuthPrincipal(r.Context())
	if !ok || principal == nil {
		return "", false
	}
	adminUUID := strings.TrimSpace(principal.AdminUUID)
	if adminUUID == "" {
		return "", false
	}
	return adminUUID, true
}
