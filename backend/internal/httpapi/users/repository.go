package users

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	dbmanager "exodus/internal/db/manager"
	"exodus/internal/httpapi/shared"
)

type UserRepository struct {
	manager *dbmanager.DatabaseManager
}

func NewUserRepository(manager *dbmanager.DatabaseManager) *UserRepository {
	return &UserRepository{manager: manager}
}

func (r *UserRepository) getAllUserRecords(ctx context.Context) ([]userRecord, error) {
	records := make([]userRecord, 0)
	err := r.manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
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

func (r *UserRepository) getUserRecordByUUID(ctx context.Context, userUUID string) (userRecord, error) {
	var record userRecord
	err := r.manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
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

func (r *UserRepository) getUserRecordsByUUIDs(ctx context.Context, userUUIDs []string) (map[string]userRecord, error) {
	clean := dedupeStrings(userUUIDs)
	records := make(map[string]userRecord, len(clean))
	if len(clean) == 0 {
		return records, nil
	}

	err := r.manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
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

func (r *UserRepository) getUsersActiveInternalSquads(ctx context.Context, userUUIDs []string) (map[string][]internalSquadResponse, error) {
	result := make(map[string][]internalSquadResponse, len(userUUIDs))
	if len(userUUIDs) == 0 {
		return result, nil
	}
	for _, userUUID := range userUUIDs {
		result[userUUID] = []internalSquadResponse{}
	}

	err := r.manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
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

func (r *UserRepository) getAllUserTags(ctx context.Context) ([]string, error) {
	tags := make([]string, 0)
	err := r.manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
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

func (r *UserRepository) replaceUserInternalSquadsTx(ctx context.Context, tx dbmanager.TxExecutor, tID int64, squadUUIDs []string) error {
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

func (r *UserRepository) resolveUserUUIDForUpdate(ctx context.Context, userUUID *string, username *string) (string, error) {
	if userUUID != nil && strings.TrimSpace(*userUUID) != "" {
		return strings.TrimSpace(*userUUID), nil
	}
	if username == nil || strings.TrimSpace(*username) == "" {
		return "", fmt.Errorf("either uuid or username must be provided")
	}

	var resolved string
	err := r.manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		err := db.QueryRowContext(ctx, `SELECT uuid FROM users WHERE username = ?`, strings.TrimSpace(*username)).Scan(&resolved)
		if errors.Is(err, sql.ErrNoRows) {
			return errUserNotFound
		}
		return err
	})
	return resolved, err
}

func (r *UserRepository) resolveUser(ctx context.Context, req resolveUserRequest) (resolveUserResponse, error) {
	var response resolveUserResponse
	clause := ""
	var arg any
	switch {
	case req.UUID != nil:
		clause = "uuid = ?"
		arg = strings.TrimSpace(*req.UUID)
	case req.ID != nil:
		clause = "t_id = ?"
		arg = *req.ID
	case req.ShortUUID != nil:
		clause = "short_uuid = ?"
		arg = strings.TrimSpace(*req.ShortUUID)
	case req.Username != nil:
		clause = "username = ?"
		arg = strings.TrimSpace(*req.Username)
	default:
		return response, fmt.Errorf("missing user lookup field")
	}

	err := r.manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		err := db.QueryRowContext(ctx, fmt.Sprintf(`
			SELECT uuid, t_id, short_uuid, username
			FROM users
			WHERE %s
		`, clause), arg).Scan(&response.UUID, &response.ID, &response.ShortUUID, &response.Username)
		if errors.Is(err, sql.ErrNoRows) {
			return errUserNotFound
		}
		return err
	})
	return response, err
}

func (r *UserRepository) createUser(ctx context.Context, userUUID, shortUUID string, req createUserRequest, credentials userProtocolCredentials, expireAt, createdAt time.Time, lastTrafficResetAt any) (int64, []string, error) {
	var tID int64
	internalSquadNodeUUIDs := make([]string, 0)

	err := r.manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		tx, err := db.BeginTx(ctx, nil)
		if err != nil {
			return err
		}

		insertErr := tx.QueryRowContext(ctx, `
			INSERT INTO users (
					uuid, short_uuid, username, status, traffic_limit_bytes, traffic_limit_strategy,
					expire_at, last_traffic_reset_at, sub_revoked_at,
					trojan_password, vless_uuid, ss_password, naive_password, shadowtls_password, hysteria2_password, anytls_password,
					description, tag, telegram_id, email,
					hwid_device_limit, external_squad_uuid, last_triggered_threshold, created_at, updated_at
				) VALUES (?, ?, ?, ?, ?, ?, ?, ?, NULL, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, 0, ?, ?)
				RETURNING t_id
			`,
			userUUID,
			shortUUID,
			strings.TrimSpace(req.Username),
			normalizeUserStatus(req.Status),
			coalesceInt64(req.TrafficLimitBytes, 0),
			normalizeTrafficStrategy(req.TrafficLimitStrategy),
			expireAt.UTC(),
			lastTrafficResetAt,
			credentials.TrojanPassword,
			credentials.VlessUUID,
			credentials.SSPassword,
			credentials.NaivePassword,
			credentials.ShadowtlsPassword,
			credentials.Hysteria2Password,
			credentials.AnytlsPassword,
			normalizeNullableString(req.Description),
			normalizeUserTag(req.Tag),
			req.TelegramID,
			normalizeNullableString(req.Email),
			req.HwidDeviceLimit,
			normalizeNullableString(req.ExternalSquadUUID),
			createdAt.UTC(),
			createdAt.UTC(),
		).Scan(&tID)
		if insertErr != nil {
			_ = tx.Rollback()
			return mapUserWriteError(insertErr)
		}

		if _, err := tx.ExecContext(ctx, `
			INSERT INTO user_traffic (
				t_id, used_traffic_bytes, lifetime_used_traffic_bytes, online_at,
				last_connected_node_uuid, first_connected_at
			) VALUES (?, 0, 0, NULL, NULL, NULL)
		`, tID); err != nil {
			_ = tx.Rollback()
			return err
		}

		if err := r.replaceUserInternalSquadsTx(ctx, tx, tID, req.ActiveInternalSquads); err != nil {
			_ = tx.Rollback()
			return err
		}
		requestedSquads := dedupeStrings(req.ActiveInternalSquads)
		if len(requestedSquads) > 0 {
			nodeUUIDs, nodeTargetsErr := r.resolveNodeUUIDsForInternalSquadsTx(ctx, tx, requestedSquads)
			if nodeTargetsErr != nil {
				_ = tx.Rollback()
				return nodeTargetsErr
			}
			internalSquadNodeUUIDs = nodeUUIDs
		}

		return tx.Commit()
	})
	return tID, internalSquadNodeUUIDs, err
}

