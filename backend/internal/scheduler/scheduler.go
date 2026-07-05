package scheduler

import (
	"context"
	"sync"
	"time"

	"exodus/internal/config"
	dbmanager "exodus/internal/db/manager"
	"exodus/internal/logger"
)

type Scheduler struct {
	manager *dbmanager.DatabaseManager
	cfg     *config.BackendConfig

	mu                  sync.Mutex
	lastRuns            map[string]string
	nodeTrafficNotified map[string]bool
}

func Start(ctx context.Context, wg *sync.WaitGroup, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) {
	if manager == nil || cfg == nil {
		return
	}

	s := &Scheduler{
		manager:             manager,
		cfg:                 cfg,
		lastRuns:            make(map[string]string),
		nodeTrafficNotified: make(map[string]bool),
	}

	cfg.Logger.RoleService(logger.RoleScheduler, logger.ServiceScheduler).Info("Scheduler initialized")
	s.logJobStates()

	wg.Add(1)
	go func() {
		defer wg.Done()
		s.run(ctx)
	}()
}

func (s *Scheduler) logJobStates() {
	if s == nil || s.cfg == nil || s.cfg.Logger == nil {
		return
	}
	jobs := []struct {
		name    string
		enabled bool
	}{
		{name: "cleanOldUsageRecords", enabled: s.cfg.Scheduler.ServiceCleanUsageHistory},
		{name: "expireUserNotifications", enabled: s.cfg.Scheduler.NotificationsEnabled},
		{name: "findUsersForThresholdNotification", enabled: s.cfg.Scheduler.NotificationsEnabled && s.cfg.Scheduler.BandwidthUsageNotificationsEnabled},
		{name: "findNotConnectedUsersNotification", enabled: s.cfg.Scheduler.NotificationsEnabled && s.cfg.Scheduler.NotConnectedUsersNotificationsEnabled},
		{name: "resetNodeTraffic", enabled: true},
		{name: "reviewNodes", enabled: true},
		{name: "vacuumTables", enabled: true},
		{name: "infraBillingNodesNotifications", enabled: true},
		{name: "trafficResetDay (00:05)", enabled: true},
		{name: "trafficResetWeek (Mon 00:15)", enabled: true},
		{name: "trafficResetMonth (1st 00:20)", enabled: true},
		{name: "srsListsCheck (every 5m)", enabled: true},
	}
	log := s.cfg.Logger.RoleService(logger.RoleScheduler, logger.ServiceJobs)
	for _, job := range jobs {
		if job.enabled {
			log.Info("Job enabled", "job", job.name)
		} else {
			log.Info("Job disabled", "job", job.name)
		}
	}
}

func (s *Scheduler) run(ctx context.Context) {
	s.cfg.Logger.RoleService(logger.RoleScheduler, logger.ServiceScheduler).Info("Scheduler started")
	s.runJob(ctx, "resetNodeTraffic", s.resetNodeTraffic)

	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			s.cfg.Logger.RoleService(logger.RoleScheduler, logger.ServiceScheduler).Info("Scheduler stopped")
			return
		case now := <-ticker.C:
			s.tick(ctx, now)
		}
	}
}

