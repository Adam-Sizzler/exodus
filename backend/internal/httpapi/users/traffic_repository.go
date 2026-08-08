package users

import (
	"context"
	"database/sql"
	"strings"
)

func (r *UserRepository) queryLimitedUserNodeUUIDsTx(ctx context.Context, tx *sql.Tx, userUUIDs []string) ([]string, error) {
	rows, err := tx.QueryContext(ctx, `
		SELECT DISTINCT cpitn.node_uuid
		FROM users u
		JOIN internal_squad_members ism ON ism.user_id = u.id
		JOIN internal_squad_inbounds isi ON isi.internal_squad_uuid = ism.internal_squad_uuid
		JOIN config_profile_inbounds_to_nodes cpitn ON cpitn.config_profile_inbound_uuid = isi.inbound_uuid
		WHERE u.status = 'LIMITED' AND u.uuid = ANY($1)
	`, dedupeStrings(userUUIDs))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanNodeUUIDRows(rows)
}

func (r *UserRepository) queryAllLimitedUserNodeUUIDsTx(ctx context.Context, tx *sql.Tx) ([]string, error) {
	rows, err := tx.QueryContext(ctx, `
		SELECT DISTINCT cpitn.node_uuid
		FROM users u
		JOIN internal_squad_members ism ON ism.user_id = u.id
		JOIN internal_squad_inbounds isi ON isi.internal_squad_uuid = ism.internal_squad_uuid
		JOIN config_profile_inbounds_to_nodes cpitn ON cpitn.config_profile_inbound_uuid = isi.inbound_uuid
		WHERE u.status = 'LIMITED'
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanNodeUUIDRows(rows)
}

func (r *UserRepository) queryReactivatedExpiredUserNodeUUIDsTx(ctx context.Context, tx *sql.Tx, userUUIDs []string, extendDays int) ([]string, error) {
	rows, err := tx.QueryContext(ctx, `
		SELECT DISTINCT cpitn.node_uuid
		FROM users u
		JOIN internal_squad_members ism ON ism.user_id = u.id
		JOIN internal_squad_inbounds isi ON isi.internal_squad_uuid = ism.internal_squad_uuid
		JOIN config_profile_inbounds_to_nodes cpitn ON cpitn.config_profile_inbound_uuid = isi.inbound_uuid
		WHERE u.status = 'EXPIRED'
		  AND u.uuid = ANY($1)
	`, dedupeStrings(userUUIDs))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanNodeUUIDRows(rows)
}

func (r *UserRepository) queryAllReactivatedExpiredUserNodeUUIDsTx(ctx context.Context, tx *sql.Tx, extendDays int) ([]string, error) {
	rows, err := tx.QueryContext(ctx, `
		SELECT DISTINCT cpitn.node_uuid
		FROM users u
		JOIN internal_squad_members ism ON ism.user_id = u.id
		JOIN internal_squad_inbounds isi ON isi.internal_squad_uuid = ism.internal_squad_uuid
		JOIN config_profile_inbounds_to_nodes cpitn ON cpitn.config_profile_inbound_uuid = isi.inbound_uuid
		WHERE u.status = 'EXPIRED'
		  AND u.expire_at + ($1::int * INTERVAL '1 day') > CURRENT_TIMESTAMP
	`, extendDays)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanNodeUUIDRows(rows)
}

func scanNodeUUIDRows(rows *sql.Rows) ([]string, error) {
	nodeUUIDs := make([]string, 0)
	for rows.Next() {
		var nodeUUID string
		if err := rows.Scan(&nodeUUID); err != nil {
			return nil, err
		}
		nodeUUIDs = append(nodeUUIDs, strings.TrimSpace(nodeUUID))
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return dedupeStrings(nodeUUIDs), nil
}

func (r *UserRepository) resetUsersTrafficByUUIDs(ctx context.Context, userUUIDs []string) (int64, []string, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, nil, err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	nodeUUIDs, nodeTargetsErr := r.queryLimitedUserNodeUUIDsTx(ctx, tx, userUUIDs)
	if nodeTargetsErr != nil {
		return 0, nil, nodeTargetsErr
	}

	result, err := tx.ExecContext(ctx, `
		UPDATE users
		SET last_traffic_reset_at = CURRENT_TIMESTAMP,
			last_triggered_threshold = 0,
			status = CASE WHEN status = 'LIMITED' THEN 'ACTIVE' ELSE status END,
			updated_at = CURRENT_TIMESTAMP
		WHERE uuid = ANY($1)
	`, dedupeStrings(userUUIDs))
	if err != nil {
		return 0, nil, err
	}
	affectedRows, _ := result.RowsAffected()

	if _, err := tx.ExecContext(ctx, `
		UPDATE user_traffic
		SET used_traffic_bytes = 0
		WHERE id IN (SELECT id FROM users WHERE uuid = ANY($1))
	`, dedupeStrings(userUUIDs)); err != nil {
		return 0, nil, err
	}

	if err := tx.Commit(); err != nil {
		return 0, nil, err
	}

	return affectedRows, nodeUUIDs, nil
}

func (r *UserRepository) resetAllUsersTraffic(ctx context.Context) (int64, []string, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, nil, err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	nodeUUIDs, nodeTargetsErr := r.queryAllLimitedUserNodeUUIDsTx(ctx, tx)
	if nodeTargetsErr != nil {
		return 0, nil, nodeTargetsErr
	}

	result, err := tx.ExecContext(ctx, `
		UPDATE users
		SET last_traffic_reset_at = CURRENT_TIMESTAMP,
			last_triggered_threshold = 0,
			status = CASE WHEN status = 'LIMITED' THEN 'ACTIVE' ELSE status END,
			updated_at = CURRENT_TIMESTAMP
	`)
	if err != nil {
		return 0, nil, err
	}
	affectedRows, _ := result.RowsAffected()

	if _, err := tx.ExecContext(ctx, `UPDATE user_traffic SET used_traffic_bytes = 0`); err != nil {
		return 0, nil, err
	}

	if err := tx.Commit(); err != nil {
		return 0, nil, err
	}

	return affectedRows, nodeUUIDs, nil
}

func (r *UserRepository) extendUsersExpirationByUUIDs(ctx context.Context, userUUIDs []string, extendDays int) (int64, []string, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, nil, err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	nodeUUIDs, nodeTargetsErr := r.queryReactivatedExpiredUserNodeUUIDsTx(ctx, tx, userUUIDs, extendDays)
	if nodeTargetsErr != nil {
		return 0, nil, nodeTargetsErr
	}

	result, err := tx.ExecContext(ctx, `
		UPDATE users
		SET expire_at = CASE
				WHEN status = 'EXPIRED' THEN CURRENT_TIMESTAMP + ($1::int * INTERVAL '1 day')
				ELSE expire_at + ($1::int * INTERVAL '1 day')
			END,
			status = CASE
				WHEN status = 'EXPIRED' THEN 'ACTIVE'
				ELSE status
			END,
			updated_at = CURRENT_TIMESTAMP
		WHERE uuid = ANY($2)
	`, extendDays, dedupeStrings(userUUIDs))
	if err != nil {
		return 0, nil, err
	}
	affectedRows, _ := result.RowsAffected()

	if err := tx.Commit(); err != nil {
		return 0, nil, err
	}

	return affectedRows, nodeUUIDs, nil
}

func (r *UserRepository) extendAllUsersExpiration(ctx context.Context, extendDays int) (int64, []string, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, nil, err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	nodeUUIDs, nodeTargetsErr := r.queryAllReactivatedExpiredUserNodeUUIDsTx(ctx, tx, extendDays)
	if nodeTargetsErr != nil {
		return 0, nil, nodeTargetsErr
	}

	result, err := tx.ExecContext(ctx, `
		UPDATE users
		SET expire_at = expire_at + ($1::int * INTERVAL '1 day'),
			status = CASE
				WHEN status = 'EXPIRED' AND expire_at + ($2::int * INTERVAL '1 day') > CURRENT_TIMESTAMP THEN 'ACTIVE'
				ELSE status
			END,
			updated_at = CURRENT_TIMESTAMP
	`, extendDays, extendDays)
	if err != nil {
		return 0, nil, err
	}
	affectedRows, _ := result.RowsAffected()

	if err := tx.Commit(); err != nil {
		return 0, nil, err
	}

	return affectedRows, nodeUUIDs, nil
}
