package externalsquads

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"exodus/internal/config"
	dbmanager "exodus/internal/db/manager"
	"exodus/internal/httpapi/shared"

	"github.com/google/uuid"
)

var (
	errExternalSquadNotFound = errors.New("external squad not found")
	errExternalSquadExists   = errors.New("external squad name already exists")
)

type ExternalSquadRecord struct {
	UUID                 string          `json:"uuid"`
	ViewPosition         int             `json:"view_position"`
	Name                 string          `json:"name"`
	SubscriptionSettings json.RawMessage `json:"subscription_settings,omitempty"`
	HostOverrides        json.RawMessage `json:"host_overrides,omitempty"`
	ResponseHeaders      json.RawMessage `json:"response_headers,omitempty"`
	HWIDSettings         json.RawMessage `json:"hwid_settings,omitempty"`
	CustomRemarks        json.RawMessage `json:"custom_remarks,omitempty"`
	SubpageConfigUUID    *string         `json:"subpage_config_uuid,omitempty"`
	CreatedAt            time.Time       `json:"created_at"`
	UpdatedAt            time.Time       `json:"updated_at"`
}

type ExternalSquadAPI struct {
	UUID                 string                  `json:"uuid"`
	ViewPosition         int                     `json:"viewPosition"`
	Name                 string                  `json:"name"`
	Info                 ExternalSquadInfo       `json:"info"`
	Templates            []ExternalSquadTemplate `json:"templates"`
	SubscriptionSettings map[string]any          `json:"subscriptionSettings"`
	HostOverrides        map[string]any          `json:"hostOverrides"`
	ResponseHeaders      map[string]string       `json:"responseHeaders"`
	HWIDSettings         map[string]any          `json:"hwidSettings"`
	CustomRemarks        map[string]any          `json:"customRemarks"`
	SubpageConfigUUID    *string                 `json:"subpageConfigUuid"`
	CreatedAt            time.Time               `json:"createdAt"`
	UpdatedAt            time.Time               `json:"updatedAt"`
}

type ExternalSquadInfo struct {
	MembersCount int `json:"membersCount"`
}

type ExternalSquadTemplate struct {
	TemplateUUID string `json:"templateUuid"`
	TemplateType string `json:"templateType"`
}

type CreateExternalSquadRequest struct {
	Name                 string            `json:"name"`
	ViewPosition         *int              `json:"viewPosition,omitempty"`
	SubscriptionSettings map[string]any    `json:"subscriptionSettings,omitempty"`
	HostOverrides        map[string]any    `json:"hostOverrides,omitempty"`
	ResponseHeaders      map[string]string `json:"responseHeaders,omitempty"`
	HWIDSettings         map[string]any    `json:"hwidSettings,omitempty"`
	CustomRemarks        map[string]any    `json:"customRemarks,omitempty"`
	SubpageConfigUUID    *string           `json:"subpageConfigUuid,omitempty"`
}

type UpdateExternalSquadRequest struct {
	UUID                 string                   `json:"uuid"`
	Name                 *string                  `json:"name,omitempty"`
	ViewPosition         *int                     `json:"viewPosition,omitempty"`
	Templates            *[]ExternalSquadTemplate `json:"templates,omitempty"`
	SubscriptionSettings json.RawMessage          `json:"subscriptionSettings,omitempty"`
	HostOverrides        json.RawMessage          `json:"hostOverrides,omitempty"`
	ResponseHeaders      json.RawMessage          `json:"responseHeaders,omitempty"`
	HWIDSettings         json.RawMessage          `json:"hwidSettings,omitempty"`
	CustomRemarks        json.RawMessage          `json:"customRemarks,omitempty"`
	SubpageConfigUUID    *string                  `json:"subpageConfigUuid,omitempty"`
}

type ReorderExternalSquadsRequest struct {
	Items []struct {
		UUID         string `json:"uuid"`
		ViewPosition int    `json:"viewPosition"`
	} `json:"items"`

	Squads []struct {
		UUID         string `json:"uuid"`
		ViewPosition int    `json:"viewPosition"`
	} `json:"squads"`
}

type BulkUsersRequest struct {
	UserUUIDs []string `json:"userUuids"`
}

