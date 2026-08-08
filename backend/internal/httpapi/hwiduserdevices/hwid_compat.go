package hwiduserdevices

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"exodus/internal/config"
	"exodus/internal/httpapi/shared"
	"exodus/internal/notifications"

	"github.com/google/uuid"
)

type hwidCompatDevice struct {
	HWID        string  `json:"hwid"`
	UserUUID    string  `json:"userUuid"`
	UserID      int64   `json:"userId"`
	Platform    *string `json:"platform"`
	OSVersion   *string `json:"osVersion"`
	DeviceModel *string `json:"deviceModel"`
	UserAgent   *string `json:"userAgent"`
	RequestIP   *string `json:"requestIp"`
	CreatedAt   string  `json:"createdAt"`
	UpdatedAt   string  `json:"updatedAt"`
}

type deleteAllUserHWIDDevicesRequest struct {
	UserID   *int64 `json:"userId,omitempty"`
	UserUUID string `json:"userUuid,omitempty"`
}

func (r *deleteAllUserHWIDDevicesRequest) UserIdentifier() string {
	if r.UserID != nil && *r.UserID > 0 {
		return strconv.FormatInt(*r.UserID, 10)
	}
	return r.UserUUID
}

type createUserHWIDDeviceRequest struct {
	HWID        string  `json:"hwid"`
	UserID      *int64  `json:"userId,omitempty"`
	UserUUID    string  `json:"userUuid,omitempty"`
	Platform    *string `json:"platform"`
	OSVersion   *string `json:"osVersion"`
	DeviceModel *string `json:"deviceModel"`
	UserAgent   *string `json:"userAgent"`
	RequestIP   *string `json:"requestIp"`
}

func (r *createUserHWIDDeviceRequest) UserIdentifier() string {
	if r.UserID != nil && *r.UserID > 0 {
		return strconv.FormatInt(*r.UserID, 10)
	}
	return r.UserUUID
}

type deleteUserHWIDDeviceRequest struct {
	HWID     string `json:"hwid"`
	UserID   *int64 `json:"userId,omitempty"`
	UserUUID string `json:"userUuid,omitempty"`
}

func (r *deleteUserHWIDDeviceRequest) UserIdentifier() string {
	if r.UserID != nil && *r.UserID > 0 {
		return strconv.FormatInt(*r.UserID, 10)
	}
	return r.UserUUID
}

type tableFilter struct {
	ID    string `json:"id"`
	Value any    `json:"value"`
}

type tableSorting struct {
	ID   string `json:"id"`
	Desc bool   `json:"desc"`
}

func HWIDCompatDevicesHandler(db *sql.DB, cfg *config.BackendConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := strings.Trim(strings.TrimPrefix(r.URL.Path, "/api/hwid/devices"), "/")

		if path == "delete" {
			if r.Method != http.MethodPost {
				shared.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
				return
			}
			handleHWIDCompatDeleteUserDevice(w, r, db, cfg)
			return
		}

		if r.Method == http.MethodPost && path == "" {
			handleHWIDCompatCreateUserDevice(w, r, db, cfg)
			return
		}

		if r.Method != http.MethodGet {
			shared.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}

		if path != "" {
			handleHWIDCompatGetUserDevices(w, r, db, cfg, path)
			return
		}

		start, size := parseStartSize(r, 0, 25, 1000)
		devices, total, err := getHWIDCompatDevices(r.Context(), db, start, size, r)
		if err != nil {
			shared.SendError(w, http.StatusInternalServerError, "failed to fetch hwid devices", err, cfg)
			return
		}

		shared.WriteJSON(w, http.StatusOK, map[string]any{
			"response": map[string]any{
				"devices": devices,
				"total":   total,
			},
		})
	}
}

func handleHWIDCompatGetUserDevices(w http.ResponseWriter, r *http.Request, db *sql.DB, cfg *config.BackendConfig, userUUID string) {
	devices, err := getHWIDCompatDevicesByUserUUID(r.Context(), db, userUUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			shared.SendError(w, http.StatusNotFound, "user not found", nil, cfg)
			return
		}
		shared.SendError(w, http.StatusInternalServerError, "failed to fetch user hwid devices", err, cfg)
		return
	}

	shared.WriteJSON(w, http.StatusOK, map[string]any{
		"response": map[string]any{
			"devices": devices,
			"total":   len(devices),
		},
	})
}

