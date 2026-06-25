package subscription

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"exodus/internal/config"
	dbmanager "exodus/internal/db/manager"
	"exodus/internal/httpapi/middleware"
	"exodus/internal/httpapi/shared"
)

// SubscriptionsHandler handles GET /api/subscriptions
func SubscriptionsHandler(manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) http.HandlerFunc {
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

		settings, err := loadSubscriptionSettings(ctx, manager, cfg)
		if err != nil {
			shared.SendError(w, http.StatusInternalServerError, "subscription settings not found", err, cfg)
			return
		}

		users, total, err := getUsersWithPagination(ctx, manager, start, size)
		if err != nil {
			shared.SendError(w, http.StatusInternalServerError, "failed to fetch users", err, cfg)
			return
		}

		subscriptions := []SubscriptionInfoResponse{}
		for _, user := range users {
			hosts, err := getHostsForUser(ctx, manager, user)
			if err != nil {
				continue
			}
			subscriptions = append(subscriptions, buildSubscriptionInfoResponse(user, settings, hosts))
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"response": map[string]interface{}{
				"total":         total,
				"subscriptions": subscriptions,
			},
		})
	}
}

// SubscriptionsByPathHandler handles /api/subscriptions/* routes.
func SubscriptionsByPathHandler(manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/api/subscriptions/")
		path = strings.Trim(path, "/")
		if path == "" {
			http.NotFound(w, r)
			return
		}

		parts := strings.Split(path, "/")
		switch parts[0] {
		case "by-username":
			if len(parts) < 2 {
				http.NotFound(w, r)
				return
			}
			handleSubscriptionByUsername(w, r, manager, cfg, parts[1])
			return
		case "by-uuid":
			if len(parts) < 2 {
				http.NotFound(w, r)
				return
			}
			handleSubscriptionByUUID(w, r, manager, cfg, parts[1])
			return
		case "by-short-uuid":
			if len(parts) < 2 {
				http.NotFound(w, r)
				return
			}
			if len(parts) >= 3 && parts[2] == "raw" {
				handleRawSubscriptionByShortUUID(w, r, manager, cfg, parts[1])
				return
			}
			handleSubscriptionByShortUUID(w, r, manager, cfg, parts[1])
			return
		case "connection-keys":
			if len(parts) < 2 {
				http.NotFound(w, r)
				return
			}
			handleConnectionKeysByUUID(w, r, manager, cfg, parts[1])
			return
		case "subpage-config":
			if len(parts) < 2 {
				http.NotFound(w, r)
				return
			}
			handleSubpageConfigByShortUUID(w, r, manager, cfg, parts[1])
			return
		default:
			http.NotFound(w, r)
			return
		}
	}
}

func handleSubscriptionByUsername(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig, username string) {
	ctx := r.Context()
	settings, err := loadSubscriptionSettings(ctx, manager, cfg)
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "subscription settings not found", err, cfg)
		return
	}
	user, err := getSubscriptionUserByUsername(ctx, manager, username)
	if err != nil {
		shared.SendError(w, http.StatusNotFound, "user not found", err, cfg)
		return
	}
	hosts, err := getHostsForUser(ctx, manager, user)
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to fetch hosts", err, cfg)
		return
	}
	response := buildSubscriptionInfoResponse(user, settings, hosts)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(map[string]interface{}{"response": response})
}

func handleSubscriptionByUUID(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig, uuid string) {
	ctx := r.Context()
	settings, err := loadSubscriptionSettings(ctx, manager, cfg)
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "subscription settings not found", err, cfg)
		return
	}
	user, err := getSubscriptionUserByUUID(ctx, manager, uuid)
	if err != nil {
		shared.SendError(w, http.StatusNotFound, "user not found", err, cfg)
		return
	}
	hosts, err := getHostsForUser(ctx, manager, user)
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to fetch hosts", err, cfg)
		return
	}
	response := buildSubscriptionInfoResponse(user, settings, hosts)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(map[string]interface{}{"response": response})
}

func handleSubscriptionByShortUUID(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig, shortUUID string) {
	ctx := r.Context()
	settings, err := loadSubscriptionSettings(ctx, manager, cfg)
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "subscription settings not found", err, cfg)
		return
	}
	user, err := getSubscriptionUserByShortUUID(ctx, manager, shortUUID)
	if err != nil {
		shared.SendError(w, http.StatusNotFound, "user not found", err, cfg)
		return
	}
	hosts, err := getHostsForUser(ctx, manager, user)
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to fetch hosts", err, cfg)
		return
	}
	response := buildSubscriptionInfoResponse(user, settings, hosts)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(map[string]interface{}{"response": response})
}

