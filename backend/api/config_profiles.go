package api

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"v2ray-stat/backend/config"
	"v2ray-stat/backend/db/manager"

	"github.com/google/uuid"
)

// ConfigProfile represents a config profile entity for API responses.
type ConfigProfile struct {
	UUID         string          `json:"uuid"`
	ViewPosition int             `json:"view_position"`
	Name         string          `json:"name"`
	Config       json.RawMessage `json:"config"` // JSON object for sing-box configuration
	CreatedAt    time.Time       `json:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at"`
}

// ConfigProfileInbound represents an inbound extracted from config profile.
type ConfigProfileInbound struct {
	UUID        string          `json:"uuid"`
	ProfileUUID string          `json:"profile_uuid"`
	Tag         string          `json:"tag"`
	Type        string          `json:"type"`
	Network     *string         `json:"network,omitempty"`
	Security    *string         `json:"security,omitempty"`
	Port        *int            `json:"port,omitempty"`
	RawInbound  json.RawMessage `json:"raw_inbound"`
}

// ConfigProfileCreateRequest represents a request to create a new config profile.
type ConfigProfileCreateRequest struct {
	ViewPosition int             `json:"view_position"`
	Name         string          `json:"name"`
	Config       json.RawMessage `json:"config"` // JSON object for sing-box configuration
}

// Validate validates the ConfigProfileCreateRequest fields.
func (r *ConfigProfileCreateRequest) Validate() error {
	if r.Name == "" {
		return fmt.Errorf("name is required")
	}

	// Validate that config is valid JSON
	if len(r.Config) == 0 {
		return fmt.Errorf("config is required")
	}

	var jsonConfig interface{}
	if err := json.Unmarshal(r.Config, &jsonConfig); err != nil {
		return fmt.Errorf("config must be valid JSON")
	}

	return nil
}

// ConfigProfileUpdateRequest represents a partial update request for a config profile.
// Only provided fields will be updated (PATCH semantics).
type ConfigProfileUpdateRequest struct {
	ViewPosition *int             `json:"view_position,omitempty"`
	Name         *string          `json:"name,omitempty"`
	Config       *json.RawMessage `json:"config,omitempty"` // JSON object for sing-box configuration
}

// Validate validates the ConfigProfileUpdateRequest fields.
func (r *ConfigProfileUpdateRequest) Validate() error {
	if r.Name != nil && *r.Name == "" {
		return fmt.Errorf("name cannot be empty")
	}

	if r.Config != nil && len(*r.Config) > 0 {
		var jsonConfig interface{}
		if err := json.Unmarshal(*r.Config, &jsonConfig); err != nil {
			return fmt.Errorf("config must be valid JSON")
		}
	}

	return nil
}

// HasUpdates checks if any field is set for update.
func (r *ConfigProfileUpdateRequest) HasUpdates() bool {
	return r.ViewPosition != nil || r.Name != nil || r.Config != nil
}

// scanConfigProfile scans a row into a ConfigProfile struct.
func scanConfigProfile(scanner RowScanner) (ConfigProfile, error) {
	var cp ConfigProfile
	var viewPosition sql.NullInt64
	var configStr sql.NullString

	err := scanner.Scan(
		&cp.UUID,
		&viewPosition,
		&cp.Name,
		&configStr,
		&cp.CreatedAt,
		&cp.UpdatedAt,
	)
	if err != nil {
		return cp, err
	}

	if viewPosition.Valid {
		cp.ViewPosition = int(viewPosition.Int64)
	}

	if configStr.Valid {
		cp.Config = json.RawMessage(configStr.String)
	}

	return cp, nil
}

// ConfigProfilesHandler handles GET/POST /api/v1/config-profiles
func ConfigProfilesHandler(manager *manager.DatabaseManager, cfg *config.BackendConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handleGetConfigProfiles(w, r, manager, cfg)
		case http.MethodPost:
			handleCreateConfigProfile(w, r, manager, cfg)
		default:
			http.Error(w, `{"error": "method not allowed"}`, http.StatusMethodNotAllowed)
		}
	}
}

// ConfigProfilesReorderHandler handles POST /api/v1/config-profiles/reorder
func ConfigProfilesReorderHandler(manager *manager.DatabaseManager, cfg *config.BackendConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
			return
		}

		var req ViewPositionReorderRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			sendError(w, http.StatusBadRequest, "invalid JSON", err, cfg)
			return
		}
		if err := req.Validate(); err != nil {
			sendError(w, http.StatusBadRequest, err.Error(), nil, cfg)
			return
		}

		err := manager.ExecuteHighPriority(func(db *sql.DB) error {
			return applyViewPositionReorder(r.Context(), db, "config_profiles", req.OrderedUUIDs, cfg)
		})
		if err != nil {
			sendError(w, http.StatusInternalServerError, "failed to reorder config profiles", err, cfg)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "config profiles reordered",
			"count":   len(req.OrderedUUIDs),
		})
	}
}

