package hosts

import (
	"encoding/json"
	"errors"
	"net/http"

	"exodus/internal/httpapi/shared"

	"github.com/google/uuid"
)

func handleGetHosts(w http.ResponseWriter, r *http.Request, service *HostService) {
	records, err := service.repo.getHosts(r.Context())
	if err != nil {
		shared.SendAPIError(w, shared.ErrGetAllHostsFailed.WithCause(err), service.cfg)
		return
	}

	uuids := make([]string, 0, len(records))
	for _, rec := range records {
		uuids = append(uuids, rec.UUID)
	}

	nodesMap, err := service.repo.getHostNodes(r.Context(), uuids)
	if err != nil {
		shared.SendAPIError(w, shared.ErrGetAllHostsFailed.WithCause(err), service.cfg)
		return
	}

	excludedMap, err := service.repo.getHostExcludedSquads(r.Context(), uuids)
	if err != nil {
		shared.SendAPIError(w, shared.ErrGetAllHostsFailed.WithCause(err), service.cfg)
		return
	}

	response := make([]HostAPI, 0, len(records))
	for _, rec := range records {
		response = append(response, mapHostRecordToAPI(rec, nodesMap[rec.UUID], excludedMap[rec.UUID]))
	}

	shared.WriteJSON(w, http.StatusOK, map[string]any{"response": response})
}

func handleGetHost(w http.ResponseWriter, r *http.Request, service *HostService, hostUUID string) {
	rec, err := service.repo.getHostByUUID(r.Context(), hostUUID)
	if err != nil {
		if errors.Is(err, errHostNotFound) {
			shared.SendAPIError(w, shared.ErrHostNotFound, service.cfg)
			return
		}
		shared.SendAPIError(w, shared.ErrGetOneHostFailed.WithCause(err), service.cfg)
		return
	}

	nodesMap, err := service.repo.getHostNodes(r.Context(), []string{hostUUID})
	if err != nil {
		shared.SendAPIError(w, shared.ErrGetOneHostFailed.WithCause(err), service.cfg)
		return
	}
	excludedMap, err := service.repo.getHostExcludedSquads(r.Context(), []string{hostUUID})
	if err != nil {
		shared.SendAPIError(w, shared.ErrGetOneHostFailed.WithCause(err), service.cfg)
		return
	}

	shared.WriteJSON(w, http.StatusOK, map[string]any{"response": mapHostRecordToAPI(rec, nodesMap[rec.UUID], excludedMap[rec.UUID])})
}

func handleCreateHost(w http.ResponseWriter, r *http.Request, service *HostService) {
	var req HostCreateRequestAPI
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		shared.SendError(w, http.StatusBadRequest, "invalid JSON", err, service.cfg)
		return
	}
	if err := validateCreateRequest(req); err != nil {
		shared.SendError(w, http.StatusBadRequest, err.Error(), nil, service.cfg)
		return
	}

	created, err := service.CreateHost(r.Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, errConfigProfileNotFound):
			shared.SendAPIError(w, shared.ErrConfigProfileNotFound, service.cfg)
		case errors.Is(err, errConfigProfileInboundNotFound):
			shared.SendAPIError(w, shared.ErrConfigProfileInboundNotFoundInProfile, service.cfg)
		case errors.Is(err, errTemplateNotFound):
			shared.SendAPIError(w, shared.ErrSubTemplateNotFound, service.cfg)
		case errors.Is(err, errTemplateTypeNotAllowed):
			shared.SendAPIError(w, shared.ErrSubTemplateTypeNotAllowed, service.cfg)
		default:
			shared.SendAPIError(w, shared.ErrCreateHostFailed.WithCause(err), service.cfg)
		}
		return
	}

	result := mapHostRecordToAPI(created, req.Nodes, req.ExcludedInternalSquads)
	shared.WriteJSON(w, http.StatusCreated, map[string]any{"response": result})
}

func handleUpdateHost(w http.ResponseWriter, r *http.Request, service *HostService) {
	var req HostUpdateRequestAPI
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		shared.SendError(w, http.StatusBadRequest, "invalid JSON", err, service.cfg)
		return
	}
	if _, err := uuid.Parse(req.UUID); err != nil {
		shared.SendError(w, http.StatusBadRequest, "invalid UUID format", nil, service.cfg)
		return
	}
	if err := validateUpdateRequest(req.hostUpdateFields); err != nil {
		shared.SendError(w, http.StatusBadRequest, err.Error(), nil, service.cfg)
		return
	}

	updated, err := service.UpdateHost(r.Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, errHostNotFound):
			shared.SendAPIError(w, shared.ErrHostNotFound, service.cfg)
		case errors.Is(err, errConfigProfileNotFound):
			shared.SendAPIError(w, shared.ErrConfigProfileNotFound, service.cfg)
		case errors.Is(err, errConfigProfileInboundNotFound):
			shared.SendAPIError(w, shared.ErrConfigProfileInboundNotFoundInProfile, service.cfg)
		case errors.Is(err, errTemplateNotFound):
			shared.SendAPIError(w, shared.ErrSubTemplateNotFound, service.cfg)
		case errors.Is(err, errTemplateTypeNotAllowed):
			shared.SendAPIError(w, shared.ErrSubTemplateTypeNotAllowed, service.cfg)
		default:
			shared.SendAPIError(w, shared.ErrUpdateHostFailed.WithCause(err), service.cfg)
		}
		return
	}

	nodesMap, _ := service.repo.getHostNodes(r.Context(), []string{req.UUID})
	excludedMap, _ := service.repo.getHostExcludedSquads(r.Context(), []string{req.UUID})

	shared.WriteJSON(w, http.StatusOK, map[string]any{"response": mapHostRecordToAPI(updated, nodesMap[req.UUID], excludedMap[req.UUID])})
}

func handleDeleteHost(w http.ResponseWriter, r *http.Request, service *HostService, hostUUID string) {
	err := service.DeleteHost(r.Context(), hostUUID)
	if err != nil {
		if errors.Is(err, errHostNotFound) {
			shared.SendAPIError(w, shared.ErrHostNotFound, service.cfg)
			return
		}
		shared.SendAPIError(w, shared.ErrDeleteHostFailed.WithCause(err), service.cfg)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
