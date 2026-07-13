package configprofiles

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	dbmanager "exodus/internal/db/manager"
	"exodus/internal/httpapi/shared"
)

type ConfigProfileRepository struct {
	manager *dbmanager.DatabaseManager
}

func NewConfigProfileRepository(manager *dbmanager.DatabaseManager) *ConfigProfileRepository {
	return &ConfigProfileRepository{manager: manager}
}

func (r *ConfigProfileRepository) getAllConfigProfileRecords(ctx context.Context) ([]configProfileRecord, error) {
	records := make([]configProfileRecord, 0)
	err := r.manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		rows, err := db.QueryContext(ctx, `
			SELECT uuid, view_position, name, config, created_at, updated_at
			FROM config_profiles
			ORDER BY view_position ASC, name ASC
		`)
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			record, scanErr := r.scanConfigProfileRecord(rows)
			if scanErr != nil {
				return scanErr
			}
			records = append(records, record)
		}
		return rows.Err()
	})
	return records, err
}

func (r *ConfigProfileRepository) getConfigProfileRecordByUUID(ctx context.Context, profileUUID string) (configProfileRecord, error) {
	var record configProfileRecord
	err := r.manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		row := db.QueryRowContext(ctx, `
			SELECT uuid, view_position, name, config, created_at, updated_at
			FROM config_profiles
			WHERE uuid = ?
		`, profileUUID)
		var scanErr error
		record, scanErr = r.scanConfigProfileRecord(row)
		if scanErr == sql.ErrNoRows {
			return errConfigProfileNotFound
		}
		return scanErr
	})
	return record, err
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
	if err := r.manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		rows, err := db.QueryContext(ctx, `
			SELECT isi.inbound_uuid, isi.internal_squad_uuid
			FROM internal_squad_inbounds isi
			JOIN config_profile_inbounds cpi ON cpi.uuid = isi.inbound_uuid
			WHERE cpi.profile_uuid = ANY(?)
		`, profileUUIDs)
		if err != nil {
			return err
		}
		defer rows.Close()
		for rows.Next() {
			var inboundUUID, squadUUID string
			if err := rows.Scan(&inboundUUID, &squadUUID); err != nil {
				return err
			}
			activeSquadsByInbound[inboundUUID] = append(activeSquadsByInbound[inboundUUID], squadUUID)
		}
		return rows.Err()
	}); err != nil {
		return nil, err
	}

	err := r.manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		rows, err := db.QueryContext(ctx, `
			SELECT uuid, profile_uuid, tag, type, network, security, port, raw_inbound
			FROM config_profile_inbounds
			WHERE profile_uuid = ANY(?)
			ORDER BY tag ASC
		`, profileUUIDs)
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			var (
				item     ConfigProfileInbound
				network  sql.NullString
				security sql.NullString
				port     sql.NullInt64
				raw      []byte
			)
			if err := rows.Scan(&item.UUID, &item.ProfileUUID, &item.Tag, &item.Type, &network, &security, &port, &raw); err != nil {
				return err
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
		return rows.Err()
	})

	return result, err
}

func (r *ConfigProfileRepository) getConfigProfileNodesMap(ctx context.Context, profileUUIDs []string) (map[string][]ConfigProfileNode, error) {
	result := make(map[string][]ConfigProfileNode, len(profileUUIDs))
	if len(profileUUIDs) == 0 {
		return result, nil
	}
	for _, profileUUID := range profileUUIDs {
		result[profileUUID] = make([]ConfigProfileNode, 0)
	}

	err := r.manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		rows, err := db.QueryContext(ctx, `
			SELECT active_config_profile_uuid, uuid, name, country_code
			FROM nodes
			WHERE active_config_profile_uuid = ANY(?)
			ORDER BY view_position ASC, name ASC
		`, profileUUIDs)
		if err != nil {
			return err
		}
		defer rows.Close()
		for rows.Next() {
			var profileUUID string
			var node ConfigProfileNode
			if err := rows.Scan(&profileUUID, &node.UUID, &node.Name, &node.CountryCode); err != nil {
				return err
			}
			result[profileUUID] = append(result[profileUUID], node)
		}
		return rows.Err()
	})
	return result, err
}

func (r *ConfigProfileRepository) createConfigProfile(ctx context.Context, profileUUID string, req createConfigProfileRequest) error {
	return r.manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		tx, err := db.BeginTx(ctx, nil)
		if err != nil {
			return err
		}

		if _, err := tx.ExecContext(ctx, `
			INSERT INTO config_profiles (uuid, name, config, created_at, updated_at)
			VALUES (?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		`, profileUUID, strings.TrimSpace(req.Name), req.Config); err != nil {
			_ = tx.Rollback()
			return err
		}

		if _, err := r.syncConfigProfileInboundsTx(ctx, tx, profileUUID, req.Config); err != nil {
			_ = tx.Rollback()
			return err
		}

		return tx.Commit()
	})
}