func (r *UserRepository) updateUserRecord(ctx context.Context, targetUUID string, record userRecord, req updateUserRequest, statusToSet string, shouldSetStatus, statusDeployRequired bool) (userRecord, []string, []string, bool, error) {
	var updatedRecord userRecord
	internalSquadsChanged := false
	internalSquadNodeUUIDs := make([]string, 0)
	statusNodeUUIDs := make([]string, 0)

	err := r.manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
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

		if statusDeployRequired {
			nodeUUIDs, nodeTargetsErr := r.resolveNodeUUIDsForUserUUIDsTx(ctx, tx, []string{targetUUID})
			if nodeTargetsErr != nil {
				_ = tx.Rollback()
				return nodeTargetsErr
			}
			statusNodeUUIDs = nodeUUIDs
		}

		if shouldSetStatus {
			add("status", statusToSet)
		}
		if req.TrafficLimitBytes != nil {
			add("traffic_limit_bytes", *req.TrafficLimitBytes)
		}
		if req.TrafficLimitStrategy != nil {
			add("traffic_limit_strategy", strings.ToUpper(strings.TrimSpace(*req.TrafficLimitStrategy)))
		}
		if req.ExpireAt != nil {
			parsed, _ := time.Parse(time.RFC3339, strings.TrimSpace(*req.ExpireAt))
			add("expire_at", parsed.UTC())
		}
		if req.Description.Set {
			if req.Description.Value == nil || strings.TrimSpace(*req.Description.Value) == "" {
				clauses = append(clauses, "description = NULL")
			} else {
				add("description", strings.TrimSpace(*req.Description.Value))
			}
		}
		if req.Tag.Set {
			if req.Tag.Value == nil || strings.TrimSpace(*req.Tag.Value) == "" {
				clauses = append(clauses, "tag = NULL")
			} else {
				add("tag", strings.ToUpper(strings.TrimSpace(*req.Tag.Value)))
			}
		}
		if req.TelegramID.Set {
			if req.TelegramID.Value == nil {
				clauses = append(clauses, "telegram_id = NULL")
			} else {
				add("telegram_id", *req.TelegramID.Value)
			}
		}
		if req.Email.Set {
			if req.Email.Value == nil || strings.TrimSpace(*req.Email.Value) == "" {
				clauses = append(clauses, "email = NULL")
			} else {
				add("email", strings.TrimSpace(*req.Email.Value))
			}
		}
		if req.HwidDeviceLimit.Set {
			if req.HwidDeviceLimit.Value == nil {
				clauses = append(clauses, "hwid_device_limit = NULL")
			} else {
				add("hwid_device_limit", *req.HwidDeviceLimit.Value)
			}
		}

		addOptionalCredential := func(field OptionalString, column string, nullable bool) {
			if !field.Set {
				return
			}
			if field.Value == nil {
				if nullable {
					clauses = append(clauses, fmt.Sprintf("%s = NULL", column))
				}
				return
			}
			add(column, strings.TrimSpace(*field.Value))
		}
		addOptionalCredential(req.TrojanPassword, "trojan_password", false)
		addOptionalCredential(req.VlessUUID, "vless_uuid", false)
		addOptionalCredential(req.SSPassword, "ss_password", false)
		addOptionalCredential(req.NaivePassword, "naive_password", false)
		addOptionalCredential(req.ShadowtlsPassword, "shadowtls_password", false)
		addOptionalCredential(req.Hysteria2Password, "hysteria2_password", false)
		addOptionalCredential(req.AnytlsPassword, "anytls_password", false)

		if req.ExternalSquadUUID.Set {
			if req.ExternalSquadUUID.Value == nil || strings.TrimSpace(*req.ExternalSquadUUID.Value) == "" {
				clauses = append(clauses, "external_squad_uuid = NULL")
			} else {
				add("external_squad_uuid", strings.TrimSpace(*req.ExternalSquadUUID.Value))
			}
		}

		if len(clauses) > 0 {
			args = append(args, targetUUID)
			query := fmt.Sprintf("UPDATE users SET %s, updated_at = CURRENT_TIMESTAMP WHERE uuid = ?", strings.Join(clauses, ", "))
			if _, err := tx.ExecContext(ctx, query, args...); err != nil {
				_ = tx.Rollback()
				return mapUserWriteError(err)
			}
		}

		if req.ActiveInternalSquads != nil {
			currentSquads, loadErr := r.getUserInternalSquadsTx(ctx, tx, record.TID)
			if loadErr != nil {
				_ = tx.Rollback()
				return loadErr
			}
			requestedSquads := dedupeStrings(*req.ActiveInternalSquads)
			if internalSquadSetsDiffer(currentSquads, requestedSquads) {
				affectedSquads := dedupeStrings(append(append([]string{}, currentSquads...), requestedSquads...))
				nodeUUIDs, nodeTargetsErr := r.resolveNodeUUIDsForInternalSquadsTx(ctx, tx, affectedSquads)
				if nodeTargetsErr != nil {
					_ = tx.Rollback()
					return nodeTargetsErr
				}
				if err := r.replaceUserInternalSquadsTx(ctx, tx, record.TID, requestedSquads); err != nil {
					_ = tx.Rollback()
					return err
				}
				internalSquadNodeUUIDs = nodeUUIDs
				internalSquadsChanged = true
			}
		}

		return tx.Commit()
	})
	if err != nil {
		return updatedRecord, nil, nil, false, err
	}

	updatedRecord, err = r.getUserRecordByUUID(ctx, targetUUID)
	return updatedRecord, statusNodeUUIDs, internalSquadNodeUUIDs, internalSquadsChanged, err
}

