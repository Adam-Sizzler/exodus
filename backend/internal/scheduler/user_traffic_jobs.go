package scheduler

import (
	"context"

	monitor "exodus/internal/nodes"
	"exodus/internal/srslists"
	"exodus/internal/userwatchdog"
)

// trafficResetDay resets traffic for users with strategy=DAY.
// Runs at 00:05 every day — mirrors remnawave cron "0 5 * * *".
func (s *Scheduler) trafficResetDay(ctx context.Context) error {
	return s.trafficResetByStrategy(ctx, "DAY")
}

// trafficResetWeek resets traffic for users with strategy=WEEK.
// Runs at 00:15 every Monday — mirrors remnawave cron "15 0 * * 1".
func (s *Scheduler) trafficResetWeek(ctx context.Context) error {
	return s.trafficResetByStrategy(ctx, "WEEK")
}

// trafficResetMonthRolling resets traffic for users with strategy=MONTH_ROLLING.
// Runs at 00:10 every day and resets users whose monthly rolling day is due.
func (s *Scheduler) trafficResetMonthRolling(ctx context.Context) error {
	return s.trafficResetByStrategy(ctx, "MONTH_ROLLING")
}

// trafficResetMonth resets traffic for users with strategy=MONTH.
// Runs at 00:20 on 1st of each month — mirrors remnawave cron "20 0 1 * *".
func (s *Scheduler) trafficResetMonth(ctx context.Context) error {
	return s.trafficResetByStrategy(ctx, "MONTH")
}

func (s *Scheduler) trafficResetByStrategy(ctx context.Context, strategy string) error {
	result, err := userwatchdog.ResetTrafficByStrategy(ctx, s.manager, strategy)
	if err != nil {
		return err
	}
	if result.Users == 0 {
		s.cfg.Logger.Debug("Traffic reset: no users to reset", "strategy", strategy)
		return nil
	}
	s.cfg.Logger.Info("Traffic reset completed", "strategy", strategy, "users", result.Users, "node_targets", len(result.NodeUUIDs))
	if len(result.NodeUUIDs) > 0 {
		monitor.RequestNodeDeploy(true, result.NodeUUIDs...)
	}
	return nil
}

// srsListsCheck checks availability of all SRS rule-set URLs.
// Runs every 5 minutes — previously managed by srslists.StartPeriodicChecker.
func (s *Scheduler) srsListsCheck(ctx context.Context) error {
	updated, err := srslists.CheckAndUpdateAvailability(ctx, s.manager, s.cfg)
	if err != nil {
		return err
	}
	if updated > 0 {
		s.cfg.Logger.Info("SRS lists availability check updated records", "updated", updated)
	} else {
		s.cfg.Logger.Debug("SRS lists availability check: no changes")
	}
	return nil
}
