package scheduler

import (
	"context"
	"database/sql"
	"math"
	"strings"
	"time"

	dbmanager "exodus/internal/db/manager"
	"exodus/internal/notifications"
)

type nodeTrafficResetTarget struct {
	UUID        string
	Bytes       int64
	ResetDay    int
	CreatedAt   time.Time
	LastResetAt sql.NullTime
	ScheduledAt time.Time
}

func (s *Scheduler) resetNodeTraffic(ctx context.Context) error {
	now := time.Now()
	targets := make([]nodeTrafficResetTarget, 0)

	err := s.manager.ExecuteLowPriority(func(db dbmanager.DBExecutor) error {
		rows, err := db.QueryContext(ctx, `
			SELECT
				n.uuid::text,
				COALESCE(n.traffic_used_bytes, 0),
				COALESCE(n.traffic_reset_day, 1),
				n.created_at,
				MAX(nth.reset_at)
			FROM nodes n
			LEFT JOIN nodes_traffic_usage_history nth ON nth.node_uuid = n.uuid
			WHERE n.is_traffic_tracking_active = true
			GROUP BY n.uuid, n.traffic_used_bytes, n.traffic_reset_day, n.created_at, n.view_position, n.name
			ORDER BY n.view_position ASC, n.name ASC
		`)
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			var target nodeTrafficResetTarget
			if scanErr := rows.Scan(&target.UUID, &target.Bytes, &target.ResetDay, &target.CreatedAt, &target.LastResetAt); scanErr != nil {
				return scanErr
			}
			scheduledAt, due := nodeTrafficResetDue(now, target.ResetDay, target.CreatedAt, target.LastResetAt)
			if !due {
				continue
			}
			target.ScheduledAt = scheduledAt
			targets = append(targets, target)
		}
		if rowsErr := rows.Err(); rowsErr != nil {
			return rowsErr
		}
		if len(targets) == 0 {
			return nil
		}

		tx, err := db.BeginTx(ctx, nil)
		if err != nil {
			return err
		}
		for _, target := range targets {
			if strings.TrimSpace(target.UUID) == "" {
				continue
			}
			if _, execErr := tx.ExecContext(ctx, `
				INSERT INTO nodes_traffic_usage_history (node_uuid, traffic_bytes, reset_at)
				VALUES (?, ?, CURRENT_TIMESTAMP)
			`, target.UUID, target.Bytes); execErr != nil {
				_ = tx.Rollback()
				return execErr
			}
			if _, execErr := tx.ExecContext(ctx, `
				UPDATE nodes
				SET traffic_used_bytes = 0, updated_at = CURRENT_TIMESTAMP
				WHERE uuid = ?
			`, target.UUID); execErr != nil {
				_ = tx.Rollback()
				return execErr
			}
		}
		return tx.Commit()
	})
	if err != nil {
		return err
	}
	if len(targets) > 0 {
		s.cfg.Logger.Info("Node traffic reset completed", "nodes", len(targets), "scheduled_at", targets[0].ScheduledAt.Format(time.RFC3339))
	}
	return nil
}

func nodeTrafficResetDue(now time.Time, resetDay int, createdAt time.Time, lastResetAt sql.NullTime) (time.Time, bool) {
	scheduledAt := latestNodeTrafficResetBoundary(now.Local(), resetDay)
	if scheduledAt.IsZero() {
		return time.Time{}, false
	}
	if createdAt.After(scheduledAt) {
		return scheduledAt, false
	}
	if lastResetAt.Valid && !lastResetAt.Time.Before(scheduledAt) {
		return scheduledAt, false
	}

	return scheduledAt, true
}

func latestNodeTrafficResetBoundary(now time.Time, resetDay int) time.Time {
	resetDay = clampResetDay(resetDay)
	local := now.Local()
	candidate := nodeTrafficResetBoundary(local.Year(), local.Month(), resetDay, local.Location())
	if local.Before(candidate) {
		previousMonth := local.AddDate(0, -1, 0)
		return nodeTrafficResetBoundary(previousMonth.Year(), previousMonth.Month(), resetDay, local.Location())
	}

	return candidate
}

