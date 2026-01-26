package api

import (
	"database/sql"
	"fmt"
	"time"
	"net/http"
	"strings"

	"v2ray-stat/common"
	"v2ray-stat/util"
	"v2ray-stat/backend/config"
	"v2ray-stat/backend/db/manager"
)

func HistoryStatsHandler(mgr *manager.DatabaseManager, cfg *config.BackendConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cfg.Logger.Debug("Starting HistoryStatsHandler request processing")

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
		name := r.URL.Query().Get("name") // Фильтр по конкретному user/source
		sum := strings.ToLower(r.URL.Query().Get("sum"))
		dateFilter := r.URL.Query().Get("date") // Умный поиск по дате

		if histType != "user" && histType != "source" {
			histType = "user"
		}
		if groupBy != "date" && groupBy != "hour" {
			groupBy = "name"
		}

		w.Header().Set("Content-Type", "text/plain; charset=utf-8")

		var statsBuilder strings.Builder

		err := buildHistoryTableStats(&statsBuilder, mgr, cfg, period, start, end, 
			histType, groupBy, name, sum, dateFilter)
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
	period, start, end, histType, groupBy, name, sum, dateFilter string,
) error {
	// Определяем таблицу и колонку в зависимости от типа
	table := "daily_user_stats"
	nameCol := "user"
	nameAlias := "User"
	if histType == "source" {
		table = "daily_source_stats"
		nameCol = "source"
		nameAlias = "Source"
	}

	// Обработка параметра date (умный поиск)
	dateCondition := ""
	var dateArgs []interface{}
	if dateFilter != "" {
		dateCondition = fmt.Sprintf("date LIKE ?")
		dateArgs = []interface{}{dateFilter + "%"}
	} else {
		// Используем старую логику с периодом
		dateFrom, dateTo, err := calculateDateRange(period, start, end)
		if err != nil {
			return err
		}
		dateCondition = "date BETWEEN ? AND ?"
		dateArgs = []interface{}{dateFrom, dateTo}
	}

	// Формируем условия WHERE
	var conditions []string
	var args []interface{}
	
	conditions = append(conditions, dateCondition)
	args = append(args, dateArgs...)
	
	if name != "" {
		conditions = append(conditions, fmt.Sprintf("%s = ?", nameCol))
		args = append(args, name)
	}
	
	whereClause := strings.Join(conditions, " AND ")
	if whereClause == "" {
		whereClause = "1=1"
	}

	// Режим SUM
	if sum == "true" {
		return mgr.ExecuteHighPriority(func(db *sql.DB) error {
			// Запрос для суммарной статистики
			query := fmt.Sprintf(`
				SELECT 
					MIN(date) || ' - ' || MAX(date) AS "Period",
					%s AS "%s",
					SUM(uplink) AS "Uplink",
					SUM(downlink) AS "Downlink",
					(SUM(uplink) + SUM(downlink)) AS "Total"
				FROM %s
				WHERE %s
				GROUP BY %s
				ORDER BY "Total" DESC
			`, nameCol, nameAlias, table, whereClause, nameCol)

			rows, err := db.Query(query, args...)
			if err != nil {
				return fmt.Errorf("sum query failed: %w", err)
			}
			defer rows.Close()

			columns, err := rows.Columns()
			if err != nil {
				return fmt.Errorf("failed to get columns: %w", err)
			}

			// Добавляем колонку для средней скорости
			columns = append(columns, "Rate_avg")
			trafficAliases := []string{"Uplink", "Downlink", "Total", "Rate_avg"}

			var data [][]string
			for rows.Next() {
				values := make([]interface{}, len(columns)-1) // -1 потому что Rate_avg нет в запросе
				valuePtrs := make([]interface{}, len(values))
				for i := range values {
					valuePtrs[i] = &values[i]
				}

				if err := rows.Scan(valuePtrs...); err != nil {
					return fmt.Errorf("failed to scan row: %w", err)
				}

				row := make([]string, len(columns))
				var totalBytes int64
				var durationHours float64
				
				for i, val := range values {
					strVal := ""
					switch v := val.(type) {
					case int64:
						if columns[i] == "Uplink" || columns[i] == "Downlink" {
							strVal = util.FormatData(float64(v), "byte")
						} else if columns[i] == "Total" {
							totalBytes = v
							strVal = util.FormatData(float64(v), "byte")
						} else {
							strVal = fmt.Sprintf("%d", v)
						}
					case string:
						strVal = v
						// Пытаемся вычислить длительность периода из строки Period
						if columns[i] == "Period" && v != "" {
							// Формат: "YYYY-MM-DD HH - YYYY-MM-DD HH"
							parts := strings.Split(v, " - ")
							if len(parts) == 2 {
								startTime, err1 := time.Parse("2006-01-02 15", parts[0])
								endTime, err2 := time.Parse("2006-01-02 15", parts[1])
								if err1 == nil && err2 == nil {
									durationHours = endTime.Sub(startTime).Hours()
									if durationHours < 1 {
										durationHours = 1
									}
								}
							}
						}
					case nil:
						strVal = ""
					default:
						strVal = fmt.Sprintf("%v", v)
					}
					row[i] = strVal
				}

				// Вычисляем среднюю скорость
				if durationHours > 0 {
					// Переводим байты в биты (умножаем на 8) и делим на количество секунд
					rateBps := float64(totalBytes) * 8 / (durationHours * 3600)
					row[len(row)-1] = util.FormatData(rateBps, "bps")
				} else {
					row[len(row)-1] = "N/A"
				}

				data = append(data, row)
			}

			// Форматируем таблицу
			formatted := formatCustomTable(columns, data, trafficAliases)
			
			header := "Summary Statistics\n\n"
			builder.WriteString(header)
			if formatted == "" {
				builder.WriteString("No data found.\n")
			} else {
				builder.WriteString(formatted)
			}

			return nil
		})
	}

	// Режим детальной статистики
	return mgr.ExecuteHighPriority(func(db *sql.DB) error {
		var query string
		var queryArgs []interface{}

		if groupBy == "date" || groupBy == "hour" {
			// Детализация по часам
			query = fmt.Sprintf(`
				SELECT date AS "Date",
				       %s AS "%s",
				       SUM(uplink)   AS "Uplink",
				       SUM(downlink) AS "Downlink",
				       (SUM(uplink) + SUM(downlink)) AS "Total"
				FROM %s
				WHERE %s
				GROUP BY date, %s
				ORDER BY date DESC, "Total" DESC
			`, nameCol, nameAlias, table, whereClause, nameCol)
			queryArgs = args
		} else {
			// Группировка по имени
			query = fmt.Sprintf(`
				SELECT ? || ' to ' || ? AS "Period",
				       %s AS "%s",
				       SUM(uplink)   AS "Uplink",
				       SUM(downlink) AS "Downlink",
				       (SUM(uplink) + SUM(downlink)) AS "Total"
				FROM %s
				WHERE %s
				GROUP BY %s
				ORDER BY "Total" DESC
			`, nameCol, nameAlias, table, whereClause, nameCol)
			
			// Получаем фактические даты из данных для отображения периода
			minMaxQuery := fmt.Sprintf(`
				SELECT MIN(date), MAX(date) FROM %s WHERE %s
			`, table, whereClause)
			
			var minDate, maxDate string
			err := db.QueryRow(minMaxQuery, args...).Scan(&minDate, &maxDate)
			if err != nil && err != sql.ErrNoRows {
				return fmt.Errorf("failed to get date range: %w", err)
			}
			
			if minDate == "" || maxDate == "" {
				minDate = "N/A"
				maxDate = "N/A"
			}
			
			queryArgs = []interface{}{minDate, maxDate}
			queryArgs = append(queryArgs, args...)
		}

		rows, err := db.Query(query, queryArgs...)
		if err != nil {
			return fmt.Errorf("query failed: %w", err)
		}
		defer rows.Close()

		trafficAliases := []string{"Uplink", "Downlink", "Total"}
		formatted, err := util.FormatTable(rows, trafficAliases, cfg)
		if err != nil {
			return fmt.Errorf("format table failed: %w", err)
		}

		// Добавляем заголовок
		var header string
		if name != "" {
			header = fmt.Sprintf("History for %s: %s\n\n", histType, name)
		} else {
			header = fmt.Sprintf("History: %s\n\n", histType)
		}
		
		if dateFilter != "" {
			header += fmt.Sprintf("Date filter: %s\n\n", dateFilter)
		}
		
		builder.WriteString(header)

		if formatted == "" {
			builder.WriteString("No data found in this period.\n")
		} else {
			builder.WriteString(formatted)
		}

		return nil
	})
}

