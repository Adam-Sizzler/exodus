package nodes

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"exodus/internal/httpapi/shared"

	"github.com/google/uuid"
)

func handleEnableNode(w http.ResponseWriter, r *http.Request, service *NodeService, nodeUUID string) {
	err := service.EnableNode(r.Context(), nodeUUID)
	if err != nil {
		if errors.Is(err, errNodeNotFound) {
			shared.SendAPIError(w, shared.ErrNodeNotFound, service.cfg)
			return
		}
		shared.SendAPIError(w, shared.ErrEnableNodeFailed.WithCause(err), service.cfg)
		return
	}
	sendUpdatedNodeResponse(w, r, service, nodeUUID)
}

func handleDisableNode(w http.ResponseWriter, r *http.Request, service *NodeService, nodeUUID string) {
	err := service.DisableNode(r.Context(), nodeUUID)
	if err != nil {
		if errors.Is(err, errNodeNotFound) {
			shared.SendAPIError(w, shared.ErrNodeNotFound, service.cfg)
			return
		}
		shared.SendAPIError(w, shared.ErrDisableNodeFailed.WithCause(err), service.cfg)
		return
	}
	sendUpdatedNodeResponse(w, r, service, nodeUUID)
}

func handleRestartNode(w http.ResponseWriter, r *http.Request, service *NodeService, nodeUUID string) {
	req, err := decodeOptionalRestartNodesRequest(r)
	if err != nil {
		shared.SendAPIError(w, shared.ErrInvalidJSON.WithCause(err), service.cfg)
		return
	}

	err = service.RestartNode(r.Context(), nodeUUID, isForceRestartRequested(req))
	if err != nil {
		if errors.Is(err, errNodeNotFound) {
			shared.SendAPIError(w, shared.ErrNodeNotFound, service.cfg)
			return
		}
		shared.SendError(w, http.StatusBadRequest, err.Error(), nil, service.cfg)
		return
	}
	w.WriteHeader(http.StatusAccepted)
}

func handleResetNodeTraffic(w http.ResponseWriter, r *http.Request, service *NodeService, nodeUUID string) {
	err := service.ResetNodeTraffic(r.Context(), nodeUUID)
	if err != nil {
		if errors.Is(err, errNodeNotFound) {
			shared.SendAPIError(w, shared.ErrNodeNotFound, service.cfg)
			return
		}
		shared.SendAPIError(w, shared.ErrResetNodeTrafficFailed.WithCause(err), service.cfg)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func handleRestartAllNodes(w http.ResponseWriter, r *http.Request, service *NodeService) {
	req, err := decodeOptionalRestartNodesRequest(r)
	if err != nil {
		shared.SendAPIError(w, shared.ErrInvalidJSON.WithCause(err), service.cfg)
		return
	}

	err = service.RestartAllNodes(r.Context(), isForceRestartRequested(req))
	if err != nil {
		if errors.Is(err, errNoEnabledNodes) {
			shared.SendError(w, http.StatusBadRequest, err.Error(), nil, service.cfg)
			return
		}
		shared.SendAPIError(w, shared.ErrRestartNodeFailed.WithCause(err), service.cfg)
		return
	}
	w.WriteHeader(http.StatusAccepted)
}

func decodeOptionalRestartNodesRequest(r *http.Request) (restartAllNodesRequest, error) {
	var req restartAllNodesRequest
	if r.Body == nil || r.ContentLength == 0 {
		return req, nil
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		if errors.Is(err, io.EOF) {
			return req, nil
		}
		return req, err
	}
	return req, nil
}

func isForceRestartRequested(req restartAllNodesRequest) bool {
	return req.ForceRestart != nil && *req.ForceRestart
}

func handleReorderNodes(w http.ResponseWriter, r *http.Request, service *NodeService) {
	var req reorderNodesRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		shared.SendAPIError(w, shared.ErrInvalidJSON.WithCause(err), service.cfg)
		return
	}
	if len(req.Nodes) == 0 {
		req.Nodes = req.Items
	}
	if len(req.Nodes) == 0 {
		shared.SendAPIError(w, shared.ErrNodesListCannotBeEmpty, service.cfg)
		return
	}
	for _, item := range req.Nodes {
		if _, err := uuid.Parse(item.UUID); err != nil {
			shared.SendAPIError(w, shared.ErrInvalidUUID, service.cfg)
			return
		}
	}

	err := service.ReorderNodes(r.Context(), req.Nodes)
	if err != nil {
		shared.SendAPIError(w, shared.ErrReorderNodesFailed.WithCause(err), service.cfg)
		return
	}

	handleGetNodes(w, r, service)
}

func sendUpdatedNodeResponse(w http.ResponseWriter, r *http.Request, service *NodeService, nodeUUID string) {
	node, err := service.repo.getNodeByUUID(r.Context(), nodeUUID)
	if err != nil {
		shared.SendAPIError(w, shared.ErrFetchUpdatedNodeFailed.WithCause(err), service.cfg)
		return
	}
	response, err := buildNodeResponses(r.Context(), service.repo, service.cfg, []nodeRecord{node})
	if err != nil {
		shared.SendAPIError(w, shared.ErrBuildNodeResponseFailed.WithCause(err), service.cfg)
		return
	}
	shared.WriteJSON(w, http.StatusOK, map[string]any{"response": response[0]})
}
