package subscriptiontemplate

import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"exodus/internal/httpapi/shared"

	"github.com/jackc/pgx/v5/pgconn"
)

const (
	subscriptionTemplatesBasePath    = "/api/subscription-templates"
	subscriptionTemplatesActionsPath = "/api/subscription-templates/actions/"
)

var subscriptionTemplateNameRegex = regexp.MustCompile(`^[A-Za-z0-9_\s-]+$`)

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
