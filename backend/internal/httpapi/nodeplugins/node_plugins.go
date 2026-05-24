package nodeplugins

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"exodus/internal/config"
	dbmanager "exodus/internal/db/manager"
	"exodus/internal/httpapi/shared"
	monitor "exodus/internal/nodes"

	"github.com/google/uuid"
)

type responseEnvelope[T any] struct {
	Response T `json:"response"`
}

type nodePlugin struct {
	UUID         string          `json:"uuid"`
	Name         string          `json:"name"`
	PluginConfig json.RawMessage `json:"pluginConfig"`
	ViewPosition int             `json:"viewPosition"`
	CreatedAt    time.Time       `json:"createdAt"`
	UpdatedAt    time.Time       `json:"updatedAt"`
}

type listPayload struct {
	NodePlugins []nodePlugin `json:"nodePlugins"`
	Total       int          `json:"total"`
}

type createRequest struct {
	Name         string          `json:"name"`
	PluginConfig json.RawMessage `json:"pluginConfig,omitempty"`
}

type updateRequest struct {
	Name         *string          `json:"name,omitempty"`
	PluginConfig *json.RawMessage `json:"pluginConfig,omitempty"`
	ViewPosition *int             `json:"viewPosition,omitempty"`
}

type reorderRequest struct {
	Items []struct {
		UUID         string `json:"uuid"`
		ViewPosition int    `json:"viewPosition"`
	} `json:"items"`
}

type cloneRequest struct {
	CloneFromUUID string  `json:"cloneFromUuid"`
	Name          *string `json:"name,omitempty"`
}

type executorRequest struct {
	Command     executorCommand `json:"command"`
	TargetNodes struct {
		Target    string   `json:"target"`
		NodeUUIDs []string `json:"nodeUuids"`
	} `json:"targetNodes"`
}

type executorCommand struct {
	Command string          `json:"command"`
	Raw     json.RawMessage `json:"-"`
}

func (c *executorCommand) UnmarshalJSON(data []byte) error {
	type alias executorCommand
	var decoded alias
	if err := json.Unmarshal(data, &decoded); err != nil {
		return err
	}
	*c = executorCommand(decoded)
	c.Raw = append(c.Raw[:0], data...)
	return nil
}

var defaultPluginConfig = json.RawMessage(`{"ingressFilter":{"enabled":false,"blockedIps":[]},"egressFilter":{"enabled":false,"blockedIps":[],"blockedPorts":[]},"haproxyAuth":{"enabled":false},"sharedLists":[]}`)

func Handler(manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/api/node-plugins")
		path = strings.Trim(path, "/")

		switch {
		case path == "":
			switch r.Method {
			case http.MethodGet:
				handleList(w, r, manager, cfg)
			case http.MethodPost:
				handleCreate(w, r, manager, cfg)
			default:
				shared.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
			}
		case path == "executor":
			handleExecutor(w, r, manager, cfg)
		case strings.HasPrefix(path, "actions/"):
			handleAction(w, r, manager, cfg, strings.TrimPrefix(path, "actions/"))
		default:
			handleByUUID(w, r, manager, cfg, path)
		}
	}
}

func handleList(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) {
	plugins, err := loadPlugins(r.Context(), manager)
	if err != nil {
		cfg.Logger.Error("Failed to load node plugins", "error", err)
		shared.WriteJSONError(w, http.StatusInternalServerError, "failed to load node plugins")
		return
	}
	shared.WriteJSON(w, http.StatusOK, responseEnvelope[listPayload]{
		Response: listPayload{NodePlugins: plugins, Total: len(plugins)},
	})
}

