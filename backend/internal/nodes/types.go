package users

import (
	"context"
	"time"
)

type deployRequest struct {
	Restart      bool
	ForceRestart bool
	NodeUUIDs    []string
}

type TagTrafficCounters struct {
	UploadBytes   int64
	DownloadBytes int64
}

type userUsageDelta struct {
	UserID       int64
	Username     string
	TotalBytes   int64
	HistoryBytes int64
}

type trafficStatsDelta struct {
	UsersOnline        int
	TotalUploadBytes   int64
	TotalDownloadBytes int64
	InboundByTag       map[string]TagTrafficCounters
	OutboundByTag      map[string]TagTrafficCounters
	UserBytesByName    map[string]int64
}

type NodeUserUsageRecorder interface {
	RecordNodeUserUsage(ctx context.Context, nodeID int64, userBytes map[int64]int64) error
}

type NodeMetricsSnapshot struct {
	NodeUUID    string
	UsersOnline int
	Inbounds    map[string]TagTrafficCounters
	Outbounds   map[string]TagTrafficCounters
	UpdatedAt   time.Time
}
