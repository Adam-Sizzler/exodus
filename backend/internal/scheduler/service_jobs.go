package scheduler

import (
	"context"

	dbmanager "exodus/internal/db/manager"
)

func (s *Scheduler) cleanOldUsageRecords(ctx context.Context) error {
	if !s.cfg.Scheduler.ServiceCleanUsageHistory {
		return nil
	}

	err := s.manager.ExecuteLowPriority(func(db dbmanager.DBExecutor) error {
		if _, err := db.ExecContext(ctx, `TRUNCATE TABLE nodes_user_usage_history`); err != nil {
			return err
		}
		_, err := db.ExecContext(ctx, `VACUUM ANALYZE nodes_user_usage_history`)
		return err
	})
	if err != nil {
		return err
	}
	s.cfg.Logger.Info("Old usage records cleaned")
	return nil
}

func (s *Scheduler) vacuumTables(ctx context.Context) error {
	err := s.manager.ExecuteLowPriority(func(db dbmanager.DBExecutor) error {
		_, err := db.ExecContext(ctx, `VACUUM ANALYZE nodes_user_usage_history`)
		return err
	})
	if err != nil {
		return err
	}
	s.cfg.Logger.Info("Usage history tables vacuumed")
	return nil
}
