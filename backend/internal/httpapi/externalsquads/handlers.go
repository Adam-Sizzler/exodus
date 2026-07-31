package externalsquads

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"exodus/internal/config"
	"exodus/internal/httpapi/shared"

	"github.com/google/uuid"
)

func ExternalSquadsHandler(db *sql.DB, cfg *config.BackendConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handleGetExternalSquads(w, r, db, cfg)
		case http.MethodPost:
			handleCreateExternalSquad(w, r, db, cfg)
		case http.MethodPatch:
			handleUpdateExternalSquad(w, r, db, cfg)
		default:
			shared.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
		}
	}
}

func ExternalSquadByUUIDHandler(db *sql.DB, cfg *config.BackendConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimSpace(trimExternalSquadsPath(r.URL.Path))

		if path == "" {
			switch r.Method {
			case http.MethodGet:
				handleGetExternalSquads(w, r, db, cfg)
			case http.MethodPost:
				handleCreateExternalSquad(w, r, db, cfg)
			case http.MethodPatch:
				handleUpdateExternalSquad(w, r, db, cfg)
			default:
				shared.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
			}
			return
		}

		parts := strings.Split(path, "/")
		squadUUID := strings.TrimSpace(parts[0])

		if _, err := uuid.Parse(squadUUID); err != nil {
			if r.Method == http.MethodPatch {
				handleUpdateExternalSquad(w, r, db, cfg)
				return
			}
			shared.SendError(w, http.StatusBadRequest, "invalid UUID format", nil, cfg)
			return
		}

		if len(parts) > 1 {
			if len(parts) == 3 && parts[1] == "bulk-actions" && parts[2] == "add-users" {
				if r.Method != http.MethodPost {
					shared.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
					return
				}
				handleBulkAddUsersToExternalSquad(w, r, db, cfg, squadUUID)
				return
			}
			if len(parts) == 3 && parts[1] == "bulk-actions" && parts[2] == "remove-users" {
				if r.Method != http.MethodDelete && r.Method != http.MethodPost {
					shared.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
					return
				}
				handleBulkRemoveUsersFromExternalSquad(w, r, db, cfg, squadUUID)
				return
			}
			shared.WriteJSONError(w, http.StatusNotFound, "endpoint not found")
			return
		}

		switch r.Method {
		case http.MethodGet:
			handleGetExternalSquadByUUID(w, r, db, cfg, squadUUID)
		case http.MethodPatch:
			handleUpdateExternalSquad(w, r, db, cfg)
		case http.MethodDelete:
			handleDeleteExternalSquad(w, r, db, cfg, squadUUID)
		default:
			shared.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
		}
	}
}

func trimExternalSquadsPath(p string) string {
	p = strings.TrimPrefix(p, "/api/external-squads")
	return strings.TrimPrefix(p, "/")
}

func ExternalSquadsReorderHandler(db *sql.DB, cfg *config.BackendConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			shared.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}

		var req ReorderExternalSquadsRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			shared.SendError(w, http.StatusBadRequest, "invalid JSON", err, cfg)
			return
		}

		if len(req.Squads) == 0 {
			req.Squads = req.Items
		}

		if len(req.Squads) == 0 {
			shared.SendError(w, http.StatusBadRequest, "items cannot be empty", nil, cfg)
			return
		}

		for _, item := range req.Squads {
			if _, err := uuid.Parse(item.UUID); err != nil {
				shared.SendError(w, http.StatusBadRequest, "invalid UUID format", nil, cfg)
				return
			}
		}

		tx, err := db.BeginTx(r.Context(), nil)
		if err != nil {
			shared.SendError(w, http.StatusInternalServerError, "failed to start transaction", err, cfg)
			return
		}
		defer func() { _ = tx.Rollback() }()

		for _, item := range req.Squads {
			if _, err := tx.ExecContext(r.Context(),
				`UPDATE external_squads SET view_position = $1 WHERE uuid = $2`,
				item.ViewPosition, item.UUID); err != nil {
				shared.SendError(w, http.StatusInternalServerError, "failed to reorder external squads", err, cfg)
				return
			}
		}

		if err := tx.Commit(); err != nil {
			shared.SendError(w, http.StatusInternalServerError, "failed to commit reorder transaction", err, cfg)
			return
		}

		handleGetExternalSquads(w, r, db, cfg)
	}
}

