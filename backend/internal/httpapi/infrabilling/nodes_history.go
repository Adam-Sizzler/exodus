package infrabilling

import (
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

func BillingNodesHandler(manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			shared.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}

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
				billingNodes = append(billingNodes, item)
			}
			return rows.Err()
		})
		if err != nil {
			shared.SendError(w, http.StatusInternalServerError, "failed to fetch infra billing nodes", err, cfg)
			return
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
			shared.SendError(w, http.StatusInternalServerError, "failed to fetch available infra billing nodes", err, cfg)
			return
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
			shared.SendError(w, http.StatusInternalServerError, "failed to fetch infra billing stats", err, cfg)
			return
		}

		shared.WriteJSON(w, http.StatusOK, map[string]any{
			"response": map[string]any{
				"totalBillingNodes":          len(billingNodes),
				"billingNodes":               billingNodes,
				"availableBillingNodes":      availableNodes,
				"totalAvailableBillingNodes": len(availableNodes),
				"stats": map[string]any{
					"upcomingNodesCount":   upcomingCount,
					"currentMonthPayments": currentMonthPayments,
					"totalSpent":           totalSpent,
				},
			},
		})
	}
}

func BillingHistoryHandler(manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			shared.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}

		start := parseIntQuery(r.URL.Query().Get("start"), 0, 0, 1_000_000)
		size := parseIntQuery(r.URL.Query().Get("size"), 50, 1, 500)
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
			shared.SendError(w, http.StatusInternalServerError, "failed to fetch infra billing history", err, cfg)
			return
		}

		shared.WriteJSON(w, http.StatusOK, map[string]any{
			"response": map[string]any{
				"records": records,
				"total":   total,
			},
		})
	}
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
