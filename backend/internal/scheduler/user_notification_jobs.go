package scheduler

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	dbmanager "exodus/internal/db/manager"
	"exodus/internal/notifications"
)

type userNotificationRecord struct {
	TID                    int64
	UUID                   string
	Username               string
	ShortUUID              string
	Status                 string
	TrafficLimitBytes      int64
	UsedTrafficBytes       int64
	ExpireAt               time.Time
	LastTriggeredThreshold int
	CreatedAt              time.Time
}

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
		{name: notifications.EventUserExpiresIn72Hours, start: now.Add(72 * time.Hour), end: now.Add(72*time.Hour + time.Minute)},
		{name: notifications.EventUserExpiresIn48Hours, start: now.Add(48 * time.Hour), end: now.Add(48*time.Hour + time.Minute)},
		{name: notifications.EventUserExpiresIn24Hours, start: now.Add(24 * time.Hour), end: now.Add(24*time.Hour + time.Minute)},
		{name: notifications.EventUserExpired24HoursAgo, start: now.Add(-24 * time.Hour), end: now.Add(-24*time.Hour + time.Minute)},
	}

	for _, window := range windows {
		users, err := s.usersByExpireAt(ctx, window.start, window.end)
		if err != nil {
			return err
		}
		if len(users) > 0 {
			s.cfg.Logger.Info("Users found for expiration notification", "event", window.name, "users", len(users))
			for _, user := range users {
				notifications.Emit(ctx, s.cfg, notifications.Event{
					Scope: notifications.ScopeUser,
					Event: window.name,
					Data:  user.notificationData(),
				})
			}
		}
	}
	return nil
}

func (s *Scheduler) usersByExpireAt(ctx context.Context, start, end time.Time) ([]userNotificationRecord, error) {
	users := make([]userNotificationRecord, 0)
	err := s.manager.ExecuteLowPriority(func(db dbmanager.DBExecutor) error {
		rows, err := db.QueryContext(ctx, `
			SELECT
				u.t_id, u.uuid::text, u.username, u.short_uuid, u.status,
				u.traffic_limit_bytes, COALESCE(ut.used_traffic_bytes, 0),
				u.expire_at, u.last_triggered_threshold, u.created_at
			FROM users u
			LEFT JOIN user_traffic ut ON ut.t_id = u.t_id
			WHERE u.expire_at >= ? AND u.expire_at < ?
			ORDER BY u.created_at ASC
		`, start, end)
		if err != nil {
			return err
		}
		defer rows.Close()
		for rows.Next() {
			var user userNotificationRecord
			if scanErr := rows.Scan(
				&user.TID,
				&user.UUID,
				&user.Username,
				&user.ShortUUID,
				&user.Status,
				&user.TrafficLimitBytes,
				&user.UsedTrafficBytes,
				&user.ExpireAt,
				&user.LastTriggeredThreshold,
				&user.CreatedAt,
			); scanErr != nil {
				return scanErr
			}
			users = append(users, user)
		}
		return rows.Err()
	})
	return users, err
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
		users, err := s.triggerThresholdNotifications(ctx, thresholds)
		if err != nil {
			return err
		}
		if len(users) == 0 {
			break
		}
		total += len(users)
		skipTelegram := len(users) >= 500
		if skipTelegram {
			notifications.Emit(ctx, s.cfg, notifications.Event{
				Scope: notifications.ScopeErrors,
				Event: notifications.EventBandwidthMaxNotification,
				Data: map[string]any{
					"batchSize":        len(users),
					"totalProcessed":   total,
					"thresholds":       thresholds,
					"telegramBatchMax": 500,
				},
			})
		}
		for _, user := range users {
			meta := map[string]any{}
			if skipTelegram {
				meta["skipTelegramNotification"] = true
			}
			notifications.Emit(ctx, s.cfg, notifications.Event{
				Scope: notifications.ScopeUser,
				Event: notifications.EventUserBandwidthThreshold,
				Data:  user.notificationData(),
				Meta:  meta,
			})
		}
		if len(users) < 5000 {
			break
		}
	}
	if total > 0 {
		s.cfg.Logger.Info("Users found for bandwidth threshold notification", "users", total, "thresholds", thresholds)
	}
	return nil
}

