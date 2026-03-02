package subscription

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"v2ray-stat/backend/panel/config"
	dbmanager "v2ray-stat/backend/panel/db/manager"
	"v2ray-stat/backend/panel/httpapi/shared"

	"github.com/google/uuid"
)

type SubscriptionPageConfig struct {
	UUID         string          `json:"uuid"`
	ViewPosition int             `json:"view_position"`
	Name         string          `json:"name"`
	Config       json.RawMessage `json:"config"`
	CreatedAt    time.Time       `json:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at"`
}

func SubscriptionPageConfigsHandler(manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
			return
		}

		ctx := r.Context()
		var configs []SubscriptionPageConfig

		err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
			rows, err := db.QueryContext(ctx, `
                SELECT uuid, view_position, name, config, created_at, updated_at
                FROM subscription_page_config
                ORDER BY view_position ASC, name ASC
            `)
			if err != nil {
				return err
			}
			defer rows.Close()

			for rows.Next() {
				var cfgItem SubscriptionPageConfig
				var viewPosition sql.NullInt64
				var configStr sql.NullString

				if err := rows.Scan(&cfgItem.UUID, &viewPosition, &cfgItem.Name, &configStr, &cfgItem.CreatedAt, &cfgItem.UpdatedAt); err != nil {
					return err
				}
				if viewPosition.Valid {
					cfgItem.ViewPosition = int(viewPosition.Int64)
				}
				if configStr.Valid {
					cfgItem.Config = json.RawMessage(configStr.String)
				}
				configs = append(configs, cfgItem)
			}

			return rows.Err()
		})

		if err != nil {
			shared.SendError(w, http.StatusInternalServerError, "failed to fetch subscription page configs", err, cfg)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"response": configs,
		})
	}
}

func SubscriptionPageConfigByUUIDHandler(manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
			return
		}

		uuidStr := strings.TrimSpace(strings.TrimPrefix(r.URL.Path, "/api/subscription-page-configs/"))
		if uuidStr == "" {
			http.NotFound(w, r)
			return
		}
		if _, err := uuid.Parse(uuidStr); err != nil {
			shared.SendError(w, http.StatusBadRequest, "invalid UUID format", err, cfg)
			return
		}

		ctx := r.Context()
		var cfgItem SubscriptionPageConfig

		err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
			row := db.QueryRowContext(ctx, `
                SELECT uuid, view_position, name, config, created_at, updated_at
                FROM subscription_page_config
                WHERE uuid = ?
                LIMIT 1
            `, uuidStr)

			var viewPosition sql.NullInt64
			var configStr sql.NullString
			if err := row.Scan(&cfgItem.UUID, &viewPosition, &cfgItem.Name, &configStr, &cfgItem.CreatedAt, &cfgItem.UpdatedAt); err != nil {
				return err
			}

			if viewPosition.Valid {
				cfgItem.ViewPosition = int(viewPosition.Int64)
			}
			if configStr.Valid {
				cfgItem.Config = json.RawMessage(configStr.String)
			}
			return nil
		})

		if err != nil {
			if err == sql.ErrNoRows {
				shared.SendError(w, http.StatusNotFound, "subscription page config not found", nil, cfg)
				return
			}
			shared.SendError(w, http.StatusInternalServerError, "failed to fetch subscription page config", err, cfg)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"response": cfgItem,
		})
	}
}
