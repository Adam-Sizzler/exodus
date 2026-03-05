package hosts

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"sort"
	"strings"

	"v2ray-stat/backend/panel/config"
	dbmanager "v2ray-stat/backend/panel/db/manager"
	"v2ray-stat/backend/panel/dbutil"
	"v2ray-stat/backend/panel/httpapi/shared"

	"github.com/google/uuid"
)

var (
	hostTagRegex          = regexp.MustCompile(`^[A-Z0-9_:]+$`)
	allowedSecurityLayers = map[string]struct{}{"DEFAULT": {}, "TLS": {}, "NONE": {}}
	allowedAlpn           = map[string]struct{}{
		"h3":             {},
		"h2":             {},
		"http/1.1":       {},
		"h2,http/1.1":    {},
		"h3,h2,http/1.1": {},
		"h3,h2":          {},
	}
	allowedFingerprints = map[string]struct{}{
		"chrome":     {},
		"firefox":    {},
		"safari":     {},
		"ios":        {},
		"android":    {},
		"edge":       {},
		"qq":         {},
		"random":     {},
		"randomized": {},
	}
	allowedTemplateTypes = map[string]struct{}{
		"XRAY_JSON":   {},
		"XRAY_BASE64": {},
		"MIHOMO":      {},
		"STASH":       {},
		"CLASH":       {},
		"SINGBOX":     {},
	}
)

var (
	errConfigProfileNotFound        = errors.New("config profile not found")
	errConfigProfileInboundNotFound = errors.New("config profile inbound not found")
	errTemplateNotFound             = errors.New("subscription template not found")
	errTemplateTypeNotAllowed       = errors.New("template type not allowed")
)

type hostRecord struct {
	UUID                     string
	ViewPosition             int
	Remark                   string
	Address                  string
	Port                     int
	Path                     *string
	SNI                      *string
	Host                     *string
	ALPN                     *string
	Fingerprint              *string
	SecurityLayer            string
	XHTTPExtraParams         json.RawMessage
	MuxParams                json.RawMessage
	SockoptParams            json.RawMessage
	IsDisabled               bool
	ServerDescription        *string
	VLESSRouteID             *int64
	AllowInsecure            bool
	ShuffleHost              bool
	MihomoX25519             bool
	XrayJSONTemplateUUID     *string
	KeepSNIBlank             bool
	Tag                      *string
	IsHidden                 bool
	OverrideSNIFromAddress   bool
	ConfigProfileUUID        *string
	ConfigProfileInboundUUID *string
	ExcludeTypes             []string
}

type HostInbound struct {
	ConfigProfileUUID        *string `json:"configProfileUuid"`
	ConfigProfileInboundUUID *string `json:"configProfileInboundUuid"`
}

type OptionalString struct {
	Set   bool
	Value *string
}

func (o *OptionalString) UnmarshalJSON(data []byte) error {
	o.Set = true
	if string(data) == "null" {
		o.Value = nil
		return nil
	}
	var value string
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	o.Value = &value
	return nil
}

type OptionalInt64 struct {
	Set   bool
	Value *int64
}

func (o *OptionalInt64) UnmarshalJSON(data []byte) error {
	o.Set = true
	if string(data) == "null" {
		o.Value = nil
		return nil
	}
	var value int64
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	o.Value = &value
	return nil
}

type OptionalJSON struct {
	Set bool
	Raw json.RawMessage
}

func (o *OptionalJSON) UnmarshalJSON(data []byte) error {
	o.Set = true
	o.Raw = append(o.Raw[:0], data...)
	return nil
}

type HostAPI struct {
	UUID                    string      `json:"uuid"`
	ViewPosition            int         `json:"viewPosition"`
	Remark                  string      `json:"remark"`
	Address                 string      `json:"address"`
	Port                    int         `json:"port"`
	Path                    *string     `json:"path"`
	SNI                     *string     `json:"sni"`
	Host                    *string     `json:"host"`
	ALPN                    *string     `json:"alpn"`
	Fingerprint             *string     `json:"fingerprint"`
	IsDisabled              bool        `json:"isDisabled"`
	SecurityLayer           string      `json:"securityLayer"`
	XHTTPExtraParams        interface{} `json:"xHttpExtraParams"`
	MuxParams               interface{} `json:"muxParams"`
	SockoptParams           interface{} `json:"sockoptParams"`
	Inbound                 HostInbound `json:"inbound"`
	ServerDescription       *string     `json:"serverDescription"`
	Tag                     *string     `json:"tag"`
	IsHidden                bool        `json:"isHidden"`
	OverrideSNIFromAddress  bool        `json:"overrideSniFromAddress"`
	KeepSNIBlank            bool        `json:"keepSniBlank"`
	VLESSRouteID            *int64      `json:"vlessRouteId"`
	AllowInsecure           bool        `json:"allowInsecure"`
	ShuffleHost             bool        `json:"shuffleHost"`
	MihomoX25519            bool        `json:"mihomoX25519"`
	Nodes                   []string    `json:"nodes"`
	XrayJSONTemplateUUID    *string     `json:"xrayJsonTemplateUuid"`
	ExcludedInternalSquads  []string    `json:"excludedInternalSquads"`
	ExcludeFromSubscription []string    `json:"excludeFromSubscriptionTypes"`
}

