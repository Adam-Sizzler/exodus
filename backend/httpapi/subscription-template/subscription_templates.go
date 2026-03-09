package subscriptiontemplate

import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"v2ray-stat/backend/config"
	"v2ray-stat/backend/db"
	dbmanager "v2ray-stat/backend/db/manager"
	"v2ray-stat/backend/httpapi/shared"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
)

const (
	subscriptionTemplatesBasePath    = "/api/subscription-templates"
	subscriptionTemplatesActionsPath = "/api/subscription-templates/actions/"
)

var subscriptionTemplateNameRegex = regexp.MustCompile(`^[A-Za-z0-9_\s-]+$`)

// SubscriptionTemplate represents a subscription template response.
type SubscriptionTemplate struct {
	UUID               string          `json:"uuid"`
	ViewPosition       int             `json:"viewPosition"`
	Name               string          `json:"name"`
	TemplateType       string          `json:"templateType"`
	TemplateJSON       json.RawMessage `json:"templateJson"`
	EncodedTemplateYML *string         `json:"encodedTemplateYaml"`
}

type subscriptionTemplatesListResponse struct {
	Total     int                    `json:"total"`
	Templates []SubscriptionTemplate `json:"templates"`
}

type subscriptionTemplateCreateRequest struct {
	Name         string `json:"name"`
	TemplateType string `json:"templateType"`
}

type subscriptionTemplateUpdateRequest struct {
	UUID               string           `json:"uuid"`
	Name               *string          `json:"name,omitempty"`
	TemplateJSON       *json.RawMessage `json:"templateJson,omitempty"`
	EncodedTemplateYML *string          `json:"encodedTemplateYaml,omitempty"`
}

type subscriptionTemplateDeleteResponse struct {
	IsDeleted bool `json:"isDeleted"`
}

type subscriptionTemplateReorderItem struct {
	UUID         string `json:"uuid"`
	ViewPosition int    `json:"viewPosition"`
}

type subscriptionTemplateReorderRequest struct {
	Items []subscriptionTemplateReorderItem `json:"items"`
}

type subscriptionTemplateRecord struct {
	UUID         string
	ViewPosition int
	Name         string
	TemplateType string
	TemplateYAML *string
	TemplateJSON json.RawMessage
}

func SubscriptionTemplatesHandler(manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handleGetSubscriptionTemplates(w, r, manager, cfg)
		case http.MethodPost:
			handleCreateSubscriptionTemplate(w, r, manager, cfg)
		case http.MethodPatch:
			handleUpdateSubscriptionTemplate(w, r, manager, cfg)
		default:
			shared.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
		}
	}
}

func SubscriptionTemplatesActionsHandler(manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			shared.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}

		path := strings.TrimPrefix(r.URL.Path, subscriptionTemplatesActionsPath)
		path = strings.Trim(path, "/")
		switch path {
		case "reorder":
			handleReorderSubscriptionTemplates(w, r, manager, cfg)
		default:
			http.NotFound(w, r)
		}
	}
}

func SubscriptionTemplateByUUIDHandler(manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uuidStr := strings.TrimSpace(strings.TrimPrefix(r.URL.Path, subscriptionTemplatesBasePath+"/"))
		if uuidStr == "" {
			switch r.Method {
			case http.MethodGet:
				handleGetSubscriptionTemplates(w, r, manager, cfg)
			case http.MethodPost:
				handleCreateSubscriptionTemplate(w, r, manager, cfg)
			case http.MethodPatch:
				handleUpdateSubscriptionTemplate(w, r, manager, cfg)
			default:
				shared.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
			}
			return
		}
		if _, err := uuid.Parse(uuidStr); err != nil {
			shared.SendError(w, http.StatusBadRequest, "invalid UUID format", nil, cfg)
			return
		}

		switch r.Method {
		case http.MethodGet:
			handleGetSubscriptionTemplateByUUID(w, r, manager, cfg, uuidStr)
		case http.MethodDelete:
			handleDeleteSubscriptionTemplate(w, r, manager, cfg, uuidStr)
		default:
			shared.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
		}
	}
}

func handleGetSubscriptionTemplates(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) {
	templates := make([]SubscriptionTemplate, 0)

	err := manager.ExecuteHighPriority(func(dbConn dbmanager.DBExecutor) error {
		rows, err := dbConn.QueryContext(r.Context(), `
			SELECT uuid, view_position, name, template_type
			FROM subscription_templates
			ORDER BY view_position ASC, template_type ASC`)
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			var rec subscriptionTemplateRecord
			if scanErr := rows.Scan(&rec.UUID, &rec.ViewPosition, &rec.Name, &rec.TemplateType); scanErr != nil {
				return scanErr
			}
			templates = append(templates, SubscriptionTemplate{
				UUID:               rec.UUID,
				ViewPosition:       rec.ViewPosition,
				Name:               rec.Name,
				TemplateType:       rec.TemplateType,
				TemplateJSON:       nil,
				EncodedTemplateYML: nil,
			})
		}
		return rows.Err()
	})
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to fetch templates", err, cfg)
		return
	}

	shared.WriteJSON(w, http.StatusOK, map[string]any{
		"response": subscriptionTemplatesListResponse{
			Total:     len(templates),
			Templates: templates,
		},
	})
}

