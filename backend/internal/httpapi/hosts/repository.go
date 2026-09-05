package hosts

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"exodus/internal/db"
	"exodus/internal/httpapi/shared"

	"github.com/google/uuid"
)

type HostRepository struct {
	db *sql.DB
}

func NewHostRepository(db *sql.DB) *HostRepository {
	return &HostRepository{db: db}
}

func (r *HostRepository) getHosts(ctx context.Context) ([]hostRecord, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT
			uuid, view_position, remark, address, port,
			path, sni, host, alpn, fingerprint, security_layer,
			xhttp_extra_params, mux_params, mapper, sockopt_params, final_mask,
			is_disabled, server_description,
			vless_route_id, pinned_peer_cert_sha256, verify_peer_cert_by_name,
			shuffle_host, mihomo_x25519, mihomo_ip_version,
			xray_json_template_uuid, keep_sni_blank,
			tags, is_hidden, override_sni_from_address,
			config_profile_uuid, config_profile_inbound_uuid,
			exclude_from_subscription_types, internal_squads_mode
		FROM hosts
		ORDER BY view_position ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var hosts []hostRecord
	for rows.Next() {
		rec, scanErr := scanHostRecord(rows)
		if scanErr != nil {
			return nil, scanErr
		}
		hosts = append(hosts, rec)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return hosts, nil
}

func (r *HostRepository) getHostByUUID(ctx context.Context, hostUUID string) (hostRecord, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT
			uuid, view_position, remark, address, port,
			path, sni, host, alpn, fingerprint, security_layer,
			xhttp_extra_params, mux_params, mapper, sockopt_params, final_mask,
			is_disabled, server_description,
			vless_route_id, pinned_peer_cert_sha256, verify_peer_cert_by_name,
			shuffle_host, mihomo_x25519, mihomo_ip_version,
			xray_json_template_uuid, keep_sni_blank,
			tags, is_hidden, override_sni_from_address,
			config_profile_uuid, config_profile_inbound_uuid,
			exclude_from_subscription_types, internal_squads_mode
		FROM hosts
		WHERE uuid = $1
	`, hostUUID)
	return scanHostRecord(row)
}

func scanHostRecord(scanner shared.RowScanner) (hostRecord, error) {
	var rec hostRecord
	var viewPosition sql.NullInt64
	var path, sni, host, alpn, fingerprint, securityLayer sql.NullString
	var serverDescription, pinnedPeerCertSha256, verifyPeerCertByName, mihomoIPVersion sql.NullString
	var xrayJSONTemplateUUID, configProfileUUID, configProfileInboundUUID, internalSquadsMode sql.NullString
	var vlessRouteID sql.NullInt64
	var isDisabled, shuffleHost, mihomoX25519, keepSNIBlank, isHidden, overrideSNIFromAddress sql.NullBool
	var xhttpExtraParams, muxParams, mapper, sockoptParams, finalMask []byte
	var tags, excludeTypes db.StringArray

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
		&mapper,
		&sockoptParams,
		&finalMask,
		&isDisabled,
		&serverDescription,
		&vlessRouteID,
		&pinnedPeerCertSha256,
		&verifyPeerCertByName,
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
		&internalSquadsMode,
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
	if securityLayer.Valid {
		rec.SecurityLayer = securityLayer.String
	} else {
		rec.SecurityLayer = "DEFAULT"
	}

	rec.XHTTPExtraParams = bytesToRawMessage(xhttpExtraParams)
	rec.MuxParams = bytesToRawMessage(muxParams)
	rec.Mapper = bytesToRawMessage(mapper)
	rec.SockoptParams = bytesToRawMessage(sockoptParams)
	rec.FinalMask = bytesToRawMessage(finalMask)

	if isDisabled.Valid {
		rec.IsDisabled = isDisabled.Bool
	}
	if serverDescription.Valid {
		rec.ServerDescription = &serverDescription.String
	}
	if vlessRouteID.Valid {
		v := int(vlessRouteID.Int64)
		rec.VlessRouteID = &v
	}
	if pinnedPeerCertSha256.Valid {
		rec.PinnedPeerCertSha256 = &pinnedPeerCertSha256.String
	}
	if verifyPeerCertByName.Valid {
		rec.VerifyPeerCertByName = &verifyPeerCertByName.String
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

	if internalSquadsMode.Valid && internalSquadsMode.String != "" {
		rec.InternalSquadsMode = internalSquadsMode.String
	} else {
		rec.InternalSquadsMode = "EXCLUDE"
	}

	rec.Tags = ensureStringSlice(tags.Slice())
	rec.ExcludeTypes = ensureStringSlice(excludeTypes.Slice())

	return rec, nil
}

func (r *HostRepository) getHostNodes(ctx context.Context, hostUUIDs []string) (map[string][]string, error) {
	result := make(map[string][]string, len(hostUUIDs))
	for _, id := range hostUUIDs {
		result[id] = []string{}
	}
	if len(hostUUIDs) == 0 {
		return result, nil
	}

	rows, err := r.db.QueryContext(ctx, `SELECT host_uuid, node_uuid FROM hosts_to_nodes WHERE host_uuid = ANY($1)`, hostUUIDs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var hostUUID, nodeUUID string
		if err := rows.Scan(&hostUUID, &nodeUUID); err != nil {
			return nil, err
		}
		result[hostUUID] = append(result[hostUUID], nodeUUID)
	}
	return result, rows.Err()
}

func (r *HostRepository) getHostInternalSquads(ctx context.Context, hostUUIDs []string) (map[string][]string, error) {
	result := make(map[string][]string, len(hostUUIDs))
	for _, id := range hostUUIDs {
		result[id] = []string{}
	}
	if len(hostUUIDs) == 0 {
		return result, nil
	}

	rows, err := r.db.QueryContext(ctx, `SELECT host_uuid, squad_uuid FROM internal_squad_host_links WHERE host_uuid = ANY($1)`, hostUUIDs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var hostUUID, squadUUID string
		if err := rows.Scan(&hostUUID, &squadUUID); err != nil {
			return nil, err
		}
		result[hostUUID] = append(result[hostUUID], squadUUID)
	}
	return result, rows.Err()
}

// replaceHostsNodesTx batch-replaces the node assignments for multiple hosts at
// once: one DELETE (ANY) + one multi-row INSERT (UNNEST cross-join), instead of
// one DELETE + one INSERT per node, repeated per host — a call this size used to
// make len(hostUUIDs) * len(nodeUUIDs) sequential round-trips.
func (r *HostRepository) replaceHostsNodesTx(ctx context.Context, tx *sql.Tx, hostUUIDs []string, nodeUUIDs []string) error {
	if len(hostUUIDs) == 0 {
		return nil
	}
	if _, err := tx.ExecContext(ctx, `DELETE FROM hosts_to_nodes WHERE host_uuid = ANY($1)`, hostUUIDs); err != nil {
		return err
	}

	cleanNodes := uniqueNonEmptyStrings(nodeUUIDs)
	if len(cleanNodes) == 0 {
		return nil
	}

	hostArg := make([]string, 0, len(hostUUIDs)*len(cleanNodes))
	nodeArg := make([]string, 0, len(hostUUIDs)*len(cleanNodes))
	for _, h := range hostUUIDs {
		for _, n := range cleanNodes {
			hostArg = append(hostArg, h)
			nodeArg = append(nodeArg, n)
		}
	}

	_, err := tx.ExecContext(ctx, `
		INSERT INTO hosts_to_nodes (host_uuid, node_uuid)
		SELECT unnest($1::uuid[]), unnest($2::uuid[])
	`, hostArg, nodeArg)
	return err
}

func (r *HostRepository) replaceHostNodesTx(ctx context.Context, tx *sql.Tx, hostUUID string, nodeUUIDs []string) error {
	if _, err := tx.ExecContext(ctx, `DELETE FROM hosts_to_nodes WHERE host_uuid = $1`, hostUUID); err != nil {
		return err
	}
	if len(nodeUUIDs) == 0 {
		return nil
	}

	for _, nodeUUID := range uniqueNonEmptyStrings(nodeUUIDs) {
		if _, err := tx.ExecContext(ctx, `INSERT INTO hosts_to_nodes (host_uuid, node_uuid) VALUES ($1, $2)`, hostUUID, nodeUUID); err != nil {
			return err
		}
	}
	return nil
}

// replaceHostsInternalSquadsTx is the batched counterpart of replaceHostInternalSquadsTx.
func (r *HostRepository) replaceHostsInternalSquadsTx(ctx context.Context, tx *sql.Tx, hostUUIDs []string, squadUUIDs []string) error {
	if len(hostUUIDs) == 0 {
		return nil
	}
	if _, err := tx.ExecContext(ctx, `DELETE FROM internal_squad_host_links WHERE host_uuid = ANY($1)`, hostUUIDs); err != nil {
		return err
	}

	cleanSquads := uniqueNonEmptyStrings(squadUUIDs)
	if len(cleanSquads) == 0 {
		return nil
	}

	hostArg := make([]string, 0, len(hostUUIDs)*len(cleanSquads))
	squadArg := make([]string, 0, len(hostUUIDs)*len(cleanSquads))
	for _, h := range hostUUIDs {
		for _, s := range cleanSquads {
			hostArg = append(hostArg, h)
			squadArg = append(squadArg, s)
		}
	}

	_, err := tx.ExecContext(ctx, `
		INSERT INTO internal_squad_host_links (host_uuid, squad_uuid)
		SELECT unnest($1::uuid[]), unnest($2::uuid[])
	`, hostArg, squadArg)
	return err
}

func (r *HostRepository) replaceHostInternalSquadsTx(ctx context.Context, tx *sql.Tx, hostUUID string, squadUUIDs []string) error {
	if _, err := tx.ExecContext(ctx, `DELETE FROM internal_squad_host_links WHERE host_uuid = $1`, hostUUID); err != nil {
		return err
	}
	if len(squadUUIDs) == 0 {
		return nil
	}

	for _, squadUUID := range uniqueNonEmptyStrings(squadUUIDs) {
		if _, err := tx.ExecContext(ctx, `INSERT INTO internal_squad_host_links (host_uuid, squad_uuid) VALUES ($1, $2)`, hostUUID, squadUUID); err != nil {
			return err
		}
	}
	return nil
}

func (r *HostRepository) getHostTags(ctx context.Context) ([]string, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT DISTINCT unnest(tags) AS tag FROM hosts ORDER BY tag ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []string
	for rows.Next() {
		var t string
		if err := rows.Scan(&t); err != nil {
			return nil, err
		}
		if strings.TrimSpace(t) != "" {
			tags = append(tags, t)
		}
	}
	return ensureStringSlice(tags), rows.Err()
}

func (r *HostRepository) ensureConfigProfileInbound(ctx context.Context, profileUUID, inboundUUID string) error {
	var exists int
	if err := r.db.QueryRowContext(ctx, `SELECT 1 FROM config_profiles WHERE uuid = $1`, profileUUID).Scan(&exists); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errConfigProfileNotFound
		}
		return err
	}
	if err := r.db.QueryRowContext(ctx, `SELECT 1 FROM config_profile_inbounds WHERE uuid = $1 AND profile_uuid = $2`, inboundUUID, profileUUID).Scan(&exists); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errConfigProfileInboundNotFound
		}
		return err
	}
	return nil
}

func (r *HostRepository) getInboundProtocolAndNetwork(ctx context.Context, inboundUUID string) (string, string, error) {
	var protocol string
	var network sql.NullString
	err := r.db.QueryRowContext(ctx, `SELECT type, network FROM config_profile_inbounds WHERE uuid = $1`, inboundUUID).Scan(&protocol, &network)
	return protocol, network.String, err
}

func (r *HostRepository) ensureXrayJSONTemplate(ctx context.Context, templateUUID string) error {
	var templateType string
	if err := r.db.QueryRowContext(ctx, `SELECT template_type FROM subscription_templates WHERE uuid = $1`, templateUUID).Scan(&templateType); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errTemplateNotFound
		}
		return err
	}
	if templateType != "XRAY_JSON" {
		return errTemplateTypeNotAllowed
	}
	return nil
}

func (r *HostRepository) createHost(ctx context.Context, hostUUID string, req HostCreateRequestAPI, xhttpExtra, mux, sockopt, finalMask []byte) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	mapperBytes, _ := normalizeJSONValue(req.Mapper, false)
	if len(mapperBytes) == 0 {
		mapperBytes = []byte("{}")
	}

	mode := "EXCLUDE"
	squads := req.ExcludedInternalSquads
	if req.InternalSquads != nil {
		if strings.ToUpper(strings.TrimSpace(req.InternalSquads.Mode)) == "ALLOW_ONLY" {
			mode = "ALLOW_ONLY"
		}
		squads = req.InternalSquads.Squads
	}

	_, err = tx.ExecContext(ctx, `
		INSERT INTO hosts (
			uuid, remark, address, port,
			path, sni, host, alpn, fingerprint, security_layer,
			xhttp_extra_params, mux_params,
			mapper, sockopt_params, final_mask,
			is_disabled, server_description,
			vless_route_id, pinned_peer_cert_sha256, verify_peer_cert_by_name,
			shuffle_host, mihomo_x25519, mihomo_ip_version,
			xray_json_template_uuid, keep_sni_blank,
			exclude_from_subscription_types, tags, is_hidden,
			override_sni_from_address, config_profile_uuid, config_profile_inbound_uuid,
			internal_squads_mode
		) VALUES (
			$1, $2, $3, $4,
			$5, $6, $7, $8, $9, $10,
			$11, $12, $13, $14, $15,
			$16, $17,
			$18, $19, $20,
			$21, $22, $23,
			$24, $25,
			$26, $27, $28,
			$29, $30, $31,
			$32
		)
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
		mapperBytes,
		sockopt,
		finalMask,
		coalesceBool(req.IsDisabled, false),
		normalizeOptionalStringAllowEmpty(req.ServerDescription),
		normalizeNullableInt(req.VlessRouteID),
		normalizeNullableString(req.PinnedPeerCertSha256),
		normalizeNullableString(req.VerifyPeerCertByName),
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
		mode,
	)
	if err != nil {
		return err
	}

	if err := r.replaceHostNodesTx(ctx, tx, hostUUID, req.Nodes); err != nil {
		return err
	}
	if err := r.replaceHostInternalSquadsTx(ctx, tx, hostUUID, squads); err != nil {
		return err
	}

	return tx.Commit()
}

