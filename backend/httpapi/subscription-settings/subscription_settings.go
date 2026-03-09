package subscriptionsettings

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"v2ray-stat/backend/config"
	dbmanager "v2ray-stat/backend/db/manager"
	"v2ray-stat/backend/httpapi/shared"

	"github.com/google/uuid"
)

type SubscriptionSettings struct {
	UUID                        string    `json:"uuid"`
	ProfileTitle                string    `json:"profile_title"`
	SupportLink                 string    `json:"support_link"`
	ProfileUpdateInterval       int       `json:"profile_update_interval"`
	Address                     string    `json:"address"`
	Port                        int       `json:"port"`
	APISchema                   string    `json:"api_schema"`
	APIPath                     string    `json:"api_path"`
	HappAnnounce                string    `json:"happ_announce"`
	HappRouting                 string    `json:"happ_routing"`
	CreatedAt                   time.Time `json:"created_at"`
	UpdatedAt                   time.Time `json:"updated_at"`
	IsProfileWebpageURLEnabled  bool      `json:"is_profile_webpage_url_enabled"`
	ServeJSONAtBaseSubscription bool      `json:"serve_json_at_base_subscription"`
	IsShowCustomRemarks         bool      `json:"is_show_custom_remarks"`
	CustomResponseHeaders       string    `json:"custom_response_headers"`
	RandomizeHosts              bool      `json:"randomize_hosts"`
	ResponseRules               string    `json:"response_rules"`
	HWIDSettings                string    `json:"hwid_settings"`
	CustomRemarks               string    `json:"custom_remarks"`
}

type SubscriptionSettingsAPI struct {
	UUID                        string            `json:"uuid"`
	ProfileTitle                string            `json:"profileTitle"`
	SupportLink                 string            `json:"supportLink"`
	ProfileUpdateInterval       int               `json:"profileUpdateInterval"`
	IsProfileWebpageUrlEnabled  bool              `json:"isProfileWebpageUrlEnabled"`
	ServeJsonAtBaseSubscription bool              `json:"serveJsonAtBaseSubscription"`
	IsShowCustomRemarks         bool              `json:"isShowCustomRemarks"`
	CustomRemarks               map[string]any    `json:"customRemarks"`
	HappAnnounce                *string           `json:"happAnnounce"`
	HappRouting                 *string           `json:"happRouting"`
	CustomResponseHeaders       map[string]string `json:"customResponseHeaders"`
	RandomizeHosts              bool              `json:"randomizeHosts"`
	ResponseRules               map[string]any    `json:"responseRules"`
	HwidSettings                map[string]any    `json:"hwidSettings"`
	CreatedAt                   time.Time         `json:"createdAt"`
	UpdatedAt                   time.Time         `json:"updatedAt"`
}

type SubscriptionSettingsUpdateRequestAPI struct {
	UUID                        string           `json:"uuid"`
	ProfileTitle                *string          `json:"profileTitle,omitempty"`
	SupportLink                 *string          `json:"supportLink,omitempty"`
	ProfileUpdateInterval       *int             `json:"profileUpdateInterval,omitempty"`
	IsProfileWebpageUrlEnabled  *bool            `json:"isProfileWebpageUrlEnabled,omitempty"`
	ServeJsonAtBaseSubscription *bool            `json:"serveJsonAtBaseSubscription,omitempty"`
	HappAnnounce                *json.RawMessage `json:"happAnnounce,omitempty"`
	HappRouting                 *json.RawMessage `json:"happRouting,omitempty"`
	IsShowCustomRemarks         *bool            `json:"isShowCustomRemarks,omitempty"`
	CustomRemarks               *json.RawMessage `json:"customRemarks,omitempty"`
	CustomResponseHeaders       *json.RawMessage `json:"customResponseHeaders,omitempty"`
	RandomizeHosts              *bool            `json:"randomizeHosts,omitempty"`
	ResponseRules               *json.RawMessage `json:"responseRules,omitempty"`
	HwidSettings                *json.RawMessage `json:"hwidSettings,omitempty"`
}

