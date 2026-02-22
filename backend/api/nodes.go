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

// RowScanner позволяет использовать одну функцию сканирования для sql.Rows и sql.Row
type RowScanner interface {
	Scan(dest ...interface{}) error
}

// Node представляет сущность узла для API
type Node struct {
	UUID                    string     `json:"uuid"`
	ID                      *int64     `json:"id,omitempty"`
	Name                    string     `json:"name"`
	Address                 string     `json:"address"`
	Port                    *int       `json:"port,omitempty"`
	APISchema               string     `json:"api_schema"`
	APIPath                 string     `json:"api_path"`
	APIMetadata             string     `json:"api_metadata"`
	ActiveConfigProfileUUID *string    `json:"active_config_profile_uuid,omitempty"`
	IsConnected             bool       `json:"is_connected"`
	IsConnecting            bool       `json:"is_connecting"`
	IsDisabled              bool       `json:"is_disabled"`
	LastStatusChange        *time.Time `json:"last_status_change,omitempty"`
	LastStatusMessage       *string    `json:"last_status_message,omitempty"`
	XrayVersion             *string    `json:"xray_version,omitempty"`
	NodeVersion             *string    `json:"node_version,omitempty"`
	XrayUptime              string     `json:"xray_uptime"`
	UsersOnline             *int       `json:"users_online,omitempty"`
	ConsumptionMultiplier   int64      `json:"consumption_multiplier"`
	IsTrafficTrackingActive bool       `json:"is_traffic_tracking_active"`
	TrafficResetDay         *int       `json:"traffic_reset_day,omitempty"`
	TrafficLimitBytes       *int64     `json:"traffic_limit_bytes,omitempty"`
	TrafficUsedBytes        *int64     `json:"traffic_used_bytes,omitempty"`
	NotifyPercent           *int       `json:"notify_percent,omitempty"`
	ProviderUUID            *string    `json:"provider_uuid,omitempty"`
	ViewPosition            int        `json:"view_position"`
	CountryCode             string     `json:"country_code"`
	Tags                    []string   `json:"tags"`
	CPUCount                *int       `json:"cpu_count,omitempty"`
	CPUModel                *string    `json:"cpu_model,omitempty"`
	TotalRAM                *string    `json:"total_ram,omitempty"`
	CreatedAt               time.Time  `json:"created_at"`
	UpdatedAt               time.Time  `json:"updated_at"`
}

// NodeCreateRequest — структура для создания ноды (POST)
type NodeCreateRequest struct {
	Name                    string   `json:"name"`
	Address                 string   `json:"address"`
	Port                    int      `json:"port"`
	APISchema               string   `json:"api_schema"`
	APIPath                 string   `json:"api_path"`
	APIMetadata             string   `json:"api_metadata"`
	IsDisabled              bool     `json:"is_disabled"`
	ConsumptionMultiplier   int64    `json:"consumption_multiplier"`
	IsTrafficTrackingActive bool     `json:"is_traffic_tracking_active"`
	TrafficResetDay         int      `json:"traffic_reset_day"`
	TrafficLimitBytes       int64    `json:"traffic_limit_bytes"`
	NotifyPercent           int      `json:"notify_percent"`
	ViewPosition            int      `json:"view_position"`
	CountryCode             string   `json:"country_code"`
	Tags                    []string `json:"tags"`
}