func handleGetSubscriptionTemplateByUUID(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig, templateUUID string) {
	var rec subscriptionTemplateRecord

	err := manager.ExecuteHighPriority(func(dbConn dbmanager.DBExecutor) error {
		row := dbConn.QueryRowContext(r.Context(), `
			SELECT uuid, view_position, name, template_type, template_yaml, template_json
			FROM subscription_templates
			WHERE uuid = ?`, templateUUID)
		return scanSubscriptionTemplateRecord(row, &rec)
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			shared.SendError(w, http.StatusNotFound, "template not found", nil, cfg)
			return
		}
		shared.SendError(w, http.StatusInternalServerError, "failed to fetch template", err, cfg)
		return
	}

	resp := mapSubscriptionTemplateRecord(rec, true)
	shared.WriteJSON(w, http.StatusOK, map[string]any{"response": resp})
}

func handleCreateSubscriptionTemplate(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) {
	var req subscriptionTemplateCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		shared.SendError(w, http.StatusBadRequest, "invalid JSON", err, cfg)
		return
	}

	req.Name = strings.TrimSpace(req.Name)
	req.TemplateType = strings.TrimSpace(strings.ToUpper(req.TemplateType))

	if err := validateTemplateName(req.Name); err != nil {
		shared.SendError(w, http.StatusBadRequest, err.Error(), nil, cfg)
		return
	}
	if req.Name == "Default" {
		shared.SendError(w, http.StatusBadRequest, "reserved template name", nil, cfg)
		return
	}
	if !isAllowedTemplateType(req.TemplateType) {
		shared.SendError(w, http.StatusBadRequest, "invalid templateType", nil, cfg)
		return
	}
	if req.TemplateType == "XRAY_BASE64" {
		shared.SendError(w, http.StatusBadRequest, "templateType not allowed", nil, cfg)
		return
	}

	defaultTmpl, ok := db.DefaultSubscriptionTemplateByType(req.TemplateType)
	if !ok {
		shared.SendError(w, http.StatusBadRequest, "unknown templateType", nil, cfg)
		return
	}

	newUUID := uuid.NewString()
	var created subscriptionTemplateRecord

	err := manager.ExecuteHighPriority(func(dbConn dbmanager.DBExecutor) error {
		viewPosition := 0
		row := dbConn.QueryRowContext(r.Context(), `SELECT COALESCE(MAX(view_position), 0) + 1 FROM subscription_templates`)
		if scanErr := row.Scan(&viewPosition); scanErr != nil {
			return scanErr
		}

		_, execErr := dbConn.ExecContext(r.Context(), `
			INSERT INTO subscription_templates (
				uuid, view_position, name, template_type, template_yaml, template_json
			) VALUES (?, ?, ?, ?, ?, ?)`,
			newUUID,
			viewPosition,
			req.Name,
			req.TemplateType,
			defaultTmpl.TemplateYAML,
			defaultTmpl.TemplateJSON,
		)
		if execErr != nil {
			return execErr
		}

		created = subscriptionTemplateRecord{
			UUID:         newUUID,
			ViewPosition: viewPosition,
			Name:         req.Name,
			TemplateType: req.TemplateType,
			TemplateYAML: defaultTmpl.TemplateYAML,
			TemplateJSON: defaultTmpl.TemplateJSON,
		}
		return nil
	})
	if err != nil {
		if isUniqueViolation(err) {
			shared.SendError(w, http.StatusConflict, "template name already exists for templateType", err, cfg)
			return
		}
		shared.SendError(w, http.StatusInternalServerError, "failed to create template", err, cfg)
		return
	}

	shared.WriteJSON(w, http.StatusCreated, map[string]any{"response": mapSubscriptionTemplateRecord(created, true)})
}

