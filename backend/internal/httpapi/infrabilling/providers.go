package infrabilling

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strings"
	"time"

	"exodus/internal/config"
	dbmanager "exodus/internal/db/manager"
	"exodus/internal/httpapi/shared"
)

type providerNode struct {
	Name    string               `json:"name"`
	Details *providerNodeDetails `json:"details"`
}

type providerNodeDetails struct {
	NodeUUID    string `json:"nodeUuid"`
	CountryCode string `json:"countryCode"`
}

type providerRecord struct {
	UUID      string         `json:"uuid"`
	Name      string         `json:"name"`
	Favicon   *string        `json:"faviconLink"`
	LoginURL  *string        `json:"loginUrl"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	History   map[string]any `json:"billingHistory"`
	Nodes     []providerNode `json:"billingNodes"`
}

type createProviderRequest struct {
	Name       string  `json:"name"`
	FaviconURL *string `json:"faviconLink"`
	LoginURL   *string `json:"loginUrl"`
}

type updateProviderRequest struct {
	UUID       string         `json:"uuid"`
	Name       optionalString `json:"name"`
	FaviconURL optionalString `json:"faviconLink"`
	LoginURL   optionalString `json:"loginUrl"`
}

type optionalString struct {
	Set   bool
	Value *string
}

func (field *optionalString) UnmarshalJSON(data []byte) error {
	field.Set = true
	if string(data) == "null" {
		field.Value = nil
		return nil
	}

	var value string
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	field.Value = &value
	return nil
}

func ProvidersHandler(manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		providerUUID := uuidFromPath(r.URL.Path, "/api/infra-billing/providers")

		switch r.Method {
		case http.MethodGet:
			if providerUUID != "" {
				writeProviderResponse(w, r, manager, cfg, providerUUID)
				return
			}
			writeProvidersResponse(w, r, manager, cfg)
		case http.MethodPost:
			handleCreateProvider(w, r, manager, cfg)
		case http.MethodPatch:
			handleUpdateProvider(w, r, manager, cfg)
		case http.MethodDelete:
			handleDeleteProvider(w, r, manager, cfg, providerUUID)
		default:
			shared.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
		}
	}
}

func handleCreateProvider(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) {
	var req createProviderRequest
	if err := decodeJSONBody(r, &req); err != nil {
		shared.WriteJSONError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	name, err := validateProviderName(req.Name)
	if err != nil {
		shared.WriteJSONError(w, http.StatusBadRequest, err.Error())
		return
	}
	faviconURL, err := normalizeOptionalURL(req.FaviconURL)
	if err != nil {
		shared.WriteJSONError(w, http.StatusBadRequest, "invalid faviconLink")
		return
	}
	loginURL, err := normalizeOptionalURL(req.LoginURL)
	if err != nil {
		shared.WriteJSONError(w, http.StatusBadRequest, "invalid loginUrl")
		return
	}

	var providerUUID string
	if err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		return db.QueryRowContext(r.Context(), `
			INSERT INTO infra_providers (name, favicon_link, login_url)
			VALUES (?, ?, ?)
			RETURNING uuid
		`, name, nullableString(faviconURL), nullableString(loginURL)).Scan(&providerUUID)
	}); err != nil {
		shared.SendError(w, http.StatusBadRequest, "failed to create infra provider", err, cfg)
		return
	}

	writeProviderResponseWithStatus(w, r, manager, cfg, providerUUID, http.StatusCreated)
}

func handleUpdateProvider(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) {
	var req updateProviderRequest
	if err := decodeJSONBody(r, &req); err != nil {
		shared.WriteJSONError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	req.UUID = strings.TrimSpace(req.UUID)
	if req.UUID == "" {
		shared.WriteJSONError(w, http.StatusBadRequest, "uuid is required")
		return
	}

	updates := make([]string, 0, 4)
	args := make([]any, 0, 5)

	if req.Name.Set {
		if req.Name.Value == nil {
			shared.WriteJSONError(w, http.StatusBadRequest, "name is required")
			return
		}
		name, err := validateProviderName(*req.Name.Value)
		if err != nil {
			shared.WriteJSONError(w, http.StatusBadRequest, err.Error())
			return
		}
		updates = append(updates, "name = ?")
		args = append(args, name)
	}

	if req.FaviconURL.Set {
		faviconURL, err := normalizeOptionalURL(req.FaviconURL.Value)
		if err != nil {
			shared.WriteJSONError(w, http.StatusBadRequest, "invalid faviconLink")
			return
		}
		updates = append(updates, "favicon_link = ?")
		args = append(args, nullableString(faviconURL))
	}

	if req.LoginURL.Set {
		loginURL, err := normalizeOptionalURL(req.LoginURL.Value)
		if err != nil {
			shared.WriteJSONError(w, http.StatusBadRequest, "invalid loginUrl")
			return
		}
		updates = append(updates, "login_url = ?")
		args = append(args, nullableString(loginURL))
	}

	if len(updates) > 0 {
		updates = append(updates, "updated_at = now()")
		args = append(args, req.UUID)
		query := "UPDATE infra_providers SET " + strings.Join(updates, ", ") + " WHERE uuid = ?"
		if err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
			result, err := db.ExecContext(r.Context(), query, args...)
			if err != nil {
				return err
			}
			affected, err := result.RowsAffected()
			if err != nil {
				return err
			}
			if affected == 0 {
				return sql.ErrNoRows
			}
			return nil
		}); err != nil {
			if err == sql.ErrNoRows {
				shared.WriteJSONError(w, http.StatusNotFound, "infra provider not found")
				return
			}
			shared.SendError(w, http.StatusBadRequest, "failed to update infra provider", err, cfg)
			return
		}
	}

	writeProviderResponse(w, r, manager, cfg, req.UUID)
}

func handleDeleteProvider(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig, providerUUID string) {
	if providerUUID == "" {
		shared.WriteJSONError(w, http.StatusBadRequest, "uuid is required")
		return
	}

	if err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		result, err := db.ExecContext(r.Context(), `DELETE FROM infra_providers WHERE uuid = ?`, providerUUID)
		if err != nil {
			return err
		}
		affected, err := result.RowsAffected()
		if err != nil {
			return err
		}
		if affected == 0 {
			return sql.ErrNoRows
		}
		return nil
	}); err != nil {
		if err == sql.ErrNoRows {
			shared.WriteJSONError(w, http.StatusNotFound, "infra provider not found")
			return
		}
		shared.SendError(w, http.StatusBadRequest, "failed to delete infra provider", err, cfg)
		return
	}

	shared.WriteJSON(w, http.StatusOK, map[string]any{"response": map[string]any{"isDeleted": true}})
}

func writeProvidersResponse(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) {
	items, err := getProviders(r.Context(), manager)
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to fetch infra providers", err, cfg)
		return
	}
	shared.WriteJSON(w, http.StatusOK, map[string]any{
		"response": map[string]any{
			"total":     len(items),
			"providers": items,
		},
	})
}

func writeProviderResponse(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig, providerUUID string) {
	writeProviderResponseWithStatus(w, r, manager, cfg, providerUUID, http.StatusOK)
}

func writeProviderResponseWithStatus(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig, providerUUID string, status int) {
	item, err := getProvider(r.Context(), manager, providerUUID)
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to fetch infra provider", err, cfg)
		return
	}
	if item == nil {
		shared.WriteJSONError(w, http.StatusNotFound, "infra provider not found")
		return
	}
	shared.WriteJSON(w, status, map[string]any{"response": item})
}

func getProviders(ctx context.Context, manager *dbmanager.DatabaseManager) ([]providerRecord, error) {
	items := make([]providerRecord, 0)
	historyByProvider := make(map[string]map[string]any)
	nodesByProvider := make(map[string][]providerNode)

	if err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		rows, err := db.QueryContext(ctx, `
			SELECT uuid, name, favicon_link, login_url, created_at, updated_at
			FROM infra_providers
			ORDER BY created_at ASC
		`)
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			var rec providerRecord
			var favicon, loginURL sql.NullString
			if scanErr := rows.Scan(&rec.UUID, &rec.Name, &favicon, &loginURL, &rec.CreatedAt, &rec.UpdatedAt); scanErr != nil {
				return scanErr
			}
			if favicon.Valid {
				rec.Favicon = &favicon.String
			}
			if loginURL.Valid {
				rec.LoginURL = &loginURL.String
			}
			items = append(items, rec)
		}
		return rows.Err()
	}); err != nil {
		return nil, err
	}

	if err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		rows, err := db.QueryContext(ctx, `
			SELECT provider_uuid, COALESCE(ROUND(SUM(amount)::numeric, 2)::float8, 0), COUNT(*)
			FROM infra_billing_history
			GROUP BY provider_uuid
		`)
		if err != nil {
			return err
		}
		defer rows.Close()
		for rows.Next() {
			var providerUUID string
			var totalAmount float64
			var totalBills int
			if scanErr := rows.Scan(&providerUUID, &totalAmount, &totalBills); scanErr != nil {
				return scanErr
			}
			historyByProvider[providerUUID] = map[string]any{
				"totalAmount": totalAmount,
				"totalBills":  totalBills,
			}
		}
		return rows.Err()
	}); err != nil {
		return nil, err
	}

	if err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		rows, err := db.QueryContext(ctx, `
			SELECT ibn.provider_uuid, ibn.node_uuid, n.name, n.country_code
			FROM infra_billing_nodes ibn
			JOIN nodes n ON n.uuid = ibn.node_uuid
			ORDER BY n.view_position ASC
		`)
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			var providerUUID string
			var node providerNode
			var details providerNodeDetails
			if scanErr := rows.Scan(&providerUUID, &details.NodeUUID, &node.Name, &details.CountryCode); scanErr != nil {
				return scanErr
			}
			node.Details = &details
			nodesByProvider[providerUUID] = append(nodesByProvider[providerUUID], node)
		}
		return rows.Err()
	}); err != nil {
		return nil, err
	}

	for i := range items {
		history := historyByProvider[items[i].UUID]
		if history == nil {
			history = map[string]any{
				"totalAmount": float64(0),
				"totalBills":  0,
			}
		}
		nodes := nodesByProvider[items[i].UUID]
		if nodes == nil {
			nodes = make([]providerNode, 0)
		}
		items[i].History = history
		items[i].Nodes = nodes
	}

	return items, nil
}

func getProvider(ctx context.Context, manager *dbmanager.DatabaseManager, providerUUID string) (*providerRecord, error) {
	providerUUID = strings.TrimSpace(providerUUID)
	if providerUUID == "" {
		return nil, nil
	}

	items, err := getProviders(ctx, manager)
	if err != nil {
		return nil, err
	}
	for i := range items {
		if items[i].UUID == providerUUID {
			return &items[i], nil
		}
	}
	return nil, nil
}

func validateProviderName(value string) (string, error) {
	name := strings.TrimSpace(value)
	if len(name) < 2 {
		return "", errProviderNameTooShort{}
	}
	if len(name) > 30 {
		return "", errProviderNameTooLong{}
	}
	return name, nil
}

type errProviderNameTooShort struct{}

func (errProviderNameTooShort) Error() string { return "name must be at least 2 characters" }

type errProviderNameTooLong struct{}

func (errProviderNameTooLong) Error() string { return "name must be less than 30 characters" }

func normalizeOptionalURL(value *string) (*string, error) {
	if value == nil {
		return nil, nil
	}

	clean := strings.TrimSpace(*value)
	if clean == "" {
		return nil, nil
	}

	parsed, err := url.ParseRequestURI(clean)
	if err != nil || parsed.Scheme == "" || parsed.Host == "" {
		if err != nil {
			return nil, err
		}
		return nil, errors.New("url must include scheme and host")
	}
	return &clean, nil
}

func nullableString(value *string) any {
	if value == nil {
		return nil
	}
	return *value
}
