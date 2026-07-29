package subscriptionpageconfigs

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"exodus/internal/config"
	"exodus/internal/httpapi/shared"
	monitor "exodus/internal/subscriptionnodes"

	"github.com/google/uuid"
)

func SubscriptionPageConfigsHandler(db *sql.DB, cfg *config.BackendConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handleGetSubscriptionPageConfigs(w, r, db, cfg)
		case http.MethodPost:
			handleCreateSubscriptionPageConfig(w, r, db, cfg)
		case http.MethodPatch:
			handleUpdateSubscriptionPageConfig(w, r, db, cfg)
		default:
			shared.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
		}
	}
}

func SubscriptionPageConfigsActionsHandler(db *sql.DB, cfg *config.BackendConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			shared.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}

		path := strings.TrimPrefix(r.URL.Path, subpageConfigsActionsPath)
		path = strings.Trim(path, "/")
		switch path {
		case "reorder":
			handleReorderSubscriptionPageConfigs(w, r, db, cfg)
		case "clone":
			handleCloneSubscriptionPageConfig(w, r, db, cfg)
		default:
			http.NotFound(w, r)
		}
	}
}

func SubscriptionPageConfigByUUIDHandler(db *sql.DB, cfg *config.BackendConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uuidStr := strings.TrimSpace(strings.TrimPrefix(r.URL.Path, subpageConfigsBasePath+"/"))
		if uuidStr == "" {
			switch r.Method {
			case http.MethodGet:
				handleGetSubscriptionPageConfigs(w, r, db, cfg)
			case http.MethodPost:
				handleCreateSubscriptionPageConfig(w, r, db, cfg)
			case http.MethodPatch:
				handleUpdateSubscriptionPageConfig(w, r, db, cfg)
			default:
				shared.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
			}
			return
		}
		if _, err := uuid.Parse(uuidStr); err != nil {
			shared.SendError(w, http.StatusBadRequest, "invalid UUID format", err, cfg)
			return
		}

		switch r.Method {
		case http.MethodGet:
			handleGetSubscriptionPageConfigByUUID(w, r, db, cfg, uuidStr)
		case http.MethodDelete:
			handleDeleteSubscriptionPageConfig(w, r, db, cfg, uuidStr)
		default:
			shared.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
		}
	}
}

func handleGetSubscriptionPageConfigs(w http.ResponseWriter, r *http.Request, db *sql.DB, cfg *config.BackendConfig) {
	ctx := r.Context()
	rows, err := db.QueryContext(ctx, `
		SELECT uuid, view_position, name, created_at, updated_at
		FROM subscription_page_config
		ORDER BY view_position ASC
	`)
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to fetch subscription page configs", err, cfg)
		return
	}
	defer rows.Close()

	configs := []SubscriptionPageConfig{}
	for rows.Next() {
		var cfgItem SubscriptionPageConfig
		var viewPosition sql.NullInt64
		if err := rows.Scan(&cfgItem.UUID, &viewPosition, &cfgItem.Name, &cfgItem.CreatedAt, &cfgItem.UpdatedAt); err != nil {
			shared.SendError(w, http.StatusInternalServerError, "failed to scan subscription page config", err, cfg)
			return
		}
		if viewPosition.Valid {
			cfgItem.ViewPosition = int(viewPosition.Int64)
		}
		configs = append(configs, cfgItem)
	}
	if err := rows.Err(); err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to iterate subscription page configs", err, cfg)
		return
	}

	shared.WriteJSON(w, http.StatusOK, map[string]any{
		"response": subscriptionPageConfigsListResponse{
			Total:   len(configs),
			Configs: configs,
		},
	})
}

func handleGetSubscriptionPageConfigByUUID(w http.ResponseWriter, r *http.Request, db *sql.DB, cfg *config.BackendConfig, uuidStr string) {
	ctx := r.Context()
	cfgItem, err := fetchSubscriptionPageConfig(ctx, db, uuidStr, true)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			shared.SendError(w, http.StatusNotFound, "subscription page config not found", nil, cfg)
			return
		}
		shared.SendError(w, http.StatusInternalServerError, "failed to fetch subscription page config", err, cfg)
		return
	}

	shared.WriteJSON(w, http.StatusOK, map[string]any{
		"response": cfgItem,
	})
}