type SubscriptionSettingsUpdateRequest struct {
	ProfileTitle                *string `json:"profile_title,omitempty"`
	SupportLink                 *string `json:"support_link,omitempty"`
	ProfileUpdateInterval       *int    `json:"profile_update_interval,omitempty"`
	Address                     *string `json:"address,omitempty"`
	Port                        *int    `json:"port,omitempty"`
	APISchema                   *string `json:"api_schema,omitempty"`
	APIPath                     *string `json:"api_path,omitempty"`
	HappAnnounce                *string `json:"happ_announce,omitempty"`
	HappRouting                 *string `json:"happ_routing,omitempty"`
	IsProfileWebpageURLEnabled  *bool   `json:"is_profile_webpage_url_enabled,omitempty"`
	ServeJSONAtBaseSubscription *bool   `json:"serve_json_at_base_subscription,omitempty"`
	IsShowCustomRemarks         *bool   `json:"is_show_custom_remarks,omitempty"`
	CustomResponseHeaders       *string `json:"custom_response_headers,omitempty"`
	RandomizeHosts              *bool   `json:"randomize_hosts,omitempty"`
	ResponseRules               *string `json:"response_rules,omitempty"`
	HWIDSettings                *string `json:"hwid_settings,omitempty"`
	CustomRemarks               *string `json:"custom_remarks,omitempty"`
}

func (r *SubscriptionSettingsUpdateRequest) Validate() error {
	if r.ProfileTitle != nil && strings.TrimSpace(*r.ProfileTitle) == "" {
		return fmt.Errorf("profile_title cannot be empty")
	}
	if r.SupportLink != nil && strings.TrimSpace(*r.SupportLink) == "" {
		return fmt.Errorf("support_link cannot be empty")
	}
	if r.ProfileUpdateInterval != nil && *r.ProfileUpdateInterval < 1 {
		return fmt.Errorf("profile_update_interval must be >= 1")
	}
	if r.Address != nil && strings.TrimSpace(*r.Address) == "" {
		return fmt.Errorf("address cannot be empty")
	}
	if r.Port != nil && (*r.Port < 1 || *r.Port > 65535) {
		return fmt.Errorf("port must be between 1 and 65535")
	}
	if r.APISchema != nil {
		schema := strings.ToLower(strings.TrimSpace(*r.APISchema))
		switch schema {
		case "grpc", "grpcs", "http", "https", "tls":
		default:
			return fmt.Errorf("api_schema must be one of: grpc, grpcs, http, https, tls")
		}
	}
	if r.CustomResponseHeaders != nil && strings.TrimSpace(*r.CustomResponseHeaders) != "" {
		if !json.Valid([]byte(*r.CustomResponseHeaders)) {
			return fmt.Errorf("custom_response_headers must be valid JSON")
		}
	}
	if r.ResponseRules != nil && strings.TrimSpace(*r.ResponseRules) != "" {
		if !json.Valid([]byte(*r.ResponseRules)) {
			return fmt.Errorf("response_rules must be valid JSON")
		}
	}
	if r.HWIDSettings != nil && strings.TrimSpace(*r.HWIDSettings) != "" {
		if !json.Valid([]byte(*r.HWIDSettings)) {
			return fmt.Errorf("hwid_settings must be valid JSON")
		}
	}
	if r.CustomRemarks != nil && strings.TrimSpace(*r.CustomRemarks) != "" {
		if !json.Valid([]byte(*r.CustomRemarks)) {
			return fmt.Errorf("custom_remarks must be valid JSON")
		}
	}
	return nil
}

func (r *SubscriptionSettingsUpdateRequest) HasUpdates() bool {
	return r.ProfileTitle != nil ||
		r.SupportLink != nil ||
		r.ProfileUpdateInterval != nil ||
		r.Address != nil ||
		r.Port != nil ||
		r.APISchema != nil ||
		r.APIPath != nil ||
		r.HappAnnounce != nil ||
		r.HappRouting != nil ||
		r.IsProfileWebpageURLEnabled != nil ||
		r.ServeJSONAtBaseSubscription != nil ||
		r.IsShowCustomRemarks != nil ||
		r.CustomResponseHeaders != nil ||
		r.RandomizeHosts != nil ||
		r.ResponseRules != nil ||
		r.HWIDSettings != nil ||
		r.CustomRemarks != nil
}