func (r *CreateExternalSquadRequest) Validate() error {
	if strings.TrimSpace(r.Name) == "" {
		return fmt.Errorf("name is required")
	}
	if len(r.Name) > 30 {
		return fmt.Errorf("name must be less than 30 characters")
	}
	if r.ViewPosition != nil && *r.ViewPosition < 0 {
		return fmt.Errorf("viewPosition must be >= 0")
	}
	return nil
}

func (r *UpdateExternalSquadRequest) Validate() error {
	if r.Name != nil && strings.TrimSpace(*r.Name) == "" {
		return fmt.Errorf("name cannot be empty")
	}
	if r.Name != nil && len(*r.Name) > 30 {
		return fmt.Errorf("name must be less than 30 characters")
	}
	if r.ViewPosition != nil && *r.ViewPosition < 0 {
		return fmt.Errorf("viewPosition must be >= 0")
	}
	return nil
}

func (r *BulkUsersRequest) Validate() error {
	if len(r.UserUUIDs) == 0 {
		return fmt.Errorf("userUuids cannot be empty")
	}
	for _, u := range r.UserUUIDs {
		if _, err := uuid.Parse(u); err != nil {
			return fmt.Errorf("invalid user UUID format: %s", u)
		}
	}
	return nil
}

func ExternalSquadsHandler(manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handleGetExternalSquads(w, r, manager, cfg)
		case http.MethodPost:
			handleCreateExternalSquad(w, r, manager, cfg)
		case http.MethodPatch:
			handleUpdateExternalSquad(w, r, manager, cfg)
		default:
			shared.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
		}
	}
}

func ExternalSquadByUUIDHandler(manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimSpace(trimExternalSquadsPath(r.URL.Path))

		if path == "" {
			switch r.Method {
			case http.MethodGet:
				handleGetExternalSquads(w, r, manager, cfg)
			case http.MethodPost:
				handleCreateExternalSquad(w, r, manager, cfg)
			case http.MethodPatch:
				handleUpdateExternalSquad(w, r, manager, cfg)
			default:
				shared.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
			}
			return
		}

		parts := strings.Split(path, "/")
		squadUUID := strings.TrimSpace(parts[0])

		if _, err := uuid.Parse(squadUUID); err != nil {
			if r.Method == http.MethodPatch {
				handleUpdateExternalSquad(w, r, manager, cfg)
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
				handleBulkAddUsersToExternalSquad(w, r, manager, cfg, squadUUID)
				return
			}
			if len(parts) == 3 && parts[1] == "bulk-actions" && parts[2] == "remove-users" {
				if r.Method != http.MethodDelete && r.Method != http.MethodPost {
					shared.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
					return
				}
				handleBulkRemoveUsersFromExternalSquad(w, r, manager, cfg, squadUUID)
				return
			}
			http.NotFound(w, r)
			return
		}

		switch r.Method {
		case http.MethodGet:
			handleGetExternalSquadByUUID(w, r, manager, cfg, squadUUID)
		case http.MethodPatch:
			handleUpdateExternalSquad(w, r, manager, cfg)
		case http.MethodDelete:
			handleDeleteExternalSquad(w, r, manager, cfg, squadUUID)
		default:
			shared.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
		}
	}
}

func trimExternalSquadsPath(path string) string {
	for _, prefix := range []string{"/api/external-squads/"} {
		if strings.HasPrefix(path, prefix) {
			return strings.Trim(strings.TrimPrefix(path, prefix), "/")
		}
	}
	return ""
}

func ExternalSquadsReorderHandler(manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) http.HandlerFunc {
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

		err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
			tx, err := db.BeginTx(r.Context(), nil)
			if err != nil {
				return err
			}

			for _, item := range req.Squads {
				if _, err := tx.ExecContext(r.Context(),
					`UPDATE external_squads SET view_position = ? WHERE uuid = ?`,
					item.ViewPosition, item.UUID); err != nil {
					_ = tx.Rollback()
					return err
				}
			}

			return tx.Commit()
		})

		if err != nil {
			shared.SendError(w, http.StatusInternalServerError, "failed to reorder external squads", err, cfg)
			return
		}

		handleGetExternalSquads(w, r, manager, cfg)
	}
}

