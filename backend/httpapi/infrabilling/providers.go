package infrabilling

import (
	"context"
	"database/sql"
	"net/http"
	"time"

	"exodus/backend/config"
	dbmanager "exodus/backend/db/manager"
	"exodus/backend/httpapi/shared"
)

type providerNode struct {
	NodeUUID    string `json:"nodeUuid"`
	Name        string `json:"name"`
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

func ProvidersHandler(manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			shared.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}

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
}

func getProviders(ctx context.Context, manager *dbmanager.DatabaseManager) ([]providerRecord, error) {
	items := make([]providerRecord, 0)
	historyByProvider := make(map[string]map[string]any)
	nodesByProvider := make(map[string][]providerNode)

	if err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		rows, err := db.QueryContext(ctx, `
			SELECT uuid, name, favicon_link, login_url, created_at, updated_at
			FROM infra_providers
			ORDER BY name ASC
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
			SELECT provider_uuid, COALESCE(SUM(amount), 0), COUNT(*)
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
			ORDER BY n.name ASC
		`)
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			var providerUUID string
			var node providerNode
			if scanErr := rows.Scan(&providerUUID, &node.NodeUUID, &node.Name, &node.CountryCode); scanErr != nil {
				return scanErr
			}
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