func (r *UserRepository) deleteUserRecord(ctx context.Context, userUUID string) ([]string, error) {
	internalSquadNodeUUIDs := make([]string, 0)
	err := r.manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		tx, err := db.BeginTx(ctx, nil)
		if err != nil {
			return err
		}

		var tID int64
		if err := tx.QueryRowContext(ctx, `SELECT t_id FROM users WHERE uuid = ?`, userUUID).Scan(&tID); err != nil {
			_ = tx.Rollback()
			if errors.Is(err, sql.ErrNoRows) {
				return errUserNotFound
			}
			return err
		}

		currentSquads, loadErr := r.getUserInternalSquadsTx(ctx, tx, tID)
		if loadErr != nil {
			_ = tx.Rollback()
			return loadErr
		}

		nodeUUIDs, nodeTargetsErr := r.resolveNodeUUIDsForInternalSquadsTx(ctx, tx, currentSquads)
		if nodeTargetsErr != nil {
			_ = tx.Rollback()
			return nodeTargetsErr
		}
		internalSquadNodeUUIDs = nodeUUIDs

		result, err := tx.ExecContext(ctx, `DELETE FROM users WHERE uuid = ?`, userUUID)
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
			return errUserNotFound
		}

		return tx.Commit()
	})
	return internalSquadNodeUUIDs, err
}

