package users

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"

	"exodus/internal/notifications"
	"exodus/internal/proto"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// receiveStream receives and processes stream data.
func (nm *NodeMonitor) receiveStream(state *nodeState) {
	for {
		resp, err := state.stream.Recv()
		if err == io.EOF {
			nm.cfg.Logger.Warn("Stream closed by node", "node", state.nodeName)
			nm.handleDisconnect(state, "Stream closed")
			return
		}
		if err != nil {
			if st, ok := status.FromError(err); ok {
				if st.Code() == codes.Canceled && state.ctx.Err() != nil {
					nm.cfg.Logger.Debug("Stream canceled", "node", state.nodeName)
					return
				}
				if st.Code() == codes.Unavailable {
					reason := "Node unavailable"
					if strings.TrimSpace(st.Message()) != "" {
						reason = fmt.Sprintf("Node unavailable: %s", st.Message())
					}
					nm.cfg.Logger.Warn(
						"Node unavailable",
						"node", state.nodeName,
						"code", st.Code().String(),
						"message", st.Message(),
					)
					nm.handleDisconnect(state, reason)
					return
				}
			}
			nm.cfg.Logger.Error("Stream error", "node", state.nodeName, "error", err)
			nm.handleDisconnect(state, fmt.Sprintf("Stream error: %v", err))
			return
		}

		nm.processResponse(state.nodeName, resp)
	}
}

// handleDisconnect handles node disconnection.
func (nm *NodeMonitor) handleDisconnect(state *nodeState, reason string) {
	state.mutex.Lock()
	wasConnected := state.isConnected
	state.isConnected = false
	state.isConnecting = false
	state.lastError = reason
	if state.stream != nil {
		state.stream = nil
	}
	if state.conn != nil {
		state.conn.Close()
		state.conn = nil
	}
	state.mutex.Unlock()

	if wasConnected {
		nm.updateConnectionStatus(state.nodeName, false, false, reason)
		nm.cfg.Logger.Warn(fmt.Sprintf("Lost connection to Node %s (%s:%d), message: %s", state.nodeName, state.address, state.port, reason))
	} else {
		nm.cfg.Logger.Warn(fmt.Sprintf("Connection attempt failed for Node %s (%s:%d), message: %s", state.nodeName, state.address, state.port, reason))
	}
}

// processResponse processes node response data.
func (nm *NodeMonitor) processResponse(nodeName string, resp *proto.NodeDataResponse) {
	switch payload := resp.Response.(type) {
	case *proto.NodeDataResponse_Stats:
		nm.cfg.Logger.Trace("Node stats received", "node", nodeName)
		nm.updateNodeRuntimeFromStats(nodeName, payload.Stats.GetStats())
	case *proto.NodeDataResponse_Users:
		nm.cfg.Logger.Trace("Node users received", "node", nodeName)
	case *proto.NodeDataResponse_LogData:
		nm.cfg.Logger.Trace("Node log data received", "node", nodeName)
	default:
		nm.cfg.Logger.Trace("Node message received", "node", nodeName)
	}
}

