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
	"net/url"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"exodus/internal/config"
	dbmanager "exodus/internal/db/manager"
	"exodus/internal/dbutil"
	"exodus/internal/httpapi/shared"
	"exodus/internal/httpapi/subscriptionresponserules"
	"exodus/internal/httpapi/subscriptionsettings"

	"github.com/iancoleman/orderedmap"
	"golang.org/x/crypto/curve25519"
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
	MuxParams                    *string
	SingboxMuxParams             *string
	ClashMuxParams               *string
	SockoptParams                *string
	IsDisabled                   bool
	ServerDescription            *string
	VLESSRouteID                 *int64
	AllowInsecure                bool
	ShuffleHost                  bool
	SelectorNodesFirst           bool
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
                       h.mux_params, h.singbox_mux_params, h.clash_mux_params, h.sockopt_params, h.is_disabled,
                       h.server_description, h.vless_route_id, h.allow_insecure, h.shuffle_host,
                       h.selector_nodes_first, h.mihomo_x25519, h.xray_json_template_uuid, h.keep_sni_blank,
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
                       h.mux_params, h.singbox_mux_params, h.clash_mux_params, h.sockopt_params, h.is_disabled,
                       h.server_description, h.vless_route_id, h.allow_insecure, h.shuffle_host,
                       h.selector_nodes_first, h.mihomo_x25519, h.xray_json_template_uuid, h.keep_sni_blank,
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
	var muxParams, singboxMuxParams, clashMuxParams, sockoptParams, serverDescription sql.NullString
	var vlessRouteID sql.NullInt64
	var xrayJSONTemplateUUID, tag, configProfileUUID, configProfileInboundUUID sql.NullString
	var inboundTag, inboundType, inboundNetwork, inboundSecurity sql.NullString
	var inboundPort sql.NullInt64
	var rawInbound sql.NullString
	var excludeTypes dbutil.StringArray
	var isDisabled, allowInsecure, shuffleHost, selectorNodesFirst, mihomoX25519, keepSNIBlank, isHidden, overrideSNIFromAddress sql.NullBool

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
		&muxParams,
		&singboxMuxParams,
		&clashMuxParams,
		&sockoptParams,
		&isDisabled,
		&serverDescription,
		&vlessRouteID,
		&allowInsecure,
		&shuffleHost,
		&selectorNodesFirst,
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
	if vlessRouteID.Valid {
		h.VLESSRouteID = &vlessRouteID.Int64
	}
	if allowInsecure.Valid {
		h.AllowInsecure = allowInsecure.Bool
	}
	if shuffleHost.Valid {
		h.ShuffleHost = shuffleHost.Bool
	}
	if selectorNodesFirst.Valid {
		h.SelectorNodesFirst = selectorNodesFirst.Bool
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

func parseJSONMapString(raw *string) map[string]interface{} {
	if raw == nil {
		return nil
	}
	trimmed := strings.TrimSpace(*raw)
	if trimmed == "" || trimmed == "null" {
		return nil
	}

	var parsed map[string]interface{}
	if err := json.Unmarshal([]byte(trimmed), &parsed); err == nil {
		if len(parsed) > 0 {
			return parsed
		}
		return nil
	}

	// Stored value can be a JSON string with YAML payload (for Clash/Mihomo editor).
	var yamlPayload string
	if err := json.Unmarshal([]byte(trimmed), &yamlPayload); err != nil {
		return nil
	}
	yamlPayload = strings.TrimSpace(yamlPayload)
	if yamlPayload == "" {
		return nil
	}

	var yamlParsed map[string]interface{}
	if err := yaml.Unmarshal([]byte(yamlPayload), &yamlParsed); err != nil {
		return nil
	}
	if len(yamlParsed) == 0 {
		return nil
	}
	return yamlParsed
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

	return marshalJSONWithTemplateTopLevelOrder(templateJSON, baseConfig)
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

	if mux := parseJSONMapString(host.MuxParams); mux != nil {
		outbound["mux"] = mux
	}

	return outbound
}

func generateSingboxConfig(templateJSON []byte, hosts []SubscriptionHost, user SubscriptionUser) (string, error) {
	baseConfig := orderedmap.New()
	if len(templateJSON) > 0 {
		if err := baseConfig.UnmarshalJSON(templateJSON); err != nil {
			baseConfig = orderedmap.New()
		}
	}

	outbounds := []interface{}{}
	if existing, ok := baseConfig.Get("outbounds"); ok {
		if items, ok := existing.([]interface{}); ok {
			outbounds = items
		}
	}

	leadingSelectorNodeTags := make([]string, 0, len(hosts))
	trailingSelectorNodeTags := make([]string, 0, len(hosts))
	for _, host := range hosts {
		// In exodus hidden hosts are not returned for singbox generation.
		if host.IsHidden {
			continue
		}

		outbound := buildSingboxOutbound(host, user)
		if outbound == nil {
			continue
		}
		outbounds = append(outbounds, outbound)

		tag := orderedMapString(*outbound, "tag")
		if tag == "" {
			continue
		}
		if host.SelectorNodesFirst {
			leadingSelectorNodeTags = append(leadingSelectorNodeTags, tag)
		} else {
			trailingSelectorNodeTags = append(trailingSelectorNodeTags, tag)
		}
	}
	baseConfig.Set("outbounds", outbounds)
	patchSingboxSelectors(baseConfig, leadingSelectorNodeTags, trailingSelectorNodeTags)

	data, err := json.MarshalIndent(baseConfig, "", "  ")
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func marshalJSONWithTemplateTopLevelOrder(templateJSON []byte, payload map[string]interface{}) (string, error) {
	templateOrder, err := extractTopLevelJSONKeys(templateJSON)
	if err != nil || len(templateOrder) == 0 {
		data, marshalErr := json.MarshalIndent(payload, "", "  ")
		if marshalErr != nil {
			return "", marshalErr
		}
		return string(data), nil
	}

	encodedValues := make(map[string][]byte, len(payload))
	for key, value := range payload {
		raw, marshalErr := json.MarshalIndent(value, "", "  ")
		if marshalErr != nil {
			return "", marshalErr
		}
		encodedValues[key] = raw
	}

	orderedKeys := make([]string, 0, len(payload))
	used := make(map[string]struct{}, len(payload))
	for _, key := range templateOrder {
		if _, exists := encodedValues[key]; !exists {
			continue
		}
		if _, exists := used[key]; exists {
			continue
		}
		orderedKeys = append(orderedKeys, key)
		used[key] = struct{}{}
	}

	remainingKeys := make([]string, 0)
	for key := range encodedValues {
		if _, exists := used[key]; !exists {
			remainingKeys = append(remainingKeys, key)
		}
	}
	sort.Strings(remainingKeys)
	orderedKeys = append(orderedKeys, remainingKeys...)

	if len(orderedKeys) == 0 {
		return "{}", nil
	}

	var builder strings.Builder
	builder.WriteString("{\n")

	for i, key := range orderedKeys {
		keyJSON, marshalErr := json.Marshal(key)
		if marshalErr != nil {
			return "", marshalErr
		}

		builder.WriteString("  ")
		builder.Write(keyJSON)
		builder.WriteString(": ")
		builder.WriteString(indentTopLevelJSONValue(encodedValues[key]))

		if i != len(orderedKeys)-1 {
			builder.WriteString(",")
		}
		builder.WriteString("\n")
	}

	builder.WriteString("}")
	return builder.String(), nil
}

func extractTopLevelJSONKeys(templateJSON []byte) ([]string, error) {
	trimmed := bytes.TrimSpace(templateJSON)
	if len(trimmed) == 0 {
		return nil, nil
	}

	decoder := json.NewDecoder(bytes.NewReader(trimmed))
	firstToken, err := decoder.Token()
	if err != nil {
		return nil, err
	}

	startDelim, ok := firstToken.(json.Delim)
	if !ok || startDelim != '{' {
		return nil, nil
	}

	keys := make([]string, 0)
	for decoder.More() {
		keyToken, tokenErr := decoder.Token()
		if tokenErr != nil {
			return nil, tokenErr
		}

		key, ok := keyToken.(string)
		if !ok {
			return nil, fmt.Errorf("invalid json object key token type")
		}
		keys = append(keys, key)

		var discard json.RawMessage
		if decodeErr := decoder.Decode(&discard); decodeErr != nil {
			return nil, decodeErr
		}
	}

	_, err = decoder.Token()
	if err != nil {
		return nil, err
	}

	return keys, nil
}

func indentTopLevelJSONValue(value []byte) string {
	text := string(value)
	if !strings.Contains(text, "\n") {
		return text
	}

	lines := strings.Split(text, "\n")
	for i := 1; i < len(lines); i++ {
		lines[i] = "  " + lines[i]
	}

	return strings.Join(lines, "\n")
}

type singboxInboundDefaults struct {
	network     string
	security    string
	path        string
	hostHeader  string
	sni         string
	alpn        string
	fingerprint string
	publicKey   string
	shortID     string
	flow        string
}

func buildSingboxOutbound(host SubscriptionHost, user SubscriptionUser) *orderedmap.OrderedMap {
	protocol := ""
	if host.InboundType != nil {
		protocol = strings.ToLower(strings.TrimSpace(*host.InboundType))
	}
	if protocol == "" {
		return nil
	}

	if protocol == "ss" {
		protocol = "shadowsocks"
	}

	if protocol != "vless" && protocol != "trojan" && protocol != "shadowsocks" {
		return nil
	}

	defaults := resolveSingboxInboundDefaults(host)

	if !isSupportedSingboxTransport(defaults.network) {
		return nil
	}

	remark := strings.TrimSpace(host.Remark)
	if remark == "" {
		remark = host.Address
	}

	outbound := orderedmap.New()
	outbound.Set("type", protocol)
	outbound.Set("tag", remark)
	outbound.Set("server", host.Address)
	outbound.Set("server_port", host.Port)

	switch protocol {
	case "vless":
		if defaults.flow == "xtls-rprx-vision" {
			outbound.Set("flow", "xtls-rprx-vision")
		}
		outbound.Set("uuid", user.VlessUUID)
	case "trojan":
		outbound.Set("password", user.TrojanPassword)
	case "shadowsocks":
		method := extractShadowsocksMethod(host.InboundRaw)
		if method == "" {
			method = "chacha20-ietf-poly1305"
		}
		outbound.Set("password", user.SSPassword)
		outbound.Set("method", method)
		outbound.Set("network", "tcp")
	}

	if defaults.security == "tls" || defaults.security == "reality" {
		tlsCfg := orderedmap.New()
		tlsCfg.Set("enabled", true)

		sni := defaults.sni
		if host.OverrideSNIFromAddress {
			sni = host.Address
		}
		if host.KeepSNIBlank {
			sni = ""
		}
		if sni != "" {
			tlsCfg.Set("server_name", sni)
		}

		if host.AllowInsecure {
			tlsCfg.Set("insecure", true)
		}

		if defaults.fingerprint != "" {
			utlsCfg := orderedmap.New()
			utlsCfg.Set("enabled", true)
			utlsCfg.Set("fingerprint", defaults.fingerprint)
			tlsCfg.Set("utls", utlsCfg)
		} else if defaults.security == "reality" {
			// exodus default for reality when fp is empty
			utlsCfg := orderedmap.New()
			utlsCfg.Set("enabled", true)
			utlsCfg.Set("fingerprint", "chrome")
			tlsCfg.Set("utls", utlsCfg)
		}

		if defaults.alpn != "" {
			raw := strings.Split(defaults.alpn, ",")
			alpn := make([]string, 0, len(raw))
			for _, v := range raw {
				v = strings.TrimSpace(v)
				if v != "" {
					alpn = append(alpn, v)
				}
			}
			if len(alpn) > 0 {
				tlsCfg.Set("alpn", alpn)
			}
		}

		if defaults.security == "reality" {
			reality := orderedmap.New()
			reality.Set("enabled", true)
			if defaults.publicKey != "" {
				reality.Set("public_key", defaults.publicKey)
			}
			if defaults.shortID != "" {
				reality.Set("short_id", defaults.shortID)
			}
			tlsCfg.Set("reality", reality)
		}

		outbound.Set("tls", tlsCfg)
	}

	if defaults.network == "ws" || defaults.network == "httpupgrade" {
		transport := orderedmap.New()
		transport.Set("type", defaults.network)

		path := defaults.path
		path, earlyData := extractEarlyDataFromPath(path)
		if path != "" {
			transport.Set("path", path)
		}

		if defaults.hostHeader != "" {
			headers := orderedmap.New()
			headers.Set("Host", defaults.hostHeader)
			transport.Set("headers", headers)
		}

		if defaults.network == "ws" && earlyData > 0 {
			transport.Set("max_early_data", earlyData)
			transport.Set("early_data_header_name", "Sec-WebSocket-Protocol")
		}

		outbound.Set("transport", transport)
	}

	if mux := parseJSONMapString(host.SingboxMuxParams); mux != nil {
		outbound.Set("multiplex", orderedMapFromMapWithPreferredOrder(mux, []string{
			"enabled",
			"protocol",
			"max_connections",
			"min_streams",
			"padding",
		}))
	}

	return outbound
}

func isSupportedSingboxTransport(network string) bool {
	switch network {
	case "", "tcp", "raw", "ws", "httpupgrade":
		return true
	default:
		return false
	}
}

func resolveSingboxInboundDefaults(host SubscriptionHost) singboxInboundDefaults {
	defaults := singboxInboundDefaults{
		network:  "tcp",
		security: "none",
	}

	raw := parseInboundRaw(host.InboundRaw)
	streamSettings := readMap(raw, "streamSettings")

	if host.InboundNetwork != nil && strings.TrimSpace(*host.InboundNetwork) != "" {
		defaults.network = strings.ToLower(strings.TrimSpace(*host.InboundNetwork))
	} else if network := readString(streamSettings, "network"); network != "" {
		defaults.network = strings.ToLower(network)
	}

	if host.InboundSecurity != nil && strings.TrimSpace(*host.InboundSecurity) != "" {
		defaults.security = strings.ToLower(strings.TrimSpace(*host.InboundSecurity))
	} else if security := readString(streamSettings, "security"); security != "" {
		defaults.security = strings.ToLower(security)
	} else {
		switch strings.ToUpper(strings.TrimSpace(host.SecurityLayer)) {
		case "TLS":
			defaults.security = "tls"
		case "NONE":
			defaults.security = "none"
		}
	}

	switch defaults.network {
	case "ws":
		wsSettings := readMap(streamSettings, "wsSettings")
		defaults.path = firstNonEmpty(
			derefString(host.Path),
			readString(wsSettings, "path"),
		)
		defaults.hostHeader = firstNonEmpty(
			derefString(host.Host),
			readString(readMap(wsSettings, "headers"), "Host"),
			readString(wsSettings, "host"),
		)
	case "httpupgrade":
		upgradeSettings := readMap(streamSettings, "httpupgradeSettings")
		defaults.path = firstNonEmpty(
			derefString(host.Path),
			readString(upgradeSettings, "path"),
		)
		defaults.hostHeader = firstNonEmpty(
			derefString(host.Host),
			readString(readMap(upgradeSettings, "headers"), "Host"),
			readString(upgradeSettings, "host"),
		)
	default:
		defaults.path = derefString(host.Path)
		defaults.hostHeader = derefString(host.Host)
	}

	switch defaults.security {
	case "tls":
		tlsSettings := readMap(streamSettings, "tlsSettings")
		defaults.sni = firstNonEmpty(derefString(host.SNI), readString(tlsSettings, "serverName"))
		defaults.fingerprint = firstNonEmpty(
			derefString(host.Fingerprint),
			readString(tlsSettings, "fingerprint"),
		)
		defaults.alpn = firstNonEmpty(
			derefString(host.ALPN),
			joinStringSlice(readStringSlice(tlsSettings, "alpn")),
		)
	case "reality":
		realitySettings := readMap(streamSettings, "realitySettings")
		defaults.sni = firstNonEmpty(
			derefString(host.SNI),
			readFirstString(readStringSlice(realitySettings, "serverNames")),
		)
		defaults.fingerprint = firstNonEmpty(
			derefString(host.Fingerprint),
			readString(realitySettings, "fingerprint"),
		)
		defaults.alpn = firstNonEmpty(derefString(host.ALPN), joinStringSlice(readStringSlice(realitySettings, "alpn")))
		defaults.shortID = readRandomString(readStringSlice(realitySettings, "shortIds"))
		defaults.publicKey = readString(realitySettings, "publicKey")
		if defaults.publicKey == "" {
			defaults.publicKey = deriveRealityPublicKey(readString(realitySettings, "privateKey"))
		}
	}

	defaults.flow = resolveVlessFlow(host, defaults)

	if defaults.sni == "" && host.OverrideSNIFromAddress {
		defaults.sni = host.Address
	}
	if host.KeepSNIBlank {
		defaults.sni = ""
	}

	return defaults
}

func resolveVlessFlow(host SubscriptionHost, defaults singboxInboundDefaults) string {
	if host.InboundType == nil || !strings.EqualFold(*host.InboundType, "vless") {
		return ""
	}

	raw := parseInboundRaw(host.InboundRaw)
	settings := readMap(raw, "settings")
	flowFromSettings := strings.TrimSpace(readString(settings, "flow"))
	if flowFromSettings == "xtls-rprx-vision" {
		return flowFromSettings
	}

	if (defaults.network == "tcp" || defaults.network == "raw") &&
		(defaults.security == "tls" || defaults.security == "reality") {
		return "xtls-rprx-vision"
	}

	return ""
}

func parseInboundRaw(raw json.RawMessage) map[string]interface{} {
	if len(raw) == 0 {
		return map[string]interface{}{}
	}
	var parsed map[string]interface{}
	if err := json.Unmarshal(raw, &parsed); err != nil || parsed == nil {
		return map[string]interface{}{}
	}
	return parsed
}

func readMap(src map[string]interface{}, key string) map[string]interface{} {
	if src == nil {
		return map[string]interface{}{}
	}
	value, ok := src[key]
	if !ok || value == nil {
		return map[string]interface{}{}
	}
	if result, ok := value.(map[string]interface{}); ok && result != nil {
		return result
	}
	return map[string]interface{}{}
}

func readString(src map[string]interface{}, key string) string {
	if src == nil {
		return ""
	}
	value, ok := src[key]
	if !ok || value == nil {
		return ""
	}
	str, ok := value.(string)
	if !ok {
		return ""
	}
	return strings.TrimSpace(str)
}

func readStringSlice(src map[string]interface{}, key string) []string {
	if src == nil {
		return nil
	}
	value, ok := src[key]
	if !ok || value == nil {
		return nil
	}

	interfaces, ok := value.([]interface{})
	if !ok {
		return nil
	}

	result := make([]string, 0, len(interfaces))
	for _, item := range interfaces {
		str, ok := item.(string)
		if !ok {
			continue
		}
		str = strings.TrimSpace(str)
		if str != "" {
			result = append(result, str)
		}
	}

	return result
}

func joinStringSlice(values []string) string {
	if len(values) == 0 {
		return ""
	}
	return strings.Join(values, ",")
}

func readFirstString(values []string) string {
	if len(values) == 0 {
		return ""
	}
	return strings.TrimSpace(values[0])
}

func readRandomString(values []string) string {
	filtered := make([]string, 0, len(values))
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value != "" {
			filtered = append(filtered, value)
		}
	}
	if len(filtered) == 0 {
		return ""
	}
	return filtered[rand.Intn(len(filtered))]
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value != "" {
			return value
		}
	}
	return ""
}

