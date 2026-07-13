package nodes

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	dbmanager "exodus/internal/db/manager"
	"exodus/internal/dbutil"
	"exodus/internal/httpapi/shared"
)

type NodeRepository struct {
	manager *dbmanager.DatabaseManager
}

func NewNodeRepository(manager *dbmanager.DatabaseManager) *NodeRepository {
	return &NodeRepository{manager: manager}
}

func (r *NodeRepository) getAllNodeRecords(ctx context.Context) ([]nodeRecord, error) {
	var nodes []nodeRecord
	err := r.manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		rows, err := db.QueryContext(ctx, `
			SELECT
				uuid, id, name, address, port, proxy_url, api_schema, api_path, grpc_auth_token, active_config_profile_uuid, active_plugin_uuid,
				is_connected, is_connecting, is_disabled, last_status_change, last_status_message,
				consumption_multiplier, node_consumption_multiplier,
				is_traffic_tracking_active, traffic_reset_day, traffic_limit_bytes, traffic_used_bytes,
				notify_percent, provider_uuid, view_position, country_code, tags, note,
				created_at, updated_at
			FROM nodes
			ORDER BY view_position ASC
		`)
		if err != nil {
			return err
		}
		defer rows.Close()
		for rows.Next() {
			node, scanErr := scanNodeRecord(rows)
			if scanErr != nil {
				return scanErr
			}
			nodes = append(nodes, node)
		}
		return rows.Err()
	})
	return nodes, err
}

func (r *NodeRepository) getNodeByUUID(ctx context.Context, nodeUUID string) (nodeRecord, error) {
	var node nodeRecord
	err := r.manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		row := db.QueryRowContext(ctx, `
			SELECT
				uuid, id, name, address, port, proxy_url, api_schema, api_path, grpc_auth_token, active_config_profile_uuid, active_plugin_uuid,
				is_connected, is_connecting, is_disabled, last_status_change, last_status_message,
				consumption_multiplier, node_consumption_multiplier,
				is_traffic_tracking_active, traffic_reset_day, traffic_limit_bytes, traffic_used_bytes,
				notify_percent, provider_uuid, view_position, country_code, tags, note,
				created_at, updated_at
			FROM nodes
			WHERE uuid = ?
		`, nodeUUID)
		var scanErr error
		node, scanErr = scanNodeRecord(row)
		return scanErr
	})
	return node, err
}

func scanNodeRecord(scanner shared.RowScanner) (nodeRecord, error) {
	var node nodeRecord
	var id sql.NullInt64
	var port sql.NullInt64
	var proxyURL sql.NullString
	var activeConfigProfileUUID sql.NullString
	var activePluginUUID sql.NullString
	var lastStatusChange sql.NullTime
	var lastStatusMessage sql.NullString
	var trafficResetDay sql.NullInt64
	var trafficLimitBytes sql.NullInt64
	var trafficUsedBytes sql.NullInt64
	var notifyPercent sql.NullInt64
	var providerUUID sql.NullString
	var tags dbutil.StringArray
	var note sql.NullString

	err := scanner.Scan(
		&node.UUID,
		&id,
		&node.Name,
		&node.Address,
		&port,
		&proxyURL,
		&node.APISchema,
		&node.APIPath,
		&node.GRPCAuthToken,
		&activeConfigProfileUUID,
		&activePluginUUID,
		&node.IsConnected,
		&node.IsConnecting,
		&node.IsDisabled,
		&lastStatusChange,
		&lastStatusMessage,
		&node.ConsumptionMultiplier,
		&node.NodeConsumptionMultiplier,
		&node.IsTrafficTrackingActive,
		&trafficResetDay,
		&trafficLimitBytes,
		&trafficUsedBytes,
		&notifyPercent,
		&providerUUID,
		&node.ViewPosition,
		&node.CountryCode,
		&tags,
		&note,
		&node.CreatedAt,
		&node.UpdatedAt,
	)
	if err != nil {
		return node, err
	}

	if id.Valid {
		node.ID = &id.Int64
	}
	if port.Valid {
		value := int(port.Int64)
		node.Port = &value
	}
	if proxyURL.Valid {
		node.ProxyURL = &proxyURL.String
	}
	if activeConfigProfileUUID.Valid {
		node.ActiveConfigProfileUUID = &activeConfigProfileUUID.String
	}
	if activePluginUUID.Valid {
		node.ActivePluginUUID = &activePluginUUID.String
	}
	if lastStatusChange.Valid {
		node.LastStatusChange = &lastStatusChange.Time
	}
	if lastStatusMessage.Valid {
		node.LastStatusMessage = &lastStatusMessage.String
	}
	if trafficResetDay.Valid {
		value := int(trafficResetDay.Int64)
		node.TrafficResetDay = &value
	}
	if trafficLimitBytes.Valid {
		node.TrafficLimitBytes = &trafficLimitBytes.Int64
	}
	if trafficUsedBytes.Valid {
		node.TrafficUsedBytes = &trafficUsedBytes.Int64
	}
	if notifyPercent.Valid {
		value := int(notifyPercent.Int64)
		node.NotifyPercent = &value
	}
	if providerUUID.Valid {
		node.ProviderUUID = &providerUUID.String
	}
	node.Tags = tags.Slice()
	if note.Valid {
		node.Note = &note.String
	}

	return node, nil
}

