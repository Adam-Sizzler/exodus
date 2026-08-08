package squads

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
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
			if len(parts) == 3 && parts[1] == "bulk-actions" {
				switch parts[2] {
				case "add-users":
					if r.Method != http.MethodPost {
						shared.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
						return
					}
					handleBulkAddUsersToInternalSquad(w, r, db, cfg, squadUUID)
					return
				case "remove-users":
					if r.Method != http.MethodDelete && r.Method != http.MethodPost {
						shared.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
						return
					}
					handleBulkRemoveUsersFromInternalSquad(w, r, db, cfg, squadUUID)
					return
				case "add-many-users":
					if r.Method != http.MethodPost {
						shared.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
						return
					}
					handleBulkAddManyUsersToInternalSquad(w, r, db, cfg, squadUUID)
					return
				case "remove-many-users":
					if r.Method != http.MethodDelete && r.Method != http.MethodPost {
						shared.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
						return
					}
					handleBulkRemoveManyUsersFromInternalSquad(w, r, db, cfg, squadUUID)
					return
				}
			}
			if len(parts) == 2 && parts[1] == "accessible-nodes" && r.Method == http.MethodGet {
				handleGetInternalSquadAccessibleNodes(w, r, service, squadUUID)
				return
			}
			if len(parts) == 2 && parts[1] == "usage" && r.Method == http.MethodGet {
				handleGetInternalSquadUsage(w, r, db, cfg, squadUUID)
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
		SELECT $1, id FROM users
		ON CONFLICT (internal_squad_uuid, user_id) DO NOTHING
	`, squadUUID)
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to add users to internal squad", err, cfg)
		return
	}

	w.WriteHeader(http.StatusAccepted)
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

	w.WriteHeader(http.StatusAccepted)
}

func handleBulkAddManyUsersToInternalSquad(w http.ResponseWriter, r *http.Request, db *sql.DB, cfg *config.BackendConfig, squadUUID string) {
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

	var req struct {
		UserIDs []int64 `json:"userIds"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || len(req.UserIDs) == 0 || len(req.UserIDs) > 1000 {
		shared.SendError(w, http.StatusBadRequest, "invalid request body (userIds array of 1..1000 required)", nil, cfg)
		return
	}

	_, err := db.ExecContext(ctx, `
		INSERT INTO internal_squad_members (internal_squad_uuid, user_id)
		SELECT $1, unnest($2::bigint[])
		ON CONFLICT (internal_squad_uuid, user_id) DO NOTHING
	`, squadUUID, req.UserIDs)
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to add users to internal squad", err, cfg)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

func handleBulkRemoveManyUsersFromInternalSquad(w http.ResponseWriter, r *http.Request, db *sql.DB, cfg *config.BackendConfig, squadUUID string) {
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

	var req struct {
		UserIDs []int64 `json:"userIds"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || len(req.UserIDs) == 0 || len(req.UserIDs) > 1000 {
		shared.SendError(w, http.StatusBadRequest, "invalid request body (userIds array of 1..1000 required)", nil, cfg)
		return
	}

	_, err := db.ExecContext(ctx, `
		DELETE FROM internal_squad_members
		WHERE internal_squad_uuid = $1 AND user_id = ANY($2::bigint[])
	`, squadUUID, req.UserIDs)
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to remove users from internal squad", err, cfg)
		return
	}

	w.WriteHeader(http.StatusAccepted)
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

func handleGetInternalSquadUsage(w http.ResponseWriter, r *http.Request, db *sql.DB, cfg *config.BackendConfig, squadUUID string) {
	q := r.URL.Query()
	limit := 250
	if limitStr := q.Get("limit"); limitStr != "" {
		if parsed, err := strconv.Atoi(limitStr); err == nil && parsed >= 1 {
			limit = parsed
		}
	}
	var cursor int64 = 0
	if cursorStr := q.Get("cursor"); cursorStr != "" {
		if parsed, err := strconv.ParseInt(cursorStr, 10, 64); err == nil && parsed >= 0 {
			cursor = parsed
		}
	}
	var minTotalBytes int64 = 0
	if minStr := q.Get("minTotalBytes"); minStr != "" {
		if parsed, err := strconv.ParseInt(minStr, 10, 64); err == nil && parsed > 0 {
			minTotalBytes = parsed
		}
	}

	var total int64
	_ = db.QueryRowContext(r.Context(), `SELECT COUNT(*) FROM internal_squad_members WHERE internal_squad_uuid = $1`, squadUUID).Scan(&total)

	rows, err := db.QueryContext(r.Context(), `
		SELECT u.id, u.username, COALESCE(SUM(nuh.total_bytes), 0) AS total_bytes
		FROM internal_squad_members ism
		JOIN users u ON ism.user_id = u.id
		LEFT JOIN internal_squad_inbounds isi ON isi.internal_squad_uuid = ism.internal_squad_uuid
		LEFT JOIN config_profile_inbounds_to_nodes cpin ON cpin.config_profile_inbound_uuid = isi.inbound_uuid
		LEFT JOIN nodes n ON n.uuid = cpin.node_uuid
		LEFT JOIN nodes_user_usage_history nuh ON nuh.user_id = u.id AND nuh.node_id = n.id
		WHERE ism.internal_squad_uuid = $1 AND u.id > $2
		GROUP BY u.id, u.username
		HAVING COALESCE(SUM(nuh.total_bytes), 0) >= $3
		ORDER BY u.id ASC
		LIMIT $4
	`, squadUUID, cursor, minTotalBytes, limit)
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to fetch squad usage", err, cfg)
		return
	}
	defer rows.Close()

	type squadUserUsageItem struct {
		ID         int64  `json:"id"`
		Username   string `json:"username"`
		TotalBytes int64  `json:"totalBytes"`
	}

	users := make([]squadUserUsageItem, 0)
	var lastID int64 = 0
	for rows.Next() {
		var item squadUserUsageItem
		if scanErr := rows.Scan(&item.ID, &item.Username, &item.TotalBytes); scanErr != nil {
			shared.SendError(w, http.StatusInternalServerError, "failed to scan squad user usage item", scanErr, cfg)
			return
		}
		users = append(users, item)
		lastID = item.ID
	}
	if err := rows.Err(); err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to fetch squad user usage items", err, cfg)
		return
	}

	var nextCursor *string
	hasMore := len(users) == limit
	if hasMore {
		cStr := fmt.Sprintf("%d", lastID)
		nextCursor = &cStr
	}

	shared.WriteJSON(w, http.StatusOK, map[string]any{
		"response": map[string]any{
			"squadUuid":  squadUUID,
			"users":      users,
			"nextCursor": nextCursor,
			"hasMore":    hasMore,
			"total":      total,
		},
	})
}

func BandwidthStatsInternalSquadsHandler(db *sql.DB, cfg *config.BackendConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/api/bandwidth-stats/internal-squads")
		path = strings.Trim(path, "/")
		parts := strings.Split(path, "/")
		if len(parts) == 0 || parts[0] == "" {
			shared.WriteJSONError(w, http.StatusNotFound, "not found")
			return
		}

		squadUUID := parts[0]
		if len(parts) == 2 && parts[1] == "usage" {
			handleGetInternalSquadUsage(w, r, db, cfg, squadUUID)
			return
		}
		if len(parts) == 4 && parts[1] == "users" && parts[3] == "usage" {
			userID, _ := strconv.ParseInt(parts[2], 10, 64)
			handleGetInternalSquadUserUsage(w, r, db, cfg, squadUUID, userID)
			return
		}

		shared.WriteJSONError(w, http.StatusNotFound, "not found")
	}
}

func handleGetInternalSquadUserUsage(w http.ResponseWriter, r *http.Request, db *sql.DB, cfg *config.BackendConfig, squadUUID string, userID int64) {
	shared.WriteJSON(w, http.StatusOK, map[string]any{
		"response": map[string]any{
			"categories": []string{},
			"data":       []int64{},
		},
	})
}