func handleRawSubscriptionByShortUUID(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig, shortUUID string) {
	ctx := r.Context()
	withDisabledHosts := strings.EqualFold(r.URL.Query().Get("withDisabledHosts"), "true")
	settings, err := loadSubscriptionSettings(ctx, manager, cfg)
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "subscription settings not found", err, cfg)
		return
	}
	user, err := getSubscriptionUserByShortUUID(ctx, manager, shortUUID)
	if err != nil {
		shared.SendError(w, http.StatusNotFound, "user not found", err, cfg)
		return
	}

	hosts, err := getHostsForUser(ctx, manager, user)
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to fetch hosts", err, cfg)
		return
	}

	hwidHeaders := extractHwidHeaders(r)
	requestIP := middleware.GetClientIP(r, cfg)
	if hwidHeaders == nil && !settings.HwidSettings.Enabled {
		hwidHeaders = extractSyntheticHwidHeaders(r, user.UUID, requestIP)
	}
	isHapp := strings.HasPrefix(strings.ToLower(r.Header.Get("User-Agent")), "happ/")
	headers := buildSubscriptionHeaders(user, settings, isHapp)

	isHwidLimited := false
	if settings.HwidSettings.Enabled {
		allowed, maxReached, notSupported := checkHwidDeviceLimit(ctx, manager, user, hwidHeaders, settings.HwidSettings)
		if !allowed {
			headers["x-hwid-limit"] = "true"
			if notSupported {
				headers["x-hwid-not-supported"] = "true"
			}
			if maxReached {
				headers["x-hwid-max-devices-reached"] = "true"
			}
			if !notSupported && !maxReached {
				headers["x-hwid-active"] = "true"
			}
			isHwidLimited = true
		}
	} else {
		if hwidHeaders != nil {
			_ = enqueueOrUpsertHwidUserDevice(ctx, manager, user.UUID, *hwidHeaders)
		}
	}

	updateSubscriptionRequest(ctx, manager, user.UUID, r.Header.Get("User-Agent"), requestIP)

	rawHosts := []RawHost{}
	if !isHwidLimited {
		for _, host := range hosts {
			if !withDisabledHosts && host.IsDisabled {
				continue
			}
			rawHosts = append(rawHosts, buildRawHost(host))
		}
	}

	response := RawSubscriptionResponse{
		User: map[string]interface{}{
			"uuid":      user.UUID,
			"shortUuid": user.ShortUUID,
			"username":  user.Username,
		},
		ConvertedUserInfo: map[string]interface{}{
			"daysLeft":            int(time.Until(user.ExpireAt).Hours() / 24),
			"trafficUsed":         shared.FormatBytes(user.UsedTrafficBytes),
			"trafficLimit":        shared.FormatBytes(user.TrafficLimitBytes),
			"lifetimeTrafficUsed": shared.FormatBytes(user.LifetimeUsedBytes),
			"isHwidLimited":       isHwidLimited,
		},
		RawHosts: rawHosts,
		Headers:  headers,
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(map[string]interface{}{"response": response})
}

func handleConnectionKeysByUUID(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig, uuid string) {
	ctx := r.Context()
	user, err := getSubscriptionUserByUUID(ctx, manager, uuid)
	if err != nil {
		shared.SendError(w, http.StatusNotFound, "user not found", err, cfg)
		return
	}

	hosts, err := getHostsForUser(ctx, manager, user)
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to fetch hosts", err, cfg)
		return
	}

	enabled := []string{}
	disabled := []string{}
	hidden := []string{}

	for _, host := range hosts {
		link, _ := buildHostLink(host, user)
		if link == "" {
			continue
		}
		if host.IsHidden && !host.IsDisabled {
			hidden = append(hidden, link)
		} else if host.IsDisabled && !host.IsHidden {
			disabled = append(disabled, link)
		} else if !host.IsDisabled && !host.IsHidden {
			enabled = append(enabled, link)
		}
	}

	response := map[string]interface{}{
		"enabledKeys":  enabled,
		"disabledKeys": disabled,
		"hiddenKeys":   hidden,
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(map[string]interface{}{"response": response})
}

type subpageConfigRequest struct {
	RequestHeaders map[string]string `json:"requestHeaders"`
}

func handleSubpageConfigByShortUUID(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig, shortUUID string) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	var req subpageConfigRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		shared.SendError(w, http.StatusBadRequest, "invalid JSON", err, cfg)
		return
	}

	ctx := r.Context()
	subpageConfigUUID, webpageAllowed, err := getSubpageConfigForUser(ctx, manager, cfg, shortUUID, req.RequestHeaders)
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to fetch subpage config", err, cfg)
		return
	}

	response := map[string]interface{}{
		"subpageConfigUuid": subpageConfigUUID,
		"webpageAllowed":    webpageAllowed,
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(map[string]interface{}{"response": response})
}
