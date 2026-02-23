package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"v2ray-stat/backend/config"
	"v2ray-stat/backend/db/manager"
)

// ==================== CONFIG PROFILES WITH INBOUNDS (FOR UI) ====================

// ConfigProfileWithInbounds represents a config profile with its inbounds for UI display.
type ConfigProfileWithInbounds struct {
	UUID         string          `json:"uuid"`
	Name         string          `json:"name"`
	ViewPosition int             `json:"view_position"`
	Config       json.RawMessage `json:"config"`
	Inbounds     []InboundInfo   `json:"inbounds"`
	CreatedAt    time.Time       `json:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at"`
}

// InboundInfo represents an inbound with its details.
type InboundInfo struct {
	UUID       string          `json:"uuid"`
	Tag        string          `json:"tag"`
	Type       string          `json:"type"`
	Network    *string         `json:"network,omitempty"`
	Security   *string         `json:"security,omitempty"`
	Port       *int            `json:"port,omitempty"`
	RawInbound json.RawMessage `json:"raw_inbound"`
}

// ConfigProfilesWithInboundsHandler handles GET /api/v1/config-profiles-with-inbounds
func ConfigProfilesWithInboundsHandler(manager *manager.DatabaseManager, cfg *config.BackendConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, `{"error": "method not allowed"}`, http.StatusMethodNotAllowed)
			return
		}

		ctx := r.Context()
		var profiles []ConfigProfileWithInbounds

		err := manager.ExecuteHighPriority(func(db *sql.DB) error {
			// First, get all config profiles
			profileQuery := `
				SELECT uuid, view_position, name, config, created_at, updated_at
				FROM config_profiles
				ORDER BY view_position ASC, name ASC`

			profileRows, err := db.QueryContext(ctx, profileQuery)
			if err != nil {
				return err
			}
			defer profileRows.Close()

			// Store profiles temporarily
			type tempProfile struct {
				UUID         string
				ViewPosition int
				Name         string
				Config       json.RawMessage
				CreatedAt    time.Time
				UpdatedAt    time.Time
			}
			var tempProfiles []tempProfile

			for profileRows.Next() {
				var tp tempProfile
				var viewPosition sql.NullInt64
				var configStr sql.NullString

				if err := profileRows.Scan(&tp.UUID, &viewPosition, &tp.Name, &configStr, &tp.CreatedAt, &tp.UpdatedAt); err != nil {
					return err
				}

				if viewPosition.Valid {
					tp.ViewPosition = int(viewPosition.Int64)
				}
				if configStr.Valid {
					tp.Config = json.RawMessage(configStr.String)
				}

				tempProfiles = append(tempProfiles, tp)
			}
			profileRows.Close()

			// Get ALL inbounds in a single query
			inboundQuery := `
				SELECT uuid, profile_uuid, tag, type, network, security, port, raw_inbound
				FROM config_profile_inbounds
				ORDER BY profile_uuid, tag`

			inboundRows, err := db.QueryContext(ctx, inboundQuery)
			if err != nil {
				return err
			}
			defer inboundRows.Close()

			// Group inbounds by profile_uuid
			inboundsByProfile := make(map[string][]InboundInfo)
			for inboundRows.Next() {
				var inbound InboundInfo
				var profileUUID string
				var network, security sql.NullString
				var port sql.NullInt64
				var rawInbound string

				if err := inboundRows.Scan(&inbound.UUID, &profileUUID, &inbound.Tag, &inbound.Type, &network, &security, &port, &rawInbound); err != nil {
					return err
				}

				if network.Valid {
					inbound.Network = &network.String
				}
				if security.Valid {
					inbound.Security = &security.String
				}
				if port.Valid {
					p := int(port.Int64)
					inbound.Port = &p
				}
				inbound.RawInbound = json.RawMessage(rawInbound)

				inboundsByProfile[profileUUID] = append(inboundsByProfile[profileUUID], inbound)
			}
			inboundRows.Close()

			// Combine profiles with their inbounds
			for _, tp := range tempProfiles {
				cp := ConfigProfileWithInbounds{
					UUID:         tp.UUID,
					ViewPosition: tp.ViewPosition,
					Name:         tp.Name,
					Config:       tp.Config,
					Inbounds:     inboundsByProfile[tp.UUID],
					CreatedAt:    tp.CreatedAt,
					UpdatedAt:    tp.UpdatedAt,
				}
				profiles = append(profiles, cp)
			}

			return nil
		})

		if err != nil {
			sendError(w, http.StatusInternalServerError, "failed to fetch config profiles with inbounds", err, cfg)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"profiles": profiles,
			"count":    len(profiles),
		})
	}
}