func handleHWIDCompatCreateUserDevice(w http.ResponseWriter, r *http.Request, db *sql.DB, cfg *config.BackendConfig) {
	var req createUserHWIDDeviceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		shared.SendError(w, http.StatusBadRequest, "invalid JSON", err, cfg)
		return
	}
	userUUID := strings.TrimSpace(req.UserUUID)
	hwid := strings.TrimSpace(req.HWID)
	if hwid == "" {
		shared.SendError(w, http.StatusBadRequest, "hwid is required", nil, cfg)
		return
	}

	var userID int64
	var fetchErr error
	if idNum, parseErr := strconv.ParseInt(userUUID, 10, 64); parseErr == nil {
		fetchErr = db.QueryRowContext(r.Context(), `SELECT id FROM users WHERE id = $1 OR uuid::text = $2 OR short_uuid = $2 OR username = $2`, idNum, userUUID).Scan(&userID)
	} else {
		fetchErr = db.QueryRowContext(r.Context(), `SELECT id FROM users WHERE uuid::text = $1 OR short_uuid = $1 OR username = $1`, userUUID).Scan(&userID)
	}
	if fetchErr != nil {
		if errors.Is(fetchErr, sql.ErrNoRows) {
			shared.SendError(w, http.StatusNotFound, "user not found", nil, cfg)
			return
		}
		shared.SendError(w, http.StatusInternalServerError, "failed to fetch user", fetchErr, cfg)
		return
	}

	var deviceExists bool
	if err := db.QueryRowContext(r.Context(), `SELECT EXISTS(SELECT 1 FROM hwid_user_devices WHERE hwid = $1 AND user_id = $2)`, hwid, userID).Scan(&deviceExists); err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to check device", err, cfg)
		return
	}
	if deviceExists {
		shared.SendError(w, http.StatusConflict, "hwid already registered for this user", nil, cfg)
		return
	}

	_, err := db.ExecContext(r.Context(), `
		INSERT INTO hwid_user_devices (hwid, user_id, platform, os_version, device_model, user_agent, request_ip, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	`, hwid, userID, req.Platform, req.OSVersion, req.DeviceModel, req.UserAgent, req.RequestIP)
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to create user hwid device", err, cfg)
		return
	}

	devices, err := getHWIDCompatDevicesByUserUUID(r.Context(), db, userUUID)
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to fetch user hwid devices", err, cfg)
		return
	}
	for _, device := range devices {
		if device.HWID == hwid && device.UserUUID == userUUID {
			emitHWIDNotification(r.Context(), cfg, notifications.EventUserHWIDDeviceAdded, hwidCompatNotificationData(device), nil)
			break
		}
	}
	shared.WriteJSON(w, http.StatusOK, map[string]any{
		"response": map[string]any{
			"devices": devices,
			"total":   len(devices),
		},
	})
}

func handleHWIDCompatDeleteUserDevice(w http.ResponseWriter, r *http.Request, db *sql.DB, cfg *config.BackendConfig) {
	var req deleteUserHWIDDeviceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		shared.SendError(w, http.StatusBadRequest, "invalid JSON", err, cfg)
		return
	}
	userUUID := strings.TrimSpace(req.UserUUID)
	hwid := strings.TrimSpace(req.HWID)
	if hwid == "" {
		shared.SendError(w, http.StatusBadRequest, "hwid is required", nil, cfg)
		return
	}

	var userID int64
	var fetchErr error
	if idNum, parseErr := strconv.ParseInt(userUUID, 10, 64); parseErr == nil {
		fetchErr = db.QueryRowContext(r.Context(), `SELECT id FROM users WHERE id = $1 OR uuid::text = $2 OR short_uuid = $2 OR username = $2`, idNum, userUUID).Scan(&userID)
	} else {
		fetchErr = db.QueryRowContext(r.Context(), `SELECT id FROM users WHERE uuid::text = $1 OR short_uuid = $1 OR username = $1`, userUUID).Scan(&userID)
	}
	if fetchErr != nil {
		if errors.Is(fetchErr, sql.ErrNoRows) {
			shared.SendError(w, http.StatusNotFound, "user not found", nil, cfg)
			return
		}
		shared.SendError(w, http.StatusInternalServerError, "failed to fetch user", fetchErr, cfg)
		return
	}

	row := db.QueryRowContext(r.Context(), `
		SELECT h.hwid, u.uuid, h.user_id, h.platform, h.os_version, h.device_model, h.user_agent, h.request_ip, h.created_at, h.updated_at
		FROM hwid_user_devices h
		JOIN users u ON u.id = h.user_id
		WHERE h.hwid = $1 AND h.user_id = $2
		LIMIT 1
	`, hwid, userID)
	deletedDevice, err := scanHWIDCompatDevice(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			shared.SendError(w, http.StatusNotFound, "hwid device not found", nil, cfg)
			return
		}
		shared.SendError(w, http.StatusInternalServerError, "failed to find hwid device", err, cfg)
		return
	}

	result, err := db.ExecContext(r.Context(), `DELETE FROM hwid_user_devices WHERE hwid = $1 AND user_id = $2`, hwid, userID)
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to delete user hwid device", err, cfg)
		return
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected == 0 {
		shared.SendError(w, http.StatusNotFound, "hwid device not found", nil, cfg)
		return
	}

	devices, err := getHWIDCompatDevicesByUserUUID(r.Context(), db, userUUID)
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to fetch user hwid devices", err, cfg)
		return
	}
	emitHWIDNotification(r.Context(), cfg, notifications.EventUserHWIDDeviceDeleted, hwidCompatNotificationData(deletedDevice), nil)
	shared.WriteJSON(w, http.StatusOK, map[string]any{
		"response": map[string]any{
			"devices": devices,
			"total":   len(devices),
		},
	})
}

