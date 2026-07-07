package nodes

import (
	"context"
	"database/sql"
	"encoding/json"
	"sort"
	"time"

	dbmanager "exodus/internal/db/manager"
	"exodus/internal/dbutil"
	"exodus/internal/httpapi/shared"
)

func getAllNodeRecords(ctx context.Context, manager *dbmanager.DatabaseManager) ([]nodeRecord, error) {
	var nodes []nodeRecord
	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
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

func getNodeByUUID(ctx context.Context, manager *dbmanager.DatabaseManager, nodeUUID string) (nodeRecord, error) {
	var node nodeRecord
	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
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

func getNodeInbounds(ctx context.Context, manager *dbmanager.DatabaseManager, nodeUUIDs []string) (map[string][]configProfileInboundResponse, error) {
	result := make(map[string][]configProfileInboundResponse)
	if len(nodeUUIDs) == 0 {
		return result, nil
	}
	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
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

func getProviders(ctx context.Context, manager *dbmanager.DatabaseManager, providerUUIDs []string) (map[string]*providerResponse, error) {
	result := make(map[string]*providerResponse)
	if len(providerUUIDs) == 0 {
		return result, nil
	}
	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
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

func replaceNodeInboundsTx(ctx context.Context, tx dbmanager.TxExecutor, nodeUUID string, inboundUUIDs []string) error {
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

func getNodeTags(ctx context.Context, manager *dbmanager.DatabaseManager) ([]string, error) {
	tags := make([]string, 0)
	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
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
