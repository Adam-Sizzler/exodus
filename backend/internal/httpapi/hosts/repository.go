package hosts

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"sort"

	dbmanager "exodus/internal/db/manager"
	"exodus/internal/dbutil"
	"exodus/internal/httpapi/shared"
)

func getHosts(ctx context.Context, manager *dbmanager.DatabaseManager) ([]hostRecord, error) {
	var hosts []hostRecord
	var err error
	err = manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		rows, err := db.QueryContext(ctx, `
            SELECT
                uuid, view_position, remark, address, port,
                path, sni, host, alpn, fingerprint, security_layer,
                xhttp_extra_params, mux_params, singbox_mux_params, clash_mux_params, sockopt_params, final_mask,
                is_disabled, server_description, override_protocol_credential, protocol_credential,
                vless_route_id, pinned_peer_cert_sha256, verify_peer_cert_by_name,
                allow_insecure, shuffle_host, mihomo_x25519, mihomo_ip_version,
                xray_json_template_uuid, keep_sni_blank,
                tags, is_hidden, override_sni_from_address,
                config_profile_uuid, config_profile_inbound_uuid,
                exclude_from_subscription_types
            FROM hosts
            ORDER BY view_position ASC
        `)
		if err != nil {
			return err
		}
		defer rows.Close()
		for rows.Next() {
			rec, scanErr := scanHostRecord(rows)
			if scanErr != nil {
				return scanErr
			}
			hosts = append(hosts, rec)
		}
		return rows.Err()
	})
	return hosts, err
}

func getHostByUUID(ctx context.Context, manager *dbmanager.DatabaseManager, hostUUID string) (hostRecord, error) {
	var host hostRecord
	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		row := db.QueryRowContext(ctx, `
            SELECT
                uuid, view_position, remark, address, port,
                path, sni, host, alpn, fingerprint, security_layer,
                xhttp_extra_params, mux_params, singbox_mux_params, clash_mux_params, sockopt_params, final_mask,
                is_disabled, server_description, override_protocol_credential, protocol_credential,
                vless_route_id, pinned_peer_cert_sha256, verify_peer_cert_by_name,
                allow_insecure, shuffle_host, mihomo_x25519, mihomo_ip_version,
                xray_json_template_uuid, keep_sni_blank,
                tags, is_hidden, override_sni_from_address,
                config_profile_uuid, config_profile_inbound_uuid,
                exclude_from_subscription_types
            FROM hosts
            WHERE uuid = ?
        `, hostUUID)
		var scanErr error
		host, scanErr = scanHostRecord(row)
		return scanErr
	})
	return host, err
}

