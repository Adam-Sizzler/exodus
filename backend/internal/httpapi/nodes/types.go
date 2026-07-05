package nodes

import (
	"encoding/json"
	"errors"
	"regexp"
	"time"

	"exodus/internal/httpapi/shared"
)

var nodeTagRegex = regexp.MustCompile(`^[A-Z0-9_:]+$`)

var (
	errConfigProfileNotFound       = errors.New("config profile not found")
	errConfigProfileInboundInvalid = errors.New("config profile inbound not found in specified profile")
	errNoEnabledNodes              = errors.New("enabled nodes not found")
)

type OptionalString = shared.OptionalString

type configProfileInboundResponse struct {
	UUID        string          `json:"uuid"`
	ProfileUUID string          `json:"profileUuid"`
	Tag         string          `json:"tag"`
	Type        string          `json:"type"`
	Network     *string         `json:"network"`
	Security    *string         `json:"security"`
	Port        *int            `json:"port"`
	RawInbound  json.RawMessage `json:"rawInbound"`
}

type providerResponse struct {
	UUID        string     `json:"uuid"`
	Name        string     `json:"name"`
	FaviconLink *string    `json:"faviconLink"`
	LoginURL    *string    `json:"loginUrl"`
	CreatedAt   *time.Time `json:"createdAt,omitempty"`
	UpdatedAt   *time.Time `json:"updatedAt,omitempty"`
}

type nodeSystemInfoResponse struct {
	Arch              string   `json:"arch"`
	CPUs              int      `json:"cpus"`
	CPUModel          string   `json:"cpuModel"`
	MemoryTotal       uint64   `json:"memoryTotal"`
	Hostname          string   `json:"hostname"`
	Platform          string   `json:"platform"`
	Release           string   `json:"release"`
	Type              string   `json:"type"`
	Version           string   `json:"version"`
	NetworkInterfaces []string `json:"networkInterfaces"`
}

type nodeSystemInterfaceResponse struct {
	Interface     string  `json:"interface"`
	RXBytesPerSec float64 `json:"rxBytesPerSec"`
	TXBytesPerSec float64 `json:"txBytesPerSec"`
	RXTotal       uint64  `json:"rxTotal"`
	TXTotal       uint64  `json:"txTotal"`
}

type nodeSystemStatsResponse struct {
	MemoryFree uint64                       `json:"memoryFree"`
	MemoryUsed uint64                       `json:"memoryUsed"`
	Uptime     float64                      `json:"uptime"`
	LoadAvg    []float64                    `json:"loadAvg"`
	Interface  *nodeSystemInterfaceResponse `json:"interface"`
}

type nodeSystemResponse struct {
	Info  nodeSystemInfoResponse  `json:"info"`
	Stats nodeSystemStatsResponse `json:"stats"`
}

type nodeVersionsResponse struct {
	Singbox string `json:"singbox"`
	Node    string `json:"node"`
}

type nodeAPI struct {
	UUID                    string                `json:"uuid"`
	Name                    string                `json:"name"`
	Address                 string                `json:"address"`
	Port                    *int                  `json:"port"`
	APISchema               string                `json:"apiSchema"`
	APIPath                 string                `json:"apiPath"`
	GRPCAuthToken           string                `json:"grpcAuthToken"`
	ActivePluginUUID        *string               `json:"activePluginUuid"`
	IsConnected             bool                  `json:"isConnected"`
	IsDisabled              bool                  `json:"isDisabled"`
	IsConnecting            bool                  `json:"isConnecting"`
	LastStatusChange        *time.Time            `json:"lastStatusChange"`
	LastStatusMessage       *string               `json:"lastStatusMessage"`
	SingboxVersion          *string               `json:"singboxVersion"`
	NodeVersion             *string               `json:"nodeVersion"`
	SingboxUptime           int64                 `json:"singboxUptime"`
	IsTrafficTrackingActive bool                  `json:"isTrafficTrackingActive"`
	TrafficResetDay         *int                  `json:"trafficResetDay"`
	TrafficLimitBytes       *int64                `json:"trafficLimitBytes"`
	TrafficUsedBytes        *int64                `json:"trafficUsedBytes"`
	NotifyPercent           *int                  `json:"notifyPercent"`
	UsersOnline             *int                  `json:"usersOnline"`
	ViewPosition            int                   `json:"viewPosition"`
	CountryCode             string                `json:"countryCode"`
	ConsumptionMultiplier   float64               `json:"consumptionMultiplier"`
	Tags                    []string              `json:"tags"`
	CPUCount                *int                  `json:"cpuCount"`
	CPUModel                *string               `json:"cpuModel"`
	TotalRAM                *string               `json:"totalRam"`
	System                  *nodeSystemResponse   `json:"system"`
	Versions                *nodeVersionsResponse `json:"versions"`
	CreatedAt               time.Time             `json:"createdAt"`
	UpdatedAt               time.Time             `json:"updatedAt"`
	ConfigProfile           struct {
		ActiveConfigProfileUUID *string                        `json:"activeConfigProfileUuid"`
		ActiveInbounds          []configProfileInboundResponse `json:"activeInbounds"`
	} `json:"configProfile"`
	ProviderUUID *string           `json:"providerUuid"`
	Provider     *providerResponse `json:"provider"`
}

