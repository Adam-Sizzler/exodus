package nodes

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"regexp"
	"sort"
	"strings"
	"time"

	"exodus/internal/config"
	dbmanager "exodus/internal/db/manager"
	"exodus/internal/dbutil"
	"exodus/internal/httpapi/shared"
	monitor "exodus/internal/nodes"
	"exodus/internal/notifications"

	"github.com/google/uuid"
)

var nodeTagRegex = regexp.MustCompile(`^[A-Z0-9_:]+$`)

var (
	errConfigProfileNotFound       = errors.New("config profile not found")
	errConfigProfileInboundInvalid = errors.New("config profile inbound not found in specified profile")
	errNoEnabledNodes              = errors.New("enabled nodes not found")
)

type OptionalString struct {
	Set   bool
	Value *string
}

func (o *OptionalString) UnmarshalJSON(data []byte) error {
	o.Set = true
	if string(data) == "null" {
		o.Value = nil
		return nil
	}
	var value string
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	o.Value = &value
	return nil
}

type configProfileInboundResponse struct {
	UUID        string          `json:"uuid"`
	ProfileUUID string          `json:"profileUuid"`
	Tag         string          `json:"tag"`
	Type        string          `json:"type"`
	Network     *string         `json:"network"`
	Security    *string         `json:"security"`
	Port        *int            `json:"port"`
	RawInbound  json.RawMessage `json:"rawInbound"`
}

type providerResponse struct {
	UUID        string     `json:"uuid"`
	Name        string     `json:"name"`
	FaviconLink *string    `json:"faviconLink"`
	LoginURL    *string    `json:"loginUrl"`
	CreatedAt   *time.Time `json:"createdAt,omitempty"`
	UpdatedAt   *time.Time `json:"updatedAt,omitempty"`
}

type nodeAPI struct {
	UUID                    string     `json:"uuid"`
	Name                    string     `json:"name"`
	Address                 string     `json:"address"`
	Port                    *int       `json:"port"`
	APISchema               string     `json:"apiSchema"`
	APIPath                 string     `json:"apiPath"`
	IsConnected             bool       `json:"isConnected"`
	IsDisabled              bool       `json:"isDisabled"`
	IsConnecting            bool       `json:"isConnecting"`
	LastStatusChange        *time.Time `json:"lastStatusChange"`
	LastStatusMessage       *string    `json:"lastStatusMessage"`
	SingboxVersion          *string    `json:"singboxVersion"`
	NodeVersion             *string    `json:"nodeVersion"`
	SingboxUptime           string     `json:"singboxUptime"`
	IsTrafficTrackingActive bool       `json:"isTrafficTrackingActive"`
	TrafficResetDay         *int       `json:"trafficResetDay"`
	TrafficLimitBytes       *int64     `json:"trafficLimitBytes"`
	TrafficUsedBytes        *int64     `json:"trafficUsedBytes"`
	NotifyPercent           *int       `json:"notifyPercent"`
	UsersOnline             *int       `json:"usersOnline"`
	ViewPosition            int        `json:"viewPosition"`
	CountryCode             string     `json:"countryCode"`
	ConsumptionMultiplier   float64    `json:"consumptionMultiplier"`
	Tags                    []string   `json:"tags"`
	CPUCount                *int       `json:"cpuCount"`
	CPUModel                *string    `json:"cpuModel"`
	TotalRAM                *string    `json:"totalRam"`
	CreatedAt               time.Time  `json:"createdAt"`
	UpdatedAt               time.Time  `json:"updatedAt"`
	ConfigProfile           struct {
		ActiveConfigProfileUUID *string                        `json:"activeConfigProfileUuid"`
		ActiveInbounds          []configProfileInboundResponse `json:"activeInbounds"`
	} `json:"configProfile"`
	ProviderUUID *string           `json:"providerUuid"`
	Provider     *providerResponse `json:"provider"`
}

type nodeRecord struct {
	ID                      *int64
	UUID                    string
	Name                    string
	Address                 string
	Port                    *int
	APISchema               string
	APIPath                 string
	ActiveConfigProfileUUID *string
	IsConnected             bool
	IsConnecting            bool
	IsDisabled              bool
	LastStatusChange        *time.Time
	LastStatusMessage       *string
	SingboxVersion          *string
	NodeVersion             *string
	SingboxUptime           string
	UsersOnline             *int
	ConsumptionMultiplier   int64
	IsTrafficTrackingActive bool
	TrafficResetDay         *int
	TrafficLimitBytes       *int64
	TrafficUsedBytes        *int64
	NotifyPercent           *int
	ProviderUUID            *string
	ViewPosition            int
	CountryCode             string
	Tags                    []string
	CPUCount                *int
	CPUModel                *string
	TotalRAM                *string
	CreatedAt               time.Time
	UpdatedAt               time.Time
}

type configProfileRefRequest struct {
	ActiveConfigProfileUUID string   `json:"activeConfigProfileUuid"`
	ActiveInbounds          []string `json:"activeInbounds"`
}

type createNodeRequest struct {
	Name                    string                  `json:"name"`
	Address                 string                  `json:"address"`
	Port                    *int                    `json:"port,omitempty"`
	APISchema               *string                 `json:"apiSchema,omitempty"`
	APIPath                 *string                 `json:"apiPath,omitempty"`
	IsTrafficTrackingActive *bool                   `json:"isTrafficTrackingActive,omitempty"`
	TrafficLimitBytes       *int64                  `json:"trafficLimitBytes,omitempty"`
	NotifyPercent           *int                    `json:"notifyPercent,omitempty"`
	TrafficResetDay         *int                    `json:"trafficResetDay,omitempty"`
	CountryCode             *string                 `json:"countryCode,omitempty"`
	ConsumptionMultiplier   *float64                `json:"consumptionMultiplier,omitempty"`
	ConfigProfile           configProfileRefRequest `json:"configProfile"`
	ProviderUUID            *string                 `json:"providerUuid,omitempty"`
	Tags                    []string                `json:"tags,omitempty"`
}

type updateNodeRequest struct {
	UUID                    string                   `json:"uuid"`
	Name                    *string                  `json:"name,omitempty"`
	Address                 *string                  `json:"address,omitempty"`
	Port                    *int                     `json:"port,omitempty"`
	APISchema               *string                  `json:"apiSchema,omitempty"`
	APIPath                 *string                  `json:"apiPath,omitempty"`
	IsTrafficTrackingActive *bool                    `json:"isTrafficTrackingActive,omitempty"`
	TrafficLimitBytes       *int64                   `json:"trafficLimitBytes,omitempty"`
	NotifyPercent           *int                     `json:"notifyPercent,omitempty"`
	TrafficResetDay         *int                     `json:"trafficResetDay,omitempty"`
	CountryCode             *string                  `json:"countryCode,omitempty"`
	ConsumptionMultiplier   *float64                 `json:"consumptionMultiplier,omitempty"`
	ConfigProfile           *configProfileRefRequest `json:"configProfile,omitempty"`
	ProviderUUID            OptionalString           `json:"providerUuid,omitempty"`
	Tags                    *[]string                `json:"tags,omitempty"`
}

type reorderNodesRequest struct {
	Nodes []struct {
		UUID         string `json:"uuid"`
		ViewPosition int    `json:"viewPosition"`
	} `json:"nodes"`
}

type restartAllNodesRequest struct {
	ForceRestart *bool `json:"forceRestart,omitempty"`
}

type bulkNodesActionsRequest struct {
	UUIDs  []string `json:"uuids"`
	Action string   `json:"action"`
}

type bulkProfileModificationRequest struct {
	UUIDs         []string                `json:"uuids"`
	ConfigProfile configProfileRefRequest `json:"configProfile"`
}