func ScanSubscriptionSettings(scanner shared.RowScanner) (SubscriptionSettings, error) {
	var s SubscriptionSettings
	var address, apiSchema, apiPath sql.NullString
	var port sql.NullInt64
	var happAnnounce, happRouting sql.NullString
	var customResponseHeaders, responseRules, hwidSettings, customRemarks sql.NullString
	var isProfileWebpageURLEnabled, serveJSONAtBaseSubscription, isShowCustomRemarks, randomizeHosts sql.NullBool

	err := scanner.Scan(
		&s.UUID,
		&s.ProfileTitle,
		&s.SupportLink,
		&s.ProfileUpdateInterval,
		&address,
		&port,
		&apiSchema,
		&apiPath,
		&happAnnounce,
		&happRouting,
		&s.CreatedAt,
		&s.UpdatedAt,
		&isProfileWebpageURLEnabled,
		&serveJSONAtBaseSubscription,
		&isShowCustomRemarks,
		&customResponseHeaders,
		&randomizeHosts,
		&responseRules,
		&hwidSettings,
		&customRemarks,
	)
	if err != nil {
		return s, err
	}

	if happAnnounce.Valid {
		s.HappAnnounce = happAnnounce.String
	}
	if address.Valid {
		s.Address = address.String
	}
	if port.Valid {
		s.Port = int(port.Int64)
	}
	if apiSchema.Valid {
		s.APISchema = apiSchema.String
	} else {
		s.APISchema = "grpc"
	}
	if apiPath.Valid {
		s.APIPath = apiPath.String
	}
	if happRouting.Valid {
		s.HappRouting = happRouting.String
	}
	if isProfileWebpageURLEnabled.Valid {
		s.IsProfileWebpageURLEnabled = isProfileWebpageURLEnabled.Bool
	}
	if serveJSONAtBaseSubscription.Valid {
		s.ServeJSONAtBaseSubscription = serveJSONAtBaseSubscription.Bool
	}
	if isShowCustomRemarks.Valid {
		s.IsShowCustomRemarks = isShowCustomRemarks.Bool
	}
	if customResponseHeaders.Valid {
		s.CustomResponseHeaders = customResponseHeaders.String
	}
	if randomizeHosts.Valid {
		s.RandomizeHosts = randomizeHosts.Bool
	}
	if responseRules.Valid {
		s.ResponseRules = responseRules.String
	}
	if hwidSettings.Valid {
		s.HWIDSettings = hwidSettings.String
	}
	if customRemarks.Valid {
		s.CustomRemarks = customRemarks.String
	}

	return s, nil
}

func SubscriptionSettingsHandler(manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handleGetSubscriptionSettingsV2RS(w, r, manager, cfg)
		case http.MethodPatch:
			handlePatchSubscriptionSettingsV2RS(w, r, manager, cfg)
		default:
			http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		}
	}
}

func SubscriptionSettingsByUUIDHandler(manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		settingsUUID := strings.TrimSpace(strings.TrimPrefix(r.URL.Path, "/api/v1/subscription-settings/"))
		if _, err := uuid.Parse(settingsUUID); err != nil {
			shared.SendError(w, http.StatusBadRequest, "invalid UUID format", nil, cfg)
			return
		}

		switch r.Method {
		case http.MethodGet:
			handleGetSubscriptionSettingsByUUID(w, r, manager, cfg, settingsUUID)
		case http.MethodPatch:
			handlePatchSubscriptionSettings(w, r, manager, cfg, settingsUUID)
		default:
			http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		}
	}
}

