package users

import (
	"context"
	"encoding/json"
	dbmanager "exodus/internal/db/manager"
	"exodus/internal/notifications"
	"exodus/internal/proto"
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"

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
	err := nm.manager.ExecuteLowPriority(func(db dbmanager.DBExecutor) error {
		var (
			nodeUUID                  string
			nodeID                    int64
			consumptionMultiplier     int64
			nodeConsumptionMultiplier int64
		)
		if err := db.QueryRow(`SELECT uuid, id, consumption_multiplier, node_consumption_multiplier FROM nodes WHERE name = ?`, nodeName).Scan(&nodeUUID, &nodeID, &consumptionMultiplier, &nodeConsumptionMultiplier); err != nil {
			return err
		}
		persistedNodeUUID = nodeUUID

		if _, execErr := db.Exec(`
			UPDATE nodes
			SET updated_at = CURRENT_TIMESTAMP
			WHERE name = ?`, nodeName); execErr != nil {
			return execErr
		}

		if trafficDelta.TotalUploadBytes > 0 || trafficDelta.TotalDownloadBytes > 0 {
			totalBytes := trafficDelta.TotalUploadBytes + trafficDelta.TotalDownloadBytes
			nodeUsageBytes := applyConsumptionMultiplier(totalBytes, nodeConsumptionMultiplier)
			if _, execErr := db.Exec(`
				INSERT INTO nodes_usage_history (node_uuid, download_bytes, upload_bytes, total_bytes)
				VALUES (?, ?, ?, ?)
				ON CONFLICT (node_uuid, created_at)
				DO UPDATE SET
					download_bytes = nodes_usage_history.download_bytes + EXCLUDED.download_bytes,
					upload_bytes = nodes_usage_history.upload_bytes + EXCLUDED.upload_bytes,
					total_bytes = nodes_usage_history.total_bytes + EXCLUDED.total_bytes,
					updated_at = now()
			`, nodeUUID, trafficDelta.TotalDownloadBytes, trafficDelta.TotalUploadBytes, totalBytes); execErr != nil {
				return execErr
			}

			if _, execErr := db.Exec(`
				UPDATE nodes
				SET traffic_used_bytes = COALESCE(traffic_used_bytes, 0) + ?, updated_at = CURRENT_TIMESTAMP
				WHERE uuid = ?
			`, nodeUsageBytes, nodeUUID); execErr != nil {
				return execErr
			}
		}

		if len(trafficDelta.UserBytesByName) == 0 {
			return nil
		}

		usernames := make([]string, 0, len(trafficDelta.UserBytesByName))
		for username := range trafficDelta.UserBytesByName {
			if strings.TrimSpace(username) != "" {
				usernames = append(usernames, username)
			}
		}
		if len(usernames) == 0 {
			return nil
		}

		rows, queryErr := db.Query(`
			SELECT t_id, username
			FROM users
			WHERE status = 'ACTIVE' AND username = ANY(?)
		`, usernames)
		if queryErr != nil {
			return queryErr
		}
		defer rows.Close()

		userIDs := make(map[string]int64, len(usernames))
		for rows.Next() {
			var (
				userID   int64
				username string
			)
			if scanErr := rows.Scan(&userID, &username); scanErr != nil {
				return scanErr
			}
			userIDs[strings.TrimSpace(username)] = userID
		}
		if rowsErr := rows.Err(); rowsErr != nil {
			return rowsErr
		}

		firstConnectedByID := make(map[int64]bool, len(userIDs))
		if len(userIDs) > 0 {
			ids := make([]int64, 0, len(userIDs))
			for _, userID := range userIDs {
				ids = append(ids, userID)
			}
			firstRows, firstErr := db.Query(`
				SELECT t_id, first_connected_at IS NOT NULL
				FROM user_traffic
				WHERE t_id = ANY(?)
			`, ids)
			if firstErr != nil {
				return firstErr
			}
			defer firstRows.Close()
			for firstRows.Next() {
				var (
					userID            int64
					hasFirstConnected bool
				)
				if scanErr := firstRows.Scan(&userID, &hasFirstConnected); scanErr != nil {
					return scanErr
				}
				firstConnectedByID[userID] = hasFirstConnected
			}
			if rowsErr := firstRows.Err(); rowsErr != nil {
				return rowsErr
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

		if len(usageDeltas) == 0 {
			return nil
		}

		bulkCtx := nm.globalCtx
		if bulkCtx == nil {
			bulkCtx = context.Background()
		}

		if execErr := bulkUpsertUserTraffic(bulkCtx, db, usageDeltas, nodeUUID); execErr != nil {
			return execErr
		}

		if execErr := nm.recordNodeUserUsageHistory(bulkCtx, db, nodeID, usageDeltas); execErr != nil {
			return execErr
		}

		return nil
	})
	if err != nil {
		nm.cfg.Logger.Warn("Failed to persist node runtime stats", "node", nodeName, "error", err)
		return
	}

	if persistedNodeUUID != "" {
		nm.updateNodeHotCache(persistedNodeUUID, singboxVersion, nodeVersion, singboxUptime, hasSingboxUptime, systemInfo, systemStats, usersOnline)
		nm.updateNodeMetricsSnapshot(persistedNodeUUID, usersOnline, trafficDelta)
	}
	for _, event := range firstConnectedEvents {
		notifications.Emit(context.Background(), nm.cfg, event)
	}
}

func (nm *NodeMonitor) updateNodeHotCache(
	nodeUUID string,
	singboxVersion string,
	nodeVersion string,
	singboxUptime int64,
	hasSingboxUptime bool,
	systemInfo json.RawMessage,
	systemStats json.RawMessage,
	usersOnline int,
) {
	if nm.hotCache == nil || strings.TrimSpace(nodeUUID) == "" {
		return
	}
	ctx := nm.globalCtx
	if ctx == nil {
		ctx = context.Background()
	}

	if err := nm.hotCache.SetUsersOnline(ctx, nodeUUID, usersOnline); err != nil {
		nm.cfg.Logger.Warn("Failed to update node users online cache", "node_uuid", nodeUUID, "error", err)
	}
	if hasSingboxUptime {
		if err := nm.hotCache.SetUptime(ctx, nodeUUID, singboxUptime); err != nil {
			nm.cfg.Logger.Warn("Failed to update node sing-box uptime cache", "node_uuid", nodeUUID, "error", err)
		}
	}
	if len(systemStats) > 0 {
		if err := nm.hotCache.SetSystemStats(ctx, nodeUUID, systemStats); err != nil {
			nm.cfg.Logger.Warn("Failed to update node system stats cache", "node_uuid", nodeUUID, "error", err)
		}
	}
	if len(systemInfo) > 0 {
		if err := nm.hotCache.SetSystemInfo(ctx, nodeUUID, systemInfo); err != nil {
			nm.cfg.Logger.Warn("Failed to update node system info cache", "node_uuid", nodeUUID, "error", err)
		}
	}
	if singboxVersion != "" || nodeVersion != "" {
		if singboxVersion == "" || nodeVersion == "" {
			hot, _ := nm.hotCache.GetOne(ctx, nodeUUID)
			if hot.Versions != nil {
				if singboxVersion == "" {
					singboxVersion = hot.Versions.Singbox
				}
				if nodeVersion == "" {
					nodeVersion = hot.Versions.Node
				}
			}
		}
		if err := nm.hotCache.SetVersions(ctx, nodeUUID, singboxVersion, nodeVersion); err != nil {
			nm.cfg.Logger.Warn("Failed to update node versions cache", "node_uuid", nodeUUID, "error", err)
		}
	}
}

type trafficStatsDelta struct {
	TotalUploadBytes   int64
	TotalDownloadBytes int64
	UserBytesByName    map[string]int64
	UsersOnline        int
	InboundByTag       map[string]TagTrafficCounters
	OutboundByTag      map[string]TagTrafficCounters
}

func extractTrafficStatsDelta(stats []*proto.Stat) trafficStatsDelta {
	delta := trafficStatsDelta{
		UserBytesByName: make(map[string]int64),
		UsersOnline:     0,
		InboundByTag:    make(map[string]TagTrafficCounters),
		OutboundByTag:   make(map[string]TagTrafficCounters),
	}
	onlineUsers := make(map[string]struct{})

	for _, stat := range stats {
		if stat == nil {
			continue
		}

		name := strings.TrimSpace(stat.GetName())
		valueRaw := strings.TrimSpace(stat.GetValue())
		if name == "" || valueRaw == "" {
			continue
		}

		value, err := strconv.ParseInt(valueRaw, 10, 64)
		if err != nil || value <= 0 {
			continue
		}

		parts := strings.Split(name, ">>>")
		if len(parts) < 3 {
			continue
		}

		switch strings.ToLower(strings.TrimSpace(parts[0])) {
		case "inbound":
			if len(parts) != 4 || !strings.EqualFold(parts[2], "traffic") {
				continue
			}
			tag := strings.TrimSpace(parts[1])
			if tag == "" {
				continue
			}
			current := delta.InboundByTag[tag]
			switch strings.ToLower(strings.TrimSpace(parts[3])) {
			case "uplink":
				current.UploadBytes += value
			case "downlink":
				current.DownloadBytes += value
			default:
				continue
			}
			delta.InboundByTag[tag] = current
		case "outbound":
			if len(parts) != 4 || !strings.EqualFold(parts[2], "traffic") {
				continue
			}
			tag := strings.TrimSpace(parts[1])
			if tag == "" {
				continue
			}
			current := delta.OutboundByTag[tag]
			switch strings.ToLower(strings.TrimSpace(parts[3])) {
			case "uplink":
				delta.TotalUploadBytes += value
				current.UploadBytes += value
			case "downlink":
				delta.TotalDownloadBytes += value
				current.DownloadBytes += value
			default:
				continue
			}
			delta.OutboundByTag[tag] = current
		case "user":
			username := strings.TrimSpace(parts[1])
			if username == "" {
				continue
			}

			if len(parts) == 4 && strings.EqualFold(parts[2], "traffic") {
				switch strings.ToLower(strings.TrimSpace(parts[3])) {
				case "uplink", "downlink":
					delta.UserBytesByName[username] += value
					onlineUsers[username] = struct{}{}
				}
			}
		}
	}

	delta.UsersOnline = len(onlineUsers)
	return delta
}

func applyConsumptionMultiplier(totalBytes int64, multiplierNano int64) int64 {
	const nanoScale = int64(1_000_000_000)

	if totalBytes <= 0 || multiplierNano <= 0 {
		return 0
	}
	if multiplierNano == nanoScale {
		return totalBytes
	}

	scaled := math.Floor((float64(multiplierNano) / float64(nanoScale)) * float64(totalBytes))
	if scaled <= 0 {
		return 0
	}
	if scaled >= float64(math.MaxInt64) {
		return math.MaxInt64
	}
	return int64(scaled)
}

func parseOptionalJSONRaw(value string) json.RawMessage {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return nil
	}
	var raw json.RawMessage
	if err := json.Unmarshal([]byte(trimmed), &raw); err != nil || len(raw) == 0 {
		return nil
	}
	return raw
}

func parseOptionalUptimeSeconds(value string) (int64, bool) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return 0, false
	}
	parsed, err := strconv.ParseInt(trimmed, 10, 64)
	if err != nil {
		return 0, false
	}
	return parsed, true
}