func handleCreateSubscriptionPageConfig(w http.ResponseWriter, r *http.Request, db *sql.DB, cfg *config.BackendConfig) {
	var req subpageConfigCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		shared.SendError(w, http.StatusBadRequest, "invalid JSON", err, cfg)
		return
	}
	name, err := normalizeSubpageConfigName(req.Name)
	if err != nil {
		shared.SendError(w, http.StatusBadRequest, err.Error(), nil, cfg)
		return
	}

	ctx := r.Context()
	defaultConfig, err := fetchDefaultSubpageConfig(ctx, db)
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to fetch default subpage config", err, cfg)
		return
	}

	row := db.QueryRowContext(ctx, `
		INSERT INTO subscription_page_config (name, config)
		VALUES ($1, $2)
		RETURNING uuid, view_position, name, config, created_at, updated_at
	`, name, string(defaultConfig))

	var created SubscriptionPageConfig
	var viewPosition sql.NullInt64
	var configStr sql.NullString
	if err := row.Scan(&created.UUID, &viewPosition, &created.Name, &configStr, &created.CreatedAt, &created.UpdatedAt); err != nil {
		if isUniqueNameError(err) {
			shared.SendError(w, http.StatusConflict, "config name already exists", err, cfg)
			return
		}
		shared.SendError(w, http.StatusInternalServerError, "failed to create subscription page config", err, cfg)
		return
	}
	if viewPosition.Valid {
		created.ViewPosition = int(viewPosition.Int64)
	}
	if configStr.Valid {
		created.Config = json.RawMessage(configStr.String)
	}

	shared.WriteJSON(w, http.StatusCreated, map[string]any{
		"response": created,
	})
	emitSubpageConfigChanged(ctx, cfg, "created", created, nil)
	targetNodeUUIDs, targetErr := getSubNodeUUIDsBySubpageConfigUUID(ctx, db, created.UUID)
	if targetErr != nil {
		cfg.Logger.Warn("Failed to resolve sub nodes for created subpage config push", "subpage_config_uuid", created.UUID, "error", targetErr)
		return
	}
	if len(targetNodeUUIDs) == 0 {
		return
	}
	monitor.RequestSubNodeSubpageConfigPush(created.UUID, created.Config, targetNodeUUIDs...)
}

func handleUpdateSubscriptionPageConfig(w http.ResponseWriter, r *http.Request, db *sql.DB, cfg *config.BackendConfig) {
	var req subpageConfigUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		shared.SendError(w, http.StatusBadRequest, "invalid JSON", err, cfg)
		return
	}
	if _, err := uuid.Parse(req.UUID); err != nil {
		shared.SendError(w, http.StatusBadRequest, "invalid UUID format", err, cfg)
		return
	}

	if req.Name == nil && req.Config == nil {
		shared.SendError(w, http.StatusBadRequest, "no fields to update", nil, cfg)
		return
	}

	updates := []string{}
	args := []any{}
	idx := 1

	if req.Name != nil {
		name, err := normalizeSubpageConfigName(*req.Name)
		if err != nil {
			shared.SendError(w, http.StatusBadRequest, err.Error(), nil, cfg)
			return
		}
		updates = append(updates, fmt.Sprintf("name = $%d", idx))
		args = append(args, name)
		idx++
	}

	if req.Config != nil {
		cleaned, err := normalizeAndValidateSubpageConfig(*req.Config)
		if err != nil {
			shared.SendError(w, http.StatusBadRequest, err.Error(), err, cfg)
			return
		}
		configJSON, err := json.Marshal(cleaned)
		if err != nil {
			shared.SendError(w, http.StatusBadRequest, "invalid config", err, cfg)
			return
		}
		updates = append(updates, fmt.Sprintf("config = $%d", idx))
		args = append(args, string(configJSON))
		idx++
	}

	updates = append(updates, "updated_at = CURRENT_TIMESTAMP")
	args = append(args, req.UUID)

	ctx := r.Context()
	query := fmt.Sprintf(`
		UPDATE subscription_page_config
		SET %s
		WHERE uuid = $%d
		RETURNING uuid, view_position, name, config, created_at, updated_at
	`, strings.Join(updates, ", "), idx)

	row := db.QueryRowContext(ctx, query, args...)
	var updated SubscriptionPageConfig
	var viewPosition sql.NullInt64
	var configStr sql.NullString
	if err := row.Scan(&updated.UUID, &viewPosition, &updated.Name, &configStr, &updated.CreatedAt, &updated.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			shared.SendError(w, http.StatusNotFound, "subscription page config not found", nil, cfg)
			return
		}
		if isUniqueNameError(err) {
			shared.SendError(w, http.StatusConflict, "config name already exists", err, cfg)
			return
		}
		shared.SendError(w, http.StatusInternalServerError, "failed to update subscription page config", err, cfg)
		return
	}
	if viewPosition.Valid {
		updated.ViewPosition = int(viewPosition.Int64)
	}
	if configStr.Valid {
		updated.Config = json.RawMessage(configStr.String)
	}

	shared.WriteJSON(w, http.StatusOK, map[string]any{
		"response": updated,
	})
	emitSubpageConfigChanged(ctx, cfg, "updated", updated, nil)
	targetNodeUUIDs, targetErr := getSubNodeUUIDsBySubpageConfigUUID(ctx, db, updated.UUID)
	if targetErr != nil {
		cfg.Logger.Warn("Failed to resolve sub nodes for updated subpage config push", "subpage_config_uuid", updated.UUID, "error", targetErr)
		return
	}
	if len(targetNodeUUIDs) == 0 {
		cfg.Logger.Debug("No linked subscription nodes for updated subpage config push", "subpage_config_uuid", updated.UUID)
		return
	}
	cfg.Logger.Info(
		"Dispatching updated subpage config push",
		"subpage_config_uuid", updated.UUID,
		"target_nodes_count", len(targetNodeUUIDs),
		"target_node_uuids", targetNodeUUIDs,
	)
	monitor.RequestSubNodeSubpageConfigPush(updated.UUID, updated.Config, targetNodeUUIDs...)
}

