package subscription

import (
	"context"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"regexp"
	"sort"
	"strings"
	"time"

	"cerberus/backend/config"
	dbmanager "cerberus/backend/db/manager"
	"cerberus/backend/dbutil"
	"cerberus/backend/httpapi/shared"
	subscriptionresponserules "cerberus/backend/httpapi/subscription-response-rules"
	subscriptionsettings "cerberus/backend/httpapi/subscription-settings"

	"gopkg.in/yaml.v3"
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

// SubscriptionSettingsParsed contains subscription settings with parsed JSON fields.
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
	HwidDeviceLimit      *int
	ExternalSquadUUID    *string
	SubLastUserAgent     *string
	SubLastOpenedAt      *time.Time
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
	SockoptParams                *string
	IsDisabled                   bool
	ServerDescription            *string
	VLESSRouteID                 *int64
	AllowInsecure                bool
	ShuffleHost                  bool
	MihomoX25519                 bool
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
	Platform    string
	OsVersion   string
	DeviceModel string
	UserAgent   string
}

// loadSubscriptionSettings loads and parses subscription settings.
func loadSubscriptionSettings(ctx context.Context, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) (SubscriptionSettingsParsed, error) {
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

	// Parse JSON fields
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
                   u.hwid_device_limit, u.external_squad_uuid,
                   u.sub_last_user_agent, u.sub_last_opened_at,
                   COALESCE(ut.used_traffic_bytes, 0), COALESCE(ut.lifetime_used_traffic_bytes, 0)
            FROM users u
            LEFT JOIN user_traffic ut ON ut.t_id = u.t_id
            WHERE %s
            LIMIT 1
        `, where)

		row := db.QueryRowContext(ctx, query, value)

		var hwidDeviceLimit sql.NullInt64
		var externalSquadUUID sql.NullString
		var subLastUserAgent sql.NullString
		var subLastOpenedAt sql.NullTime
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
			&hwidDeviceLimit,
			&externalSquadUUID,
			&subLastUserAgent,
			&subLastOpenedAt,
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
		if subLastUserAgent.Valid {
			v := subLastUserAgent.String
			user.SubLastUserAgent = &v
		}
		if subLastOpenedAt.Valid {
			v := subLastOpenedAt.Time
			user.SubLastOpenedAt = &v
		}
		return nil
	})
	if err != nil {
		return user, err
	}

	return user, nil
}

func getUserSquadUUIDs(ctx context.Context, manager *dbmanager.DatabaseManager, userTID int64) ([]string, error) {
	var squads []string
	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		rows, err := db.QueryContext(ctx, `
            SELECT internal_squad_uuid
            FROM internal_squad_members
            WHERE user_id = ?
        `, userTID)
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			var uuid string
			if err := rows.Scan(&uuid); err != nil {
				return err
			}
			squads = append(squads, uuid)
		}
		return rows.Err()
	})
	if err != nil {
		return nil, err
	}
	return squads, nil
}

func getHostsForUser(ctx context.Context, manager *dbmanager.DatabaseManager, user SubscriptionUser) ([]SubscriptionHost, error) {
	squads, err := getUserSquadUUIDs(ctx, manager, user.TID)
	if err != nil {
		return nil, err
	}

	var hosts []SubscriptionHost

	err = manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		var rows *sql.Rows
		if len(squads) == 0 {
			rows, err = db.QueryContext(ctx, `
                SELECT h.uuid, h.view_position, h.remark, h.address, h.port,
                       h.path, h.sni, h.host, h.alpn, h.fingerprint, h.security_layer,
                       h.xhttp_extra_params, h.mux_params, h.sockopt_params, h.is_disabled,
                       h.server_description, h.vless_route_id, h.allow_insecure, h.shuffle_host,
                       h.mihomo_x25519, h.xray_json_template_uuid, h.keep_sni_blank,
                       h.exclude_from_subscription_types, h.tag, h.is_hidden, h.override_sni_from_address,
                       h.config_profile_uuid, h.config_profile_inbound_uuid,
                       cpi.tag, cpi.type, cpi.network, cpi.security, cpi.port, cpi.raw_inbound
                FROM hosts h
                LEFT JOIN config_profile_inbounds cpi ON h.config_profile_inbound_uuid = cpi.uuid
                ORDER BY h.view_position ASC, h.remark ASC
            `)
		} else {
			rows, err = db.QueryContext(ctx, `
                SELECT DISTINCT h.uuid, h.view_position, h.remark, h.address, h.port,
                       h.path, h.sni, h.host, h.alpn, h.fingerprint, h.security_layer,
                       h.xhttp_extra_params, h.mux_params, h.sockopt_params, h.is_disabled,
                       h.server_description, h.vless_route_id, h.allow_insecure, h.shuffle_host,
                       h.mihomo_x25519, h.xray_json_template_uuid, h.keep_sni_blank,
                       h.exclude_from_subscription_types, h.tag, h.is_hidden, h.override_sni_from_address,
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
		}
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