func (nm *NodeMonitor) updateNodeRuntimeFromStats(nodeName string, stats []*proto.Stat) {
	if len(stats) == 0 {
		return
	}

	values := make(map[string]string, len(stats))
	for _, stat := range stats {
		if stat == nil {
			continue
		}
		key := strings.ToLower(strings.TrimSpace(stat.GetName()))
		if key == "" {
			continue
		}
		values[key] = strings.TrimSpace(stat.GetValue())
	}

	trafficDelta := extractTrafficStatsDelta(stats)

	coreStatus := strings.ToLower(strings.TrimSpace(values["core_status"]))
	coreError := strings.TrimSpace(values["core_error"])
	switch coreStatus {
	case "running", "ok", "healthy":
		nm.updateConnectionStatus(nodeName, true, false, "")
	case "error", "failed", "unhealthy", "stopped":
		message := "Core error"
		if coreError != "" {
			message = fmt.Sprintf("Core error: %s", coreError)
		}
		nm.updateConnectionStatus(nodeName, false, false, message)
	}

	singboxVersion := firstNonEmptyString(values["singbox_version"])
	nodeVersion := firstNonEmptyString(values["node_version"])
	singboxUptime, hasSingboxUptime := parseOptionalUptimeSeconds(values["singbox_uptime"])
	systemInfo := parseOptionalJSONRaw(values["system_info"])
	systemStats := parseOptionalJSONRaw(values["system_stats"])
	usersOnline := trafficDelta.UsersOnline

	persistedNodeUUID := ""
	firstConnectedEvents := make([]notifications.Event, 0)

	var (
		nodeUUID                  string
		nodeID                    int64
		consumptionMultiplier     int64
		nodeConsumptionMultiplier int64
	)
	if err := nm.db.QueryRow(`SELECT uuid, id, consumption_multiplier, node_consumption_multiplier FROM nodes WHERE name = $1`, nodeName).Scan(&nodeUUID, &nodeID, &consumptionMultiplier, &nodeConsumptionMultiplier); err != nil {
		nm.cfg.Logger.Warn("Failed to query node from DB", "node", nodeName, "error", err)
		return
	}
	persistedNodeUUID = nodeUUID

	if _, execErr := nm.db.Exec(`
		UPDATE nodes
		SET updated_at = CURRENT_TIMESTAMP
		WHERE name = $1`, nodeName); execErr != nil {
		nm.cfg.Logger.Warn("Failed to update node updated_at", "node", nodeName, "error", execErr)
		return
	}

	if trafficDelta.TotalUploadBytes > 0 || trafficDelta.TotalDownloadBytes > 0 {
		totalBytes := trafficDelta.TotalUploadBytes + trafficDelta.TotalDownloadBytes
		nodeUsageBytes := applyConsumptionMultiplier(totalBytes, nodeConsumptionMultiplier)
		if _, execErr := nm.db.Exec(`
			INSERT INTO nodes_usage_history (node_uuid, download_bytes, upload_bytes, total_bytes)
			VALUES ($1, $2, $3, $4)
			ON CONFLICT (node_uuid, created_at)
			DO UPDATE SET
				download_bytes = nodes_usage_history.download_bytes + EXCLUDED.download_bytes,
				upload_bytes = nodes_usage_history.upload_bytes + EXCLUDED.upload_bytes,
				total_bytes = nodes_usage_history.total_bytes + EXCLUDED.total_bytes,
				updated_at = now()
		`, nodeUUID, trafficDelta.TotalDownloadBytes, trafficDelta.TotalUploadBytes, totalBytes); execErr != nil {
			nm.cfg.Logger.Warn("Failed to insert node usage history", "node", nodeName, "error", execErr)
			return
		}

		if _, execErr := nm.db.Exec(`
			UPDATE nodes
			SET traffic_used_bytes = COALESCE(traffic_used_bytes, 0) + $1, updated_at = CURRENT_TIMESTAMP
			WHERE uuid = $2
		`, nodeUsageBytes, nodeUUID); execErr != nil {
			nm.cfg.Logger.Warn("Failed to update node traffic used bytes", "node", nodeName, "error", execErr)
			return
		}
	}

	if len(trafficDelta.UserBytesByName) > 0 {
		usernames := make([]string, 0, len(trafficDelta.UserBytesByName))
		for username := range trafficDelta.UserBytesByName {
			if strings.TrimSpace(username) != "" {
				usernames = append(usernames, username)
			}
		}
		if len(usernames) > 0 {
			rows, queryErr := nm.db.Query(`
				SELECT t_id, username
				FROM users
				WHERE status = 'ACTIVE' AND username = ANY($1)
			`, usernames)
			if queryErr == nil {
				userIDs := make(map[string]int64, len(usernames))
				for rows.Next() {
					var (
						userID   int64
						username string
					)
					if scanErr := rows.Scan(&userID, &username); scanErr == nil {
						userIDs[strings.TrimSpace(username)] = userID
					}
				}
				rows.Close()

				firstConnectedByID := make(map[int64]bool, len(userIDs))
				if len(userIDs) > 0 {
					ids := make([]int64, 0, len(userIDs))
					for _, userID := range userIDs {
						ids = append(ids, userID)
					}
					firstRows, firstErr := nm.db.Query(`
						SELECT t_id, first_connected_at IS NOT NULL
						FROM user_traffic
						WHERE t_id = ANY($1)
					`, ids)
					if firstErr == nil {
						for firstRows.Next() {
							var (
								userID            int64
								hasFirstConnected bool
							)
							if scanErr := firstRows.Scan(&userID, &hasFirstConnected); scanErr == nil {
								firstConnectedByID[userID] = hasFirstConnected
							}
						}
						firstRows.Close()
					}
				}

				usageDeltas := make([]userUsageDelta, 0, len(trafficDelta.UserBytesByName))
				for username, rawBytes := range trafficDelta.UserBytesByName {
					username = strings.TrimSpace(username)
					userID, ok := userIDs[username]
					if !ok {
						continue
					}
					if nm.cfg != nil && nm.cfg.Redis.UserUsageIgnoreBelowBytes > 0 && rawBytes < nm.cfg.Redis.UserUsageIgnoreBelowBytes {
						continue
					}

					effectiveBytes := applyConsumptionMultiplier(rawBytes, consumptionMultiplier)
					if effectiveBytes <= 0 {
						continue
					}

					usageDeltas = append(usageDeltas, userUsageDelta{
						UserID:       userID,
						Username:     username,
						TotalBytes:   effectiveBytes,
						HistoryBytes: rawBytes,
					})

					if !firstConnectedByID[userID] {
						firstConnectedEvents = append(firstConnectedEvents, notifications.Event{
							Scope: notifications.ScopeUser,
							Event: notifications.EventUserFirstConnected,
							Data: map[string]any{
								"tId":      userID,
								"username": username,
								"nodeUuid": nodeUUID,
								"nodeName": nodeName,
							},
						})
						firstConnectedByID[userID] = true
					}
				}

				if len(usageDeltas) > 0 {
					bulkCtx := nm.globalCtx
					if bulkCtx == nil {
						bulkCtx = context.Background()
					}

					if execErr := bulkUpsertUserTraffic(bulkCtx, nm.db, usageDeltas, nodeUUID); execErr != nil {
						nm.cfg.Logger.Warn("Failed to upsert user traffic", "node", nodeName, "error", execErr)
					}
					if recordErr := nm.recordNodeUserUsageHistory(bulkCtx, nm.db, nodeID, usageDeltas); recordErr != nil {
						nm.cfg.Logger.Warn("Failed to record node user usage history", "node", nodeName, "error", recordErr)
					}
				}
			}
		}
	}

	for _, event := range firstConnectedEvents {
		notifications.Emit(context.Background(), nm.cfg, event)
	}

	nm.updateNodeMetricsSnapshot(persistedNodeUUID, usersOnline, trafficDelta)
	nm.updateHotCacheNodeRuntime(nodeName, persistedNodeUUID, singboxVersion, nodeVersion, hasSingboxUptime, singboxUptime, usersOnline, systemInfo, systemStats, trafficDelta)
}

