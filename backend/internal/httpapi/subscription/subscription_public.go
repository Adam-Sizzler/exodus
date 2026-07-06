package subscription

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"maps"
	"net/http"
	"strings"

	"exodus/internal/config"
	dbmanager "exodus/internal/db/manager"
	"exodus/internal/httpapi/middleware"
	"exodus/internal/httpapi/shared"
	"exodus/internal/logger"
)

func SubscriptionPublicHandler(manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) http.HandlerFunc {
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

		// /api/sub/:shortUuid/info
		if len(parts) == 2 && parts[1] == "info" {
			handlePublicSubscriptionInfo(w, r, manager, cfg, parts[0])
			return
		}

		// /api/sub/outline/:shortUuid/:type/:encodedTag
		if len(parts) >= 4 && parts[0] == "outline" {
			handlePublicOutlineSubscription(w, r, manager, cfg, parts)
			return
		}

		// /api/sub/:shortUuid/:clientType
		if len(parts) == 2 {
			handlePublicSubscription(w, r, manager, cfg, parts[0], parts[1])
			return
		}

		// /api/sub/:shortUuid
		if len(parts) == 1 {
			handlePublicSubscription(w, r, manager, cfg, parts[0], "")
			return
		}

		http.NotFound(w, r)
	}
}

func handlePublicSubscriptionInfo(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig, shortUUID string) {
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

func handlePublicOutlineSubscription(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig, parts []string) {
	if len(parts) < 4 {
		http.NotFound(w, r)
		return
	}

	shortUUID := parts[1]
	typ := parts[2]
	encodedTag := parts[3]

	if typ != "ss" || encodedTag == "" {
		http.NotFound(w, r)
		return
	}

	handlePublicSubscription(w, r, manager, cfg, shortUUID, "")
}

func handlePublicSubscription(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig, shortUUID, clientType string) {
	ctx := r.Context()
	log := cfg.Logger.RoleService(logger.RoleAPI, logger.ServiceHTTP)

	var headerDump strings.Builder
	for name, values := range r.Header {
		fmt.Fprintf(&headerDump, "  %s: %s\n", name, strings.Join(values, ", "))
	}
	log.Debug("Public subscription request", "short_uuid", shortUUID, "headers", headerDump.String())

	log.Debug("Loading subscription settings")
	settings, err := loadSubscriptionSettings(ctx, manager, cfg)
	if err != nil {
		log.Warn("Failed to load subscription settings", "error", err)
		shared.SendError(w, http.StatusInternalServerError, "subscription settings not found", err, cfg)
		return
	}
	log.Debug("Subscription settings loaded", "hwid_enabled", settings.HwidSettings.Enabled)

	log.Debug("Looking up subscription user", "short_uuid", shortUUID)
	user, err := getSubscriptionUserByShortUUID(ctx, manager, shortUUID)
	if err != nil {
		log.Warn("Failed to load subscription user", "short_uuid", shortUUID, "error", err)
		shared.SendError(w, http.StatusNotFound, "user not found", err, cfg)
		return
	}
	log.Debug("Subscription user found", "uuid", user.UUID, "status", user.Status, "hwid_limit", user.HwidDeviceLimit)

	var externalOverrides *ExternalSquadOverrides
	if user.ExternalSquadUUID != nil && *user.ExternalSquadUUID != "" {
		log.Debug("Loading external squad overrides", "squad_uuid", *user.ExternalSquadUUID)
		overrides, err := loadExternalSquadOverrides(ctx, manager, *user.ExternalSquadUUID, cfg)
		if err != nil {
			log.Warn("Failed to load external squad overrides, continuing without them", "error", err)
		} else {
			externalOverrides = overrides
			settings = applyExternalSquadOverrides(settings, externalOverrides)
			log.Debug("Applied external squad overrides",
				"has_host_overrides", len(externalOverrides.HostOverrides) > 0,
				"has_templates", len(externalOverrides.Templates) > 0)
		}
	}

	headersForMatch := r.Header.Clone()
	headersForMatch.Set("x-exodus-injected-short-uuid", shortUUID)
	if clientType != "" {
		headersForMatch.Set("x-exodus-injected-client-type", clientType)
	}

	log.Debug("Matching subscription response rules", "client_type", clientType)
	matchResult := matchResponseRulesDetailed(settings.ResponseRules, headersForMatch, clientType)
	log.Debug("Subscription response rule matched", "matched", matchResult.Matched, "response_type", matchResult.ResponseType)

	if !matchResult.Matched || matchResult.ResponseType == "" {
		log.Debug("Response rule did not match, returning forbidden")
		shared.SendError(w, http.StatusForbidden, "forbidden", nil, cfg)
		return
	}

	responseType := matchResult.ResponseType
	if responseType == "" {
		responseType = defaultResponseType
	}
	log.Debug("Effective subscription response type", "response_type", responseType)

	var extraHeaders map[string]string
	var overrideTemplateName string
	ignoreServeJsonAtBaseSubscription := false

	if matchResult.MatchedRule != nil && matchResult.MatchedRule.ResponseModifications != nil {
		mods := matchResult.MatchedRule.ResponseModifications
		if len(mods.Headers) > 0 {
			if mods.ApplyHeadersToEnd {
				extraHeaders = map[string]string{}
				for _, header := range mods.Headers {
					extraHeaders[header.Key] = header.Value
				}
			} else {
				for _, header := range mods.Headers {
					w.Header().Set(header.Key, header.Value)
				}
			}
		}
		if mods.SubscriptionTemplate != nil {
			overrideTemplateName = strings.TrimSpace(*mods.SubscriptionTemplate)
		}
		if mods.IgnoreServeJsonAtBaseSubscription {
			ignoreServeJsonAtBaseSubscription = true
		}
	}

	switch responseType {
	case responseTypeBlock:
		log.Debug("Response type BLOCK, returning forbidden")
		shared.SendError(w, http.StatusForbidden, "forbidden", nil, cfg)
		return
	case responseTypeStatus404:
		log.Debug("Response type 404, returning not found")
		http.NotFound(w, r)
		return
	case responseTypeStatus451:
		log.Debug("Response type 451, returning unavailable for legal reasons")
		http.Error(w, "Unavailable For Legal Reasons", http.StatusUnavailableForLegalReasons)
		return
	case responseTypeSocketDrop:
		log.Debug("Response type SOCKET_DROP, dropping connection")
		if hj, ok := w.(http.Hijacker); ok {
			conn, _, err := hj.Hijack()
			if err == nil {
				_ = conn.Close()
				return
			}
		}
		return
	}

	if responseType == responseTypeXrayBase64 && settings.Raw.ServeJSONAtBaseSubscription && !ignoreServeJsonAtBaseSubscription {
		log.Debug("Switching response type from XRAY_BASE64 to XRAY_JSON")
		responseType = responseTypeXrayJSON
	}

	if responseType == responseTypeBrowser {
		log.Debug("Generating browser subscription info page")
		hosts, err := getHostsForUser(ctx, manager, user)
		if err != nil {
			log.Warn("Failed to fetch hosts for browser subscription", "error", err)
			shared.SendError(w, http.StatusInternalServerError, "failed to fetch hosts", err, cfg)
			return
		}
		response := buildSubscriptionInfoResponse(user, settings, hosts)
		for key, value := range extraHeaders {
			w.Header().Set(key, value)
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(w).Encode(response)
		return
	}

	log.Debug("Fetching subscription hosts")
	hosts, err := getHostsForUser(ctx, manager, user)
	if err != nil {
		log.Warn("Failed to fetch subscription hosts", "error", err)
		shared.SendError(w, http.StatusInternalServerError, "failed to fetch hosts", err, cfg)
		return
	}
	log.Debug("Subscription hosts fetched", "hosts", len(hosts))

	if externalOverrides != nil && len(externalOverrides.HostOverrides) > 0 {
		hosts = applyHostOverrides(hosts, externalOverrides.HostOverrides)
		log.Debug("Applied host overrides to subscription hosts")
	}

	shuffleHostsIfNeeded(hosts, settings)

	log.Debug("Extracting HWID headers")
	hwidHeaders := extractHwidHeaders(r)
	requestIP := middleware.GetClientIP(r, cfg)
	if hwidHeaders == nil && !settings.HwidSettings.Enabled {
		hwidHeaders = extractSyntheticHwidHeaders(r, user.UUID, requestIP)
		if hwidHeaders != nil {
			log.Debug("Synthetic HWID generated", "hwid", hwidHeaders.Hwid, "platform", hwidHeaders.Platform, "model", hwidHeaders.DeviceModel, "user_agent", hwidHeaders.UserAgent)
		}
	}
	if hwidHeaders != nil && hwidHeaders.RequestIP == nil {
		hwidHeaders.RequestIP = stringPtrIfNotEmpty(requestIP)
	}
	log.Debug("HWID headers extracted", "hwid_headers", hwidHeaders)

	isHapp := strings.HasPrefix(strings.ToLower(r.Header.Get("User-Agent")), "happ/")
	headers := buildSubscriptionHeaders(user, settings, isHapp)
	maps.Copy(headers, extraHeaders)

	if settings.HwidSettings.Enabled {
		log.Debug("Checking HWID device limit")
		allowed, maxReached, notSupported := checkHwidDeviceLimit(ctx, manager, user, hwidHeaders, settings.HwidSettings)
		log.Debug("HWID device limit checked", "allowed", allowed, "max_reached", maxReached, "not_supported", notSupported)
		if !allowed {
			headers["x-hwid-limit"] = "true"
			if notSupported {
				headers["x-hwid-not-supported"] = "true"
			}
			if maxReached {
				headers["x-hwid-max-devices-reached"] = "true"
			}
			if !notSupported && !maxReached {
				// device is known but limit is active (e.g. per-user override)
				headers["x-hwid-active"] = "true"
			}
			if settings.HwidSettings.MaxDevicesAnnounce != nil {
				headers["announce"] = "base64:" + base64EncodeSafe(*settings.HwidSettings.MaxDevicesAnnounce)
			}
			if settings.Raw.IsShowCustomRemarks {
				body := strings.Join(buildHwidRemarks(settings.CustomRemarks, maxReached, notSupported), "\n")
				writeSubscriptionResponse(w, headers, "text/plain", body)
				return
			}
			writeSubscriptionResponse(w, headers, "text/plain", "")
			return
		}
	} else {
		if hwidHeaders != nil {
			log.Debug("HWID disabled, inserting device record")
			if err := enqueueOrUpsertHwidUserDevice(ctx, manager, user.TID, *hwidHeaders); err != nil {
				log.Warn("Failed to upsert HWID user device", "error", err)
			} else {
				log.Debug("HWID device record upserted")
			}
		}
	}

	filteredHosts := filterHostsForResponseType(hosts, responseType, false)

	if settings.Raw.IsShowCustomRemarks {
		statusRemarks := buildStatusRemarks(settings.CustomRemarks, user.Status)
		if len(statusRemarks) > 0 {
			writeSubscriptionResponse(w, headers, "text/plain", strings.Join(statusRemarks, "\n"))
			return
		}
		if len(filteredHosts) == 0 && len(settings.CustomRemarks.EmptyHosts) > 0 {
			writeSubscriptionResponse(w, headers, "text/plain", strings.Join(settings.CustomRemarks.EmptyHosts, "\n"))
			return
		}
	}

	templateType := responseTypeToTemplateType(responseType)
	var templateData []byte
	var templateLoaded bool

	normalizeKey := func(k string) string {
		k = strings.ReplaceAll(k, "_", "")
		k = strings.ReplaceAll(k, "-", "")
		return strings.ToUpper(k)
	}

	if externalOverrides != nil {
		targetKey := normalizeKey(templateType)
		var templateName string
		for k, v := range externalOverrides.Templates {
			if normalizeKey(k) == targetKey {
				templateName = v
				break
			}
		}

		if templateName != "" {
			log.Debug("Using external squad template", "template_type", templateType, "template_name", templateName)
			var err error
			_, templateData, err = getSubscriptionTemplateByName(ctx, manager, templateName)
			if err == nil && len(templateData) > 0 {
				templateLoaded = true
			} else {
				log.Warn("Failed to load external squad template, falling back to global", "template_name", templateName, "error", err)
			}
		}
	}

	if !templateLoaded {
		templateData, _ = getSubscriptionTemplate(ctx, manager, templateType)
	}

	if overrideTemplateName != "" {
		log.Debug("Applying template override from response rules", "template_name", overrideTemplateName)
		if overrideType, overrideData, err := getSubscriptionTemplateByName(ctx, manager, overrideTemplateName); err == nil {
			if strings.EqualFold(normalizeKey(overrideType), normalizeKey(templateType)) && len(overrideData) > 0 {
				templateData = overrideData
				log.Debug("Template from response rule override applied successfully")
			} else {
				log.Warn("Override template type mismatch or empty", "expected", templateType, "actual", overrideType)
			}
		} else {
			log.Warn("Failed to load override template by name", "name", overrideTemplateName, "error", err)
		}
	}

	subscription, err := generateSubscriptionContent(responseType, templateData, filteredHosts, user)
	if err != nil {
		log.Warn("Failed to generate subscription content", "error", err)
		shared.SendError(w, http.StatusInternalServerError, "failed to generate subscription", err, cfg)
		return
	}

	log.Debug("Updating subscription request history")
	updateSubscriptionRequest(ctx, manager, user.UUID, user.TID, r.Header.Get("User-Agent"), requestIP)

	log.Debug("Writing subscription response", "bytes", len(subscription.Body))
	writeSubscriptionResponse(w, headers, subscription.ContentType, subscription.Body)
}

func mapClientTypeToResponseType(clientType string) string {
	switch strings.ToLower(clientType) {
	case "json", "v2ray-json":
		return responseTypeXrayJSON
	case "mihomo":
		return responseTypeMihomo
	case "stash":
		return responseTypeStash
	case "clash":
		return responseTypeClash
	case "singbox":
		return responseTypeSingbox
	case "xray-json":
		return responseTypeXrayJSON
	default:
		return responseTypeBlock
	}
}

func writeSubscriptionResponse(w http.ResponseWriter, headers map[string]string, contentType, body string) {
	for key, value := range headers {
		if value == "" {
			continue
		}
		w.Header().Set(key, value)
	}
	if contentType != "" {
		w.Header().Set("Content-Type", contentType)
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(body))
}

func base64EncodeSafe(value string) string {
	return base64.StdEncoding.EncodeToString([]byte(value))
}

func buildHwidRemarks(remarks CustomRemarks, maxReached, notSupported bool) []string {
	if notSupported && len(remarks.HWIDNotSupported) > 0 {
		return remarks.HWIDNotSupported
	}
	if maxReached && len(remarks.HWIDMaxDevicesExceeded) > 0 {
		return remarks.HWIDMaxDevicesExceeded
	}
	return []string{"Subscription limited"}
}

func buildStatusRemarks(remarks CustomRemarks, status string) []string {
	switch strings.ToUpper(status) {
	case "EXPIRED":
		return remarks.ExpiredUsers
	case "LIMITED":
		return remarks.LimitedUsers
	case "DISABLED":
		return remarks.DisabledUsers
	default:
		return nil
	}
}