func deriveRealityPublicKey(privateKey string) string {
	privateKey = strings.TrimSpace(privateKey)
	if privateKey == "" {
		return ""
	}

	raw, ok := decodeBase64Any(privateKey)
	if !ok || len(raw) != 32 {
		return ""
	}

	var scalar [32]byte
	copy(scalar[:], raw)
	var public [32]byte
	curve25519.ScalarBaseMult(&public, &scalar)

	return base64.RawURLEncoding.EncodeToString(public[:])
}

func decodeBase64Any(value string) ([]byte, bool) {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil, false
	}

	encodings := []*base64.Encoding{
		base64.RawStdEncoding,
		base64.StdEncoding,
		base64.RawURLEncoding,
		base64.URLEncoding,
	}

	for _, encoding := range encodings {
		if decoded, err := encoding.DecodeString(value); err == nil {
			return decoded, true
		}
	}

	return nil, false
}

func patchSingboxSelectors(baseConfig *orderedmap.OrderedMap, preferredHostNodeTags, regularHostNodeTags []string) {
	rawValue, ok := baseConfig.Get("outbounds")
	if !ok {
		return
	}

	rawOutbounds, ok := rawValue.([]interface{})
	if !ok {
		return
	}

	knownHostTags := appendUniqueStrings(append([]string(nil), preferredHostNodeTags...), regularHostNodeTags...)
	knownHostSet := make(map[string]struct{}, len(knownHostTags))
	for _, tag := range knownHostTags {
		knownHostSet[tag] = struct{}{}
	}

	allNodeTags := make([]string, 0, len(knownHostTags))
	urltestTags := make([]string, 0, len(rawOutbounds))
	for _, item := range rawOutbounds {
		ob, ok := orderedMapValue(item)
		if !ok {
			continue
		}
		typ := orderedMapString(ob, "type")
		tag := orderedMapString(ob, "tag")
		if tag == "" {
			continue
		}
		if typ == "urltest" {
			urltestTags = append(urltestTags, tag)
			continue
		}
		if _, isHostNode := knownHostSet[tag]; isHostNode {
			allNodeTags = append(allNodeTags, tag)
		}
	}

	for index, item := range rawOutbounds {
		ob, ok := orderedMapValue(item)
		if !ok {
			continue
		}
		typ := orderedMapString(ob, "type")
		switch typ {
		case "urltest":
			ob.Set("outbounds", append([]string(nil), allNodeTags...))
		case "selector":
			existingEntries := orderedMapStrings(ob, "outbounds")
			middleEntries := make([]string, 0, len(existingEntries))
			for _, entry := range existingEntries {
				if _, isHostNode := knownHostSet[entry]; isHostNode {
					continue
				}
				middleEntries = append(middleEntries, entry)
			}
			if len(middleEntries) == 0 {
				middleEntries = append(middleEntries, urltestTags...)
			}

			selectorTags := make([]string, 0, len(preferredHostNodeTags)+len(middleEntries)+len(regularHostNodeTags))
			selectorTags = append(selectorTags, preferredHostNodeTags...)
			selectorTags = append(selectorTags, middleEntries...)
			selectorTags = append(selectorTags, regularHostNodeTags...)
			ob.Set("outbounds", appendUniqueStrings(nil, selectorTags...))
		}
		rawOutbounds[index] = ob
	}

	baseConfig.Set("outbounds", rawOutbounds)
}