func handleGetSubscriptionSettings(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) {
	var settings SubscriptionSettings
	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		row := db.QueryRowContext(r.Context(), `
			SELECT
				uuid, profile_title, support_link, profile_update_interval,
				address, port, api_schema, api_path,
				happ_announce, happ_routing, created_at, updated_at,
				is_profile_webpage_url_enabled, serve_json_at_base_subscription,
				is_show_custom_remarks, custom_response_headers, randomize_hosts,
				response_rules, hwid_settings, custom_remarks
			FROM subscription_settings
			ORDER BY created_at ASC
			LIMIT 1`)
		var scanErr error
		settings, scanErr = ScanSubscriptionSettings(row)
		return scanErr
	})
	if err != nil {
		if err == sql.ErrNoRows {
			shared.SendError(w, http.StatusNotFound, "subscription settings not found", nil, cfg)
		} else {
			shared.SendError(w, http.StatusInternalServerError, "failed to fetch subscription settings", err, cfg)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{"settings": settings})
}

func handleGetSubscriptionSettingsV2RS(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) {
	var settings SubscriptionSettings
	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		row := db.QueryRowContext(r.Context(), `
            SELECT
                uuid, profile_title, support_link, profile_update_interval,
                address, port, api_schema, api_path,
                happ_announce, happ_routing,
                created_at, updated_at,
                is_profile_webpage_url_enabled, serve_json_at_base_subscription,
                is_show_custom_remarks, custom_response_headers,
                randomize_hosts, response_rules, hwid_settings, custom_remarks
            FROM subscription_settings
            ORDER BY created_at DESC
            LIMIT 1`)
		var scanErr error
		settings, scanErr = ScanSubscriptionSettings(row)
		return scanErr
	})
	if err != nil {
		shared.SendError(w, http.StatusNotFound, "subscription settings not found", err, cfg)
		return
	}

	apiSettings, err := convertSubscriptionSettingsToAPI(settings)
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to parse subscription settings", err, cfg)
		return
	}

	shared.WriteJSON(w, http.StatusOK, map[string]any{"response": apiSettings})
}

func handlePatchSubscriptionSettingsV2RS(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) {
	var req SubscriptionSettingsUpdateRequestAPI
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		shared.SendError(w, http.StatusBadRequest, "invalid JSON", err, cfg)
		return
	}
	if _, err := uuid.Parse(req.UUID); err != nil {
		shared.SendError(w, http.StatusBadRequest, "invalid UUID format", nil, cfg)
		return
	}

	if err := validateSubscriptionSettingsUpdate(req); err != nil {
		shared.SendError(w, http.StatusBadRequest, err.Error(), nil, cfg)
		return
	}

	updates := []string{}
	args := []any{}
	add := func(column string, value any) {
		updates = append(updates, fmt.Sprintf("%s = ?", column))
		args = append(args, value)
	}

	if req.ProfileTitle != nil {
		add("profile_title", strings.TrimSpace(*req.ProfileTitle))
	}
	if req.SupportLink != nil {
		add("support_link", strings.TrimSpace(*req.SupportLink))
	}
	if req.ProfileUpdateInterval != nil {
		add("profile_update_interval", *req.ProfileUpdateInterval)
	}
	if req.IsProfileWebpageUrlEnabled != nil {
		add("is_profile_webpage_url_enabled", *req.IsProfileWebpageUrlEnabled)
	}
	if req.ServeJsonAtBaseSubscription != nil {
		add("serve_json_at_base_subscription", *req.ServeJsonAtBaseSubscription)
	}
	if req.IsShowCustomRemarks != nil {
		add("is_show_custom_remarks", *req.IsShowCustomRemarks)
	}
	if req.RandomizeHosts != nil {
		add("randomize_hosts", *req.RandomizeHosts)
	}

	if set, val, err := parseOptionalString(req.HappAnnounce, 200); err != nil {
		shared.SendError(w, http.StatusBadRequest, err.Error(), nil, cfg)
		return
	} else if set {
		add("happ_announce", val)
	}
	if set, val, err := parseOptionalString(req.HappRouting, 0); err != nil {
		shared.SendError(w, http.StatusBadRequest, err.Error(), nil, cfg)
		return
	} else if set {
		add("happ_routing", val)
	}

	if set, val, err := parseOptionalJSONMap(req.CustomRemarks, true); err != nil {
		shared.SendError(w, http.StatusBadRequest, err.Error(), nil, cfg)
		return
	} else if set {
		add("custom_remarks", val)
	}

	if set, val, err := parseOptionalJSONMap(req.ResponseRules, true); err != nil {
		shared.SendError(w, http.StatusBadRequest, err.Error(), nil, cfg)
		return
	} else if set {
		add("response_rules", val)
	}

	if set, val, err := parseOptionalJSONMap(req.HwidSettings, true); err != nil {
		shared.SendError(w, http.StatusBadRequest, err.Error(), nil, cfg)
		return
	} else if set {
		add("hwid_settings", val)
	}

	if set, val, err := parseOptionalHeaders(req.CustomResponseHeaders); err != nil {
		shared.SendError(w, http.StatusBadRequest, err.Error(), nil, cfg)
		return
	} else if set {
		add("custom_response_headers", val)
	}

	if len(updates) == 0 {
		shared.SendError(w, http.StatusBadRequest, "no fields to update", nil, cfg)
		return
	}

	updates = append(updates, "updated_at = CURRENT_TIMESTAMP")
	args = append(args, req.UUID)

	var settings SubscriptionSettings
	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		query := fmt.Sprintf(`
            UPDATE subscription_settings
            SET %s
            WHERE uuid = ?
            RETURNING uuid, profile_title, support_link, profile_update_interval,
                address, port, api_schema, api_path,
                happ_announce, happ_routing,
                created_at, updated_at,
                is_profile_webpage_url_enabled, serve_json_at_base_subscription,
                is_show_custom_remarks, custom_response_headers,
                randomize_hosts, response_rules, hwid_settings, custom_remarks
        `, strings.Join(updates, ", "))
		row := db.QueryRowContext(r.Context(), query, args...)
		var scanErr error
		settings, scanErr = ScanSubscriptionSettings(row)
		return scanErr
	})
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to update subscription settings", err, cfg)
		return
	}

	apiSettings, err := convertSubscriptionSettingsToAPI(settings)
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to parse subscription settings", err, cfg)
		return
	}

	shared.WriteJSON(w, http.StatusOK, map[string]any{"response": apiSettings})
}

