package users

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"

	"cerberus/backend/config"
	dbmanager "cerberus/backend/db/manager"
	"cerberus/backend/httpapi/shared"
	monitor "cerberus/backend/nodes"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
)

var (
	userUsernameRegex = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
	userTagRegex      = regexp.MustCompile(`^[A-Z0-9_]+$`)

	errUserNotFound        = errors.New("user not found")
	errUsernameExists      = errors.New("username already exists")
	errShortUUIDExists     = errors.New("short uuid already exists")
	errVLESSUUIDExists     = errors.New("vless uuid already exists")
	errExternalSquadAbsent = errors.New("external squad not found")
)

type OptionalString struct {
	Set   bool
	Value *string
}

func (o *OptionalString) UnmarshalJSON(data []byte) error {
	o.Set = true
	if string(data) == "null" {
		o.Value = nil
		return nil
	}
	var value string
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	o.Value = &value
	return nil
}

type OptionalInt struct {
	Set   bool
	Value *int
}

func (o *OptionalInt) UnmarshalJSON(data []byte) error {
	o.Set = true
	if string(data) == "null" {
		o.Value = nil
		return nil
	}
	var value int
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	o.Value = &value
	return nil
}

type OptionalInt64 struct {
	Set   bool
	Value *int64
}

func (o *OptionalInt64) UnmarshalJSON(data []byte) error {
	o.Set = true
	if string(data) == "null" {
		o.Value = nil
		return nil
	}
	var value int64
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	o.Value = &value
	return nil
}

type internalSquadResponse struct {
	UUID string `json:"uuid"`
	Name string `json:"name"`
}

type userTrafficResponse struct {
	UsedTrafficBytes         int64      `json:"usedTrafficBytes"`
	LifetimeUsedTrafficBytes int64      `json:"lifetimeUsedTrafficBytes"`
	OnlineAt                 *time.Time `json:"onlineAt"`
	FirstConnectedAt         *time.Time `json:"firstConnectedAt"`
	LastConnectedNodeUUID    *string    `json:"lastConnectedNodeUuid"`
}

type userAPI struct {
	UUID                   string                  `json:"uuid"`
	ID                     int64                   `json:"id"`
	ShortUUID              string                  `json:"shortUuid"`
	Username               string                  `json:"username"`
	Status                 string                  `json:"status"`
	TrafficLimitBytes      int64                   `json:"trafficLimitBytes"`
	TrafficLimitStrategy   string                  `json:"trafficLimitStrategy"`
	ExpireAt               time.Time               `json:"expireAt"`
	TelegramID             *int64                  `json:"telegramId"`
	Email                  *string                 `json:"email"`
	Description            *string                 `json:"description"`
	Tag                    *string                 `json:"tag"`
	HwidDeviceLimit        *int                    `json:"hwidDeviceLimit"`
	ExternalSquadUUID      *string                 `json:"externalSquadUuid"`
	TrojanPassword         string                  `json:"trojanPassword"`
	VlessUUID              string                  `json:"vlessUuid"`
	SSPassword             string                  `json:"ssPassword"`
	LastTriggeredThreshold int                     `json:"lastTriggeredThreshold"`
	SubRevokedAt           *time.Time              `json:"subRevokedAt"`
	SubLastUserAgent       *string                 `json:"subLastUserAgent"`
	SubLastOpenedAt        *time.Time              `json:"subLastOpenedAt"`
	LastTrafficResetAt     *time.Time              `json:"lastTrafficResetAt"`
	CreatedAt              time.Time               `json:"createdAt"`
	UpdatedAt              time.Time               `json:"updatedAt"`
	SubscriptionURL        string                  `json:"subscriptionUrl"`
	ActiveInternalSquads   []internalSquadResponse `json:"activeInternalSquads"`
	UserTraffic            userTrafficResponse     `json:"userTraffic"`
}

type userRecord struct {
	TID                      int64
	UUID                     string
	ShortUUID                string
	Username                 string
	Status                   string
	TrafficLimitBytes        int64
	TrafficLimitStrategy     string
	ExpireAt                 time.Time
	SubLastUserAgent         *string
	SubLastOpenedAt          *time.Time
	LastTrafficResetAt       *time.Time
	SubRevokedAt             *time.Time
	TrojanPassword           string
	VlessUUID                string
	SSPassword               string
	Description              *string
	Tag                      *string
	TelegramID               *int64
	Email                    *string
	HwidDeviceLimit          *int
	ExternalSquadUUID        *string
	LastTriggeredThreshold   int
	CreatedAt                time.Time
	UpdatedAt                time.Time
	UsedTrafficBytes         int64
	LifetimeUsedTrafficBytes int64
	OnlineAt                 *time.Time
	LastConnectedNodeUUID    *string
	FirstConnectedAt         *time.Time
}

type createUserRequest struct {
	UUID                 *string  `json:"uuid,omitempty"`
	Username             string   `json:"username"`
	Status               *string  `json:"status,omitempty"`
	ShortUUID            *string  `json:"shortUuid,omitempty"`
	TrojanPassword       *string  `json:"trojanPassword,omitempty"`
	VlessUUID            *string  `json:"vlessUuid,omitempty"`
	SSPassword           *string  `json:"ssPassword,omitempty"`
	TrafficLimitBytes    *int64   `json:"trafficLimitBytes,omitempty"`
	TrafficLimitStrategy *string  `json:"trafficLimitStrategy,omitempty"`
	ExpireAt             string   `json:"expireAt"`
	CreatedAt            *string  `json:"createdAt,omitempty"`
	LastTrafficResetAt   *string  `json:"lastTrafficResetAt,omitempty"`
	Description          *string  `json:"description,omitempty"`
	Tag                  *string  `json:"tag,omitempty"`
	TelegramID           *int64   `json:"telegramId,omitempty"`
	Email                *string  `json:"email,omitempty"`
	HwidDeviceLimit      *int     `json:"hwidDeviceLimit,omitempty"`
	ActiveInternalSquads []string `json:"activeInternalSquads,omitempty"`
	ExternalSquadUUID    *string  `json:"externalSquadUuid,omitempty"`
}

