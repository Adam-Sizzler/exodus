package bandwidthstats

import (
	"database/sql"
	"hash/fnv"
	"net/http"
	"strconv"
	"strings"
	"time"

	"cerberus/backend/config"
	dbmanager "cerberus/backend/db/manager"
	"cerberus/backend/httpapi/shared"
)

var palette = []string{
	"#3b82f6", "#06b6d4", "#22c55e", "#eab308", "#f97316", "#ef4444", "#8b5cf6", "#ec4899",
}

type nodeRealtimeUsage struct {
	NodeUUID         string `json:"nodeUuid"`
	NodeName         string `json:"nodeName"`
	CountryCode      string `json:"countryCode"`
	DownloadBytes    int64  `json:"downloadBytes"`
	UploadBytes      int64  `json:"uploadBytes"`
	TotalBytes       int64  `json:"totalBytes"`
	DownloadSpeedBps int64  `json:"downloadSpeedBps"`
	UploadSpeedBps   int64  `json:"uploadSpeedBps"`
	TotalSpeedBps    int64  `json:"totalSpeedBps"`
}

type usageSeries struct {
	UUID        string  `json:"uuid"`
	Name        string  `json:"name"`
	Color       string  `json:"color"`
	CountryCode string  `json:"countryCode"`
	Total       int64   `json:"total"`
	Data        []int64 `json:"data"`
}

type topNode struct {
	UUID        string `json:"uuid"`
	Name        string `json:"name"`
	Color       string `json:"color"`
	CountryCode string `json:"countryCode"`
	Total       int64  `json:"total"`
}

type topUser struct {
	Color    string `json:"color"`
	Username string `json:"username"`
	Total    int64  `json:"total"`
}

type legacyNodeUserUsage struct {
	Date     string `json:"date"`
	NodeUUID string `json:"nodeUuid"`
	UserUUID string `json:"userUuid"`
	Username string `json:"username"`
	Total    int64  `json:"total"`
}

type legacyUserUsage struct {
	Date        string `json:"date"`
	UserUUID    string `json:"userUuid"`
	NodeUUID    string `json:"nodeUuid"`
	NodeName    string `json:"nodeName"`
	CountryCode string `json:"countryCode"`
	Total       int64  `json:"total"`
}

func NodesHandler(manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.Header().Set("Allow", http.MethodGet)
			shared.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}

		path := strings.TrimPrefix(r.URL.Path, "/api/bandwidth-stats/nodes")
		path = strings.Trim(path, "/")
		switch {
		case path == "":
			handleGetNodesUsage(w, r, manager, cfg)
		case path == "realtime":
			handleGetNodesRealtimeUsage(w, r, manager, cfg)
		case strings.HasSuffix(path, "/users/legacy"):
			nodeUUID := strings.TrimSuffix(path, "/users/legacy")
			handleGetLegacyNodeUsersUsage(w, r, manager, cfg, nodeUUID)
		case strings.HasSuffix(path, "/users"):
			nodeUUID := strings.TrimSuffix(path, "/users")
			handleGetNodeUsersUsage(w, r, manager, cfg, nodeUUID)
		default:
			shared.WriteJSONError(w, http.StatusNotFound, "not found")
		}
	}
}

func UsersHandler(manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.Header().Set("Allow", http.MethodGet)
			shared.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}

		path := strings.TrimPrefix(r.URL.Path, "/api/bandwidth-stats/users")
		path = strings.Trim(path, "/")
		if path == "" {
			shared.WriteJSONError(w, http.StatusNotFound, "not found")
			return
		}

		if strings.HasSuffix(path, "/legacy") {
			userUUID := strings.TrimSuffix(path, "/legacy")
			handleGetLegacyUserUsage(w, r, manager, cfg, userUUID)
			return
		}
		handleGetUserUsage(w, r, manager, cfg, path)
	}
}

