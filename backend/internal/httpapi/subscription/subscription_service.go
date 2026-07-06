package subscription

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"regexp"
	"sort"
	"strings"
	"time"

	"exodus/internal/config"
	dbmanager "exodus/internal/db/manager"
	"exodus/internal/dbutil"
	"exodus/internal/httpapi/shared"
	"exodus/internal/httpapi/subscriptionresponserules"
	"exodus/internal/httpapi/subscriptionsettings"
	"exodus/internal/jobqueue"
	"exodus/internal/logger"

	"github.com/google/uuid"
)

const (
	responseTypeBrowser    = "BROWSER"
	responseTypeXrayBase64 = "XRAY_BASE64"
	responseTypeXrayJSON   = "XRAY_JSON"
	responseTypeMihomo     = "MIHOMO"
	responseTypeStash      = "STASH"
	responseTypeClash      = "CLASH"
	responseTypeSingbox    = "SINGBOX"
	responseTypeBlock      = "BLOCK"
	responseTypeStatus404  = "STATUS_CODE_404"
	responseTypeStatus451  = "STATUS_CODE_451"
	responseTypeSocketDrop = "SOCKET_DROP"
)

var defaultResponseType = responseTypeXrayBase64

const defaultSubpageConfigUUID = "00000000-0000-0000-0000-000000000000"

type ExternalSquadOverrides struct {
	SubscriptionSettings *subscriptionsettings.SubscriptionSettings `json:"subscription_settings"`
	HostOverrides        map[string]HostOverride                    `json:"host_overrides"`
	ResponseHeaders      map[string]string                          `json:"response_headers"`
	HwidSettings         *HwidSettings                              `json:"hwid_settings"`
	CustomRemarks        *CustomRemarks                             `json:"custom_remarks"`
	Templates            map[string]string
}

type HwidSettingsInput struct {
	Enabled             *bool `json:"enabled"`
	FallbackDeviceLimit *int  `json:"fallbackDeviceLimit"`
	MaxDevicesAnnounce  *int  `json:"maxDevicesAnnounce"`
}

type UpdateExternalSquadInput struct {
	UUID              string           `json:"uuid"`
	Name              *string          `json:"name"`
	SubpageConfigUUID *string          `json:"subpageConfigUuid"`
	CustomRemarks     *json.RawMessage `json:"customRemarks"`
	HwidSettings      json.RawMessage  `json:"hwidSettings"`
}
type HostOverride struct {
	Address *string `json:"address"`
	Port    *int    `json:"port"`
	Remark  *string `json:"remark"`
	SNI     *string `json:"sni"`
	Host    *string `json:"host"`
	Path    *string `json:"path"`
}

type SubscriptionSettingsParsed struct {
	Raw subscriptionsettings.SubscriptionSettings

	CustomResponseHeaders map[string]string
	ResponseRules         *subscriptionresponserules.Config
	HwidSettings          HwidSettings
	CustomRemarks         CustomRemarks
}

type HwidSettings struct {
	Enabled             bool    `json:"enabled"`
	FallbackDeviceLimit int     `json:"fallbackDeviceLimit"`
	MaxDevicesAnnounce  *string `json:"maxDevicesAnnounce"`
}

type CustomRemarks struct {
	ExpiredUsers           []string `json:"expiredUsers"`
	LimitedUsers           []string `json:"limitedUsers"`
	DisabledUsers          []string `json:"disabledUsers"`
	EmptyHosts             []string `json:"emptyHosts"`
	HWIDMaxDevicesExceeded []string `json:"HWIDMaxDevicesExceeded"`
	HWIDNotSupported       []string `json:"HWIDNotSupported"`
}

type SubscriptionUser struct {
	TID                  int64
	UUID                 string
	ShortUUID            string
	Username             string
	Status               string
	TrafficLimitBytes    int64
	TrafficLimitStrategy string
	ExpireAt             time.Time
	TrojanPassword       string
	VlessUUID            string
	SSPassword           string
	NaivePassword        string
	ShadowtlsPassword    string
	Hysteria2Password    string
	AnytlsPassword       string
	HwidDeviceLimit      *int
	ExternalSquadUUID    *string
	UsedTrafficBytes     int64
	LifetimeUsedBytes    int64
}

type SubscriptionHost struct {
	UUID                         string
	ViewPosition                 int
	Remark                       string
	Address                      string
	Port                         int
	Path                         *string
	SNI                          *string
	Host                         *string
	ALPN                         *string
	Fingerprint                  *string
	SecurityLayer                string
	XHTTPExtraParams             *string
	MuxParams                    *string
	SingboxMuxParams             *string
	ClashMuxParams               *string
	SockoptParams                *string
	IsDisabled                   bool
	ServerDescription            *string
	OverrideProtocolCredential   bool
	ProtocolCredential           *string
	AllowInsecure                bool
	ShuffleHost                  bool
	MihomoX25519                 bool
	MihomoIPVersion              *string
	XrayJSONTemplateUUID         *string
	KeepSNIBlank                 bool
	ExcludeFromSubscriptionTypes []string
	Tag                          *string
	IsHidden                     bool
	OverrideSNIFromAddress       bool
	ConfigProfileUUID            *string
	ConfigProfileInboundUUID     *string

	InboundTag      *string
	InboundType     *string
	InboundNetwork  *string
	InboundSecurity *string
	InboundPort     *int
	InboundRaw      json.RawMessage
}

type RawHost struct {
	UUID          string  `json:"uuid"`
	Remark        string  `json:"remark"`
	Address       string  `json:"address"`
	Port          int     `json:"port"`
	Protocol      string  `json:"protocol"`
	Network       *string `json:"network,omitempty"`
	Security      *string `json:"security,omitempty"`
	Path          *string `json:"path,omitempty"`
	SNI           *string `json:"sni,omitempty"`
	Host          *string `json:"host,omitempty"`
	ALPN          *string `json:"alpn,omitempty"`
	Fingerprint   *string `json:"fingerprint,omitempty"`
	AllowInsecure bool    `json:"allow_insecure"`
	IsDisabled    bool    `json:"is_disabled"`
	IsHidden      bool    `json:"is_hidden"`
}

