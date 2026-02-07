package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"v2ray-stat/backend/config"
	"v2ray-stat/backend/db/manager"
)

// HistoryRecord описывает одну строку статистики
type HistoryRecord struct {
	Date     string  `json:"date,omitempty"`   // Используется в обычном режиме
	Period   string  `json:"period,omitempty"` // Используется в sum=1
	Name     string  `json:"name"`             // Имя юзера или источника
	Rate     float64 `json:"rate_bps"`         // Скорость в бит/с
	Uplink   int64   `json:"uplink_bytes"`
	Downlink int64   `json:"downlink_bytes"`
	Total    int64   `json:"total_bytes"`
}

func HistoryJSONHandler(mgr *manager.DatabaseManager, cfg *config.BackendConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Параметры запроса аналогичны текстовой версии
		histType := strings.ToLower(r.URL.Query().Get("type"))
		if histType != "source" {
			histType = "user"
		}

		filterName := r.URL.Query().Get("name")
		dateInput := r.URL.Query().Get("date")
		isSumMode := r.URL.Query().Get("sum") == "1" || strings.ToLower(r.URL.Query().Get("sum")) == "true"
		limit := r.URL.Query().Get("limit")
		if limit == "" {
			limit = "200"
		}

		records, err := fetchHistoryData(mgr, histType, filterName, dateInput, isSumMode, limit)
		if err != nil {
			cfg.Logger.Error("Failed to fetch history JSON", "error", err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(records)
	}
}

func fetchHistoryData(
	mgr *manager.DatabaseManager,
	histType, filterName, dateInput string,
	isSumMode bool,
	limit string,
) ([]HistoryRecord, error) {
	tableName := "daily_user_stats"
	colName := "user"
	if histType == "source" {
		tableName = "daily_source_stats"
		colName = "source"
	}

	start, end := resolveSmartDate(dateInput)
	var args []any
	whereClause := "WHERE date BETWEEN ? AND ?"
	args = append(args, start, end)

	if filterName != "" {
		whereClause += fmt.Sprintf(" AND %s = ?", colName)
		args = append(args, filterName)
	}

	var query string
	if isSumMode {
		durationCalc := `(
            strftime('%s', substr(MAX(date), 1, 10) || ' ' || substr(MAX(date), 12, 2) || ':00:00') - 
            strftime('%s', substr(MIN(date), 1, 10) || ' ' || substr(MIN(date), 12, 2) || ':00:00') + 3600
        )`
		query = fmt.Sprintf(`
            SELECT Period, Name, (Total * 8.0 / Duration) AS Rate, Uplink, Downlink, Total
            FROM (
                SELECT 
                    MIN(date) || ' - ' || MAX(date) AS Period,
                    %[1]s AS Name,
                    SUM(uplink) AS Uplink,
                    SUM(downlink) AS Downlink,
                    SUM(uplink + downlink) AS Total,
                    CAST(%[2]s AS FLOAT) AS Duration
                FROM %[3]s
                %[4]s
                GROUP BY %[1]s
            )
            ORDER BY Total DESC
            LIMIT %[5]s
        `, colName, durationCalc, tableName, whereClause, limit)
	} else {
		query = fmt.Sprintf(`
            SELECT 
                date AS Date,
                %[1]s AS Name,
                ((uplink + downlink) * 8.0 / 3600.0) AS Rate,
                uplink AS Uplink,
                downlink AS Downlink,
                (uplink + downlink) AS Total
            FROM %[2]s
            %[3]s
            ORDER BY date DESC
            LIMIT %[4]s
        `, colName, tableName, whereClause, limit)
	}

	var results []HistoryRecord
	err := mgr.ExecuteHighPriority(func(db *sql.DB) error {
		rows, err := db.Query(query, args...)
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			var r HistoryRecord
			if isSumMode {
				if err := rows.Scan(&r.Period, &r.Name, &r.Rate, &r.Uplink, &r.Downlink, &r.Total); err != nil {
					return err
				}
			} else {
				if err := rows.Scan(&r.Date, &r.Name, &r.Rate, &r.Uplink, &r.Downlink, &r.Total); err != nil {
					return err
				}
			}
			results = append(results, r)
		}
		return nil
	})

	return results, err
}