func orderedMapValue(value interface{}) (orderedmap.OrderedMap, bool) {
	switch typed := value.(type) {
	case orderedmap.OrderedMap:
		return typed, true
	case *orderedmap.OrderedMap:
		if typed == nil {
			return orderedmap.OrderedMap{}, false
		}
		return *typed, true
	default:
		return orderedmap.OrderedMap{}, false
	}
}

func orderedMapString(obj orderedmap.OrderedMap, key string) string {
	value, ok := obj.Get(key)
	if !ok || value == nil {
		return ""
	}
	str, ok := value.(string)
	if !ok {
		return ""
	}
	return strings.TrimSpace(str)
}

func orderedMapStrings(obj orderedmap.OrderedMap, key string) []string {
	value, ok := obj.Get(key)
	if !ok || value == nil {
		return nil
	}

	switch typed := value.(type) {
	case []string:
		return appendUniqueStrings(nil, typed...)
	case []interface{}:
		result := make([]string, 0, len(typed))
		for _, item := range typed {
			str, ok := item.(string)
			if !ok {
				continue
			}
			result = append(result, str)
		}
		return appendUniqueStrings(nil, result...)
	default:
		return nil
	}
}

func orderedMapFromMapWithPreferredOrder(values map[string]interface{}, preferred []string) orderedmap.OrderedMap {
	obj := orderedmap.New()
	used := make(map[string]struct{}, len(values))

	for _, key := range preferred {
		value, exists := values[key]
		if !exists {
			continue
		}
		obj.Set(key, orderedJSONValue(value))
		used[key] = struct{}{}
	}

	remainingKeys := make([]string, 0, len(values))
	for key := range values {
		if _, exists := used[key]; exists {
			continue
		}
		remainingKeys = append(remainingKeys, key)
	}
	sort.Strings(remainingKeys)

	for _, key := range remainingKeys {
		obj.Set(key, orderedJSONValue(values[key]))
	}

	return *obj
}

