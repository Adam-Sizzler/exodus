package squads

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"exodus/internal/config"
	dbmanager "exodus/internal/db/manager"
	"exodus/internal/httpapi/shared"
	monitor "exodus/internal/nodes"

	"github.com/google/uuid"
)

var errInternalSquadNotFound = errors.New("internal squad not found")

// InternalSquad represents an internal squad entity for API responses.
type InternalSquad struct {
	UUID         string    `json:"uuid"`
	ViewPosition int       `json:"view_position"`
	Name         string    `json:"name"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type InternalSquadInfo struct {
	MembersCount  int `json:"membersCount"`
	InboundsCount int `json:"inboundsCount"`
}

type InternalSquadInboundAPI struct {
	UUID        string          `json:"uuid"`
	ProfileUUID string          `json:"profileUuid"`
	Tag         string          `json:"tag"`
	Type        string          `json:"type"`
	Network     *string         `json:"network"`
	Security    *string         `json:"security"`
	Port        *int            `json:"port"`
	RawInbound  json.RawMessage `json:"rawInbound"`
}

type InternalSquadAPI struct {
	UUID         string                    `json:"uuid"`
	ViewPosition int                       `json:"viewPosition"`
	Name         string                    `json:"name"`
	Info         InternalSquadInfo         `json:"info"`
	Inbounds     []InternalSquadInboundAPI `json:"inbounds"`
	CreatedAt    time.Time                 `json:"createdAt"`
	UpdatedAt    time.Time                 `json:"updatedAt"`
}

type InternalSquadAccessibleNode struct {
	UUID              string   `json:"uuid"`
	NodeName          string   `json:"nodeName"`
	CountryCode       string   `json:"countryCode"`
	ConfigProfileUUID string   `json:"configProfileUuid"`
	ConfigProfileName string   `json:"configProfileName"`
	ActiveInbounds    []string `json:"activeInbounds"`
}

// InternalSquadCreateRequest represents a request to create a new internal squad.
type InternalSquadCreateRequest struct {
	ViewPosition int    `json:"viewPosition"`
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
	UUID         string   `json:"uuid,omitempty"`
	ViewPosition *int     `json:"viewPosition,omitempty"`
	Name         *string  `json:"name,omitempty"`
	Inbounds     []string `json:"inbounds,omitempty"`
}

// Validate validates the InternalSquadUpdateRequest fields.
func (r *InternalSquadUpdateRequest) Validate() error {
	if r.Name != nil && *r.Name == "" {
		return fmt.Errorf("name cannot be empty")
	}
	for _, inboundUUID := range r.Inbounds {
		if _, err := uuid.Parse(strings.TrimSpace(inboundUUID)); err != nil {
			return fmt.Errorf("invalid inbound UUID format")
		}
	}
	return nil
}

// HasUpdates checks if any field is set for update.
func (r *InternalSquadUpdateRequest) HasUpdates() bool {
	return r.ViewPosition != nil || r.Name != nil || r.Inbounds != nil
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

// InternalSquadsHandler handles GET/POST /api/internal-squads
func InternalSquadsHandler(manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handleGetInternalSquads(w, r, manager, cfg)
		case http.MethodPost:
			handleCreateInternalSquad(w, r, manager, cfg)
		case http.MethodPatch:
			var req InternalSquadUpdateRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				shared.SendError(w, http.StatusBadRequest, "invalid JSON", err, cfg)
				return
			}
			if _, err := uuid.Parse(strings.TrimSpace(req.UUID)); err != nil {
				shared.SendError(w, http.StatusBadRequest, "invalid UUID format", nil, cfg)
				return
			}
			// Re-marshal and inject body back for shared update path.
			payload, _ := json.Marshal(req)
			r.Body = io.NopCloser(strings.NewReader(string(payload)))
			handlePatchInternalSquad(w, r, manager, cfg, req.UUID)
		default:
			http.Error(w, `{"error": "method not allowed"}`, http.StatusMethodNotAllowed)
		}
	}
}

// InternalSquadsReorderHandler handles POST /api/internal-squads/reorder
func InternalSquadsReorderHandler(manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
			return
		}

		var req struct {
			Items []struct {
				UUID         string `json:"uuid"`
				ViewPosition int    `json:"viewPosition"`
			} `json:"items"`
			OrderedUUIDs []string `json:"ordered_uuids"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			shared.SendError(w, http.StatusBadRequest, "invalid JSON", err, cfg)
			return
		}
		if len(req.Items) == 0 && len(req.OrderedUUIDs) == 0 {
			shared.SendError(w, http.StatusBadRequest, "items cannot be empty", nil, cfg)
			return
		}

		err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
			if len(req.Items) > 0 {
				for _, item := range req.Items {
					if _, err := uuid.Parse(item.UUID); err != nil {
						return fmt.Errorf("invalid UUID format")
					}
					if _, err := db.ExecContext(r.Context(),
						`UPDATE internal_squads SET view_position = ?, updated_at = CURRENT_TIMESTAMP WHERE uuid = ?`,
						item.ViewPosition, item.UUID); err != nil {
						return err
					}
				}
				return nil
			}
			if len(req.OrderedUUIDs) > 0 {
				return shared.ApplyViewPositionReorder(r.Context(), db, "internal_squads", req.OrderedUUIDs, cfg)
			}
			return nil
		})
		if err != nil {
			shared.SendError(w, http.StatusInternalServerError, "failed to reorder internal squads", err, cfg)
			return
		}

		handleGetInternalSquads(w, r, manager, cfg)
	}
}

