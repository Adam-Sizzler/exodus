package subscription

import (
	dbmanager "exodus/internal/db/manager"

	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"exodus/internal/config"
	"exodus/internal/dbutil"
	"exodus/internal/httpapi/shared"
	"exodus/internal/httpapi/subscriptionresponserules"
	"exodus/internal/httpapi/subscriptionsettings"
	"exodus/internal/jobqueue"
	"exodus/internal/logger"
	"fmt"
	"net/http"
	"strings"
	"time"
)

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

func scanSubscriptionHost(scanner shared.RowScanner) (SubscriptionHost, error) {
	var h SubscriptionHost
	var viewPosition sql.NullInt64
	var path, sni, host, alpn, fingerprint, securityLayer sql.NullString
	var xhttpExtraParams, muxParams, singboxMuxParams, clashMuxParams, singboxCustomParams, mihomoCustomParams, sockoptParams, serverDescription, protocolCredential sql.NullString
	var xrayJSONTemplateUUID, mihomoIPVersion, configProfileUUID, configProfileInboundUUID sql.NullString
	var inboundTag, inboundType, inboundNetwork, inboundSecurity sql.NullString
	var inboundPort sql.NullInt64
	var rawInbound sql.NullString
	var excludeTypes, hostTags dbutil.StringArray
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