func (nm *NodeMonitor) updateHotCacheNodeRuntime(
	nodeName string,
	nodeUUID string,
	singboxVersion string,
	nodeVersion string,
	hasSingboxUptime bool,
	singboxUptime int64,
	usersOnline int,
	systemInfo json.RawMessage,
	systemStats json.RawMessage,
	delta trafficStatsDelta,
) {
	if nm.hotCache == nil {
		return
	}
	ctx := nm.globalCtx
	if ctx == nil {
		ctx = context.Background()
	}

	if nm.hotCache != nil && strings.TrimSpace(nodeUUID) != "" {
		if systemInfo != nil && json.Valid(systemInfo) {
			_ = nm.hotCache.SetSystemInfo(ctx, nodeUUID, systemInfo)
		}
		if systemStats != nil && json.Valid(systemStats) {
			_ = nm.hotCache.SetSystemStats(ctx, nodeUUID, systemStats)
		}
		if singboxVersion != "" || nodeVersion != "" {
			_ = nm.hotCache.SetVersions(ctx, nodeUUID, singboxVersion, nodeVersion)
		}
		if hasSingboxUptime {
			_ = nm.hotCache.SetUptime(ctx, nodeUUID, singboxUptime)
		}
		if usersOnline >= 0 {
			_ = nm.hotCache.SetUsersOnline(ctx, nodeUUID, usersOnline)
		}
	}
}