// handleGetConfigProfiles handles GET /api/v1/config-profiles
func handleGetConfigProfiles(w http.ResponseWriter, r *http.Request, manager *manager.DatabaseManager, cfg *config.BackendConfig) {
	ctx := r.Context()
	var profiles []ConfigProfile

	err := manager.ExecuteHighPriority(func(db *sql.DB) error {
		query := `
			SELECT uuid, view_position, name, config, created_at, updated_at
			FROM config_profiles
			ORDER BY view_position ASC, name ASC`

		rows, err := db.QueryContext(ctx, query)
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			cp, err := scanConfigProfile(rows)
			if err != nil {
				return err
			}
			profiles = append(profiles, cp)
		}
		return rows.Err()
	})

	if err != nil {
		sendError(w, http.StatusInternalServerError, "failed to fetch config profiles", err, cfg)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"profiles": profiles,
		"count":    len(profiles),
	})
}

// handleCreateConfigProfile handles POST /api/v1/config-profiles
func handleCreateConfigProfile(w http.ResponseWriter, r *http.Request, manager *manager.DatabaseManager, cfg *config.BackendConfig) {
	var req ConfigProfileCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, http.StatusBadRequest, "invalid JSON", err, cfg)
		return
	}

	// Validate required fields
	if err := req.Validate(); err != nil {
		sendError(w, http.StatusBadRequest, err.Error(), nil, cfg)
		return
	}

	// Generate UUID
	profileUUID := uuid.New().String()

	// Convert config JSON to string for storage
	configStr := string(req.Config)

	ctx := r.Context()
	err := manager.ExecuteHighPriority(func(db *sql.DB) error {
		query := `
			INSERT INTO config_profiles (
				uuid, view_position, name, config, created_at, updated_at
			) VALUES (?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)`

		_, err := db.ExecContext(ctx, query,
			profileUUID, req.ViewPosition, req.Name, configStr)
		if err != nil {
			return err
		}

		// Sync inbounds after creating profile
		_, err = syncConfigProfileInbounds(ctx, db, profileUUID, req.Config, cfg)
		return err
	})

	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			sendError(w, http.StatusConflict, "name already exists", err, cfg)
			return
		}
		sendError(w, http.StatusInternalServerError, "failed to create config profile", err, cfg)
		return
	}

	// Return created profile UUID
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "config profile created",
		"uuid":    profileUUID,
	})
}

// ConfigProfileByUUIDHandler handles GET/PATCH/DELETE /api/v1/config-profiles/{uuid}
func ConfigProfileByUUIDHandler(manager *manager.DatabaseManager, cfg *config.BackendConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/api/v1/config-profiles/")
		profileUUID := strings.TrimSpace(path)

		if _, err := uuid.Parse(profileUUID); err != nil {
			sendError(w, http.StatusBadRequest, "invalid UUID format", nil, cfg)
			return
		}

		switch r.Method {
		case http.MethodGet:
			handleGetConfigProfile(w, r, manager, cfg, profileUUID)
		case http.MethodPatch:
			handlePatchConfigProfile(w, r, manager, cfg, profileUUID)
		case http.MethodDelete:
			handleDeleteConfigProfile(w, r, manager, cfg, profileUUID)
		default:
			http.Error(w, `{"error": "method not allowed"}`, http.StatusMethodNotAllowed)
		}
	}
}