func handleGetExternalSquads(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) {
	records, err := getExternalSquads(r.Context(), manager)
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

		_ = manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
			_ = db.QueryRowContext(r.Context(), `SELECT COUNT(*) FROM users WHERE external_squad_uuid = ?`, api.UUID).Scan(&api.Info.MembersCount)
			rows, err := db.QueryContext(r.Context(), `SELECT template_uuid, template_type FROM external_squads_templates WHERE external_squad_uuid = ?`, api.UUID)
			if err == nil {
				defer rows.Close()
				for rows.Next() {
					var t ExternalSquadTemplate
					if err := rows.Scan(&t.TemplateUUID, &t.TemplateType); err == nil {
						api.Templates = append(api.Templates, t)
					}
				}
			}
			return nil
		})

		result = append(result, api)
	}

	shared.WriteJSON(w, http.StatusOK, map[string]any{
		"response": map[string]any{
			"total":          len(result),
			"externalSquads": result,
		},
	})
}

func handleGetExternalSquadByUUID(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig, squadUUID string) {
	record, err := getExternalSquadByUUID(r.Context(), manager, squadUUID)
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

	_ = manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		_ = db.QueryRowContext(r.Context(), `SELECT COUNT(*) FROM users WHERE external_squad_uuid = ?`, api.UUID).Scan(&api.Info.MembersCount)
		rows, err := db.QueryContext(r.Context(), `SELECT template_uuid, template_type FROM external_squads_templates WHERE external_squad_uuid = ?`, api.UUID)
		if err == nil {
			defer rows.Close()
			for rows.Next() {
				var t ExternalSquadTemplate
				if err := rows.Scan(&t.TemplateUUID, &t.TemplateType); err == nil {
					api.Templates = append(api.Templates, t)
				}
			}
		}
		return nil
	})

	shared.WriteJSON(w, http.StatusOK, map[string]any{"response": api})
}

func handleCreateExternalSquad(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) {
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

	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		tx, err := db.BeginTx(r.Context(), nil)
		if err != nil {
			return err
		}

		subSettings, err := marshalJSON(req.SubscriptionSettings)
		if err != nil {
			_ = tx.Rollback()
			return err
		}
		hostOverrides, err := marshalJSON(req.HostOverrides)
		if err != nil {
			_ = tx.Rollback()
			return err
		}
		respHeaders, err := marshalJSON(req.ResponseHeaders)
		if err != nil {
			_ = tx.Rollback()
			return err
		}
		hwidSettings, err := marshalJSON(req.HWIDSettings)
		if err != nil {
			_ = tx.Rollback()
			return err
		}
		customRemarks, err := marshalJSON(req.CustomRemarks)
		if err != nil {
			_ = tx.Rollback()
			return err
		}

		_, err = tx.ExecContext(r.Context(), `
			INSERT INTO external_squads (
				uuid, view_position, name,
				subscription_settings, host_overrides, response_headers,
				hwid_settings, custom_remarks, subpage_config_uuid
			) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
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
			_ = tx.Rollback()
			if isUniqueViolation(err, "external_squads_name_key") {
				return errExternalSquadExists
			}
			return err
		}

		return tx.Commit()
	})

	if err != nil {
		if errors.Is(err, errExternalSquadExists) {
			shared.SendError(w, http.StatusConflict, "external squad name already exists", nil, cfg)
			return
		}
		shared.SendError(w, http.StatusInternalServerError, "failed to create external squad", err, cfg)
		return
	}

	handleGetExternalSquadByUUID(w, r, manager, cfg, squadUUID)
}

func handleUpdateExternalSquad(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) {
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		shared.SendError(w, http.StatusBadRequest, "failed to read body", err, cfg)
		return
	}

	cfg.Logger.Info("UpdateExternalSquad: received raw body", "body", string(bodyBytes))

	var req UpdateExternalSquadRequest
	if err := json.Unmarshal(bodyBytes, &req); err != nil {
		cfg.Logger.Error("UpdateExternalSquad: failed to unmarshal JSON", "error", err)
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

	cfg.Logger.Info("UpdateExternalSquad: parsed request struct",
		"uuid", req.UUID,
		"has_name", req.Name != nil,
		"has_view_position", req.ViewPosition != nil,
		"has_subscription_settings", req.SubscriptionSettings != nil,
		"has_host_overrides", req.HostOverrides != nil,
		"has_response_headers", req.ResponseHeaders != nil,
		"has_hwid_settings", req.HWIDSettings != nil,
		"has_custom_remarks", req.CustomRemarks != nil,
		"has_subpage_config", req.SubpageConfigUUID != nil,
		"has_templates", req.Templates != nil,
	)

	err = manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
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
				WHERE uuid = ?
			`, strings.Join(clauses, ", "))

			if _, err := tx.ExecContext(r.Context(), query, args...); err != nil {
				cfg.Logger.Error("UpdateExternalSquad: SQL execution failed", "query", query, "error", err)
				_ = tx.Rollback()
				if isUniqueViolation(err, "external_squads_name_key") {
					return errExternalSquadExists
				}
				return err
			}
		} else if req.Templates != nil {
			_, err := tx.ExecContext(r.Context(), `
				UPDATE external_squads 
				SET updated_at = CURRENT_TIMESTAMP 
				WHERE uuid = ?
			`, req.UUID)
			if err != nil {
				_ = tx.Rollback()
				return err
			}
		}

		if req.Templates != nil {
			_, err = tx.ExecContext(r.Context(), `DELETE FROM external_squads_templates WHERE external_squad_uuid = ?`, req.UUID)
			if err != nil {
				_ = tx.Rollback()
				return err
			}

			for _, t := range *req.Templates {
				_, err = tx.ExecContext(r.Context(), `
					INSERT INTO external_squads_templates (external_squad_uuid, template_uuid, template_type)
					VALUES (?, ?, ?)
				`, req.UUID, t.TemplateUUID, t.TemplateType)
				if err != nil {
					_ = tx.Rollback()
					return err
				}
			}
		}

		if len(clauses) == 0 && req.Templates == nil {
			_ = tx.Rollback()
			return fmt.Errorf("no fields to update")
		}

		return tx.Commit()
	})

	if err != nil {
		if errors.Is(err, errExternalSquadExists) {
			shared.SendError(w, http.StatusConflict, "external squad name already exists", nil, cfg)
			return
		}
		if err.Error() == "no fields to update" {
			shared.SendError(w, http.StatusBadRequest, "no fields to update", nil, cfg)
			return
		}
		shared.SendError(w, http.StatusInternalServerError, "failed to update external squad", err, cfg)
		return
	}

	handleGetExternalSquadByUUID(w, r, manager, cfg, req.UUID)
}