type SubscriptionInfoUser struct {
	ShortUUID                string    `json:"shortUuid"`
	DaysLeft                 int       `json:"daysLeft"`
	TrafficUsed              string    `json:"trafficUsed"`
	TrafficLimit             string    `json:"trafficLimit"`
	LifetimeTrafficUsed      string    `json:"lifetimeTrafficUsed"`
	TrafficUsedBytes         string    `json:"trafficUsedBytes"`
	TrafficLimitBytes        string    `json:"trafficLimitBytes"`
	LifetimeTrafficUsedBytes string    `json:"lifetimeTrafficUsedBytes"`
	Username                 string    `json:"username"`
	ExpiresAt                time.Time `json:"expiresAt"`
	IsActive                 bool      `json:"isActive"`
	UserStatus               string    `json:"userStatus"`
	TrafficLimitStrategy     string    `json:"trafficLimitStrategy"`
}

type SubscriptionInfoResponse struct {
	IsFound         bool                 `json:"isFound"`
	User            SubscriptionInfoUser `json:"user"`
	Links           []string             `json:"links"`
	SSConfLinks     map[string]string    `json:"ssConfLinks"`
	SubscriptionURL string               `json:"subscriptionUrl"`
}

type RawSubscriptionResponse struct {
	User              map[string]interface{} `json:"user"`
	ConvertedUserInfo map[string]interface{} `json:"convertedUserInfo"`
	RawHosts          []RawHost              `json:"rawHosts"`
	Headers           map[string]string      `json:"headers"`
}

type SubscriptionHeaders struct {
	ContentDisposition     string            `json:"content-disposition,omitempty"`
	ProfileTitle           string            `json:"profile-title,omitempty"`
	ProfileUpdateInterval  string            `json:"profile-update-interval,omitempty"`
	ProfileWebPageURL      string            `json:"profile-web-page-url,omitempty"`
	SubscriptionUserInfo   string            `json:"subscription-userinfo,omitempty"`
	SubscriptionRefillDate string            `json:"subscription-refill-date,omitempty"`
	SupportURL             string            `json:"support-url,omitempty"`
	Announce               string            `json:"announce,omitempty"`
	Routing                string            `json:"routing,omitempty"`
	Extra                  map[string]string `json:"-"`
}

type SubscriptionWithConfig struct {
	Headers     map[string]string
	ContentType string
	Body        string
}

type HwidHeaders struct {
	Hwid        string
	Platform    *string
	OsVersion   *string
	DeviceModel *string
	UserAgent   *string
	RequestIP   *string
	Synthetic   bool
}

func loadSubscriptionSettings(ctx context.Context, manager *dbmanager.DatabaseManager, _ *config.BackendConfig) (SubscriptionSettingsParsed, error) {
	var parsed SubscriptionSettingsParsed

	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		row := db.QueryRowContext(ctx, `
            SELECT uuid, profile_title, support_link, profile_update_interval,
                   address, port, api_schema, api_path,
                   happ_announce, happ_routing, created_at, updated_at,
                   is_profile_webpage_url_enabled, serve_json_at_base_subscription,
                   is_show_custom_remarks, custom_response_headers, randomize_hosts,
                   response_rules, hwid_settings, custom_remarks
            FROM subscription_settings
            ORDER BY created_at ASC
            LIMIT 1
        `)

		settings, err := subscriptionsettings.ScanSubscriptionSettings(row)
		if err != nil {
			return err
		}

		parsed.Raw = settings
		parsed.CustomResponseHeaders = map[string]string{}
		parsed.CustomRemarks = CustomRemarks{}
		parsed.HwidSettings = HwidSettings{Enabled: false, FallbackDeviceLimit: 999}
		return nil
	})
	if err != nil {
		return parsed, err
	}

	if strings.TrimSpace(parsed.Raw.CustomResponseHeaders) != "" {
		_ = json.Unmarshal([]byte(parsed.Raw.CustomResponseHeaders), &parsed.CustomResponseHeaders)
	}

	if strings.TrimSpace(parsed.Raw.ResponseRules) != "" {
		var rules subscriptionresponserules.Config
		if err := json.Unmarshal([]byte(parsed.Raw.ResponseRules), &rules); err == nil {
			parsed.ResponseRules = &rules
		}
	}

	if strings.TrimSpace(parsed.Raw.HWIDSettings) != "" {
		var hwid HwidSettings
		if err := json.Unmarshal([]byte(parsed.Raw.HWIDSettings), &hwid); err == nil {
			parsed.HwidSettings = hwid
		}
	}

	if strings.TrimSpace(parsed.Raw.CustomRemarks) != "" {
		var remarks CustomRemarks
		if err := json.Unmarshal([]byte(parsed.Raw.CustomRemarks), &remarks); err == nil {
			parsed.CustomRemarks = remarks
		}
	}

	return parsed, nil
}

