package subscription

import (
	"context"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"maps"
	"net/http"
	"strings"

	"exodus/internal/config"
	"exodus/internal/httpapi/middleware"
	"exodus/internal/httpapi/shared"
	"exodus/internal/httpapi/subscriptionresponserules"
	"exodus/internal/logger"
)

func SubscriptionPublicHandler(db, backgroundDB *sql.DB, cfg *config.BackendConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
			return
		}

		path := strings.TrimPrefix(r.URL.Path, "/api/sub/")
		path = strings.Trim(path, "/")
		if path == "" {
			http.NotFound(w, r)
			return
		}

		parts := strings.Split(path, "/")
		if len(parts) == 0 {
			http.NotFound(w, r)
			return
		}

		if len(parts) == 2 && parts[1] == "info" {
			handlePublicSubscriptionInfo(w, r, db, backgroundDB, cfg, parts[0])
			return
		}

		if len(parts) >= 4 && parts[0] == "outline" {
			handlePublicOutlineSubscription(w, r, db, backgroundDB, cfg, parts)
			return
		}

		if len(parts) == 2 {
			handlePublicSubscription(w, r, db, backgroundDB, cfg, parts[0], parts[1])
			return
		}

		if len(parts) == 1 {
			handlePublicSubscription(w, r, db, backgroundDB, cfg, parts[0], "")
			return
		}

		http.NotFound(w, r)
	}
}

func handlePublicSubscriptionInfo(w http.ResponseWriter, r *http.Request, db, backgroundDB *sql.DB, cfg *config.BackendConfig, shortUUID string) {
	ctx := r.Context()
	renderService := NewRenderService(db, backgroundDB, cfg)

	user, err := getSubscriptionUserByShortUUID(ctx, db, shortUUID)
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

func handlePublicOutlineSubscription(w http.ResponseWriter, r *http.Request, db, backgroundDB *sql.DB, cfg *config.BackendConfig, parts []string) {
	ctx := r.Context()
	log := cfg.Logger.RoleService(logger.RoleAPI, logger.ServiceHTTP)

	shortUUID := parts[1]
	reqType := parts[2]
	encodedTag := parts[3]

	tagBytes, err := base64.RawURLEncoding.DecodeString(encodedTag)
	if err != nil {
		tagBytes, err = base64.StdEncoding.DecodeString(encodedTag)
		if err != nil {
			shared.SendError(w, http.StatusBadRequest, "invalid encoded tag", err, cfg)
			return
		}
	}
	targetTag := string(tagBytes)

	user, err := getSubscriptionUserByShortUUID(ctx, db, shortUUID)
	if err != nil {
		if errorsIsNoRows(err) {
			shared.SendError(w, http.StatusNotFound, "user not found", nil, cfg)
			return
		}
		shared.SendError(w, http.StatusInternalServerError, "failed to fetch user", err, cfg)
		return
	}

	hosts, err := getHostsForUser(ctx, db, user)
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to fetch hosts", err, cfg)
		return
	}

	var targetHost *SubscriptionHost
	for _, h := range hosts {
		if h.InboundTag != nil && *h.InboundTag == targetTag {
			targetHost = &h
			break
		}
	}
	if targetHost == nil {
		shared.SendError(w, http.StatusNotFound, "host for tag not found", nil, cfg)
		return
	}

	settings, err := loadSubscriptionSettings(ctx, db, cfg)
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to load settings", err, cfg)
		return
	}

	renderService := NewRenderService(db, backgroundDB, cfg)

	if reqType == "ss" || reqType == "shadowsocks" {
		links, err := NewXrayGenerator(cfg).GenerateLinks(user, []SubscriptionHost{*targetHost}, settings)
		if err != nil || len(links) == 0 {
			shared.SendError(w, http.StatusInternalServerError, "failed to generate shadowsocks link", err, cfg)
			return
		}
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		_, _ = w.Write([]byte(links[0]))
		return
	}

	content, contentType, headers, err := renderService.RenderUserSubscription(
		ctx, user, r.UserAgent(), reqType, middleware.GetClientIP(r, cfg), ExtractHwidHeaders(r),
	)
	if err != nil {
		log.Warn("Failed to render outline subscription", "short_uuid", shortUUID, "error", err)
		shared.SendError(w, http.StatusBadRequest, err.Error(), nil, cfg)
		return
	}

	for k, v := range headers {
		w.Header().Set(k, v)
	}
	if contentType != "" {
		w.Header().Set("Content-Type", contentType)
	}
	_, _ = w.Write(content)
}

