package scheduler

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	dbmanager "exodus/internal/db/manager"
)

func (s *Scheduler) findUsersForExpireNotifications(ctx context.Context) error {
	if !s.cfg.Scheduler.NotificationsEnabled {
		return nil
	}

	now := time.Now().UTC().Truncate(time.Minute)
	windows := []struct {
		name  string
		start time.Time
		end   time.Time
	}{
		{name: "user.expire_notify_expires_in_72_hours", start: now.Add(72 * time.Hour), end: now.Add(72*time.Hour + time.Minute)},
		{name: "user.expire_notify_expires_in_48_hours", start: now.Add(48 * time.Hour), end: now.Add(48*time.Hour + time.Minute)},
		{name: "user.expire_notify_expires_in_24_hours", start: now.Add(24 * time.Hour), end: now.Add(24*time.Hour + time.Minute)},
		{name: "user.expire_notify_expired_24_hours_ago", start: now.Add(-24 * time.Hour), end: now.Add(-24*time.Hour + time.Minute)},
	}

	for _, window := range windows {
		count, err := s.countUsersByExpireAt(ctx, window.start, window.end)
		if err != nil {
			return err
		}
		if count > 0 {
			s.cfg.Logger.Info("Users found for expiration notification", "event", window.name, "users", count)
		}
	}
	return nil
}

func (s *Scheduler) countUsersByExpireAt(ctx context.Context, start, end time.Time) (int, error) {
	count := 0
	err := s.manager.ExecuteLowPriority(func(db dbmanager.DBExecutor) error {
		return db.QueryRowContext(ctx, `
			SELECT COUNT(*)
			FROM users
			WHERE expire_at >= ? AND expire_at < ?
		`, start, end).Scan(&count)
	})
	return count, err
}

func (s *Scheduler) findUsersForThresholdNotification(ctx context.Context) error {
	if !s.cfg.Scheduler.NotificationsEnabled || !s.cfg.Scheduler.BandwidthUsageNotificationsEnabled {
		return nil
	}

	thresholds := normalizeThresholds(s.cfg.Scheduler.BandwidthUsageNotificationsThreshold)
	if len(thresholds) == 0 {
		s.cfg.Logger.Warn("Bandwidth usage notifications enabled without valid thresholds")
		return nil
	}

	total := 0
	for {
		count, err := s.triggerThresholdNotifications(ctx, thresholds)
		if err != nil {
			return err
		}
		if count == 0 {
			break
		}
		total += count
		if count < 5000 {
			break
		}
	}
	if total > 0 {
		s.cfg.Logger.Info("Users found for bandwidth threshold notification", "users", total, "thresholds", thresholds)
	}
	return nil
}

func (s *Scheduler) triggerThresholdNotifications(ctx context.Context, thresholds []int) (int, error) {
	if len(thresholds) == 0 {
		return 0, nil
	}

	values := make([]string, 0, len(thresholds))
	args := make([]any, 0, len(thresholds))
	for _, threshold := range thresholds {
		values = append(values, "(?)")
		args = append(args, threshold)
	}

	query := fmt.Sprintf(`
		WITH thresholds(pct) AS (VALUES %s),
		candidates AS (
			SELECT
				u.t_id,
				MIN(u.created_at) AS created_at_for_order,
				MAX(t.pct) AS new_threshold
			FROM users u
			INNER JOIN user_traffic ut ON ut.t_id = u.t_id
			INNER JOIN thresholds t
				ON u.status = 'ACTIVE'
				AND u.traffic_limit_bytes > 0
				AND ut.used_traffic_bytes >= (u.traffic_limit_bytes * t.pct / 100)
				AND u.last_triggered_threshold < t.pct
			GROUP BY u.t_id
			ORDER BY created_at_for_order
			LIMIT 5000
		)
		UPDATE users AS u
		SET last_triggered_threshold = c.new_threshold,
		    updated_at = CURRENT_TIMESTAMP
		FROM candidates c
		WHERE u.t_id = c.t_id
		RETURNING u.t_id
	`, strings.Join(values, ","))

	count := 0
	err := s.manager.ExecuteLowPriority(func(db dbmanager.DBExecutor) error {
		rows, err := db.QueryContext(ctx, query, args...)
		if err != nil {
			return err
		}
		defer rows.Close()
		var scanErr error
		count, scanErr = scanCount(rows)
		return scanErr
	})
	return count, err
}

func (s *Scheduler) findNotConnectedUsersNotification(ctx context.Context) error {
	if !s.cfg.Scheduler.NotificationsEnabled || !s.cfg.Scheduler.NotConnectedUsersNotificationsEnabled {
		return nil
	}

	intervals := normalizeNotificationIntervals(s.cfg.Scheduler.NotConnectedUsersNotificationsAfterHours)
	if len(intervals) == 0 {
		s.cfg.Logger.Warn("Not-connected users notifications enabled without valid intervals")
		return nil
	}

	now := time.Now().UTC()
	for _, interval := range intervals {
		target := now.Add(-time.Duration(interval) * time.Hour)
		start := target.Add(-10 * time.Minute)
		end := target

		count, err := s.countNotConnectedUsers(ctx, start, end)
		if err != nil {
			return err
		}
		if count > 0 {
			s.cfg.Logger.Info("Users found for not-connected notification", "users", count, "after_hours", interval)
		}
	}
	return nil
}

func (s *Scheduler) countNotConnectedUsers(ctx context.Context, start, end time.Time) (int, error) {
	count := 0
	err := s.manager.ExecuteLowPriority(func(db dbmanager.DBExecutor) error {
		return db.QueryRowContext(ctx, `
			SELECT COUNT(*)
			FROM users u
			INNER JOIN user_traffic ut ON ut.t_id = u.t_id
			WHERE u.status = 'ACTIVE'
			  AND ut.first_connected_at IS NULL
			  AND ut.online_at IS NULL
			  AND u.created_at >= ?
			  AND u.created_at < ?
		`, start, end).Scan(&count)
	})
	return count, err
}

func normalizeThresholds(input []int) []int {
	if len(input) == 0 {
		return nil
	}
	seen := make(map[int]struct{}, len(input))
	result := make([]int, 0, len(input))
	for _, value := range input {
		if value < 25 || value > 95 {
			continue
		}
		if _, exists := seen[value]; exists {
			continue
		}
		seen[value] = struct{}{}
		result = append(result, value)
	}
	sort.Ints(result)
	return result
}

func normalizeNotificationIntervals(input []int) []int {
	if len(input) == 0 {
		return nil
	}
	seen := make(map[int]struct{}, len(input))
	result := make([]int, 0, len(input))
	for _, value := range input {
		if value <= 0 {
			continue
		}
		if _, exists := seen[value]; exists {
			continue
		}
		seen[value] = struct{}{}
		result = append(result, value)
	}
	sort.Ints(result)
	return result
}