func loadExternalSquadOverrides(ctx context.Context, manager *dbmanager.DatabaseManager, squadUUID string, cfg *config.BackendConfig) (*ExternalSquadOverrides, error) {
	log := cfg.Logger.RoleService(logger.RoleAPI, logger.ServiceHTTP)
	overrides := &ExternalSquadOverrides{
		Templates: make(map[string]string),
	}

	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		var subscriptionSettingsJSON, hostOverridesJSON, responseHeadersJSON, hwidSettingsJSON, customRemarksJSON sql.NullString

		query := `SELECT subscription_settings, host_overrides, response_headers, hwid_settings, custom_remarks
				  FROM external_squads WHERE uuid = ? LIMIT 1`
		row := db.QueryRowContext(ctx, query, squadUUID)

		if err := row.Scan(&subscriptionSettingsJSON, &hostOverridesJSON, &responseHeadersJSON, &hwidSettingsJSON, &customRemarksJSON); err != nil {
			return err
		}

		if subscriptionSettingsJSON.Valid && subscriptionSettingsJSON.String != "" {
			var ss subscriptionsettings.SubscriptionSettings
			if err := json.Unmarshal([]byte(subscriptionSettingsJSON.String), &ss); err == nil {
				overrides.SubscriptionSettings = &ss
				log.Debug("Loaded subscription_settings override")
			}
		}
		if hostOverridesJSON.Valid && hostOverridesJSON.String != "" {
			var ho map[string]HostOverride
			if err := json.Unmarshal([]byte(hostOverridesJSON.String), &ho); err == nil {
				overrides.HostOverrides = ho
				log.Debug("Loaded host_overrides override", "count", len(ho))
			}
		}
		if responseHeadersJSON.Valid && responseHeadersJSON.String != "" {
			var rh map[string]string
			if err := json.Unmarshal([]byte(responseHeadersJSON.String), &rh); err == nil {
				overrides.ResponseHeaders = rh
				log.Debug("Loaded response_headers override")
			}
		}
		if hwidSettingsJSON.Valid && hwidSettingsJSON.String != "" {
			var hs HwidSettings
			if err := json.Unmarshal([]byte(hwidSettingsJSON.String), &hs); err == nil {
				overrides.HwidSettings = &hs
				log.Debug("Loaded hwid_settings override")
			}
		}
		if customRemarksJSON.Valid && customRemarksJSON.String != "" {
			var cr CustomRemarks
			if err := json.Unmarshal([]byte(customRemarksJSON.String), &cr); err == nil {
				overrides.CustomRemarks = &cr
				log.Debug("Loaded custom_remarks override")
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	err = manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		rows, err := db.QueryContext(ctx, `
			SELECT t.name, est.template_type
			FROM external_squads_templates est
			JOIN subscription_templates t ON t.uuid = est.template_uuid
			WHERE est.external_squad_uuid = ?
		`, squadUUID)
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			var templateName, templateType string
			if err := rows.Scan(&templateName, &templateType); err != nil {
				return err
			}
			overrides.Templates[strings.ToUpper(templateType)] = templateName
			log.Debug("Loaded template override", "type", templateType, "name", templateName)
		}
		return rows.Err()
	})
	if err != nil {
		log.Warn("Failed to load external squad templates", "error", err)
	}

	return overrides, nil
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

func getSubscriptionUserByShortUUID(ctx context.Context, manager *dbmanager.DatabaseManager, shortUUID string) (SubscriptionUser, error) {
	return getSubscriptionUserByField(ctx, manager, "short_uuid", shortUUID)
}

func getSubscriptionUserByUUID(ctx context.Context, manager *dbmanager.DatabaseManager, userUUID string) (SubscriptionUser, error) {
	return getSubscriptionUserByField(ctx, manager, "uuid", userUUID)
}

func getSubscriptionUserByUsername(ctx context.Context, manager *dbmanager.DatabaseManager, username string) (SubscriptionUser, error) {
	return getSubscriptionUserByField(ctx, manager, "username", username)
}

func getSubscriptionUserByField(ctx context.Context, manager *dbmanager.DatabaseManager, field string, value string) (SubscriptionUser, error) {
	var user SubscriptionUser

	var where string
	switch field {
	case "short_uuid":
		where = "u.short_uuid = ?"
	case "uuid":
		where = "u.uuid = ?"
	case "username":
		where = "u.username = ?"
	default:
		return user, fmt.Errorf("unsupported search field")
	}

	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		query := fmt.Sprintf(`
            SELECT u.t_id, u.uuid, u.short_uuid, u.username, u.status,
                   u.traffic_limit_bytes, u.traffic_limit_strategy, u.expire_at,
	                   u.trojan_password, u.vless_uuid, u.ss_password,
	                   u.naive_password, u.shadowtls_password, u.hysteria2_password, u.anytls_password,
	                   u.hwid_device_limit, u.external_squad_uuid,
                   COALESCE(ut.used_traffic_bytes, 0), COALESCE(ut.lifetime_used_traffic_bytes, 0)
            FROM users u
            LEFT JOIN user_traffic ut ON ut.t_id = u.t_id
            WHERE %s
            LIMIT 1
        `, where)

		row := db.QueryRowContext(ctx, query, value)

		var hwidDeviceLimit sql.NullInt64
		var externalSquadUUID sql.NullString
		var naivePassword, shadowtlsPassword, hysteria2Password, anytlsPassword sql.NullString
		if err := row.Scan(
			&user.TID,
			&user.UUID,
			&user.ShortUUID,
			&user.Username,
			&user.Status,
			&user.TrafficLimitBytes,
			&user.TrafficLimitStrategy,
			&user.ExpireAt,
			&user.TrojanPassword,
			&user.VlessUUID,
			&user.SSPassword,
			&naivePassword,
			&shadowtlsPassword,
			&hysteria2Password,
			&anytlsPassword,
			&hwidDeviceLimit,
			&externalSquadUUID,
			&user.UsedTrafficBytes,
			&user.LifetimeUsedBytes,
		); err != nil {
			return err
		}

		if hwidDeviceLimit.Valid {
			v := int(hwidDeviceLimit.Int64)
			user.HwidDeviceLimit = &v
		}
		if externalSquadUUID.Valid {
			v := externalSquadUUID.String
			user.ExternalSquadUUID = &v
		}
		user.NaivePassword = nullableSQLString(naivePassword)
		user.ShadowtlsPassword = nullableSQLString(shadowtlsPassword)
		user.Hysteria2Password = nullableSQLString(hysteria2Password)
		user.AnytlsPassword = nullableSQLString(anytlsPassword)
		return nil
	})
	if err != nil {
		return user, err
	}

	return user, nil
}