type updateUserRequest struct {
	UUID                 *string        `json:"uuid,omitempty"`
	Username             *string        `json:"username,omitempty"`
	Status               *string        `json:"status,omitempty"`
	TrafficLimitBytes    *int64         `json:"trafficLimitBytes,omitempty"`
	TrafficLimitStrategy *string        `json:"trafficLimitStrategy,omitempty"`
	ExpireAt             *string        `json:"expireAt,omitempty"`
	Description          OptionalString `json:"description,omitempty"`
	Tag                  OptionalString `json:"tag,omitempty"`
	TelegramID           OptionalInt64  `json:"telegramId,omitempty"`
	Email                OptionalString `json:"email,omitempty"`
	HwidDeviceLimit      OptionalInt    `json:"hwidDeviceLimit,omitempty"`
	ActiveInternalSquads *[]string      `json:"activeInternalSquads,omitempty"`
	ExternalSquadUUID    OptionalString `json:"externalSquadUuid,omitempty"`
}

type bulkDeleteUsersRequest struct {
	UUIDs []string `json:"uuids"`
}

func UsersHandler(manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handleGetUsers(w, r, manager, cfg)
		case http.MethodPost:
			handleCreateUser(w, r, manager, cfg)
		case http.MethodPatch:
			handleUpdateUser(w, r, manager, cfg)
		default:
			shared.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
		}
	}
}

func UserByUUIDHandler(manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := trimUsersPath(r.URL.Path, "/")
		if path == "" {
			switch r.Method {
			case http.MethodGet:
				handleGetUsers(w, r, manager, cfg)
			case http.MethodPost:
				handleCreateUser(w, r, manager, cfg)
			case http.MethodPatch:
				handleUpdateUser(w, r, manager, cfg)
			default:
				shared.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
			}
			return
		}

		parts := strings.Split(path, "/")
		userUUID := parts[0]
		if _, err := uuid.Parse(userUUID); err != nil {
			shared.SendError(w, http.StatusBadRequest, "invalid UUID format", nil, cfg)
			return
		}

		if len(parts) >= 3 && parts[1] == "actions" && r.Method == http.MethodPost {
			switch parts[2] {
			case "enable":
				handleEnableUser(w, r, manager, cfg, userUUID)
			case "disable":
				handleDisableUser(w, r, manager, cfg, userUUID)
			case "reset-traffic":
				handleResetUserTraffic(w, r, manager, cfg, userUUID)
			case "revoke":
				handleRevokeUserSubscription(w, r, manager, cfg, userUUID)
			default:
				http.NotFound(w, r)
			}
			return
		}

		if len(parts) != 1 {
			http.NotFound(w, r)
			return
		}

		switch r.Method {
		case http.MethodGet:
			handleGetUser(w, r, manager, cfg, userUUID)
		case http.MethodDelete:
			handleDeleteUser(w, r, manager, cfg, userUUID)
		default:
			shared.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
		}
	}
}

func UsersBulkHandler(manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			shared.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}
		path := trimUsersPath(r.URL.Path, "/bulk/")
		switch path {
		case "delete":
			handleBulkDeleteUsers(w, r, manager, cfg)
		default:
			http.NotFound(w, r)
		}
	}
}

func UsersTagsHandler(manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			shared.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}

		tags, err := getAllUserTags(r.Context(), manager)
		if err != nil {
			shared.SendError(w, http.StatusInternalServerError, "failed to fetch user tags", err, cfg)
			return
		}

		shared.WriteJSON(w, http.StatusOK, map[string]any{
			"response": map[string]any{
				"tags": tags,
			},
		})
	}
}

func trimUsersPath(path string, suffix string) string {
	for _, prefix := range []string{"/api/users"} {
		if strings.HasPrefix(path, prefix+suffix) {
			return strings.Trim(strings.TrimPrefix(path, prefix+suffix), "/")
		}
	}
	return strings.Trim(path, "/")
}

func handleGetUsers(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) {
	records, err := getAllUserRecords(r.Context(), manager)
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to fetch users", err, cfg)
		return
	}

	response, err := buildUserResponses(r.Context(), manager, records, resolveUsersSubscriptionBase(r, cfg))
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to build users response", err, cfg)
		return
	}

	shared.WriteJSON(w, http.StatusOK, map[string]any{
		"response": map[string]any{
			"users": response,
			"total": len(response),
		},
	})
}

func handleGetUser(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig, userUUID string) {
	record, err := getUserRecordByUUID(r.Context(), manager, userUUID)
	if err != nil {
		if errors.Is(err, errUserNotFound) {
			shared.SendError(w, http.StatusNotFound, "user not found", nil, cfg)
			return
		}
		shared.SendError(w, http.StatusInternalServerError, "failed to fetch user", err, cfg)
		return
	}

	response, err := buildUserResponses(r.Context(), manager, []userRecord{record}, resolveUsersSubscriptionBase(r, cfg))
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to build user response", err, cfg)
		return
	}

	shared.WriteJSON(w, http.StatusOK, map[string]any{"response": response[0]})
}

