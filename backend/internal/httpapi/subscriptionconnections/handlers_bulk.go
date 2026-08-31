package subscriptionconnections

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"exodus/internal/httpapi/shared"

	"github.com/google/uuid"
)

func handleEnableNode(w http.ResponseWriter, r *http.Request, service *SubscriptionConnectionService, nodeUUID string) {
	_, err := service.EnableNode(r.Context(), nodeUUID)
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

func handleDisableNode(w http.ResponseWriter, r *http.Request, service *SubscriptionConnectionService, nodeUUID string) {
	_, err := service.DisableNode(r.Context(), nodeUUID)
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

func handleRestartNode(w http.ResponseWriter, r *http.Request, service *SubscriptionConnectionService, nodeUUID string) {
	err := service.RestartNode(r.Context(), nodeUUID)
	if err != nil {
		if errors.Is(err, errNodeNotFound) {
			shared.SendAPIError(w, shared.ErrNodeNotFound, service.cfg)
			return
		}
		shared.SendAPIError(w, shared.ErrRestartNodeFailed.WithCause(err), service.cfg)
		return
	}

	shared.WriteJSON(w, http.StatusOK, map[string]any{"response": map[string]any{"eventSent": true}})
}

func handleResetNodeTraffic(w http.ResponseWriter, r *http.Request, service *SubscriptionConnectionService, nodeUUID string) {
	err := service.ResetNodeTraffic(r.Context(), nodeUUID)
	if err != nil {
		if errors.Is(err, errNodeNotFound) {
			shared.SendAPIError(w, shared.ErrNodeNotFound, service.cfg)
			return
		}
		shared.SendAPIError(w, shared.ErrResetNodeTrafficFailed.WithCause(err), service.cfg)
		return
	}

	shared.WriteJSON(w, http.StatusOK, map[string]any{"response": map[string]any{"eventSent": true}})
}

func handleRestartAllNodes(w http.ResponseWriter, r *http.Request, service *SubscriptionConnectionService) {
	var req restartAllNodesRequest
	if len(strings.TrimSpace(r.Header.Get("Content-Length"))) != 0 || r.ContentLength > 0 {
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			shared.SendError(w, http.StatusBadRequest, "invalid JSON", err, service.cfg)
			return
		}
	}

	err := service.RestartAllNodes(r.Context())
	if err != nil {
		if errors.Is(err, errNoEnabledNodes) {
			shared.SendError(w, http.StatusBadRequest, err.Error(), nil, service.cfg)
			return
		}
		shared.SendAPIError(w, shared.ErrRestartNodeFailed.WithCause(err), service.cfg)
		return
	}

	shared.WriteJSON(w, http.StatusOK, map[string]any{"response": map[string]any{"eventSent": true}})
}

func handleReorderNodes(w http.ResponseWriter, r *http.Request, service *SubscriptionConnectionService) {
	var req reorderNodesRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		shared.SendError(w, http.StatusBadRequest, "invalid JSON", err, service.cfg)
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
			shared.SendError(w, http.StatusBadRequest, "invalid UUID format", nil, service.cfg)
			return
		}
	}

	err := service.ReorderNodes(r.Context(), req)
	if err != nil {
		shared.SendAPIError(w, shared.ErrReorderNodesFailed.WithCause(err), service.cfg)
		return
	}

	handleGetNodes(w, r, service)
}

func handleBulkProfileModification(w http.ResponseWriter, r *http.Request, service *SubscriptionConnectionService) {
	var req bulkProfileModificationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		shared.SendError(w, http.StatusBadRequest, "invalid JSON", err, service.cfg)
		return
	}
	if err := validateUUIDs(req.UUIDs); err != nil {
		shared.SendError(w, http.StatusBadRequest, err.Error(), nil, service.cfg)
		return
	}

	err := service.BulkProfileModification(r.Context(), req)
	if err != nil {
		shared.SendAPIError(w, shared.ErrModifyNodesProfileFailed.WithCause(err), service.cfg)
		return
	}

	shared.WriteJSON(w, http.StatusOK, map[string]any{"response": map[string]any{"eventSent": true}})
}

func handleBulkNodesActions(w http.ResponseWriter, r *http.Request, service *SubscriptionConnectionService) {
	var req bulkNodesActionsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		shared.SendError(w, http.StatusBadRequest, "invalid JSON", err, service.cfg)
		return
	}
	if err := validateUUIDs(req.UUIDs); err != nil {
		shared.SendError(w, http.StatusBadRequest, err.Error(), nil, service.cfg)
		return
	}

	err := service.BulkNodesActions(r.Context(), req)
	if err != nil {
		shared.SendAPIError(w, shared.ErrNodeBulkActionFailed.WithCause(err), service.cfg)
		return
	}

	shared.WriteJSON(w, http.StatusOK, map[string]any{"response": map[string]any{"eventSent": true}})
}

func sendUpdatedNodeResponse(w http.ResponseWriter, r *http.Request, service *SubscriptionConnectionService, nodeUUID string) {
	node, err := service.repo.getNodeByUUID(r.Context(), nodeUUID)
	if err != nil {
		shared.SendAPIError(w, shared.ErrFetchUpdatedNodeFailed.WithCause(err), service.cfg)
		return
	}
	providerUUIDs := make([]string, 0, 1)
	if node.ProviderUUID != nil && *node.ProviderUUID != "" {
		providerUUIDs = append(providerUUIDs, *node.ProviderUUID)
	}
	providersMap, err := service.repo.getProviders(r.Context(), providerUUIDs)
	if err != nil {
		shared.SendAPIError(w, shared.ErrFetchUpdatedNodeFailed.WithCause(err), service.cfg)
		return
	}
	response := buildNodeResponses([]nodeRecord{node}, providersMap)
	shared.WriteJSON(w, http.StatusOK, map[string]any{"response": response[0]})
}