func getHostsForUser(ctx context.Context, manager *dbmanager.DatabaseManager, user SubscriptionUser) ([]SubscriptionHost, error) {
	var hosts []SubscriptionHost

	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		rows, err := db.QueryContext(ctx, `
                SELECT DISTINCT h.uuid, h.view_position, h.remark, h.address, h.port,
                       h.path, h.sni, h.host, h.alpn, h.fingerprint, h.security_layer,
                       h.xhttp_extra_params, h.mux_params, h.singbox_mux_params, h.clash_mux_params, h.sockopt_params, h.is_disabled,
                       h.server_description, h.override_protocol_credential, h.protocol_credential, h.allow_insecure, h.shuffle_host,
                       h.mihomo_x25519, h.mihomo_ip_version, h.xray_json_template_uuid, h.keep_sni_blank,
                       h.exclude_from_subscription_types, h.tags, h.is_hidden, h.override_sni_from_address,
                       h.config_profile_uuid, h.config_profile_inbound_uuid,
                       cpi.tag, cpi.type, cpi.network, cpi.security, cpi.port, cpi.raw_inbound
                FROM internal_squad_members ism
                JOIN internal_squad_inbounds isi ON ism.internal_squad_uuid = isi.internal_squad_uuid
                JOIN config_profile_inbounds cpi ON isi.inbound_uuid = cpi.uuid
                JOIN hosts h ON h.config_profile_inbound_uuid = cpi.uuid
                LEFT JOIN internal_squad_host_exclusions ihe
                    ON ihe.host_uuid = h.uuid AND ihe.squad_uuid = ism.internal_squad_uuid
                WHERE ism.user_id = ? AND ihe.host_uuid IS NULL
                ORDER BY h.view_position ASC, h.remark ASC
            `, user.TID)
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			host, err := scanSubscriptionHost(rows)
			if err != nil {
				return err
			}
			hosts = append(hosts, host)
		}
		return rows.Err()
	})

	if err != nil {
		return nil, err
	}

	return hosts, nil
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

func scanSubscriptionHost(scanner shared.RowScanner) (SubscriptionHost, error) {
	var h SubscriptionHost
	var viewPosition sql.NullInt64
	var path, sni, host, alpn, fingerprint, securityLayer sql.NullString
	var xhttpExtraParams, muxParams, singboxMuxParams, clashMuxParams, sockoptParams, serverDescription, protocolCredential sql.NullString
	var xrayJSONTemplateUUID, mihomoIPVersion, configProfileUUID, configProfileInboundUUID sql.NullString
	var inboundTag, inboundType, inboundNetwork, inboundSecurity sql.NullString
	var inboundPort sql.NullInt64
	var rawInbound sql.NullString
	var excludeTypes, hostTags dbutil.StringArray
	var isDisabled, overrideProtocolCredential, allowInsecure, shuffleHost, mihomoX25519, keepSNIBlank, isHidden, overrideSNIFromAddress sql.NullBool

	err := scanner.Scan(
		&h.UUID,
		&viewPosition,
		&h.Remark,
		&h.Address,
		&h.Port,
		&path,
		&sni,
		&host,
		&alpn,
		&fingerprint,
		&securityLayer,
		&xhttpExtraParams,
		&muxParams,
		&singboxMuxParams,
		&clashMuxParams,
		&sockoptParams,
		&isDisabled,
		&serverDescription,
		&overrideProtocolCredential,
		&protocolCredential,
		&allowInsecure,
		&shuffleHost,
		&mihomoX25519,
		&mihomoIPVersion,
		&xrayJSONTemplateUUID,
		&keepSNIBlank,
		&excludeTypes,
		&hostTags,
		&isHidden,
		&overrideSNIFromAddress,
		&configProfileUUID,
		&configProfileInboundUUID,
		&inboundTag,
		&inboundType,
		&inboundNetwork,
		&inboundSecurity,
		&inboundPort,
		&rawInbound,
	)
	if err != nil {
		return h, err
	}

	if viewPosition.Valid {
		h.ViewPosition = int(viewPosition.Int64)
	}
	if path.Valid {
		h.Path = &path.String
	}
	if sni.Valid {
		h.SNI = &sni.String
	}
	if host.Valid {
		h.Host = &host.String
	}
	if alpn.Valid {
		h.ALPN = &alpn.String
	}
	if fingerprint.Valid {
		h.Fingerprint = &fingerprint.String
	}
	if securityLayer.Valid && securityLayer.String != "" {
		h.SecurityLayer = securityLayer.String
	} else {
		h.SecurityLayer = "DEFAULT"
	}
	if xhttpExtraParams.Valid {
		h.XHTTPExtraParams = &xhttpExtraParams.String
	}
	if muxParams.Valid {
		h.MuxParams = &muxParams.String
	}
	if singboxMuxParams.Valid {
		h.SingboxMuxParams = &singboxMuxParams.String
	}
	if clashMuxParams.Valid {
		h.ClashMuxParams = &clashMuxParams.String
	}
	if sockoptParams.Valid {
		h.SockoptParams = &sockoptParams.String
	}
	if isDisabled.Valid {
		h.IsDisabled = isDisabled.Bool
	}
	if serverDescription.Valid {
		h.ServerDescription = &serverDescription.String
	}
	if overrideProtocolCredential.Valid {
		h.OverrideProtocolCredential = overrideProtocolCredential.Bool
	}
	if protocolCredential.Valid {
		h.ProtocolCredential = &protocolCredential.String
	}
	if allowInsecure.Valid {
		h.AllowInsecure = allowInsecure.Bool
	}
	if shuffleHost.Valid {
		h.ShuffleHost = shuffleHost.Bool
	}
	if mihomoX25519.Valid {
		h.MihomoX25519 = mihomoX25519.Bool
	}
	if mihomoIPVersion.Valid {
		h.MihomoIPVersion = &mihomoIPVersion.String
	}
	if xrayJSONTemplateUUID.Valid {
		h.XrayJSONTemplateUUID = &xrayJSONTemplateUUID.String
	}
	if keepSNIBlank.Valid {
		h.KeepSNIBlank = keepSNIBlank.Bool
	}
	if firstTag := firstHostTag(hostTags.Slice()); firstTag != nil {
		h.Tag = firstTag
	}
	if isHidden.Valid {
		h.IsHidden = isHidden.Bool
	}
	if overrideSNIFromAddress.Valid {
		h.OverrideSNIFromAddress = overrideSNIFromAddress.Bool
	}
	if configProfileUUID.Valid {
		h.ConfigProfileUUID = &configProfileUUID.String
	}
	if configProfileInboundUUID.Valid {
		h.ConfigProfileInboundUUID = &configProfileInboundUUID.String
	}

	h.ExcludeFromSubscriptionTypes = excludeTypes.Slice()

	if inboundTag.Valid {
		h.InboundTag = &inboundTag.String
	}
	if inboundType.Valid {
		h.InboundType = &inboundType.String
	}
	if inboundNetwork.Valid {
		h.InboundNetwork = &inboundNetwork.String
	}
	if inboundSecurity.Valid {
		h.InboundSecurity = &inboundSecurity.String
	}
	if inboundPort.Valid {
		p := int(inboundPort.Int64)
		h.InboundPort = &p
	}
	if rawInbound.Valid {
		h.InboundRaw = json.RawMessage(rawInbound.String)
	}

	return h, nil
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

func extractHwidHeaders(r *http.Request) *HwidHeaders {
	hwid := strings.TrimSpace(r.Header.Get("X-HWID"))
	if hwid == "" {
		return nil
	}
	userAgent := firstNonEmptyHeader(r, "User-Agent", "X-HWID-User-Agent")
	platform := firstNonEmptyLowerHeader(r, "X-Device-OS", "X-HWID-Platform")
	osVersion := firstNonEmptyHeader(r, "X-Ver-OS", "X-HWID-OS-Version")
	deviceModel := firstNonEmptyHeader(r, "X-Device-Model", "X-HWID-Device-Model")
	platform, osVersion, deviceModel, userAgent = normalizeHwidMetadata(platform, osVersion, deviceModel, userAgent)

	h := &HwidHeaders{
		Hwid:        hwid,
		Platform:    platform,
		OsVersion:   osVersion,
		DeviceModel: deviceModel,
		UserAgent:   userAgent,
	}
	return h
}

func extractSyntheticHwidHeaders(r *http.Request, userUUID, requestIP string) *HwidHeaders {
	userAgent := strings.TrimSpace(r.Header.Get("User-Agent"))
	platform := firstNonEmptyLowerHeader(r, "X-Device-OS", "X-HWID-Platform")
	osVersion := firstNonEmptyHeader(r, "X-Ver-OS", "X-HWID-OS-Version")
	deviceModel := firstNonEmptyHeader(r, "X-Device-Model", "X-HWID-Device-Model")

	hasMetadata := userAgent != "" || platform != nil || osVersion != nil || deviceModel != nil
	if !hasMetadata {
		return nil
	}
	platform, osVersion, deviceModel, userAgentPtr := normalizeHwidMetadata(
		platform,
		osVersion,
		deviceModel,
		stringPtrIfNotEmpty(userAgent),
	)

	signature := strings.Join([]string{
		"exodus:synthetic-hwid:v1",
		"ua=" + strings.ToLower(ptrString(userAgentPtr)),
		"platform=" + strings.ToLower(ptrString(platform)),
		"os=" + strings.ToLower(ptrString(osVersion)),
		"model=" + strings.ToLower(ptrString(deviceModel)),
	}, "|")

	return &HwidHeaders{
		Hwid:        deterministicSyntheticHwid(userUUID, signature),
		Platform:    platform,
		OsVersion:   osVersion,
		DeviceModel: deviceModel,
		UserAgent:   userAgentPtr,
		RequestIP:   stringPtrIfNotEmpty(requestIP),
		Synthetic:   true,
	}
}

func normalizeHwidMetadata(platform, osVersion, deviceModel, userAgent *string) (*string, *string, *string, *string) {
	normalizedUserAgent := stringPtrIfNotEmpty(ptrString(userAgent))
	normalizedPlatform := lowerStringPtr(platform)
	if normalizedPlatform == nil {
		if inferred := inferPlatformFromUserAgent(ptrString(normalizedUserAgent)); inferred != "" {
			normalizedPlatform = &inferred
		}
	}
	if deviceModel == nil {
		deviceModel = stringPtrIfNotEmpty("unknown")
	}
	return normalizedPlatform, stringPtrIfNotEmpty(ptrString(osVersion)), stringPtrIfNotEmpty(ptrString(deviceModel)), normalizedUserAgent
}

func deterministicSyntheticHwid(userUUID, signature string) string {
	namespace, err := uuid.Parse(strings.TrimSpace(userUUID))
	if err != nil {
		namespace = uuid.NameSpaceOID
	}
	return uuid.NewSHA1(namespace, []byte(signature)).String()
}

func firstNonEmptyHeader(r *http.Request, names ...string) *string {
	for _, name := range names {
		value := strings.TrimSpace(r.Header.Get(name))
		if value == "" {
			continue
		}
		return &value
	}
	return nil
}

func firstNonEmptyLowerHeader(r *http.Request, names ...string) *string {
	value := firstNonEmptyHeader(r, names...)
	if value == nil {
		return nil
	}
	return lowerStringPtr(value)
}

func stringPtrIfNotEmpty(value string) *string {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil
	}
	return &value
}

