package hosts

import (
	"database/sql"
	"net/http"
	"strings"

	"exodus/internal/config"
	"exodus/internal/httpapi/shared"

	"github.com/google/uuid"
)

func HostsHandler(db *sql.DB, cfg *config.BackendConfig) http.HandlerFunc {
	repo := NewHostRepository(db)
	service := NewHostService(repo, cfg)
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handleGetHosts(w, r, service)
		case http.MethodPost:
			handleCreateHost(w, r, service)
		case http.MethodPatch:
			handleUpdateHost(w, r, service)
		default:
			shared.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
		}
	}
}

func HostByUUIDHandler(db *sql.DB, cfg *config.BackendConfig) http.HandlerFunc {
	repo := NewHostRepository(db)
	service := NewHostService(repo, cfg)
	return func(w http.ResponseWriter, r *http.Request) {
		uuidStr := strings.TrimSpace(strings.TrimPrefix(r.URL.Path, "/api/hosts/"))
		if uuidStr == "" {
			switch r.Method {
			case http.MethodGet:
				handleGetHosts(w, r, service)
			case http.MethodPost:
				handleCreateHost(w, r, service)
			case http.MethodPatch:
				handleUpdateHost(w, r, service)
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
			handleGetHost(w, r, service, uuidStr)
		case http.MethodDelete:
			handleDeleteHost(w, r, service, uuidStr)
		default:
			shared.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
		}
	}
}

func HostsActionsHandler(db *sql.DB, cfg *config.BackendConfig) http.HandlerFunc {
	repo := NewHostRepository(db)
	service := NewHostService(repo, cfg)
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			shared.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}
		path := strings.TrimPrefix(r.URL.Path, "/api/hosts/actions/")
		path = strings.Trim(path, "/")
		switch path {
		case "reorder":
			handleReorderHosts(w, r, service)
		default:
			http.NotFound(w, r)
		}
	}
}

func HostsBulkHandler(db *sql.DB, cfg *config.BackendConfig) http.HandlerFunc {
	repo := NewHostRepository(db)
	service := NewHostService(repo, cfg)
	return func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/api/hosts/bulk/")
		path = strings.Trim(path, "/")

		requireMethod := func(method string) bool {
			if r.Method != method {
				shared.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
				return false
			}
			return true
		}

		switch path {
		case "enable":
			if !requireMethod(http.MethodPost) {
				return
			}
			handleBulkEnableHosts(w, r, service)
		case "disable":
			if !requireMethod(http.MethodPost) {
				return
			}
			handleBulkDisableHosts(w, r, service)
		case "delete":
			if !requireMethod(http.MethodPost) {
				return
			}
			handleBulkDeleteHosts(w, r, service)
		case "set-inbound":
			if !requireMethod(http.MethodPost) {
				return
			}
			handleBulkSetInbound(w, r, service)
		case "set-port":
			if !requireMethod(http.MethodPost) {
				return
			}
			handleBulkSetPort(w, r, service)
		case "update":
			if !requireMethod(http.MethodPatch) {
				return
			}
			handleBulkUpdateHosts(w, r, service)
		default:
			http.NotFound(w, r)
		}
	}
}

func HostsTagsHandler(db *sql.DB, cfg *config.BackendConfig) http.HandlerFunc {
	repo := NewHostRepository(db)
	service := NewHostService(repo, cfg)
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			shared.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}
		tags, err := service.repo.getHostTags(r.Context())
		if err != nil {
			shared.SendError(w, http.StatusInternalServerError, "failed to fetch host tags", err, cfg)
			return
		}
		shared.WriteJSON(w, http.StatusOK, map[string]any{"response": map[string]any{"tags": tags}})
	}
}
