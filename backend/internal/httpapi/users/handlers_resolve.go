package users

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"exodus/internal/httpapi/shared"
)

// handleResolveUser godoc
// @Summary      Resolve a user
// @Description  Resolve user basic info by ID, shortUUID, or username
// @Tags         Users Controller
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body      resolveUserRequest  true  "Identifier to resolve"
// @Success      200   {object}  object{response=resolveUserResponse}
// @Failure      400   {object}  shared.ErrorResponse
// @Failure      404   {object}  shared.ErrorResponse
// @Failure      500   {object}  shared.ErrorResponse
// @Router       /users/resolve [post]
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
		return fmt.Errorf("exactly one of id, shortUuid, or username must be provided")
	}
	return nil
}