func NodesHandler(manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handleGetNodes(w, r, manager, cfg)
		case http.MethodPost:
			handleCreateNode(w, r, manager, cfg)
		case http.MethodPatch:
			handleUpdateNode(w, r, manager, cfg)
		default:
			shared.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
		}
	}
}

func NodeByUUIDHandler(manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := trimNodesPath(r.URL.Path, "/")
		if path == "" {
			switch r.Method {
			case http.MethodGet:
				handleGetNodes(w, r, manager, cfg)
			case http.MethodPost:
				handleCreateNode(w, r, manager, cfg)
			case http.MethodPatch:
				handleUpdateNode(w, r, manager, cfg)
			default:
				shared.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
			}
			return
		}

		parts := strings.Split(path, "/")
		if len(parts) == 0 {
			http.NotFound(w, r)
			return
		}
		nodeUUID := parts[0]
		if _, err := uuid.Parse(nodeUUID); err != nil {
			shared.SendError(w, http.StatusBadRequest, "invalid UUID format", nil, cfg)
			return
		}

		if len(parts) >= 3 && parts[1] == "actions" && r.Method == http.MethodPost {
			switch parts[2] {
			case "enable":
				handleEnableNode(w, r, manager, cfg, nodeUUID)
			case "disable":
				handleDisableNode(w, r, manager, cfg, nodeUUID)
			case "restart":
				handleRestartNode(w, r, manager, cfg, nodeUUID)
			case "reset-traffic":
				handleResetNodeTraffic(w, r, manager, cfg, nodeUUID)
			default:
				http.NotFound(w, r)
			}
			return
		}

		if len(parts) != 1 {
			http.NotFound(w, r)
			return
		}

		switch r.Method {
		case http.MethodGet:
			handleGetNode(w, r, manager, cfg, nodeUUID)
		case http.MethodDelete:
			handleDeleteNode(w, r, manager, cfg, nodeUUID)
		default:
			shared.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
		}
	}
}

func NodesActionsHandler(manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			shared.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}
		path := trimNodesPath(r.URL.Path, "/actions/")
		switch path {
		case "restart-all":
			handleRestartAllNodes(w, r, manager, cfg)
		case "reorder":
			handleReorderNodes(w, r, manager, cfg)
		default:
			http.NotFound(w, r)
		}
	}
}

func NodesBulkActionsHandler(manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			shared.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}

		path := trimNodesPath(r.URL.Path, "/bulk-actions")
		switch path {
		case "":
			handleBulkNodesActions(w, r, manager, cfg)
		case "profile-modification":
			handleBulkProfileModification(w, r, manager, cfg)
		default:
			http.NotFound(w, r)
		}
	}
}

func NodesTagsHandler(manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			shared.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}
		tags, err := getNodeTags(r.Context(), manager)
		if err != nil {
			shared.SendError(w, http.StatusInternalServerError, "failed to fetch node tags", err, cfg)
			return
		}
		shared.WriteJSON(w, http.StatusOK, map[string]any{
			"response": map[string]any{
				"tags": tags,
			},
		})
	}
}

func trimNodesPath(path string, suffix string) string {
	for _, prefix := range []string{"/api/nodes"} {
		if strings.HasPrefix(path, prefix+suffix) {
			return strings.Trim(strings.TrimPrefix(path, prefix+suffix), "/")
		}
	}
	return strings.Trim(path, "/")
}

func handleGetNodes(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) {
	nodes, err := getAllNodeRecords(r.Context(), manager)
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to fetch nodes", err, cfg)
		return
	}

	response, err := buildNodeResponses(r.Context(), manager, nodes)
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to build node response", err, cfg)
		return
	}

	shared.WriteJSON(w, http.StatusOK, map[string]any{"response": response})
}

func handleGetNode(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig, nodeUUID string) {
	node, err := getNodeByUUID(r.Context(), manager, nodeUUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			shared.SendError(w, http.StatusNotFound, "node not found", nil, cfg)
			return
		}
		shared.SendError(w, http.StatusInternalServerError, "failed to fetch node", err, cfg)
		return
	}

	response, err := buildNodeResponses(r.Context(), manager, []nodeRecord{node})
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to build node response", err, cfg)
		return
	}

	shared.WriteJSON(w, http.StatusOK, map[string]any{"response": response[0]})
}

func handleCreateNode(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) {
	var req createNodeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		shared.SendError(w, http.StatusBadRequest, "invalid JSON", err, cfg)
		return
	}
	if err := validateCreateRequest(req); err != nil {
		shared.SendError(w, http.StatusBadRequest, err.Error(), nil, cfg)
		return
	}
	if err := ensureConfigProfileInbounds(r.Context(), manager, req.ConfigProfile.ActiveConfigProfileUUID, req.ConfigProfile.ActiveInbounds); err != nil {
		handleConfigProfileValidationError(w, err, cfg)
		return
	}
	if req.ProviderUUID != nil && strings.TrimSpace(*req.ProviderUUID) != "" {
		if _, err := uuid.Parse(strings.TrimSpace(*req.ProviderUUID)); err != nil {
			shared.SendError(w, http.StatusBadRequest, "invalid providerUuid", nil, cfg)
			return
		}
	}

	nodeUUID := uuid.NewString()
	now := time.Now().UTC()
	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		tx, err := db.BeginTx(r.Context(), nil)
		if err != nil {
			return err
		}

		_, err = tx.ExecContext(r.Context(), `
			INSERT INTO nodes (
				uuid, name, address, port, api_schema, api_path, active_config_profile_uuid,
				is_connected, is_connecting, is_disabled, last_status_change, last_status_message,
				singbox_version, node_version, singbox_uptime, users_online, consumption_multiplier,
				is_traffic_tracking_active, traffic_reset_day, traffic_limit_bytes, traffic_used_bytes,
				notify_percent, provider_uuid, country_code, tags, created_at, updated_at
			) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		`,
			nodeUUID,
			strings.TrimSpace(req.Name),
			strings.TrimSpace(req.Address),
			req.Port,
			normalizeAPISchema(req.APISchema),
			normalizeAPIPath(req.APIPath),
			req.ConfigProfile.ActiveConfigProfileUUID,
			false,
			false,
			false,
			nil,
			nil,
			nil,
			nil,
			"0",
			0,
			toNanoMultiplier(coalesceFloat(req.ConsumptionMultiplier, 1)),
			coalesceBool(req.IsTrafficTrackingActive, false),
			coalesceInt(req.TrafficResetDay, 1),
			coalesceInt64(req.TrafficLimitBytes, 0),
			0,
			coalesceInt(req.NotifyPercent, 0),
			normalizeNullableString(req.ProviderUUID),
			normalizeCountryCode(req.CountryCode),
			normalizeTags(req.Tags),
			now,
			now,
		)
		if err != nil {
			_ = tx.Rollback()
			return err
		}

		if err := replaceNodeInboundsTx(r.Context(), tx, nodeUUID, req.ConfigProfile.ActiveInbounds); err != nil {
			_ = tx.Rollback()
			return err
		}

		return tx.Commit()
	})
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to create node", err, cfg)
		return
	}

	monitor.RequestNodeSync()
	monitor.RequestNodeDeploy(true, nodeUUID)
	node, err := getNodeByUUID(r.Context(), manager, nodeUUID)
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to fetch created node", err, cfg)
		return
	}
	response, err := buildNodeResponses(r.Context(), manager, []nodeRecord{node})
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to build node response", err, cfg)
		return
	}
	emitNodeNotification(r.Context(), cfg, notifications.EventNodeCreated, node, nil)
	shared.WriteJSON(w, http.StatusCreated, map[string]any{"response": response[0]})
}

