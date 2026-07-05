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

	errUserNotFound        = errors.New("user not found")
	errUsernameExists      = errors.New("username already exists")
	errShortUUIDExists     = errors.New("short uuid already exists")
	errVLESSUUIDExists     = errors.New("vless uuid already exists")
	errExternalSquadAbsent = errors.New("external squad not found")
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

type resolveUserRequest struct {
	UUID      *string `json:"uuid,omitempty"`
	ID        *int64  `json:"id,omitempty"`
	ShortUUID *string `json:"shortUuid,omitempty"`
	Username  *string `json:"username,omitempty"`
}

type resolveUserResponse struct {
	UUID      string `json:"uuid"`
	ID        int64  `json:"id"`
	ShortUUID string `json:"shortUuid"`
	Username  string `json:"username"`
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