func handleDeleteSubscriptionPageConfig(w http.ResponseWriter, r *http.Request, db *sql.DB, cfg *config.BackendConfig, uuidStr string) {
	if uuidStr == defaultSubpageConfigUUID {
		shared.SendError(w, http.StatusBadRequest, "reserved config cannot be deleted", nil, cfg)
		return
	}

	ctx := r.Context()
	targetNodeUUIDs, targetErr := getSubNodeUUIDsBySubpageConfigUUID(ctx, db, uuidStr)
	if targetErr != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to resolve linked subscription nodes", targetErr, cfg)
		return
	}

	result, err := db.ExecContext(ctx, `DELETE FROM subscription_page_config WHERE uuid = $1`, uuidStr)
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to delete subscription page config", err, cfg)
		return
	}
	ra, err := result.RowsAffected()
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to read rows affected", err, cfg)
		return
	}
	if ra == 0 {
		shared.SendError(w, http.StatusNotFound, "subscription page config not found", nil, cfg)
		return
	}

	shared.WriteJSON(w, http.StatusOK, map[string]any{
		"response": subpageConfigDeleteResponse{IsDeleted: true},
	})
	emitSubpageConfigChanged(ctx, cfg, "deleted", SubscriptionPageConfig{UUID: uuidStr}, nil)
	if len(targetNodeUUIDs) == 0 {
		return
	}
	monitor.RequestSubNodeSubpageConfigPush(uuidStr, nil, targetNodeUUIDs...)
}