func handleGetExternalSquads(w http.ResponseWriter, r *http.Request, db *sql.DB, cfg *config.BackendConfig) {
	records, err := getExternalSquads(r.Context(), db)
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to fetch external squads", err, cfg)
		return
	}

	result := make([]ExternalSquadAPI, 0, len(records))
	for _, rec := range records {
		api, err := convertExternalSquadToAPI(rec)
		if err != nil {
			cfg.Logger.Error("Failed to convert external squad", "uuid", rec.UUID, "error", err)
			continue
		}

		_ = db.QueryRowContext(r.Context(), `SELECT COUNT(*) FROM users WHERE external_squad_uuid = $1`, api.UUID).Scan(&api.Info.MembersCount)
		rows, err := db.QueryContext(r.Context(), `SELECT template_uuid, template_type FROM external_squads_templates WHERE external_squad_uuid = $1`, api.UUID)
		if err == nil {
			for rows.Next() {
				var t ExternalSquadTemplate
				if err := rows.Scan(&t.TemplateUUID, &t.TemplateType); err == nil {
					api.Templates = append(api.Templates, t)
				}
			}
			rows.Close()
		}

		result = append(result, api)
	}

	shared.WriteJSON(w, http.StatusOK, map[string]any{
		"response": map[string]any{
			"total":          len(result),
			"externalSquads": result,
		},
	})
}

func handleGetExternalSquadByUUID(w http.ResponseWriter, r *http.Request, db *sql.DB, cfg *config.BackendConfig, squadUUID string) {
	record, err := getExternalSquadByUUID(r.Context(), db, squadUUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			shared.SendError(w, http.StatusNotFound, "external squad not found", nil, cfg)
			return
		}
		shared.SendError(w, http.StatusInternalServerError, "failed to fetch external squad", err, cfg)
		return
	}

	api, err := convertExternalSquadToAPI(record)
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to convert external squad", err, cfg)
		return
	}

	_ = db.QueryRowContext(r.Context(), `SELECT COUNT(*) FROM users WHERE external_squad_uuid = $1`, api.UUID).Scan(&api.Info.MembersCount)
	rows, err := db.QueryContext(r.Context(), `SELECT template_uuid, template_type FROM external_squads_templates WHERE external_squad_uuid = $1`, api.UUID)
	if err == nil {
		for rows.Next() {
			var t ExternalSquadTemplate
			if err := rows.Scan(&t.TemplateUUID, &t.TemplateType); err == nil {
				api.Templates = append(api.Templates, t)
			}
		}
		rows.Close()
	}

	shared.WriteJSON(w, http.StatusOK, map[string]any{"response": api})
}

