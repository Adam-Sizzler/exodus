package squads

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"

	"exodus/internal/config"
	"exodus/internal/httpapi/shared"

	"github.com/google/uuid"
)

func InternalSquadsHandler(db *sql.DB, cfg *config.BackendConfig) http.HandlerFunc {
	repo := NewSquadRepository(db)
	service := NewSquadService(repo, cfg)
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handleGetInternalSquads(w, r, service)
		case http.MethodPost:
			handleCreateInternalSquad(w, r, service)
		case http.MethodPatch:
			handleUpdateInternalSquad(w, r, service)
		default:
			shared.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
		}
	}
}

func InternalSquadsReorderHandler(db *sql.DB, cfg *config.BackendConfig) http.HandlerFunc {
	repo := NewSquadRepository(db)
	service := NewSquadService(repo, cfg)
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			shared.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}
		handleReorderInternalSquads(w, r, service)
	}
}

func InternalSquadByUUIDHandler(db *sql.DB, cfg *config.BackendConfig) http.HandlerFunc {
	repo := NewSquadRepository(db)
	service := NewSquadService(repo, cfg)
	return func(w http.ResponseWriter, r *http.Request) {
		path := strings.Trim(strings.TrimPrefix(r.URL.Path, "/api/internal-squads/"), "/")
		if path == "" {
			InternalSquadsHandler(db, cfg)(w, r)
			return
		}
		parts := strings.Split(path, "/")
		squadUUID := parts[0]
		if _, err := uuid.Parse(squadUUID); err != nil {
			shared.SendError(w, http.StatusBadRequest, "invalid UUID format", nil, cfg)
			return
		}

		if len(parts) > 1 {
			if len(parts) == 3 && parts[1] == "bulk-actions" && parts[2] == "add-users" {
				if r.Method != http.MethodPost {
					shared.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
					return
				}
				handleBulkAddUsersToInternalSquad(w, r, db, cfg, squadUUID)
				return
			}
			if len(parts) == 3 && parts[1] == "bulk-actions" && parts[2] == "remove-users" {
				if r.Method != http.MethodDelete && r.Method != http.MethodPost {
					shared.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
					return
				}
				handleBulkRemoveUsersFromInternalSquad(w, r, db, cfg, squadUUID)
				return
			}
			if len(parts) == 2 && parts[1] == "accessible-nodes" && r.Method == http.MethodGet {
				handleGetInternalSquadAccessibleNodes(w, r, service, squadUUID)
				return
			}
			http.NotFound(w, r)
			return
		}

		switch r.Method {
		case http.MethodGet:
			handleGetInternalSquad(w, r, service, squadUUID)
		case http.MethodDelete:
			handleDeleteInternalSquad(w, r, service, squadUUID)
		default:
			shared.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
		}
	}
}

func handleBulkAddUsersToInternalSquad(w http.ResponseWriter, r *http.Request, db *sql.DB, cfg *config.BackendConfig, squadUUID string) {
	ctx := r.Context()
	var exists int
	if err := db.QueryRowContext(ctx, `SELECT 1 FROM internal_squads WHERE uuid = $1`, squadUUID).Scan(&exists); err != nil {
		if err == sql.ErrNoRows {
			shared.SendError(w, http.StatusNotFound, "internal squad not found", nil, cfg)
			return
		}
		shared.SendError(w, http.StatusInternalServerError, "failed to check internal squad", err, cfg)
		return
	}

	_, err := db.ExecContext(ctx, `
		INSERT INTO internal_squad_members (internal_squad_uuid, user_id)
		SELECT $1, t_id FROM users
		ON CONFLICT (internal_squad_uuid, user_id) DO NOTHING
	`, squadUUID)
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to add users to internal squad", err, cfg)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(map[string]any{
		"response": map[string]any{
			"eventSent": true,
		},
	})
}

func handleBulkRemoveUsersFromInternalSquad(w http.ResponseWriter, r *http.Request, db *sql.DB, cfg *config.BackendConfig, squadUUID string) {
	ctx := r.Context()
	var exists int
	if err := db.QueryRowContext(ctx, `SELECT 1 FROM internal_squads WHERE uuid = $1`, squadUUID).Scan(&exists); err != nil {
		if err == sql.ErrNoRows {
			shared.SendError(w, http.StatusNotFound, "internal squad not found", nil, cfg)
			return
		}
		shared.SendError(w, http.StatusInternalServerError, "failed to check internal squad", err, cfg)
		return
	}

	_, err := db.ExecContext(ctx, `
		DELETE FROM internal_squad_members WHERE internal_squad_uuid = $1
	`, squadUUID)
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to remove users from internal squad", err, cfg)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(map[string]any{
		"response": map[string]any{
			"eventSent": true,
		},
	})
}

func InboundAssignmentsHandler(db *sql.DB, cfg *config.BackendConfig) http.HandlerFunc {
	repo := NewSquadRepository(db)
	service := NewSquadService(repo, cfg)
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handleGetInboundAssignments(w, r, service)
		case http.MethodPost:
			handleSetInboundAssignments(w, r, service)
		default:
			http.Error(w, `{"error": "method not allowed"}`, http.StatusMethodNotAllowed)
		}
	}
}

func ConfigProfilesWithInboundsHandler(db *sql.DB, cfg *config.BackendConfig) http.HandlerFunc {
	repo := NewSquadRepository(db)
	service := NewSquadService(repo, cfg)
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, `{"error": "method not allowed"}`, http.StatusMethodNotAllowed)
			return
		}
		handleGetConfigProfilesWithInbounds(w, r, service)
	}
}
