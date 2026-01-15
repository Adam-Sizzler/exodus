package common

import (
	"time"
	"v2ray-stat/logger"
)

var TimeLocation *time.Location

func InitTimezone(tz string, logger *logger.Logger) {
	loc, _ := time.LoadLocation(tz)
	logger.Info("Timezone set successfully", "timezone", tz)
	TimeLocation = loc
}

func GetLocalUnix() int64 {
	now := time.Now()
	if TimeLocation != nil {
		now = now.In(TimeLocation)
	}
	return now.Unix()
}