func handleCreateExternalSquad(w http.ResponseWriter, r *http.Request, db *sql.DB, cfg *config.BackendConfig) {
	var req CreateExternalSquadRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		shared.SendError(w, http.StatusBadRequest, "invalid JSON", err, cfg)
		return
	}

	if err := req.Validate(); err != nil {
		shared.SendError(w, http.StatusBadRequest, err.Error(), nil, cfg)
		return
	}

	squadUUID := uuid.NewString()

	tx, err := db.BeginTx(r.Context(), nil)
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to start transaction", err, cfg)
		return
	}
	defer func() { _ = tx.Rollback() }()

	subSettings, err := marshalJSON(req.SubscriptionSettings)
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to marshal settings", err, cfg)
		return
	}
	hostOverrides, err := marshalJSON(req.HostOverrides)
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to marshal overrides", err, cfg)
		return
	}
	respHeaders, err := marshalJSON(req.ResponseHeaders)
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to marshal headers", err, cfg)
		return
	}
	hwidSettings, err := marshalJSON(req.HWIDSettings)
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to marshal hwid settings", err, cfg)
		return
	}
	customRemarks, err := marshalJSON(req.CustomRemarks)
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to marshal custom remarks", err, cfg)
		return
	}

	_, err = tx.ExecContext(r.Context(), `
		INSERT INTO external_squads (
			uuid, view_position, name,
			subscription_settings, host_overrides, response_headers,
			hwid_settings, custom_remarks, subpage_config_uuid
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`,
		squadUUID,
		coalesceInt(req.ViewPosition, 0),
		strings.TrimSpace(req.Name),
		subSettings,
		hostOverrides,
		respHeaders,
		hwidSettings,
		customRemarks,
		normalizeStringPtr(req.SubpageConfigUUID),
	)
	if err != nil {
		if isUniqueViolation(err, "external_squads_name_key") {
			shared.SendError(w, http.StatusConflict, "external squad name already exists", nil, cfg)
			return
		}
		shared.SendError(w, http.StatusInternalServerError, "failed to create external squad", err, cfg)
		return
	}

	if err := tx.Commit(); err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to commit external squad creation", err, cfg)
		return
	}

	handleGetExternalSquadByUUID(w, r, db, cfg, squadUUID)
}

func handleUpdateExternalSquad(w http.ResponseWriter, r *http.Request, db *sql.DB, cfg *config.BackendConfig) {
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		shared.SendError(w, http.StatusBadRequest, "failed to read body", err, cfg)
		return
	}

	var req UpdateExternalSquadRequest
	if err := json.Unmarshal(bodyBytes, &req); err != nil {
		shared.SendError(w, http.StatusBadRequest, "invalid JSON", err, cfg)
		return
	}

	if err := req.Validate(); err != nil {
		shared.SendError(w, http.StatusBadRequest, err.Error(), nil, cfg)
		return
	}

	if _, err := uuid.Parse(req.UUID); err != nil {
		shared.SendError(w, http.StatusBadRequest, "invalid UUID format", nil, cfg)
		return
	}

	tx, err := db.BeginTx(r.Context(), nil)
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to start transaction", err, cfg)
		return
	}
	defer func() { _ = tx.Rollback() }()

	clauses := make([]string, 0)
	args := make([]any, 0)
	idx := 1

	add := func(column string, value any) {
		clauses = append(clauses, fmt.Sprintf("%s = $%d", column, idx))
		args = append(args, value)
		idx++
	}

	if req.Name != nil {
		add("name", strings.TrimSpace(*req.Name))
	}
	if req.ViewPosition != nil {
		add("view_position", *req.ViewPosition)
	}

	if req.SubscriptionSettings != nil {
		if string(req.SubscriptionSettings) == "null" {
			add("subscription_settings", nil)
		} else {
			add("subscription_settings", string(req.SubscriptionSettings))
		}
	}
	if req.HostOverrides != nil {
		if string(req.HostOverrides) == "null" {
			add("host_overrides", nil)
		} else {
			add("host_overrides", string(req.HostOverrides))
		}
	}
	if req.ResponseHeaders != nil {
		if string(req.ResponseHeaders) == "null" {
			add("response_headers", nil)
		} else {
			add("response_headers", string(req.ResponseHeaders))
		}
	}
	if req.HWIDSettings != nil {
		if string(req.HWIDSettings) == "null" {
			add("hwid_settings", nil)
		} else {
			add("hwid_settings", string(req.HWIDSettings))
		}
	}
	if req.CustomRemarks != nil {
		if string(req.CustomRemarks) == "null" {
			add("custom_remarks", nil)
		} else {
			add("custom_remarks", string(req.CustomRemarks))
		}
	}
	if req.SubpageConfigUUID != nil {
		add("subpage_config_uuid", normalizeStringPtr(req.SubpageConfigUUID))
	}

	if len(clauses) > 0 {
		clauses = append(clauses, "updated_at = CURRENT_TIMESTAMP")
		args = append(args, req.UUID)

		query := fmt.Sprintf(`
			UPDATE external_squads
			SET %s
			WHERE uuid = $%d
		`, strings.Join(clauses, ", "), idx)

		if _, err := tx.ExecContext(r.Context(), query, args...); err != nil {
			if isUniqueViolation(err, "external_squads_name_key") {
				shared.SendError(w, http.StatusConflict, "external squad name already exists", nil, cfg)
				return
			}
			shared.SendError(w, http.StatusInternalServerError, "failed to update external squad", err, cfg)
			return
		}
	} else if req.Templates != nil {
		_, err := tx.ExecContext(r.Context(), `
			UPDATE external_squads 
			SET updated_at = CURRENT_TIMESTAMP 
			WHERE uuid = $1
		`, req.UUID)
		if err != nil {
			shared.SendError(w, http.StatusInternalServerError, "failed to update timestamp", err, cfg)
			return
		}
	}

	if req.Templates != nil {
		_, err = tx.ExecContext(r.Context(), `DELETE FROM external_squads_templates WHERE external_squad_uuid = $1`, req.UUID)
		if err != nil {
			shared.SendError(w, http.StatusInternalServerError, "failed to delete existing templates", err, cfg)
			return
		}

		for _, t := range *req.Templates {
			_, err = tx.ExecContext(r.Context(), `
				INSERT INTO external_squads_templates (external_squad_uuid, template_uuid, template_type)
				VALUES ($1, $2, $3)
			`, req.UUID, t.TemplateUUID, t.TemplateType)
			if err != nil {
				shared.SendError(w, http.StatusInternalServerError, "failed to insert template", err, cfg)
				return
			}
		}
	}

	if len(clauses) == 0 && req.Templates == nil {
		shared.SendError(w, http.StatusBadRequest, "no fields to update", nil, cfg)
		return
	}

	if err := tx.Commit(); err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to commit external squad update", err, cfg)
		return
	}

	handleGetExternalSquadByUUID(w, r, db, cfg, req.UUID)
}