func handleUpdateNode(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) {
	var req updateNodeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		shared.SendError(w, http.StatusBadRequest, "invalid JSON", err, cfg)
		return
	}
	if _, err := uuid.Parse(req.UUID); err != nil {
		shared.SendError(w, http.StatusBadRequest, "invalid UUID format", nil, cfg)
		return
	}
	if err := validateUpdateRequest(req); err != nil {
		shared.SendError(w, http.StatusBadRequest, err.Error(), nil, cfg)
		return
	}
	if req.ConfigProfile != nil {
		if err := ensureConfigProfileInbounds(r.Context(), manager, req.ConfigProfile.ActiveConfigProfileUUID, req.ConfigProfile.ActiveInbounds); err != nil {
			handleConfigProfileValidationError(w, err, cfg)
			return
		}
	}
	if req.ProviderUUID.Set && req.ProviderUUID.Value != nil && strings.TrimSpace(*req.ProviderUUID.Value) != "" {
		if _, err := uuid.Parse(strings.TrimSpace(*req.ProviderUUID.Value)); err != nil {
			shared.SendError(w, http.StatusBadRequest, "invalid providerUuid", nil, cfg)
			return
		}
	}

	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		tx, err := db.BeginTx(r.Context(), nil)
		if err != nil {
			return err
		}

		clauses := make([]string, 0)
		args := make([]any, 0)
		add := func(column string, value any) {
			clauses = append(clauses, fmt.Sprintf("%s = ?", column))
			args = append(args, value)
		}

		if req.Name != nil {
			add("name", strings.TrimSpace(*req.Name))
		}
		if req.Address != nil {
			add("address", strings.TrimSpace(*req.Address))
		}
		if req.Port != nil {
			add("port", *req.Port)
		}
		if req.APISchema != nil {
			add("api_schema", normalizeAPISchema(req.APISchema))
		}
		if req.APIPath != nil {
			add("api_path", normalizeAPIPath(req.APIPath))
		}
		if req.IsTrafficTrackingActive != nil {
			add("is_traffic_tracking_active", *req.IsTrafficTrackingActive)
		}
		if req.TrafficLimitBytes != nil {
			add("traffic_limit_bytes", *req.TrafficLimitBytes)
		}
		if req.NotifyPercent != nil {
			add("notify_percent", *req.NotifyPercent)
		}
		if req.TrafficResetDay != nil {
			add("traffic_reset_day", *req.TrafficResetDay)
		}
		if req.CountryCode != nil {
			add("country_code", strings.ToUpper(strings.TrimSpace(*req.CountryCode)))
		}
		if req.ConsumptionMultiplier != nil {
			add("consumption_multiplier", toNanoMultiplier(*req.ConsumptionMultiplier))
		}
		if req.Tags != nil {
			add("tags", normalizeTags(*req.Tags))
		}
		if req.ProviderUUID.Set {
			if req.ProviderUUID.Value == nil || strings.TrimSpace(*req.ProviderUUID.Value) == "" {
				clauses = append(clauses, "provider_uuid = NULL")
			} else {
				add("provider_uuid", strings.TrimSpace(*req.ProviderUUID.Value))
			}
		}
		if req.ConfigProfile != nil {
			add("active_config_profile_uuid", req.ConfigProfile.ActiveConfigProfileUUID)
		}

		if len(clauses) > 0 {
			args = append(args, req.UUID)
			query := fmt.Sprintf("UPDATE nodes SET %s, updated_at = CURRENT_TIMESTAMP WHERE uuid = ?", strings.Join(clauses, ", "))
			result, err := tx.ExecContext(r.Context(), query, args...)
			if err != nil {
				_ = tx.Rollback()
				return err
			}
			rows, err := result.RowsAffected()
			if err != nil {
				_ = tx.Rollback()
				return err
			}
			if rows == 0 {
				_ = tx.Rollback()
				return sql.ErrNoRows
			}
		}

		if req.ConfigProfile != nil {
			if err := replaceNodeInboundsTx(r.Context(), tx, req.UUID, req.ConfigProfile.ActiveInbounds); err != nil {
				_ = tx.Rollback()
				return err
			}
		}

		return tx.Commit()
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			shared.SendError(w, http.StatusNotFound, "node not found", nil, cfg)
			return
		}
		shared.SendError(w, http.StatusInternalServerError, "failed to update node", err, cfg)
		return
	}

	monitor.RequestNodeSync()
	monitor.RequestNodeDeploy(true, req.UUID)
	node, err := getNodeByUUID(r.Context(), manager, req.UUID)
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to fetch updated node", err, cfg)
		return
	}
	response, err := buildNodeResponses(r.Context(), manager, []nodeRecord{node})
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to build node response", err, cfg)
		return
	}
	emitNodeNotification(r.Context(), cfg, notifications.EventNodeModified, node, nil)
	shared.WriteJSON(w, http.StatusOK, map[string]any{"response": response[0]})
}

func handleDeleteNode(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig, nodeUUID string) {
	node, nodeErr := getNodeByUUID(r.Context(), manager, nodeUUID)
	if nodeErr != nil && !errors.Is(nodeErr, sql.ErrNoRows) {
		shared.SendError(w, http.StatusInternalServerError, "failed to fetch node", nodeErr, cfg)
		return
	}
	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		result, err := db.ExecContext(r.Context(), `DELETE FROM nodes WHERE uuid = ?`, nodeUUID)
		if err != nil {
			return err
		}
		rows, err := result.RowsAffected()
		if err != nil {
			return err
		}
		if rows == 0 {
			return sql.ErrNoRows
		}
		return nil
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			shared.SendError(w, http.StatusNotFound, "node not found", nil, cfg)
			return
		}
		shared.SendError(w, http.StatusInternalServerError, "failed to delete node", err, cfg)
		return
	}

	monitor.RequestNodeSync()
	monitor.RequestNodeDeploy(true, nodeUUID)
	if nodeErr == nil {
		emitNodeNotification(r.Context(), cfg, notifications.EventNodeDeleted, node, nil)
	}
	shared.WriteJSON(w, http.StatusOK, map[string]any{"response": map[string]any{"isDeleted": true}})
}

func handleEnableNode(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig, nodeUUID string) {
	node, err := getNodeByUUID(r.Context(), manager, nodeUUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			shared.SendError(w, http.StatusNotFound, "node not found", nil, cfg)
			return
		}
		shared.SendError(w, http.StatusInternalServerError, "failed to fetch node", err, cfg)
		return
	}

	inboundsMap, err := getNodeInbounds(r.Context(), manager, []string{nodeUUID})
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to fetch node inbounds", err, cfg)
		return
	}

	err = manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		if node.ActiveConfigProfileUUID == nil || len(inboundsMap[nodeUUID]) == 0 {
			_, execErr := db.ExecContext(r.Context(), `
				UPDATE nodes
				SET is_disabled = true, active_config_profile_uuid = NULL, is_connecting = false,
					is_connected = false, last_status_message = NULL, last_status_change = ?, users_online = 0
				WHERE uuid = ?
			`, time.Now().UTC(), nodeUUID)
			return execErr
		}
		_, execErr := db.ExecContext(r.Context(), `UPDATE nodes SET is_disabled = false, updated_at = CURRENT_TIMESTAMP WHERE uuid = ?`, nodeUUID)
		return execErr
	})
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to enable node", err, cfg)
		return
	}

	monitor.RequestNodeSync()
	monitor.RequestNodeDeploy(true, nodeUUID)
	if updated, loadErr := getNodeByUUID(r.Context(), manager, nodeUUID); loadErr == nil {
		emitNodeNotification(r.Context(), cfg, notifications.EventNodeEnabled, updated, nil)
	}
	sendUpdatedNodeResponse(w, r, manager, cfg, nodeUUID)
}