func (r *NodeRepository) getNodeInbounds(ctx context.Context, nodeUUIDs []string) (map[string][]configProfileInboundResponse, error) {
	result := make(map[string][]configProfileInboundResponse)
	if len(nodeUUIDs) == 0 {
		return result, nil
	}
	err := r.manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		rows, err := db.QueryContext(ctx, `
			SELECT
				cpitn.node_uuid,
				cpi.uuid, cpi.profile_uuid, cpi.tag, cpi.type, cpi.network, cpi.security, cpi.port, cpi.raw_inbound
			FROM config_profile_inbounds_to_nodes cpitn
			JOIN config_profile_inbounds cpi ON cpi.uuid = cpitn.config_profile_inbound_uuid
			WHERE cpitn.node_uuid = ANY(?)
			ORDER BY cpi.tag ASC
		`, nodeUUIDs)
		if err != nil {
			return err
		}
		defer rows.Close()
		for rows.Next() {
			var nodeUUID string
			var inbound configProfileInboundResponse
			var network sql.NullString
			var security sql.NullString
			var port sql.NullInt64
			var rawInbound []byte
			if err := rows.Scan(
				&nodeUUID,
				&inbound.UUID,
				&inbound.ProfileUUID,
				&inbound.Tag,
				&inbound.Type,
				&network,
				&security,
				&port,
				&rawInbound,
			); err != nil {
				return err
			}
			if network.Valid {
				inbound.Network = &network.String
			}
			if security.Valid {
				inbound.Security = &security.String
			}
			if port.Valid {
				value := int(port.Int64)
				inbound.Port = &value
			}
			inbound.RawInbound = json.RawMessage(rawInbound)
			result[nodeUUID] = append(result[nodeUUID], inbound)
		}
		return rows.Err()
	})
	return result, err
}

func (r *NodeRepository) getProviders(ctx context.Context, providerUUIDs []string) (map[string]*providerResponse, error) {
	result := make(map[string]*providerResponse)
	if len(providerUUIDs) == 0 {
		return result, nil
	}
	err := r.manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		rows, err := db.QueryContext(ctx, `
			SELECT uuid, name, favicon_link, login_url, created_at, updated_at
			FROM infra_providers
			WHERE uuid = ANY(?)
		`, providerUUIDs)
		if err != nil {
			return err
		}
		defer rows.Close()
		for rows.Next() {
			var item providerResponse
			var favicon sql.NullString
			var loginURL sql.NullString
			var createdAt time.Time
			var updatedAt time.Time
			if err := rows.Scan(&item.UUID, &item.Name, &favicon, &loginURL, &createdAt, &updatedAt); err != nil {
				return err
			}
			if favicon.Valid {
				item.FaviconLink = &favicon.String
			}
			if loginURL.Valid {
				item.LoginURL = &loginURL.String
			}
			item.CreatedAt = &createdAt
			item.UpdatedAt = &updatedAt
			result[item.UUID] = &item
		}
		return rows.Err()
	})
	return result, err
}

func (r *NodeRepository) replaceNodeInboundsTx(ctx context.Context, tx dbmanager.TxExecutor, nodeUUID string, inboundUUIDs []string) error {
	if _, err := tx.ExecContext(ctx, `DELETE FROM config_profile_inbounds_to_nodes WHERE node_uuid = ?`, nodeUUID); err != nil {
		return err
	}
	for _, inboundUUID := range dedupeStrings(inboundUUIDs) {
		if _, err := tx.ExecContext(ctx, `
			INSERT INTO config_profile_inbounds_to_nodes (config_profile_inbound_uuid, node_uuid)
			VALUES (?, ?)
		`, inboundUUID, nodeUUID); err != nil {
			return err
		}
	}
	return nil
}