func handleCreate(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) {
	var req createRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		shared.WriteJSONError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	name := strings.TrimSpace(req.Name)
	if name == "" {
		shared.WriteJSONError(w, http.StatusBadRequest, "name is required")
		return
	}
	configJSON, err := normalizePluginConfig(req.PluginConfig)
	if err != nil {
		shared.WriteJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	plugin, err := createPlugin(r.Context(), manager, name, configJSON)
	if err != nil {
		cfg.Logger.Error("Failed to create node plugin", "error", err)
		shared.WriteJSONError(w, http.StatusInternalServerError, "failed to create node plugin")
		return
	}
	shared.WriteJSON(w, http.StatusCreated, responseEnvelope[nodePlugin]{Response: plugin})
}

func handleByUUID(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig, rawPath string) {
	parts := strings.Split(strings.Trim(rawPath, "/"), "/")
	if len(parts) == 0 || parts[0] == "" {
		shared.WriteJSONError(w, http.StatusNotFound, "not found")
		return
	}
	pluginUUID := parts[0]
	if _, err := uuid.Parse(pluginUUID); err != nil {
		shared.WriteJSONError(w, http.StatusBadRequest, "invalid uuid")
		return
	}

	if len(parts) > 1 {
		shared.WriteJSONError(w, http.StatusNotFound, "not found")
		return
	}

	switch r.Method {
	case http.MethodGet:
		plugin, err := loadPluginByUUID(r.Context(), manager, pluginUUID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				shared.WriteJSONError(w, http.StatusNotFound, "node plugin not found")
				return
			}
			cfg.Logger.Error("Failed to load node plugin", "uuid", pluginUUID, "error", err)
			shared.WriteJSONError(w, http.StatusInternalServerError, "failed to load node plugin")
			return
		}
		shared.WriteJSON(w, http.StatusOK, responseEnvelope[nodePlugin]{Response: plugin})
	case http.MethodPatch:
		handleUpdate(w, r, manager, cfg, pluginUUID)
	case http.MethodDelete:
		handleDelete(w, r, manager, cfg, pluginUUID)
	default:
		shared.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func handleUpdate(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig, pluginUUID string) {
	var req updateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		shared.WriteJSONError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	var name *string
	if req.Name != nil {
		trimmed := strings.TrimSpace(*req.Name)
		if trimmed == "" {
			shared.WriteJSONError(w, http.StatusBadRequest, "name cannot be empty")
			return
		}
		name = &trimmed
	}

	var configJSON *json.RawMessage
	if req.PluginConfig != nil {
		normalized, err := normalizePluginConfig(*req.PluginConfig)
		if err != nil {
			shared.WriteJSONError(w, http.StatusBadRequest, err.Error())
			return
		}
		configJSON = &normalized
	}

	plugin, err := updatePlugin(r.Context(), manager, pluginUUID, name, configJSON, req.ViewPosition)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			shared.WriteJSONError(w, http.StatusNotFound, "node plugin not found")
			return
		}
		cfg.Logger.Error("Failed to update node plugin", "uuid", pluginUUID, "error", err)
		shared.WriteJSONError(w, http.StatusInternalServerError, "failed to update node plugin")
		return
	}
	shared.WriteJSON(w, http.StatusOK, responseEnvelope[nodePlugin]{Response: plugin})
}

func handleDelete(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig, pluginUUID string) {
	if err := deletePlugin(r.Context(), manager, pluginUUID); err != nil {
		cfg.Logger.Error("Failed to delete node plugin", "uuid", pluginUUID, "error", err)
		shared.WriteJSONError(w, http.StatusInternalServerError, "failed to delete node plugin")
		return
	}
	shared.WriteJSON(w, http.StatusOK, responseEnvelope[map[string]bool]{Response: map[string]bool{"isDeleted": true}})
}

func handleExecutor(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) {
	if r.Method != http.MethodPost {
		shared.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var req executorRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		shared.WriteJSONError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	command := strings.TrimSpace(req.Command.Command)
	switch command {
	case "blockIps", "unblockIps", "recreateTables":
	default:
		shared.WriteJSONError(w, http.StatusBadRequest, "unsupported executor command")
		return
	}
	if req.TargetNodes.Target != "specificNodes" {
		shared.WriteJSONError(w, http.StatusBadRequest, "targetNodes.target must be specificNodes")
		return
	}
	targetNodeUUIDs := normalizeUUIDList(req.TargetNodes.NodeUUIDs)
	if len(targetNodeUUIDs) == 0 {
		shared.WriteJSONError(w, http.StatusBadRequest, "nodeUuids are required")
		return
	}
	for _, nodeUUID := range targetNodeUUIDs {
		if _, err := uuid.Parse(nodeUUID); err != nil {
			shared.WriteJSONError(w, http.StatusBadRequest, "invalid node uuid")
			return
		}
	}

	if err := ensureNodesExist(r.Context(), manager, targetNodeUUIDs); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			shared.WriteJSONError(w, http.StatusBadRequest, "one or more nodes were not found")
			return
		}
		cfg.Logger.Error("Failed to validate node plugin executor targets", "error", err)
		shared.WriteJSONError(w, http.StatusInternalServerError, "failed to validate targets")
		return
	}

	cfg.Logger.Info(
		"Node plugin executor command accepted",
		"command", command,
		"nodes", strings.Join(targetNodeUUIDs, ","),
	)
	monitor.RequestNodeDeploy(true, targetNodeUUIDs...)
	shared.WriteJSON(w, http.StatusOK, responseEnvelope[map[string]bool]{Response: map[string]bool{"eventSent": true}})
}

