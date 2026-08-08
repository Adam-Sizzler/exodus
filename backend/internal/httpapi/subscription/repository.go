package subscription

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"

	"exodus/internal/config"
	"exodus/internal/db"
	"exodus/internal/httpapi/shared"
	"exodus/internal/httpapi/subscriptionresponserules"
	"exodus/internal/httpapi/subscriptionsettings"
	"exodus/internal/jobqueue"
	"exodus/internal/logger"
)

func loadSubscriptionSettings(ctx context.Context, dbConn *sql.DB, _ *config.BackendConfig) (SubscriptionSettingsParsed, error) {
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

	return parsed, nil
}

func loadExternalSquadOverrides(ctx context.Context, dbConn *sql.DB, squadUUID string, cfg *config.BackendConfig) (*ExternalSquadOverrides, error) {
	log := cfg.Logger.RoleService(logger.RoleAPI, logger.ServiceHTTP)
	overrides := &ExternalSquadOverrides{
		Templates: make(map[string]string),
	}

	var subscriptionSettingsJSON, hostOverridesJSON, responseHeadersJSON, hwidSettingsJSON, customRemarksJSON sql.NullString

	query := `SELECT subscription_settings, host_overrides, response_headers, hwid_settings, custom_remarks
			  FROM external_squads WHERE uuid = $1 LIMIT 1`
	row := dbConn.QueryRowContext(ctx, query, squadUUID)

	if err := row.Scan(&subscriptionSettingsJSON, &hostOverridesJSON, &responseHeadersJSON, &hwidSettingsJSON, &customRemarksJSON); err != nil {
		return nil, err
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
	} else {
		log.Warn("Failed to load external squad templates", "error", err)
	}

	return overrides, nil
}

