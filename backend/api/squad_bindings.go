package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"v2ray-stat/backend/config"
	"v2ray-stat/backend/db/manager"

	"github.com/google/uuid"
)

// ==================== CONFIG PROFILE INBOUNDS TO NODES ====================

// ConfigProfileInboundToNode represents a binding between an inbound and a node.
type ConfigProfileInboundToNode struct {
	ConfigProfileInboundUUID string    `json:"config_profile_inbound_uuid"`
	NodeUUID                 string    `json:"node_uuid"`
	CreatedAt                time.Time `json:"created_at"`
}

// InboundAssignmentRequest represents a request to assign inbounds to a node.
type InboundAssignmentRequest struct {
	NodeUUID   string   `json:"node_uuid"`
	InboundUUIDs []string `json:"inbound_uuids"`
}

// Validate validates the InboundAssignmentRequest.
func (r *InboundAssignmentRequest) Validate() error {
	if r.NodeUUID == "" {
		return fmt.Errorf("node_uuid is required")
	}
	if _, err := uuid.Parse(r.NodeUUID); err != nil {
		return fmt.Errorf("invalid node_uuid format")
	}
	for _, inboundUUID := range r.InboundUUIDs {
		if _, err := uuid.Parse(inboundUUID); err != nil {
			return fmt.Errorf("invalid inbound_uuid format: %s", inboundUUID)
		}
	}
	return nil
}

// InboundAssignmentsHandler handles GET/POST /api/v1/inbound-assignments
func InboundAssignmentsHandler(manager *manager.DatabaseManager, cfg *config.BackendConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handleGetInboundAssignments(w, r, manager, cfg)
		case http.MethodPost:
			handleSetInboundAssignments(w, r, manager, cfg)
		default:
			http.Error(w, `{"error": "method not allowed"}`, http.StatusMethodNotAllowed)
		}
	}
}

// handleGetInboundAssignments handles GET /api/v1/inbound-assignments?node_uuid={uuid}
func handleGetInboundAssignments(w http.ResponseWriter, r *http.Request, manager *manager.DatabaseManager, cfg *config.BackendConfig) {
	ctx := r.Context()
	nodeUUID := r.URL.Query().Get("node_uuid")

	var assignments []ConfigProfileInboundToNode

	err := manager.ExecuteHighPriority(func(db *sql.DB) error {
		var query string
		var args []interface{}

		if nodeUUID != "" {
			query = `
				SELECT config_profile_inbound_uuid, node_uuid, created_at
				FROM config_profile_inbounds_to_nodes
				WHERE node_uuid = ?
				ORDER BY config_profile_inbound_uuid`
			args = []interface{}{nodeUUID}
		} else {
			query = `
				SELECT config_profile_inbound_uuid, node_uuid, created_at
				FROM config_profile_inbounds_to_nodes
				ORDER BY node_uuid, config_profile_inbound_uuid`
		}

		rows, err := db.QueryContext(ctx, query, args...)
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			var a ConfigProfileInboundToNode
			if err := rows.Scan(&a.ConfigProfileInboundUUID, &a.NodeUUID, &a.CreatedAt); err != nil {
				return err
			}
			assignments = append(assignments, a)
		}
		return rows.Err()
	})

	if err != nil {
		sendError(w, http.StatusInternalServerError, "failed to fetch inbound assignments", err, cfg)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"assignments": assignments,
		"count":       len(assignments),
	})
}

