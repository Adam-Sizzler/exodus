package subscription

/*
CROSS-CUTTING RULES / НЕЯВНЫЕ ЗАВИСИМОСТИ:
1. Использование host.Path:
   Значение `host.Path` используется как для WebSocket (`ws.path`), так и для gRPC (`grpc.service_name`).
   В билдерах (Mihomo, Xray, Singbox) при проверке пути gRPC нужно использовать `host.Path`.
2. Обработка Reality:
   Для совместимости с xray/mihomo билдерами, Reality принудительно схлопывается в `"tls"`.
*/

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"sort"
	"strings"
	"time"

	"exodus/internal/config"
	dbmanager "exodus/internal/db/manager"
	"exodus/internal/httpapi/shared"
	"exodus/internal/httpapi/subscriptionresponserules"
	"exodus/internal/httpapi/subscriptionsettings"
)

type RenderService struct {
	manager *dbmanager.DatabaseManager
	cfg     *config.BackendConfig
}

func NewRenderService(manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) *RenderService {
	return &RenderService{manager: manager, cfg: cfg}
}

func applyExternalSquadOverrides(base SubscriptionSettingsParsed, overrides *ExternalSquadOverrides) SubscriptionSettingsParsed {
	if overrides == nil {
		return base
	}

	if overrides.SubscriptionSettings != nil {
		base.Raw = mergeSubscriptionSettings(base.Raw, *overrides.SubscriptionSettings)

		if strings.TrimSpace(base.Raw.ResponseRules) != "" {
			var rules subscriptionresponserules.Config
			if err := json.Unmarshal([]byte(base.Raw.ResponseRules), &rules); err == nil {
				base.ResponseRules = &rules
			}
		}

		if strings.TrimSpace(base.Raw.CustomResponseHeaders) != "" {
			merged := map[string]string{}
			if err := json.Unmarshal([]byte(base.Raw.CustomResponseHeaders), &merged); err == nil {
				base.CustomResponseHeaders = merged
			}
		}
	}

	if overrides.HwidSettings != nil {
		base.HwidSettings = *overrides.HwidSettings
	}

	if overrides.CustomRemarks != nil {
		base.CustomRemarks = *overrides.CustomRemarks
	}

	if len(overrides.ResponseHeaders) > 0 {
		base.CustomResponseHeaders = overrides.ResponseHeaders
	}

	return base
}

func mergeSubscriptionSettings(base, override subscriptionsettings.SubscriptionSettings) subscriptionsettings.SubscriptionSettings {
	result := base
	if override.ProfileTitle != "" {
		result.ProfileTitle = override.ProfileTitle
	}
	if override.SupportLink != "" {
		result.SupportLink = override.SupportLink
	}
	if override.ProfileUpdateInterval != 0 {
		result.ProfileUpdateInterval = override.ProfileUpdateInterval
	}
	if override.Address != "" {
		result.Address = override.Address
	}
	if override.Port != 0 {
		result.Port = override.Port
	}
	if override.APISchema != "" {
		result.APISchema = override.APISchema
	}
	if override.APIPath != "" {
		result.APIPath = override.APIPath
	}
	if override.HappAnnounce != "" {
		result.HappAnnounce = override.HappAnnounce
	}
	if override.HappRouting != "" {
		result.HappRouting = override.HappRouting
	}
	if override.IsProfileWebpageURLEnabled != result.IsProfileWebpageURLEnabled {
		result.IsProfileWebpageURLEnabled = override.IsProfileWebpageURLEnabled
	}
	if override.ServeJSONAtBaseSubscription != result.ServeJSONAtBaseSubscription {
		result.ServeJSONAtBaseSubscription = override.ServeJSONAtBaseSubscription
	}
	if override.IsShowCustomRemarks != result.IsShowCustomRemarks {
		result.IsShowCustomRemarks = override.IsShowCustomRemarks
	}
	if override.RandomizeHosts != result.RandomizeHosts {
		result.RandomizeHosts = override.RandomizeHosts
	}
	if override.CustomResponseHeaders != "" {
		result.CustomResponseHeaders = override.CustomResponseHeaders
	}
	if override.ResponseRules != "" {
		result.ResponseRules = override.ResponseRules
	}
	if override.HWIDSettings != "" {
		result.HWIDSettings = override.HWIDSettings
	}
	if override.CustomRemarks != "" {
		result.CustomRemarks = override.CustomRemarks
	}
	return result
}

func applyHostOverrides(hosts []SubscriptionHost, overrides map[string]HostOverride) []SubscriptionHost {
	if len(overrides) == 0 {
		return hosts
	}
	for i := range hosts {
		h := &hosts[i]
		if override, ok := overrides[h.UUID]; ok {
			if override.Address != nil {
				h.Address = *override.Address
			}
			if override.Port != nil {
				h.Port = *override.Port
			}
			if override.Remark != nil {
				h.Remark = *override.Remark
			}
			if override.SNI != nil {
				h.SNI = override.SNI
			}
			if override.Host != nil {
				h.Host = override.Host
			}
			if override.Path != nil {
				h.Path = override.Path
			}
		}
	}
	return hosts
}

