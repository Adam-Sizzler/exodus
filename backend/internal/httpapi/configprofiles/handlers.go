package configprofiles

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"exodus/internal/config"
	"exodus/internal/httpapi/shared"
	"exodus/internal/util"

	"github.com/google/uuid"
)

func handleGetConfigProfiles(w http.ResponseWriter, r *http.Request, service *ConfigProfileService) {
	records, err := service.repo.getAllConfigProfileRecords(r.Context())
	if err != nil {
		shared.SendAPIError(w, shared.ErrGetConfigProfilesFailed.WithCause(err), service.cfg)
		return
	}

	response, err := buildConfigProfileResponses(r.Context(), service, records)
	if err != nil {
		shared.SendAPIError(w, shared.ErrGetConfigProfilesFailed.WithCause(err), service.cfg)
		return
	}

	shared.WriteJSON(w, http.StatusOK, map[string]any{
		"response": map[string]any{
			"total":          len(response),
			"configProfiles": response,
		},
	})
}

func handleGetConfigProfile(w http.ResponseWriter, r *http.Request, service *ConfigProfileService, profileUUID string) {
	record, err := service.repo.getConfigProfileRecordByUUID(r.Context(), profileUUID)
	if err != nil {
		if errors.Is(err, errConfigProfileNotFound) {
			shared.SendAPIError(w, shared.ErrConfigProfileNotFound, service.cfg)
			return
		}
		shared.SendAPIError(w, shared.ErrGetConfigProfileByUUIDFailed.WithCause(err), service.cfg)
		return
	}

	response, err := buildConfigProfileResponses(r.Context(), service, []configProfileRecord{record})
	if err != nil {
		shared.SendAPIError(w, shared.ErrGetConfigProfileByUUIDFailed.WithCause(err), service.cfg)
		return
	}
	shared.WriteJSON(w, http.StatusOK, map[string]any{"response": response[0]})
}

func handleGetComputedConfigProfile(w http.ResponseWriter, r *http.Request, service *ConfigProfileService, profileUUID string) {
	handleGetConfigProfile(w, r, service, profileUUID)
}

func handleGetConfigProfileInbounds(w http.ResponseWriter, r *http.Request, service *ConfigProfileService, profileUUID string) {
	if _, err := service.repo.getConfigProfileRecordByUUID(r.Context(), profileUUID); err != nil {
		if errors.Is(err, errConfigProfileNotFound) {
			shared.SendAPIError(w, shared.ErrConfigProfileNotFound, service.cfg)
			return
		}
		shared.SendAPIError(w, shared.ErrGetConfigProfileByUUIDFailed.WithCause(err), service.cfg)
		return
	}

	inbounds, err := service.repo.getConfigProfileInboundsMap(r.Context(), []string{profileUUID})
	if err != nil {
		shared.SendAPIError(w, shared.ErrGetConfigProfileByUUIDFailed.WithCause(err), service.cfg)
		return
	}

	shared.WriteJSON(w, http.StatusOK, map[string]any{
		"response": map[string]any{
			"total":    len(inbounds[profileUUID]),
			"inbounds": inbounds[profileUUID],
		},
	})
}

func handleGetAllInbounds(w http.ResponseWriter, r *http.Request, service *ConfigProfileService) {
	records, err := service.repo.getAllConfigProfileRecords(r.Context())
	if err != nil {
		shared.SendAPIError(w, shared.ErrGetConfigProfilesFailed.WithCause(err), service.cfg)
		return
	}
	profileUUIDs := make([]string, 0, len(records))
	for _, record := range records {
		profileUUIDs = append(profileUUIDs, record.UUID)
	}
	inboundsMap, err := service.repo.getConfigProfileInboundsMap(r.Context(), profileUUIDs)
	if err != nil {
		shared.SendAPIError(w, shared.ErrGetConfigProfilesFailed.WithCause(err), service.cfg)
		return
	}
	all := make([]ConfigProfileInbound, 0)
	for _, profileUUID := range profileUUIDs {
		all = append(all, inboundsMap[profileUUID]...)
	}

	shared.WriteJSON(w, http.StatusOK, map[string]any{
		"response": map[string]any{
			"total":    len(all),
			"inbounds": all,
		},
	})
}

func handleCreateConfigProfile(w http.ResponseWriter, r *http.Request, service *ConfigProfileService) {
	var req createConfigProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		shared.SendError(w, http.StatusBadRequest, "invalid JSON", err, service.cfg)
		return
	}
	if err := validateCreateConfigProfileRequest(req); err != nil {
		shared.SendError(w, http.StatusBadRequest, err.Error(), nil, service.cfg)
		return
	}

	result, err := service.CreateConfigProfile(r.Context(), req)
	if err != nil {
		handleConfigProfileWriteError(w, err, service.cfg)
		return
	}

	response, err := buildConfigProfileResponses(r.Context(), service, []configProfileRecord{result})
	if err != nil {
		shared.SendAPIError(w, shared.ErrBuildConfigProfileResponseFailed.WithCause(err), service.cfg)
		return
	}
	shared.WriteJSON(w, http.StatusCreated, map[string]any{"response": response[0]})
}