func (s *Scheduler) triggerThresholdNotifications(ctx context.Context, thresholds []int) ([]userNotificationRecord, error) {
	if len(thresholds) == 0 {
		return nil, nil
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
		RETURNING
			u.t_id,
			u.uuid::text,
			u.username,
			u.short_uuid,
			u.status,
			u.traffic_limit_bytes,
			COALESCE((SELECT ut.used_traffic_bytes FROM user_traffic ut WHERE ut.t_id = u.t_id), 0),
			u.expire_at,
			u.last_triggered_threshold,
			u.created_at
	`, strings.Join(values, ","))

	users := make([]userNotificationRecord, 0)
	err := s.manager.ExecuteLowPriority(func(db dbmanager.DBExecutor) error {
		rows, err := db.QueryContext(ctx, query, args...)
		if err != nil {
			return err
		}
		defer rows.Close()
		for rows.Next() {
			var user userNotificationRecord
			if scanErr := rows.Scan(
				&user.TID,
				&user.UUID,
				&user.Username,
				&user.ShortUUID,
				&user.Status,
				&user.TrafficLimitBytes,
				&user.UsedTrafficBytes,
				&user.ExpireAt,
				&user.LastTriggeredThreshold,
				&user.CreatedAt,
			); scanErr != nil {
				return scanErr
			}
			users = append(users, user)
		}
		return rows.Err()
	})
	return users, err
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

		users, err := s.notConnectedUsers(ctx, start, end)
		if err != nil {
			return err
		}
		if len(users) > 0 {
			s.cfg.Logger.Info("Users found for not-connected notification", "users", len(users), "after_hours", interval)
			skipTelegram := len(users) >= 500
			for _, user := range users {
				meta := map[string]any{"notConnectedAfterHours": interval}
				if skipTelegram {
					meta["skipTelegramNotification"] = true
				}
				notifications.Emit(ctx, s.cfg, notifications.Event{
					Scope: notifications.ScopeUser,
					Event: notifications.EventUserNotConnected,
					Data:  user.notificationData(),
					Meta:  meta,
				})
			}
		}
	}
	return nil
}

func (s *Scheduler) notConnectedUsers(ctx context.Context, start, end time.Time) ([]userNotificationRecord, error) {
	users := make([]userNotificationRecord, 0)
	err := s.manager.ExecuteLowPriority(func(db dbmanager.DBExecutor) error {
		rows, err := db.QueryContext(ctx, `
			SELECT
				u.t_id, u.uuid::text, u.username, u.short_uuid, u.status,
				u.traffic_limit_bytes, COALESCE(ut.used_traffic_bytes, 0),
				u.expire_at, u.last_triggered_threshold, u.created_at
			FROM users u
			INNER JOIN user_traffic ut ON ut.t_id = u.t_id
			WHERE u.status = 'ACTIVE'
			  AND ut.first_connected_at IS NULL
			  AND ut.online_at IS NULL
			  AND u.created_at >= ?
			  AND u.created_at < ?
			ORDER BY u.created_at ASC
		`, start, end)
		if err != nil {
			return err
		}
		defer rows.Close()
		for rows.Next() {
			var user userNotificationRecord
			if scanErr := rows.Scan(
				&user.TID,
				&user.UUID,
				&user.Username,
				&user.ShortUUID,
				&user.Status,
				&user.TrafficLimitBytes,
				&user.UsedTrafficBytes,
				&user.ExpireAt,
				&user.LastTriggeredThreshold,
				&user.CreatedAt,
			); scanErr != nil {
				return scanErr
			}
			users = append(users, user)
		}
		return rows.Err()
	})
	return users, err
}

func (u userNotificationRecord) notificationData() map[string]any {
	return map[string]any{
		"tId":                    u.TID,
		"uuid":                   u.UUID,
		"username":               u.Username,
		"shortUuid":              u.ShortUUID,
		"status":                 u.Status,
		"trafficLimitBytes":      u.TrafficLimitBytes,
		"usedTrafficBytes":       u.UsedTrafficBytes,
		"expireAt":               u.ExpireAt.UTC().Format(time.RFC3339),
		"lastTriggeredThreshold": u.LastTriggeredThreshold,
		"createdAt":              u.CreatedAt.UTC().Format(time.RFC3339),
		"userTraffic": map[string]any{
			"usedTrafficBytes": u.UsedTrafficBytes,
		},
	}
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