// NodeUpdateRequest — структура для частичного обновления (PATCH)
type NodeUpdateRequest struct {
	Name                    *string   `json:"name,omitempty"`
	Address                 *string   `json:"address,omitempty"`
	Port                    *int      `json:"port,omitempty"`
	APISchema               *string   `json:"api_schema,omitempty"`
	APIPath                 *string   `json:"api_path,omitempty"`
	APIMetadata             *string   `json:"api_metadata,omitempty"`
	ActiveConfigProfileUUID *string   `json:"active_config_profile_uuid,omitempty"`
	IsDisabled              *bool     `json:"is_disabled,omitempty"`
	ConsumptionMultiplier   *int64    `json:"consumption_multiplier,omitempty"`
	IsTrafficTrackingActive *bool     `json:"is_traffic_tracking_active,omitempty"`
	TrafficResetDay         *int      `json:"traffic_reset_day,omitempty"`
	TrafficLimitBytes       *int64    `json:"traffic_limit_bytes,omitempty"`
	NotifyPercent           *int      `json:"notify_percent,omitempty"`
	ProviderUUID            *string   `json:"provider_uuid,omitempty"`
	ViewPosition            *int      `json:"view_position,omitempty"`
	CountryCode             *string   `json:"country_code,omitempty"`
	Tags                    *[]string `json:"tags,omitempty"`
}

// scanNode — вспомогательная функция, убирающая дублирование Scan во всем файле
func scanNode(scanner RowScanner) (Node, error) {
	var n Node
	var lastStatusChange sql.NullTime
	var lastStatusMessage, xrayVersion, nodeVersion, providerUUID, activeConfigProfileUUID, tagsJSON sql.NullString
	var usersOnline, trafficResetDay, notifyPercent, cpuCount, trafficLimitBytes, trafficUsedBytes, port, id sql.NullInt64
	var cpuModel, totalRam sql.NullString

	err := scanner.Scan(
		&n.UUID, &id, &n.Name, &n.Address, &port, &n.APISchema, &n.APIPath, &n.APIMetadata, &activeConfigProfileUUID,
		&n.IsConnected, &n.IsConnecting, &n.IsDisabled, &lastStatusChange, &lastStatusMessage,
		&xrayVersion, &nodeVersion, &n.XrayUptime, &usersOnline, &n.ConsumptionMultiplier,
		&n.IsTrafficTrackingActive, &trafficResetDay, &trafficLimitBytes, &trafficUsedBytes,
		&notifyPercent, &providerUUID, &n.ViewPosition, &n.CountryCode, &tagsJSON,
		&cpuCount, &cpuModel, &totalRam, &n.CreatedAt, &n.UpdatedAt,
	)
	if err != nil {
		return n, err
	}

	// Обработка nullable полей
	if id.Valid {
		n.ID = &id.Int64
	}
	if port.Valid {
		p := int(port.Int64)
		n.Port = &p
	}
	if activeConfigProfileUUID.Valid {
		n.ActiveConfigProfileUUID = &activeConfigProfileUUID.String
	}
	if lastStatusChange.Valid {
		n.LastStatusChange = &lastStatusChange.Time
	}
	if lastStatusMessage.Valid {
		n.LastStatusMessage = &lastStatusMessage.String
	}
	if xrayVersion.Valid {
		n.XrayVersion = &xrayVersion.String
	}
	if nodeVersion.Valid {
		n.NodeVersion = &nodeVersion.String
	}
	if usersOnline.Valid {
		u := int(usersOnline.Int64)
		n.UsersOnline = &u
	}
	if trafficResetDay.Valid {
		t := int(trafficResetDay.Int64)
		n.TrafficResetDay = &t
	}
	if trafficLimitBytes.Valid {
		n.TrafficLimitBytes = &trafficLimitBytes.Int64
	}
	if trafficUsedBytes.Valid {
		n.TrafficUsedBytes = &trafficUsedBytes.Int64
	}
	if notifyPercent.Valid {
		np := int(notifyPercent.Int64)
		n.NotifyPercent = &np
	}
	if providerUUID.Valid {
		n.ProviderUUID = &providerUUID.String
	}
	if cpuCount.Valid {
		c := int(cpuCount.Int64)
		n.CPUCount = &c
	}
	if cpuModel.Valid {
		n.CPUModel = &cpuModel.String
	}
	if totalRam.Valid {
		n.TotalRAM = &totalRam.String
	}

	if tagsJSON.Valid && tagsJSON.String != "" {
		if err := json.Unmarshal([]byte(tagsJSON.String), &n.Tags); err != nil {
			n.Tags = []string{}
		}
	} else {
		n.Tags = []string{}
	}

	return n, nil
}

