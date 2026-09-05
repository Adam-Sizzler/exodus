package subscription

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/http"
	"net/url"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"exodus/internal/config"
	"exodus/internal/db"
	"exodus/internal/httpapi/externalsquads"
	"exodus/internal/httpapi/shared"
	"exodus/internal/httpapi/subscriptionresponserules"
	"exodus/internal/httpapi/subscriptionsettings"
	"exodus/internal/jobqueue"
	"exodus/internal/logger"
	"exodus/internal/util"
)

var (
	subSettingsCacheLock sync.RWMutex
	subSettingsCached    *SubscriptionSettingsParsed
	subSettingsCacheTime time.Time

	squadOverridesCacheLock sync.RWMutex
	squadOverridesCache     = make(map[string]cachedSquadOverride)

	subNodeBaseLock sync.RWMutex
	subNodeBaseVal  string
	subNodeBaseExp  time.Time

	subTemplateLock      sync.RWMutex
	subTemplateTypeCache = make(map[string]cachedTemplate)
	subTemplateNameCache = make(map[string]cachedNamedTemplate)
	subTemplateUUIDCache = make(map[string]cachedNamedTemplate)
)

type cachedSquadOverride struct {
	overrides *ExternalSquadOverrides
	expiresAt time.Time
}

type cachedTemplate struct {
	data      []byte
	expiresAt time.Time
}

type cachedNamedTemplate struct {
	templateType string
	data         []byte
	expiresAt    time.Time
}

const (
	subSettingsCacheTTL    = 1 * time.Hour
	squadOverridesCacheTTL = 1 * time.Hour
	subNodeBaseTTL         = 30 * time.Second
	subTemplateCacheTTL    = 5 * time.Minute
)

func init() {
	subscriptionsettings.OnSettingsUpdated = InvalidateSubscriptionSettingsCache
	externalsquads.OnSquadUpdated = InvalidateExternalSquadCache
}

// InvalidateSubscriptionSettingsCache clears the cached subscription settings.
func InvalidateSubscriptionSettingsCache() {
	subSettingsCacheLock.Lock()
	subSettingsCached = nil
	subSettingsCacheLock.Unlock()
}

// InvalidateExternalSquadCache clears the cached squad overrides.
func InvalidateExternalSquadCache(squadUUID string) {
	squadOverridesCacheLock.Lock()
	if squadUUID == "" {
		squadOverridesCache = make(map[string]cachedSquadOverride)
	} else {
		delete(squadOverridesCache, squadUUID)
	}
	squadOverridesCacheLock.Unlock()
}

func loadSubscriptionSettings(ctx context.Context, dbConn *sql.DB, _ *config.BackendConfig) (SubscriptionSettingsParsed, error) {
	subSettingsCacheLock.RLock()
	if subSettingsCached != nil && time.Since(subSettingsCacheTime) < subSettingsCacheTTL {
		cached := *subSettingsCached
		subSettingsCacheLock.RUnlock()
		return cached, nil
	}
	subSettingsCacheLock.RUnlock()

	var parsed SubscriptionSettingsParsed

	row := dbConn.QueryRowContext(ctx, `
		SELECT uuid, address, port, api_schema, api_path,
			   serve_json_at_base_subscription, is_show_custom_remarks, custom_remarks,
			   custom_response_headers, randomize_hosts, response_rules, hwid_settings,
			   created_at, updated_at
		FROM subscription_settings
		ORDER BY created_at ASC
		LIMIT 1
	`)

	settings, err := subscriptionsettings.ScanSubscriptionSettings(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			_, _ = dbConn.ExecContext(ctx, `
				INSERT INTO subscription_settings (
					uuid, address, port, api_schema, api_path,
					serve_json_at_base_subscription, is_show_custom_remarks, custom_remarks,
					custom_response_headers, randomize_hosts, response_rules, hwid_settings
				) VALUES (
					'00000000-0000-0000-0000-000000000000', '', 9263, 'grpc', '',
					false, true, '{}'::jsonb,
					'{"profile-title":"exEncodeBase64:exodus","support-url":"https://github.com","profile-update-interval":"12"}'::jsonb,
					false, '[]'::jsonb, '{}'::jsonb
				) ON CONFLICT DO NOTHING
			`)
			rowRetry := dbConn.QueryRowContext(ctx, `
				SELECT uuid, address, port, api_schema, api_path,
					   serve_json_at_base_subscription, is_show_custom_remarks, custom_remarks,
					   custom_response_headers, randomize_hosts, response_rules, hwid_settings,
					   created_at, updated_at
				FROM subscription_settings
				ORDER BY created_at ASC
				LIMIT 1
			`)
			if sRetry, errRetry := subscriptionsettings.ScanSubscriptionSettings(rowRetry); errRetry == nil {
				settings = sRetry
				err = nil
			}
		}
		if err != nil {
			return parsed, err
		}
	}

	parsed.Raw = settings
	parsed.CustomResponseHeaders = map[string]string{}
	parsed.CustomRemarks = CustomRemarks{}
	parsed.HwidSettings = HwidSettings{Enabled: false, FallbackDeviceLimit: 999}

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

	subSettingsCacheLock.Lock()
	subSettingsCached = &parsed
	subSettingsCacheTime = time.Now()
	subSettingsCacheLock.Unlock()

	return parsed, nil
}