func (r *NodeRepository) getNodeTags(ctx context.Context) ([]string, error) {
	tags := make([]string, 0)
	err := r.manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		rows, err := db.QueryContext(ctx, `SELECT DISTINCT unnest(tags) AS tag FROM nodes ORDER BY tag ASC`)
		if err != nil {
			return err
		}
		defer rows.Close()
		for rows.Next() {
			var tag string
			if err := rows.Scan(&tag); err != nil {
				return err
			}
			if tag != "" {
				tags = append(tags, tag)
			}
		}
		return rows.Err()
	})
	if err != nil {
		return nil, err
	}
	sort.Strings(tags)
	return dedupeStrings(tags), nil
}

func (r *NodeRepository) createNode(ctx context.Context, nodeUUID string, req createNodeRequest, grpcAuthToken string, now time.Time) error {
	return r.manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		tx, err := db.BeginTx(ctx, nil)
		if err != nil {
			return err
		}

		_, err = tx.ExecContext(ctx, `
			INSERT INTO nodes (
				uuid, name, address, port, proxy_url, api_schema, api_path, grpc_auth_token, active_config_profile_uuid, active_plugin_uuid,
				is_connected, is_connecting, is_disabled, last_status_change, last_status_message,
				consumption_multiplier, node_consumption_multiplier,
				is_traffic_tracking_active, traffic_reset_day, traffic_limit_bytes, traffic_used_bytes,
				notify_percent, provider_uuid, country_code, tags, note, created_at, updated_at
			) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		`,
			nodeUUID,
			strings.TrimSpace(req.Name),
			strings.TrimSpace(req.Address),
			req.Port,
			normalizeNullableString(req.ProxyURL),
			normalizeAPISchema(req.APISchema),
			normalizeAPIPath(req.APIPath),
			grpcAuthToken,
			req.ConfigProfile.ActiveConfigProfileUUID,
			normalizeNullableString(req.ActivePluginUUID),
			false,
			false,
			false,
			nil,
			nil,
			toNanoMultiplier(coalesceFloat(req.ConsumptionMultiplier, 1)),
			toNanoMultiplier(coalesceFloat(req.NodeConsumptionMultiplier, 1)),
			coalesceBool(req.IsTrafficTrackingActive, false),
			coalesceInt(req.TrafficResetDay, 1),
			coalesceInt64(req.TrafficLimitBytes, 0),
			0,
			coalesceInt(req.NotifyPercent, 0),
			normalizeNullableString(req.ProviderUUID),
			normalizeCountryCode(req.CountryCode),
			normalizeTags(req.Tags),
			normalizeNullableString(req.Note),
			now,
			now,
		)
		if err != nil {
			_ = tx.Rollback()
			errStr := err.Error()
			if strings.Contains(errStr, "nodes_name_key") {
				return fmt.Errorf("node with this name already exists")
			}
			if strings.Contains(errStr, "nodes_address_key") {
				return fmt.Errorf("node with this address already exists")
			}
			return err
		}

		if err := r.replaceNodeInboundsTx(ctx, tx, nodeUUID, req.ConfigProfile.ActiveInbounds); err != nil {
			_ = tx.Rollback()
			return err
		}

		return tx.Commit()
	})
}