// handleSetInboundAssignments handles POST /api/v1/inbound-assignments
// Replaces all inbounds for a node with the provided list.
func handleSetInboundAssignments(w http.ResponseWriter, r *http.Request, manager *manager.DatabaseManager, cfg *config.BackendConfig) {
	var req InboundAssignmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, http.StatusBadRequest, "invalid JSON", err, cfg)
		return
	}

	if err := req.Validate(); err != nil {
		sendError(w, http.StatusBadRequest, err.Error(), nil, cfg)
		return
	}

	ctx := r.Context()
	err := manager.ExecuteHighPriority(func(db *sql.DB) error {
		// Verify node exists
		var nodeID int64
		err := db.QueryRowContext(ctx, "SELECT id FROM nodes WHERE uuid = ?", req.NodeUUID).Scan(&nodeID)
		if err != nil {
			if err == sql.ErrNoRows {
				return fmt.Errorf("node not found")
			}
			return err
		}

		// Delete existing assignments for this node
		_, err = db.ExecContext(ctx, "DELETE FROM config_profile_inbounds_to_nodes WHERE node_uuid = ?", req.NodeUUID)
		if err != nil {
			return fmt.Errorf("failed to clear existing assignments: %w", err)
		}

		// Insert new assignments
		for _, inboundUUID := range req.InboundUUIDs {
			// Verify inbound exists
			var inboundID string
			err := db.QueryRowContext(ctx, "SELECT uuid FROM config_profile_inbounds WHERE uuid = ?", inboundUUID).Scan(&inboundID)
			if err != nil {
				if err == sql.ErrNoRows {
					return fmt.Errorf("inbound not found: %s", inboundUUID)
				}
				return err
			}

			_, err = db.ExecContext(ctx,
				"INSERT INTO config_profile_inbounds_to_nodes (config_profile_inbound_uuid, node_uuid) VALUES (?, ?)",
				inboundUUID, req.NodeUUID)
			if err != nil {
				return fmt.Errorf("failed to insert assignment: %w", err)
			}
		}

		return nil
	})

	if err != nil {
		if err.Error() == "node not found" {
			sendError(w, http.StatusNotFound, "node not found", err, cfg)
			return
		}
		if strings.Contains(err.Error(), "inbound not found") {
			sendError(w, http.StatusBadRequest, err.Error(), err, cfg)
			return
		}
		sendError(w, http.StatusInternalServerError, "failed to set inbound assignments", err, cfg)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":       "inbound assignments updated",
		"node_uuid":     req.NodeUUID,
		"inbounds_count": len(req.InboundUUIDs),
	})
}

// ==================== INTERNAL SQUAD INBOUNDS ====================

// InternalSquadInbound represents a binding between a squad and an inbound.
type InternalSquadInbound struct {
	InternalSquadUUID string    `json:"internal_squad_uuid"`
	InboundUUID       string    `json:"inbound_uuid"`
}

// SquadInboundsRequest represents a request to set inbounds for a squad.
type SquadInboundsRequest struct {
	SquadUUID   string   `json:"squad_uuid"`
	InboundUUIDs []string `json:"inbound_uuids"`
}

// Validate validates the SquadInboundsRequest.
func (r *SquadInboundsRequest) Validate() error {
	if r.SquadUUID == "" {
		return fmt.Errorf("squad_uuid is required")
	}
	if _, err := uuid.Parse(r.SquadUUID); err != nil {
		return fmt.Errorf("invalid squad_uuid format")
	}
	for _, inboundUUID := range r.InboundUUIDs {
		if _, err := uuid.Parse(inboundUUID); err != nil {
			return fmt.Errorf("invalid inbound_uuid format: %s", inboundUUID)
		}
	}
	return nil
}

// SquadInboundsHandler handles GET/POST /api/v1/squad-inbounds
func SquadInboundsHandler(manager *manager.DatabaseManager, cfg *config.BackendConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handleGetSquadInbounds(w, r, manager, cfg)
		case http.MethodPost:
			handleSetSquadInbounds(w, r, manager, cfg)
		default:
			http.Error(w, `{"error": "method not allowed"}`, http.StatusMethodNotAllowed)
		}
	}
}

// handleGetSquadInbounds handles GET /api/v1/squad-inbounds?squad_uuid={uuid}
func handleGetSquadInbounds(w http.ResponseWriter, r *http.Request, manager *manager.DatabaseManager, cfg *config.BackendConfig) {
	ctx := r.Context()
	squadUUID := r.URL.Query().Get("squad_uuid")

	var squadInbounds []InternalSquadInbound

	err := manager.ExecuteHighPriority(func(db *sql.DB) error {
		var query string
		var args []interface{}

		if squadUUID != "" {
			query = `
				SELECT internal_squad_uuid, inbound_uuid
				FROM internal_squad_inbounds
				WHERE internal_squad_uuid = ?
				ORDER BY inbound_uuid`
			args = []interface{}{squadUUID}
		} else {
			query = `
				SELECT internal_squad_uuid, inbound_uuid
				FROM internal_squad_inbounds
				ORDER BY internal_squad_uuid, inbound_uuid`
		}

		rows, err := db.QueryContext(ctx, query, args...)
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			var si InternalSquadInbound
			if err := rows.Scan(&si.InternalSquadUUID, &si.InboundUUID); err != nil {
				return err
			}
			squadInbounds = append(squadInbounds, si)
		}
		return rows.Err()
	})

	if err != nil {
		sendError(w, http.StatusInternalServerError, "failed to fetch squad inbounds", err, cfg)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"squad_inbounds": squadInbounds,
		"count":          len(squadInbounds),
	})
}

