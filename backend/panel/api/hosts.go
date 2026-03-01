package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"v2ray-stat/backend/panel/config"
	dbmanager "v2ray-stat/backend/panel/db/manager"

	"github.com/google/uuid"
)

// Host represents a host entry from the hosts table.
type Host struct {
	UUID                     string  `json:"uuid"`
	ViewPosition             int     `json:"view_position"`
	Remark                   string  `json:"remark"`
	Address                  string  `json:"address"`
	Port                     int     `json:"port"`
	Path                     *string `json:"path,omitempty"`
	SNI                      *string `json:"sni,omitempty"`
	Host                     *string `json:"host,omitempty"`
	ALPN                     *string `json:"alpn,omitempty"`
	Fingerprint              *string `json:"fingerprint,omitempty"`
	SecurityLayer            string  `json:"security_layer"`
	XHTTPExtraParams         *string `json:"xhttp_extra_params,omitempty"`
	MuxParams                *string `json:"mux_params,omitempty"`
	SockoptParams            *string `json:"sockopt_params,omitempty"`
	IsDisabled               bool    `json:"is_disabled"`
	ServerDescription        *string `json:"server_description,omitempty"`
	VLESSRouteID             *int64  `json:"vless_route_id,omitempty"`
	AllowInsecure            bool    `json:"allow_insecure"`
	ShuffleHost              bool    `json:"shuffle_host"`
	MihomoX25519             bool    `json:"mihomo_x25519"`
	XrayJSONTemplateUUID     *string `json:"xray_json_template_uuid,omitempty"`
	KeepSNIBlank             bool    `json:"keep_sni_blank"`
	Tag                      *string `json:"tag,omitempty"`
	IsHidden                 bool    `json:"is_hidden"`
	OverrideSNIFromAddress   bool    `json:"override_sni_from_address"`
	ConfigProfileUUID        *string `json:"config_profile_uuid,omitempty"`
	ConfigProfileInboundUUID *string `json:"config_profile_inbound_uuid,omitempty"`
	ConfigProfileName        *string `json:"config_profile_name,omitempty"`
	ConfigProfileInboundTag  *string `json:"config_profile_inbound_tag,omitempty"`
}

// HostCreateRequest describes POST /api/v1/hosts payload.
type HostCreateRequest struct {
	ViewPosition             int     `json:"view_position"`
	Remark                   string  `json:"remark"`
	Address                  string  `json:"address"`
	Port                     int     `json:"port"`
	Path                     *string `json:"path,omitempty"`
	SNI                      *string `json:"sni,omitempty"`
	Host                     *string `json:"host,omitempty"`
	ALPN                     *string `json:"alpn,omitempty"`
	Fingerprint              *string `json:"fingerprint,omitempty"`
	SecurityLayer            string  `json:"security_layer"`
	XHTTPExtraParams         *string `json:"xhttp_extra_params,omitempty"`
	MuxParams                *string `json:"mux_params,omitempty"`
	SockoptParams            *string `json:"sockopt_params,omitempty"`
	IsDisabled               bool    `json:"is_disabled"`
	ServerDescription        *string `json:"server_description,omitempty"`
	VLESSRouteID             *int64  `json:"vless_route_id,omitempty"`
	AllowInsecure            bool    `json:"allow_insecure"`
	ShuffleHost              bool    `json:"shuffle_host"`
	MihomoX25519             bool    `json:"mihomo_x25519"`
	XrayJSONTemplateUUID     *string `json:"xray_json_template_uuid,omitempty"`
	KeepSNIBlank             bool    `json:"keep_sni_blank"`
	Tag                      *string `json:"tag,omitempty"`
	IsHidden                 bool    `json:"is_hidden"`
	OverrideSNIFromAddress   bool    `json:"override_sni_from_address"`
	ConfigProfileUUID        *string `json:"config_profile_uuid,omitempty"`
	ConfigProfileInboundUUID *string `json:"config_profile_inbound_uuid,omitempty"`
}