// NodesHandler — GET/POST /api/v1/nodes
func NodesHandler(manager *manager.DatabaseManager, cfg *config.BackendConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handleGetNodes(w, r, manager, cfg)
		case http.MethodPost:
			handleCreateNode(w, r, manager, cfg)
		case http.MethodDelete:
			handleDeleteAllNodes(w, r, manager, cfg)
		default:
			http.Error(w, `{"error": "method not allowed"}`, http.StatusMethodNotAllowed)
		}
	}
}

// handleGetNodes handles GET /api/v1/nodes
func handleGetNodes(w http.ResponseWriter, r *http.Request, manager *manager.DatabaseManager, cfg *config.BackendConfig) {
	// Check if help requested
	if r.URL.Query().Has("help") {
		sendNodesHelp(w)
		return
	}

	ctx := r.Context()
	var nodes []Node

	err := manager.ExecuteHighPriority(func(db *sql.DB) error {
		query := `
			SELECT
				uuid, id, name, address, port, api_schema, api_path, api_metadata, active_config_profile_uuid,
				is_connected, is_connecting, is_disabled, last_status_change, last_status_message,
				xray_version, node_version, xray_uptime, users_online, consumption_multiplier,
				is_traffic_tracking_active, traffic_reset_day, traffic_limit_bytes, traffic_used_bytes,
				notify_percent, provider_uuid, view_position, country_code, tags,
				cpu_count, cpu_model, total_ram, created_at, updated_at
			FROM nodes
			ORDER BY view_position ASC, name ASC`

		rows, err := db.QueryContext(ctx, query)
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			n, err := scanNode(rows)
			if err != nil {
				return err
			}
			nodes = append(nodes, n)
		}
		return rows.Err()
	})

	if err != nil {
		sendError(w, http.StatusInternalServerError, "failed to fetch nodes", err, cfg)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(map[string]interface{}{"nodes": nodes, "count": len(nodes)})
}

