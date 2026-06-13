package hosts

import (
	"net/http"
	"strings"

	"exodus/internal/config"
	dbmanager "exodus/internal/db/manager"
	"exodus/internal/httpapi/shared"

	"github.com/google/uuid"
)

func HostsHandler(manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handleGetHosts(w, r, manager, cfg)
		case http.MethodPost:
			handleCreateHost(w, r, manager, cfg)
		case http.MethodPatch:
			handleUpdateHost(w, r, manager, cfg)
		default:
			shared.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
		}
	}
}

func HostByUUIDHandler(manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uuidStr := strings.TrimSpace(strings.TrimPrefix(r.URL.Path, "/api/hosts/"))
		if uuidStr == "" {
			switch r.Method {
			case http.MethodGet:
				handleGetHosts(w, r, manager, cfg)
			case http.MethodPost:
				handleCreateHost(w, r, manager, cfg)
			case http.MethodPatch:
				handleUpdateHost(w, r, manager, cfg)
			default:
				shared.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
			}
			return
		}
		if _, err := uuid.Parse(uuidStr); err != nil {
			shared.SendError(w, http.StatusBadRequest, "invalid UUID format", nil, cfg)
			return
		}

		switch r.Method {
		case http.MethodGet:
			handleGetHost(w, r, manager, cfg, uuidStr)
		case http.MethodDelete:
			handleDeleteHost(w, r, manager, cfg, uuidStr)
		default:
			shared.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
		}
	}
}

func HostsActionsHandler(manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			shared.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}
		path := strings.TrimPrefix(r.URL.Path, "/api/hosts/actions/")
		path = strings.Trim(path, "/")
		switch path {
		case "reorder":
			handleReorderHosts(w, r, manager, cfg)
		default:
			http.NotFound(w, r)
		}
	}
}

func HostsBulkHandler(manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			shared.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}

		path := strings.TrimPrefix(r.URL.Path, "/api/hosts/bulk/")
		path = strings.Trim(path, "/")
		switch path {
		case "enable":
			handleBulkEnableHosts(w, r, manager, cfg)
		case "disable":
			handleBulkDisableHosts(w, r, manager, cfg)
		case "delete":
			handleBulkDeleteHosts(w, r, manager, cfg)
		case "set-inbound":
			handleBulkSetInbound(w, r, manager, cfg)
		case "set-port":
			handleBulkSetPort(w, r, manager, cfg)
		default:
			http.NotFound(w, r)
		}
	}
}

func HostsTagsHandler(manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			shared.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}
		tags, err := getHostTags(r.Context(), manager)
		if err != nil {
			shared.SendError(w, http.StatusInternalServerError, "failed to fetch host tags", err, cfg)
			return
		}
		shared.WriteJSON(w, http.StatusOK, map[string]any{"response": map[string]any{"tags": tags}})
	}
}