func firstHostTag(tags []string) *string {
	for _, tag := range tags {
		trimmed := strings.TrimSpace(tag)
		if trimmed == "" {
			continue
		}
		return &trimmed
	}
	return nil
}

func matchResponseRules(rules *subscriptionresponserules.Config, headers http.Header) string {
	result := matchResponseRulesDetailed(rules, headers, "")
	if result.Matched && result.ResponseType != "" {
		return result.ResponseType
	}
	return defaultResponseType
}

func matchResponseRulesDetailed(rules *subscriptionresponserules.Config, headers http.Header, overrideClientType string) subscriptionresponserules.MatchResult {
	return subscriptionresponserules.MatchRulesDetailed(rules, headers, overrideClientType, mapClientTypeToResponseType, defaultResponseType)
}

func (s *RenderService) resolveSubscriptionURL(settings SubscriptionSettingsParsed, shortUUID string) string {
	schema := strings.ToLower(strings.TrimSpace(settings.Raw.APISchema))
	if schema == "" {
		schema = "https"
	}

	host := strings.TrimSpace(settings.Raw.Address)
	if host == "" {
		host = "localhost"
	}

	port := settings.Raw.Port
	portPart := ""
	if port > 0 && !isDefaultPort(schema, port) {
		portPart = fmt.Sprintf(":%d", port)
	}

	path := strings.TrimSpace(settings.Raw.APIPath)
	if path != "" {
		path = strings.TrimPrefix(path, "/")
		return fmt.Sprintf("%s://%s%s/%s/%s", schema, host, portPart, path, shortUUID)
	}

	return fmt.Sprintf("%s://%s%s/%s", schema, host, portPart, shortUUID)
}

func isDefaultPort(schema string, port int) bool {
	switch schema {
	case "http":
		return port == 80
	case "https", "grpcs", "tls":
		return port == 443
	default:
		return false
	}
}

func getSubscriptionUserInfo(user SubscriptionUser) map[string]int64 {
	expire := user.ExpireAt.Unix()
	if user.ExpireAt.Year() == 2099 {
		expire = 0
	}

	return map[string]int64{
		"upload":   0,
		"download": user.UsedTrafficBytes,
		"total":    user.TrafficLimitBytes,
		"expire":   expire,
	}
}

func getSubscriptionRefillDate(strategy string) string {
	return getSubscriptionRefillDateAt(strategy, time.Now())
}

func getSubscriptionRefillDateAt(strategy string, now time.Time) string {
	now = now.Local()
	switch strings.ToUpper(strategy) {
	case "DAY":
		now = now.AddDate(0, 0, 1)
		now = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		return fmt.Sprintf("%d", now.Unix())
	case "WEEK":
		offset := (int(time.Monday) - int(now.Weekday()) + 7) % 7
		if offset == 0 {
			offset = 7
		}
		now = now.AddDate(0, 0, offset)
		now = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		return fmt.Sprintf("%d", now.Unix())
	case "MONTH":
		now = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		now = now.AddDate(0, 1, 0)
		return fmt.Sprintf("%d", now.Unix())
	default:
		return ""
	}
}

func formatTemplateValue(value string, user SubscriptionUser, _ SubscriptionSettingsParsed, subscriptionURL string) string {
	replacer := strings.NewReplacer(
		"{USERNAME}", user.Username,
		"{{username}}", user.Username,
		"{SHORT_UUID}", user.ShortUUID,
		"{{shortUuid}}", user.ShortUUID,
		"{SUBSCRIPTION_URL}", subscriptionURL,
		"{{subscriptionUrl}}", subscriptionURL,
	)
	return replacer.Replace(value)
}

func (s *RenderService) buildSubscriptionHeaders(user SubscriptionUser, settings SubscriptionSettingsParsed, isHapp bool) map[string]string {
	headers := map[string]string{}
	subscriptionURL := s.resolveSubscriptionURL(settings, user.ShortUUID)

	headers["content-disposition"] = fmt.Sprintf("attachment; filename=%s", user.Username)
	headers["support-url"] = settings.Raw.SupportLink

	profileTitle := formatTemplateValue(settings.Raw.ProfileTitle, user, settings, subscriptionURL)
	headers["profile-title"] = fmt.Sprintf("base64:%s", base64.StdEncoding.EncodeToString([]byte(profileTitle)))
	headers["profile-update-interval"] = fmt.Sprintf("%d", settings.Raw.ProfileUpdateInterval)

	userInfo := getSubscriptionUserInfo(user)
	parts := []string{}
	for key, val := range userInfo {
		parts = append(parts, fmt.Sprintf("%s=%d", key, val))
	}
	sort.Strings(parts)
	headers["subscription-userinfo"] = strings.Join(parts, "; ")

	if settings.Raw.HappAnnounce != "" {
		announce := formatTemplateValue(settings.Raw.HappAnnounce, user, settings, subscriptionURL)
		headers["announce"] = fmt.Sprintf("base64:%s", base64.StdEncoding.EncodeToString([]byte(announce)))
	}

	if isHapp && settings.Raw.HappRouting != "" {
		headers["routing"] = settings.Raw.HappRouting
	}

	if settings.Raw.IsProfileWebpageURLEnabled {
		headers["profile-web-page-url"] = subscriptionURL
	}

	if refillDate := getSubscriptionRefillDate(user.TrafficLimitStrategy); refillDate != "" {
		headers["subscription-refill-date"] = refillDate
	}

	for key, value := range settings.CustomResponseHeaders {
		headers[key] = formatTemplateValue(value, user, settings, subscriptionURL)
	}

	return headers
}

