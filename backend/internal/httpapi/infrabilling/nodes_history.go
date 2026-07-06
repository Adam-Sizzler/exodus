package infrabilling

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"exodus/internal/config"
	dbmanager "exodus/internal/db/manager"
	"exodus/internal/httpapi/shared"
)

type billingNodeProvider struct {
	UUID       string  `json:"uuid"`
	Name       string  `json:"name"`
	LoginURL   *string `json:"loginUrl"`
	FaviconURL *string `json:"faviconLink"`
}

type billingNodeNode struct {
	UUID        string `json:"uuid"`
	Name        string `json:"name"`
	CountryCode string `json:"countryCode"`
}

type billingNodeRecord struct {
	UUID          string              `json:"uuid"`
	NodeUUID      string              `json:"nodeUuid"`
	Name          *string             `json:"name"`
	ProviderUUID  string              `json:"providerUuid"`
	Provider      billingNodeProvider `json:"provider"`
	Node          billingNodeNode     `json:"node"`
	NextBillingAt time.Time           `json:"nextBillingAt"`
	CreatedAt     time.Time           `json:"createdAt"`
	UpdatedAt     time.Time           `json:"updatedAt"`
}

type availableNodeRecord struct {
	UUID        string `json:"uuid"`
	Name        string `json:"name"`
	CountryCode string `json:"countryCode"`
}

type billingHistoryProvider struct {
	UUID       string  `json:"uuid"`
	Name       string  `json:"name"`
	FaviconURL *string `json:"faviconLink"`
}

type billingHistoryRecord struct {
	UUID         string                 `json:"uuid"`
	ProviderUUID string                 `json:"providerUuid"`
	Amount       float64                `json:"amount"`
	BilledAt     time.Time              `json:"billedAt"`
	Provider     billingHistoryProvider `json:"provider"`
}

type createBillingNodeRequest struct {
	NodeUUID      string  `json:"nodeUuid"`
	ProviderUUID  string  `json:"providerUuid"`
	NextBillingAt *string `json:"nextBillingAt"`
}

type updateBillingNodeRequest struct {
	UUIDs         []string `json:"uuids"`
	NextBillingAt string   `json:"nextBillingAt"`
}

type createBillingHistoryRequest struct {
	ProviderUUID string  `json:"providerUuid"`
	Amount       float64 `json:"amount"`
	BilledAt     string  `json:"billedAt"`
}

func BillingNodesHandler(manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			writeBillingNodesResponse(w, r, manager, cfg, http.StatusOK)
		case http.MethodPost:
			handleCreateBillingNode(w, r, manager, cfg)
		case http.MethodPatch:
			handleUpdateBillingNodes(w, r, manager, cfg)
		case http.MethodDelete:
			handleDeleteBillingNode(w, r, manager, cfg)
		default:
			shared.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
		}
	}
}

func BillingHistoryHandler(manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			start, size := historyPaginationFromQuery(r)
			writeBillingHistoryResponse(w, r, manager, cfg, http.StatusOK, start, size)
		case http.MethodPost:
			handleCreateBillingHistory(w, r, manager, cfg)
		case http.MethodDelete:
			handleDeleteBillingHistory(w, r, manager, cfg)
		default:
			shared.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
		}
	}
}

func handleCreateBillingNode(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) {
	var req createBillingNodeRequest
	if err := decodeJSONBody(r, &req); err != nil {
		shared.WriteJSONError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if strings.TrimSpace(req.NodeUUID) == "" || strings.TrimSpace(req.ProviderUUID) == "" {
		shared.WriteJSONError(w, http.StatusBadRequest, "nodeUuid and providerUuid are required")
		return
	}

	nextBillingAt := time.Now().UTC()
	if req.NextBillingAt != nil && strings.TrimSpace(*req.NextBillingAt) != "" {
		parsed, err := parseAPITime(*req.NextBillingAt)
		if err != nil {
			shared.WriteJSONError(w, http.StatusBadRequest, "invalid nextBillingAt")
			return
		}
		nextBillingAt = parsed
	}

	if err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		_, err := db.ExecContext(r.Context(), `
			INSERT INTO infra_billing_nodes (node_uuid, provider_uuid, next_billing_at)
			VALUES (?, ?, ?)
		`, req.NodeUUID, req.ProviderUUID, nextBillingAt)
		return err
	}); err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to create infra billing node", err, cfg)
		return
	}

	writeBillingNodesResponse(w, r, manager, cfg, http.StatusCreated)
}