func handleDisableNode(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig, nodeUUID string) {
	node, err := getNodeByUUID(r.Context(), manager, nodeUUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			shared.SendError(w, http.StatusNotFound, "node not found", nil, cfg)
			return
		}
		shared.SendError(w, http.StatusInternalServerError, "failed to fetch node", err, cfg)
		return
	}

	inboundsMap, err := getNodeInbounds(r.Context(), manager, []string{nodeUUID})
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to fetch node inbounds", err, cfg)
		return
	}

	err = manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		if node.ActiveConfigProfileUUID == nil || len(inboundsMap[nodeUUID]) == 0 {
			if _, execErr := db.ExecContext(r.Context(), `UPDATE nodes SET active_config_profile_uuid = NULL WHERE uuid = ?`, nodeUUID); execErr != nil {
				return execErr
			}
		}
		_, execErr := db.ExecContext(r.Context(), `
			UPDATE nodes
			SET is_disabled = true, is_connecting = false, is_connected = false,
				last_status_message = NULL, last_status_change = ?, users_online = 0,
				updated_at = CURRENT_TIMESTAMP
			WHERE uuid = ?
		`, time.Now().UTC(), nodeUUID)
		return execErr
	})
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to disable node", err, cfg)
		return
	}

	monitor.RequestNodeSync()
	monitor.RequestNodeDeploy(true, nodeUUID)
	if updated, loadErr := getNodeByUUID(r.Context(), manager, nodeUUID); loadErr == nil {
		emitNodeNotification(r.Context(), cfg, notifications.EventNodeDisabled, updated, nil)
	}
	sendUpdatedNodeResponse(w, r, manager, cfg, nodeUUID)
}

func handleRestartNode(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig, nodeUUID string) {
	req, err := decodeOptionalRestartNodesRequest(r)
	if err != nil {
		shared.SendError(w, http.StatusBadRequest, "invalid JSON", err, cfg)
		return
	}

	node, err := getNodeByUUID(r.Context(), manager, nodeUUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			shared.SendError(w, http.StatusNotFound, "node not found", nil, cfg)
			return
		}
		shared.SendError(w, http.StatusInternalServerError, "failed to fetch node", err, cfg)
		return
	}
	if node.IsDisabled {
		shared.SendError(w, http.StatusBadRequest, "node is disabled", nil, cfg)
		return
	}

	monitor.RequestNodeSync()
	monitor.RequestNodeDeployWithForce(true, isForceRestartRequested(req), nodeUUID)
	shared.WriteJSON(w, http.StatusOK, map[string]any{"response": map[string]any{"eventSent": true}})
}

func handleResetNodeTraffic(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig, nodeUUID string) {
	node, err := getNodeByUUID(r.Context(), manager, nodeUUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			shared.SendError(w, http.StatusNotFound, "node not found", nil, cfg)
			return
		}
		shared.SendError(w, http.StatusInternalServerError, "failed to fetch node", err, cfg)
		return
	}

	err = manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		tx, err := db.BeginTx(r.Context(), nil)
		if err != nil {
			return err
		}
		_, err = tx.ExecContext(r.Context(), `
			INSERT INTO nodes_traffic_usage_history (node_uuid, traffic_bytes, reset_at)
			VALUES (?, ?, ?)
		`, nodeUUID, coalesceInt64Ptr(node.TrafficUsedBytes), time.Now().UTC())
		if err != nil {
			_ = tx.Rollback()
			return err
		}
		_, err = tx.ExecContext(r.Context(), `
			UPDATE nodes SET traffic_used_bytes = 0, updated_at = CURRENT_TIMESTAMP WHERE uuid = ?
		`, nodeUUID)
		if err != nil {
			_ = tx.Rollback()
			return err
		}
		return tx.Commit()
	})
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to reset node traffic", err, cfg)
		return
	}

	shared.WriteJSON(w, http.StatusOK, map[string]any{"response": map[string]any{"eventSent": true}})
}

func handleRestartAllNodes(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) {
	req, err := decodeOptionalRestartNodesRequest(r)
	if err != nil {
		shared.SendError(w, http.StatusBadRequest, "invalid JSON", err, cfg)
		return
	}

	var enabledCount int
	err = manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		return db.QueryRowContext(r.Context(), `SELECT COUNT(*) FROM nodes WHERE is_disabled = false`).Scan(&enabledCount)
	})
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to inspect nodes", err, cfg)
		return
	}
	if enabledCount == 0 {
		shared.SendError(w, http.StatusBadRequest, errNoEnabledNodes.Error(), nil, cfg)
		return
	}

	monitor.RequestNodeSync()
	monitor.RequestNodeDeployWithForce(true, isForceRestartRequested(req))
	shared.WriteJSON(w, http.StatusOK, map[string]any{"response": map[string]any{"eventSent": true}})
}

func decodeOptionalRestartNodesRequest(r *http.Request) (restartAllNodesRequest, error) {
	var req restartAllNodesRequest
	if r.Body == nil || r.ContentLength == 0 {
		return req, nil
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		if errors.Is(err, io.EOF) {
			return req, nil
		}
		return req, err
	}
	return req, nil
}

func isForceRestartRequested(req restartAllNodesRequest) bool {
	return req.ForceRestart != nil && *req.ForceRestart
}

func handleReorderNodes(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) {
	var req reorderNodesRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		shared.SendError(w, http.StatusBadRequest, "invalid JSON", err, cfg)
		return
	}
	if len(req.Nodes) == 0 {
		shared.SendError(w, http.StatusBadRequest, "nodes cannot be empty", nil, cfg)
		return
	}
	for _, item := range req.Nodes {
		if _, err := uuid.Parse(item.UUID); err != nil {
			shared.SendError(w, http.StatusBadRequest, "invalid UUID format", nil, cfg)
			return
		}
	}

	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		tx, err := db.BeginTx(r.Context(), nil)
		if err != nil {
			return err
		}
		for _, item := range req.Nodes {
			if _, err := tx.ExecContext(r.Context(), `UPDATE nodes SET view_position = ? WHERE uuid = ?`, item.ViewPosition, item.UUID); err != nil {
				_ = tx.Rollback()
				return err
			}
		}
		if _, err := tx.ExecContext(r.Context(), `SELECT setval('nodes_view_position_seq', (SELECT COALESCE(MAX(view_position), 0) FROM nodes) + 1)`); err != nil {
			_ = tx.Rollback()
			return err
		}
		return tx.Commit()
	})
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to reorder nodes", err, cfg)
		return
	}

	handleGetNodes(w, r, manager, cfg)
}

func handleBulkProfileModification(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) {
	var req bulkProfileModificationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		shared.SendError(w, http.StatusBadRequest, "invalid JSON", err, cfg)
		return
	}
	if err := validateUUIDs(req.UUIDs); err != nil {
		shared.SendError(w, http.StatusBadRequest, err.Error(), nil, cfg)
		return
	}
	if err := ensureConfigProfileInbounds(r.Context(), manager, req.ConfigProfile.ActiveConfigProfileUUID, req.ConfigProfile.ActiveInbounds); err != nil {
		handleConfigProfileValidationError(w, err, cfg)
		return
	}

	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		tx, err := db.BeginTx(r.Context(), nil)
		if err != nil {
			return err
		}
		for _, nodeUUID := range req.UUIDs {
			if _, err := tx.ExecContext(r.Context(), `UPDATE nodes SET active_config_profile_uuid = ?, updated_at = CURRENT_TIMESTAMP WHERE uuid = ?`, req.ConfigProfile.ActiveConfigProfileUUID, nodeUUID); err != nil {
				_ = tx.Rollback()
				return err
			}
			if err := replaceNodeInboundsTx(r.Context(), tx, nodeUUID, req.ConfigProfile.ActiveInbounds); err != nil {
				_ = tx.Rollback()
				return err
			}
		}
		return tx.Commit()
	})
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to modify nodes profile", err, cfg)
		return
	}

	monitor.RequestNodeSync()
	monitor.RequestNodeDeploy(true, req.UUIDs...)
	shared.WriteJSON(w, http.StatusOK, map[string]any{"response": map[string]any{"eventSent": true}})
}

