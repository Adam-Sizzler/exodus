package hwiduserdevices

import (
	"context"
	"database/sql"
	"net/http"
	"strconv"
	"strings"

	"v2ray-stat/backend/config"
	dbmanager "v2ray-stat/backend/db/manager"
	"v2ray-stat/backend/httpapi/shared"
)

type hwidCompatDevice struct {
	HWID        string  `json:"hwid"`
	UserUUID    string  `json:"userUuid"`
	Platform    *string `json:"platform"`
	OSVersion   *string `json:"osVersion"`
	DeviceModel *string `json:"deviceModel"`
	UserAgent   *string `json:"userAgent"`
	CreatedAt   string  `json:"createdAt"`
	UpdatedAt   string  `json:"updatedAt"`
}

func HWIDCompatDevicesHandler(manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			shared.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}

		start, size := parseStartSize(r, 0, 25, 1000)
		devices, total, err := getHWIDCompatDevices(r.Context(), manager, start, size)
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

func HWIDCompatStatsHandler(manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			shared.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}

		type kv struct {
			Name  string `json:"platform,omitempty"`
			Count int    `json:"count"`
		}
		byPlatform := make([]map[string]any, 0)
		byApp := make([]map[string]any, 0)
		var totalDevices int
		var uniqueDevices int
		var uniqueUsers int

		err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
			row := db.QueryRowContext(r.Context(), `
				SELECT COUNT(*), COUNT(DISTINCT hwid), COUNT(DISTINCT user_uuid)
				FROM hwid_user_devices
			`)
			return row.Scan(&totalDevices, &uniqueDevices, &uniqueUsers)
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
				"byApp":      byApp,
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
			if err := db.QueryRowContext(r.Context(), `SELECT COUNT(DISTINCT user_uuid) FROM hwid_user_devices`).Scan(&total); err != nil {
				return err
			}
			rows, err := db.QueryContext(r.Context(), `
				SELECT u.uuid, u.t_id, u.username, COUNT(*) AS devices_count
				FROM hwid_user_devices h
				JOIN users u ON u.uuid = h.user_uuid
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

func getHWIDCompatDevices(ctx context.Context, manager *dbmanager.DatabaseManager, start, size int) ([]hwidCompatDevice, int, error) {
	items := make([]hwidCompatDevice, 0)
	total := 0
	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		if err := db.QueryRowContext(ctx, `SELECT COUNT(*) FROM hwid_user_devices`).Scan(&total); err != nil {
			return err
		}

		rows, err := db.QueryContext(ctx, `
			SELECT h.hwid, h.user_uuid, h.platform, h.os_version, h.device_model, h.user_agent, h.created_at, h.updated_at
			FROM hwid_user_devices h
			ORDER BY h.created_at DESC
			OFFSET ? LIMIT ?
		`, start, size)
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			var item hwidCompatDevice
			var createdAt, updatedAt sql.NullTime
			if scanErr := rows.Scan(&item.HWID, &item.UserUUID, &item.Platform, &item.OSVersion, &item.DeviceModel, &item.UserAgent, &createdAt, &updatedAt); scanErr != nil {
				return scanErr
			}
			if createdAt.Valid {
				item.CreatedAt = createdAt.Time.UTC().Format("2006-01-02T15:04:05.000Z")
			}
			if updatedAt.Valid {
				item.UpdatedAt = updatedAt.Time.UTC().Format("2006-01-02T15:04:05.000Z")
			}
			items = append(items, item)
		}
		return rows.Err()
	})
	return items, total, err
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
