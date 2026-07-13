package subscriptionconnections

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"

	dbmanager "exodus/internal/db/manager"
	"exodus/internal/dbutil"
	"exodus/internal/httpapi/shared"
)

type SubscriptionConnectionRepository struct {
	manager *dbmanager.DatabaseManager
}

func NewSubscriptionConnectionRepository(manager *dbmanager.DatabaseManager) *SubscriptionConnectionRepository {
	return &SubscriptionConnectionRepository{manager: manager}
}

func (r *SubscriptionConnectionRepository) getAllNodeRecords(ctx context.Context) ([]nodeRecord, error) {
	var nodes []nodeRecord
	err := r.manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		rows, err := db.QueryContext(ctx, `
				SELECT
					n.uuid, n.id, n.name, n.address, n.public_domain, n.port, n.api_schema, n.api_path, n.grpc_auth_token, sns.subpage_config_uuid,
					n.is_connected, n.is_connecting, n.is_disabled,
					n.last_status_change, n.last_status_message,
					n.provider_uuid, n.view_position, n.tags, n.created_at, n.updated_at
			FROM sub_nodes n
			LEFT JOIN sub_nodes_to_subscription_page_config sns ON sns.node_uuid = n.uuid
			ORDER BY n.view_position ASC
		`)
		if err != nil {
			return err
		}
		defer rows.Close()
		for rows.Next() {
			node, scanErr := r.scanNodeRecord(rows)
			if scanErr != nil {
				return scanErr
			}
			nodes = append(nodes, node)
		}
		return rows.Err()
	})
	return nodes, err
}

func (r *SubscriptionConnectionRepository) getNodeByUUID(ctx context.Context, nodeUUID string) (nodeRecord, error) {
	var node nodeRecord
	err := r.manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		row := db.QueryRowContext(ctx, `
				SELECT
					n.uuid, n.id, n.name, n.address, n.public_domain, n.port, n.api_schema, n.api_path, n.grpc_auth_token, sns.subpage_config_uuid,
					n.is_connected, n.is_connecting, n.is_disabled,
					n.last_status_change, n.last_status_message,
					n.provider_uuid, n.view_position, n.tags, n.created_at, n.updated_at
			FROM sub_nodes n
			LEFT JOIN sub_nodes_to_subscription_page_config sns ON sns.node_uuid = n.uuid
			WHERE n.uuid = ?
		`, nodeUUID)
		var scanErr error
		node, scanErr = r.scanNodeRecord(row)
		return scanErr
	})
	return node, err
}

func (r *SubscriptionConnectionRepository) scanNodeRecord(scanner shared.RowScanner) (nodeRecord, error) {
	var node nodeRecord
	var id sql.NullInt64
	var port sql.NullInt64
	var publicDomain sql.NullString
	var lastStatusChange sql.NullTime
	var lastStatusMessage sql.NullString
	var providerUUID sql.NullString
	var subpageConfigUUID sql.NullString
	var tags dbutil.StringArray

	err := scanner.Scan(
		&node.UUID,
		&id,
		&node.Name,
		&node.Address,
		&publicDomain,
		&port,
		&node.APISchema,
		&node.APIPath,
		&node.GRPCAuthToken,
		&subpageConfigUUID,
		&node.IsConnected,
		&node.IsConnecting,
		&node.IsDisabled,
		&lastStatusChange,
		&lastStatusMessage,
		&providerUUID,
		&node.ViewPosition,
		&tags,
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
	if publicDomain.Valid {
		value := strings.TrimSpace(publicDomain.String)
		if value != "" {
			node.PublicDomain = &value
		}
	}
	if lastStatusChange.Valid {
		ts := lastStatusChange.Time
		node.LastStatusChange = &ts
	}
	if lastStatusMessage.Valid {
		msg := lastStatusMessage.String
		node.LastStatusMessage = &msg
	}
	if providerUUID.Valid {
		node.ProviderUUID = &providerUUID.String
	}
	if subpageConfigUUID.Valid {
		value := strings.TrimSpace(subpageConfigUUID.String)
		if value != "" {
			node.SubpageConfigUUID = &value
		}
	}
	node.APISchema = normalizeSubNodeSchema(node.APISchema)
	node.SingboxVersion = nil
	node.NodeVersion = nil
	node.SingboxUptime = "0"
	node.CPUCount = nil
	node.CPUModel = nil
	node.TotalRAM = nil
	node.Tags = tags.Slice()
	node.CountryCode = "XX"
	node.ConsumptionMultiplier = 1_000_000_000
	node.IsTrafficTrackingActive = false

	return node, nil
}

func (r *SubscriptionConnectionRepository) getProviders(ctx context.Context, providerUUIDs []string) (map[string]*providerResponse, error) {
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

func (r *SubscriptionConnectionRepository) getNodeTags(ctx context.Context) ([]string, error) {
	tags := make([]string, 0)
	err := r.manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		rows, err := db.QueryContext(ctx, `SELECT DISTINCT unnest(tags) AS tag FROM sub_nodes ORDER BY tag ASC`)
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

func (r *SubscriptionConnectionRepository) subpageConfigExists(ctx context.Context, configUUID string) (bool, error) {
	exists := false
	err := r.manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		var count int
		if err := db.QueryRowContext(ctx, `SELECT COUNT(*) FROM subscription_page_config WHERE uuid = ?`, configUUID).Scan(&count); err != nil {
			return err
		}
		exists = count > 0
		return nil
	})
	return exists, err
}

func (r *SubscriptionConnectionRepository) fetchSubpageConfigRaw(ctx context.Context, configUUID string) ([]byte, error) {
	var payload string
	err := r.manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		return db.QueryRowContext(ctx, `
			SELECT config
			FROM subscription_page_config
			WHERE uuid = ?
			LIMIT 1
		`, configUUID).Scan(&payload)
	})
	if err != nil {
		return nil, err
	}

	raw := []byte(strings.TrimSpace(payload))
	if len(raw) == 0 || !json.Valid(raw) {
		return nil, fmt.Errorf("invalid subpage config payload")
	}

	return raw, nil
}

