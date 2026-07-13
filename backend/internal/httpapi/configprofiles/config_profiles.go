package configprofiles

import (
	"net/http"
	"strings"

	"exodus/internal/config"
	dbmanager "exodus/internal/db/manager"
	"exodus/internal/httpapi/shared"

	"github.com/google/uuid"
)

func ConfigProfilesHandler(manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) http.HandlerFunc {
	repo := NewConfigProfileRepository(manager)
	service := NewConfigProfileService(repo, cfg)
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handleGetConfigProfiles(w, r, service)
		case http.MethodPost:
			handleCreateConfigProfile(w, r, service)
		case http.MethodPatch:
			handleUpdateConfigProfile(w, r, service)
		default:
			shared.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
		}
	}
}

func ConfigProfileByUUIDHandler(manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) http.HandlerFunc {
	repo := NewConfigProfileRepository(manager)
	service := NewConfigProfileService(repo, cfg)
	return func(w http.ResponseWriter, r *http.Request) {
		path := trimConfigProfilesPath(r.URL.Path, "/")
		if path == "" {
			switch r.Method {
			case http.MethodGet:
				handleGetConfigProfiles(w, r, service)
			case http.MethodPost:
				handleCreateConfigProfile(w, r, service)
			case http.MethodPatch:
				handleUpdateConfigProfile(w, r, service)
			default:
				shared.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
			}
			return
		}
		parts := strings.Split(path, "/")
		profileUUID := parts[0]
		if _, err := uuid.Parse(profileUUID); err != nil {
			shared.SendError(w, http.StatusBadRequest, "invalid UUID format", nil, cfg)
			return
		}

		if len(parts) == 2 && r.Method == http.MethodGet {
			switch parts[1] {
			case "inbounds":
				handleGetConfigProfileInbounds(w, r, service, profileUUID)
			case "computed-config":
				handleGetComputedConfigProfile(w, r, service, profileUUID)
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
			handleGetConfigProfile(w, r, service, profileUUID)
		case http.MethodDelete:
			handleDeleteConfigProfile(w, r, service, profileUUID)
		default:
			shared.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
		}
	}
}

func ConfigProfilesActionsHandler(manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) http.HandlerFunc {
	repo := NewConfigProfileRepository(manager)
	service := NewConfigProfileService(repo, cfg)
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			shared.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}

		path := trimConfigProfilesPath(r.URL.Path, "/actions/")
		switch path {
		case "reorder":
			handleReorderConfigProfiles(w, r, service)
		default:
			http.NotFound(w, r)
		}
	}
}

func ConfigProfilesInboundsHandler(manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) http.HandlerFunc {
	repo := NewConfigProfileRepository(manager)
	service := NewConfigProfileService(repo, cfg)
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			shared.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}
		handleGetAllInbounds(w, r, service)
	}
}

func ConfigProfileSnippetsHandler(manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) http.HandlerFunc {
	repo := NewConfigProfileRepository(manager)
	service := NewConfigProfileService(repo, cfg)
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handleGetConfigProfileSnippets(w, r, service)
		case http.MethodPost:
			handleCreateConfigProfileSnippet(w, r, service)
		case http.MethodPatch:
			handleUpdateConfigProfileSnippet(w, r, service)
		case http.MethodDelete:
			handleDeleteConfigProfileSnippet(w, r, service)
		default:
			shared.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
		}
	}
}

func trimConfigProfilesPath(path string, suffix string) string {
	for _, prefix := range []string{"/api/config-profiles"} {
		if strings.HasPrefix(path, prefix+suffix) {
			return strings.Trim(strings.TrimPrefix(path, prefix+suffix), "/")
		}
	}
	return strings.Trim(path, "/")
}
