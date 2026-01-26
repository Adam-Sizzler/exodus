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