// HostUpdateRequest describes PATCH /api/v1/hosts/{uuid} payload.
type HostUpdateRequest struct {
	ViewPosition             *int    `json:"view_position,omitempty"`
	Remark                   *string `json:"remark,omitempty"`
	Address                  *string `json:"address,omitempty"`
	Port                     *int    `json:"port,omitempty"`
	Path                     *string `json:"path,omitempty"`
	SNI                      *string `json:"sni,omitempty"`
	Host                     *string `json:"host,omitempty"`
	ALPN                     *string `json:"alpn,omitempty"`
	Fingerprint              *string `json:"fingerprint,omitempty"`
	SecurityLayer            *string `json:"security_layer,omitempty"`
	XHTTPExtraParams         *string `json:"xhttp_extra_params,omitempty"`
	MuxParams                *string `json:"mux_params,omitempty"`
	SockoptParams            *string `json:"sockopt_params,omitempty"`
	IsDisabled               *bool   `json:"is_disabled,omitempty"`
	ServerDescription        *string `json:"server_description,omitempty"`
	VLESSRouteID             *int64  `json:"vless_route_id,omitempty"`
	AllowInsecure            *bool   `json:"allow_insecure,omitempty"`
	ShuffleHost              *bool   `json:"shuffle_host,omitempty"`
	MihomoX25519             *bool   `json:"mihomo_x25519,omitempty"`
	XrayJSONTemplateUUID     *string `json:"xray_json_template_uuid,omitempty"`
	KeepSNIBlank             *bool   `json:"keep_sni_blank,omitempty"`
	Tag                      *string `json:"tag,omitempty"`
	IsHidden                 *bool   `json:"is_hidden,omitempty"`
	OverrideSNIFromAddress   *bool   `json:"override_sni_from_address,omitempty"`
	ConfigProfileUUID        *string `json:"config_profile_uuid,omitempty"`
	ConfigProfileInboundUUID *string `json:"config_profile_inbound_uuid,omitempty"`
}

type HostNodeAssignment struct {
	HostUUID string `json:"host_uuid"`
	NodeUUID string `json:"node_uuid"`
}

type HostNodeAssignmentRequest struct {
	HostUUID  string   `json:"host_uuid"`
	NodeUUIDs []string `json:"node_uuids"`
}

type HostNodeAssignmentDeleteRequest struct {
	HostUUID  string   `json:"host_uuid"`
	NodeUUIDs []string `json:"node_uuids"`
}

func scanHost(scanner RowScanner) (Host, error) {
	var h Host
	var viewPosition sql.NullInt64
	var path, sni, host, alpn, fingerprint, securityLayer sql.NullString
	var xhttpExtraParams, muxParams, sockoptParams, serverDescription sql.NullString
	var vlessRouteID sql.NullInt64
	var xrayJSONTemplateUUID, tag, configProfileUUID, configProfileInboundUUID sql.NullString
	var configProfileName, configProfileInboundTag sql.NullString
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
		&tag,
		&isHidden,
		&overrideSNIFromAddress,
		&configProfileUUID,
		&configProfileInboundUUID,
		&configProfileName,
		&configProfileInboundTag,
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
		h.SecurityLayer = "none"
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
	if configProfileName.Valid {
		h.ConfigProfileName = &configProfileName.String
	}
	if configProfileInboundTag.Valid {
		h.ConfigProfileInboundTag = &configProfileInboundTag.String
	}

	return h, nil
}

func (r *HostCreateRequest) Validate() error {
	if strings.TrimSpace(r.Remark) == "" {
		return fmt.Errorf("remark is required")
	}
	if strings.TrimSpace(r.Address) == "" {
		return fmt.Errorf("address is required")
	}
	if r.Port < 1 || r.Port > 65535 {
		return fmt.Errorf("invalid port")
	}
	if strings.TrimSpace(r.SecurityLayer) == "" {
		r.SecurityLayer = "none"
	}
	switch strings.ToLower(strings.TrimSpace(r.SecurityLayer)) {
	case "none", "tls", "reality":
	default:
		return fmt.Errorf("security_layer must be one of: none, tls, reality")
	}

	for _, field := range []*string{r.ConfigProfileUUID, r.ConfigProfileInboundUUID, r.XrayJSONTemplateUUID} {
		if field != nil && *field != "" {
			if _, err := uuid.Parse(*field); err != nil {
				return fmt.Errorf("invalid UUID in reference field")
			}
		}
	}
	return nil
}