func lowerStringPtr(value *string) *string {
	if value == nil {
		return nil
	}
	lowered := strings.ToLower(strings.TrimSpace(*value))
	if lowered == "" {
		return nil
	}
	return &lowered
}

func ptrString(value *string) string {
	if value == nil {
		return ""
	}
	return strings.TrimSpace(*value)
}

func inferClientAppFromUserAgent(userAgent string) string {
	userAgent = strings.TrimSpace(userAgent)
	if userAgent == "" {
		return ""
	}
	for _, sep := range []string{"/", " ", "(", ";"} {
		if idx := strings.Index(userAgent, sep); idx > 0 {
			return strings.TrimSpace(userAgent[:idx])
		}
	}
	return userAgent
}

func inferKnownClientPlatform(client string) string {
	client = strings.ToLower(strings.TrimSpace(client))

	switch client {
	case "sfa",
		"sfatv",
		"sfandroidtv",
		"v2rayng",
		"exclave",
		"nekoboxforandroid",
		"matsuri",
		"sagernet",
		"clashforandroid",
		"clashmetaforandroid",
		"cmfa":
		return "android"

	case "sfi",
		"streisand",
		"v2box",
		"rabbithole",
		"shadowrocket":
		return "ios"

	case "sft":
		return "tvos"

	case "sfw",
		"v2rayn":
		return "windows"

	case "sfm",
		"v2rayu",
		"v2rayx",
		"v2rayxs",
		"clashx":
		return "macos"

	case "sfl":
		return "linux"

	default:
		return ""
	}
}