func handleDeleteExternalSquad(w http.ResponseWriter, r *http.Request, db *sql.DB, cfg *config.BackendConfig, squadUUID string) {
	result, err := db.ExecContext(r.Context(), `DELETE FROM external_squads WHERE uuid = $1`, squadUUID)
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to delete external squad", err, cfg)
		return
	}
	rows, err := result.RowsAffected()
	if err != nil || rows == 0 {
		shared.SendError(w, http.StatusNotFound, "external squad not found", nil, cfg)
		return
	}

	shared.WriteJSON(w, http.StatusOK, map[string]any{"response": map[string]any{"isDeleted": true}})
}

func handleBulkAddUsersToExternalSquad(w http.ResponseWriter, r *http.Request, db *sql.DB, cfg *config.BackendConfig, squadUUID string) {
	var req BulkUsersRequest
	bodyBytes, err := io.ReadAll(r.Body)

	hasSpecificUsers := false
	if err == nil && len(bodyBytes) > 0 {
		if jsonErr := json.Unmarshal(bodyBytes, &req); jsonErr == nil && len(req.UserUUIDs) > 0 {
			if validateErr := req.Validate(); validateErr == nil {
				hasSpecificUsers = true
			}
		}
	}

	var affected int64
	var exists int
	if err := db.QueryRowContext(r.Context(), `SELECT 1 FROM external_squads WHERE uuid = $1`, squadUUID).Scan(&exists); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			shared.SendError(w, http.StatusNotFound, "external squad not found", nil, cfg)
			return
		}
		shared.SendError(w, http.StatusInternalServerError, "failed to check external squad", err, cfg)
		return
	}

	if hasSpecificUsers {
		placeholders := make([]string, len(req.UserUUIDs))
		args := make([]any, 0, len(req.UserUUIDs)+1)
		args = append(args, squadUUID)

		for i, u := range req.UserUUIDs {
			placeholders[i] = fmt.Sprintf("$%d", i+2)
			args = append(args, u)
		}

		query := fmt.Sprintf(`
			UPDATE users
			SET external_squad_uuid = $1, updated_at = CURRENT_TIMESTAMP
			WHERE uuid IN (%s)
		`, strings.Join(placeholders, ", "))

		result, err := db.ExecContext(r.Context(), query, args...)
		if err != nil {
			shared.SendError(w, http.StatusInternalServerError, "failed to add users to external squad", err, cfg)
			return
		}
		affected, _ = result.RowsAffected()
	} else {
		result, err := db.ExecContext(r.Context(), `
			UPDATE users
			SET external_squad_uuid = $1::uuid, updated_at = CURRENT_TIMESTAMP
			WHERE external_squad_uuid IS DISTINCT FROM $1::uuid
		`, squadUUID)
		if err != nil {
			shared.SendError(w, http.StatusInternalServerError, "failed to add users to external squad", err, cfg)
			return
		}
		affected, _ = result.RowsAffected()
	}

	cfg.Logger.Info("Users added to external squad", "squad_uuid", squadUUID, "affected_rows", affected)
	shared.WriteJSON(w, http.StatusOK, map[string]any{"response": map[string]any{"eventSent": true}})
}

