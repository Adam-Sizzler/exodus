package subscriptionpageconfigs

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"exodus/internal/config"
	"exodus/internal/db"
	dbmanager "exodus/internal/db/manager"
	"exodus/internal/httpapi/shared"
	"exodus/internal/notifications"
	monitor "exodus/internal/subscriptionnodes"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
)

const (
	subpageConfigsBasePath    = "/api/subscription-page-configs"
	subpageConfigsActionsPath = "/api/subscription-page-configs/actions/"
	defaultSubpageConfigUUID  = "00000000-0000-0000-0000-000000000000"
)

var subpageConfigNameRegex = regexp.MustCompile(`^[A-Za-z0-9_\s-]+$`)

// SubscriptionPageConfig represents a subscription page config in API responses.
type SubscriptionPageConfig struct {
	UUID         string          `json:"uuid"`
	ViewPosition int             `json:"viewPosition"`
	Name         string          `json:"name"`
	Config       json.RawMessage `json:"config"`
	CreatedAt    time.Time       `json:"createdAt"`
	UpdatedAt    time.Time       `json:"updatedAt"`
}

type subscriptionPageConfigsListResponse struct {
	Total   int                      `json:"total"`
	Configs []SubscriptionPageConfig `json:"configs"`
}

type subpageConfigCreateRequest struct {
	Name string `json:"name"`
}

type subpageConfigUpdateRequest struct {
	UUID   string           `json:"uuid"`
	Name   *string          `json:"name,omitempty"`
	Config *json.RawMessage `json:"config,omitempty"`
}

type subpageConfigDeleteResponse struct {
	IsDeleted bool `json:"isDeleted"`
}

type subpageConfigReorderItem struct {
	UUID         string `json:"uuid"`
	ViewPosition int    `json:"viewPosition"`
}

type subpageConfigReorderRequest struct {
	Items []subpageConfigReorderItem `json:"items"`
}

type subpageConfigCloneRequest struct {
	CloneFromUUID string `json:"cloneFromUuid"`
}

func SubscriptionPageConfigsHandler(manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handleGetSubscriptionPageConfigs(w, r, manager, cfg)
		case http.MethodPost:
			handleCreateSubscriptionPageConfig(w, r, manager, cfg)
		case http.MethodPatch:
			handleUpdateSubscriptionPageConfig(w, r, manager, cfg)
		default:
			shared.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
		}
	}
}

func SubscriptionPageConfigsActionsHandler(manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			shared.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}

		path := strings.TrimPrefix(r.URL.Path, subpageConfigsActionsPath)
		path = strings.Trim(path, "/")
		switch path {
		case "reorder":
			handleReorderSubscriptionPageConfigs(w, r, manager, cfg)
		case "clone":
			handleCloneSubscriptionPageConfig(w, r, manager, cfg)
		default:
			http.NotFound(w, r)
		}
	}
}

func SubscriptionPageConfigByUUIDHandler(manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uuidStr := strings.TrimSpace(strings.TrimPrefix(r.URL.Path, subpageConfigsBasePath+"/"))
		if uuidStr == "" {
			switch r.Method {
			case http.MethodGet:
				handleGetSubscriptionPageConfigs(w, r, manager, cfg)
			case http.MethodPost:
				handleCreateSubscriptionPageConfig(w, r, manager, cfg)
			case http.MethodPatch:
				handleUpdateSubscriptionPageConfig(w, r, manager, cfg)
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
			handleGetSubscriptionPageConfigByUUID(w, r, manager, cfg, uuidStr)
		case http.MethodDelete:
			handleDeleteSubscriptionPageConfig(w, r, manager, cfg, uuidStr)
		default:
			shared.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
		}
	}
}

func handleGetSubscriptionPageConfigs(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) {
	ctx := r.Context()
	configs := []SubscriptionPageConfig{}

	err := manager.ExecuteHighPriority(func(dbExec dbmanager.DBExecutor) error {
		rows, err := dbExec.QueryContext(ctx, `
			SELECT uuid, view_position, name, created_at, updated_at
			FROM subscription_page_config
			ORDER BY view_position ASC
		`)
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			var cfgItem SubscriptionPageConfig
			var viewPosition sql.NullInt64
			if err := rows.Scan(&cfgItem.UUID, &viewPosition, &cfgItem.Name, &cfgItem.CreatedAt, &cfgItem.UpdatedAt); err != nil {
				return err
			}
			if viewPosition.Valid {
				cfgItem.ViewPosition = int(viewPosition.Int64)
			}
			configs = append(configs, cfgItem)
		}
		return rows.Err()
	})
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to fetch subscription page configs", err, cfg)
		return
	}

	shared.WriteJSON(w, http.StatusOK, map[string]any{
		"response": subscriptionPageConfigsListResponse{
			Total:   len(configs),
			Configs: configs,
		},
	})
}