// sendNodesHelp returns API documentation
func sendNodesHelp(w http.ResponseWriter) {
	help := map[string]interface{}{
		"description": "V2Ray Nodes Management API",
		"endpoints": map[string]interface{}{
			"GET /api/v1/nodes": map[string]interface{}{
				"description":   "Get all nodes",
				"query_params":  "?help - show this help message",
				"response":      "List of all nodes with their statistics",
			},
			"GET /api/v1/nodes/{uuid}": map[string]interface{}{
				"description": "Get single node by UUID",
			},
			"POST /api/v1/nodes": map[string]interface{}{
				"description": "Create a new node",
				"required_fields": []string{"name", "address", "port", "country_code"},
				"optional_fields": []string{"api_schema", "api_path", "api_metadata", "is_disabled", "consumption_multiplier", "is_traffic_tracking_active", "traffic_reset_day", "traffic_limit_bytes", "notify_percent", "view_position", "tags"},
				"example": map[string]interface{}{
					"name":                     "Germany-Frankfurt-01",
					"address":                  "192.168.1.100",
					"port":                     9253,
					"api_schema":               "grpc",
					"api_path":                 "",
					"is_disabled":              false,
					"consumption_multiplier":   100,
					"is_traffic_tracking_active": true,
					"traffic_reset_day":        1,
					"traffic_limit_bytes":      1099511627776,
					"notify_percent":           80,
					"view_position":            10,
					"country_code":             "DE",
					"tags":                     []string{"premium", "fast"},
				},
			},
			"PATCH /api/v1/nodes/{uuid}": map[string]interface{}{
				"description": "Update node (partial update)",
				"note":        "Send only fields you want to update",
				"example":     map[string]interface{}{"port": 8443, "is_disabled": true},
			},
			"DELETE /api/v1/nodes/{uuid}": map[string]interface{}{
				"description": "Delete a specific node",
			},
		},
		"response_fields": map[string]string{
			"uuid":                         "Node unique identifier",
			"name":                         "Node display name",
			"address":                      "Node IP address or domain",
			"port":                         "Node gRPC port",
			"api_schema":                   "API protocol (grpc, https, http)",
			"api_path":                     "API endpoint path",
			"is_connected":                 "Node is currently connected",
			"is_disabled":                  "Node is administratively disabled",
			"users_online":                 "Number of active users on this node",
			"traffic_used_bytes":           "Total traffic used through this node",
			"traffic_limit_bytes":          "Traffic limit (0 = unlimited)",
			"consumption_multiplier":       "Traffic multiplier percentage",
			"is_traffic_tracking_active":   "Traffic tracking enabled",
			"last_status_message":          "Last connection status message",
			"last_status_change":           "Last status change timestamp",
			"country_code":                 "2-letter country code",
			"tags":                         "Node tags array",
			"created_at":                   "Node creation timestamp",
			"updated_at":                   "Node last update timestamp",
		},
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(help)
}

// handleCreateNode handles POST /api/v1/nodes
func handleCreateNode(w http.ResponseWriter, r *http.Request, manager *manager.DatabaseManager, cfg *config.BackendConfig) {
	var req NodeCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, http.StatusBadRequest, "invalid JSON", err, cfg)
		return
	}

	// Validate required fields
	if req.Name == "" {
		sendError(w, http.StatusBadRequest, "name is required", nil, cfg)
		return
	}
	if req.Address == "" {
		sendError(w, http.StatusBadRequest, "address is required", nil, cfg)
		return
	}
	if req.Port < 1 || req.Port > 65535 {
		sendError(w, http.StatusBadRequest, "invalid port", nil, cfg)
		return
	}
	if len(req.CountryCode) != 2 {
		sendError(w, http.StatusBadRequest, "country_code must be 2 letters", nil, cfg)
		return
	}

	// Generate UUID
	nodeUUID := uuid.New().String()

	// Marshal tags to JSON
	tagsJSON, err := json.Marshal(req.Tags)
	if err != nil {
		sendError(w, http.StatusInternalServerError, "failed to marshal tags", err, cfg)
		return
	}

	ctx := r.Context()
	err = manager.ExecuteHighPriority(func(db *sql.DB) error {
		query := `
			INSERT INTO nodes (
				uuid, name, address, port, api_schema, api_path, api_metadata,
				is_disabled, consumption_multiplier, is_traffic_tracking_active,
				traffic_reset_day, traffic_limit_bytes, notify_percent,
				view_position, country_code, tags,
				created_at, updated_at
			) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)`

		_, err := db.ExecContext(ctx, query,
			nodeUUID, req.Name, req.Address, req.Port, req.APISchema, req.APIPath, req.APIMetadata,
			req.IsDisabled, req.ConsumptionMultiplier, req.IsTrafficTrackingActive,
			req.TrafficResetDay, req.TrafficLimitBytes, req.NotifyPercent,
			req.ViewPosition, req.CountryCode, string(tagsJSON))

		return err
	})

	if err != nil {
		sendError(w, http.StatusInternalServerError, "failed to create node", err, cfg)
		return
	}

	// Return created node
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "node created",
		"uuid":    nodeUUID,
	})
}

// NodeByUUIDHandler — GET/PATCH/DELETE /api/v1/nodes/{uuid}
func NodeByUUIDHandler(manager *manager.DatabaseManager, cfg *config.BackendConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/api/v1/nodes/")
		nodeUUID := strings.TrimSpace(path)

		if _, err := uuid.Parse(nodeUUID); err != nil {
			sendError(w, http.StatusBadRequest, "invalid UUID format", nil, cfg)
			return
		}

		switch r.Method {
		case http.MethodGet:
			handleGetNode(w, r, manager, cfg, nodeUUID)
		case http.MethodPatch:
			handlePatchNode(w, r, manager, cfg, nodeUUID)
		case http.MethodDelete:
			handleDeleteNode(w, r, manager, cfg, nodeUUID)
		default:
			http.Error(w, `{"error": "method not allowed"}`, http.StatusMethodNotAllowed)
		}
	}
}

