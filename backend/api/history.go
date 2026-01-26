package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"v2ray-stat/backend/config"
	"v2ray-stat/backend/db/manager"
	"v2ray-stat/common"
)

type StatEntry struct {
	Date     string `json:"date"` // Диапазон или конкретный час
	Name     string `json:"name"` // user или source
	Uplink   int64  `json:"uplink"`
	Downlink int64  `json:"downlink"`
}

type HistoryResponse struct {
	Stats   []StatEntry `json:"stats"`
	Period  string      `json:"period"`   // "hour", "day", etc.
	From    string      `json:"from"`     // YYYY-MM-DD HH
	To      string      `json:"to"`       // YYYY-MM-DD HH
	GroupBy string      `json:"group_by"` // "name" (SUM) или "date" (по часам)
	Type    string      `json:"type"`     // "user" или "source"
	Error   string      `json:"error,omitempty"`
}

// HistoryHandler обрабатывает запросы статистики трафика
func HistoryHandler(mgr *manager.DatabaseManager, cfg *config.BackendConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Парсинг параметров
		period := strings.ToLower(r.URL.Query().Get("period"))    // hour, day, week, month, year
		start := r.URL.Query().Get("start")                       // YYYY-MM-DD или YYYY-MM-DD HH
		end := r.URL.Query().Get("end")                           // YYYY-MM-DD или YYYY-MM-DD HH
		histType := strings.ToLower(r.URL.Query().Get("type"))    // user или source (default: user)
		groupBy := strings.ToLower(r.URL.Query().Get("group_by")) // name (default, SUM), date (по часам)

		if histType != "user" && histType != "source" {
			histType = "user"
		}
		if groupBy != "date" {
			groupBy = "name" // default: SUM по именам
		}

		// Вычисление диапазона дат
		dateFrom, dateTo, err := calculateDateRange(period, start, end)
		if err != nil {
			resp := HistoryResponse{Error: err.Error()}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(resp)
			return
		}

		table := "daily_user_stats"
		nameCol := "user"
		if histType == "source" {
			table = "daily_source_stats"
			nameCol = "source"
		}

		var stats []StatEntry
		resp := HistoryResponse{
			Period:  period,
			From:    dateFrom,
			To:      dateTo,
			GroupBy: groupBy,
			Type:    histType,
		}

		err = mgr.ExecuteHighPriority(func(db *sql.DB) error {
			var rows *sql.Rows
			var query string

			if groupBy == "date" {
				// Детализация по часам (GROUP BY date, name)
				query = fmt.Sprintf(`
					SELECT date, %s AS name, SUM(uplink) AS total_uplink, SUM(downlink) AS total_downlink
					FROM %s
					WHERE date BETWEEN ? AND ?
					GROUP BY date, %s
					ORDER BY date DESC, (total_uplink + total_downlink) DESC
				`, nameCol, table, nameCol)
			} else {
				// SUM по именам за весь период (GROUP BY name)
				query = fmt.Sprintf(`
					SELECT %s AS name, SUM(uplink) AS total_uplink, SUM(downlink) AS total_downlink
					FROM %s
					WHERE date BETWEEN ? AND ?
					GROUP BY %s
					ORDER BY (total_uplink + total_downlink) DESC
				`, nameCol, table, nameCol)
			}

			var queryErr error
			rows, queryErr = db.QueryContext(r.Context(), query, dateFrom, dateTo)
			if queryErr != nil {
				return fmt.Errorf("query failed: %w", queryErr)
			}
			defer rows.Close()

			for rows.Next() {
				var s StatEntry
				if groupBy == "date" {
					if err := rows.Scan(&s.Date, &s.Name, &s.Uplink, &s.Downlink); err != nil {
						return err
					}
				} else {
					if err := rows.Scan(&s.Name, &s.Uplink, &s.Downlink); err != nil {
						return err
					}
					s.Date = fmt.Sprintf("%s to %s", dateFrom, dateTo)
				}
				stats = append(stats, s)
			}

			return rows.Err()
		})

		if err != nil {
			cfg.Logger.Error("History API failed", "error", err, "period", period, "type", histType)
			resp.Error = "Internal Server Error"
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			resp.Stats = stats
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}

func calculateDateRange(period, start, end string) (string, string, error) {
	currentTime := common.GetLocalUnix()
	now := time.Unix(currentTime, 0)
	layoutHour := "2006-01-02 15" // YYYY-MM-DD HH

	var fromTime, toTime time.Time
	if start != "" && end != "" {
		// Custom интервал: парсим start/end (поддержка YYYY-MM-DD или YYYY-MM-DD HH)
		var err error
		fromTime, err = time.Parse("2006-01-02", start)
		if err != nil {
			fromTime, err = time.Parse(layoutHour, start)
			if err != nil {
				return "", "", fmt.Errorf("invalid start format: %w", err)
			}
		} else {
			// Если только дата, добавляем 00 час
			fromTime = time.Date(fromTime.Year(), fromTime.Month(), fromTime.Day(), 0, 0, 0, 0, now.Location())
		}

		toTime, err = time.Parse("2006-01-02", end)
		if err != nil {
			toTime, err = time.Parse(layoutHour, end)
			if err != nil {
				return "", "", fmt.Errorf("invalid end format: %w", err)
			}
		} else {
			// Если только дата, добавляем 23 час
			toTime = time.Date(toTime.Year(), toTime.Month(), toTime.Day(), 23, 0, 0, 0, now.Location())
		}
	} else {
		// Предустановленные периоды
		switch period {
		case "hour":
			fromTime = now.Add(-1 * time.Hour).Truncate(time.Hour)
			toTime = now.Truncate(time.Hour)
		case "day":
			fromTime = now.Add(-24 * time.Hour).Truncate(24 * time.Hour)
			toTime = now.Truncate(24 * time.Hour)
		case "week":
			fromTime = now.AddDate(0, 0, -7).Truncate(24 * time.Hour)
			toTime = now.Truncate(24 * time.Hour)
		case "month":
			fromTime = now.AddDate(0, -1, 0).Truncate(24 * time.Hour)
			toTime = now.Truncate(24 * time.Hour)
		case "year":
			fromTime = now.AddDate(-1, 0, 0).Truncate(24 * time.Hour)
			toTime = now.Truncate(24 * time.Hour)
		default:
			return "", "", fmt.Errorf("invalid period: use 'hour', 'day', 'week', 'month', 'year'")
		}
	}

	if fromTime.After(toTime) {
		return "", "", fmt.Errorf("start date cannot be after end date")
	}

	return fromTime.Format(layoutHour), toTime.Format(layoutHour), nil
}
