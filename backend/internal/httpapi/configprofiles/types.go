package configprofiles

import (
	"encoding/json"
	"errors"
	"time"
)

var (
	errConfigProfileNotFound        = errors.New("config profile not found")
	errConfigProfileSnippetNotFound = errors.New("config profile snippet not found")
)

const (
	errMessageConfigProfileNameAlreadyExists = "Config profile name already exists in database. Config profile names must be unique."
	errMessageInboundTagsMustBeUnique        = "Inbounds with same tag already exists in database. Inbound tags must be unique."
)

type ConfigProfileInbound struct {
	UUID         string          `json:"uuid"`
	ProfileUUID  string          `json:"profileUuid"`
	Tag          string          `json:"tag"`
	Type         string          `json:"type"`
	Network      *string         `json:"network"`
	Security     *string         `json:"security"`
	Port         *int            `json:"port"`
	RawInbound   json.RawMessage `json:"rawInbound"`
	ActiveSquads []string        `json:"activeSquads"`
}

type ConfigProfileNode struct {
	UUID        string `json:"uuid"`
	Name        string `json:"name"`
	CountryCode string `json:"countryCode"`
}

type ConfigProfile struct {
	UUID         string                 `json:"uuid"`
	ViewPosition int                    `json:"viewPosition"`
	Name         string                 `json:"name"`
	Config       json.RawMessage        `json:"config"`
	Inbounds     []ConfigProfileInbound `json:"inbounds"`
	Nodes        []ConfigProfileNode    `json:"nodes"`
	CreatedAt    time.Time              `json:"createdAt"`
	UpdatedAt    time.Time              `json:"updatedAt"`
}

type configProfileRecord struct {
	UUID         string
	ViewPosition int
	Name         string
	Config       json.RawMessage
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type createConfigProfileRequest struct {
	Name   string          `json:"name"`
	Config json.RawMessage `json:"config"`
}

type updateConfigProfileRequest struct {
	UUID   string           `json:"uuid"`
	Name   *string          `json:"name,omitempty"`
	Config *json.RawMessage `json:"config,omitempty"`
}

type reorderConfigProfilesItem struct {
	UUID         string `json:"uuid"`
	ViewPosition int    `json:"viewPosition"`
}

type reorderConfigProfilesRequest struct {
	Items []reorderConfigProfilesItem `json:"items"`
}

type ConfigProfileSnippet struct {
	Name      string          `json:"name"`
	Snippet   json.RawMessage `json:"snippet"`
	CreatedAt time.Time       `json:"createdAt"`
}

type configProfileSnippetRequest struct {
	Name    string          `json:"name"`
	Snippet json.RawMessage `json:"snippet"`
}