func handleUpdateSubscriptionTemplate(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) {
	var req subscriptionTemplateUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		shared.SendError(w, http.StatusBadRequest, "invalid JSON", err, cfg)
		return
	}
	if _, err := uuid.Parse(req.UUID); err != nil {
		shared.SendError(w, http.StatusBadRequest, "invalid UUID format", nil, cfg)
		return
	}

	if req.Name != nil {
		trimmed := strings.TrimSpace(*req.Name)
		if err := validateTemplateName(trimmed); err != nil {
			shared.SendError(w, http.StatusBadRequest, err.Error(), nil, cfg)
			return
		}
		if trimmed == "Default" {
			shared.SendError(w, http.StatusBadRequest, "reserved template name", nil, cfg)
			return
		}
		req.Name = &trimmed
	}

	if req.TemplateJSON != nil {
		if err := ensureJSONObject(*req.TemplateJSON); err != nil {
			shared.SendError(w, http.StatusBadRequest, err.Error(), nil, cfg)
			return
		}
	}

	if req.TemplateJSON != nil && req.EncodedTemplateYML != nil {
		shared.SendError(w, http.StatusBadRequest, "templateJson and encodedTemplateYaml cannot be updated together", nil, cfg)
		return
	}

	var template subscriptionTemplateRecord
	var decodedYAML *string

	err := manager.ExecuteHighPriority(func(dbConn dbmanager.DBExecutor) error {
		row := dbConn.QueryRowContext(r.Context(), `
			SELECT uuid, view_position, name, template_type, template_yaml, template_json
			FROM subscription_templates
			WHERE uuid = ?`, req.UUID)
		if scanErr := scanSubscriptionTemplateRecord(row, &template); scanErr != nil {
			return scanErr
		}

		isYamlTemplate := template.TemplateType == "MIHOMO" || template.TemplateType == "STASH" || template.TemplateType == "CLASH"
		isJsonTemplate := template.TemplateType == "XRAY_JSON" || template.TemplateType == "SINGBOX"

		if isYamlTemplate && req.TemplateJSON != nil {
			return fmt.Errorf("templateJson not allowed for YAML template")
		}
		if isJsonTemplate && req.EncodedTemplateYML != nil {
			return fmt.Errorf("encodedTemplateYaml not allowed for JSON template")
		}

		if req.EncodedTemplateYML != nil {
			decoded, decodeErr := base64.StdEncoding.DecodeString(*req.EncodedTemplateYML)
			if decodeErr != nil {
				return decodeErr
			}
			decodedStr := string(decoded)
			decodedYAML = &decodedStr
		}

		var clauses []string
		var args []any
		add := func(column string, value any) {
			clauses = append(clauses, fmt.Sprintf("%s = ?", column))
			args = append(args, value)
		}

		if req.Name != nil {
			add("name", *req.Name)
		}
		if req.TemplateJSON != nil {
			add("template_json", []byte(*req.TemplateJSON))
		}
		if req.EncodedTemplateYML != nil {
			add("template_yaml", decodedYAML)
		}

		if len(clauses) == 0 {
			return fmt.Errorf("no fields to update")
		}

		args = append(args, req.UUID)
		query := fmt.Sprintf("UPDATE subscription_templates SET %s, updated_at = CURRENT_TIMESTAMP WHERE uuid = ?", strings.Join(clauses, ", "))

		result, execErr := dbConn.ExecContext(r.Context(), query, args...)
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

		if req.Name != nil {
			template.Name = *req.Name
		}
		if req.TemplateJSON != nil {
			template.TemplateJSON = *req.TemplateJSON
		}
		if req.EncodedTemplateYML != nil {
			template.TemplateYAML = decodedYAML
		}

		return nil
	})
	if err != nil {
		var corrupt base64.CorruptInputError
		switch {
		case errors.Is(err, sql.ErrNoRows):
			shared.SendError(w, http.StatusNotFound, "template not found", nil, cfg)
		case isUniqueViolation(err):
			shared.SendError(w, http.StatusConflict, "template name already exists for templateType", err, cfg)
		case strings.Contains(err.Error(), "not allowed"):
			shared.SendError(w, http.StatusBadRequest, err.Error(), nil, cfg)
		case strings.Contains(err.Error(), "no fields to update"):
			shared.SendError(w, http.StatusBadRequest, "no fields to update", nil, cfg)
		case errors.As(err, &corrupt):
			shared.SendError(w, http.StatusBadRequest, "encodedTemplateYaml must be base64", err, cfg)
		default:
			shared.SendError(w, http.StatusInternalServerError, "update failed", err, cfg)
		}
		return
	}

	shared.WriteJSON(w, http.StatusOK, map[string]any{"response": mapSubscriptionTemplateRecord(template, true)})
}

