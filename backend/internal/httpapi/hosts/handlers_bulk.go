package hosts

import (
	"encoding/json"
	"errors"
	"net/http"

	"exodus/internal/httpapi/shared"

	"github.com/google/uuid"
)

func handleBulkUpdateHosts(w http.ResponseWriter, r *http.Request, service *HostService) {
	var req HostBulkUpdateRequestAPI
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		shared.SendError(w, http.StatusBadRequest, "invalid JSON", err, service.cfg)
		return
	}
	if len(req.Uuids) == 0 {
		shared.SendError(w, http.StatusBadRequest, "uuids cannot be empty", nil, service.cfg)
		return
	}
	for _, hostUUID := range req.Uuids {
		if _, err := uuid.Parse(hostUUID); err != nil {
			shared.SendError(w, http.StatusBadRequest, "invalid UUID format", nil, service.cfg)
			return
		}
	}
	if err := validateUpdateRequest(req.hostUpdateFields); err != nil {
		shared.SendError(w, http.StatusBadRequest, err.Error(), nil, service.cfg)
		return
	}

	err := service.BulkUpdateHosts(r.Context(), req)
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

	w.WriteHeader(http.StatusNoContent)
}

func handleReorderHosts(w http.ResponseWriter, r *http.Request, service *HostService) {
	var req reorderHostsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		shared.SendError(w, http.StatusBadRequest, "invalid JSON", err, service.cfg)
		return
	}
	if len(req.Hosts) == 0 {
		req.Hosts = req.Items
	}
	if len(req.Hosts) == 0 {
		shared.SendError(w, http.StatusBadRequest, "hosts cannot be empty", nil, service.cfg)
		return
	}
	for _, item := range req.Hosts {
		if _, err := uuid.Parse(item.UUID); err != nil {
			shared.SendError(w, http.StatusBadRequest, "invalid UUID format", nil, service.cfg)
			return
		}
	}

	err := service.ReorderHosts(r.Context(), req)
	if err != nil {
		shared.SendAPIError(w, shared.ErrUpdateHostFailed.WithCause(err), service.cfg)
		return
	}

	shared.WriteJSON(w, http.StatusOK, map[string]any{"response": map[string]any{"isUpdated": true}})
}

type cloneHostRequest struct {
	CloneFromUUID string `json:"cloneFromUuid"`
}

func handleCloneHost(w http.ResponseWriter, r *http.Request, service *HostService) {
	var req cloneHostRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		shared.SendError(w, http.StatusBadRequest, "invalid JSON", err, service.cfg)
		return
	}
	if _, err := uuid.Parse(req.CloneFromUUID); err != nil {
		shared.SendError(w, http.StatusBadRequest, "invalid cloneFromUuid format", nil, service.cfg)
		return
	}

	cloned, nodes, squads, err := service.CloneHost(r.Context(), req.CloneFromUUID)
	if err != nil {
		if errors.Is(err, errHostNotFound) {
			shared.SendAPIError(w, shared.ErrHostNotFound, service.cfg)
			return
		}
		shared.SendAPIError(w, shared.ErrCreateHostFailed.WithCause(err), service.cfg)
		return
	}

	result := mapHostRecordToAPI(cloned, nodes, squads)
	shared.WriteJSON(w, http.StatusCreated, map[string]any{"response": result})
}

func handleBulkEnableHosts(w http.ResponseWriter, r *http.Request, service *HostService) {
	var req bulkUUIDsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		shared.SendError(w, http.StatusBadRequest, "invalid JSON", err, service.cfg)
		return
	}
	if len(req.UUIDs) == 0 {
		shared.SendError(w, http.StatusBadRequest, "uuids cannot be empty", nil, service.cfg)
		return
	}

	err := service.BulkEnableHosts(r.Context(), req)
	if err != nil {
		shared.SendAPIError(w, shared.ErrUpdateHostFailed.WithCause(err), service.cfg)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func handleBulkDisableHosts(w http.ResponseWriter, r *http.Request, service *HostService) {
	var req bulkUUIDsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		shared.SendError(w, http.StatusBadRequest, "invalid JSON", err, service.cfg)
		return
	}
	if len(req.UUIDs) == 0 {
		shared.SendError(w, http.StatusBadRequest, "uuids cannot be empty", nil, service.cfg)
		return
	}

	err := service.BulkDisableHosts(r.Context(), req)
	if err != nil {
		shared.SendAPIError(w, shared.ErrUpdateHostFailed.WithCause(err), service.cfg)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func handleBulkDeleteHosts(w http.ResponseWriter, r *http.Request, service *HostService) {
	var req bulkUUIDsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		shared.SendError(w, http.StatusBadRequest, "invalid JSON", err, service.cfg)
		return
	}
	if len(req.UUIDs) == 0 {
		shared.SendError(w, http.StatusBadRequest, "uuids cannot be empty", nil, service.cfg)
		return
	}

	err := service.BulkDeleteHosts(r.Context(), req)
	if err != nil {
		shared.SendAPIError(w, shared.ErrDeleteHostFailed.WithCause(err), service.cfg)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func handleBulkSetInbound(w http.ResponseWriter, r *http.Request, service *HostService) {
	var req setInboundRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		shared.SendError(w, http.StatusBadRequest, "invalid JSON", err, service.cfg)
		return
	}
	if len(req.UUIDs) == 0 {
		shared.SendError(w, http.StatusBadRequest, "uuids cannot be empty", nil, service.cfg)
		return
	}
	if _, err := uuid.Parse(req.ConfigProfileUUID); err != nil {
		shared.SendError(w, http.StatusBadRequest, "invalid configProfileUuid", nil, service.cfg)
		return
	}
	if _, err := uuid.Parse(req.ConfigProfileInboundUUID); err != nil {
		shared.SendError(w, http.StatusBadRequest, "invalid configProfileInboundUuid", nil, service.cfg)
		return
	}

	err := service.BulkSetInbound(r.Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, errConfigProfileNotFound):
			shared.SendAPIError(w, shared.ErrConfigProfileNotFound, service.cfg)
		case errors.Is(err, errConfigProfileInboundNotFound):
			shared.SendAPIError(w, shared.ErrConfigProfileInboundNotFoundInProfile, service.cfg)
		default:
			shared.SendAPIError(w, shared.ErrUpdateHostFailed.WithCause(err), service.cfg)
		}
		return
	}

	handleGetHosts(w, r, service)
}

func handleBulkSetPort(w http.ResponseWriter, r *http.Request, service *HostService) {
	var req setPortRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		shared.SendError(w, http.StatusBadRequest, "invalid JSON", err, service.cfg)
		return
	}
	if len(req.UUIDs) == 0 {
		shared.SendError(w, http.StatusBadRequest, "uuids cannot be empty", nil, service.cfg)
		return
	}
	if req.Port < 1 || req.Port > 65535 {
		shared.SendError(w, http.StatusBadRequest, "invalid port", nil, service.cfg)
		return
	}

	err := service.BulkSetPort(r.Context(), req)
	if err != nil {
		shared.SendAPIError(w, shared.ErrUpdateHostFailed.WithCause(err), service.cfg)
		return
	}

	handleGetHosts(w, r, service)
}