func handlePublicSubscription(w http.ResponseWriter, r *http.Request, db, backgroundDB *sql.DB, cfg *config.BackendConfig, shortUUID string, clientType string) {
	ctx := r.Context()
	log := cfg.Logger.RoleService(logger.RoleAPI, logger.ServiceHTTP)

	user, err := getSubscriptionUserByShortUUID(ctx, db, shortUUID)
	if err != nil {
		if errorsIsNoRows(err) {
			shared.SendError(w, http.StatusNotFound, "user not found", nil, cfg)
			return
		}
		shared.SendError(w, http.StatusInternalServerError, "failed to fetch user", err, cfg)
		return
	}

	renderService := NewRenderService(db, backgroundDB, cfg)

	requestHeaders := make(map[string]string)
	for k, v := range r.Header {
		if len(v) > 0 {
			requestHeaders[k] = v[0]
		}
	}

	subpageConfigUUID, webpageAllowed, _ := getSubpageConfigForUser(ctx, db, cfg, shortUUID, requestHeaders)

	if webpageAllowed {
		_ = subpageConfigUUID
	}

	content, contentType, headers, err := renderService.RenderUserSubscription(
		ctx, user, r.UserAgent(), clientType, middleware.GetClientIP(r, cfg), ExtractHwidHeaders(r),
	)
	if err != nil {
		if err == ErrHwidLimitExceeded {
			shared.SendError(w, http.StatusForbidden, "HWID limit exceeded", nil, cfg)
			return
		}
		if err == ErrHwidRequired {
			shared.SendError(w, http.StatusBadRequest, "HWID header required", nil, cfg)
			return
		}
		if err == ErrUserDisabled {
			shared.SendError(w, http.StatusForbidden, "user is disabled or expired", nil, cfg)
			return
		}
		if err == ErrNoHosts {
			shared.SendError(w, http.StatusNotFound, "no active hosts available", nil, cfg)
			return
		}
		log.Warn("Failed to render subscription", "short_uuid", shortUUID, "error", err)
		shared.SendError(w, http.StatusInternalServerError, "failed to render subscription", err, cfg)
		return
	}

	for k, v := range headers {
		w.Header().Set(k, v)
	}
	if contentType != "" {
		w.Header().Set("Content-Type", contentType)
	}
	_, _ = w.Write(content)
}

func ExtractHwidHeaders(r *http.Request) *HwidHeaders {
	hwid := r.Header.Get("X-Hwid")
	if hwid == "" {
		hwid = r.Header.Get("Hwid")
	}
	if hwid == "" {
		return nil
	}

	platform := r.Header.Get("X-Hwid-Platform")
	if platform == "" {
		platform = r.Header.Get("Hwid-Platform")
	}
	osVersion := r.Header.Get("X-Hwid-Os-Version")
	if osVersion == "" {
		osVersion = r.Header.Get("Hwid-Os-Version")
	}
	deviceModel := r.Header.Get("X-Hwid-Device-Model")
	if deviceModel == "" {
		deviceModel = r.Header.Get("Hwid-Device-Model")
	}

	var pPtr, osPtr, modelPtr *string
	if platform != "" {
		pPtr = &platform
	}
	if osVersion != "" {
		osPtr = &osVersion
	}
	if deviceModel != "" {
		modelPtr = &deviceModel
	}

	userAgent := r.UserAgent()
	var uaPtr *string
	if userAgent != "" {
		uaPtr = &userAgent
	}

	return &HwidHeaders{
		Hwid:        hwid,
		Platform:    pPtr,
		OsVersion:   osPtr,
		DeviceModel: modelPtr,
		UserAgent:   uaPtr,
		RequestIP:   nil,
	}
}

func (s *RenderService) GetSubscriptionUserByShortUUID(ctx context.Context, shortUUID string) (SubscriptionUser, error) {
	return getSubscriptionUserByShortUUID(ctx, s.db, shortUUID)
}