func (r *NodeRepository) updateNode(ctx context.Context, req updateNodeRequest, grpcAuthToken *string) error {
	return r.manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		tx, err := db.BeginTx(ctx, nil)
		if err != nil {
			return err
		}

		clauses := make([]string, 0)
		args := make([]any, 0)
		add := func(column string, value any) {
			clauses = append(clauses, fmt.Sprintf("%s = ?", column))
			args = append(args, value)
		}

		if req.Name != nil {
			add("name", strings.TrimSpace(*req.Name))
		}
		if req.Address != nil {
			add("address", strings.TrimSpace(*req.Address))
		}
		if req.Port != nil {
			add("port", *req.Port)
		}
		if req.ProxyURL.Set {
			if req.ProxyURL.Value == nil || strings.TrimSpace(*req.ProxyURL.Value) == "" {
				clauses = append(clauses, "proxy_url = NULL")
			} else {
				add("proxy_url", strings.TrimSpace(*req.ProxyURL.Value))
			}
		}
		if req.APISchema != nil {
			add("api_schema", normalizeAPISchema(req.APISchema))
		}
		if req.APIPath != nil {
			add("api_path", normalizeAPIPath(req.APIPath))
		}
		if grpcAuthToken != nil {
			add("grpc_auth_token", *grpcAuthToken)
		}
		if req.IsTrafficTrackingActive != nil {
			add("is_traffic_tracking_active", *req.IsTrafficTrackingActive)
		}
		if req.TrafficLimitBytes != nil {
			add("traffic_limit_bytes", *req.TrafficLimitBytes)
		}
		if req.NotifyPercent != nil {
			add("notify_percent", *req.NotifyPercent)
		}
		if req.TrafficResetDay != nil {
			add("traffic_reset_day", *req.TrafficResetDay)
		}
		if req.CountryCode != nil {
			add("country_code", strings.ToUpper(strings.TrimSpace(*req.CountryCode)))
		}
		if req.ConsumptionMultiplier != nil {
			add("consumption_multiplier", toNanoMultiplier(*req.ConsumptionMultiplier))
		}
		if req.NodeConsumptionMultiplier != nil {
			add("node_consumption_multiplier", toNanoMultiplier(*req.NodeConsumptionMultiplier))
		}
		if req.Tags != nil {
			add("tags", normalizeTags(*req.Tags))
		}
		if req.Note.Set {
			if req.Note.Value == nil || strings.TrimSpace(*req.Note.Value) == "" {
				clauses = append(clauses, "note = NULL")
			} else {
				add("note", strings.TrimSpace(*req.Note.Value))
			}
		}
		if req.ProviderUUID.Set {
			if req.ProviderUUID.Value == nil || strings.TrimSpace(*req.ProviderUUID.Value) == "" {
				clauses = append(clauses, "provider_uuid = NULL")
			} else {
				add("provider_uuid", strings.TrimSpace(*req.ProviderUUID.Value))
			}
		}
		if req.ActivePluginUUID.Set {
			if req.ActivePluginUUID.Value == nil || strings.TrimSpace(*req.ActivePluginUUID.Value) == "" {
				clauses = append(clauses, "active_plugin_uuid = NULL")
			} else {
				add("active_plugin_uuid", strings.TrimSpace(*req.ActivePluginUUID.Value))
			}
		}
		if req.ConfigProfile != nil {
			add("active_config_profile_uuid", req.ConfigProfile.ActiveConfigProfileUUID)
		}

		if len(clauses) > 0 {
			args = append(args, req.UUID)
			query := fmt.Sprintf("UPDATE nodes SET %s, updated_at = CURRENT_TIMESTAMP WHERE uuid = ?", strings.Join(clauses, ", "))
			result, err := tx.ExecContext(ctx, query, args...)
			if err != nil {
				_ = tx.Rollback()
				errStr := err.Error()
				if strings.Contains(errStr, "nodes_name_key") {
					return fmt.Errorf("node with this name already exists")
				}
				if strings.Contains(errStr, "nodes_address_key") {
					return fmt.Errorf("node with this address already exists")
				}
				return err
			}
			rows, err := result.RowsAffected()
			if err != nil {
				_ = tx.Rollback()
				return err
			}
			if rows == 0 {
				_ = tx.Rollback()
				return sql.ErrNoRows
			}
		}

		if req.ConfigProfile != nil {
			if err := r.replaceNodeInboundsTx(ctx, tx, req.UUID, req.ConfigProfile.ActiveInbounds); err != nil {
				_ = tx.Rollback()
				return err
			}
		}

		return tx.Commit()
	})
}

func (r *NodeRepository) deleteNode(ctx context.Context, nodeUUID string) error {
	return r.manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		result, err := db.ExecContext(ctx, `DELETE FROM nodes WHERE uuid = ?`, nodeUUID)
		if err != nil {
			return err
		}
		rows, err := result.RowsAffected()
		if err != nil {
			return err
		}
		if rows == 0 {
			return sql.ErrNoRows
		}
		return nil
	})
}

func (r *NodeRepository) resetNodeTraffic(ctx context.Context, nodeUUID string, node nodeRecord) error {
	return r.manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		tx, err := db.BeginTx(ctx, nil)
		if err != nil {
			return err
		}
		_, err = tx.ExecContext(ctx, `
			INSERT INTO nodes_traffic_usage_history (node_uuid, traffic_bytes, reset_at)
			VALUES (?, ?, ?)
		`, nodeUUID, coalesceInt64Ptr(node.TrafficUsedBytes), time.Now().UTC())
		if err != nil {
			_ = tx.Rollback()
			return err
		}
		if _, err := tx.ExecContext(ctx, `UPDATE nodes SET traffic_used_bytes = 0, updated_at = CURRENT_TIMESTAMP WHERE uuid = ?`, nodeUUID); err != nil {
			_ = tx.Rollback()
			return err
		}
		return tx.Commit()
	})
}

