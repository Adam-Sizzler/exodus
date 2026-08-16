package scheduler

import (
	"context"
	"database/sql"
	"strings"
	"sync"
	"time"
)

var (
	nodeDeployMu       sync.RWMutex
	nodeDeployCallback func(restart bool, nodeUUIDs ...string)
)

// SetNodeDeployCallback registers a callback to trigger node deploy on user state changes.
func SetNodeDeployCallback(cb func(restart bool, nodeUUIDs ...string)) {
	nodeDeployMu.Lock()
	defer nodeDeployMu.Unlock()
	nodeDeployCallback = cb
}

func triggerNodeDeploy(restart bool, nodeUUIDs ...string) {
	nodeDeployMu.RLock()
	cb := nodeDeployCallback
	nodeDeployMu.RUnlock()
	if cb != nil && len(nodeUUIDs) > 0 {
		cb(restart, nodeUUIDs...)
	}
}

type StatusUpdateResult struct {
	Users     int64
	NodeUUIDs []string
}

func (s *Scheduler) runExpiredUsersReview(ctx context.Context) error {
	result, err := UpdateExpiredUsers(ctx, s.db)
	s.handleStatusUpdateResult("expired", result, err)
	return err
}

func (s *Scheduler) runExceededUsersReview(ctx context.Context) error {
	result, err := UpdateExceededTrafficUsers(ctx, s.db)
	s.handleStatusUpdateResult("limited", result, err)
	return err
}

func (s *Scheduler) handleStatusUpdateResult(statusName string, result StatusUpdateResult, err error) {
	if err != nil {
		s.cfg.Logger.Warn("User status review failed", "status", statusName, "error", err)
		return
	}
	if result.Users == 0 {
		s.cfg.Logger.Debug("User status review found no users", "status", statusName)
		return
	}

	s.cfg.Logger.Info("User status review updated users", "status", statusName, "users", result.Users, "node_targets", len(result.NodeUUIDs))
	if len(result.NodeUUIDs) > 0 {
		triggerNodeDeploy(true, result.NodeUUIDs...)
	}
}

func UpdateExpiredUsers(ctx context.Context, db *sql.DB) (StatusUpdateResult, error) {
	return updateUsersAndCollectNodes(ctx, db, `
		WITH affected_users AS (
			UPDATE users
			SET status = 'EXPIRED', updated_at = CURRENT_TIMESTAMP
			WHERE status IN ('ACTIVE', 'LIMITED')
			  AND expire_at < CURRENT_TIMESTAMP
			RETURNING id
		),
		affected_total AS (
			SELECT COUNT(*)::bigint AS total FROM affected_users
		),
		affected_nodes AS (
			SELECT DISTINCT cpitn.node_uuid::text AS node_uuid
			FROM affected_users au
			JOIN internal_squad_members ism ON ism.user_id = au.id
			JOIN internal_squad_inbounds isi ON isi.internal_squad_uuid = ism.internal_squad_uuid
			JOIN config_profile_inbounds_to_nodes cpitn ON cpitn.config_profile_inbound_uuid = isi.inbound_uuid
		)
		SELECT affected_total.total, affected_nodes.node_uuid
		FROM affected_total
		LEFT JOIN affected_nodes ON TRUE
	`)
}

func UpdateExceededTrafficUsers(ctx context.Context, db *sql.DB) (StatusUpdateResult, error) {
	return updateUsersAndCollectNodes(ctx, db, `
		WITH affected_users AS (
			UPDATE users AS u
			SET status = 'LIMITED', updated_at = CURRENT_TIMESTAMP
			FROM user_traffic AS ut
			WHERE ut.id = u.id
			  AND u.status = 'ACTIVE'
			  AND u.traffic_limit_bytes <> 0
			  AND ut.used_traffic_bytes >= u.traffic_limit_bytes
			RETURNING u.id
		),
		affected_total AS (
			SELECT COUNT(*)::bigint AS total FROM affected_users
		),
		affected_nodes AS (
			SELECT DISTINCT cpitn.node_uuid::text AS node_uuid
			FROM affected_users au
			JOIN internal_squad_members ism ON ism.user_id = au.id
			JOIN internal_squad_inbounds isi ON isi.internal_squad_uuid = ism.internal_squad_uuid
			JOIN config_profile_inbounds_to_nodes cpitn ON cpitn.config_profile_inbound_uuid = isi.inbound_uuid
		)
		SELECT affected_total.total, affected_nodes.node_uuid
		FROM affected_total
		LEFT JOIN affected_nodes ON TRUE
	`)
}

func ResetTrafficByStrategy(ctx context.Context, db *sql.DB, strategy string) (StatusUpdateResult, error) {
	return ResetTrafficByStrategyAt(ctx, db, strategy, time.Now())
}