// ==================== SQUAD DETAILS WITH INBOUNDS AND MEMBERS ====================

// SquadDetails represents a squad with its inbounds and members for UI display.
type SquadDetails struct {
	UUID         string          `json:"uuid"`
	Name         string          `json:"name"`
	ViewPosition int             `json:"view_position"`
	Inbounds     []InboundInfo   `json:"inbounds"`
	Members      []SquadMemberInfo `json:"members"`
	CreatedAt    time.Time       `json:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at"`
}

// SquadMemberInfo represents a squad member with user details.
type SquadMemberInfo struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
}

// SquadDetailsHandler handles GET /api/v1/squad-details/{uuid}
func SquadDetailsHandler(manager *manager.DatabaseManager, cfg *config.BackendConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, `{"error": "method not allowed"}`, http.StatusMethodNotAllowed)
			return
		}

		ctx := r.Context()

		// Extract UUID from path
		path := r.URL.Path
		lastSlash := 0
		for i := len(path) - 1; i >= 0; i-- {
			if path[i] == '/' {
				lastSlash = i
				break
			}
		}
		squadUUID := path[lastSlash+1:]

		if squadUUID == "" {
			sendError(w, http.StatusBadRequest, "squad UUID is required", nil, cfg)
			return
		}

		var squad SquadDetails

		err := manager.ExecuteHighPriority(func(db *sql.DB) error {
			// Get squad details
			query := `
				SELECT uuid, view_position, name, created_at, updated_at
				FROM internal_squads
				WHERE uuid = ?`

			var viewPosition sql.NullInt64
			err := db.QueryRowContext(ctx, query, squadUUID).Scan(
				&squad.UUID, &viewPosition, &squad.Name, &squad.CreatedAt, &squad.UpdatedAt)
			if err != nil {
				if err == sql.ErrNoRows {
					return fmt.Errorf("squad not found")
				}
				return err
			}

			if viewPosition.Valid {
				squad.ViewPosition = int(viewPosition.Int64)
			}

			// Get inbounds for this squad using JOIN
			inboundQuery := `
				SELECT i.uuid, i.tag, i.type, i.network, i.security, i.port, i.raw_inbound
				FROM internal_squad_inbounds si
				JOIN config_profile_inbounds i ON si.inbound_uuid = i.uuid
				WHERE si.internal_squad_uuid = ?
				ORDER BY i.tag`

			rows, err := db.QueryContext(ctx, inboundQuery, squadUUID)
			if err != nil {
				return err
			}
			defer rows.Close()

			for rows.Next() {
				var inbound InboundInfo
				var network, security sql.NullString
				var port sql.NullInt64
				var rawInbound string

				if err := rows.Scan(&inbound.UUID, &inbound.Tag, &inbound.Type, &network, &security, &port, &rawInbound); err != nil {
					return err
				}

				if network.Valid {
					inbound.Network = &network.String
				}
				if security.Valid {
					inbound.Security = &security.String
				}
				if port.Valid {
					p := int(port.Int64)
					inbound.Port = &p
				}
				inbound.RawInbound = json.RawMessage(rawInbound)

				squad.Inbounds = append(squad.Inbounds, inbound)
			}

			// Get members for this squad
			memberQuery := `
				SELECT m.user_id, u.username
				FROM internal_squad_members m
				JOIN users u ON m.user_id = u.t_id
				WHERE m.internal_squad_uuid = ?
				ORDER BY u.username`

			memberRows, err := db.QueryContext(ctx, memberQuery, squadUUID)
			if err != nil {
				return err
			}
			defer memberRows.Close()

			for memberRows.Next() {
				var member SquadMemberInfo
				if err := memberRows.Scan(&member.UserID, &member.Username); err != nil {
					return err
				}
				squad.Members = append(squad.Members, member)
			}

			return nil
		})

		if err != nil {
			if err.Error() == "squad not found" {
				sendError(w, http.StatusNotFound, "squad not found", nil, cfg)
				return
			}
			sendError(w, http.StatusInternalServerError, "failed to fetch squad details", err, cfg)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"squad": squad,
		})
	}
}

// ==================== ALL SQUADS SUMMARY (FOR USER ASSIGNMENT UI) ====================

// SquadSummary represents a squad with basic info for selection UI.
type SquadSummary struct {
	UUID         string `json:"uuid"`
	Name         string `json:"name"`
	ViewPosition int    `json:"view_position"`
	MembersCount int    `json:"members_count"`
	InboundsCount int   `json:"inbounds_count"`
}

