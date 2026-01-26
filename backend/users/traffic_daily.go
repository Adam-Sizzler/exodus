package users

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"time"
	"v2ray-stat/backend/config"
	"v2ray-stat/backend/db/manager"
	"v2ray-stat/common"
)

var (
	accumulatedUserTraffic   = make(map[string][2]int64)
	accumulatedSourceTraffic = make(map[string][2]int64)
	accumMutex               sync.Mutex
)

// aggregateHourlyStats агрегирует diff за последние 10 мин и добавляет к текущему часу
func aggregateHourlyStats(manager *manager.DatabaseManager, cfg *config.BackendConfig) error {
	currentTime := time.Now().In(common.TimeLocation)
	//currentTime := common.GetLocalUnix()
	currentHour := currentTime.Format("2006-01-02-15")

	err := manager.ExecuteHighPriority(func(db *sql.DB) error {
		tx, err := db.BeginTx(context.Background(), nil)
		if err != nil {
			return fmt.Errorf("start transaction for aggregation: %w", err)
		}
		defer tx.Rollback()

		accumMutex.Lock()
		for user, acc := range accumulatedUserTraffic {
			if acc[0] == 0 && acc[1] == 0 {
				continue
			}
			_, err := tx.Exec(`
				INSERT INTO daily_user_stats (date, user, uplink, downlink)
				VALUES (?, ?, ?, ?)
				ON CONFLICT(date, user) DO UPDATE SET
					uplink = uplink + ?,
					downlink = downlink + ?`,
				currentHour, user, acc[0], acc[1], acc[0], acc[1])
			if err != nil {
				accumMutex.Unlock()
				return fmt.Errorf("update daily_user_stats for %s: %w", user, err)
			}
		}
		accumulatedUserTraffic = make(map[string][2]int64)

		for source, acc := range accumulatedSourceTraffic {
			if acc[0] == 0 && acc[1] == 0 {
				continue
			}
			_, err := tx.Exec(`
				INSERT INTO daily_source_stats (date, source, uplink, downlink)
				VALUES (?, ?, ?, ?)
				ON CONFLICT(date, source) DO UPDATE SET
					uplink = uplink + ?,
					downlink = downlink + ?`,
				currentHour, source, acc[0], acc[1], acc[0], acc[1])
			if err != nil {
				accumMutex.Unlock()
				return fmt.Errorf("update daily_source_stats for %s: %w", source, err)
			}
		}
		accumulatedSourceTraffic = make(map[string][2]int64)
		accumMutex.Unlock()

		return tx.Commit()
	})
	if err != nil {
		cfg.Logger.Error("Failed to aggregate hourly stats", "error", err)
		return err
	}

	cfg.Logger.Debug("Aggregated hourly stats successfully")
	return nil
}
