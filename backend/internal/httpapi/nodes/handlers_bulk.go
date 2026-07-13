package nodes

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"exodus/internal/httpapi/shared"
)

func handleBulkProfileModification(w http.ResponseWriter, r *http.Request, service *NodeService) {
	var req bulkProfileModificationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		shared.SendError(w, http.StatusBadRequest, "invalid JSON", err, service.cfg)
		return
	}
	if err := validateUUIDs(req.UUIDs); err != nil {
		shared.SendError(w, http.StatusBadRequest, err.Error(), nil, service.cfg)
		return
	}
	if err := service.repo.ensureConfigProfileInbounds(r.Context(), req.ConfigProfile.ActiveConfigProfileUUID, req.ConfigProfile.ActiveInbounds); err != nil {
		handleConfigProfileValidationError(w, err, service.cfg)
		return
	}

	err := service.BulkProfileModification(r.Context(), req)
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to modify nodes profile", err, service.cfg)
		return
	}

	shared.WriteJSON(w, http.StatusOK, map[string]any{"response": map[string]any{"eventSent": true}})
}

func handleBulkNodesUpdate(w http.ResponseWriter, r *http.Request, service *NodeService) {
	var req bulkUpdateNodesRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		shared.SendError(w, http.StatusBadRequest, "invalid JSON", err, service.cfg)
		return
	}
	if err := validateBulkUpdateRequest(req); err != nil {
		shared.SendError(w, http.StatusBadRequest, err.Error(), nil, service.cfg)
		return
	}

	clauses := make([]string, 0, 5)
	args := make([]any, 0, 6)
	add := func(column string, value any) {
		clauses = append(clauses, fmt.Sprintf("%s = ?", column))
		args = append(args, value)
	}

	fields := req.Fields
	if fields.CountryCode != nil {
		add("country_code", strings.ToUpper(strings.TrimSpace(*fields.CountryCode)))
	}
	if fields.ConsumptionMultiplier != nil {
		add("consumption_multiplier", toNanoMultiplier(*fields.ConsumptionMultiplier))
	}
	if fields.NodeConsumptionMultiplier != nil {
		add("node_consumption_multiplier", toNanoMultiplier(*fields.NodeConsumptionMultiplier))
	}
	if fields.ProviderUUID.Set {
		if fields.ProviderUUID.Value == nil || strings.TrimSpace(*fields.ProviderUUID.Value) == "" {
			clauses = append(clauses, "provider_uuid = NULL")
		} else {
			add("provider_uuid", strings.TrimSpace(*fields.ProviderUUID.Value))
		}
	}
	if fields.Tags != nil {
		add("tags", normalizeTags(*fields.Tags))
	}
	if fields.Note.Set {
		if fields.Note.Value == nil || strings.TrimSpace(*fields.Note.Value) == "" {
			clauses = append(clauses, "note = NULL")
		} else {
			add("note", strings.TrimSpace(*fields.Note.Value))
		}
	}
	if fields.ActivePluginUUID.Set {
		if fields.ActivePluginUUID.Value == nil || strings.TrimSpace(*fields.ActivePluginUUID.Value) == "" {
			clauses = append(clauses, "active_plugin_uuid = NULL")
		} else {
			add("active_plugin_uuid", strings.TrimSpace(*fields.ActivePluginUUID.Value))
		}
	}

	if len(clauses) > 0 {
		err := service.BulkNodesUpdate(r.Context(), req, clauses, args)
		if err != nil {
			shared.SendError(w, http.StatusInternalServerError, "failed to update nodes", err, service.cfg)
			return
		}
	}

	shared.WriteJSON(w, http.StatusOK, map[string]any{"response": map[string]any{"eventSent": true}})
}

func handleBulkNodesActions(w http.ResponseWriter, r *http.Request, service *NodeService) {
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
		if errors.Is(err, errNodeNotFound) {
			shared.SendError(w, http.StatusNotFound, "node not found", nil, service.cfg)
			return
		}
		shared.SendError(w, http.StatusInternalServerError, "failed to perform bulk actions", err, service.cfg)
		return
	}

	shared.WriteJSON(w, http.StatusOK, map[string]any{"response": map[string]any{"eventSent": true}})
}
