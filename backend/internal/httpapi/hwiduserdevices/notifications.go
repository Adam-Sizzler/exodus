package hwiduserdevices

import (
	"context"

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

func hwidCompatNotificationData(device hwidCompatDevice) map[string]any {
	return map[string]any{
		"hwid":        device.HWID,
		"userId":      device.UserID,
		"platform":    device.Platform,
		"osVersion":   device.OSVersion,
		"deviceModel": device.DeviceModel,
		"userAgent":   device.UserAgent,
		"createdAt":   device.CreatedAt,
		"updatedAt":   device.UpdatedAt,
	}
}