func handleBulkNodesActions(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) {
	var req bulkNodesActionsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		shared.SendError(w, http.StatusBadRequest, "invalid JSON", err, cfg)
		return
	}
	if err := validateUUIDs(req.UUIDs); err != nil {
		shared.SendError(w, http.StatusBadRequest, err.Error(), nil, cfg)
		return
	}

	for _, nodeUUID := range req.UUIDs {
		switch req.Action {
		case "ENABLE":
			if err := performEnableAction(r.Context(), manager, nodeUUID); err != nil {
				shared.SendError(w, http.StatusInternalServerError, "failed to enable nodes", err, cfg)
				return
			}
		case "DISABLE":
			if err := performDisableAction(r.Context(), manager, nodeUUID); err != nil {
				shared.SendError(w, http.StatusInternalServerError, "failed to disable nodes", err, cfg)
				return
			}
		case "RESTART":
			if _, err := getNodeByUUID(r.Context(), manager, nodeUUID); err != nil {
				shared.SendError(w, http.StatusInternalServerError, "failed to restart nodes", err, cfg)
				return
			}
		case "RESET_TRAFFIC":
			if err := performResetTrafficAction(r.Context(), manager, nodeUUID); err != nil {
				shared.SendError(w, http.StatusInternalServerError, "failed to reset nodes traffic", err, cfg)
				return
			}
		default:
			shared.SendError(w, http.StatusBadRequest, "invalid bulk action", nil, cfg)
			return
		}
	}

	monitor.RequestNodeSync()
	monitor.RequestNodeDeploy(true, req.UUIDs...)
	switch req.Action {
	case "ENABLE":
		emitNodesByUUIDsNotification(r.Context(), manager, cfg, notifications.EventNodeEnabled, req.UUIDs)
	case "DISABLE":
		emitNodesByUUIDsNotification(r.Context(), manager, cfg, notifications.EventNodeDisabled, req.UUIDs)
	case "RESET_TRAFFIC":
		emitNodesByUUIDsNotification(r.Context(), manager, cfg, notifications.EventNodeModified, req.UUIDs)
	}
	shared.WriteJSON(w, http.StatusOK, map[string]any{"response": map[string]any{"eventSent": true}})
}

func performEnableAction(ctx context.Context, manager *dbmanager.DatabaseManager, nodeUUID string) error {
	node, err := getNodeByUUID(ctx, manager, nodeUUID)
	if err != nil {
		return err
	}
	inboundsMap, err := getNodeInbounds(ctx, manager, []string{nodeUUID})
	if err != nil {
		return err
	}
	return manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		if node.ActiveConfigProfileUUID == nil || len(inboundsMap[nodeUUID]) == 0 {
			_, execErr := db.ExecContext(ctx, `
				UPDATE nodes
				SET is_disabled = true, active_config_profile_uuid = NULL, is_connecting = false,
					is_connected = false, last_status_message = NULL, last_status_change = ?, users_online = 0
				WHERE uuid = ?
			`, time.Now().UTC(), nodeUUID)
			return execErr
		}
		_, execErr := db.ExecContext(ctx, `UPDATE nodes SET is_disabled = false, updated_at = CURRENT_TIMESTAMP WHERE uuid = ?`, nodeUUID)
		return execErr
	})
}

func performDisableAction(ctx context.Context, manager *dbmanager.DatabaseManager, nodeUUID string) error {
	node, err := getNodeByUUID(ctx, manager, nodeUUID)
	if err != nil {
		return err
	}
	inboundsMap, err := getNodeInbounds(ctx, manager, []string{nodeUUID})
	if err != nil {
		return err
	}
	return manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		if node.ActiveConfigProfileUUID == nil || len(inboundsMap[nodeUUID]) == 0 {
			if _, execErr := db.ExecContext(ctx, `UPDATE nodes SET active_config_profile_uuid = NULL WHERE uuid = ?`, nodeUUID); execErr != nil {
				return execErr
			}
		}
		_, execErr := db.ExecContext(ctx, `
			UPDATE nodes
			SET is_disabled = true, is_connecting = false, is_connected = false,
				last_status_message = NULL, last_status_change = ?, users_online = 0,
				updated_at = CURRENT_TIMESTAMP
			WHERE uuid = ?
		`, time.Now().UTC(), nodeUUID)
		return execErr
	})
}

func performResetTrafficAction(ctx context.Context, manager *dbmanager.DatabaseManager, nodeUUID string) error {
	node, err := getNodeByUUID(ctx, manager, nodeUUID)
	if err != nil {
		return err
	}
	return manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		tx, err := db.BeginTx(ctx, nil)
		if err != nil {
			return err
		}
		_, err = tx.ExecContext(ctx, `
			INSERT INTO nodes_traffic_usage_history (node_uuid, traffic_bytes, reset_at)
			VALUES (?, ?, ?)
		`, nodeUUID, coalesceInt64Ptr(node.TrafficUsedBytes), time.Now().UTC())
		if err != nil {
			_ = tx.Rollback()
			return err
		}
		if _, err := tx.ExecContext(ctx, `UPDATE nodes SET traffic_used_bytes = 0, updated_at = CURRENT_TIMESTAMP WHERE uuid = ?`, nodeUUID); err != nil {
			_ = tx.Rollback()
			return err
		}
		return tx.Commit()
	})
}

func sendUpdatedNodeResponse(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig, nodeUUID string) {
	node, err := getNodeByUUID(r.Context(), manager, nodeUUID)
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to fetch updated node", err, cfg)
		return
	}
	response, err := buildNodeResponses(r.Context(), manager, []nodeRecord{node})
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to build node response", err, cfg)
		return
	}
	shared.WriteJSON(w, http.StatusOK, map[string]any{"response": response[0]})
}

func validateCreateRequest(req createNodeRequest) error {
	if len(strings.TrimSpace(req.Name)) < 3 || len(strings.TrimSpace(req.Name)) > 30 {
		return fmt.Errorf("name must be between 3 and 30 characters")
	}
	if len(strings.TrimSpace(req.Address)) < 2 {
		return fmt.Errorf("address must be at least 2 characters")
	}
	if req.Port != nil && (*req.Port < 1 || *req.Port > 65535) {
		return fmt.Errorf("invalid port")
	}
	if req.APISchema != nil && strings.TrimSpace(*req.APISchema) == "" {
		return fmt.Errorf("apiSchema cannot be empty")
	}
	if req.APISchema != nil && !isAllowedNodeAPISchema(*req.APISchema) {
		return fmt.Errorf("apiSchema must be one of: mtls, tls")
	}
	if req.TrafficLimitBytes != nil && *req.TrafficLimitBytes < 0 {
		return fmt.Errorf("trafficLimitBytes must be greater than or equal to 0")
	}
	if req.NotifyPercent != nil && (*req.NotifyPercent < 0 || *req.NotifyPercent > 100) {
		return fmt.Errorf("notifyPercent must be between 0 and 100")
	}
	if req.TrafficResetDay != nil && (*req.TrafficResetDay < 1 || *req.TrafficResetDay > 31) {
		return fmt.Errorf("trafficResetDay must be between 1 and 31")
	}
	if req.ConsumptionMultiplier != nil && (*req.ConsumptionMultiplier < 0 || *req.ConsumptionMultiplier > 100) {
		return fmt.Errorf("consumptionMultiplier must be between 0 and 100")
	}
	if _, err := uuid.Parse(req.ConfigProfile.ActiveConfigProfileUUID); err != nil {
		return fmt.Errorf("invalid activeConfigProfileUuid")
	}
	if err := validateUUIDs(req.ConfigProfile.ActiveInbounds); err != nil {
		return err
	}
	if err := validateTags(req.Tags); err != nil {
		return err
	}
	if req.CountryCode != nil {
		code := strings.TrimSpace(*req.CountryCode)
		if code != "" && len(code) != 2 {
			return fmt.Errorf("countryCode must be 2 characters")
		}
	}
	return nil
}

