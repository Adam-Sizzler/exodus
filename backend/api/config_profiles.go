package api

import (
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

// ConfigProfile represents a config profile entity for API responses and requests.
type ConfigProfile struct {
	UUID         string    `json:"uuid"`
	ViewPosition int       `json:"view_position"`
	Name         string    `json:"name"`
	Config       string    `json:"config"` // JSON string for sing-box configuration
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// ConfigProfileCreateRequest represents a request to create a new config profile.
type ConfigProfileCreateRequest struct {
	ViewPosition int    `json:"view_position"`
	Name         string `json:"name"`
	Config       string `json:"config"` // JSON string for sing-box configuration
}

// Validate validates the ConfigProfileCreateRequest fields.
func (r *ConfigProfileCreateRequest) Validate() error {
	if r.Name == "" {
		return fmt.Errorf("name is required")
	}

	// Validate that config is valid JSON
	if r.Config == "" {
		return fmt.Errorf("config is required")
	}

	var jsonConfig interface{}
	if err := json.Unmarshal([]byte(r.Config), &jsonConfig); err != nil {
		return fmt.Errorf("config must be valid JSON")
	}

	return nil
}

// ConfigProfileUpdateRequest represents a partial update request for a config profile.
// Only provided fields will be updated (PATCH semantics).
type ConfigProfileUpdateRequest struct {
	ViewPosition *int    `json:"view_position,omitempty"`
	Name         *string `json:"name,omitempty"`
	Config       *string `json:"config,omitempty"` // JSON string for sing-box configuration
}

// Validate validates the ConfigProfileUpdateRequest fields.
func (r *ConfigProfileUpdateRequest) Validate() error {
	if r.Name != nil && *r.Name == "" {
		return fmt.Errorf("name cannot be empty")
	}

	if r.Config != nil && *r.Config != "" {
		var jsonConfig interface{}
		if err := json.Unmarshal([]byte(*r.Config), &jsonConfig); err != nil {
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

	err := scanner.Scan(
		&cp.UUID,
		&viewPosition,
		&cp.Name,
		&cp.Config,
		&cp.CreatedAt,
		&cp.UpdatedAt,
	)
	if err != nil {
		return cp, err
	}

	if viewPosition.Valid {
		cp.ViewPosition = int(viewPosition.Int64)
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

// handleGetConfigProfiles handles GET /api/v1/config-profiles
func handleGetConfigProfiles(w http.ResponseWriter, r *http.Request, manager *manager.DatabaseManager, cfg *config.BackendConfig) {
	// Check if help requested
	if r.URL.Query().Has("help") {
		sendConfigProfilesHelp(w)
		return
	}

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

// sendConfigProfilesHelp returns API documentation
func sendConfigProfilesHelp(w http.ResponseWriter) {
	help := map[string]interface{}{
		"description": "V2Ray Config Profiles Management API (sing-box JSON configurations)",
		"endpoints": map[string]interface{}{
			"GET /api/v1/config-profiles": map[string]interface{}{
				"description":  "Get all config profiles",
				"query_params": "?help - show this help message",
				"response":     "List of all config profiles with their JSON configurations",
			},
			"GET /api/v1/config-profiles/{uuid}": map[string]interface{}{
				"description": "Get single config profile by UUID",
			},
			"POST /api/v1/config-profiles": map[string]interface{}{
				"description":   "Create a new config profile",
				"required_fields": []string{"name", "config"},
				"optional_fields": []string{"view_position"},
				"example": map[string]interface{}{
					"name":           "default-profile",
					"view_position":  0,
					"config": map[string]interface{}{
						"log": map[string]interface{}{
							"level": "info",
						},
						"dns": map[string]interface{}{
							"servers": []string{"8.8.8.8", "1.1.1.1"},
						},
						"inbounds": []map[string]interface{}{
							{
								"type": "mixed",
								"tag":  "mixed-in",
								"listen": "0.0.0.0",
								"listen_port": 2080,
							},
						},
						"outbounds": []map[string]interface{}{
							{
								"type": "direct",
								"tag":  "direct",
							},
						},
						"route": map[string]interface{}{
							"rules": []map[string]interface{}{},
						},
					},
				},
			},
			"PATCH /api/v1/config-profiles/{uuid}": map[string]interface{}{
				"description": "Update config profile (partial update)",
				"note":        "Send only fields you want to update",
				"example":     map[string]interface{}{"name": "updated-profile", "view_position": 5},
			},
			"DELETE /api/v1/config-profiles/{uuid}": map[string]interface{}{
				"description": "Delete a specific config profile",
			},
		},
		"response_fields": map[string]string{
			"uuid":           "Config profile unique identifier (auto-generated)",
			"name":           "Config profile name (unique)",
			"view_position":  "Display order position",
			"config":         "sing-box configuration as JSON string",
			"created_at":     "Profile creation timestamp",
			"updated_at":     "Profile last update timestamp",
		},
		"sing_box_config_example": `{
  "log": { "level": "info" },
  "dns": { "servers": ["8.8.8.8"] },
  "inbounds": [
    { "type": "mixed", "tag": "mixed-in", "listen": "0.0.0.0", "listen_port": 2080 }
  ],
  "outbounds": [
    { "type": "direct", "tag": "direct" }
  ],
  "route": { "rules": [] }
}`,
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(help)
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

	ctx := r.Context()
	err := manager.ExecuteHighPriority(func(db *sql.DB) error {
		query := `
			INSERT INTO config_profiles (
				uuid, view_position, name, config, created_at, updated_at
			) VALUES (?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)`

		_, err := db.ExecContext(ctx, query,
			profileUUID, req.ViewPosition, req.Name, req.Config)

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

	// Return created profile
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
		add("config", *req.Config)
	}

	if len(clauses) == 0 {
		sendError(w, http.StatusBadRequest, "no fields to update", nil, cfg)
		return
	}

	args = append(args, profileUUID)
	query := fmt.Sprintf("UPDATE config_profiles SET %s, updated_at = CURRENT_TIMESTAMP WHERE uuid = ?", strings.Join(clauses, ", "))

	err := manager.ExecuteHighPriority(func(db *sql.DB) error {
		result, err := db.ExecContext(r.Context(), query, args...)
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
