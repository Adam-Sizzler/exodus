package subscriptionrequesthistory

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"exodus/internal/config"
	dbmanager "exodus/internal/db/manager"
	"exodus/internal/httpapi/shared"
)

type historyRecord struct {
	ID        int64   `json:"id"`
	UserID    int64   `json:"userId"`
	RequestIP *string `json:"requestIp"`
	UserAgent *string `json:"userAgent"`
	RequestAt string  `json:"requestAt"`
}

type tableFilter struct {
	ID    string `json:"id"`
	Value any    `json:"value"`
}

type tableSorting struct {
	ID   string `json:"id"`
	Desc bool   `json:"desc"`
}

func SubscriptionRequestHistoryHandler(manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			shared.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}

		start := parseInt(r.URL.Query().Get("start"), 0, 0, 1_000_000)
		size := parseInt(r.URL.Query().Get("size"), 25, 1, 500)
		columns := map[string]string{
			"id":        "id",
			"userId":    "user_id",
			"requestIp": "request_ip",
			"userAgent": "user_agent",
			"requestAt": "request_at",
		}
		whereSQL, whereArgs := buildTableWhereClause(r, columns)
		orderSQL := buildTableOrderClause(r, columns, "request_at DESC")

		records := make([]historyRecord, 0)
		total := 0
		err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
			countQuery := `SELECT COUNT(*) FROM user_subscription_request_history` + whereSQL
			if err := db.QueryRowContext(r.Context(), countQuery, whereArgs...).Scan(&total); err != nil {
				return err
			}

			args := append(append([]any{}, whereArgs...), start, size)
			query := `
				SELECT id, user_id, request_ip, user_agent, request_at
				FROM user_subscription_request_history
			` + whereSQL + orderSQL + ` OFFSET ? LIMIT ?`
			rows, err := db.QueryContext(r.Context(), query, args...)
			if err != nil {
				return err
			}
			defer rows.Close()
			for rows.Next() {
				var item historyRecord
				var requestAt sql.NullTime
				if scanErr := rows.Scan(&item.ID, &item.UserID, &item.RequestIP, &item.UserAgent, &requestAt); scanErr != nil {
					return scanErr
				}
				if requestAt.Valid {
					item.RequestAt = requestAt.Time.UTC().Format("2006-01-02T15:04:05.000Z")
				}
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
				WHERE request_at >= NOW() - INTERVAL '48 hours'
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

func buildTableWhereClause(r *http.Request, columns map[string]string) (string, []any) {
	filters := parseTableFilters(r.URL.Query().Get("filters"))
	if len(filters) == 0 {
		return "", nil
	}

	parts := make([]string, 0, len(filters))
	args := make([]any, 0, len(filters))
	for _, filter := range filters {
		column, ok := columns[filter.ID]
		if !ok {
			continue
		}
		value := tableFilterValue(filter.Value)
		if value == "" {
			continue
		}
		parts = append(parts, "LOWER(COALESCE("+column+"::text, '')) LIKE LOWER(?)")
		args = append(args, "%"+value+"%")
	}
	if len(parts) == 0 {
		return "", nil
	}
	return " WHERE " + strings.Join(parts, " AND "), args
}

func buildTableOrderClause(r *http.Request, columns map[string]string, fallback string) string {
	sorting := parseTableSorting(r.URL.Query().Get("sorting"))
	if len(sorting) == 0 {
		return " ORDER BY " + fallback
	}

	parts := make([]string, 0, len(sorting))
	for _, sort := range sorting {
		column, ok := columns[sort.ID]
		if !ok {
			continue
		}
		direction := "ASC"
		if sort.Desc {
			direction = "DESC"
		}
		parts = append(parts, column+" "+direction+" NULLS LAST")
	}
	if len(parts) == 0 {
		return " ORDER BY " + fallback
	}
	return " ORDER BY " + strings.Join(parts, ", ")
}

func parseTableFilters(raw string) []tableFilter {
	raw = strings.TrimSpace(raw)
	if raw == "" || raw == "[]" {
		return nil
	}
	var filters []tableFilter
	if err := json.Unmarshal([]byte(raw), &filters); err != nil {
		return nil
	}
	return filters
}

func parseTableSorting(raw string) []tableSorting {
	raw = strings.TrimSpace(raw)
	if raw == "" || raw == "[]" {
		return nil
	}
	var sorting []tableSorting
	if err := json.Unmarshal([]byte(raw), &sorting); err != nil {
		return nil
	}
	return sorting
}

func tableFilterValue(value any) string {
	if value == nil {
		return ""
	}
	switch v := value.(type) {
	case string:
		return strings.TrimSpace(v)
	default:
		return strings.TrimSpace(fmt.Sprint(v))
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