func (r *UserRepository) updateUserStatus(ctx context.Context, userUUID string, status string) ([]string, error) {
	nodeUUIDs := make([]string, 0)
	err := r.manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		tx, err := db.BeginTx(ctx, nil)
		if err != nil {
			return err
		}

		targets, nodeTargetsErr := r.resolveNodeUUIDsForUserUUIDsTx(ctx, tx, []string{userUUID})
		if nodeTargetsErr != nil {
			_ = tx.Rollback()
			return nodeTargetsErr
		}
		nodeUUIDs = targets

		result, err := tx.ExecContext(ctx, `
			UPDATE users
			SET status = ?, updated_at = CURRENT_TIMESTAMP
			WHERE uuid = ?
		`, status, userUUID)
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
			return errUserNotFound
		}
		return tx.Commit()
	})
	return nodeUUIDs, err
}

func (r *UserRepository) deleteUsersRecord(ctx context.Context, uuids []string) ([]string, error) {
	internalSquadNodeUUIDs := make([]string, 0)
	err := r.manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		tx, err := db.BeginTx(ctx, nil)
		if err != nil {
			return err
		}

		nodeUUIDs, nodeTargetsErr := r.resolveNodeUUIDsForUserUUIDsTx(ctx, tx, uuids)
		if nodeTargetsErr != nil {
			_ = tx.Rollback()
			return nodeTargetsErr
		}
		internalSquadNodeUUIDs = nodeUUIDs

		if _, err := tx.ExecContext(ctx, `DELETE FROM users WHERE uuid = ANY(?)`, uuids); err != nil {
			_ = tx.Rollback()
			return err
		}

		return tx.Commit()
	})
	return internalSquadNodeUUIDs, err
}

func (r *UserRepository) deleteUsersByStatus(ctx context.Context, status string) (int64, []string, error) {
	var affectedRows int64
	var internalSquadNodeUUIDs []string
	err := r.manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		tx, err := db.BeginTx(ctx, nil)
		if err != nil {
			return err
		}

		rows, queryErr := tx.QueryContext(ctx,
			`SELECT DISTINCT cpitn.node_uuid
			   FROM users u
			   JOIN internal_squad_members ism ON ism.user_id = u.t_id
			   JOIN internal_squad_inbounds isi ON isi.internal_squad_uuid = ism.internal_squad_uuid
			   JOIN config_profile_inbounds_to_nodes cpitn ON cpitn.config_profile_inbound_uuid = isi.inbound_uuid
			  WHERE u.status = ?`, status)
		if queryErr != nil {
			_ = tx.Rollback()
			return queryErr
		}
		nodeUUIDs := make([]string, 0)
		for rows.Next() {
			var nodeUUID string
			if scanErr := rows.Scan(&nodeUUID); scanErr != nil {
				_ = rows.Close()
				_ = tx.Rollback()
				return scanErr
			}
			nodeUUIDs = append(nodeUUIDs, nodeUUID)
		}
		_ = rows.Close()
		if rowsErr := rows.Err(); rowsErr != nil {
			_ = tx.Rollback()
			return rowsErr
		}
		internalSquadNodeUUIDs = nodeUUIDs

		result, execErr := tx.ExecContext(ctx,
			`DELETE FROM users WHERE status = ?`, status)
		if execErr != nil {
			_ = tx.Rollback()
			return execErr
		}
		n, _ := result.RowsAffected()
		affectedRows = n

		return tx.Commit()
	})
	return affectedRows, internalSquadNodeUUIDs, err
}