func (r *HostUpdateRequest) HasUpdates() bool {
	return r.ViewPosition != nil || r.Remark != nil || r.Address != nil || r.Port != nil ||
		r.Path != nil || r.SNI != nil || r.Host != nil || r.ALPN != nil ||
		r.Fingerprint != nil || r.SecurityLayer != nil || r.XHTTPExtraParams != nil ||
		r.MuxParams != nil || r.SockoptParams != nil || r.IsDisabled != nil ||
		r.ServerDescription != nil || r.VLESSRouteID != nil || r.AllowInsecure != nil ||
		r.ShuffleHost != nil || r.MihomoX25519 != nil || r.XrayJSONTemplateUUID != nil ||
		r.KeepSNIBlank != nil || r.Tag != nil || r.IsHidden != nil ||
		r.OverrideSNIFromAddress != nil || r.ConfigProfileUUID != nil || r.ConfigProfileInboundUUID != nil
}

func (r *HostUpdateRequest) Validate() error {
	if r.Remark != nil && strings.TrimSpace(*r.Remark) == "" {
		return fmt.Errorf("remark cannot be empty")
	}
	if r.Address != nil && strings.TrimSpace(*r.Address) == "" {
		return fmt.Errorf("address cannot be empty")
	}
	if r.Port != nil && (*r.Port < 1 || *r.Port > 65535) {
		return fmt.Errorf("invalid port")
	}
	if r.SecurityLayer != nil {
		switch strings.ToLower(strings.TrimSpace(*r.SecurityLayer)) {
		case "none", "tls", "reality":
		default:
			return fmt.Errorf("security_layer must be one of: none, tls, reality")
		}
	}
	for _, field := range []*string{r.ConfigProfileUUID, r.ConfigProfileInboundUUID, r.XrayJSONTemplateUUID} {
		if field != nil && *field != "" {
			if _, err := uuid.Parse(*field); err != nil {
				return fmt.Errorf("invalid UUID in reference field")
			}
		}
	}
	return nil
}

// HostsHandler handles GET/POST/DELETE /api/v1/hosts.
func HostsHandler(manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handleGetHosts(w, r, manager, cfg)
		case http.MethodPost:
			handleCreateHost(w, r, manager, cfg)
		case http.MethodDelete:
			handleDeleteAllHosts(w, r, manager, cfg)
		default:
			http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		}
	}
}

// HostsReorderHandler handles POST /api/v1/hosts/reorder
func HostsReorderHandler(manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
			return
		}

		var req ViewPositionReorderRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			sendError(w, http.StatusBadRequest, "invalid JSON", err, cfg)
			return
		}
		if err := req.Validate(); err != nil {
			sendError(w, http.StatusBadRequest, err.Error(), nil, cfg)
			return
		}

		err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
			return applyViewPositionReorder(r.Context(), db, "hosts", req.OrderedUUIDs, cfg)
		})
		if err != nil {
			sendError(w, http.StatusInternalServerError, "failed to reorder hosts", err, cfg)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "hosts reordered",
			"count":   len(req.OrderedUUIDs),
		})
	}
}

func handleGetHosts(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) {
	ctx := r.Context()
	var hosts []Host

	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		rows, err := db.QueryContext(ctx, `
			SELECT
				h.uuid, h.view_position, h.remark, h.address, h.port,
				h.path, h.sni, h.host, h.alpn, h.fingerprint, h.security_layer,
				h.xhttp_extra_params, h.mux_params, h.sockopt_params,
				h.is_disabled, h.server_description, h.vless_route_id,
				h.allow_insecure, h.shuffle_host, h.mihomo_x25519,
				h.xray_json_template_uuid, h.keep_sni_blank, h.tag, h.is_hidden,
				h.override_sni_from_address, h.config_profile_uuid, h.config_profile_inbound_uuid,
				cp.name, cpi.tag
			FROM hosts h
			LEFT JOIN config_profiles cp ON h.config_profile_uuid = cp.uuid
			LEFT JOIN config_profile_inbounds cpi ON h.config_profile_inbound_uuid = cpi.uuid
			ORDER BY h.view_position ASC, h.remark ASC
		`)
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			h, err := scanHost(rows)
			if err != nil {
				return err
			}
			hosts = append(hosts, h)
		}
		return rows.Err()
	})

	if err != nil {
		sendError(w, http.StatusInternalServerError, "failed to fetch hosts", err, cfg)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(map[string]interface{}{"hosts": hosts, "count": len(hosts)})
}

