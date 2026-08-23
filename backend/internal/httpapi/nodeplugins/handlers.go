package nodeplugins

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"exodus/internal/config"
	"exodus/internal/httpapi/shared"
	monitor "exodus/internal/nodes"

	"github.com/google/uuid"
)

// Handler godoc
// @Summary      Manage node plugins
// @Description  List, create (201), update, delete (204) node plugins or reorder
// @Tags         Node Plugins Controller
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        uuid  path      string  false  "Node plugin UUID" format(uuid)
// @Param        body  body      object  false  "Node plugin payload"
// @Success      200   {object}  map[string]any
// @Success      201   {object}  map[string]any
// @Success      204
// @Failure      400   {object}  shared.ErrorResponse
// @Failure      404   {object}  shared.ErrorResponse
// @Failure      500   {object}  shared.ErrorResponse
// @Router       /node-plugins [get]
// @Router       /node-plugins [post]
// @Router       /node-plugins [patch]
// @Router       /node-plugins/{uuid} [get]
// @Router       /node-plugins/{uuid} [delete]
// @Router       /node-plugins/actions/reorder [post]
func Handler(db *sql.DB, cfg *config.BackendConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if cfg != nil {
			path = strings.TrimPrefix(path, cfg.Backend.Trimmed())
		}
		path = strings.TrimPrefix(path, "/api/node-plugins")
		path = strings.Trim(path, "/")

		switch {
		case path == "":
			switch r.Method {
			case http.MethodGet:
				handleList(w, r, db, cfg)
			case http.MethodPost:
				handleCreate(w, r, db, cfg)
			case http.MethodPatch:
				handleUpdate(w, r, db, cfg, "")
			default:
				shared.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
			}
		case path == "executor":
			handleExecutor(w, r, db, cfg)
		case strings.HasPrefix(path, "actions/"):
			handleAction(w, r, db, cfg, strings.TrimPrefix(path, "actions/"))
		case path == "shared-lists" || strings.HasPrefix(path, "shared-lists/"):
			handleSharedLists(w, r, db, cfg, strings.TrimPrefix(strings.TrimPrefix(path, "shared-lists"), "/"))
		default:
			handleByUUID(w, r, db, cfg, path)
		}
	}
}

func handleList(w http.ResponseWriter, r *http.Request, db *sql.DB, cfg *config.BackendConfig) {
	plugins, err := loadPlugins(r.Context(), db)
	if err != nil {
		cfg.Logger.Error("Failed to load node plugins", "error", err)
		shared.WriteJSONError(w, http.StatusInternalServerError, "failed to load node plugins")
		return
	}
	shared.WriteJSON(w, http.StatusOK, responseEnvelope[listPayload]{
		Response: listPayload{NodePlugins: plugins, Total: len(plugins)},
	})
}