func nodeTrafficResetBoundary(year int, month time.Month, resetDay int, location *time.Location) time.Time {
	if location == nil {
		location = time.Local
	}

	day := min(clampResetDay(resetDay), lastDayOfMonth(time.Date(year, month, 1, 0, 0, 0, 0, location)))
	return time.Date(year, month, day, 1, 0, 0, 0, location)
}

func clampResetDay(day int) int {
	if day < 1 {
		return 1
	}
	if day > 31 {
		return 31
	}

	return day
}

type nodeReviewRecord struct {
	UUID              string
	Name              string
	Address           string
	Port              int
	TrafficUsedBytes  int64
	TrafficLimitBytes int64
	TrafficResetDay   int
	NotifyPercent     int
}

func (s *Scheduler) reviewNodes(ctx context.Context) error {
	nodes := make([]nodeReviewRecord, 0)
	err := s.manager.ExecuteLowPriority(func(db dbmanager.DBExecutor) error {
		rows, err := db.QueryContext(ctx, `
			SELECT uuid::text, name, address, COALESCE(port, 0), COALESCE(traffic_used_bytes, 0), COALESCE(traffic_limit_bytes, 0), COALESCE(traffic_reset_day, 1), COALESCE(notify_percent, 0)
			FROM nodes
			WHERE is_traffic_tracking_active = true
			  AND is_disabled = false
			  AND COALESCE(notify_percent, 0) > 0
			  AND COALESCE(traffic_limit_bytes, 0) > 0
			ORDER BY view_position ASC, name ASC
		`)
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			var node nodeReviewRecord
			if scanErr := rows.Scan(&node.UUID, &node.Name, &node.Address, &node.Port, &node.TrafficUsedBytes, &node.TrafficLimitBytes, &node.TrafficResetDay, &node.NotifyPercent); scanErr != nil {
				return scanErr
			}
			nodes = append(nodes, node)
		}
		return rows.Err()
	})
	if err != nil {
		return err
	}

	for _, node := range nodes {
		if node.TrafficLimitBytes <= 0 || node.NotifyPercent <= 0 {
			continue
		}
		currentPercent := int(math.Floor((float64(node.TrafficUsedBytes) / float64(node.TrafficLimitBytes)) * 100))
		s.mu.Lock()
		alreadyNotified := s.nodeTrafficNotified[node.UUID]
		if currentPercent >= node.NotifyPercent && !alreadyNotified {
			s.nodeTrafficNotified[node.UUID] = true
			s.mu.Unlock()
			s.cfg.Logger.Warn(
				"Node traffic threshold reached",
				"node_uuid", node.UUID,
				"node", node.Name,
				"percent", currentPercent,
				"threshold", node.NotifyPercent,
			)

			notifications.Emit(ctx, s.cfg, notifications.Event{
				Scope: notifications.ScopeNode,
				Event: notifications.EventNodeTrafficNotify,
				Data: map[string]any{
					"uuid":              node.UUID,
					"name":              node.Name,
					"address":           node.Address,
					"port":              node.Port,
					"trafficUsedBytes":  node.TrafficUsedBytes,
					"trafficLimitBytes": node.TrafficLimitBytes,
					"trafficResetDay":   node.TrafficResetDay,
					"notifyPercent":     node.NotifyPercent,
				},
			})
			continue
		}
		if currentPercent < node.NotifyPercent && alreadyNotified {
			delete(s.nodeTrafficNotified, node.UUID)
		}
		s.mu.Unlock()
	}
	return nil
}

func nullableStringFromSQL(value sql.NullString) string {
	if !value.Valid {
		return ""
	}
	return strings.TrimSpace(value.String)
}
