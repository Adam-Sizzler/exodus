package subscriptionsettings

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"

	"exodus/internal/httpapi/shared"
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
	if r.ProfileUpdateInterval != nil && *r.ProfileUpdateInterval <= 0 {
		return fmt.Errorf("profile_update_interval must be greater than 0")
	}
	if r.Port != nil && (*r.Port < 1 || *r.Port > 65535) {
		return fmt.Errorf("port must be between 1 and 65535")
	}
	if r.APISchema != nil {
		schema := strings.ToLower(strings.TrimSpace(*r.APISchema))
		if schema != "http" && schema != "https" {
			return fmt.Errorf("api_schema must be 'http' or 'https'")
		}
	}
	return nil
}

func (s SubscriptionSettings) ProfileTitleValue() string {
	return strings.TrimSpace(s.ProfileTitle)
}

func (s SubscriptionSettings) SupportLinkValue() string {
	return strings.TrimSpace(s.SupportLink)
}

func (s SubscriptionSettings) ProfileUpdateIntervalValue() int {
	if s.ProfileUpdateInterval <= 0 {
		return 12
	}
	return s.ProfileUpdateInterval
}

func ScanSubscriptionSettings(scanner shared.RowScanner) (SubscriptionSettings, error) {
	var s SubscriptionSettings
	var (
		profileTitle, supportLink, address, apiSchema, apiPath sql.NullString
		happAnnounce, happRouting, customHeaders               sql.NullString
		responseRules, hwidSettings, customRemarks            sql.NullString
		profileUpdateInterval, port                             sql.NullInt64
		isProfileWebpageURLEnabled, serveJSONAtBase             sql.NullBool
		isShowCustomRemarks, randomizeHosts                    sql.NullBool
	)

	err := scanner.Scan(
		&s.UUID,
		&profileTitle,
		&supportLink,
		&profileUpdateInterval,
		&address,
		&port,
		&apiSchema,
		&apiPath,
		&happAnnounce,
		&happRouting,
		&s.CreatedAt,
		&s.UpdatedAt,
		&isProfileWebpageURLEnabled,
		&serveJSONAtBase,
		&isShowCustomRemarks,
		&customHeaders,
		&randomizeHosts,
		&responseRules,
		&hwidSettings,
		&customRemarks,
	)
	if err != nil {
		return s, err
	}

	if profileTitle.Valid {
		s.ProfileTitle = profileTitle.String
	}
	if supportLink.Valid {
		s.SupportLink = supportLink.String
	}
	if profileUpdateInterval.Valid {
		s.ProfileUpdateInterval = int(profileUpdateInterval.Int64)
	}
	if address.Valid {
		s.Address = address.String
	}
	if port.Valid {
		s.Port = int(port.Int64)
	}
	if apiSchema.Valid {
		s.APISchema = apiSchema.String
	}
	if apiPath.Valid {
		s.APIPath = apiPath.String
	}
	if happAnnounce.Valid {
		s.HappAnnounce = happAnnounce.String
	}
	if happRouting.Valid {
		s.HappRouting = happRouting.String
	}
	if isProfileWebpageURLEnabled.Valid {
		s.IsProfileWebpageURLEnabled = isProfileWebpageURLEnabled.Bool
	}
	if serveJSONAtBase.Valid {
		s.ServeJSONAtBaseSubscription = serveJSONAtBase.Bool
	}
	if isShowCustomRemarks.Valid {
		s.IsShowCustomRemarks = isShowCustomRemarks.Bool
	}
	if customHeaders.Valid {
		s.CustomResponseHeaders = customHeaders.String
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

func validateSubscriptionSettingsUpdate(req SubscriptionSettingsUpdateRequestAPI) error {
	if req.ProfileTitle != nil && strings.TrimSpace(*req.ProfileTitle) == "" {
		return fmt.Errorf("profileTitle cannot be empty")
	}
	if req.SupportLink != nil && strings.TrimSpace(*req.SupportLink) == "" {
		return fmt.Errorf("supportLink cannot be empty")
	}
	if req.ProfileUpdateInterval != nil && *req.ProfileUpdateInterval <= 0 {
		return fmt.Errorf("profileUpdateInterval must be greater than 0")
	}
	return nil
}

func parseOptionalString(raw *json.RawMessage, maxLen int) (bool, string, error) {
	if raw == nil {
		return false, "", nil
	}

	var str string
	if err := json.Unmarshal(*raw, &str); err != nil {
		return false, "", fmt.Errorf("must be a string")
	}

	str = strings.TrimSpace(str)
	if maxLen > 0 && len(str) > maxLen {
		return false, "", fmt.Errorf("length must not exceed %d characters", maxLen)
	}

	return true, str, nil
}

func parseOptionalJSONMap(raw *json.RawMessage, allowEmpty bool) (bool, string, error) {
	if raw == nil {
		return false, "", nil
	}

	var obj map[string]any
	if err := json.Unmarshal(*raw, &obj); err != nil {
		return false, "", fmt.Errorf("must be a valid JSON object")
	}

	if !allowEmpty && len(obj) == 0 {
		return false, "", fmt.Errorf("object cannot be empty")
	}

	bytes, err := json.Marshal(obj)
	if err != nil {
		return false, "", fmt.Errorf("failed to encode JSON object")
	}

	return true, string(bytes), nil
}

func parseOptionalHeaders(raw *json.RawMessage) (bool, string, error) {
	if raw == nil {
		return false, "", nil
	}

	var obj map[string]string
	if err := json.Unmarshal(*raw, &obj); err != nil {
		return false, "", fmt.Errorf("must be a valid JSON object mapping string keys to string values")
	}

	validHeaderName := regexp.MustCompile(`^[A-Za-z0-9-]+$`)
	clean := make(map[string]string, len(obj))
	for k, v := range obj {
		key := strings.TrimSpace(k)
		val := strings.TrimSpace(v)
		if key == "" || !validHeaderName.MatchString(key) {
			return false, "", fmt.Errorf("invalid header name: %s", k)
		}
		if val == "" {
			return false, "", fmt.Errorf("header value for %s cannot be empty", key)
		}
		clean[key] = val
	}

	bytes, err := json.Marshal(clean)
	if err != nil {
		return false, "", fmt.Errorf("failed to encode headers")
	}

	return true, string(bytes), nil
}

func parseJSONMap(raw string) (map[string]any, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return map[string]any{}, nil
	}
	var out map[string]any
	if err := json.Unmarshal([]byte(raw), &out); err != nil {
		return nil, err
	}
	return out, nil
}

func parseJSONHeaders(raw string) (map[string]string, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return map[string]string{}, nil
	}
	var out map[string]string
	if err := json.Unmarshal([]byte(raw), &out); err != nil {
		return nil, err
	}
	return out, nil
}
