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
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"exodus/internal/config"
	dbmanager "exodus/internal/db/manager"
	"exodus/internal/httpapi/shared"
	monitor "exodus/internal/nodes"
	"exodus/internal/notifications"

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
	NaivePassword          string                  `json:"naivePassword"`
	ShadowtlsPassword      string                  `json:"shadowtlsPassword"`
	Hysteria2Password      string                  `json:"hysteria2Password"`
	AnytlsPassword         string                  `json:"anytlsPassword"`
	LastTriggeredThreshold int                     `json:"lastTriggeredThreshold"`
	SubRevokedAt           *time.Time              `json:"subRevokedAt"`
	LastTrafficResetAt     *time.Time              `json:"lastTrafficResetAt"`
	CreatedAt              time.Time               `json:"createdAt"`
	UpdatedAt              time.Time               `json:"updatedAt"`
	SubscriptionURL        string                  `json:"subscriptionUrl"`
	ActiveInternalSquads   []internalSquadResponse `json:"activeInternalSquads"`
	UserTraffic            userTrafficResponse     `json:"userTraffic"`
}

type userSubscriptionRequestHistoryRecord struct {
	ID        int64   `json:"id"`
	UserUUID  string  `json:"userUuid"`
	RequestIP *string `json:"requestIp"`
	UserAgent *string `json:"userAgent"`
	RequestAt string  `json:"requestAt"`
}

type userAccessibleNode struct {
	UUID              string                `json:"uuid"`
	NodeName          string                `json:"nodeName"`
	CountryCode       string                `json:"countryCode"`
	ConfigProfileUUID string                `json:"configProfileUuid"`
	ConfigProfileName string                `json:"configProfileName"`
	ActiveSquads      []userAccessibleSquad `json:"activeSquads"`
}

type userAccessibleSquad struct {
	SquadName      string   `json:"squadName"`
	ActiveInbounds []string `json:"activeInbounds"`
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
	LastTrafficResetAt       *time.Time
	SubRevokedAt             *time.Time
	TrojanPassword           string
	VlessUUID                string
	SSPassword               string
	NaivePassword            *string
	ShadowtlsPassword        *string
	Hysteria2Password        *string
	AnytlsPassword           *string
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
	NaivePassword        *string  `json:"naivePassword,omitempty"`
	ShadowtlsPassword    *string  `json:"shadowtlsPassword,omitempty"`
	Hysteria2Password    *string  `json:"hysteria2Password,omitempty"`
	AnytlsPassword       *string  `json:"anytlsPassword,omitempty"`
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
	TrojanPassword       OptionalString `json:"trojanPassword,omitempty"`
	VlessUUID            OptionalString `json:"vlessUuid,omitempty"`
	SSPassword           OptionalString `json:"ssPassword,omitempty"`
	NaivePassword        OptionalString `json:"naivePassword,omitempty"`
	ShadowtlsPassword    OptionalString `json:"shadowtlsPassword,omitempty"`
	Hysteria2Password    OptionalString `json:"hysteria2Password,omitempty"`
	AnytlsPassword       OptionalString `json:"anytlsPassword,omitempty"`
	ActiveInternalSquads *[]string      `json:"activeInternalSquads,omitempty"`
	ExternalSquadUUID    OptionalString `json:"externalSquadUuid,omitempty"`
}

type bulkDeleteUsersRequest struct {
	UUIDs []string `json:"uuids"`
}

type bulkExtendExpirationDateRequest struct {
	UUIDs      []string `json:"uuids"`
	ExtendDays int      `json:"extendDays"`
}

type bulkUpdateUsersRequest struct {
	UUIDs  []string              `json:"uuids"`
	Fields bulkUpdateUsersFields `json:"fields"`
}

type bulkUpdateUsersSquadsRequest struct {
	UUIDs                []string `json:"uuids"`
	ActiveInternalSquads []string `json:"activeInternalSquads"`
}

type bulkUpdateUsersFields struct {
	Status               *string        `json:"status,omitempty"`
	TrafficLimitBytes    *int64         `json:"trafficLimitBytes,omitempty"`
	TrafficLimitStrategy *string        `json:"trafficLimitStrategy,omitempty"`
	ExpireAt             *string        `json:"expireAt,omitempty"`
	Description          OptionalString `json:"description,omitempty"`
	Tag                  OptionalString `json:"tag,omitempty"`
	TelegramID           OptionalInt64  `json:"telegramId,omitempty"`
	Email                OptionalString `json:"email,omitempty"`
	HwidDeviceLimit      OptionalInt    `json:"hwidDeviceLimit,omitempty"`
	ExternalSquadUUID    OptionalString `json:"externalSquadUuid,omitempty"`
}

type bulkAllExtendExpirationDateRequest struct {
	ExtendDays int `json:"extendDays"`
}

type bulkAllUpdateUsersRequest struct {
	Status               *string        `json:"status,omitempty"`
	TrafficLimitBytes    *int64         `json:"trafficLimitBytes,omitempty"`
	TrafficLimitStrategy *string        `json:"trafficLimitStrategy,omitempty"`
	ExpireAt             *string        `json:"expireAt,omitempty"`
	Description          OptionalString `json:"description,omitempty"`
	Tag                  OptionalString `json:"tag,omitempty"`
	TelegramID           OptionalInt64  `json:"telegramId,omitempty"`
	Email                OptionalString `json:"email,omitempty"`
	HwidDeviceLimit      OptionalInt    `json:"hwidDeviceLimit,omitempty"`
}

type usersTableFilter struct {
	ID    string `json:"id"`
	Value any    `json:"value"`
}

type usersTableSorting struct {
	ID   string `json:"id"`
	Desc bool   `json:"desc"`
}

