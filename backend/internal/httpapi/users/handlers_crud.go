package users

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"exodus/internal/httpapi/shared"
)

// handleGetUsers godoc
// @Summary      Get all users using offset-based pagination
// @Description  Get all users using offset-based pagination with filtering and sorting options
// @Tags         Users Controller
// @Produce      json
// @Security     BearerAuth
// @Param        start    query     int     false  "Pagination start index"
// @Param        size     query     int     false  "Page size"
// @Param        filters  query     string  false  "JSON filters array"
// @Param        sorting  query     string  false  "JSON sorting array"
// @Success      200      {object}  UsersListResponseEnvelope
// @Failure      400      {object}  shared.ErrorResponse
// @Failure      500      {object}  shared.ErrorResponse
// @Router       /users [get]
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

// handleGetUsersStream godoc
// @Summary      Get all users using cursor-based (keyset) pagination with filtering options
// @Description  Stream users with cursor-based pagination and filters
// @Tags         Users Controller
// @Produce      json
// @Security     BearerAuth
// @Param        size                  query     int     false  "Batch size (1-1000, default 250)"
// @Param        cursor                query     int     false  "Pagination cursor"
// @Param        telegramId            query     int     false  "Filter by Telegram ID"
// @Param        email                 query     string  false  "Filter by email"
// @Param        tag                   query     string  false  "Filter by tag"
// @Param        status                query     string  false  "Filter by status (ACTIVE, DISABLED, LIMITED, EXPIRED)"
// @Param        trafficLimitStrategy  query     string  false  "Filter by strategy"
// @Param        externalSquadUuid     query     string  false  "Filter by external squad UUID" format(uuid)
// @Success      200                   {object}  UsersStreamResponseEnvelope
// @Failure      500                   {object}  shared.ErrorResponse
// @Router       /users/stream [get]
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

// handleGetUser godoc
// @Summary      Get user by ID
// @Description  Get single user details by user ID
// @Tags         Users Controller
// @Produce      json
// @Security     BearerAuth
// @Param        userId  path      int  true  "Numeric User ID"
// @Success      200     {object}  UserResponseEnvelope
// @Failure      400     {object}  shared.ErrorResponse
// @Failure      404     {object}  shared.ErrorResponse
// @Failure      500     {object}  shared.ErrorResponse
// @Router       /users/{userId} [get]
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

// handleCreateUser godoc
// @Summary      Create a new user
// @Description  Create a new user account
// @Tags         Users Controller
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body      createUserRequest  true  "User creation parameters"
// @Success      201   {object}  UserResponseEnvelope
// @Failure      400   {object}  shared.ErrorResponse
// @Failure      409   {object}  shared.ErrorResponse
// @Failure      500   {object}  shared.ErrorResponse
// @Router       /users [post]
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

// handleUpdateUser godoc
// @Summary      Update a user
// @Description  Update user profile and protocol fields
// @Tags         Users Controller
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body      updateUserRequest  true  "User update fields"
// @Success      200   {object}  UserResponseEnvelope
// @Failure      400   {object}  shared.ErrorResponse
// @Failure      404   {object}  shared.ErrorResponse
// @Failure      500   {object}  shared.ErrorResponse
// @Router       /users [patch]
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

// handleGetUserSubscriptionRequestHistory godoc
// @Summary      Get user subscription request history, recent 24 records
// @Description  Get history of subscription page accesses by user ID
// @Tags         Users Controller
// @Produce      json
// @Security     BearerAuth
// @Param        userId  path      int  true  "Numeric User ID"
// @Success      200     {object}  UserSubscriptionRequestHistoryResponseEnvelope
// @Failure      400     {object}  shared.ErrorResponse
// @Failure      404     {object}  shared.ErrorResponse
// @Failure      500     {object}  shared.ErrorResponse
// @Router       /users/{userId}/subscription-request-history [get]
func handleGetUserSubscriptionRequestHistory(w http.ResponseWriter, r *http.Request, service *UserService, userUUID string) {
	userID, parseErr := strconv.ParseInt(strings.TrimSpace(userUUID), 10, 64)
	if parseErr != nil {
		shared.SendError(w, http.StatusBadRequest, "userId must be numeric", parseErr, service.cfg)
		return
	}
	records, err := service.repo.getUserSubscriptionRequestHistory(r.Context(), userID)
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

// handleGetUserAccessibleNodes godoc
// @Summary      Get user accessible nodes
// @Description  Get all active nodes accessible to user through assigned squads
// @Tags         Users Controller
// @Produce      json
// @Security     BearerAuth
// @Param        userId  path      int  true  "Numeric User ID"
// @Success      200     {object}  UserAccessibleNodesResponseEnvelope
// @Failure      400     {object}  shared.ErrorResponse
// @Failure      404     {object}  shared.ErrorResponse
// @Failure      500     {object}  shared.ErrorResponse
// @Router       /users/{userId}/accessible-nodes [get]
func handleGetUserAccessibleNodes(w http.ResponseWriter, r *http.Request, service *UserService, userUUID string) {
	userID, parseErr := strconv.ParseInt(strings.TrimSpace(userUUID), 10, 64)
	if parseErr != nil {
		shared.SendError(w, http.StatusBadRequest, "userId must be numeric", parseErr, service.cfg)
		return
	}
	activeNodes, err := service.repo.getUserAccessibleNodes(r.Context(), userID)
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
			"activeNodes": activeNodes,
		},
	})
}

// handleDeleteUser godoc
// @Summary      Delete user
// @Description  Delete user by user ID
// @Tags         Users Controller
// @Security     BearerAuth
// @Param        userId  path  int  true  "Numeric User ID"
// @Success      204
// @Failure      400  {object}  shared.ErrorResponse
// @Failure      404  {object}  shared.ErrorResponse
// @Failure      500  {object}  shared.ErrorResponse
// @Router       /users/{userId} [delete]
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