func scanSubscriptionHost(scanner shared.RowScanner) (SubscriptionHost, error) {
	var h SubscriptionHost
	var viewPosition sql.NullInt64
	var path, sni, host, alpn, fingerprint, securityLayer sql.NullString
	var xhttpExtraParams, muxParams, sockoptParams, serverDescription sql.NullString
	var vlessRouteID sql.NullInt64
	var xrayJSONTemplateUUID, tag, configProfileUUID, configProfileInboundUUID sql.NullString
	var inboundTag, inboundType, inboundNetwork, inboundSecurity sql.NullString
	var inboundPort sql.NullInt64
	var rawInbound sql.NullString
	var excludeTypes dbutil.StringArray
	var isDisabled, allowInsecure, shuffleHost, mihomoX25519, keepSNIBlank, isHidden, overrideSNIFromAddress sql.NullBool

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
		&sockoptParams,
		&isDisabled,
		&serverDescription,
		&vlessRouteID,
		&allowInsecure,
		&shuffleHost,
		&mihomoX25519,
		&xrayJSONTemplateUUID,
		&keepSNIBlank,
		&excludeTypes,
		&tag,
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
	if sockoptParams.Valid {
		h.SockoptParams = &sockoptParams.String
	}
	if isDisabled.Valid {
		h.IsDisabled = isDisabled.Bool
	}
	if serverDescription.Valid {
		h.ServerDescription = &serverDescription.String
	}
	if vlessRouteID.Valid {
		h.VLESSRouteID = &vlessRouteID.Int64
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
	if xrayJSONTemplateUUID.Valid {
		h.XrayJSONTemplateUUID = &xrayJSONTemplateUUID.String
	}
	if keepSNIBlank.Valid {
		h.KeepSNIBlank = keepSNIBlank.Bool
	}
	if tag.Valid {
		h.Tag = &tag.String
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

	return &HwidHeaders{
		Hwid:        hwid,
		Platform:    strings.TrimSpace(r.Header.Get("X-HWID-Platform")),
		OsVersion:   strings.TrimSpace(r.Header.Get("X-HWID-OS-Version")),
		DeviceModel: strings.TrimSpace(r.Header.Get("X-HWID-Device-Model")),
		UserAgent:   strings.TrimSpace(r.Header.Get("X-HWID-User-Agent")),
	}
}

func checkHwidDeviceLimit(ctx context.Context, manager *dbmanager.DatabaseManager, user SubscriptionUser, hwid *HwidHeaders, settings HwidSettings) (bool, bool, bool) {
	// returns (isAllowed, maxDeviceReached, hwidNotSupported)
	if user.HwidDeviceLimit != nil && *user.HwidDeviceLimit == 0 {
		if hwid != nil {
			_ = upsertHwidUserDevice(ctx, manager, user.UUID, *hwid)
		}
		return true, false, false
	}

	if hwid == nil {
		return false, false, true
	}

	exists, err := hwidDeviceExists(ctx, manager, user.UUID, hwid.Hwid)
	if err == nil && exists {
		_ = upsertHwidUserDevice(ctx, manager, user.UUID, *hwid)
		return true, false, false
	}

	count, err := countHwidDevices(ctx, manager, user.UUID)
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

	if err := upsertHwidUserDevice(ctx, manager, user.UUID, *hwid); err != nil {
		return false, true, false
	}

	return true, false, false
}

func countHwidDevices(ctx context.Context, manager *dbmanager.DatabaseManager, userUUID string) (int, error) {
	var count int
	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		return db.QueryRowContext(ctx, `SELECT COUNT(*) FROM hwid_user_devices WHERE user_uuid = ?`, userUUID).Scan(&count)
	})
	if err != nil {
		return 0, err
	}
	return count, nil
}