// handleSetSquadInbounds handles POST /api/v1/squad-inbounds
// Replaces all inbounds for a squad with the provided list.
func handleSetSquadInbounds(w http.ResponseWriter, r *http.Request, manager *manager.DatabaseManager, cfg *config.BackendConfig) {
	var req SquadInboundsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, http.StatusBadRequest, "invalid JSON", err, cfg)
		return
	}

	if err := req.Validate(); err != nil {
		sendError(w, http.StatusBadRequest, err.Error(), nil, cfg)
		return
	}

	ctx := r.Context()
	err := manager.ExecuteHighPriority(func(db *sql.DB) error {
		// Verify squad exists
		var squadID string
		err := db.QueryRowContext(ctx, "SELECT uuid FROM internal_squads WHERE uuid = ?", req.SquadUUID).Scan(&squadID)
		if err != nil {
			if err == sql.ErrNoRows {
				return fmt.Errorf("squad not found")
			}
			return err
		}

		// Delete existing inbounds for this squad
		_, err = db.ExecContext(ctx, "DELETE FROM internal_squad_inbounds WHERE internal_squad_uuid = ?", req.SquadUUID)
		if err != nil {
			return fmt.Errorf("failed to clear existing inbounds: %w", err)
		}

		// Insert new inbounds
		for _, inboundUUID := range req.InboundUUIDs {
			// Verify inbound exists
			var inboundID string
			err := db.QueryRowContext(ctx, "SELECT uuid FROM config_profile_inbounds WHERE uuid = ?", inboundUUID).Scan(&inboundID)
			if err != nil {
				if err == sql.ErrNoRows {
					return fmt.Errorf("inbound not found: %s", inboundUUID)
				}
				return err
			}

			_, err = db.ExecContext(ctx,
				"INSERT INTO internal_squad_inbounds (internal_squad_uuid, inbound_uuid) VALUES (?, ?)",
				inboundUUID, req.SquadUUID)
			if err != nil {
				return fmt.Errorf("failed to insert inbound: %w", err)
			}
		}

		return nil
	})

	if err != nil {
		if err.Error() == "squad not found" {
			sendError(w, http.StatusNotFound, "squad not found", err, cfg)
			return
		}
		if strings.Contains(err.Error(), "inbound not found") {
			sendError(w, http.StatusBadRequest, err.Error(), err, cfg)
			return
		}
		sendError(w, http.StatusInternalServerError, "failed to set squad inbounds", err, cfg)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":        "squad inbounds updated",
		"squad_uuid":     req.SquadUUID,
		"inbounds_count": len(req.InboundUUIDs),
	})
}

// ==================== INTERNAL SQUAD MEMBERS ====================

// InternalSquadMember represents a binding between a squad and a user.
type InternalSquadMember struct {
	InternalSquadUUID string `json:"internal_squad_uuid"`
	UserID            int64  `json:"user_id"`
	Username          string `json:"username,omitempty"`
}

// SquadMembersRequest represents a request to set members for a squad.
type SquadMembersRequest struct {
	SquadUUID string  `json:"squad_uuid"`
	UserIDs   []int64 `json:"user_ids"`
}

// Validate validates the SquadMembersRequest.
func (r *SquadMembersRequest) Validate() error {
	if r.SquadUUID == "" {
		return fmt.Errorf("squad_uuid is required")
	}
	if _, err := uuid.Parse(r.SquadUUID); err != nil {
		return fmt.Errorf("invalid squad_uuid format")
	}
	for _, userID := range r.UserIDs {
		if userID <= 0 {
			return fmt.Errorf("invalid user_id: %d", userID)
		}
	}
	return nil
}