func handleCreateUser(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) {
	var req createUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		shared.SendError(w, http.StatusBadRequest, "invalid JSON", err, cfg)
		return
	}

	if err := validateCreateUserRequest(req); err != nil {
		shared.SendError(w, http.StatusBadRequest, err.Error(), nil, cfg)
		return
	}

	userUUID := coalesceUUID(req.UUID)
	shortUUID := coalesceShortUUID(req.ShortUUID)
	trojanPassword := coalesceRandomString(req.TrojanPassword, 16)
	vlessUUID := coalesceUUID(req.VlessUUID)
	ssPassword := coalesceRandomString(req.SSPassword, 16)

	if trojanPassword == "" || ssPassword == "" {
		shared.SendError(w, http.StatusInternalServerError, "failed to generate user credentials", nil, cfg)
		return
	}

	expireAt, _ := time.Parse(time.RFC3339, req.ExpireAt)
	createdAt := time.Now().UTC()
	if req.CreatedAt != nil && strings.TrimSpace(*req.CreatedAt) != "" {
		createdAt, _ = time.Parse(time.RFC3339, strings.TrimSpace(*req.CreatedAt))
	}
	var lastTrafficResetAt any
	if req.LastTrafficResetAt != nil && strings.TrimSpace(*req.LastTrafficResetAt) != "" {
		if parsed, err := time.Parse(time.RFC3339, strings.TrimSpace(*req.LastTrafficResetAt)); err == nil {
			lastTrafficResetAt = parsed
		}
	}

	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		tx, err := db.BeginTx(r.Context(), nil)
		if err != nil {
			return err
		}

		var tID int64
		insertErr := tx.QueryRowContext(r.Context(), `
			INSERT INTO users (
				uuid, short_uuid, username, status, traffic_limit_bytes, traffic_limit_strategy,
				expire_at, sub_last_user_agent, sub_last_opened_at, last_traffic_reset_at, sub_revoked_at,
				trojan_password, vless_uuid, ss_password, description, tag, telegram_id, email,
				hwid_device_limit, external_squad_uuid, last_triggered_threshold, created_at, updated_at
			) VALUES (?, ?, ?, ?, ?, ?, ?, NULL, NULL, ?, NULL, ?, ?, ?, ?, ?, ?, ?, ?, ?, 0, ?, ?)
			RETURNING t_id
		`,
			userUUID,
			shortUUID,
			strings.TrimSpace(req.Username),
			normalizeUserStatus(req.Status),
			coalesceInt64(req.TrafficLimitBytes, 0),
			normalizeTrafficStrategy(req.TrafficLimitStrategy),
			expireAt.UTC(),
			lastTrafficResetAt,
			trojanPassword,
			vlessUUID,
			ssPassword,
			normalizeNullableString(req.Description),
			normalizeUserTag(req.Tag),
			req.TelegramID,
			normalizeNullableString(req.Email),
			req.HwidDeviceLimit,
			normalizeNullableString(req.ExternalSquadUUID),
			createdAt.UTC(),
			createdAt.UTC(),
		).Scan(&tID)
		if insertErr != nil {
			_ = tx.Rollback()
			return mapUserWriteError(insertErr)
		}

		if _, err := tx.ExecContext(r.Context(), `
			INSERT INTO user_traffic (
				t_id, used_traffic_bytes, lifetime_used_traffic_bytes, online_at,
				last_connected_node_uuid, first_connected_at
			) VALUES (?, 0, 0, NULL, NULL, NULL)
		`, tID); err != nil {
			_ = tx.Rollback()
			return err
		}

		if err := replaceUserInternalSquadsTx(r.Context(), tx, tID, req.ActiveInternalSquads); err != nil {
			_ = tx.Rollback()
			return err
		}

		return tx.Commit()
	})
	if err != nil {
		handleUserWriteError(w, err, cfg)
		return
	}

	record, err := getUserRecordByUUID(r.Context(), manager, userUUID)
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to fetch created user", err, cfg)
		return
	}
	response, err := buildUserResponses(r.Context(), manager, []userRecord{record}, resolveUsersSubscriptionBase(r, cfg))
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to build created user response", err, cfg)
		return
	}

	monitor.RequestNodeDeploy(true)
	shared.WriteJSON(w, http.StatusCreated, map[string]any{"response": response[0]})
}

func handleUpdateUser(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) {
	var req updateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		shared.SendError(w, http.StatusBadRequest, "invalid JSON", err, cfg)
		return
	}
	if err := validateUpdateUserRequest(req); err != nil {
		shared.SendError(w, http.StatusBadRequest, err.Error(), nil, cfg)
		return
	}

	targetUUID, err := resolveUserUUIDForUpdate(r.Context(), manager, req.UUID, req.Username)
	if err != nil {
		if errors.Is(err, errUserNotFound) {
			shared.SendError(w, http.StatusNotFound, "user not found", nil, cfg)
			return
		}
		shared.SendError(w, http.StatusBadRequest, err.Error(), nil, cfg)
		return
	}

	record, err := getUserRecordByUUID(r.Context(), manager, targetUUID)
	if err != nil {
		if errors.Is(err, errUserNotFound) {
			shared.SendError(w, http.StatusNotFound, "user not found", nil, cfg)
			return
		}
		shared.SendError(w, http.StatusInternalServerError, "failed to fetch user", err, cfg)
		return
	}

	internalSquadsChanged := false
	internalSquadNodeUUIDs := make([]string, 0)
	err = manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		tx, err := db.BeginTx(r.Context(), nil)
		if err != nil {
			return err
		}

		clauses := make([]string, 0)
		args := make([]any, 0)
		add := func(column string, value any) {
			clauses = append(clauses, fmt.Sprintf("%s = ?", column))
			args = append(args, value)
		}

		if req.Status != nil {
			add("status", strings.ToUpper(strings.TrimSpace(*req.Status)))
		}
		if req.TrafficLimitBytes != nil {
			add("traffic_limit_bytes", *req.TrafficLimitBytes)
		}
		if req.TrafficLimitStrategy != nil {
			add("traffic_limit_strategy", strings.ToUpper(strings.TrimSpace(*req.TrafficLimitStrategy)))
		}
		if req.ExpireAt != nil {
			parsed, _ := time.Parse(time.RFC3339, strings.TrimSpace(*req.ExpireAt))
			add("expire_at", parsed.UTC())
		}
		if req.Description.Set {
			if req.Description.Value == nil || strings.TrimSpace(*req.Description.Value) == "" {
				clauses = append(clauses, "description = NULL")
			} else {
				add("description", strings.TrimSpace(*req.Description.Value))
			}
		}
		if req.Tag.Set {
			if req.Tag.Value == nil || strings.TrimSpace(*req.Tag.Value) == "" {
				clauses = append(clauses, "tag = NULL")
			} else {
				add("tag", strings.ToUpper(strings.TrimSpace(*req.Tag.Value)))
			}
		}
		if req.TelegramID.Set {
			if req.TelegramID.Value == nil {
				clauses = append(clauses, "telegram_id = NULL")
			} else {
				add("telegram_id", *req.TelegramID.Value)
			}
		}
		if req.Email.Set {
			if req.Email.Value == nil || strings.TrimSpace(*req.Email.Value) == "" {
				clauses = append(clauses, "email = NULL")
			} else {
				add("email", strings.TrimSpace(*req.Email.Value))
			}
		}
		if req.HwidDeviceLimit.Set {
			if req.HwidDeviceLimit.Value == nil {
				clauses = append(clauses, "hwid_device_limit = NULL")
			} else {
				add("hwid_device_limit", *req.HwidDeviceLimit.Value)
			}
		}
		if req.ExternalSquadUUID.Set {
			if req.ExternalSquadUUID.Value == nil || strings.TrimSpace(*req.ExternalSquadUUID.Value) == "" {
				clauses = append(clauses, "external_squad_uuid = NULL")
			} else {
				add("external_squad_uuid", strings.TrimSpace(*req.ExternalSquadUUID.Value))
			}
		}

		if len(clauses) > 0 {
			args = append(args, targetUUID)
			query := fmt.Sprintf("UPDATE users SET %s, updated_at = CURRENT_TIMESTAMP WHERE uuid = ?", strings.Join(clauses, ", "))
			if _, err := tx.ExecContext(r.Context(), query, args...); err != nil {
				_ = tx.Rollback()
				return mapUserWriteError(err)
			}
		}

		if req.ActiveInternalSquads != nil {
			currentSquads, loadErr := getUserInternalSquadsTx(r.Context(), tx, record.TID)
			if loadErr != nil {
				_ = tx.Rollback()
				return loadErr
			}
			requestedSquads := dedupeStrings(*req.ActiveInternalSquads)
			if internalSquadSetsDiffer(currentSquads, requestedSquads) {
				affectedSquads := dedupeStrings(append(append([]string{}, currentSquads...), requestedSquads...))
				nodeUUIDs, nodeTargetsErr := resolveNodeUUIDsForInternalSquadsTx(r.Context(), tx, affectedSquads)
				if nodeTargetsErr != nil {
					_ = tx.Rollback()
					return nodeTargetsErr
				}
				if err := replaceUserInternalSquadsTx(r.Context(), tx, record.TID, requestedSquads); err != nil {
					_ = tx.Rollback()
					return err
				}
				internalSquadNodeUUIDs = nodeUUIDs
				internalSquadsChanged = true
			}
		}

		return tx.Commit()
	})
	if err != nil {
		handleUserWriteError(w, err, cfg)
		return
	}

	updatedRecord, err := getUserRecordByUUID(r.Context(), manager, targetUUID)
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to fetch updated user", err, cfg)
		return
	}
	response, err := buildUserResponses(r.Context(), manager, []userRecord{updatedRecord}, resolveUsersSubscriptionBase(r, cfg))
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to build updated user response", err, cfg)
		return
	}

	if internalSquadsChanged && len(internalSquadNodeUUIDs) > 0 {
		monitor.RequestNodeDeploy(true, internalSquadNodeUUIDs...)
	}
	shared.WriteJSON(w, http.StatusOK, map[string]any{"response": response[0]})
}