func scanHostRecord(scanner shared.RowScanner) (hostRecord, error) {
	var rec hostRecord
	var viewPosition sql.NullInt64
	var path, sni, host, alpn, fingerprint, securityLayer sql.NullString
	var serverDescription, protocolCredential, pinnedPeerCertSha256, verifyPeerCertByName, mihomoIPVersion sql.NullString
	var xrayJSONTemplateUUID, configProfileUUID, configProfileInboundUUID sql.NullString
	var vlessRouteID sql.NullInt64
	var isDisabled, overrideProtocolCredential, allowInsecure, shuffleHost, mihomoX25519, keepSNIBlank, isHidden, overrideSNIFromAddress sql.NullBool
	var xhttpExtraParams, muxParams, singboxMuxParams, sockoptParams, finalMask []byte
	var clashMuxParams sql.NullString
	var tags, excludeTypes dbutil.StringArray

	err := scanner.Scan(
		&rec.UUID,
		&viewPosition,
		&rec.Remark,
		&rec.Address,
		&rec.Port,
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
		&finalMask,
		&isDisabled,
		&serverDescription,
		&overrideProtocolCredential,
		&protocolCredential,
		&vlessRouteID,
		&pinnedPeerCertSha256,
		&verifyPeerCertByName,
		&allowInsecure,
		&shuffleHost,
		&mihomoX25519,
		&mihomoIPVersion,
		&xrayJSONTemplateUUID,
		&keepSNIBlank,
		&tags,
		&isHidden,
		&overrideSNIFromAddress,
		&configProfileUUID,
		&configProfileInboundUUID,
		&excludeTypes,
	)
	if err != nil {
		return rec, err
	}

	if viewPosition.Valid {
		rec.ViewPosition = int(viewPosition.Int64)
	}
	if path.Valid {
		rec.Path = &path.String
	}
	if sni.Valid {
		rec.SNI = &sni.String
	}
	if host.Valid {
		rec.Host = &host.String
	}
	if alpn.Valid {
		rec.ALPN = &alpn.String
	}
	if fingerprint.Valid {
		rec.Fingerprint = &fingerprint.String
	}
	if securityLayer.Valid && securityLayer.String != "" {
		rec.SecurityLayer = securityLayer.String
	} else {
		rec.SecurityLayer = "DEFAULT"
	}
	if len(xhttpExtraParams) > 0 {
		rec.XHTTPExtraParams = json.RawMessage(xhttpExtraParams)
	}
	if len(muxParams) > 0 {
		rec.MuxParams = json.RawMessage(muxParams)
	}
	if len(singboxMuxParams) > 0 {
		rec.SingboxMuxParams = json.RawMessage(singboxMuxParams)
	}
	if clashMuxParams.Valid && clashMuxParams.String != "" {
		rec.ClashMuxParams = &clashMuxParams.String
	}
	if len(sockoptParams) > 0 {
		rec.SockoptParams = json.RawMessage(sockoptParams)
	}
	if len(finalMask) > 0 {
		rec.FinalMask = json.RawMessage(finalMask)
	}
	if isDisabled.Valid {
		rec.IsDisabled = isDisabled.Bool
	}
	if serverDescription.Valid {
		rec.ServerDescription = &serverDescription.String
	}
	if overrideProtocolCredential.Valid {
		rec.OverrideProtocolCredential = overrideProtocolCredential.Bool
	}
	if protocolCredential.Valid {
		rec.ProtocolCredential = &protocolCredential.String
	}
	if vlessRouteID.Valid {
		value := int(vlessRouteID.Int64)
		rec.VlessRouteID = &value
	}
	if pinnedPeerCertSha256.Valid {
		rec.PinnedPeerCertSha256 = &pinnedPeerCertSha256.String
	}
	if verifyPeerCertByName.Valid {
		rec.VerifyPeerCertByName = &verifyPeerCertByName.String
	}
	if allowInsecure.Valid {
		rec.AllowInsecure = allowInsecure.Bool
	}
	if shuffleHost.Valid {
		rec.ShuffleHost = shuffleHost.Bool
	}
	if mihomoX25519.Valid {
		rec.MihomoX25519 = mihomoX25519.Bool
	}
	if mihomoIPVersion.Valid {
		rec.MihomoIPVersion = &mihomoIPVersion.String
	}
	if xrayJSONTemplateUUID.Valid {
		rec.XrayJSONTemplateUUID = &xrayJSONTemplateUUID.String
	}
	if keepSNIBlank.Valid {
		rec.KeepSNIBlank = keepSNIBlank.Bool
	}
	if tags != nil {
		rec.Tags = tags.Slice()
	}
	if isHidden.Valid {
		rec.IsHidden = isHidden.Bool
	}
	if overrideSNIFromAddress.Valid {
		rec.OverrideSNIFromAddress = overrideSNIFromAddress.Bool
	}
	if configProfileUUID.Valid {
		rec.ConfigProfileUUID = &configProfileUUID.String
	}
	if configProfileInboundUUID.Valid {
		rec.ConfigProfileInboundUUID = &configProfileInboundUUID.String
	}
	if excludeTypes != nil {
		rec.ExcludeTypes = excludeTypes.Slice()
	}

	return rec, nil
}