func inferPlatformFromUserAgent(userAgent string) string {
	lower := strings.ToLower(userAgent)
	if lower == "" {
		return ""
	}

	if platform := inferKnownClientPlatform(inferClientAppFromUserAgent(lower)); platform != "" {
		return platform
	}

	if idx := strings.Index(lower, "platform/"); idx >= 0 {
		rest := userAgent[idx+len("platform/"):]
		for i, r := range rest {
			if !(r == '-' || r == '_' || r == '.' || r == '/' || r >= '0' && r <= '9' || r >= 'A' && r <= 'Z' || r >= 'a' && r <= 'z') {
				rest = rest[:i]
				break
			}
		}
		if rest = strings.Trim(rest, "/ "); rest != "" {
			return strings.ToLower(rest)
		}
	}

	switch {
	case strings.Contains(lower, "windows"):
		return "windows"
	case strings.Contains(lower, "android"):
		return "android"
	case strings.Contains(lower, "iphone") || strings.Contains(lower, "ipad") || strings.Contains(lower, "ios"):
		return "ios"
	case strings.Contains(lower, "mac os") || strings.Contains(lower, "macos") || strings.Contains(lower, "macintosh") || strings.Contains(lower, "darwin"):
		return "macos"
	case strings.Contains(lower, "linux"):
		return "linux"
	default:
		return ""
	}
}

func checkHwidDeviceLimit(ctx context.Context, manager *dbmanager.DatabaseManager, user SubscriptionUser, hwid *HwidHeaders, settings HwidSettings) (bool, bool, bool) {
	if user.HwidDeviceLimit != nil && *user.HwidDeviceLimit == 0 {
		if hwid != nil {
			_ = enqueueOrUpsertHwidUserDevice(ctx, manager, user.TID, *hwid)
		}
		return true, false, false
	}

	if hwid == nil {
		return false, false, true
	}

	exists, err := hwidDeviceExists(ctx, manager, user.TID, hwid.Hwid)
	if err == nil && exists {
		_ = enqueueOrUpsertHwidUserDevice(ctx, manager, user.TID, *hwid)
		return true, false, false
	}

	count, err := countHwidDevices(ctx, manager, user.TID)
	if err != nil {
		return false, true, false
	}

	limit := settings.FallbackDeviceLimit
	if user.HwidDeviceLimit != nil {
		limit = *user.HwidDeviceLimit
	}

	if count >= limit {
		return false, true, false
	}

	if err := upsertHwidUserDevice(ctx, manager, user.TID, *hwid); err != nil {
		return false, true, false
	}

	return true, false, false
}

func countHwidDevices(ctx context.Context, manager *dbmanager.DatabaseManager, userID int64) (int, error) {
	var count int
	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		return db.QueryRowContext(ctx, `SELECT COUNT(*) FROM hwid_user_devices WHERE user_id = ?`, userID).Scan(&count)
	})
	if err != nil {
		return 0, err
	}
	return count, nil
}

func hwidDeviceExists(ctx context.Context, manager *dbmanager.DatabaseManager, userID int64, hwid string) (bool, error) {
	var exists bool
	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		var tmp int
		err := db.QueryRowContext(ctx, `SELECT 1 FROM hwid_user_devices WHERE user_id = ? AND hwid = ?`, userID, hwid).Scan(&tmp)
		if err == sql.ErrNoRows {
			exists = false
			return nil
		}
		if err != nil {
			return err
		}
		exists = true
		return nil
	})
	if err != nil {
		return false, err
	}
	return exists, nil
}

func upsertHwidUserDevice(ctx context.Context, manager *dbmanager.DatabaseManager, userID int64, hwid HwidHeaders) error {
	hwid.Platform = lowerStringPtr(hwid.Platform)
	return manager.ExecuteLowPriority(func(db dbmanager.DBExecutor) error {
		_, err := db.ExecContext(ctx, `
            INSERT INTO hwid_user_devices (hwid, user_id, platform, os_version, device_model, user_agent, request_ip)
            VALUES (?, ?, ?, ?, ?, ?, ?)
            ON CONFLICT (hwid, user_id)
            DO UPDATE SET
                platform = EXCLUDED.platform,
                os_version = EXCLUDED.os_version,
                device_model = EXCLUDED.device_model,
                user_agent = EXCLUDED.user_agent,
                request_ip = COALESCE(EXCLUDED.request_ip, hwid_user_devices.request_ip),
                updated_at = now()
        `, hwid.Hwid, userID, hwid.Platform, hwid.OsVersion, hwid.DeviceModel, hwid.UserAgent, hwid.RequestIP)
		return err
	})
}

func enqueueOrUpsertHwidUserDevice(ctx context.Context, manager *dbmanager.DatabaseManager, userID int64, hwid HwidHeaders) error {
	hwid.Platform = lowerStringPtr(hwid.Platform)
	queued, err := jobqueue.EnqueueUpsertHwidDevice(ctx, jobqueue.UpsertHwidDevicePayload{
		UserID:      userID,
		Hwid:        hwid.Hwid,
		Platform:    hwid.Platform,
		OsVersion:   hwid.OsVersion,
		DeviceModel: hwid.DeviceModel,
		UserAgent:   hwid.UserAgent,
		RequestIP:   hwid.RequestIP,
	})
	if err == nil && queued {
		return nil
	}
	return upsertHwidUserDevice(ctx, manager, userID, hwid)
}

func updateSubscriptionRequest(ctx context.Context, manager *dbmanager.DatabaseManager, userUUID string, userID int64, userAgent, requestIP string) {
	updateQueued, updateErr := jobqueue.EnqueueUpdateUserSubscription(ctx, jobqueue.UpdateUserSubscriptionPayload{
		UserUUID:  userUUID,
		UserAgent: userAgent,
	})
	recordQueued, recordErr := jobqueue.EnqueueAddSubscriptionRequestRecord(ctx, jobqueue.AddSubscriptionRequestRecordPayload{
		UserID:    userID,
		RequestIP: requestIP,
		UserAgent: userAgent,
	})
	if updateErr == nil && recordErr == nil && updateQueued && recordQueued {
		return
	}

	go func() {
		jobCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		_ = manager.ExecuteLowPriority(func(db dbmanager.DBExecutor) error {
			if _, err := db.ExecContext(jobCtx, `
	            INSERT INTO user_subscription_request_history (user_id, request_ip, user_agent)
	            VALUES (?, ?, ?)
	        `, userID, requestIP, userAgent); err != nil {
				return err
			}

			_, err := db.ExecContext(jobCtx, `
	            DELETE FROM user_subscription_request_history
	            WHERE user_id = ?
	              AND id NOT IN (
                  SELECT id
                  FROM user_subscription_request_history
                  WHERE user_id = ?
                  ORDER BY request_at DESC, id DESC
	                  LIMIT 24
	              )
	        `, userID, userID)
			return err
		})
	}()
}

