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
		nm.updateConnectionStatus(nodeName, false, false, coreError)
	}

	singboxVersion := firstNonEmptyString(values["singbox_version"])
	nodeVersion := firstNonEmptyString(values["node_version"])
	singboxUptime, hasSingboxUptime := parseOptionalUptimeSeconds(values["singbox_uptime"])
	systemInfo := parseOptionalJSONRaw(values["system_info"])
	systemStats := parseOptionalJSONRaw(values["system_stats"])
	usersOnline := trafficDelta.UsersOnline

	persistedNodeUUID := ""
	firstConnectedEvents := make([]notifications.Event, 0)

	nodeUUID, nodeID, consumptionMultiplier, nodeConsumptionMultiplier, err := nm.getNodeMetadata(nodeName)
	if err != nil {
		nm.cfg.Logger.Warn("Failed to query node metadata", "node", nodeName, "error", err)
		return
	}
	persistedNodeUUID = nodeUUID

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
	} else {
		if _, execErr := nm.db.Exec(`
			UPDATE nodes
			SET updated_at = CURRENT_TIMESTAMP
			WHERE name = $1`, nodeName); execErr != nil {
			nm.cfg.Logger.Warn("Failed to update node updated_at", "node", nodeName, "error", execErr)
			return
		}
	}

	if len(trafficDelta.UserBytesByName) > 0 {
		var (
			userIDs            = make(map[string]int64, len(trafficDelta.UserBytesByName))
			userIDsLower       = make(map[string]int64)
			firstConnectedByID = make(map[int64]bool, len(trafficDelta.UserBytesByName))
			legacyUsernames    []string
		)

		// Step 1: Fast-path for numeric user IDs (strconv, zero DB lookup for user resolution)
		for rawKey := range trafficDelta.UserBytesByName {
			trimmed := strings.TrimSpace(rawKey)
			if trimmed == "" {
				continue
			}
			if parsedID, err := strconv.ParseInt(trimmed, 10, 64); err == nil && parsedID > 0 {
				userIDs[trimmed] = parsedID
			} else {
				legacyUsernames = append(legacyUsernames, trimmed)
			}
		}

		// Step 2: Fallback for legacy string usernames (if older node has not refreshed config yet)
		if len(legacyUsernames) > 0 {
			rows, queryErr := nm.db.Query(`
				SELECT u.id, u.username, COALESCE(ut.first_connected_at IS NOT NULL, false)
				FROM users u
				LEFT JOIN user_traffic ut ON u.id = ut.id
				WHERE u.status = 'ACTIVE' AND u.username = ANY($1)
			`, legacyUsernames)
			if queryErr == nil {
				for rows.Next() {
					var (
						userID            int64
						username          string
						hasFirstConnected bool
					)
					if scanErr := rows.Scan(&userID, &username, &hasFirstConnected); scanErr == nil {
						trimmed := strings.TrimSpace(username)
						userIDs[trimmed] = userID
						userIDsLower[strings.ToLower(trimmed)] = userID
						firstConnectedByID[userID] = hasFirstConnected
					}
				}
				_ = rows.Err()
				rows.Close()
			}

			// Fallback: check case-insensitively only for usernames not matched by exact lookup
			if len(userIDsLower) < len(legacyUsernames) {
				missingLower := make([]string, 0, len(legacyUsernames)-len(userIDsLower))
				for _, un := range legacyUsernames {
					if _, ok := userIDs[un]; !ok {
						missingLower = append(missingLower, strings.ToLower(un))
					}
				}
				if len(missingLower) > 0 {
					fbRows, fbErr := nm.db.Query(`
						SELECT u.id, u.username, COALESCE(ut.first_connected_at IS NOT NULL, false)
						FROM users u
						LEFT JOIN user_traffic ut ON u.id = ut.id
						WHERE u.status = 'ACTIVE' AND lower(u.username) = ANY($1)
					`, missingLower)
					if fbErr == nil {
						for fbRows.Next() {
							var (
								userID            int64
								username          string
								hasFirstConnected bool
							)
							if scanErr := fbRows.Scan(&userID, &username, &hasFirstConnected); scanErr == nil {
								trimmed := strings.TrimSpace(username)
								userIDs[trimmed] = userID
								userIDsLower[strings.ToLower(trimmed)] = userID
								firstConnectedByID[userID] = hasFirstConnected
							}
						}
						_ = fbRows.Err()
						fbRows.Close()
					}
				}
			}
		}

		// Step 3: Check firstConnected status for numeric userIDs if not already populated
		var numericIDsWithoutFirstConnected []int64
		for _, id := range userIDs {
			if _, ok := firstConnectedByID[id]; !ok {
				numericIDsWithoutFirstConnected = append(numericIDsWithoutFirstConnected, id)
			}
		}
		if len(numericIDsWithoutFirstConnected) > 0 {
			fcRows, fcErr := nm.db.Query(`
				SELECT id, COALESCE(first_connected_at IS NOT NULL, false)
				FROM user_traffic
				WHERE id = ANY($1)
			`, numericIDsWithoutFirstConnected)
			if fcErr == nil {
				for fcRows.Next() {
					var (
						id                int64
						hasFirstConnected bool
					)
					if scanErr := fcRows.Scan(&id, &hasFirstConnected); scanErr == nil {
						firstConnectedByID[id] = hasFirstConnected
					}
				}
				_ = fcRows.Err()
				fcRows.Close()
			}
		}

		usageDeltas := make([]userUsageDelta, 0, len(trafficDelta.UserBytesByName))
		for username, rawBytes := range trafficDelta.UserBytesByName {
			username = strings.TrimSpace(username)
			userID, ok := userIDs[username]
			if !ok {
				userID, ok = userIDsLower[strings.ToLower(username)]
			}
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
			nm.cfg.Logger.Trace("Recorded user traffic delta", "node", nodeName, "user", username, "bytes", effectiveBytes)

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

	for _, event := range firstConnectedEvents {
		notifications.Emit(context.Background(), nm.cfg, event)
	}

	nm.updateNodeMetricsSnapshot(persistedNodeUUID, usersOnline, trafficDelta)
	nm.updateHotCacheNodeRuntime(nodeName, persistedNodeUUID, singboxVersion, nodeVersion, hasSingboxUptime, singboxUptime, usersOnline, systemInfo, systemStats, trafficDelta)
}

func (nm *NodeMonitor) updateHotCacheNodeRuntime(
	_ string,
	nodeUUID string,
	singboxVersion string,
	nodeVersion string,
	hasSingboxUptime bool,
	singboxUptime int64,
	usersOnline int,
	systemInfo json.RawMessage,
	systemStats json.RawMessage,
	_ trafficStatsDelta,
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
		rawKey := strings.TrimSpace(stat.GetName())
		if rawKey == "" {
			continue
		}
		valStr := strings.TrimSpace(stat.GetValue())
		val, _ := strconv.ParseInt(valStr, 10, 64)

		if strings.HasPrefix(rawKey, "user>>>") {
			rest := rawKey[7:]
			idx := strings.Index(rest, ">>>")
			if idx > 0 && strings.HasPrefix(rest[idx+3:], "traffic>>>") {
				username := rest[:idx]
				if val > 0 {
					delta.UserBytesByName[username] += val
				}
			}
			continue
		}

		if strings.HasPrefix(rawKey, "inbound>>>") {
			rest := rawKey[10:]
			idx := strings.Index(rest, ">>>")
			if idx > 0 && strings.HasPrefix(rest[idx+3:], "traffic>>>") {
				tag := rest[:idx]
				direction := strings.ToLower(rest[idx+3+10:])
				counters := delta.InboundByTag[tag]
				switch direction {
				case "uplink":
					counters.UploadBytes += val
				case "downlink":
					counters.DownloadBytes += val
				}
				delta.InboundByTag[tag] = counters
			}
			continue
		}

		if strings.HasPrefix(rawKey, "outbound>>>") {
			rest := rawKey[11:]
			idx := strings.Index(rest, ">>>")
			if idx > 0 && strings.HasPrefix(rest[idx+3:], "traffic>>>") {
				tag := rest[:idx]
				direction := strings.ToLower(rest[idx+3+10:])
				switch direction {
				case "uplink":
					delta.TotalUploadBytes += val
				case "downlink":
					delta.TotalDownloadBytes += val
				}
				counters := delta.OutboundByTag[tag]
				switch direction {
				case "uplink":
					counters.UploadBytes += val
				case "downlink":
					counters.DownloadBytes += val
				}
				delta.OutboundByTag[tag] = counters
			}
			continue
		}

		key := strings.ToLower(rawKey)
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
				switch direction {
				case "down", "download":
					counters.DownloadBytes += val
				case "up", "upload":
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
				switch direction {
				case "down", "download":
					counters.DownloadBytes += val
				case "up", "upload":
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
