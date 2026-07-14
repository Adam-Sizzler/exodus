package users

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	dbmanager "exodus/internal/db/manager"
	"exodus/internal/notifications"
)

// updateConnectionStatus updates node connection status in database (only on change).
func (nm *NodeMonitor) updateConnectionStatus(nodeName string, isConnected, isConnecting bool, message string) {
	message, messageDBValue := optionalStatusMessage(message)

	var notificationEvent notifications.Event
	err := nm.manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		var (
			nodeUUID          string
			nodeAddress       string
			nodePort          sql.NullInt64
			currentConnected  bool
			currentConnecting bool
			currentMessage    sql.NullString
		)

		err := db.QueryRow(`SELECT uuid, address, port, is_connected, is_connecting, last_status_message FROM nodes WHERE name = ?`, nodeName).
			Scan(&nodeUUID, &nodeAddress, &nodePort, &currentConnected, &currentConnecting, &currentMessage)
		if err != nil {
			if err == sql.ErrNoRows {
				nm.cfg.Logger.Debug("Node not found in DB", "node", nodeName)
				return nil
			}
			return fmt.Errorf("query node status: %w", err)
		}

		msgStr := ""
		if currentMessage.Valid {
			msgStr = currentMessage.String
		}

		if currentConnected == isConnected && currentConnecting == isConnecting && msgStr == message {
			if !isConnected {
				return nm.clearDisconnectedNodeRuntimeFields(context.Background(), nodeUUID)
			}
			return nil
		}

		query := `
			UPDATE nodes
			SET is_connected = ?,
			    is_connecting = ?,
			    last_status_message = ?,
			    last_status_change = CURRENT_TIMESTAMP`
		args := []any{isConnected, isConnecting, messageDBValue}

		query += ` WHERE name = ?`
		args = append(args, nodeName)

		if _, execErr := db.Exec(query, args...); execErr != nil {
			return fmt.Errorf("update node status: %w", execErr)
		}

		if !isConnected {
			if err := nm.clearDisconnectedNodeRuntimeFields(context.Background(), nodeUUID); err != nil {
				nm.cfg.Logger.Warn("Failed to clear runtime cache fields on node status change", "node", nodeName, "error", err)
			}
		}

		if currentConnected != isConnected {
			eventName := notifications.EventNodeConnectionRestored
			if !isConnected {
				eventName = notifications.EventNodeConnectionLost
			}
			notificationEvent = notifications.Event{
				Scope: notifications.ScopeNode,
				Event: eventName,
				Data: map[string]any{
					"uuid":        nodeUUID,
					"name":        nodeName,
					"address":     nodeAddress,
					"port":        nodePort.Int64,
					"isConnected": isConnected,
					"message":     message,
				},
			}
		}

		return nil
	})

	if err != nil {
		nm.cfg.Logger.Warn("Failed to update node status in DB", "node", nodeName, "error", err)
		return
	}

	if notificationEvent.Event != "" {
		notifications.Emit(context.Background(), nm.cfg, notificationEvent)
	}
}

func (nm *NodeMonitor) recordNodeUserUsageHistory(ctx context.Context, db dbmanager.DBExecutor, nodeID int64, usageDeltas []userUsageDelta) error {
	if len(usageDeltas) == 0 {
		return nil
	}
	if nm.usageRecorder != nil {
		userBytes := make(map[int64]int64, len(usageDeltas))
		for _, delta := range usageDeltas {
			if delta.UserID <= 0 || delta.HistoryBytes <= 0 {
				continue
			}
			userBytes[delta.UserID] += delta.HistoryBytes
		}
		if len(userBytes) == 0 {
			return nil
		}
		if err := nm.usageRecorder.RecordNodeUserUsage(ctx, nodeID, userBytes); err == nil {
			return nil
		} else if nm.cfg != nil && nm.cfg.Logger != nil {
			nm.cfg.Logger.Warn("Failed to enqueue node user usage history in Redis, falling back to direct database write", "error", err)
		}
	}
	return bulkUpsertNodeUserUsageHistory(ctx, db, nodeID, usageDeltas)
}

func bulkUpsertUserTraffic(ctx context.Context, db dbmanager.DBExecutor, usageDeltas []userUsageDelta, nodeUUID string) error {
	const chunkSize = 1000

	for start := 0; start < len(usageDeltas); start += chunkSize {
		end := min(start+chunkSize, len(usageDeltas))
		chunk := usageDeltas[start:end]

		var query strings.Builder
		args := make([]any, 0, len(chunk)*3)

		query.WriteString(`
			INSERT INTO user_traffic (
				t_id, used_traffic_bytes, lifetime_used_traffic_bytes,
				online_at, last_connected_node_uuid, first_connected_at
			)
			SELECT
				v.t_id,
				v.total_bytes,
				v.total_bytes,
				now(),
				v.last_connected_node_uuid,
				now()
			FROM (VALUES `)

		for i, delta := range chunk {
			if i > 0 {
				query.WriteString(", ")
			}
			query.WriteString("(?::bigint, ?::bigint, ?::uuid)")
			args = append(args, delta.UserID, delta.TotalBytes, nodeUUID)
		}

		query.WriteString(`) AS v(t_id, total_bytes, last_connected_node_uuid)
			ON CONFLICT (t_id)
			DO UPDATE SET
				used_traffic_bytes = user_traffic.used_traffic_bytes + EXCLUDED.used_traffic_bytes,
				lifetime_used_traffic_bytes = user_traffic.lifetime_used_traffic_bytes + EXCLUDED.lifetime_used_traffic_bytes,
				online_at = now(),
				last_connected_node_uuid = EXCLUDED.last_connected_node_uuid,
				first_connected_at = COALESCE(user_traffic.first_connected_at, now())
		`)

		if _, err := db.ExecContext(ctx, query.String(), args...); err != nil {
			return err
		}
	}

	return nil
}

func bulkUpsertNodeUserUsageHistory(ctx context.Context, db dbmanager.DBExecutor, nodeID int64, usageDeltas []userUsageDelta) error {
	const chunkSize = 1000

	for start := 0; start < len(usageDeltas); start += chunkSize {
		end := min(start+chunkSize, len(usageDeltas))
		chunk := usageDeltas[start:end]

		var query strings.Builder
		args := make([]any, 0, len(chunk)*3)

		query.WriteString(`
			INSERT INTO nodes_user_usage_history (node_id, user_id, total_bytes)
			VALUES `)

		for i, delta := range chunk {
			if i > 0 {
				query.WriteString(", ")
			}
			query.WriteString("(?::bigint, ?::bigint, ?::bigint)")
			args = append(args, nodeID, delta.UserID, delta.HistoryBytes)
		}

		query.WriteString(`
			ON CONFLICT (node_id, created_at, user_id)
			DO UPDATE SET
				total_bytes = nodes_user_usage_history.total_bytes + EXCLUDED.total_bytes,
				updated_at = now()
		`)

		if _, err := db.ExecContext(ctx, query.String(), args...); err != nil {
			return err
		}
	}

	return nil
}

func (nm *NodeMonitor) clearDisconnectedNodeRuntimeFields(ctx context.Context, nodeUUID string) error {
	if nm.hotCache == nil {
		return nil
	}
	if err := nm.hotCache.DeleteTransient(ctx, nodeUUID); err != nil {
		return fmt.Errorf("clear disconnected node runtime fields: %w", err)
	}
	return nil
}

func optionalStatusMessage(message string) (string, any) {
	trimmed := strings.TrimSpace(message)
	if trimmed == "" {
		return "", nil
	}
	return trimmed, trimmed
}