func handleCreateHost(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) {
	var req HostCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, http.StatusBadRequest, "invalid JSON", err, cfg)
		return
	}

	if err := req.Validate(); err != nil {
		sendError(w, http.StatusBadRequest, err.Error(), nil, cfg)
		return
	}

	hostUUID := uuid.New().String()
	ctx := r.Context()

	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		_, err := db.ExecContext(ctx, `
			INSERT INTO hosts (
				uuid, view_position, remark, address, port,
				path, sni, host, alpn, fingerprint, security_layer,
				xhttp_extra_params, mux_params, sockopt_params,
				is_disabled, server_description, vless_route_id,
				allow_insecure, shuffle_host, mihomo_x25519,
				xray_json_template_uuid, keep_sni_blank, tag, is_hidden,
				override_sni_from_address, config_profile_uuid, config_profile_inbound_uuid
				) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
			`,
			hostUUID, req.ViewPosition, req.Remark, req.Address, req.Port,
			normalizeNullableString(req.Path),
			normalizeNullableString(req.SNI),
			normalizeNullableString(req.Host),
			normalizeNullableString(req.ALPN),
			normalizeNullableString(req.Fingerprint),
			strings.ToLower(strings.TrimSpace(req.SecurityLayer)),
			normalizeNullableString(req.XHTTPExtraParams),
			normalizeNullableString(req.MuxParams),
			normalizeNullableString(req.SockoptParams),
			req.IsDisabled,
			normalizeNullableString(req.ServerDescription),
			req.VLESSRouteID,
			req.AllowInsecure,
			req.ShuffleHost,
			req.MihomoX25519,
			normalizeNullableString(req.XrayJSONTemplateUUID),
			req.KeepSNIBlank,
			normalizeNullableString(req.Tag),
			req.IsHidden,
			req.OverrideSNIFromAddress,
			normalizeNullableString(req.ConfigProfileUUID),
			normalizeNullableString(req.ConfigProfileInboundUUID),
		)
		return err
	})

	if err != nil {
		sendError(w, http.StatusInternalServerError, "failed to create host", err, cfg)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{"message": "host created", "uuid": hostUUID})
}

// HostByUUIDHandler handles GET/PATCH/DELETE /api/v1/hosts/{uuid}.
func HostByUUIDHandler(manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		hostUUID := strings.TrimSpace(strings.TrimPrefix(r.URL.Path, "/api/v1/hosts/"))
		if _, err := uuid.Parse(hostUUID); err != nil {
			sendError(w, http.StatusBadRequest, "invalid UUID format", nil, cfg)
			return
		}

		switch r.Method {
		case http.MethodGet:
			handleGetHost(w, r, manager, cfg, hostUUID)
		case http.MethodPatch:
			handlePatchHost(w, r, manager, cfg, hostUUID)
		case http.MethodDelete:
			handleDeleteHost(w, r, manager, cfg, hostUUID)
		default:
			http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		}
	}
}