// handleGetInternalSquads handles GET /api/internal-squads
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

	response := make([]InternalSquadAPI, 0, len(squads))
	for _, squad := range squads {
		item, err := buildInternalSquadResponse(ctx, manager, squad)
		if err != nil {
			shared.SendError(w, http.StatusInternalServerError, "failed to build internal squads response", err, cfg)
			return
		}
		response = append(response, item)
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"response": map[string]interface{}{
			"total":          len(response),
			"internalSquads": response,
		},
	})
}

// handleCreateInternalSquad handles POST /api/internal-squads
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

	squad, getErr := getInternalSquadByUUID(r.Context(), manager, squadUUID)
	if getErr != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to fetch created internal squad", getErr, cfg)
		return
	}
	response, respErr := buildInternalSquadResponse(r.Context(), manager, squad)
	if respErr != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to build created internal squad response", respErr, cfg)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	cfg.Logger.Info("Internal squad created", "uuid", squadUUID, "name", req.Name)
	monitor.RequestNodeDeploy(true)
	json.NewEncoder(w).Encode(map[string]interface{}{"response": response})
}

// InternalSquadByUUIDHandler handles GET/PATCH/DELETE /api/internal-squads/{uuid}
func InternalSquadByUUIDHandler(manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimSpace(trimInternalSquadsPath(r.URL.Path))
		if path == "" {
			switch r.Method {
			case http.MethodGet:
				handleGetInternalSquads(w, r, manager, cfg)
			case http.MethodPost:
				handleCreateInternalSquad(w, r, manager, cfg)
			case http.MethodPatch:
				var req InternalSquadUpdateRequest
				if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
					shared.SendError(w, http.StatusBadRequest, "invalid JSON", err, cfg)
					return
				}
				if _, err := uuid.Parse(strings.TrimSpace(req.UUID)); err != nil {
					shared.SendError(w, http.StatusBadRequest, "invalid UUID format", nil, cfg)
					return
				}
				payload, _ := json.Marshal(req)
				r.Body = io.NopCloser(strings.NewReader(string(payload)))
				handlePatchInternalSquad(w, r, manager, cfg, req.UUID)
			default:
				http.Error(w, `{"error": "method not allowed"}`, http.StatusMethodNotAllowed)
			}
			return
		}

		parts := strings.Split(path, "/")
		squadUUID := strings.TrimSpace(parts[0])
		if _, err := uuid.Parse(squadUUID); err != nil {
			shared.SendError(w, http.StatusBadRequest, "invalid UUID format", nil, cfg)
			return
		}
		if len(parts) > 1 {
			if len(parts) == 2 && parts[1] == "accessible-nodes" {
				if r.Method != http.MethodGet {
					http.Error(w, `{"error": "method not allowed"}`, http.StatusMethodNotAllowed)
					return
				}
				handleGetInternalSquadAccessibleNodes(w, r, manager, cfg, squadUUID)
				return
			}
			if len(parts) == 3 && parts[1] == "bulk-actions" && parts[2] == "add-users" {
				if r.Method != http.MethodPost {
					http.Error(w, `{"error": "method not allowed"}`, http.StatusMethodNotAllowed)
					return
				}
				handleBulkAddUsersToInternalSquad(w, r, manager, cfg, squadUUID)
				return
			}
			if len(parts) == 3 && parts[1] == "bulk-actions" && parts[2] == "remove-users" {
				if r.Method != http.MethodDelete {
					http.Error(w, `{"error": "method not allowed"}`, http.StatusMethodNotAllowed)
					return
				}
				handleBulkRemoveUsersFromInternalSquad(w, r, manager, cfg, squadUUID)
				return
			}
			http.NotFound(w, r)
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

func handleGetInternalSquadAccessibleNodes(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig, squadUUID string) {
	nodes := make([]InternalSquadAccessibleNode, 0)
	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		var exists int
		if err := db.QueryRowContext(r.Context(), `SELECT 1 FROM internal_squads WHERE uuid = ?`, squadUUID).Scan(&exists); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return errInternalSquadNotFound
			}
			return err
		}

		rows, err := db.QueryContext(r.Context(), `
			SELECT
				n.uuid,
				n.name,
				n.country_code,
				cp.uuid,
				cp.name,
				cpi.tag
			FROM internal_squad_inbounds isi
			INNER JOIN config_profile_inbounds cpi ON cpi.uuid = isi.inbound_uuid
			INNER JOIN config_profiles cp ON cp.uuid = cpi.profile_uuid
			INNER JOIN config_profile_inbounds_to_nodes cpin
				ON cpin.config_profile_inbound_uuid = cpi.uuid
			INNER JOIN nodes n
				ON n.uuid = cpin.node_uuid
				AND n.active_config_profile_uuid = cp.uuid
			WHERE isi.internal_squad_uuid = ?
			ORDER BY n.view_position ASC, n.name ASC, cpi.tag ASC
		`, squadUUID)
		if err != nil {
			return err
		}
		defer rows.Close()

		indexByNode := make(map[string]int)
		inboundSeenByNode := make(map[string]map[string]bool)
		for rows.Next() {
			var nodeUUID, nodeName, countryCode, profileUUID, profileName, inboundTag string
			if err := rows.Scan(&nodeUUID, &nodeName, &countryCode, &profileUUID, &profileName, &inboundTag); err != nil {
				return err
			}
			idx, ok := indexByNode[nodeUUID]
			if !ok {
				nodes = append(nodes, InternalSquadAccessibleNode{
					UUID:              nodeUUID,
					NodeName:          nodeName,
					CountryCode:       countryCode,
					ConfigProfileUUID: profileUUID,
					ConfigProfileName: profileName,
					ActiveInbounds:    make([]string, 0),
				})
				idx = len(nodes) - 1
				indexByNode[nodeUUID] = idx
				inboundSeenByNode[nodeUUID] = make(map[string]bool)
			}
			if !inboundSeenByNode[nodeUUID][inboundTag] {
				nodes[idx].ActiveInbounds = append(nodes[idx].ActiveInbounds, inboundTag)
				inboundSeenByNode[nodeUUID][inboundTag] = true
			}
		}
		return rows.Err()
	})
	if err != nil {
		if errors.Is(err, errInternalSquadNotFound) {
			shared.SendError(w, http.StatusNotFound, "internal squad not found", nil, cfg)
			return
		}
		shared.SendError(w, http.StatusInternalServerError, "failed to fetch internal squad accessible nodes", err, cfg)
		return
	}

	shared.WriteJSON(w, http.StatusOK, map[string]any{
		"response": map[string]any{
			"squadUuid":       squadUUID,
			"accessibleNodes": nodes,
		},
	})
}