func (r *NodeRepository) reorderNodes(ctx context.Context, items []reorderNodeItem) error {
	return r.manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		tx, err := db.BeginTx(ctx, nil)
		if err != nil {
			return err
		}
		for _, item := range items {
			if _, err := tx.ExecContext(ctx, `UPDATE nodes SET view_position = ? WHERE uuid = ?`, item.ViewPosition, item.UUID); err != nil {
				_ = tx.Rollback()
				return err
			}
		}
		if _, err := tx.ExecContext(ctx, `SELECT setval('nodes_view_position_seq', (SELECT COALESCE(MAX(view_position), 0) FROM nodes) + 1)`); err != nil {
			_ = tx.Rollback()
			return err
		}
		return tx.Commit()
	})
}

func (r *NodeRepository) enableNodeRecord(ctx context.Context, nodeUUID string, node nodeRecord, inbounds []configProfileInboundResponse) error {
	return r.manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		if node.ActiveConfigProfileUUID == nil || len(inbounds) == 0 {
			_, execErr := db.ExecContext(ctx, `
				UPDATE nodes
				SET is_disabled = true, active_config_profile_uuid = NULL, is_connecting = false,
					is_connected = false, last_status_message = NULL, last_status_change = ?
				WHERE uuid = ?
			`, time.Now().UTC(), nodeUUID)
			return execErr
		}
		_, execErr := db.ExecContext(ctx, `UPDATE nodes SET is_disabled = false, updated_at = CURRENT_TIMESTAMP WHERE uuid = ?`, nodeUUID)
		return execErr
	})
}

func (r *NodeRepository) disableNodeRecord(ctx context.Context, nodeUUID string, node nodeRecord, inbounds []configProfileInboundResponse) error {
	return r.manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		if node.ActiveConfigProfileUUID == nil || len(inbounds) == 0 {
			if _, execErr := db.ExecContext(ctx, `UPDATE nodes SET active_config_profile_uuid = NULL WHERE uuid = ?`, nodeUUID); execErr != nil {
				return execErr
			}
		}
		_, execErr := db.ExecContext(ctx, `
			UPDATE nodes
			SET is_disabled = true, is_connecting = false, is_connected = false,
				last_status_message = NULL, last_status_change = ?,
				updated_at = CURRENT_TIMESTAMP
			WHERE uuid = ?
		`, time.Now().UTC(), nodeUUID)
		return execErr
	})
}

func (r *NodeRepository) bulkProfileModification(ctx context.Context, uuids []string, activeConfigProfileUUID string, activeInbounds []string) error {
	return r.manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		tx, err := db.BeginTx(ctx, nil)
		if err != nil {
			return err
		}
		for _, nodeUUID := range uuids {
			if _, err := tx.ExecContext(ctx, `UPDATE nodes SET active_config_profile_uuid = ?, updated_at = CURRENT_TIMESTAMP WHERE uuid = ?`, activeConfigProfileUUID, nodeUUID); err != nil {
				_ = tx.Rollback()
				return err
			}
			if err := r.replaceNodeInboundsTx(ctx, tx, nodeUUID, activeInbounds); err != nil {
				_ = tx.Rollback()
				return err
			}
		}
		return tx.Commit()
	})
}

func (r *NodeRepository) bulkUpdateNodes(ctx context.Context, uuids []string, clauses []string, args []any) error {
	if len(clauses) == 0 {
		return nil
	}
	query := fmt.Sprintf("UPDATE nodes SET %s, updated_at = CURRENT_TIMESTAMP WHERE uuid = ANY(?)", strings.Join(clauses, ", "))
	args = append(args, uuids)
	return r.manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		_, execErr := db.ExecContext(ctx, query, args...)
		return execErr
	})
}

func (r *NodeRepository) ensureConfigProfileInbounds(ctx context.Context, profileUUID string, inboundUUIDs []string) error {
	return r.manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		var exists int
		if err := db.QueryRowContext(ctx, `SELECT 1 FROM config_profiles WHERE uuid = ?`, profileUUID).Scan(&exists); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return errConfigProfileNotFound
			}
			return err
		}

		found := make(map[string]struct{}, len(inboundUUIDs))
		rows, err := db.QueryContext(ctx, `
			SELECT uuid
			FROM config_profile_inbounds
			WHERE profile_uuid = ? AND uuid = ANY(?)
		`, profileUUID, inboundUUIDs)
		if err != nil {
			return err
		}
		defer rows.Close()
		for rows.Next() {
			var inboundUUID string
			if err := rows.Scan(&inboundUUID); err != nil {
				return err
			}
			found[inboundUUID] = struct{}{}
		}
		if err := rows.Err(); err != nil {
			return err
		}
		for _, inboundUUID := range inboundUUIDs {
			if _, ok := found[inboundUUID]; !ok {
				return errConfigProfileInboundInvalid
			}
		}
		return nil
	})
}
