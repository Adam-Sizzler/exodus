package subscriptionconnections

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"sort"
	"strings"
	"time"

	"cerberus/backend/config"
	dbmanager "cerberus/backend/db/manager"
	"cerberus/backend/dbutil"
	"cerberus/backend/httpapi/shared"
	monitor "cerberus/backend/subscriptionnodes"

	"github.com/google/uuid"
)

var nodeTagRegex = regexp.MustCompile(`^[A-Z0-9_:]+$`)

var (
	errNoEnabledNodes = errors.New("enabled nodes not found")
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
	SubpageConfigUUID       *string    `json:"subpageConfigUuid,omitempty"`
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
	SubpageConfigUUID       *string
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
	Name              string   `json:"name"`
	Address           string   `json:"address"`
	Port              *int     `json:"port,omitempty"`
	APISchema         *string  `json:"apiSchema,omitempty"`
	APIPath           *string  `json:"apiPath,omitempty"`
	SubpageConfigUUID *string  `json:"subpageConfigUuid,omitempty"`
	ProviderUUID      *string  `json:"providerUuid,omitempty"`
	Tags              []string `json:"tags,omitempty"`
}

type updateNodeRequest struct {
	UUID              string         `json:"uuid"`
	Name              *string        `json:"name,omitempty"`
	Address           *string        `json:"address,omitempty"`
	Port              *int           `json:"port,omitempty"`
	APISchema         *string        `json:"apiSchema,omitempty"`
	APIPath           *string        `json:"apiPath,omitempty"`
	SubpageConfigUUID OptionalString `json:"subpageConfigUuid,omitempty"`
	ProviderUUID      OptionalString `json:"providerUuid,omitempty"`
	Tags              *[]string      `json:"tags,omitempty"`
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
	for _, prefix := range []string{"/api/subscription-connections"} {
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
	if req.ProviderUUID != nil && strings.TrimSpace(*req.ProviderUUID) != "" {
		if _, err := uuid.Parse(strings.TrimSpace(*req.ProviderUUID)); err != nil {
			shared.SendError(w, http.StatusBadRequest, "invalid providerUuid", nil, cfg)
			return
		}
	}
	schema := normalizeAPISchema(req.APISchema)
	subpageConfigUUID := normalizeNullableUUID(req.SubpageConfigUUID)
	if subpageConfigUUID != nil {
		if _, err := uuid.Parse(*subpageConfigUUID); err != nil {
			shared.SendError(w, http.StatusBadRequest, "invalid subpageConfigUuid", nil, cfg)
			return
		}
		exists, err := subpageConfigExists(r.Context(), manager, *subpageConfigUUID)
		if err != nil {
			shared.SendError(w, http.StatusInternalServerError, "failed to validate subpage config", err, cfg)
			return
		}
		if !exists {
			shared.SendError(w, http.StatusBadRequest, "subpageConfigUuid not found", nil, cfg)
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
			INSERT INTO sub_nodes (
				uuid, name, address, port, api_schema, api_path,
				is_connected, is_connecting, is_disabled, provider_uuid, tags, created_at, updated_at
			) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		`,
			nodeUUID,
			strings.TrimSpace(req.Name),
			strings.TrimSpace(req.Address),
			req.Port,
			schema,
			normalizeAPIPath(req.APIPath),
			false,
			false,
			false,
			normalizeNullableString(req.ProviderUUID),
			normalizeTags(req.Tags),
			now,
			now,
		)
		if err != nil {
			_ = tx.Rollback()
			return err
		}
		if subpageConfigUUID != nil {
			if _, err := tx.ExecContext(r.Context(), `
				INSERT INTO sub_nodes_to_subscription_page_config (node_uuid, subpage_config_uuid)
				VALUES (?, ?)
				ON CONFLICT (node_uuid) DO UPDATE
				SET subpage_config_uuid = EXCLUDED.subpage_config_uuid
			`, nodeUUID, *subpageConfigUUID); err != nil {
				_ = tx.Rollback()
				return err
			}
		}
		return tx.Commit()
	})
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to create node", err, cfg)
		return
	}
	monitor.RequestSubNodeSync()
	if subpageConfigUUID != nil {
		subpageConfig, err := fetchSubpageConfigRaw(r.Context(), manager, *subpageConfigUUID)
		if err != nil {
			cfg.Logger.Warn("Failed to preload subpage config for created subscription node", "node_uuid", nodeUUID, "subpage_config_uuid", *subpageConfigUUID, "error", err)
		} else {
			monitor.RequestSubNodeSubpageConfigPush(*subpageConfigUUID, subpageConfig, nodeUUID)
		}
	}

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
	if req.SubpageConfigUUID.Set && req.SubpageConfigUUID.Value != nil && strings.TrimSpace(*req.SubpageConfigUUID.Value) != "" {
		if _, err := uuid.Parse(strings.TrimSpace(*req.SubpageConfigUUID.Value)); err != nil {
			shared.SendError(w, http.StatusBadRequest, "invalid subpageConfigUuid", nil, cfg)
			return
		}
	}
	if req.ProviderUUID.Set && req.ProviderUUID.Value != nil && strings.TrimSpace(*req.ProviderUUID.Value) != "" {
		if _, err := uuid.Parse(strings.TrimSpace(*req.ProviderUUID.Value)); err != nil {
			shared.SendError(w, http.StatusBadRequest, "invalid providerUuid", nil, cfg)
			return
		}
	}

	currentNode, err := getNodeByUUID(r.Context(), manager, req.UUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			shared.SendError(w, http.StatusNotFound, "node not found", nil, cfg)
			return
		}
		shared.SendError(w, http.StatusInternalServerError, "failed to fetch node", err, cfg)
		return
	}

	finalSchema := normalizeSubNodeSchema(currentNode.APISchema)
	if req.APISchema != nil {
		finalSchema = normalizeAPISchema(req.APISchema)
	}

	var finalSubpageConfigUUID *string
	if currentNode.SubpageConfigUUID != nil {
		trimmed := strings.TrimSpace(*currentNode.SubpageConfigUUID)
		if trimmed != "" {
			finalSubpageConfigUUID = &trimmed
		}
	}
	if req.SubpageConfigUUID.Set {
		finalSubpageConfigUUID = normalizeNullableUUID(req.SubpageConfigUUID.Value)
	}
	if finalSubpageConfigUUID != nil {
		exists, err := subpageConfigExists(r.Context(), manager, *finalSubpageConfigUUID)
		if err != nil {
			shared.SendError(w, http.StatusInternalServerError, "failed to validate subpage config", err, cfg)
			return
		}
		if !exists {
			shared.SendError(w, http.StatusBadRequest, "subpageConfigUuid not found", nil, cfg)
			return
		}
	}

	err = manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
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
			add("api_schema", finalSchema)
		}
		if req.APIPath != nil {
			add("api_path", normalizeAPIPath(req.APIPath))
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

		if len(clauses) > 0 {
			args = append(args, req.UUID)
			query := fmt.Sprintf("UPDATE sub_nodes SET %s, updated_at = CURRENT_TIMESTAMP WHERE uuid = ?", strings.Join(clauses, ", "))
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
		if req.SubpageConfigUUID.Set {
			if finalSubpageConfigUUID == nil {
				if _, err := tx.ExecContext(r.Context(), `DELETE FROM sub_nodes_to_subscription_page_config WHERE node_uuid = ?`, req.UUID); err != nil {
					_ = tx.Rollback()
					return err
				}
			} else {
				if _, err := tx.ExecContext(r.Context(), `
					INSERT INTO sub_nodes_to_subscription_page_config (node_uuid, subpage_config_uuid)
					VALUES (?, ?)
					ON CONFLICT (node_uuid) DO UPDATE
					SET subpage_config_uuid = EXCLUDED.subpage_config_uuid
				`, req.UUID, *finalSubpageConfigUUID); err != nil {
					_ = tx.Rollback()
					return err
				}
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
	monitor.RequestSubNodeSync()
	monitor.RequestSubNodeDeploy(req.UUID)
	previousSubpageConfigUUID := normalizeNullableUUID(currentNode.SubpageConfigUUID)
	if previousSubpageConfigUUID != nil {
		if finalSubpageConfigUUID == nil || *finalSubpageConfigUUID != *previousSubpageConfigUUID {
			monitor.RequestSubNodeSubpageConfigPush(*previousSubpageConfigUUID, nil, req.UUID)
		}
	}
	if finalSubpageConfigUUID != nil {
		subpageConfig, err := fetchSubpageConfigRaw(r.Context(), manager, *finalSubpageConfigUUID)
		if err != nil {
			cfg.Logger.Warn("Failed to push selected subpage config to subscription node", "node_uuid", req.UUID, "subpage_config_uuid", *finalSubpageConfigUUID, "error", err)
		} else {
			monitor.RequestSubNodeSubpageConfigPush(*finalSubpageConfigUUID, subpageConfig, req.UUID)
		}
	}

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
	shared.WriteJSON(w, http.StatusOK, map[string]any{"response": response[0]})
}

func handleDeleteNode(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig, nodeUUID string) {
	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		result, err := db.ExecContext(r.Context(), `DELETE FROM sub_nodes WHERE uuid = ?`, nodeUUID)
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
	monitor.RequestSubNodeSync()

	shared.WriteJSON(w, http.StatusOK, map[string]any{"response": map[string]any{"isDeleted": true}})
}

func handleEnableNode(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig, nodeUUID string) {
	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		result, execErr := db.ExecContext(r.Context(), `UPDATE sub_nodes SET is_disabled = false, updated_at = CURRENT_TIMESTAMP WHERE uuid = ?`, nodeUUID)
		if execErr != nil {
			return execErr
		}
		rows, rowsErr := result.RowsAffected()
		if rowsErr != nil {
			return rowsErr
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
		shared.SendError(w, http.StatusInternalServerError, "failed to enable node", err, cfg)
		return
	}
	monitor.RequestSubNodeSync()

	sendUpdatedNodeResponse(w, r, manager, cfg, nodeUUID)
}

func handleDisableNode(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig, nodeUUID string) {
	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		result, execErr := db.ExecContext(r.Context(), `
			UPDATE sub_nodes
			SET is_disabled = true, is_connecting = false, is_connected = false, updated_at = CURRENT_TIMESTAMP
			WHERE uuid = ?
		`, nodeUUID)
		if execErr != nil {
			return execErr
		}
		rows, rowsErr := result.RowsAffected()
		if rowsErr != nil {
			return rowsErr
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
		shared.SendError(w, http.StatusInternalServerError, "failed to disable node", err, cfg)
		return
	}
	monitor.RequestSubNodeSync()

	sendUpdatedNodeResponse(w, r, manager, cfg, nodeUUID)
}

func handleRestartNode(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig, nodeUUID string) {
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
	monitor.RequestSubNodeDeploy(nodeUUID)

	shared.WriteJSON(w, http.StatusOK, map[string]any{"response": map[string]any{"eventSent": true}})
}

func handleResetNodeTraffic(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig, nodeUUID string) {
	_, err := getNodeByUUID(r.Context(), manager, nodeUUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			shared.SendError(w, http.StatusNotFound, "node not found", nil, cfg)
			return
		}
		shared.SendError(w, http.StatusInternalServerError, "failed to reset node traffic", err, cfg)
		return
	}

	shared.WriteJSON(w, http.StatusOK, map[string]any{"response": map[string]any{"eventSent": true}})
}

func handleRestartAllNodes(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) {
	var req restartAllNodesRequest
	if len(strings.TrimSpace(r.Header.Get("Content-Length"))) != 0 || r.ContentLength > 0 {
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			shared.SendError(w, http.StatusBadRequest, "invalid JSON", err, cfg)
			return
		}
	}

	var enabledCount int
	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		return db.QueryRowContext(r.Context(), `SELECT COUNT(*) FROM sub_nodes WHERE is_disabled = false`).Scan(&enabledCount)
	})
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to inspect nodes", err, cfg)
		return
	}
	if enabledCount == 0 {
		shared.SendError(w, http.StatusBadRequest, errNoEnabledNodes.Error(), nil, cfg)
		return
	}
	monitor.RequestSubNodeDeploy()

	_ = req
	shared.WriteJSON(w, http.StatusOK, map[string]any{"response": map[string]any{"eventSent": true}})
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
			if _, err := tx.ExecContext(r.Context(), `UPDATE sub_nodes SET view_position = ? WHERE uuid = ?`, item.ViewPosition, item.UUID); err != nil {
				_ = tx.Rollback()
				return err
			}
		}
		if _, err := tx.ExecContext(r.Context(), `SELECT setval('sub_nodes_view_position_seq', (SELECT COALESCE(MAX(view_position), 0) FROM sub_nodes) + 1)`); err != nil {
			_ = tx.Rollback()
			return err
		}
		return tx.Commit()
	})
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to reorder nodes", err, cfg)
		return
	}
	monitor.RequestSubNodeSync()

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
	switch req.Action {
	case "RESTART":
		monitor.RequestSubNodeDeploy(req.UUIDs...)
	default:
		monitor.RequestSubNodeSync()
	}

	shared.WriteJSON(w, http.StatusOK, map[string]any{"response": map[string]any{"eventSent": true}})
}

func performEnableAction(ctx context.Context, manager *dbmanager.DatabaseManager, nodeUUID string) error {
	return manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		result, err := db.ExecContext(ctx, `UPDATE sub_nodes SET is_disabled = false, updated_at = CURRENT_TIMESTAMP WHERE uuid = ?`, nodeUUID)
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
}

func performDisableAction(ctx context.Context, manager *dbmanager.DatabaseManager, nodeUUID string) error {
	return manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		result, err := db.ExecContext(ctx, `
			UPDATE sub_nodes
			SET is_disabled = true, is_connecting = false, is_connected = false, updated_at = CURRENT_TIMESTAMP
			WHERE uuid = ?
		`, nodeUUID)
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
}

func performResetTrafficAction(ctx context.Context, manager *dbmanager.DatabaseManager, nodeUUID string) error {
	_, err := getNodeByUUID(ctx, manager, nodeUUID)
	return err
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
	if req.APISchema != nil {
		switch strings.ToLower(strings.TrimSpace(*req.APISchema)) {
		case "mtls", "tls":
		default:
			return fmt.Errorf("apiSchema must be mtls or tls")
		}
	}
	if err := validateTags(req.Tags); err != nil {
		return err
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
	if req.APISchema != nil {
		switch strings.ToLower(strings.TrimSpace(*req.APISchema)) {
		case "mtls", "tls":
		default:
			return fmt.Errorf("apiSchema must be mtls or tls")
		}
	}
	if req.Tags != nil {
		if err := validateTags(*req.Tags); err != nil {
			return err
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

func buildNodeResponses(ctx context.Context, manager *dbmanager.DatabaseManager, records []nodeRecord) ([]nodeAPI, error) {
	providerUUIDs := make([]string, 0, len(records))
	for _, record := range records {
		if record.ProviderUUID != nil && *record.ProviderUUID != "" {
			providerUUIDs = append(providerUUIDs, *record.ProviderUUID)
		}
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
		item.APISchema = normalizeSubNodeSchema(record.APISchema)
		item.APIPath = record.APIPath
		item.SubpageConfigUUID = record.SubpageConfigUUID
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
		item.ConsumptionMultiplier = 1
		item.Tags = ensureStringSlice(record.Tags)
		item.CPUCount = record.CPUCount
		item.CPUModel = record.CPUModel
		item.TotalRAM = record.TotalRAM
		item.CreatedAt = record.CreatedAt
		item.UpdatedAt = record.UpdatedAt
		item.ConfigProfile.ActiveConfigProfileUUID = nil
		item.ConfigProfile.ActiveInbounds = ensureInboundSlice(nil)
		item.ProviderUUID = record.ProviderUUID
		if record.ProviderUUID != nil {
			item.Provider = providersMap[*record.ProviderUUID]
		}
		response = append(response, item)
	}

	return response, nil
}

func getAllNodeRecords(ctx context.Context, manager *dbmanager.DatabaseManager) ([]nodeRecord, error) {
	var nodes []nodeRecord
	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		rows, err := db.QueryContext(ctx, `
			SELECT
				n.uuid, n.id, n.name, n.address, n.port, n.api_schema, n.api_path, sns.subpage_config_uuid,
				n.is_connected, n.is_connecting, n.is_disabled,
				n.last_status_change, n.last_status_message,
				n.provider_uuid, n.view_position, n.tags, n.created_at, n.updated_at
			FROM sub_nodes n
			LEFT JOIN sub_nodes_to_subscription_page_config sns ON sns.node_uuid = n.uuid
			ORDER BY n.view_position ASC
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
				n.uuid, n.id, n.name, n.address, n.port, n.api_schema, n.api_path, sns.subpage_config_uuid,
				n.is_connected, n.is_connecting, n.is_disabled,
				n.last_status_change, n.last_status_message,
				n.provider_uuid, n.view_position, n.tags, n.created_at, n.updated_at
			FROM sub_nodes n
			LEFT JOIN sub_nodes_to_subscription_page_config sns ON sns.node_uuid = n.uuid
			WHERE n.uuid = ?
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
	var lastStatusChange sql.NullTime
	var lastStatusMessage sql.NullString
	var providerUUID sql.NullString
	var subpageConfigUUID sql.NullString
	var tags dbutil.StringArray

	err := scanner.Scan(
		&node.UUID,
		&id,
		&node.Name,
		&node.Address,
		&port,
		&node.APISchema,
		&node.APIPath,
		&subpageConfigUUID,
		&node.IsConnected,
		&node.IsConnecting,
		&node.IsDisabled,
		&lastStatusChange,
		&lastStatusMessage,
		&providerUUID,
		&node.ViewPosition,
		&tags,
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
	if lastStatusChange.Valid {
		ts := lastStatusChange.Time
		node.LastStatusChange = &ts
	}
	if lastStatusMessage.Valid {
		msg := lastStatusMessage.String
		node.LastStatusMessage = &msg
	}
	if providerUUID.Valid {
		node.ProviderUUID = &providerUUID.String
	}
	if subpageConfigUUID.Valid {
		value := strings.TrimSpace(subpageConfigUUID.String)
		if value != "" {
			node.SubpageConfigUUID = &value
		}
	}
	node.APISchema = normalizeSubNodeSchema(node.APISchema)
	node.SingboxVersion = nil
	node.NodeVersion = nil
	node.SingboxUptime = "0"
	node.CPUCount = nil
	node.CPUModel = nil
	node.TotalRAM = nil
	node.Tags = tags.Slice()
	node.CountryCode = "XX"
	node.ConsumptionMultiplier = 1_000_000_000
	node.IsTrafficTrackingActive = false

	return node, nil
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

func getNodeTags(ctx context.Context, manager *dbmanager.DatabaseManager) ([]string, error) {
	tags := make([]string, 0)
	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		rows, err := db.QueryContext(ctx, `SELECT DISTINCT unnest(tags) AS tag FROM sub_nodes ORDER BY tag ASC`)
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

func normalizeNullableUUID(value *string) *string {
	if value == nil {
		return nil
	}
	trimmed := strings.TrimSpace(*value)
	if trimmed == "" {
		return nil
	}
	return &trimmed
}

func normalizeAPISchema(value *string) string {
	if value == nil {
		return "mtls"
	}
	normalized := strings.ToLower(strings.TrimSpace(*value))
	if normalized == "" {
		return "mtls"
	}
	switch normalized {
	case "tls":
		return "tls"
	case "mtls":
		return "mtls"
	default:
		return "mtls"
	}
}

func normalizeSubNodeSchema(value string) string {
	normalized := strings.ToLower(strings.TrimSpace(value))
	switch normalized {
	case "tls":
		return "tls"
	case "mtls":
		return "mtls"
	default:
		return "mtls"
	}
}

func normalizeAPIPath(value *string) string {
	if value == nil {
		return "/"
	}
	path := strings.TrimSpace(*value)
	if path == "" {
		return "/"
	}
	if !strings.HasPrefix(path, "/") {
		return "/" + path
	}
	return path
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

func subpageConfigExists(ctx context.Context, manager *dbmanager.DatabaseManager, configUUID string) (bool, error) {
	exists := false
	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		var count int
		if err := db.QueryRowContext(ctx, `SELECT COUNT(*) FROM subscription_page_config WHERE uuid = ?`, configUUID).Scan(&count); err != nil {
			return err
		}
		exists = count > 0
		return nil
	})
	return exists, err
}

func fetchSubpageConfigRaw(ctx context.Context, manager *dbmanager.DatabaseManager, configUUID string) ([]byte, error) {
	var payload string
	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		return db.QueryRowContext(ctx, `
			SELECT config
			FROM subscription_page_config
			WHERE uuid = ?
			LIMIT 1
		`, configUUID).Scan(&payload)
	})
	if err != nil {
		return nil, err
	}

	raw := []byte(strings.TrimSpace(payload))
	if len(raw) == 0 || !json.Valid(raw) {
		return nil, fmt.Errorf("invalid subpage config payload")
	}

	return raw, nil
}
