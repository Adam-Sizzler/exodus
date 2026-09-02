package squads

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"exodus/internal/httpapi/shared"

	"github.com/google/uuid"
)

func handleGetInternalSquads(w http.ResponseWriter, r *http.Request, service *SquadService) {
	squads, err := service.repo.getSquads(r.Context())
	if err != nil {
		shared.SendAPIError(w, shared.ErrGetInternalSquadsFailed.WithCause(err), service.cfg)
		return
	}

	response := make([]InternalSquadAPI, 0, len(squads))
	for _, squad := range squads {
		membersCount, err := service.repo.getSquadMembersCount(r.Context(), squad.UUID)
		if err != nil {
			shared.SendAPIError(w, shared.ErrGetInternalSquadsFailed.WithCause(err), service.cfg)
			return
		}
		inbounds, err := service.repo.getSquadInbounds(r.Context(), squad.UUID)
		if err != nil {
			shared.SendAPIError(w, shared.ErrGetInternalSquadsFailed.WithCause(err), service.cfg)
			return
		}
		item := buildInternalSquadResponse(squad, membersCount, inbounds)
		response = append(response, item)
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"response": map[string]interface{}{
			"total":          len(response),
			"internalSquads": response,
		},
	})
}

func handleCreateInternalSquad(w http.ResponseWriter, r *http.Request, service *SquadService) {
	var req InternalSquadCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		shared.SendError(w, http.StatusBadRequest, "invalid JSON", err, service.cfg)
		return
	}
	if err := req.Validate(); err != nil {
		shared.SendError(w, http.StatusBadRequest, err.Error(), nil, service.cfg)
		return
	}

	result, err := service.CreateSquad(r.Context(), req)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") || strings.Contains(err.Error(), "already exists") {
			shared.SendAPIError(w, shared.ErrInternalSquadNameAlreadyExists, service.cfg)
			return
		}
		shared.SendAPIError(w, shared.ErrCreateInternalSquadFailed.WithCause(err), service.cfg)
		return
	}

	service.cfg.Logger.Info("Internal squad created", "uuid", result.UUID, "name", result.Name)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"response": result,
	})
}

func handleUpdateInternalSquad(w http.ResponseWriter, r *http.Request, service *SquadService) {
	var req InternalSquadUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		shared.SendError(w, http.StatusBadRequest, "invalid JSON", err, service.cfg)
		return
	}
	if req.UUID == "" {
		shared.SendError(w, http.StatusBadRequest, "uuid is required", nil, service.cfg)
		return
	}
	if _, err := uuid.Parse(req.UUID); err != nil {
		shared.SendError(w, http.StatusBadRequest, "invalid UUID format", nil, service.cfg)
		return
	}
	if err := req.Validate(); err != nil {
		shared.SendError(w, http.StatusBadRequest, err.Error(), nil, service.cfg)
		return
	}

	if !req.HasUpdates() {
		handleGetInternalSquad(w, r, service, req.UUID)
		return
	}

	result, err := service.UpdateSquad(r.Context(), req)
	if err != nil {
		if errors.Is(err, errSquadNotFound) {
			shared.SendAPIError(w, shared.ErrInternalSquadNotFound, service.cfg)
		} else if strings.Contains(strings.ToLower(err.Error()), "inbound not found") {
			shared.SendAPIError(w, shared.ErrConfigProfileInboundNotFoundInProfile.WithCause(err), service.cfg)
		} else {
			if strings.Contains(err.Error(), "UNIQUE constraint failed") || strings.Contains(err.Error(), "already exists") {
				shared.SendAPIError(w, shared.ErrInternalSquadNameAlreadyExists, service.cfg)
			} else {
				shared.SendAPIError(w, shared.ErrUpdateInternalSquadFailed.WithCause(err), service.cfg)
			}
		}
		return
	}

	service.cfg.Logger.Info(
		"Internal squad updated",
		"uuid", req.UUID,
		"name_updated", req.Name != nil,
		"view_position_updated", req.ViewPosition != nil,
		"inbounds_updated", req.Inbounds != nil,
		"inbounds_count", len(req.Inbounds),
	)

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"response": result,
	})
}

