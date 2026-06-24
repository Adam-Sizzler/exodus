package users

import (
	"context"

	dbmanager "exodus/internal/db/manager"
)

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