type usersTableQuery struct {
	Start       int
	Size        int
	Filters     []usersTableFilter
	FilterModes map[string]string
	Sorting     []usersTableSorting
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

		if len(parts) == 2 && parts[1] == "subscription-request-history" {
			if r.Method != http.MethodGet {
				shared.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
				return
			}
			handleGetUserSubscriptionRequestHistory(w, r, manager, cfg, userUUID)
			return
		}

		if len(parts) == 2 && parts[1] == "accessible-nodes" {
			if r.Method != http.MethodGet {
				shared.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
				return
			}
			handleGetUserAccessibleNodes(w, r, manager, cfg, userUUID)
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
		case "reset-traffic":
			handleBulkResetUsersTraffic(w, r, manager, cfg)
		case "update":
			handleBulkUpdateUsers(w, r, manager, cfg)
		case "update-squads":
			handleBulkUpdateUsersSquads(w, r, manager, cfg)
		case "extend-expiration-date":
			handleBulkExtendUsersExpirationDate(w, r, manager, cfg)
		case "all/reset-traffic":
			handleBulkAllResetUsersTraffic(w, r, manager, cfg)
		case "all/extend-expiration-date":
			handleBulkAllExtendUsersExpirationDate(w, r, manager, cfg)
		case "all/update":
			handleBulkAllUpdateUsers(w, r, manager, cfg)
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
	query, err := parseUsersTableQuery(r)
	if err != nil {
		shared.SendError(w, http.StatusBadRequest, "invalid table query", err, cfg)
		return
	}

	records, err := getAllUserRecords(r.Context(), manager)
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to fetch users", err, cfg)
		return
	}

	response, err := buildUserResponses(r.Context(), manager, records, resolveUsersSubscriptionBase(r.Context(), manager, r, cfg))
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to build users response", err, cfg)
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

	response, err := buildUserResponses(r.Context(), manager, []userRecord{record}, resolveUsersSubscriptionBase(r.Context(), manager, r, cfg))
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
	naivePassword := coalesceRandomString(req.NaivePassword, 16)
	shadowtlsPassword := coalesceRandomString(req.ShadowtlsPassword, 16)
	hysteria2Password := coalesceRandomString(req.Hysteria2Password, 16)
	anytlsPassword := coalesceRandomString(req.AnytlsPassword, 16)

	if trojanPassword == "" || ssPassword == "" || naivePassword == "" || shadowtlsPassword == "" || hysteria2Password == "" || anytlsPassword == "" {
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

	internalSquadNodeUUIDs := make([]string, 0)
	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		tx, err := db.BeginTx(r.Context(), nil)
		if err != nil {
			return err
		}

		var tID int64
		insertErr := tx.QueryRowContext(r.Context(), `
			INSERT INTO users (
					uuid, short_uuid, username, status, traffic_limit_bytes, traffic_limit_strategy,
					expire_at, last_traffic_reset_at, sub_revoked_at,
					trojan_password, vless_uuid, ss_password, naive_password, shadowtls_password, hysteria2_password, anytls_password,
					description, tag, telegram_id, email,
					hwid_device_limit, external_squad_uuid, last_triggered_threshold, created_at, updated_at
				) VALUES (?, ?, ?, ?, ?, ?, ?, ?, NULL, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, 0, ?, ?)
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
			naivePassword,
			shadowtlsPassword,
			hysteria2Password,
			anytlsPassword,
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
		requestedSquads := dedupeStrings(req.ActiveInternalSquads)
		if len(requestedSquads) > 0 {
			nodeUUIDs, nodeTargetsErr := resolveNodeUUIDsForInternalSquadsTx(r.Context(), tx, requestedSquads)
			if nodeTargetsErr != nil {
				_ = tx.Rollback()
				return nodeTargetsErr
			}
			internalSquadNodeUUIDs = nodeUUIDs
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
	response, err := buildUserResponses(r.Context(), manager, []userRecord{record}, resolveUsersSubscriptionBase(r.Context(), manager, r, cfg))
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to build created user response", err, cfg)
		return
	}

	if strings.EqualFold(normalizeUserStatus(req.Status), "ACTIVE") && len(internalSquadNodeUUIDs) > 0 {
		monitor.RequestNodeDeploy(true, internalSquadNodeUUIDs...)
	}
	emitUserNotification(r.Context(), manager, cfg, notifications.EventUserCreated, record, nil)
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
	statusNodeUUIDs := make([]string, 0)
	statusToSet, shouldSetStatus := plannedUserStatusForUpdate(record, req, time.Now().UTC())
	statusDeployRequired := shouldSetStatus && userConfigPresenceChanges(record.Status, statusToSet)
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

		if statusDeployRequired {
			nodeUUIDs, nodeTargetsErr := resolveNodeUUIDsForUserUUIDsTx(r.Context(), tx, []string{targetUUID})
			if nodeTargetsErr != nil {
				_ = tx.Rollback()
				return nodeTargetsErr
			}
			statusNodeUUIDs = nodeUUIDs
		}

		if shouldSetStatus {
			add("status", statusToSet)
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
		addOptionalCredential := func(field OptionalString, column string, nullable bool) {
			if !field.Set {
				return
			}
			if field.Value == nil {
				if nullable {
					clauses = append(clauses, fmt.Sprintf("%s = NULL", column))
				}
				return
			}
			add(column, strings.TrimSpace(*field.Value))
		}
		addOptionalCredential(req.TrojanPassword, "trojan_password", false)
		addOptionalCredential(req.VlessUUID, "vless_uuid", false)
		addOptionalCredential(req.SSPassword, "ss_password", false)
		addOptionalCredential(req.NaivePassword, "naive_password", true)
		addOptionalCredential(req.ShadowtlsPassword, "shadowtls_password", true)
		addOptionalCredential(req.Hysteria2Password, "hysteria2_password", true)
		addOptionalCredential(req.AnytlsPassword, "anytls_password", true)
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
	response, err := buildUserResponses(r.Context(), manager, []userRecord{updatedRecord}, resolveUsersSubscriptionBase(r.Context(), manager, r, cfg))
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to build updated user response", err, cfg)
		return
	}

	deployNodeUUIDs := dedupeStrings(append(statusNodeUUIDs, internalSquadNodeUUIDs...))
	if (internalSquadsChanged || statusDeployRequired) && len(deployNodeUUIDs) > 0 {
		monitor.RequestNodeDeploy(true, deployNodeUUIDs...)
	}
	emitUserNotification(r.Context(), manager, cfg, notifications.EventUserModified, updatedRecord, nil)
	if statusChanged := userStatusChangedNotification(record.Status, updatedRecord.Status); statusChanged != "" {
		emitUserNotification(r.Context(), manager, cfg, statusChanged, updatedRecord, nil)
	}
	shared.WriteJSON(w, http.StatusOK, map[string]any{"response": response[0]})
}

func handleGetUserSubscriptionRequestHistory(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig, userUUID string) {
	records := make([]userSubscriptionRequestHistoryRecord, 0)
	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		var exists bool
		if err := db.QueryRowContext(r.Context(), `SELECT EXISTS(SELECT 1 FROM users WHERE uuid = ?)`, userUUID).Scan(&exists); err != nil {
			return err
		}
		if !exists {
			return errUserNotFound
		}

		rows, err := db.QueryContext(r.Context(), `
			SELECT id, user_uuid, request_ip, user_agent, request_at
			FROM user_subscription_request_history
			WHERE user_uuid = ?
			ORDER BY request_at DESC
			LIMIT 24
		`, userUUID)
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			var item userSubscriptionRequestHistoryRecord
			var requestAt time.Time
			if scanErr := rows.Scan(&item.ID, &item.UserUUID, &item.RequestIP, &item.UserAgent, &requestAt); scanErr != nil {
				return scanErr
			}
			item.RequestAt = requestAt.UTC().Format("2006-01-02T15:04:05.000Z")
			records = append(records, item)
		}
		return rows.Err()
	})
	if err != nil {
		if errors.Is(err, errUserNotFound) {
			shared.SendError(w, http.StatusNotFound, "user not found", nil, cfg)
			return
		}
		shared.SendError(w, http.StatusInternalServerError, "failed to fetch user subscription request history", err, cfg)
		return
	}

	shared.WriteJSON(w, http.StatusOK, map[string]any{
		"response": map[string]any{
			"records": records,
			"total":   len(records),
		},
	})
}

func handleGetUserAccessibleNodes(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig, userUUID string) {
	activeNodes := make([]userAccessibleNode, 0)
	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		var userID int64
		if err := db.QueryRowContext(r.Context(), `SELECT t_id FROM users WHERE uuid = ?`, userUUID).Scan(&userID); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return errUserNotFound
			}
			return err
		}

		rows, err := db.QueryContext(r.Context(), `
			SELECT
				n.uuid,
				n.name,
				n.country_code,
				cp.uuid,
				cp.name,
				sq.uuid,
				sq.name,
				cpi.tag
			FROM nodes n
			INNER JOIN config_profiles cp ON cp.uuid = n.active_config_profile_uuid
			INNER JOIN config_profile_inbounds cpi ON cpi.profile_uuid = cp.uuid
			INNER JOIN config_profile_inbounds_to_nodes cpin
				ON cpin.config_profile_inbound_uuid = cpi.uuid
				AND cpin.node_uuid = n.uuid
			INNER JOIN internal_squad_inbounds isi ON isi.inbound_uuid = cpi.uuid
			INNER JOIN internal_squads sq ON sq.uuid = isi.internal_squad_uuid
			INNER JOIN internal_squad_members ism
				ON ism.internal_squad_uuid = sq.uuid
				AND ism.user_id = ?
			ORDER BY n.view_position ASC, sq.view_position ASC, cpi.tag ASC
		`, userID)
		if err != nil {
			return err
		}
		defer rows.Close()

		nodeIndexes := make(map[string]int)
		squadIndexesByNode := make(map[string]map[string]int)
		for rows.Next() {
			var nodeUUID, nodeName, countryCode, profileUUID, profileName string
			var squadUUID, squadName, inboundTag string
			if scanErr := rows.Scan(&nodeUUID, &nodeName, &countryCode, &profileUUID, &profileName, &squadUUID, &squadName, &inboundTag); scanErr != nil {
				return scanErr
			}

			nodeIndex, ok := nodeIndexes[nodeUUID]
			if !ok {
				activeNodes = append(activeNodes, userAccessibleNode{
					UUID:              nodeUUID,
					NodeName:          nodeName,
					CountryCode:       countryCode,
					ConfigProfileUUID: profileUUID,
					ConfigProfileName: profileName,
					ActiveSquads:      make([]userAccessibleSquad, 0),
				})
				nodeIndex = len(activeNodes) - 1
				nodeIndexes[nodeUUID] = nodeIndex
				squadIndexesByNode[nodeUUID] = make(map[string]int)
			}

			squadIndexes := squadIndexesByNode[nodeUUID]
			squadIndex, ok := squadIndexes[squadUUID]
			if !ok {
				activeNodes[nodeIndex].ActiveSquads = append(activeNodes[nodeIndex].ActiveSquads, userAccessibleSquad{
					SquadName:      squadName,
					ActiveInbounds: make([]string, 0),
				})
				squadIndex = len(activeNodes[nodeIndex].ActiveSquads) - 1
				squadIndexes[squadUUID] = squadIndex
			}

			activeNodes[nodeIndex].ActiveSquads[squadIndex].ActiveInbounds = append(
				activeNodes[nodeIndex].ActiveSquads[squadIndex].ActiveInbounds,
				inboundTag,
			)
		}
		return rows.Err()
	})
	if err != nil {
		if errors.Is(err, errUserNotFound) {
			shared.SendError(w, http.StatusNotFound, "user not found", nil, cfg)
			return
		}
		shared.SendError(w, http.StatusInternalServerError, "failed to fetch user accessible nodes", err, cfg)
		return
	}

	shared.WriteJSON(w, http.StatusOK, map[string]any{
		"response": map[string]any{
			"userUuid":    userUUID,
			"activeNodes": activeNodes,
		},
	})
}

func handleDeleteUser(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig, userUUID string) {
	record, recordErr := getUserRecordByUUID(r.Context(), manager, userUUID)
	if recordErr != nil && !errors.Is(recordErr, errUserNotFound) {
		shared.SendError(w, http.StatusInternalServerError, "failed to fetch user", recordErr, cfg)
		return
	}
	internalSquadNodeUUIDs := make([]string, 0)
	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		tx, err := db.BeginTx(r.Context(), nil)
		if err != nil {
			return err
		}

		var tID int64
		if err := tx.QueryRowContext(r.Context(), `SELECT t_id FROM users WHERE uuid = ?`, userUUID).Scan(&tID); err != nil {
			_ = tx.Rollback()
			if errors.Is(err, sql.ErrNoRows) {
				return errUserNotFound
			}
			return err
		}

		currentSquads, loadErr := getUserInternalSquadsTx(r.Context(), tx, tID)
		if loadErr != nil {
			_ = tx.Rollback()
			return loadErr
		}

		nodeUUIDs, nodeTargetsErr := resolveNodeUUIDsForInternalSquadsTx(r.Context(), tx, currentSquads)
		if nodeTargetsErr != nil {
			_ = tx.Rollback()
			return nodeTargetsErr
		}
		internalSquadNodeUUIDs = nodeUUIDs

		result, err := tx.ExecContext(r.Context(), `DELETE FROM users WHERE uuid = ?`, userUUID)
		if err != nil {
			_ = tx.Rollback()
			return err
		}
		rows, err := result.RowsAffected()
		if err != nil {
			_ = tx.Rollback()
			return err
		}
		if rows == 0 {
			_ = tx.Rollback()
			return errUserNotFound
		}

		return tx.Commit()
	})
	if err != nil {
		if errors.Is(err, errUserNotFound) {
			shared.SendError(w, http.StatusNotFound, "user not found", nil, cfg)
			return
		}
		shared.SendError(w, http.StatusInternalServerError, "failed to delete user", err, cfg)
		return
	}

	if len(internalSquadNodeUUIDs) > 0 {
		monitor.RequestNodeDeploy(true, internalSquadNodeUUIDs...)
	}
	if recordErr == nil {
		emitUserNotification(r.Context(), manager, cfg, notifications.EventUserDeleted, record, nil)
	}
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

func resolveNodeUUIDsForUserUUIDsTx(ctx context.Context, tx dbmanager.TxExecutor, userUUIDs []string) ([]string, error) {
	cleanUserUUIDs := dedupeStrings(userUUIDs)
	if len(cleanUserUUIDs) == 0 {
		return []string{}, nil
	}

	rows, err := tx.QueryContext(ctx, `
		SELECT DISTINCT cpitn.node_uuid
		FROM users u
		JOIN internal_squad_members ism ON ism.user_id = u.t_id
		JOIN internal_squad_inbounds isi ON isi.internal_squad_uuid = ism.internal_squad_uuid
		JOIN config_profile_inbounds_to_nodes cpitn ON cpitn.config_profile_inbound_uuid = isi.inbound_uuid
		WHERE u.uuid = ANY(?)
	`, cleanUserUUIDs)
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

	notificationRecords, err := getUserRecordsByUUIDs(r.Context(), manager, req.UUIDs)
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to fetch users", err, cfg)
		return
	}

	internalSquadNodeUUIDs := make([]string, 0)
	err = manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		tx, err := db.BeginTx(r.Context(), nil)
		if err != nil {
			return err
		}

		nodeUUIDs, nodeTargetsErr := resolveNodeUUIDsForUserUUIDsTx(r.Context(), tx, req.UUIDs)
		if nodeTargetsErr != nil {
			_ = tx.Rollback()
			return nodeTargetsErr
		}
		internalSquadNodeUUIDs = nodeUUIDs

		if _, err := tx.ExecContext(r.Context(), `DELETE FROM users WHERE uuid = ANY(?)`, req.UUIDs); err != nil {
			_ = tx.Rollback()
			return err
		}

		return tx.Commit()
	})
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to delete users", err, cfg)
		return
	}

	if len(internalSquadNodeUUIDs) > 0 {
		monitor.RequestNodeDeploy(true, internalSquadNodeUUIDs...)
	}
	emitUsersNotificationFromRecords(r.Context(), manager, cfg, notifications.EventUserDeleted, req.UUIDs, notificationRecords)
	shared.WriteJSON(w, http.StatusOK, map[string]any{"response": map[string]any{"isDeleted": true}})
}

func handleBulkResetUsersTraffic(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) {
	var req bulkDeleteUsersRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		shared.SendError(w, http.StatusBadRequest, "invalid JSON", err, cfg)
		return
	}
	if err := validateUUIDList(req.UUIDs); err != nil {
		shared.SendError(w, http.StatusBadRequest, err.Error(), nil, cfg)
		return
	}

	affectedRows, nodeUUIDs, err := resetUsersTrafficByUUIDs(r.Context(), manager, req.UUIDs)
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to reset users traffic", err, cfg)
		return
	}
	if len(nodeUUIDs) > 0 {
		monitor.RequestNodeDeploy(true, nodeUUIDs...)
	}
	emitUsersByUUIDsNotification(r.Context(), manager, cfg, notifications.EventUserTrafficReset, req.UUIDs)

	shared.WriteJSON(w, http.StatusOK, map[string]any{"response": map[string]any{"affectedRows": affectedRows}})
}

func handleBulkAllResetUsersTraffic(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) {
	affectedRows, nodeUUIDs, err := resetAllUsersTraffic(r.Context(), manager)
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to reset all users traffic", err, cfg)
		return
	}
	if len(nodeUUIDs) > 0 {
		monitor.RequestNodeDeploy(true, nodeUUIDs...)
	}
	emitBulkSummaryNotification(r.Context(), cfg, notifications.EventUserTrafficReset, affectedRows)

	shared.WriteJSON(w, http.StatusOK, map[string]any{"response": map[string]any{"affectedRows": affectedRows}})
}

func handleBulkExtendUsersExpirationDate(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) {
	var req bulkExtendExpirationDateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		shared.SendError(w, http.StatusBadRequest, "invalid JSON", err, cfg)
		return
	}
	if err := validateUUIDList(req.UUIDs); err != nil {
		shared.SendError(w, http.StatusBadRequest, err.Error(), nil, cfg)
		return
	}
	if err := validateExtendDays(req.ExtendDays); err != nil {
		shared.SendError(w, http.StatusBadRequest, err.Error(), nil, cfg)
		return
	}

	affectedRows, nodeUUIDs, err := extendUsersExpirationByUUIDs(r.Context(), manager, req.UUIDs, req.ExtendDays)
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to extend users expiration date", err, cfg)
		return
	}
	if len(nodeUUIDs) > 0 {
		monitor.RequestNodeDeploy(true, nodeUUIDs...)
	}
	emitUsersByUUIDsNotification(r.Context(), manager, cfg, notifications.EventUserModified, req.UUIDs)

	shared.WriteJSON(w, http.StatusOK, map[string]any{"response": map[string]any{"affectedRows": affectedRows}})
}

func handleBulkAllExtendUsersExpirationDate(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) {
	var req bulkAllExtendExpirationDateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		shared.SendError(w, http.StatusBadRequest, "invalid JSON", err, cfg)
		return
	}
	if err := validateExtendDays(req.ExtendDays); err != nil {
		shared.SendError(w, http.StatusBadRequest, err.Error(), nil, cfg)
		return
	}

	affectedRows, nodeUUIDs, err := extendAllUsersExpiration(r.Context(), manager, req.ExtendDays)
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to extend all users expiration date", err, cfg)
		return
	}
	if len(nodeUUIDs) > 0 {
		monitor.RequestNodeDeploy(true, nodeUUIDs...)
	}
	emitBulkSummaryNotification(r.Context(), cfg, notifications.EventUserModified, affectedRows)

	shared.WriteJSON(w, http.StatusOK, map[string]any{"response": map[string]any{"affectedRows": affectedRows}})
}

func handleBulkUpdateUsers(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) {
	var req bulkUpdateUsersRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		shared.SendError(w, http.StatusBadRequest, "invalid JSON", err, cfg)
		return
	}
	if err := validateUUIDList(req.UUIDs); err != nil {
		shared.SendError(w, http.StatusBadRequest, err.Error(), nil, cfg)
		return
	}
	if err := validateBulkUpdateUsersFields(req.Fields); err != nil {
		shared.SendError(w, http.StatusBadRequest, err.Error(), nil, cfg)
		return
	}

	clauses, args := buildBulkUpdateUserClauses(req.Fields)
	if len(clauses) == 0 {
		shared.SendError(w, http.StatusBadRequest, "at least one field must be provided", nil, cfg)
		return
	}

	cleanUUIDs := dedupeStrings(req.UUIDs)
	var affectedRows int64
	nodeUUIDs := make([]string, 0)
	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		tx, err := db.BeginTx(r.Context(), nil)
		if err != nil {
			return err
		}

		targets, nodeTargetsErr := resolveNodeUUIDsForUserUUIDsTx(r.Context(), tx, cleanUUIDs)
		if nodeTargetsErr != nil {
			_ = tx.Rollback()
			return nodeTargetsErr
		}
		nodeUUIDs = targets

		queryArgs := append(args, cleanUUIDs)
		query := fmt.Sprintf("UPDATE users SET %s, updated_at = CURRENT_TIMESTAMP WHERE uuid = ANY(?)", strings.Join(clauses, ", "))
		result, execErr := tx.ExecContext(r.Context(), query, queryArgs...)
		if execErr != nil {
			_ = tx.Rollback()
			return mapUserWriteError(execErr)
		}
		rows, rowsErr := result.RowsAffected()
		if rowsErr != nil {
			_ = tx.Rollback()
			return rowsErr
		}
		affectedRows = rows

		return tx.Commit()
	})
	if err != nil {
		handleUserWriteError(w, err, cfg)
		return
	}

	if affectedRows > 0 {
		if len(nodeUUIDs) > 0 {
			monitor.RequestNodeDeploy(true, nodeUUIDs...)
		}
		emitUsersByUUIDsNotification(r.Context(), manager, cfg, notifications.EventUserModified, cleanUUIDs)
	}

	shared.WriteJSON(w, http.StatusOK, map[string]any{"response": map[string]any{"affectedRows": affectedRows}})
}

func handleBulkUpdateUsersSquads(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) {
	var req bulkUpdateUsersSquadsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		shared.SendError(w, http.StatusBadRequest, "invalid JSON", err, cfg)
		return
	}
	if err := validateUUIDList(req.UUIDs); err != nil {
		shared.SendError(w, http.StatusBadRequest, err.Error(), nil, cfg)
		return
	}
	if err := validateUUIDListAllowEmpty(req.ActiveInternalSquads); err != nil {
		shared.SendError(w, http.StatusBadRequest, "invalid activeInternalSquads", err, cfg)
		return
	}

	cleanUserUUIDs := dedupeStrings(req.UUIDs)
	requestedSquads := dedupeStrings(req.ActiveInternalSquads)
	var affectedRows int64
	nodeUUIDs := make([]string, 0)
	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		tx, err := db.BeginTx(r.Context(), nil)
		if err != nil {
			return err
		}

		targets, nodeTargetsErr := resolveNodeUUIDsForUserUUIDsTx(r.Context(), tx, cleanUserUUIDs)
		if nodeTargetsErr != nil {
			_ = tx.Rollback()
			return nodeTargetsErr
		}
		squadTargets, squadTargetsErr := resolveNodeUUIDsForInternalSquadsTx(r.Context(), tx, requestedSquads)
		if squadTargetsErr != nil {
			_ = tx.Rollback()
			return squadTargetsErr
		}
		nodeUUIDs = dedupeStrings(append(targets, squadTargets...))

		rows, err := tx.QueryContext(r.Context(), `SELECT t_id FROM users WHERE uuid = ANY(?)`, cleanUserUUIDs)
		if err != nil {
			_ = tx.Rollback()
			return err
		}
		userIDs := make([]int64, 0, len(cleanUserUUIDs))
		for rows.Next() {
			var userID int64
			if err := rows.Scan(&userID); err != nil {
				_ = rows.Close()
				_ = tx.Rollback()
				return err
			}
			userIDs = append(userIDs, userID)
		}
		if err := rows.Err(); err != nil {
			_ = rows.Close()
			_ = tx.Rollback()
			return err
		}
		_ = rows.Close()

		for _, userID := range userIDs {
			if err := replaceUserInternalSquadsTx(r.Context(), tx, userID, requestedSquads); err != nil {
				_ = tx.Rollback()
				return err
			}
		}
		if len(userIDs) > 0 {
			if _, err := tx.ExecContext(r.Context(), `UPDATE users SET updated_at = CURRENT_TIMESTAMP WHERE t_id = ANY(?)`, userIDs); err != nil {
				_ = tx.Rollback()
				return err
			}
		}
		affectedRows = int64(len(userIDs))

		return tx.Commit()
	})
	if err != nil {
		handleUserWriteError(w, err, cfg)
		return
	}

	if affectedRows > 0 {
		if len(nodeUUIDs) > 0 {
			monitor.RequestNodeDeploy(true, nodeUUIDs...)
		}
		emitUsersByUUIDsNotification(r.Context(), manager, cfg, notifications.EventUserModified, cleanUserUUIDs)
	}

	shared.WriteJSON(w, http.StatusOK, map[string]any{"response": map[string]any{"affectedRows": affectedRows}})
}

func handleBulkAllUpdateUsers(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) {
	var req bulkAllUpdateUsersRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		shared.SendError(w, http.StatusBadRequest, "invalid JSON", err, cfg)
		return
	}
	if err := validateBulkAllUpdateUsersRequest(req); err != nil {
		shared.SendError(w, http.StatusBadRequest, err.Error(), nil, cfg)
		return
	}

	clauses, args := buildBulkUpdateUserClauses(bulkUpdateUsersFields{
		Status:               req.Status,
		TrafficLimitBytes:    req.TrafficLimitBytes,
		TrafficLimitStrategy: req.TrafficLimitStrategy,
		ExpireAt:             req.ExpireAt,
		Description:          req.Description,
		Tag:                  req.Tag,
		TelegramID:           req.TelegramID,
		Email:                req.Email,
		HwidDeviceLimit:      req.HwidDeviceLimit,
	})

	var affectedRows int64
	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		query := fmt.Sprintf("UPDATE users SET %s, updated_at = CURRENT_TIMESTAMP", strings.Join(clauses, ", "))
		result, execErr := db.ExecContext(r.Context(), query, args...)
		if execErr != nil {
			return mapUserWriteError(execErr)
		}
		rows, rowsErr := result.RowsAffected()
		if rowsErr != nil {
			return rowsErr
		}
		affectedRows = rows
		return nil
	})
	if err != nil {
		handleUserWriteError(w, err, cfg)
		return
	}

	if affectedRows > 0 {
		monitor.RequestNodeDeploy(true)
		emitBulkSummaryNotification(r.Context(), cfg, notifications.EventUserModified, affectedRows)
	}
	shared.WriteJSON(w, http.StatusOK, map[string]any{"response": map[string]any{"eventSent": affectedRows > 0}})
}

func handleEnableUser(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig, userUUID string) {
	nodeUUIDs, err := updateUserStatus(r.Context(), manager, userUUID, "ACTIVE")
	if err != nil {
		handleUserActionError(w, err, cfg, "failed to enable user")
		return
	}
	if len(nodeUUIDs) > 0 {
		monitor.RequestNodeDeploy(true, nodeUUIDs...)
	}
	if record, loadErr := getUserRecordByUUID(r.Context(), manager, userUUID); loadErr == nil {
		emitUserNotification(r.Context(), manager, cfg, notifications.EventUserEnabled, record, nil)
	}
	sendUpdatedUserResponse(w, r, manager, cfg, userUUID)
}

func handleDisableUser(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig, userUUID string) {
	nodeUUIDs, err := updateUserStatus(r.Context(), manager, userUUID, "DISABLED")
	if err != nil {
		handleUserActionError(w, err, cfg, "failed to disable user")
		return
	}
	if len(nodeUUIDs) > 0 {
		monitor.RequestNodeDeploy(true, nodeUUIDs...)
	}
	if record, loadErr := getUserRecordByUUID(r.Context(), manager, userUUID); loadErr == nil {
		emitUserNotification(r.Context(), manager, cfg, notifications.EventUserDisabled, record, nil)
	}
	sendUpdatedUserResponse(w, r, manager, cfg, userUUID)
}

func handleResetUserTraffic(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig, userUUID string) {
	nodeUUIDs := make([]string, 0)
	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		tx, err := db.BeginTx(r.Context(), nil)
		if err != nil {
			return err
		}

		var (
			tID    int64
			status string
		)
		if err := tx.QueryRowContext(r.Context(), `SELECT t_id, status FROM users WHERE uuid = ?`, userUUID).Scan(&tID, &status); err != nil {
			_ = tx.Rollback()
			if errors.Is(err, sql.ErrNoRows) {
				return errUserNotFound
			}
			return err
		}

		if strings.EqualFold(status, "LIMITED") {
			targets, nodeTargetsErr := resolveNodeUUIDsForUserUUIDsTx(r.Context(), tx, []string{userUUID})
			if nodeTargetsErr != nil {
				_ = tx.Rollback()
				return nodeTargetsErr
			}
			nodeUUIDs = targets
		}

		if _, err := tx.ExecContext(r.Context(), `
			UPDATE users
			SET last_traffic_reset_at = CURRENT_TIMESTAMP,
			    last_triggered_threshold = 0,
			    status = CASE WHEN status = 'LIMITED' THEN 'ACTIVE' ELSE status END,
			    updated_at = CURRENT_TIMESTAMP
			WHERE uuid = ?
		`, userUUID); err != nil {
			_ = tx.Rollback()
			return err
		}

		if _, err := tx.ExecContext(r.Context(), `
			INSERT INTO user_traffic (t_id, used_traffic_bytes, lifetime_used_traffic_bytes)
			VALUES (?, 0, 0)
			ON CONFLICT (t_id)
			DO UPDATE SET used_traffic_bytes = 0
		`, tID); err != nil {
			_ = tx.Rollback()
			return err
		}

		return tx.Commit()
	})
	if err != nil {
		handleUserActionError(w, err, cfg, "failed to reset user traffic")
		return
	}
	if len(nodeUUIDs) > 0 {
		monitor.RequestNodeDeploy(true, nodeUUIDs...)
	}
	if record, loadErr := getUserRecordByUUID(r.Context(), manager, userUUID); loadErr == nil {
		emitUserNotification(r.Context(), manager, cfg, notifications.EventUserTrafficReset, record, nil)
		if strings.EqualFold(record.Status, "ACTIVE") {
			emitUserNotification(r.Context(), manager, cfg, notifications.EventUserEnabled, record, map[string]any{"reason": "traffic_reset"})
		}
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
	if record, loadErr := getUserRecordByUUID(r.Context(), manager, userUUID); loadErr == nil {
		emitUserNotification(r.Context(), manager, cfg, notifications.EventUserRevoked, record, nil)
	}
	sendUpdatedUserResponse(w, r, manager, cfg, userUUID)
}

func sendUpdatedUserResponse(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig, userUUID string) {
	record, err := getUserRecordByUUID(r.Context(), manager, userUUID)
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to fetch updated user", err, cfg)
		return
	}
	response, err := buildUserResponses(r.Context(), manager, []userRecord{record}, resolveUsersSubscriptionBase(r.Context(), manager, r, cfg))
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

func updateUserStatus(ctx context.Context, manager *dbmanager.DatabaseManager, userUUID string, status string) ([]string, error) {
	nodeUUIDs := make([]string, 0)
	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		tx, err := db.BeginTx(ctx, nil)
		if err != nil {
			return err
		}

		targets, nodeTargetsErr := resolveNodeUUIDsForUserUUIDsTx(ctx, tx, []string{userUUID})
		if nodeTargetsErr != nil {
			_ = tx.Rollback()
			return nodeTargetsErr
		}
		nodeUUIDs = targets

		result, err := tx.ExecContext(ctx, `
			UPDATE users
			SET status = ?, updated_at = CURRENT_TIMESTAMP
			WHERE uuid = ?
		`, status, userUUID)
		if err != nil {
			_ = tx.Rollback()
			return err
		}
		rows, err := result.RowsAffected()
		if err != nil {
			_ = tx.Rollback()
			return err
		}
		if rows == 0 {
			_ = tx.Rollback()
			return errUserNotFound
		}
		return tx.Commit()
	})
	return nodeUUIDs, err
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

func queryLimitedUserNodeUUIDsTx(ctx context.Context, tx dbmanager.TxExecutor, userUUIDs []string) ([]string, error) {
	rows, err := tx.QueryContext(ctx, `
		SELECT DISTINCT cpitn.node_uuid
		FROM users u
		JOIN internal_squad_members ism ON ism.user_id = u.t_id
		JOIN internal_squad_inbounds isi ON isi.internal_squad_uuid = ism.internal_squad_uuid
		JOIN config_profile_inbounds_to_nodes cpitn ON cpitn.config_profile_inbound_uuid = isi.inbound_uuid
		WHERE u.status = 'LIMITED' AND u.uuid = ANY(?)
	`, dedupeStrings(userUUIDs))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanNodeUUIDRows(rows)
}

func queryAllLimitedUserNodeUUIDsTx(ctx context.Context, tx dbmanager.TxExecutor) ([]string, error) {
	rows, err := tx.QueryContext(ctx, `
		SELECT DISTINCT cpitn.node_uuid
		FROM users u
		JOIN internal_squad_members ism ON ism.user_id = u.t_id
		JOIN internal_squad_inbounds isi ON isi.internal_squad_uuid = ism.internal_squad_uuid
		JOIN config_profile_inbounds_to_nodes cpitn ON cpitn.config_profile_inbound_uuid = isi.inbound_uuid
		WHERE u.status = 'LIMITED'
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanNodeUUIDRows(rows)
}

func queryReactivatedExpiredUserNodeUUIDsTx(ctx context.Context, tx dbmanager.TxExecutor, userUUIDs []string, extendDays int) ([]string, error) {
	rows, err := tx.QueryContext(ctx, `
		SELECT DISTINCT cpitn.node_uuid
		FROM users u
		JOIN internal_squad_members ism ON ism.user_id = u.t_id
		JOIN internal_squad_inbounds isi ON isi.internal_squad_uuid = ism.internal_squad_uuid
		JOIN config_profile_inbounds_to_nodes cpitn ON cpitn.config_profile_inbound_uuid = isi.inbound_uuid
		WHERE u.status = 'EXPIRED'
		  AND u.uuid = ANY(?)
		  AND u.expire_at + (?::int * INTERVAL '1 day') > CURRENT_TIMESTAMP
	`, dedupeStrings(userUUIDs), extendDays)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanNodeUUIDRows(rows)
}

func queryAllReactivatedExpiredUserNodeUUIDsTx(ctx context.Context, tx dbmanager.TxExecutor, extendDays int) ([]string, error) {
	rows, err := tx.QueryContext(ctx, `
		SELECT DISTINCT cpitn.node_uuid
		FROM users u
		JOIN internal_squad_members ism ON ism.user_id = u.t_id
		JOIN internal_squad_inbounds isi ON isi.internal_squad_uuid = ism.internal_squad_uuid
		JOIN config_profile_inbounds_to_nodes cpitn ON cpitn.config_profile_inbound_uuid = isi.inbound_uuid
		WHERE u.status = 'EXPIRED'
		  AND u.expire_at + (?::int * INTERVAL '1 day') > CURRENT_TIMESTAMP
	`, extendDays)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanNodeUUIDRows(rows)
}

func scanNodeUUIDRows(rows *sql.Rows) ([]string, error) {
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

func resetUsersTrafficByUUIDs(ctx context.Context, manager *dbmanager.DatabaseManager, userUUIDs []string) (int64, []string, error) {
	var affectedRows int64
	nodeUUIDs := make([]string, 0)

	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		tx, err := db.BeginTx(ctx, nil)
		if err != nil {
			return err
		}

		targets, nodeTargetsErr := queryLimitedUserNodeUUIDsTx(ctx, tx, userUUIDs)
		if nodeTargetsErr != nil {
			_ = tx.Rollback()
			return nodeTargetsErr
		}
		nodeUUIDs = targets

		result, err := tx.ExecContext(ctx, `
			UPDATE users
			SET last_traffic_reset_at = CURRENT_TIMESTAMP,
			    last_triggered_threshold = 0,
			    status = CASE WHEN status = 'LIMITED' THEN 'ACTIVE' ELSE status END,
			    updated_at = CURRENT_TIMESTAMP
			WHERE uuid = ANY(?)
		`, dedupeStrings(userUUIDs))
		if err != nil {
			_ = tx.Rollback()
			return err
		}
		affectedRows, _ = result.RowsAffected()

		if _, err := tx.ExecContext(ctx, `
			UPDATE user_traffic
			SET used_traffic_bytes = 0
			WHERE t_id IN (SELECT t_id FROM users WHERE uuid = ANY(?))
		`, dedupeStrings(userUUIDs)); err != nil {
			_ = tx.Rollback()
			return err
		}

		return tx.Commit()
	})

	return affectedRows, nodeUUIDs, err
}

func resetAllUsersTraffic(ctx context.Context, manager *dbmanager.DatabaseManager) (int64, []string, error) {
	var affectedRows int64
	nodeUUIDs := make([]string, 0)

	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		tx, err := db.BeginTx(ctx, nil)
		if err != nil {
			return err
		}

		targets, nodeTargetsErr := queryAllLimitedUserNodeUUIDsTx(ctx, tx)
		if nodeTargetsErr != nil {
			_ = tx.Rollback()
			return nodeTargetsErr
		}
		nodeUUIDs = targets

		result, err := tx.ExecContext(ctx, `
			UPDATE users
			SET last_traffic_reset_at = CURRENT_TIMESTAMP,
			    last_triggered_threshold = 0,
			    status = CASE WHEN status = 'LIMITED' THEN 'ACTIVE' ELSE status END,
			    updated_at = CURRENT_TIMESTAMP
		`)
		if err != nil {
			_ = tx.Rollback()
			return err
		}
		affectedRows, _ = result.RowsAffected()

		if _, err := tx.ExecContext(ctx, `UPDATE user_traffic SET used_traffic_bytes = 0`); err != nil {
			_ = tx.Rollback()
			return err
		}

		return tx.Commit()
	})

	return affectedRows, nodeUUIDs, err
}