func (r *ConfigProfileRepository) updateConfigProfile(ctx context.Context, profileUUID string, clauses []string, args []any, updateConfig *json.RawMessage) error {
	return r.manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		tx, err := db.BeginTx(ctx, nil)
		if err != nil {
			return err
		}

		if len(clauses) > 0 {
			args = append(args, profileUUID)
			result, err := tx.ExecContext(ctx, fmt.Sprintf(`
				UPDATE config_profiles
				SET %s, updated_at = CURRENT_TIMESTAMP
				WHERE uuid = ?
			`, strings.Join(clauses, ", ")), args...)
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
				return errConfigProfileNotFound
			}
		}

		if updateConfig != nil {
			if _, err := r.syncConfigProfileInboundsTx(ctx, tx, profileUUID, *updateConfig); err != nil {
				_ = tx.Rollback()
				return err
			}
		}

		return tx.Commit()
	})
}

func (r *ConfigProfileRepository) deleteConfigProfile(ctx context.Context, profileUUID string) error {
	return r.manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		result, err := db.ExecContext(ctx, `DELETE FROM config_profiles WHERE uuid = ?`, profileUUID)
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
	})
}

func (r *ConfigProfileRepository) reorderConfigProfiles(ctx context.Context, items []reorderConfigProfilesItem) error {
	return r.manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		tx, err := db.BeginTx(ctx, nil)
		if err != nil {
			return err
		}
		for _, item := range items {
			if _, err := tx.ExecContext(ctx, `UPDATE config_profiles SET view_position = ? WHERE uuid = ?`, item.ViewPosition, item.UUID); err != nil {
				_ = tx.Rollback()
				return err
			}
		}
		if _, err := tx.ExecContext(ctx, `SELECT setval('config_profiles_view_position_seq', (SELECT COALESCE(MAX(view_position), 0) FROM config_profiles) + 1)`); err != nil {
			_ = tx.Rollback()
			return err
		}
		return tx.Commit()
	})
}

func (r *ConfigProfileRepository) syncConfigProfileInboundsTx(ctx context.Context, db dbmanager.TxExecutor, profileUUID string, configJSON json.RawMessage) (int, error) {
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

		if _, err := db.ExecContext(ctx, `
			INSERT INTO config_profile_inbounds (
				uuid, profile_uuid, tag, type, network, security, port, raw_inbound
			) VALUES (?, ?, ?, ?, ?, ?, ?, ?)
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
		if _, err := db.ExecContext(ctx, `
			DELETE FROM config_profile_inbounds
			WHERE profile_uuid = ? AND NOT (tag = ANY(?))
		`, profileUUID, currentTags); err != nil {
			return 0, err
		}
	} else if _, err := db.ExecContext(ctx, `DELETE FROM config_profile_inbounds WHERE profile_uuid = ?`, profileUUID); err != nil {
		return 0, err
	}

	return len(inbounds), nil
}

func (r *ConfigProfileRepository) getSnippets(ctx context.Context) ([]ConfigProfileSnippet, error) {
	var snippets []ConfigProfileSnippet
	err := r.manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		rows, err := db.QueryContext(ctx, `
			SELECT name, snippet, created_at
			FROM config_profile_snippets
			ORDER BY name ASC`)
		if err != nil {
			return err
		}
		defer rows.Close()

		result := make([]ConfigProfileSnippet, 0)
		for rows.Next() {
			var snippet ConfigProfileSnippet
			if scanErr := rows.Scan(&snippet.Name, &snippet.Snippet, &snippet.CreatedAt); scanErr != nil {
				return scanErr
			}
			result = append(result, snippet)
		}
		if err := rows.Err(); err != nil {
			return err
		}
		snippets = result
		return nil
	})
	return snippets, err
}

func (r *ConfigProfileRepository) createSnippet(ctx context.Context, req configProfileSnippetRequest) error {
	return r.manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		_, err := db.ExecContext(ctx, `
			INSERT INTO config_profile_snippets (name, snippet)
			VALUES (?, ?::jsonb)`, req.Name, string(req.Snippet))
		return err
	})
}

func (r *ConfigProfileRepository) updateSnippet(ctx context.Context, req configProfileSnippetRequest) error {
	return r.manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		res, err := db.ExecContext(ctx, `
			UPDATE config_profile_snippets
			SET snippet = ?::jsonb
			WHERE name = ?`, string(req.Snippet), req.Name)
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
	})
}

func (r *ConfigProfileRepository) deleteSnippet(ctx context.Context, name string) error {
	return r.manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		res, err := db.ExecContext(ctx, `DELETE FROM config_profile_snippets WHERE name = ?`, name)
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
	})
}
