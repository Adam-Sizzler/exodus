package subscriptionconnections

import (
	"encoding/json"
	"errors"
	"net/http"

	"exodus/internal/httpapi/shared"

	"github.com/google/uuid"
)

func handleGetNodes(w http.ResponseWriter, r *http.Request, service *SubscriptionConnectionService) {
	nodes, err := service.repo.getAllNodeRecords(r.Context())
	if err != nil {
		shared.SendAPIError(w, shared.ErrGetAllNodesFailed.WithCause(err), service.cfg)
		return
	}

	providerUUIDs := make([]string, 0, len(nodes))
	for _, record := range nodes {
		if record.ProviderUUID != nil && *record.ProviderUUID != "" {
			providerUUIDs = append(providerUUIDs, *record.ProviderUUID)
		}
	}

	providersMap, err := service.repo.getProviders(r.Context(), dedupeStrings(providerUUIDs))
	if err != nil {
		shared.SendAPIError(w, shared.ErrGetAllNodesFailed.WithCause(err), service.cfg)
		return
	}

	response := buildNodeResponses(nodes, providersMap)
	shared.WriteJSON(w, http.StatusOK, map[string]any{"response": response})
}

func handleGetNode(w http.ResponseWriter, r *http.Request, service *SubscriptionConnectionService, nodeUUID string) {
	node, err := service.repo.getNodeByUUID(r.Context(), nodeUUID)
	if err != nil {
		if errors.Is(err, errNodeNotFound) {
			shared.SendAPIError(w, shared.ErrNodeNotFound, service.cfg)
			return
		}
		shared.SendAPIError(w, shared.ErrGetOneNodeFailed.WithCause(err), service.cfg)
		return
	}

	providerUUIDs := make([]string, 0, 1)
	if node.ProviderUUID != nil && *node.ProviderUUID != "" {
		providerUUIDs = append(providerUUIDs, *node.ProviderUUID)
	}

	providersMap, err := service.repo.getProviders(r.Context(), providerUUIDs)
	if err != nil {
		shared.SendAPIError(w, shared.ErrGetOneNodeFailed.WithCause(err), service.cfg)
		return
	}

	response := buildNodeResponses([]nodeRecord{node}, providersMap)
	shared.WriteJSON(w, http.StatusOK, map[string]any{"response": response[0]})
}

func handleCreateNode(w http.ResponseWriter, r *http.Request, service *SubscriptionConnectionService) {
	var req createNodeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		shared.SendError(w, http.StatusBadRequest, "invalid JSON", err, service.cfg)
		return
	}
	if err := validateCreateRequest(req); err != nil {
		shared.SendError(w, http.StatusBadRequest, err.Error(), nil, service.cfg)
		return
	}

	created, err := service.CreateNode(r.Context(), req)
	if err != nil {
		shared.SendAPIError(w, shared.ErrCreateNodeFailed.WithCause(err), service.cfg)
		return
	}

	providerUUIDs := make([]string, 0, 1)
	if created.ProviderUUID != nil && *created.ProviderUUID != "" {
		providerUUIDs = append(providerUUIDs, *created.ProviderUUID)
	}

	providersMap, err := service.repo.getProviders(r.Context(), providerUUIDs)
	if err != nil {
		shared.SendAPIError(w, shared.ErrCreateNodeFailed.WithCause(err), service.cfg)
		return
	}

	response := buildNodeResponses([]nodeRecord{created}, providersMap)
	shared.WriteJSON(w, http.StatusCreated, map[string]any{"response": response[0]})
}

func handleUpdateNode(w http.ResponseWriter, r *http.Request, service *SubscriptionConnectionService) {
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

	updated, err := service.UpdateNode(r.Context(), req)
	if err != nil {
		if errors.Is(err, errNodeNotFound) {
			shared.SendAPIError(w, shared.ErrNodeNotFound, service.cfg)
			return
		}
		shared.SendAPIError(w, shared.ErrUpdateNodeFailed.WithCause(err), service.cfg)
		return
	}

	providerUUIDs := make([]string, 0, 1)
	if updated.ProviderUUID != nil && *updated.ProviderUUID != "" {
		providerUUIDs = append(providerUUIDs, *updated.ProviderUUID)
	}

	providersMap, err := service.repo.getProviders(r.Context(), providerUUIDs)
	if err != nil {
		shared.SendAPIError(w, shared.ErrUpdateNodeFailed.WithCause(err), service.cfg)
		return
	}

	response := buildNodeResponses([]nodeRecord{updated}, providersMap)
	shared.WriteJSON(w, http.StatusOK, map[string]any{"response": response[0]})
}

func handleDeleteNode(w http.ResponseWriter, r *http.Request, service *SubscriptionConnectionService, nodeUUID string) {
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
