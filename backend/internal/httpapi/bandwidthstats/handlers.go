package bandwidthstats

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"exodus/internal/config"
	"exodus/internal/httpapi/shared"
)

func NodesHandler(db *sql.DB, cfg *config.BackendConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/api/bandwidth-stats/nodes")
		path = strings.Trim(path, "/")

		if path == "usage" {
			if r.Method != http.MethodPost {
				w.Header().Set("Allow", http.MethodPost)
				shared.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
				return
			}
			handlePostNodesUsage(w, r, db, cfg)
			return
		}

		if path == "users" {
			if r.Method != http.MethodPost {
				w.Header().Set("Allow", http.MethodPost)
				shared.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
				return
			}
			handleGetNodesUsersUsage(w, r, db, cfg)
			return
		}

		if r.Method != http.MethodGet {
			w.Header().Set("Allow", http.MethodGet)
			shared.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}

		switch {
		case path == "":
			handleGetNodesUsage(w, r, db, cfg)
		case path == "realtime":
			handleGetNodesRealtimeUsage(w, r, db, cfg)
		case strings.HasSuffix(path, "/users"):
			nodeUUID := strings.TrimSuffix(path, "/users")
			handleGetNodeUsersUsage(w, r, db, cfg, nodeUUID)
		default:
			shared.WriteJSONError(w, http.StatusNotFound, "not found")
		}
	}
}

func UsersHandler(db *sql.DB, cfg *config.BackendConfig) http.HandlerFunc {
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

		handleGetUserUsage(w, r, db, cfg, path)
	}
}

func handleGetNodesRealtimeUsage(w http.ResponseWriter, r *http.Request, db *sql.DB, cfg *config.BackendConfig) {
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
		shared.SendError(w, http.StatusInternalServerError, "failed to fetch nodes realtime usage", err, cfg)
		return
	}
	defer rows.Close()

	items := make([]nodeRealtimeUsage, 0)
	for rows.Next() {
		var it nodeRealtimeUsage
		if scanErr := rows.Scan(
			&it.NodeUUID, &it.NodeName, &it.CountryCode,
			&it.DownloadBytes, &it.UploadBytes, &it.TotalBytes,
			&it.DownloadSpeedBps, &it.UploadSpeedBps, &it.TotalSpeedBps,
		); scanErr != nil {
			shared.SendError(w, http.StatusInternalServerError, "failed to scan nodes realtime usage", scanErr, cfg)
			return
		}
		items = append(items, it)
	}
	if err := rows.Err(); err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to fetch nodes realtime usage", err, cfg)
		return
	}

	shared.WriteJSON(w, http.StatusOK, map[string]any{"response": items})
}

