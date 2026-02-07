package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"
	"time"
	"v2ray-stat/backend/config"
	"v2ray-stat/backend/db/manager"
	"v2ray-stat/util"
)

func HistoryStatsHandler(mgr *manager.DatabaseManager, cfg *config.BackendConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
			return
		}

		// Тип: user (по умолчанию) или source
		histType := strings.ToLower(r.URL.Query().Get("type"))
		if histType != "source" {
			histType = "user"
		}

		// Фильтр по имени (конкретный юзер)
		filterName := r.URL.Query().Get("name")

		// Умная дата: "2026", "2026-01", "2026-01-26"
		dateInput := r.URL.Query().Get("date")

		// Режим суммы: если передано sum=1 или sum=true
		isSumMode := r.URL.Query().Get("sum") == "1" || strings.ToLower(r.URL.Query().Get("sum")) == "true"

		limit := r.URL.Query().Get("limit")
		if limit == "" {
			limit = "500"
		}

		w.Header().Set("Content-Type", "text/plain; charset=utf-8")

		var statsBuilder strings.Builder

		// Вызываем построитель
		err := buildHistoryTableStats(&statsBuilder, mgr, cfg, histType, filterName, dateInput, isSumMode, limit)
		if err != nil {
			cfg.Logger.Error("Failed to build history table", "error", err)
			fmt.Fprintf(w, "Error: %v\n", err)
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
	histType, filterName, dateInput string,
	isSumMode bool,
	limit string,
) error {
	tableName := "daily_user_stats"
	colName := "user"
	if histType == "source" {
		tableName = "daily_source_stats"
		colName = "source"
	}

	start, end := resolveSmartDate(dateInput)

	var query string
	var args []any
	whereClause := "WHERE date BETWEEN ? AND ?"
	args = append(args, start, end)

	if filterName != "" {
		whereClause += fmt.Sprintf(" AND %s = ?", colName)
		args = append(args, filterName)
	}

	// Колонки, которые FormatTable превратит в GiB/Mbps
	trafficAliases := []string{"Uplink", "Downlink", "Total", "Rate"}

	if isSumMode {
		// Используем подзапрос, чтобы сначала все сложить, а потом посчитать Rate
		// Поддержка формата даты YYYY-MM-DD-HH или YYYY-MM-DD HH
		durationCalc := `(
			strftime('%s', substr(MAX(date), 1, 10) || ' ' || substr(MAX(date), 12, 2) || ':00:00') - 
			strftime('%s', substr(MIN(date), 1, 10) || ' ' || substr(MIN(date), 12, 2) || ':00:00') + 3600
		)`

		query = fmt.Sprintf(`
			SELECT Period %[2]s (Total * 8.0 / Duration) AS "Rate", Uplink, Downlink, Total
			FROM (
				SELECT 
					MIN(date) || ' - ' || MAX(date) AS "Period",
					%[1]s AS "%[2]s",
					SUM(uplink) AS "Uplink",
					SUM(downlink) AS "Downlink",
					SUM(uplink + downlink) AS "Total",
					CAST(%[3]s AS FLOAT) AS "Duration"
				FROM %[4]s
				%[5]s
				GROUP BY %[1]s
			)
			ORDER BY Total DESC
			LIMIT %[6]s
		`, colName, strings.Title(colName), durationCalc, tableName, whereClause, limit)

	} else {
		// Обычный режим (почасовой)
		query = fmt.Sprintf(`
			SELECT 
				date AS "Date",
				%[1]s AS "%[2]s",
				((uplink + downlink) * 8.0 / 3600.0) AS "Rate",
				uplink AS "Uplink",
				downlink AS "Downlink",
				(uplink + downlink) AS "Total"
			FROM %[3]s
			%[4]s
			ORDER BY date ASC
			LIMIT %[5]s
		`, colName, strings.Title(colName), tableName, whereClause, limit)
	}

	return mgr.ExecuteHighPriority(func(db *sql.DB) error {
		rows, err := db.Query(query, args...)
		if err != nil {
			return fmt.Errorf("query failed: %w", err)
		}
		defer rows.Close()

		formatted, err := util.FormatTable(rows, trafficAliases, cfg)
		if err != nil {
			return fmt.Errorf("format table failed: %w", err)
		}

		if formatted == "" {
			builder.WriteString("No data found.\n")
		} else {
			builder.WriteString(formatted)
		}
		return nil
	})
}

func resolveSmartDate(inputDate string) (string, string) {
	inputDate = strings.TrimSpace(inputDate)

	// ВАЖНОЕ ИСПРАВЛЕНИЕ:
	// Если дата не указана, берем диапазон от "начала времен" до далекого будущего,
	// чтобы SQL вернул вообще ВСЕ записи.
	if inputDate == "" {
		return "0000-01-01 00", "9999-12-31 23"
	}

	// Поддержка диапазона через запятую: "2026-01-26 03,2026-01-26 09"
	if strings.Contains(inputDate, ",") {
		parts := strings.Split(inputDate, ",")
		return strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])
	}

	// 1. Только год: "2026"
	if len(inputDate) == 4 {
		return inputDate + "-01-01 00", inputDate + "-12-31 23"
	}

	// 2. Год и месяц: "2026-01"
	if len(inputDate) == 7 {
		t, err := time.Parse("2006-01", inputDate)
		if err == nil {
			lastDay := t.AddDate(0, 1, -1).Day()
			return fmt.Sprintf("%s-01 00", inputDate), fmt.Sprintf("%s-%02d 23", inputDate, lastDay)
		}
	}

	// 3. Полная дата с часом: "2026-01-26 03" (13 символов)
	// Важно: в базе используется пробел, а не тире
	if len(inputDate) == 13 {
		return inputDate, inputDate
	}

	// 4. Полная дата: "2026-01-26" -> "2026-01-26 00" ... "2026-01-26 23"
	if len(inputDate) >= 10 {
		baseDate := inputDate[:10]
		return baseDate + " 00", baseDate + " 23"
	}

	return inputDate, inputDate
}
