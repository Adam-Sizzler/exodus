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
	dbmanager "exodus/internal/db/manager"
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
	UserUUID string `json:"userUuid"`
}

type createUserHWIDDeviceRequest struct {
	HWID        string  `json:"hwid"`
	UserUUID    string  `json:"userUuid"`
	Platform    *string `json:"platform"`
	OSVersion   *string `json:"osVersion"`
	DeviceModel *string `json:"deviceModel"`
	UserAgent   *string `json:"userAgent"`
	RequestIP   *string `json:"requestIp"`
}

type deleteUserHWIDDeviceRequest struct {
	HWID     string `json:"hwid"`
	UserUUID string `json:"userUuid"`
}

type tableFilter struct {
	ID    string `json:"id"`
	Value any    `json:"value"`
}

type tableSorting struct {
	ID   string `json:"id"`
	Desc bool   `json:"desc"`
}

func HWIDCompatDevicesHandler(manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := strings.Trim(strings.TrimPrefix(r.URL.Path, "/api/hwid/devices"), "/")

		if path == "delete" {
			if r.Method != http.MethodPost {
				shared.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
				return
			}
			handleHWIDCompatDeleteUserDevice(w, r, manager, cfg)
			return
		}

		if r.Method == http.MethodPost && path == "" {
			handleHWIDCompatCreateUserDevice(w, r, manager, cfg)
			return
		}

		if r.Method != http.MethodGet {
			shared.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}

		if path != "" {
			handleHWIDCompatGetUserDevices(w, r, manager, cfg, path)
			return
		}

		start, size := parseStartSize(r, 0, 25, 1000)
		devices, total, err := getHWIDCompatDevices(r.Context(), manager, start, size, r)
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

func handleHWIDCompatGetUserDevices(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig, userUUID string) {
	if _, err := uuid.Parse(userUUID); err != nil {
		shared.SendError(w, http.StatusBadRequest, "invalid userUuid format", nil, cfg)
		return
	}

	devices, err := getHWIDCompatDevicesByUserUUID(r.Context(), manager, userUUID)
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

func handleHWIDCompatCreateUserDevice(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) {
	var req createUserHWIDDeviceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		shared.SendError(w, http.StatusBadRequest, "invalid JSON", err, cfg)
		return
	}
	userUUID := strings.TrimSpace(req.UserUUID)
	hwid := strings.TrimSpace(req.HWID)
	if _, err := uuid.Parse(userUUID); err != nil {
		shared.SendError(w, http.StatusBadRequest, "invalid userUuid format", nil, cfg)
		return
	}
	if hwid == "" {
		shared.SendError(w, http.StatusBadRequest, "hwid is required", nil, cfg)
		return
	}

	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		var userID int64
		if err := db.QueryRowContext(r.Context(), `SELECT t_id FROM users WHERE uuid = ?`, userUUID).Scan(&userID); err != nil {
			return err
		}
		var deviceExists bool
		if err := db.QueryRowContext(r.Context(), `SELECT EXISTS(SELECT 1 FROM hwid_user_devices WHERE hwid = ? AND user_id = ?)`, hwid, userID).Scan(&deviceExists); err != nil {
			return err
		}
		if deviceExists {
			return errHWIDDeviceExists
		}
		_, err := db.ExecContext(r.Context(), `
			INSERT INTO hwid_user_devices (hwid, user_id, platform, os_version, device_model, user_agent, request_ip, created_at, updated_at)
			VALUES (?, ?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		`, hwid, userID, req.Platform, req.OSVersion, req.DeviceModel, req.UserAgent, req.RequestIP)
		return err
	})
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			shared.SendError(w, http.StatusNotFound, "user not found", nil, cfg)
		case errors.Is(err, errHWIDDeviceExists):
			shared.SendError(w, http.StatusConflict, "hwid already registered for this user", nil, cfg)
		default:
			shared.SendError(w, http.StatusInternalServerError, "failed to create user hwid device", err, cfg)
		}
		return
	}

	devices, err := getHWIDCompatDevicesByUserUUID(r.Context(), manager, userUUID)
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