func handleBulkAddUsersToInternalSquad(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig, squadUUID string) {
	var affected int64
	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		var exists int
		if err := db.QueryRowContext(r.Context(), `SELECT 1 FROM internal_squads WHERE uuid = ?`, squadUUID).Scan(&exists); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return errInternalSquadNotFound
			}
			return err
		}
		result, err := db.ExecContext(r.Context(), `
			INSERT INTO internal_squad_members (internal_squad_uuid, user_id)
			SELECT ?::uuid, t_id
			FROM users
			ON CONFLICT (internal_squad_uuid, user_id) DO NOTHING
		`, squadUUID)
		if err != nil {
			return err
		}
		affected, _ = result.RowsAffected()
		return nil
	})
	if err != nil {
		if errors.Is(err, errInternalSquadNotFound) {
			shared.SendError(w, http.StatusNotFound, "internal squad not found", nil, cfg)
			return
		}
		shared.SendError(w, http.StatusInternalServerError, "failed to add users to internal squad", err, cfg)
		return
	}

	cfg.Logger.Info("Users added to internal squad", "squad_uuid", squadUUID, "affected_rows", affected)
	monitor.RequestNodeDeploy(true)
	shared.WriteJSON(w, http.StatusOK, map[string]any{"response": map[string]any{"eventSent": true}})
}

