package templates

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"v2ray-stat/backend/panel/config"
	dbmanager "v2ray-stat/backend/panel/db/manager"
	"v2ray-stat/backend/panel/httpapi/shared"

	"github.com/google/uuid"
)

type SubscriptionTemplate struct {
	UUID         string    `json:"uuid"`
	TemplateType string    `json:"template_type"`
	TemplateYAML *string   `json:"template_yaml,omitempty"`
	TemplateJSON *string   `json:"template_json,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Name         string    `json:"name"`
	ViewPosition int       `json:"view_position"`
}

type SubscriptionTemplateUpdateRequest struct {
	TemplateYAML *string `json:"template_yaml,omitempty"`
	TemplateJSON *string `json:"template_json,omitempty"`
	Name         *string `json:"name,omitempty"`
	ViewPosition *int    `json:"view_position,omitempty"`
}

type SubscriptionTemplateCreateRequest struct {
	TemplateType string  `json:"template_type"`
	TemplateYAML *string `json:"template_yaml,omitempty"`
	TemplateJSON *string `json:"template_json,omitempty"`
	Name         string  `json:"name"`
	ViewPosition *int    `json:"view_position,omitempty"`
}

func (r *SubscriptionTemplateUpdateRequest) Validate() error {
	if r.Name != nil && strings.TrimSpace(*r.Name) == "" {
		return fmt.Errorf("name cannot be empty")
	}
	if r.TemplateJSON != nil && strings.TrimSpace(*r.TemplateJSON) != "" {
		if !json.Valid([]byte(*r.TemplateJSON)) {
			return fmt.Errorf("template_json must be valid JSON")
		}
	}
	if r.ViewPosition != nil && *r.ViewPosition < 0 {
		return fmt.Errorf("view_position must be >= 0")
	}
	return nil
}

func (r *SubscriptionTemplateCreateRequest) Validate() error {
	r.TemplateType = strings.TrimSpace(strings.ToUpper(r.TemplateType))
	r.Name = strings.TrimSpace(r.Name)

	if r.TemplateType == "" {
		return fmt.Errorf("template_type is required")
	}
	if !isAllowedTemplateType(r.TemplateType) {
		return fmt.Errorf("template_type must be one of: XRAY_JSON, MIHOMO, STASH, CLASH, SINGBOX")
	}
	if r.Name == "" {
		return fmt.Errorf("name is required")
	}
	if r.TemplateJSON != nil && strings.TrimSpace(*r.TemplateJSON) != "" {
		if !json.Valid([]byte(*r.TemplateJSON)) {
			return fmt.Errorf("template_json must be valid JSON")
		}
	}
	if r.ViewPosition != nil && *r.ViewPosition < 0 {
		return fmt.Errorf("view_position must be >= 0")
	}
	return nil
}

func (r *SubscriptionTemplateUpdateRequest) HasUpdates() bool {
	return r.TemplateYAML != nil || r.TemplateJSON != nil || r.Name != nil || r.ViewPosition != nil
}

func scanSubscriptionTemplate(scanner shared.RowScanner) (SubscriptionTemplate, error) {
	var t SubscriptionTemplate
	var templateYAML, templateJSON sql.NullString
	var viewPosition sql.NullInt64

	err := scanner.Scan(
		&t.UUID,
		&t.TemplateType,
		&templateYAML,
		&templateJSON,
		&t.CreatedAt,
		&t.UpdatedAt,
		&t.Name,
		&viewPosition,
	)
	if err != nil {
		return t, err
	}

	if templateYAML.Valid {
		t.TemplateYAML = &templateYAML.String
	}
	if templateJSON.Valid {
		t.TemplateJSON = &templateJSON.String
	}
	if viewPosition.Valid {
		t.ViewPosition = int(viewPosition.Int64)
	}

	return t, nil
}

func SubscriptionTemplatesHandler(manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handleGetSubscriptionTemplates(w, r, manager, cfg)
		case http.MethodPost:
			handleCreateSubscriptionTemplate(w, r, manager, cfg)
		default:
			http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		}
	}
}

func SubscriptionTemplatesReorderHandler(manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
			return
		}

		var req shared.ViewPositionReorderRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			shared.SendError(w, http.StatusBadRequest, "invalid JSON", err, cfg)
			return
		}
		if err := req.Validate(); err != nil {
			shared.SendError(w, http.StatusBadRequest, err.Error(), nil, cfg)
			return
		}

		err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
			return shared.ApplyViewPositionReorder(r.Context(), db, "subscription_templates", req.OrderedUUIDs, cfg)
		})
		if err != nil {
			shared.SendError(w, http.StatusInternalServerError, "failed to reorder templates", err, cfg)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "templates reordered",
			"count":   len(req.OrderedUUIDs),
		})
	}
}

func SubscriptionTemplateByUUIDHandler(manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		templateUUID := strings.TrimSpace(strings.TrimPrefix(r.URL.Path, "/api/v1/templates/"))
		if _, err := uuid.Parse(templateUUID); err != nil {
			shared.SendError(w, http.StatusBadRequest, "invalid UUID format", nil, cfg)
			return
		}

		switch r.Method {
		case http.MethodGet:
			handleGetSubscriptionTemplateByUUID(w, r, manager, cfg, templateUUID)
		case http.MethodPatch:
			handlePatchSubscriptionTemplate(w, r, manager, cfg, templateUUID)
		case http.MethodDelete:
			handleDeleteSubscriptionTemplate(w, r, manager, cfg, templateUUID)
		default:
			http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		}
	}
}

func handleGetSubscriptionTemplates(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) {
	templateType := strings.TrimSpace(r.URL.Query().Get("template_type"))
	var templates []SubscriptionTemplate

	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		baseQuery := `
			SELECT
				uuid, template_type, template_yaml, template_json,
				created_at, updated_at, name, view_position
			FROM subscription_templates`

		var rows *sql.Rows
		var err error
		if templateType != "" {
			rows, err = db.QueryContext(r.Context(), baseQuery+` WHERE template_type = ? ORDER BY view_position ASC, template_type ASC`, templateType)
		} else {
			rows, err = db.QueryContext(r.Context(), baseQuery+` ORDER BY view_position ASC, template_type ASC`)
		}
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			t, scanErr := scanSubscriptionTemplate(rows)
			if scanErr != nil {
				return scanErr
			}
			templates = append(templates, t)
		}
		return rows.Err()
	})
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to fetch templates", err, cfg)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"templates": templates,
		"count":     len(templates),
	})
}

func handleGetSubscriptionTemplateByUUID(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig, templateUUID string) {
	var template SubscriptionTemplate

	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		row := db.QueryRowContext(r.Context(), `
			SELECT
				uuid, template_type, template_yaml, template_json,
				created_at, updated_at, name, view_position
			FROM subscription_templates
			WHERE uuid = ?`, templateUUID)
		var scanErr error
		template, scanErr = scanSubscriptionTemplate(row)
		return scanErr
	})
	if err != nil {
		if err == sql.ErrNoRows {
			shared.SendError(w, http.StatusNotFound, "template not found", nil, cfg)
		} else {
			shared.SendError(w, http.StatusInternalServerError, "failed to fetch template", err, cfg)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{"template": template})
}

func handlePatchSubscriptionTemplate(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig, templateUUID string) {
	var req SubscriptionTemplateUpdateRequest
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

	if req.TemplateYAML != nil {
		add("template_yaml", *req.TemplateYAML)
	}
	if req.TemplateJSON != nil {
		add("template_json", *req.TemplateJSON)
	}
	if req.Name != nil {
		add("name", strings.TrimSpace(*req.Name))
	}
	if req.ViewPosition != nil {
		add("view_position", *req.ViewPosition)
	}

	args = append(args, templateUUID)
	query := fmt.Sprintf("UPDATE subscription_templates SET %s, updated_at = CURRENT_TIMESTAMP WHERE uuid = ?", strings.Join(clauses, ", "))

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
			shared.SendError(w, http.StatusNotFound, "template not found", nil, cfg)
		} else if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			shared.SendError(w, http.StatusConflict, "name already exists", err, cfg)
		} else {
			shared.SendError(w, http.StatusInternalServerError, "update failed", err, cfg)
		}
		return
	}

	handleGetSubscriptionTemplateByUUID(w, r, manager, cfg, templateUUID)
}

func handleCreateSubscriptionTemplate(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) {
	var req SubscriptionTemplateCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		shared.SendError(w, http.StatusBadRequest, "invalid JSON", err, cfg)
		return
	}
	if err := req.Validate(); err != nil {
		shared.SendError(w, http.StatusBadRequest, err.Error(), nil, cfg)
		return
	}

	templateUUID := uuid.NewString()
	viewPosition := 0
	if req.ViewPosition != nil {
		viewPosition = *req.ViewPosition
	}

	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		if req.ViewPosition == nil {
			row := db.QueryRowContext(r.Context(), `SELECT COALESCE(MAX(view_position), -1) + 1 FROM subscription_templates WHERE template_type = ?`, req.TemplateType)
			if scanErr := row.Scan(&viewPosition); scanErr != nil {
				return scanErr
			}
		}

		_, execErr := db.ExecContext(r.Context(), `
			INSERT INTO subscription_templates (
				uuid, view_position, name, template_type, template_yaml, template_json
			) VALUES (?, ?, ?, ?, ?, ?)`,
			templateUUID,
			viewPosition,
			req.Name,
			req.TemplateType,
			req.TemplateYAML,
			req.TemplateJSON,
		)
		return execErr
	})
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			shared.SendError(w, http.StatusConflict, "template with this name already exists for template_type", err, cfg)
		} else {
			shared.SendError(w, http.StatusInternalServerError, "failed to create template", err, cfg)
		}
		return
	}

	w.WriteHeader(http.StatusCreated)
	handleGetSubscriptionTemplateByUUID(w, r, manager, cfg, templateUUID)
}

func handleDeleteSubscriptionTemplate(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig, templateUUID string) {
	var templateType, templateName string
	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		return db.QueryRowContext(r.Context(), `SELECT template_type, name FROM subscription_templates WHERE uuid = ?`, templateUUID).Scan(&templateType, &templateName)
	})
	if err != nil {
		if err == sql.ErrNoRows {
			shared.SendError(w, http.StatusNotFound, "template not found", nil, cfg)
		} else {
			shared.SendError(w, http.StatusInternalServerError, "failed to find template", err, cfg)
		}
		return
	}

	err = manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		_, execErr := db.ExecContext(r.Context(), `DELETE FROM subscription_templates WHERE uuid = ?`, templateUUID)
		return execErr
	})
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to delete template", err, cfg)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"message":       "template deleted",
		"uuid":          templateUUID,
		"template_type": templateType,
		"name":          templateName,
	})
}

func isAllowedTemplateType(v string) bool {
	switch v {
	case "XRAY_JSON", "MIHOMO", "STASH", "CLASH", "SINGBOX":
		return true
	default:
		return false
	}
}