func handleReorderSubscriptionPageConfigs(w http.ResponseWriter, r *http.Request, db *sql.DB, cfg *config.BackendConfig) {
	var req subpageConfigReorderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		shared.SendError(w, http.StatusBadRequest, "invalid JSON", err, cfg)
		return
	}
	if len(req.Items) == 0 {
		shared.SendError(w, http.StatusBadRequest, "items is required", nil, cfg)
		return
	}

	seen := map[string]struct{}{}
	for _, item := range req.Items {
		if _, err := uuid.Parse(item.UUID); err != nil {
			shared.SendError(w, http.StatusBadRequest, "invalid UUID format", err, cfg)
			return
		}
		if _, ok := seen[item.UUID]; ok {
			shared.SendError(w, http.StatusBadRequest, "duplicate uuid in items", nil, cfg)
			return
		}
		seen[item.UUID] = struct{}{}
	}

	ctx := r.Context()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "begin tx failed", err, cfg)
		return
	}
	defer func() {
		_ = tx.Rollback()
	}()

	for _, item := range req.Items {
		if _, err := tx.ExecContext(ctx, `UPDATE subscription_page_config SET view_position = $1 WHERE uuid = $2`, item.ViewPosition, item.UUID); err != nil {
			shared.SendError(w, http.StatusInternalServerError, "failed to reorder subscription page configs", err, cfg)
			return
		}
	}
	if _, err := tx.ExecContext(ctx, `SELECT setval('subscription_page_config_view_position_seq', (SELECT COALESCE(MAX(view_position), 0) FROM subscription_page_config) + 1)`); err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to update sequence", err, cfg)
		return
	}
	if err := tx.Commit(); err != nil {
		shared.SendError(w, http.StatusInternalServerError, "commit reorder failed", err, cfg)
		return
	}

	emitSubpageConfigChanged(ctx, cfg, "reordered", SubscriptionPageConfig{}, map[string]any{"items": len(req.Items)})
	handleGetSubscriptionPageConfigs(w, r, db, cfg)
}

func handleCloneSubscriptionPageConfig(w http.ResponseWriter, r *http.Request, db *sql.DB, cfg *config.BackendConfig) {
	var req subpageConfigCloneRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		shared.SendError(w, http.StatusBadRequest, "invalid JSON", err, cfg)
		return
	}
	if _, err := uuid.Parse(req.CloneFromUUID); err != nil {
		shared.SendError(w, http.StatusBadRequest, "invalid UUID format", err, cfg)
		return
	}

	ctx := r.Context()
	var cfgItem SubscriptionPageConfig
	row := db.QueryRowContext(ctx, `
		SELECT uuid, view_position, name, config, created_at, updated_at
		FROM subscription_page_config
		WHERE uuid = $1
		LIMIT 1
	`, req.CloneFromUUID)

	var viewPosition sql.NullInt64
	var configStr sql.NullString
	if err := row.Scan(&cfgItem.UUID, &viewPosition, &cfgItem.Name, &configStr, &cfgItem.CreatedAt, &cfgItem.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			shared.SendError(w, http.StatusNotFound, "subscription page config not found", nil, cfg)
			return
		}
		shared.SendError(w, http.StatusInternalServerError, "failed to fetch subscription page config for cloning", err, cfg)
		return
	}
	if viewPosition.Valid {
		cfgItem.ViewPosition = int(viewPosition.Int64)
	}
	if configStr.Valid {
		cfgItem.Config = json.RawMessage(configStr.String)
	}

	cloneName := fmt.Sprintf("Clone %s", randomSuffix(5))
	insertRow := db.QueryRowContext(ctx, `
		INSERT INTO subscription_page_config (name, config)
		VALUES ($1, $2)
		RETURNING uuid, view_position, name, config, created_at, updated_at
	`, cloneName, string(cfgItem.Config))

	var created SubscriptionPageConfig
	var insertViewPosition sql.NullInt64
	var insertConfigStr sql.NullString
	if err := insertRow.Scan(&created.UUID, &insertViewPosition, &created.Name, &insertConfigStr, &created.CreatedAt, &created.UpdatedAt); err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to clone subscription page config", err, cfg)
		return
	}
	if insertViewPosition.Valid {
		created.ViewPosition = int(insertViewPosition.Int64)
	}
	if insertConfigStr.Valid {
		created.Config = json.RawMessage(insertConfigStr.String)
	}

	shared.WriteJSON(w, http.StatusOK, map[string]any{
		"response": created,
	})
	emitSubpageConfigChanged(ctx, cfg, "cloned", created, map[string]any{"cloneFromUuid": req.CloneFromUUID})
	targetNodeUUIDs, targetErr := getSubNodeUUIDsBySubpageConfigUUID(ctx, db, created.UUID)
	if targetErr != nil {
		cfg.Logger.Warn("Failed to resolve sub nodes for cloned subpage config push", "subpage_config_uuid", created.UUID, "error", targetErr)
		return
	}
	if len(targetNodeUUIDs) == 0 {
		return
	}
	monitor.RequestSubNodeSubpageConfigPush(created.UUID, created.Config, targetNodeUUIDs...)
}
