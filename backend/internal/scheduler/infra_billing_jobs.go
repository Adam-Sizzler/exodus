package scheduler

import (
	"context"
	"database/sql"
	"time"

	dbmanager "exodus/internal/db/manager"
)

type infraBillingNotificationWindow struct {
	Name  string
	Start time.Time
	End   time.Time
}

type infraBillingNotificationRecord struct {
	NodeName      string
	ProviderName  string
	LoginURL      string
	NextBillingAt time.Time
}

func (s *Scheduler) infraBillingNodesNotifications(ctx context.Context) error {
	if !s.cfg.Scheduler.NotificationsEnabled {
		return nil
	}

	now := time.Now().Local()
	windows := []infraBillingNotificationWindow{
		{Name: "INFRA_BILLING_NODE_PAYMENT_IN_7_DAYS", Start: startOfDay(now.AddDate(0, 0, 7)), End: endOfDayExclusive(now.AddDate(0, 0, 7))},
		{Name: "INFRA_BILLING_NODE_PAYMENT_IN_48HRS", Start: startOfDay(now.AddDate(0, 0, 2)), End: endOfDayExclusive(now.AddDate(0, 0, 2))},
		{Name: "INFRA_BILLING_NODE_PAYMENT_IN_24HRS", Start: startOfDay(now.AddDate(0, 0, 1)), End: endOfDayExclusive(now.AddDate(0, 0, 1))},
		{Name: "INFRA_BILLING_NODE_PAYMENT_DUE_TODAY", Start: startOfDay(now), End: endOfDayExclusive(now)},
		{Name: "INFRA_BILLING_NODE_PAYMENT_OVERDUE_24HRS", Start: startOfDay(now.AddDate(0, 0, -1)), End: endOfDayExclusive(now.AddDate(0, 0, -1))},
		{Name: "INFRA_BILLING_NODE_PAYMENT_OVERDUE_48HRS", Start: startOfDay(now.AddDate(0, 0, -2)), End: endOfDayExclusive(now.AddDate(0, 0, -2))},
		{Name: "INFRA_BILLING_NODE_PAYMENT_OVERDUE_7_DAYS", Start: startOfDay(now.AddDate(0, 0, -7)), End: endOfDayExclusive(now.AddDate(0, 0, -7))},
	}

	total := 0
	for _, window := range windows {
		items, err := s.getInfraBillingNotifications(ctx, window)
		if err != nil {
			return err
		}
		if len(items) == 0 {
			continue
		}
		total += len(items)
		for _, item := range items {
			s.cfg.Logger.Info(
				"Infra billing notification",
				"event", window.Name,
				"node", item.NodeName,
				"provider", item.ProviderName,
				"login_url", item.LoginURL,
				"next_billing_at", item.NextBillingAt.Format(time.RFC3339),
			)
		}
	}
	if total > 0 {
		s.cfg.Logger.Info("Infra billing notification scan completed", "notifications", total)
	}
	return nil
}

func (s *Scheduler) getInfraBillingNotifications(ctx context.Context, window infraBillingNotificationWindow) ([]infraBillingNotificationRecord, error) {
	items := make([]infraBillingNotificationRecord, 0)
	err := s.manager.ExecuteLowPriority(func(db dbmanager.DBExecutor) error {
		rows, err := db.QueryContext(ctx, `
			SELECT n.name, ip.name, ip.login_url, ibn.next_billing_at
			FROM infra_billing_nodes ibn
			INNER JOIN nodes n ON n.uuid = ibn.node_uuid
			INNER JOIN infra_providers ip ON ip.uuid = ibn.provider_uuid
			WHERE ibn.next_billing_at >= ?
			  AND ibn.next_billing_at < ?
			ORDER BY ibn.next_billing_at ASC, n.name ASC
		`, window.Start, window.End)
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			var (
				item     infraBillingNotificationRecord
				loginURL sql.NullString
			)
			if scanErr := rows.Scan(&item.NodeName, &item.ProviderName, &loginURL, &item.NextBillingAt); scanErr != nil {
				return scanErr
			}
			item.LoginURL = nullableStringFromSQL(loginURL)
			if item.LoginURL == "" {
				item.LoginURL = "https://docs.rw"
			}
			items = append(items, item)
		}
		return rows.Err()
	})
	return items, err
}