func HWIDCompatDeleteAllUserDevicesHandler(db *sql.DB, cfg *config.BackendConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			shared.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}

		var req deleteAllUserHWIDDevicesRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			shared.SendError(w, http.StatusBadRequest, "invalid JSON", err, cfg)
			return
		}
		userUUID := strings.TrimSpace(req.UserUUID)
		if _, err := uuid.Parse(userUUID); err != nil {
			shared.SendError(w, http.StatusBadRequest, "invalid userUuid format", nil, cfg)
			return
		}

		var userID int64
		if err := db.QueryRowContext(r.Context(), `SELECT id FROM users WHERE uuid = $1`, userUUID).Scan(&userID); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				shared.SendError(w, http.StatusNotFound, "user not found", nil, cfg)
				return
			}
			shared.SendError(w, http.StatusInternalServerError, "failed to fetch user", err, cfg)
			return
		}

		result, err := db.ExecContext(r.Context(), `DELETE FROM hwid_user_devices WHERE user_id = $1`, userID)
		if err != nil {
			shared.SendError(w, http.StatusInternalServerError, "failed to delete user hwid devices", err, cfg)
			return
		}
		deletedCount, _ := result.RowsAffected()

		if deletedCount > 0 {
			emitHWIDNotification(r.Context(), cfg, notifications.EventUserHWIDDeviceDeleted, map[string]any{
				"userUuid":    userUUID,
				"deletedAll":  true,
				"deletedRows": deletedCount,
			}, map[string]any{"bulk": true})
		}
		shared.WriteJSON(w, http.StatusOK, map[string]any{
			"response": map[string]any{
				"devices": []hwidCompatDevice{},
				"total":   0,
			},
		})
	}
}

