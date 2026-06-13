package nodes

import (
	"context"
	"encoding/json"
	"strings"

	dbmanager "exodus/internal/db/manager"
)

func buildNodeResponses(ctx context.Context, manager *dbmanager.DatabaseManager, records []nodeRecord) ([]nodeAPI, error) {
	nodeUUIDs := make([]string, 0, len(records))
	providerUUIDs := make([]string, 0, len(records))
	for _, record := range records {
		nodeUUIDs = append(nodeUUIDs, record.UUID)
		if record.ProviderUUID != nil && *record.ProviderUUID != "" {
			providerUUIDs = append(providerUUIDs, *record.ProviderUUID)
		}
	}

	inboundsMap, err := getNodeInbounds(ctx, manager, nodeUUIDs)
	if err != nil {
		return nil, err
	}
	providersMap, err := getProviders(ctx, manager, dedupeStrings(providerUUIDs))
	if err != nil {
		return nil, err
	}

	response := make([]nodeAPI, 0, len(records))
	for _, record := range records {
		var item nodeAPI
		item.UUID = record.UUID
		item.Name = record.Name
		item.Address = record.Address
		item.Port = record.Port
		item.APISchema = normalizeAPISchema(&record.APISchema)
		item.APIPath = normalizeAPIPath(&record.APIPath)
		item.GRPCAuthToken = record.GRPCAuthToken
		item.ActivePluginUUID = record.ActivePluginUUID
		item.IsConnected = record.IsConnected
		item.IsDisabled = record.IsDisabled
		item.IsConnecting = record.IsConnecting
		item.LastStatusChange = record.LastStatusChange
		item.LastStatusMessage = record.LastStatusMessage
		item.SingboxVersion = record.SingboxVersion
		item.NodeVersion = record.NodeVersion
		item.SingboxUptime = record.SingboxUptime
		item.IsTrafficTrackingActive = record.IsTrafficTrackingActive
		item.TrafficResetDay = record.TrafficResetDay
		item.TrafficLimitBytes = record.TrafficLimitBytes
		item.TrafficUsedBytes = record.TrafficUsedBytes
		item.NotifyPercent = record.NotifyPercent
		item.UsersOnline = record.UsersOnline
		item.ViewPosition = record.ViewPosition
		item.CountryCode = record.CountryCode
		item.ConsumptionMultiplier = fromNanoMultiplier(record.ConsumptionMultiplier)
		item.Tags = ensureStringSlice(record.Tags)
		item.CPUCount = record.CPUCount
		item.CPUModel = record.CPUModel
		item.TotalRAM = record.TotalRAM
		item.System = buildNodeSystem(record.SystemInfoRaw, record.SystemStatsRaw)
		item.Versions = buildNodeVersions(record.SingboxVersion, record.NodeVersion)
		if !record.IsConnected || record.IsConnecting || record.IsDisabled {
			usersOnline := 0
			item.SingboxUptime = "0"
			item.UsersOnline = &usersOnline
			item.CPUCount = nil
			item.CPUModel = nil
			item.TotalRAM = nil
			item.System = nil
		}
		item.CreatedAt = record.CreatedAt
		item.UpdatedAt = record.UpdatedAt
		item.ConfigProfile.ActiveConfigProfileUUID = record.ActiveConfigProfileUUID
		item.ConfigProfile.ActiveInbounds = ensureInboundSlice(inboundsMap[record.UUID])
		item.ProviderUUID = record.ProviderUUID
		if record.ProviderUUID != nil {
			item.Provider = providersMap[*record.ProviderUUID]
		}
		response = append(response, item)
	}

	return response, nil
}

func buildNodeSystem(infoRaw []byte, statsRaw []byte) *nodeSystemResponse {
	if len(infoRaw) == 0 || len(statsRaw) == 0 {
		return nil
	}

	var info nodeSystemInfoResponse
	if err := json.Unmarshal(infoRaw, &info); err != nil {
		return nil
	}
	var stats nodeSystemStatsResponse
	if err := json.Unmarshal(statsRaw, &stats); err != nil {
		return nil
	}
	if stats.LoadAvg == nil {
		stats.LoadAvg = []float64{0, 0, 0}
	}
	if info.NetworkInterfaces == nil {
		info.NetworkInterfaces = []string{}
	}

	return &nodeSystemResponse{
		Info:  info,
		Stats: stats,
	}
}

func buildNodeVersions(singboxVersion *string, nodeVersion *string) *nodeVersionsResponse {
	if singboxVersion == nil && nodeVersion == nil {
		return nil
	}
	versions := &nodeVersionsResponse{
		Singbox: "unknown",
		Node:    "unknown",
	}
	if singboxVersion != nil && strings.TrimSpace(*singboxVersion) != "" {
		versions.Singbox = strings.TrimSpace(*singboxVersion)
	}
	if nodeVersion != nil && strings.TrimSpace(*nodeVersion) != "" {
		versions.Node = strings.TrimSpace(*nodeVersion)
	}
	return versions
}