func convertSubscriptionSettingsToAPI(settings SubscriptionSettings) (SubscriptionSettingsAPI, error) {
	customRemarks, err := parseJSONMap(settings.CustomRemarks)
	if err != nil {
		return SubscriptionSettingsAPI{}, err
	}
	customHeaders, err := parseJSONHeaders(settings.CustomResponseHeaders)
	if err != nil {
		return SubscriptionSettingsAPI{}, err
	}
	responseRules, err := parseJSONMap(settings.ResponseRules)
	if err != nil {
		return SubscriptionSettingsAPI{}, err
	}
	hwidSettings, err := parseJSONMap(settings.HWIDSettings)
	if err != nil {
		return SubscriptionSettingsAPI{}, err
	}

	var happAnnounce *string
	if strings.TrimSpace(settings.HappAnnounce) != "" {
		val := settings.HappAnnounce
		happAnnounce = &val
	}
	var happRouting *string
	if strings.TrimSpace(settings.HappRouting) != "" {
		val := settings.HappRouting
		happRouting = &val
	}

	return SubscriptionSettingsAPI{
		UUID:                        settings.UUID,
		ProfileTitle:                settings.ProfileTitle,
		SupportLink:                 settings.SupportLink,
		ProfileUpdateInterval:       settings.ProfileUpdateInterval,
		IsProfileWebpageUrlEnabled:  settings.IsProfileWebpageURLEnabled,
		ServeJsonAtBaseSubscription: settings.ServeJSONAtBaseSubscription,
		IsShowCustomRemarks:         settings.IsShowCustomRemarks,
		CustomRemarks:               customRemarks,
		HappAnnounce:                happAnnounce,
		HappRouting:                 happRouting,
		CustomResponseHeaders:       customHeaders,
		RandomizeHosts:              settings.RandomizeHosts,
		ResponseRules:               responseRules,
		HwidSettings:                hwidSettings,
		CreatedAt:                   settings.CreatedAt,
		UpdatedAt:                   settings.UpdatedAt,
	}, nil
}