func handleDeleteExternalSquad(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig, squadUUID string) {
	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		result, err := db.ExecContext(r.Context(), `DELETE FROM external_squads WHERE uuid = ?`, squadUUID)
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
			shared.SendError(w, http.StatusNotFound, "external squad not found", nil, cfg)
			return
		}
		shared.SendError(w, http.StatusInternalServerError, "failed to delete external squad", err, cfg)
		return
	}

	shared.WriteJSON(w, http.StatusOK, map[string]any{"response": map[string]any{"isDeleted": true}})
}

func handleBulkAddUsersToExternalSquad(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig, squadUUID string) {
	var req BulkUsersRequest
	bodyBytes, err := io.ReadAll(r.Body)

	hasSpecificUsers := false
	if err == nil && len(bodyBytes) > 0 {
		// Пытаемся распарсить только если тело не пустое
		if jsonErr := json.Unmarshal(bodyBytes, &req); jsonErr == nil && len(req.UserUUIDs) > 0 {
			if validateErr := req.Validate(); validateErr == nil {
				hasSpecificUsers = true
			}
		}
	}

	var affected int64
	err = manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		var exists int
		if err := db.QueryRowContext(r.Context(), `SELECT 1 FROM external_squads WHERE uuid = ?`, squadUUID).Scan(&exists); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return errExternalSquadNotFound
			}
			return err
		}

		if hasSpecificUsers {
			// Логика из NEW: обновляем только переданных юзеров
			placeholders := make([]string, len(req.UserUUIDs))
			args := make([]any, 0, len(req.UserUUIDs)+1)
			args = append(args, squadUUID)

			for i, u := range req.UserUUIDs {
				placeholders[i] = "?"
				args = append(args, u)
			}

			query := fmt.Sprintf(`
				UPDATE users
				SET external_squad_uuid = ?, updated_at = CURRENT_TIMESTAMP
				WHERE uuid IN (%s)
			`, strings.Join(placeholders, ", "))

			result, err := db.ExecContext(r.Context(), query, args...)
			if err != nil {
				return err
			}
			affected, _ = result.RowsAffected()
		} else {
			// Логика из OLD: привязываем всех, у кого отличается
			result, err := db.ExecContext(r.Context(), `
				UPDATE users
				SET external_squad_uuid = ?::uuid, updated_at = CURRENT_TIMESTAMP
				WHERE external_squad_uuid IS DISTINCT FROM ?::uuid
			`, squadUUID, squadUUID)
			if err != nil {
				return err
			}
			affected, _ = result.RowsAffected()
		}
		return nil
	})
	if err != nil {
		if errors.Is(err, errExternalSquadNotFound) {
			shared.SendError(w, http.StatusNotFound, "external squad not found", nil, cfg)
			return
		}
		shared.SendError(w, http.StatusInternalServerError, "failed to add users to external squad", err, cfg)
		return
	}

	cfg.Logger.Info("Users added to external squad", "squad_uuid", squadUUID, "affected_rows", affected)
	shared.WriteJSON(w, http.StatusOK, map[string]any{"response": map[string]any{"eventSent": true}})
}