func applyConsumptionMultiplier(bytes int64, multiplier int64) int64 {
	if bytes <= 0 || multiplier <= 0 {
		return 0
	}
	return (bytes * multiplier) / 1_000_000_000
}

func parseOptionalUptimeSeconds(raw string) (int64, bool) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return 0, false
	}
	val, err := strconv.ParseInt(trimmed, 10, 64)
	if err != nil || val < 0 {
		return 0, false
	}
	return val, true
}

func parseOptionalJSONRaw(raw string) json.RawMessage {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return nil
	}
	b := []byte(trimmed)
	if !json.Valid(b) {
		return nil
	}
	return json.RawMessage(b)
}

func extractTrafficStatsDelta(stats []*proto.Stat) trafficStatsDelta {
	var delta trafficStatsDelta
	delta.InboundByTag = make(map[string]TagTrafficCounters)
	delta.OutboundByTag = make(map[string]TagTrafficCounters)
	delta.UserBytesByName = make(map[string]int64)

	for _, stat := range stats {
		if stat == nil {
			continue
		}
		key := strings.ToLower(strings.TrimSpace(stat.GetName()))
		valStr := strings.TrimSpace(stat.GetValue())
		val, _ := strconv.ParseInt(valStr, 10, 64)

		if strings.Contains(key, ">>>") {
			parts := strings.Split(key, ">>>")
			if len(parts) == 4 && parts[2] == "traffic" {
				category := parts[0]
				tagOrUser := parts[1]
				direction := parts[3]
				switch category {
				case "outbound":
					if direction == "uplink" {
						delta.TotalUploadBytes += val
					} else if direction == "downlink" {
						delta.TotalDownloadBytes += val
					}
					counters := delta.OutboundByTag[tagOrUser]
					if direction == "uplink" {
						counters.UploadBytes += val
					} else if direction == "downlink" {
						counters.DownloadBytes += val
					}
					delta.OutboundByTag[tagOrUser] = counters
				case "inbound":
					counters := delta.InboundByTag[tagOrUser]
					if direction == "uplink" {
						counters.UploadBytes += val
					} else if direction == "downlink" {
						counters.DownloadBytes += val
					}
					delta.InboundByTag[tagOrUser] = counters
				case "user":
					if val > 0 {
						delta.UserBytesByName[tagOrUser] += val
					}
				}
			}
			continue
		}

		switch {
		case key == "total_download_bytes":
			delta.TotalDownloadBytes = val
		case key == "total_upload_bytes":
			delta.TotalUploadBytes = val
		case key == "users_online":
			if val >= 0 && val <= math.MaxInt {
				delta.UsersOnline = int(val)
			}
		case strings.HasPrefix(key, "user_"):
			username := strings.TrimPrefix(key, "user_")
			username = strings.TrimSuffix(username, "_bytes")
			username = strings.TrimSpace(username)
			if username != "" && val > 0 {
				delta.UserBytesByName[username] += val
			}
		case strings.HasPrefix(key, "inbound_"):
			parts := strings.Split(key, "_")
			if len(parts) >= 3 {
				direction := parts[len(parts)-1]
				tag := strings.Join(parts[1:len(parts)-1], "_")
				counters := delta.InboundByTag[tag]
				if direction == "down" || direction == "download" {
					counters.DownloadBytes += val
				} else if direction == "up" || direction == "upload" {
					counters.UploadBytes += val
				}
				delta.InboundByTag[tag] = counters
			}
		case strings.HasPrefix(key, "outbound_"):
			parts := strings.Split(key, "_")
			if len(parts) >= 3 {
				direction := parts[len(parts)-1]
				tag := strings.Join(parts[1:len(parts)-1], "_")
				counters := delta.OutboundByTag[tag]
				if direction == "down" || direction == "download" {
					counters.DownloadBytes += val
				} else if direction == "up" || direction == "upload" {
					counters.UploadBytes += val
				}
				delta.OutboundByTag[tag] = counters
			}
		}
	}
	if delta.UsersOnline == 0 {
		delta.UsersOnline = len(delta.UserBytesByName)
	}
	return delta
}