func (s *RenderService) filterHostsForResponseType(hosts []SubscriptionHost, responseType string, includeDisabled bool) []SubscriptionHost {
	filtered := make([]SubscriptionHost, 0, len(hosts))
	for _, host := range hosts {
		if !includeDisabled && host.IsDisabled {
			continue
		}
		if responseType != "" && len(host.ExcludeFromSubscriptionTypes) > 0 {
			exclude := false
			for _, t := range host.ExcludeFromSubscriptionTypes {
				if strings.EqualFold(t, responseType) {
					exclude = true
					break
				}
			}
			if exclude {
				continue
			}
		}
		filtered = append(filtered, host)
	}

	return filtered
}

func (s *RenderService) generateSubscriptionContent(responseType string, templateData []byte, hosts []SubscriptionHost, user SubscriptionUser) (SubscriptionWithConfig, error) {
	switch responseType {
	case responseTypeXrayJSON:
		body, err := generateXrayJSONConfig(templateData, hosts, user)
		return SubscriptionWithConfig{Body: body, ContentType: "application/json"}, err
	case responseTypeMihomo, responseTypeStash, responseTypeClash:
		body, err := generateYAMLConfig(templateData, hosts, user)
		return SubscriptionWithConfig{Body: body, ContentType: "text/yaml"}, err
	case responseTypeSingbox:
		body, err := generateSingboxConfig(templateData, hosts, user)
		return SubscriptionWithConfig{Body: body, ContentType: "application/json"}, err
	case responseTypeXrayBase64:
		links, _ := buildSubscriptionLinks(hosts, user)
		joined := strings.Join(links, "\n")
		encoded := base64.StdEncoding.EncodeToString([]byte(joined))
		return SubscriptionWithConfig{Body: encoded, ContentType: "text/plain"}, nil
	default:
		links, _ := buildSubscriptionLinks(hosts, user)
		joined := strings.Join(links, "\n")
		encoded := base64.StdEncoding.EncodeToString([]byte(joined))
		return SubscriptionWithConfig{Body: encoded, ContentType: "text/plain"}, nil
	}
}

func (s *RenderService) shuffleHostsIfNeeded(hosts []SubscriptionHost, settings SubscriptionSettingsParsed) {
	if settings.Raw.RandomizeHosts {
		rand.Shuffle(len(hosts), func(i, j int) {
			hosts[i], hosts[j] = hosts[j], hosts[i]
		})
	}
}

func responseTypeToTemplateType(responseType string) string {
	switch responseType {
	case responseTypeXrayJSON:
		return responseTypeXrayJSON
	case responseTypeMihomo:
		return responseTypeMihomo
	case responseTypeStash:
		return responseTypeStash
	case responseTypeClash:
		return responseTypeClash
	case responseTypeSingbox:
		return responseTypeSingbox
	default:
		return responseTypeXrayJSON
	}
}

func (s *RenderService) buildSubscriptionInfoResponse(user SubscriptionUser, settings SubscriptionSettingsParsed, hosts []SubscriptionHost) SubscriptionInfoResponse {
	filtered := s.filterHostsForResponseType(hosts, "", false)
	links, ssConfLinks := buildSubscriptionLinks(filtered, user)

	daysLeft := int(time.Until(user.ExpireAt).Hours() / 24)
	if daysLeft < 0 {
		daysLeft = 0
	}

	infoUser := SubscriptionInfoUser{
		ShortUUID:                user.ShortUUID,
		DaysLeft:                 daysLeft,
		TrafficUsed:              shared.FormatBytes(user.UsedTrafficBytes),
		TrafficLimit:             shared.FormatBytes(user.TrafficLimitBytes),
		LifetimeTrafficUsed:      shared.FormatBytes(user.LifetimeUsedBytes),
		TrafficUsedBytes:         fmt.Sprintf("%d", user.UsedTrafficBytes),
		TrafficLimitBytes:        fmt.Sprintf("%d", user.TrafficLimitBytes),
		LifetimeTrafficUsedBytes: fmt.Sprintf("%d", user.LifetimeUsedBytes),
		Username:                 user.Username,
		ExpiresAt:                user.ExpireAt,
		IsActive:                 strings.EqualFold(user.Status, "ACTIVE"),
		UserStatus:               user.Status,
		TrafficLimitStrategy:     user.TrafficLimitStrategy,
	}

	return SubscriptionInfoResponse{
		IsFound:         true,
		User:            infoUser,
		Links:           links,
		SSConfLinks:     ssConfLinks,
		SubscriptionURL: s.resolveSubscriptionURL(settings, user.ShortUUID),
	}
}