func (r *HostRepository) updateHost(ctx context.Context, hostUUID string, clauses []string, args []any, nodes []string, squads []string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	if len(clauses) > 0 {
		args = append(args, hostUUID)
		query := fmt.Sprintf("UPDATE hosts SET %s WHERE uuid = $%d", strings.Join(clauses, ", "), len(args))
		result, err := tx.ExecContext(ctx, query, args...)
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
	}

	if nodes != nil {
		if err := r.replaceHostNodesTx(ctx, tx, hostUUID, nodes); err != nil {
			return err
		}
	}
	if squads != nil {
		if err := r.replaceHostInternalSquadsTx(ctx, tx, hostUUID, squads); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *HostRepository) deleteHost(ctx context.Context, hostUUID string) error {
	result, err := r.db.ExecContext(ctx, `DELETE FROM hosts WHERE uuid = $1`, hostUUID)
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
}

func (r *HostRepository) reorderHosts(ctx context.Context, items []reorderHostItem) error {
	if len(items) == 0 {
		return nil
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	// Single batched UPDATE via UNNEST instead of one round-trip per host.
	uuids := make([]string, len(items))
	positions := make([]int32, len(items))
	for i, item := range items {
		uuids[i] = item.UUID
		positions[i] = int32(item.ViewPosition)
	}

	if _, err := tx.ExecContext(ctx, `
		UPDATE hosts AS h
		SET view_position = v.view_position
		FROM (
			SELECT unnest($1::uuid[]) AS uuid, unnest($2::int[]) AS view_position
		) AS v
		WHERE h.uuid = v.uuid
	`, uuids, positions); err != nil {
		return err
	}

	if _, err := tx.ExecContext(ctx, `SELECT setval('hosts_view_position_seq', (SELECT COALESCE(MAX(view_position), 0) FROM hosts) + 1)`); err != nil {
		return err
	}
	return tx.Commit()
}

func (r *HostRepository) bulkUpdateHostsEnabled(ctx context.Context, uuids []string, enabled bool) error {
	_, err := r.db.ExecContext(ctx, `UPDATE hosts SET is_disabled = $1 WHERE uuid = ANY($2)`, !enabled, uuids)
	return err
}

func (r *HostRepository) bulkDeleteHosts(ctx context.Context, uuids []string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM hosts WHERE uuid = ANY($1)`, uuids)
	return err
}

func (r *HostRepository) bulkSetInbound(ctx context.Context, uuids []string, configProfileUUID, configProfileInboundUUID string) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE hosts
		SET config_profile_uuid = $1, config_profile_inbound_uuid = $2
		WHERE uuid = ANY($3)
	`, configProfileUUID, configProfileInboundUUID, uuids)
	return err
}

func (r *HostRepository) bulkSetPort(ctx context.Context, uuids []string, port int) error {
	_, err := r.db.ExecContext(ctx, `UPDATE hosts SET port = $1 WHERE uuid = ANY($2)`, port, uuids)
	return err
}

func (r *HostRepository) bulkUpdateHosts(ctx context.Context, uuids []string, clauses []string, args []any, nodes []string, squads []string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	if len(clauses) > 0 {
		args = append(args, uuids)
		query := fmt.Sprintf("UPDATE hosts SET %s WHERE uuid = ANY($%d)", strings.Join(clauses, ", "), len(args))
		if _, err := tx.ExecContext(ctx, query, args...); err != nil {
			return err
		}
	}

	if nodes != nil {
		if err := r.replaceHostsNodesTx(ctx, tx, uuids, nodes); err != nil {
			return err
		}
	}
	if squads != nil {
		if err := r.replaceHostsInternalSquadsTx(ctx, tx, uuids, squads); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *HostRepository) cloneHost(ctx context.Context, cloneFromUUID string) (hostRecord, []string, []string, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return hostRecord{}, nil, nil, err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	orig, err := r.getHostByUUID(ctx, cloneFromUUID)
	if err != nil {
		return hostRecord{}, nil, nil, err
	}

	nodesMap, err := r.getHostNodes(ctx, []string{cloneFromUUID})
	if err != nil {
		return hostRecord{}, nil, nil, err
	}
	squadsMap, err := r.getHostInternalSquads(ctx, []string{cloneFromUUID})
	if err != nil {
		return hostRecord{}, nil, nil, err
	}
	nodes := nodesMap[cloneFromUUID]
	squads := squadsMap[cloneFromUUID]

	if _, err := tx.ExecContext(ctx, `UPDATE hosts SET view_position = view_position + 1 WHERE view_position > $1`, orig.ViewPosition); err != nil {
		return hostRecord{}, nil, nil, err
	}

	cloneUUID := uuid.NewString()
	newPosition := orig.ViewPosition + 1
	newRemark := cloneString(orig.Remark)

	_, err = tx.ExecContext(ctx, `
		INSERT INTO hosts (
			uuid, view_position, remark, address, port,
			path, sni, host, alpn, fingerprint, security_layer,
			xhttp_extra_params, mux_params,
			mapper, sockopt_params, final_mask,
			is_disabled, server_description,
			vless_route_id, pinned_peer_cert_sha256, verify_peer_cert_by_name,
			shuffle_host, mihomo_x25519, mihomo_ip_version,
			xray_json_template_uuid, keep_sni_blank,
			exclude_from_subscription_types, tags, is_hidden,
			override_sni_from_address, config_profile_uuid, config_profile_inbound_uuid,
			internal_squads_mode
		) VALUES (
			$1, $2, $3, $4, $5,
			$6, $7, $8, $9, $10, $11,
			$12, $13,
			$14, $15, $16,
			$17, $18,
			$19, $20, $21,
			$22, $23, $24,
			$25, $26,
			$27, $28, $29,
			$30, $31, $32,
			$33
		)
	`,
		cloneUUID,
		newPosition,
		newRemark,
		orig.Address,
		orig.Port,
		orig.Path,
		orig.SNI,
		orig.Host,
		orig.ALPN,
		orig.Fingerprint,
		orig.SecurityLayer,
		orig.XHTTPExtraParams,
		orig.MuxParams,
		orig.Mapper,
		orig.SockoptParams,
		orig.FinalMask,
		true,
		orig.ServerDescription,
		orig.VlessRouteID,
		orig.PinnedPeerCertSha256,
		orig.VerifyPeerCertByName,
		orig.ShuffleHost,
		orig.MihomoX25519,
		orig.MihomoIPVersion,
		orig.XrayJSONTemplateUUID,
		orig.KeepSNIBlank,
		ensureStringSlice(orig.ExcludeTypes),
		ensureStringSlice(orig.Tags),
		orig.IsHidden,
		orig.OverrideSNIFromAddress,
		orig.ConfigProfileUUID,
		orig.ConfigProfileInboundUUID,
		orig.InternalSquadsMode,
	)
	if err != nil {
		return hostRecord{}, nil, nil, err
	}

	if err := r.replaceHostNodesTx(ctx, tx, cloneUUID, nodes); err != nil {
		return hostRecord{}, nil, nil, err
	}
	if err := r.replaceHostInternalSquadsTx(ctx, tx, cloneUUID, squads); err != nil {
		return hostRecord{}, nil, nil, err
	}

	if _, err := tx.ExecContext(ctx, `SELECT setval('hosts_view_position_seq', (SELECT COALESCE(MAX(view_position), 0) FROM hosts) + 1)`); err != nil {
		return hostRecord{}, nil, nil, err
	}

	if err := tx.Commit(); err != nil {
		return hostRecord{}, nil, nil, err
	}

	cloned, err := r.getHostByUUID(ctx, cloneUUID)
	return cloned, nodes, squads, err
}