func handleUpdateBillingNodes(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) {
	var req updateBillingNodeRequest
	if err := decodeJSONBody(r, &req); err != nil {
		shared.WriteJSONError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if len(req.UUIDs) == 0 {
		shared.WriteJSONError(w, http.StatusBadRequest, "uuids are required")
		return
	}
	nextBillingAt, err := parseAPITime(req.NextBillingAt)
	if err != nil {
		shared.WriteJSONError(w, http.StatusBadRequest, "invalid nextBillingAt")
		return
	}

	if err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		for _, uuid := range req.UUIDs {
			uuid = strings.TrimSpace(uuid)
			if uuid == "" {
				continue
			}
			if _, execErr := db.ExecContext(r.Context(), `
				UPDATE infra_billing_nodes
				SET next_billing_at = ?, updated_at = now()
				WHERE uuid = ?
			`, nextBillingAt, uuid); execErr != nil {
				return execErr
			}
		}
		return nil
	}); err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to update infra billing node", err, cfg)
		return
	}

	writeBillingNodesResponse(w, r, manager, cfg, http.StatusOK)
}

func handleDeleteBillingNode(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) {
	uuid := uuidFromPath(r.URL.Path, "/api/infra-billing/nodes")
	if uuid == "" {
		shared.WriteJSONError(w, http.StatusBadRequest, "uuid is required")
		return
	}
	if err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		_, err := db.ExecContext(r.Context(), `DELETE FROM infra_billing_nodes WHERE uuid = ?`, uuid)
		return err
	}); err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to delete infra billing node", err, cfg)
		return
	}

	writeBillingNodesResponse(w, r, manager, cfg, http.StatusOK)
}

func handleCreateBillingHistory(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) {
	var req createBillingHistoryRequest
	if err := decodeJSONBody(r, &req); err != nil {
		shared.WriteJSONError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if strings.TrimSpace(req.ProviderUUID) == "" {
		shared.WriteJSONError(w, http.StatusBadRequest, "providerUuid is required")
		return
	}
	if req.Amount < 0 {
		shared.WriteJSONError(w, http.StatusBadRequest, "amount must be greater than or equal to 0")
		return
	}
	billedAt, err := parseAPITime(req.BilledAt)
	if err != nil {
		shared.WriteJSONError(w, http.StatusBadRequest, "invalid billedAt")
		return
	}

	if err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		_, err := db.ExecContext(r.Context(), `
			INSERT INTO infra_billing_history (provider_uuid, amount, billed_at)
			VALUES (?, ?, ?)
		`, req.ProviderUUID, req.Amount, billedAt)
		return err
	}); err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to create infra billing history record", err, cfg)
		return
	}

	writeBillingHistoryResponse(w, r, manager, cfg, http.StatusCreated, 0, 50)
}

func handleDeleteBillingHistory(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) {
	uuid := uuidFromPath(r.URL.Path, "/api/infra-billing/history")
	if uuid == "" {
		shared.WriteJSONError(w, http.StatusBadRequest, "uuid is required")
		return
	}
	if err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		_, err := db.ExecContext(r.Context(), `DELETE FROM infra_billing_history WHERE uuid = ?`, uuid)
		return err
	}); err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to delete infra billing history record", err, cfg)
		return
	}

	writeBillingHistoryResponse(w, r, manager, cfg, http.StatusOK, 0, 50)
}

func writeBillingNodesResponse(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig, status int) {
	response, err := getBillingNodesResponse(r, manager)
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to fetch infra billing nodes", err, cfg)
		return
	}
	shared.WriteJSON(w, status, map[string]any{"response": response})
}

func writeBillingHistoryResponse(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig, status int, start, size int) {
	response, err := getBillingHistoryResponse(r, manager, start, size)
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to fetch infra billing history", err, cfg)
		return
	}
	shared.WriteJSON(w, status, map[string]any{"response": response})
}