func (r *UserRepository) bulkUpdateUsers(ctx context.Context, cleanUUIDs []string, clauses []string, args []any) (int64, []string, error) {
	var affectedRows int64
	nodeUUIDs := make([]string, 0)
	err := r.manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		tx, err := db.BeginTx(ctx, nil)
		if err != nil {
			return err
		}

		targets, nodeTargetsErr := r.resolveNodeUUIDsForUserUUIDsTx(ctx, tx, cleanUUIDs)
		if nodeTargetsErr != nil {
			_ = tx.Rollback()
			return nodeTargetsErr
		}
		nodeUUIDs = targets

		queryArgs := append(args, cleanUUIDs)
		query := fmt.Sprintf("UPDATE users SET %s, updated_at = CURRENT_TIMESTAMP WHERE uuid = ANY(?)", strings.Join(clauses, ", "))
		result, execErr := tx.ExecContext(ctx, query, queryArgs...)
		if execErr != nil {
			_ = tx.Rollback()
			return mapUserWriteError(execErr)
		}
		rows, rowsErr := result.RowsAffected()
		if rowsErr != nil {
			_ = tx.Rollback()
			return rowsErr
		}
		affectedRows = rows

		return tx.Commit()
	})
	return affectedRows, nodeUUIDs, err
}

func (r *UserRepository) bulkUpdateUsersSquads(ctx context.Context, cleanUserUUIDs []string, requestedSquads []string) (int64, []string, error) {
	var affectedRows int64
	nodeUUIDs := make([]string, 0)
	err := r.manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		tx, err := db.BeginTx(ctx, nil)
		if err != nil {
			return err
		}

		targets, nodeTargetsErr := r.resolveNodeUUIDsForUserUUIDsTx(ctx, tx, cleanUserUUIDs)
		if nodeTargetsErr != nil {
			_ = tx.Rollback()
			return nodeTargetsErr
		}
		squadTargets, squadTargetsErr := r.resolveNodeUUIDsForInternalSquadsTx(ctx, tx, requestedSquads)
		if squadTargetsErr != nil {
			_ = tx.Rollback()
			return squadTargetsErr
		}
		nodeUUIDs = dedupeStrings(append(targets, squadTargets...))

		rows, err := tx.QueryContext(ctx, `SELECT t_id FROM users WHERE uuid = ANY(?)`, cleanUserUUIDs)
		if err != nil {
			_ = tx.Rollback()
			return err
		}
		userIDs := make([]int64, 0, len(cleanUserUUIDs))
		for rows.Next() {
			var userID int64
			if err := rows.Scan(&userID); err != nil {
				_ = rows.Close()
				_ = tx.Rollback()
				return err
			}
			userIDs = append(userIDs, userID)
		}
		if err := rows.Err(); err != nil {
			_ = rows.Close()
			_ = tx.Rollback()
			return err
		}
		_ = rows.Close()

		for _, userID := range userIDs {
			if err := r.replaceUserInternalSquadsTx(ctx, tx, userID, requestedSquads); err != nil {
				_ = tx.Rollback()
				return err
			}
		}
		if len(userIDs) > 0 {
			if _, err := tx.ExecContext(ctx, `UPDATE users SET updated_at = CURRENT_TIMESTAMP WHERE t_id = ANY(?)`, userIDs); err != nil {
				_ = tx.Rollback()
				return err
			}
		}
		affectedRows = int64(len(userIDs))

		return tx.Commit()
	})
	return affectedRows, nodeUUIDs, err
}

func (r *UserRepository) bulkAllUpdateUsers(ctx context.Context, clauses []string, args []any) (int64, error) {
	var affectedRows int64
	err := r.manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		query := fmt.Sprintf("UPDATE users SET %s, updated_at = CURRENT_TIMESTAMP", strings.Join(clauses, ", "))
		result, execErr := db.ExecContext(ctx, query, args...)
		if execErr != nil {
			return mapUserWriteError(execErr)
		}
		rows, rowsErr := result.RowsAffected()
		if rowsErr != nil {
			return rowsErr
		}
		affectedRows = rows
		return nil
	})
	return affectedRows, err
}

