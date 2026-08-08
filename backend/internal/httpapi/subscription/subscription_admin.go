package subscription

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"exodus/internal/config"
	"exodus/internal/httpapi/shared"
)

// SubscriptionsHandler handles GET /api/subscriptions
func SubscriptionsHandler(db, backgroundDB *sql.DB, cfg *config.BackendConfig) http.HandlerFunc {
	renderService := NewRenderService(db, backgroundDB, cfg)
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
			return
		}

		start, _ := strconv.Atoi(r.URL.Query().Get("start"))
		size, _ := strconv.Atoi(r.URL.Query().Get("size"))
		if start < 0 {
			start = 0
		}
		if size <= 0 {
			size = 25
		}

		ctx := r.Context()

		settings, err := loadSubscriptionSettings(ctx, db, cfg)
		if err != nil {
			shared.SendError(w, http.StatusInternalServerError, "subscription settings not found", err, cfg)
			return
		}

		users, total, err := getUsersWithPagination(ctx, db, start, size)
		if err != nil {
			shared.SendError(w, http.StatusInternalServerError, "failed to fetch users", err, cfg)
			return
		}

		subscriptions := []SubscriptionInfoResponse{}
		for _, user := range users {
			hosts, err := getHostsForUser(ctx, db, user)
			if err != nil {
				continue
			}
			subscriptions = append(subscriptions, renderService.buildSubscriptionInfoResponse(user, settings, hosts))
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"response": map[string]interface{}{
				"subscriptions": subscriptions,
				"total":         total,
				"start":         start,
				"size":          size,
			},
		})
	}
}

// SubscriptionByUUIDHandler handles GET /api/subscriptions/:uuid and GET /api/subscriptions/connection-keys/:uuid
func SubscriptionByUUIDHandler(db, backgroundDB *sql.DB, cfg *config.BackendConfig) http.HandlerFunc {
	renderService := NewRenderService(db, backgroundDB, cfg)
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
			return
		}

		path := strings.Trim(strings.TrimPrefix(r.URL.Path, "/api/subscriptions/"), "/")
		if path == "" {
			http.NotFound(w, r)
			return
		}

		if strings.HasPrefix(path, "connection-keys/") {
			userUUID := strings.TrimPrefix(path, "connection-keys/")
			userUUID = strings.TrimSuffix(strings.TrimSpace(userUUID), "/raw")
			handleGetConnectionKeysByUUID(w, r, db, cfg, userUUID)
			return
		}

		path = strings.TrimSuffix(strings.TrimSpace(path), "/raw")

		if strings.HasPrefix(path, "by-id/") {
			path = strings.TrimPrefix(path, "by-id/")
		} else if strings.HasPrefix(path, "by-uuid/") {
			path = strings.TrimPrefix(path, "by-uuid/")
		} else if strings.HasPrefix(path, "by-short-uuid/") {
			path = strings.TrimPrefix(path, "by-short-uuid/")
		} else if strings.HasPrefix(path, "by-username/") {
			path = strings.TrimPrefix(path, "by-username/")
		}

		userUUID := path
		ctx := r.Context()

		user, err := getSubscriptionUserByUUID(ctx, db, userUUID)
		if err != nil {
			if errorsIsNoRows(err) {
				shared.SendError(w, http.StatusNotFound, "user not found", nil, cfg)
				return
			}
			shared.SendError(w, http.StatusInternalServerError, "failed to fetch user", err, cfg)
			return
		}

		settings, err := loadSubscriptionSettings(ctx, db, cfg)
		if err != nil {
			shared.SendError(w, http.StatusInternalServerError, "subscription settings not found", err, cfg)
			return
		}

		hosts, err := getHostsForUser(ctx, db, user)
		if err != nil {
			shared.SendError(w, http.StatusInternalServerError, "failed to fetch hosts", err, cfg)
			return
		}

		info := renderService.buildSubscriptionInfoResponse(user, settings, hosts)

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"response": info,
		})
	}
}

func handleGetConnectionKeysByUUID(w http.ResponseWriter, r *http.Request, db *sql.DB, cfg *config.BackendConfig, userUUID string) {
	ctx := r.Context()
	user, err := getSubscriptionUserByUUID(ctx, db, userUUID)
	if err != nil {
		if errorsIsNoRows(err) {
			shared.SendError(w, http.StatusNotFound, "user not found", nil, cfg)
			return
		}
		shared.SendError(w, http.StatusInternalServerError, "failed to fetch user", err, cfg)
		return
	}

	enabledKeys := []string{}
	hiddenKeys := []string{}
	disabledKeys := []string{}

	if user.VlessUUID != "" {
		enabledKeys = append(enabledKeys, "vless")
	}
	if user.SSPassword != "" {
		enabledKeys = append(enabledKeys, "shadowsocks")
	}
	if user.TrojanPassword != "" {
		enabledKeys = append(enabledKeys, "trojan")
	}
	if user.NaivePassword != "" {
		enabledKeys = append(enabledKeys, "naive")
	}
	if user.ShadowtlsPassword != "" {
		enabledKeys = append(enabledKeys, "shadowtls")
	}
	if user.Hysteria2Password != "" {
		enabledKeys = append(enabledKeys, "hysteria2")
	}
	if user.AnytlsPassword != "" {
		enabledKeys = append(enabledKeys, "anytls")
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"response": map[string]interface{}{
			"enabledKeys":  enabledKeys,
			"hiddenKeys":   hiddenKeys,
			"disabledKeys": disabledKeys,
		},
	})
}

func errorsIsNoRows(err error) bool {
	return err != nil && (err == sql.ErrNoRows || strings.Contains(err.Error(), "no rows in result set"))
}