func handleDeleteUser(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig, userUUID string) {
	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		result, err := db.ExecContext(r.Context(), `DELETE FROM users WHERE uuid = ?`, userUUID)
		if err != nil {
			return err
		}
		rows, err := result.RowsAffected()
		if err != nil {
			return err
		}
		if rows == 0 {
			return errUserNotFound
		}
		return nil
	})
	if err != nil {
		if errors.Is(err, errUserNotFound) {
			shared.SendError(w, http.StatusNotFound, "user not found", nil, cfg)
			return
		}
		shared.SendError(w, http.StatusInternalServerError, "failed to delete user", err, cfg)
		return
	}

	monitor.RequestNodeDeploy(true)
	shared.WriteJSON(w, http.StatusOK, map[string]any{"response": map[string]any{"isDeleted": true}})
}

func getUserInternalSquadsTx(ctx context.Context, tx dbmanager.TxExecutor, tID int64) ([]string, error) {
	rows, err := tx.QueryContext(ctx, `SELECT internal_squad_uuid FROM internal_squad_members WHERE user_id = ?`, tID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	current := make([]string, 0)
	for rows.Next() {
		var squadUUID string
		if err := rows.Scan(&squadUUID); err != nil {
			return nil, err
		}
		current = append(current, strings.TrimSpace(squadUUID))
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return dedupeStrings(current), nil
}

func internalSquadSetsDiffer(current []string, requested []string) bool {
	nextSet := make(map[string]struct{})
	for _, squadUUID := range dedupeStrings(requested) {
		nextSet[squadUUID] = struct{}{}
	}
	if len(current) != len(nextSet) {
		return true
	}
	for _, squadUUID := range current {
		if _, ok := nextSet[squadUUID]; !ok {
			return true
		}
	}
	return false
}

func resolveNodeUUIDsForInternalSquadsTx(ctx context.Context, tx dbmanager.TxExecutor, squadUUIDs []string) ([]string, error) {
	cleanSquadUUIDs := dedupeStrings(squadUUIDs)
	if len(cleanSquadUUIDs) == 0 {
		return []string{}, nil
	}

	rows, err := tx.QueryContext(ctx, `
		SELECT DISTINCT cpitn.node_uuid
		FROM internal_squad_inbounds isi
		JOIN config_profile_inbounds_to_nodes cpitn ON cpitn.config_profile_inbound_uuid = isi.inbound_uuid
		WHERE isi.internal_squad_uuid = ANY(?)
	`, cleanSquadUUIDs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	nodeUUIDs := make([]string, 0)
	for rows.Next() {
		var nodeUUID string
		if err := rows.Scan(&nodeUUID); err != nil {
			return nil, err
		}
		nodeUUIDs = append(nodeUUIDs, strings.TrimSpace(nodeUUID))
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return dedupeStrings(nodeUUIDs), nil
}

func handleBulkDeleteUsers(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) {
	var req bulkDeleteUsersRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		shared.SendError(w, http.StatusBadRequest, "invalid JSON", err, cfg)
		return
	}
	if err := validateUUIDList(req.UUIDs); err != nil {
		shared.SendError(w, http.StatusBadRequest, err.Error(), nil, cfg)
		return
	}

	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		_, err := db.ExecContext(r.Context(), `DELETE FROM users WHERE uuid = ANY(?)`, req.UUIDs)
		return err
	})
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to delete users", err, cfg)
		return
	}

	monitor.RequestNodeDeploy(true)
	shared.WriteJSON(w, http.StatusOK, map[string]any{"response": map[string]any{"isDeleted": true}})
}

func handleEnableUser(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig, userUUID string) {
	if err := updateUserStatus(r.Context(), manager, userUUID, "ACTIVE"); err != nil {
		handleUserActionError(w, err, cfg, "failed to enable user")
		return
	}
	monitor.RequestNodeDeploy(true)
	sendUpdatedUserResponse(w, r, manager, cfg, userUUID)
}

func handleDisableUser(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig, userUUID string) {
	if err := updateUserStatus(r.Context(), manager, userUUID, "DISABLED"); err != nil {
		handleUserActionError(w, err, cfg, "failed to disable user")
		return
	}
	monitor.RequestNodeDeploy(true)
	sendUpdatedUserResponse(w, r, manager, cfg, userUUID)
}