func handleGetNodesRealtimeUsage(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) {
	items := make([]nodeRealtimeUsage, 0)
	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		rows, err := db.QueryContext(r.Context(), `
WITH nodes_latest_updates AS (
	SELECT
		node_uuid,
		SUM(download_bytes) AS current_download_bytes,
		SUM(upload_bytes) AS current_upload_bytes,
		SUM(total_bytes) AS current_total_bytes,
		MAX(updated_at) AS latest_update_time
	FROM nodes_usage_history
	WHERE created_at = date_trunc('hour', NOW())
	GROUP BY node_uuid
)
SELECT
	n.uuid AS node_uuid,
	n.name AS node_name,
	n.country_code,
	l.current_download_bytes,
	l.current_upload_bytes,
	l.current_total_bytes,
	COALESCE(CAST(l.current_download_bytes / NULLIF(EXTRACT(EPOCH FROM (l.latest_update_time - date_trunc('hour', l.latest_update_time))), 0) AS BIGINT), 0) AS download_speed_bps,
	COALESCE(CAST(l.current_upload_bytes / NULLIF(EXTRACT(EPOCH FROM (l.latest_update_time - date_trunc('hour', l.latest_update_time))), 0) AS BIGINT), 0) AS upload_speed_bps,
	COALESCE(CAST(l.current_total_bytes / NULLIF(EXTRACT(EPOCH FROM (l.latest_update_time - date_trunc('hour', l.latest_update_time))), 0) AS BIGINT), 0) AS total_speed_bps
FROM nodes_latest_updates l
JOIN nodes n ON n.uuid = l.node_uuid
ORDER BY total_speed_bps DESC
		`)
		if err != nil {
			return err
		}
		defer rows.Close()
		for rows.Next() {
			var it nodeRealtimeUsage
			if scanErr := rows.Scan(
				&it.NodeUUID, &it.NodeName, &it.CountryCode,
				&it.DownloadBytes, &it.UploadBytes, &it.TotalBytes,
				&it.DownloadSpeedBps, &it.UploadSpeedBps, &it.TotalSpeedBps,
			); scanErr != nil {
				return scanErr
			}
			items = append(items, it)
		}
		return rows.Err()
	})
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to fetch nodes realtime usage", err, cfg)
		return
	}
	shared.WriteJSON(w, http.StatusOK, map[string]any{"response": items})
}

