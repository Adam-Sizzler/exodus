package subscriptionsettings

import (
	"database/sql"
	"encoding/json"
	"errors"
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
			uuid, address, port, api_schema, api_path,
			serve_json_at_base_subscription, is_show_custom_remarks, custom_remarks,
			custom_response_headers, randomize_hosts, response_rules, hwid_settings,
			created_at, updated_at
		FROM subscription_settings
		ORDER BY created_at DESC
		LIMIT 1`)
	settings, err := ScanSubscriptionSettings(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			_, _ = db.ExecContext(r.Context(), `
				INSERT INTO subscription_settings (
					uuid, address, port, api_schema, api_path,
					serve_json_at_base_subscription, is_show_custom_remarks, custom_remarks,
					custom_response_headers, randomize_hosts, response_rules, hwid_settings
				) VALUES (
					'00000000-0000-0000-0000-000000000000', '', 9263, 'grpc', '',
					false, true, '{}'::jsonb,
					'{"profile-title":"exEncodeBase64:exodus","support-url":"https://github.com","profile-update-interval":"12"}'::jsonb,
					false, '[]'::jsonb, '{}'::jsonb
				) ON CONFLICT DO NOTHING
			`)
			rowRetry := db.QueryRowContext(r.Context(), `
				SELECT
					uuid, address, port, api_schema, api_path,
					serve_json_at_base_subscription, is_show_custom_remarks, custom_remarks,
					custom_response_headers, randomize_hosts, response_rules, hwid_settings,
					created_at, updated_at
				FROM subscription_settings
				ORDER BY created_at DESC
				LIMIT 1`)
			if sRetry, errRetry := ScanSubscriptionSettings(rowRetry); errRetry == nil {
				settings = sRetry
				err = nil
			}
		}
		if err != nil {
			shared.SendError(w, http.StatusNotFound, "subscription settings not found", err, cfg)
			return
		}
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

	// Fetch current headers so we can update header keys inside custom_response_headers
	var currentHeadersRaw sql.NullString
	_ = db.QueryRowContext(r.Context(), `SELECT custom_response_headers::text FROM subscription_settings WHERE uuid = $1`, req.UUID).Scan(&currentHeadersRaw)
	headersMap := make(map[string]string)
	if currentHeadersRaw.Valid && currentHeadersRaw.String != "" {
		_ = json.Unmarshal([]byte(currentHeadersRaw.String), &headersMap)
	}

	headersModified := false
	if req.ProfileTitle != nil {
		headersMap["profile-title"] = "exEncodeBase64:" + strings.TrimSpace(*req.ProfileTitle)
		headersModified = true
	}
	if req.SupportLink != nil {
		headersMap["support-url"] = strings.TrimSpace(*req.SupportLink)
		headersModified = true
	}
	if req.ProfileUpdateInterval != nil {
		headersMap["profile-update-interval"] = fmt.Sprintf("%d", *req.ProfileUpdateInterval)
		headersModified = true
	}
	if req.IsProfileWebpageUrlEnabled != nil {
		if *req.IsProfileWebpageUrlEnabled {
			headersMap["profile-web-page-url"] = "{{SUBSCRIPTION_URL}}"
		} else {
			delete(headersMap, "profile-web-page-url")
		}
		headersModified = true
	}
	if set, val, err := parseOptionalString(req.HappAnnounce, 200); err != nil {
		shared.SendError(w, http.StatusBadRequest, err.Error(), nil, cfg)
		return
	} else if set {
		if val != "" {
			headersMap["announce"] = "exEncodeBase64:" + val
		} else {
			delete(headersMap, "announce")
		}
		headersModified = true
	}
	if set, val, err := parseOptionalString(req.HappRouting, 0); err != nil {
		shared.SendError(w, http.StatusBadRequest, err.Error(), nil, cfg)
		return
	} else if set {
		if val != "" {
			headersMap["routing"] = val
		} else {
			delete(headersMap, "routing")
		}
		headersModified = true
	}

	if set, val, err := parseOptionalHeaders(req.CustomResponseHeaders); err != nil {
		shared.SendError(w, http.StatusBadRequest, err.Error(), nil, cfg)
		return
	} else if set {
		userHeaders := make(map[string]string)
		if val != "" && val != "{}" && val != "null" {
			var rawMap map[string]string
			if err := json.Unmarshal([]byte(val), &rawMap); err == nil {
				for k, v := range rawMap {
					trimmedKey := strings.ToLower(strings.TrimSpace(k))
					if trimmedKey != "" {
						userHeaders[trimmedKey] = v
					}
				}
			}
		}
		headersMap = userHeaders
		headersModified = true
	}

	updates := []string{}
	args := []any{}
	idx := 1
	add := func(column string, value any) {
		updates = append(updates, fmt.Sprintf("%s = $%d", column, idx))
		args = append(args, value)
		idx++
	}

	if headersModified {
		hBytes, _ := json.Marshal(headersMap)
		add("custom_response_headers", string(hBytes))
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
		RETURNING uuid, address, port, api_schema, api_path,
			serve_json_at_base_subscription, is_show_custom_remarks, custom_remarks,
			custom_response_headers, randomize_hosts, response_rules, hwid_settings,
			created_at, updated_at
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
