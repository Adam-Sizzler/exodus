package subscriptionrequesthistory

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"cerberus/backend/config"
	dbmanager "cerberus/backend/db/manager"
	"cerberus/backend/httpapi/shared"
)

type historyRecord struct {
	ID        int64   `json:"id"`
	UserUUID  string  `json:"userUuid"`
	RequestIP *string `json:"requestIp"`
	UserAgent *string `json:"userAgent"`
	RequestAt string  `json:"requestAt"`
}

func SubscriptionRequestHistoryHandler(manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			shared.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}

		start := parseInt(r.URL.Query().Get("start"), 0, 0, 1_000_000)
		size := parseInt(r.URL.Query().Get("size"), 25, 1, 500)

		records := make([]historyRecord, 0)
		total := 0
		err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
			if err := db.QueryRowContext(r.Context(), `SELECT COUNT(*) FROM user_subscription_request_history`).Scan(&total); err != nil {
				return err
			}

			rows, err := db.QueryContext(r.Context(), `
				SELECT id, user_uuid, request_ip, user_agent, request_at
				FROM user_subscription_request_history
				ORDER BY request_at DESC
				OFFSET ? LIMIT ?
			`, start, size)
			if err != nil {
				return err
			}
			defer rows.Close()
			for rows.Next() {
				var item historyRecord
				var requestAt time.Time
				if scanErr := rows.Scan(&item.ID, &item.UserUUID, &item.RequestIP, &item.UserAgent, &requestAt); scanErr != nil {
					return scanErr
				}
				item.RequestAt = requestAt.UTC().Format("2006-01-02T15:04:05.000Z")
				records = append(records, item)
			}
			return rows.Err()
		})
		if err != nil {
			shared.SendError(w, http.StatusInternalServerError, "failed to fetch subscription request history", err, cfg)
			return
		}

		shared.WriteJSON(w, http.StatusOK, map[string]any{
			"response": map[string]any{
				"records": records,
				"total":   total,
			},
		})
	}
}

func SubscriptionRequestHistoryStatsHandler(manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			shared.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}

		byParsedApp := make([]map[string]any, 0)
		hourly := make([]map[string]any, 0)

		err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
			rows, err := db.QueryContext(r.Context(), `
				SELECT
					COALESCE(NULLIF(SPLIT_PART(COALESCE(user_agent, ''), '/', 1), ''), 'Unknown') AS app,
					COUNT(*) AS count
				FROM user_subscription_request_history
				GROUP BY app
				ORDER BY count DESC
			`)
			if err != nil {
				return err
			}
			defer rows.Close()
			for rows.Next() {
				var app string
				var count int
				if scanErr := rows.Scan(&app, &count); scanErr != nil {
					return scanErr
				}
				byParsedApp = append(byParsedApp, map[string]any{"app": app, "count": count})
			}
			if err := rows.Err(); err != nil {
				return err
			}

			rows2, err := db.QueryContext(r.Context(), `
				SELECT date_trunc('hour', request_at) AS date_time, COUNT(*) AS request_count
				FROM user_subscription_request_history
				WHERE request_at >= NOW() - INTERVAL '24 hours'
				GROUP BY date_time
				ORDER BY date_time ASC
			`)
			if err != nil {
				return err
			}
			defer rows2.Close()
			for rows2.Next() {
				var dt time.Time
				var count int
				if scanErr := rows2.Scan(&dt, &count); scanErr != nil {
					return scanErr
				}
				hourly = append(hourly, map[string]any{
					"dateTime":     dt.UTC().Format("2006-01-02T15:04:05.000Z"),
					"requestCount": count,
				})
			}
			return rows2.Err()
		})
		if err != nil {
			shared.SendError(w, http.StatusInternalServerError, "failed to fetch subscription request history stats", err, cfg)
			return
		}

		shared.WriteJSON(w, http.StatusOK, map[string]any{
			"response": map[string]any{
				"byParsedApp":        byParsedApp,
				"hourlyRequestStats": hourly,
			},
		})
	}
}

func parseInt(raw string, def, min, max int) int {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return def
	}
	v, err := strconv.Atoi(raw)
	if err != nil {
		return def
	}
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}
