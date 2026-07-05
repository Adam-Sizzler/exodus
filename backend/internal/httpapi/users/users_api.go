package users

import (
	"net/http"
	"strings"

	"exodus/internal/config"
	dbmanager "exodus/internal/db/manager"
	"exodus/internal/httpapi/shared"

	"github.com/google/uuid"
)

func UsersHandler(manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handleGetUsers(w, r, manager, cfg)
		case http.MethodPost:
			handleCreateUser(w, r, manager, cfg)
		case http.MethodPatch:
			handleUpdateUser(w, r, manager, cfg)
		default:
			shared.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
		}
	}
}

func UserByUUIDHandler(manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := trimUsersPath(r.URL.Path, "/")
		if path == "" {
			switch r.Method {
			case http.MethodGet:
				handleGetUsers(w, r, manager, cfg)
			case http.MethodPost:
				handleCreateUser(w, r, manager, cfg)
			case http.MethodPatch:
				handleUpdateUser(w, r, manager, cfg)
			default:
				shared.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
			}
			return
		}

		parts := strings.Split(path, "/")
		if len(parts) == 1 && parts[0] == "resolve" {
			if r.Method != http.MethodPost {
				shared.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
				return
			}
			handleResolveUser(w, r, manager, cfg)
			return
		}

		userUUID := parts[0]
		if _, err := uuid.Parse(userUUID); err != nil {
			shared.SendError(w, http.StatusBadRequest, "invalid UUID format", nil, cfg)
			return
		}

		if len(parts) >= 3 && parts[1] == "actions" && r.Method == http.MethodPost {
			switch parts[2] {
			case "enable":
				handleEnableUser(w, r, manager, cfg, userUUID)
			case "disable":
				handleDisableUser(w, r, manager, cfg, userUUID)
			case "reset-traffic":
				handleResetUserTraffic(w, r, manager, cfg, userUUID)
			case "revoke":
				handleRevokeUserSubscription(w, r, manager, cfg, userUUID)
			default:
				http.NotFound(w, r)
			}
			return
		}

		if len(parts) == 2 && parts[1] == "subscription-request-history" {
			if r.Method != http.MethodGet {
				shared.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
				return
			}
			handleGetUserSubscriptionRequestHistory(w, r, manager, cfg, userUUID)
			return
		}

		if len(parts) == 2 && parts[1] == "accessible-nodes" {
			if r.Method != http.MethodGet {
				shared.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
				return
			}
			handleGetUserAccessibleNodes(w, r, manager, cfg, userUUID)
			return
		}

		if len(parts) != 1 {
			http.NotFound(w, r)
			return
		}

		switch r.Method {
		case http.MethodGet:
			handleGetUser(w, r, manager, cfg, userUUID)
		case http.MethodDelete:
			handleDeleteUser(w, r, manager, cfg, userUUID)
		default:
			shared.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
		}
	}
}

func UsersBulkHandler(manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			shared.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}
		path := trimUsersPath(r.URL.Path, "/bulk/")
		switch path {
		case "delete":
			handleBulkDeleteUsers(w, r, manager, cfg)
		case "delete-by-status":
			handleBulkDeleteUsersByStatus(w, r, manager, cfg)
		case "reset-traffic":
			handleBulkResetUsersTraffic(w, r, manager, cfg)
		case "update":
			handleBulkUpdateUsers(w, r, manager, cfg)
		case "update-squads":
			handleBulkUpdateUsersSquads(w, r, manager, cfg)
		case "extend-expiration-date":
			handleBulkExtendUsersExpirationDate(w, r, manager, cfg)
		case "all/reset-traffic":
			handleBulkAllResetUsersTraffic(w, r, manager, cfg)
		case "all/extend-expiration-date":
			handleBulkAllExtendUsersExpirationDate(w, r, manager, cfg)
		case "all/update":
			handleBulkAllUpdateUsers(w, r, manager, cfg)
		default:
			http.NotFound(w, r)
		}
	}
}

func UsersTagsHandler(manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			shared.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}

		tags, err := getAllUserTags(r.Context(), manager)
		if err != nil {
			shared.SendError(w, http.StatusInternalServerError, "failed to fetch user tags", err, cfg)
			return
		}

		shared.WriteJSON(w, http.StatusOK, map[string]any{
			"response": map[string]any{
				"tags": tags,
			},
		})
	}
}

func trimUsersPath(path string, suffix string) string {
	for _, prefix := range []string{"/api/users"} {
		if strings.HasPrefix(path, prefix+suffix) {
			return strings.Trim(strings.TrimPrefix(path, prefix+suffix), "/")
		}
	}
	return strings.Trim(path, "/")
}