func (r *UserRepository) getUserSubscriptionRequestHistory(ctx context.Context, userUUID string) ([]userSubscriptionRequestHistoryRecord, error) {
	records := make([]userSubscriptionRequestHistoryRecord, 0)
	err := r.manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		var exists bool
		if err := db.QueryRowContext(ctx, `SELECT EXISTS(SELECT 1 FROM users WHERE uuid = ?)`, userUUID).Scan(&exists); err != nil {
			return err
		}
		if !exists {
			return errUserNotFound
		}

		rows, err := db.QueryContext(ctx, `
			SELECT id, user_id, request_ip, user_agent, request_at
			FROM user_subscription_request_history
			WHERE user_id = (SELECT t_id FROM users WHERE uuid = ?)
			ORDER BY request_at DESC
			LIMIT 24
		`, userUUID)
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			var item userSubscriptionRequestHistoryRecord
			var requestAt time.Time
			if scanErr := rows.Scan(&item.ID, &item.UserID, &item.RequestIP, &item.UserAgent, &requestAt); scanErr != nil {
				return scanErr
			}
			item.RequestAt = requestAt.UTC().Format("2006-01-02T15:04:05.000Z")
			records = append(records, item)
		}
		return rows.Err()
	})
	return records, err
}

func (r *UserRepository) getUserAccessibleNodes(ctx context.Context, userUUID string) ([]userAccessibleNode, error) {
	activeNodes := make([]userAccessibleNode, 0)
	err := r.manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		var userID int64
		if err := db.QueryRowContext(ctx, `SELECT t_id FROM users WHERE uuid = ?`, userUUID).Scan(&userID); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return errUserNotFound
			}
			return err
		}

		rows, err := db.QueryContext(ctx, `
			SELECT
				n.uuid,
				n.name,
				n.country_code,
				cp.uuid,
				cp.name,
				sq.uuid,
				sq.name,
				cpi.tag
			FROM nodes n
			INNER JOIN config_profiles cp ON cp.uuid = n.active_config_profile_uuid
			INNER JOIN config_profile_inbounds cpi ON cpi.profile_uuid = cp.uuid
			INNER JOIN config_profile_inbounds_to_nodes cpin
				ON cpin.config_profile_inbound_uuid = cpi.uuid
				AND cpin.node_uuid = n.uuid
			INNER JOIN internal_squad_inbounds isi ON isi.inbound_uuid = cpi.uuid
			INNER JOIN internal_squads sq ON sq.uuid = isi.internal_squad_uuid
			INNER JOIN internal_squad_members ism
				ON ism.internal_squad_uuid = sq.uuid
				AND ism.user_id = ?
			ORDER BY n.view_position ASC, sq.view_position ASC, cpi.tag ASC
		`, userID)
		if err != nil {
			return err
		}
		defer rows.Close()

		nodeIndexes := make(map[string]int)
		squadIndexesByNode := make(map[string]map[string]int)
		for rows.Next() {
			var nodeUUID, nodeName, countryCode, profileUUID, profileName string
			var squadUUID, squadName, inboundTag string
			if scanErr := rows.Scan(&nodeUUID, &nodeName, &countryCode, &profileUUID, &profileName, &squadUUID, &squadName, &inboundTag); scanErr != nil {
				return scanErr
			}

			nodeIndex, ok := nodeIndexes[nodeUUID]
			if !ok {
				activeNodes = append(activeNodes, userAccessibleNode{
					UUID:              nodeUUID,
					NodeName:          nodeName,
					CountryCode:       countryCode,
					ConfigProfileUUID: profileUUID,
					ConfigProfileName: profileName,
					ActiveSquads:      make([]userAccessibleSquad, 0),
				})
				nodeIndex = len(activeNodes) - 1
				nodeIndexes[nodeUUID] = nodeIndex
				squadIndexesByNode[nodeUUID] = make(map[string]int)
			}

			squadIndexes := squadIndexesByNode[nodeUUID]
			squadIndex, ok := squadIndexes[squadUUID]
			if !ok {
				activeNodes[nodeIndex].ActiveSquads = append(activeNodes[nodeIndex].ActiveSquads, userAccessibleSquad{
					SquadName:      squadName,
					ActiveInbounds: make([]string, 0),
				})
				squadIndex = len(activeNodes[nodeIndex].ActiveSquads) - 1
				squadIndexes[squadUUID] = squadIndex
			}

			activeNodes[nodeIndex].ActiveSquads[squadIndex].ActiveInbounds = append(
				activeNodes[nodeIndex].ActiveSquads[squadIndex].ActiveInbounds,
				inboundTag,
			)
		}
		return rows.Err()
	})
	return activeNodes, err
}
