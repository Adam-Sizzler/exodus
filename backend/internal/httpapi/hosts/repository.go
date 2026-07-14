package hosts

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strings"

	dbmanager "exodus/internal/db/manager"
	"exodus/internal/dbutil"
	"exodus/internal/httpapi/shared"
)

type HostRepository struct {
	manager *dbmanager.DatabaseManager
}

func NewHostRepository(manager *dbmanager.DatabaseManager) *HostRepository {
	return &HostRepository{manager: manager}
}

func (r *HostRepository) getHosts(ctx context.Context) ([]hostRecord, error) {
	var hosts []hostRecord
	var err error
	err = r.manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		rows, err := db.QueryContext(ctx, `
            SELECT
                uuid, view_position, remark, address, port,
                path, sni, host, alpn, fingerprint, security_layer,
                xhttp_extra_params, mux_params, singbox_mux_params, clash_mux_params, singbox_custom_params, mihomo_custom_params, sockopt_params, final_mask,
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

func (r *HostRepository) getHostByUUID(ctx context.Context, hostUUID string) (hostRecord, error) {
	var host hostRecord
	err := r.manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		row := db.QueryRowContext(ctx, `
            SELECT
                uuid, view_position, remark, address, port,
                path, sni, host, alpn, fingerprint, security_layer,
                xhttp_extra_params, mux_params, singbox_mux_params, clash_mux_params, singbox_custom_params, mihomo_custom_params, sockopt_params, final_mask,
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
	var xhttpExtraParams, muxParams, singboxMuxParams, singboxCustomParams, sockoptParams, finalMask []byte
	var clashMuxParams, mihomoCustomParams sql.NullString
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
		&singboxCustomParams,
		&mihomoCustomParams,
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
	if len(singboxCustomParams) > 0 {
		rec.SingboxCustomParams = json.RawMessage(singboxCustomParams)
	}
	if mihomoCustomParams.Valid && mihomoCustomParams.String != "" {
		rec.MihomoCustomParams = &mihomoCustomParams.String
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

func (r *HostRepository) getHostNodes(ctx context.Context, hostUUIDs []string) (map[string][]string, error) {
	result := make(map[string][]string)
	if len(hostUUIDs) == 0 {
		return result, nil
	}

	err := r.manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
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

func (r *HostRepository) getHostExcludedSquads(ctx context.Context, hostUUIDs []string) (map[string][]string, error) {
	result := make(map[string][]string)
	if len(hostUUIDs) == 0 {
		return result, nil
	}

	err := r.manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
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

func (r *HostRepository) replaceHostNodesTx(ctx context.Context, tx dbmanager.TxExecutor, hostUUID string, nodes []string) error {
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

func (r *HostRepository) replaceHostExcludedSquadsTx(ctx context.Context, tx dbmanager.TxExecutor, hostUUID string, squads []string) error {
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

func (r *HostRepository) getHostTags(ctx context.Context) ([]string, error) {
	tags := make([]string, 0)
	var err error
	err = r.manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
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

func (r *HostRepository) ensureConfigProfileInbound(ctx context.Context, profileUUID string, inboundUUID string) error {
	return r.manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
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

func (r *HostRepository) getInboundProtocolAndNetwork(ctx context.Context, inboundUUID string) (string, string, error) {
	var protocol string
	var network sql.NullString
	err := r.manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		return db.QueryRowContext(ctx, `SELECT type, network FROM config_profile_inbounds WHERE uuid = ?`, inboundUUID).Scan(&protocol, &network)
	})
	return protocol, network.String, err
}

func (r *HostRepository) ensureXrayJSONTemplate(ctx context.Context, templateUUID string) error {
	return r.manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
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

func (r *HostRepository) createHost(ctx context.Context, hostUUID string, req HostCreateRequestAPI, xhttpExtra, mux, singboxMux []byte, clashMux *string, sockopt, finalMask []byte) error {
	return r.manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		tx, err := db.BeginTx(ctx, nil)
		if err != nil {
			return err
		}

		_, err = tx.ExecContext(ctx, `
            INSERT INTO hosts (
                uuid, remark, address, port,
                path, sni, host, alpn, fingerprint, security_layer,
                xhttp_extra_params, mux_params, singbox_mux_params, clash_mux_params, sockopt_params, final_mask,
                is_disabled, server_description, override_protocol_credential, protocol_credential,
                vless_route_id, pinned_peer_cert_sha256, verify_peer_cert_by_name,
                allow_insecure, shuffle_host, mihomo_x25519, mihomo_ip_version,
                xray_json_template_uuid, keep_sni_blank,
                exclude_from_subscription_types, tags, is_hidden,
                override_sni_from_address, config_profile_uuid, config_profile_inbound_uuid
            ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
        `,
			hostUUID,
			req.Remark,
			strings.TrimSpace(req.Address),
			req.Port,
			normalizeOptionalStringAllowEmpty(req.Path),
			normalizeOptionalStringAllowEmpty(req.SNI),
			normalizeOptionalStringAllowEmpty(req.Host),
			normalizeOptionalStringAllowEmpty(req.ALPN),
			normalizeOptionalStringAllowEmpty(req.Fingerprint),
			normalizeSecurityLayer(req.SecurityLayer),
			xhttpExtra,
			mux,
			singboxMux,
			clashMux,
			sockopt,
			finalMask,
			coalesceBool(req.IsDisabled, false),
			normalizeOptionalStringAllowEmpty(req.ServerDescription),
			coalesceBool(req.OverrideProtocolCredential, false),
			normalizeProtocolCredentialForCreate(req.OverrideProtocolCredential, req.ProtocolCredential),
			normalizeNullableInt(req.VlessRouteID),
			normalizeNullableString(req.PinnedPeerCertSha256),
			normalizeNullableString(req.VerifyPeerCertByName),
			coalesceBool(req.AllowInsecure, false),
			coalesceBool(req.ShuffleHost, false),
			coalesceBool(req.MihomoX25519, false),
			normalizeMihomoIPVersion(req.MihomoIPVersion),
			normalizeOptionalStringAllowEmpty(req.XrayJSONTemplateUUID),
			coalesceBool(req.KeepSNIBlank, false),
			ensureStringSlice(req.ExcludeFromSubscription),
			normalizeTags(req.Tags),
			coalesceBool(req.IsHidden, false),
			coalesceBool(req.OverrideSNIFromAddress, false),
			normalizeOptionalStringAllowEmpty(req.Inbound.ConfigProfileUUID),
			normalizeOptionalStringAllowEmpty(req.Inbound.ConfigProfileInboundUUID),
		)
		if err != nil {
			_ = tx.Rollback()
			return err
		}

		if err := r.replaceHostNodesTx(ctx, tx, hostUUID, req.Nodes); err != nil {
			_ = tx.Rollback()
			return err
		}
		if err := r.replaceHostExcludedSquadsTx(ctx, tx, hostUUID, req.ExcludedInternalSquads); err != nil {
			_ = tx.Rollback()
			return err
		}

		return tx.Commit()
	})
}

func (r *HostRepository) updateHost(ctx context.Context, hostUUID string, clauses []string, args []any, nodes []string, excludedInternalSquads []string) error {
	return r.manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		tx, err := db.BeginTx(ctx, nil)
		if err != nil {
			return err
		}

		if len(clauses) > 0 {
			args = append(args, hostUUID)
			query := fmt.Sprintf("UPDATE hosts SET %s WHERE uuid = ?", strings.Join(clauses, ", "))
			result, err := tx.ExecContext(ctx, query, args...)
			if err != nil {
				_ = tx.Rollback()
				return err
			}
			rows, err := result.RowsAffected()
			if err != nil {
				_ = tx.Rollback()
				return err
			}
			if rows == 0 {
				_ = tx.Rollback()
				return sql.ErrNoRows
			}
		}

		if nodes != nil {
			if err := r.replaceHostNodesTx(ctx, tx, hostUUID, nodes); err != nil {
				_ = tx.Rollback()
				return err
			}
		}
		if excludedInternalSquads != nil {
			if err := r.replaceHostExcludedSquadsTx(ctx, tx, hostUUID, excludedInternalSquads); err != nil {
				_ = tx.Rollback()
				return err
			}
		}

		return tx.Commit()
	})
}

func (r *HostRepository) deleteHost(ctx context.Context, hostUUID string) error {
	return r.manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		result, err := db.ExecContext(ctx, `DELETE FROM hosts WHERE uuid = ?`, hostUUID)
		if err != nil {
			return err
		}
		rows, err := result.RowsAffected()
		if err != nil {
			return err
		}
		if rows == 0 {
			return sql.ErrNoRows
		}
		return nil
	})
}

func (r *HostRepository) reorderHosts(ctx context.Context, items []reorderHostItem) error {
	return r.manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		tx, err := db.BeginTx(ctx, nil)
		if err != nil {
			return err
		}
		for _, item := range items {
			if _, err := tx.ExecContext(ctx, `UPDATE hosts SET view_position = ? WHERE uuid = ?`, item.ViewPosition, item.UUID); err != nil {
				_ = tx.Rollback()
				return err
			}
		}
		if _, err := tx.ExecContext(ctx, `SELECT setval('hosts_view_position_seq', (SELECT COALESCE(MAX(view_position), 0) FROM hosts) + 1)`); err != nil {
			_ = tx.Rollback()
			return err
		}
		return tx.Commit()
	})
}