func hwidDeviceExists(ctx context.Context, manager *dbmanager.DatabaseManager, userUUID, hwid string) (bool, error) {
	var exists bool
	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		var tmp int
		err := db.QueryRowContext(ctx, `SELECT 1 FROM hwid_user_devices WHERE user_uuid = ? AND hwid = ?`, userUUID, hwid).Scan(&tmp)
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

func upsertHwidUserDevice(ctx context.Context, manager *dbmanager.DatabaseManager, userUUID string, hwid HwidHeaders) error {
	return manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		_, err := db.ExecContext(ctx, `
            INSERT INTO hwid_user_devices (hwid, user_uuid, platform, os_version, device_model, user_agent)
            VALUES (?, ?, ?, ?, ?, ?)
            ON CONFLICT (hwid, user_uuid)
            DO UPDATE SET platform = EXCLUDED.platform, os_version = EXCLUDED.os_version, device_model = EXCLUDED.device_model, user_agent = EXCLUDED.user_agent, updated_at = now()
        `, hwid.Hwid, userUUID, hwid.Platform, hwid.OsVersion, hwid.DeviceModel, hwid.UserAgent)
		return err
	})
}

func updateSubscriptionRequest(ctx context.Context, manager *dbmanager.DatabaseManager, userUUID, userAgent, requestIP string) {
	_ = manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		_, _ = db.ExecContext(ctx, `
            UPDATE users SET sub_last_opened_at = now(), sub_last_user_agent = ? WHERE uuid = ?
        `, userAgent, userUUID)

		_, _ = db.ExecContext(ctx, `
            INSERT INTO user_subscription_request_history (user_uuid, request_ip, user_agent)
            VALUES (?, ?, ?)
        `, userUUID, requestIP, userAgent)

		return nil
	})
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
		regexp.MustCompile(`^io\\.github\\.saeeddev94\\.xray/`),
		regexp.MustCompile(`^v2rayNG/(\\d+\\.\\d+\\.\\d+)`),
		regexp.MustCompile(`^v2rayN/(\\d+\\.\\d+\\.\\d+)`),
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
	now := time.Now()
	switch strings.ToUpper(strategy) {
	case "DAY":
		now = now.AddDate(0, 0, 1)
		now = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		return fmt.Sprintf("%d", now.Unix())
	case "WEEK":
		offset := 8 - int(now.Weekday())
		now = now.AddDate(0, 0, offset)
		now = time.Date(now.Year(), now.Month(), now.Day(), 0, 5, 0, 0, now.Location())
		return fmt.Sprintf("%d", now.Unix())
	case "MONTH":
		now = time.Date(now.Year(), now.Month(), 1, 0, 10, 0, 0, now.Location())
		now = now.AddDate(0, 1, 0)
		return fmt.Sprintf("%d", now.Unix())
	default:
		return ""
	}
}

