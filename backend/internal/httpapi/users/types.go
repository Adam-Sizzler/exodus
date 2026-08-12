package users

import (
	"errors"
	"regexp"
	"time"

	"exodus/internal/httpapi/shared"
)

var (
	userUsernameRegex = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
	userTagRegex      = regexp.MustCompile(`^[A-Z0-9_]+$`)

	errUserNotFound    = errors.New("user not found")
	errUsernameExists  = errors.New("username already exists")
	errShortUUIDExists = errors.New("short uuid already exists")
	errVLESSUUIDExists = errors.New("vless uuid already exists")
)

type OptionalString = shared.OptionalString

type OptionalInt = shared.OptionalInt

type OptionalInt64 = shared.OptionalInt64

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

// UserResponseEnvelope wraps single user response.
type UserResponseEnvelope struct {
	Response userAPI `json:"response"`
}

// UsersListResponse wraps paginated user list.
type UsersListResponse struct {
	Users []userAPI `json:"users"`
	Total int       `json:"total"`
}

// UsersListResponseEnvelope wraps UsersListResponse.
type UsersListResponseEnvelope struct {
	Response UsersListResponse `json:"response"`
}

// UsersStreamResponse wraps streaming users response.
type UsersStreamResponse struct {
	Users      []userAPI `json:"users"`
	Total      int       `json:"total"`
	NextCursor *int64    `json:"nextCursor"`
}

// UsersStreamResponseEnvelope wraps UsersStreamResponse.
type UsersStreamResponseEnvelope struct {
	Response UsersStreamResponse `json:"response"`
}

// UserTagsResponseEnvelope wraps user tags response.
type UserTagsResponseEnvelope struct {
	Response struct {
		Tags []string `json:"tags"`
	} `json:"response"`
}

// UserAccessibleNodesResponse wraps accessible nodes for a user.
type UserAccessibleNodesResponse struct {
	UserID      int64                `json:"userId"`
	ActiveNodes []userAccessibleNode `json:"activeNodes"`
}

// UserAccessibleNodesResponseEnvelope wraps UserAccessibleNodesResponse.
type UserAccessibleNodesResponseEnvelope struct {
	Response UserAccessibleNodesResponse `json:"response"`
}

// UserSubscriptionRequestHistoryResponse wraps request history for a user.
type UserSubscriptionRequestHistoryResponse struct {
	Records []userSubscriptionRequestHistoryRecord `json:"records"`
	Total   int                                    `json:"total"`
}

// UserSubscriptionRequestHistoryResponseEnvelope wraps UserSubscriptionRequestHistoryResponse.
type UserSubscriptionRequestHistoryResponseEnvelope struct {
	Response UserSubscriptionRequestHistoryResponse `json:"response"`
}

type userSubscriptionRequestHistoryRecord struct {
	ID              int64   `json:"id"`
	UserID          int64   `json:"userId"`
	SRRResponseType string  `json:"srrResponseType"`
	SRRRuleName     *string `json:"srrRuleName"`
	RequestIP       *string `json:"requestIp"`
	UserAgent       *string `json:"userAgent"`
	RequestAt       string  `json:"requestAt"`
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
	ID                   *int64         `json:"id,omitempty"`
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

type resolveUserRequest struct {
	ID        *int64  `json:"id,omitempty"`
	ShortUUID *string `json:"shortUuid,omitempty"`
	Username  *string `json:"username,omitempty"`
}

type resolveUserResponse struct {
	ID        int64  `json:"id"`
	ShortUUID string `json:"shortUuid"`
	Username  string `json:"username"`
}

type ResolveUserResponseEnvelope struct {
	Response resolveUserResponse `json:"response"`
}

type bulkDeleteUsersRequest struct {
	UserIDs []int64 `json:"userIds,omitempty"`
}

type bulkRevokeUsersSubscriptionRequest struct {
	UserIDs []int64 `json:"userIds,omitempty"`
}

type extendUserExpirationRequest struct {
	Days       int `json:"days,omitempty"`
	ExtendDays int `json:"extendDays,omitempty"`
}

type bulkDeleteUsersByStatusRequest struct {
	Status string `json:"status"`
}

type bulkExtendExpirationDateRequest struct {
	UserIDs    []int64 `json:"userIds,omitempty"`
	ExtendDays int     `json:"extendDays"`
}

type bulkUpdateUsersRequest struct {
	UserIDs []int64               `json:"userIds,omitempty"`
	Fields  bulkUpdateUsersFields `json:"fields"`
}

type bulkUpdateUsersSquadsRequest struct {
	UserIDs              []int64  `json:"userIds,omitempty"`
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
