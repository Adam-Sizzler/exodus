package users

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	dbmanager "exodus/internal/db/manager"
	"exodus/internal/httpapi/shared"
)

func getAllUserRecords(ctx context.Context, manager *dbmanager.DatabaseManager) ([]userRecord, error) {
	records := make([]userRecord, 0)
	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		rows, err := db.QueryContext(ctx, `
			SELECT
				u.t_id, u.uuid, u.short_uuid, u.username, u.status, u.traffic_limit_bytes,
				u.traffic_limit_strategy, u.expire_at, u.last_traffic_reset_at,
					u.sub_revoked_at, u.trojan_password, u.vless_uuid, u.ss_password,
					u.naive_password, u.shadowtls_password, u.hysteria2_password, u.anytls_password,
					u.description, u.tag, u.telegram_id, u.email, u.hwid_device_limit, u.external_squad_uuid,
				u.last_triggered_threshold, u.created_at, u.updated_at,
				COALESCE(ut.used_traffic_bytes, 0), COALESCE(ut.lifetime_used_traffic_bytes, 0),
				ut.online_at, ut.last_connected_node_uuid, ut.first_connected_at
			FROM users u
			LEFT JOIN user_traffic ut ON ut.t_id = u.t_id
			ORDER BY u.t_id DESC
		`)
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			record, scanErr := scanUserRecord(rows)
			if scanErr != nil {
				return scanErr
			}
			records = append(records, record)
		}
		return rows.Err()
	})
	return records, err
}

func getUserRecordByUUID(ctx context.Context, manager *dbmanager.DatabaseManager, userUUID string) (userRecord, error) {
	var record userRecord
	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		row := db.QueryRowContext(ctx, `
			SELECT
				u.t_id, u.uuid, u.short_uuid, u.username, u.status, u.traffic_limit_bytes,
				u.traffic_limit_strategy, u.expire_at, u.last_traffic_reset_at,
					u.sub_revoked_at, u.trojan_password, u.vless_uuid, u.ss_password,
					u.naive_password, u.shadowtls_password, u.hysteria2_password, u.anytls_password,
					u.description, u.tag, u.telegram_id, u.email, u.hwid_device_limit, u.external_squad_uuid,
				u.last_triggered_threshold, u.created_at, u.updated_at,
				COALESCE(ut.used_traffic_bytes, 0), COALESCE(ut.lifetime_used_traffic_bytes, 0),
				ut.online_at, ut.last_connected_node_uuid, ut.first_connected_at
			FROM users u
			LEFT JOIN user_traffic ut ON ut.t_id = u.t_id
			WHERE u.uuid = ?
		`, userUUID)
		var scanErr error
		record, scanErr = scanUserRecord(row)
		if errors.Is(scanErr, sql.ErrNoRows) {
			return errUserNotFound
		}
		return scanErr
	})
	return record, err
}

func getUserRecordsByUUIDs(ctx context.Context, manager *dbmanager.DatabaseManager, userUUIDs []string) (map[string]userRecord, error) {
	clean := dedupeStrings(userUUIDs)
	records := make(map[string]userRecord, len(clean))
	if len(clean) == 0 {
		return records, nil
	}

	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		rows, err := db.QueryContext(ctx, `
			SELECT
				u.t_id, u.uuid, u.short_uuid, u.username, u.status, u.traffic_limit_bytes,
				u.traffic_limit_strategy, u.expire_at, u.last_traffic_reset_at,
					u.sub_revoked_at, u.trojan_password, u.vless_uuid, u.ss_password,
					u.naive_password, u.shadowtls_password, u.hysteria2_password, u.anytls_password,
					u.description, u.tag, u.telegram_id, u.email, u.hwid_device_limit, u.external_squad_uuid,
				u.last_triggered_threshold, u.created_at, u.updated_at,
				COALESCE(ut.used_traffic_bytes, 0), COALESCE(ut.lifetime_used_traffic_bytes, 0),
				ut.online_at, ut.last_connected_node_uuid, ut.first_connected_at
			FROM users u
			LEFT JOIN user_traffic ut ON ut.t_id = u.t_id
			WHERE u.uuid = ANY(?)
		`, clean)
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			record, scanErr := scanUserRecord(rows)
			if scanErr != nil {
				return scanErr
			}
			records[record.UUID] = record
		}
		return rows.Err()
	})
	return records, err
}

