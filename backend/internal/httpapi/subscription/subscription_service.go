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
	"context"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"sort"
	"strings"
	"time"

	"exodus/internal/config"
	"exodus/internal/httpapi/subscriptionresponserules"
	"exodus/internal/httpapi/subscriptionsettings"
	"exodus/internal/logger"
)

type RenderService struct {
	db           *sql.DB
	backgroundDB *sql.DB
	cfg          *config.BackendConfig
}

func NewRenderService(db, backgroundDB *sql.DB, cfg *config.BackendConfig) *RenderService {
	return &RenderService{db: db, backgroundDB: backgroundDB, cfg: cfg}
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
	}

	if len(overrides.HostOverrides) > 0 {
		base.HostOverrides = overrides.HostOverrides
	}

	if len(overrides.ResponseHeaders) > 0 {
		if base.ResponseHeaders == nil {
			base.ResponseHeaders = make(map[string]string)
		}
		for k, v := range overrides.ResponseHeaders {
			base.ResponseHeaders[k] = v
		}
	}

	if overrides.HwidSettings != nil {
		base.HwidSettings = *overrides.HwidSettings
	}

	if overrides.CustomRemarks != nil {
		base.CustomRemarks = *overrides.CustomRemarks
	}

	return base
}

func mergeSubscriptionSettings(base, squad subscriptionsettings.SubscriptionSettings) subscriptionsettings.SubscriptionSettings {
	if strings.TrimSpace(squad.ProfileTitle) != "" {
		base.ProfileTitle = squad.ProfileTitle
	}
	if strings.TrimSpace(squad.SupportLink) != "" {
		base.SupportLink = squad.SupportLink
	}
	if squad.ProfileUpdateInterval > 0 {
		base.ProfileUpdateInterval = squad.ProfileUpdateInterval
	}
	if strings.TrimSpace(squad.HappAnnounce) != "" {
		base.HappAnnounce = squad.HappAnnounce
	}
	if strings.TrimSpace(squad.HappRouting) != "" {
		base.HappRouting = squad.HappRouting
	}
	base.IsProfileWebpageURLEnabled = squad.IsProfileWebpageURLEnabled
	base.ServeJSONAtBaseSubscription = squad.ServeJSONAtBaseSubscription
	base.IsShowCustomRemarks = squad.IsShowCustomRemarks
	base.RandomizeHosts = squad.RandomizeHosts
	return base
}

