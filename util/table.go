package util

import (
	"database/sql"
	"fmt"
	"slices"
	"strings"
	"time"
	"v2ray-stat/backend/config"
	"v2ray-stat/common"
)

func Contains(slice []string, item string) bool {
	return slices.Contains(slice, item)
}

func AppendStats(builder *strings.Builder, content string) {
	builder.WriteString(content)
}

func FormatTable(rows *sql.Rows, trafficColumns []string, cfg *config.BackendConfig) (string, error) {
	columns, err := rows.Columns()
	if err != nil {
		return "", fmt.Errorf("failed to get column names: %v", err)
	}

	maxWidths := make([]int, len(columns))
	for i, col := range columns {
		maxWidths[i] = len(col)
	}

	var data [][]string
	for rows.Next() {
		values := make([]any, len(columns))
		valuePtrs := make([]any, len(columns))
		for i := range columns {
			valuePtrs[i] = &values[i]
		}

		if err := rows.Scan(valuePtrs...); err != nil {
			return "", fmt.Errorf("failed to scan row: %v", err)
		}

		row := make([]string, len(columns))
		for i, val := range values {
			strVal := ""
			columnName := columns[i]

			var floatVal float64
			isNum := false
			switch v := val.(type) {
			case int64:
				floatVal = float64(v)
				isNum = true
			case float64:
				floatVal = v
				isNum = true
			case []byte:
				strVal = string(v)
			case string:
				strVal = v
			case nil:
				strVal = ""
			default:
				strVal = fmt.Sprintf("%v", v)
			}

			// Если это число и колонка относится к трафику/скорости
			if isNum && Contains(trafficColumns, columnName) {
				var unit string
				if columnName == "Rate" {
					unit = cfg.Monitor.RateUnit
					if unit == "" {
						unit = "bps"
					}
				} else {
					unit = cfg.Monitor.TrafficUnit
					if unit == "" {
						unit = "byte"
					}
				}
				strVal = FormatData(floatVal, unit)
			} else if isNum && !Contains(trafficColumns, columnName) {
				// Обработка дат, если они пришли как int64
				if (columnName == "created" || columnName == "sub_end" || columnName == "last_seen") && floatVal > 0 {
					strVal = time.Unix(int64(floatVal), 0).In(common.TimeLocation).Format("2006-01-02 15:04")
				} else {
					strVal = fmt.Sprintf("%.0f", floatVal)
				}
			}

			// Ограничение длины для таблицы
			if len(strVal) > 255 {
				strVal = strVal[:255]
			}
			row[i] = strVal
			if len(strVal) > maxWidths[i] {
				maxWidths[i] = len(strVal)
			}
		}
		data = append(data, row)
	}

	// Сборка таблицы
	var table strings.Builder

	// Заголовки
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
			if Contains(trafficColumns, columns[i]) {
				// Числа выравниваем по правому краю для красоты
				table.WriteString(fmt.Sprintf("%*s  ", maxWidths[i], val))
			} else {
				table.WriteString(fmt.Sprintf("%-*s", maxWidths[i]+2, val))
			}
		}
		table.WriteString("\n")
	}

	return table.String(), nil
}
