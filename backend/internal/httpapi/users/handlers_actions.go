package users

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"exodus/internal/config"
	"exodus/internal/httpapi/shared"
)

type revokeUserSubscriptionRequest struct {
	RevokeOnlyPasswords bool `json:"revokeOnlyPasswords"`
}

func handleEnableUser(w http.ResponseWriter, r *http.Request, service *UserService, userUUID string) {
	err := service.EnableUser(r.Context(), userUUID)
	if err != nil {
		handleUserActionError(w, err, service.cfg, "failed to enable user")
		return
	}
	sendUpdatedUserResponse(w, r, service, userUUID)
}

func handleDisableUser(w http.ResponseWriter, r *http.Request, service *UserService, userUUID string) {
	err := service.DisableUser(r.Context(), userUUID)
	if err != nil {
		handleUserActionError(w, err, service.cfg, "failed to disable user")
		return
	}
	sendUpdatedUserResponse(w, r, service, userUUID)
}

func handleResetUserTraffic(w http.ResponseWriter, r *http.Request, service *UserService, userUUID string) {
	err := service.ResetUserTraffic(r.Context(), userUUID)
	if err != nil {
		handleUserActionError(w, err, service.cfg, "failed to reset user traffic")
		return
	}
	sendUpdatedUserResponse(w, r, service, userUUID)
}

func handleRevokeUserSubscription(w http.ResponseWriter, r *http.Request, service *UserService, userUUID string) {
	req := revokeUserSubscriptionRequest{}
	if r.Body != nil {
		body, _ := io.ReadAll(io.LimitReader(r.Body, 1<<20))
		if len(strings.TrimSpace(string(body))) > 0 {
			if err := json.Unmarshal(body, &req); err != nil {
				shared.SendError(w, http.StatusBadRequest, "invalid JSON", err, service.cfg)
				return
			}
		}
	}

	err := service.RevokeUserSubscription(r.Context(), userUUID, req)
	if err != nil {
		handleUserWriteError(w, err, service.cfg)
		return
	}
	sendUpdatedUserResponse(w, r, service, userUUID)
}

func sendUpdatedUserResponse(w http.ResponseWriter, r *http.Request, service *UserService, userUUID string) {
	record, err := service.repo.getUserRecordByUUID(r.Context(), userUUID)
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to fetch updated user", err, service.cfg)
		return
	}
	response, err := buildUserResponses(r.Context(), service.repo, []userRecord{record}, resolveUsersSubscriptionBase(r.Context(), service.repo.manager, r, service.cfg))
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to build updated user response", err, service.cfg)
		return
	}
	shared.WriteJSON(w, http.StatusOK, map[string]any{"response": response[0]})
}

func handleUserActionError(w http.ResponseWriter, err error, cfg *config.BackendConfig, message string) {
	if errors.Is(err, errUserNotFound) {
		shared.SendError(w, http.StatusNotFound, "user not found", nil, cfg)
		return
	}
	shared.SendError(w, http.StatusInternalServerError, message, err, cfg)
}

func plannedUserStatusForUpdate(record userRecord, req updateUserRequest, now time.Time) (string, bool) {
	if req.Status != nil {
		return strings.ToUpper(strings.TrimSpace(*req.Status)), true
	}

	if req.TrafficLimitBytes != nil && strings.EqualFold(record.Status, "LIMITED") {
		if *req.TrafficLimitBytes == 0 || *req.TrafficLimitBytes > record.TrafficLimitBytes {
			return "ACTIVE", true
		}
	}

	if req.ExpireAt != nil && strings.EqualFold(record.Status, "EXPIRED") {
		parsed, err := time.Parse(time.RFC3339, strings.TrimSpace(*req.ExpireAt))
		if err == nil {
			newExpireAt := parsed.UTC()
			if !newExpireAt.Equal(record.ExpireAt.UTC()) && newExpireAt.After(now.UTC()) {
				return "ACTIVE", true
			}
		}
	}

	return "", false
}

func userConfigPresenceChanges(previousStatus string, nextStatus string) bool {
	previousActive := strings.EqualFold(previousStatus, "ACTIVE")
	nextActive := strings.EqualFold(nextStatus, "ACTIVE")
	return previousActive != nextActive
}

func validateExtendDays(days int) error {
	if days < 1 || days > 9999 {
		return fmt.Errorf("extendDays must be between 1 and 9999")
	}
	return nil
}