func handleUpdateConfigProfile(w http.ResponseWriter, r *http.Request, service *ConfigProfileService) {
	var req updateConfigProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		shared.SendError(w, http.StatusBadRequest, "invalid JSON", err, service.cfg)
		return
	}
	if _, err := uuid.Parse(strings.TrimSpace(req.UUID)); err != nil {
		shared.SendError(w, http.StatusBadRequest, "invalid UUID format", nil, service.cfg)
		return
	}
	if err := validateUpdateConfigProfileRequest(req); err != nil {
		shared.SendError(w, http.StatusBadRequest, err.Error(), nil, service.cfg)
		return
	}

	result, err := service.UpdateConfigProfile(r.Context(), req)
	if err != nil {
		handleConfigProfileWriteError(w, err, service.cfg)
		return
	}

	response, err := buildConfigProfileResponses(r.Context(), service, []configProfileRecord{result})
	if err != nil {
		shared.SendAPIError(w, shared.ErrBuildConfigProfileResponseFailed.WithCause(err), service.cfg)
		return
	}
	shared.WriteJSON(w, http.StatusOK, map[string]any{"response": response[0]})
}

func handleDeleteConfigProfile(w http.ResponseWriter, r *http.Request, service *ConfigProfileService, profileUUID string) {
	err := service.DeleteConfigProfile(r.Context(), profileUUID)
	if err != nil {
		handleConfigProfileWriteError(w, err, service.cfg)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func handleReorderConfigProfiles(w http.ResponseWriter, r *http.Request, service *ConfigProfileService) {
	var req reorderConfigProfilesRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		shared.SendError(w, http.StatusBadRequest, "invalid JSON", err, service.cfg)
		return
	}
	if len(req.Items) == 0 {
		shared.SendError(w, http.StatusBadRequest, "items cannot be empty", nil, service.cfg)
		return
	}
	for _, item := range req.Items {
		if _, err := uuid.Parse(item.UUID); err != nil {
			shared.SendError(w, http.StatusBadRequest, "invalid UUID format", nil, service.cfg)
			return
		}
	}

	err := service.ReorderConfigProfiles(r.Context(), req)
	if err != nil {
		shared.SendAPIError(w, shared.ErrReorderConfigProfilesFailed.WithCause(err), service.cfg)
		return
	}

	handleGetConfigProfiles(w, r, service)
}

func handleConfigProfileWriteError(w http.ResponseWriter, err error, cfg *config.BackendConfig) {
	switch {
	case errors.Is(err, errConfigProfileNotFound):
		shared.SendAPIError(w, shared.ErrConfigProfileNotFound, cfg)
	case strings.Contains(err.Error(), "no fields to update"),
		strings.Contains(err.Error(), "duplicate inbound tag"),
		strings.Contains(err.Error(), "all inbounds must have a non-empty tag"),
		strings.Contains(err.Error(), "is not allowed in inbound tag"),
		strings.Contains(err.Error(), "must be between 2 and 30 characters"),
		strings.Contains(err.Error(), "config must be valid JSON"),
		strings.Contains(err.Error(), "inbounds must be an array"),
		strings.Contains(err.Error(), "failed to parse config JSON"):
		shared.SendError(w, http.StatusBadRequest, err.Error(), nil, cfg)
	default:
		switch {
		case util.IsUniqueViolation(err, "config_profiles_name_key"):
			shared.SendAPIError(w, shared.ErrConfigProfileNameAlreadyExists, cfg)
		case util.IsUniqueViolation(err):
			// Matches the original behavior: any other unique violation on
			// this write (named "config_profile_inbounds_tag_key", any other
			// constraint, or a non-Postgres test double) is reported as a
			// duplicate inbound tag rather than a generic failure.
			shared.SendAPIError(w, shared.ErrInboundTagsMustBeUnique, cfg)
		default:
			shared.SendAPIError(w, shared.ErrUpdateConfigProfileFailed.WithCause(err), cfg)
		}
	}
}

func buildConfigProfileResponses(ctx context.Context, service *ConfigProfileService, records []configProfileRecord) ([]ConfigProfile, error) {
	profileUUIDs := make([]string, 0, len(records))
	for _, record := range records {
		profileUUIDs = append(profileUUIDs, record.UUID)
	}

	inboundsMap, err := service.repo.getConfigProfileInboundsMap(ctx, profileUUIDs)
	if err != nil {
		return nil, err
	}
	nodesMap, err := service.repo.getConfigProfileNodesMap(ctx, profileUUIDs)
	if err != nil {
		return nil, err
	}

	response := make([]ConfigProfile, 0, len(records))
	for _, record := range records {
		inbounds := inboundsMap[record.UUID]
		if inbounds == nil {
			inbounds = make([]ConfigProfileInbound, 0)
		}
		nodes := nodesMap[record.UUID]
		if nodes == nil {
			nodes = make([]ConfigProfileNode, 0)
		}

		response = append(response, ConfigProfile{
			UUID:         record.UUID,
			ViewPosition: record.ViewPosition,
			Name:         record.Name,
			Config:       record.Config,
			Inbounds:     inbounds,
			Nodes:        nodes,
			CreatedAt:    record.CreatedAt,
			UpdatedAt:    record.UpdatedAt,
		})
	}
	return response, nil
}

func handleGetConfigProfileSnippets(w http.ResponseWriter, r *http.Request, service *ConfigProfileService) {
	writeConfigProfileSnippetsResponse(w, r, service, http.StatusOK)
}

func handleCreateConfigProfileSnippet(w http.ResponseWriter, r *http.Request, service *ConfigProfileService) {
	var req configProfileSnippetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		shared.SendError(w, http.StatusBadRequest, "invalid request body", err, service.cfg)
		return
	}

	req.Name = strings.TrimSpace(req.Name)
	if req.Name == "" {
		shared.SendError(w, http.StatusBadRequest, "name is required", nil, service.cfg)
		return
	}
	if len(req.Name) > 255 {
		shared.SendError(w, http.StatusBadRequest, "name must be 255 characters or less", nil, service.cfg)
		return
	}
	if err := validateConfigProfileSnippet(req.Snippet); err != nil {
		if errors.Is(err, errConfigProfileSnippetEmpty) {
			shared.SendError(w, http.StatusBadRequest, "Snippet cannot be empty", nil, service.cfg)
			return
		}
		if errors.Is(err, errConfigProfileSnippetContainsEmptyObjects) {
			shared.SendError(w, http.StatusBadRequest, "Snippet cannot contain empty objects", nil, service.cfg)
			return
		}
		shared.SendError(w, http.StatusBadRequest, "snippet must be valid JSON", nil, service.cfg)
		return
	}

	if err := service.CreateSnippet(r.Context(), req); err != nil {
		if util.IsUniqueViolation(err) {
			shared.SendAPIError(w, shared.ErrSnippetNameAlreadyExists, service.cfg)
			return
		}
		shared.SendAPIError(w, shared.ErrCreateConfigProfileFailed.WithCause(err), service.cfg)
		return
	}

	writeConfigProfileSnippetsResponse(w, r, service, http.StatusCreated)
}

