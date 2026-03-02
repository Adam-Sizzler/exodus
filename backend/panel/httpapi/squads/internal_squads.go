package squads

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"v2ray-stat/backend/panel/config"
	dbmanager "v2ray-stat/backend/panel/db/manager"
	"v2ray-stat/backend/panel/httpapi/shared"

	"github.com/google/uuid"
)

// InternalSquad represents an internal squad entity for API responses.
type InternalSquad struct {
	UUID         string    `json:"uuid"`
	ViewPosition int       `json:"view_position"`
	Name         string    `json:"name"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// InternalSquadCreateRequest represents a request to create a new internal squad.
type InternalSquadCreateRequest struct {
	ViewPosition int    `json:"view_position"`
	Name         string `json:"name"`
}

// Validate validates the InternalSquadCreateRequest fields.
func (r *InternalSquadCreateRequest) Validate() error {
	if r.Name == "" {
		return fmt.Errorf("name is required")
	}
	return nil
}

// InternalSquadUpdateRequest represents a partial update request for an internal squad.
// Only provided fields will be updated (PATCH semantics).
type InternalSquadUpdateRequest struct {
	ViewPosition *int    `json:"view_position,omitempty"`
	Name         *string `json:"name,omitempty"`
}

// Validate validates the InternalSquadUpdateRequest fields.
func (r *InternalSquadUpdateRequest) Validate() error {
	if r.Name != nil && *r.Name == "" {
		return fmt.Errorf("name cannot be empty")
	}
	return nil
}

// HasUpdates checks if any field is set for update.
func (r *InternalSquadUpdateRequest) HasUpdates() bool {
	return r.ViewPosition != nil || r.Name != nil
}

// scanInternalSquad scans a row into an InternalSquad struct.
func scanInternalSquad(scanner shared.RowScanner) (InternalSquad, error) {
	var squad InternalSquad
	var viewPosition sql.NullInt64

	err := scanner.Scan(
		&squad.UUID,
		&viewPosition,
		&squad.Name,
		&squad.CreatedAt,
		&squad.UpdatedAt,
	)
	if err != nil {
		return squad, err
	}

	if viewPosition.Valid {
		squad.ViewPosition = int(viewPosition.Int64)
	}

	return squad, nil
}

// InternalSquadsHandler handles GET/POST /api/v1/internal-squads
func InternalSquadsHandler(manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handleGetInternalSquads(w, r, manager, cfg)
		case http.MethodPost:
			handleCreateInternalSquad(w, r, manager, cfg)
		default:
			http.Error(w, `{"error": "method not allowed"}`, http.StatusMethodNotAllowed)
		}
	}
}

// InternalSquadsReorderHandler handles POST /api/v1/internal-squads/reorder
func InternalSquadsReorderHandler(manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
			return
		}

		var req shared.ViewPositionReorderRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			shared.SendError(w, http.StatusBadRequest, "invalid JSON", err, cfg)
			return
		}
		if err := req.Validate(); err != nil {
			shared.SendError(w, http.StatusBadRequest, err.Error(), nil, cfg)
			return
		}

		err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
			return shared.ApplyViewPositionReorder(r.Context(), db, "internal_squads", req.OrderedUUIDs, cfg)
		})
		if err != nil {
			shared.SendError(w, http.StatusInternalServerError, "failed to reorder internal squads", err, cfg)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "internal squads reordered",
			"count":   len(req.OrderedUUIDs),
		})
	}
}

// handleGetInternalSquads handles GET /api/v1/internal-squads
func handleGetInternalSquads(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) {
	ctx := r.Context()
	var squads []InternalSquad

	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		query := `
			SELECT uuid, view_position, name, created_at, updated_at
			FROM internal_squads
			ORDER BY view_position ASC, name ASC`

		rows, err := db.QueryContext(ctx, query)
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			squad, err := scanInternalSquad(rows)
			if err != nil {
				return err
			}
			squads = append(squads, squad)
		}
		return rows.Err()
	})

	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to fetch internal squads", err, cfg)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"squads": squads,
		"count":  len(squads),
	})
}

// handleCreateInternalSquad handles POST /api/v1/internal-squads
func handleCreateInternalSquad(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) {
	var req InternalSquadCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		shared.SendError(w, http.StatusBadRequest, "invalid JSON", err, cfg)
		return
	}

	// Validate required fields
	if err := req.Validate(); err != nil {
		shared.SendError(w, http.StatusBadRequest, err.Error(), nil, cfg)
		return
	}

	// Generate UUID
	squadUUID := uuid.New().String()

	ctx := r.Context()
	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		query := `
			INSERT INTO internal_squads (
				uuid, view_position, name, created_at, updated_at
			) VALUES (?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)`

		_, err := db.ExecContext(ctx, query,
			squadUUID, req.ViewPosition, req.Name)

		return err
	})

	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			shared.SendError(w, http.StatusConflict, "name already exists", err, cfg)
			return
		}
		shared.SendError(w, http.StatusInternalServerError, "failed to create internal squad", err, cfg)
		return
	}

	// Return created squad UUID
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "internal squad created",
		"uuid":    squadUUID,
	})
}

// InternalSquadByUUIDHandler handles GET/PATCH/DELETE /api/v1/internal-squads/{uuid}
func InternalSquadByUUIDHandler(manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/api/v1/internal-squads/")
		squadUUID := strings.TrimSpace(path)

		if _, err := uuid.Parse(squadUUID); err != nil {
			shared.SendError(w, http.StatusBadRequest, "invalid UUID format", nil, cfg)
			return
		}

		switch r.Method {
		case http.MethodGet:
			handleGetInternalSquad(w, r, manager, cfg, squadUUID)
		case http.MethodPatch:
			handlePatchInternalSquad(w, r, manager, cfg, squadUUID)
		case http.MethodDelete:
			handleDeleteInternalSquad(w, r, manager, cfg, squadUUID)
		default:
			http.Error(w, `{"error": "method not allowed"}`, http.StatusMethodNotAllowed)
		}
	}
}

func handleGetInternalSquad(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig, squadUUID string) {
	var squad InternalSquad
	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		query := `SELECT uuid, view_position, name, created_at, updated_at
				  FROM internal_squads WHERE uuid = ?`
		row := db.QueryRowContext(r.Context(), query, squadUUID)
		var scanErr error
		squad, scanErr = scanInternalSquad(row)
		return scanErr
	})

	if err != nil {
		if err == sql.ErrNoRows {
			shared.SendError(w, http.StatusNotFound, "internal squad not found", nil, cfg)
		} else {
			shared.SendError(w, http.StatusInternalServerError, "failed to fetch internal squad", err, cfg)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(map[string]interface{}{"squad": squad})
}

func handlePatchInternalSquad(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig, squadUUID string) {
	var req InternalSquadUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		shared.SendError(w, http.StatusBadRequest, "invalid JSON", err, cfg)
		return
	}

	if err := req.Validate(); err != nil {
		shared.SendError(w, http.StatusBadRequest, err.Error(), nil, cfg)
		return
	}

	if !req.HasUpdates() {
		shared.SendError(w, http.StatusBadRequest, "no fields to update", nil, cfg)
		return
	}

	var clauses []string
	var args []interface{}

	// Dynamic clause building
	add := func(col string, val interface{}) {
		clauses = append(clauses, fmt.Sprintf("%s = ?", col))
		args = append(args, val)
	}

	if req.ViewPosition != nil {
		add("view_position", *req.ViewPosition)
	}
	if req.Name != nil {
		add("name", *req.Name)
	}

	if len(clauses) == 0 {
		shared.SendError(w, http.StatusBadRequest, "no fields to update", nil, cfg)
		return
	}

	args = append(args, squadUUID)
	query := fmt.Sprintf("UPDATE internal_squads SET %s, updated_at = CURRENT_TIMESTAMP WHERE uuid = ?", strings.Join(clauses, ", "))

	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		result, err := db.ExecContext(r.Context(), query, args...)
		if err != nil {
			return err
		}
		rowsAffected, err := result.RowsAffected()
		if err != nil {
			return err
		}
		if rowsAffected == 0 {
			return sql.ErrNoRows
		}
		return nil
	})

	if err != nil {
		if err == sql.ErrNoRows {
			shared.SendError(w, http.StatusNotFound, "internal squad not found", nil, cfg)
		} else {
			if strings.Contains(err.Error(), "UNIQUE constraint failed") {
				shared.SendError(w, http.StatusConflict, "name already exists", err, cfg)
			} else {
				shared.SendError(w, http.StatusInternalServerError, "update failed", err, cfg)
			}
		}
		return
	}

	handleGetInternalSquad(w, r, manager, cfg, squadUUID) // Return updated squad
}

func handleDeleteInternalSquad(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig, squadUUID string) {
	ctx := r.Context()

	// Get squad name for logging
	var squadName string
	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		return db.QueryRowContext(ctx, "SELECT name FROM internal_squads WHERE uuid = ?", squadUUID).Scan(&squadName)
	})

	if err != nil {
		if err == sql.ErrNoRows {
			shared.SendError(w, http.StatusNotFound, "internal squad not found", nil, cfg)
		} else {
			shared.SendError(w, http.StatusInternalServerError, "failed to find internal squad", err, cfg)
		}
		return
	}

	err = manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		_, err := db.ExecContext(ctx, "DELETE FROM internal_squads WHERE uuid = ?", squadUUID)
		return err
	})

	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to delete internal squad", err, cfg)
		return
	}

	cfg.Logger.Info("Internal squad deleted", "uuid", squadUUID, "name", squadName)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "internal squad deleted",
		"uuid":    squadUUID,
		"name":    squadName,
	})
}