func handleResetUserTraffic(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig, userUUID string) {
	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		_, err := db.ExecContext(r.Context(), `
			UPDATE users
			SET last_traffic_reset_at = CURRENT_TIMESTAMP, last_triggered_threshold = 0, updated_at = CURRENT_TIMESTAMP
			WHERE uuid = ?
		`, userUUID)
		if err != nil {
			return err
		}

		result, err := db.ExecContext(r.Context(), `
			UPDATE user_traffic
			SET used_traffic_bytes = 0
			WHERE t_id = (SELECT t_id FROM users WHERE uuid = ?)
		`, userUUID)
		if err != nil {
			return err
		}
		rows, err := result.RowsAffected()
		if err != nil {
			return err
		}
		if rows == 0 {
			return errUserNotFound
		}
		return nil
	})
	if err != nil {
		handleUserActionError(w, err, cfg, "failed to reset user traffic")
		return
	}
	sendUpdatedUserResponse(w, r, manager, cfg, userUUID)
}

func handleRevokeUserSubscription(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig, userUUID string) {
	shortUUID := generateSubscriptionShortUUID()
	if shortUUID == "" {
		shared.SendError(w, http.StatusInternalServerError, "failed to generate short uuid", nil, cfg)
		return
	}

	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		result, err := db.ExecContext(r.Context(), `
			UPDATE users
			SET short_uuid = ?, sub_revoked_at = CURRENT_TIMESTAMP, updated_at = CURRENT_TIMESTAMP
			WHERE uuid = ?
		`, shortUUID, userUUID)
		if err != nil {
			return mapUserWriteError(err)
		}
		rows, err := result.RowsAffected()
		if err != nil {
			return err
		}
		if rows == 0 {
			return errUserNotFound
		}
		return nil
	})
	if err != nil {
		handleUserWriteError(w, err, cfg)
		return
	}
	monitor.RequestNodeDeploy(true)
	sendUpdatedUserResponse(w, r, manager, cfg, userUUID)
}

func sendUpdatedUserResponse(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig, userUUID string) {
	record, err := getUserRecordByUUID(r.Context(), manager, userUUID)
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to fetch updated user", err, cfg)
		return
	}
	response, err := buildUserResponses(r.Context(), manager, []userRecord{record}, resolveUsersSubscriptionBase(r, cfg))
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to build updated user response", err, cfg)
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

func updateUserStatus(ctx context.Context, manager *dbmanager.DatabaseManager, userUUID string, status string) error {
	return manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		result, err := db.ExecContext(ctx, `
			UPDATE users
			SET status = ?, updated_at = CURRENT_TIMESTAMP
			WHERE uuid = ?
		`, status, userUUID)
		if err != nil {
			return err
		}
		rows, err := result.RowsAffected()
		if err != nil {
			return err
		}
		if rows == 0 {
			return errUserNotFound
		}
		return nil
	})
}

func validateCreateUserRequest(req createUserRequest) error {
	username := strings.TrimSpace(req.Username)
	if len(username) < 3 || len(username) > 36 {
		return fmt.Errorf("username must be between 3 and 36 characters")
	}
	if !userUsernameRegex.MatchString(username) {
		return fmt.Errorf("username can only contain letters, numbers, underscores and dashes")
	}
	if req.Status != nil && !isValidUserStatus(*req.Status) {
		return fmt.Errorf("invalid status")
	}
	if req.ExpireAt == "" {
		return fmt.Errorf("expireAt is required")
	}
	if _, err := time.Parse(time.RFC3339, req.ExpireAt); err != nil {
		return fmt.Errorf("expireAt must be RFC3339")
	}
	if req.CreatedAt != nil && strings.TrimSpace(*req.CreatedAt) != "" {
		if _, err := time.Parse(time.RFC3339, strings.TrimSpace(*req.CreatedAt)); err != nil {
			return fmt.Errorf("createdAt must be RFC3339")
		}
	}
	if req.LastTrafficResetAt != nil && strings.TrimSpace(*req.LastTrafficResetAt) != "" {
		if _, err := time.Parse(time.RFC3339, strings.TrimSpace(*req.LastTrafficResetAt)); err != nil {
			return fmt.Errorf("lastTrafficResetAt must be RFC3339")
		}
	}
	if req.TrafficLimitBytes != nil && *req.TrafficLimitBytes < 0 {
		return fmt.Errorf("trafficLimitBytes must be non-negative")
	}
	if req.TrafficLimitStrategy != nil && !isValidTrafficStrategy(*req.TrafficLimitStrategy) {
		return fmt.Errorf("invalid trafficLimitStrategy")
	}
	if req.Tag != nil && strings.TrimSpace(*req.Tag) != "" {
		tag := strings.ToUpper(strings.TrimSpace(*req.Tag))
		if len(tag) > 16 || !userTagRegex.MatchString(tag) {
			return fmt.Errorf("tag can only contain uppercase letters, numbers, underscores and be up to 16 chars")
		}
	}
	if req.Email != nil && strings.TrimSpace(*req.Email) != "" && !strings.Contains(strings.TrimSpace(*req.Email), "@") {
		return fmt.Errorf("invalid email")
	}
	if req.HwidDeviceLimit != nil && *req.HwidDeviceLimit < 0 {
		return fmt.Errorf("hwidDeviceLimit must be non-negative")
	}
	if req.TelegramID != nil && *req.TelegramID < 0 {
		return fmt.Errorf("telegramId must be non-negative")
	}
	if err := validateUUIDListAllowEmpty(req.ActiveInternalSquads); err != nil {
		return err
	}
	if req.UUID != nil && strings.TrimSpace(*req.UUID) != "" {
		if _, err := uuid.Parse(strings.TrimSpace(*req.UUID)); err != nil {
			return fmt.Errorf("invalid uuid")
		}
	}
	if req.VlessUUID != nil && strings.TrimSpace(*req.VlessUUID) != "" {
		if _, err := uuid.Parse(strings.TrimSpace(*req.VlessUUID)); err != nil {
			return fmt.Errorf("invalid vlessUuid")
		}
	}
	if req.ExternalSquadUUID != nil && strings.TrimSpace(*req.ExternalSquadUUID) != "" {
		if _, err := uuid.Parse(strings.TrimSpace(*req.ExternalSquadUUID)); err != nil {
			return fmt.Errorf("invalid externalSquadUuid")
		}
	}
	return nil
}