func handleGetHost(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig, hostUUID string) {
	var host Host
	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		row := db.QueryRowContext(r.Context(), `
			SELECT
				h.uuid, h.view_position, h.remark, h.address, h.port,
				h.path, h.sni, h.host, h.alpn, h.fingerprint, h.security_layer,
				h.xhttp_extra_params, h.mux_params, h.sockopt_params,
				h.is_disabled, h.server_description, h.vless_route_id,
				h.allow_insecure, h.shuffle_host, h.mihomo_x25519,
				h.xray_json_template_uuid, h.keep_sni_blank, h.tag, h.is_hidden,
				h.override_sni_from_address, h.config_profile_uuid, h.config_profile_inbound_uuid,
				cp.name, cpi.tag
			FROM hosts h
			LEFT JOIN config_profiles cp ON h.config_profile_uuid = cp.uuid
			LEFT JOIN config_profile_inbounds cpi ON h.config_profile_inbound_uuid = cpi.uuid
			WHERE h.uuid = ?
		`, hostUUID)
		var err error
		host, err = scanHost(row)
		return err
	})
	if err != nil {
		if err == sql.ErrNoRows {
			sendError(w, http.StatusNotFound, "host not found", nil, cfg)
		} else {
			sendError(w, http.StatusInternalServerError, "failed to fetch host", err, cfg)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(map[string]interface{}{"host": host})
}

func handlePatchHost(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig, hostUUID string) {
	var req HostUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, http.StatusBadRequest, "invalid JSON", err, cfg)
		return
	}
	if !req.HasUpdates() {
		sendError(w, http.StatusBadRequest, "no fields to update", nil, cfg)
		return
	}
	if err := req.Validate(); err != nil {
		sendError(w, http.StatusBadRequest, err.Error(), nil, cfg)
		return
	}

	var clauses []string
	var args []interface{}
	add := func(col string, value interface{}) {
		clauses = append(clauses, fmt.Sprintf("%s = ?", col))
		args = append(args, value)
	}
	addNullable := func(col string, value *string) {
		if value == nil {
			return
		}
		if *value == "" {
			clauses = append(clauses, fmt.Sprintf("%s = NULL", col))
			return
		}
		add(col, *value)
	}

	if req.ViewPosition != nil {
		add("view_position", *req.ViewPosition)
	}
	if req.Remark != nil {
		add("remark", strings.TrimSpace(*req.Remark))
	}
	if req.Address != nil {
		add("address", strings.TrimSpace(*req.Address))
	}
	if req.Port != nil {
		add("port", *req.Port)
	}
	addNullable("path", req.Path)
	addNullable("sni", req.SNI)
	addNullable("host", req.Host)
	addNullable("alpn", req.ALPN)
	addNullable("fingerprint", req.Fingerprint)
	if req.SecurityLayer != nil {
		add("security_layer", strings.ToLower(strings.TrimSpace(*req.SecurityLayer)))
	}
	addNullable("xhttp_extra_params", req.XHTTPExtraParams)
	addNullable("mux_params", req.MuxParams)
	addNullable("sockopt_params", req.SockoptParams)
	if req.IsDisabled != nil {
		add("is_disabled", *req.IsDisabled)
	}
	addNullable("server_description", req.ServerDescription)
	if req.VLESSRouteID != nil {
		add("vless_route_id", *req.VLESSRouteID)
	}
	if req.AllowInsecure != nil {
		add("allow_insecure", *req.AllowInsecure)
	}
	if req.ShuffleHost != nil {
		add("shuffle_host", *req.ShuffleHost)
	}
	if req.MihomoX25519 != nil {
		add("mihomo_x25519", *req.MihomoX25519)
	}
	addNullable("xray_json_template_uuid", req.XrayJSONTemplateUUID)
	if req.KeepSNIBlank != nil {
		add("keep_sni_blank", *req.KeepSNIBlank)
	}
	addNullable("tag", req.Tag)
	if req.IsHidden != nil {
		add("is_hidden", *req.IsHidden)
	}
	if req.OverrideSNIFromAddress != nil {
		add("override_sni_from_address", *req.OverrideSNIFromAddress)
	}
	addNullable("config_profile_uuid", req.ConfigProfileUUID)
	addNullable("config_profile_inbound_uuid", req.ConfigProfileInboundUUID)

	if len(clauses) == 0 {
		sendError(w, http.StatusBadRequest, "no fields to update", nil, cfg)
		return
	}

	args = append(args, hostUUID)
	query := fmt.Sprintf("UPDATE hosts SET %s WHERE uuid = ?", strings.Join(clauses, ", "))

	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		result, err := db.ExecContext(r.Context(), query, args...)
		if err != nil {
			return err
		}
		rowsAffected, err := result.RowsAffected()
		if err != nil {
			return err
		}
		if rowsAffected == 0 {
			return sql.ErrNoRows
		}
		return nil
	})
	if err != nil {
		if err == sql.ErrNoRows {
			sendError(w, http.StatusNotFound, "host not found", nil, cfg)
		} else {
			sendError(w, http.StatusInternalServerError, "update failed", err, cfg)
		}
		return
	}

	handleGetHost(w, r, manager, cfg, hostUUID)
}