func validateUpdateRequest(req updateNodeRequest) error {
	if req.Name != nil {
		name := strings.TrimSpace(*req.Name)
		if len(name) < 3 || len(name) > 30 {
			return fmt.Errorf("name must be between 3 and 30 characters")
		}
	}
	if req.Address != nil && len(strings.TrimSpace(*req.Address)) < 2 {
		return fmt.Errorf("address must be at least 2 characters")
	}
	if req.Port != nil && (*req.Port < 1 || *req.Port > 65535) {
		return fmt.Errorf("invalid port")
	}
	if req.APISchema != nil && strings.TrimSpace(*req.APISchema) == "" {
		return fmt.Errorf("apiSchema cannot be empty")
	}
	if req.APISchema != nil && !isAllowedNodeAPISchema(*req.APISchema) {
		return fmt.Errorf("apiSchema must be one of: mtls, tls")
	}
	if req.TrafficLimitBytes != nil && *req.TrafficLimitBytes < 0 {
		return fmt.Errorf("trafficLimitBytes must be greater than or equal to 0")
	}
	if req.NotifyPercent != nil && (*req.NotifyPercent < 0 || *req.NotifyPercent > 100) {
		return fmt.Errorf("notifyPercent must be between 0 and 100")
	}
	if req.TrafficResetDay != nil && (*req.TrafficResetDay < 1 || *req.TrafficResetDay > 31) {
		return fmt.Errorf("trafficResetDay must be between 1 and 31")
	}
	if req.ConsumptionMultiplier != nil && (*req.ConsumptionMultiplier < 0 || *req.ConsumptionMultiplier > 100) {
		return fmt.Errorf("consumptionMultiplier must be between 0 and 100")
	}
	if req.ConfigProfile != nil {
		if _, err := uuid.Parse(req.ConfigProfile.ActiveConfigProfileUUID); err != nil {
			return fmt.Errorf("invalid activeConfigProfileUuid")
		}
		if err := validateUUIDs(req.ConfigProfile.ActiveInbounds); err != nil {
			return err
		}
	}
	if req.Tags != nil {
		if err := validateTags(*req.Tags); err != nil {
			return err
		}
	}
	if req.CountryCode != nil {
		code := strings.TrimSpace(*req.CountryCode)
		if code != "" && len(code) != 2 {
			return fmt.Errorf("countryCode must be 2 characters")
		}
	}
	return nil
}

func validateUUIDs(values []string) error {
	if len(values) == 0 {
		return fmt.Errorf("uuids cannot be empty")
	}
	for _, value := range values {
		if _, err := uuid.Parse(value); err != nil {
			return fmt.Errorf("invalid uuid value")
		}
	}
	return nil
}

func validateTags(tags []string) error {
	if len(tags) > 10 {
		return fmt.Errorf("maximum 10 tags")
	}
	for _, tag := range tags {
		if !nodeTagRegex.MatchString(tag) {
			return fmt.Errorf("tag can only contain uppercase letters, numbers, underscores and colons")
		}
		if len(tag) > 36 {
			return fmt.Errorf("each tag must be less than 36 characters")
		}
	}
	return nil
}

func handleConfigProfileValidationError(w http.ResponseWriter, err error, cfg *config.BackendConfig) {
	switch {
	case errors.Is(err, errConfigProfileNotFound):
		shared.SendError(w, http.StatusBadRequest, "config profile not found", nil, cfg)
	case errors.Is(err, errConfigProfileInboundInvalid):
		shared.SendError(w, http.StatusBadRequest, "config profile inbound not found in specified profile", nil, cfg)
	default:
		shared.SendError(w, http.StatusInternalServerError, "failed to validate config profile inbounds", err, cfg)
	}
}

func ensureConfigProfileInbounds(ctx context.Context, manager *dbmanager.DatabaseManager, profileUUID string, inboundUUIDs []string) error {
	return manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		var exists int
		if err := db.QueryRowContext(ctx, `SELECT 1 FROM config_profiles WHERE uuid = ?`, profileUUID).Scan(&exists); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return errConfigProfileNotFound
			}
			return err
		}

		found := make(map[string]struct{}, len(inboundUUIDs))
		rows, err := db.QueryContext(ctx, `
			SELECT uuid
			FROM config_profile_inbounds
			WHERE profile_uuid = ? AND uuid = ANY(?)
		`, profileUUID, inboundUUIDs)
		if err != nil {
			return err
		}
		defer rows.Close()
		for rows.Next() {
			var inboundUUID string
			if err := rows.Scan(&inboundUUID); err != nil {
				return err
			}
			found[inboundUUID] = struct{}{}
		}
		if err := rows.Err(); err != nil {
			return err
		}
		for _, inboundUUID := range inboundUUIDs {
			if _, ok := found[inboundUUID]; !ok {
				return errConfigProfileInboundInvalid
			}
		}
		return nil
	})
}

func buildNodeResponses(ctx context.Context, manager *dbmanager.DatabaseManager, records []nodeRecord) ([]nodeAPI, error) {
	nodeUUIDs := make([]string, 0, len(records))
	providerUUIDs := make([]string, 0, len(records))
	for _, record := range records {
		nodeUUIDs = append(nodeUUIDs, record.UUID)
		if record.ProviderUUID != nil && *record.ProviderUUID != "" {
			providerUUIDs = append(providerUUIDs, *record.ProviderUUID)
		}
	}

	inboundsMap, err := getNodeInbounds(ctx, manager, nodeUUIDs)
	if err != nil {
		return nil, err
	}
	providersMap, err := getProviders(ctx, manager, dedupeStrings(providerUUIDs))
	if err != nil {
		return nil, err
	}

	response := make([]nodeAPI, 0, len(records))
	for _, record := range records {
		var item nodeAPI
		item.UUID = record.UUID
		item.Name = record.Name
		item.Address = record.Address
		item.Port = record.Port
		item.APISchema = normalizeAPISchema(&record.APISchema)
		item.APIPath = normalizeAPIPath(&record.APIPath)
		item.IsConnected = record.IsConnected
		item.IsDisabled = record.IsDisabled
		item.IsConnecting = record.IsConnecting
		item.LastStatusChange = record.LastStatusChange
		item.LastStatusMessage = record.LastStatusMessage
		item.SingboxVersion = record.SingboxVersion
		item.NodeVersion = record.NodeVersion
		item.SingboxUptime = record.SingboxUptime
		item.IsTrafficTrackingActive = record.IsTrafficTrackingActive
		item.TrafficResetDay = record.TrafficResetDay
		item.TrafficLimitBytes = record.TrafficLimitBytes
		item.TrafficUsedBytes = record.TrafficUsedBytes
		item.NotifyPercent = record.NotifyPercent
		item.UsersOnline = record.UsersOnline
		item.ViewPosition = record.ViewPosition
		item.CountryCode = record.CountryCode
		item.ConsumptionMultiplier = fromNanoMultiplier(record.ConsumptionMultiplier)
		item.Tags = ensureStringSlice(record.Tags)
		item.CPUCount = record.CPUCount
		item.CPUModel = record.CPUModel
		item.TotalRAM = record.TotalRAM
		item.CreatedAt = record.CreatedAt
		item.UpdatedAt = record.UpdatedAt
		item.ConfigProfile.ActiveConfigProfileUUID = record.ActiveConfigProfileUUID
		item.ConfigProfile.ActiveInbounds = ensureInboundSlice(inboundsMap[record.UUID])
		item.ProviderUUID = record.ProviderUUID
		if record.ProviderUUID != nil {
			item.Provider = providersMap[*record.ProviderUUID]
		}
		response = append(response, item)
	}

	return response, nil
}

