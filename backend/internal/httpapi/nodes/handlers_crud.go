package nodes

import (
	"encoding/json"
	"errors"
	"net/http"

	"exodus/internal/httpapi/shared"

	"github.com/google/uuid"
)

func handleGetNodes(w http.ResponseWriter, r *http.Request, service *NodeService) {
	nodes, err := service.repo.getAllNodeRecords(r.Context())
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to fetch nodes", err, service.cfg)
		return
	}

	response, err := buildNodeResponses(r.Context(), service.repo, service.cfg, nodes)
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to build node response", err, service.cfg)
		return
	}

	shared.WriteJSON(w, http.StatusOK, map[string]any{"response": response})
}

func handleGetNode(w http.ResponseWriter, r *http.Request, service *NodeService, nodeUUID string) {
	node, err := service.repo.getNodeByUUID(r.Context(), nodeUUID)
	if err != nil {
		if errors.Is(err, errNodeNotFound) {
			shared.SendError(w, http.StatusNotFound, "node not found", nil, service.cfg)
			return
		}
		shared.SendError(w, http.StatusInternalServerError, "failed to fetch node", err, service.cfg)
		return
	}

	response, err := buildNodeResponses(r.Context(), service.repo, service.cfg, []nodeRecord{node})
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to build node response", err, service.cfg)
		return
	}

	shared.WriteJSON(w, http.StatusOK, map[string]any{"response": response[0]})
}

func handleCreateNode(w http.ResponseWriter, r *http.Request, service *NodeService) {
	var req createNodeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		shared.SendError(w, http.StatusBadRequest, "invalid JSON", err, service.cfg)
		return
	}
	if err := validateCreateRequest(req); err != nil {
		shared.SendError(w, http.StatusBadRequest, err.Error(), nil, service.cfg)
		return
	}
	if err := service.repo.ensureConfigProfileInbounds(r.Context(), req.ConfigProfile.ActiveConfigProfileUUID, req.ConfigProfile.ActiveInbounds); err != nil {
		handleConfigProfileValidationError(w, err, service.cfg)
		return
	}

	node, err := service.CreateNode(r.Context(), req)
	if err != nil {
		errStr := err.Error()
		if errors.Is(err, errNodeNameExists) || errors.Is(err, errNodeAddressExists) {
			shared.SendError(w, http.StatusBadRequest, errStr, err, service.cfg)
		} else {
			shared.SendError(w, http.StatusInternalServerError, "failed to create node", err, service.cfg)
		}
		return
	}

	response, err := buildNodeResponses(r.Context(), service.repo, service.cfg, []nodeRecord{node})
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to build node response", err, service.cfg)
		return
	}
	shared.WriteJSON(w, http.StatusCreated, map[string]any{"response": response[0]})
}

func handleUpdateNode(w http.ResponseWriter, r *http.Request, service *NodeService) {
	var req updateNodeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		shared.SendError(w, http.StatusBadRequest, "invalid JSON", err, service.cfg)
		return
	}
	if _, err := uuid.Parse(req.UUID); err != nil {
		shared.SendError(w, http.StatusBadRequest, "invalid UUID format", nil, service.cfg)
		return
	}
	if err := validateUpdateRequest(req); err != nil {
		shared.SendError(w, http.StatusBadRequest, err.Error(), nil, service.cfg)
		return
	}
	if req.ConfigProfile != nil {
		if err := service.repo.ensureConfigProfileInbounds(r.Context(), req.ConfigProfile.ActiveConfigProfileUUID, req.ConfigProfile.ActiveInbounds); err != nil {
			handleConfigProfileValidationError(w, err, service.cfg)
			return
		}
	}

	node, err := service.UpdateNode(r.Context(), req)
	if err != nil {
		if errors.Is(err, errNodeNotFound) {
			shared.SendError(w, http.StatusNotFound, "node not found", nil, service.cfg)
			return
		}
		errStr := err.Error()
		if errors.Is(err, errNodeNameExists) || errors.Is(err, errNodeAddressExists) {
			shared.SendError(w, http.StatusBadRequest, errStr, err, service.cfg)
		} else {
			shared.SendError(w, http.StatusInternalServerError, "failed to update node", err, service.cfg)
		}
		return
	}

	response, err := buildNodeResponses(r.Context(), service.repo, service.cfg, []nodeRecord{node})
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to build node response", err, service.cfg)
		return
	}
	shared.WriteJSON(w, http.StatusOK, map[string]any{"response": response[0]})
}

func handleDeleteNode(w http.ResponseWriter, r *http.Request, service *NodeService, nodeUUID string) {
	err := service.DeleteNode(r.Context(), nodeUUID)
	if err != nil {
		if errors.Is(err, errNodeNotFound) {
			shared.SendError(w, http.StatusNotFound, "node not found", nil, service.cfg)
			return
		}
		shared.SendError(w, http.StatusInternalServerError, "failed to delete node", err, service.cfg)
		return
	}

	shared.WriteJSON(w, http.StatusOK, map[string]any{"response": map[string]any{"isDeleted": true}})
}
