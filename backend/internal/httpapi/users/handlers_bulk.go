package users

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"exodus/internal/httpapi/shared"
)

func resolveBulkTargetUUIDs(ctx context.Context, repo *UserRepository, userIDs []int64) ([]string, error) {
	if len(userIDs) == 0 {
		return nil, fmt.Errorf("userIds cannot be empty")
	}
	if len(userIDs) > 500 {
		return nil, fmt.Errorf("userIds cannot contain more than 500 items")
	}
	return repo.resolveUUIDsByUserIDs(ctx, userIDs)
}

func handleBulkDeleteUsers(w http.ResponseWriter, r *http.Request, service *UserService) {
	var req bulkDeleteUsersRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		shared.SendError(w, http.StatusBadRequest, "invalid JSON", err, service.cfg)
		return
	}
	targets, err := resolveBulkTargetUUIDs(r.Context(), service.repo, req.UserIDs)
	if err != nil {
		shared.SendError(w, http.StatusBadRequest, err.Error(), nil, service.cfg)
		return
	}

	err = service.BulkDeleteUsers(r.Context(), targets)
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to delete users", err, service.cfg)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func handleBulkResetUsersTraffic(w http.ResponseWriter, r *http.Request, service *UserService) {
	var req bulkDeleteUsersRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		shared.SendError(w, http.StatusBadRequest, "invalid JSON", err, service.cfg)
		return
	}
	targets, err := resolveBulkTargetUUIDs(r.Context(), service.repo, req.UserIDs)
	if err != nil {
		shared.SendError(w, http.StatusBadRequest, err.Error(), nil, service.cfg)
		return
	}

	affectedRows, err := service.BulkResetUsersTraffic(r.Context(), targets)
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to reset users traffic", err, service.cfg)
		return
	}

	shared.WriteJSON(w, http.StatusOK, map[string]any{"response": map[string]any{"affectedRows": affectedRows}})
}

func handleBulkDeleteUsersByStatus(w http.ResponseWriter, r *http.Request, service *UserService) {
	var req struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		shared.SendError(w, http.StatusBadRequest, "invalid JSON", err, service.cfg)
		return
	}

	status := strings.ToUpper(strings.TrimSpace(req.Status))
	if !isValidUserStatus(status) {
		shared.SendError(w, http.StatusBadRequest, "invalid status", nil, service.cfg)
		return
	}

	affectedRows, err := service.BulkDeleteUsersByStatus(r.Context(), status)
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to delete users by status", err, service.cfg)
		return
	}

	shared.WriteJSON(w, http.StatusOK, map[string]any{"response": map[string]any{"affectedRows": affectedRows}})
}

func handleBulkAllResetUsersTraffic(w http.ResponseWriter, r *http.Request, service *UserService) {
	affectedRows, err := service.BulkAllResetUsersTraffic(r.Context())
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to reset all users traffic", err, service.cfg)
		return
	}

	shared.WriteJSON(w, http.StatusOK, map[string]any{"response": map[string]any{"eventSent": affectedRows > 0}})
}

func handleBulkExtendUsersExpirationDate(w http.ResponseWriter, r *http.Request, service *UserService) {
	var req bulkExtendExpirationDateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		shared.SendError(w, http.StatusBadRequest, "invalid JSON", err, service.cfg)
		return
	}
	targets, err := resolveBulkTargetUUIDs(r.Context(), service.repo, req.UserIDs)
	if err != nil {
		shared.SendError(w, http.StatusBadRequest, err.Error(), nil, service.cfg)
		return
	}
	if err := validateExtendDays(req.ExtendDays); err != nil {
		shared.SendError(w, http.StatusBadRequest, err.Error(), nil, service.cfg)
		return
	}

	_, err = service.BulkExtendUsersExpirationDate(r.Context(), targets, req.ExtendDays)
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to extend users expiration date", err, service.cfg)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func handleBulkAllExtendUsersExpirationDate(w http.ResponseWriter, r *http.Request, service *UserService) {
	var req bulkAllExtendExpirationDateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		shared.SendError(w, http.StatusBadRequest, "invalid JSON", err, service.cfg)
		return
	}
	if err := validateExtendDays(req.ExtendDays); err != nil {
		shared.SendError(w, http.StatusBadRequest, err.Error(), nil, service.cfg)
		return
	}

	affectedRows, err := service.BulkAllExtendUsersExpirationDate(r.Context(), req.ExtendDays)
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to extend all users expiration date", err, service.cfg)
		return
	}

	shared.WriteJSON(w, http.StatusOK, map[string]any{"response": map[string]any{"eventSent": affectedRows > 0}})
}

func handleBulkUpdateUsers(w http.ResponseWriter, r *http.Request, service *UserService) {
	var req bulkUpdateUsersRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		shared.SendError(w, http.StatusBadRequest, "invalid JSON", err, service.cfg)
		return
	}
	targets, err := resolveBulkTargetUUIDs(r.Context(), service.repo, req.UserIDs)
	if err != nil {
		shared.SendError(w, http.StatusBadRequest, err.Error(), nil, service.cfg)
		return
	}
	if err := validateBulkUpdateUsersFields(req.Fields); err != nil {
		shared.SendError(w, http.StatusBadRequest, err.Error(), nil, service.cfg)
		return
	}

	affectedRows, err := service.BulkUpdateUsers(r.Context(), targets, req.Fields)
	if err != nil {
		handleUserWriteError(w, err, service.cfg)
		return
	}

	shared.WriteJSON(w, http.StatusOK, map[string]any{"response": map[string]any{"affectedRows": affectedRows}})
}

func handleBulkUpdateUsersSquads(w http.ResponseWriter, r *http.Request, service *UserService) {
	var req bulkUpdateUsersSquadsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		shared.SendError(w, http.StatusBadRequest, "invalid JSON", err, service.cfg)
		return
	}
	targets, err := resolveBulkTargetUUIDs(r.Context(), service.repo, req.UserIDs)
	if err != nil {
		shared.SendError(w, http.StatusBadRequest, err.Error(), nil, service.cfg)
		return
	}
	if err := validateUUIDListAllowEmpty(req.ActiveInternalSquads); err != nil {
		shared.SendError(w, http.StatusBadRequest, "invalid activeInternalSquads", err, service.cfg)
		return
	}

	_, err = service.BulkUpdateUsersSquads(r.Context(), targets, req.ActiveInternalSquads)
	if err != nil {
		handleUserWriteError(w, err, service.cfg)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func handleBulkAllUpdateUsers(w http.ResponseWriter, r *http.Request, service *UserService) {
	var req bulkAllUpdateUsersRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		shared.SendError(w, http.StatusBadRequest, "invalid JSON", err, service.cfg)
		return
	}
	if err := validateBulkAllUpdateUsersRequest(req); err != nil {
		shared.SendError(w, http.StatusBadRequest, err.Error(), nil, service.cfg)
		return
	}

	affectedRows, err := service.BulkAllUpdateUsers(r.Context(), req)
	if err != nil {
		handleUserWriteError(w, err, service.cfg)
		return
	}

	shared.WriteJSON(w, http.StatusOK, map[string]any{"response": map[string]any{"eventSent": affectedRows > 0}})
}