func handleGetNodesUsage(w http.ResponseWriter, r *http.Request, db *sql.DB, cfg *config.BackendConfig) {
	startDate, endDate, dates, ok := parseDateRange(w, r)
	if !ok {
		return
	}
	topLimit := parsePositiveIntWithDefault(r.URL.Query().Get("topNodesLimit"), 20)

	sparkRows, err := db.QueryContext(r.Context(), `
WITH daily_traffic AS (
	SELECT DATE_TRUNC('day', created_at AT TIME ZONE 'UTC')::date AS date, SUM(total_bytes) AS bytes
	FROM nodes_usage_history
	WHERE created_at >= $1 AND created_at <= $2
	GROUP BY DATE_TRUNC('day', created_at AT TIME ZONE 'UTC')
)
SELECT COALESCE(dt.bytes, 0) AS value
FROM unnest($3::date[]) WITH ORDINALITY AS d(date, ord)
LEFT JOIN daily_traffic dt ON dt.date = d.date
ORDER BY d.ord
	`, startDate, endDate, pgDateArrayLiteral(dates))
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to fetch nodes sparkline", err, cfg)
		return
	}
	defer sparkRows.Close()

	sparkline := make([]int64, 0, len(dates))
	for sparkRows.Next() {
		var v int64
		if scanErr := sparkRows.Scan(&v); scanErr != nil {
			shared.SendError(w, http.StatusInternalServerError, "failed to scan nodes sparkline", scanErr, cfg)
			return
		}
		sparkline = append(sparkline, v)
	}
	if err := sparkRows.Err(); err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to fetch nodes sparkline", err, cfg)
		return
	}

	seriesRows, err := db.QueryContext(r.Context(), `
WITH daily_usage AS (
	SELECT
		n.uuid, n.name, n.country_code,
		DATE_TRUNC('day', h.created_at)::date AS date,
		SUM(h.total_bytes) AS bytes
	FROM nodes n
	INNER JOIN nodes_usage_history h ON h.node_uuid = n.uuid
	WHERE h.created_at >= $1 AND h.created_at <= $2
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
CROSS JOIN unnest($3::date[]) WITH ORDINALITY AS d(date, ord)
LEFT JOIN daily_usage du ON du.uuid = nt.uuid AND du.date = d.date
GROUP BY nt.uuid, nt.name, nt.country_code, nt.total_bytes
ORDER BY nt.total_bytes DESC
	`, startDate, endDate, pgDateArrayLiteral(dates))
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to fetch nodes usage", err, cfg)
		return
	}
	defer seriesRows.Close()

	series := make([]usageSeries, 0)
	for seriesRows.Next() {
		var s usageSeries
		var dataRaw string
		if scanErr := seriesRows.Scan(&s.UUID, &s.Name, &s.CountryCode, &s.Total, &dataRaw); scanErr != nil {
			shared.SendError(w, http.StatusInternalServerError, "failed to scan nodes usage", scanErr, cfg)
			return
		}
		s.Color = colorFromUUID(s.UUID)
		s.Data = parsePgBigintArray(dataRaw)
		series = append(series, s)
	}
	if err := seriesRows.Err(); err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to fetch nodes usage", err, cfg)
		return
	}

	topRows, err := db.QueryContext(r.Context(), `
SELECT n.uuid, n.name, n.country_code, COALESCE(SUM(h.total_bytes), 0) AS total
FROM nodes n
INNER JOIN nodes_usage_history h ON h.node_uuid = n.uuid
WHERE h.created_at >= $1 AND h.created_at <= $2
GROUP BY n.uuid, n.name, n.country_code
ORDER BY total DESC
LIMIT $3
	`, startDate, endDate, topLimit)
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to fetch top nodes", err, cfg)
		return
	}
	defer topRows.Close()

	topNodes := make([]topNode, 0)
	for topRows.Next() {
		var t topNode
		if scanErr := topRows.Scan(&t.UUID, &t.Name, &t.CountryCode, &t.Total); scanErr != nil {
			shared.SendError(w, http.StatusInternalServerError, "failed to scan top nodes", scanErr, cfg)
			return
		}
		t.Color = colorFromUUID(t.UUID)
		topNodes = append(topNodes, t)
	}
	if err := topRows.Err(); err != nil {
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

func handleGetNodeUsersUsage(w http.ResponseWriter, r *http.Request, db *sql.DB, cfg *config.BackendConfig, nodeUUID string) {
	startDate, endDate, dates, ok := parseDateRange(w, r)
	if !ok {
		return
	}
	topLimit := parsePositiveIntWithDefault(r.URL.Query().Get("topUsersLimit"), 100)

	var nodeID int64
	err := db.QueryRowContext(r.Context(), `SELECT id FROM nodes WHERE uuid = $1`, nodeUUID).Scan(&nodeID)
	if err != nil {
		if err == sql.ErrNoRows {
			shared.WriteJSONError(w, http.StatusNotFound, "node not found")
			return
		}
		shared.SendError(w, http.StatusInternalServerError, "failed to fetch node", err, cfg)
		return
	}

	sparkRows, err := db.QueryContext(r.Context(), `
WITH daily_traffic AS (
	SELECT created_at::date AS date, SUM(total_bytes) AS bytes
	FROM nodes_user_usage_history
	WHERE node_id = $1 AND created_at >= $2::date AND created_at <= $3::date
	GROUP BY created_at
)
SELECT COALESCE(dt.bytes, 0) AS value
FROM unnest($4::date[]) WITH ORDINALITY AS d(date, ord)
LEFT JOIN daily_traffic dt ON dt.date = d.date::date
ORDER BY d.ord
	`, nodeID, startDate, endDate, pgDateArrayLiteral(dates))
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to fetch node users sparkline", err, cfg)
		return
	}
	defer sparkRows.Close()

	sparkline := make([]int64, 0, len(dates))
	for sparkRows.Next() {
		var v int64
		if scanErr := sparkRows.Scan(&v); scanErr != nil {
			shared.SendError(w, http.StatusInternalServerError, "failed to scan node users sparkline", scanErr, cfg)
			return
		}
		sparkline = append(sparkline, v)
	}
	if err := sparkRows.Err(); err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to fetch node users sparkline", err, cfg)
		return
	}

	topRows, err := db.QueryContext(r.Context(), `
SELECT u.uuid, u.username, COALESCE(SUM(nuh.total_bytes), 0) AS total
FROM users u
INNER JOIN nodes_user_usage_history nuh ON nuh.user_id = u.id
WHERE nuh.node_id = $1 AND nuh.created_at >= $2 AND nuh.created_at <= $3
GROUP BY u.uuid, u.username
ORDER BY total DESC
LIMIT $4
	`, nodeID, startDate, endDate, topLimit)
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to fetch top users", err, cfg)
		return
	}
	defer topRows.Close()

	topUsers := make([]topUser, 0)
	for topRows.Next() {
		var userUUID, username string
		var total int64
		if scanErr := topRows.Scan(&userUUID, &username, &total); scanErr != nil {
			shared.SendError(w, http.StatusInternalServerError, "failed to scan top users", scanErr, cfg)
			return
		}
		topUsers = append(topUsers, topUser{
			Color:    colorFromUUID(userUUID),
			Username: username,
			Total:    total,
		})
	}
	if err := topRows.Err(); err != nil {
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

func handleGetNodesUsersUsage(w http.ResponseWriter, r *http.Request, db *sql.DB, cfg *config.BackendConfig) {
	startDate, endDate, dates, ok := parseDateRange(w, r)
	if !ok {
		return
	}
	topLimit := parsePositiveIntWithDefault(r.URL.Query().Get("topUsersLimit"), 100)

	var req getNodesUsersUsageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		shared.WriteJSONError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if len(req.NodesUUIDs) == 0 {
		shared.WriteJSONError(w, http.StatusBadRequest, "nodesUuids: must be at least 1 node UUID")
		return
	}

	nodeRows, err := db.QueryContext(r.Context(), `SELECT id FROM nodes WHERE uuid = ANY($1)`, req.NodesUUIDs)
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to fetch nodes", err, cfg)
		return
	}
	defer nodeRows.Close()

	var nodeIDs []int64
	for nodeRows.Next() {
		var id int64
		if scanErr := nodeRows.Scan(&id); scanErr != nil {
			shared.SendError(w, http.StatusInternalServerError, "failed to scan node ID", scanErr, cfg)
			return
		}
		nodeIDs = append(nodeIDs, id)
	}
	if err := nodeRows.Err(); err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to fetch nodes", err, cfg)
		return
	}
	if len(nodeIDs) == 0 {
		shared.WriteJSONError(w, http.StatusNotFound, "nodes not found")
		return
	}

	sparkRows, err := db.QueryContext(r.Context(), `
WITH daily_traffic AS (
	SELECT created_at::date AS date, SUM(total_bytes) AS bytes
	FROM nodes_user_usage_history
	WHERE node_id = ANY($1) AND created_at >= $2::date AND created_at <= $3::date
	GROUP BY created_at
)
SELECT COALESCE(dt.bytes, 0) AS value
FROM unnest($4::date[]) WITH ORDINALITY AS d(date, ord)
LEFT JOIN daily_traffic dt ON dt.date = d.date::date
ORDER BY d.ord
	`, nodeIDs, startDate, endDate, pgDateArrayLiteral(dates))
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to fetch nodes users sparkline", err, cfg)
		return
	}
	defer sparkRows.Close()

	sparkline := make([]int64, 0, len(dates))
	for sparkRows.Next() {
		var v int64
		if scanErr := sparkRows.Scan(&v); scanErr != nil {
			shared.SendError(w, http.StatusInternalServerError, "failed to scan nodes users sparkline", scanErr, cfg)
			return
		}
		sparkline = append(sparkline, v)
	}
	if err := sparkRows.Err(); err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to fetch nodes users sparkline", err, cfg)
		return
	}

	topRows, err := db.QueryContext(r.Context(), `
SELECT u.uuid, u.username, COALESCE(SUM(nuh.total_bytes), 0) AS total
FROM users u
INNER JOIN nodes_user_usage_history nuh ON nuh.user_id = u.id
WHERE nuh.node_id = ANY($1) AND nuh.created_at >= $2 AND nuh.created_at <= $3
GROUP BY u.uuid, u.username
ORDER BY total DESC
LIMIT $4
	`, nodeIDs, startDate, endDate, topLimit)
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to fetch top users", err, cfg)
		return
	}
	defer topRows.Close()

	topUsers := make([]topUser, 0)
	for topRows.Next() {
		var userUUID, username string
		var total int64
		if scanErr := topRows.Scan(&userUUID, &username, &total); scanErr != nil {
			shared.SendError(w, http.StatusInternalServerError, "failed to scan top users", scanErr, cfg)
			return
		}
		topUsers = append(topUsers, topUser{
			Color:    colorFromUUID(userUUID),
			Username: username,
			Total:    total,
		})
	}
	if err := topRows.Err(); err != nil {
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

func handleGetUserUsage(w http.ResponseWriter, r *http.Request, db *sql.DB, cfg *config.BackendConfig, userUUID string) {
	startDate, endDate, dates, ok := parseDateRange(w, r)
	if !ok {
		return
	}
	topLimit := parsePositiveIntWithDefault(r.URL.Query().Get("topNodesLimit"), 20)

	var userID int64
	var err error
	if idNum, parseErr := strconv.ParseInt(userUUID, 10, 64); parseErr == nil {
		err = db.QueryRowContext(r.Context(), `SELECT id FROM users WHERE id = $1 OR uuid::text = $2 OR short_uuid = $2 OR username = $2`, idNum, userUUID).Scan(&userID)
	} else {
		err = db.QueryRowContext(r.Context(), `SELECT id FROM users WHERE uuid::text = $1 OR short_uuid = $1 OR username = $1`, userUUID).Scan(&userID)
	}
	if err != nil {
		if err == sql.ErrNoRows {
			shared.WriteJSONError(w, http.StatusNotFound, "user not found")
			return
		}
		shared.SendError(w, http.StatusInternalServerError, "failed to fetch user", err, cfg)
		return
	}

	sparkRows, err := db.QueryContext(r.Context(), `
WITH daily_traffic AS (
	SELECT created_at::date AS date, SUM(total_bytes) AS bytes
	FROM nodes_user_usage_history
	WHERE user_id = $1 AND created_at >= $2::date AND created_at <= $3::date
	GROUP BY created_at
)
SELECT COALESCE(dt.bytes, 0) AS value
FROM unnest($4::date[]) WITH ORDINALITY AS d(date, ord)
LEFT JOIN daily_traffic dt ON dt.date = d.date::date
ORDER BY d.ord
	`, userID, startDate, endDate, pgDateArrayLiteral(dates))
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to fetch user sparkline", err, cfg)
		return
	}
	defer sparkRows.Close()

	sparkline := make([]int64, 0, len(dates))
	for sparkRows.Next() {
		var v int64
		if scanErr := sparkRows.Scan(&v); scanErr != nil {
			shared.SendError(w, http.StatusInternalServerError, "failed to scan user sparkline", scanErr, cfg)
			return
		}
		sparkline = append(sparkline, v)
	}
	if err := sparkRows.Err(); err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to fetch user sparkline", err, cfg)
		return
	}

	seriesRows, err := db.QueryContext(r.Context(), `
WITH daily_usage AS (
	SELECT
		n.uuid, n.name, n.country_code,
		nuh.created_at::date AS date,
		SUM(nuh.total_bytes) AS bytes
	FROM nodes n
	INNER JOIN nodes_user_usage_history nuh ON nuh.node_id = n.id
	WHERE nuh.user_id = $1 AND nuh.created_at >= $2::date AND nuh.created_at <= $3::date
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
CROSS JOIN unnest($4::date[]) WITH ORDINALITY AS d(date, ord)
LEFT JOIN daily_usage du ON du.uuid = nt.uuid AND du.date = d.date::date
GROUP BY nt.uuid, nt.name, nt.country_code, nt.total_bytes
ORDER BY nt.total_bytes DESC
	`, userID, startDate, endDate, pgDateArrayLiteral(dates))
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to fetch user nodes series", err, cfg)
		return
	}
	defer seriesRows.Close()

	series := make([]usageSeries, 0)
	for seriesRows.Next() {
		var s usageSeries
		var dataRaw string
		if scanErr := seriesRows.Scan(&s.UUID, &s.Name, &s.CountryCode, &s.Total, &dataRaw); scanErr != nil {
			shared.SendError(w, http.StatusInternalServerError, "failed to scan user nodes series", scanErr, cfg)
			return
		}
		s.Color = colorFromUUID(s.UUID)
		s.Data = parsePgBigintArray(dataRaw)
		series = append(series, s)
	}
	if err := seriesRows.Err(); err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to fetch user nodes series", err, cfg)
		return
	}

	topRows, err := db.QueryContext(r.Context(), `
SELECT n.uuid, n.name, n.country_code, COALESCE(SUM(nuh.total_bytes), 0) AS total
FROM nodes n
INNER JOIN nodes_user_usage_history nuh ON nuh.node_id = n.id
WHERE nuh.user_id = $1 AND nuh.created_at >= $2 AND nuh.created_at <= $3
GROUP BY n.uuid, n.name, n.country_code
ORDER BY total DESC
LIMIT $4
	`, userID, startDate, endDate, topLimit)
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to fetch user top nodes", err, cfg)
		return
	}
	defer topRows.Close()

	topNodes := make([]topNode, 0)
	for topRows.Next() {
		var t topNode
		if scanErr := topRows.Scan(&t.UUID, &t.Name, &t.CountryCode, &t.Total); scanErr != nil {
			shared.SendError(w, http.StatusInternalServerError, "failed to scan user top nodes", scanErr, cfg)
			return
		}
		t.Color = colorFromUUID(t.UUID)
		topNodes = append(topNodes, t)
	}
	if err := topRows.Err(); err != nil {
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

func handlePostNodesUsage(w http.ResponseWriter, r *http.Request, db *sql.DB, cfg *config.BackendConfig) {
	startDate, endDate, _, ok := parseDateRange(w, r)
	if !ok {
		return
	}
	minTotalBytesStr := r.URL.Query().Get("minTotalBytes")
	var minTotalBytes int64 = 0
	if minTotalBytesStr != "" {
		if parsed, err := strconv.ParseInt(minTotalBytesStr, 10, 64); err == nil && parsed > 0 {
			minTotalBytes = parsed
		}
	}

	var req struct {
		NodesUUIDs []string `json:"nodesUuids"`
	}
	if r.Body != nil {
		_ = json.NewDecoder(r.Body).Decode(&req)
	}

	if len(req.NodesUUIDs) == 0 {
		shared.WriteJSON(w, http.StatusOK, map[string]any{"response": []any{}})
		return
	}

	rows, err := db.QueryContext(r.Context(), `
		SELECT u.id, u.username, COALESCE(SUM(nuh.total_bytes), 0) AS total_bytes
		FROM nodes_user_usage_history nuh
		JOIN nodes n ON nuh.node_id = n.id
		JOIN users u ON nuh.user_id = u.id
		WHERE n.uuid = ANY($1::text[]) AND nuh.created_at >= $2 AND nuh.created_at <= $3
		GROUP BY u.id, u.username
		HAVING SUM(nuh.total_bytes) >= $4
		ORDER BY total_bytes DESC
	`, pgDateArrayLiteral(req.NodesUUIDs), startDate, endDate, minTotalBytes)
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to fetch nodes usage", err, cfg)
		return
	}
	defer rows.Close()

	type nodeUserUsageItem struct {
		UserID     int64  `json:"userId"`
		Username   string `json:"username"`
		TotalBytes int64  `json:"totalBytes"`
	}

	items := make([]nodeUserUsageItem, 0)
	for rows.Next() {
		var item nodeUserUsageItem
		if scanErr := rows.Scan(&item.UserID, &item.Username, &item.TotalBytes); scanErr != nil {
			shared.SendError(w, http.StatusInternalServerError, "failed to scan node user usage item", scanErr, cfg)
			return
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to fetch node user usage items", err, cfg)
		return
	}

	shared.WriteJSON(w, http.StatusOK, map[string]any{"response": items})
}