type HostCreateRequestAPI struct {
	Inbound                 HostInbound      `json:"inbound"`
	Remark                  string           `json:"remark"`
	Address                 string           `json:"address"`
	Port                    int              `json:"port"`
	Path                    *string          `json:"path,omitempty"`
	SNI                     *string          `json:"sni,omitempty"`
	Host                    *string          `json:"host,omitempty"`
	ALPN                    *string          `json:"alpn,omitempty"`
	Fingerprint             *string          `json:"fingerprint,omitempty"`
	IsDisabled              *bool            `json:"isDisabled,omitempty"`
	SecurityLayer           *string          `json:"securityLayer,omitempty"`
	XHTTPExtraParams        *json.RawMessage `json:"xHttpExtraParams,omitempty"`
	MuxParams               *json.RawMessage `json:"muxParams,omitempty"`
	SockoptParams           *json.RawMessage `json:"sockoptParams,omitempty"`
	ServerDescription       *string          `json:"serverDescription,omitempty"`
	Tag                     *string          `json:"tag,omitempty"`
	IsHidden                *bool            `json:"isHidden,omitempty"`
	OverrideSNIFromAddress  *bool            `json:"overrideSniFromAddress,omitempty"`
	KeepSNIBlank            *bool            `json:"keepSniBlank,omitempty"`
	AllowInsecure           *bool            `json:"allowInsecure,omitempty"`
	VLESSRouteID            *int64           `json:"vlessRouteId,omitempty"`
	ShuffleHost             *bool            `json:"shuffleHost,omitempty"`
	MihomoX25519            *bool            `json:"mihomoX25519,omitempty"`
	Nodes                   []string         `json:"nodes,omitempty"`
	XrayJSONTemplateUUID    *string          `json:"xrayJsonTemplateUuid,omitempty"`
	ExcludedInternalSquads  []string         `json:"excludedInternalSquads,omitempty"`
	ExcludeFromSubscription []string         `json:"excludeFromSubscriptionTypes,omitempty"`
}

type HostUpdateRequestAPI struct {
	UUID                    string         `json:"uuid"`
	Inbound                 *HostInbound   `json:"inbound,omitempty"`
	Remark                  OptionalString `json:"remark,omitempty"`
	Address                 OptionalString `json:"address,omitempty"`
	Port                    *int           `json:"port,omitempty"`
	Path                    OptionalString `json:"path,omitempty"`
	SNI                     OptionalString `json:"sni,omitempty"`
	Host                    OptionalString `json:"host,omitempty"`
	ALPN                    OptionalString `json:"alpn,omitempty"`
	Fingerprint             OptionalString `json:"fingerprint,omitempty"`
	IsDisabled              *bool          `json:"isDisabled,omitempty"`
	SecurityLayer           *string        `json:"securityLayer,omitempty"`
	XHTTPExtraParams        OptionalJSON   `json:"xHttpExtraParams,omitempty"`
	MuxParams               OptionalJSON   `json:"muxParams,omitempty"`
	SockoptParams           OptionalJSON   `json:"sockoptParams,omitempty"`
	ServerDescription       OptionalString `json:"serverDescription,omitempty"`
	Tag                     OptionalString `json:"tag,omitempty"`
	IsHidden                *bool          `json:"isHidden,omitempty"`
	OverrideSNIFromAddress  *bool          `json:"overrideSniFromAddress,omitempty"`
	KeepSNIBlank            *bool          `json:"keepSniBlank,omitempty"`
	AllowInsecure           *bool          `json:"allowInsecure,omitempty"`
	VLESSRouteID            OptionalInt64  `json:"vlessRouteId,omitempty"`
	ShuffleHost             *bool          `json:"shuffleHost,omitempty"`
	MihomoX25519            *bool          `json:"mihomoX25519,omitempty"`
	Nodes                   []string       `json:"nodes,omitempty"`
	XrayJSONTemplateUUID    OptionalString `json:"xrayJsonTemplateUuid,omitempty"`
	ExcludedInternalSquads  []string       `json:"excludedInternalSquads,omitempty"`
	ExcludeFromSubscription []string       `json:"excludeFromSubscriptionTypes,omitempty"`
}