func handleGetNode(w http.ResponseWriter, r *http.Request, manager *manager.DatabaseManager, cfg *config.BackendConfig, nodeUUID string) {
	var node Node
	err := manager.ExecuteHighPriority(func(db *sql.DB) error {
		query := `SELECT 
					uuid, id, name, address, port, api_schema, api_path, api_metadata, active_config_profile_uuid,
					is_connected, is_connecting, is_disabled, last_status_change, last_status_message,
					xray_version, node_version, xray_uptime, users_online, consumption_multiplier,
					is_traffic_tracking_active, traffic_reset_day, traffic_limit_bytes, traffic_used_bytes,
					notify_percent, provider_uuid, view_position, country_code, tags,
					cpu_count, cpu_model, total_ram, created_at, updated_at
				  FROM nodes WHERE uuid = ?`
		row := db.QueryRowContext(r.Context(), query, nodeUUID)
		var scanErr error
		node, scanErr = scanNode(row)
		return scanErr
	})

	if err != nil {
		if err == sql.ErrNoRows {
			sendError(w, http.StatusNotFound, "node not found", nil, cfg)
		} else {
			sendError(w, http.StatusInternalServerError, "failed to fetch node", err, cfg)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(map[string]interface{}{"node": node})
}

func handlePatchNode(w http.ResponseWriter, r *http.Request, manager *manager.DatabaseManager, cfg *config.BackendConfig, nodeUUID string) {
	var req NodeUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, http.StatusBadRequest, "invalid JSON", err, cfg)
		return
	}

	if err := req.Validate(); err != nil {
		sendError(w, http.StatusBadRequest, err.Error(), nil, cfg)
		return
	}

	var clauses []string
	var args []interface{}

	// Динамическая сборка запроса
	add := func(col string, val interface{}) {
		clauses = append(clauses, fmt.Sprintf("%s = ?", col))
		args = append(args, val)
	}

	if req.Name != nil {
		add("name", *req.Name)
	}
	if req.Address != nil {
		add("address", *req.Address)
	}
	if req.Port != nil {
		add("port", *req.Port)
	}
	if req.APISchema != nil {
		add("api_schema", *req.APISchema)
	}
	if req.APIPath != nil {
		add("api_path", *req.APIPath)
	}
	if req.APIMetadata != nil {
		add("api_metadata", *req.APIMetadata)
	}
	if req.IsDisabled != nil {
		add("is_disabled", *req.IsDisabled)
	}
	if req.ConsumptionMultiplier != nil {
		add("consumption_multiplier", *req.ConsumptionMultiplier)
	}
	if req.IsTrafficTrackingActive != nil {
		add("is_traffic_tracking_active", *req.IsTrafficTrackingActive)
	}
	if req.TrafficResetDay != nil {
		add("traffic_reset_day", *req.TrafficResetDay)
	}
	if req.TrafficLimitBytes != nil {
		add("traffic_limit_bytes", *req.TrafficLimitBytes)
	}
	if req.NotifyPercent != nil {
		add("notify_percent", *req.NotifyPercent)
	}
	if req.ViewPosition != nil {
		add("view_position", *req.ViewPosition)
	}
	if req.CountryCode != nil {
		add("country_code", *req.CountryCode)
	}

	if req.ActiveConfigProfileUUID != nil {
		if *req.ActiveConfigProfileUUID == "" {
			clauses = append(clauses, "active_config_profile_uuid = NULL")
		} else {
			add("active_config_profile_uuid", *req.ActiveConfigProfileUUID)
		}
	}
	if req.ProviderUUID != nil {
		if *req.ProviderUUID == "" {
			clauses = append(clauses, "provider_uuid = NULL")
		} else {
			add("provider_uuid", *req.ProviderUUID)
		}
	}
	if req.Tags != nil {
		b, _ := json.Marshal(*req.Tags)
		add("tags", string(b))
	}

	if len(clauses) == 0 {
		sendError(w, http.StatusBadRequest, "no fields to update", nil, cfg)
		return
	}

	args = append(args, nodeUUID)
	query := fmt.Sprintf("UPDATE nodes SET %s, updated_at = CURRENT_TIMESTAMP WHERE uuid = ?", strings.Join(clauses, ", "))

	err := manager.ExecuteHighPriority(func(db *sql.DB) error {
		result, err := db.ExecContext(r.Context(), query, args...)
		if err != nil {
			return err
		}
		// Check if any row was actually updated
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
			sendError(w, http.StatusNotFound, "node not found", nil, cfg)
		} else {
			sendError(w, http.StatusInternalServerError, "update failed", err, cfg)
		}
		return
	}

	handleGetNode(w, r, manager, cfg, nodeUUID) // Возвращаем обновленный объект
}