func handleGetNodesUsage(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) {
	startDate, endDate, dates, ok := parseDateRange(w, r)
	if !ok {
		return
	}
	topLimit := parsePositiveIntWithDefault(r.URL.Query().Get("topNodesLimit"), 20)

	sparkline := make([]int64, 0, len(dates))
	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		rows, err := db.QueryContext(r.Context(), `
WITH daily_traffic AS (
	SELECT DATE_TRUNC('day', created_at AT TIME ZONE 'UTC')::date AS date, SUM(total_bytes) AS bytes
	FROM nodes_usage_history
	WHERE created_at >= ? AND created_at <= ?
	GROUP BY DATE_TRUNC('day', created_at AT TIME ZONE 'UTC')
)
SELECT COALESCE(dt.bytes, 0) AS value
FROM unnest(?::date[]) WITH ORDINALITY AS d(date, ord)
LEFT JOIN daily_traffic dt ON dt.date = d.date
ORDER BY d.ord
		`, startDate, endDate, pgDateArrayLiteral(dates))
		if err != nil {
			return err
		}
		defer rows.Close()
		for rows.Next() {
			var v int64
			if scanErr := rows.Scan(&v); scanErr != nil {
				return scanErr
			}
			sparkline = append(sparkline, v)
		}
		return rows.Err()
	})
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to fetch nodes sparkline", err, cfg)
		return
	}

	series := make([]usageSeries, 0)
	err = manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		rows, err := db.QueryContext(r.Context(), `
WITH daily_usage AS (
	SELECT
		n.uuid, n.name, n.country_code,
		DATE_TRUNC('day', h.created_at)::date AS date,
		SUM(h.total_bytes) AS bytes
	FROM nodes n
	INNER JOIN nodes_usage_history h ON h.node_uuid = n.uuid
	WHERE h.created_at >= ? AND h.created_at <= ?
	GROUP BY n.uuid, n.name, n.country_code, DATE_TRUNC('day', h.created_at)
),
nodes_with_totals AS (
	SELECT uuid, name, country_code, SUM(bytes) AS total_bytes
	FROM daily_usage
	GROUP BY uuid, name, country_code
)
SELECT
	nt.uuid, nt.name, nt.country_code, nt.total_bytes,
	ARRAY_AGG(COALESCE(du.bytes, 0) ORDER BY d.ord) AS data
FROM nodes_with_totals nt
CROSS JOIN unnest(?::date[]) WITH ORDINALITY AS d(date, ord)
LEFT JOIN daily_usage du ON du.uuid = nt.uuid AND du.date = d.date
GROUP BY nt.uuid, nt.name, nt.country_code, nt.total_bytes
ORDER BY nt.total_bytes DESC
		`, startDate, endDate, pgDateArrayLiteral(dates))
		if err != nil {
			return err
		}
		defer rows.Close()
		for rows.Next() {
			var s usageSeries
			var dataRaw string
			if scanErr := rows.Scan(&s.UUID, &s.Name, &s.CountryCode, &s.Total, &dataRaw); scanErr != nil {
				return scanErr
			}
			s.Color = colorFromUUID(s.UUID)
			s.Data = parsePgBigintArray(dataRaw)
			series = append(series, s)
		}
		return rows.Err()
	})
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to fetch nodes usage", err, cfg)
		return
	}

	topNodes := make([]topNode, 0)
	err = manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		rows, err := db.QueryContext(r.Context(), `
SELECT n.uuid, n.name, n.country_code, COALESCE(SUM(h.total_bytes), 0) AS total
FROM nodes n
INNER JOIN nodes_usage_history h ON h.node_uuid = n.uuid
WHERE h.created_at >= ? AND h.created_at <= ?
GROUP BY n.uuid, n.name, n.country_code
ORDER BY total DESC
LIMIT ?
		`, startDate, endDate, topLimit)
		if err != nil {
			return err
		}
		defer rows.Close()
		for rows.Next() {
			var t topNode
			if scanErr := rows.Scan(&t.UUID, &t.Name, &t.CountryCode, &t.Total); scanErr != nil {
				return scanErr
			}
			t.Color = colorFromUUID(t.UUID)
			topNodes = append(topNodes, t)
		}
		return rows.Err()
	})
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to fetch top nodes", err, cfg)
		return
	}

	shared.WriteJSON(w, http.StatusOK, map[string]any{
		"response": map[string]any{
			"categories":    dates,
			"sparklineData": sparkline,
			"topNodes":      topNodes,
			"series":        series,
		},
	})
}

