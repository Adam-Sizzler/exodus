package configprofiles

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"exodus/internal/config"
	"exodus/internal/httpapi/shared"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
)

func handleGetConfigProfiles(w http.ResponseWriter, r *http.Request, service *ConfigProfileService) {
	records, err := service.repo.getAllConfigProfileRecords(r.Context())
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to fetch config profiles", err, service.cfg)
		return
	}

	response, err := buildConfigProfileResponses(service, records)
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to build config profiles response", err, service.cfg)
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
			shared.SendError(w, http.StatusNotFound, "config profile not found", nil, service.cfg)
			return
		}
		shared.SendError(w, http.StatusInternalServerError, "failed to fetch config profile", err, service.cfg)
		return
	}

	response, err := buildConfigProfileResponses(service, []configProfileRecord{record})
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to build config profile response", err, service.cfg)
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
			shared.SendError(w, http.StatusNotFound, "config profile not found", nil, service.cfg)
			return
		}
		shared.SendError(w, http.StatusInternalServerError, "failed to fetch config profile", err, service.cfg)
		return
	}

	inbounds, err := service.repo.getConfigProfileInboundsMap(r.Context(), []string{profileUUID})
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to fetch config profile inbounds", err, service.cfg)
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
		shared.SendError(w, http.StatusInternalServerError, "failed to fetch config profiles", err, service.cfg)
		return
	}
	profileUUIDs := make([]string, 0, len(records))
	for _, record := range records {
		profileUUIDs = append(profileUUIDs, record.UUID)
	}
	inboundsMap, err := service.repo.getConfigProfileInboundsMap(r.Context(), profileUUIDs)
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to fetch inbounds", err, service.cfg)
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

	response, err := buildConfigProfileResponses(service, []configProfileRecord{result})
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to build config profile response", err, service.cfg)
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

	response, err := buildConfigProfileResponses(service, []configProfileRecord{result})
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to build config profile response", err, service.cfg)
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
		shared.SendError(w, http.StatusInternalServerError, "failed to reorder config profiles", err, service.cfg)
		return
	}

	handleGetConfigProfiles(w, r, service)
}

func handleConfigProfileWriteError(w http.ResponseWriter, err error, cfg *config.BackendConfig) {
	switch {
	case errors.Is(err, errConfigProfileNotFound):
		shared.SendError(w, http.StatusNotFound, "config profile not found", nil, cfg)
	case strings.Contains(err.Error(), "no fields to update"):
		shared.SendError(w, http.StatusBadRequest, "no fields to update", nil, cfg)
	default:
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			switch pgErr.ConstraintName {
			case "config_profiles_name_key":
				shared.SendError(w, http.StatusConflict, errMessageConfigProfileNameAlreadyExists, nil, cfg)
			case "config_profile_inbounds_tag_key":
				shared.SendError(w, http.StatusConflict, errMessageInboundTagsMustBeUnique, nil, cfg)
			default:
				shared.SendError(w, http.StatusConflict, errMessageInboundTagsMustBeUnique, nil, cfg)
			}
			return
		}
		if strings.Contains(err.Error(), "config_profiles_name_key") || strings.Contains(err.Error(), "config_profiles.name") {
			shared.SendError(w, http.StatusConflict, errMessageConfigProfileNameAlreadyExists, nil, cfg)
			return
		}
		if strings.Contains(err.Error(), "config_profile_inbounds_tag_key") || strings.Contains(err.Error(), "config_profile_inbounds.tag") {
			shared.SendError(w, http.StatusConflict, errMessageInboundTagsMustBeUnique, nil, cfg)
			return
		}
		if strings.Contains(err.Error(), "UNIQUE constraint failed") || strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			shared.SendError(w, http.StatusConflict, errMessageInboundTagsMustBeUnique, nil, cfg)
			return
		}
		shared.SendError(w, http.StatusInternalServerError, "failed to write config profile", err, cfg)
	}
}

func buildConfigProfileResponses(service *ConfigProfileService, records []configProfileRecord) ([]ConfigProfile, error) {
	profileUUIDs := make([]string, 0, len(records))
	for _, record := range records {
		profileUUIDs = append(profileUUIDs, record.UUID)
	}

	inboundsMap, err := service.repo.getConfigProfileInboundsMap(rContext(service), profileUUIDs)
	if err != nil {
		return nil, err
	}
	nodesMap, err := service.repo.getConfigProfileNodesMap(rContext(service), profileUUIDs)
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

func rContext(s *ConfigProfileService) interface {
	context.Context
} {
	return context.Background()
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
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			shared.SendError(w, http.StatusBadRequest, "Snippet name already exists", nil, service.cfg)
			return
		}
		shared.SendError(w, http.StatusInternalServerError, "Create snippet error", err, service.cfg)
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
			shared.SendError(w, http.StatusNotFound, "Snippet not found", nil, service.cfg)
			return
		}
		shared.SendError(w, http.StatusInternalServerError, "Update snippet error", err, service.cfg)
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
			shared.SendError(w, http.StatusNotFound, "Snippet not found", nil, service.cfg)
			return
		}
		shared.SendError(w, http.StatusInternalServerError, "Delete snippet by name error", err, service.cfg)
		return
	}

	writeConfigProfileSnippetsResponse(w, r, service, http.StatusOK)
}

func writeConfigProfileSnippetsResponse(w http.ResponseWriter, r *http.Request, service *ConfigProfileService, status int) {
	snippets, err := service.repo.getSnippets(r.Context())
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "Get snippets error", err, service.cfg)
		return
	}

	shared.WriteJSON(w, status, map[string]any{
		"response": map[string]any{
			"total":    len(snippets),
			"snippets": snippets,
		},
	})
}