func handleGetConfigProfile(w http.ResponseWriter, r *http.Request, manager *manager.DatabaseManager, cfg *config.BackendConfig, profileUUID string) {
	var profile ConfigProfile
	err := manager.ExecuteHighPriority(func(db *sql.DB) error {
		query := `SELECT uuid, view_position, name, config, created_at, updated_at
				  FROM config_profiles WHERE uuid = ?`
		row := db.QueryRowContext(r.Context(), query, profileUUID)
		var scanErr error
		profile, scanErr = scanConfigProfile(row)
		return scanErr
	})

	if err != nil {
		if err == sql.ErrNoRows {
			sendError(w, http.StatusNotFound, "config profile not found", nil, cfg)
		} else {
			sendError(w, http.StatusInternalServerError, "failed to fetch config profile", err, cfg)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(map[string]interface{}{"profile": profile})
}

func handlePatchConfigProfile(w http.ResponseWriter, r *http.Request, manager *manager.DatabaseManager, cfg *config.BackendConfig, profileUUID string) {
	var req ConfigProfileUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, http.StatusBadRequest, "invalid JSON", err, cfg)
		return
	}

	if err := req.Validate(); err != nil {
		sendError(w, http.StatusBadRequest, err.Error(), nil, cfg)
		return
	}

	if !req.HasUpdates() {
		sendError(w, http.StatusBadRequest, "no fields to update", nil, cfg)
		return
	}

	var clauses []string
	var args []interface{}

	// Dynamic clause building
	add := func(col string, val interface{}) {
		clauses = append(clauses, fmt.Sprintf("%s = ?", col))
		args = append(args, val)
	}

	if req.ViewPosition != nil {
		add("view_position", *req.ViewPosition)
	}
	if req.Name != nil {
		add("name", *req.Name)
	}
	if req.Config != nil {
		add("config", string(*req.Config))
	}

	if len(clauses) == 0 {
		sendError(w, http.StatusBadRequest, "no fields to update", nil, cfg)
		return
	}

	args = append(args, profileUUID)
	query := fmt.Sprintf("UPDATE config_profiles SET %s, updated_at = CURRENT_TIMESTAMP WHERE uuid = ?", strings.Join(clauses, ", "))

	ctx := r.Context()
	configChanged := req.Config != nil
	var newConfig json.RawMessage
	if configChanged {
		newConfig = *req.Config
	}

	err := manager.ExecuteHighPriority(func(db *sql.DB) error {
		result, err := db.ExecContext(ctx, query, args...)
		if err != nil {
			return err
		}
		rowsAffected, err := result.RowsAffected()
		if err != nil {
			return err
		}
		if rowsAffected == 0 {
			return sql.ErrNoRows
		}

		// Sync inbounds if config was updated
		if configChanged {
			_, err = syncConfigProfileInbounds(ctx, db, profileUUID, newConfig, cfg)
			return err
		}
		return nil
	})

	if err != nil {
		if err == sql.ErrNoRows {
			sendError(w, http.StatusNotFound, "config profile not found", nil, cfg)
		} else {
			if strings.Contains(err.Error(), "UNIQUE constraint failed") {
				sendError(w, http.StatusConflict, "name already exists", err, cfg)
			} else {
				sendError(w, http.StatusInternalServerError, "update failed", err, cfg)
			}
		}
		return
	}

	handleGetConfigProfile(w, r, manager, cfg, profileUUID) // Return updated profile
}

func handleDeleteConfigProfile(w http.ResponseWriter, r *http.Request, manager *manager.DatabaseManager, cfg *config.BackendConfig, profileUUID string) {
	ctx := r.Context()

	// Get profile name for logging
	var profileName string
	err := manager.ExecuteHighPriority(func(db *sql.DB) error {
		return db.QueryRowContext(ctx, "SELECT name FROM config_profiles WHERE uuid = ?", profileUUID).Scan(&profileName)
	})

	if err != nil {
		if err == sql.ErrNoRows {
			sendError(w, http.StatusNotFound, "config profile not found", nil, cfg)
		} else {
			sendError(w, http.StatusInternalServerError, "failed to find config profile", err, cfg)
		}
		return
	}

	err = manager.ExecuteHighPriority(func(db *sql.DB) error {
		_, err := db.ExecContext(ctx, "DELETE FROM config_profiles WHERE uuid = ?", profileUUID)
		return err
	})

	if err != nil {
		sendError(w, http.StatusInternalServerError, "failed to delete config profile", err, cfg)
		return
	}

	cfg.Logger.Info("Config profile deleted", "uuid", profileUUID, "name", profileName)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "config profile deleted",
		"uuid":    profileUUID,
		"name":    profileName,
	})
}