func orderedJSONValue(value interface{}) interface{} {
	switch typed := value.(type) {
	case map[string]interface{}:
		obj := orderedMapFromMapWithPreferredOrder(typed, nil)
		return obj
	case []interface{}:
		items := make([]interface{}, 0, len(typed))
		for _, item := range typed {
			items = append(items, orderedJSONValue(item))
		}
		return items
	default:
		return value
	}
}

func derefString(v *string) string {
	if v == nil {
		return ""
	}
	return strings.TrimSpace(*v)
}

func extractEarlyDataFromPath(path string) (string, int) {
	if path == "" || !strings.Contains(path, "?ed=") {
		return path, 0
	}

	parts := strings.SplitN(path, "?ed=", 2)
	cleanPath := strings.TrimSpace(parts[0])
	edPart := strings.TrimSpace(parts[1])
	if idx := strings.Index(edPart, "/"); idx >= 0 {
		edPart = edPart[:idx]
	}
	n, err := strconv.Atoi(edPart)
	if err != nil || n <= 0 {
		return cleanPath, 0
	}

	return cleanPath, n
}

func generateYAMLConfig(templateYAML []byte, hosts []SubscriptionHost, user SubscriptionUser) (string, error) {
	var root yaml.Node
	if len(templateYAML) > 0 {
		if err := yaml.Unmarshal(templateYAML, &root); err != nil {
			root = yaml.Node{}
		}
	}

	topLevelSpacing := extractYAMLTopLevelSpacing(templateYAML)
	config := ensureYAMLDocumentMappingNode(&root)
	proxiesNode := ensureYAMLMappingSequenceValue(config, "proxies")

	proxyNames := []string{}
	leadingSelectorProxyNames := []string{}
	trailingSelectorProxyNames := []string{}
	for _, host := range hosts {
		proxy := buildMihomoProxy(host, user)
		if proxy == nil {
			continue
		}
		proxiesNode.Content = append(proxiesNode.Content, buildOrderedYAMLValueNode("proxy", proxy))
		if name, ok := proxy["name"].(string); ok && name != "" {
			proxyNames = append(proxyNames, name)
			if host.SelectorNodesFirst {
				leadingSelectorProxyNames = append(leadingSelectorProxyNames, name)
			} else {
				trailingSelectorProxyNames = append(trailingSelectorProxyNames, name)
			}
		}
	}

	groupsNode := ensureYAMLMappingSequenceValue(config, "proxy-groups")
	for _, group := range groupsNode.Content {
		if group == nil || group.Kind != yaml.MappingNode {
			continue
		}

		groupProxies := ensureYAMLMappingSequenceValue(group, "proxies")
		groupType := strings.ToLower(strings.TrimSpace(yamlMappingString(group, "type")))
		switch groupType {
		case "select":
			existingEntries := yamlSequenceStrings(groupProxies)
			middleEntries := make([]string, 0, len(existingEntries))
			hostNames := make(map[string]struct{}, len(proxyNames))
			for _, name := range proxyNames {
				hostNames[name] = struct{}{}
			}
			for _, entry := range existingEntries {
				if _, isHostName := hostNames[entry]; isHostName {
					continue
				}
				middleEntries = append(middleEntries, entry)
			}
			finalEntries := make([]string, 0, len(leadingSelectorProxyNames)+len(middleEntries)+len(trailingSelectorProxyNames))
			finalEntries = append(finalEntries, leadingSelectorProxyNames...)
			finalEntries = append(finalEntries, middleEntries...)
			finalEntries = append(finalEntries, trailingSelectorProxyNames...)
			setYAMLSequenceStrings(groupProxies, finalEntries)
		case "url-test", "urltest":
			setYAMLSequenceStrings(groupProxies, proxyNames)
		default:
			finalEntries := appendUniqueStrings(yamlSequenceStrings(groupProxies), proxyNames...)
			setYAMLSequenceStrings(groupProxies, finalEntries)
		}
	}

	var buf bytes.Buffer
	encoder := yaml.NewEncoder(&buf)
	encoder.SetIndent(2)
	err := encoder.Encode(config)
	closeErr := encoder.Close()
	if err != nil {
		return "", err
	}
	if closeErr != nil {
		return "", closeErr
	}

	rendered := applyYAMLTopLevelSpacing(buf.String(), topLevelSpacing)
	return rendered, nil
}