func getHostNodes(ctx context.Context, manager *dbmanager.DatabaseManager, hostUUIDs []string) (map[string][]string, error) {
	result := make(map[string][]string)
	if len(hostUUIDs) == 0 {
		return result, nil
	}

	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		rows, err := db.QueryContext(ctx, `SELECT host_uuid, node_uuid FROM hosts_to_nodes WHERE host_uuid = ANY(?)`, hostUUIDs)
		if err != nil {
			return err
		}
		defer rows.Close()
		for rows.Next() {
			var hostUUID, nodeUUID string
			if err := rows.Scan(&hostUUID, &nodeUUID); err != nil {
				return err
			}
			result[hostUUID] = append(result[hostUUID], nodeUUID)
		}
		return rows.Err()
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func getHostExcludedSquads(ctx context.Context, manager *dbmanager.DatabaseManager, hostUUIDs []string) (map[string][]string, error) {
	result := make(map[string][]string)
	if len(hostUUIDs) == 0 {
		return result, nil
	}

	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		rows, err := db.QueryContext(ctx, `SELECT host_uuid, squad_uuid FROM internal_squad_host_exclusions WHERE host_uuid = ANY(?)`, hostUUIDs)
		if err != nil {
			return err
		}
		defer rows.Close()
		for rows.Next() {
			var hostUUID, squadUUID string
			if err := rows.Scan(&hostUUID, &squadUUID); err != nil {
				return err
			}
			result[hostUUID] = append(result[hostUUID], squadUUID)
		}
		return rows.Err()
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func replaceHostNodesTx(ctx context.Context, tx dbmanager.TxExecutor, hostUUID string, nodes []string) error {
	if _, err := tx.ExecContext(ctx, `DELETE FROM hosts_to_nodes WHERE host_uuid = ?`, hostUUID); err != nil {
		return err
	}
	for _, nodeUUID := range dedupeStrings(nodes) {
		if nodeUUID == "" {
			continue
		}
		if _, err := tx.ExecContext(ctx, `INSERT INTO hosts_to_nodes (host_uuid, node_uuid) VALUES (?, ?)`, hostUUID, nodeUUID); err != nil {
			return err
		}
	}
	return nil
}

func replaceHostExcludedSquadsTx(ctx context.Context, tx dbmanager.TxExecutor, hostUUID string, squads []string) error {
	if _, err := tx.ExecContext(ctx, `DELETE FROM internal_squad_host_exclusions WHERE host_uuid = ?`, hostUUID); err != nil {
		return err
	}
	for _, squadUUID := range dedupeStrings(squads) {
		if squadUUID == "" {
			continue
		}
		if _, err := tx.ExecContext(ctx, `INSERT INTO internal_squad_host_exclusions (host_uuid, squad_uuid) VALUES (?, ?)`, hostUUID, squadUUID); err != nil {
			return err
		}
	}
	return nil
}

func getHostTags(ctx context.Context, manager *dbmanager.DatabaseManager) ([]string, error) {
	tags := make([]string, 0)
	var err error
	err = manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		rows, err := db.QueryContext(ctx, `SELECT DISTINCT unnest(tags) AS tag FROM hosts ORDER BY tag ASC`)
		if err != nil {
			return err
		}
		defer rows.Close()
		for rows.Next() {
			var tag string
			if err := rows.Scan(&tag); err != nil {
				return err
			}
			tags = append(tags, tag)
		}
		return rows.Err()
	})
	if err != nil {
		return nil, err
	}
	sort.Strings(tags)
	return tags, nil
}

func ensureConfigProfileInbound(ctx context.Context, manager *dbmanager.DatabaseManager, profileUUID string, inboundUUID string) error {
	return manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		var exists int
		if err := db.QueryRowContext(ctx, `SELECT 1 FROM config_profiles WHERE uuid = ?`, profileUUID).Scan(&exists); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return errConfigProfileNotFound
			}
			return err
		}
		if err := db.QueryRowContext(ctx, `SELECT 1 FROM config_profile_inbounds WHERE uuid = ? AND profile_uuid = ?`, inboundUUID, profileUUID).Scan(&exists); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return errConfigProfileInboundNotFound
			}
			return err
		}
		return nil
	})
}

func ensureXrayJSONTemplate(ctx context.Context, manager *dbmanager.DatabaseManager, templateUUID string) error {
	return manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		var templateType string
		if err := db.QueryRowContext(ctx, `SELECT template_type FROM subscription_templates WHERE uuid = ?`, templateUUID).Scan(&templateType); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return errTemplateNotFound
			}
			return err
		}
		if templateType != "XRAY_JSON" {
			return errTemplateTypeNotAllowed
		}
		return nil
	})
}