func extendUsersExpirationByUUIDs(ctx context.Context, manager *dbmanager.DatabaseManager, userUUIDs []string, extendDays int) (int64, []string, error) {
	var affectedRows int64
	nodeUUIDs := make([]string, 0)

	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		tx, err := db.BeginTx(ctx, nil)
		if err != nil {
			return err
		}

		targets, nodeTargetsErr := queryReactivatedExpiredUserNodeUUIDsTx(ctx, tx, userUUIDs, extendDays)
		if nodeTargetsErr != nil {
			_ = tx.Rollback()
			return nodeTargetsErr
		}
		nodeUUIDs = targets

		result, err := tx.ExecContext(ctx, `
			UPDATE users
			SET expire_at = expire_at + (?::int * INTERVAL '1 day'),
			    status = CASE
			        WHEN status = 'EXPIRED' AND expire_at + (?::int * INTERVAL '1 day') > CURRENT_TIMESTAMP THEN 'ACTIVE'
			        ELSE status
			    END,
			    updated_at = CURRENT_TIMESTAMP
			WHERE uuid = ANY(?)
		`, extendDays, extendDays, dedupeStrings(userUUIDs))
		if err != nil {
			_ = tx.Rollback()
			return err
		}
		affectedRows, _ = result.RowsAffected()

		return tx.Commit()
	})

	return affectedRows, nodeUUIDs, err
}