func (s *RenderService) RenderUserSubscription(
	ctx context.Context,
	user SubscriptionUser,
	userAgent string,
	requestedType string,
	requestIP string,
	hwid *HwidHeaders,
) ([]byte, string, map[string]string, error) {
	if user.Status != "ACTIVE" || (!user.ExpireAt.IsZero() && user.ExpireAt.Before(time.Now())) {
		return nil, "", nil, ErrUserDisabled
	}

	settings, err := loadSubscriptionSettings(ctx, s.db, s.cfg)
	if err != nil {
		return nil, "", nil, err
	}

	squadOverrides, _ := loadExternalSquadOverrides(ctx, s.db, ptrString(user.ExternalSquadUUID), s.cfg)
	settings = applyExternalSquadOverrides(settings, squadOverrides)

	hosts, err := getHostsForUser(ctx, s.db, user)
	if err != nil {
		return nil, "", nil, err
	}
	if len(hosts) == 0 {
		return nil, "", nil, ErrNoHosts
	}

	if len(settings.HostOverrides) > 0 {
		hosts = applyHostOverrides(hosts, settings.HostOverrides)
	}

	if settings.Raw.RandomizeHosts {
		rand.Shuffle(len(hosts), func(i, j int) {
			hosts[i], hosts[j] = hosts[j], hosts[i]
		})
	}

	if settings.HwidSettings.Enabled {
		ok, _, limitReached := checkHwidDeviceLimit(ctx, s.db, user, hwid, settings.HwidSettings)
		if !ok {
			if limitReached {
				return nil, "", nil, ErrHwidLimitExceeded
			}
			return nil, "", nil, ErrHwidRequired
		}
	} else if hwid != nil {
		_ = enqueueOrUpsertHwidUserDevice(ctx, s.db, user.TID, *hwid)
	}

	updateSubscriptionRequest(ctx, s.backgroundDB, user.UUID, user.TID, userAgent, requestIP)

	reqType := strings.ToUpper(strings.TrimSpace(requestedType))
	if reqType == "" {
		h := http.Header{}
		if userAgent != "" {
			h.Set("User-Agent", userAgent)
		}
		if settings.ResponseRules != nil {
			reqType = matchResponseRules(settings.ResponseRules, h)
		}
	}
	if reqType == "" {
		reqType = defaultResponseType
	}

	if s.cfg != nil && s.cfg.Logger != nil {
		s.cfg.Logger.RoleService(logger.RoleAPI, logger.ServiceHTTP).Debug(
			"Rendering user subscription",
			"username", user.Username,
			"short_uuid", user.ShortUUID,
			"user_agent", userAgent,
			"client_ip", requestIP,
			"requested_type", requestedType,
			"resolved_type", reqType,
			"used_traffic_bytes", user.UsedTrafficBytes,
			"traffic_limit_bytes", user.TrafficLimitBytes,
			"active_hosts", len(hosts),
		)
	}

	templateType := strings.ToLower(reqType)
	templateContent, _ := getSubscriptionTemplate(ctx, s.db, templateType)

	var outputContent string
	var contentType string

	switch reqType {
	case responseTypeMihomo, responseTypeClash, responseTypeStash:
		gen := NewMihomoGenerator(s.cfg)
		out, err := gen.Generate([]byte(templateContent), user, hosts, settings)
		if err != nil {
			return nil, "", nil, err
		}
		outputContent = out
		contentType = "text/yaml; charset=utf-8"
	case responseTypeSingbox:
		gen := NewSingboxGenerator(s.cfg)
		out, err := gen.Generate([]byte(templateContent), user, hosts, settings)
		if err != nil {
			return nil, "", nil, err
		}
		outputContent = out
		contentType = "application/json; charset=utf-8"
	case responseTypeXrayJSON:
		gen := NewXrayGenerator(s.cfg)
		out, err := gen.GenerateJSON([]byte(templateContent), user, hosts, settings)
		if err != nil {
			return nil, "", nil, err
		}
		outputContent = out
		contentType = "application/json; charset=utf-8"
	default:
		gen := NewXrayGenerator(s.cfg)
		links, err := gen.GenerateLinks(user, hosts, settings)
		if err != nil {
			return nil, "", nil, err
		}
		outputContent = base64.StdEncoding.EncodeToString([]byte(strings.Join(links, "\n")))
		contentType = "text/plain; charset=utf-8"
	}

	responseHeaders := buildResponseHeaders(user, settings, contentType)
	return []byte(outputContent), contentType, responseHeaders, nil
}

func (s *RenderService) buildSubscriptionInfoResponse(
	user SubscriptionUser,
	settings SubscriptionSettingsParsed,
	hosts []SubscriptionHost,
) SubscriptionInfoResponse {
	links, ssConfLinks := buildSubscriptionLinks(hosts, user)

	domain := strings.TrimSpace(settings.Raw.Address)
	if domain == "" {
		domain = "panel.exodus.dev"
	}
	scheme := strings.TrimSpace(settings.Raw.APISchema)
	if scheme == "" {
		scheme = "https"
	}
	apiPath := strings.Trim(strings.TrimSpace(settings.Raw.APIPath), "/")
	if apiPath == "" {
		apiPath = "api/sub"
	}
	subURL := fmt.Sprintf("%s://%s/%s/%s", scheme, domain, apiPath, user.ShortUUID)

	return SubscriptionInfoResponse{
		IsFound: true,
		User: SubscriptionInfoUser{
			ShortUUID:            user.ShortUUID,
			Username:             user.Username,
			TrafficLimitBytes:    fmt.Sprintf("%d", user.TrafficLimitBytes),
			TrafficUsedBytes:     fmt.Sprintf("%d", user.UsedTrafficBytes),
			ExpiresAt:            user.ExpireAt,
			IsActive:             user.Status == "ACTIVE",
			UserStatus:           user.Status,
			TrafficLimitStrategy: user.TrafficLimitStrategy,
		},
		Links:           links,
		SSConfLinks:     ssConfLinks,
		SubscriptionURL: subURL,
	}
}

func sortHosts(hosts []SubscriptionHost) {
	sort.Slice(hosts, func(i, j int) bool {
		if hosts[i].ViewPosition != hosts[j].ViewPosition {
			return hosts[i].ViewPosition < hosts[j].ViewPosition
		}
		return hosts[i].Remark < hosts[j].Remark
	})
}