func handleGetNodeUsersUsage(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig, nodeUUID string) {
	startDate, endDate, dates, ok := parseDateRange(w, r)
	if !ok {
		return
	}
	topLimit := parsePositiveIntWithDefault(r.URL.Query().Get("topUsersLimit"), 100)

	var nodeID int64
	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		return db.QueryRowContext(r.Context(), `SELECT id FROM nodes WHERE uuid = ?`, nodeUUID).Scan(&nodeID)
	})
	if err != nil {
		if err == sql.ErrNoRows {
			shared.WriteJSONError(w, http.StatusNotFound, "node not found")
			return
		}
		shared.SendError(w, http.StatusInternalServerError, "failed to fetch node", err, cfg)
		return
	}

	sparkline := make([]int64, 0, len(dates))
	err = manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		rows, err := db.QueryContext(r.Context(), `
WITH daily_traffic AS (
	SELECT created_at::date AS date, SUM(total_bytes) AS bytes
	FROM nodes_user_usage_history
	WHERE node_id = ? AND created_at >= ?::date AND created_at <= ?::date
	GROUP BY created_at
)
SELECT COALESCE(dt.bytes, 0) AS value
FROM unnest(?::date[]) WITH ORDINALITY AS d(date, ord)
LEFT JOIN daily_traffic dt ON dt.date = d.date::date
ORDER BY d.ord
		`, nodeID, startDate, endDate, pgDateArrayLiteral(dates))
		if err != nil {
			return err
		}
		defer rows.Close()
		for rows.Next() {
			var v int64
			if scanErr := rows.Scan(&v); scanErr != nil {
				return scanErr
			}
			sparkline = append(sparkline, v)
		}
		return rows.Err()
	})
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to fetch node users sparkline", err, cfg)
		return
	}

	topUsers := make([]topUser, 0)
	err = manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		rows, err := db.QueryContext(r.Context(), `
SELECT u.uuid, u.username, COALESCE(SUM(nuh.total_bytes), 0) AS total
FROM users u
INNER JOIN nodes_user_usage_history nuh ON nuh.user_id = u.t_id
WHERE nuh.node_id = ? AND nuh.created_at >= ? AND nuh.created_at <= ?
GROUP BY u.uuid, u.username
ORDER BY total DESC
LIMIT ?
		`, nodeID, startDate, endDate, topLimit)
		if err != nil {
			return err
		}
		defer rows.Close()
		for rows.Next() {
			var userUUID, username string
			var total int64
			if scanErr := rows.Scan(&userUUID, &username, &total); scanErr != nil {
				return scanErr
			}
			topUsers = append(topUsers, topUser{
				Color:    colorFromUUID(userUUID),
				Username: username,
				Total:    total,
			})
		}
		return rows.Err()
	})
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to fetch top users", err, cfg)
		return
	}

	shared.WriteJSON(w, http.StatusOK, map[string]any{
		"response": map[string]any{
			"categories":    dates,
			"sparklineData": sparkline,
			"topUsers":      topUsers,
		},
	})
}

func handleGetUserUsage(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig, userUUID string) {
	startDate, endDate, dates, ok := parseDateRange(w, r)
	if !ok {
		return
	}
	topLimit := parsePositiveIntWithDefault(r.URL.Query().Get("topNodesLimit"), 20)

	var userID int64
	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		return db.QueryRowContext(r.Context(), `SELECT t_id FROM users WHERE uuid = ?`, userUUID).Scan(&userID)
	})
	if err != nil {
		if err == sql.ErrNoRows {
			shared.WriteJSONError(w, http.StatusNotFound, "user not found")
			return
		}
		shared.SendError(w, http.StatusInternalServerError, "failed to fetch user", err, cfg)
		return
	}

	sparkline := make([]int64, 0, len(dates))
	err = manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		rows, err := db.QueryContext(r.Context(), `
WITH daily_traffic AS (
	SELECT created_at::date AS date, SUM(total_bytes) AS bytes
	FROM nodes_user_usage_history
	WHERE user_id = ? AND created_at >= ?::date AND created_at <= ?::date
	GROUP BY created_at
)
SELECT COALESCE(dt.bytes, 0) AS value
FROM unnest(?::date[]) WITH ORDINALITY AS d(date, ord)
LEFT JOIN daily_traffic dt ON dt.date = d.date::date
ORDER BY d.ord
		`, userID, startDate, endDate, pgDateArrayLiteral(dates))
		if err != nil {
			return err
		}
		defer rows.Close()
		for rows.Next() {
			var v int64
			if scanErr := rows.Scan(&v); scanErr != nil {
				return scanErr
			}
			sparkline = append(sparkline, v)
		}
		return rows.Err()
	})
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to fetch user sparkline", err, cfg)
		return
	}

	series := make([]usageSeries, 0)
	err = manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		rows, err := db.QueryContext(r.Context(), `
WITH daily_usage AS (
	SELECT
		n.uuid, n.name, n.country_code,
		nuh.created_at::date AS date,
		SUM(nuh.total_bytes) AS bytes
	FROM nodes n
	INNER JOIN nodes_user_usage_history nuh ON nuh.node_id = n.id
	WHERE nuh.user_id = ? AND nuh.created_at >= ?::date AND nuh.created_at <= ?::date
	GROUP BY n.uuid, n.name, n.country_code, nuh.created_at
),
nodes_with_totals AS (
	SELECT uuid, name, country_code, SUM(bytes) AS total_bytes
	FROM daily_usage
	GROUP BY uuid, name, country_code
)
SELECT
	nt.uuid, nt.name, nt.country_code, nt.total_bytes,
	ARRAY_AGG(COALESCE(du.bytes, 0) ORDER BY d.ord) AS data
FROM nodes_with_totals nt
CROSS JOIN unnest(?::date[]) WITH ORDINALITY AS d(date, ord)
LEFT JOIN daily_usage du ON du.uuid = nt.uuid AND du.date = d.date::date
GROUP BY nt.uuid, nt.name, nt.country_code, nt.total_bytes
ORDER BY nt.total_bytes DESC
		`, userID, startDate, endDate, pgDateArrayLiteral(dates))
		if err != nil {
			return err
		}
		defer rows.Close()
		for rows.Next() {
			var s usageSeries
			var dataRaw string
			if scanErr := rows.Scan(&s.UUID, &s.Name, &s.CountryCode, &s.Total, &dataRaw); scanErr != nil {
				return scanErr
			}
			s.Color = colorFromUUID(s.UUID)
			s.Data = parsePgBigintArray(dataRaw)
			series = append(series, s)
		}
		return rows.Err()
	})
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to fetch user nodes series", err, cfg)
		return
	}

	topNodes := make([]topNode, 0)
	err = manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		rows, err := db.QueryContext(r.Context(), `
SELECT n.uuid, n.name, n.country_code, COALESCE(SUM(nuh.total_bytes), 0) AS total
FROM nodes n
INNER JOIN nodes_user_usage_history nuh ON nuh.node_id = n.id
WHERE nuh.user_id = ? AND nuh.created_at >= ? AND nuh.created_at <= ?
GROUP BY n.uuid, n.name, n.country_code
ORDER BY total DESC
LIMIT ?
		`, userID, startDate, endDate, topLimit)
		if err != nil {
			return err
		}
		defer rows.Close()
		for rows.Next() {
			var t topNode
			if scanErr := rows.Scan(&t.UUID, &t.Name, &t.CountryCode, &t.Total); scanErr != nil {
				return scanErr
			}
			t.Color = colorFromUUID(t.UUID)
			topNodes = append(topNodes, t)
		}
		return rows.Err()
	})
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to fetch user top nodes", err, cfg)
		return
	}

	shared.WriteJSON(w, http.StatusOK, map[string]any{
		"response": map[string]any{
			"categories":    dates,
			"sparklineData": sparkline,
			"topNodes":      topNodes,
			"series":        series,
		},
	})
}