func handleCreate(w http.ResponseWriter, r *http.Request, db *sql.DB, cfg *config.BackendConfig) {
	var req createRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		shared.WriteJSONError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	name := strings.TrimSpace(req.Name)
	if name == "" {
		shared.WriteJSONError(w, http.StatusBadRequest, "name is required")
		return
	}
	configJSON, err := normalizePluginConfig(req.PluginConfig)
	if err != nil {
		shared.WriteJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	plugin, err := createPlugin(r.Context(), db, name, configJSON)
	if err != nil {
		cfg.Logger.Error("Failed to create node plugin", "error", err)
		shared.WriteJSONError(w, http.StatusInternalServerError, "failed to create node plugin")
		return
	}
	shared.WriteJSON(w, http.StatusCreated, responseEnvelope[nodePlugin]{Response: plugin})
}

func handleByUUID(w http.ResponseWriter, r *http.Request, db *sql.DB, cfg *config.BackendConfig, rawPath string) {
	parts := strings.Split(strings.Trim(rawPath, "/"), "/")
	if len(parts) == 0 || parts[0] == "" {
		shared.WriteJSONError(w, http.StatusNotFound, "not found")
		return
	}
	pluginUUID := parts[0]
	if _, err := uuid.Parse(pluginUUID); err != nil {
		shared.WriteJSONError(w, http.StatusBadRequest, "invalid uuid")
		return
	}

	if len(parts) > 1 {
		shared.WriteJSONError(w, http.StatusNotFound, "not found")
		return
	}

	switch r.Method {
	case http.MethodGet:
		plugin, err := loadPluginByUUID(r.Context(), db, pluginUUID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				shared.WriteJSONError(w, http.StatusNotFound, "node plugin not found")
				return
			}
			cfg.Logger.Error("Failed to load node plugin", "uuid", pluginUUID, "error", err)
			shared.WriteJSONError(w, http.StatusInternalServerError, "failed to load node plugin")
			return
		}
		shared.WriteJSON(w, http.StatusOK, responseEnvelope[nodePlugin]{Response: plugin})
	case http.MethodPatch:
		handleUpdate(w, r, db, cfg, pluginUUID)
	case http.MethodDelete:
		handleDelete(w, r, db, cfg, pluginUUID)
	default:
		shared.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func handleUpdate(w http.ResponseWriter, r *http.Request, db *sql.DB, cfg *config.BackendConfig, urlUUID string) {
	var req updateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		shared.WriteJSONError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	pluginUUID := urlUUID
	if pluginUUID == "" {
		if req.UUID == nil || *req.UUID == "" {
			shared.WriteJSONError(w, http.StatusBadRequest, "uuid is required")
			return
		}
		pluginUUID = *req.UUID
	}

	if _, err := uuid.Parse(pluginUUID); err != nil {
		shared.WriteJSONError(w, http.StatusBadRequest, "invalid uuid")
		return
	}

	var name *string
	if req.Name != nil {
		trimmed := strings.TrimSpace(*req.Name)
		if trimmed == "" {
			shared.WriteJSONError(w, http.StatusBadRequest, "name cannot be empty")
			return
		}
		name = &trimmed
	}

	var configJSON *json.RawMessage
	if req.PluginConfig != nil {
		normalized, err := normalizePluginConfig(*req.PluginConfig)
		if err != nil {
			shared.WriteJSONError(w, http.StatusBadRequest, err.Error())
			return
		}
		configJSON = &normalized
	}

	plugin, err := updatePlugin(r.Context(), db, pluginUUID, name, configJSON, req.ViewPosition)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			shared.WriteJSONError(w, http.StatusNotFound, "node plugin not found")
			return
		}
		cfg.Logger.Error("Failed to update node plugin", "uuid", pluginUUID, "error", err)
		shared.WriteJSONError(w, http.StatusInternalServerError, "failed to update node plugin")
		return
	}
	shared.WriteJSON(w, http.StatusOK, responseEnvelope[nodePlugin]{Response: plugin})
}