func handleUpdateConfigProfileSnippet(w http.ResponseWriter, r *http.Request, service *ConfigProfileService) {
	var req configProfileSnippetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		shared.SendError(w, http.StatusBadRequest, "invalid request body", err, service.cfg)
		return
	}

	req.Name = strings.TrimSpace(req.Name)
	if req.Name == "" {
		shared.SendError(w, http.StatusBadRequest, "name is required", nil, service.cfg)
		return
	}
	if len(req.Name) > 255 {
		shared.SendError(w, http.StatusBadRequest, "name must be 255 characters or less", nil, service.cfg)
		return
	}
	if err := validateConfigProfileSnippet(req.Snippet); err != nil {
		if errors.Is(err, errConfigProfileSnippetEmpty) {
			shared.SendError(w, http.StatusBadRequest, "Snippet cannot be empty", nil, service.cfg)
			return
		}
		if errors.Is(err, errConfigProfileSnippetContainsEmptyObjects) {
			shared.SendError(w, http.StatusBadRequest, "Snippet cannot contain empty objects", nil, service.cfg)
			return
		}
		shared.SendError(w, http.StatusBadRequest, "snippet must be valid JSON", nil, service.cfg)
		return
	}

	if err := service.UpdateSnippet(r.Context(), req); err != nil {
		if errors.Is(err, errConfigProfileSnippetNotFound) {
			shared.SendAPIError(w, shared.ErrSnippetNotFound, service.cfg)
			return
		}
		shared.SendAPIError(w, shared.ErrUpdateConfigProfileFailed.WithCause(err), service.cfg)
		return
	}

	writeConfigProfileSnippetsResponse(w, r, service, http.StatusOK)
}

func handleDeleteConfigProfileSnippet(w http.ResponseWriter, r *http.Request, service *ConfigProfileService) {
	var req struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		shared.SendError(w, http.StatusBadRequest, "invalid request body", err, service.cfg)
		return
	}

	req.Name = strings.TrimSpace(req.Name)
	if req.Name == "" {
		shared.SendError(w, http.StatusBadRequest, "name is required", nil, service.cfg)
		return
	}

	if err := service.DeleteSnippet(r.Context(), req.Name); err != nil {
		if errors.Is(err, errConfigProfileSnippetNotFound) {
			shared.SendAPIError(w, shared.ErrSnippetNotFound, service.cfg)
			return
		}
		shared.SendAPIError(w, shared.ErrDeleteConfigProfileFailed.WithCause(err), service.cfg)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func writeConfigProfileSnippetsResponse(w http.ResponseWriter, r *http.Request, service *ConfigProfileService, status int) {
	snippets, err := service.repo.getSnippets(r.Context())
	if err != nil {
		shared.SendAPIError(w, shared.ErrGetConfigProfilesFailed.WithCause(err), service.cfg)
		return
	}

	shared.WriteJSON(w, status, map[string]any{
		"response": map[string]any{
			"total":    len(snippets),
			"snippets": snippets,
		},
	})
}