func handleBulkRemoveUsersFromExternalSquad(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig, squadUUID string) {
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
	err = manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		var exists int
		if err := db.QueryRowContext(r.Context(), `SELECT 1 FROM external_squads WHERE uuid = ?`, squadUUID).Scan(&exists); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return errExternalSquadNotFound
			}
			return err
		}

		if hasSpecificUsers {
			// Логика из NEW: отвязываем только указанных юзеров
			placeholders := make([]string, len(req.UserUUIDs))
			args := make([]any, 0, len(req.UserUUIDs)+1)

			for i, u := range req.UserUUIDs {
				placeholders[i] = "?"
				args = append(args, u)
			}
			args = append(args, squadUUID)

			query := fmt.Sprintf(`
				UPDATE users
				SET external_squad_uuid = NULL, updated_at = CURRENT_TIMESTAMP
				WHERE uuid IN (%s) AND external_squad_uuid = ?
			`, strings.Join(placeholders, ", "))

			result, err := db.ExecContext(r.Context(), query, args...)
			if err != nil {
				return err
			}
			affected, _ = result.RowsAffected()
		} else {
			// Логика из OLD: отвязываем всех юзеров именно от этого сквада
			result, err := db.ExecContext(r.Context(), `
				UPDATE users
				SET external_squad_uuid = NULL, updated_at = CURRENT_TIMESTAMP
				WHERE external_squad_uuid = ?
			`, squadUUID)
			if err != nil {
				return err
			}
			affected, _ = result.RowsAffected()
		}
		return nil
	})
	if err != nil {
		if errors.Is(err, errExternalSquadNotFound) {
			shared.SendError(w, http.StatusNotFound, "external squad not found", nil, cfg)
			return
		}
		shared.SendError(w, http.StatusInternalServerError, "failed to remove users from external squad", err, cfg)
		return
	}

	cfg.Logger.Info("Users removed from external squad", "squad_uuid", squadUUID, "affected_rows", affected)

	shared.WriteJSON(w, http.StatusOK, map[string]any{
		"response": map[string]any{
			"eventSent": true,
		},
	})
}

func getExternalSquads(ctx context.Context, manager *dbmanager.DatabaseManager) ([]ExternalSquadRecord, error) {
	records := make([]ExternalSquadRecord, 0)

	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		rows, err := db.QueryContext(ctx, `
			SELECT uuid, view_position, name,
				subscription_settings, host_overrides, response_headers,
				hwid_settings, custom_remarks, subpage_config_uuid,
				created_at, updated_at
			FROM external_squads
			ORDER BY view_position ASC
		`)
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			rec, err := scanExternalSquad(rows)
			if err != nil {
				return err
			}
			records = append(records, rec)
		}

		return rows.Err()
	})

	return records, err
}