func handleBulkRemoveUsersFromExternalSquad(w http.ResponseWriter, r *http.Request, db *sql.DB, cfg *config.BackendConfig, squadUUID string) {
	var req BulkUsersRequest
	bodyBytes, err := io.ReadAll(r.Body)

	hasSpecificUsers := false
	if err == nil && len(bodyBytes) > 0 {
		if jsonErr := json.Unmarshal(bodyBytes, &req); jsonErr == nil && len(req.UserUUIDs) > 0 {
			if validateErr := req.Validate(); validateErr == nil {
				hasSpecificUsers = true
			}
		}
	}

	var affected int64
	var exists int
	if err := db.QueryRowContext(r.Context(), `SELECT 1 FROM external_squads WHERE uuid = $1`, squadUUID).Scan(&exists); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			shared.SendError(w, http.StatusNotFound, "external squad not found", nil, cfg)
			return
		}
		shared.SendError(w, http.StatusInternalServerError, "failed to check external squad", err, cfg)
		return
	}

	if hasSpecificUsers {
		placeholders := make([]string, len(req.UserUUIDs))
		args := make([]any, 0, len(req.UserUUIDs)+1)

		for i, u := range req.UserUUIDs {
			placeholders[i] = fmt.Sprintf("$%d", i+1)
			args = append(args, u)
		}
		args = append(args, squadUUID)

		query := fmt.Sprintf(`
			UPDATE users
			SET external_squad_uuid = NULL, updated_at = CURRENT_TIMESTAMP
			WHERE uuid IN (%s) AND external_squad_uuid = $%d
		`, strings.Join(placeholders, ", "), len(req.UserUUIDs)+1)

		result, err := db.ExecContext(r.Context(), query, args...)
		if err != nil {
			shared.SendError(w, http.StatusInternalServerError, "failed to remove users from external squad", err, cfg)
			return
		}
		affected, _ = result.RowsAffected()
	} else {
		result, err := db.ExecContext(r.Context(), `
			UPDATE users
			SET external_squad_uuid = NULL, updated_at = CURRENT_TIMESTAMP
			WHERE external_squad_uuid = $1
		`, squadUUID)
		if err != nil {
			shared.SendError(w, http.StatusInternalServerError, "failed to remove users from external squad", err, cfg)
			return
		}
		affected, _ = result.RowsAffected()
	}

	cfg.Logger.Info("Users removed from external squad", "squad_uuid", squadUUID, "affected_rows", affected)

	shared.WriteJSON(w, http.StatusOK, map[string]any{
		"response": map[string]any{
			"eventSent": true,
		},
	})
}