func handleGetSubscriptionPageConfigByUUID(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig, uuidStr string) {
	ctx := r.Context()
	cfgItem, err := fetchSubscriptionPageConfig(ctx, manager, uuidStr, true)
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

func handleCreateSubscriptionPageConfig(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) {
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
	var created SubscriptionPageConfig

	err = manager.ExecuteHighPriority(func(dbExec dbmanager.DBExecutor) error {
		defaultConfig, err := fetchDefaultSubpageConfig(ctx, dbExec)
		if err != nil {
			return err
		}

		row := dbExec.QueryRowContext(ctx, `
			INSERT INTO subscription_page_config (name, config)
			VALUES (?, ?)
			RETURNING uuid, view_position, name, config, created_at, updated_at
		`, name, string(defaultConfig))

		var viewPosition sql.NullInt64
		var configStr sql.NullString
		if err := row.Scan(&created.UUID, &viewPosition, &created.Name, &configStr, &created.CreatedAt, &created.UpdatedAt); err != nil {
			return err
		}
		if viewPosition.Valid {
			created.ViewPosition = int(viewPosition.Int64)
		}
		if configStr.Valid {
			created.Config = json.RawMessage(configStr.String)
		}
		return nil
	})
	if err != nil {
		if isUniqueNameError(err) {
			shared.SendError(w, http.StatusConflict, "config name already exists", err, cfg)
			return
		}
		shared.SendError(w, http.StatusInternalServerError, "failed to create subscription page config", err, cfg)
		return
	}

	shared.WriteJSON(w, http.StatusCreated, map[string]any{
		"response": created,
	})
	emitSubpageConfigChanged(ctx, cfg, "created", created, nil)
	targetNodeUUIDs, targetErr := getSubNodeUUIDsBySubpageConfigUUID(ctx, manager, created.UUID)
	if targetErr != nil {
		cfg.Logger.Warn("Failed to resolve sub nodes for created subpage config push", "subpage_config_uuid", created.UUID, "error", targetErr)
		return
	}
	if len(targetNodeUUIDs) == 0 {
		return
	}
	monitor.RequestSubNodeSubpageConfigPush(created.UUID, created.Config, targetNodeUUIDs...)
}

func handleUpdateSubscriptionPageConfig(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) {
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

	if req.Name != nil {
		name, err := normalizeSubpageConfigName(*req.Name)
		if err != nil {
			shared.SendError(w, http.StatusBadRequest, err.Error(), nil, cfg)
			return
		}
		updates = append(updates, "name = ?")
		args = append(args, name)
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
		updates = append(updates, "config = ?")
		args = append(args, string(configJSON))
	}

	updates = append(updates, "updated_at = CURRENT_TIMESTAMP")
	args = append(args, req.UUID)

	ctx := r.Context()
	var updated SubscriptionPageConfig

	err := manager.ExecuteHighPriority(func(dbExec dbmanager.DBExecutor) error {
		query := fmt.Sprintf(`
			UPDATE subscription_page_config
			SET %s
			WHERE uuid = ?
			RETURNING uuid, view_position, name, config, created_at, updated_at
		`, strings.Join(updates, ", "))

		row := dbExec.QueryRowContext(ctx, query, args...)
		var viewPosition sql.NullInt64
		var configStr sql.NullString
		if err := row.Scan(&updated.UUID, &viewPosition, &updated.Name, &configStr, &updated.CreatedAt, &updated.UpdatedAt); err != nil {
			return err
		}
		if viewPosition.Valid {
			updated.ViewPosition = int(viewPosition.Int64)
		}
		if configStr.Valid {
			updated.Config = json.RawMessage(configStr.String)
		}
		return nil
	})
	if err != nil {
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

	shared.WriteJSON(w, http.StatusOK, map[string]any{
		"response": updated,
	})
	emitSubpageConfigChanged(ctx, cfg, "updated", updated, nil)
	targetNodeUUIDs, targetErr := getSubNodeUUIDsBySubpageConfigUUID(ctx, manager, updated.UUID)
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

func handleDeleteSubscriptionPageConfig(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig, uuidStr string) {
	if uuidStr == defaultSubpageConfigUUID {
		shared.SendError(w, http.StatusBadRequest, "reserved config cannot be deleted", nil, cfg)
		return
	}

	ctx := r.Context()
	deleted := false

	targetNodeUUIDs, targetErr := getSubNodeUUIDsBySubpageConfigUUID(ctx, manager, uuidStr)
	if targetErr != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to resolve linked subscription nodes", targetErr, cfg)
		return
	}

	err := manager.ExecuteHighPriority(func(dbExec dbmanager.DBExecutor) error {
		result, err := dbExec.ExecContext(ctx, `DELETE FROM subscription_page_config WHERE uuid = ?`, uuidStr)
		if err != nil {
			return err
		}
		ra, err := result.RowsAffected()
		if err != nil {
			return err
		}
		deleted = ra > 0
		return nil
	})
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to delete subscription page config", err, cfg)
		return
	}
	if !deleted {
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

func handleReorderSubscriptionPageConfigs(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) {
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
	err := manager.ExecuteHighPriority(func(dbExec dbmanager.DBExecutor) error {
		tx, err := dbExec.BeginTx(ctx, nil)
		if err != nil {
			return err
		}
		for _, item := range req.Items {
			if _, err := tx.ExecContext(ctx, `UPDATE subscription_page_config SET view_position = ? WHERE uuid = ?`, item.ViewPosition, item.UUID); err != nil {
				_ = tx.Rollback()
				return err
			}
		}
		if _, err := tx.ExecContext(ctx, `SELECT setval('subscription_page_config_view_position_seq', (SELECT COALESCE(MAX(view_position), 0) FROM subscription_page_config) + 1)`); err != nil {
			_ = tx.Rollback()
			return err
		}
		return tx.Commit()
	})
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to reorder subscription page configs", err, cfg)
		return
	}

	emitSubpageConfigChanged(ctx, cfg, "reordered", SubscriptionPageConfig{}, map[string]any{"items": len(req.Items)})
	// Return refreshed list
	handleGetSubscriptionPageConfigs(w, r, manager, cfg)
}

func handleCloneSubscriptionPageConfig(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) {
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
	var created SubscriptionPageConfig

	err := manager.ExecuteHighPriority(func(dbExec dbmanager.DBExecutor) error {
		var cfgItem SubscriptionPageConfig
		row := dbExec.QueryRowContext(ctx, `
			SELECT uuid, view_position, name, config, created_at, updated_at
			FROM subscription_page_config
			WHERE uuid = ?
			LIMIT 1
		`, req.CloneFromUUID)

		var viewPosition sql.NullInt64
		var configStr sql.NullString
		if err := row.Scan(&cfgItem.UUID, &viewPosition, &cfgItem.Name, &configStr, &cfgItem.CreatedAt, &cfgItem.UpdatedAt); err != nil {
			return err
		}
		if viewPosition.Valid {
			cfgItem.ViewPosition = int(viewPosition.Int64)
		}
		if configStr.Valid {
			cfgItem.Config = json.RawMessage(configStr.String)
		}

		cloneName := fmt.Sprintf("Clone %s", randomSuffix(5))
		insertRow := dbExec.QueryRowContext(ctx, `
			INSERT INTO subscription_page_config (name, config)
			VALUES (?, ?)
			RETURNING uuid, view_position, name, config, created_at, updated_at
		`, cloneName, string(cfgItem.Config))

		var insertViewPosition sql.NullInt64
		var insertConfigStr sql.NullString
		if err := insertRow.Scan(&created.UUID, &insertViewPosition, &created.Name, &insertConfigStr, &created.CreatedAt, &created.UpdatedAt); err != nil {
			return err
		}
		if insertViewPosition.Valid {
			created.ViewPosition = int(insertViewPosition.Int64)
		}
		if insertConfigStr.Valid {
			created.Config = json.RawMessage(insertConfigStr.String)
		}
		return nil
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			shared.SendError(w, http.StatusNotFound, "subscription page config not found", nil, cfg)
			return
		}
		shared.SendError(w, http.StatusInternalServerError, "failed to clone subscription page config", err, cfg)
		return
	}

	shared.WriteJSON(w, http.StatusOK, map[string]any{
		"response": created,
	})
	emitSubpageConfigChanged(ctx, cfg, "cloned", created, map[string]any{"cloneFromUuid": req.CloneFromUUID})
	targetNodeUUIDs, targetErr := getSubNodeUUIDsBySubpageConfigUUID(ctx, manager, created.UUID)
	if targetErr != nil {
		cfg.Logger.Warn("Failed to resolve sub nodes for cloned subpage config push", "subpage_config_uuid", created.UUID, "error", targetErr)
		return
	}
	if len(targetNodeUUIDs) == 0 {
		return
	}
	monitor.RequestSubNodeSubpageConfigPush(created.UUID, created.Config, targetNodeUUIDs...)
}

func fetchSubscriptionPageConfig(ctx context.Context, manager *dbmanager.DatabaseManager, uuidStr string, withConfig bool) (SubscriptionPageConfig, error) {
	var cfgItem SubscriptionPageConfig

	err := manager.ExecuteHighPriority(func(dbExec dbmanager.DBExecutor) error {
		query := `
			SELECT uuid, view_position, name, created_at, updated_at
			FROM subscription_page_config
			WHERE uuid = ?
			LIMIT 1
		`
		if withConfig {
			query = `
				SELECT uuid, view_position, name, config, created_at, updated_at
				FROM subscription_page_config
				WHERE uuid = ?
				LIMIT 1
			`
		}

		row := dbExec.QueryRowContext(ctx, query, uuidStr)
		var viewPosition sql.NullInt64
		var configStr sql.NullString
		if withConfig {
			if err := row.Scan(&cfgItem.UUID, &viewPosition, &cfgItem.Name, &configStr, &cfgItem.CreatedAt, &cfgItem.UpdatedAt); err != nil {
				return err
			}
			if configStr.Valid {
				cfgItem.Config = json.RawMessage(configStr.String)
			}
		} else {
			if err := row.Scan(&cfgItem.UUID, &viewPosition, &cfgItem.Name, &cfgItem.CreatedAt, &cfgItem.UpdatedAt); err != nil {
				return err
			}
		}

		if viewPosition.Valid {
			cfgItem.ViewPosition = int(viewPosition.Int64)
		}
		return nil
	})

	return cfgItem, err
}

func emitSubpageConfigChanged(ctx context.Context, cfg *config.BackendConfig, action string, item SubscriptionPageConfig, extra map[string]any) {
	data := map[string]any{
		"action": action,
	}
	if item.UUID != "" {
		data["uuid"] = item.UUID
	}
	if item.Name != "" {
		data["name"] = item.Name
	}
	if item.ViewPosition != 0 {
		data["viewPosition"] = item.ViewPosition
	}
	for key, value := range extra {
		data[key] = value
	}
	notifications.Emit(ctx, cfg, notifications.Event{
		Scope: notifications.ScopeService,
		Event: notifications.EventServiceSubpageChanged,
		Data:  data,
	})
}

func fetchDefaultSubpageConfig(ctx context.Context, dbExec dbmanager.DBExecutor) (json.RawMessage, error) {
	row := dbExec.QueryRowContext(ctx, `SELECT config FROM subscription_page_config WHERE uuid = ?`, defaultSubpageConfigUUID)
	var cfgStr sql.NullString
	if err := row.Scan(&cfgStr); err == nil && cfgStr.Valid {
		return json.RawMessage(cfgStr.String), nil
	}

	// fallback to embedded default
	return json.RawMessage(db.DefaultSubscriptionPageConfig), nil
}

func getSubNodeUUIDsBySubpageConfigUUID(ctx context.Context, manager *dbmanager.DatabaseManager, subpageConfigUUID string) ([]string, error) {
	nodeUUIDs := make([]string, 0)
	err := manager.ExecuteHighPriority(func(dbExec dbmanager.DBExecutor) error {
		rows, err := dbExec.QueryContext(ctx, `
			SELECT node_uuid
			FROM sub_nodes_to_subscription_page_config
			WHERE subpage_config_uuid = ?
		`, subpageConfigUUID)
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			var nodeUUID string
			if err := rows.Scan(&nodeUUID); err != nil {
				return err
			}
			nodeUUID = strings.TrimSpace(nodeUUID)
			if nodeUUID != "" {
				nodeUUIDs = append(nodeUUIDs, nodeUUID)
			}
		}

		return rows.Err()
	})
	if err != nil {
		return nil, err
	}

	return nodeUUIDs, nil
}

func normalizeSubpageConfigName(name string) (string, error) {
	name = strings.TrimSpace(name)
	if len(name) < 2 {
		return "", fmt.Errorf("name must be at least 2 characters")
	}
	if len(name) > 30 {
		return "", fmt.Errorf("name must be less than 30 characters")
	}
	if !subpageConfigNameRegex.MatchString(name) {
		return "", fmt.Errorf("name can only contain letters, numbers, underscores, dashes and spaces")
	}
	return name, nil
}

func randomSuffix(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyz0123456789"
	buf := make([]byte, n)
	if _, err := rand.Read(buf); err != nil {
		return "00000"
	}
	for i := range buf {
		buf[i] = letters[int(buf[i])%len(letters)]
	}
	return string(buf)
}

func isUniqueNameError(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "23505"
	}
	return false
}