// handleDeleteNode handles DELETE /api/v1/nodes/{uuid}
func handleDeleteNode(w http.ResponseWriter, r *http.Request, manager *manager.DatabaseManager, cfg *config.BackendConfig, nodeUUID string) {
	ctx := r.Context()
	
	// Get node name for logging
	var nodeName string
	err := manager.ExecuteHighPriority(func(db *sql.DB) error {
		return db.QueryRowContext(ctx, "SELECT name FROM nodes WHERE uuid = ?", nodeUUID).Scan(&nodeName)
	})
	
	if err != nil {
		if err == sql.ErrNoRows {
			sendError(w, http.StatusNotFound, "node not found", nil, cfg)
		} else {
			sendError(w, http.StatusInternalServerError, "failed to find node", err, cfg)
		}
		return
	}

	err = manager.ExecuteHighPriority(func(db *sql.DB) error {
		_, err := db.ExecContext(ctx, "DELETE FROM nodes WHERE uuid = ?", nodeUUID)
		return err
	})

	if err != nil {
		sendError(w, http.StatusInternalServerError, "failed to delete node", err, cfg)
		return
	}

	cfg.Logger.Info("Node deleted", "uuid", nodeUUID, "name", nodeName)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":  "node deleted",
		"uuid":     nodeUUID,
		"node_name": nodeName,
	})
}

// handleDeleteAllNodes handles DELETE /api/v1/nodes (delete all nodes)
func handleDeleteAllNodes(w http.ResponseWriter, r *http.Request, manager *manager.DatabaseManager, cfg *config.BackendConfig) {
	ctx := r.Context()
	
	// Require confirmation
	if r.URL.Query().Get("confirm") != "true" {
		sendError(w, http.StatusBadRequest, "confirmation required. Use DELETE /api/v1/nodes?confirm=true", nil, cfg)
		return
	}

	var deletedCount int
	err := manager.ExecuteHighPriority(func(db *sql.DB) error {
		result, err := db.ExecContext(ctx, "DELETE FROM nodes")
		if err != nil {
			return err
		}
		count, err := result.RowsAffected()
		if err != nil {
			return err
		}
		deletedCount = int(count)
		return nil
	})

	if err != nil {
		sendError(w, http.StatusInternalServerError, "failed to delete nodes", err, cfg)
		return
	}

	cfg.Logger.Info("All nodes deleted", "count", deletedCount)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "all nodes deleted",
		"count":   deletedCount,
	})
}

// Вспомогательные методы валидации и ошибок
func (r *NodeUpdateRequest) Validate() error {
	if r.Name != nil && *r.Name == "" {
		return fmt.Errorf("name cannot be empty")
	}
	if r.Port != nil && (*r.Port < 1 || *r.Port > 65535) {
		return fmt.Errorf("invalid port")
	}
	if r.CountryCode != nil && len(*r.CountryCode) != 2 {
		return fmt.Errorf("country_code must be 2 letters")
	}
	return nil
}

