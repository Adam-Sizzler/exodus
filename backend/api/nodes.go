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
	ID                      int64      `json:"id,omitempty"`
	Name                    string     `json:"name"`
	Address                 string     `json:"address"`
	Port                    *int       `json:"port,omitempty"`
	APISchema               string     `json:"api_schema"`   // grpc, https, http
	APIPath                 string     `json:"api_path"`     // путь, например /grpcnode
	APIMetadata             string     `json:"api_metadata"` // JSON строка с путями к сертификатам
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
	var usersOnline, trafficResetDay, notifyPercent, cpuCount, trafficLimitBytes, trafficUsedBytes, port sql.NullInt64
	var cpuModel, totalRam sql.NullString

	err := scanner.Scan(
		&n.UUID, &n.ID, &n.Name, &n.Address, &port, &n.APISchema, &n.APIPath, &n.APIMetadata, &activeConfigProfileUUID,
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

// NodesHandler — GET /api/v1/nodes
func NodesHandler(manager *manager.DatabaseManager, cfg *config.BackendConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, `{"error": "method not allowed"}`, http.StatusMethodNotAllowed)
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
}

// NodeByUUIDHandler — GET/PATCH /api/v1/nodes/{uuid}
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
		_, err := db.ExecContext(r.Context(), query, args...)
		return err
	})

	if err != nil {
		sendError(w, http.StatusInternalServerError, "update failed", err, cfg)
		return
	}

	handleGetNode(w, r, manager, cfg, nodeUUID) // Возвращаем обновленный объект
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