func getSubscriptionUserByShortUUID(ctx context.Context, dbConn *sql.DB, shortUUID string) (SubscriptionUser, error) {
	return getSubscriptionUserByField(ctx, dbConn, "short_uuid", shortUUID)
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
			   u.trojan_password, u.vless_uuid, u.ss_password,
			   u.naive_password, u.shadowtls_password, u.hysteria2_password, u.anytls_password,
			   u.hwid_device_limit, u.external_squad_uuid,
			   COALESCE(ut.used_traffic_bytes, 0), COALESCE(ut.lifetime_used_traffic_bytes, 0)
		FROM users u
		LEFT JOIN user_traffic ut ON ut.id = u.id
		WHERE %s
		LIMIT 1
	`, where)

	row := dbConn.QueryRowContext(ctx, query, value)

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
		return user, err
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
	rows, err := dbConn.QueryContext(ctx, `
		SELECT DISTINCT h.uuid, h.view_position, h.remark, h.address, h.port,
			   h.path, h.sni, h.host, h.alpn, h.fingerprint, h.security_layer,
			   h.xhttp_extra_params, h.mux_params, h.singbox_mux_params, h.clash_mux_params, h.singbox_custom_params, h.mihomo_custom_params, h.sockopt_params, h.is_disabled,
			   h.server_description, h.override_protocol_credential, h.protocol_credential, h.shuffle_host,
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
		WHERE ism.user_id = $1 AND ihe.host_uuid IS NULL AND NOT COALESCE(h.is_disabled, false) AND NOT COALESCE(h.is_hidden, false)
		ORDER BY h.view_position ASC, h.remark ASC
	`, user.TID)
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
	var xhttpExtraParams, muxParams, singboxMuxParams, clashMuxParams, singboxCustomParams, mihomoCustomParams, sockoptParams, serverDescription, protocolCredential sql.NullString
	var xrayJSONTemplateUUID, mihomoIPVersion, configProfileUUID, configProfileInboundUUID sql.NullString
	var inboundTag, inboundType, inboundNetwork, inboundSecurity sql.NullString
	var inboundPort sql.NullInt64
	var rawInbound sql.NullString
	var excludeTypes, hostTags db.StringArray
	var isDisabled, overrideProtocolCredential, shuffleHost, mihomoX25519, keepSNIBlank, isHidden, overrideSNIFromAddress sql.NullBool

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
		&singboxCustomParams,
		&mihomoCustomParams,
		&sockoptParams,
		&isDisabled,
		&serverDescription,
		&overrideProtocolCredential,
		&protocolCredential,
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
	if singboxCustomParams.Valid {
		b := json.RawMessage(singboxCustomParams.String)
		h.SingboxCustomParams = &b
	}
	if mihomoCustomParams.Valid {
		h.MihomoCustomParams = &mihomoCustomParams.String
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
	if tags := hostTags.Slice(); len(tags) > 0 && strings.TrimSpace(tags[0]) != "" {
		firstTag := strings.TrimSpace(tags[0])
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

func checkHwidDeviceLimit(ctx context.Context, dbConn *sql.DB, user SubscriptionUser, hwid *HwidHeaders, settings HwidSettings) (bool, bool, bool) {
	if user.HwidDeviceLimit != nil && *user.HwidDeviceLimit == 0 {
		if hwid != nil {
			_ = enqueueOrUpsertHwidUserDevice(ctx, dbConn, user.TID, *hwid)
		}
		return true, false, false
	}

	if hwid == nil {
		return false, false, true
	}

	exists, err := hwidDeviceExists(ctx, dbConn, user.TID, hwid.Hwid)
	if err == nil && exists {
		_ = enqueueOrUpsertHwidUserDevice(ctx, dbConn, user.TID, *hwid)
		return true, false, false
	}

	count, err := countHwidDevices(ctx, dbConn, user.TID)
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

	if err := upsertHwidUserDevice(ctx, dbConn, user.TID, *hwid); err != nil {
		return false, true, false
	}

	return true, false, false
}

func countHwidDevices(ctx context.Context, dbConn *sql.DB, userID int64) (int, error) {
	var count int
	err := dbConn.QueryRowContext(ctx, `SELECT COUNT(*) FROM hwid_user_devices WHERE user_id = $1`, userID).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
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

func updateSubscriptionRequest(ctx context.Context, dbConn *sql.DB, userUUID string, userID int64, userAgent, requestIP string) {
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

		_, _ = dbConn.ExecContext(jobCtx, `
			INSERT INTO user_subscription_request_history (user_id, request_ip, user_agent)
			VALUES ($1, $2, $3)
		`, userID, requestIP, userAgent)

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
	return templateData, nil
}

func getSubscriptionTemplateByName(ctx context.Context, dbConn *sql.DB, name string) (string, []byte, error) {
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

func buildResponseHeaders(user SubscriptionUser, settings SubscriptionSettingsParsed, contentType string) map[string]string {
	headers := make(map[string]string)
	headers["content-disposition"] = fmt.Sprintf("attachment; filename=%s", user.Username)

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
	subscriptionURL := fmt.Sprintf("%s://%s/%s/%s", scheme, domain, apiPath, user.ShortUUID)

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

	if settings.Raw.HappRouting != "" {
		headers["routing"] = settings.Raw.HappRouting
	}

	if settings.Raw.IsProfileWebpageURLEnabled {
		headers["profile-web-page-url"] = subscriptionURL
	}

	if refillDate := getSubscriptionRefillDate(user.TrafficLimitStrategy); refillDate != "" {
		headers["subscription-refill-date"] = refillDate
	}

	for k, v := range settings.CustomResponseHeaders {
		headers[k] = formatTemplateValue(v, user, settings, subscriptionURL)
	}
	for k, v := range settings.ResponseHeaders {
		headers[k] = formatTemplateValue(v, user, settings, subscriptionURL)
	}
	return headers
}

func formatTemplateValue(value string, user SubscriptionUser, _ SubscriptionSettingsParsed, subscriptionURL string) string {
	shouldBase64 := false
	if strings.HasPrefix(value, "exEncodeBase64:") {
		shouldBase64 = true
		value = strings.TrimPrefix(value, "exEncodeBase64:")
	}

	replacer := strings.NewReplacer(
		"{USERNAME}", user.Username,
		"{{username}}", user.Username,
		"{SHORT_UUID}", user.ShortUUID,
		"{{shortUuid}}", user.ShortUUID,
		"{SUBSCRIPTION_URL}", subscriptionURL,
		"{{subscriptionUrl}}", subscriptionURL,
	)
	res := replacer.Replace(value)

	if shouldBase64 {
		return "base64:" + base64.StdEncoding.EncodeToString([]byte(res))
	}
	return res
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

func getSubscriptionRefillDate(strategy string) string {
	now := time.Now().Local()
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

func firstHostTag(host SubscriptionHost) string {
	if host.Remark != "" {
		return host.Remark
	}
	return host.Address
}