func handleDeleteHost(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig, hostUUID string) {
	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		result, err := db.ExecContext(r.Context(), "DELETE FROM hosts WHERE uuid = ?", hostUUID)
		if err != nil {
			return err
		}
		rowsAffected, err := result.RowsAffected()
		if err != nil {
			return err
		}
		if rowsAffected == 0 {
			return sql.ErrNoRows
		}
		return nil
	})
	if err != nil {
		if err == sql.ErrNoRows {
			sendError(w, http.StatusNotFound, "host not found", nil, cfg)
		} else {
			sendError(w, http.StatusInternalServerError, "failed to delete host", err, cfg)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(map[string]interface{}{"message": "host deleted", "uuid": hostUUID})
}

func handleDeleteAllHosts(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) {
	if rawUUIDs := r.URL.Query().Get("uuids"); rawUUIDs != "" {
		uuids, err := parseUUIDCSV(rawUUIDs)
		if err != nil {
			sendError(w, http.StatusBadRequest, err.Error(), nil, cfg)
			return
		}

		var deleted int64
		err = manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
			tx, err := db.BeginTx(r.Context(), nil)
			if err != nil {
				return err
			}
			for _, id := range uuids {
				res, execErr := tx.ExecContext(r.Context(), "DELETE FROM hosts WHERE uuid = ?", id)
				if execErr != nil {
					_ = tx.Rollback()
					return execErr
				}
				rows, rowsErr := res.RowsAffected()
				if rowsErr != nil {
					_ = tx.Rollback()
					return rowsErr
				}
				deleted += rows
			}
			return tx.Commit()
		})
		if err != nil {
			sendError(w, http.StatusInternalServerError, "failed to delete selected hosts", err, cfg)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(w).Encode(map[string]interface{}{"message": "selected hosts deleted", "count": deleted})
		return
	}

	if r.URL.Query().Get("confirm") != "true" {
		sendError(w, http.StatusBadRequest, "confirmation required. use DELETE /api/v1/hosts?confirm=true", nil, cfg)
		return
	}

	var deleted int64
	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		result, err := db.ExecContext(r.Context(), "DELETE FROM hosts")
		if err != nil {
			return err
		}
		deleted, err = result.RowsAffected()
		return err
	})
	if err != nil {
		sendError(w, http.StatusInternalServerError, "failed to delete hosts", err, cfg)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(map[string]interface{}{"message": "all hosts deleted", "count": deleted})
}

// HostNodeAssignmentsHandler handles GET/POST/DELETE /api/v1/hosts-to-nodes.
func HostNodeAssignmentsHandler(manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handleGetHostNodeAssignments(w, r, manager, cfg)
		case http.MethodPost:
			handleSetHostNodeAssignments(w, r, manager, cfg)
		case http.MethodDelete:
			handleDeleteHostNodeAssignments(w, r, manager, cfg)
		default:
			http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		}
	}
}

func handleGetHostNodeAssignments(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) {
	hostUUID := strings.TrimSpace(r.URL.Query().Get("host_uuid"))
	nodeUUID := strings.TrimSpace(r.URL.Query().Get("node_uuid"))
	if hostUUID != "" {
		if _, err := uuid.Parse(hostUUID); err != nil {
			sendError(w, http.StatusBadRequest, "invalid host_uuid", nil, cfg)
			return
		}
	}
	if nodeUUID != "" {
		if _, err := uuid.Parse(nodeUUID); err != nil {
			sendError(w, http.StatusBadRequest, "invalid node_uuid", nil, cfg)
			return
		}
	}

	assignments := make([]HostNodeAssignment, 0)
	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		query := "SELECT host_uuid, node_uuid FROM hosts_to_nodes"
		args := []interface{}{}
		conditions := make([]string, 0, 2)
		if hostUUID != "" {
			conditions = append(conditions, "host_uuid = ?")
			args = append(args, hostUUID)
		}
		if nodeUUID != "" {
			conditions = append(conditions, "node_uuid = ?")
			args = append(args, nodeUUID)
		}
		if len(conditions) > 0 {
			query += " WHERE " + strings.Join(conditions, " AND ")
		}
		query += " ORDER BY host_uuid, node_uuid"

		rows, err := db.QueryContext(r.Context(), query, args...)
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			var a HostNodeAssignment
			if err := rows.Scan(&a.HostUUID, &a.NodeUUID); err != nil {
				return err
			}
			assignments = append(assignments, a)
		}
		return rows.Err()
	})
	if err != nil {
		sendError(w, http.StatusInternalServerError, "failed to fetch host-node assignments", err, cfg)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(map[string]interface{}{"assignments": assignments, "count": len(assignments)})
}