func ResetTrafficByStrategyAt(ctx context.Context, db *sql.DB, strategy string, now time.Time) (StatusUpdateResult, error) {
	normalizedStrategy := strings.ToUpper(strings.TrimSpace(strategy))
	result := StatusUpdateResult{NodeUUIDs: []string{}}

	boundary, ok := resetPeriodBoundary(normalizedStrategy, now)
	if !ok {
		return result, nil
	}

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return result, err
	}
	defer func() { _ = tx.Rollback() }()

	nodeUUIDs, err := queryLimitedUserNodeUUIDsByStrategyTx(ctx, tx, normalizedStrategy, boundary, now)
	if err != nil {
		return result, err
	}

	err = tx.QueryRowContext(ctx, `
		WITH affected_users AS (
			UPDATE users
			SET last_traffic_reset_at = CURRENT_TIMESTAMP,
			    last_triggered_threshold = 0,
			    status = CASE WHEN status = 'LIMITED' THEN 'ACTIVE' ELSE status END,
			    updated_at = CURRENT_TIMESTAMP
				WHERE traffic_limit_strategy = $1
				  AND COALESCE(last_traffic_reset_at, created_at) < $2
				  AND (
				      $3 <> 'MONTH_ROLLING'
				      OR (
				          (created_at + interval '1 month')::date <= $4::date
				          AND LEAST(
				              EXTRACT(DAY FROM created_at),
				              EXTRACT(DAY FROM date_trunc('month', $5::timestamp) + interval '1 month - 1 day')
				          ) = EXTRACT(DAY FROM $6::timestamp)
				      )
				  )
				RETURNING id
			),
		reset_traffic AS (
			INSERT INTO user_traffic (id, used_traffic_bytes, lifetime_used_traffic_bytes)
			SELECT id, 0, 0 FROM affected_users
			ON CONFLICT (id)
			DO UPDATE SET used_traffic_bytes = 0
			RETURNING id
		)
		SELECT COUNT(*)::bigint FROM reset_traffic
	`, normalizedStrategy, boundary, normalizedStrategy, now, now, now).Scan(&result.Users)
	if err != nil {
		return result, err
	}

	result.NodeUUIDs = nodeUUIDs
	if err := tx.Commit(); err != nil {
		return result, err
	}
	return result, nil
}

func resetPeriodBoundary(strategy string, now time.Time) (time.Time, bool) {
	local := now.Local()
	dayStart := time.Date(local.Year(), local.Month(), local.Day(), 0, 0, 0, 0, local.Location())

	switch strings.ToUpper(strings.TrimSpace(strategy)) {
	case "DAY":
		return dayStart, true
	case "WEEK":
		daysSinceMonday := (int(local.Weekday()) - int(time.Monday) + 7) % 7
		return dayStart.AddDate(0, 0, -daysSinceMonday), true
	case "MONTH":
		return time.Date(local.Year(), local.Month(), 1, 0, 0, 0, 0, local.Location()), true
	case "MONTH_ROLLING":
		return dayStart, true
	default:
		return time.Time{}, false
	}
}

func updateUsersAndCollectNodes(ctx context.Context, db *sql.DB, query string) (StatusUpdateResult, error) {
	result := StatusUpdateResult{NodeUUIDs: []string{}}

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return result, err
	}
	defer rows.Close()

	nodeUUIDs := make([]string, 0)
	for rows.Next() {
		var (
			users    int64
			nodeUUID sql.NullString
		)
		if err := rows.Scan(&users, &nodeUUID); err != nil {
			return result, err
		}
		result.Users = users
		if nodeUUID.Valid {
			nodeUUIDs = append(nodeUUIDs, strings.TrimSpace(nodeUUID.String))
		}
	}
	if err := rows.Err(); err != nil {
		return result, err
	}
	result.NodeUUIDs = dedupeStrings(nodeUUIDs)
	return result, nil
}

func queryLimitedUserNodeUUIDsByStrategyTx(ctx context.Context, tx *sql.Tx, strategy string, boundary time.Time, now time.Time) ([]string, error) {
	rows, err := tx.QueryContext(ctx, `
		SELECT DISTINCT cpitn.node_uuid::text AS node_uuid
		FROM users u
		JOIN internal_squad_members ism ON ism.user_id = u.id
		JOIN internal_squad_inbounds isi ON isi.internal_squad_uuid = ism.internal_squad_uuid
		JOIN config_profile_inbounds_to_nodes cpitn ON cpitn.config_profile_inbound_uuid = isi.inbound_uuid
		WHERE u.status = 'LIMITED'
		  AND u.traffic_limit_strategy = $1
		  AND COALESCE(u.last_traffic_reset_at, u.created_at) < $2
		  AND (
		      $3 <> 'MONTH_ROLLING'
		      OR (
		          (u.created_at + interval '1 month')::date <= $4::date
		          AND LEAST(
		              EXTRACT(DAY FROM u.created_at),
		              EXTRACT(DAY FROM date_trunc('month', $5::timestamp) + interval '1 month - 1 day')
		          ) = EXTRACT(DAY FROM $6::timestamp)
		      )
		  )
	`, strategy, boundary, strategy, now, now, now)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	nodeUUIDs := make([]string, 0)
	for rows.Next() {
		var nodeUUID string
		if err := rows.Scan(&nodeUUID); err != nil {
			return nil, err
		}
		nodeUUIDs = append(nodeUUIDs, nodeUUID)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return dedupeStrings(nodeUUIDs), nil
}

func dedupeStrings(values []string) []string {
	result := make([]string, 0, len(values))
	seen := make(map[string]struct{}, len(values))
	for _, value := range values {
		trimmed := strings.TrimSpace(value)
		if trimmed == "" {
			continue
		}
		if _, ok := seen[trimmed]; ok {
			continue
		}
		seen[trimmed] = struct{}{}
		result = append(result, trimmed)
	}
	return result
}