func handleGetLegacyNodeUsersUsage(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig, nodeUUID string) {
	startDate, endDate, _, ok := parseDateRange(w, r)
	if !ok {
		return
	}
	items := make([]legacyNodeUserUsage, 0)
	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		rows, err := db.QueryContext(r.Context(), `
SELECT
	DATE(h.created_at) AS date,
	n.uuid AS node_uuid,
	u.uuid AS user_uuid,
	u.username,
	COALESCE(SUM(h.total_bytes), 0) AS total
FROM nodes_user_usage_history h
JOIN users u ON h.user_id = u.t_id
JOIN nodes n ON h.node_id = n.id
WHERE n.uuid = ? AND h.created_at >= ? AND h.created_at <= ?
GROUP BY DATE(h.created_at), n.uuid, u.uuid, u.username
ORDER BY DATE(h.created_at) ASC
		`, nodeUUID, startDate, endDate)
		if err != nil {
			return err
		}
		defer rows.Close()
		for rows.Next() {
			var item legacyNodeUserUsage
			if scanErr := rows.Scan(&item.Date, &item.NodeUUID, &item.UserUUID, &item.Username, &item.Total); scanErr != nil {
				return scanErr
			}
			items = append(items, item)
		}
		return rows.Err()
	})
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to fetch legacy node users usage", err, cfg)
		return
	}
	shared.WriteJSON(w, http.StatusOK, map[string]any{"response": items})
}