func handleHWIDCompatDeleteUserDevice(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) {
	var req deleteUserHWIDDeviceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		shared.SendError(w, http.StatusBadRequest, "invalid JSON", err, cfg)
		return
	}
	userUUID := strings.TrimSpace(req.UserUUID)
	hwid := strings.TrimSpace(req.HWID)
	if _, err := uuid.Parse(userUUID); err != nil {
		shared.SendError(w, http.StatusBadRequest, "invalid userUuid format", nil, cfg)
		return
	}
	if hwid == "" {
		shared.SendError(w, http.StatusBadRequest, "hwid is required", nil, cfg)
		return
	}

	var deletedDevice hwidCompatDevice
	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		var userID int64
		if err := db.QueryRowContext(r.Context(), `SELECT t_id FROM users WHERE uuid = ?`, userUUID).Scan(&userID); err != nil {
			return err
		}
		row := db.QueryRowContext(r.Context(), `
			SELECT h.hwid, u.uuid, h.user_id, h.platform, h.os_version, h.device_model, h.user_agent, h.request_ip, h.created_at, h.updated_at
			FROM hwid_user_devices h
			JOIN users u ON u.t_id = h.user_id
			WHERE h.hwid = ? AND h.user_id = ?
			LIMIT 1
		`, hwid, userID)
		device, scanErr := scanHWIDCompatDevice(row)
		if scanErr != nil {
			if errors.Is(scanErr, sql.ErrNoRows) {
				return errHWIDDeviceNotFound
			}
			return scanErr
		}
		deletedDevice = device
		result, err := db.ExecContext(r.Context(), `DELETE FROM hwid_user_devices WHERE hwid = ? AND user_id = ?`, hwid, userID)
		if err != nil {
			return err
		}
		rowsAffected, err := result.RowsAffected()
		if err != nil {
			return err
		}
		if rowsAffected == 0 {
			return errHWIDDeviceNotFound
		}
		return nil
	})
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			shared.SendError(w, http.StatusNotFound, "user not found", nil, cfg)
		case errors.Is(err, errHWIDDeviceNotFound):
			shared.SendError(w, http.StatusNotFound, "hwid device not found", nil, cfg)
		default:
			shared.SendError(w, http.StatusInternalServerError, "failed to delete user hwid device", err, cfg)
		}
		return
	}

	devices, err := getHWIDCompatDevicesByUserUUID(r.Context(), manager, userUUID)
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

func HWIDCompatDeleteAllUserDevicesHandler(manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) http.HandlerFunc {
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

		var deletedCount int64
		err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
			var userID int64
			if err := db.QueryRowContext(r.Context(), `SELECT t_id FROM users WHERE uuid = ?`, userUUID).Scan(&userID); err != nil {
				return err
			}
			result, err := db.ExecContext(r.Context(), `DELETE FROM hwid_user_devices WHERE user_id = ?`, userID)
			if err != nil {
				return err
			}
			deletedCount, err = result.RowsAffected()
			return err
		})
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				shared.SendError(w, http.StatusNotFound, "user not found", nil, cfg)
				return
			}
			shared.SendError(w, http.StatusInternalServerError, "failed to delete user hwid devices", err, cfg)
			return
		}

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