func isJsonSubscriptionSupported(userAgent string) bool {
	if userAgent == "" {
		return false
	}
	clients := []*regexp.Regexp{
		regexp.MustCompile(`^[Ss]treisand`),
		regexp.MustCompile(`^Happ/`),
		regexp.MustCompile(`^ktor-client`),
		regexp.MustCompile(`^V2Box`),
		regexp.MustCompile(`^io\.github\.saeeddev94\.xray/`),
		regexp.MustCompile(`^v2rayNG/(\d+\.\d+\.\d+)`),
		regexp.MustCompile(`^v2rayN/(\d+\.\d+\.\d+)`),
	}
	for _, re := range clients {
		if re.MatchString(userAgent) {
			return true
		}
	}
	return false
}

func resolveSubscriptionURL(settings SubscriptionSettingsParsed, shortUUID string) string {
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

func buildSubscriptionHeaders(user SubscriptionUser, settings SubscriptionSettingsParsed, isHapp bool) map[string]string {
	headers := map[string]string{}
	subscriptionURL := resolveSubscriptionURL(settings, user.ShortUUID)

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

func filterHostsForResponseType(hosts []SubscriptionHost, responseType string, includeDisabled bool) []SubscriptionHost {
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

func generateSubscriptionContent(responseType string, templateData []byte, hosts []SubscriptionHost, user SubscriptionUser) (SubscriptionWithConfig, error) {
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

func shuffleHostsIfNeeded(hosts []SubscriptionHost, settings SubscriptionSettingsParsed) {
	if settings.Raw.RandomizeHosts {
		rand.Shuffle(len(hosts), func(i, j int) {
			hosts[i], hosts[j] = hosts[j], hosts[i]
		})
	}
}

func getSubscriptionTemplate(ctx context.Context, manager *dbmanager.DatabaseManager, templateType string) ([]byte, error) {
	var templateData []byte
	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		row := db.QueryRowContext(ctx, `
            SELECT template_yaml, template_json
            FROM subscription_templates
            WHERE template_type = ?
            ORDER BY view_position ASC
            LIMIT 1
        `, templateType)

		var templateYAML sql.NullString
		var templateJSON sql.NullString
		if err := row.Scan(&templateYAML, &templateJSON); err != nil {
			return err
		}

		if templateType == responseTypeXrayJSON || templateType == responseTypeSingbox {
			if templateJSON.Valid {
				templateData = []byte(templateJSON.String)
			}
		} else {
			if templateYAML.Valid {
				templateData = []byte(templateYAML.String)
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return templateData, nil
}

func getSubscriptionTemplateByName(ctx context.Context, manager *dbmanager.DatabaseManager, name string) (string, []byte, error) {
	var templateType string
	var templateData []byte
	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		row := db.QueryRowContext(ctx, `
            SELECT template_type, template_yaml, template_json
            FROM subscription_templates
            WHERE name = ?
            LIMIT 1
        `, name)

		var templateYAML sql.NullString
		var templateJSON sql.NullString
		if err := row.Scan(&templateType, &templateYAML, &templateJSON); err != nil {
			return err
		}

		if templateType == responseTypeXrayJSON || templateType == responseTypeSingbox {
			if templateJSON.Valid {
				templateData = []byte(templateJSON.String)
			}
		} else {
			if templateYAML.Valid {
				templateData = []byte(templateYAML.String)
			}
		}
		return nil
	})
	if err != nil {
		return "", nil, err
	}
	return templateType, templateData, nil
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

func buildSubscriptionInfoResponse(user SubscriptionUser, settings SubscriptionSettingsParsed, hosts []SubscriptionHost) SubscriptionInfoResponse {
	filtered := filterHostsForResponseType(hosts, "", false)
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
		SubscriptionURL: resolveSubscriptionURL(settings, user.ShortUUID),
	}
}

func getUsersWithPagination(ctx context.Context, manager *dbmanager.DatabaseManager, start, size int) ([]SubscriptionUser, int, error) {
	users := []SubscriptionUser{}
	total := 0

	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		if err := db.QueryRowContext(ctx, `SELECT COUNT(*) FROM users`).Scan(&total); err != nil {
			return err
		}

		rows, err := db.QueryContext(ctx, `
            SELECT u.t_id, u.uuid, u.short_uuid, u.username, u.status,
                   u.traffic_limit_bytes, u.traffic_limit_strategy, u.expire_at,
	                   u.trojan_password, u.vless_uuid, u.ss_password,
	                   u.naive_password, u.shadowtls_password, u.hysteria2_password, u.anytls_password,
	                   u.hwid_device_limit, u.external_squad_uuid,
                   COALESCE(ut.used_traffic_bytes, 0), COALESCE(ut.lifetime_used_traffic_bytes, 0)
            FROM users u
            LEFT JOIN user_traffic ut ON ut.t_id = u.t_id
            ORDER BY u.created_at DESC
            LIMIT ? OFFSET ?
        `, size, start)
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			var user SubscriptionUser
			var hwidDeviceLimit sql.NullInt64
			var externalSquadUUID sql.NullString
			var naivePassword, shadowtlsPassword, hysteria2Password, anytlsPassword sql.NullString

			if err := rows.Scan(
				&user.TID,
				&user.UUID,
				&user.ShortUUID,
				&user.Username,
				&user.Status,
				&user.TrafficLimitBytes,
				&user.TrafficLimitStrategy,
				&user.ExpireAt,
				&user.TrojanPassword,
				&user.VlessUUID,
				&user.SSPassword,
				&naivePassword,
				&shadowtlsPassword,
				&hysteria2Password,
				&anytlsPassword,
				&hwidDeviceLimit,
				&externalSquadUUID,
				&user.UsedTrafficBytes,
				&user.LifetimeUsedBytes,
			); err != nil {
				return err
			}

			if hwidDeviceLimit.Valid {
				v := int(hwidDeviceLimit.Int64)
				user.HwidDeviceLimit = &v
			}
			if externalSquadUUID.Valid {
				v := externalSquadUUID.String
				user.ExternalSquadUUID = &v
			}
			user.NaivePassword = nullableSQLString(naivePassword)
			user.ShadowtlsPassword = nullableSQLString(shadowtlsPassword)
			user.Hysteria2Password = nullableSQLString(hysteria2Password)
			user.AnytlsPassword = nullableSQLString(anytlsPassword)

			users = append(users, user)
		}

		return rows.Err()
	})

	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func getSubpageConfigForUser(ctx context.Context, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig, shortUUID string, requestHeaders map[string]string) (string, bool, error) {
	log := cfg.Logger.RoleService(logger.RoleAPI, logger.ServiceHTTP)
	user, err := getSubscriptionUserByShortUUID(ctx, manager, shortUUID)
	if err != nil {
		return "", false, err
	}

	subpageConfigUUID := ""

	if user.ExternalSquadUUID != nil {
		var squadCustomRemarks sql.NullString
		var squadIsHwidLimited bool
		var squadHwidMaxDevices int

		err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
			return db.QueryRowContext(ctx, `
			SELECT custom_remarks, is_hwid_limited, hwid_max_devices 
			FROM external_squads 
			WHERE uuid = ?`,
				*user.ExternalSquadUUID).Scan(&squadCustomRemarks, squadIsHwidLimited, squadHwidMaxDevices)
		})

		if err == nil {
			log.Debug("Applied external squad overrides: is_hwid_limited=%v hwid_max_devices=%d", squadIsHwidLimited, squadHwidMaxDevices)

			_ = squadCustomRemarks

			log.Debug(fmt.Sprintf(" Applied external squad overrides: is_hwid_limited=%v hwid_max_devices=%d", squadIsHwidLimited, squadHwidMaxDevices))
		} else {
			log.Error(fmt.Sprintf("Failed to load external squad overrides: %v", err))
		}
	}

	if subpageConfigUUID == "" {
		subpageConfigUUID = defaultSubpageConfigUUID
	}

	settings, err := loadSubscriptionSettings(ctx, manager, cfg)
	if err != nil {
		return subpageConfigUUID, false, err
	}

	webpageAllowed := false
	if settings.ResponseRules != nil {
		header := http.Header{}
		for key, value := range requestHeaders {
			header.Set(key, value)
		}
		responseType := matchResponseRules(settings.ResponseRules, header)
		webpageAllowed = responseType == responseTypeBrowser
	}

	return subpageConfigUUID, webpageAllowed, nil
}

