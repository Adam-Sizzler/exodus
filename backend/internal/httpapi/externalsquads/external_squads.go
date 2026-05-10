package externalsquads

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
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
	UUID                 string            `json:"uuid"`
	Name                 *string           `json:"name,omitempty"`
	ViewPosition         *int              `json:"viewPosition,omitempty"`
	SubscriptionSettings map[string]any    `json:"subscriptionSettings,omitempty"`
	HostOverrides        map[string]any    `json:"hostOverrides,omitempty"`
	ResponseHeaders      map[string]string `json:"responseHeaders,omitempty"`
	HWIDSettings         map[string]any    `json:"hwidSettings,omitempty"`
	CustomRemarks        map[string]any    `json:"customRemarks,omitempty"`
	SubpageConfigUUID    *string           `json:"subpageConfigUuid,omitempty"`
}

type ReorderExternalSquadsRequest struct {
	Items []struct {
		UUID         string `json:"uuid"`
		ViewPosition int    `json:"viewPosition"`
	} `json:"items"`

	// Backward compatibility with legacy payload.
	Squads []struct {
		UUID         string `json:"uuid"`
		ViewPosition int    `json:"viewPosition"`
	} `json:"squads"`
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
				if r.Method != http.MethodDelete {
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

			if _, err := tx.ExecContext(r.Context(),
				`SELECT setval('external_squads_view_position_seq', (SELECT COALESCE(MAX(view_position), 0) FROM external_squads) + 1)`); err != nil {
				_ = tx.Rollback()
				return err
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

	created, err := getExternalSquadByUUID(r.Context(), manager, squadUUID)
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to fetch created external squad", err, cfg)
		return
	}

	api, err := convertExternalSquadToAPI(created)
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to convert external squad", err, cfg)
		return
	}

	shared.WriteJSON(w, http.StatusCreated, map[string]any{"response": api})
}

func handleUpdateExternalSquad(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) {
	var req UpdateExternalSquadRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
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

		if req.Name != nil {
			add("name", strings.TrimSpace(*req.Name))
		}
		if req.ViewPosition != nil {
			add("view_position", *req.ViewPosition)
		}
		if req.SubscriptionSettings != nil {
			val, err := marshalJSON(req.SubscriptionSettings)
			if err != nil {
				_ = tx.Rollback()
				return err
			}
			add("subscription_settings", val)
		}
		if req.HostOverrides != nil {
			val, err := marshalJSON(req.HostOverrides)
			if err != nil {
				_ = tx.Rollback()
				return err
			}
			add("host_overrides", val)
		}
		if req.ResponseHeaders != nil {
			val, err := marshalJSON(req.ResponseHeaders)
			if err != nil {
				_ = tx.Rollback()
				return err
			}
			add("response_headers", val)
		}
		if req.HWIDSettings != nil {
			val, err := marshalJSON(req.HWIDSettings)
			if err != nil {
				_ = tx.Rollback()
				return err
			}
			add("hwid_settings", val)
		}
		if req.CustomRemarks != nil {
			val, err := marshalJSON(req.CustomRemarks)
			if err != nil {
				_ = tx.Rollback()
				return err
			}
			add("custom_remarks", val)
		}
		if req.SubpageConfigUUID != nil {
			add("subpage_config_uuid", normalizeStringPtr(req.SubpageConfigUUID))
		}

		if len(clauses) == 0 {
			_ = tx.Rollback()
			return fmt.Errorf("no fields to update")
		}

		clauses = append(clauses, "updated_at = CURRENT_TIMESTAMP")
		args = append(args, req.UUID)

		query := fmt.Sprintf(`
			UPDATE external_squads
			SET %s
			WHERE uuid = ?
			RETURNING uuid, view_position, name,
				subscription_settings, host_overrides, response_headers,
				hwid_settings, custom_remarks, subpage_config_uuid,
				created_at, updated_at
		`, strings.Join(clauses, ", "))

		row := tx.QueryRowContext(r.Context(), query, args...)
		var rec ExternalSquadRecord
		var subSettings, hostOverrides, respHeaders, hwidSettings, customRemarks sql.NullString
		var subpageConfigUUID sql.NullString

		scanErr := row.Scan(
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
		if scanErr != nil {
			_ = tx.Rollback()
			if errors.Is(scanErr, sql.ErrNoRows) {
				return errExternalSquadNotFound
			}
			return scanErr
		}

		rec.SubscriptionSettings = parseJSONRaw(subSettings)
		rec.HostOverrides = parseJSONRaw(hostOverrides)
		rec.ResponseHeaders = parseJSONRaw(respHeaders)
		rec.HWIDSettings = parseJSONRaw(hwidSettings)
		rec.CustomRemarks = parseJSONRaw(customRemarks)
		if subpageConfigUUID.Valid {
			rec.SubpageConfigUUID = &subpageConfigUUID.String
		}

		return tx.Commit()
	})

	if err != nil {
		if errors.Is(err, errExternalSquadNotFound) {
			shared.SendError(w, http.StatusNotFound, "external squad not found", nil, cfg)
			return
		}
		if errors.Is(err, errExternalSquadExists) {
			shared.SendError(w, http.StatusConflict, "external squad name already exists", nil, cfg)
			return
		}
		shared.SendError(w, http.StatusInternalServerError, "failed to update external squad", err, cfg)
		return
	}

	updated, getErr := getExternalSquadByUUID(r.Context(), manager, req.UUID)
	if getErr != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to fetch updated external squad", getErr, cfg)
		return
	}

	api, err := convertExternalSquadToAPI(updated)
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to convert external squad", err, cfg)
		return
	}

	shared.WriteJSON(w, http.StatusOK, map[string]any{"response": api})
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
	var affected int64
	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		var exists int
		if err := db.QueryRowContext(r.Context(), `SELECT 1 FROM external_squads WHERE uuid = ?`, squadUUID).Scan(&exists); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return errExternalSquadNotFound
			}
			return err
		}
		result, err := db.ExecContext(r.Context(), `
			UPDATE users
			SET external_squad_uuid = ?::uuid, updated_at = CURRENT_TIMESTAMP
			WHERE external_squad_uuid IS DISTINCT FROM ?::uuid
		`, squadUUID, squadUUID)
		if err != nil {
			return err
		}
		affected, _ = result.RowsAffected()
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
	var affected int64
	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		var exists int
		if err := db.QueryRowContext(r.Context(), `SELECT 1 FROM external_squads WHERE uuid = ?`, squadUUID).Scan(&exists); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return errExternalSquadNotFound
			}
			return err
		}
		result, err := db.ExecContext(r.Context(), `
			UPDATE users
			SET external_squad_uuid = NULL, updated_at = CURRENT_TIMESTAMP
			WHERE external_squad_uuid = ?
		`, squadUUID)
		if err != nil {
			return err
		}
		affected, _ = result.RowsAffected()
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
	shared.WriteJSON(w, http.StatusOK, map[string]any{"response": map[string]any{"eventSent": true}})
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
	// Check for unique constraint violation
	// This is a simplified check, may need to be enhanced based on DB driver
	return strings.Contains(err.Error(), "duplicate key") ||
		strings.Contains(err.Error(), "Unique constraint") ||
		strings.Contains(err.Error(), constraint)
}