type nodeRecord struct {
	ID                      *int64
	UUID                    string
	Name                    string
	Address                 string
	Port                    *int
	APISchema               string
	APIPath                 string
	GRPCAuthToken           string
	ActiveConfigProfileUUID *string
	ActivePluginUUID        *string
	IsConnected             bool
	IsConnecting            bool
	IsDisabled              bool
	LastStatusChange        *time.Time
	LastStatusMessage       *string
	ConsumptionMultiplier   int64
	IsTrafficTrackingActive bool
	TrafficResetDay         *int
	TrafficLimitBytes       *int64
	TrafficUsedBytes        *int64
	NotifyPercent           *int
	ProviderUUID            *string
	ViewPosition            int
	CountryCode             string
	Tags                    []string
	CreatedAt               time.Time
	UpdatedAt               time.Time
}

type configProfileRefRequest struct {
	ActiveConfigProfileUUID string   `json:"activeConfigProfileUuid"`
	ActiveInbounds          []string `json:"activeInbounds"`
}

type createNodeRequest struct {
	Name                    string                  `json:"name"`
	Address                 string                  `json:"address"`
	Port                    *int                    `json:"port,omitempty"`
	APISchema               *string                 `json:"apiSchema,omitempty"`
	APIPath                 *string                 `json:"apiPath,omitempty"`
	GRPCAuthToken           *string                 `json:"grpcAuthToken,omitempty"`
	ActivePluginUUID        *string                 `json:"activePluginUuid,omitempty"`
	IsTrafficTrackingActive *bool                   `json:"isTrafficTrackingActive,omitempty"`
	TrafficLimitBytes       *int64                  `json:"trafficLimitBytes,omitempty"`
	NotifyPercent           *int                    `json:"notifyPercent,omitempty"`
	TrafficResetDay         *int                    `json:"trafficResetDay,omitempty"`
	CountryCode             *string                 `json:"countryCode,omitempty"`
	ConsumptionMultiplier   *float64                `json:"consumptionMultiplier,omitempty"`
	ConfigProfile           configProfileRefRequest `json:"configProfile"`
	ProviderUUID            *string                 `json:"providerUuid,omitempty"`
	Tags                    []string                `json:"tags,omitempty"`
}

type updateNodeRequest struct {
	UUID                    string                   `json:"uuid"`
	Name                    *string                  `json:"name,omitempty"`
	Address                 *string                  `json:"address,omitempty"`
	Port                    *int                     `json:"port,omitempty"`
	APISchema               *string                  `json:"apiSchema,omitempty"`
	APIPath                 *string                  `json:"apiPath,omitempty"`
	GRPCAuthToken           *string                  `json:"grpcAuthToken,omitempty"`
	ActivePluginUUID        OptionalString           `json:"activePluginUuid,omitempty"`
	IsTrafficTrackingActive *bool                    `json:"isTrafficTrackingActive,omitempty"`
	TrafficLimitBytes       *int64                   `json:"trafficLimitBytes,omitempty"`
	NotifyPercent           *int                     `json:"notifyPercent,omitempty"`
	TrafficResetDay         *int                     `json:"trafficResetDay,omitempty"`
	CountryCode             *string                  `json:"countryCode,omitempty"`
	ConsumptionMultiplier   *float64                 `json:"consumptionMultiplier,omitempty"`
	ConfigProfile           *configProfileRefRequest `json:"configProfile,omitempty"`
	ProviderUUID            OptionalString           `json:"providerUuid,omitempty"`
	Tags                    *[]string                `json:"tags,omitempty"`
}

type reorderNodesRequest struct {
	Nodes []struct {
		UUID         string `json:"uuid"`
		ViewPosition int    `json:"viewPosition"`
	} `json:"nodes"`
}

type restartAllNodesRequest struct {
	ForceRestart *bool `json:"forceRestart,omitempty"`
}

type bulkNodesActionsRequest struct {
	UUIDs  []string `json:"uuids"`
	Action string   `json:"action"`
}

type bulkUpdateNodeFieldsRequest struct {
	CountryCode           *string        `json:"countryCode,omitempty"`
	ConsumptionMultiplier *float64       `json:"consumptionMultiplier,omitempty"`
	ProviderUUID          OptionalString `json:"providerUuid,omitempty"`
	Tags                  *[]string      `json:"tags,omitempty"`
	ActivePluginUUID      OptionalString `json:"activePluginUuid,omitempty"`
}

type bulkUpdateNodesRequest struct {
	UUIDs  []string                    `json:"uuids"`
	Fields bulkUpdateNodeFieldsRequest `json:"fields"`
}

type bulkProfileModificationRequest struct {
	UUIDs         []string                `json:"uuids"`
	ConfigProfile configProfileRefRequest `json:"configProfile"`
}