func emitNodeNotification(ctx context.Context, cfg *config.BackendConfig, event string, record nodeRecord, meta map[string]any) {
	notifications.Emit(ctx, cfg, notifications.Event{
		Scope: notifications.ScopeNode,
		Event: event,
		Data:  nodeRecordNotificationData(record),
		Meta:  meta,
	})
}

func emitNodesByUUIDsNotification(ctx context.Context, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig, event string, nodeUUIDs []string) {
	clean := dedupeStrings(nodeUUIDs)
	if len(clean) == 0 {
		return
	}
	skipTelegram := len(clean) >= 500
	for _, nodeUUID := range clean {
		meta := map[string]any{"bulk": true}
		if skipTelegram {
			meta["skipTelegramNotification"] = true
		}
		record, err := getNodeByUUID(ctx, manager, nodeUUID)
		if err == nil {
			emitNodeNotification(ctx, cfg, event, record, meta)
			continue
		}
		notifications.Emit(ctx, cfg, notifications.Event{
			Scope: notifications.ScopeNode,
			Event: event,
			Data:  map[string]any{"uuid": nodeUUID},
			Meta:  meta,
		})
	}
}

func nodeRecordNotificationData(record nodeRecord) map[string]any {
	return map[string]any{
		"id":                      record.ID,
		"uuid":                    record.UUID,
		"name":                    record.Name,
		"address":                 record.Address,
		"port":                    record.Port,
		"apiSchema":               record.APISchema,
		"apiPath":                 record.APIPath,
		"activeConfigProfileUuid": record.ActiveConfigProfileUUID,
		"isConnected":             record.IsConnected,
		"isConnecting":            record.IsConnecting,
		"isDisabled":              record.IsDisabled,
		"lastStatusChange":        optionalTimeString(record.LastStatusChange),
		"lastStatusMessage":       record.LastStatusMessage,
		"singboxVersion":          record.SingboxVersion,
		"nodeVersion":             record.NodeVersion,
		"singboxUptime":           record.SingboxUptime,
		"usersOnline":             record.UsersOnline,
		"consumptionMultiplier":   fromNanoMultiplier(record.ConsumptionMultiplier),
		"isTrafficTrackingActive": record.IsTrafficTrackingActive,
		"trafficResetDay":         record.TrafficResetDay,
		"trafficLimitBytes":       record.TrafficLimitBytes,
		"trafficUsedBytes":        record.TrafficUsedBytes,
		"notifyPercent":           record.NotifyPercent,
		"providerUuid":            record.ProviderUUID,
		"viewPosition":            record.ViewPosition,
		"countryCode":             record.CountryCode,
		"tags":                    ensureStringSlice(record.Tags),
		"createdAt":               record.CreatedAt.UTC().Format(time.RFC3339),
		"updatedAt":               record.UpdatedAt.UTC().Format(time.RFC3339),
	}
}

func optionalTimeString(value *time.Time) *string {
	if value == nil {
		return nil
	}
	formatted := value.UTC().Format(time.RFC3339)
	return &formatted
}

func getAllNodeRecords(ctx context.Context, manager *dbmanager.DatabaseManager) ([]nodeRecord, error) {
	var nodes []nodeRecord
	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		rows, err := db.QueryContext(ctx, `
			SELECT
				uuid, id, name, address, port, api_schema, api_path, active_config_profile_uuid,
				is_connected, is_connecting, is_disabled, last_status_change, last_status_message,
				singbox_version, node_version, singbox_uptime, users_online, consumption_multiplier,
				is_traffic_tracking_active, traffic_reset_day, traffic_limit_bytes, traffic_used_bytes,
				notify_percent, provider_uuid, view_position, country_code, tags,
				cpu_count, cpu_model, total_ram, created_at, updated_at
			FROM nodes
			ORDER BY view_position ASC
		`)
		if err != nil {
			return err
		}
		defer rows.Close()
		for rows.Next() {
			node, scanErr := scanNodeRecord(rows)
			if scanErr != nil {
				return scanErr
			}
			nodes = append(nodes, node)
		}
		return rows.Err()
	})
	return nodes, err
}

func getNodeByUUID(ctx context.Context, manager *dbmanager.DatabaseManager, nodeUUID string) (nodeRecord, error) {
	var node nodeRecord
	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		row := db.QueryRowContext(ctx, `
			SELECT
				uuid, id, name, address, port, api_schema, api_path, active_config_profile_uuid,
				is_connected, is_connecting, is_disabled, last_status_change, last_status_message,
				singbox_version, node_version, singbox_uptime, users_online, consumption_multiplier,
				is_traffic_tracking_active, traffic_reset_day, traffic_limit_bytes, traffic_used_bytes,
				notify_percent, provider_uuid, view_position, country_code, tags,
				cpu_count, cpu_model, total_ram, created_at, updated_at
			FROM nodes
			WHERE uuid = ?
		`, nodeUUID)
		var scanErr error
		node, scanErr = scanNodeRecord(row)
		return scanErr
	})
	return node, err
}

func scanNodeRecord(scanner shared.RowScanner) (nodeRecord, error) {
	var node nodeRecord
	var id sql.NullInt64
	var port sql.NullInt64
	var activeConfigProfileUUID sql.NullString
	var lastStatusChange sql.NullTime
	var lastStatusMessage sql.NullString
	var singboxVersion sql.NullString
	var nodeVersion sql.NullString
	var usersOnline sql.NullInt64
	var trafficResetDay sql.NullInt64
	var trafficLimitBytes sql.NullInt64
	var trafficUsedBytes sql.NullInt64
	var notifyPercent sql.NullInt64
	var providerUUID sql.NullString
	var cpuCount sql.NullInt64
	var cpuModel sql.NullString
	var totalRAM sql.NullString
	var tags dbutil.StringArray

	err := scanner.Scan(
		&node.UUID,
		&id,
		&node.Name,
		&node.Address,
		&port,
		&node.APISchema,
		&node.APIPath,
		&activeConfigProfileUUID,
		&node.IsConnected,
		&node.IsConnecting,
		&node.IsDisabled,
		&lastStatusChange,
		&lastStatusMessage,
		&singboxVersion,
		&nodeVersion,
		&node.SingboxUptime,
		&usersOnline,
		&node.ConsumptionMultiplier,
		&node.IsTrafficTrackingActive,
		&trafficResetDay,
		&trafficLimitBytes,
		&trafficUsedBytes,
		&notifyPercent,
		&providerUUID,
		&node.ViewPosition,
		&node.CountryCode,
		&tags,
		&cpuCount,
		&cpuModel,
		&totalRAM,
		&node.CreatedAt,
		&node.UpdatedAt,
	)
	if err != nil {
		return node, err
	}

	if id.Valid {
		node.ID = &id.Int64
	}
	if port.Valid {
		value := int(port.Int64)
		node.Port = &value
	}
	if activeConfigProfileUUID.Valid {
		node.ActiveConfigProfileUUID = &activeConfigProfileUUID.String
	}
	if lastStatusChange.Valid {
		node.LastStatusChange = &lastStatusChange.Time
	}
	if lastStatusMessage.Valid {
		node.LastStatusMessage = &lastStatusMessage.String
	}
	if singboxVersion.Valid {
		node.SingboxVersion = &singboxVersion.String
	}
	if nodeVersion.Valid {
		node.NodeVersion = &nodeVersion.String
	}
	if usersOnline.Valid {
		value := int(usersOnline.Int64)
		node.UsersOnline = &value
	}
	if trafficResetDay.Valid {
		value := int(trafficResetDay.Int64)
		node.TrafficResetDay = &value
	}
	if trafficLimitBytes.Valid {
		node.TrafficLimitBytes = &trafficLimitBytes.Int64
	}
	if trafficUsedBytes.Valid {
		node.TrafficUsedBytes = &trafficUsedBytes.Int64
	}
	if notifyPercent.Valid {
		value := int(notifyPercent.Int64)
		node.NotifyPercent = &value
	}
	if providerUUID.Valid {
		node.ProviderUUID = &providerUUID.String
	}
	if cpuCount.Valid {
		value := int(cpuCount.Int64)
		node.CPUCount = &value
	}
	if cpuModel.Valid {
		node.CPUModel = &cpuModel.String
	}
	if totalRAM.Valid {
		node.TotalRAM = &totalRAM.String
	}
	node.Tags = tags.Slice()

	return node, nil
}