// SquadMembersHandler handles GET/POST /api/v1/squad-members
func SquadMembersHandler(manager *manager.DatabaseManager, cfg *config.BackendConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handleGetSquadMembers(w, r, manager, cfg)
		case http.MethodPost:
			handleSetSquadMembers(w, r, manager, cfg)
		default:
			http.Error(w, `{"error": "method not allowed"}`, http.StatusMethodNotAllowed)
		}
	}
}

// handleGetSquadMembers handles GET /api/v1/squad-members?squad_uuid={uuid}
func handleGetSquadMembers(w http.ResponseWriter, r *http.Request, manager *manager.DatabaseManager, cfg *config.BackendConfig) {
	ctx := r.Context()
	squadUUID := r.URL.Query().Get("squad_uuid")

	var squadMembers []InternalSquadMember

	err := manager.ExecuteHighPriority(func(db *sql.DB) error {
		var query string
		var args []interface{}

		if squadUUID != "" {
			query = `
				SELECT m.internal_squad_uuid, m.user_id, u.username
				FROM internal_squad_members m
				JOIN users u ON m.user_id = u.t_id
				WHERE m.internal_squad_uuid = ?
				ORDER BY m.user_id`
			args = []interface{}{squadUUID}
		} else {
			query = `
				SELECT m.internal_squad_uuid, m.user_id, u.username
				FROM internal_squad_members m
				JOIN users u ON m.user_id = u.t_id
				ORDER BY m.internal_squad_uuid, m.user_id`
		}

		rows, err := db.QueryContext(ctx, query, args...)
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			var sm InternalSquadMember
			if err := rows.Scan(&sm.InternalSquadUUID, &sm.UserID, &sm.Username); err != nil {
				return err
			}
			squadMembers = append(squadMembers, sm)
		}
		return rows.Err()
	})

	if err != nil {
		sendError(w, http.StatusInternalServerError, "failed to fetch squad members", err, cfg)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"squad_members": squadMembers,
		"count":         len(squadMembers),
	})
}

// handleSetSquadMembers handles POST /api/v1/squad-members
// Replaces all members for a squad with the provided list.
func handleSetSquadMembers(w http.ResponseWriter, r *http.Request, manager *manager.DatabaseManager, cfg *config.BackendConfig) {
	var req SquadMembersRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, http.StatusBadRequest, "invalid JSON", err, cfg)
		return
	}

	if err := req.Validate(); err != nil {
		sendError(w, http.StatusBadRequest, err.Error(), nil, cfg)
		return
	}

	ctx := r.Context()
	err := manager.ExecuteHighPriority(func(db *sql.DB) error {
		// Verify squad exists
		var squadID string
		err := db.QueryRowContext(ctx, "SELECT uuid FROM internal_squads WHERE uuid = ?", req.SquadUUID).Scan(&squadID)
		if err != nil {
			if err == sql.ErrNoRows {
				return fmt.Errorf("squad not found")
			}
			return err
		}

		// Delete existing members for this squad
		_, err = db.ExecContext(ctx, "DELETE FROM internal_squad_members WHERE internal_squad_uuid = ?", req.SquadUUID)
		if err != nil {
			return fmt.Errorf("failed to clear existing members: %w", err)
		}

		// Insert new members
		for _, userID := range req.UserIDs {
			// Verify user exists
			var exists bool
			err := db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM users WHERE t_id = ?)", userID).Scan(&exists)
			if err != nil {
				return err
			}
			if !exists {
				return fmt.Errorf("user not found: %d", userID)
			}

			_, err = db.ExecContext(ctx,
				"INSERT INTO internal_squad_members (internal_squad_uuid, user_id) VALUES (?, ?)",
				req.SquadUUID, userID)
			if err != nil {
				return fmt.Errorf("failed to insert member: %w", err)
			}
		}

		return nil
	})

	if err != nil {
		if err.Error() == "squad not found" {
			sendError(w, http.StatusNotFound, "squad not found", err, cfg)
			return
		}
		if strings.Contains(err.Error(), "user not found") {
			sendError(w, http.StatusBadRequest, err.Error(), err, cfg)
			return
		}
		sendError(w, http.StatusInternalServerError, "failed to set squad members", err, cfg)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":     "squad members updated",
		"squad_uuid":  req.SquadUUID,
		"members_count": len(req.UserIDs),
	})
}