func formatTemplateValue(value string, user SubscriptionUser, settings SubscriptionSettingsParsed, subscriptionURL string) string {
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

func buildRawHost(host SubscriptionHost) RawHost {
	protocol := ""
	if host.InboundType != nil {
		protocol = *host.InboundType
	}

	return RawHost{
		UUID:          host.UUID,
		Remark:        host.Remark,
		Address:       host.Address,
		Port:          host.Port,
		Protocol:      protocol,
		Network:       host.InboundNetwork,
		Security:      host.InboundSecurity,
		Path:          host.Path,
		SNI:           host.SNI,
		Host:          host.Host,
		ALPN:          host.ALPN,
		Fingerprint:   host.Fingerprint,
		AllowInsecure: host.AllowInsecure,
		IsDisabled:    host.IsDisabled,
		IsHidden:      host.IsHidden,
	}
}

func buildSubscriptionLinks(hosts []SubscriptionHost, user SubscriptionUser) ([]string, map[string]string) {
	links := []string{}
	ssConfLinks := map[string]string{}

	for _, host := range hosts {
		link, protocol := buildHostLink(host, user)
		if link == "" {
			continue
		}
		links = append(links, link)

		if protocol == "shadowsocks" || protocol == "ss" {
			remark := host.Remark
			if remark == "" {
				remark = host.Address
			}
			encoded := base64.RawURLEncoding.EncodeToString([]byte(remark))
			domain := host.Address
			ssConfLinks[remark] = fmt.Sprintf("ssconf://%s/%s/ss/%s#%s", domain, user.ShortUUID, encoded, url.QueryEscape(remark))
		}
	}

	return links, ssConfLinks
}

func buildHostLink(host SubscriptionHost, user SubscriptionUser) (string, string) {
	protocol := ""
	if host.InboundType != nil {
		protocol = strings.ToLower(*host.InboundType)
	}
	if protocol == "" {
		return "", ""
	}

	switch protocol {
	case "vless":
		return buildVlessLink(host, user), protocol
	case "trojan":
		return buildTrojanLink(host, user), protocol
	case "shadowsocks", "ss":
		return buildShadowsocksLink(host, user), protocol
	case "vmess":
		return buildVmessLink(host, user), protocol
	default:
		return "", protocol
	}
}

func buildVlessLink(host SubscriptionHost, user SubscriptionUser) string {
	if user.VlessUUID == "" {
		return ""
	}
	params := url.Values{}
	params.Set("encryption", "none")
	applyTransportParams(&params, host)

	remark := host.Remark
	if remark == "" {
		remark = host.Address
	}
	return fmt.Sprintf("vless://%s@%s:%d?%s#%s", user.VlessUUID, host.Address, host.Port, params.Encode(), url.QueryEscape(remark))
}

func buildTrojanLink(host SubscriptionHost, user SubscriptionUser) string {
	if user.TrojanPassword == "" {
		return ""
	}
	params := url.Values{}
	applyTransportParams(&params, host)

	remark := host.Remark
	if remark == "" {
		remark = host.Address
	}
	return fmt.Sprintf("trojan://%s@%s:%d?%s#%s", url.QueryEscape(user.TrojanPassword), host.Address, host.Port, params.Encode(), url.QueryEscape(remark))
}

func buildShadowsocksLink(host SubscriptionHost, user SubscriptionUser) string {
	method := extractShadowsocksMethod(host.InboundRaw)
	if method == "" {
		method = "aes-128-gcm"
	}
	if user.SSPassword == "" {
		return ""
	}

	creds := fmt.Sprintf("%s:%s", method, user.SSPassword)
	encoded := base64.RawURLEncoding.EncodeToString([]byte(creds))

	remark := host.Remark
	if remark == "" {
		remark = host.Address
	}
	return fmt.Sprintf("ss://%s@%s:%d#%s", encoded, host.Address, host.Port, url.QueryEscape(remark))
}

func buildVmessLink(host SubscriptionHost, user SubscriptionUser) string {
	// VMess is not supported yet in this implementation.
	return ""
}

func applyTransportParams(params *url.Values, host SubscriptionHost) {
	network := "tcp"
	if host.InboundNetwork != nil && *host.InboundNetwork != "" {
		network = *host.InboundNetwork
	}

	security := "none"
	if host.InboundSecurity != nil && *host.InboundSecurity != "" {
		security = *host.InboundSecurity
	} else {
		switch strings.ToUpper(host.SecurityLayer) {
		case "TLS":
			security = "tls"
		case "NONE":
			security = "none"
		}
	}

	params.Set("type", network)
	if security != "" {
		params.Set("security", security)
	}

	sni := ""
	if host.SNI != nil {
		sni = *host.SNI
	}
	if sni == "" && host.OverrideSNIFromAddress {
		sni = host.Address
	}

	if sni != "" {
		params.Set("sni", sni)
	}

	if host.ALPN != nil && *host.ALPN != "" {
		params.Set("alpn", *host.ALPN)
	}
	if host.Fingerprint != nil && *host.Fingerprint != "" {
		params.Set("fp", *host.Fingerprint)
	}
	if host.AllowInsecure {
		params.Set("allowInsecure", "1")
	}

	if host.Path != nil && *host.Path != "" {
		params.Set("path", *host.Path)
	}
	if host.Host != nil && *host.Host != "" {
		params.Set("host", *host.Host)
	}
}

func extractShadowsocksMethod(raw json.RawMessage) string {
	if len(raw) == 0 {
		return ""
	}

	var obj map[string]interface{}
	if err := json.Unmarshal(raw, &obj); err != nil {
		return ""
	}

	if settings, ok := obj["settings"].(map[string]interface{}); ok {
		if method, ok := settings["method"].(string); ok {
			return method
		}
	}
	if method, ok := obj["method"].(string); ok {
		return method
	}

	return ""
}

func generateXrayJSONConfig(templateJSON []byte, hosts []SubscriptionHost, user SubscriptionUser) (string, error) {
	baseConfig := map[string]interface{}{}
	if len(templateJSON) > 0 {
		if err := json.Unmarshal(templateJSON, &baseConfig); err != nil {
			baseConfig = map[string]interface{}{}
		}
	}

	outbounds := []interface{}{}
	if existing, ok := baseConfig["outbounds"].([]interface{}); ok {
		outbounds = existing
	}

	for _, host := range hosts {
		outbound := buildXrayOutbound(host, user)
		if outbound != nil {
			outbounds = append(outbounds, outbound)
		}
	}

	baseConfig["outbounds"] = outbounds

	data, err := json.MarshalIndent(baseConfig, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func buildXrayOutbound(host SubscriptionHost, user SubscriptionUser) map[string]interface{} {
	protocol := ""
	if host.InboundType != nil {
		protocol = strings.ToLower(*host.InboundType)
	}
	if protocol == "" {
		return nil
	}

	remark := host.Remark
	if remark == "" {
		remark = host.Address
	}

	network := "tcp"
	if host.InboundNetwork != nil && *host.InboundNetwork != "" {
		network = *host.InboundNetwork
	}

	security := "none"
	if host.InboundSecurity != nil && *host.InboundSecurity != "" {
		security = *host.InboundSecurity
	} else {
		switch strings.ToUpper(host.SecurityLayer) {
		case "TLS":
			security = "tls"
		case "NONE":
			security = "none"
		}
	}

	streamSettings := map[string]interface{}{
		"network":  network,
		"security": security,
	}

	sni := ""
	if host.SNI != nil {
		sni = *host.SNI
	}
	if sni == "" && host.OverrideSNIFromAddress {
		sni = host.Address
	}

	if security == "tls" {
		tlsSettings := map[string]interface{}{
			"allowInsecure": host.AllowInsecure,
		}
		if sni != "" {
			tlsSettings["serverName"] = sni
		}
		if host.ALPN != nil && *host.ALPN != "" {
			tlsSettings["alpn"] = strings.Split(*host.ALPN, ",")
		}
		if host.Fingerprint != nil && *host.Fingerprint != "" {
			tlsSettings["fingerprint"] = *host.Fingerprint
		}
		streamSettings["tlsSettings"] = tlsSettings
	}

	if network == "ws" {
		wsSettings := map[string]interface{}{}
		if host.Path != nil && *host.Path != "" {
			wsSettings["path"] = *host.Path
		}
		if host.Host != nil && *host.Host != "" {
			wsSettings["headers"] = map[string]interface{}{"Host": *host.Host}
		}
		streamSettings["wsSettings"] = wsSettings
	}

	if network == "grpc" {
		grpcSettings := map[string]interface{}{}
		if host.Path != nil && *host.Path != "" {
			grpcSettings["serviceName"] = *host.Path
		}
		streamSettings["grpcSettings"] = grpcSettings
	}

	outbound := map[string]interface{}{
		"tag":            remark,
		"protocol":       protocol,
		"streamSettings": streamSettings,
	}

	switch protocol {
	case "vless":
		outbound["settings"] = map[string]interface{}{
			"vnext": []interface{}{
				map[string]interface{}{
					"address": host.Address,
					"port":    host.Port,
					"users": []interface{}{
						map[string]interface{}{
							"id":         user.VlessUUID,
							"encryption": "none",
						},
					},
				},
			},
		}
	case "trojan":
		outbound["settings"] = map[string]interface{}{
			"servers": []interface{}{
				map[string]interface{}{
					"address":  host.Address,
					"port":     host.Port,
					"password": user.TrojanPassword,
				},
			},
		}
	case "shadowsocks", "ss":
		method := extractShadowsocksMethod(host.InboundRaw)
		if method == "" {
			method = "aes-128-gcm"
		}
		outbound["protocol"] = "shadowsocks"
		outbound["settings"] = map[string]interface{}{
			"servers": []interface{}{
				map[string]interface{}{
					"address":  host.Address,
					"port":     host.Port,
					"method":   method,
					"password": user.SSPassword,
				},
			},
		}
	default:
		return nil
	}

	return outbound
}

func generateYAMLConfig(templateYAML []byte, hosts []SubscriptionHost, user SubscriptionUser) (string, error) {
	var config map[string]interface{}
	if len(templateYAML) > 0 {
		if err := yaml.Unmarshal(templateYAML, &config); err != nil {
			config = map[string]interface{}{}
		}
	} else {
		config = map[string]interface{}{}
	}

	proxies := []interface{}{}
	proxyNames := []string{}
	for _, host := range hosts {
		proxy := buildMihomoProxy(host, user)
		if proxy == nil {
			continue
		}
		proxies = append(proxies, proxy)
		if name, ok := proxy["name"].(string); ok && name != "" {
			proxyNames = append(proxyNames, name)
		}
	}

	config["proxies"] = proxies

	if groups, ok := config["proxy-groups"].([]interface{}); ok {
		for _, group := range groups {
			groupMap, ok := group.(map[string]interface{})
			if !ok {
				continue
			}
			if proxiesVal, exists := groupMap["proxies"]; exists {
				switch v := proxiesVal.(type) {
				case []interface{}:
					if len(v) == 0 {
						groupMap["proxies"] = proxyNames
					}
				case string, nil:
					groupMap["proxies"] = proxyNames
				default:
					// leave as-is
				}
			}
		}
	}

	data, err := yaml.Marshal(config)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func buildMihomoProxy(host SubscriptionHost, user SubscriptionUser) map[string]interface{} {
	protocol := ""
	if host.InboundType != nil {
		protocol = strings.ToLower(*host.InboundType)
	}
	if protocol == "" {
		return nil
	}

	name := host.Remark
	if name == "" {
		name = host.Address
	}

	proxy := map[string]interface{}{
		"name":   name,
		"type":   protocol,
		"server": host.Address,
		"port":   host.Port,
		"udp":    true,
	}

	network := "tcp"
	if host.InboundNetwork != nil && *host.InboundNetwork != "" {
		network = *host.InboundNetwork
	}

	if protocol == "vless" {
		proxy["uuid"] = user.VlessUUID
	} else if protocol == "trojan" {
		proxy["password"] = user.TrojanPassword
	} else if protocol == "shadowsocks" || protocol == "ss" {
		proxy["type"] = "ss"
		proxy["password"] = user.SSPassword
		method := extractShadowsocksMethod(host.InboundRaw)
		if method == "" {
			method = "aes-128-gcm"
		}
		proxy["cipher"] = method
	}

	security := "none"
	if host.InboundSecurity != nil && *host.InboundSecurity != "" {
		security = *host.InboundSecurity
	} else {
		switch strings.ToUpper(host.SecurityLayer) {
		case "TLS":
			security = "tls"
		case "NONE":
			security = "none"
		}
	}

	if security == "tls" {
		proxy["tls"] = true
		proxy["skip-cert-verify"] = host.AllowInsecure
		if host.SNI != nil && *host.SNI != "" {
			proxy["servername"] = *host.SNI
		}
	}

	if host.Fingerprint != nil && *host.Fingerprint != "" {
		proxy["client-fingerprint"] = *host.Fingerprint
	}

	if host.ALPN != nil && *host.ALPN != "" {
		proxy["alpn"] = strings.Split(*host.ALPN, ",")
	}

	if network != "" {
		proxy["network"] = network
	}

	if network == "ws" {
		wsOpts := map[string]interface{}{}
		if host.Path != nil && *host.Path != "" {
			wsOpts["path"] = *host.Path
		}
		headers := map[string]interface{}{}
		if host.Host != nil && *host.Host != "" {
			headers["Host"] = *host.Host
		}
		if len(headers) > 0 {
			wsOpts["headers"] = headers
		}
		if len(wsOpts) > 0 {
			proxy["ws-opts"] = wsOpts
		}
	}

	if network == "grpc" {
		grpcOpts := map[string]interface{}{}
		if host.Path != nil && *host.Path != "" {
			grpcOpts["grpc-service-name"] = *host.Path
		}
		if len(grpcOpts) > 0 {
			proxy["grpc-opts"] = grpcOpts
		}
	}

	return proxy
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
		body, err := generateXrayJSONConfig(templateData, hosts, user)
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
                   u.hwid_device_limit, u.external_squad_uuid,
                   u.sub_last_user_agent, u.sub_last_opened_at,
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
			var subLastUserAgent sql.NullString
			var subLastOpenedAt sql.NullTime

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
				&hwidDeviceLimit,
				&externalSquadUUID,
				&subLastUserAgent,
				&subLastOpenedAt,
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
			if subLastUserAgent.Valid {
				v := subLastUserAgent.String
				user.SubLastUserAgent = &v
			}
			if subLastOpenedAt.Valid {
				v := subLastOpenedAt.Time
				user.SubLastOpenedAt = &v
			}

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
	user, err := getSubscriptionUserByShortUUID(ctx, manager, shortUUID)
	if err != nil {
		return "", false, err
	}

	subpageConfigUUID := ""

	if user.ExternalSquadUUID != nil {
		_ = manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
			return db.QueryRowContext(ctx, `SELECT subpage_config_uuid FROM external_squads WHERE uuid = ?`, *user.ExternalSquadUUID).Scan(&subpageConfigUUID)
		})
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