// parseConfigInbounds extracts inbounds from a config JSON and returns them as ConfigProfileInbound structs.
// The tag field is used as the unique identifier for inbounds within a profile.
func parseConfigInbounds(profileUUID string, configJSON json.RawMessage) ([]ConfigProfileInbound, error) {
	// Parse the config JSON to extract inbounds
	var configData map[string]interface{}
	if err := json.Unmarshal(configJSON, &configData); err != nil {
		return nil, fmt.Errorf("failed to parse config JSON: %w", err)
	}

	// Get inbounds array
	inboundsRaw, ok := configData["inbounds"]
	if !ok {
		// No inbounds in config, return empty slice
		return []ConfigProfileInbound{}, nil
	}

	inboundsArray, ok := inboundsRaw.([]interface{})
	if !ok {
		return nil, fmt.Errorf("inbounds must be an array")
	}

	var inbounds []ConfigProfileInbound
	seenTags := make(map[string]bool)

	for _, inboundRaw := range inboundsArray {
		inboundMap, ok := inboundRaw.(map[string]interface{})
		if !ok {
			continue
		}

		// Extract tag (required, unique within profile)
		tagRaw, ok := inboundMap["tag"]
		if !ok {
			continue // Skip inbounds without tag
		}
		tag, ok := tagRaw.(string)
		if !ok || tag == "" {
			continue
		}

		// Check for duplicate tags within the same profile
		if seenTags[tag] {
			continue // Skip duplicate tags, keep first occurrence
		}
		seenTags[tag] = true

		// Extract type/protocol
		typeStr := ""
		if typeRaw, ok := inboundMap["type"].(string); ok {
			typeStr = typeRaw
		} else if protocolRaw, ok := inboundMap["protocol"].(string); ok {
			typeStr = protocolRaw
		}

		// Extract network
		var network *string
		if networkRaw, ok := inboundMap["network"].(string); ok && networkRaw != "" {
			network = &networkRaw
		}

		// Extract security
		var security *string
		if securityRaw, ok := inboundMap["security"].(string); ok && securityRaw != "" {
			security = &securityRaw
		}

		// Extract port
		var port *int
		if portRaw, ok := inboundMap["listen_port"].(float64); ok {
			p := int(portRaw)
			port = &p
		} else if portRaw, ok := inboundMap["port"].(float64); ok {
			p := int(portRaw)
			port = &p
		}

		// Convert inbound back to JSON for raw_inbound storage
		rawInbound, err := json.Marshal(inboundMap)
		if err != nil {
			continue
		}

		inbounds = append(inbounds, ConfigProfileInbound{
			UUID:        uuid.New().String(),
			ProfileUUID: profileUUID,
			Tag:         tag,
			Type:        typeStr,
			Network:     network,
			Security:    security,
			Port:        port,
			RawInbound:  rawInbound,
		})
	}

	return inbounds, nil
}

// syncConfigProfileInbounds synchronizes inbounds for a config profile.
// It performs:
// 1. Delete existing inbounds for the profile (CASCADE handles related tables)
// 2. Parse new inbounds from config JSON
// 3. Insert new inbounds
// Returns the number of inbounds synced.
func syncConfigProfileInbounds(ctx context.Context, db *sql.DB, profileUUID string, configJSON json.RawMessage, cfg *config.BackendConfig) (int, error) {
	// Parse inbounds from config
	inbounds, err := parseConfigInbounds(profileUUID, configJSON)
	if err != nil {
		cfg.Logger.Error("Failed to parse config inbounds", "profile_uuid", profileUUID, "error", err)
		return 0, fmt.Errorf("failed to parse inbounds: %w", err)
	}

	// Delete existing inbounds for this profile (CASCADE handles related tables)
	_, err = db.ExecContext(ctx, "DELETE FROM config_profile_inbounds WHERE profile_uuid = ?", profileUUID)
	if err != nil {
		cfg.Logger.Error("Failed to delete existing inbounds", "profile_uuid", profileUUID, "error", err)
		return 0, fmt.Errorf("failed to delete existing inbounds: %w", err)
	}

	// Insert new inbounds
	for _, inbound := range inbounds {
		query := `
			INSERT INTO config_profile_inbounds (
				uuid, profile_uuid, tag, type, network, security, port, raw_inbound
			) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`

		var networkVal, securityVal interface{}
		if inbound.Network != nil {
			networkVal = *inbound.Network
		}
		if inbound.Security != nil {
			securityVal = *inbound.Security
		}
		var portVal interface{}
		if inbound.Port != nil {
			portVal = *inbound.Port
		}

		_, err := db.ExecContext(ctx, query,
			inbound.UUID,
			inbound.ProfileUUID,
			inbound.Tag,
			inbound.Type,
			networkVal,
			securityVal,
			portVal,
			inbound.RawInbound,
		)
		if err != nil {
			cfg.Logger.Error("Failed to insert inbound", "profile_uuid", profileUUID, "tag", inbound.Tag, "error", err)
			return 0, fmt.Errorf("failed to insert inbound %s: %w", inbound.Tag, err)
		}
	}

	cfg.Logger.Debug("Synced config profile inbounds", "profile_uuid", profileUUID, "inbounds_count", len(inbounds))
	return len(inbounds), nil
}
