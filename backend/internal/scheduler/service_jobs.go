package scheduler

import (
	"context"
)

func (s *Scheduler) cleanOldUsageRecords(ctx context.Context) error {
	if !s.cfg.Scheduler.ServiceCleanUsageHistory {
		return nil
	}

	if _, err := s.db.ExecContext(ctx, `TRUNCATE TABLE nodes_user_usage_history`); err != nil {
		return err
	}
	if _, err := s.db.ExecContext(ctx, `VACUUM ANALYZE nodes_user_usage_history`); err != nil {
		return err
	}

	s.cfg.Logger.Info("Old usage records cleaned")
	return nil
}

func (s *Scheduler) vacuumTables(ctx context.Context) error {
	if _, err := s.db.ExecContext(ctx, `VACUUM ANALYZE nodes_user_usage_history`); err != nil {
		return err
	}
	s.cfg.Logger.Info("Usage history tables vacuumed")
	return nil
}
