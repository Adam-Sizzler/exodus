package subscription

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strings"

	"cerberus/backend/config"
	dbmanager "cerberus/backend/db/manager"
	"cerberus/backend/httpapi/middleware"
	"cerberus/backend/httpapi/shared"
)

// SubscriptionPublicHandler handles public subscription endpoints under /api/sub.
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

	headersForMatch := r.Header.Clone()
	headersForMatch.Set("x-cerberus-injected-short-uuid", shortUUID)
	if clientType != "" {
		headersForMatch.Set("x-cerberus-injected-client-type", clientType)
	}

	matchResult := matchResponseRulesDetailed(settings.ResponseRules, headersForMatch, clientType)
	if !matchResult.Matched || matchResult.ResponseType == "" {
		shared.SendError(w, http.StatusForbidden, "forbidden", nil, cfg)
		return
	}

	responseType := matchResult.ResponseType
	if responseType == "" {
		responseType = defaultResponseType
	}

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
		shared.SendError(w, http.StatusForbidden, "forbidden", nil, cfg)
		return
	case responseTypeStatus404:
		http.NotFound(w, r)
		return
	case responseTypeStatus451:
		http.Error(w, "Unavailable For Legal Reasons", http.StatusUnavailableForLegalReasons)
		return
	case responseTypeSocketDrop:
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
		responseType = responseTypeXrayJSON
	}

	if responseType == responseTypeBrowser {
		hosts, err := getHostsForUser(ctx, manager, user)
		if err != nil {
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

	hosts, err := getHostsForUser(ctx, manager, user)
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to fetch hosts", err, cfg)
		return
	}

	shuffleHostsIfNeeded(hosts, settings)

	hwidHeaders := extractHwidHeaders(r)
	isHapp := strings.HasPrefix(strings.ToLower(r.Header.Get("User-Agent")), "happ/")
	headers := buildSubscriptionHeaders(user, settings, isHapp)
	for key, value := range extraHeaders {
		headers[key] = value
	}

	if settings.HwidSettings.Enabled {
		allowed, maxReached, notSupported := checkHwidDeviceLimit(ctx, manager, user, hwidHeaders, settings.HwidSettings)
		if !allowed {
			headers["x-hwid-limit"] = "true"
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
			_ = upsertHwidUserDevice(ctx, manager, user.UUID, *hwidHeaders)
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
	templateData, _ := getSubscriptionTemplate(ctx, manager, templateType)
	if overrideTemplateName != "" {
		if overrideType, overrideData, err := getSubscriptionTemplateByName(ctx, manager, overrideTemplateName); err == nil {
			if strings.EqualFold(overrideType, templateType) && len(overrideData) > 0 {
				templateData = overrideData
			}
		}
	}

	subscription, err := generateSubscriptionContent(responseType, templateData, filteredHosts, user)
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to generate subscription", err, cfg)
		return
	}

	updateSubscriptionRequest(ctx, manager, user.UUID, r.Header.Get("User-Agent"), middleware.GetClientIP(r, cfg))

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
	if value == "" {
		return ""
	}
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
