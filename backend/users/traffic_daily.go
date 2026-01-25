package users

import (
	"database/sql"
	"time"
	"v2ray-stat/backend/config"
	"v2ray-stat/backend/db/manager"
	"v2ray-stat/common"
)

// aggregateHourlyStats агрегирует diff за последние 10 мин и добавляет к текущему часу
func aggregateHourlyStats(manager *manager.DatabaseManager, cfg *config.BackendConfig) error {
	currentTime := common.GetLocalUnix()
	currentHour := time.Unix(currentTime, 0).Format("2006-01-02 15") // YYYY-MM-DD HH

	// Вычисляем текущие SUM по всем нодам
	currentUserSums := make(map[string][2]int64)
	currentSourceSums := make(map[string][2]int64)

	err := manager.ExecuteHighPriority(func(db *sql.DB) error {
		// Для users
		rows, err := db.Query("SELECT user, SUM(uplink), SUM(downlink) FROM user_traffic GROUP BY user")
		if err != nil {
			return err
		}
		defer rows.Close()
		for rows.Next() {
			var user string
			var uplink, downlink int64
			if err := rows.Scan(&user, &uplink, &downlink); err != nil {
				return err
			}
			currentUserSums[user] = [2]int64{uplink, downlink}
		}

		// Для sources
		rows, err = db.Query("SELECT source, SUM(uplink), SUM(downlink) FROM bound_traffic GROUP BY source")
		if err != nil {
			return err
		}
		defer rows.Close()
		for rows.Next() {
			var source string
			var uplink, downlink int64
			if err := rows.Scan(&source, &uplink, &downlink); err != nil {
				return err
			}
			currentSourceSums[source] = [2]int64{uplink, downlink}
		}
		return nil
	})
	if err != nil {
		return err
	}

	// Вычисляем и добавляем diff
	sumsMutex.Lock()
	defer sumsMutex.Unlock()

	// Для users
	for user, current := range currentUserSums {
		prev := previousUserSums[user]
		uplinkDelta := current[0] - prev[0]
		downlinkDelta := current[1] - prev[1]
		if uplinkDelta > 0 || downlinkDelta > 0 {
			err := manager.ExecuteHighPriority(func(db *sql.DB) error {
				query := `
					INSERT INTO daily_user_stats (date, user, uplink, downlink)
					VALUES (?, ?, ?, ?)
					ON CONFLICT(date, user) DO UPDATE SET
						uplink = uplink + excluded.uplink,
						downlink = downlink + excluded.downlink
				`
				_, err := db.Exec(query, currentHour, user, uplinkDelta, downlinkDelta)
				return err
			})
			if err != nil {
				cfg.Logger.Error("Failed to add user hourly delta", "user", user, "error", err)
			}
		}
		previousUserSums[user] = current
	}

	// Для sources (аналогично)
	for source, current := range currentSourceSums {
		prev := previousSourceSums[source]
		uplinkDelta := current[0] - prev[0]
		downlinkDelta := current[1] - prev[1]
		if uplinkDelta > 0 || downlinkDelta > 0 {
			err := manager.ExecuteHighPriority(func(db *sql.DB) error {
				query := `
					INSERT INTO daily_source_stats (date, source, uplink, downlink)
					VALUES (?, ?, ?, ?)
					ON CONFLICT(date, source) DO UPDATE SET
						uplink = uplink + excluded.uplink,
						downlink = downlink + excluded.downlink
				`
				_, err := db.Exec(query, currentHour, source, uplinkDelta, downlinkDelta)
				return err
			})
			if err != nil {
				cfg.Logger.Error("Failed to add source hourly delta", "source", source, "error", err)
			}
		}
		previousSourceSums[source] = current
	}

	cfg.Logger.Info("Hourly aggregation completed", "hour", currentHour)
	return nil
}
