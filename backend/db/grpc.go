package db

import (
	"database/sql"
	"fmt"

	"exodus/backend/config"
	dbmanager "exodus/backend/db/manager"
	"exodus/backend/dbutil"
)

// DBNode represents a node loaded from database.
type DBNode struct {
	UUID                    string
	Name                    string
	Address                 string
	Port                    int
	APISchema               string
	APIPath                 string
	IsDisabled              bool
	ConsumptionMultiplier   int64
	IsTrafficTrackingActive bool
	TrafficResetDay         int
	TrafficLimitBytes       int64
	NotifyPercent           int
	ViewPosition            int
	CountryCode             string
	Tags                    []string
}

// LoadNodesFromDB loads all active nodes from the database.
func LoadNodesFromDB(manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) ([]DBNode, error) {
	var nodes []DBNode

	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		query := `
			SELECT uuid, name, address, port, api_schema, api_path,
			       is_disabled, consumption_multiplier, is_traffic_tracking_active,
			       traffic_reset_day, traffic_limit_bytes, notify_percent,
			       view_position, country_code, tags
			FROM nodes
			WHERE is_disabled = false
			ORDER BY view_position ASC, name ASC`

		rows, err := db.Query(query)
		if err != nil {
			return fmt.Errorf("query nodes: %w", err)
		}
		defer rows.Close()

		for rows.Next() {
			var n DBNode
			var port sql.NullInt64
			var apiSchema, apiPath, countryCode sql.NullString
			var tags dbutil.StringArray
			var consumptionMultiplier, trafficLimitBytes, trafficResetDay, notifyPercent, viewPosition sql.NullInt64
			var isTrafficTrackingActive sql.NullBool

			err := rows.Scan(
				&n.UUID, &n.Name, &n.Address, &port, &apiSchema, &apiPath,
				&n.IsDisabled, &consumptionMultiplier, &isTrafficTrackingActive,
				&trafficResetDay, &trafficLimitBytes, &notifyPercent, &viewPosition,
				&countryCode, &tags,
			)
			if err != nil {
				return fmt.Errorf("scan node: %w", err)
			}

			if port.Valid {
				n.Port = int(port.Int64)
			} else {
				n.Port = 9253
			}
			if apiSchema.Valid {
				n.APISchema = apiSchema.String
			} else {
				n.APISchema = "grpc"
			}
			if apiPath.Valid {
				n.APIPath = apiPath.String
			} else {
				n.APIPath = ""
			}
			if consumptionMultiplier.Valid {
				n.ConsumptionMultiplier = consumptionMultiplier.Int64
			} else {
				n.ConsumptionMultiplier = 100
			}
			if isTrafficTrackingActive.Valid {
				n.IsTrafficTrackingActive = isTrafficTrackingActive.Bool
			} else {
				n.IsTrafficTrackingActive = true
			}
			if trafficResetDay.Valid {
				n.TrafficResetDay = int(trafficResetDay.Int64)
			} else {
				n.TrafficResetDay = 1
			}
			if trafficLimitBytes.Valid {
				n.TrafficLimitBytes = trafficLimitBytes.Int64
			}
			if notifyPercent.Valid {
				n.NotifyPercent = int(notifyPercent.Int64)
			} else {
				n.NotifyPercent = 80
			}
			if viewPosition.Valid {
				n.ViewPosition = int(viewPosition.Int64)
			}
			if countryCode.Valid {
				n.CountryCode = countryCode.String
			}

			if len(tags) > 0 {
				n.Tags = tags.Slice()
			} else {
				n.Tags = []string{}
			}

			nodes = append(nodes, n)
		}

		return rows.Err()
	})

	if err != nil {
		return nil, err
	}

	cfg.Logger.Info("Loaded nodes from database", "count", len(nodes))
	return nodes, nil
}
