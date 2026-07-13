package userwatchdog

import (
	"context"
	"database/sql"
	"strings"
	"sync"
	"time"

	"exodus/internal/config"
	dbmanager "exodus/internal/db/manager"
	monitor "exodus/internal/nodes"
)

const (
	expiredUsersInterval  = 30 * time.Second
	exceededUsersInterval = 45 * time.Second
)

type StatusUpdateResult struct {
	Users     int64
	NodeUUIDs []string
}

func Start(ctx context.Context, wg *sync.WaitGroup, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) {
	if manager == nil || cfg == nil {
		return
	}

	wg.Add(1)
	go func() {
		defer wg.Done()

		cfg.Logger.Info("User status watchdog started", "expired_interval", expiredUsersInterval.String(), "exceeded_interval", exceededUsersInterval.String())
		runExpiredUsersReview(ctx, manager, cfg)
		runExceededUsersReview(ctx, manager, cfg)

		expiredTicker := time.NewTicker(expiredUsersInterval)
		defer expiredTicker.Stop()
		exceededTicker := time.NewTicker(exceededUsersInterval)
		defer exceededTicker.Stop()

		for {
			select {
			case <-ctx.Done():
				cfg.Logger.Info("User status watchdog stopped")
				return
			case <-expiredTicker.C:
				runExpiredUsersReview(ctx, manager, cfg)
			case <-exceededTicker.C:
				runExceededUsersReview(ctx, manager, cfg)
			}
		}
	}()
}

func runExpiredUsersReview(ctx context.Context, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) {
	result, err := UpdateExpiredUsers(ctx, manager)
	handleStatusUpdateResult(cfg, "expired", result, err)
}

func runExceededUsersReview(ctx context.Context, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) {
	result, err := UpdateExceededTrafficUsers(ctx, manager)
	handleStatusUpdateResult(cfg, "limited", result, err)
}

func handleStatusUpdateResult(cfg *config.BackendConfig, statusName string, result StatusUpdateResult, err error) {
	if err != nil {
		cfg.Logger.Warn("User status review failed", "status", statusName, "error", err)
		return
	}
	if result.Users == 0 {
		cfg.Logger.Debug("User status review found no users", "status", statusName)
		return
	}

	cfg.Logger.Info("User status review updated users", "status", statusName, "users", result.Users, "node_targets", len(result.NodeUUIDs))
	if len(result.NodeUUIDs) > 0 {
		monitor.RequestNodeDeploy(true, result.NodeUUIDs...)
	}
}

func UpdateExpiredUsers(ctx context.Context, manager *dbmanager.DatabaseManager) (StatusUpdateResult, error) {
	return updateUsersAndCollectNodes(ctx, manager, `
		WITH affected_users AS (
			UPDATE users
			SET status = 'EXPIRED', updated_at = CURRENT_TIMESTAMP
			WHERE status IN ('ACTIVE', 'LIMITED')
			  AND expire_at < CURRENT_TIMESTAMP
			RETURNING t_id
		),
		affected_total AS (
			SELECT COUNT(*)::bigint AS total FROM affected_users
		),
		affected_nodes AS (
			SELECT DISTINCT cpitn.node_uuid::text AS node_uuid
			FROM affected_users au
			JOIN internal_squad_members ism ON ism.user_id = au.t_id
			JOIN internal_squad_inbounds isi ON isi.internal_squad_uuid = ism.internal_squad_uuid
			JOIN config_profile_inbounds_to_nodes cpitn ON cpitn.config_profile_inbound_uuid = isi.inbound_uuid
		)
		SELECT affected_total.total, affected_nodes.node_uuid
		FROM affected_total
		LEFT JOIN affected_nodes ON TRUE
	`)
}

func UpdateExceededTrafficUsers(ctx context.Context, manager *dbmanager.DatabaseManager) (StatusUpdateResult, error) {
	return updateUsersAndCollectNodes(ctx, manager, `
		WITH affected_users AS (
			UPDATE users AS u
			SET status = 'LIMITED', updated_at = CURRENT_TIMESTAMP
			FROM user_traffic AS ut
			WHERE ut.t_id = u.t_id
			  AND u.status = 'ACTIVE'
			  AND u.traffic_limit_bytes <> 0
			  AND ut.used_traffic_bytes >= u.traffic_limit_bytes
			RETURNING u.t_id
		),
		affected_total AS (
			SELECT COUNT(*)::bigint AS total FROM affected_users
		),
		affected_nodes AS (
			SELECT DISTINCT cpitn.node_uuid::text AS node_uuid
			FROM affected_users au
			JOIN internal_squad_members ism ON ism.user_id = au.t_id
			JOIN internal_squad_inbounds isi ON isi.internal_squad_uuid = ism.internal_squad_uuid
			JOIN config_profile_inbounds_to_nodes cpitn ON cpitn.config_profile_inbound_uuid = isi.inbound_uuid
		)
		SELECT affected_total.total, affected_nodes.node_uuid
		FROM affected_total
		LEFT JOIN affected_nodes ON TRUE
	`)
}

