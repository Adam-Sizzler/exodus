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
	resetUsersInterval    = time.Minute
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

		cfg.Logger.Info("User status watchdog started", "expired_interval", expiredUsersInterval.String(), "exceeded_interval", exceededUsersInterval.String(), "reset_interval", resetUsersInterval.String())
		runExpiredUsersReview(ctx, manager, cfg)
		runExceededUsersReview(ctx, manager, cfg)
		runScheduledTrafficReset(ctx, manager, cfg, time.Now())

		expiredTicker := time.NewTicker(expiredUsersInterval)
		defer expiredTicker.Stop()
		exceededTicker := time.NewTicker(exceededUsersInterval)
		defer exceededTicker.Stop()
		resetTicker := time.NewTicker(resetUsersInterval)
		defer resetTicker.Stop()

		for {
			select {
			case <-ctx.Done():
				cfg.Logger.Info("User status watchdog stopped")
				return
			case <-expiredTicker.C:
				runExpiredUsersReview(ctx, manager, cfg)
			case <-exceededTicker.C:
				runExceededUsersReview(ctx, manager, cfg)
			case now := <-resetTicker.C:
				runScheduledTrafficReset(ctx, manager, cfg, now)
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

func runScheduledTrafficReset(ctx context.Context, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig, now time.Time) {
	for _, strategy := range scheduledResetStrategies(now) {
		result, err := ResetTrafficByStrategy(ctx, manager, strategy)
		if err != nil {
			cfg.Logger.Warn("Scheduled user traffic reset failed", "strategy", strategy, "error", err)
			continue
		}
		cfg.Logger.Info("Scheduled user traffic reset completed", "strategy", strategy, "users", result.Users, "node_targets", len(result.NodeUUIDs))
		if len(result.NodeUUIDs) > 0 {
			monitor.RequestNodeDeploy(true, result.NodeUUIDs...)
		}
	}
}

func scheduledResetStrategies(now time.Time) []string {
	local := now.Local()
	strategies := make([]string, 0, 3)

	if local.Hour() == 0 && local.Minute() == 3 {
		strategies = append(strategies, "DAY")
	}
	if local.Weekday() == time.Monday && local.Hour() == 0 && local.Minute() == 8 {
		strategies = append(strategies, "WEEK")
	}
	if local.Day() == 1 && local.Hour() == 0 && local.Minute() == 10 {
		strategies = append(strategies, "MONTH")
	}

	return strategies
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
	normalizedStrategy := strings.ToUpper(strings.TrimSpace(strategy))
	result := StatusUpdateResult{NodeUUIDs: []string{}}

	err := manager.ExecuteLowPriority(func(db dbmanager.DBExecutor) error {
		tx, err := db.BeginTx(ctx, nil)
		if err != nil {
			return err
		}

		nodeUUIDs, err := queryLimitedUserNodeUUIDsByStrategyTx(ctx, tx, normalizedStrategy)
		if err != nil {
			_ = tx.Rollback()
			return err
		}

		updateResult, err := tx.ExecContext(ctx, `
			UPDATE users
			SET last_traffic_reset_at = CURRENT_TIMESTAMP,
			    last_triggered_threshold = 0,
			    status = CASE WHEN status = 'LIMITED' THEN 'ACTIVE' ELSE status END,
			    updated_at = CURRENT_TIMESTAMP
			WHERE traffic_limit_strategy = ?
		`, normalizedStrategy)
		if err != nil {
			_ = tx.Rollback()
			return err
		}

		if _, err := tx.ExecContext(ctx, `
			UPDATE user_traffic
			SET used_traffic_bytes = 0
			WHERE t_id IN (SELECT t_id FROM users WHERE traffic_limit_strategy = ?)
		`, normalizedStrategy); err != nil {
			_ = tx.Rollback()
			return err
		}

		affectedRows, _ := updateResult.RowsAffected()
		result.Users = affectedRows
		result.NodeUUIDs = nodeUUIDs
		return tx.Commit()
	})

	return result, err
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

func queryLimitedUserNodeUUIDsByStrategyTx(ctx context.Context, tx dbmanager.TxExecutor, strategy string) ([]string, error) {
	rows, err := tx.QueryContext(ctx, `
		SELECT DISTINCT cpitn.node_uuid::text AS node_uuid
		FROM users u
		JOIN internal_squad_members ism ON ism.user_id = u.t_id
		JOIN internal_squad_inbounds isi ON isi.internal_squad_uuid = ism.internal_squad_uuid
		JOIN config_profile_inbounds_to_nodes cpitn ON cpitn.config_profile_inbound_uuid = isi.inbound_uuid
		WHERE u.status = 'LIMITED'
		  AND u.traffic_limit_strategy = ?
	`, strategy)
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
