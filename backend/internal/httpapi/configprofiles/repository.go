package configprofiles

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	exodusdb "exodus/internal/db"
	"exodus/internal/httpapi/shared"
)

type ConfigProfileRepository struct {
	db *sql.DB
}

func NewConfigProfileRepository(db *sql.DB) *ConfigProfileRepository {
	return &ConfigProfileRepository{db: db}
}

func (r *ConfigProfileRepository) getAllConfigProfileRecords(ctx context.Context) ([]configProfileRecord, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT uuid, view_position, name, config, created_at, updated_at
		FROM config_profiles
		ORDER BY view_position ASC, name ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	records := make([]configProfileRecord, 0)
	for rows.Next() {
		record, scanErr := r.scanConfigProfileRecord(rows)
		if scanErr != nil {
			return nil, scanErr
		}
		records = append(records, record)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return records, nil
}

func (r *ConfigProfileRepository) getConfigProfileRecordByUUID(ctx context.Context, profileUUID string) (configProfileRecord, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT uuid, view_position, name, config, created_at, updated_at
		FROM config_profiles
		WHERE uuid = $1
	`, profileUUID)
	record, scanErr := r.scanConfigProfileRecord(row)
	if scanErr == sql.ErrNoRows {
		return record, errConfigProfileNotFound
	}
	return record, scanErr
}

func (r *ConfigProfileRepository) scanConfigProfileRecord(scanner shared.RowScanner) (configProfileRecord, error) {
	var record configProfileRecord
	var viewPosition sql.NullInt64
	var configRaw []byte
	if err := scanner.Scan(&record.UUID, &viewPosition, &record.Name, &configRaw, &record.CreatedAt, &record.UpdatedAt); err != nil {
		return record, err
	}
	if viewPosition.Valid {
		record.ViewPosition = int(viewPosition.Int64)
	}
	record.Config = json.RawMessage(configRaw)
	return record, nil
}

func (r *ConfigProfileRepository) getConfigProfileInboundsMap(ctx context.Context, profileUUIDs []string) (map[string][]ConfigProfileInbound, error) {
	result := make(map[string][]ConfigProfileInbound, len(profileUUIDs))
	if len(profileUUIDs) == 0 {
		return result, nil
	}
	for _, profileUUID := range profileUUIDs {
		result[profileUUID] = make([]ConfigProfileInbound, 0)
	}

	activeSquadsByInbound := make(map[string][]string)
	rows, err := r.db.QueryContext(ctx, `
		SELECT isi.inbound_uuid, isi.internal_squad_uuid
		FROM internal_squad_inbounds isi
		JOIN config_profile_inbounds cpi ON cpi.uuid = isi.inbound_uuid
		WHERE cpi.profile_uuid = ANY($1)
	`, profileUUIDs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var inboundUUID, squadUUID string
		if err := rows.Scan(&inboundUUID, &squadUUID); err != nil {
			return nil, err
		}
		activeSquadsByInbound[inboundUUID] = append(activeSquadsByInbound[inboundUUID], squadUUID)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	inboundRows, err := r.db.QueryContext(ctx, `
		SELECT uuid, profile_uuid, tag, type, network, security, port, raw_inbound
		FROM config_profile_inbounds
		WHERE profile_uuid = ANY($1)
		ORDER BY tag ASC
	`, profileUUIDs)
	if err != nil {
		return nil, err
	}
	defer inboundRows.Close()

	for inboundRows.Next() {
		var (
			item     ConfigProfileInbound
			network  sql.NullString
			security sql.NullString
			port     sql.NullInt64
			raw      []byte
		)
		if err := inboundRows.Scan(&item.UUID, &item.ProfileUUID, &item.Tag, &item.Type, &network, &security, &port, &raw); err != nil {
			return nil, err
		}
		if network.Valid {
			item.Network = &network.String
		}
		if security.Valid {
			item.Security = &security.String
		}
		if port.Valid {
			value := int(port.Int64)
			item.Port = &value
		}
		item.RawInbound = json.RawMessage(raw)
		item.ActiveSquads = dedupeStrings(activeSquadsByInbound[item.UUID])
		result[item.ProfileUUID] = append(result[item.ProfileUUID], item)
	}
	if err := inboundRows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func (r *ConfigProfileRepository) getConfigProfileNodesMap(ctx context.Context, profileUUIDs []string) (map[string][]ConfigProfileNode, error) {
	result := make(map[string][]ConfigProfileNode, len(profileUUIDs))
	if len(profileUUIDs) == 0 {
		return result, nil
	}
	for _, profileUUID := range profileUUIDs {
		result[profileUUID] = make([]ConfigProfileNode, 0)
	}

	rows, err := r.db.QueryContext(ctx, `
		SELECT active_config_profile_uuid, uuid, name, country_code
		FROM nodes
		WHERE active_config_profile_uuid = ANY($1)
		ORDER BY view_position ASC, name ASC
	`, profileUUIDs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var profileUUID string
		var node ConfigProfileNode
		if err := rows.Scan(&profileUUID, &node.UUID, &node.Name, &node.CountryCode); err != nil {
			return nil, err
		}
		result[profileUUID] = append(result[profileUUID], node)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return result, nil
}

func (r *ConfigProfileRepository) createConfigProfile(ctx context.Context, profileUUID string, req createConfigProfileRequest) error {
	return exodusdb.WithRetryTx(ctx, r.db, func(tx *sql.Tx) error {
		if _, err := tx.ExecContext(ctx, `
			INSERT INTO config_profiles (uuid, name, config, created_at, updated_at)
			VALUES ($1, $2, $3, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		`, profileUUID, strings.TrimSpace(req.Name), req.Config); err != nil {
			return err
		}

		if _, err := r.syncConfigProfileInboundsTx(ctx, tx, profileUUID, req.Config); err != nil {
			return err
		}
		return nil
	})
}

func (r *ConfigProfileRepository) updateConfigProfile(ctx context.Context, profileUUID string, clauses []string, args []any, updateConfig *json.RawMessage) error {
	return exodusdb.WithRetryTx(ctx, r.db, func(tx *sql.Tx) error {
		if len(clauses) > 0 {
			txArgs := append(append([]any{}, args...), profileUUID)
			result, err := tx.ExecContext(ctx, fmt.Sprintf(`
				UPDATE config_profiles
				SET %s, updated_at = CURRENT_TIMESTAMP
				WHERE uuid = $%d
			`, strings.Join(clauses, ", "), len(txArgs)), txArgs...)
			if err != nil {
				return err
			}
			rows, err := result.RowsAffected()
			if err != nil {
				return err
			}
			if rows == 0 {
				return errConfigProfileNotFound
			}
		}

		if updateConfig != nil {
			if _, err := r.syncConfigProfileInboundsTx(ctx, tx, profileUUID, *updateConfig); err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *ConfigProfileRepository) deleteConfigProfile(ctx context.Context, profileUUID string) error {
	result, err := r.db.ExecContext(ctx, `DELETE FROM config_profiles WHERE uuid = $1`, profileUUID)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errConfigProfileNotFound
	}
	return nil
}

func (r *ConfigProfileRepository) reorderConfigProfiles(ctx context.Context, items []reorderConfigProfilesItem) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	for _, item := range items {
		if _, err := tx.ExecContext(ctx, `UPDATE config_profiles SET view_position = $1 WHERE uuid = $2`, item.ViewPosition, item.UUID); err != nil {
			return err
		}
	}
	if _, err := tx.ExecContext(ctx, `SELECT setval('config_profiles_view_position_seq', (SELECT COALESCE(MAX(view_position), 0) FROM config_profiles) + 1)`); err != nil {
		return err
	}
	return tx.Commit()
}

func (r *ConfigProfileRepository) syncConfigProfileInboundsTx(ctx context.Context, tx *sql.Tx, profileUUID string, configJSON json.RawMessage) (int, error) {
	inbounds, err := parseConfigInbounds(profileUUID, configJSON)
	if err != nil {
		return 0, err
	}

	currentTags := make([]string, 0, len(inbounds))
	for _, inbound := range inbounds {
		currentTags = append(currentTags, inbound.Tag)
	}

	for _, inbound := range inbounds {
		var networkVal, securityVal, portVal any
		if inbound.Network != nil {
			networkVal = *inbound.Network
		}
		if inbound.Security != nil {
			securityVal = *inbound.Security
		}
		if inbound.Port != nil {
			portVal = *inbound.Port
		}

		if _, err := tx.ExecContext(ctx, `
			INSERT INTO config_profile_inbounds (
				uuid, profile_uuid, tag, type, network, security, port, raw_inbound
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
			ON CONFLICT (tag) DO UPDATE SET
				profile_uuid = EXCLUDED.profile_uuid,
				type         = EXCLUDED.type,
				network      = EXCLUDED.network,
				security     = EXCLUDED.security,
				port         = EXCLUDED.port,
				raw_inbound  = EXCLUDED.raw_inbound
		`, inbound.UUID, inbound.ProfileUUID, inbound.Tag, inbound.Type, networkVal, securityVal, portVal, inbound.RawInbound); err != nil {
			return 0, err
		}
	}

	if len(currentTags) > 0 {
		if _, err := tx.ExecContext(ctx, `
			DELETE FROM config_profile_inbounds
			WHERE profile_uuid = $1 AND NOT (tag = ANY($2))
		`, profileUUID, currentTags); err != nil {
			return 0, err
		}
	} else if _, err := tx.ExecContext(ctx, `DELETE FROM config_profile_inbounds WHERE profile_uuid = $1`, profileUUID); err != nil {
		return 0, err
	}

	return len(inbounds), nil
}

func (r *ConfigProfileRepository) getSnippets(ctx context.Context) ([]ConfigProfileSnippet, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT name, snippet, created_at
		FROM config_profile_snippets
		ORDER BY name ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]ConfigProfileSnippet, 0)
	for rows.Next() {
		var snippet ConfigProfileSnippet
		if scanErr := rows.Scan(&snippet.Name, &snippet.Snippet, &snippet.CreatedAt); scanErr != nil {
			return nil, scanErr
		}
		result = append(result, snippet)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return result, nil
}

func (r *ConfigProfileRepository) createSnippet(ctx context.Context, req configProfileSnippetRequest) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO config_profile_snippets (name, snippet)
		VALUES ($1, $2::jsonb)`, req.Name, string(req.Snippet))
	return err
}

func (r *ConfigProfileRepository) updateSnippet(ctx context.Context, req configProfileSnippetRequest) error {
	res, err := r.db.ExecContext(ctx, `
		UPDATE config_profile_snippets
		SET snippet = $1::jsonb
		WHERE name = $2`, string(req.Snippet), req.Name)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return errConfigProfileSnippetNotFound
	}
	return nil
}

func (r *ConfigProfileRepository) deleteSnippet(ctx context.Context, name string) error {
	res, err := r.db.ExecContext(ctx, `DELETE FROM config_profile_snippets WHERE name = $1`, name)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return errConfigProfileSnippetNotFound
	}
	return nil
}