// AllSquadsSummaryHandler handles GET /api/v1/squads-summary
func AllSquadsSummaryHandler(manager *manager.DatabaseManager, cfg *config.BackendConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, `{"error": "method not allowed"}`, http.StatusMethodNotAllowed)
			return
		}

		ctx := r.Context()
		var squads []SquadSummary

		err := manager.ExecuteHighPriority(func(db *sql.DB) error {
			query := `
				SELECT 
					s.uuid,
					s.view_position,
					s.name,
					COUNT(DISTINCT m.user_id) as members_count,
					COUNT(DISTINCT si.inbound_uuid) as inbounds_count
				FROM internal_squads s
				LEFT JOIN internal_squad_members m ON s.uuid = m.internal_squad_uuid
				LEFT JOIN internal_squad_inbounds si ON s.uuid = si.internal_squad_uuid
				GROUP BY s.uuid, s.view_position, s.name
				ORDER BY s.view_position ASC, s.name ASC`

			rows, err := db.QueryContext(ctx, query)
			if err != nil {
				return err
			}
			defer rows.Close()

			for rows.Next() {
				var squad SquadSummary
				var viewPosition sql.NullInt64
				var membersCount, inboundsCount sql.NullInt64

				if err := rows.Scan(&squad.UUID, &viewPosition, &squad.Name, &membersCount, &inboundsCount); err != nil {
					return err
				}

				if viewPosition.Valid {
					squad.ViewPosition = int(viewPosition.Int64)
				}
				if membersCount.Valid {
					squad.MembersCount = int(membersCount.Int64)
				}
				if inboundsCount.Valid {
					squad.InboundsCount = int(inboundsCount.Int64)
				}

				squads = append(squads, squad)
			}

			return rows.Err()
		})

		if err != nil {
			sendError(w, http.StatusInternalServerError, "failed to fetch squads summary", err, cfg)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"squads": squads,
			"count":  len(squads),
		})
	}
}

// ==================== NODE WITH CONFIG PROFILE INFO ====================

// NodeWithConfig represents a node with its config profile info.
type NodeWithConfig struct {
	UUID                     string  `json:"uuid"`
	ID                       int64   `json:"id"`
	Name                     string  `json:"name"`
	Address                  string  `json:"address"`
	Port                     int     `json:"port"`
	ActiveConfigProfileUUID  *string `json:"active_config_profile_uuid,omitempty"`
	ConfigProfileName        *string `json:"config_profile_name,omitempty"`
	IsConnected              bool    `json:"is_connected"`
	IsDisabled               bool    `json:"is_disabled"`
	ViewPosition             int     `json:"view_position"`
	CountryCode              string  `json:"country_code"`
	Tags                     string  `json:"tags"` // JSON array
}

// NodesWithConfigHandler handles GET /api/v1/nodes-with-config
func NodesWithConfigHandler(manager *manager.DatabaseManager, cfg *config.BackendConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, `{"error": "method not allowed"}`, http.StatusMethodNotAllowed)
			return
		}

		ctx := r.Context()
		var nodes []NodeWithConfig

		err := manager.ExecuteHighPriority(func(db *sql.DB) error {
			query := `
				SELECT 
					n.uuid,
					n.id,
					n.name,
					n.address,
					n.port,
					n.active_config_profile_uuid,
					COALESCE(cp.name, '') as config_profile_name,
					n.is_connected,
					n.is_disabled,
					n.view_position,
					n.country_code,
					COALESCE(n.tags, '[]') as tags
				FROM nodes n
				LEFT JOIN config_profiles cp ON n.active_config_profile_uuid = cp.uuid
				ORDER BY n.view_position ASC, n.name ASC`

			rows, err := db.QueryContext(ctx, query)
			if err != nil {
				return err
			}
			defer rows.Close()

			for rows.Next() {
				var node NodeWithConfig
				var activeConfigProfileUUID, configProfileName sql.NullString
				var tags string

				if err := rows.Scan(&node.UUID, &node.ID, &node.Name, &node.Address, &node.Port,
					&activeConfigProfileUUID, &configProfileName, &node.IsConnected, &node.IsDisabled,
					&node.ViewPosition, &node.CountryCode, &tags); err != nil {
					return err
				}

				if activeConfigProfileUUID.Valid && activeConfigProfileUUID.String != "" {
					node.ActiveConfigProfileUUID = &activeConfigProfileUUID.String
					if configProfileName.Valid && configProfileName.String != "" {
						node.ConfigProfileName = &configProfileName.String
					}
				}
				node.Tags = tags

				nodes = append(nodes, node)
			}

			return rows.Err()
		})

		if err != nil {
			sendError(w, http.StatusInternalServerError, "failed to fetch nodes with config", err, cfg)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"nodes": nodes,
			"count": len(nodes),
		})
	}
}
