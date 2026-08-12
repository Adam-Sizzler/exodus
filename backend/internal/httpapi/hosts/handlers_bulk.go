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
			shared.SendError(w, http.StatusNotFound, "host not found", nil, service.cfg)
		case errors.Is(err, errConfigProfileNotFound):
			shared.SendError(w, http.StatusBadRequest, "config profile not found", nil, service.cfg)
		case errors.Is(err, errConfigProfileInboundNotFound):
			shared.SendError(w, http.StatusBadRequest, "config profile inbound not found in specified profile", nil, service.cfg)
		case errors.Is(err, errTemplateNotFound):
			shared.SendError(w, http.StatusBadRequest, "subscription template not found", nil, service.cfg)
		case errors.Is(err, errTemplateTypeNotAllowed):
			shared.SendError(w, http.StatusBadRequest, "template type not allowed", nil, service.cfg)
		default:
			shared.SendError(w, http.StatusInternalServerError, "failed to update hosts", err, service.cfg)
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
		shared.SendError(w, http.StatusInternalServerError, "failed to reorder hosts", err, service.cfg)
		return
	}

	shared.WriteJSON(w, http.StatusOK, map[string]any{"response": map[string]any{"isUpdated": true}})
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
		shared.SendError(w, http.StatusInternalServerError, "failed to enable hosts", err, service.cfg)
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
		shared.SendError(w, http.StatusInternalServerError, "failed to disable hosts", err, service.cfg)
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
		shared.SendError(w, http.StatusInternalServerError, "failed to delete hosts", err, service.cfg)
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
			shared.SendError(w, http.StatusBadRequest, "config profile not found", nil, service.cfg)
		case errors.Is(err, errConfigProfileInboundNotFound):
			shared.SendError(w, http.StatusBadRequest, "config profile inbound not found in specified profile", nil, service.cfg)
		default:
			shared.SendError(w, http.StatusInternalServerError, "failed to validate config profile inbound", err, service.cfg)
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
		shared.SendError(w, http.StatusInternalServerError, "failed to set port", err, service.cfg)
		return
	}

	handleGetHosts(w, r, service)
}
