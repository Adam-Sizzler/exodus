package subscriptionnodes

import (
	"database/sql"
	"fmt"
	"strings"
)

func (sm *SubNodeMonitor) loadActiveNodes() ([]dbSubNode, error) {
	rows, err := sm.db.Query(`
		SELECT n.uuid, n.name, n.address, n.port, n.api_schema, n.api_path, n.grpc_auth_token,
		       sns.subpage_config_uuid
		FROM sub_nodes n
		LEFT JOIN sub_nodes_to_subscription_page_config sns ON sns.node_uuid = n.uuid
		WHERE n.is_disabled = false
		ORDER BY n.view_position ASC, n.name ASC
	`)
	if err != nil {
		return nil, fmt.Errorf("query sub_nodes: %w", err)
	}
	defer rows.Close()

	nodes := make([]dbSubNode, 0)
	for rows.Next() {
		var (
			n             dbSubNode
			port          sql.NullInt64
			grpcAuthToken sql.NullString
			subpageConfig sql.NullString
		)
		if err := rows.Scan(&n.UUID, &n.Name, &n.Address, &port, &n.APISchema, &n.APIPath, &grpcAuthToken, &subpageConfig); err != nil {
			return nil, fmt.Errorf("scan sub_node: %w", err)
		}
		if port.Valid {
			n.Port = int(port.Int64)
		} else {
			n.Port = 2222
		}
		n.APISchema = normalizeSubSchema(n.APISchema)
		n.GRPCAuthToken = strings.TrimSpace(grpcAuthToken.String)
		if subpageConfig.Valid {
			n.SubpageConfigUUID = normalizeAssignedSubpageConfigUUID(subpageConfig.String)
		}
		nodes = append(nodes, n)
	}
	return nodes, rows.Err()
}

func (sm *SubNodeMonitor) updateConnectionStatus(nodeName string, isConnected, isConnecting bool, message string) {
	var (
		currentConnected  bool
		currentConnecting bool
		currentMessage    sql.NullString
	)
	err := sm.db.QueryRow(`SELECT is_connected, is_connecting, last_status_message FROM sub_nodes WHERE name = $1`, nodeName).
		Scan(&currentConnected, &currentConnecting, &currentMessage)
	if err != nil {
		if err != sql.ErrNoRows {
			sm.cfg.Logger.Warn("Failed to query subscription node status", "node", nodeName, "error", err)
		}
		return
	}

	msgStr := ""
	if currentMessage.Valid {
		msgStr = currentMessage.String
	}
	if currentConnected == isConnected && currentConnecting == isConnecting && msgStr == message {
		return
	}

	_, err = sm.db.Exec(`
		UPDATE sub_nodes
		SET is_connected = $1,
		    is_connecting = $2,
		    last_status_message = $3,
		    last_status_change = CURRENT_TIMESTAMP,
		    updated_at = CURRENT_TIMESTAMP
		WHERE name = $4
	`, isConnected, isConnecting, message, nodeName)
	if err != nil {
		sm.cfg.Logger.Warn("Failed to update subscription node status", "node", nodeName, "error", err)
	}
}