type reorderHostsRequest struct {
	Hosts []struct {
		UUID         string `json:"uuid"`
		ViewPosition int    `json:"viewPosition"`
	} `json:"hosts"`
}

type bulkUUIDsRequest struct {
	UUIDs []string `json:"uuids"`
}

type setInboundRequest struct {
	UUIDs                    []string `json:"uuids"`
	ConfigProfileUUID        string   `json:"configProfileUuid"`
	ConfigProfileInboundUUID string   `json:"configProfileInboundUuid"`
}

type setPortRequest struct {
	UUIDs []string `json:"uuids"`
	Port  int      `json:"port"`
}

func HostsHandler(manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handleGetHosts(w, r, manager, cfg)
		case http.MethodPost:
			handleCreateHost(w, r, manager, cfg)
		case http.MethodPatch:
			handleUpdateHost(w, r, manager, cfg)
		default:
			shared.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
		}
	}
}

func HostByUUIDHandler(manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uuidStr := strings.TrimSpace(strings.TrimPrefix(r.URL.Path, "/api/hosts/"))
		if _, err := uuid.Parse(uuidStr); err != nil {
			shared.SendError(w, http.StatusBadRequest, "invalid UUID format", nil, cfg)
			return
		}

		switch r.Method {
		case http.MethodGet:
			handleGetHost(w, r, manager, cfg, uuidStr)
		case http.MethodDelete:
			handleDeleteHost(w, r, manager, cfg, uuidStr)
		default:
			shared.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
		}
	}
}

func HostsActionsHandler(manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			shared.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}
		path := strings.TrimPrefix(r.URL.Path, "/api/hosts/actions/")
		path = strings.Trim(path, "/")
		switch path {
		case "reorder":
			handleReorderHosts(w, r, manager, cfg)
		default:
			http.NotFound(w, r)
		}
	}
}

func HostsBulkHandler(manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			shared.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}

		path := strings.TrimPrefix(r.URL.Path, "/api/hosts/bulk/")
		path = strings.Trim(path, "/")
		switch path {
		case "enable":
			handleBulkEnableHosts(w, r, manager, cfg)
		case "disable":
			handleBulkDisableHosts(w, r, manager, cfg)
		case "delete":
			handleBulkDeleteHosts(w, r, manager, cfg)
		case "set-inbound":
			handleBulkSetInbound(w, r, manager, cfg)
		case "set-port":
			handleBulkSetPort(w, r, manager, cfg)
		default:
			http.NotFound(w, r)
		}
	}
}

