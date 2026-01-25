package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"

	"v2ray-stat/backend/config"
	"v2ray-stat/backend/db/manager"
	"v2ray-stat/util"
)

func HistoryStatsHandler(mgr *manager.DatabaseManager, cfg *config.BackendConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cfg.Logger.Debug("Starting HistoryStatsHandler request processing")

		// Только GET
		if r.Method != http.MethodGet {
			http.Error(w, "Invalid method. Use GET", http.StatusMethodNotAllowed)
			return
		}

		// Парсинг параметров
		period := strings.ToLower(r.URL.Query().Get("period"))
		start := r.URL.Query().Get("start")
		end := r.URL.Query().Get("end")
		histType := strings.ToLower(r.URL.Query().Get("type"))
		groupBy := strings.ToLower(r.URL.Query().Get("group_by"))

		if histType != "user" && histType != "source" {
			histType = "user"
		}
		if groupBy != "date" {
			groupBy = "name"
		}

		w.Header().Set("Content-Type", "text/plain; charset=utf-8")

		var statsBuilder strings.Builder

		err := buildHistoryTableStats(&statsBuilder, mgr, cfg, period, start, end, histType, groupBy)
		if err != nil {
			cfg.Logger.Error("Failed to build history table", "error", err)
			fmt.Fprintf(w, "Error retrieving history statistics: %v\n", err)
			return
		}

		if statsBuilder.Len() == 0 {
			fmt.Fprintln(w, "No history data found.")
		} else {
			fmt.Fprint(w, statsBuilder.String())
		}
	}
}

func buildHistoryTableStats(
	builder *strings.Builder,
	mgr *manager.DatabaseManager,
	cfg *config.BackendConfig,
	period, start, end, histType, groupBy string,
) error {
	dateFrom, dateTo, err := calculateDateRange(period, start, end)
	if err != nil {
		return err
	}

	table := "daily_user_stats"
	nameCol := "user"
	nameAlias := "User"
	if histType == "source" {
		table = "daily_source_stats"
		nameCol = "source"
		nameAlias = "Source"
	}

	// Колонки, которые нужно форматировать как трафик (MiB, GiB и т.д.)
	trafficAliases := []string{"Uplink", "Downlink", "Total"}

	var query string
	var args []interface{}

	if groupBy == "date" {
		// Детализация по часам — Date = реальный час
		query = fmt.Sprintf(`
			SELECT date AS "Date",
			       %s AS "%s",
			       SUM(uplink)   AS "Uplink",
			       SUM(downlink) AS "Downlink",
			       (SUM(uplink) + SUM(downlink)) AS "Total"
			FROM %s
			WHERE date BETWEEN ? AND ?
			GROUP BY date, %s
			ORDER BY date DESC, "Total" DESC
		`, nameCol, nameAlias, table, nameCol)
		args = []interface{}{dateFrom, dateTo}
	} else {
		// Сумма за период — Date = диапазон
		query = fmt.Sprintf(`
			SELECT ? || ' to ' || ? AS "Date",
			       %s AS "%s",
			       SUM(uplink)   AS "Uplink",
			       SUM(downlink) AS "Downlink",
			       (SUM(uplink) + SUM(downlink)) AS "Total"
			FROM %s
			WHERE date BETWEEN ? AND ?
			GROUP BY %s
			ORDER BY "Total" DESC
		`, nameCol, nameAlias, table, nameCol)
		args = []interface{}{dateFrom, dateTo, dateFrom, dateTo}
	}

	return mgr.ExecuteHighPriority(func(db *sql.DB) error {
		rows, err := db.Query(query, args...)
		if err != nil {
			return fmt.Errorf("query failed: %w", err)
		}
		defer rows.Close()

		// Форматируем таблицу через вашу утилиту
		formatted, err := util.FormatTable(rows, trafficAliases, cfg)
		if err != nil {
			return fmt.Errorf("format table failed: %w", err)
		}

		// Добавляем красивый заголовок
		header := fmt.Sprintf("History: %s (%s to %s)\n\n", histType, dateFrom, dateTo)
		builder.WriteString(header)

		if formatted == "" {
			builder.WriteString("No data found in this period.\n")
		} else {
			builder.WriteString(formatted)
		}

		return nil
	})
}