func handleAction(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig, action string) {
	if r.Method != http.MethodPost {
		shared.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	switch strings.Trim(action, "/") {
	case "reorder":
		var req reorderRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			shared.WriteJSONError(w, http.StatusBadRequest, "invalid JSON body")
			return
		}
		if len(req.Items) == 0 {
			shared.WriteJSONError(w, http.StatusBadRequest, "items are required")
			return
		}
		if err := reorderPlugins(r.Context(), manager, req); err != nil {
			cfg.Logger.Error("Failed to reorder node plugins", "error", err)
			shared.WriteJSONError(w, http.StatusInternalServerError, "failed to reorder node plugins")
			return
		}
		plugins, err := loadPlugins(r.Context(), manager)
		if err != nil {
			cfg.Logger.Error("Failed to load reordered node plugins", "error", err)
			shared.WriteJSONError(w, http.StatusInternalServerError, "failed to load node plugins")
			return
		}
		shared.WriteJSON(w, http.StatusOK, responseEnvelope[listPayload]{Response: listPayload{NodePlugins: plugins, Total: len(plugins)}})
	case "clone":
		var req cloneRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			shared.WriteJSONError(w, http.StatusBadRequest, "invalid JSON body")
			return
		}
		if _, err := uuid.Parse(strings.TrimSpace(req.CloneFromUUID)); err != nil {
			shared.WriteJSONError(w, http.StatusBadRequest, "invalid cloneFromUuid")
			return
		}
		plugin, err := clonePlugin(r.Context(), manager, req)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				shared.WriteJSONError(w, http.StatusNotFound, "node plugin not found")
				return
			}
			cfg.Logger.Error("Failed to clone node plugin", "error", err)
			shared.WriteJSONError(w, http.StatusInternalServerError, "failed to clone node plugin")
			return
		}
		shared.WriteJSON(w, http.StatusCreated, responseEnvelope[nodePlugin]{Response: plugin})
	default:
		shared.WriteJSONError(w, http.StatusNotFound, "not found")
	}
}

func normalizePluginConfig(raw json.RawMessage) (json.RawMessage, error) {
	if len(raw) == 0 || string(raw) == "null" {
		return append(json.RawMessage(nil), defaultPluginConfig...), nil
	}

	var obj map[string]any
	if err := json.Unmarshal(raw, &obj); err != nil {
		return nil, fmt.Errorf("pluginConfig must be valid JSON object")
	}
	if obj == nil {
		return append(json.RawMessage(nil), defaultPluginConfig...), nil
	}
	if _, ok := obj["torrentBlocker"]; ok {
		return nil, fmt.Errorf("torrentBlocker plugin is not supported by sing-box core")
	}
	if _, ok := obj["connectionDrop"]; ok {
		return nil, fmt.Errorf("connectionDrop plugin is not supported by sing-box core")
	}

	defaults := map[string]any{}
	_ = json.Unmarshal(defaultPluginConfig, &defaults)
	for key, value := range defaults {
		if _, ok := obj[key]; !ok {
			obj[key] = value
		}
	}

	normalized, err := json.Marshal(obj)
	if err != nil {
		return nil, fmt.Errorf("pluginConfig cannot be encoded")
	}
	return normalized, nil
}

func loadPlugins(ctx context.Context, manager *dbmanager.DatabaseManager) ([]nodePlugin, error) {
	plugins := make([]nodePlugin, 0)
	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		rows, err := db.QueryContext(ctx, `
			SELECT uuid::text, name, plugin_config::text, view_position, created_at, updated_at
			FROM node_plugins
			ORDER BY view_position ASC, created_at ASC
		`)
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			plugin, scanErr := scanPlugin(rows)
			if scanErr != nil {
				return scanErr
			}
			plugins = append(plugins, plugin)
		}
		return rows.Err()
	})
	return plugins, err
}

func loadPluginByUUID(ctx context.Context, manager *dbmanager.DatabaseManager, pluginUUID string) (nodePlugin, error) {
	var plugin nodePlugin
	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		row := db.QueryRowContext(ctx, `
			SELECT uuid::text, name, plugin_config::text, view_position, created_at, updated_at
			FROM node_plugins
			WHERE uuid::text = ?
		`, pluginUUID)
		return scanPluginRow(row, &plugin)
	})
	return plugin, err
}

func createPlugin(ctx context.Context, manager *dbmanager.DatabaseManager, name string, configJSON json.RawMessage) (nodePlugin, error) {
	var plugin nodePlugin
	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		row := db.QueryRowContext(ctx, `
			INSERT INTO node_plugins (name, plugin_config)
			VALUES (?, ?::jsonb)
			RETURNING uuid::text, name, plugin_config::text, view_position, created_at, updated_at
		`, name, string(configJSON))
		return scanPluginRow(row, &plugin)
	})
	return plugin, err
}