func ensureYAMLDocumentMappingNode(root *yaml.Node) *yaml.Node {
	if root.Kind != yaml.DocumentNode {
		*root = yaml.Node{Kind: yaml.DocumentNode}
	}
	if len(root.Content) == 0 || root.Content[0] == nil || root.Content[0].Kind != yaml.MappingNode {
		root.Content = []*yaml.Node{{Kind: yaml.MappingNode, Tag: "!!map"}}
	}
	return root.Content[0]
}

func ensureYAMLMappingSequenceValue(mapping *yaml.Node, key string) *yaml.Node {
	if mapping == nil || mapping.Kind != yaml.MappingNode {
		return &yaml.Node{Kind: yaml.SequenceNode, Tag: "!!seq"}
	}

	for i := 0; i+1 < len(mapping.Content); i += 2 {
		keyNode := mapping.Content[i]
		if keyNode == nil || keyNode.Value != key {
			continue
		}

		valueNode := mapping.Content[i+1]
		if valueNode == nil || valueNode.Kind != yaml.SequenceNode {
			valueNode = &yaml.Node{Kind: yaml.SequenceNode, Tag: "!!seq"}
			mapping.Content[i+1] = valueNode
		}
		return valueNode
	}

	valueNode := &yaml.Node{Kind: yaml.SequenceNode, Tag: "!!seq"}
	mapping.Content = append(mapping.Content, newYAMLScalarNode(key), valueNode)
	return valueNode
}