func handleSetHostNodeAssignments(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) {
	var req HostNodeAssignmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, http.StatusBadRequest, "invalid JSON", err, cfg)
		return
	}
	req.NodeUUIDs = dedupeStrings(req.NodeUUIDs)

	if _, err := uuid.Parse(req.HostUUID); err != nil {
		sendError(w, http.StatusBadRequest, "invalid host_uuid", nil, cfg)
		return
	}
	for _, nodeUUID := range req.NodeUUIDs {
		if _, err := uuid.Parse(nodeUUID); err != nil {
			sendError(w, http.StatusBadRequest, "invalid node UUID in node_uuids", nil, cfg)
			return
		}
	}

	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		tx, err := db.BeginTx(r.Context(), nil)
		if err != nil {
			return err
		}

		if _, err := tx.ExecContext(r.Context(), "DELETE FROM hosts_to_nodes WHERE host_uuid = ?", req.HostUUID); err != nil {
			_ = tx.Rollback()
			return err
		}

		for _, nodeUUID := range req.NodeUUIDs {
			if _, err := tx.ExecContext(r.Context(), "INSERT INTO hosts_to_nodes (host_uuid, node_uuid) VALUES (?, ?)", req.HostUUID, nodeUUID); err != nil {
				_ = tx.Rollback()
				return err
			}
		}

		return tx.Commit()
	})
	if err != nil {
		sendError(w, http.StatusInternalServerError, "failed to set host-node assignments", err, cfg)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":     "host-node assignments updated",
		"host_uuid":   req.HostUUID,
		"nodes_count": len(req.NodeUUIDs),
	})
}

func handleDeleteHostNodeAssignments(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) {
	var req HostNodeAssignmentDeleteRequest
	if r.ContentLength > 0 {
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			sendError(w, http.StatusBadRequest, "invalid JSON", err, cfg)
			return
		}
	}

	if req.HostUUID == "" {
		req.HostUUID = strings.TrimSpace(r.URL.Query().Get("host_uuid"))
	}
	if len(req.NodeUUIDs) == 0 {
		nodeUUIDs, err := parseOptionalUUIDCSV(r.URL.Query().Get("node_uuids"))
		if err != nil {
			sendError(w, http.StatusBadRequest, "invalid node_uuids query parameter", err, cfg)
			return
		}
		req.NodeUUIDs = nodeUUIDs
	}
	req.NodeUUIDs = dedupeStrings(req.NodeUUIDs)

	if _, err := uuid.Parse(req.HostUUID); err != nil {
		sendError(w, http.StatusBadRequest, "invalid host_uuid", nil, cfg)
		return
	}
	for _, nodeUUID := range req.NodeUUIDs {
		if _, err := uuid.Parse(nodeUUID); err != nil {
			sendError(w, http.StatusBadRequest, "invalid node UUID in node_uuids", nil, cfg)
			return
		}
	}

	var deletedCount int64
	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		if len(req.NodeUUIDs) == 0 {
			result, err := db.ExecContext(r.Context(), "DELETE FROM hosts_to_nodes WHERE host_uuid = ?", req.HostUUID)
			if err != nil {
				return err
			}
			rowsAffected, err := result.RowsAffected()
			if err != nil {
				return err
			}
			deletedCount = rowsAffected
			return nil
		}

		tx, err := db.BeginTx(r.Context(), nil)
		if err != nil {
			return err
		}

		for _, nodeUUID := range req.NodeUUIDs {
			result, execErr := tx.ExecContext(
				r.Context(),
				"DELETE FROM hosts_to_nodes WHERE host_uuid = ? AND node_uuid = ?",
				req.HostUUID,
				nodeUUID,
			)
			if execErr != nil {
				_ = tx.Rollback()
				return execErr
			}
			rowsAffected, rowsErr := result.RowsAffected()
			if rowsErr != nil {
				_ = tx.Rollback()
				return rowsErr
			}
			deletedCount += rowsAffected
		}

		return tx.Commit()
	})
	if err != nil {
		sendError(w, http.StatusInternalServerError, "failed to delete host-node assignments", err, cfg)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":       "host-node assignments deleted",
		"host_uuid":     req.HostUUID,
		"deleted_count": deletedCount,
		"all_nodes":     len(req.NodeUUIDs) == 0,
	})
}

func parseOptionalUUIDCSV(raw string) ([]string, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil, nil
	}
	return parseUUIDCSV(raw)
}

func dedupeStrings(values []string) []string {
	if len(values) < 2 {
		return values
	}

	seen := make(map[string]struct{}, len(values))
	out := make([]string, 0, len(values))
	for _, value := range values {
		if _, exists := seen[value]; exists {
			continue
		}
		seen[value] = struct{}{}
		out = append(out, value)
	}

	return out
}

func normalizeNullableString(s *string) interface{} {
	if s == nil || *s == "" {
		return nil
	}
	return *s
}