func handleDeleteInternalSquad(w http.ResponseWriter, r *http.Request, service *SquadService, squadUUID string) {
	_, err := service.DeleteSquad(r.Context(), squadUUID)
	if err != nil {
		if errors.Is(err, errSquadNotFound) {
			shared.SendAPIError(w, shared.ErrInternalSquadNotFound, service.cfg)
		} else {
			shared.SendAPIError(w, shared.ErrDeleteInternalSquadFailed.WithCause(err), service.cfg)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func handleGetInternalSquad(w http.ResponseWriter, r *http.Request, service *SquadService, squadUUID string) {
	squad, err := service.repo.getSquadByUUID(r.Context(), squadUUID)
	if err != nil {
		if errors.Is(err, errSquadNotFound) {
			shared.SendAPIError(w, shared.ErrInternalSquadNotFound, service.cfg)
		} else {
			shared.SendAPIError(w, shared.ErrGetInternalSquadByUUIDFailed.WithCause(err), service.cfg)
		}
		return
	}

	membersCount, err := service.repo.getSquadMembersCount(r.Context(), squadUUID)
	if err != nil {
		shared.SendAPIError(w, shared.ErrGetInternalSquadByUUIDFailed.WithCause(err), service.cfg)
		return
	}
	inbounds, err := service.repo.getSquadInbounds(r.Context(), squadUUID)
	if err != nil {
		shared.SendAPIError(w, shared.ErrGetInternalSquadByUUIDFailed.WithCause(err), service.cfg)
		return
	}

	result := buildInternalSquadResponse(squad, membersCount, inbounds)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"response": result,
	})
}

func handleGetInternalSquadAccessibleNodes(w http.ResponseWriter, r *http.Request, service *SquadService, squadUUID string) {
	nodes, err := service.repo.getSquadAccessibleNodes(r.Context(), squadUUID)
	if err != nil {
		if errors.Is(err, errInternalSquadNotFound) {
			shared.SendAPIError(w, shared.ErrInternalSquadNotFound, service.cfg)
		} else {
			shared.SendAPIError(w, shared.ErrGetInternalSquadAccessibleNodesFailed.WithCause(err), service.cfg)
		}
		return
	}

	if nodes == nil {
		nodes = []InternalSquadAccessibleNode{}
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"response": map[string]interface{}{
			"squadUuid":       squadUUID,
			"accessibleNodes": nodes,
		},
	})
}

func handleReorderInternalSquads(w http.ResponseWriter, r *http.Request, service *SquadService) {
	var req reorderSquadsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		shared.SendError(w, http.StatusBadRequest, "invalid JSON", err, service.cfg)
		return
	}
	if len(req.Squads) == 0 {
		req.Squads = req.Items
	}
	if len(req.Squads) == 0 {
		shared.SendError(w, http.StatusBadRequest, "items cannot be empty", nil, service.cfg)
		return
	}
	for _, item := range req.Squads {
		if _, err := uuid.Parse(item.UUID); err != nil {
			shared.SendError(w, http.StatusBadRequest, "invalid UUID format", nil, service.cfg)
			return
		}
	}

	err := service.ReorderSquads(r.Context(), req)
	if err != nil {
		shared.SendAPIError(w, shared.ErrReorderInternalSquadsFailed.WithCause(err), service.cfg)
		return
	}

	handleGetInternalSquads(w, r, service)
}

func handleGetInboundAssignments(w http.ResponseWriter, r *http.Request, service *SquadService) {
	nodeUUID := r.URL.Query().Get("node_uuid")
	assignments, err := service.repo.getInboundAssignments(r.Context(), nodeUUID)
	if err != nil {
		shared.SendAPIError(w, shared.ErrGetInternalSquadByUUIDFailed.WithCause(err), service.cfg)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"response": assignments,
	})
}

func handleSetInboundAssignments(w http.ResponseWriter, r *http.Request, service *SquadService) {
	var req InboundAssignmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		shared.SendError(w, http.StatusBadRequest, "invalid JSON", err, service.cfg)
		return
	}
	if err := req.Validate(); err != nil {
		shared.SendError(w, http.StatusBadRequest, err.Error(), nil, service.cfg)
		return
	}

	result, err := service.SetInboundAssignments(r.Context(), req)
	if err != nil {
		shared.SendAPIError(w, shared.ErrUpdateInternalSquadFailed.WithCause(err), service.cfg)
		return
	}

	service.cfg.Logger.Info(
		"Config profile inbound assignments for node updated",
		"node_uuid", req.NodeUUID,
		"inbound_uuids_count", len(req.InboundUUIDs),
	)

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"response": result,
	})
}

func handleGetConfigProfilesWithInbounds(w http.ResponseWriter, r *http.Request, service *SquadService) {
	profiles, err := service.repo.getConfigProfilesWithInbounds(r.Context())
	if err != nil {
		shared.SendAPIError(w, shared.ErrGetInternalSquadsFailed.WithCause(err), service.cfg)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"profiles": profiles,
		"count":    len(profiles),
	})
}

func handleGetInternalSquadTags(w http.ResponseWriter, r *http.Request, service *SquadService) {
	tags, err := service.GetTags(r.Context())
	if err != nil {
		shared.SendAPIError(w, shared.ErrGetInternalSquadsFailed.WithCause(err), service.cfg)
		return
	}
	shared.WriteJSON(w, http.StatusOK, map[string]any{
		"response": map[string]any{
			"tags": tags,
		},
	})
}

func handleSetInternalSquadTags(w http.ResponseWriter, r *http.Request, service *SquadService) {
	var req shared.SetEntityTagsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		shared.SendError(w, http.StatusBadRequest, "invalid JSON", err, service.cfg)
		return
	}
	if _, err := uuid.Parse(strings.TrimSpace(req.UUID)); err != nil {
		shared.SendError(w, http.StatusBadRequest, "invalid UUID format", nil, service.cfg)
		return
	}

	tags, err := service.SetTags(r.Context(), req.UUID, req.Tags)
	if err != nil {
		if errors.Is(err, errInternalSquadNotFound) {
			shared.SendAPIError(w, shared.ErrInternalSquadNotFound, service.cfg)
			return
		}
		shared.SendAPIError(w, shared.ErrUpdateInternalSquadFailed.WithCause(err), service.cfg)
		return
	}

	shared.WriteJSON(w, http.StatusOK, map[string]any{
		"response": map[string]any{
			"uuid": req.UUID,
			"tags": tags,
		},
	})
}