func parseJSONMap(raw string) (map[string]any, error) {
	if strings.TrimSpace(raw) == "" {
		return nil, nil
	}
	var obj map[string]any
	if err := json.Unmarshal([]byte(raw), &obj); err != nil {
		return nil, err
	}
	return obj, nil
}

func parseJSONHeaders(raw string) (map[string]string, error) {
	if strings.TrimSpace(raw) == "" {
		return nil, nil
	}
	var obj map[string]string
	if err := json.Unmarshal([]byte(raw), &obj); err != nil {
		return nil, err
	}
	return obj, nil
}

func parseOptionalString(raw *json.RawMessage, maxLen int) (bool, sql.NullString, error) {
	if raw == nil {
		return false, sql.NullString{}, nil
	}
	if string(*raw) == "null" {
		return true, sql.NullString{Valid: false}, nil
	}
	var value string
	if err := json.Unmarshal(*raw, &value); err != nil {
		return true, sql.NullString{}, fmt.Errorf("invalid string value")
	}
	value = strings.TrimSpace(value)
	if maxLen > 0 && len(value) > maxLen {
		return true, sql.NullString{}, fmt.Errorf("value must be less than %d characters", maxLen)
	}
	return true, sql.NullString{String: value, Valid: true}, nil
}

func parseOptionalJSONMap(raw *json.RawMessage, allowNull bool) (bool, sql.NullString, error) {
	if raw == nil {
		return false, sql.NullString{}, nil
	}
	if string(*raw) == "null" {
		if allowNull {
			return true, sql.NullString{Valid: false}, nil
		}
		return true, sql.NullString{}, fmt.Errorf("null is not allowed")
	}
	var obj map[string]any
	if err := json.Unmarshal(*raw, &obj); err != nil {
		return true, sql.NullString{}, fmt.Errorf("invalid JSON object")
	}
	normalized, err := json.Marshal(obj)
	if err != nil {
		return true, sql.NullString{}, fmt.Errorf("invalid JSON object")
	}
	return true, sql.NullString{String: string(normalized), Valid: true}, nil
}

func parseOptionalHeaders(raw *json.RawMessage) (bool, sql.NullString, error) {
	if raw == nil {
		return false, sql.NullString{}, nil
	}
	if string(*raw) == "null" {
		return true, sql.NullString{Valid: false}, nil
	}
	var obj map[string]string
	if err := json.Unmarshal(*raw, &obj); err != nil {
		return true, sql.NullString{}, fmt.Errorf("customResponseHeaders must be an object")
	}
	for key := range obj {
		if !regexp.MustCompile(`^[a-zA-Z0-9_-]+$`).MatchString(key) {
			return true, sql.NullString{}, fmt.Errorf("invalid header name: %s", key)
		}
	}
	normalized, err := json.Marshal(obj)
	if err != nil {
		return true, sql.NullString{}, fmt.Errorf("invalid customResponseHeaders")
	}
	return true, sql.NullString{String: string(normalized), Valid: true}, nil
}

func validateSubscriptionSettingsUpdate(req SubscriptionSettingsUpdateRequestAPI) error {
	if req.ProfileUpdateInterval != nil && *req.ProfileUpdateInterval < 1 {
		return fmt.Errorf("profileUpdateInterval must be >= 1")
	}
	return nil
}

func handleGetSubscriptionSettingsByUUID(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig, settingsUUID string) {
	var settings SubscriptionSettings
	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		row := db.QueryRowContext(r.Context(), `
			SELECT
				uuid, profile_title, support_link, profile_update_interval,
				address, port, api_schema, api_path,
				happ_announce, happ_routing, created_at, updated_at,
				is_profile_webpage_url_enabled, serve_json_at_base_subscription,
				is_show_custom_remarks, custom_response_headers, randomize_hosts,
				response_rules, hwid_settings, custom_remarks
			FROM subscription_settings
			WHERE uuid = ?`, settingsUUID)
		var scanErr error
		settings, scanErr = ScanSubscriptionSettings(row)
		return scanErr
	})
	if err != nil {
		if err == sql.ErrNoRows {
			shared.SendError(w, http.StatusNotFound, "subscription settings not found", nil, cfg)
		} else {
			shared.SendError(w, http.StatusInternalServerError, "failed to fetch subscription settings", err, cfg)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{"settings": settings})
}