func HWIDCompatStatsHandler(manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			shared.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}

		byPlatform := make([]map[string]any, 0)
		byAppByPlatform := map[string][]map[string]any{}
		var totalDevices int
		var uniqueDevices int
		var uniqueUsers int

		err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
			row := db.QueryRowContext(r.Context(), `
				SELECT COUNT(*), COUNT(DISTINCT hwid), COUNT(DISTINCT user_id)
				FROM hwid_user_devices
			`)
			if err := row.Scan(&totalDevices, &uniqueDevices, &uniqueUsers); err != nil {
				return err
			}

			rows, err := db.QueryContext(r.Context(), `
				SELECT COALESCE(NULLIF(platform, ''), 'Unknown') AS platform, COUNT(*) AS count
				FROM hwid_user_devices
				GROUP BY platform
				ORDER BY count DESC
			`)
			if err != nil {
				return err
			}
			defer rows.Close()
			for rows.Next() {
				var platform string
				var count int
				if scanErr := rows.Scan(&platform, &count); scanErr != nil {
					return scanErr
				}
				byPlatform = append(byPlatform, map[string]any{"platform": platform, "count": count, "byApp": []map[string]any{}})
			}
			if err := rows.Err(); err != nil {
				return err
			}

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
				return err
			}
			defer rows2.Close()
			for rows2.Next() {
				var platform string
				var app string
				var count int
				if scanErr := rows2.Scan(&platform, &app, &count); scanErr != nil {
					return scanErr
				}
				if strings.HasPrefix(app, "https:") {
					continue
				}
				byAppByPlatform[platform] = append(byAppByPlatform[platform], map[string]any{"app": app, "count": count})
			}
			if err := rows2.Err(); err != nil {
				return err
			}

			for _, platform := range byPlatform {
				name, _ := platform["platform"].(string)
				if byApp, ok := byAppByPlatform[name]; ok {
					platform["byApp"] = byApp
				}
			}
			return nil
		})
		if err != nil {
			shared.SendError(w, http.StatusInternalServerError, "failed to fetch hwid stats", err, cfg)
			return
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

func HWIDCompatTopUsersHandler(manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) http.HandlerFunc {
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
		users := make([]item, 0)
		total := 0

		err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
			if err := db.QueryRowContext(r.Context(), `SELECT COUNT(DISTINCT user_id) FROM hwid_user_devices`).Scan(&total); err != nil {
				return err
			}
			rows, err := db.QueryContext(r.Context(), `
				SELECT u.uuid, u.t_id, u.username, COUNT(*) AS devices_count
				FROM hwid_user_devices h
				JOIN users u ON u.t_id = h.user_id
				GROUP BY u.uuid, u.t_id, u.username
				ORDER BY devices_count DESC, u.t_id ASC
				OFFSET ? LIMIT ?
			`, start, size)
			if err != nil {
				return err
			}
			defer rows.Close()
			for rows.Next() {
				var v item
				if scanErr := rows.Scan(&v.UserUUID, &v.ID, &v.Username, &v.DevicesCount); scanErr != nil {
					return scanErr
				}
				users = append(users, v)
			}
			return rows.Err()
		})
		if err != nil {
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

func getHWIDCompatDevices(ctx context.Context, manager *dbmanager.DatabaseManager, start, size int, r *http.Request) ([]hwidCompatDevice, int, error) {
	items := make([]hwidCompatDevice, 0)
	total := 0
	columns := map[string]string{
		"hwid":        "h.hwid",
		"userUuid":    "u.uuid",
		"userId":      "u.t_id",
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

	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		countQuery := `SELECT COUNT(*) FROM hwid_user_devices h JOIN users u ON u.t_id = h.user_id` + whereSQL
		if err := db.QueryRowContext(ctx, countQuery, whereArgs...).Scan(&total); err != nil {
			return err
		}

		args := append(append([]any{}, whereArgs...), start, size)
		query := `
			SELECT h.hwid, u.uuid, h.user_id, h.platform, h.os_version, h.device_model, h.user_agent, h.request_ip, h.created_at, h.updated_at
			FROM hwid_user_devices h
			JOIN users u ON u.t_id = h.user_id
		` + whereSQL + orderSQL + ` OFFSET ? LIMIT ?`
		rows, err := db.QueryContext(ctx, query, args...)
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			item, scanErr := scanHWIDCompatDevice(rows)
			if scanErr != nil {
				return scanErr
			}
			items = append(items, item)
		}
		return rows.Err()
	})
	return items, total, err
}

func getHWIDCompatDevicesByUserUUID(ctx context.Context, manager *dbmanager.DatabaseManager, userUUID string) ([]hwidCompatDevice, error) {
	items := make([]hwidCompatDevice, 0)
	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		var exists bool
		if err := db.QueryRowContext(ctx, `SELECT EXISTS(SELECT 1 FROM users WHERE uuid = ?)`, userUUID).Scan(&exists); err != nil {
			return err
		}
		if !exists {
			return sql.ErrNoRows
		}
		rows, err := db.QueryContext(ctx, `
			SELECT h.hwid, u.uuid, h.user_id, h.platform, h.os_version, h.device_model, h.user_agent, h.request_ip, h.created_at, h.updated_at
			FROM hwid_user_devices h
			JOIN users u ON u.t_id = h.user_id
			WHERE h.user_id = (SELECT t_id FROM users WHERE uuid = ?)
			ORDER BY h.created_at DESC
		`, userUUID)
		if err != nil {
			return err
		}
		defer rows.Close()
		for rows.Next() {
			item, scanErr := scanHWIDCompatDevice(rows)
			if scanErr != nil {
				return scanErr
			}
			items = append(items, item)
		}
		return rows.Err()
	})
	return items, err
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
	for _, filter := range filters {
		column, ok := columns[filter.ID]
		if !ok {
			continue
		}
		value := tableFilterValue(filter.Value)
		if value == "" {
			continue
		}
		parts = append(parts, "LOWER(COALESCE("+column+"::text, '')) LIKE LOWER(?)")
		args = append(args, "%"+value+"%")
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
