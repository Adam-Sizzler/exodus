package users

import (
	"context"
	"strings"
	"time"

	"exodus/internal/config"
	"exodus/internal/notifications"
)

func emitUserNotification(ctx context.Context, repo *UserRepository, cfg *config.BackendConfig, event string, record userRecord, meta map[string]any) {
	data := userRecordNotificationData(record)
	if userNotificationNeedsInternalSquads(event) {
		enrichUserNotificationInternalSquads(ctx, repo, record.UUID, data)
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

func enrichUserNotificationInternalSquads(ctx context.Context, repo *UserRepository, userUUID string, data map[string]any) {
	if repo == nil || strings.TrimSpace(userUUID) == "" || data == nil {
		return
	}

	squadsByUser, err := repo.getUsersActiveInternalSquads(ctx, []string{userUUID})
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

func emitUsersByUUIDsNotification(ctx context.Context, repo *UserRepository, cfg *config.BackendConfig, event string, userUUIDs []string) {
	clean := dedupeStrings(userUUIDs)
	if len(clean) == 0 {
		return
	}
	records, err := repo.getUserRecordsByUUIDs(ctx, clean)
	if err != nil {
		emitUsersNotificationFromRecords(ctx, repo, cfg, event, clean, nil)
		return
	}
	emitUsersNotificationFromRecords(ctx, repo, cfg, event, clean, records)
}

func emitUsersNotificationFromRecords(ctx context.Context, repo *UserRepository, cfg *config.BackendConfig, event string, userUUIDs []string, records map[string]userRecord) {
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
			emitUserNotification(ctx, repo, cfg, event, record, meta)
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
		"naivePassword":          protocolCredentialString(record.NaivePassword, ""),
		"shadowtlsPassword":      protocolCredentialString(record.ShadowtlsPassword, ""),
		"hysteria2Password":      protocolCredentialString(record.Hysteria2Password, ""),
		"anytlsPassword":         protocolCredentialString(record.AnytlsPassword, ""),
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