func sendError(w http.ResponseWriter, code int, msg string, err error, cfg *config.BackendConfig) {
	if err != nil {
		cfg.Logger.Error(msg, "error", err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}

// formatBytes converts bytes to human readable format
func formatBytes(bytes int64) string {
	if bytes == 0 {
		return "0 B"
	}
	const unit = 1024
	sizes := []string{"B", "KB", "MB", "GB", "TB", "PB"}
	i := int64(0)
	for int64(bytes) >= unit && i < int64(len(sizes)-1) {
		bytes /= unit
		i++
	}
	return fmt.Sprintf("%.2f %s", float64(bytes)/float64(1024), sizes[i])
}

// NodeSummary represents a simplified node response for frontend
type NodeSummary struct {
	UUID           string `json:"uuid"`
	Name           string `json:"name"`
	Address        string `json:"address"`
	Port           int    `json:"port"`
	APISchema      string `json:"api_schema"`
	APIPath        string `json:"api_path"`
	UsersOnline    int    `json:"users_online"`
	IsActive       bool   `json:"is_active"`
	IsConnected    bool   `json:"is_connected"`
	IsDisabled     bool   `json:"is_disabled"`
	TrafficUsed    string `json:"traffic_used"`
	TrafficLimit   string `json:"traffic_limit"`
	TrafficUsedRaw int64  `json:"traffic_used_bytes"`
	TrafficLimitRaw int64 `json:"traffic_limit_bytes"`
	CountryCode    string `json:"country_code"`
	Tags           []string `json:"tags"`
	LastStatusMsg  string `json:"last_status_message"`
	CreatedAt      string `json:"created_at"`
}

// handleGetNodesSummary returns simplified node list for frontend
func handleGetNodesSummary(w http.ResponseWriter, r *http.Request, manager *manager.DatabaseManager, cfg *config.BackendConfig) {
	ctx := r.Context()
	var summaries []NodeSummary

	err := manager.ExecuteHighPriority(func(db *sql.DB) error {
		query := `
			SELECT uuid, name, address, port, api_schema, api_path, users_online,
			       is_connected, is_disabled, traffic_used_bytes, traffic_limit_bytes,
			       country_code, tags, last_status_message, created_at
			FROM nodes
			WHERE is_disabled = 0
			ORDER BY view_position ASC, name ASC`

		rows, err := db.QueryContext(ctx, query)
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			var s NodeSummary
			var tagsJSON string
			err := rows.Scan(
				&s.UUID, &s.Name, &s.Address, &s.Port, &s.APISchema, &s.APIPath,
				&s.UsersOnline, &s.IsConnected, &s.IsDisabled, &s.TrafficUsedRaw,
				&s.TrafficLimitRaw, &s.CountryCode, &tagsJSON, &s.LastStatusMsg, &s.CreatedAt,
			)
			if err != nil {
				return err
			}

			s.IsActive = s.IsConnected && !s.IsDisabled
			s.TrafficUsed = formatBytes(s.TrafficUsedRaw)
			s.TrafficLimit = formatBytes(s.TrafficLimitRaw)

			if tagsJSON != "" && tagsJSON != "[]" {
				json.Unmarshal([]byte(tagsJSON), &s.Tags)
			} else {
				s.Tags = []string{}
			}

			summaries = append(summaries, s)
		}
		return rows.Err()
	})

	if err != nil {
		sendError(w, http.StatusInternalServerError, "failed to fetch nodes", err, cfg)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(map[string]interface{}{"nodes": summaries, "count": len(summaries)})
}

// NodesSummaryHandler returns simplified node list for frontend
func NodesSummaryHandler(manager *manager.DatabaseManager, cfg *config.BackendConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, `{"error": "method not allowed"}`, http.StatusMethodNotAllowed)
			return
		}
		handleGetNodesSummary(w, r, manager, cfg)
	}
}