func (r *SubscriptionConnectionRepository) countEnabledNodes(ctx context.Context) (int, error) {
	var enabledCount int
	err := r.manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		return db.QueryRowContext(ctx, `SELECT COUNT(*) FROM sub_nodes WHERE is_disabled = false`).Scan(&enabledCount)
	})
	return enabledCount, err
}

func (r *SubscriptionConnectionRepository) createNode(ctx context.Context, nodeUUID string, req createNodeRequest, schema, grpcAuthToken string, subpageConfigUUID *string, now time.Time) error {
	return r.manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		tx, err := db.BeginTx(ctx, nil)
		if err != nil {
			return err
		}

		_, err = tx.ExecContext(ctx, `
				INSERT INTO sub_nodes (
					uuid, name, address, public_domain, port, api_schema, api_path, grpc_auth_token,
					is_connected, is_connecting, is_disabled, provider_uuid, tags, created_at, updated_at
				) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
			`,
			nodeUUID,
			strings.TrimSpace(req.Name),
			strings.TrimSpace(req.Address),
			normalizePublicDomain(req.PublicDomain),
			req.Port,
			schema,
			normalizeAPIPath(req.APIPath),
			grpcAuthToken,
			false,
			false,
			false,
			normalizeNullableString(req.ProviderUUID),
			normalizeTags(req.Tags),
			now,
			now,
		)
		if err != nil {
			_ = tx.Rollback()
			return err
		}
		if subpageConfigUUID != nil {
			if _, err := tx.ExecContext(ctx, `
				INSERT INTO sub_nodes_to_subscription_page_config (node_uuid, subpage_config_uuid)
				VALUES (?, ?)
				ON CONFLICT (node_uuid) DO UPDATE
				SET subpage_config_uuid = EXCLUDED.subpage_config_uuid
			`, nodeUUID, *subpageConfigUUID); err != nil {
				_ = tx.Rollback()
				return err
			}
		}
		return tx.Commit()
	})
}

func (r *SubscriptionConnectionRepository) updateNode(ctx context.Context, nodeUUID string, clauses []string, args []any, subpageConfigSet bool, finalSubpageConfigUUID *string) error {
	return r.manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		tx, err := db.BeginTx(ctx, nil)
		if err != nil {
			return err
		}

		if len(clauses) > 0 {
			args = append(args, nodeUUID)
			query := fmt.Sprintf("UPDATE sub_nodes SET %s, updated_at = CURRENT_TIMESTAMP WHERE uuid = ?", strings.Join(clauses, ", "))
			result, err := tx.ExecContext(ctx, query, args...)
			if err != nil {
				_ = tx.Rollback()
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
		if subpageConfigSet {
			if finalSubpageConfigUUID == nil {
				if _, err := tx.ExecContext(ctx, `DELETE FROM sub_nodes_to_subscription_page_config WHERE node_uuid = ?`, nodeUUID); err != nil {
					_ = tx.Rollback()
					return err
				}
			} else {
				if _, err := tx.ExecContext(ctx, `
					INSERT INTO sub_nodes_to_subscription_page_config (node_uuid, subpage_config_uuid)
					VALUES (?, ?)
					ON CONFLICT (node_uuid) DO UPDATE
					SET subpage_config_uuid = EXCLUDED.subpage_config_uuid
				`, nodeUUID, *finalSubpageConfigUUID); err != nil {
					_ = tx.Rollback()
					return err
				}
			}
		}

		return tx.Commit()
	})
}

func (r *SubscriptionConnectionRepository) deleteNode(ctx context.Context, nodeUUID string) error {
	return r.manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		result, err := db.ExecContext(ctx, `DELETE FROM sub_nodes WHERE uuid = ?`, nodeUUID)
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

func (r *SubscriptionConnectionRepository) setNodeDisabled(ctx context.Context, nodeUUID string, disabled bool) error {
	return r.manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		var result sql.Result
		var err error
		if disabled {
			result, err = db.ExecContext(ctx, `
				UPDATE sub_nodes
				SET is_disabled = true, is_connecting = false, is_connected = false, updated_at = CURRENT_TIMESTAMP
				WHERE uuid = ?
			`, nodeUUID)
		} else {
			result, err = db.ExecContext(ctx, `
				UPDATE sub_nodes
				SET is_disabled = false, updated_at = CURRENT_TIMESTAMP
				WHERE uuid = ?
			`, nodeUUID)
		}
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

func (r *SubscriptionConnectionRepository) reorderNodes(ctx context.Context, items []reorderNodeItem) error {
	return r.manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		tx, err := db.BeginTx(ctx, nil)
		if err != nil {
			return err
		}
		for _, item := range items {
			if _, err := tx.ExecContext(ctx, `UPDATE sub_nodes SET view_position = ? WHERE uuid = ?`, item.ViewPosition, item.UUID); err != nil {
				_ = tx.Rollback()
				return err
			}
		}
		if _, err := tx.ExecContext(ctx, `SELECT setval('sub_nodes_view_position_seq', (SELECT COALESCE(MAX(view_position), 0) FROM sub_nodes) + 1)`); err != nil {
			_ = tx.Rollback()
			return err
		}
		return tx.Commit()
	})
}
