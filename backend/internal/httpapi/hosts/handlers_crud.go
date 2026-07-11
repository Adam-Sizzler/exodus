package hosts

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"exodus/internal/config"
	dbmanager "exodus/internal/db/manager"
	"exodus/internal/httpapi/shared"

	"github.com/google/uuid"
)

func handleGetHosts(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) {
	records, err := getHosts(r.Context(), manager)
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to fetch hosts", err, cfg)
		return
	}

	ids := make([]string, 0, len(records))
	for _, rec := range records {
		ids = append(ids, rec.UUID)
	}

	nodesMap, err := getHostNodes(r.Context(), manager, ids)
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to fetch host nodes", err, cfg)
		return
	}
	excludedMap, err := getHostExcludedSquads(r.Context(), manager, ids)
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to fetch host exclusions", err, cfg)
		return
	}

	result := make([]HostAPI, 0, len(records))
	for _, rec := range records {
		result = append(result, mapHostRecordToAPI(rec, nodesMap[rec.UUID], excludedMap[rec.UUID]))
	}

	shared.WriteJSON(w, http.StatusOK, map[string]any{"response": result})
}

func handleGetHost(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig, hostUUID string) {
	rec, err := getHostByUUID(r.Context(), manager, hostUUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			shared.SendError(w, http.StatusNotFound, "host not found", nil, cfg)
			return
		}
		shared.SendError(w, http.StatusInternalServerError, "failed to fetch host", err, cfg)
		return
	}

	nodesMap, err := getHostNodes(r.Context(), manager, []string{hostUUID})
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to fetch host nodes", err, cfg)
		return
	}
	excludedMap, err := getHostExcludedSquads(r.Context(), manager, []string{hostUUID})
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to fetch host exclusions", err, cfg)
		return
	}

	shared.WriteJSON(w, http.StatusOK, map[string]any{"response": mapHostRecordToAPI(rec, nodesMap[rec.UUID], excludedMap[rec.UUID])})
}

