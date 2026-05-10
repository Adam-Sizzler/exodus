package hwiduserdevices

import (
	"context"
	"time"

	"exodus/internal/config"
	"exodus/internal/notifications"
)

func emitHWIDNotification(ctx context.Context, cfg *config.BackendConfig, event string, data map[string]any, meta map[string]any) {
	notifications.Emit(ctx, cfg, notifications.Event{
		Scope: notifications.ScopeUserHWIDDevices,
		Event: event,
		Data:  data,
		Meta:  meta,
	})
}

func hwidRecordNotificationData(record HWIDUserDeviceRecord) map[string]any {
	return map[string]any{
		"uuid":        record.UUID,
		"userTId":     record.UserTID,
		"hwid":        record.HWID,
		"deviceName":  record.DeviceName,
		"firstSeenAt": record.FirstSeenAt.UTC().Format(time.RFC3339),
		"lastSeenAt":  record.LastSeenAt.UTC().Format(time.RFC3339),
		"createdAt":   record.CreatedAt.UTC().Format(time.RFC3339),
		"updatedAt":   record.UpdatedAt.UTC().Format(time.RFC3339),
	}
}

func hwidCompatNotificationData(device hwidCompatDevice) map[string]any {
	return map[string]any{
		"hwid":        device.HWID,
		"userUuid":    device.UserUUID,
		"platform":    device.Platform,
		"osVersion":   device.OSVersion,
		"deviceModel": device.DeviceModel,
		"userAgent":   device.UserAgent,
		"createdAt":   device.CreatedAt,
		"updatedAt":   device.UpdatedAt,
	}
}
