package nodes

import (
	"context"
	"time"

	"exodus/internal/config"
	dbmanager "exodus/internal/db/manager"
	"exodus/internal/nodehotcache"
	"exodus/internal/notifications"
)

func emitNodeNotification(ctx context.Context, cfg *config.BackendConfig, event string, record nodeRecord, meta map[string]any) {
	notifications.Emit(ctx, cfg, notifications.Event{
		Scope: notifications.ScopeNode,
		Event: event,
		Data:  nodeRecordNotificationData(ctx, cfg, record),
		Meta:  meta,
	})
}

func emitNodesByUUIDsNotification(ctx context.Context, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig, event string, nodeUUIDs []string) {
	clean := dedupeStrings(nodeUUIDs)
	if len(clean) == 0 {
		return
	}
	skipTelegram := len(clean) >= 500
	for _, nodeUUID := range clean {
		meta := map[string]any{"bulk": true}
		if skipTelegram {
			meta["skipTelegramNotification"] = true
		}
		record, err := getNodeByUUID(ctx, manager, nodeUUID)
		if err == nil {
			emitNodeNotification(ctx, cfg, event, record, meta)
			continue
		}
		notifications.Emit(ctx, cfg, notifications.Event{
			Scope: notifications.ScopeNode,
			Event: event,
			Data:  map[string]any{"uuid": nodeUUID},
			Meta:  meta,
		})
	}
}

func nodeRecordNotificationData(ctx context.Context, cfg *config.BackendConfig, record nodeRecord) map[string]any {
	hot, _ := nodehotcache.Default(cfg).GetOne(ctx, record.UUID)
	var singboxVersion *string
	var nodeVersion *string
	if hot.Versions != nil {
		singboxVersion = stringPtrIfNotEmpty(hot.Versions.Singbox)
		nodeVersion = stringPtrIfNotEmpty(hot.Versions.Node)
	}
	return map[string]any{
		"id":                      record.ID,
		"uuid":                    record.UUID,
		"name":                    record.Name,
		"address":                 record.Address,
		"port":                    record.Port,
		"apiSchema":               record.APISchema,
		"apiPath":                 record.APIPath,
		"activeConfigProfileUuid": record.ActiveConfigProfileUUID,
		"activePluginUuid":        record.ActivePluginUUID,
		"isConnected":             record.IsConnected,
		"isConnecting":            record.IsConnecting,
		"isDisabled":              record.IsDisabled,
		"lastStatusChange":        optionalTimeString(record.LastStatusChange),
		"lastStatusMessage":       record.LastStatusMessage,
		"singboxVersion":          singboxVersion,
		"nodeVersion":             nodeVersion,
		"singboxUptime":           hot.SingboxUptime,
		"usersOnline":             hot.UsersOnline,
		"consumptionMultiplier":   fromNanoMultiplier(record.ConsumptionMultiplier),
		"isTrafficTrackingActive": record.IsTrafficTrackingActive,
		"trafficResetDay":         record.TrafficResetDay,
		"trafficLimitBytes":       record.TrafficLimitBytes,
		"trafficUsedBytes":        record.TrafficUsedBytes,
		"notifyPercent":           record.NotifyPercent,
		"providerUuid":            record.ProviderUUID,
		"viewPosition":            record.ViewPosition,
		"countryCode":             record.CountryCode,
		"tags":                    ensureStringSlice(record.Tags),
		"createdAt":               record.CreatedAt.UTC().Format(time.RFC3339),
		"updatedAt":               record.UpdatedAt.UTC().Format(time.RFC3339),
	}
}

func optionalTimeString(value *time.Time) *string {
	if value == nil {
		return nil
	}
	formatted := value.UTC().Format(time.RFC3339)
	return &formatted
}