func getNodeInbounds(ctx context.Context, manager *dbmanager.DatabaseManager, nodeUUIDs []string) (map[string][]configProfileInboundResponse, error) {
	result := make(map[string][]configProfileInboundResponse)
	if len(nodeUUIDs) == 0 {
		return result, nil
	}
	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		rows, err := db.QueryContext(ctx, `
			SELECT
				cpitn.node_uuid,
				cpi.uuid, cpi.profile_uuid, cpi.tag, cpi.type, cpi.network, cpi.security, cpi.port, cpi.raw_inbound
			FROM config_profile_inbounds_to_nodes cpitn
			JOIN config_profile_inbounds cpi ON cpi.uuid = cpitn.config_profile_inbound_uuid
			WHERE cpitn.node_uuid = ANY(?)
			ORDER BY cpi.tag ASC
		`, nodeUUIDs)
		if err != nil {
			return err
		}
		defer rows.Close()
		for rows.Next() {
			var nodeUUID string
			var inbound configProfileInboundResponse
			var network sql.NullString
			var security sql.NullString
			var port sql.NullInt64
			var rawInbound []byte
			if err := rows.Scan(
				&nodeUUID,
				&inbound.UUID,
				&inbound.ProfileUUID,
				&inbound.Tag,
				&inbound.Type,
				&network,
				&security,
				&port,
				&rawInbound,
			); err != nil {
				return err
			}
			if network.Valid {
				inbound.Network = &network.String
			}
			if security.Valid {
				inbound.Security = &security.String
			}
			if port.Valid {
				value := int(port.Int64)
				inbound.Port = &value
			}
			inbound.RawInbound = json.RawMessage(rawInbound)
			result[nodeUUID] = append(result[nodeUUID], inbound)
		}
		return rows.Err()
	})
	return result, err
}

func getProviders(ctx context.Context, manager *dbmanager.DatabaseManager, providerUUIDs []string) (map[string]*providerResponse, error) {
	result := make(map[string]*providerResponse)
	if len(providerUUIDs) == 0 {
		return result, nil
	}
	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		rows, err := db.QueryContext(ctx, `
			SELECT uuid, name, favicon_link, login_url, created_at, updated_at
			FROM infra_providers
			WHERE uuid = ANY(?)
		`, providerUUIDs)
		if err != nil {
			return err
		}
		defer rows.Close()
		for rows.Next() {
			var item providerResponse
			var favicon sql.NullString
			var loginURL sql.NullString
			var createdAt time.Time
			var updatedAt time.Time
			if err := rows.Scan(&item.UUID, &item.Name, &favicon, &loginURL, &createdAt, &updatedAt); err != nil {
				return err
			}
			if favicon.Valid {
				item.FaviconLink = &favicon.String
			}
			if loginURL.Valid {
				item.LoginURL = &loginURL.String
			}
			item.CreatedAt = &createdAt
			item.UpdatedAt = &updatedAt
			result[item.UUID] = &item
		}
		return rows.Err()
	})
	return result, err
}

func replaceNodeInboundsTx(ctx context.Context, tx dbmanager.TxExecutor, nodeUUID string, inboundUUIDs []string) error {
	if _, err := tx.ExecContext(ctx, `DELETE FROM config_profile_inbounds_to_nodes WHERE node_uuid = ?`, nodeUUID); err != nil {
		return err
	}
	for _, inboundUUID := range dedupeStrings(inboundUUIDs) {
		if _, err := tx.ExecContext(ctx, `
			INSERT INTO config_profile_inbounds_to_nodes (config_profile_inbound_uuid, node_uuid)
			VALUES (?, ?)
		`, inboundUUID, nodeUUID); err != nil {
			return err
		}
	}
	return nil
}

func getNodeTags(ctx context.Context, manager *dbmanager.DatabaseManager) ([]string, error) {
	tags := make([]string, 0)
	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		rows, err := db.QueryContext(ctx, `SELECT DISTINCT unnest(tags) AS tag FROM nodes ORDER BY tag ASC`)
		if err != nil {
			return err
		}
		defer rows.Close()
		for rows.Next() {
			var tag string
			if err := rows.Scan(&tag); err != nil {
				return err
			}
			if tag != "" {
				tags = append(tags, tag)
			}
		}
		return rows.Err()
	})
	if err != nil {
		return nil, err
	}
	sort.Strings(tags)
	return dedupeStrings(tags), nil
}

func normalizeCountryCode(value *string) string {
	if value == nil || strings.TrimSpace(*value) == "" {
		return "XX"
	}
	return strings.ToUpper(strings.TrimSpace(*value))
}

func normalizeTags(tags []string) []string {
	if tags == nil {
		return []string{}
	}
	normalized := make([]string, 0, len(tags))
	for _, tag := range tags {
		tag = strings.ToUpper(strings.TrimSpace(tag))
		if tag == "" {
			continue
		}
		normalized = append(normalized, tag)
	}
	return dedupeStrings(normalized)
}

func normalizeNullableString(value *string) any {
	if value == nil || strings.TrimSpace(*value) == "" {
		return nil
	}
	return strings.TrimSpace(*value)
}

func coalesceBool(value *bool, fallback bool) bool {
	if value == nil {
		return fallback
	}
	return *value
}

func coalesceInt(value *int, fallback int) int {
	if value == nil {
		return fallback
	}
	return *value
}

func coalesceInt64(value *int64, fallback int64) int64 {
	if value == nil {
		return fallback
	}
	return *value
}

func coalesceFloat(value *float64, fallback float64) float64 {
	if value == nil {
		return fallback
	}
	return *value
}

func coalesceInt64Ptr(value *int64) int64 {
	if value == nil {
		return 0
	}
	return *value
}

func toNanoMultiplier(value float64) int64 {
	return int64(math.Round(value * 1_000_000_000))
}

func normalizeAPISchema(value *string) string {
	if value == nil {
		return "mtls"
	}
	normalized := strings.ToLower(strings.TrimSpace(*value))
	switch normalized {
	case "tls":
		return "tls"
	case "mtls", "grpc", "grpcs", "https", "":
		return "mtls"
	default:
		return "mtls"
	}
}

func normalizeAPIPath(value *string) string {
	if value == nil {
		return "/"
	}
	trimmed := strings.TrimSpace(*value)
	if trimmed == "" || trimmed == "/" {
		return "/"
	}
	return "/" + strings.Trim(trimmed, "/")
}

func isAllowedNodeAPISchema(value string) bool {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "mtls", "tls", "grpc", "grpcs", "https":
		return true
	default:
		return false
	}
}

func fromNanoMultiplier(value int64) float64 {
	return float64(value) / 1_000_000_000
}

func ensureStringSlice(values []string) []string {
	if values == nil {
		return []string{}
	}
	return values
}

func ensureInboundSlice(values []configProfileInboundResponse) []configProfileInboundResponse {
	if values == nil {
		return []configProfileInboundResponse{}
	}
	return values
}

func dedupeStrings(values []string) []string {
	if len(values) < 2 {
		return values
	}
	seen := make(map[string]struct{}, len(values))
	result := make([]string, 0, len(values))
	for _, value := range values {
		if _, ok := seen[value]; ok {
			continue
		}
		seen[value] = struct{}{}
		result = append(result, value)
	}
	return result
}