func (r *HostRepository) bulkUpdateHostsEnabled(ctx context.Context, uuids []string, enabled bool) error {
	return r.manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		_, err := db.ExecContext(ctx, `UPDATE hosts SET is_disabled = ? WHERE uuid = ANY(?)`, !enabled, uuids)
		return err
	})
}

func (r *HostRepository) bulkDeleteHosts(ctx context.Context, uuids []string) error {
	return r.manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		_, err := db.ExecContext(ctx, `DELETE FROM hosts WHERE uuid = ANY(?)`, uuids)
		return err
	})
}

func (r *HostRepository) bulkSetInbound(ctx context.Context, uuids []string, configProfileUUID, configProfileInboundUUID string) error {
	return r.manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		_, err := db.ExecContext(ctx, `
            UPDATE hosts
            SET config_profile_uuid = ?, config_profile_inbound_uuid = ?
            WHERE uuid = ANY(?)
        `, configProfileUUID, configProfileInboundUUID, uuids)
		return err
	})
}

func (r *HostRepository) bulkSetPort(ctx context.Context, uuids []string, port int) error {
	return r.manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		_, err := db.ExecContext(ctx, `UPDATE hosts SET port = ? WHERE uuid = ANY(?)`, port, uuids)
		return err
	})
}

func (r *HostRepository) bulkUpdateHosts(ctx context.Context, uuids []string, clauses []string, args []any, nodes []string, excludedInternalSquads []string) error {
	return r.manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		tx, err := db.BeginTx(ctx, nil)
		if err != nil {
			return err
		}

		if len(clauses) > 0 {
			args = append(args, uuids)
			query := fmt.Sprintf("UPDATE hosts SET %s WHERE uuid = ANY(?)", strings.Join(clauses, ", "))
			if _, err := tx.ExecContext(ctx, query, args...); err != nil {
				_ = tx.Rollback()
				return err
			}
		}

		if nodes != nil {
			for _, hostUUID := range uuids {
				if err := r.replaceHostNodesTx(ctx, tx, hostUUID, nodes); err != nil {
					_ = tx.Rollback()
					return err
				}
			}
		}
		if excludedInternalSquads != nil {
			for _, hostUUID := range uuids {
				if err := r.replaceHostExcludedSquadsTx(ctx, tx, hostUUID, excludedInternalSquads); err != nil {
					_ = tx.Rollback()
					return err
				}
			}
		}

		return tx.Commit()
	})
}