func loadExternalSquadOverrides(ctx context.Context, dbConn *sql.DB, squadUUID string, cfg *config.BackendConfig) (*ExternalSquadOverrides, error) {
	if strings.TrimSpace(squadUUID) == "" {
		return nil, nil
	}

	squadOverridesCacheLock.RLock()
	if cached, ok := squadOverridesCache[squadUUID]; ok && time.Now().Before(cached.expiresAt) {
		squadOverridesCacheLock.RUnlock()
		return cached.overrides, nil
	}
	squadOverridesCacheLock.RUnlock()

	log := cfg.Logger.RoleService(logger.RoleAPI, logger.ServiceHTTP)
	overrides := &ExternalSquadOverrides{
		Templates: make(map[string]string),
	}

	var subscriptionSettingsJSON, hostOverridesJSON, responseHeadersAddJSON, hwidSettingsJSON, customRemarksJSON sql.NullString
	var responseHeadersRemoveRaw sql.NullString

	query := `SELECT subscription_settings, host_overrides, response_headers_add,
			  array_to_json(COALESCE(response_headers_remove, ARRAY[]::text[]))::text AS response_headers_remove,
			  hwid_settings, custom_remarks
			  FROM external_squads WHERE uuid = $1 LIMIT 1`
	row := dbConn.QueryRowContext(ctx, query, squadUUID)

	if err := row.Scan(&subscriptionSettingsJSON, &hostOverridesJSON, &responseHeadersAddJSON, &responseHeadersRemoveRaw, &hwidSettingsJSON, &customRemarksJSON); err != nil {
		return nil, err
	}
	if responseHeadersRemoveRaw.Valid {
		overrides.ResponseHeadersRemove = shared.ParsePgTextArray(responseHeadersRemoveRaw.String)
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
	if responseHeadersAddJSON.Valid && responseHeadersAddJSON.String != "" {
		var rh map[string]string
		if err := json.Unmarshal([]byte(responseHeadersAddJSON.String), &rh); err == nil {
			overrides.ResponseHeaders = rh
			log.Debug("Loaded response_headers_add override")
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

	rows, err := dbConn.QueryContext(ctx, `
		SELECT t.name, est.template_type
		FROM external_squads_templates est
		JOIN subscription_templates t ON t.uuid = est.template_uuid
		WHERE est.external_squad_uuid = $1
	`, squadUUID)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var templateName, templateType string
			if err := rows.Scan(&templateName, &templateType); err != nil {
				break
			}
			overrides.Templates[strings.ToUpper(templateType)] = templateName
			log.Debug("Loaded template override", "type", templateType, "name", templateName)
		}
		if rowsErr := rows.Err(); rowsErr != nil {
			log.Warn("Error iterating external squad templates", "error", rowsErr)
		}
	} else {
		log.Warn("Failed to load external squad templates", "error", err)
	}

	squadOverridesCacheLock.Lock()
	squadOverridesCache[squadUUID] = cachedSquadOverride{
		overrides: overrides,
		expiresAt: time.Now().Add(squadOverridesCacheTTL),
	}
	squadOverridesCacheLock.Unlock()

	return overrides, nil
}

func getSubscriptionUserByShortUUID(ctx context.Context, dbConn *sql.DB, shortUUID string) (SubscriptionUser, error) {
	return getSubscriptionUserByField(ctx, dbConn, "short_uuid", shortUUID)
}

func getSubscriptionUserByID(ctx context.Context, dbConn *sql.DB, userID int64) (SubscriptionUser, error) {
	return getSubscriptionUserByField(ctx, dbConn, "id", userID)
}

func getSubscriptionUserByUUID(ctx context.Context, dbConn *sql.DB, userUUID string) (SubscriptionUser, error) {
	return getSubscriptionUserByField(ctx, dbConn, "uuid", userUUID)
}

func getSubscriptionUserByUsername(ctx context.Context, dbConn *sql.DB, username string) (SubscriptionUser, error) {
	return getSubscriptionUserByField(ctx, dbConn, "username", username)
}

func getSubscriptionUserByField(ctx context.Context, dbConn *sql.DB, field string, value any) (SubscriptionUser, error) {
	var user SubscriptionUser

	var where string
	switch field {
	case "id":
		where = "u.id = $1"
	case "short_uuid":
		where = "u.short_uuid = $1"
	case "uuid":
		where = "u.uuid::text = $1 OR u.short_uuid = $1 OR u.username = $1"
	case "username":
		where = "u.username = $1"
	default:
		return user, fmt.Errorf("unsupported search field")
	}

	query := fmt.Sprintf(`
		SELECT u.id, u.uuid, u.short_uuid, u.username, u.status,
			   u.traffic_limit_bytes, u.traffic_limit_strategy, u.expire_at,
			   u.last_traffic_reset_at, u.created_at, u.updated_at, u.sub_revoked_at,
			   u.last_triggered_threshold, u.description, u.tag, u.telegram_id, u.email,
			   u.trojan_password, u.vless_uuid, u.ss_password,
			   u.naive_password, u.shadowtls_password, u.hysteria2_password, u.anytls_password,
			   u.hwid_device_limit, u.external_squad_uuid,
			   COALESCE(ut.used_traffic_bytes, 0), COALESCE(ut.lifetime_used_traffic_bytes, 0),
			   ut.online_at, ut.last_connected_node_uuid, ut.first_connected_at
		FROM users u
		LEFT JOIN user_traffic ut ON ut.id = u.id
		WHERE %s
		LIMIT 1
	`, where)

	row := dbConn.QueryRowContext(ctx, query, value)

	var lastTrafficReset, subRevokedAt, onlineAt, firstConnectedAt sql.NullTime
	var updatedAt sql.NullTime
	var description, tag, email, lastConnectedNodeUUID sql.NullString
	var telegramID sql.NullInt64
	var hwidDeviceLimit sql.NullInt64
	var lastTriggeredThreshold sql.NullInt64
	var externalSquadUUID sql.NullString
	var naivePassword, shadowtlsPassword, hysteria2Password, anytlsPassword sql.NullString
	if err := row.Scan(
		&user.ID,
		&user.UUID,
		&user.ShortUUID,
		&user.Username,
		&user.Status,
		&user.TrafficLimitBytes,
		&user.TrafficLimitStrategy,
		&user.ExpireAt,
		&lastTrafficReset,
		&user.CreatedAt,
		&updatedAt,
		&subRevokedAt,
		&lastTriggeredThreshold,
		&description,
		&tag,
		&telegramID,
		&email,
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
		&onlineAt,
		&lastConnectedNodeUUID,
		&firstConnectedAt,
	); err != nil {
		return user, err
	}

	if lastTrafficReset.Valid {
		user.LastTrafficResetAt = &lastTrafficReset.Time
	}
	if subRevokedAt.Valid {
		user.SubRevokedAt = &subRevokedAt.Time
	}
	if updatedAt.Valid {
		user.UpdatedAt = updatedAt.Time
	} else {
		user.UpdatedAt = user.CreatedAt
	}
	if lastTriggeredThreshold.Valid {
		user.LastTriggeredThreshold = int(lastTriggeredThreshold.Int64)
	}
	if onlineAt.Valid {
		user.OnlineAt = &onlineAt.Time
	}
	if firstConnectedAt.Valid {
		user.FirstConnectedAt = &firstConnectedAt.Time
	}
	if lastConnectedNodeUUID.Valid {
		user.LastConnectedNodeUUID = &lastConnectedNodeUUID.String
	}
	if description.Valid {
		user.Description = &description.String
	}
	if tag.Valid {
		user.Tag = &tag.String
	}
	if telegramID.Valid {
		user.TelegramID = &telegramID.Int64
	}
	if email.Valid {
		user.Email = &email.String
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

	return user, nil
}

func getHostsForUser(ctx context.Context, dbConn *sql.DB, user SubscriptionUser) ([]SubscriptionHost, error) {
	return getHostsForUserWithOptions(ctx, dbConn, user, false, false)
}

func getHostsForUserWithOptions(ctx context.Context, dbConn *sql.DB, user SubscriptionUser, withDisabled, withHidden bool) ([]SubscriptionHost, error) {
	whereClause := `ism.user_id = $1 AND (
		(COALESCE(h.internal_squads_mode, 'EXCLUDE') = 'ALLOW_ONLY' AND ishl.host_uuid IS NOT NULL)
		OR
		(COALESCE(h.internal_squads_mode, 'EXCLUDE') != 'ALLOW_ONLY' AND ishl.host_uuid IS NULL)
	)`
	if !withDisabled {
		whereClause += " AND NOT COALESCE(h.is_disabled, false)"
	}
	if !withHidden {
		whereClause += " AND NOT COALESCE(h.is_hidden, false)"
	}

	query := fmt.Sprintf(`
		SELECT DISTINCT h.uuid, h.view_position, h.remark, h.address, h.port,
			   h.path, h.sni, h.host, h.alpn, h.fingerprint, h.security_layer,
			   h.xhttp_extra_params, h.mux_params, h.mapper, h.sockopt_params, h.final_mask, h.is_disabled,
			   h.server_description, h.shuffle_host,
			   h.mihomo_x25519, h.mihomo_ip_version, h.xray_json_template_uuid, h.keep_sni_blank,
			   h.exclude_from_subscription_types, h.tags, h.is_hidden, h.override_sni_from_address,
			   h.config_profile_uuid, h.config_profile_inbound_uuid,
			   h.pinned_peer_cert_sha256, h.verify_peer_cert_by_name,
			   cpi.tag, cpi.type, cpi.network, cpi.security, cpi.port, cpi.raw_inbound
		FROM internal_squad_members ism
		JOIN internal_squad_inbounds isi ON ism.internal_squad_uuid = isi.internal_squad_uuid
		JOIN config_profile_inbounds cpi ON isi.inbound_uuid = cpi.uuid
		JOIN hosts h ON h.config_profile_inbound_uuid = cpi.uuid
		LEFT JOIN internal_squad_host_links ishl
			ON ishl.host_uuid = h.uuid AND ishl.squad_uuid = ism.internal_squad_uuid
		WHERE %s
		ORDER BY h.view_position ASC, h.remark ASC
	`, whereClause)

	rows, err := dbConn.QueryContext(ctx, query, user.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var hosts []SubscriptionHost
	for rows.Next() {
		host, err := scanSubscriptionHost(rows)
		if err != nil {
			return nil, err
		}
		hosts = append(hosts, host)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return hosts, nil
}

func scanSubscriptionHost(scanner shared.RowScanner) (SubscriptionHost, error) {
	var h SubscriptionHost
	var viewPosition sql.NullInt64
	var path, sni, host, alpn, fingerprint, securityLayer sql.NullString
	var xhttpExtraParams, muxParams, mapper, sockoptParams, finalMask, serverDescription sql.NullString
	var xrayJSONTemplateUUID, mihomoIPVersion, configProfileUUID, configProfileInboundUUID, pinnedPeerCertSha256, verifyPeerCertByName sql.NullString
	var inboundTag, inboundType, inboundNetwork, inboundSecurity sql.NullString
	var inboundPort sql.NullInt64
	var rawInbound sql.NullString
	var excludeTypes, hostTags db.StringArray
	var isDisabled, shuffleHost, mihomoX25519, keepSNIBlank, isHidden, overrideSNIFromAddress sql.NullBool

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
		&mapper,
		&sockoptParams,
		&finalMask,
		&isDisabled,
		&serverDescription,
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
		&pinnedPeerCertSha256,
		&verifyPeerCertByName,
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

	if pinnedPeerCertSha256.Valid {
		h.PinnedPeerCertSha256 = &pinnedPeerCertSha256.String
	}
	if verifyPeerCertByName.Valid {
		h.VerifyPeerCertByName = &verifyPeerCertByName.String
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
	if mapper.Valid && strings.TrimSpace(mapper.String) != "" {
		h.Mapper = ParseHostMapper([]byte(mapper.String))
	}
	if sockoptParams.Valid {
		h.SockoptParams = &sockoptParams.String
	}
	if finalMask.Valid {
		h.FinalMask = &finalMask.String
	}
	if isDisabled.Valid {
		h.IsDisabled = isDisabled.Bool
	}
	if serverDescription.Valid {
		h.ServerDescription = &serverDescription.String
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
	h.Tags = hostTags.Slice()
	if len(h.Tags) > 0 && strings.TrimSpace(h.Tags[0]) != "" {
		firstTag := strings.TrimSpace(h.Tags[0])
		h.Tag = &firstTag
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

// checkHwidDeviceLimit ports upstream's checkHwidDeviceLimit(): it never
// treats "device limit reached" and "no X-HWID header sent" as the same
// outcome, and DB failures are surfaced as a real error rather than folded
// into one of the two device-limit reasons.
func checkHwidDeviceLimit(ctx context.Context, dbConn *sql.DB, user SubscriptionUser, hwid *HwidHeaders, settings HwidSettings) (HwidCheckupResult, error) {
	if user.HwidDeviceLimit != nil && *user.HwidDeviceLimit == 0 {
		if hwid != nil {
			_ = enqueueOrUpsertHwidUserDevice(ctx, dbConn, user.ID, *hwid)
		}
		return HwidCheckupResult{Allowed: true, LimitBypassed: true}, nil
	}

	if hwid == nil {
		return HwidCheckupResult{Allowed: false, HwidNotSupported: true}, nil
	}

	exists, err := hwidDeviceExists(ctx, dbConn, user.ID, hwid.Hwid)
	if err == nil && exists {
		_ = enqueueOrUpsertHwidUserDevice(ctx, dbConn, user.ID, *hwid)
		return HwidCheckupResult{Allowed: true}, nil
	}

	limit := settings.FallbackDeviceLimit
	if user.HwidDeviceLimit != nil {
		limit = *user.HwidDeviceLimit
	}

	allowed, err := createHwidDeviceWithAdvisoryLock(ctx, dbConn, user.ID, *hwid, limit)
	if err != nil {
		return HwidCheckupResult{}, fmt.Errorf("create hwid device: %w", err)
	}
	if !allowed {
		return HwidCheckupResult{Allowed: false, MaxDeviceReached: true}, nil
	}

	return HwidCheckupResult{Allowed: true}, nil
}

// hwidLockPrefix (900000000n): added to
// the user ID to form the pg_advisory_xact_lock key, keeping this lock's
// numeric key space clear of the fixed key internal/db/migrations.go uses
// for its own advisory lock (2203092601 — far outside this per-user range
// for any realistic user ID).
const hwidLockPrefix int64 = 900000000

// createHwidDeviceWithAdvisoryLock serializes device creation per user: count-then-insert is a
// classic check-then-act race — two concurrent requests for the same user
// with two different new HWIDs can both pass a plain "count < limit" check
// before either INSERT commits, silently exceeding the device limit under
// concurrent load. pg_advisory_xact_lock(userID) serializes this per user
// (other users are unaffected) and releases automatically at transaction
// end, so a crashed/canceled request can't leave the lock held.
func createHwidDeviceWithAdvisoryLock(ctx context.Context, dbConn *sql.DB, userID int64, hwid HwidHeaders, deviceLimit int) (bool, error) {
	tx, err := dbConn.BeginTx(ctx, nil)
	if err != nil {
		return false, err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	if _, err := tx.ExecContext(ctx, `SELECT pg_advisory_xact_lock($1)`, hwidLockPrefix+userID); err != nil {
		return false, fmt.Errorf("acquire hwid advisory lock: %w", err)
	}

	// Check if this device already exists for this user. If it does, update its
	// metadata and allow access regardless of deviceLimit (matches upstream v3.4.2+ fix).
	var exists bool
	err = tx.QueryRowContext(ctx, `SELECT EXISTS(SELECT 1 FROM hwid_user_devices WHERE hwid = $1 AND user_id = $2)`, hwid.Hwid, userID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("check hwid device exists: %w", err)
	}

	hwid.Platform = lowerStringPtr(hwid.Platform)

	if exists {
		if _, err := tx.ExecContext(ctx, `
			UPDATE hwid_user_devices SET
				platform = COALESCE($3, platform),
				os_version = COALESCE($4, os_version),
				device_model = COALESCE($5, device_model),
				user_agent = COALESCE($6, user_agent),
				request_ip = COALESCE($7, request_ip),
				updated_at = now()
			WHERE hwid = $1 AND user_id = $2
		`, hwid.Hwid, userID, hwid.Platform, hwid.OsVersion, hwid.DeviceModel, hwid.UserAgent, hwid.RequestIP); err != nil {
			return false, fmt.Errorf("update hwid device: %w", err)
		}

		if err := tx.Commit(); err != nil {
			return false, fmt.Errorf("commit hwid device tx: %w", err)
		}

		return true, nil
	}

	// Device is new: enforce deviceLimit
	var count int
	if err := tx.QueryRowContext(ctx, `SELECT COUNT(*) FROM hwid_user_devices WHERE user_id = $1`, userID).Scan(&count); err != nil {
		return false, fmt.Errorf("count hwid devices: %w", err)
	}

	if count >= deviceLimit {
		return false, nil
	}

	if _, err := tx.ExecContext(ctx, `
		INSERT INTO hwid_user_devices (hwid, user_id, platform, os_version, device_model, user_agent, request_ip)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, hwid.Hwid, userID, hwid.Platform, hwid.OsVersion, hwid.DeviceModel, hwid.UserAgent, hwid.RequestIP); err != nil {
		return false, fmt.Errorf("insert hwid device: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return false, fmt.Errorf("commit hwid device tx: %w", err)
	}

	return true, nil
}

func hwidDeviceExists(ctx context.Context, dbConn *sql.DB, userID int64, hwid string) (bool, error) {
	var tmp int
	err := dbConn.QueryRowContext(ctx, `SELECT 1 FROM hwid_user_devices WHERE user_id = $1 AND hwid = $2`, userID, hwid).Scan(&tmp)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func upsertHwidUserDevice(ctx context.Context, dbConn *sql.DB, userID int64, hwid HwidHeaders) error {
	hwid.Platform = lowerStringPtr(hwid.Platform)
	_, err := dbConn.ExecContext(ctx, `
		INSERT INTO hwid_user_devices (hwid, user_id, platform, os_version, device_model, user_agent, request_ip)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (hwid, user_id)
		DO UPDATE SET
			platform = COALESCE(EXCLUDED.platform, hwid_user_devices.platform),
			os_version = COALESCE(EXCLUDED.os_version, hwid_user_devices.os_version),
			device_model = COALESCE(EXCLUDED.device_model, hwid_user_devices.device_model),
			user_agent = COALESCE(EXCLUDED.user_agent, hwid_user_devices.user_agent),
			request_ip = COALESCE(EXCLUDED.request_ip, hwid_user_devices.request_ip),
			updated_at = now()
	`, hwid.Hwid, userID, hwid.Platform, hwid.OsVersion, hwid.DeviceModel, hwid.UserAgent, hwid.RequestIP)
	return err
}

func enqueueOrUpsertHwidUserDevice(ctx context.Context, dbConn *sql.DB, userID int64, hwid HwidHeaders) error {
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
	return upsertHwidUserDevice(ctx, dbConn, userID, hwid)
}

func updateSubscriptionRequest(ctx context.Context, dbConn *sql.DB, userUUID string, userID int64, userAgent, requestIP, responseType, ruleName string) {
	if responseType == "" {
		responseType = "UNKNOWN"
	}
	var ruleVal *string
	if strings.TrimSpace(ruleName) != "" {
		trimmed := strings.TrimSpace(ruleName)
		ruleVal = &trimmed
	}

	updateQueued, updateErr := jobqueue.EnqueueUpdateUserSubscription(ctx, jobqueue.UpdateUserSubscriptionPayload{
		UserUUID:  userUUID,
		UserAgent: userAgent,
	})
	recordQueued, recordErr := jobqueue.EnqueueAddSubscriptionRequestRecord(ctx, jobqueue.AddSubscriptionRequestRecordPayload{
		UserID:          userID,
		RequestIP:       requestIP,
		UserAgent:       userAgent,
		SRRResponseType: responseType,
		SRRRuleName:     ruleVal,
	})
	if updateErr == nil && recordErr == nil && updateQueued && recordQueued {
		return
	}

	go func() {
		jobCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		_, _ = dbConn.ExecContext(jobCtx, `
			INSERT INTO user_subscription_request_history (user_id, srr_response_type, srr_rule_name, request_ip, user_agent)
			VALUES ($1, $2, $3, $4, $5)
		`, userID, responseType, ruleVal, requestIP, userAgent)

		_, _ = dbConn.ExecContext(jobCtx, `
			DELETE FROM user_subscription_request_history
			WHERE user_id = $1
			  AND id NOT IN (
				  SELECT id
				  FROM user_subscription_request_history
				  WHERE user_id = $2
				  ORDER BY request_at DESC, id DESC
				  LIMIT 24
			  )
		`, userID, userID)
	}()
}

func getSubscriptionTemplate(ctx context.Context, dbConn *sql.DB, templateType string) ([]byte, error) {
	upperType := strings.ToUpper(strings.TrimSpace(templateType))

	subTemplateLock.RLock()
	if cached, ok := subTemplateTypeCache[upperType]; ok && time.Now().Before(cached.expiresAt) {
		subTemplateLock.RUnlock()
		return cached.data, nil
	}
	subTemplateLock.RUnlock()

	row := dbConn.QueryRowContext(ctx, `
		SELECT template_yaml, template_json
		FROM subscription_templates
		WHERE UPPER(template_type) = $1
		ORDER BY view_position ASC
		LIMIT 1
	`, upperType)

	var templateYAML sql.NullString
	var templateJSON sql.NullString
	if err := row.Scan(&templateYAML, &templateJSON); err != nil {
		return nil, err
	}

	var templateData []byte
	if upperType == responseTypeXrayJSON || upperType == responseTypeSingbox {
		if templateJSON.Valid {
			templateData = []byte(templateJSON.String)
		}
	} else {
		if templateYAML.Valid {
			templateData = []byte(templateYAML.String)
		}
	}

	subTemplateLock.Lock()
	subTemplateTypeCache[upperType] = cachedTemplate{
		data:      templateData,
		expiresAt: time.Now().Add(subTemplateCacheTTL),
	}
	subTemplateLock.Unlock()

	return templateData, nil
}

func getSubscriptionTemplateByName(ctx context.Context, dbConn *sql.DB, name string) (string, []byte, error) {
	subTemplateLock.RLock()
	if cached, ok := subTemplateNameCache[name]; ok && time.Now().Before(cached.expiresAt) {
		subTemplateLock.RUnlock()
		return cached.templateType, cached.data, nil
	}
	subTemplateLock.RUnlock()

	row := dbConn.QueryRowContext(ctx, `
		SELECT template_type, template_yaml, template_json
		FROM subscription_templates
		WHERE name = $1
		LIMIT 1
	`, name)

	var templateType string
	var templateYAML sql.NullString
	var templateJSON sql.NullString
	if err := row.Scan(&templateType, &templateYAML, &templateJSON); err != nil {
		return "", nil, err
	}

	upperType := strings.ToUpper(strings.TrimSpace(templateType))
	var templateData []byte
	if upperType == responseTypeXrayJSON || upperType == responseTypeSingbox {
		if templateJSON.Valid {
			templateData = []byte(templateJSON.String)
		}
	} else {
		if templateYAML.Valid {
			templateData = []byte(templateYAML.String)
		}
	}

	subTemplateLock.Lock()
	subTemplateNameCache[name] = cachedNamedTemplate{
		templateType: templateType,
		data:         templateData,
		expiresAt:    time.Now().Add(subTemplateCacheTTL),
	}
	subTemplateLock.Unlock()

	return templateType, templateData, nil
}

func getSubscriptionTemplateByUUID(ctx context.Context, dbConn *sql.DB, uuidStr string) (string, []byte, error) {
	subTemplateLock.RLock()
	if cached, ok := subTemplateUUIDCache[uuidStr]; ok && time.Now().Before(cached.expiresAt) {
		subTemplateLock.RUnlock()
		return cached.templateType, cached.data, nil
	}
	subTemplateLock.RUnlock()

	row := dbConn.QueryRowContext(ctx, `
		SELECT template_type, template_yaml, template_json
		FROM subscription_templates
		WHERE uuid = $1
		LIMIT 1
	`, uuidStr)

	var templateType string
	var templateYAML sql.NullString
	var templateJSON sql.NullString
	if err := row.Scan(&templateType, &templateYAML, &templateJSON); err != nil {
		return "", nil, err
	}

	upperType := strings.ToUpper(strings.TrimSpace(templateType))
	var templateData []byte
	if upperType == responseTypeXrayJSON || upperType == responseTypeSingbox {
		if templateJSON.Valid {
			templateData = []byte(templateJSON.String)
		}
	} else {
		if templateYAML.Valid {
			templateData = []byte(templateYAML.String)
		}
	}

	subTemplateLock.Lock()
	subTemplateUUIDCache[uuidStr] = cachedNamedTemplate{
		templateType: templateType,
		data:         templateData,
		expiresAt:    time.Now().Add(subTemplateCacheTTL),
	}
	subTemplateLock.Unlock()

	return templateType, templateData, nil
}

func getUsersWithPagination(ctx context.Context, dbConn *sql.DB, start, size int) ([]SubscriptionUser, int, error) {
	var total int
	if err := dbConn.QueryRowContext(ctx, `SELECT COUNT(*) FROM users`).Scan(&total); err != nil {
		return nil, 0, err
	}

	rows, err := dbConn.QueryContext(ctx, `
		SELECT u.id, u.uuid, u.short_uuid, u.username, u.status,
			   u.traffic_limit_bytes, u.traffic_limit_strategy, u.expire_at,
			   u.trojan_password, u.vless_uuid, u.ss_password,
			   u.naive_password, u.shadowtls_password, u.hysteria2_password, u.anytls_password,
			   u.hwid_device_limit, u.external_squad_uuid,
			   COALESCE(ut.used_traffic_bytes, 0), COALESCE(ut.lifetime_used_traffic_bytes, 0)
		FROM users u
		LEFT JOIN user_traffic ut ON ut.id = u.id
		ORDER BY u.created_at DESC
		LIMIT $1 OFFSET $2
	`, size, start)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	users := []SubscriptionUser{}
	for rows.Next() {
		var user SubscriptionUser
		var hwidDeviceLimit sql.NullInt64
		var externalSquadUUID sql.NullString
		var naivePassword, shadowtlsPassword, hysteria2Password, anytlsPassword sql.NullString

		if err := rows.Scan(
			&user.ID,
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
			return nil, 0, err
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
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func getSubpageConfigForUser(ctx context.Context, dbConn *sql.DB, cfg *config.BackendConfig, shortUUID string, requestHeaders map[string]string) (string, bool, error) {
	log := cfg.Logger.RoleService(logger.RoleAPI, logger.ServiceHTTP)
	user, err := getSubscriptionUserByShortUUID(ctx, dbConn, shortUUID)
	if err != nil {
		return "", false, err
	}

	subpageConfigUUID := ""

	if user.ExternalSquadUUID != nil {
		var squadSubpageUUID sql.NullString

		err := dbConn.QueryRowContext(ctx, `
			SELECT subpage_config_uuid 
			FROM external_squads 
			WHERE uuid = $1`,
			*user.ExternalSquadUUID).Scan(&squadSubpageUUID)

		if err == nil && squadSubpageUUID.Valid {
			subpageConfigUUID = squadSubpageUUID.String
		} else if err != nil && !errors.Is(err, sql.ErrNoRows) {
			log.Error(fmt.Sprintf("Failed to load external squad subpage config: %v", err))
		}
	}

	if subpageConfigUUID == "" {
		subpageConfigUUID = defaultSubpageConfigUUID
	}

	settings, err := loadSubscriptionSettings(ctx, dbConn, cfg)
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

	log.Debug(
		"Resolved subpage config for user",
		"short_uuid", shortUUID,
		"subpage_config_uuid", subpageConfigUUID,
		"webpage_allowed", webpageAllowed,
	)

	return subpageConfigUUID, webpageAllowed, nil
}

func UpdateExternalSquad(ctx context.Context, dbConn *sql.DB, squadUUID string, input UpdateExternalSquadInput) error {
	var currentName string
	var currentSubpageConfigUUID sql.NullString
	var currentCustomRemarks sql.NullString
	var currentHwidSettingsRaw []byte

	err := dbConn.QueryRowContext(ctx,
		`SELECT name, subpage_config_uuid, custom_remarks, hwid_settings FROM external_squads WHERE uuid = $1`,
		squadUUID).Scan(&currentName, &currentSubpageConfigUUID, &currentCustomRemarks, &currentHwidSettingsRaw)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("external squad not found")
		}
		return fmt.Errorf("failed to fetch current external squad: %w", err)
	}

	var columns []string
	var args []interface{}
	idx := 1

	if input.Name != nil {
		columns = append(columns, fmt.Sprintf("name = $%d", idx))
		args = append(args, *input.Name)
		idx++
	}

	if input.SubpageConfigUUID != nil {
		columns = append(columns, fmt.Sprintf("subpage_config_uuid = $%d", idx))
		if *input.SubpageConfigUUID == "" {
			args = append(args, nil)
		} else {
			args = append(args, *input.SubpageConfigUUID)
		}
		idx++
	}

	if input.CustomRemarks != nil {
		columns = append(columns, fmt.Sprintf("custom_remarks = $%d", idx))
		if len(*input.CustomRemarks) == 0 || string(*input.CustomRemarks) == "null" {
			args = append(args, nil)
		} else {
			args = append(args, string(*input.CustomRemarks))
		}
		idx++
	}

	if len(input.HwidSettings) > 0 {
		raw := strings.TrimSpace(string(input.HwidSettings))
		if raw == "null" {
			columns = append(columns, fmt.Sprintf("hwid_settings = $%d", idx))
			args = append(args, nil)
			idx++
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

			columns = append(columns, fmt.Sprintf("hwid_settings = $%d", idx))
			args = append(args, string(updatedHwidRaw))
			idx++
		}
	}

	if len(columns) == 0 {
		return fmt.Errorf("no fields to update")
	}

	args = append(args, squadUUID)
	query := fmt.Sprintf("UPDATE external_squads SET %s WHERE uuid = $%d", strings.Join(columns, ", "), idx)

	_, err = dbConn.ExecContext(ctx, query, args...)
	return err
}

func applyHostOverrides(hosts []SubscriptionHost, overrides map[string]HostOverride) []SubscriptionHost {
	if len(overrides) == 0 {
		return hosts
	}
	result := make([]SubscriptionHost, len(hosts))
	for i, h := range hosts {
		override, ok := overrides[h.UUID]
		if !ok {
			result[i] = h
			continue
		}
		if override.Address != nil && strings.TrimSpace(*override.Address) != "" {
			h.Address = *override.Address
		}
		if override.Port != nil && *override.Port > 0 {
			h.Port = *override.Port
		}
		if override.Remark != nil && strings.TrimSpace(*override.Remark) != "" {
			h.Remark = *override.Remark
		}
		if override.SNI != nil && strings.TrimSpace(*override.SNI) != "" {
			h.SNI = override.SNI
		}
		if override.Host != nil && strings.TrimSpace(*override.Host) != "" {
			h.Host = override.Host
		}
		if override.Path != nil && strings.TrimSpace(*override.Path) != "" {
			h.Path = override.Path
		}
		result[i] = h
	}
	return result
}

func resolveSubscriptionBaseFromNode(ctx context.Context, dbConn *sql.DB) string {
	if dbConn == nil {
		return ""
	}

	subNodeBaseLock.RLock()
	if time.Now().Before(subNodeBaseExp) {
		val := subNodeBaseVal
		subNodeBaseLock.RUnlock()
		return val
	}
	subNodeBaseLock.RUnlock()

	var domain sql.NullString
	var apiPath sql.NullString
	row := dbConn.QueryRowContext(ctx, `
		SELECT
			COALESCE(NULLIF(BTRIM(public_domain), ''), NULLIF(BTRIM(address), '')) AS domain,
			COALESCE(NULLIF(BTRIM(api_path), ''), '/') AS api_path
		FROM sub_nodes
		ORDER BY is_disabled ASC, view_position ASC, created_at ASC
		LIMIT 1
	`)

	scanErr := row.Scan(&domain, &apiPath)
	if errors.Is(scanErr, sql.ErrNoRows) || scanErr != nil || !domain.Valid {
		subNodeBaseLock.Lock()
		subNodeBaseVal = ""
		subNodeBaseExp = time.Now().Add(subNodeBaseTTL)
		subNodeBaseLock.Unlock()
		return ""
	}

	nodeDomain := strings.TrimSpace(strings.Split(domain.String, ",")[0])
	if nodeDomain == "" {
		subNodeBaseLock.Lock()
		subNodeBaseVal = ""
		subNodeBaseExp = time.Now().Add(subNodeBaseTTL)
		subNodeBaseLock.Unlock()
		return ""
	}

	if !strings.Contains(nodeDomain, "://") {
		nodeDomain = "https://" + nodeDomain
	}

	parsedDomain, parseErr := url.Parse(nodeDomain)
	if parseErr != nil || strings.TrimSpace(parsedDomain.Host) == "" {
		subNodeBaseLock.Lock()
		subNodeBaseVal = ""
		subNodeBaseExp = time.Now().Add(subNodeBaseTTL)
		subNodeBaseLock.Unlock()
		return ""
	}

	parsedDomain.Path = ""
	parsedDomain.RawQuery = ""
	parsedDomain.Fragment = ""
	parsedDomain.User = nil

	base := strings.TrimRight(parsedDomain.String(), "/")
	path := normalizeSubscriptionAPIPath(apiPath.String)
	res := base + path

	subNodeBaseLock.Lock()
	subNodeBaseVal = res
	subNodeBaseExp = time.Now().Add(subNodeBaseTTL)
	subNodeBaseLock.Unlock()

	return res
}

func normalizeSubscriptionAPIPath(value string) string {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" || trimmed == "/" {
		return "/"
	}

	return "/" + strings.Trim(trimmed, "/") + "/"
}

func resolveSubscriptionURL(ctx context.Context, dbConn *sql.DB, user SubscriptionUser, settings SubscriptionSettingsParsed) string {
	if dbConn != nil {
		if base := resolveSubscriptionBaseFromNode(ctx, dbConn); base != "" {
			return base + user.ShortUUID
		}
	}

	domain := strings.TrimSpace(settings.Raw.Address)
	if domain == "" {
		domain = "panel.exodus.dev"
	}
	scheme := strings.TrimSpace(settings.Raw.APISchema)
	if scheme != "http" && scheme != "https" {
		scheme = "https"
	}
	apiPath := strings.Trim(strings.TrimSpace(settings.Raw.APIPath), "/")
	if apiPath == "" {
		apiPath = "api/sub"
	}
	return fmt.Sprintf("%s://%s/%s/%s", scheme, domain, apiPath, user.ShortUUID)
}

// buildResponseHeaders takes an already-resolved subscriptionURL instead of
// (ctx, dbConn) so callers that already computed it for host remarks (or for
// the raw/debug endpoint) don't pay for a second resolveSubscriptionURL DB
// round-trip per request.
func buildResponseHeaders(user SubscriptionUser, settings SubscriptionSettingsParsed, contentType string, subscriptionURL string) map[string]string {
	headers := make(map[string]string)
	if contentType != "" {
		headers["content-type"] = contentType
	}
	headers["content-disposition"] = fmt.Sprintf("attachment; filename=%s", user.Username)

	userInfo := getSubscriptionUserInfo(user)
	parts := []string{}
	for key, val := range userInfo {
		parts = append(parts, fmt.Sprintf("%s=%d", key, val))
	}
	sort.Strings(parts)
	headers["subscription-userinfo"] = strings.Join(parts, "; ")

	if refillDate := getSubscriptionRefillDate(user.TrafficLimitStrategy, user.CreatedAt); refillDate != "" {
		headers["subscription-refill-date"] = refillDate
	}

	if settings.HwidSettings.Enabled {
		headers["x-hwid-active"] = "true"
	}

	if len(settings.CustomResponseHeaders) > 0 {
		for k, v := range settings.CustomResponseHeaders {
			headers[strings.ToLower(strings.TrimSpace(k))] = formatTemplateValue(v, user, settings, subscriptionURL)
		}
	} else {
		title := settings.Raw.ProfileTitle
		if title == "" {
			title = user.Username
		} else {
			title = formatTemplateValue(title, user, settings, subscriptionURL)
		}
		headers["profile-title"] = fmt.Sprintf("base64:%s", base64.StdEncoding.EncodeToString([]byte(title)))

		if settings.Raw.SupportLink != "" {
			headers["support-url"] = settings.Raw.SupportLink
		}

		interval := settings.Raw.ProfileUpdateInterval
		if interval <= 0 {
			interval = 24
		}
		headers["profile-update-interval"] = fmt.Sprintf("%d", interval)

		if settings.Raw.HappAnnounce != "" {
			announce := formatTemplateValue(settings.Raw.HappAnnounce, user, settings, subscriptionURL)
			headers["announce"] = fmt.Sprintf("base64:%s", base64.StdEncoding.EncodeToString([]byte(announce)))
		}

		if settings.Raw.HappRouting != "" {
			headers["routing"] = settings.Raw.HappRouting
		}

		if settings.Raw.IsProfileWebpageURLEnabled {
			headers["profile-web-page-url"] = subscriptionURL
		}
	}

	for k, v := range settings.ResponseHeaders {
		headers[strings.ToLower(strings.TrimSpace(k))] = formatTemplateValue(v, user, settings, subscriptionURL)
	}
	for _, k := range settings.ResponseHeadersRemove {
		delete(headers, strings.ToLower(strings.TrimSpace(k)))
	}
	return headers
}

var templateRegex = regexp.MustCompile(`\{\{(\w+)(?::([^{}]*))?\}\}`)

func parseTemplateArgs(rawArgs string) map[string]string {
	args := make(map[string]string)
	for _, pair := range strings.Split(rawArgs, "|") {
		if idx := strings.Index(pair, "="); idx != -1 {
			args[strings.TrimSpace(pair[:idx])] = pair[idx+1:]
		}
	}
	return args
}

// formatTemplateValue is used for values that additionally support the
// exEncodeBase64:/rwEncodeBase64: prefix (currently: custom response headers).
// The actual {{VAR}} substitution lives in resolveTemplateVariables so that
// every other caller (e.g. host remarks) shares the exact same variable set
// without duplicating the switch below.
func formatTemplateValue(value string, user SubscriptionUser, settings SubscriptionSettingsParsed, subscriptionURL string) string {
	shouldBase64 := false
	if strings.HasPrefix(value, "exEncodeBase64:") {
		shouldBase64 = true
		value = strings.TrimPrefix(value, "exEncodeBase64:")
	} else if strings.HasPrefix(value, "rwEncodeBase64:") {
		shouldBase64 = true
		value = strings.TrimPrefix(value, "rwEncodeBase64:")
	}

	res := resolveTemplateVariables(value, user, settings, subscriptionURL)

	if shouldBase64 {
		return "base64:" + base64.StdEncoding.EncodeToString([]byte(res))
	}
	return res
}

// formatTemplateTrafficBytes renders a traffic value for {{TRAFFIC_USED}} /
// {{TRAFFIC_LEFT}} / {{TOTAL_TRAFFIC}} placeholders using util.FormatBytes,
// so the unit auto-scales past 1024 GiB into TiB/PiB instead of staying
// pinned to GiB.
func formatTemplateTrafficBytes(bytes int64) string {
	if bytes < 0 {
		bytes = 0
	}
	return util.FormatBytes(bytes)
}

var dayjsBracketRegex = regexp.MustCompile(`\[([^\]]*)\]`)

func convertDayjsToGoFormat(layout string) string {
	if layout == "" {
		return "02.01.2006"
	}

	// Handle bracketed escape literals like [at], [UTC] in dayjs
	var escapes []string
	layout = dayjsBracketRegex.ReplaceAllStringFunc(layout, func(m string) string {
		content := m[1 : len(m)-1]
		escapes = append(escapes, content)
		return fmt.Sprintf("\x00%d\x00", len(escapes)-1)
	})

	r := strings.NewReplacer(
		"YYYY", "2006",
		"YY", "06",
		"MMMM", "January",
		"MMM", "Jan",
		"MM", "01",
		"M", "1",
		"dddd", "Monday",
		"ddd", "Mon",
		"DD", "02",
		"D", "2",
		"HH", "15",
		"H", "15",
		"hh", "03",
		"h", "3",
		"mm", "04",
		"m", "4",
		"ss", "05",
		"s", "5",
		"SSS", ".000",
		"A", "PM",
		"a", "pm",
		"ZZ", "-0700",
		"Z", "-07:00",
	)
	res := r.Replace(layout)

	for i, esc := range escapes {
		res = strings.ReplaceAll(res, fmt.Sprintf("\x00%d\x00", i), esc)
	}

	return res
}

func formatTemplateDate(t *time.Time, args map[string]string) string {
	if t == nil || t.IsZero() {
		return ""
	}
	fmtStr := args["format"]
	layout := convertDayjsToGoFormat(fmtStr)
	return t.Format(layout)
}

func getNextTrafficResetAt(strategy string, createdAt time.Time) *time.Time {
	now := time.Now().UTC()
	createdAt = createdAt.UTC()

	switch strategy {
	case "DAY":
		startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 5, 0, 0, time.UTC)
		if startOfDay.After(now) {
			return &startOfDay
		}
		next := startOfDay.AddDate(0, 0, 1)
		return &next

	case "WEEK":
		daysUntilMonday := (1 - int(now.Weekday()) + 7) % 7
		startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 15, 0, 0, time.UTC)
		next := startOfDay.AddDate(0, 0, daysUntilMonday)
		if next.After(now) {
			return &next
		}
		nextWeek := next.AddDate(0, 0, 7)
		return &nextWeek

	case "MONTH":
		startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 20, 0, 0, time.UTC)
		if startOfMonth.After(now) {
			return &startOfMonth
		}
		nextMonth := startOfMonth.AddDate(0, 1, 0)
		return &nextMonth

	case "MONTH_ROLLING":
		anchorDay := createdAt.Day()
		daysInMonth := func(year int, month time.Month) int {
			return time.Date(year, month+1, 0, 0, 0, 0, 0, time.UTC).Day()
		}

		resetAtIn := func(t time.Time) time.Time {
			dim := daysInMonth(t.Year(), t.Month())
			day := anchorDay
			if day > dim {
				day = dim
			}
			return time.Date(t.Year(), t.Month(), day, 0, 10, 0, 0, time.UTC)
		}

		next := resetAtIn(now)
		var resolved time.Time
		if next.After(now) {
			resolved = next
		} else {
			resolved = resetAtIn(now.AddDate(0, 1, 0))
		}

		firstReset := resetAtIn(createdAt.AddDate(0, 1, 0))
		if resolved.Before(firstReset) {
			return &firstReset
		}
		return &resolved

	default:
		return nil
	}
}

// resolveTemplateVariables replaces every supported {{VAR}} / {{VAR:k=v|...}}
// placeholder in value using the user/settings context. This mirrors
// upstream's TemplateEngine.replace() for the variable set and syntax, but
// NOT for TOTAL_TRAFFIC/TRAFFIC_USED/TRAFFIC_LEFT unit scaling - see
// formatTemplateTrafficBytes. This is the single place where the list of
// supported template variables is defined - reused by both response headers
// (formatTemplateValue) and host remarks (resolveHostRemarks).
func resolveTemplateVariables(value string, user SubscriptionUser, settings SubscriptionSettingsParsed, subscriptionURL string) string {
	trafficLeft := int64(0)
	if user.TrafficLimitBytes > 0 {
		if user.TrafficLimitBytes > user.UsedTrafficBytes {
			trafficLeft = user.TrafficLimitBytes - user.UsedTrafficBytes
		}
	}

	daysLeft := int64(0)
	if !user.ExpireAt.IsZero() && user.ExpireAt.After(time.Now()) {
		daysLeft = int64(math.Max(0, time.Until(user.ExpireAt).Hours()/24))
	}

	lastResetUnix := int64(0)
	if user.LastTrafficResetAt != nil && !user.LastTrafficResetAt.IsZero() {
		lastResetUnix = user.LastTrafficResetAt.Unix()
	}

	nextResetAt := getNextTrafficResetAt(user.TrafficLimitStrategy, user.CreatedAt)
	nextResetUnix := int64(0)
	if nextResetAt != nil && !nextResetAt.IsZero() {
		nextResetUnix = nextResetAt.Unix()
	}

	createdAtUnix := int64(0)
	if !user.CreatedAt.IsZero() {
		createdAtUnix = user.CreatedAt.Unix()
	}

	expireUnix := int64(0)
	if !user.ExpireAt.IsZero() {
		expireUnix = user.ExpireAt.Unix()
	}

	email := ""
	if user.Email != nil {
		email = *user.Email
	}

	tag := ""
	if user.Tag != nil {
		tag = *user.Tag
	}

	telegramID := ""
	if user.TelegramID != nil {
		telegramID = strconv.FormatInt(*user.TelegramID, 10)
	}

	description := ""
	if user.Description != nil {
		description = *user.Description
	}

	hwidLimit := 0
	if user.HwidDeviceLimit != nil {
		hwidLimit = *user.HwidDeviceLimit
	} else {
		hwidLimit = settings.HwidSettings.FallbackDeviceLimit
	}

	userStatusLabel := ""
	if len(user.Status) > 0 {
		userStatusLabel = strings.ToUpper(user.Status[:1]) + strings.ToLower(user.Status[1:])
	}

	res := templateRegex.ReplaceAllStringFunc(value, func(match string) string {
		submatches := templateRegex.FindStringSubmatch(match)
		if len(submatches) < 2 {
			return match
		}
		key := submatches[1]
		var args map[string]string
		if len(submatches) >= 3 && submatches[2] != "" {
			args = parseTemplateArgs(submatches[2])
		} else {
			args = make(map[string]string)
		}

		switch key {
		case "DAYS_LEFT":
			return strconv.FormatInt(daysLeft, 10)
		case "TRAFFIC_USED":
			return formatTemplateTrafficBytes(user.UsedTrafficBytes)
		case "TRAFFIC_LEFT":
			return formatTemplateTrafficBytes(trafficLeft)
		case "TOTAL_TRAFFIC":
			return formatTemplateTrafficBytes(user.TrafficLimitBytes)
		case "STATUS":
			if val, ok := args[user.Status]; ok {
				return val
			}
			return userStatusLabel
		case "USERNAME":
			return user.Username
		case "EMAIL":
			return email
		case "TELEGRAM_ID":
			return telegramID
		case "SUBSCRIPTION_URL":
			return subscriptionURL
		case "TAG":
			return tag
		case "EXPIRE_UNIX":
			return strconv.FormatInt(expireUnix, 10)
		case "SHORT_UUID":
			return user.ShortUUID
		case "ID":
			return strconv.FormatInt(user.ID, 10)
		case "TRAFFIC_USED_BYTES":
			return strconv.FormatInt(user.UsedTrafficBytes, 10)
		case "TRAFFIC_LEFT_BYTES":
			return strconv.FormatInt(trafficLeft, 10)
		case "TOTAL_TRAFFIC_BYTES":
			return strconv.FormatInt(user.TrafficLimitBytes, 10)
		case "RESET_STRATEGY":
			if val, ok := args[user.TrafficLimitStrategy]; ok {
				return val
			}
			return user.TrafficLimitStrategy
		case "LIFETIME_USED_BYTES":
			return strconv.FormatInt(user.LifetimeUsedBytes, 10)
		case "CREATED_AT_UNIX":
			return strconv.FormatInt(createdAtUnix, 10)
		case "LAST_TRAFFIC_RESET_AT_UNIX":
			return strconv.FormatInt(lastResetUnix, 10)
		case "LAST_TRAFFIC_RESET_AT":
			return formatTemplateDate(user.LastTrafficResetAt, args)
		case "NEXT_TRAFFIC_RESET_AT_UNIX":
			return strconv.FormatInt(nextResetUnix, 10)
		case "NEXT_TRAFFIC_RESET_AT":
			return formatTemplateDate(nextResetAt, args)
		case "SS_HWID_LIMIT":
			return strconv.Itoa(hwidLimit)
		case "DESCRIPTION":
			return description
		default:
			return match
		}
	})

	return res
}

// resolveHostRemarks applies {{VAR}} template substitution to every host's
// Remark (in place) and deduplicates the results, exactly like upstream's
// resolve-proxy-config.service.ts does before handing hosts off to any of
// the format-specific generators (xray/singbox/mihomo). Call this once,
// right after the final host list for a user is assembled, so every
// generator - and the raw/debug JSON view - sees the same resolved remark
// without each of them re-implementing template substitution.
func resolveHostRemarks(hosts []SubscriptionHost, user SubscriptionUser, settings SubscriptionSettingsParsed, subscriptionURL string) {
	knownRemarks := make(map[string]int, len(hosts))
	for i := range hosts {
		hosts[i].Remark = deduplicateRemark(
			resolveTemplateVariables(hosts[i].Remark, user, settings, subscriptionURL),
			knownRemarks,
		)
	}
}

// deduplicateRemark ports upstream's deduplicateRemark(): if a remark was
// already seen for this user's host list, it appends a " ^~N~^" suffix so
// clients don't end up with multiple nodes sharing an identical name.
func deduplicateRemark(remark string, knownRemarks map[string]int) string {
	currentCount := knownRemarks[remark]
	knownRemarks[remark] = currentCount + 1

	if currentCount == 0 {
		return remark
	}

	hasExistingSuffix := strings.Contains(remark, "^~") && strings.HasSuffix(remark, "~^")
	suffix := currentCount + 1
	if hasExistingSuffix {
		suffix = currentCount
	}
	return fmt.Sprintf("%s ^~%d~^", remark, suffix)
}

func getSubscriptionUserInfo(user SubscriptionUser) map[string]int64 {
	expire := user.ExpireAt.Unix()
	if user.ExpireAt.IsZero() || user.ExpireAt.Year() <= 1 || user.ExpireAt.Year() == 2099 {
		expire = 0
	}

	return map[string]int64{
		"upload":   0,
		"download": user.UsedTrafficBytes,
		"total":    user.TrafficLimitBytes,
		"expire":   expire,
	}
}

func getSubscriptionRefillDate(strategy string, createdAt time.Time) string {
	next := getNextTrafficResetAt(strategy, createdAt)
	if next == nil {
		return ""
	}
	return fmt.Sprintf("%d", next.Unix())
}