func HWIDCompatStatsHandler(db *sql.DB, cfg *config.BackendConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			shared.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}

		var totalDevices, uniqueDevices, uniqueUsers int
		row := db.QueryRowContext(r.Context(), `
			SELECT COUNT(*), COUNT(DISTINCT hwid), COUNT(DISTINCT user_id)
			FROM hwid_user_devices
		`)
		if err := row.Scan(&totalDevices, &uniqueDevices, &uniqueUsers); err != nil {
			shared.SendError(w, http.StatusInternalServerError, "failed to fetch hwid stats", err, cfg)
			return
		}

		byPlatform := make([]map[string]any, 0)
		rows, err := db.QueryContext(r.Context(), `
			SELECT COALESCE(NULLIF(platform, ''), 'Unknown') AS platform, COUNT(*) AS count
			FROM hwid_user_devices
			GROUP BY platform
			ORDER BY count DESC
		`)
		if err != nil {
			shared.SendError(w, http.StatusInternalServerError, "failed to fetch hwid platform stats", err, cfg)
			return
		}
		defer rows.Close()

		for rows.Next() {
			var platform string
			var count int
			if scanErr := rows.Scan(&platform, &count); scanErr != nil {
				shared.SendError(w, http.StatusInternalServerError, "failed to scan hwid platform stats", scanErr, cfg)
				return
			}
			byPlatform = append(byPlatform, map[string]any{"platform": platform, "count": count, "byApp": []map[string]any{}})
		}
		if err := rows.Err(); err != nil {
			shared.SendError(w, http.StatusInternalServerError, "failed to fetch hwid platform stats", err, cfg)
			return
		}

		byAppByPlatform := map[string][]map[string]any{}
		rows2, err := db.QueryContext(r.Context(), `
			SELECT
				COALESCE(NULLIF(platform, ''), 'Unknown') AS platform,
				COALESCE(NULLIF(SPLIT_PART(user_agent, '/', 1), ''), 'Unknown') AS app,
				COUNT(*) AS count
			FROM hwid_user_devices
			GROUP BY platform, app
			ORDER BY platform ASC, count DESC
		`)
		if err != nil {
			shared.SendError(w, http.StatusInternalServerError, "failed to fetch hwid app stats", err, cfg)
			return
		}
		defer rows2.Close()

		for rows2.Next() {
			var platform string
			var app string
			var count int
			if scanErr := rows2.Scan(&platform, &app, &count); scanErr != nil {
				shared.SendError(w, http.StatusInternalServerError, "failed to scan hwid app stats", scanErr, cfg)
				return
			}
			if strings.HasPrefix(app, "https:") {
				continue
			}
			byAppByPlatform[platform] = append(byAppByPlatform[platform], map[string]any{"app": app, "count": count})
		}
		if err := rows2.Err(); err != nil {
			shared.SendError(w, http.StatusInternalServerError, "failed to fetch hwid app stats", err, cfg)
			return
		}

		for _, platform := range byPlatform {
			name, _ := platform["platform"].(string)
			if byApp, ok := byAppByPlatform[name]; ok {
				platform["byApp"] = byApp
			}
		}

		avg := 0.0
		if uniqueUsers > 0 {
			avg = float64(totalDevices) / float64(uniqueUsers)
		}

		shared.WriteJSON(w, http.StatusOK, map[string]any{
			"response": map[string]any{
				"byPlatform": byPlatform,
				"stats": map[string]any{
					"totalUniqueDevices":        uniqueDevices,
					"totalHwidDevices":          totalDevices,
					"averageHwidDevicesPerUser": avg,
				},
			},
		})
	}
}

func HWIDCompatTopUsersHandler(db *sql.DB, cfg *config.BackendConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			shared.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}

		start, size := parseStartSize(r, 0, 5, 100)
		type item struct {
			UserUUID     string `json:"userUuid"`
			ID           int64  `json:"id"`
			Username     string `json:"username"`
			DevicesCount int    `json:"devicesCount"`
		}
		var total int
		if err := db.QueryRowContext(r.Context(), `SELECT COUNT(DISTINCT user_id) FROM hwid_user_devices`).Scan(&total); err != nil {
			shared.SendError(w, http.StatusInternalServerError, "failed to fetch top users count", err, cfg)
			return
		}

		rows, err := db.QueryContext(r.Context(), `
			SELECT u.uuid, u.id, u.username, COUNT(*) AS devices_count
			FROM hwid_user_devices h
			JOIN users u ON u.id = h.user_id
			GROUP BY u.uuid, u.id, u.username
			ORDER BY devices_count DESC, u.id ASC
			OFFSET $1 LIMIT $2
		`, start, size)
		if err != nil {
			shared.SendError(w, http.StatusInternalServerError, "failed to fetch top users by hwid devices", err, cfg)
			return
		}
		defer rows.Close()

		users := make([]item, 0)
		for rows.Next() {
			var v item
			if scanErr := rows.Scan(&v.UserUUID, &v.ID, &v.Username, &v.DevicesCount); scanErr != nil {
				shared.SendError(w, http.StatusInternalServerError, "failed to scan top users by hwid devices", scanErr, cfg)
				return
			}
			users = append(users, v)
		}
		if err := rows.Err(); err != nil {
			shared.SendError(w, http.StatusInternalServerError, "failed to fetch top users by hwid devices", err, cfg)
			return
		}

		shared.WriteJSON(w, http.StatusOK, map[string]any{
			"response": map[string]any{
				"users": users,
				"total": total,
			},
		})
	}
}

