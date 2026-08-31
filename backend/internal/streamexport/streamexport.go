package streamexport

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	StreamKeyUserUsage            = "ioraw:export:user_usage"
	StreamKeySubscriptionRequests = "ioraw:export:subscription_requests"

	UserUsageStreamMessageVersion           = "1"
	SubscriptionRequestStreamMessageVersion = "1"
)

type UserUsageEntry struct {
	UserID     int64
	TotalBytes int64
}

type SubscriptionRequestExport struct {
	UserID          int64
	RequestAt       time.Time
	SSRResponseType string
	RequestIP       string
	UserAgent       string
	SRRRuleName     *string
}

// ExportUserUsageBatch exports a batch of user traffic deltas to Redis stream ioraw:export:user_usage
func ExportUserUsageBatch(ctx context.Context, client *redis.Client, enabled bool, maxLen int, nodeID int64, entries []UserUsageEntry) error {
	if !enabled || client == nil || len(entries) == 0 {
		return nil
	}
	if maxLen <= 0 {
		maxLen = 3000
	}

	parts := make([]string, 0, len(entries))
	for _, entry := range entries {
		if entry.UserID > 0 && entry.TotalBytes > 0 {
			parts = append(parts, fmt.Sprintf("%d:%d", entry.UserID, entry.TotalBytes))
		}
	}
	if len(parts) == 0 {
		return nil
	}

	values := map[string]interface{}{
		"v":       UserUsageStreamMessageVersion,
		"nodeId":  strconv.FormatInt(nodeID, 10),
		"ts":      time.Now().UTC().Format(time.RFC3339Nano),
		"records": strings.Join(parts, ";"),
	}

	return client.XAdd(ctx, &redis.XAddArgs{
		Stream: StreamKeyUserUsage,
		MaxLen: int64(maxLen),
		Approx: true,
		Values: values,
	}).Err()
}

// ExportSubscriptionRequest exports a subscription request event to Redis stream ioraw:export:subscription_requests
func ExportSubscriptionRequest(ctx context.Context, client *redis.Client, enabled bool, maxLen int, req SubscriptionRequestExport) error {
	if !enabled || client == nil || req.UserID <= 0 {
		return nil
	}
	if maxLen <= 0 {
		maxLen = 3000
	}

	reqAt := req.RequestAt
	if reqAt.IsZero() {
		reqAt = time.Now().UTC()
	}

	ssrType := req.SSRResponseType
	if ssrType == "" {
		ssrType = "UNKNOWN"
	}

	values := map[string]interface{}{
		"v":               SubscriptionRequestStreamMessageVersion,
		"userId":          strconv.FormatInt(req.UserID, 10),
		"requestAt":       reqAt.Format(time.RFC3339Nano),
		"ssrResponseType": ssrType,
	}

	if req.RequestIP != "" {
		values["requestIp"] = req.RequestIP
	}
	if req.UserAgent != "" {
		values["userAgent"] = req.UserAgent
	}
	if req.SRRRuleName != nil && *req.SRRRuleName != "" {
		values["srrRuleName"] = *req.SRRRuleName
	}

	return client.XAdd(ctx, &redis.XAddArgs{
		Stream: StreamKeySubscriptionRequests,
		MaxLen: int64(maxLen),
		Approx: true,
		Values: values,
	}).Err()
}
