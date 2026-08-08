package nodeplugins

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

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
	UUID         *string          `json:"uuid,omitempty"`
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

var defaultPluginConfig = json.RawMessage(`{"ingressFilter":{"enabled":false,"blockedIps":[]},"egressFilter":{"enabled":false,"blockedIps":[],"blockedPorts":[]},"haproxyAuth":{"inboundTags":[]},"sharedLists":[]}`)

const haproxyAllInboundTags = "*"

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

	if _, ok := obj["sharedLists"]; !ok {
		obj["sharedLists"] = []any{}
	}

	haproxyAuth, err := normalizeHaproxyAuthConfig(obj["haproxyAuth"])
	if err != nil {
		return nil, err
	}
	obj["haproxyAuth"] = haproxyAuth

	normalized, err := json.Marshal(obj)
	if err != nil {
		return nil, fmt.Errorf("pluginConfig cannot be encoded")
	}
	return normalized, nil
}

func normalizeHaproxyAuthConfig(raw any) (map[string]any, error) {
	result := map[string]any{"inboundTags": []string{}}
	if raw == nil {
		return result, nil
	}

	obj, ok := raw.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("haproxyAuth must be a JSON object")
	}

	if rawTags, ok := obj["inboundTags"]; ok {
		values, ok := rawTags.([]any)
		if !ok {
			return nil, fmt.Errorf("haproxyAuth.inboundTags must be an array")
		}

		tags := make([]string, 0, len(values))
		for _, item := range values {
			value, ok := item.(string)
			if !ok {
				return nil, fmt.Errorf("haproxyAuth.inboundTags items must be strings")
			}
			tags = append(tags, value)
		}
		result["inboundTags"] = tags
		return result, nil
	}

	if enabled, ok := obj["enabled"].(bool); ok && enabled {
		result["inboundTags"] = []string{"*"}
	}
	return result, nil
}

func normalizeHaproxyInboundTags(raw []string) []string {
	seen := make(map[string]struct{}, len(raw))
	result := make([]string, 0, len(raw))
	for _, item := range raw {
		value := strings.TrimSpace(item)
		if value == "" {
			continue
		}
		if value == haproxyAllInboundTags {
			return []string{haproxyAllInboundTags}
		}
		if _, ok := seen[value]; ok {
			continue
		}
		seen[value] = struct{}{}
		result = append(result, value)
	}
	return result
}

func loadPlugins(ctx context.Context, db *sql.DB) ([]nodePlugin, error) {
	rows, err := db.QueryContext(ctx, `
		SELECT uuid::text, name, plugin_config::text, view_position, created_at, updated_at
		FROM node_plugin
		ORDER BY view_position ASC, created_at ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	plugins := make([]nodePlugin, 0)
	for rows.Next() {
		plugin, scanErr := scanPlugin(rows)
		if scanErr != nil {
			return nil, scanErr
		}
		plugins = append(plugins, plugin)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return plugins, nil
}

func loadPluginByUUID(ctx context.Context, db *sql.DB, pluginUUID string) (nodePlugin, error) {
	var plugin nodePlugin
	row := db.QueryRowContext(ctx, `
		SELECT uuid::text, name, plugin_config::text, view_position, created_at, updated_at
		FROM node_plugin
		WHERE uuid::text = $1
	`, pluginUUID)
	err := scanPluginRow(row, &plugin)
	return plugin, err
}

func createPlugin(ctx context.Context, db *sql.DB, name string, configJSON json.RawMessage) (nodePlugin, error) {
	var plugin nodePlugin
	row := db.QueryRowContext(ctx, `
		INSERT INTO node_plugin (name, plugin_config)
		VALUES ($1, $2::jsonb)
		RETURNING uuid::text, name, plugin_config::text, view_position, created_at, updated_at
	`, name, string(configJSON))
	err := scanPluginRow(row, &plugin)
	return plugin, err
}

func updatePlugin(ctx context.Context, db *sql.DB, pluginUUID string, name *string, configJSON *json.RawMessage, viewPosition *int) (nodePlugin, error) {
	current, err := loadPluginByUUID(ctx, db, pluginUUID)
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
	row := db.QueryRowContext(ctx, `
		UPDATE node_plugin
		SET name = $1, plugin_config = $2::jsonb, view_position = $3, updated_at = CURRENT_TIMESTAMP
		WHERE uuid::text = $4
		RETURNING uuid::text, name, plugin_config::text, view_position, created_at, updated_at
	`, nextName, string(nextConfig), nextPosition, pluginUUID)
	err = scanPluginRow(row, &plugin)
	return plugin, err
}

func deletePlugin(ctx context.Context, db *sql.DB, pluginUUID string) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback()
	}()
	if _, err = tx.ExecContext(ctx, `UPDATE nodes SET active_plugin_uuid = NULL WHERE active_plugin_uuid::text = $1`, pluginUUID); err != nil {
		return err
	}
	if _, err = tx.ExecContext(ctx, `DELETE FROM node_plugin WHERE uuid::text = $1`, pluginUUID); err != nil {
		return err
	}
	return tx.Commit()
}

func reorderPlugins(ctx context.Context, db *sql.DB, req reorderRequest) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback()
	}()
	for _, item := range req.Items {
		if _, err := uuid.Parse(item.UUID); err != nil {
			return fmt.Errorf("invalid uuid")
		}
		if _, err := tx.ExecContext(ctx, `UPDATE node_plugin SET view_position = $1, updated_at = CURRENT_TIMESTAMP WHERE uuid::text = $2`, item.ViewPosition, item.UUID); err != nil {
			return err
		}
	}
	return tx.Commit()
}

func clonePlugin(ctx context.Context, db *sql.DB, req cloneRequest) (nodePlugin, error) {
	source, err := loadPluginByUUID(ctx, db, strings.TrimSpace(req.CloneFromUUID))
	if err != nil {
		return nodePlugin{}, err
	}
	name := source.Name + " copy"
	if req.Name != nil && strings.TrimSpace(*req.Name) != "" {
		name = strings.TrimSpace(*req.Name)
	} else {
		name = fmt.Sprintf("%s copy %s", source.Name, time.Now().UTC().Format("20060102150405"))
	}
	return createPlugin(ctx, db, name, source.PluginConfig)
}

func ensureNodesExist(ctx context.Context, db *sql.DB, nodeUUIDs []string) error {
	var count int
	if err := db.QueryRowContext(ctx, `SELECT COUNT(*) FROM nodes WHERE uuid::text = ANY($1)`, nodeUUIDs).Scan(&count); err != nil {
		return err
	}
	if count != len(nodeUUIDs) {
		return sql.ErrNoRows
	}
	return nil
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