func scanUserRecord(scanner shared.RowScanner) (userRecord, error) {
	var record userRecord
	var (
		lastTrafficReset sql.NullTime
		subRevokedAt     sql.NullTime
		naivePassword    sql.NullString
		shadowtlsPass    sql.NullString
		hysteria2Pass    sql.NullString
		anytlsPass       sql.NullString
		description      sql.NullString
		tag              sql.NullString
		telegramID       sql.NullInt64
		email            sql.NullString
		hwidDeviceLimit  sql.NullInt64
		externalSquad    sql.NullString
		onlineAt         sql.NullTime
		lastNodeUUID     sql.NullString
		firstConnectedAt sql.NullTime
	)

	err := scanner.Scan(
		&record.TID,
		&record.UUID,
		&record.ShortUUID,
		&record.Username,
		&record.Status,
		&record.TrafficLimitBytes,
		&record.TrafficLimitStrategy,
		&record.ExpireAt,
		&lastTrafficReset,
		&subRevokedAt,
		&record.TrojanPassword,
		&record.VlessUUID,
		&record.SSPassword,
		&naivePassword,
		&shadowtlsPass,
		&hysteria2Pass,
		&anytlsPass,
		&description,
		&tag,
		&telegramID,
		&email,
		&hwidDeviceLimit,
		&externalSquad,
		&record.LastTriggeredThreshold,
		&record.CreatedAt,
		&record.UpdatedAt,
		&record.UsedTrafficBytes,
		&record.LifetimeUsedTrafficBytes,
		&onlineAt,
		&lastNodeUUID,
		&firstConnectedAt,
	)
	if err != nil {
		return record, err
	}

	if lastTrafficReset.Valid {
		record.LastTrafficResetAt = &lastTrafficReset.Time
	}
	if subRevokedAt.Valid {
		record.SubRevokedAt = &subRevokedAt.Time
	}
	if naivePassword.Valid {
		record.NaivePassword = &naivePassword.String
	}
	if shadowtlsPass.Valid {
		record.ShadowtlsPassword = &shadowtlsPass.String
	}
	if hysteria2Pass.Valid {
		record.Hysteria2Password = &hysteria2Pass.String
	}
	if anytlsPass.Valid {
		record.AnytlsPassword = &anytlsPass.String
	}
	if description.Valid {
		record.Description = &description.String
	}
	if tag.Valid {
		record.Tag = &tag.String
	}
	if telegramID.Valid {
		value := telegramID.Int64
		record.TelegramID = &value
	}
	if email.Valid {
		record.Email = &email.String
	}
	if hwidDeviceLimit.Valid {
		value := int(hwidDeviceLimit.Int64)
		record.HwidDeviceLimit = &value
	}
	if externalSquad.Valid {
		record.ExternalSquadUUID = &externalSquad.String
	}
	if onlineAt.Valid {
		record.OnlineAt = &onlineAt.Time
	}
	if lastNodeUUID.Valid {
		record.LastConnectedNodeUUID = &lastNodeUUID.String
	}
	if firstConnectedAt.Valid {
		record.FirstConnectedAt = &firstConnectedAt.Time
	}

	return record, nil
}

func getUsersActiveInternalSquads(ctx context.Context, manager *dbmanager.DatabaseManager, userUUIDs []string) (map[string][]internalSquadResponse, error) {
	result := make(map[string][]internalSquadResponse, len(userUUIDs))
	if len(userUUIDs) == 0 {
		return result, nil
	}
	for _, userUUID := range userUUIDs {
		result[userUUID] = []internalSquadResponse{}
	}

	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		rows, err := db.QueryContext(ctx, `
			SELECT u.uuid, s.uuid, s.name
			FROM users u
			INNER JOIN internal_squad_members ism ON ism.user_id = u.t_id
			INNER JOIN internal_squads s ON s.uuid = ism.internal_squad_uuid
			WHERE u.uuid = ANY(?)
			ORDER BY s.view_position ASC, s.name ASC
		`, userUUIDs)
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			var userUUID, squadUUID, squadName string
			if err := rows.Scan(&userUUID, &squadUUID, &squadName); err != nil {
				return err
			}
			result[userUUID] = append(result[userUUID], internalSquadResponse{
				UUID: squadUUID,
				Name: squadName,
			})
		}
		return rows.Err()
	})

	return result, err
}

func getAllUserTags(ctx context.Context, manager *dbmanager.DatabaseManager) ([]string, error) {
	tags := make([]string, 0)
	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		rows, err := db.QueryContext(ctx, `
			SELECT DISTINCT tag
			FROM users
			WHERE tag IS NOT NULL AND tag <> ''
			ORDER BY tag ASC
		`)
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			var tag string
			if err := rows.Scan(&tag); err != nil {
				return err
			}
			tags = append(tags, tag)
		}
		return rows.Err()
	})
	return tags, err
}

func replaceUserInternalSquadsTx(ctx context.Context, tx dbmanager.TxExecutor, tID int64, squadUUIDs []string) error {
	if _, err := tx.ExecContext(ctx, `DELETE FROM internal_squad_members WHERE user_id = ?`, tID); err != nil {
		return err
	}
	for _, squadUUID := range dedupeStrings(squadUUIDs) {
		if _, err := tx.ExecContext(ctx, `
			INSERT INTO internal_squad_members (internal_squad_uuid, user_id)
			VALUES (?, ?)
			ON CONFLICT (internal_squad_uuid, user_id) DO NOTHING
		`, squadUUID, tID); err != nil {
			return err
		}
	}
	return nil
}

func resolveUserUUIDForUpdate(ctx context.Context, manager *dbmanager.DatabaseManager, userUUID *string, username *string) (string, error) {
	if userUUID != nil && strings.TrimSpace(*userUUID) != "" {
		return strings.TrimSpace(*userUUID), nil
	}
	if username == nil || strings.TrimSpace(*username) == "" {
		return "", fmt.Errorf("either uuid or username must be provided")
	}

	var resolved string
	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		err := db.QueryRowContext(ctx, `SELECT uuid FROM users WHERE username = ?`, strings.TrimSpace(*username)).Scan(&resolved)
		if errors.Is(err, sql.ErrNoRows) {
			return errUserNotFound
		}
		return err
	})
	return resolved, err
}