func handleGetLegacyUserUsage(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig, userUUID string) {
	startDate, endDate, _, ok := parseDateRange(w, r)
	if !ok {
		return
	}
	items := make([]legacyUserUsage, 0)
	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		rows, err := db.QueryContext(r.Context(), `
SELECT
	DATE(h.created_at) AS date,
	u.uuid AS user_uuid,
	n.uuid AS node_uuid,
	n.name AS node_name,
	n.country_code,
	COALESCE(SUM(h.total_bytes), 0) AS total
FROM nodes_user_usage_history h
JOIN nodes n ON h.node_id = n.id
JOIN users u ON h.user_id = u.t_id
WHERE u.uuid = ? AND h.created_at >= ? AND h.created_at <= ?
GROUP BY DATE(h.created_at), u.uuid, n.uuid, n.name, n.country_code
ORDER BY DATE(h.created_at) ASC
		`, userUUID, startDate, endDate)
		if err != nil {
			return err
		}
		defer rows.Close()
		for rows.Next() {
			var item legacyUserUsage
			if scanErr := rows.Scan(&item.Date, &item.UserUUID, &item.NodeUUID, &item.NodeName, &item.CountryCode, &item.Total); scanErr != nil {
				return scanErr
			}
			items = append(items, item)
		}
		return rows.Err()
	})
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to fetch legacy user usage", err, cfg)
		return
	}
	shared.WriteJSON(w, http.StatusOK, map[string]any{"response": items})
}

func parseDateRange(w http.ResponseWriter, r *http.Request) (time.Time, time.Time, []string, bool) {
	start := strings.TrimSpace(r.URL.Query().Get("start"))
	end := strings.TrimSpace(r.URL.Query().Get("end"))
	if start == "" || end == "" {
		shared.WriteJSONError(w, http.StatusBadRequest, "start and end are required")
		return time.Time{}, time.Time{}, nil, false
	}

	startDate, err := time.Parse("2006-01-02", start)
	if err != nil {
		shared.WriteJSONError(w, http.StatusBadRequest, "invalid start date")
		return time.Time{}, time.Time{}, nil, false
	}
	endDate, err := time.Parse("2006-01-02", end)
	if err != nil {
		shared.WriteJSONError(w, http.StatusBadRequest, "invalid end date")
		return time.Time{}, time.Time{}, nil, false
	}
	if endDate.Before(startDate) {
		shared.WriteJSONError(w, http.StatusBadRequest, "end date must be >= start date")
		return time.Time{}, time.Time{}, nil, false
	}

	startDate = startDate.UTC().Truncate(24 * time.Hour)
	endDate = endDate.UTC().Truncate(24 * time.Hour).Add(24*time.Hour - time.Nanosecond)
	dates := dateRange(startDate, endDate)
	return startDate, endDate, dates, true
}

func dateRange(start, end time.Time) []string {
	out := make([]string, 0, int(end.Sub(start)/(24*time.Hour))+1)
	cur := time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, time.UTC)
	last := time.Date(end.Year(), end.Month(), end.Day(), 0, 0, 0, 0, time.UTC)
	for !cur.After(last) {
		out = append(out, cur.Format("2006-01-02"))
		cur = cur.Add(24 * time.Hour)
	}
	return out
}

func pgDateArrayLiteral(dates []string) string {
	return "{" + strings.Join(dates, ",") + "}"
}

func parsePositiveIntWithDefault(raw string, fallback int) int {
	v, err := strconv.Atoi(raw)
	if err != nil || v < 1 {
		return fallback
	}
	return v
}

func parsePgBigintArray(v string) []int64 {
	raw := strings.TrimSpace(v)
	raw = strings.TrimPrefix(raw, "{")
	raw = strings.TrimSuffix(raw, "}")
	if raw == "" {
		return []int64{}
	}
	parts := strings.Split(raw, ",")
	out := make([]int64, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(strings.Trim(p, `"`))
		n, err := strconv.ParseInt(p, 10, 64)
		if err != nil {
			out = append(out, 0)
			continue
		}
		out = append(out, n)
	}
	return out
}

func colorFromUUID(id string) string {
	hasher := fnv.New32a()
	_, _ = hasher.Write([]byte(id))
	return palette[hasher.Sum32()%uint32(len(palette))]
}