func handlePatchSubscriptionSettings(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig, settingsUUID string) {
	var req SubscriptionSettingsUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		shared.SendError(w, http.StatusBadRequest, "invalid JSON", err, cfg)
		return
	}

	if err := req.Validate(); err != nil {
		shared.SendError(w, http.StatusBadRequest, err.Error(), nil, cfg)
		return
	}

	if !req.HasUpdates() {
		shared.SendError(w, http.StatusBadRequest, "no fields to update", nil, cfg)
		return
	}

	var clauses []string
	var args []interface{}
	add := func(col string, val interface{}) {
		clauses = append(clauses, fmt.Sprintf("%s = ?", col))
		args = append(args, val)
	}

	if req.ProfileTitle != nil {
		add("profile_title", strings.TrimSpace(*req.ProfileTitle))
	}
	if req.SupportLink != nil {
		add("support_link", strings.TrimSpace(*req.SupportLink))
	}
	if req.ProfileUpdateInterval != nil {
		add("profile_update_interval", *req.ProfileUpdateInterval)
	}
	if req.Address != nil {
		add("address", strings.TrimSpace(*req.Address))
	}
	if req.Port != nil {
		add("port", *req.Port)
	}
	if req.APISchema != nil {
		add("api_schema", strings.ToLower(strings.TrimSpace(*req.APISchema)))
	}
	if req.APIPath != nil {
		add("api_path", strings.TrimSpace(*req.APIPath))
	}
	if req.HappAnnounce != nil {
		add("happ_announce", *req.HappAnnounce)
	}
	if req.HappRouting != nil {
		add("happ_routing", *req.HappRouting)
	}
	if req.IsProfileWebpageURLEnabled != nil {
		add("is_profile_webpage_url_enabled", *req.IsProfileWebpageURLEnabled)
	}
	if req.ServeJSONAtBaseSubscription != nil {
		add("serve_json_at_base_subscription", *req.ServeJSONAtBaseSubscription)
	}
	if req.IsShowCustomRemarks != nil {
		add("is_show_custom_remarks", *req.IsShowCustomRemarks)
	}
	if req.CustomResponseHeaders != nil {
		add("custom_response_headers", strings.TrimSpace(*req.CustomResponseHeaders))
	}
	if req.RandomizeHosts != nil {
		add("randomize_hosts", *req.RandomizeHosts)
	}
	if req.ResponseRules != nil {
		add("response_rules", strings.TrimSpace(*req.ResponseRules))
	}
	if req.HWIDSettings != nil {
		add("hwid_settings", strings.TrimSpace(*req.HWIDSettings))
	}
	if req.CustomRemarks != nil {
		add("custom_remarks", strings.TrimSpace(*req.CustomRemarks))
	}

	args = append(args, settingsUUID)
	query := fmt.Sprintf("UPDATE subscription_settings SET %s, updated_at = CURRENT_TIMESTAMP WHERE uuid = ?", strings.Join(clauses, ", "))

	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		result, execErr := db.ExecContext(r.Context(), query, args...)
		if execErr != nil {
			return execErr
		}
		rowsAffected, raErr := result.RowsAffected()
		if raErr != nil {
			return raErr
		}
		if rowsAffected == 0 {
			return sql.ErrNoRows
		}
		return nil
	})
	if err != nil {
		if err == sql.ErrNoRows {
			shared.SendError(w, http.StatusNotFound, "subscription settings not found", nil, cfg)
		} else {
			shared.SendError(w, http.StatusInternalServerError, "update failed", err, cfg)
		}
		return
	}

	handleGetSubscriptionSettingsByUUID(w, r, manager, cfg, settingsUUID)
}