func handleDeleteSubscriptionTemplate(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig, templateUUID string) {
	var templateName string

	err := manager.ExecuteHighPriority(func(dbConn dbmanager.DBExecutor) error {
		row := dbConn.QueryRowContext(r.Context(), `SELECT name FROM subscription_templates WHERE uuid = ?`, templateUUID)
		return row.Scan(&templateName)
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			shared.SendError(w, http.StatusNotFound, "template not found", nil, cfg)
			return
		}
		shared.SendError(w, http.StatusInternalServerError, "failed to find template", err, cfg)
		return
	}

	if templateName == "Default" {
		shared.SendError(w, http.StatusBadRequest, "reserved template cannot be deleted", nil, cfg)
		return
	}

	err = manager.ExecuteHighPriority(func(dbConn dbmanager.DBExecutor) error {
		_, execErr := dbConn.ExecContext(r.Context(), `DELETE FROM subscription_templates WHERE uuid = ?`, templateUUID)
		return execErr
	})
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to delete template", err, cfg)
		return
	}

	shared.WriteJSON(w, http.StatusOK, map[string]any{
		"response": subscriptionTemplateDeleteResponse{IsDeleted: true},
	})
}

func handleReorderSubscriptionTemplates(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) {
	var req subscriptionTemplateReorderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		shared.SendError(w, http.StatusBadRequest, "invalid JSON", err, cfg)
		return
	}
	if len(req.Items) == 0 {
		shared.SendError(w, http.StatusBadRequest, "items cannot be empty", nil, cfg)
		return
	}
	for _, item := range req.Items {
		if _, err := uuid.Parse(item.UUID); err != nil {
			shared.SendError(w, http.StatusBadRequest, "invalid UUID format", nil, cfg)
			return
		}
	}

	err := manager.ExecuteHighPriority(func(dbConn dbmanager.DBExecutor) error {
		tx, err := dbConn.BeginTx(r.Context(), nil)
		if err != nil {
			return err
		}
		defer func() {
			_ = tx.Rollback()
		}()

		for _, item := range req.Items {
			if _, execErr := tx.ExecContext(r.Context(), `UPDATE subscription_templates SET view_position = ? WHERE uuid = ?`, item.ViewPosition, item.UUID); execErr != nil {
				return execErr
			}
		}

		if err := tx.Commit(); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to reorder templates", err, cfg)
		return
	}

	// Return updated list without content
	handleGetSubscriptionTemplates(w, r, manager, cfg)
}

func scanSubscriptionTemplateRecord(scanner shared.RowScanner, dest *subscriptionTemplateRecord) error {
	var templateYAML sql.NullString
	var templateJSONBytes []byte
	if err := scanner.Scan(&dest.UUID, &dest.ViewPosition, &dest.Name, &dest.TemplateType, &templateYAML, &templateJSONBytes); err != nil {
		return err
	}

	if templateYAML.Valid {
		val := templateYAML.String
		dest.TemplateYAML = &val
	}

	if len(templateJSONBytes) > 0 {
		dest.TemplateJSON = json.RawMessage(templateJSONBytes)
	} else {
		dest.TemplateJSON = nil
	}

	return nil
}

func mapSubscriptionTemplateRecord(rec subscriptionTemplateRecord, includeContent bool) SubscriptionTemplate {
	var encodedYAML *string
	var templateJSON json.RawMessage
	if includeContent {
		templateJSON = rec.TemplateJSON
		if rec.TemplateYAML != nil {
			encoded := base64.StdEncoding.EncodeToString([]byte(*rec.TemplateYAML))
			encodedYAML = &encoded
		}
	}

	return SubscriptionTemplate{
		UUID:               rec.UUID,
		ViewPosition:       rec.ViewPosition,
		Name:               rec.Name,
		TemplateType:       rec.TemplateType,
		TemplateJSON:       templateJSON,
		EncodedTemplateYML: encodedYAML,
	}
}

func validateTemplateName(name string) error {
	if name == "" {
		return fmt.Errorf("name is required")
	}
	if len(name) < 2 {
		return fmt.Errorf("name must be at least 2 characters")
	}
	if len(name) > 255 {
		return fmt.Errorf("name must be less than 255 characters")
	}
	if !subscriptionTemplateNameRegex.MatchString(name) {
		return fmt.Errorf("name can only contain letters, numbers, underscores, dashes and spaces")
	}
	return nil
}

func ensureJSONObject(raw json.RawMessage) error {
	if len(raw) == 0 {
		return fmt.Errorf("templateJson must be a JSON object")
	}
	var value any
	if err := json.Unmarshal(raw, &value); err != nil {
		return fmt.Errorf("templateJson must be valid JSON")
	}
	if value == nil {
		return fmt.Errorf("templateJson must be a JSON object")
	}
	if _, ok := value.(map[string]any); !ok {
		return fmt.Errorf("templateJson must be a JSON object")
	}
	return nil
}

func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "23505"
	}
	return strings.Contains(err.Error(), "UNIQUE constraint failed")
}

func isAllowedTemplateType(v string) bool {
	switch v {
	case "XRAY_JSON", "XRAY_BASE64", "MIHOMO", "STASH", "CLASH", "SINGBOX":
		return true
	default:
		return false
	}
}