func (s *Scheduler) tick(ctx context.Context, now time.Time) {
	local := now.Local()

	// Every hour: review proxy nodes status.
	if local.Minute() == 0 && s.shouldRun("reviewNodes", local.Format("2006-01-02T15")) {
		s.runJob(ctx, "reviewNodes", s.reviewNodes)
	}

	// Every day at 01:00: reset accumulated node traffic counters.
	// Mirrors remnawave EVERY_DAY_AT_1AM (cron: 0 1 * * *).
	if local.Hour() == 1 && local.Minute() == 0 && s.shouldRun("resetNodeTraffic", local.Format("2006-01-02")) {
		s.runJob(ctx, "resetNodeTraffic", s.resetNodeTraffic)
	}

	// Traffic reset schedules — mirrors remnawave cron jobs:
	//   DAY:           0 5 * * *   (every day at 00:05)
	//   MONTH_ROLLING: 10 0 * * *  (every day at 00:10)
	//   WEEK:          15 0 * * 1  (Monday at 00:15)
	//   MONTH:         20 0 1 * *  (1st of month at 00:20)
	if local.Hour() == 0 && local.Minute() == 5 && s.shouldRun("trafficResetDay", local.Format("2006-01-02")) {
		s.runJob(ctx, "trafficResetDay", s.trafficResetDay)
	}
	if local.Hour() == 0 && local.Minute() == 10 && s.shouldRun("trafficResetMonthRolling", local.Format("2006-01-02")) {
		s.runJob(ctx, "trafficResetMonthRolling", s.trafficResetMonthRolling)
	}
	if local.Weekday() == time.Monday && local.Hour() == 0 && local.Minute() == 15 && s.shouldRun("trafficResetWeek", local.Format("2006-01-02")) {
		s.runJob(ctx, "trafficResetWeek", s.trafficResetWeek)
	}
	if local.Day() == 1 && local.Hour() == 0 && local.Minute() == 20 && s.shouldRun("trafficResetMonth", local.Format("2006-01")) {
		s.runJob(ctx, "trafficResetMonth", s.trafficResetMonth)
	}

	// Every minute: expire user notifications.
	if s.shouldRun("expireUserNotifications", local.Format("2006-01-02T15:04")) {
		s.runJob(ctx, "expireUserNotifications", s.findUsersForExpireNotifications)
	}

	// Every 5 minutes: bandwidth threshold notifications.
	if local.Minute()%5 == 0 && s.shouldRun("findUsersForThresholdNotification", local.Format("2006-01-02T15:04")) {
		s.runJob(ctx, "findUsersForThresholdNotification", s.findUsersForThresholdNotification)
	}

	// Every 5 minutes: check SRS lists availability.
	if local.Minute()%5 == 0 && s.shouldRun("srsListsCheck", local.Format("2006-01-02T15:04")) {
		s.runJob(ctx, "srsListsCheck", s.srsListsCheck)
	}

	// Every 10 minutes: not-connected users notifications.
	if local.Minute()%10 == 0 && s.shouldRun("findNotConnectedUsersNotification", local.Format("2006-01-02T15:04")) {
		s.runJob(ctx, "findNotConnectedUsersNotification", s.findNotConnectedUsersNotification)
	}

	// Monday 00:30: clean old usage history records.
	if local.Weekday() == time.Monday && local.Hour() == 0 && local.Minute() == 30 && s.shouldRun("cleanOldUsageRecords", local.Format("2006-01-02")) {
		s.runJob(ctx, "cleanOldUsageRecords", s.cleanOldUsageRecords)
	}

	// Monday 00:45: vacuum tables.
	if local.Weekday() == time.Monday && local.Hour() == 0 && local.Minute() == 45 && s.shouldRun("vacuumTables", local.Format("2006-01-02")) {
		s.runJob(ctx, "vacuumTables", s.vacuumTables)
	}

	// Every day at 17:00: infra billing nodes notifications.
	if local.Hour() == 17 && local.Minute() == 0 && s.shouldRun("infraBillingNodesNotifications", local.Format("2006-01-02")) {
		s.runJob(ctx, "infraBillingNodesNotifications", s.infraBillingNodesNotifications)
	}
}

func (s *Scheduler) shouldRun(name, slot string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	key := name + ":" + slot
	if s.lastRuns[name] == key {
		return false
	}
	s.lastRuns[name] = key
	return true
}

func (s *Scheduler) runJob(ctx context.Context, name string, fn func(context.Context) error) {
	if ctx.Err() != nil {
		return
	}
	start := time.Now()
	if err := fn(ctx); err != nil {
		s.cfg.Logger.RoleService(logger.RoleScheduler, logger.ServiceJobs).Warn("Scheduler job failed", "job", name, "error", err, "duration", time.Since(start).String())
		return
	}
	s.cfg.Logger.RoleService(logger.RoleScheduler, logger.ServiceJobs).Debug("Scheduler job completed", "job", name, "duration", time.Since(start).String())
}

func startOfDay(t time.Time) time.Time {
	local := t.Local()
	return time.Date(local.Year(), local.Month(), local.Day(), 0, 0, 0, 0, local.Location())
}

func endOfDayExclusive(t time.Time) time.Time {
	return startOfDay(t).AddDate(0, 0, 1)
}

func lastDayOfMonth(t time.Time) int {
	return time.Date(t.Year(), t.Month()+1, 0, 0, 0, 0, 0, t.Location()).Day()
}