func ResetTrafficByStrategy(ctx context.Context, manager *dbmanager.DatabaseManager, strategy string) (StatusUpdateResult, error) {
	return ResetTrafficByStrategyAt(ctx, manager, strategy, time.Now())
}

func ResetTrafficByStrategyAt(ctx context.Context, manager *dbmanager.DatabaseManager, strategy string, now time.Time) (StatusUpdateResult, error) {
	normalizedStrategy := strings.ToUpper(strings.TrimSpace(strategy))
	result := StatusUpdateResult{NodeUUIDs: []string{}}

	boundary, ok := resetPeriodBoundary(normalizedStrategy, now)
	if !ok {
		return result, nil
	}

	err := manager.ExecuteLowPriority(func(db dbmanager.DBExecutor) error {
		tx, err := db.BeginTx(ctx, nil)
		if err != nil {
			return err
		}

		nodeUUIDs, err := queryLimitedUserNodeUUIDsByStrategyTx(ctx, tx, normalizedStrategy, boundary, now)
		if err != nil {
			_ = tx.Rollback()
			return err
		}

		err = tx.QueryRowContext(ctx, `
			WITH affected_users AS (
				UPDATE users
				SET last_traffic_reset_at = CURRENT_TIMESTAMP,
				    last_triggered_threshold = 0,
				    status = CASE WHEN status = 'LIMITED' THEN 'ACTIVE' ELSE status END,
				    updated_at = CURRENT_TIMESTAMP
					WHERE traffic_limit_strategy = ?
					  AND COALESCE(last_traffic_reset_at, created_at) < ?
					  AND (
					      ? <> 'MONTH_ROLLING'
					      OR (
					          (created_at + interval '1 month')::date <= ?::date
					          AND LEAST(
					              EXTRACT(DAY FROM created_at),
					              EXTRACT(DAY FROM date_trunc('month', ?::timestamp) + interval '1 month - 1 day')
					          ) = EXTRACT(DAY FROM ?::timestamp)
					      )
					  )
					RETURNING t_id
				),
			reset_traffic AS (
				INSERT INTO user_traffic (t_id, used_traffic_bytes, lifetime_used_traffic_bytes)
				SELECT t_id, 0, 0 FROM affected_users
				ON CONFLICT (t_id)
				DO UPDATE SET used_traffic_bytes = 0
				RETURNING t_id
			)
			SELECT COUNT(*)::bigint FROM reset_traffic
			`, normalizedStrategy, boundary, normalizedStrategy, now, now, now).Scan(&result.Users)
		if err != nil {
			_ = tx.Rollback()
			return err
		}
		result.NodeUUIDs = nodeUUIDs
		return tx.Commit()
	})

	return result, err
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

func updateUsersAndCollectNodes(ctx context.Context, manager *dbmanager.DatabaseManager, query string) (StatusUpdateResult, error) {
	result := StatusUpdateResult{NodeUUIDs: []string{}}

	err := manager.ExecuteLowPriority(func(db dbmanager.DBExecutor) error {
		rows, err := db.QueryContext(ctx, query)
		if err != nil {
			return err
		}
		defer rows.Close()

		nodeUUIDs := make([]string, 0)
		for rows.Next() {
			var (
				users    int64
				nodeUUID sql.NullString
			)
			if err := rows.Scan(&users, &nodeUUID); err != nil {
				return err
			}
			result.Users = users
			if nodeUUID.Valid {
				nodeUUIDs = append(nodeUUIDs, strings.TrimSpace(nodeUUID.String))
			}
		}
		if err := rows.Err(); err != nil {
			return err
		}
		result.NodeUUIDs = dedupeStrings(nodeUUIDs)
		return nil
	})

	return result, err
}

func queryLimitedUserNodeUUIDsByStrategyTx(ctx context.Context, tx dbmanager.TxExecutor, strategy string, boundary time.Time, now time.Time) ([]string, error) {
	rows, err := tx.QueryContext(ctx, `
			SELECT DISTINCT cpitn.node_uuid::text AS node_uuid
			FROM users u
		JOIN internal_squad_members ism ON ism.user_id = u.t_id
		JOIN internal_squad_inbounds isi ON isi.internal_squad_uuid = ism.internal_squad_uuid
		JOIN config_profile_inbounds_to_nodes cpitn ON cpitn.config_profile_inbound_uuid = isi.inbound_uuid
			WHERE u.status = 'LIMITED'
			  AND u.traffic_limit_strategy = ?
			  AND COALESCE(u.last_traffic_reset_at, u.created_at) < ?
			  AND (
			      ? <> 'MONTH_ROLLING'
			      OR (
			          (u.created_at + interval '1 month')::date <= ?::date
			          AND LEAST(
			              EXTRACT(DAY FROM u.created_at),
			              EXTRACT(DAY FROM date_trunc('month', ?::timestamp) + interval '1 month - 1 day')
			          ) = EXTRACT(DAY FROM ?::timestamp)
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