func validateUpdateUserRequest(req updateUserRequest) error {
	if req.UUID == nil && req.Username == nil {
		return fmt.Errorf("either uuid or username must be provided")
	}
	if req.UUID != nil && strings.TrimSpace(*req.UUID) != "" {
		if _, err := uuid.Parse(strings.TrimSpace(*req.UUID)); err != nil {
			return fmt.Errorf("invalid uuid")
		}
	}
	if req.Status != nil && !isValidUserStatus(*req.Status) {
		return fmt.Errorf("invalid status")
	}
	if req.TrafficLimitBytes != nil && *req.TrafficLimitBytes < 0 {
		return fmt.Errorf("trafficLimitBytes must be non-negative")
	}
	if req.TrafficLimitStrategy != nil && !isValidTrafficStrategy(*req.TrafficLimitStrategy) {
		return fmt.Errorf("invalid trafficLimitStrategy")
	}
	if req.ExpireAt != nil {
		if _, err := time.Parse(time.RFC3339, strings.TrimSpace(*req.ExpireAt)); err != nil {
			return fmt.Errorf("expireAt must be RFC3339")
		}
	}
	if req.Tag.Set && req.Tag.Value != nil && strings.TrimSpace(*req.Tag.Value) != "" {
		tag := strings.ToUpper(strings.TrimSpace(*req.Tag.Value))
		if len(tag) > 16 || !userTagRegex.MatchString(tag) {
			return fmt.Errorf("tag can only contain uppercase letters, numbers, underscores and be up to 16 chars")
		}
	}
	if req.Email.Set && req.Email.Value != nil && strings.TrimSpace(*req.Email.Value) != "" && !strings.Contains(strings.TrimSpace(*req.Email.Value), "@") {
		return fmt.Errorf("invalid email")
	}
	if req.HwidDeviceLimit.Set && req.HwidDeviceLimit.Value != nil && *req.HwidDeviceLimit.Value < 0 {
		return fmt.Errorf("hwidDeviceLimit must be non-negative")
	}
	if req.TelegramID.Set && req.TelegramID.Value != nil && *req.TelegramID.Value < 0 {
		return fmt.Errorf("telegramId must be non-negative")
	}
	if req.ActiveInternalSquads != nil {
		if err := validateUUIDListAllowEmpty(*req.ActiveInternalSquads); err != nil {
			return err
		}
	}
	if req.ExternalSquadUUID.Set && req.ExternalSquadUUID.Value != nil && strings.TrimSpace(*req.ExternalSquadUUID.Value) != "" {
		if _, err := uuid.Parse(strings.TrimSpace(*req.ExternalSquadUUID.Value)); err != nil {
			return fmt.Errorf("invalid externalSquadUuid")
		}
	}
	return nil
}

func getAllUserRecords(ctx context.Context, manager *dbmanager.DatabaseManager) ([]userRecord, error) {
	records := make([]userRecord, 0)
	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		rows, err := db.QueryContext(ctx, `
			SELECT
				u.t_id, u.uuid, u.short_uuid, u.username, u.status, u.traffic_limit_bytes,
				u.traffic_limit_strategy, u.expire_at, u.sub_last_user_agent, u.sub_last_opened_at,
				u.last_traffic_reset_at, u.sub_revoked_at, u.trojan_password, u.vless_uuid, u.ss_password,
				u.description, u.tag, u.telegram_id, u.email, u.hwid_device_limit, u.external_squad_uuid,
				u.last_triggered_threshold, u.created_at, u.updated_at,
				COALESCE(ut.used_traffic_bytes, 0), COALESCE(ut.lifetime_used_traffic_bytes, 0),
				ut.online_at, ut.last_connected_node_uuid, ut.first_connected_at
			FROM users u
			LEFT JOIN user_traffic ut ON ut.t_id = u.t_id
			ORDER BY u.t_id DESC
		`)
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			record, scanErr := scanUserRecord(rows)
			if scanErr != nil {
				return scanErr
			}
			records = append(records, record)
		}
		return rows.Err()
	})
	return records, err
}

func getUserRecordByUUID(ctx context.Context, manager *dbmanager.DatabaseManager, userUUID string) (userRecord, error) {
	var record userRecord
	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		row := db.QueryRowContext(ctx, `
			SELECT
				u.t_id, u.uuid, u.short_uuid, u.username, u.status, u.traffic_limit_bytes,
				u.traffic_limit_strategy, u.expire_at, u.sub_last_user_agent, u.sub_last_opened_at,
				u.last_traffic_reset_at, u.sub_revoked_at, u.trojan_password, u.vless_uuid, u.ss_password,
				u.description, u.tag, u.telegram_id, u.email, u.hwid_device_limit, u.external_squad_uuid,
				u.last_triggered_threshold, u.created_at, u.updated_at,
				COALESCE(ut.used_traffic_bytes, 0), COALESCE(ut.lifetime_used_traffic_bytes, 0),
				ut.online_at, ut.last_connected_node_uuid, ut.first_connected_at
			FROM users u
			LEFT JOIN user_traffic ut ON ut.t_id = u.t_id
			WHERE u.uuid = ?
		`, userUUID)
		var scanErr error
		record, scanErr = scanUserRecord(row)
		if errors.Is(scanErr, sql.ErrNoRows) {
			return errUserNotFound
		}
		return scanErr
	})
	return record, err
}