// Вспомогательная функция для форматирования таблицы с custom данными
func formatCustomTable(columns []string, data [][]string, trafficColumns []string) string {
	if len(data) == 0 {
		return ""
	}

	// Вычисляем максимальные ширины колонок
	maxWidths := make([]int, len(columns))
	for i, col := range columns {
		maxWidths[i] = len(col)
	}

	for _, row := range data {
		for i, val := range row {
			if len(val) > maxWidths[i] {
				maxWidths[i] = len(val)
			}
		}
	}

	var table strings.Builder

	// Заголовок
	for i, col := range columns {
		table.WriteString(fmt.Sprintf("%-*s", maxWidths[i]+2, col))
	}
	table.WriteString("\n")

	// Разделитель
	for _, width := range maxWidths {
		table.WriteString(strings.Repeat("-", width) + "  ")
	}
	table.WriteString("\n")

	// Данные
	for _, row := range data {
		for i, val := range row {
			// Выравнивание для числовых колонок
			if contains(trafficColumns, columns[i]) {
				table.WriteString(fmt.Sprintf("%*s  ", maxWidths[i], val))
			} else {
				table.WriteString(fmt.Sprintf("%-*s", maxWidths[i]+2, val))
			}
		}
		table.WriteString("\n")
	}

	return table.String()
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// calculateDateRange вычисляет диапазон для hourly статистики
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