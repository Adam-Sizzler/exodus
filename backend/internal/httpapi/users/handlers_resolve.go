package users

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"exodus/internal/httpapi/shared"

	"github.com/google/uuid"
)

func handleResolveUser(w http.ResponseWriter, r *http.Request, service *UserService) {
	var req resolveUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		shared.SendError(w, http.StatusBadRequest, "invalid JSON", err, service.cfg)
		return
	}
	if err := validateResolveUserRequest(req); err != nil {
		shared.SendError(w, http.StatusBadRequest, err.Error(), nil, service.cfg)
		return
	}

	user, err := service.repo.resolveUser(r.Context(), req)
	if err != nil {
		if errors.Is(err, errUserNotFound) {
			shared.SendError(w, http.StatusNotFound, "user not found", nil, service.cfg)
			return
		}
		shared.SendError(w, http.StatusInternalServerError, "failed to resolve user", err, service.cfg)
		return
	}

	shared.WriteJSON(w, http.StatusOK, map[string]any{"response": user})
}

func validateResolveUserRequest(req resolveUserRequest) error {
	provided := 0
	if req.UUID != nil {
		provided++
		if _, err := uuid.Parse(strings.TrimSpace(*req.UUID)); err != nil {
			return fmt.Errorf("invalid uuid")
		}
	}
	if req.ID != nil {
		provided++
	}
	if req.ShortUUID != nil {
		provided++
	}
	if req.Username != nil {
		provided++
	}
	if provided != 1 {
		return fmt.Errorf("exactly one of uuid, id, shortUuid, or username must be provided")
	}
	return nil
}