func scanUserRecord(scanner shared.RowScanner) (userRecord, error) {
	var record userRecord
	var (
		subLastUserAgent sql.NullString
		subLastOpenedAt  sql.NullTime
		lastTrafficReset sql.NullTime
		subRevokedAt     sql.NullTime
		description      sql.NullString
		tag              sql.NullString
		telegramID       sql.NullInt64
		email            sql.NullString
		hwidDeviceLimit  sql.NullInt64
		externalSquad    sql.NullString
		onlineAt         sql.NullTime
		lastNodeUUID     sql.NullString
		firstConnectedAt sql.NullTime
	)

	err := scanner.Scan(
		&record.TID,
		&record.UUID,
		&record.ShortUUID,
		&record.Username,
		&record.Status,
		&record.TrafficLimitBytes,
		&record.TrafficLimitStrategy,
		&record.ExpireAt,
		&subLastUserAgent,
		&subLastOpenedAt,
		&lastTrafficReset,
		&subRevokedAt,
		&record.TrojanPassword,
		&record.VlessUUID,
		&record.SSPassword,
		&description,
		&tag,
		&telegramID,
		&email,
		&hwidDeviceLimit,
		&externalSquad,
		&record.LastTriggeredThreshold,
		&record.CreatedAt,
		&record.UpdatedAt,
		&record.UsedTrafficBytes,
		&record.LifetimeUsedTrafficBytes,
		&onlineAt,
		&lastNodeUUID,
		&firstConnectedAt,
	)
	if err != nil {
		return record, err
	}

	if subLastUserAgent.Valid {
		record.SubLastUserAgent = &subLastUserAgent.String
	}
	if subLastOpenedAt.Valid {
		record.SubLastOpenedAt = &subLastOpenedAt.Time
	}
	if lastTrafficReset.Valid {
		record.LastTrafficResetAt = &lastTrafficReset.Time
	}
	if subRevokedAt.Valid {
		record.SubRevokedAt = &subRevokedAt.Time
	}
	if description.Valid {
		record.Description = &description.String
	}
	if tag.Valid {
		record.Tag = &tag.String
	}
	if telegramID.Valid {
		value := telegramID.Int64
		record.TelegramID = &value
	}
	if email.Valid {
		record.Email = &email.String
	}
	if hwidDeviceLimit.Valid {
		value := int(hwidDeviceLimit.Int64)
		record.HwidDeviceLimit = &value
	}
	if externalSquad.Valid {
		record.ExternalSquadUUID = &externalSquad.String
	}
	if onlineAt.Valid {
		record.OnlineAt = &onlineAt.Time
	}
	if lastNodeUUID.Valid {
		record.LastConnectedNodeUUID = &lastNodeUUID.String
	}
	if firstConnectedAt.Valid {
		record.FirstConnectedAt = &firstConnectedAt.Time
	}

	return record, nil
}

func buildUserResponses(ctx context.Context, manager *dbmanager.DatabaseManager, records []userRecord, subscriptionBase string) ([]userAPI, error) {
	userUUIDs := make([]string, 0, len(records))
	for _, record := range records {
		userUUIDs = append(userUUIDs, record.UUID)
	}

	activeSquadsMap, err := getUsersActiveInternalSquads(ctx, manager, userUUIDs)
	if err != nil {
		return nil, err
	}

	response := make([]userAPI, 0, len(records))
	for _, record := range records {
		activeSquads := activeSquadsMap[record.UUID]
		if activeSquads == nil {
			activeSquads = []internalSquadResponse{}
		}
		response = append(response, userAPI{
			UUID:                   record.UUID,
			ID:                     record.TID,
			ShortUUID:              record.ShortUUID,
			Username:               record.Username,
			Status:                 record.Status,
			TrafficLimitBytes:      record.TrafficLimitBytes,
			TrafficLimitStrategy:   record.TrafficLimitStrategy,
			ExpireAt:               record.ExpireAt,
			TelegramID:             record.TelegramID,
			Email:                  record.Email,
			Description:            record.Description,
			Tag:                    record.Tag,
			HwidDeviceLimit:        record.HwidDeviceLimit,
			ExternalSquadUUID:      record.ExternalSquadUUID,
			TrojanPassword:         record.TrojanPassword,
			VlessUUID:              record.VlessUUID,
			SSPassword:             record.SSPassword,
			LastTriggeredThreshold: record.LastTriggeredThreshold,
			SubRevokedAt:           record.SubRevokedAt,
			SubLastUserAgent:       record.SubLastUserAgent,
			SubLastOpenedAt:        record.SubLastOpenedAt,
			LastTrafficResetAt:     record.LastTrafficResetAt,
			CreatedAt:              record.CreatedAt,
			UpdatedAt:              record.UpdatedAt,
			SubscriptionURL:        subscriptionBase + record.ShortUUID,
			ActiveInternalSquads:   activeSquads,
			UserTraffic: userTrafficResponse{
				UsedTrafficBytes:         record.UsedTrafficBytes,
				LifetimeUsedTrafficBytes: record.LifetimeUsedTrafficBytes,
				OnlineAt:                 record.OnlineAt,
				FirstConnectedAt:         record.FirstConnectedAt,
				LastConnectedNodeUUID:    record.LastConnectedNodeUUID,
			},
		})
	}

	return response, nil
}

func getUsersActiveInternalSquads(ctx context.Context, manager *dbmanager.DatabaseManager, userUUIDs []string) (map[string][]internalSquadResponse, error) {
	result := make(map[string][]internalSquadResponse, len(userUUIDs))
	if len(userUUIDs) == 0 {
		return result, nil
	}
	for _, userUUID := range userUUIDs {
		result[userUUID] = []internalSquadResponse{}
	}

	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		rows, err := db.QueryContext(ctx, `
			SELECT u.uuid, s.uuid, s.name
			FROM users u
			INNER JOIN internal_squad_members ism ON ism.user_id = u.t_id
			INNER JOIN internal_squads s ON s.uuid = ism.internal_squad_uuid
			WHERE u.uuid = ANY(?)
			ORDER BY s.view_position ASC, s.name ASC
		`, userUUIDs)
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			var userUUID, squadUUID, squadName string
			if err := rows.Scan(&userUUID, &squadUUID, &squadName); err != nil {
				return err
			}
			result[userUUID] = append(result[userUUID], internalSquadResponse{
				UUID: squadUUID,
				Name: squadName,
			})
		}
		return rows.Err()
	})

	return result, err
}

func getAllUserTags(ctx context.Context, manager *dbmanager.DatabaseManager) ([]string, error) {
	tags := make([]string, 0)
	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		rows, err := db.QueryContext(ctx, `
			SELECT DISTINCT tag
			FROM users
			WHERE tag IS NOT NULL AND tag <> ''
			ORDER BY tag ASC
		`)
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			var tag string
			if err := rows.Scan(&tag); err != nil {
				return err
			}
			tags = append(tags, tag)
		}
		return rows.Err()
	})
	return tags, err
}

func replaceUserInternalSquadsTx(ctx context.Context, tx dbmanager.TxExecutor, tID int64, squadUUIDs []string) error {
	if _, err := tx.ExecContext(ctx, `DELETE FROM internal_squad_members WHERE user_id = ?`, tID); err != nil {
		return err
	}
	for _, squadUUID := range dedupeStrings(squadUUIDs) {
		if _, err := tx.ExecContext(ctx, `
			INSERT INTO internal_squad_members (internal_squad_uuid, user_id)
			VALUES (?, ?)
			ON CONFLICT (internal_squad_uuid, user_id) DO NOTHING
		`, squadUUID, tID); err != nil {
			return err
		}
	}
	return nil
}