func UpdateExternalSquad(ctx context.Context, manager *dbmanager.DatabaseManager, squadUUID string, input UpdateExternalSquadInput) error {
	var currentName string
	var currentSubpageConfigUUID sql.NullString
	var currentCustomRemarks sql.NullString
	var currentHwidSettingsRaw []byte

	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		return db.QueryRowContext(ctx,
			`SELECT name, subpage_config_uuid, custom_remarks, hwid_settings FROM external_squads WHERE uuid = ?`,
			squadUUID).Scan(&currentName, &currentSubpageConfigUUID, &currentCustomRemarks, &currentHwidSettingsRaw)
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("external squad not found")
		}
		return fmt.Errorf("failed to fetch current external squad: %w", err)
	}

	var columns []string
	var args []interface{}

	if input.Name != nil {
		columns = append(columns, "name = ?")
		args = append(args, *input.Name)
	}

	if input.SubpageConfigUUID != nil {
		columns = append(columns, "subpage_config_uuid = ?")
		if *input.SubpageConfigUUID == "" {
			args = append(args, nil)
		} else {
			args = append(args, *input.SubpageConfigUUID)
		}
	}

	if input.CustomRemarks != nil {
		columns = append(columns, "custom_remarks = ?")
		if len(*input.CustomRemarks) == 0 || string(*input.CustomRemarks) == "null" {
			args = append(args, nil)
		} else {
			args = append(args, string(*input.CustomRemarks))
		}
	}

	if len(input.HwidSettings) > 0 {
		raw := strings.TrimSpace(string(input.HwidSettings))
		if raw == "null" {
			columns = append(columns, "hwid_settings = ?")
			args = append(args, nil)
		} else {
			var hwidInput HwidSettingsInput
			if err := json.Unmarshal(input.HwidSettings, &hwidInput); err != nil {
				return fmt.Errorf("invalid hwidSettings: %w", err)
			}

			var hwid struct {
				Enabled             bool `json:"enabled"`
				MaxDevicesAnnounce  *int `json:"maxDevicesAnnounce"`
				FallbackDeviceLimit int  `json:"fallbackDeviceLimit"`
			}
			hwid.FallbackDeviceLimit = 999
			if len(currentHwidSettingsRaw) > 0 && !bytes.Equal(currentHwidSettingsRaw, []byte("null")) {
				if err := json.Unmarshal(currentHwidSettingsRaw, &hwid); err != nil {
					return fmt.Errorf("failed to unmarshal current hwid settings from database: %w", err)
				}
			}

			if hwidInput.Enabled != nil {
				hwid.Enabled = *hwidInput.Enabled
			}
			if hwidInput.FallbackDeviceLimit != nil {
				hwid.FallbackDeviceLimit = *hwidInput.FallbackDeviceLimit
			}
			if hwidInput.MaxDevicesAnnounce != nil {
				hwid.MaxDevicesAnnounce = hwidInput.MaxDevicesAnnounce
			}

			updatedHwidRaw, err := json.Marshal(hwid)
			if err != nil {
				return fmt.Errorf("failed to marshal merged hwid settings: %w", err)
			}

			columns = append(columns, "hwid_settings = ?")
			args = append(args, string(updatedHwidRaw))
		}
	}

	if len(columns) == 0 {
		return fmt.Errorf("no fields to update")
	}

	args = append(args, squadUUID)
	query := fmt.Sprintf("UPDATE external_squads SET %s WHERE uuid = ?", strings.Join(columns, ", "))

	err = manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		_, err := db.ExecContext(ctx, query, args...)
		return err
	})
	return err
}