func yamlMappingString(mapping *yaml.Node, key string) string {
	if mapping == nil || mapping.Kind != yaml.MappingNode {
		return ""
	}

	for i := 0; i+1 < len(mapping.Content); i += 2 {
		keyNode := mapping.Content[i]
		if keyNode == nil || keyNode.Value != key {
			continue
		}

		valueNode := mapping.Content[i+1]
		if valueNode == nil || valueNode.Kind != yaml.ScalarNode {
			return ""
		}
		return strings.TrimSpace(valueNode.Value)
	}

	return ""
}

func yamlSequenceStrings(sequence *yaml.Node) []string {
	if sequence == nil || sequence.Kind != yaml.SequenceNode {
		return nil
	}

	values := make([]string, 0, len(sequence.Content))
	for _, item := range sequence.Content {
		if item == nil || item.Kind != yaml.ScalarNode {
			continue
		}
		value := strings.TrimSpace(item.Value)
		if value == "" {
			continue
		}
		values = append(values, value)
	}

	return values
}

func setYAMLSequenceStrings(sequence *yaml.Node, values []string) {
	if sequence == nil {
		return
	}
	sequence.Kind = yaml.SequenceNode
	sequence.Tag = "!!seq"
	sequence.Content = sequence.Content[:0]
	for _, value := range values {
		sequence.Content = append(sequence.Content, newYAMLScalarNode(value))
	}
}