func resolveUserUUIDForUpdate(ctx context.Context, manager *dbmanager.DatabaseManager, userUUID *string, username *string) (string, error) {
	if userUUID != nil && strings.TrimSpace(*userUUID) != "" {
		return strings.TrimSpace(*userUUID), nil
	}
	if username == nil || strings.TrimSpace(*username) == "" {
		return "", fmt.Errorf("either uuid or username must be provided")
	}

	var resolved string
	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		err := db.QueryRowContext(ctx, `SELECT uuid FROM users WHERE username = ?`, strings.TrimSpace(*username)).Scan(&resolved)
		if errors.Is(err, sql.ErrNoRows) {
			return errUserNotFound
		}
		return err
	})
	return resolved, err
}

func resolveUsersSubscriptionBase(r *http.Request, cfg *config.BackendConfig) string {
	if domain := strings.TrimSpace(os.Getenv("FRONT_END_DOMAIN")); domain != "" {
		first := strings.TrimSpace(strings.Split(domain, ",")[0])
		if parsed, err := url.Parse(first); err == nil && parsed.Host != "" {
			return strings.TrimRight(parsed.String(), "/") + "/"
		}
		if strings.HasPrefix(first, "http://") || strings.HasPrefix(first, "https://") {
			return strings.TrimRight(first, "/") + "/"
		}
		return "https://" + strings.TrimRight(first, "/") + "/"
	}

	scheme := "https"
	if cfg != nil && cfg.Panel.AllowInsecureHTTP {
		scheme = "http"
	}
	return scheme + "://" + r.Host + "/"
}

func handleUserWriteError(w http.ResponseWriter, err error, cfg *config.BackendConfig) {
	switch {
	case errors.Is(err, errUsernameExists):
		shared.SendError(w, http.StatusConflict, "username already exists", nil, cfg)
	case errors.Is(err, errShortUUIDExists):
		shared.SendError(w, http.StatusConflict, "short uuid already exists", nil, cfg)
	case errors.Is(err, errVLESSUUIDExists):
		shared.SendError(w, http.StatusConflict, "vless uuid already exists", nil, cfg)
	case errors.Is(err, errUserNotFound):
		shared.SendError(w, http.StatusNotFound, "user not found", nil, cfg)
	default:
		shared.SendError(w, http.StatusInternalServerError, "failed to write user", err, cfg)
	}
}

func mapUserWriteError(err error) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		switch pgErr.ConstraintName {
		case "users_username_key":
			return errUsernameExists
		case "users_short_uuid_key":
			return errShortUUIDExists
		case "users_vless_uuid_key":
			return errVLESSUUIDExists
		}
	}

	message := err.Error()
	switch {
	case strings.Contains(message, "users_username_key"), strings.Contains(message, "UNIQUE constraint failed: users.username"):
		return errUsernameExists
	case strings.Contains(message, "users_short_uuid_key"), strings.Contains(message, "UNIQUE constraint failed: users.short_uuid"):
		return errShortUUIDExists
	case strings.Contains(message, "users_vless_uuid_key"), strings.Contains(message, "UNIQUE constraint failed: users.vless_uuid"):
		return errVLESSUUIDExists
	default:
		return err
	}
}

func validateUUIDList(values []string) error {
	if len(values) == 0 {
		return fmt.Errorf("uuids cannot be empty")
	}
	return validateUUIDListAllowEmpty(values)
}

func validateUUIDListAllowEmpty(values []string) error {
	for _, value := range values {
		if _, err := uuid.Parse(strings.TrimSpace(value)); err != nil {
			return fmt.Errorf("invalid uuid value")
		}
	}
	return nil
}

func dedupeStrings(values []string) []string {
	result := make([]string, 0, len(values))
	seen := make(map[string]struct{}, len(values))
	for _, value := range values {
		trimmed := strings.TrimSpace(value)
		if trimmed == "" {
			continue
		}
		if _, ok := seen[trimmed]; ok {
			continue
		}
		seen[trimmed] = struct{}{}
		result = append(result, trimmed)
	}
	return result
}

func normalizeUserStatus(value *string) string {
	if value == nil || strings.TrimSpace(*value) == "" {
		return "ACTIVE"
	}
	return strings.ToUpper(strings.TrimSpace(*value))
}

func normalizeTrafficStrategy(value *string) string {
	if value == nil || strings.TrimSpace(*value) == "" {
		return "NO_RESET"
	}
	return strings.ToUpper(strings.TrimSpace(*value))
}

func isValidUserStatus(value string) bool {
	switch strings.ToUpper(strings.TrimSpace(value)) {
	case "ACTIVE", "DISABLED", "LIMITED", "EXPIRED":
		return true
	default:
		return false
	}
}

func isValidTrafficStrategy(value string) bool {
	switch strings.ToUpper(strings.TrimSpace(value)) {
	case "NO_RESET", "DAY", "WEEK", "MONTH":
		return true
	default:
		return false
	}
}

func normalizeNullableString(value *string) any {
	if value == nil || strings.TrimSpace(*value) == "" {
		return nil
	}
	return strings.TrimSpace(*value)
}

func normalizeUserTag(value *string) any {
	if value == nil || strings.TrimSpace(*value) == "" {
		return nil
	}
	return strings.ToUpper(strings.TrimSpace(*value))
}

func coalesceInt64(value *int64, fallback int64) int64 {
	if value == nil {
		return fallback
	}
	return *value
}

func coalesceUUID(value *string) string {
	if value != nil && strings.TrimSpace(*value) != "" {
		return strings.TrimSpace(*value)
	}
	return uuid.NewString()
}

func coalesceShortUUID(value *string) string {
	if value != nil && strings.TrimSpace(*value) != "" {
		return strings.TrimSpace(*value)
	}
	return generateSubscriptionShortUUID()
}

func coalesceRandomString(value *string, length int) string {
	if value != nil && strings.TrimSpace(*value) != "" {
		return strings.TrimSpace(*value)
	}
	return generateRandomString(length)
}

func generateRandomString(length int) string {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return ""
	}
	return hex.EncodeToString(bytes)[:length]
}

func generateSubscriptionShortUUID() string {
	bytes := make([]byte, 12)
	if _, err := rand.Read(bytes); err != nil {
		return ""
	}
	return base64.RawURLEncoding.EncodeToString(bytes)
}