func handleCreateHost(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) {
	var req HostCreateRequestAPI
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		shared.SendError(w, http.StatusBadRequest, "invalid JSON", err, cfg)
		return
	}
	if err := validateCreateRequest(req); err != nil {
		shared.SendError(w, http.StatusBadRequest, err.Error(), nil, cfg)
		return
	}
	if err := ensureConfigProfileInbound(r.Context(), manager, *req.Inbound.ConfigProfileUUID, *req.Inbound.ConfigProfileInboundUUID); err != nil {
		switch {
		case errors.Is(err, errConfigProfileNotFound):
			shared.SendError(w, http.StatusBadRequest, "config profile not found", nil, cfg)
			return
		case errors.Is(err, errConfigProfileInboundNotFound):
			shared.SendError(w, http.StatusBadRequest, "config profile inbound not found in specified profile", nil, cfg)
			return
		default:
			shared.SendError(w, http.StatusInternalServerError, "failed to validate config profile inbound", err, cfg)
			return
		}
	}
	if req.XrayJSONTemplateUUID != nil && *req.XrayJSONTemplateUUID != "" {
		if err := ensureXrayJSONTemplate(r.Context(), manager, *req.XrayJSONTemplateUUID); err != nil {
			switch {
			case errors.Is(err, errTemplateNotFound):
				shared.SendError(w, http.StatusBadRequest, "subscription template not found", nil, cfg)
				return
			case errors.Is(err, errTemplateTypeNotAllowed):
				shared.SendError(w, http.StatusBadRequest, "template type not allowed", nil, cfg)
				return
			default:
				shared.SendError(w, http.StatusInternalServerError, "failed to validate subscription template", err, cfg)
				return
			}
		}
	}

	hostUUID := uuid.NewString()

	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		tx, err := db.BeginTx(r.Context(), nil)
		if err != nil {
			return err
		}

		xhttpExtra, err := normalizeJSONValue(req.XHTTPExtraParams, true)
		if err != nil {
			_ = tx.Rollback()
			return err
		}
		mux, err := normalizeJSONValue(req.MuxParams, true)
		if err != nil {
			_ = tx.Rollback()
			return err
		}
		singboxMux, err := normalizeJSONValue(req.SingboxMuxParams, true)
		if err != nil {
			_ = tx.Rollback()
			return err
		}
		clashMux, err := normalizeClashMuxYAML(req.ClashMuxParams)
		if err != nil {
			_ = tx.Rollback()
			return err
		}
		sockopt, err := normalizeJSONValue(req.SockoptParams, true)
		if err != nil {
			_ = tx.Rollback()
			return err
		}
		finalMask, err := normalizeJSONValue(req.FinalMask, true)
		if err != nil {
			_ = tx.Rollback()
			return err
		}

		_, err = tx.ExecContext(r.Context(), `
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

		if err := replaceHostNodesTx(r.Context(), tx, hostUUID, req.Nodes); err != nil {
			_ = tx.Rollback()
			return err
		}
		if err := replaceHostExcludedSquadsTx(r.Context(), tx, hostUUID, req.ExcludedInternalSquads); err != nil {
			_ = tx.Rollback()
			return err
		}

		return tx.Commit()
	})
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to create host", err, cfg)
		return
	}

	created, err := getHostByUUID(r.Context(), manager, hostUUID)
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to fetch created host", err, cfg)
		return
	}
	result := mapHostRecordToAPI(created, req.Nodes, req.ExcludedInternalSquads)
	shared.WriteJSON(w, http.StatusCreated, map[string]any{"response": result})
}

// buildHostUpdateClauses turns a partial hostUpdateFields payload into the
// list of "column = ?" SQL clauses (and matching args) to apply to the
// hosts table. It is pure/DB-transaction-agnostic on purpose so both the
// single-host update handler and the bulk-update handler can share it
// without either owning a *sql.Tx.
func buildHostUpdateClauses(fields hostUpdateFields) ([]string, []any, error) {
	clauses := make([]string, 0)
	args := make([]any, 0)
	add := func(column string, value any) {
		clauses = append(clauses, fmt.Sprintf("%s = ?", column))
		args = append(args, value)
	}
	addOptionalString := func(column string, value *string) {
		if value == nil {
			return
		}
		add(column, strings.TrimSpace(*value))
	}

	if fields.Remark.Set {
		if fields.Remark.Value == nil {
			return nil, nil, fmt.Errorf("remark cannot be null")
		}
		add("remark", strings.TrimSpace(*fields.Remark.Value))
	}
	if fields.Address.Set {
		if fields.Address.Value == nil {
			return nil, nil, fmt.Errorf("address cannot be null")
		}
		add("address", strings.TrimSpace(*fields.Address.Value))
	}
	if fields.Port != nil {
		add("port", *fields.Port)
	}
	if fields.Path.Set {
		if fields.Path.Value == nil {
			return nil, nil, fmt.Errorf("path cannot be null")
		}
		add("path", strings.TrimSpace(*fields.Path.Value))
	}
	if fields.SNI.Set {
		if fields.SNI.Value == nil {
			return nil, nil, fmt.Errorf("sni cannot be null")
		}
		add("sni", strings.TrimSpace(*fields.SNI.Value))
	}
	if fields.Host.Set {
		if fields.Host.Value == nil {
			return nil, nil, fmt.Errorf("host cannot be null")
		}
		add("host", strings.TrimSpace(*fields.Host.Value))
	}
	if fields.ALPN.Set {
		if fields.ALPN.Value == nil {
			clauses = append(clauses, "alpn = NULL")
		} else {
			add("alpn", strings.TrimSpace(*fields.ALPN.Value))
		}
	}
	if fields.Fingerprint.Set {
		if fields.Fingerprint.Value == nil {
			clauses = append(clauses, "fingerprint = NULL")
		} else {
			add("fingerprint", strings.TrimSpace(*fields.Fingerprint.Value))
		}
	}
	if fields.SecurityLayer != nil {
		add("security_layer", normalizeSecurityLayer(fields.SecurityLayer))
	}

	if set, val, err := normalizeOptionalJSONField(fields.XHTTPExtraParams, true); err != nil {
		return nil, nil, err
	} else if set {
		if val == nil {
			clauses = append(clauses, "xhttp_extra_params = NULL")
		} else {
			add("xhttp_extra_params", val)
		}
	}
	if set, val, err := normalizeOptionalJSONField(fields.MuxParams, true); err != nil {
		return nil, nil, err
	} else if set {
		if val == nil {
			clauses = append(clauses, "mux_params = NULL")
		} else {
			add("mux_params", val)
		}
	}
	if set, val, err := normalizeOptionalJSONField(fields.SingboxMuxParams, true); err != nil {
		return nil, nil, err
	} else if set {
		if val == nil {
			clauses = append(clauses, "singbox_mux_params = NULL")
		} else {
			add("singbox_mux_params", val)
		}
	}
	if set, val, err := normalizeOptionalClashMuxYAML(fields.ClashMuxParams); err != nil {
		return nil, nil, err
	} else if set {
		if val == nil {
			clauses = append(clauses, "clash_mux_params = NULL")
		} else {
			add("clash_mux_params", val)
		}
	}
	if set, val, err := normalizeOptionalJSONField(fields.SockoptParams, true); err != nil {
		return nil, nil, err
	} else if set {
		if val == nil {
			clauses = append(clauses, "sockopt_params = NULL")
		} else {
			add("sockopt_params", val)
		}
	}
	if set, val, err := normalizeOptionalJSONField(fields.FinalMask, true); err != nil {
		return nil, nil, err
	} else if set {
		if val == nil {
			clauses = append(clauses, "final_mask = NULL")
		} else {
			add("final_mask", val)
		}
	}

	if fields.IsDisabled != nil {
		add("is_disabled", *fields.IsDisabled)
	}
	if fields.ServerDescription.Set {
		if fields.ServerDescription.Value == nil {
			clauses = append(clauses, "server_description = NULL")
		} else {
			add("server_description", strings.TrimSpace(*fields.ServerDescription.Value))
		}
	}
	protocolCredentialCleared := false
	if fields.OverrideProtocolCredential != nil {
		add("override_protocol_credential", *fields.OverrideProtocolCredential)
		if !*fields.OverrideProtocolCredential {
			clauses = append(clauses, "protocol_credential = NULL")
			protocolCredentialCleared = true
		}
	}
	if fields.ProtocolCredential.Set && !protocolCredentialCleared {
		normalizedCredential := normalizeProtocolCredentialPointer(fields.ProtocolCredential.Value)
		if normalizedCredential == nil {
			clauses = append(clauses, "protocol_credential = NULL")
		} else {
			add("protocol_credential", *normalizedCredential)
		}
	}
	if fields.VlessRouteID.Set {
		if fields.VlessRouteID.Value == nil {
			clauses = append(clauses, "vless_route_id = NULL")
		} else {
			add("vless_route_id", *fields.VlessRouteID.Value)
		}
	}
	if fields.PinnedPeerCertSha256.Set {
		if fields.PinnedPeerCertSha256.Value == nil {
			clauses = append(clauses, "pinned_peer_cert_sha256 = NULL")
		} else {
			add("pinned_peer_cert_sha256", strings.TrimSpace(*fields.PinnedPeerCertSha256.Value))
		}
	}
	if fields.VerifyPeerCertByName.Set {
		if fields.VerifyPeerCertByName.Value == nil {
			clauses = append(clauses, "verify_peer_cert_by_name = NULL")
		} else {
			add("verify_peer_cert_by_name", strings.TrimSpace(*fields.VerifyPeerCertByName.Value))
		}
	}
	if fields.AllowInsecure != nil {
		add("allow_insecure", *fields.AllowInsecure)
	}
	if fields.ShuffleHost != nil {
		add("shuffle_host", *fields.ShuffleHost)
	}
	if fields.MihomoX25519 != nil {
		add("mihomo_x25519", *fields.MihomoX25519)
	}
	if fields.MihomoIPVersion.Set {
		if fields.MihomoIPVersion.Value == nil {
			clauses = append(clauses, "mihomo_ip_version = NULL")
		} else {
			add("mihomo_ip_version", normalizeMihomoIPVersion(fields.MihomoIPVersion.Value))
		}
	}
	if fields.XrayJSONTemplateUUID.Set {
		if fields.XrayJSONTemplateUUID.Value == nil {
			clauses = append(clauses, "xray_json_template_uuid = NULL")
		} else {
			add("xray_json_template_uuid", strings.TrimSpace(*fields.XrayJSONTemplateUUID.Value))
		}
	}
	if fields.KeepSNIBlank != nil {
		add("keep_sni_blank", *fields.KeepSNIBlank)
	}
	if fields.Tags != nil {
		add("tags", normalizeTags(fields.Tags))
	}
	if fields.IsHidden != nil {
		add("is_hidden", *fields.IsHidden)
	}
	if fields.OverrideSNIFromAddress != nil {
		add("override_sni_from_address", *fields.OverrideSNIFromAddress)
	}
	if fields.Inbound != nil {
		addOptionalString("config_profile_uuid", fields.Inbound.ConfigProfileUUID)
		addOptionalString("config_profile_inbound_uuid", fields.Inbound.ConfigProfileInboundUUID)
	}
	if fields.ExcludeFromSubscription != nil {
		add("exclude_from_subscription_types", fields.ExcludeFromSubscription)
	}

	return clauses, args, nil
}

func handleUpdateHost(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) {
	var req HostUpdateRequestAPI
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		shared.SendError(w, http.StatusBadRequest, "invalid JSON", err, cfg)
		return
	}
	if _, err := uuid.Parse(req.UUID); err != nil {
		shared.SendError(w, http.StatusBadRequest, "invalid UUID format", nil, cfg)
		return
	}
	if err := validateUpdateRequest(req.hostUpdateFields); err != nil {
		shared.SendError(w, http.StatusBadRequest, err.Error(), nil, cfg)
		return
	}
	if req.Inbound != nil {
		if err := ensureConfigProfileInbound(r.Context(), manager, *req.Inbound.ConfigProfileUUID, *req.Inbound.ConfigProfileInboundUUID); err != nil {
			switch {
			case errors.Is(err, errConfigProfileNotFound):
				shared.SendError(w, http.StatusBadRequest, "config profile not found", nil, cfg)
				return
			case errors.Is(err, errConfigProfileInboundNotFound):
				shared.SendError(w, http.StatusBadRequest, "config profile inbound not found in specified profile", nil, cfg)
				return
			default:
				shared.SendError(w, http.StatusInternalServerError, "failed to validate config profile inbound", err, cfg)
				return
			}
		}
	}
	if req.XrayJSONTemplateUUID.Set && req.XrayJSONTemplateUUID.Value != nil && *req.XrayJSONTemplateUUID.Value != "" {
		if err := ensureXrayJSONTemplate(r.Context(), manager, *req.XrayJSONTemplateUUID.Value); err != nil {
			switch {
			case errors.Is(err, errTemplateNotFound):
				shared.SendError(w, http.StatusBadRequest, "subscription template not found", nil, cfg)
				return
			case errors.Is(err, errTemplateTypeNotAllowed):
				shared.SendError(w, http.StatusBadRequest, "template type not allowed", nil, cfg)
				return
			default:
				shared.SendError(w, http.StatusInternalServerError, "failed to validate subscription template", err, cfg)
				return
			}
		}
	}

	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		tx, err := db.BeginTx(r.Context(), nil)
		if err != nil {
			return err
		}

		clauses, args, err := buildHostUpdateClauses(req.hostUpdateFields)
		if err != nil {
			_ = tx.Rollback()
			return err
		}

		if len(clauses) > 0 {
			args = append(args, req.UUID)
			query := fmt.Sprintf("UPDATE hosts SET %s WHERE uuid = ?", strings.Join(clauses, ", "))
			result, err := tx.ExecContext(r.Context(), query, args...)
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

		if req.Nodes != nil {
			if err := replaceHostNodesTx(r.Context(), tx, req.UUID, req.Nodes); err != nil {
				_ = tx.Rollback()
				return err
			}
		}
		if req.ExcludedInternalSquads != nil {
			if err := replaceHostExcludedSquadsTx(r.Context(), tx, req.UUID, req.ExcludedInternalSquads); err != nil {
				_ = tx.Rollback()
				return err
			}
		}

		return tx.Commit()
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			shared.SendError(w, http.StatusNotFound, "host not found", nil, cfg)
			return
		}
		shared.SendError(w, http.StatusInternalServerError, "failed to update host", err, cfg)
		return
	}

	host, err := getHostByUUID(r.Context(), manager, req.UUID)
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to fetch updated host", err, cfg)
		return
	}

	nodesMap, _ := getHostNodes(r.Context(), manager, []string{req.UUID})
	excludedMap, _ := getHostExcludedSquads(r.Context(), manager, []string{req.UUID})

	shared.WriteJSON(w, http.StatusOK, map[string]any{"response": mapHostRecordToAPI(host, nodesMap[req.UUID], excludedMap[req.UUID])})
}

func handleDeleteHost(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig, hostUUID string) {
	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		result, err := db.ExecContext(r.Context(), `DELETE FROM hosts WHERE uuid = ?`, hostUUID)
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
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			shared.SendError(w, http.StatusNotFound, "host not found", nil, cfg)
			return
		}
		shared.SendError(w, http.StatusInternalServerError, "failed to delete host", err, cfg)
		return
	}

	shared.WriteJSON(w, http.StatusOK, map[string]any{"response": map[string]any{"isDeleted": true}})
}