func appendUniqueStrings(base []string, extra ...string) []string {
	result := make([]string, 0, len(base)+len(extra))
	seen := make(map[string]struct{}, len(base)+len(extra))

	for _, value := range base {
		value = strings.TrimSpace(value)
		if value == "" {
			continue
		}
		if _, exists := seen[value]; exists {
			continue
		}
		seen[value] = struct{}{}
		result = append(result, value)
	}

	for _, value := range extra {
		value = strings.TrimSpace(value)
		if value == "" {
			continue
		}
		if _, exists := seen[value]; exists {
			continue
		}
		seen[value] = struct{}{}
		result = append(result, value)
	}

	return result
}

func buildOrderedYAMLValueNode(parentKey string, value interface{}) *yaml.Node {
	switch v := value.(type) {
	case map[string]interface{}:
		return buildOrderedYAMLMappingNode(parentKey, v)
	case []string:
		node := &yaml.Node{Kind: yaml.SequenceNode, Tag: "!!seq"}
		for _, item := range v {
			node.Content = append(node.Content, newYAMLScalarNode(item))
		}
		return node
	case []interface{}:
		node := &yaml.Node{Kind: yaml.SequenceNode, Tag: "!!seq"}
		for _, item := range v {
			node.Content = append(node.Content, buildOrderedYAMLValueNode("", item))
		}
		return node
	case string:
		return newYAMLScalarNode(v)
	case bool:
		return &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!bool", Value: strconv.FormatBool(v)}
	case int:
		return &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!int", Value: strconv.Itoa(v)}
	case int64:
		return &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!int", Value: strconv.FormatInt(v, 10)}
	case float64:
		return &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!float", Value: strconv.FormatFloat(v, 'f', -1, 64)}
	case nil:
		return &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!null"}
	default:
		return newYAMLScalarNode(fmt.Sprint(v))
	}
}