func getBillingNodesResponse(r *http.Request, manager *dbmanager.DatabaseManager) (map[string]any, error) {
	billingNodes := make([]billingNodeRecord, 0)
	availableNodes := make([]availableNodeRecord, 0)
	upcomingCount := 0
	currentMonthPayments := float64(0)
	totalSpent := float64(0)

	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		rows, err := db.QueryContext(r.Context(), `
			SELECT
				ibn.uuid, ibn.node_uuid, ibn.provider_uuid, ibn.next_billing_at, ibn.created_at, ibn.updated_at,
				ip.uuid, ip.name, ip.login_url, ip.favicon_link,
				n.uuid, n.name, n.country_code
			FROM infra_billing_nodes ibn
			JOIN infra_providers ip ON ip.uuid = ibn.provider_uuid
			JOIN nodes n ON n.uuid = ibn.node_uuid
			ORDER BY ibn.next_billing_at ASC, n.name ASC
		`)
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			var item billingNodeRecord
			if scanErr := rows.Scan(
				&item.UUID, &item.NodeUUID, &item.ProviderUUID, &item.NextBillingAt, &item.CreatedAt, &item.UpdatedAt,
				&item.Provider.UUID, &item.Provider.Name, &item.Provider.LoginURL, &item.Provider.FaviconURL,
				&item.Node.UUID, &item.Node.Name, &item.Node.CountryCode,
			); scanErr != nil {
				return scanErr
			}
			item.Name = &item.Node.Name
			billingNodes = append(billingNodes, item)
		}
		return rows.Err()
	})
	if err != nil {
		return nil, err
	}

	err = manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		rows, err := db.QueryContext(r.Context(), `
			SELECT n.uuid, n.name, n.country_code
			FROM nodes n
			LEFT JOIN infra_billing_nodes ibn ON ibn.node_uuid = n.uuid
			WHERE ibn.node_uuid IS NULL
			ORDER BY n.name ASC
		`)
		if err != nil {
			return err
		}
		defer rows.Close()
		for rows.Next() {
			var item availableNodeRecord
			if scanErr := rows.Scan(&item.UUID, &item.Name, &item.CountryCode); scanErr != nil {
				return scanErr
			}
			availableNodes = append(availableNodes, item)
		}
		return rows.Err()
	})
	if err != nil {
		return nil, err
	}

	err = manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		if scanErr := db.QueryRowContext(r.Context(), `
			SELECT COUNT(*) FROM infra_billing_nodes
			WHERE next_billing_at <= (NOW() + INTERVAL '7 days')
		`).Scan(&upcomingCount); scanErr != nil {
			return scanErr
		}
		if scanErr := db.QueryRowContext(r.Context(), `
			SELECT COALESCE(SUM(amount), 0)
			FROM infra_billing_history
			WHERE billed_at >= date_trunc('month', NOW())
		`).Scan(&currentMonthPayments); scanErr != nil {
			return scanErr
		}
		return db.QueryRowContext(r.Context(), `
			SELECT COALESCE(SUM(amount), 0)
			FROM infra_billing_history
		`).Scan(&totalSpent)
	})
	if err != nil {
		return nil, err
	}

	return map[string]any{
		"totalBillingNodes":          len(billingNodes),
		"billingNodes":               billingNodes,
		"availableBillingNodes":      availableNodes,
		"totalAvailableBillingNodes": len(availableNodes),
		"stats": map[string]any{
			"upcomingNodesCount":   upcomingCount,
			"currentMonthPayments": currentMonthPayments,
			"totalSpent":           totalSpent,
		},
	}, nil
}

func getBillingHistoryResponse(r *http.Request, manager *dbmanager.DatabaseManager, start, size int) (map[string]any, error) {
	records := make([]billingHistoryRecord, 0)
	total := 0

	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		if scanErr := db.QueryRowContext(r.Context(), `SELECT COUNT(*) FROM infra_billing_history`).Scan(&total); scanErr != nil {
			return scanErr
		}

		rows, err := db.QueryContext(r.Context(), `
			SELECT
				ibh.uuid, ibh.provider_uuid, ibh.amount, ibh.billed_at,
				ip.uuid, ip.name, ip.favicon_link
			FROM infra_billing_history ibh
			JOIN infra_providers ip ON ip.uuid = ibh.provider_uuid
			ORDER BY ibh.billed_at DESC
			OFFSET ? LIMIT ?
		`, start, size)
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			var item billingHistoryRecord
			if scanErr := rows.Scan(
				&item.UUID, &item.ProviderUUID, &item.Amount, &item.BilledAt,
				&item.Provider.UUID, &item.Provider.Name, &item.Provider.FaviconURL,
			); scanErr != nil {
				return scanErr
			}
			records = append(records, item)
		}
		return rows.Err()
	})
	if err != nil {
		return nil, err
	}

	return map[string]any{
		"records": records,
		"total":   total,
	}, nil
}

func historyPaginationFromQuery(r *http.Request) (int, int) {
	start := parseIntQuery(r.URL.Query().Get("start"), 0, 0, 1_000_000)
	size := parseIntQuery(r.URL.Query().Get("size"), 50, 1, 500)
	return start, size
}

func decodeJSONBody(r *http.Request, target any) error {
	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	return decoder.Decode(target)
}

func parseAPITime(value string) (time.Time, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return time.Time{}, strconv.ErrSyntax
	}
	if parsed, err := time.Parse(time.RFC3339Nano, value); err == nil {
		return parsed, nil
	}
	return time.Parse("2006-01-02", value)
}

func uuidFromPath(path, base string) string {
	rest := strings.TrimPrefix(path, base)
	rest = strings.Trim(rest, "/")
	if rest == "" || strings.Contains(rest, "/") {
		return ""
	}
	return rest
}

func parseIntQuery(raw string, def, min, max int) int {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return def
	}
	v, err := strconv.Atoi(raw)
	if err != nil {
		return def
	}
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}