func (s *RenderService) GetHostsForUser(ctx context.Context, user SubscriptionUser) ([]SubscriptionHost, error) {
	return getHostsForUser(ctx, s.db, user)
}

func (s *RenderService) LoadSubscriptionSettings(ctx context.Context) (SubscriptionSettingsParsed, error) {
	return loadSubscriptionSettings(ctx, s.db, s.cfg)
}

func (s *RenderService) LoadExternalSquadOverrides(ctx context.Context, squadUUID string) (*ExternalSquadOverrides, error) {
	return loadExternalSquadOverrides(ctx, s.db, squadUUID, s.cfg)
}

func (s *RenderService) MergeSettings(base SubscriptionSettingsParsed, squad *ExternalSquadOverrides) SubscriptionSettingsParsed {
	return applyExternalSquadOverrides(base, squad)
}

func (s *RenderService) ApplyHostOverrides(hosts []SubscriptionHost, overrides map[string]HostOverride) []SubscriptionHost {
	return applyHostOverrides(hosts, overrides)
}

func (s *RenderService) CheckHwidDeviceLimit(ctx context.Context, user SubscriptionUser, hwid *HwidHeaders, settings HwidSettings) (bool, bool, bool) {
	return checkHwidDeviceLimit(ctx, s.db, user, hwid, settings)
}

func (s *RenderService) UpdateSubscriptionRequest(ctx context.Context, userUUID string, userID int64, userAgent, requestIP string) {
	updateSubscriptionRequest(ctx, s.backgroundDB, userUUID, userID, userAgent, requestIP)
}

func (s *RenderService) GetSubscriptionTemplate(ctx context.Context, templateType string) ([]byte, error) {
	return getSubscriptionTemplate(ctx, s.db, templateType)
}

func (s *RenderService) GetSubscriptionTemplateByName(ctx context.Context, name string) (string, []byte, error) {
	return getSubscriptionTemplateByName(ctx, s.db, name)
}

func (s *RenderService) BuildResponseHeaders(user SubscriptionUser, settings SubscriptionSettingsParsed, contentType string) map[string]string {
	return buildResponseHeaders(user, settings, contentType)
}

func (s *RenderService) MatchResponseRules(rules *subscriptionresponserules.Config, header http.Header) string {
	return matchResponseRules(rules, header)
}

func (s *RenderService) DetectClientType(userAgent string) string {
	return detectClientType(userAgent)
}

func (s *RenderService) BuildSubscriptionInfoResponse(user SubscriptionUser, settings SubscriptionSettingsParsed, hosts []SubscriptionHost) SubscriptionInfoResponse {
	return s.buildSubscriptionInfoResponse(user, settings, hosts)
}

func (s *RenderService) CopyMap(src map[string]string) map[string]string {
	if src == nil {
		return nil
	}
	dst := make(map[string]string, len(src))
	maps.Copy(dst, src)
	return dst
}

func SubpageConfigPublicHandler(db, backgroundDB *sql.DB, cfg *config.BackendConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			shared.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}

		path := strings.TrimPrefix(r.URL.Path, "/api/subscriptions/subpage-config/")
		shortUUID := strings.Trim(path, "/")
		if shortUUID == "" {
			shared.SendError(w, http.StatusBadRequest, "shortUuid is required", nil, cfg)
			return
		}

		var req struct {
			RequestHeaders map[string]string `json:"requestHeaders"`
		}
		_ = json.NewDecoder(r.Body).Decode(&req)
		if req.RequestHeaders == nil {
			req.RequestHeaders = make(map[string]string)
		}

		subpageConfigUUID, webpageAllowed, err := getSubpageConfigForUser(r.Context(), db, cfg, shortUUID, req.RequestHeaders)
		if err != nil {
			if errorsIsNoRows(err) {
				shared.SendError(w, http.StatusNotFound, "user not found", nil, cfg)
				return
			}
			shared.SendError(w, http.StatusInternalServerError, "failed to get subpage config", err, cfg)
			return
		}

		var uuidPtr *string
		if subpageConfigUUID != "" {
			uuidPtr = &subpageConfigUUID
		}

		shared.WriteJSON(w, http.StatusOK, map[string]any{
			"response": map[string]any{
				"subpageConfigUuid": uuidPtr,
				"webpageAllowed":    webpageAllowed,
			},
		})
	}
}