func handleBulkRemoveUsersFromInternalSquad(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig, squadUUID string) {
	var affected int64
	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		var exists int
		if err := db.QueryRowContext(r.Context(), `SELECT 1 FROM internal_squads WHERE uuid = ?`, squadUUID).Scan(&exists); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return errInternalSquadNotFound
			}
			return err
		}
		result, err := db.ExecContext(r.Context(), `DELETE FROM internal_squad_members WHERE internal_squad_uuid = ?`, squadUUID)
		if err != nil {
			return err
		}
		affected, _ = result.RowsAffected()
		return nil
	})
	if err != nil {
		if errors.Is(err, errInternalSquadNotFound) {
			shared.SendError(w, http.StatusNotFound, "internal squad not found", nil, cfg)
			return
		}
		shared.SendError(w, http.StatusInternalServerError, "failed to remove users from internal squad", err, cfg)
		return
	}

	cfg.Logger.Info("Users removed from internal squad", "squad_uuid", squadUUID, "affected_rows", affected)
	monitor.RequestNodeDeploy(true)
	shared.WriteJSON(w, http.StatusOK, map[string]any{"response": map[string]any{"eventSent": true}})
}

func trimInternalSquadsPath(path string) string {
	for _, prefix := range []string{"/api/internal-squads/"} {
		if strings.HasPrefix(path, prefix) {
			return strings.Trim(strings.TrimPrefix(path, prefix), "/")
		}
	}
	return ""
}

func handleGetInternalSquad(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig, squadUUID string) {
	squad, err := getInternalSquadByUUID(r.Context(), manager, squadUUID)

	if err != nil {
		if err == sql.ErrNoRows {
			shared.SendError(w, http.StatusNotFound, "internal squad not found", nil, cfg)
		} else {
			shared.SendError(w, http.StatusInternalServerError, "failed to fetch internal squad", err, cfg)
		}
		return
	}

	response, respErr := buildInternalSquadResponse(r.Context(), manager, squad)
	if respErr != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to build internal squad response", respErr, cfg)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(map[string]interface{}{"response": response})
}

func handlePatchInternalSquad(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig, squadUUID string) {
	var req InternalSquadUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		shared.SendError(w, http.StatusBadRequest, "invalid JSON", err, cfg)
		return
	}

	effectiveSquadUUID := strings.TrimSpace(squadUUID)
	if effectiveSquadUUID == "" {
		effectiveSquadUUID = strings.TrimSpace(req.UUID)
	}
	if _, err := uuid.Parse(effectiveSquadUUID); err != nil {
		shared.SendError(w, http.StatusBadRequest, "invalid UUID format", nil, cfg)
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

	query := ""
	if len(clauses) > 0 {
		args = append(args, effectiveSquadUUID)
		query = fmt.Sprintf("UPDATE internal_squads SET %s, updated_at = CURRENT_TIMESTAMP WHERE uuid = ?", strings.Join(clauses, ", "))
	}

	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		tx, err := db.BeginTx(r.Context(), nil)
		if err != nil {
			return err
		}

		if len(clauses) > 0 {
			result, err := tx.ExecContext(r.Context(), query, args...)
			if err != nil {
				_ = tx.Rollback()
				return err
			}
			rowsAffected, err := result.RowsAffected()
			if err != nil {
				_ = tx.Rollback()
				return err
			}
			if rowsAffected == 0 {
				_ = tx.Rollback()
				return sql.ErrNoRows
			}
		} else {
			var exists int
			if err := tx.QueryRowContext(r.Context(), `SELECT 1 FROM internal_squads WHERE uuid = ?`, effectiveSquadUUID).Scan(&exists); err != nil {
				_ = tx.Rollback()
				return err
			}
		}

		if req.Inbounds != nil {
			if _, err := tx.ExecContext(r.Context(), `DELETE FROM internal_squad_inbounds WHERE internal_squad_uuid = ?`, effectiveSquadUUID); err != nil {
				_ = tx.Rollback()
				return err
			}

			seen := make(map[string]struct{}, len(req.Inbounds))
			for _, inboundUUID := range req.Inbounds {
				cleanInboundUUID := strings.TrimSpace(inboundUUID)
				if cleanInboundUUID == "" {
					continue
				}
				if _, ok := seen[cleanInboundUUID]; ok {
					continue
				}
				seen[cleanInboundUUID] = struct{}{}

				var inboundExists int
				if err := tx.QueryRowContext(r.Context(), `SELECT 1 FROM config_profile_inbounds WHERE uuid = ?`, cleanInboundUUID).Scan(&inboundExists); err != nil {
					_ = tx.Rollback()
					if err == sql.ErrNoRows {
						return fmt.Errorf("inbound not found")
					}
					return err
				}

				if _, err := tx.ExecContext(
					r.Context(),
					`INSERT INTO internal_squad_inbounds (internal_squad_uuid, inbound_uuid) VALUES (?, ?)`,
					effectiveSquadUUID,
					cleanInboundUUID,
				); err != nil {
					_ = tx.Rollback()
					return err
				}
			}
		}

		if err := tx.Commit(); err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		if err == sql.ErrNoRows {
			shared.SendError(w, http.StatusNotFound, "internal squad not found", nil, cfg)
		} else if strings.Contains(strings.ToLower(err.Error()), "inbound not found") {
			shared.SendError(w, http.StatusBadRequest, "one or more inbounds not found", err, cfg)
		} else {
			if strings.Contains(err.Error(), "UNIQUE constraint failed") {
				shared.SendError(w, http.StatusConflict, "name already exists", err, cfg)
			} else {
				shared.SendError(w, http.StatusInternalServerError, "update failed", err, cfg)
			}
		}
		return
	}

	cfg.Logger.Info(
		"Internal squad updated",
		"uuid", effectiveSquadUUID,
		"name_updated", req.Name != nil,
		"view_position_updated", req.ViewPosition != nil,
		"inbounds_updated", req.Inbounds != nil,
		"inbounds_count", len(req.Inbounds),
	)
	monitor.RequestNodeDeploy(true)
	handleGetInternalSquad(w, r, manager, cfg, effectiveSquadUUID) // Return updated squad
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
	monitor.RequestNodeDeploy(true)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"response": map[string]any{
			"isDeleted": true,
		},
	})
}

