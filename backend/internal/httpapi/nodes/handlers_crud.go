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
		shared.SendAPIError(w, shared.ErrGetAllNodesFailed.WithCause(err), service.cfg)
		return
	}

	response, err := buildNodeResponses(r.Context(), service.repo, service.cfg, nodes)
	if err != nil {
		shared.SendAPIError(w, shared.ErrBuildNodeResponseFailed.WithCause(err), service.cfg)
		return
	}

	shared.WriteJSON(w, http.StatusOK, map[string]any{"response": response})
}

func handleGetNode(w http.ResponseWriter, r *http.Request, service *NodeService, nodeUUID string) {
	node, err := service.repo.getNodeByUUID(r.Context(), nodeUUID)
	if err != nil {
		if errors.Is(err, errNodeNotFound) {
			shared.SendAPIError(w, shared.ErrNodeNotFound, service.cfg)
			return
		}
		shared.SendAPIError(w, shared.ErrGetOneNodeFailed.WithCause(err), service.cfg)
		return
	}

	response, err := buildNodeResponses(r.Context(), service.repo, service.cfg, []nodeRecord{node})
	if err != nil {
		shared.SendAPIError(w, shared.ErrBuildNodeResponseFailed.WithCause(err), service.cfg)
		return
	}

	shared.WriteJSON(w, http.StatusOK, map[string]any{"response": response[0]})
}

func handleCreateNode(w http.ResponseWriter, r *http.Request, service *NodeService) {
	var req createNodeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		shared.SendAPIError(w, shared.ErrInvalidJSON.WithCause(err), service.cfg)
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
		switch {
		case errors.Is(err, errNodeNameExists):
			shared.SendAPIError(w, shared.ErrNodeNameAlreadyExists, service.cfg)
		case errors.Is(err, errNodeAddressExists):
			shared.SendAPIError(w, shared.ErrNodeAddressAlreadyExists, service.cfg)
		default:
			shared.SendAPIError(w, shared.ErrCreateNodeFailed.WithCause(err), service.cfg)
		}
		return
	}

	response, err := buildNodeResponses(r.Context(), service.repo, service.cfg, []nodeRecord{node})
	if err != nil {
		shared.SendAPIError(w, shared.ErrBuildNodeResponseFailed.WithCause(err), service.cfg)
		return
	}
	shared.WriteJSON(w, http.StatusCreated, map[string]any{"response": response[0]})
}

func handleUpdateNode(w http.ResponseWriter, r *http.Request, service *NodeService) {
	var req updateNodeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		shared.SendAPIError(w, shared.ErrInvalidJSON.WithCause(err), service.cfg)
		return
	}
	if _, err := uuid.Parse(req.UUID); err != nil {
		shared.SendAPIError(w, shared.ErrInvalidUUID, service.cfg)
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
			shared.SendAPIError(w, shared.ErrNodeNotFound, service.cfg)
			return
		}
		switch {
		case errors.Is(err, errNodeNameExists):
			shared.SendAPIError(w, shared.ErrNodeNameAlreadyExists, service.cfg)
		case errors.Is(err, errNodeAddressExists):
			shared.SendAPIError(w, shared.ErrNodeAddressAlreadyExists, service.cfg)
		default:
			shared.SendAPIError(w, shared.ErrUpdateNodeFailed.WithCause(err), service.cfg)
		}
		return
	}

	response, err := buildNodeResponses(r.Context(), service.repo, service.cfg, []nodeRecord{node})
	if err != nil {
		shared.SendAPIError(w, shared.ErrBuildNodeResponseFailed.WithCause(err), service.cfg)
		return
	}
	shared.WriteJSON(w, http.StatusOK, map[string]any{"response": response[0]})
}

func handleDeleteNode(w http.ResponseWriter, r *http.Request, service *NodeService, nodeUUID string) {
	err := service.DeleteNode(r.Context(), nodeUUID)
	if err != nil {
		if errors.Is(err, errNodeNotFound) {
			shared.SendAPIError(w, shared.ErrNodeNotFound, service.cfg)
			return
		}
		shared.SendAPIError(w, shared.ErrDeleteNodeFailed.WithCause(err), service.cfg)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
