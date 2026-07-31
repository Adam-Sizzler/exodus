package subscriptionconnections

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"

	"exodus/internal/db"
	"exodus/internal/httpapi/shared"
)

type SubscriptionConnectionRepository struct {
	db *sql.DB
}

func NewSubscriptionConnectionRepository(db *sql.DB) *SubscriptionConnectionRepository {
	return &SubscriptionConnectionRepository{db: db}
}

func (r *SubscriptionConnectionRepository) getAllNodeRecords(ctx context.Context) ([]nodeRecord, error) {
	rows, err := r.db.QueryContext(ctx, `
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
		return nil, err
	}
	defer rows.Close()

	var nodes []nodeRecord
	for rows.Next() {
		node, scanErr := r.scanNodeRecord(rows)
		if scanErr != nil {
			return nil, scanErr
		}
		nodes = append(nodes, node)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return nodes, nil
}

func (r *SubscriptionConnectionRepository) getNodeByUUID(ctx context.Context, nodeUUID string) (nodeRecord, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT
			n.uuid, n.id, n.name, n.address, n.public_domain, n.port, n.api_schema, n.api_path, n.grpc_auth_token, sns.subpage_config_uuid,
			n.is_connected, n.is_connecting, n.is_disabled,
			n.last_status_change, n.last_status_message,
			n.provider_uuid, n.view_position, n.tags, n.created_at, n.updated_at
		FROM sub_nodes n
		LEFT JOIN sub_nodes_to_subscription_page_config sns ON sns.node_uuid = n.uuid
		WHERE n.uuid = $1
	`, nodeUUID)
	return r.scanNodeRecord(row)
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
	var tags db.StringArray

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
	rows, err := r.db.QueryContext(ctx, `
		SELECT uuid, name, favicon_link, login_url, created_at, updated_at
		FROM infra_providers
		WHERE uuid = ANY($1)
	`, providerUUIDs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var item providerResponse
		var favicon sql.NullString
		var loginURL sql.NullString
		var createdAt time.Time
		var updatedAt time.Time
		if err := rows.Scan(&item.UUID, &item.Name, &favicon, &loginURL, &createdAt, &updatedAt); err != nil {
			return nil, err
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
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return result, nil
}

func (r *SubscriptionConnectionRepository) getNodeTags(ctx context.Context) ([]string, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT DISTINCT unnest(tags) AS tag FROM sub_nodes ORDER BY tag ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tags := make([]string, 0)
	for rows.Next() {
		var tag string
		if err := rows.Scan(&tag); err != nil {
			return nil, err
		}
		if tag != "" {
			tags = append(tags, tag)
		}
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	sort.Strings(tags)
	return dedupeStrings(tags), nil
}

func (r *SubscriptionConnectionRepository) subpageConfigExists(ctx context.Context, configUUID string) (bool, error) {
	var count int
	if err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM subscription_page_config WHERE uuid = $1`, configUUID).Scan(&count); err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *SubscriptionConnectionRepository) fetchSubpageConfigRaw(ctx context.Context, configUUID string) ([]byte, error) {
	var payload string
	err := r.db.QueryRowContext(ctx, `
		SELECT config
		FROM subscription_page_config
		WHERE uuid = $1
		LIMIT 1
	`, configUUID).Scan(&payload)
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
	err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM sub_nodes WHERE is_disabled = false`).Scan(&enabledCount)
	return enabledCount, err
}

func (r *SubscriptionConnectionRepository) createNode(ctx context.Context, nodeUUID string, req createNodeRequest, schema, grpcAuthToken string, subpageConfigUUID *string, now time.Time) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	_, err = tx.ExecContext(ctx, `
		INSERT INTO sub_nodes (
			uuid, name, address, public_domain, port, api_schema, api_path, grpc_auth_token,
			is_connected, is_connecting, is_disabled, provider_uuid, tags, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
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
		return err
	}
	if subpageConfigUUID != nil {
		if _, err := tx.ExecContext(ctx, `
			INSERT INTO sub_nodes_to_subscription_page_config (node_uuid, subpage_config_uuid)
			VALUES ($1, $2)
			ON CONFLICT (node_uuid) DO UPDATE
			SET subpage_config_uuid = EXCLUDED.subpage_config_uuid
		`, nodeUUID, *subpageConfigUUID); err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (r *SubscriptionConnectionRepository) updateNode(ctx context.Context, nodeUUID string, clauses []string, args []any, subpageConfigSet bool, finalSubpageConfigUUID *string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	if len(clauses) > 0 {
		args = append(args, nodeUUID)
		query := fmt.Sprintf("UPDATE sub_nodes SET %s, updated_at = CURRENT_TIMESTAMP WHERE uuid = $%d", strings.Join(clauses, ", "), len(args))
		result, err := tx.ExecContext(ctx, query, args...)
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
	}
	if subpageConfigSet {
		if finalSubpageConfigUUID == nil {
			if _, err := tx.ExecContext(ctx, `DELETE FROM sub_nodes_to_subscription_page_config WHERE node_uuid = $1`, nodeUUID); err != nil {
				return err
			}
		} else {
			if _, err := tx.ExecContext(ctx, `
				INSERT INTO sub_nodes_to_subscription_page_config (node_uuid, subpage_config_uuid)
				VALUES ($1, $2)
				ON CONFLICT (node_uuid) DO UPDATE
				SET subpage_config_uuid = EXCLUDED.subpage_config_uuid
			`, nodeUUID, *finalSubpageConfigUUID); err != nil {
				return err
			}
		}
	}

	return tx.Commit()
}

func (r *SubscriptionConnectionRepository) deleteNode(ctx context.Context, nodeUUID string) error {
	result, err := r.db.ExecContext(ctx, `DELETE FROM sub_nodes WHERE uuid = $1`, nodeUUID)
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
}

func (r *SubscriptionConnectionRepository) setNodeDisabled(ctx context.Context, nodeUUID string, disabled bool) error {
	var result sql.Result
	var err error
	if disabled {
		result, err = r.db.ExecContext(ctx, `
			UPDATE sub_nodes
			SET is_disabled = true, is_connecting = false, is_connected = false, updated_at = CURRENT_TIMESTAMP
			WHERE uuid = $1
		`, nodeUUID)
	} else {
		result, err = r.db.ExecContext(ctx, `
			UPDATE sub_nodes
			SET is_disabled = false, updated_at = CURRENT_TIMESTAMP
			WHERE uuid = $1
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
}

func (r *SubscriptionConnectionRepository) reorderNodes(ctx context.Context, items []reorderNodeItem) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	for _, item := range items {
		if _, err := tx.ExecContext(ctx, `UPDATE sub_nodes SET view_position = $1 WHERE uuid = $2`, item.ViewPosition, item.UUID); err != nil {
			return err
		}
	}
	if _, err := tx.ExecContext(ctx, `SELECT setval('sub_nodes_view_position_seq', (SELECT COALESCE(MAX(view_position), 0) FROM sub_nodes) + 1)`); err != nil {
		return err
	}
	return tx.Commit()
}
