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
			shared.SendError(w, http.StatusNotFound, "node not found", nil, service.cfg)
			return
		}
		shared.SendError(w, http.StatusInternalServerError, "failed to enable node", err, service.cfg)
		return
	}
	sendUpdatedNodeResponse(w, r, service, nodeUUID)
}

func handleDisableNode(w http.ResponseWriter, r *http.Request, service *NodeService, nodeUUID string) {
	err := service.DisableNode(r.Context(), nodeUUID)
	if err != nil {
		if errors.Is(err, errNodeNotFound) {
			shared.SendError(w, http.StatusNotFound, "node not found", nil, service.cfg)
			return
		}
		shared.SendError(w, http.StatusInternalServerError, "failed to disable node", err, service.cfg)
		return
	}
	sendUpdatedNodeResponse(w, r, service, nodeUUID)
}

func handleRestartNode(w http.ResponseWriter, r *http.Request, service *NodeService, nodeUUID string) {
	req, err := decodeOptionalRestartNodesRequest(r)
	if err != nil {
		shared.SendError(w, http.StatusBadRequest, "invalid JSON", err, service.cfg)
		return
	}

	err = service.RestartNode(r.Context(), nodeUUID, isForceRestartRequested(req))
	if err != nil {
		if errors.Is(err, errNodeNotFound) {
			shared.SendError(w, http.StatusNotFound, "node not found", nil, service.cfg)
			return
		}
		shared.SendError(w, http.StatusBadRequest, err.Error(), nil, service.cfg)
		return
	}
	shared.WriteJSON(w, http.StatusOK, map[string]any{"response": map[string]any{"eventSent": true}})
}

func handleResetNodeTraffic(w http.ResponseWriter, r *http.Request, service *NodeService, nodeUUID string) {
	err := service.ResetNodeTraffic(r.Context(), nodeUUID)
	if err != nil {
		if errors.Is(err, errNodeNotFound) {
			shared.SendError(w, http.StatusNotFound, "node not found", nil, service.cfg)
			return
		}
		shared.SendError(w, http.StatusInternalServerError, "failed to reset node traffic", err, service.cfg)
		return
	}

	shared.WriteJSON(w, http.StatusOK, map[string]any{"response": map[string]any{"eventSent": true}})
}

func handleRestartAllNodes(w http.ResponseWriter, r *http.Request, service *NodeService) {
	req, err := decodeOptionalRestartNodesRequest(r)
	if err != nil {
		shared.SendError(w, http.StatusBadRequest, "invalid JSON", err, service.cfg)
		return
	}

	err = service.RestartAllNodes(r.Context(), isForceRestartRequested(req))
	if err != nil {
		if errors.Is(err, errNoEnabledNodes) {
			shared.SendError(w, http.StatusBadRequest, err.Error(), nil, service.cfg)
			return
		}
		shared.SendError(w, http.StatusInternalServerError, "failed to restart all nodes", err, service.cfg)
		return
	}
	shared.WriteJSON(w, http.StatusOK, map[string]any{"response": map[string]any{"eventSent": true}})
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
		shared.SendError(w, http.StatusBadRequest, "invalid JSON", err, service.cfg)
		return
	}
	if len(req.Nodes) == 0 {
		shared.SendError(w, http.StatusBadRequest, "nodes cannot be empty", nil, service.cfg)
		return
	}
	for _, item := range req.Nodes {
		if _, err := uuid.Parse(item.UUID); err != nil {
			shared.SendError(w, http.StatusBadRequest, "invalid UUID format", nil, service.cfg)
			return
		}
	}

	err := service.ReorderNodes(r.Context(), req.Nodes)
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to reorder nodes", err, service.cfg)
		return
	}

	handleGetNodes(w, r, service)
}

func sendUpdatedNodeResponse(w http.ResponseWriter, r *http.Request, service *NodeService, nodeUUID string) {
	node, err := service.repo.getNodeByUUID(r.Context(), nodeUUID)
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to fetch updated node", err, service.cfg)
		return
	}
	response, err := buildNodeResponses(r.Context(), service.repo, service.cfg, []nodeRecord{node})
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to build node response", err, service.cfg)
		return
	}
	shared.WriteJSON(w, http.StatusOK, map[string]any{"response": response[0]})
}
