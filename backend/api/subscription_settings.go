package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"v2ray-stat/backend/config"
	"v2ray-stat/backend/db/manager"

	"github.com/google/uuid"
)

type SubscriptionSettings struct {
	UUID                        string    `json:"uuid"`
	ProfileTitle                string    `json:"profile_title"`
	SupportLink                 string    `json:"support_link"`
	ProfileUpdateInterval       int       `json:"profile_update_interval"`
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

type SubscriptionSettingsUpdateRequest struct {
	ProfileTitle                *string `json:"profile_title,omitempty"`
	SupportLink                 *string `json:"support_link,omitempty"`
	ProfileUpdateInterval       *int    `json:"profile_update_interval,omitempty"`
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

func scanSubscriptionSettings(scanner RowScanner) (SubscriptionSettings, error) {
	var s SubscriptionSettings
	var happAnnounce, happRouting sql.NullString
	var customResponseHeaders, responseRules, hwidSettings, customRemarks sql.NullString
	var isProfileWebpageURLEnabled, serveJSONAtBaseSubscription, isShowCustomRemarks, randomizeHosts sql.NullBool

	err := scanner.Scan(
		&s.UUID,
		&s.ProfileTitle,
		&s.SupportLink,
		&s.ProfileUpdateInterval,
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

func SubscriptionSettingsHandler(manager *manager.DatabaseManager, cfg *config.BackendConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
			return
		}
		handleGetSubscriptionSettings(w, r, manager, cfg)
	}
}

func SubscriptionSettingsByUUIDHandler(manager *manager.DatabaseManager, cfg *config.BackendConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		settingsUUID := strings.TrimSpace(strings.TrimPrefix(r.URL.Path, "/api/v1/subscription-settings/"))
		if _, err := uuid.Parse(settingsUUID); err != nil {
			sendError(w, http.StatusBadRequest, "invalid UUID format", nil, cfg)
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

func handleGetSubscriptionSettings(w http.ResponseWriter, r *http.Request, manager *manager.DatabaseManager, cfg *config.BackendConfig) {
	var settings SubscriptionSettings
	err := manager.ExecuteHighPriority(func(db *sql.DB) error {
		row := db.QueryRowContext(r.Context(), `
			SELECT
				uuid, profile_title, support_link, profile_update_interval,
				happ_announce, happ_routing, created_at, updated_at,
				is_profile_webpage_url_enabled, serve_json_at_base_subscription,
				is_show_custom_remarks, custom_response_headers, randomize_hosts,
				response_rules, hwid_settings, custom_remarks
			FROM subscription_settings
			ORDER BY created_at ASC
			LIMIT 1`)
		var scanErr error
		settings, scanErr = scanSubscriptionSettings(row)
		return scanErr
	})
	if err != nil {
		if err == sql.ErrNoRows {
			sendError(w, http.StatusNotFound, "subscription settings not found", nil, cfg)
		} else {
			sendError(w, http.StatusInternalServerError, "failed to fetch subscription settings", err, cfg)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{"settings": settings})
}

func handleGetSubscriptionSettingsByUUID(w http.ResponseWriter, r *http.Request, manager *manager.DatabaseManager, cfg *config.BackendConfig, settingsUUID string) {
	var settings SubscriptionSettings
	err := manager.ExecuteHighPriority(func(db *sql.DB) error {
		row := db.QueryRowContext(r.Context(), `
			SELECT
				uuid, profile_title, support_link, profile_update_interval,
				happ_announce, happ_routing, created_at, updated_at,
				is_profile_webpage_url_enabled, serve_json_at_base_subscription,
				is_show_custom_remarks, custom_response_headers, randomize_hosts,
				response_rules, hwid_settings, custom_remarks
			FROM subscription_settings
			WHERE uuid = ?`, settingsUUID)
		var scanErr error
		settings, scanErr = scanSubscriptionSettings(row)
		return scanErr
	})
	if err != nil {
		if err == sql.ErrNoRows {
			sendError(w, http.StatusNotFound, "subscription settings not found", nil, cfg)
		} else {
			sendError(w, http.StatusInternalServerError, "failed to fetch subscription settings", err, cfg)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{"settings": settings})
}

func handlePatchSubscriptionSettings(w http.ResponseWriter, r *http.Request, manager *manager.DatabaseManager, cfg *config.BackendConfig, settingsUUID string) {
	var req SubscriptionSettingsUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, http.StatusBadRequest, "invalid JSON", err, cfg)
		return
	}

	if err := req.Validate(); err != nil {
		sendError(w, http.StatusBadRequest, err.Error(), nil, cfg)
		return
	}

	if !req.HasUpdates() {
		sendError(w, http.StatusBadRequest, "no fields to update", nil, cfg)
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

	err := manager.ExecuteHighPriority(func(db *sql.DB) error {
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
			sendError(w, http.StatusNotFound, "subscription settings not found", nil, cfg)
		} else {
			sendError(w, http.StatusInternalServerError, "update failed", err, cfg)
		}
		return
	}

	handleGetSubscriptionSettingsByUUID(w, r, manager, cfg, settingsUUID)
}
