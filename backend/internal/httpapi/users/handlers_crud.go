package users

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"exodus/internal/httpapi/shared"
)

func handleGetUsers(w http.ResponseWriter, r *http.Request, service *UserService) {
	query, err := parseUsersTableQuery(r)
	if err != nil {
		shared.SendError(w, http.StatusBadRequest, "invalid table query", err, service.cfg)
		return
	}

	records, err := service.repo.getAllUserRecords(r.Context())
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to fetch users", err, service.cfg)
		return
	}

	response, err := buildUserResponses(r.Context(), service.repo, records, resolveUsersSubscriptionBase(r.Context(), service.repo.db, r, service.cfg))
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to build users response", err, service.cfg)
		return
	}

	response = filterUsersTableResponse(response, query.Filters, query.FilterModes)
	sortUsersTableResponse(response, query.Sorting)
	total := len(response)
	response = paginateUsersTableResponse(response, query.Start, query.Size)

	shared.WriteJSON(w, http.StatusOK, map[string]any{
		"response": map[string]any{
			"users": response,
			"total": total,
		},
	})
}

func handleGetUsersStream(w http.ResponseWriter, r *http.Request, service *UserService) {
	q := r.URL.Query()
	size := 250
	if sizeStr := q.Get("size"); sizeStr != "" {
		if parsed, err := strconv.Atoi(sizeStr); err == nil && parsed >= 1 {
			size = parsed
			if size > 1000 {
				size = 1000
			}
		}
	}
	var cursor int64 = 0
	if cursorStr := q.Get("cursor"); cursorStr != "" {
		if parsed, err := strconv.ParseInt(cursorStr, 10, 64); err == nil && parsed >= 0 {
			cursor = parsed
		}
	}

	telegramID := q.Get("telegramId")
	email := q.Get("email")
	tag := q.Get("tag")
	status := q.Get("status")
	trafficLimitStrategy := q.Get("trafficLimitStrategy")
	externalSquadUUID := q.Get("externalSquadUuid")

	records, nextCursor, total, err := service.repo.getUsersStream(r.Context(), cursor, size, telegramID, email, tag, status, trafficLimitStrategy, externalSquadUUID)
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to stream users", err, service.cfg)
		return
	}

	response, err := buildUserResponses(r.Context(), service.repo, records, resolveUsersSubscriptionBase(r.Context(), service.repo.db, r, service.cfg))
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to build users response", err, service.cfg)
		return
	}

	shared.WriteJSON(w, http.StatusOK, map[string]any{
		"response": map[string]any{
			"users":      response,
			"total":      total,
			"nextCursor": nextCursor,
		},
	})
}

func handleGetUser(w http.ResponseWriter, r *http.Request, service *UserService, userUUID string) {
	record, err := service.repo.getUserRecordByUUID(r.Context(), userUUID)
	if err != nil {
		if errors.Is(err, errUserNotFound) {
			shared.SendError(w, http.StatusNotFound, "user not found", nil, service.cfg)
			return
		}
		shared.SendError(w, http.StatusInternalServerError, "failed to fetch user", err, service.cfg)
		return
	}

	response, err := buildUserResponses(r.Context(), service.repo, []userRecord{record}, resolveUsersSubscriptionBase(r.Context(), service.repo.db, r, service.cfg))
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to build user response", err, service.cfg)
		return
	}

	shared.WriteJSON(w, http.StatusOK, map[string]any{"response": response[0]})
}

func handleCreateUser(w http.ResponseWriter, r *http.Request, service *UserService) {
	var req createUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		shared.SendError(w, http.StatusBadRequest, "invalid JSON", err, service.cfg)
		return
	}

	if err := validateCreateUserRequest(req); err != nil {
		shared.SendError(w, http.StatusBadRequest, err.Error(), nil, service.cfg)
		return
	}

	record, err := service.CreateUser(r.Context(), req)
	if err != nil {
		handleUserWriteError(w, err, service.cfg)
		return
	}

	response, err := buildUserResponses(r.Context(), service.repo, []userRecord{record}, resolveUsersSubscriptionBase(r.Context(), service.repo.db, r, service.cfg))
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to build created user response", err, service.cfg)
		return
	}

	shared.WriteJSON(w, http.StatusCreated, map[string]any{"response": response[0]})
}

func handleUpdateUser(w http.ResponseWriter, r *http.Request, service *UserService) {
	var req updateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		shared.SendError(w, http.StatusBadRequest, "invalid JSON", err, service.cfg)
		return
	}
	if err := validateUpdateUserRequest(req); err != nil {
		shared.SendError(w, http.StatusBadRequest, err.Error(), nil, service.cfg)
		return
	}

	updatedRecord, err := service.UpdateUser(r.Context(), req)
	if err != nil {
		if errors.Is(err, errUserNotFound) {
			shared.SendError(w, http.StatusNotFound, "user not found", nil, service.cfg)
			return
		}
		handleUserWriteError(w, err, service.cfg)
		return
	}

	response, err := buildUserResponses(r.Context(), service.repo, []userRecord{updatedRecord}, resolveUsersSubscriptionBase(r.Context(), service.repo.db, r, service.cfg))
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to build updated user response", err, service.cfg)
		return
	}

	shared.WriteJSON(w, http.StatusOK, map[string]any{"response": response[0]})
}

func handleGetUserSubscriptionRequestHistory(w http.ResponseWriter, r *http.Request, service *UserService, userUUID string) {
	records, err := service.repo.getUserSubscriptionRequestHistory(r.Context(), userUUID)
	if err != nil {
		if errors.Is(err, errUserNotFound) {
			shared.SendError(w, http.StatusNotFound, "user not found", nil, service.cfg)
			return
		}
		shared.SendError(w, http.StatusInternalServerError, "failed to fetch user subscription request history", err, service.cfg)
		return
	}

	shared.WriteJSON(w, http.StatusOK, map[string]any{
		"response": map[string]any{
			"records": records,
			"total":   len(records),
		},
	})
}

func handleGetUserAccessibleNodes(w http.ResponseWriter, r *http.Request, service *UserService, userUUID string) {
	userID, canonicalUUID, activeNodes, err := service.repo.getUserAccessibleNodes(r.Context(), userUUID)
	if err != nil {
		if errors.Is(err, errUserNotFound) {
			shared.SendError(w, http.StatusNotFound, "user not found", nil, service.cfg)
			return
		}
		shared.SendError(w, http.StatusInternalServerError, "failed to fetch user accessible nodes", err, service.cfg)
		return
	}

	shared.WriteJSON(w, http.StatusOK, map[string]any{
		"response": map[string]any{
			"userId":      userID,
			"userUuid":    canonicalUUID,
			"activeNodes": activeNodes,
		},
	})
}

func handleDeleteUser(w http.ResponseWriter, r *http.Request, service *UserService, userUUID string) {
	err := service.DeleteUser(r.Context(), userUUID)
	if err != nil {
		if errors.Is(err, errUserNotFound) {
			shared.SendError(w, http.StatusNotFound, "user not found", nil, service.cfg)
			return
		}
		shared.SendError(w, http.StatusInternalServerError, "failed to delete user", err, service.cfg)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
