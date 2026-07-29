package subscriptionsettings

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"exodus/internal/config"
	"exodus/internal/httpapi/shared"

	"github.com/google/uuid"
)

func SubscriptionSettingsHandler(db *sql.DB, cfg *config.BackendConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handleGetSubscriptionSettingsEXODUS(w, r, db, cfg)
		case http.MethodPatch:
			handlePatchSubscriptionSettingsEXODUS(w, r, db, cfg)
		default:
			http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		}
	}
}

func handleGetSubscriptionSettingsEXODUS(w http.ResponseWriter, r *http.Request, db *sql.DB, cfg *config.BackendConfig) {
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
	settings, err := ScanSubscriptionSettings(row)
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

func handlePatchSubscriptionSettingsEXODUS(w http.ResponseWriter, r *http.Request, db *sql.DB, cfg *config.BackendConfig) {
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
	idx := 1
	add := func(column string, value any) {
		updates = append(updates, fmt.Sprintf("%s = $%d", column, idx))
		args = append(args, value)
		idx++
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

	query := fmt.Sprintf(`
		UPDATE subscription_settings
		SET %s
		WHERE uuid = $%d
		RETURNING uuid, profile_title, support_link, profile_update_interval,
			address, port, api_schema, api_path,
			happ_announce, happ_routing,
			created_at, updated_at,
			is_profile_webpage_url_enabled, serve_json_at_base_subscription,
			is_show_custom_remarks, custom_response_headers,
			randomize_hosts, response_rules, hwid_settings, custom_remarks
	`, strings.Join(updates, ", "), idx)
	row := db.QueryRowContext(r.Context(), query, args...)
	settings, err := ScanSubscriptionSettings(row)
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