func HostsTagsHandler(manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			shared.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}
		tags, err := getHostTags(r.Context(), manager)
		if err != nil {
			shared.SendError(w, http.StatusInternalServerError, "failed to fetch host tags", err, cfg)
			return
		}
		shared.WriteJSON(w, http.StatusOK, map[string]any{"response": map[string]any{"tags": tags}})
	}
}

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

		xhttp, err := normalizeJSONValue(req.XHTTPExtraParams, false)
		if err != nil {
			_ = tx.Rollback()
			return err
		}
		mux, err := normalizeJSONValue(req.MuxParams, true)
		if err != nil {
			_ = tx.Rollback()
			return err
		}
		sockopt, err := normalizeJSONValue(req.SockoptParams, true)
		if err != nil {
			_ = tx.Rollback()
			return err
		}

		_, err = tx.ExecContext(r.Context(), `
            INSERT INTO hosts (
                uuid, remark, address, port,
                path, sni, host, alpn, fingerprint, security_layer,
                xhttp_extra_params, mux_params, sockopt_params,
                is_disabled, server_description, vless_route_id,
                allow_insecure, shuffle_host, mihomo_x25519,
                xray_json_template_uuid, keep_sni_blank,
                exclude_from_subscription_types, tag, is_hidden,
                override_sni_from_address, config_profile_uuid, config_profile_inbound_uuid
            ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
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
			xhttp,
			mux,
			sockopt,
			coalesceBool(req.IsDisabled, false),
			normalizeOptionalStringAllowEmpty(req.ServerDescription),
			req.VLESSRouteID,
			coalesceBool(req.AllowInsecure, false),
			coalesceBool(req.ShuffleHost, false),
			coalesceBool(req.MihomoX25519, false),
			normalizeOptionalStringAllowEmpty(req.XrayJSONTemplateUUID),
			coalesceBool(req.KeepSNIBlank, false),
			ensureStringSlice(req.ExcludeFromSubscription),
			normalizeOptionalStringAllowEmpty(req.Tag),
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
	if err := validateUpdateRequest(req); err != nil {
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

		if req.Remark.Set {
			if req.Remark.Value == nil {
				_ = tx.Rollback()
				return fmt.Errorf("remark cannot be null")
			}
			add("remark", strings.TrimSpace(*req.Remark.Value))
		}
		if req.Address.Set {
			if req.Address.Value == nil {
				_ = tx.Rollback()
				return fmt.Errorf("address cannot be null")
			}
			add("address", strings.TrimSpace(*req.Address.Value))
		}
		if req.Port != nil {
			add("port", *req.Port)
		}
		if req.Path.Set {
			if req.Path.Value == nil {
				_ = tx.Rollback()
				return fmt.Errorf("path cannot be null")
			}
			add("path", strings.TrimSpace(*req.Path.Value))
		}
		if req.SNI.Set {
			if req.SNI.Value == nil {
				_ = tx.Rollback()
				return fmt.Errorf("sni cannot be null")
			}
			add("sni", strings.TrimSpace(*req.SNI.Value))
		}
		if req.Host.Set {
			if req.Host.Value == nil {
				_ = tx.Rollback()
				return fmt.Errorf("host cannot be null")
			}
			add("host", strings.TrimSpace(*req.Host.Value))
		}
		if req.ALPN.Set {
			if req.ALPN.Value == nil {
				clauses = append(clauses, "alpn = NULL")
			} else {
				add("alpn", strings.TrimSpace(*req.ALPN.Value))
			}
		}
		if req.Fingerprint.Set {
			if req.Fingerprint.Value == nil {
				clauses = append(clauses, "fingerprint = NULL")
			} else {
				add("fingerprint", strings.TrimSpace(*req.Fingerprint.Value))
			}
		}
		if req.SecurityLayer != nil {
			add("security_layer", normalizeSecurityLayer(req.SecurityLayer))
		}

		if set, val, err := normalizeOptionalJSONField(req.XHTTPExtraParams, false); err != nil {
			_ = tx.Rollback()
			return err
		} else if set {
			if val == nil {
				clauses = append(clauses, "xhttp_extra_params = NULL")
			} else {
				add("xhttp_extra_params", val)
			}
		}
		if set, val, err := normalizeOptionalJSONField(req.MuxParams, true); err != nil {
			_ = tx.Rollback()
			return err
		} else if set {
			if val == nil {
				clauses = append(clauses, "mux_params = NULL")
			} else {
				add("mux_params", val)
			}
		}
		if set, val, err := normalizeOptionalJSONField(req.SockoptParams, true); err != nil {
			_ = tx.Rollback()
			return err
		} else if set {
			if val == nil {
				clauses = append(clauses, "sockopt_params = NULL")
			} else {
				add("sockopt_params", val)
			}
		}

		if req.IsDisabled != nil {
			add("is_disabled", *req.IsDisabled)
		}
		if req.ServerDescription.Set {
			if req.ServerDescription.Value == nil {
				clauses = append(clauses, "server_description = NULL")
			} else {
				add("server_description", strings.TrimSpace(*req.ServerDescription.Value))
			}
		}
		if req.VLESSRouteID.Set {
			if req.VLESSRouteID.Value == nil {
				clauses = append(clauses, "vless_route_id = NULL")
			} else {
				add("vless_route_id", *req.VLESSRouteID.Value)
			}
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
		if req.XrayJSONTemplateUUID.Set {
			if req.XrayJSONTemplateUUID.Value == nil {
				clauses = append(clauses, "xray_json_template_uuid = NULL")
			} else {
				add("xray_json_template_uuid", strings.TrimSpace(*req.XrayJSONTemplateUUID.Value))
			}
		}
		if req.KeepSNIBlank != nil {
			add("keep_sni_blank", *req.KeepSNIBlank)
		}
		if req.Tag.Set {
			if req.Tag.Value == nil {
				clauses = append(clauses, "tag = NULL")
			} else {
				add("tag", strings.TrimSpace(*req.Tag.Value))
			}
		}
		if req.IsHidden != nil {
			add("is_hidden", *req.IsHidden)
		}
		if req.OverrideSNIFromAddress != nil {
			add("override_sni_from_address", *req.OverrideSNIFromAddress)
		}
		if req.Inbound != nil {
			addOptionalString("config_profile_uuid", req.Inbound.ConfigProfileUUID)
			addOptionalString("config_profile_inbound_uuid", req.Inbound.ConfigProfileInboundUUID)
		}
		if req.ExcludeFromSubscription != nil {
			add("exclude_from_subscription_types", req.ExcludeFromSubscription)
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

func handleReorderHosts(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) {
	var req reorderHostsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		shared.SendError(w, http.StatusBadRequest, "invalid JSON", err, cfg)
		return
	}
	if len(req.Hosts) == 0 {
		shared.SendError(w, http.StatusBadRequest, "hosts cannot be empty", nil, cfg)
		return
	}
	for _, item := range req.Hosts {
		if _, err := uuid.Parse(item.UUID); err != nil {
			shared.SendError(w, http.StatusBadRequest, "invalid UUID format", nil, cfg)
			return
		}
	}

	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		tx, err := db.BeginTx(r.Context(), nil)
		if err != nil {
			return err
		}
		for _, item := range req.Hosts {
			if _, err := tx.ExecContext(r.Context(), `UPDATE hosts SET view_position = ? WHERE uuid = ?`, item.ViewPosition, item.UUID); err != nil {
				_ = tx.Rollback()
				return err
			}
		}
		if _, err := tx.ExecContext(r.Context(), `SELECT setval('hosts_view_position_seq', (SELECT COALESCE(MAX(view_position), 0) FROM hosts) + 1)`); err != nil {
			_ = tx.Rollback()
			return err
		}
		return tx.Commit()
	})
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to reorder hosts", err, cfg)
		return
	}

	shared.WriteJSON(w, http.StatusOK, map[string]any{"response": map[string]any{"isUpdated": true}})
}

func handleBulkEnableHosts(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) {
	bulkUpdateHostsEnabled(w, r, manager, cfg, true)
}

func handleBulkDisableHosts(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) {
	bulkUpdateHostsEnabled(w, r, manager, cfg, false)
}

func handleBulkDeleteHosts(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) {
	var req bulkUUIDsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		shared.SendError(w, http.StatusBadRequest, "invalid JSON", err, cfg)
		return
	}
	if len(req.UUIDs) == 0 {
		shared.SendError(w, http.StatusBadRequest, "uuids cannot be empty", nil, cfg)
		return
	}

	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		_, err := db.ExecContext(r.Context(), `DELETE FROM hosts WHERE uuid = ANY(?)`, req.UUIDs)
		return err
	})
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to delete hosts", err, cfg)
		return
	}

	handleGetHosts(w, r, manager, cfg)
}

func handleBulkSetInbound(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) {
	var req setInboundRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		shared.SendError(w, http.StatusBadRequest, "invalid JSON", err, cfg)
		return
	}
	if len(req.UUIDs) == 0 {
		shared.SendError(w, http.StatusBadRequest, "uuids cannot be empty", nil, cfg)
		return
	}
	if _, err := uuid.Parse(req.ConfigProfileUUID); err != nil {
		shared.SendError(w, http.StatusBadRequest, "invalid configProfileUuid", nil, cfg)
		return
	}
	if _, err := uuid.Parse(req.ConfigProfileInboundUUID); err != nil {
		shared.SendError(w, http.StatusBadRequest, "invalid configProfileInboundUuid", nil, cfg)
		return
	}
	if err := ensureConfigProfileInbound(r.Context(), manager, req.ConfigProfileUUID, req.ConfigProfileInboundUUID); err != nil {
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

	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		_, err := db.ExecContext(r.Context(), `
            UPDATE hosts
            SET config_profile_uuid = ?, config_profile_inbound_uuid = ?
            WHERE uuid = ANY(?)
        `, req.ConfigProfileUUID, req.ConfigProfileInboundUUID, req.UUIDs)
		return err
	})
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to set inbound", err, cfg)
		return
	}

	handleGetHosts(w, r, manager, cfg)
}

func handleBulkSetPort(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) {
	var req setPortRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		shared.SendError(w, http.StatusBadRequest, "invalid JSON", err, cfg)
		return
	}
	if len(req.UUIDs) == 0 {
		shared.SendError(w, http.StatusBadRequest, "uuids cannot be empty", nil, cfg)
		return
	}
	if req.Port < 1 || req.Port > 65535 {
		shared.SendError(w, http.StatusBadRequest, "invalid port", nil, cfg)
		return
	}

	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		_, err := db.ExecContext(r.Context(), `UPDATE hosts SET port = ? WHERE uuid = ANY(?)`, req.Port, req.UUIDs)
		return err
	})
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to set port", err, cfg)
		return
	}

	handleGetHosts(w, r, manager, cfg)
}

func bulkUpdateHostsEnabled(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig, enabled bool) {
	var req bulkUUIDsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		shared.SendError(w, http.StatusBadRequest, "invalid JSON", err, cfg)
		return
	}
	if len(req.UUIDs) == 0 {
		shared.SendError(w, http.StatusBadRequest, "uuids cannot be empty", nil, cfg)
		return
	}

	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		_, err := db.ExecContext(r.Context(), `UPDATE hosts SET is_disabled = ? WHERE uuid = ANY(?)`, !enabled, req.UUIDs)
		return err
	})
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to update hosts", err, cfg)
		return
	}

	handleGetHosts(w, r, manager, cfg)
}

func validateCreateRequest(req HostCreateRequestAPI) error {
	if len(req.Remark) < 1 {
		return fmt.Errorf("remark must be at least 1 character")
	}
	if len(req.Remark) > 40 {
		return fmt.Errorf("remark must be less than 40 characters")
	}
	if req.Port < 1 || req.Port > 65535 {
		return fmt.Errorf("invalid port")
	}
	if req.Inbound.ConfigProfileUUID == nil || req.Inbound.ConfigProfileInboundUUID == nil {
		return fmt.Errorf("inbound configProfileUuid and configProfileInboundUuid are required")
	}
	if _, err := uuid.Parse(*req.Inbound.ConfigProfileUUID); err != nil {
		return fmt.Errorf("invalid configProfileUuid")
	}
	if _, err := uuid.Parse(*req.Inbound.ConfigProfileInboundUUID); err != nil {
		return fmt.Errorf("invalid configProfileInboundUuid")
	}
	if req.Tag != nil && *req.Tag != "" {
		if !hostTagRegex.MatchString(*req.Tag) {
			return fmt.Errorf("invalid tag format")
		}
		if len(*req.Tag) > 32 {
			return fmt.Errorf("tag must be less than 32 characters")
		}
	} else if req.Tag != nil && *req.Tag == "" {
		return fmt.Errorf("tag cannot be empty")
	}
	if req.ServerDescription != nil && len(*req.ServerDescription) > 30 {
		return fmt.Errorf("serverDescription must be less than 30 characters")
	}
	if req.VLESSRouteID != nil && (*req.VLESSRouteID < 0 || *req.VLESSRouteID > 65535) {
		return fmt.Errorf("vlessRouteId must be between 0 and 65535")
	}
	if req.SecurityLayer != nil {
		if strings.TrimSpace(*req.SecurityLayer) == "" {
			return fmt.Errorf("invalid securityLayer")
		}
		if _, ok := allowedSecurityLayers[strings.ToUpper(*req.SecurityLayer)]; !ok {
			return fmt.Errorf("invalid securityLayer")
		}
	}
	if req.ALPN != nil && *req.ALPN != "" {
		if _, ok := allowedAlpn[*req.ALPN]; !ok {
			return fmt.Errorf("invalid alpn")
		}
	} else if req.ALPN != nil && *req.ALPN == "" {
		return fmt.Errorf("invalid alpn")
	}
	if req.Fingerprint != nil && *req.Fingerprint != "" {
		if _, ok := allowedFingerprints[*req.Fingerprint]; !ok {
			return fmt.Errorf("invalid fingerprint")
		}
	} else if req.Fingerprint != nil && *req.Fingerprint == "" {
		return fmt.Errorf("invalid fingerprint")
	}
	if req.XrayJSONTemplateUUID != nil && *req.XrayJSONTemplateUUID != "" {
		if _, err := uuid.Parse(*req.XrayJSONTemplateUUID); err != nil {
			return fmt.Errorf("invalid xrayJsonTemplateUuid")
		}
	} else if req.XrayJSONTemplateUUID != nil && *req.XrayJSONTemplateUUID == "" {
		return fmt.Errorf("invalid xrayJsonTemplateUuid")
	}
	if err := validateUUIDList(req.Nodes); err != nil {
		return err
	}
	if err := validateUUIDList(req.ExcludedInternalSquads); err != nil {
		return err
	}
	if err := validateTemplateTypes(req.ExcludeFromSubscription); err != nil {
		return err
	}
	return nil
}

func validateUpdateRequest(req HostUpdateRequestAPI) error {
	if req.Remark.Set {
		if req.Remark.Value == nil {
			return fmt.Errorf("remark cannot be null")
		}
		if len(*req.Remark.Value) > 40 {
			return fmt.Errorf("remark must be less than 40 characters")
		}
	}
	if req.Address.Set && req.Address.Value == nil {
		return fmt.Errorf("address cannot be null")
	}
	if req.Path.Set && req.Path.Value == nil {
		return fmt.Errorf("path cannot be null")
	}
	if req.SNI.Set && req.SNI.Value == nil {
		return fmt.Errorf("sni cannot be null")
	}
	if req.Host.Set && req.Host.Value == nil {
		return fmt.Errorf("host cannot be null")
	}
	if req.Port != nil && (*req.Port < 1 || *req.Port > 65535) {
		return fmt.Errorf("invalid port")
	}
	if req.Tag.Set {
		if req.Tag.Value == nil {
			// nullable
		} else if *req.Tag.Value == "" {
			return fmt.Errorf("tag cannot be empty")
		} else if !hostTagRegex.MatchString(*req.Tag.Value) {
			return fmt.Errorf("invalid tag format")
		} else if len(*req.Tag.Value) > 32 {
			return fmt.Errorf("tag must be less than 32 characters")
		}
	}
	if req.ServerDescription.Set && req.ServerDescription.Value != nil {
		if len(*req.ServerDescription.Value) > 30 {
			return fmt.Errorf("serverDescription must be less than 30 characters")
		}
	}
	if req.VLESSRouteID.Set && req.VLESSRouteID.Value != nil {
		if *req.VLESSRouteID.Value < 0 || *req.VLESSRouteID.Value > 65535 {
			return fmt.Errorf("vlessRouteId must be between 0 and 65535")
		}
	}
	if req.SecurityLayer != nil {
		if strings.TrimSpace(*req.SecurityLayer) == "" {
			return fmt.Errorf("invalid securityLayer")
		}
		if _, ok := allowedSecurityLayers[strings.ToUpper(*req.SecurityLayer)]; !ok {
			return fmt.Errorf("invalid securityLayer")
		}
	}
	if req.ALPN.Set {
		if req.ALPN.Value == nil {
			// nullable
		} else if *req.ALPN.Value == "" {
			return fmt.Errorf("invalid alpn")
		} else if _, ok := allowedAlpn[*req.ALPN.Value]; !ok {
			return fmt.Errorf("invalid alpn")
		}
	}
	if req.Fingerprint.Set {
		if req.Fingerprint.Value == nil {
			// nullable
		} else if *req.Fingerprint.Value == "" {
			return fmt.Errorf("invalid fingerprint")
		} else if _, ok := allowedFingerprints[*req.Fingerprint.Value]; !ok {
			return fmt.Errorf("invalid fingerprint")
		}
	}
	if req.XrayJSONTemplateUUID.Set {
		if req.XrayJSONTemplateUUID.Value == nil {
			// nullable
		} else if *req.XrayJSONTemplateUUID.Value == "" {
			return fmt.Errorf("invalid xrayJsonTemplateUuid")
		} else if _, err := uuid.Parse(*req.XrayJSONTemplateUUID.Value); err != nil {
			return fmt.Errorf("invalid xrayJsonTemplateUuid")
		}
	}
	if req.Inbound != nil {
		if req.Inbound.ConfigProfileUUID == nil || req.Inbound.ConfigProfileInboundUUID == nil {
			return fmt.Errorf("inbound configProfileUuid and configProfileInboundUuid are required")
		}
		if _, err := uuid.Parse(*req.Inbound.ConfigProfileUUID); err != nil {
			return fmt.Errorf("invalid configProfileUuid")
		}
		if _, err := uuid.Parse(*req.Inbound.ConfigProfileInboundUUID); err != nil {
			return fmt.Errorf("invalid configProfileInboundUuid")
		}
	}
	if err := validateUUIDList(req.Nodes); err != nil {
		return err
	}
	if err := validateUUIDList(req.ExcludedInternalSquads); err != nil {
		return err
	}
	if err := validateTemplateTypes(req.ExcludeFromSubscription); err != nil {
		return err
	}
	return nil
}

func validateUUIDList(values []string) error {
	for _, value := range values {
		if _, err := uuid.Parse(value); err != nil {
			return fmt.Errorf("invalid uuid value")
		}
	}
	return nil
}

func validateTemplateTypes(values []string) error {
	for _, value := range values {
		if _, ok := allowedTemplateTypes[value]; !ok {
			return fmt.Errorf("invalid subscription template type")
		}
	}
	return nil
}

func normalizeOptionalStringAllowEmpty(value *string) interface{} {
	if value == nil {
		return nil
	}
	return strings.TrimSpace(*value)
}

func normalizeSecurityLayer(value *string) string {
	if value == nil || strings.TrimSpace(*value) == "" {
		return "DEFAULT"
	}
	upper := strings.ToUpper(strings.TrimSpace(*value))
	if _, ok := allowedSecurityLayers[upper]; ok {
		return upper
	}
	return "DEFAULT"
}

func normalizeJSONField(raw *json.RawMessage, emptyObjectAsNull bool) (bool, []byte, error) {
	if raw == nil {
		return false, nil, nil
	}
	trimmed := strings.TrimSpace(string(*raw))
	if trimmed == "" || trimmed == "null" {
		return true, nil, nil
	}
	if !json.Valid(*raw) {
		return true, nil, fmt.Errorf("invalid JSON payload")
	}
	if emptyObjectAsNull {
		var obj map[string]any
		if err := json.Unmarshal(*raw, &obj); err == nil {
			if len(obj) == 0 {
				return true, nil, nil
			}
		}
	}
	return true, []byte(*raw), nil
}

func normalizeJSONValue(raw *json.RawMessage, emptyObjectAsNull bool) ([]byte, error) {
	_, val, err := normalizeJSONField(raw, emptyObjectAsNull)
	return val, err
}

func normalizeOptionalJSONField(raw OptionalJSON, emptyObjectAsNull bool) (bool, []byte, error) {
	if !raw.Set {
		return false, nil, nil
	}
	value := json.RawMessage(raw.Raw)
	return normalizeJSONField(&value, emptyObjectAsNull)
}

func coalesceBool(value *bool, fallback bool) bool {
	if value == nil {
		return fallback
	}
	return *value
}

func mapHostRecordToAPI(rec hostRecord, nodes []string, excluded []string) HostAPI {
	return HostAPI{
		UUID:                   rec.UUID,
		ViewPosition:           rec.ViewPosition,
		Remark:                 rec.Remark,
		Address:                rec.Address,
		Port:                   rec.Port,
		Path:                   rec.Path,
		SNI:                    rec.SNI,
		Host:                   rec.Host,
		ALPN:                   rec.ALPN,
		Fingerprint:            rec.Fingerprint,
		IsDisabled:             rec.IsDisabled,
		SecurityLayer:          rec.SecurityLayer,
		XHTTPExtraParams:       parseJSONAny(rec.XHTTPExtraParams),
		MuxParams:              parseJSONAny(rec.MuxParams),
		SockoptParams:          parseJSONAny(rec.SockoptParams),
		Inbound:                HostInbound{ConfigProfileUUID: rec.ConfigProfileUUID, ConfigProfileInboundUUID: rec.ConfigProfileInboundUUID},
		ServerDescription:      rec.ServerDescription,
		Tag:                    rec.Tag,
		IsHidden:               rec.IsHidden,
		OverrideSNIFromAddress: rec.OverrideSNIFromAddress,
		KeepSNIBlank:           rec.KeepSNIBlank,
		VLESSRouteID:           rec.VLESSRouteID,
		AllowInsecure:          rec.AllowInsecure,
		ShuffleHost:            rec.ShuffleHost,
		MihomoX25519:           rec.MihomoX25519,
		Nodes:                  ensureStringSlice(nodes),
		XrayJSONTemplateUUID:   rec.XrayJSONTemplateUUID,
		ExcludedInternalSquads: ensureStringSlice(excluded),
		ExcludeFromSubscription: func() []string {
			if len(rec.ExcludeTypes) == 0 {
				return []string{}
			}
			return rec.ExcludeTypes
		}(),
	}
}

func ensureStringSlice(values []string) []string {
	if values == nil {
		return []string{}
	}
	return values
}

func parseJSONAny(raw json.RawMessage) interface{} {
	if len(raw) == 0 {
		return nil
	}
	var value interface{}
	if err := json.Unmarshal(raw, &value); err != nil {
		return nil
	}
	return value
}

func getHosts(ctx context.Context, manager *dbmanager.DatabaseManager) ([]hostRecord, error) {
	var hosts []hostRecord
	var err error
	err = manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		rows, err := db.QueryContext(ctx, `
            SELECT
                uuid, view_position, remark, address, port,
                path, sni, host, alpn, fingerprint, security_layer,
                xhttp_extra_params, mux_params, sockopt_params,
                is_disabled, server_description, vless_route_id,
                allow_insecure, shuffle_host, mihomo_x25519,
                xray_json_template_uuid, keep_sni_blank,
                tag, is_hidden, override_sni_from_address,
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
                xhttp_extra_params, mux_params, sockopt_params,
                is_disabled, server_description, vless_route_id,
                allow_insecure, shuffle_host, mihomo_x25519,
                xray_json_template_uuid, keep_sni_blank,
                tag, is_hidden, override_sni_from_address,
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
	var serverDescription, tag sql.NullString
	var vlessRouteID sql.NullInt64
	var xrayJSONTemplateUUID, configProfileUUID, configProfileInboundUUID sql.NullString
	var isDisabled, allowInsecure, shuffleHost, mihomoX25519, keepSNIBlank, isHidden, overrideSNIFromAddress sql.NullBool
	var xhttpExtraParams, muxParams, sockoptParams []byte
	var excludeTypes dbutil.StringArray

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
	if len(sockoptParams) > 0 {
		rec.SockoptParams = json.RawMessage(sockoptParams)
	}
	if isDisabled.Valid {
		rec.IsDisabled = isDisabled.Bool
	}
	if serverDescription.Valid {
		rec.ServerDescription = &serverDescription.String
	}
	if vlessRouteID.Valid {
		rec.VLESSRouteID = &vlessRouteID.Int64
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
	if xrayJSONTemplateUUID.Valid {
		rec.XrayJSONTemplateUUID = &xrayJSONTemplateUUID.String
	}
	if keepSNIBlank.Valid {
		rec.KeepSNIBlank = keepSNIBlank.Bool
	}
	if tag.Valid {
		rec.Tag = &tag.String
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
	var tags []string
	var err error
	err = manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		rows, err := db.QueryContext(ctx, `SELECT DISTINCT tag FROM hosts WHERE tag IS NOT NULL AND tag <> '' ORDER BY tag ASC`)
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