func getHWIDCompatDevices(ctx context.Context, db *sql.DB, start, size int, r *http.Request) ([]hwidCompatDevice, int, error) {
	columns := map[string]string{
		"hwid":        "h.hwid",
		"userUuid":    "u.uuid",
		"userId":      "u.id",
		"platform":    "h.platform",
		"osVersion":   "h.os_version",
		"deviceModel": "h.device_model",
		"userAgent":   "h.user_agent",
		"requestIp":   "h.request_ip",
		"createdAt":   "h.created_at",
		"updatedAt":   "h.updated_at",
	}
	whereSQL, whereArgs := buildTableWhereClause(r, columns)
	orderSQL := buildTableOrderClause(r, columns, "h.created_at DESC")

	var total int
	countQuery := `SELECT COUNT(*) FROM hwid_user_devices h JOIN users u ON u.id = h.user_id` + whereSQL
	if err := db.QueryRowContext(ctx, countQuery, whereArgs...).Scan(&total); err != nil {
		return nil, 0, err
	}

	args := append(append([]any{}, whereArgs...), start, size)
	query := fmt.Sprintf(`
		SELECT h.hwid, u.uuid, h.user_id, h.platform, h.os_version, h.device_model, h.user_agent, h.request_ip, h.created_at, h.updated_at
		FROM hwid_user_devices h
		JOIN users u ON u.id = h.user_id
	`+whereSQL+orderSQL+` OFFSET $%d LIMIT $%d`, len(whereArgs)+1, len(whereArgs)+2)

	rows, err := db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	items := make([]hwidCompatDevice, 0)
	for rows.Next() {
		item, scanErr := scanHWIDCompatDevice(rows)
		if scanErr != nil {
			return nil, 0, scanErr
		}
		items = append(items, item)
	}
	return items, total, rows.Err()
}

func getHWIDCompatDevicesByUserUUID(ctx context.Context, db *sql.DB, userIdentifier string) ([]hwidCompatDevice, error) {
	var userTID int64
	var fetchErr error
	if idNum, parseErr := strconv.ParseInt(userIdentifier, 10, 64); parseErr == nil {
		fetchErr = db.QueryRowContext(ctx, `SELECT id FROM users WHERE id = $1 OR uuid::text = $2 OR short_uuid = $2 OR username = $2`, idNum, userIdentifier).Scan(&userTID)
	} else {
		fetchErr = db.QueryRowContext(ctx, `SELECT id FROM users WHERE uuid::text = $1 OR short_uuid = $1 OR username = $1`, userIdentifier).Scan(&userTID)
	}
	if fetchErr != nil {
		return nil, fetchErr
	}

	rows, err := db.QueryContext(ctx, `
		SELECT h.hwid, u.uuid, h.user_id, h.platform, h.os_version, h.device_model, h.user_agent, h.request_ip, h.created_at, h.updated_at
		FROM hwid_user_devices h
		JOIN users u ON u.id = h.user_id
		WHERE h.user_id = $1
		ORDER BY h.created_at DESC
	`, userTID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]hwidCompatDevice, 0)
	for rows.Next() {
		item, scanErr := scanHWIDCompatDevice(rows)
		if scanErr != nil {
			return nil, scanErr
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

type hwidCompatScanner interface {
	Scan(dest ...any) error
}

func scanHWIDCompatDevice(scanner hwidCompatScanner) (hwidCompatDevice, error) {
	var item hwidCompatDevice
	var createdAt, updatedAt sql.NullTime
	if err := scanner.Scan(&item.HWID, &item.UserUUID, &item.UserID, &item.Platform, &item.OSVersion, &item.DeviceModel, &item.UserAgent, &item.RequestIP, &createdAt, &updatedAt); err != nil {
		return item, err
	}
	if createdAt.Valid {
		item.CreatedAt = createdAt.Time.UTC().Format("2006-01-02T15:04:05.000Z")
	}
	if updatedAt.Valid {
		item.UpdatedAt = updatedAt.Time.UTC().Format("2006-01-02T15:04:05.000Z")
	}
	return item, nil
}

func buildTableWhereClause(r *http.Request, columns map[string]string) (string, []any) {
	filters := parseTableFilters(r.URL.Query().Get("filters"))
	if len(filters) == 0 {
		return "", nil
	}

	parts := make([]string, 0, len(filters))
	args := make([]any, 0, len(filters))
	idx := 1
	for _, filter := range filters {
		column, ok := columns[filter.ID]
		if !ok {
			continue
		}
		value := tableFilterValue(filter.Value)
		if value == "" {
			continue
		}
		parts = append(parts, fmt.Sprintf("LOWER(COALESCE(%s::text, '')) LIKE LOWER($%d)", column, idx))
		args = append(args, "%"+value+"%")
		idx++
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

func parseStartSize(r *http.Request, defaultStart, defaultSize, maxSize int) (int, int) {
	start := defaultStart
	size := defaultSize
	if raw := strings.TrimSpace(r.URL.Query().Get("start")); raw != "" {
		if v, err := strconv.Atoi(raw); err == nil && v >= 0 {
			start = v
		}
	}
	if raw := strings.TrimSpace(r.URL.Query().Get("size")); raw != "" {
		if v, err := strconv.Atoi(raw); err == nil && v > 0 {
			size = v
		}
	}
	if size > maxSize {
		size = maxSize
	}
	return start, size
}
