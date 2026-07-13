package squads

import (
	"net/http"
	"strings"

	"exodus/internal/config"
	dbmanager "exodus/internal/db/manager"
	"exodus/internal/httpapi/shared"

	"github.com/google/uuid"
)

func InternalSquadsHandler(manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) http.HandlerFunc {
	repo := NewSquadRepository(manager)
	service := NewSquadService(repo, cfg)
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handleGetInternalSquads(w, r, service)
		case http.MethodPost:
			handleCreateInternalSquad(w, r, service)
		case http.MethodPatch:
			handleUpdateInternalSquad(w, r, service)
		default:
			shared.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
		}
	}
}

func InternalSquadsReorderHandler(manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) http.HandlerFunc {
	repo := NewSquadRepository(manager)
	service := NewSquadService(repo, cfg)
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			shared.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}
		handleReorderInternalSquads(w, r, service)
	}
}

func InternalSquadByUUIDHandler(manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) http.HandlerFunc {
	repo := NewSquadRepository(manager)
	service := NewSquadService(repo, cfg)
	return func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/api/internal-squads/")
		parts := strings.Split(path, "/")
		if len(parts) == 0 || parts[0] == "" {
			http.NotFound(w, r)
			return
		}
		squadUUID := parts[0]
		if _, err := uuid.Parse(squadUUID); err != nil {
			shared.SendError(w, http.StatusBadRequest, "invalid UUID format", nil, cfg)
			return
		}

		if len(parts) == 2 && parts[1] == "accessible-nodes" && r.Method == http.MethodGet {
			handleGetInternalSquadAccessibleNodes(w, r, service, squadUUID)
			return
		}

		if len(parts) != 1 {
			http.NotFound(w, r)
			return
		}

		switch r.Method {
		case http.MethodGet:
			handleGetInternalSquad(w, r, service, squadUUID)
		case http.MethodDelete:
			handleDeleteInternalSquad(w, r, service, squadUUID)
		default:
			shared.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
		}
	}
}

func InboundAssignmentsHandler(manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) http.HandlerFunc {
	repo := NewSquadRepository(manager)
	service := NewSquadService(repo, cfg)
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handleGetInboundAssignments(w, r, service)
		case http.MethodPost:
			handleSetInboundAssignments(w, r, service)
		default:
			http.Error(w, `{"error": "method not allowed"}`, http.StatusMethodNotAllowed)
		}
	}
}

func ConfigProfilesWithInboundsHandler(manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) http.HandlerFunc {
	repo := NewSquadRepository(manager)
	service := NewSquadService(repo, cfg)
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, `{"error": "method not allowed"}`, http.StatusMethodNotAllowed)
			return
		}
		handleGetConfigProfilesWithInbounds(w, r, service)
	}
}