func extendAllUsersExpiration(ctx context.Context, manager *dbmanager.DatabaseManager, extendDays int) (int64, []string, error) {
	var affectedRows int64
	nodeUUIDs := make([]string, 0)

	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		tx, err := db.BeginTx(ctx, nil)
		if err != nil {
			return err
		}

		targets, nodeTargetsErr := queryAllReactivatedExpiredUserNodeUUIDsTx(ctx, tx, extendDays)
		if nodeTargetsErr != nil {
			_ = tx.Rollback()
			return nodeTargetsErr
		}
		nodeUUIDs = targets

		result, err := tx.ExecContext(ctx, `
			UPDATE users
			SET expire_at = expire_at + (?::int * INTERVAL '1 day'),
			    status = CASE
			        WHEN status = 'EXPIRED' AND expire_at + (?::int * INTERVAL '1 day') > CURRENT_TIMESTAMP THEN 'ACTIVE'
			        ELSE status
			    END,
			    updated_at = CURRENT_TIMESTAMP
		`, extendDays, extendDays)
		if err != nil {
			_ = tx.Rollback()
			return err
		}
		affectedRows, _ = result.RowsAffected()

		return tx.Commit()
	})

	return affectedRows, nodeUUIDs, err
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
	for name, value := range map[string]*string{
		"trojanPassword":    req.TrojanPassword,
		"ssPassword":        req.SSPassword,
		"naivePassword":     req.NaivePassword,
		"shadowtlsPassword": req.ShadowtlsPassword,
		"hysteria2Password": req.Hysteria2Password,
		"anytlsPassword":    req.AnytlsPassword,
	} {
		if value != nil && len(strings.TrimSpace(*value)) > 256 {
			return fmt.Errorf("%s must be less than 256 characters", name)
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
	if req.VlessUUID.Set {
		if req.VlessUUID.Value == nil || strings.TrimSpace(*req.VlessUUID.Value) == "" {
			return fmt.Errorf("vlessUuid cannot be empty")
		}
		if _, err := uuid.Parse(strings.TrimSpace(*req.VlessUUID.Value)); err != nil {
			return fmt.Errorf("invalid vlessUuid")
		}
	}
	for name, field := range map[string]OptionalString{
		"trojanPassword":    req.TrojanPassword,
		"ssPassword":        req.SSPassword,
		"naivePassword":     req.NaivePassword,
		"shadowtlsPassword": req.ShadowtlsPassword,
		"hysteria2Password": req.Hysteria2Password,
		"anytlsPassword":    req.AnytlsPassword,
	} {
		if !field.Set {
			continue
		}
		if field.Value == nil {
			if name == "trojanPassword" || name == "ssPassword" {
				return fmt.Errorf("%s cannot be empty", name)
			}
			continue
		}
		if (name == "trojanPassword" || name == "ssPassword") && strings.TrimSpace(*field.Value) == "" {
			return fmt.Errorf("%s cannot be empty", name)
		}
		if len(strings.TrimSpace(*field.Value)) > 256 {
			return fmt.Errorf("%s must be less than 256 characters", name)
		}
	}
	if req.ExternalSquadUUID.Set && req.ExternalSquadUUID.Value != nil && strings.TrimSpace(*req.ExternalSquadUUID.Value) != "" {
		if _, err := uuid.Parse(strings.TrimSpace(*req.ExternalSquadUUID.Value)); err != nil {
			return fmt.Errorf("invalid externalSquadUuid")
		}
	}
	return nil
}

func validateBulkUpdateUsersFields(fields bulkUpdateUsersFields) error {
	hasUpdate := fields.Status != nil ||
		fields.TrafficLimitBytes != nil ||
		fields.TrafficLimitStrategy != nil ||
		fields.ExpireAt != nil ||
		fields.Description.Set ||
		fields.Tag.Set ||
		fields.TelegramID.Set ||
		fields.Email.Set ||
		fields.HwidDeviceLimit.Set ||
		fields.ExternalSquadUUID.Set
	if !hasUpdate {
		return fmt.Errorf("at least one field must be provided")
	}
	if fields.Status != nil {
		status := strings.ToUpper(strings.TrimSpace(*fields.Status))
		if !isValidUserStatus(status) {
			return fmt.Errorf("invalid status")
		}
		if status == "EXPIRED" || status == "LIMITED" {
			return fmt.Errorf("invalid status")
		}
	}
	if fields.TrafficLimitBytes != nil && *fields.TrafficLimitBytes < 0 {
		return fmt.Errorf("trafficLimitBytes must be non-negative")
	}
	if fields.TrafficLimitStrategy != nil && !isValidTrafficStrategy(*fields.TrafficLimitStrategy) {
		return fmt.Errorf("invalid trafficLimitStrategy")
	}
	if fields.ExpireAt != nil {
		parsed, err := time.Parse(time.RFC3339, strings.TrimSpace(*fields.ExpireAt))
		if err != nil {
			return fmt.Errorf("expireAt must be RFC3339")
		}
		if !parsed.After(time.Now().UTC()) {
			return fmt.Errorf("expireAt must be in the future")
		}
	}
	if fields.Tag.Set && fields.Tag.Value != nil && strings.TrimSpace(*fields.Tag.Value) != "" {
		tag := strings.ToUpper(strings.TrimSpace(*fields.Tag.Value))
		if len(tag) > 16 || !userTagRegex.MatchString(tag) {
			return fmt.Errorf("tag can only contain uppercase letters, numbers, underscores and be up to 16 chars")
		}
	}
	if fields.Email.Set && fields.Email.Value != nil && strings.TrimSpace(*fields.Email.Value) != "" && !strings.Contains(strings.TrimSpace(*fields.Email.Value), "@") {
		return fmt.Errorf("invalid email")
	}
	if fields.HwidDeviceLimit.Set && fields.HwidDeviceLimit.Value != nil && *fields.HwidDeviceLimit.Value < 0 {
		return fmt.Errorf("hwidDeviceLimit must be non-negative")
	}
	if fields.TelegramID.Set && fields.TelegramID.Value != nil && *fields.TelegramID.Value < 0 {
		return fmt.Errorf("telegramId must be non-negative")
	}
	if fields.ExternalSquadUUID.Set && fields.ExternalSquadUUID.Value != nil && strings.TrimSpace(*fields.ExternalSquadUUID.Value) != "" {
		if _, err := uuid.Parse(strings.TrimSpace(*fields.ExternalSquadUUID.Value)); err != nil {
			return fmt.Errorf("invalid externalSquadUuid")
		}
	}
	return nil
}

func buildBulkUpdateUserClauses(fields bulkUpdateUsersFields) ([]string, []any) {
	clauses := make([]string, 0)
	args := make([]any, 0)
	add := func(column string, value any) {
		clauses = append(clauses, fmt.Sprintf("%s = ?", column))
		args = append(args, value)
	}

	if fields.Status != nil {
		add("status", strings.ToUpper(strings.TrimSpace(*fields.Status)))
	}
	if fields.TrafficLimitBytes != nil {
		add("traffic_limit_bytes", *fields.TrafficLimitBytes)
	}
	if fields.TrafficLimitStrategy != nil {
		add("traffic_limit_strategy", strings.ToUpper(strings.TrimSpace(*fields.TrafficLimitStrategy)))
	}
	if fields.ExpireAt != nil {
		parsed, _ := time.Parse(time.RFC3339, strings.TrimSpace(*fields.ExpireAt))
		add("expire_at", parsed.UTC())
	}
	if fields.Description.Set {
		if fields.Description.Value == nil || strings.TrimSpace(*fields.Description.Value) == "" {
			clauses = append(clauses, "description = NULL")
		} else {
			add("description", strings.TrimSpace(*fields.Description.Value))
		}
	}
	if fields.Tag.Set {
		if fields.Tag.Value == nil || strings.TrimSpace(*fields.Tag.Value) == "" {
			clauses = append(clauses, "tag = NULL")
		} else {
			add("tag", strings.ToUpper(strings.TrimSpace(*fields.Tag.Value)))
		}
	}
	if fields.TelegramID.Set {
		if fields.TelegramID.Value == nil {
			clauses = append(clauses, "telegram_id = NULL")
		} else {
			add("telegram_id", *fields.TelegramID.Value)
		}
	}
	if fields.Email.Set {
		if fields.Email.Value == nil || strings.TrimSpace(*fields.Email.Value) == "" {
			clauses = append(clauses, "email = NULL")
		} else {
			add("email", strings.TrimSpace(*fields.Email.Value))
		}
	}
	if fields.HwidDeviceLimit.Set {
		if fields.HwidDeviceLimit.Value == nil {
			clauses = append(clauses, "hwid_device_limit = NULL")
		} else {
			add("hwid_device_limit", *fields.HwidDeviceLimit.Value)
		}
	}
	if fields.ExternalSquadUUID.Set {
		if fields.ExternalSquadUUID.Value == nil || strings.TrimSpace(*fields.ExternalSquadUUID.Value) == "" {
			clauses = append(clauses, "external_squad_uuid = NULL")
		} else {
			add("external_squad_uuid", strings.TrimSpace(*fields.ExternalSquadUUID.Value))
		}
	}

	return clauses, args
}

func validateBulkAllUpdateUsersRequest(req bulkAllUpdateUsersRequest) error {
	hasUpdate := req.Status != nil ||
		req.TrafficLimitBytes != nil ||
		req.TrafficLimitStrategy != nil ||
		req.ExpireAt != nil ||
		req.Description.Set ||
		req.Tag.Set ||
		req.TelegramID.Set ||
		req.Email.Set ||
		req.HwidDeviceLimit.Set
	if !hasUpdate {
		return fmt.Errorf("at least one field must be provided")
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
		parsed, err := time.Parse(time.RFC3339, strings.TrimSpace(*req.ExpireAt))
		if err != nil {
			return fmt.Errorf("expireAt must be RFC3339")
		}
		if !parsed.After(time.Now().UTC()) {
			return fmt.Errorf("expireAt must be in the future")
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
	return nil
}

func parseUsersTableQuery(r *http.Request) (usersTableQuery, error) {
	values := r.URL.Query()
	query := usersTableQuery{
		Start:       0,
		Size:        25,
		Filters:     []usersTableFilter{},
		FilterModes: map[string]string{},
		Sorting:     []usersTableSorting{},
	}

	if raw := strings.TrimSpace(values.Get("start")); raw != "" {
		start, err := strconv.Atoi(raw)
		if err != nil {
			return query, err
		}
		if start > 0 {
			query.Start = start
		}
	}
	if raw := strings.TrimSpace(values.Get("size")); raw != "" {
		size, err := strconv.Atoi(raw)
		if err != nil {
			return query, err
		}
		if size > 0 {
			query.Size = size
		}
		if query.Size > 1000 {
			query.Size = 1000
		}
	}
	if raw := strings.TrimSpace(values.Get("filters")); raw != "" {
		if err := json.Unmarshal([]byte(raw), &query.Filters); err != nil {
			return query, err
		}
	}
	if raw := strings.TrimSpace(values.Get("filterModes")); raw != "" {
		if err := json.Unmarshal([]byte(raw), &query.FilterModes); err != nil {
			return query, err
		}
	}
	if raw := strings.TrimSpace(values.Get("sorting")); raw != "" {
		if err := json.Unmarshal([]byte(raw), &query.Sorting); err != nil {
			return query, err
		}
	}

	return query, nil
}

func filterUsersTableResponse(users []userAPI, filters []usersTableFilter, modes map[string]string) []userAPI {
	if len(filters) == 0 {
		return users
	}

	filtered := make([]userAPI, 0, len(users))
	for _, user := range users {
		matches := true
		for _, filter := range filters {
			if !matchesUsersTableFilter(user, filter, modes[filter.ID]) {
				matches = false
				break
			}
		}
		if matches {
			filtered = append(filtered, user)
		}
	}

	return filtered
}

func matchesUsersTableFilter(user userAPI, filter usersTableFilter, mode string) bool {
	if filter.Value == nil {
		return true
	}
	values := normalizeUsersTableFilterValues(filter.Value)
	if len(values) == 0 {
		return true
	}

	if filter.ID == "activeInternalSquads" {
		for _, value := range values {
			for _, squad := range user.ActiveInternalSquads {
				if strings.EqualFold(squad.UUID, value) || strings.EqualFold(squad.Name, value) {
					return true
				}
			}
		}
		return false
	}

	field := usersTableFieldValue(user, filter.ID)
	if field == nil {
		return false
	}

	if isNumericUsersTableFilterMode(mode) {
		return matchesUsersTableNumericFilter(field, values, mode)
	}

	fieldText := strings.ToLower(strings.TrimSpace(fmt.Sprint(field)))
	if fieldText == "" {
		return false
	}
	if mode == "" {
		mode = "contains"
	}

	for _, value := range values {
		needle := strings.ToLower(strings.TrimSpace(value))
		if needle == "" {
			continue
		}
		switch mode {
		case "equals", "equalsString":
			if fieldText == needle {
				return true
			}
		case "startsWith":
			if strings.HasPrefix(fieldText, needle) {
				return true
			}
		case "endsWith":
			if strings.HasSuffix(fieldText, needle) {
				return true
			}
		default:
			if strings.Contains(fieldText, needle) {
				return true
			}
		}
	}

	return false
}

func normalizeUsersTableFilterValues(value any) []string {
	switch typed := value.(type) {
	case []any:
		values := make([]string, 0, len(typed))
		for _, item := range typed {
			if item == nil {
				continue
			}
			values = append(values, strings.TrimSpace(fmt.Sprint(item)))
		}
		return values
	case []string:
		return typed
	case string:
		if strings.TrimSpace(typed) == "" {
			return nil
		}
		return []string{typed}
	case float64:
		return []string{strconv.FormatFloat(typed, 'f', -1, 64)}
	case int:
		return []string{strconv.Itoa(typed)}
	case int64:
		return []string{strconv.FormatInt(typed, 10)}
	default:
		text := strings.TrimSpace(fmt.Sprint(typed))
		if text == "" || text == "<nil>" {
			return nil
		}
		return []string{text}
	}
}

func usersTableFieldValue(user userAPI, columnID string) any {
	switch columnID {
	case "uuid":
		return user.UUID
	case "id":
		return user.ID
	case "shortUuid":
		return user.ShortUUID
	case "username":
		return user.Username
	case "status":
		return user.Status
	case "trafficLimitBytes":
		return user.TrafficLimitBytes
	case "trafficLimitStrategy":
		return user.TrafficLimitStrategy
	case "expireAt":
		return user.ExpireAt
	case "telegramId":
		if user.TelegramID == nil {
			return nil
		}
		return *user.TelegramID
	case "email":
		if user.Email == nil {
			return nil
		}
		return *user.Email
	case "description":
		if user.Description == nil {
			return nil
		}
		return *user.Description
	case "tag":
		if user.Tag == nil {
			return nil
		}
		return *user.Tag
	case "hwidDeviceLimit":
		if user.HwidDeviceLimit == nil {
			return nil
		}
		return *user.HwidDeviceLimit
	case "externalSquadUuid":
		if user.ExternalSquadUUID == nil {
			return nil
		}
		return *user.ExternalSquadUUID
	case "trojanPassword":
		return user.TrojanPassword
	case "vlessUuid":
		return user.VlessUUID
	case "ssPassword":
		return user.SSPassword
	case "naivePassword":
		return user.NaivePassword
	case "shadowtlsPassword":
		return user.ShadowtlsPassword
	case "hysteria2Password":
		return user.Hysteria2Password
	case "anytlsPassword":
		return user.AnytlsPassword
	case "createdAt":
		return user.CreatedAt
	case "updatedAt":
		return user.UpdatedAt
	case "subRevokedAt":
		return user.SubRevokedAt
	case "lastTrafficResetAt":
		return user.LastTrafficResetAt
	case "nodeName", "userTraffic.lastConnectedNodeUuid":
		if user.UserTraffic.LastConnectedNodeUUID == nil {
			return nil
		}
		return *user.UserTraffic.LastConnectedNodeUUID
	case "usedTrafficBytes", "userTraffic.usedTrafficBytes":
		return user.UserTraffic.UsedTrafficBytes
	case "userTraffic.lifetimeUsedTrafficBytes":
		return user.UserTraffic.LifetimeUsedTrafficBytes
	case "userTraffic.onlineAt":
		return user.UserTraffic.OnlineAt
	case "userTraffic.firstConnectedAt":
		return user.UserTraffic.FirstConnectedAt
	default:
		return nil
	}
}

func isNumericUsersTableFilterMode(mode string) bool {
	switch mode {
	case "equals", "greaterThan", "greaterThanOrEqualTo", "lessThan", "lessThanOrEqualTo", "between", "betweenInclusive":
		return true
	default:
		return false
	}
}

func matchesUsersTableNumericFilter(field any, values []string, mode string) bool {
	fieldValue, ok := usersTableFloatValue(field)
	if !ok {
		return false
	}
	parse := func(index int) (float64, bool) {
		if index >= len(values) || strings.TrimSpace(values[index]) == "" {
			return 0, false
		}
		value, err := strconv.ParseFloat(strings.TrimSpace(values[index]), 64)
		return value, err == nil
	}

	switch mode {
	case "greaterThan":
		value, ok := parse(0)
		return ok && fieldValue > value
	case "greaterThanOrEqualTo":
		value, ok := parse(0)
		return ok && fieldValue >= value
	case "lessThan":
		value, ok := parse(0)
		return ok && fieldValue < value
	case "lessThanOrEqualTo":
		value, ok := parse(0)
		return ok && fieldValue <= value
	case "between", "betweenInclusive":
		min, minOK := parse(0)
		max, maxOK := parse(1)
		if minOK && maxOK {
			return fieldValue >= min && fieldValue <= max
		}
		if minOK {
			return fieldValue >= min
		}
		if maxOK {
			return fieldValue <= max
		}
		return true
	default:
		value, ok := parse(0)
		return ok && fieldValue == value
	}
}

func usersTableFloatValue(value any) (float64, bool) {
	switch typed := value.(type) {
	case int:
		return float64(typed), true
	case int64:
		return float64(typed), true
	case float64:
		return typed, true
	case *int:
		if typed == nil {
			return 0, false
		}
		return float64(*typed), true
	case *int64:
		if typed == nil {
			return 0, false
		}
		return float64(*typed), true
	default:
		parsed, err := strconv.ParseFloat(strings.TrimSpace(fmt.Sprint(value)), 64)
		return parsed, err == nil
	}
}

func sortUsersTableResponse(users []userAPI, sorting []usersTableSorting) {
	if len(sorting) == 0 {
		return
	}

	sort.SliceStable(users, func(i, j int) bool {
		for _, sortRule := range sorting {
			comparison := compareUsersTableValues(usersTableFieldValue(users[i], sortRule.ID), usersTableFieldValue(users[j], sortRule.ID))
			if comparison == 0 {
				continue
			}
			if sortRule.Desc {
				return comparison > 0
			}
			return comparison < 0
		}
		return false
	})
}

func compareUsersTableValues(left any, right any) int {
	leftTime, leftIsTime := usersTableTimeValue(left)
	rightTime, rightIsTime := usersTableTimeValue(right)
	if leftIsTime || rightIsTime {
		if !leftIsTime && rightIsTime {
			return -1
		}
		if leftIsTime && !rightIsTime {
			return 1
		}
		if leftTime.Before(rightTime) {
			return -1
		}
		if leftTime.After(rightTime) {
			return 1
		}
		return 0
	}

	leftFloat, leftIsFloat := usersTableFloatValue(left)
	rightFloat, rightIsFloat := usersTableFloatValue(right)
	if leftIsFloat || rightIsFloat {
		if !leftIsFloat && rightIsFloat {
			return -1
		}
		if leftIsFloat && !rightIsFloat {
			return 1
		}
		if leftFloat < rightFloat {
			return -1
		}
		if leftFloat > rightFloat {
			return 1
		}
		return 0
	}

	leftText := strings.ToLower(strings.TrimSpace(fmt.Sprint(left)))
	rightText := strings.ToLower(strings.TrimSpace(fmt.Sprint(right)))
	if leftText < rightText {
		return -1
	}
	if leftText > rightText {
		return 1
	}
	return 0
}

func usersTableTimeValue(value any) (time.Time, bool) {
	switch typed := value.(type) {
	case time.Time:
		return typed, true
	case *time.Time:
		if typed == nil {
			return time.Time{}, false
		}
		return *typed, true
	default:
		return time.Time{}, false
	}
}

func paginateUsersTableResponse(users []userAPI, start int, size int) []userAPI {
	if start < 0 {
		start = 0
	}
	if size <= 0 {
		size = 25
	}
	if start >= len(users) {
		return []userAPI{}
	}
	end := start + size
	if end > len(users) {
		end = len(users)
	}
	return users[start:end]
}

func getAllUserRecords(ctx context.Context, manager *dbmanager.DatabaseManager) ([]userRecord, error) {
	records := make([]userRecord, 0)
	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		rows, err := db.QueryContext(ctx, `
			SELECT
				u.t_id, u.uuid, u.short_uuid, u.username, u.status, u.traffic_limit_bytes,
				u.traffic_limit_strategy, u.expire_at, u.last_traffic_reset_at,
					u.sub_revoked_at, u.trojan_password, u.vless_uuid, u.ss_password,
					u.naive_password, u.shadowtls_password, u.hysteria2_password, u.anytls_password,
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
				u.traffic_limit_strategy, u.expire_at, u.last_traffic_reset_at,
					u.sub_revoked_at, u.trojan_password, u.vless_uuid, u.ss_password,
					u.naive_password, u.shadowtls_password, u.hysteria2_password, u.anytls_password,
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

func getUserRecordsByUUIDs(ctx context.Context, manager *dbmanager.DatabaseManager, userUUIDs []string) (map[string]userRecord, error) {
	clean := dedupeStrings(userUUIDs)
	records := make(map[string]userRecord, len(clean))
	if len(clean) == 0 {
		return records, nil
	}

	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		rows, err := db.QueryContext(ctx, `
			SELECT
				u.t_id, u.uuid, u.short_uuid, u.username, u.status, u.traffic_limit_bytes,
				u.traffic_limit_strategy, u.expire_at, u.last_traffic_reset_at,
					u.sub_revoked_at, u.trojan_password, u.vless_uuid, u.ss_password,
					u.naive_password, u.shadowtls_password, u.hysteria2_password, u.anytls_password,
					u.description, u.tag, u.telegram_id, u.email, u.hwid_device_limit, u.external_squad_uuid,
				u.last_triggered_threshold, u.created_at, u.updated_at,
				COALESCE(ut.used_traffic_bytes, 0), COALESCE(ut.lifetime_used_traffic_bytes, 0),
				ut.online_at, ut.last_connected_node_uuid, ut.first_connected_at
			FROM users u
			LEFT JOIN user_traffic ut ON ut.t_id = u.t_id
			WHERE u.uuid = ANY(?)
		`, clean)
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			record, scanErr := scanUserRecord(rows)
			if scanErr != nil {
				return scanErr
			}
			records[record.UUID] = record
		}
		return rows.Err()
	})
	return records, err
}

func scanUserRecord(scanner shared.RowScanner) (userRecord, error) {
	var record userRecord
	var (
		lastTrafficReset sql.NullTime
		subRevokedAt     sql.NullTime
		naivePassword    sql.NullString
		shadowtlsPass    sql.NullString
		hysteria2Pass    sql.NullString
		anytlsPass       sql.NullString
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
		&lastTrafficReset,
		&subRevokedAt,
		&record.TrojanPassword,
		&record.VlessUUID,
		&record.SSPassword,
		&naivePassword,
		&shadowtlsPass,
		&hysteria2Pass,
		&anytlsPass,
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

	if lastTrafficReset.Valid {
		record.LastTrafficResetAt = &lastTrafficReset.Time
	}
	if subRevokedAt.Valid {
		record.SubRevokedAt = &subRevokedAt.Time
	}
	if naivePassword.Valid {
		record.NaivePassword = &naivePassword.String
	}
	if shadowtlsPass.Valid {
		record.ShadowtlsPassword = &shadowtlsPass.String
	}
	if hysteria2Pass.Valid {
		record.Hysteria2Password = &hysteria2Pass.String
	}
	if anytlsPass.Valid {
		record.AnytlsPassword = &anytlsPass.String
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
			NaivePassword:          protocolCredentialString(record.NaivePassword, record.TrojanPassword),
			ShadowtlsPassword:      protocolCredentialString(record.ShadowtlsPassword, record.SSPassword),
			Hysteria2Password:      protocolCredentialString(record.Hysteria2Password, record.TrojanPassword),
			AnytlsPassword:         protocolCredentialString(record.AnytlsPassword, record.TrojanPassword),
			LastTriggeredThreshold: record.LastTriggeredThreshold,
			SubRevokedAt:           record.SubRevokedAt,
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

func emitUserNotification(ctx context.Context, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig, event string, record userRecord, meta map[string]any) {
	data := userRecordNotificationData(record)
	if userNotificationNeedsInternalSquads(event) {
		enrichUserNotificationInternalSquads(ctx, manager, record.UUID, data)
	}

	notifications.Emit(ctx, cfg, notifications.Event{
		Scope: notifications.ScopeUser,
		Event: event,
		Data:  data,
		Meta:  meta,
	})
}

func userNotificationNeedsInternalSquads(event string) bool {
	switch event {
	case notifications.EventUserCreated, notifications.EventUserModified, notifications.EventUserRevoked:
		return true
	default:
		return false
	}
}

func enrichUserNotificationInternalSquads(ctx context.Context, manager *dbmanager.DatabaseManager, userUUID string, data map[string]any) {
	if manager == nil || strings.TrimSpace(userUUID) == "" || data == nil {
		return
	}

	squadsByUser, err := getUsersActiveInternalSquads(ctx, manager, []string{userUUID})
	if err != nil {
		return
	}

	squads := squadsByUser[userUUID]
	names := make([]string, 0, len(squads))
	for _, squad := range squads {
		name := strings.TrimSpace(squad.Name)
		if name != "" {
			names = append(names, name)
		}
	}

	data["activeInternalSquads"] = names
	data["internalSquads"] = names
}

func emitUsersByUUIDsNotification(ctx context.Context, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig, event string, userUUIDs []string) {
	clean := dedupeStrings(userUUIDs)
	if len(clean) == 0 {
		return
	}
	records, err := getUserRecordsByUUIDs(ctx, manager, clean)
	if err != nil {
		emitUsersNotificationFromRecords(ctx, manager, cfg, event, clean, nil)
		return
	}
	emitUsersNotificationFromRecords(ctx, manager, cfg, event, clean, records)
}

func emitUsersNotificationFromRecords(ctx context.Context, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig, event string, userUUIDs []string, records map[string]userRecord) {
	clean := dedupeStrings(userUUIDs)
	if len(clean) == 0 {
		return
	}
	skipTelegram := len(clean) >= 500
	for _, userUUID := range clean {
		meta := map[string]any{"bulk": true}
		if skipTelegram {
			meta["skipTelegramNotification"] = true
		}
		if record, ok := records[userUUID]; ok {
			emitUserNotification(ctx, manager, cfg, event, record, meta)
			continue
		}
		notifications.Emit(ctx, cfg, notifications.Event{
			Scope: notifications.ScopeUser,
			Event: event,
			Data:  map[string]any{"uuid": userUUID},
			Meta:  meta,
		})
	}
}

func emitBulkSummaryNotification(ctx context.Context, cfg *config.BackendConfig, event string, affectedRows int64) {
	if affectedRows <= 0 {
		return
	}
	notifications.Emit(ctx, cfg, notifications.Event{
		Scope: notifications.ScopeUser,
		Event: event,
		Data: map[string]any{
			"affectedRows": affectedRows,
		},
		Meta: map[string]any{
			"bulk":                     true,
			"skipTelegramNotification": affectedRows >= 500,
		},
	})
}

func userRecordNotificationData(record userRecord) map[string]any {
	return map[string]any{
		"tId":                    record.TID,
		"uuid":                   record.UUID,
		"shortUuid":              record.ShortUUID,
		"username":               record.Username,
		"status":                 record.Status,
		"trafficLimitBytes":      record.TrafficLimitBytes,
		"trafficLimitStrategy":   record.TrafficLimitStrategy,
		"expireAt":               record.ExpireAt.UTC().Format(time.RFC3339),
		"telegramId":             record.TelegramID,
		"email":                  record.Email,
		"description":            record.Description,
		"tag":                    record.Tag,
		"hwidDeviceLimit":        record.HwidDeviceLimit,
		"externalSquadUuid":      record.ExternalSquadUUID,
		"trojanPassword":         record.TrojanPassword,
		"vlessUuid":              record.VlessUUID,
		"ssPassword":             record.SSPassword,
		"naivePassword":          protocolCredentialString(record.NaivePassword, record.TrojanPassword),
		"shadowtlsPassword":      protocolCredentialString(record.ShadowtlsPassword, record.SSPassword),
		"hysteria2Password":      protocolCredentialString(record.Hysteria2Password, record.TrojanPassword),
		"anytlsPassword":         protocolCredentialString(record.AnytlsPassword, record.TrojanPassword),
		"lastTriggeredThreshold": record.LastTriggeredThreshold,
		"subRevokedAt":           optionalTimeString(record.SubRevokedAt),
		"lastTrafficResetAt":     optionalTimeString(record.LastTrafficResetAt),
		"createdAt":              record.CreatedAt.UTC().Format(time.RFC3339),
		"updatedAt":              record.UpdatedAt.UTC().Format(time.RFC3339),
		"userTraffic": map[string]any{
			"usedTrafficBytes":         record.UsedTrafficBytes,
			"lifetimeUsedTrafficBytes": record.LifetimeUsedTrafficBytes,
			"onlineAt":                 optionalTimeString(record.OnlineAt),
			"firstConnectedAt":         optionalTimeString(record.FirstConnectedAt),
			"lastConnectedNodeUuid":    record.LastConnectedNodeUUID,
		},
	}
}

func optionalTimeString(value *time.Time) *string {
	if value == nil {
		return nil
	}
	formatted := value.UTC().Format(time.RFC3339)
	return &formatted
}

func userStatusChangedNotification(previous, next string) string {
	if strings.EqualFold(previous, next) {
		return ""
	}
	switch strings.ToUpper(strings.TrimSpace(next)) {
	case "ACTIVE":
		return notifications.EventUserEnabled
	case "DISABLED":
		return notifications.EventUserDisabled
	case "LIMITED":
		return notifications.EventUserLimited
	case "EXPIRED":
		return notifications.EventUserExpired
	default:
		return ""
	}
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

func resolveUsersSubscriptionBase(ctx context.Context, manager *dbmanager.DatabaseManager, r *http.Request, cfg *config.BackendConfig) string {
	if base := resolveUsersSubscriptionBaseFromNode(ctx, manager); base != "" {
		return base
	}

	return resolveUsersSubscriptionBaseFallback(r, cfg)
}

func resolveUsersSubscriptionBaseFromNode(ctx context.Context, manager *dbmanager.DatabaseManager) string {
	if manager == nil {
		return ""
	}

	var domain sql.NullString
	var apiPath sql.NullString
	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		row := db.QueryRowContext(ctx, `
			SELECT
				COALESCE(NULLIF(BTRIM(public_domain), ''), NULLIF(BTRIM(address), '')) AS domain,
				COALESCE(NULLIF(BTRIM(api_path), ''), '/') AS api_path
			FROM sub_nodes
			ORDER BY is_disabled ASC, view_position ASC, created_at ASC
			LIMIT 1
		`)

		scanErr := row.Scan(&domain, &apiPath)
		if errors.Is(scanErr, sql.ErrNoRows) {
			return nil
		}

		return scanErr
	})
	if err != nil || !domain.Valid {
		return ""
	}

	nodeDomain := strings.TrimSpace(strings.Split(domain.String, ",")[0])
	if nodeDomain == "" {
		return ""
	}

	if !strings.Contains(nodeDomain, "://") {
		nodeDomain = "https://" + nodeDomain
	}

	parsedDomain, parseErr := url.Parse(nodeDomain)
	if parseErr != nil || strings.TrimSpace(parsedDomain.Host) == "" {
		return ""
	}

	parsedDomain.Path = ""
	parsedDomain.RawQuery = ""
	parsedDomain.Fragment = ""
	parsedDomain.User = nil

	base := strings.TrimRight(parsedDomain.String(), "/")
	path := normalizeUsersSubscriptionAPIPath(apiPath.String)
	return base + path
}

func normalizeUsersSubscriptionAPIPath(value string) string {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" || trimmed == "/" {
		return "/"
	}

	return "/" + strings.Trim(trimmed, "/") + "/"
}

func resolveUsersSubscriptionBaseFallback(r *http.Request, cfg *config.BackendConfig) string {
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

func protocolCredentialString(value *string, fallback string) string {
	if value != nil && strings.TrimSpace(*value) != "" {
		return strings.TrimSpace(*value)
	}
	return strings.TrimSpace(fallback)
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