func getInternalSquadByUUID(ctx context.Context, manager *dbmanager.DatabaseManager, squadUUID string) (InternalSquad, error) {
	var squad InternalSquad
	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		query := `SELECT uuid, view_position, name, created_at, updated_at
				  FROM internal_squads WHERE uuid = ?`
		row := db.QueryRowContext(ctx, query, squadUUID)
		var scanErr error
		squad, scanErr = scanInternalSquad(row)
		return scanErr
	})
	return squad, err
}

func buildInternalSquadResponse(ctx context.Context, manager *dbmanager.DatabaseManager, squad InternalSquad) (InternalSquadAPI, error) {
	response := InternalSquadAPI{
		UUID:         squad.UUID,
		ViewPosition: squad.ViewPosition,
		Name:         squad.Name,
		Info: InternalSquadInfo{
			MembersCount:  0,
			InboundsCount: 0,
		},
		Inbounds:  make([]InternalSquadInboundAPI, 0),
		CreatedAt: squad.CreatedAt,
		UpdatedAt: squad.UpdatedAt,
	}

	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		if err := db.QueryRowContext(ctx,
			`SELECT COUNT(*) FROM internal_squad_members WHERE internal_squad_uuid = ?`,
			squad.UUID).Scan(&response.Info.MembersCount); err != nil {
			return err
		}

		rows, err := db.QueryContext(ctx, `
			SELECT cpi.uuid, cpi.profile_uuid, cpi.tag, cpi.type, cpi.network, cpi.security, cpi.port, cpi.raw_inbound
			FROM internal_squad_inbounds isi
			JOIN config_profile_inbounds cpi ON cpi.uuid = isi.inbound_uuid
			WHERE isi.internal_squad_uuid = ?
			ORDER BY cpi.tag ASC
		`, squad.UUID)
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			var inbound InternalSquadInboundAPI
			var network, security, rawInbound sql.NullString
			var port sql.NullInt64
			if scanErr := rows.Scan(
				&inbound.UUID,
				&inbound.ProfileUUID,
				&inbound.Tag,
				&inbound.Type,
				&network,
				&security,
				&port,
				&rawInbound,
			); scanErr != nil {
				return scanErr
			}
			if network.Valid {
				inbound.Network = &network.String
			}
			if security.Valid {
				inbound.Security = &security.String
			}
			if port.Valid {
				p := int(port.Int64)
				inbound.Port = &p
			}
			if rawInbound.Valid {
				inbound.RawInbound = json.RawMessage(rawInbound.String)
			}
			response.Inbounds = append(response.Inbounds, inbound)
		}
		return rows.Err()
	})

	response.Info.InboundsCount = len(response.Inbounds)
	return response, err
}
