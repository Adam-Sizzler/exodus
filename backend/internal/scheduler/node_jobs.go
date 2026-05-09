package scheduler

import (
	"context"
	"database/sql"
	"fmt"
	"math"
	"strings"
	"time"

	dbmanager "exodus/internal/db/manager"
)

type nodeTrafficResetTarget struct {
	UUID  string
	Bytes int64
}

func (s *Scheduler) resetNodeTraffic(ctx context.Context) error {
	now := time.Now()
	day := now.Local().Day()
	lastDay := lastDayOfMonth(now.Local())
	targets := make([]nodeTrafficResetTarget, 0)

	err := s.manager.ExecuteLowPriority(func(db dbmanager.DBExecutor) error {
		rows, err := db.QueryContext(ctx, `
			SELECT uuid::text, COALESCE(traffic_used_bytes, 0)
			FROM nodes
			WHERE is_traffic_tracking_active = true
			  AND (
			    COALESCE(traffic_reset_day, 1) = ?
			    OR (COALESCE(traffic_reset_day, 1) > ? AND ? = ?)
			  )
			ORDER BY view_position ASC, name ASC
		`, day, lastDay, day, lastDay)
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			var target nodeTrafficResetTarget
			if scanErr := rows.Scan(&target.UUID, &target.Bytes); scanErr != nil {
				return scanErr
			}
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
		s.cfg.Logger.Info("Node traffic reset completed", "nodes", len(targets))
	}
	return nil
}

type nodeReviewRecord struct {
	UUID              string
	Name              string
	TrafficUsedBytes  int64
	TrafficLimitBytes int64
	NotifyPercent     int
}

func (s *Scheduler) reviewNodes(ctx context.Context) error {
	nodes := make([]nodeReviewRecord, 0)
	err := s.manager.ExecuteLowPriority(func(db dbmanager.DBExecutor) error {
		rows, err := db.QueryContext(ctx, `
			SELECT uuid::text, name, COALESCE(traffic_used_bytes, 0), COALESCE(traffic_limit_bytes, 0), COALESCE(notify_percent, 0)
			FROM nodes
			WHERE is_traffic_tracking_active = true
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
			if scanErr := rows.Scan(&node.UUID, &node.Name, &node.TrafficUsedBytes, &node.TrafficLimitBytes, &node.NotifyPercent); scanErr != nil {
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

func scanCount(rows *sql.Rows) (int, error) {
	count := 0
	for rows.Next() {
		count++
	}
	if err := rows.Err(); err != nil {
		return count, fmt.Errorf("scan rows: %w", err)
	}
	return count, nil
}