func handleDelete(w http.ResponseWriter, r *http.Request, db *sql.DB, cfg *config.BackendConfig, pluginUUID string) {
	if err := deletePlugin(r.Context(), db, pluginUUID); err != nil {
		cfg.Logger.Error("Failed to delete node plugin", "uuid", pluginUUID, "error", err)
		shared.WriteJSONError(w, http.StatusInternalServerError, "failed to delete node plugin")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func handleExecutor(w http.ResponseWriter, r *http.Request, db *sql.DB, cfg *config.BackendConfig) {
	if r.Method != http.MethodPost {
		shared.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var req executorRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		shared.WriteJSONError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	command := strings.TrimSpace(req.Command.Command)
	switch command {
	case "blockIps", "unblockIps", "recreateTables":
	default:
		shared.WriteJSONError(w, http.StatusBadRequest, "unsupported executor command")
		return
	}
	if req.TargetNodes.Target != "specificNodes" {
		shared.WriteJSONError(w, http.StatusBadRequest, "targetNodes.target must be specificNodes")
		return
	}
	targetNodeUUIDs := normalizeUUIDList(req.TargetNodes.NodeUUIDs)
	if len(targetNodeUUIDs) == 0 {
		shared.WriteJSONError(w, http.StatusBadRequest, "nodeUuids are required")
		return
	}
	for _, nodeUUID := range targetNodeUUIDs {
		if _, err := uuid.Parse(nodeUUID); err != nil {
			shared.WriteJSONError(w, http.StatusBadRequest, "invalid node uuid")
			return
		}
	}

	if err := ensureNodesExist(r.Context(), db, targetNodeUUIDs); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			shared.WriteJSONError(w, http.StatusBadRequest, "one or more nodes were not found")
			return
		}
		cfg.Logger.Error("Failed to validate node plugin executor targets", "error", err)
		shared.WriteJSONError(w, http.StatusInternalServerError, "failed to validate targets")
		return
	}

	cfg.Logger.Info(
		"Node plugin executor command accepted",
		"command", command,
		"nodes", strings.Join(targetNodeUUIDs, ","),
	)
	if err := monitor.RequestNodePluginExecutor(req.Command.Raw, targetNodeUUIDs...); err != nil {
		cfg.Logger.Warn("Failed to send node plugin executor command", "command", command, "error", err)
		shared.WriteJSONError(w, http.StatusBadGateway, err.Error())
		return
	}
	w.WriteHeader(http.StatusAccepted)
}

func handleAction(w http.ResponseWriter, r *http.Request, db *sql.DB, cfg *config.BackendConfig, action string) {
	if r.Method != http.MethodPost {
		shared.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	switch strings.Trim(action, "/") {
	case "reorder":
		var req reorderRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			shared.WriteJSONError(w, http.StatusBadRequest, "invalid JSON body")
			return
		}
		if len(req.Items) == 0 {
			shared.WriteJSONError(w, http.StatusBadRequest, "items are required")
			return
		}
		if err := reorderPlugins(r.Context(), db, req); err != nil {
			cfg.Logger.Error("Failed to reorder node plugins", "error", err)
			shared.WriteJSONError(w, http.StatusInternalServerError, "failed to reorder node plugins")
			return
		}
		plugins, err := loadPlugins(r.Context(), db)
		if err != nil {
			cfg.Logger.Error("Failed to load reordered node plugins", "error", err)
			shared.WriteJSONError(w, http.StatusInternalServerError, "failed to load node plugins")
			return
		}
		shared.WriteJSON(w, http.StatusOK, responseEnvelope[listPayload]{Response: listPayload{NodePlugins: plugins, Total: len(plugins)}})
	case "clone":
		var req cloneRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			shared.WriteJSONError(w, http.StatusBadRequest, "invalid JSON body")
			return
		}
		if _, err := uuid.Parse(strings.TrimSpace(req.CloneFromUUID)); err != nil {
			shared.WriteJSONError(w, http.StatusBadRequest, "invalid cloneFromUuid")
			return
		}
		plugin, err := clonePlugin(r.Context(), db, req)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				shared.WriteJSONError(w, http.StatusNotFound, "node plugin not found")
				return
			}
			cfg.Logger.Error("Failed to clone node plugin", "error", err)
			shared.WriteJSONError(w, http.StatusInternalServerError, "failed to clone node plugin")
			return
		}
		shared.WriteJSON(w, http.StatusCreated, responseEnvelope[nodePlugin]{Response: plugin})
	case "sync":
		var req struct {
			UUID string `json:"uuid"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			shared.WriteJSONError(w, http.StatusBadRequest, "invalid JSON body")
			return
		}
		pluginUUID := strings.TrimSpace(req.UUID)
		if _, err := uuid.Parse(pluginUUID); err != nil {
			shared.WriteJSONError(w, http.StatusBadRequest, "invalid uuid")
			return
		}
		if err := syncPlugin(r.Context(), db, pluginUUID); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				shared.WriteJSONError(w, http.StatusNotFound, "node plugin not found")
				return
			}
			cfg.Logger.Error("Failed to sync node plugin", "uuid", pluginUUID, "error", err)
			shared.WriteJSONError(w, http.StatusInternalServerError, "failed to sync node plugin")
			return
		}
		w.WriteHeader(http.StatusAccepted)
	default:
		shared.WriteJSONError(w, http.StatusNotFound, "not found")
	}
}