func getExternalSquadByUUID(ctx context.Context, manager *dbmanager.DatabaseManager, squadUUID string) (ExternalSquadRecord, error) {
	var record ExternalSquadRecord

	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		row := db.QueryRowContext(ctx, `
			SELECT uuid, view_position, name,
				subscription_settings, host_overrides, response_headers,
				hwid_settings, custom_remarks, subpage_config_uuid,
				created_at, updated_at
			FROM external_squads
			WHERE uuid = ?
		`, squadUUID)

		rec, err := scanExternalSquad(row)
		if err != nil {
			return err
		}
		record = rec
		return nil
	})

	return record, err
}

func scanExternalSquad(scanner shared.RowScanner) (ExternalSquadRecord, error) {
	var rec ExternalSquadRecord
	var subSettings, hostOverrides, respHeaders, hwidSettings, customRemarks sql.NullString
	var subpageConfigUUID sql.NullString

	err := scanner.Scan(
		&rec.UUID,
		&rec.ViewPosition,
		&rec.Name,
		&subSettings,
		&hostOverrides,
		&respHeaders,
		&hwidSettings,
		&customRemarks,
		&subpageConfigUUID,
		&rec.CreatedAt,
		&rec.UpdatedAt,
	)
	if err != nil {
		return rec, err
	}

	rec.SubscriptionSettings = parseJSONRaw(subSettings)
	rec.HostOverrides = parseJSONRaw(hostOverrides)
	rec.ResponseHeaders = parseJSONRaw(respHeaders)
	rec.HWIDSettings = parseJSONRaw(hwidSettings)
	rec.CustomRemarks = parseJSONRaw(customRemarks)
	if subpageConfigUUID.Valid {
		rec.SubpageConfigUUID = &subpageConfigUUID.String
	}

	return rec, nil
}

func convertExternalSquadToAPI(rec ExternalSquadRecord) (ExternalSquadAPI, error) {
	api := ExternalSquadAPI{
		UUID:              rec.UUID,
		ViewPosition:      rec.ViewPosition,
		Name:              rec.Name,
		Info:              ExternalSquadInfo{MembersCount: 0},
		Templates:         make([]ExternalSquadTemplate, 0),
		CreatedAt:         rec.CreatedAt,
		UpdatedAt:         rec.UpdatedAt,
		SubpageConfigUUID: rec.SubpageConfigUUID,
	}

	var err error
	if len(rec.SubscriptionSettings) > 0 {
		api.SubscriptionSettings, err = parseJSONMap(string(rec.SubscriptionSettings))
		if err != nil {
			return api, err
		}
	}
	if len(rec.HostOverrides) > 0 {
		api.HostOverrides, err = parseJSONMap(string(rec.HostOverrides))
		if err != nil {
			return api, err
		}
	}
	if len(rec.ResponseHeaders) > 0 {
		api.ResponseHeaders, err = parseJSONHeaders(string(rec.ResponseHeaders))
		if err != nil {
			return api, err
		}
	}
	if len(rec.HWIDSettings) > 0 {
		api.HWIDSettings, err = parseJSONMap(string(rec.HWIDSettings))
		if err != nil {
			return api, err
		}
	}
	if len(rec.CustomRemarks) > 0 {
		api.CustomRemarks, err = parseJSONMap(string(rec.CustomRemarks))
		if err != nil {
			return api, err
		}
	}

	return api, nil
}

func parseJSONRaw(raw sql.NullString) json.RawMessage {
	if !raw.Valid || raw.String == "" {
		return nil
	}
	return json.RawMessage(raw.String)
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

func marshalJSON(v any) (sql.NullString, error) {
	if v == nil {
		return sql.NullString{Valid: false}, nil
	}
	data, err := json.Marshal(v)
	if err != nil {
		return sql.NullString{}, err
	}
	return sql.NullString{String: string(data), Valid: true}, nil
}

func coalesceInt(v *int, fallback int) int {
	if v == nil {
		return fallback
	}
	return *v
}

func normalizeStringPtr(v *string) interface{} {
	if v == nil {
		return nil
	}
	if strings.TrimSpace(*v) == "" {
		return nil
	}
	return strings.TrimSpace(*v)
}

func isUniqueViolation(err error, constraint string) bool {
	return strings.Contains(err.Error(), "duplicate key") ||
		strings.Contains(err.Error(), "Unique constraint") ||
		strings.Contains(err.Error(), constraint)
}
