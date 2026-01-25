package history

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"v2ray-stat/backend/config"
	"v2ray-stat/backend/db/manager"
)

type StatsRepository struct {
	mgr *manager.DatabaseManager
	cfg *config.BackendConfig
}

func NewStatsRepository(mgr *manager.DatabaseManager, cfg *config.BackendConfig) *StatsRepository {
	return &StatsRepository{mgr: mgr, cfg: cfg}
}

// AddUserTraffic добавляет дельту трафика пользователя за текущий интервал
func (r *StatsRepository) AddUserTraffic(ctx context.Context, user string, uplinkDelta, downlinkDelta int64) error {
	if uplinkDelta == 0 && downlinkDelta == 0 {
		return nil // Нет изменений — skip
	}
	today := time.Now().In(time.FixedZone("CET", 3600)).Format("2006-01-02") // CET для Geneva, CH
	return r.mgr.ExecuteHighPriority(func(db *sql.DB) error {
		query := `
			INSERT INTO daily_user_stats (date, user, uplink, downlink)
			VALUES (?, ?, ?, ?)
			ON CONFLICT(date, user) DO UPDATE SET
				uplink = uplink + excluded.uplink,
				downlink = downlink + excluded.downlink
		`
		_, err := db.ExecContext(ctx, query, today, user, uplinkDelta, downlinkDelta)
		if err != nil {
			return fmt.Errorf("failed to upsert daily user stats: %w", err)
		}
		r.cfg.Logger.Debug("Added user traffic delta", "user", user, "date", today, "uplink_delta", uplinkDelta, "downlink_delta", downlinkDelta)
		return nil
	})
}

// AddSourceTraffic добавляет дельту трафика источника за текущий интервал
func (r *StatsRepository) AddSourceTraffic(ctx context.Context, source string, uplinkDelta, downlinkDelta int64) error {
	if uplinkDelta == 0 && downlinkDelta == 0 {
		return nil
	}
	today := time.Now().In(time.FixedZone("CET", 3600)).Format("2006-01-02")
	return r.mgr.ExecuteHighPriority(func(db *sql.DB) error {
		query := `
			INSERT INTO daily_source_stats (date, source, uplink, downlink)
			VALUES (?, ?, ?, ?)
			ON CONFLICT(date, source) DO UPDATE SET
				uplink = uplink + excluded.uplink,
				downlink = downlink + excluded.downlink
		`
		_, err := db.ExecContext(ctx, query, today, source, uplinkDelta, downlinkDelta)
		if err != nil {
			return fmt.Errorf("failed to upsert daily source stats: %w", err)
		}
		r.cfg.Logger.Debug("Added source traffic delta", "source", source, "date", today, "uplink_delta", uplinkDelta, "downlink_delta", downlinkDelta)
		return nil
	})
}
