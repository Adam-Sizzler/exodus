package subscriptionconnections

import (
	"net/http"
	"strings"

	"exodus/internal/config"
	dbmanager "exodus/internal/db/manager"
	"exodus/internal/httpapi/shared"

	"github.com/google/uuid"
)

func NodesHandler(manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) http.HandlerFunc {
	repo := NewSubscriptionConnectionRepository(manager)
	service := NewSubscriptionConnectionService(repo, cfg)
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handleGetNodes(w, r, service)
		case http.MethodPost:
			handleCreateNode(w, r, service)
		case http.MethodPatch:
			handleUpdateNode(w, r, service)
		default:
			shared.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
		}
	}
}

func NodeByUUIDHandler(manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) http.HandlerFunc {
	repo := NewSubscriptionConnectionRepository(manager)
	service := NewSubscriptionConnectionService(repo, cfg)
	return func(w http.ResponseWriter, r *http.Request) {
		path := trimNodesPath(r.URL.Path, "/")
		if path == "" {
			switch r.Method {
			case http.MethodGet:
				handleGetNodes(w, r, service)
			case http.MethodPost:
				handleCreateNode(w, r, service)
			case http.MethodPatch:
				handleUpdateNode(w, r, service)
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
				handleEnableNode(w, r, service, nodeUUID)
			case "disable":
				handleDisableNode(w, r, service, nodeUUID)
			case "restart":
				handleRestartNode(w, r, service, nodeUUID)
			case "reset-traffic":
				handleResetNodeTraffic(w, r, service, nodeUUID)
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
			handleGetNode(w, r, service, nodeUUID)
		case http.MethodDelete:
			handleDeleteNode(w, r, service, nodeUUID)
		default:
			shared.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
		}
	}
}

func NodesActionsHandler(manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) http.HandlerFunc {
	repo := NewSubscriptionConnectionRepository(manager)
	service := NewSubscriptionConnectionService(repo, cfg)
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			shared.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}
		path := trimNodesPath(r.URL.Path, "/actions/")
		switch path {
		case "restart-all":
			handleRestartAllNodes(w, r, service)
		case "reorder":
			handleReorderNodes(w, r, service)
		default:
			http.NotFound(w, r)
		}
	}
}

func NodesBulkActionsHandler(manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) http.HandlerFunc {
	repo := NewSubscriptionConnectionRepository(manager)
	service := NewSubscriptionConnectionService(repo, cfg)
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			shared.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}

		path := trimNodesPath(r.URL.Path, "/bulk-actions")
		switch path {
		case "":
			handleBulkNodesActions(w, r, service)
		case "profile-modification":
			handleBulkProfileModification(w, r, service)
		default:
			http.NotFound(w, r)
		}
	}
}

func NodesTagsHandler(manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) http.HandlerFunc {
	repo := NewSubscriptionConnectionRepository(manager)
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			shared.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}
		tags, err := repo.getNodeTags(r.Context())
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
