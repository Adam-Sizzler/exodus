package common

import (
	"time"
	"v2ray-stat/logger"
)

var TimeLocation *time.Location

func InitTimezone(tz string, logger *logger.Logger) {
	loc, err := time.LoadLocation(tz)
	if err != nil {
		logger.Error("Failed to load timezone", "tz", tz, "error", err)
		return // fallback to UTC
	}
	logger.Info("Timezone set successfully", "tz", tz)
	TimeLocation = loc
}

func GetLocalUnix() int64 {
	now := time.Now()
	if TimeLocation == nil {
		return now.Unix()
	}
	_, offset := now.In(TimeLocation).Zone()
	return now.Unix() + int64(offset)
}