func buildOrderedYAMLMappingNode(parentKey string, values map[string]interface{}) *yaml.Node {
	node := &yaml.Node{Kind: yaml.MappingNode, Tag: "!!map"}
	usedKeys := make(map[string]struct{}, len(values))

	for _, key := range preferredYAMLKeyOrder(parentKey) {
		value, exists := values[key]
		if !exists {
			continue
		}
		appendYAMLMappingEntry(node, key, buildOrderedYAMLValueNode(key, value))
		usedKeys[key] = struct{}{}
	}

	remainingKeys := make([]string, 0, len(values))
	for key := range values {
		if _, exists := usedKeys[key]; exists {
			continue
		}
		remainingKeys = append(remainingKeys, key)
	}
	sort.Strings(remainingKeys)

	for _, key := range remainingKeys {
		appendYAMLMappingEntry(node, key, buildOrderedYAMLValueNode(key, values[key]))
	}

	return node
}

func preferredYAMLKeyOrder(parentKey string) []string {
	switch parentKey {
	case "proxy":
		return []string{
			"name",
			"type",
			"server",
			"port",
			"udp",
			"network",
			"uuid",
			"password",
			"cipher",
			"tls",
			"skip-cert-verify",
			"servername",
			"client-fingerprint",
			"alpn",
			"packet-encoding",
			"flow",
			"encryption",
			"ws-opts",
			"grpc-opts",
			"smux",
		}
	case "ws-opts":
		return []string{"path", "headers", "max-early-data", "early-data-header-name", "v2ray-http-upgrade", "v2ray-http-upgrade-fast-open"}
	case "headers":
		return []string{"Host"}
	case "grpc-opts":
		return []string{"grpc-service-name"}
	case "smux":
		return []string{"enabled", "protocol", "max-connections", "padding"}
	default:
		return nil
	}
}

func appendYAMLMappingEntry(node *yaml.Node, key string, value *yaml.Node) {
	node.Content = append(node.Content, newYAMLScalarNode(key), value)
}

func newYAMLScalarNode(value string) *yaml.Node {
	return &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!str", Value: value}
}

func extractYAMLTopLevelSpacing(templateYAML []byte) map[string]bool {
	if len(templateYAML) == 0 {
		return nil
	}

	lines := strings.Split(strings.ReplaceAll(string(templateYAML), "\r\n", "\n"), "\n")
	spacing := map[string]bool{}
	previousBlank := false

	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			previousBlank = true
			continue
		}

		if key, ok := extractYAMLTopLevelKey(line); ok && previousBlank {
			spacing[key] = true
		}

		previousBlank = false
	}

	return spacing
}

func applyYAMLTopLevelSpacing(rendered string, spacing map[string]bool) string {
	if rendered == "" || len(spacing) == 0 {
		return rendered
	}

	hasTrailingNewline := strings.HasSuffix(rendered, "\n")
	lines := strings.Split(strings.TrimSuffix(rendered, "\n"), "\n")
	output := make([]string, 0, len(lines)+len(spacing))

	for index, line := range lines {
		if key, ok := extractYAMLTopLevelKey(line); ok && index > 0 && spacing[key] {
			if len(output) > 0 && strings.TrimSpace(output[len(output)-1]) != "" {
				output = append(output, "")
			}
		}
		output = append(output, line)
	}

	result := strings.Join(output, "\n")
	if hasTrailingNewline {
		result += "\n"
	}

	return result
}

func extractYAMLTopLevelKey(line string) (string, bool) {
	if line == "" {
		return "", false
	}
	if line[0] == ' ' || line[0] == '\t' || strings.HasPrefix(strings.TrimSpace(line), "- ") {
		return "", false
	}

	index := strings.IndexByte(line, ':')
	if index <= 0 {
		return "", false
	}

	key := strings.TrimSpace(line[:index])
	if key == "" {
		return "", false
	}

	return key, true
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

	if mux := parseJSONMapString(host.ClashMuxParams); mux != nil {
		proxy["smux"] = mux
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