func updatePlugin(ctx context.Context, manager *dbmanager.DatabaseManager, pluginUUID string, name *string, configJSON *json.RawMessage, viewPosition *int) (nodePlugin, error) {
	current, err := loadPluginByUUID(ctx, manager, pluginUUID)
	if err != nil {
		return nodePlugin{}, err
	}
	nextName := current.Name
	if name != nil {
		nextName = *name
	}
	nextConfig := current.PluginConfig
	if configJSON != nil {
		nextConfig = *configJSON
	}
	nextPosition := current.ViewPosition
	if viewPosition != nil {
		nextPosition = *viewPosition
	}

	var plugin nodePlugin
	err = manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		row := db.QueryRowContext(ctx, `
			UPDATE node_plugins
			SET name = ?, plugin_config = ?::jsonb, view_position = ?, updated_at = CURRENT_TIMESTAMP
			WHERE uuid::text = ?
			RETURNING uuid::text, name, plugin_config::text, view_position, created_at, updated_at
		`, nextName, string(nextConfig), nextPosition, pluginUUID)
		return scanPluginRow(row, &plugin)
	})
	return plugin, err
}

func deletePlugin(ctx context.Context, manager *dbmanager.DatabaseManager, pluginUUID string) error {
	return manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		tx, err := db.BeginTx(ctx, nil)
		if err != nil {
			return err
		}
		defer tx.Rollback()
		if _, err = tx.ExecContext(ctx, `UPDATE nodes SET active_plugin_uuid = NULL WHERE active_plugin_uuid::text = ?`, pluginUUID); err != nil {
			return err
		}
		if _, err = tx.ExecContext(ctx, `DELETE FROM node_plugins WHERE uuid::text = ?`, pluginUUID); err != nil {
			return err
		}
		return tx.Commit()
	})
}

func reorderPlugins(ctx context.Context, manager *dbmanager.DatabaseManager, req reorderRequest) error {
	return manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		tx, err := db.BeginTx(ctx, nil)
		if err != nil {
			return err
		}
		defer tx.Rollback()
		for _, item := range req.Items {
			if _, err := uuid.Parse(item.UUID); err != nil {
				return fmt.Errorf("invalid uuid")
			}
			if _, err := tx.ExecContext(ctx, `UPDATE node_plugins SET view_position = ?, updated_at = CURRENT_TIMESTAMP WHERE uuid::text = ?`, item.ViewPosition, item.UUID); err != nil {
				return err
			}
		}
		return tx.Commit()
	})
}

func clonePlugin(ctx context.Context, manager *dbmanager.DatabaseManager, req cloneRequest) (nodePlugin, error) {
	source, err := loadPluginByUUID(ctx, manager, strings.TrimSpace(req.CloneFromUUID))
	if err != nil {
		return nodePlugin{}, err
	}
	name := source.Name + " copy"
	if req.Name != nil && strings.TrimSpace(*req.Name) != "" {
		name = strings.TrimSpace(*req.Name)
	} else {
		name = fmt.Sprintf("%s copy %s", source.Name, time.Now().UTC().Format("20060102150405"))
	}
	return createPlugin(ctx, manager, name, source.PluginConfig)
}

func ensureNodesExist(ctx context.Context, manager *dbmanager.DatabaseManager, nodeUUIDs []string) error {
	return manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		var count int
		if err := db.QueryRowContext(ctx, `SELECT COUNT(*) FROM nodes WHERE uuid::text = ANY(?)`, nodeUUIDs).Scan(&count); err != nil {
			return err
		}
		if count != len(nodeUUIDs) {
			return sql.ErrNoRows
		}
		return nil
	})
}

func normalizeUUIDList(raw []string) []string {
	seen := make(map[string]struct{}, len(raw))
	result := make([]string, 0, len(raw))
	for _, item := range raw {
		value := strings.TrimSpace(item)
		if value == "" {
			continue
		}
		if _, exists := seen[value]; exists {
			continue
		}
		seen[value] = struct{}{}
		result = append(result, value)
	}
	return result
}

type pluginScanner interface {
	Scan(dest ...any) error
}

func scanPlugin(rows *sql.Rows) (nodePlugin, error) {
	var plugin nodePlugin
	err := scanPluginRow(rows, &plugin)
	return plugin, err
}

func scanPluginRow(scanner pluginScanner, plugin *nodePlugin) error {
	var rawConfig string
	if err := scanner.Scan(
		&plugin.UUID,
		&plugin.Name,
		&rawConfig,
		&plugin.ViewPosition,
		&plugin.CreatedAt,
		&plugin.UpdatedAt,
	); err != nil {
		return err
	}
	plugin.PluginConfig = json.RawMessage(rawConfig)
	return nil
}
