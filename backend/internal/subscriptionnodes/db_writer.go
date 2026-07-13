package subscriptionnodes

import (
	"database/sql"
	"fmt"
	"strings"

	dbmanager "exodus/internal/db/manager"
)

func (sm *SubNodeMonitor) loadActiveNodes() ([]dbSubNode, error) {
	nodes := make([]dbSubNode, 0)
	err := sm.manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		rows, err := db.Query(`
			SELECT n.uuid, n.name, n.address, n.port, n.api_schema, n.api_path, n.grpc_auth_token,
			       sns.subpage_config_uuid
			FROM sub_nodes n
			LEFT JOIN sub_nodes_to_subscription_page_config sns ON sns.node_uuid = n.uuid
			WHERE n.is_disabled = false
			ORDER BY n.view_position ASC, n.name ASC
		`)
		if err != nil {
			return fmt.Errorf("query sub_nodes: %w", err)
		}
		defer rows.Close()

		for rows.Next() {
			var (
				n             dbSubNode
				port          sql.NullInt64
				grpcAuthToken sql.NullString
				subpageConfig sql.NullString
			)
			if err := rows.Scan(&n.UUID, &n.Name, &n.Address, &port, &n.APISchema, &n.APIPath, &grpcAuthToken, &subpageConfig); err != nil {
				return fmt.Errorf("scan sub_node: %w", err)
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
		return rows.Err()
	})
	return nodes, err
}

func (sm *SubNodeMonitor) updateConnectionStatus(nodeName string, isConnected, isConnecting bool, message string) {
	err := sm.manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		var (
			currentConnected  bool
			currentConnecting bool
			currentMessage    sql.NullString
		)
		err := db.QueryRow(`SELECT is_connected, is_connecting, last_status_message FROM sub_nodes WHERE name = ?`, nodeName).
			Scan(&currentConnected, &currentConnecting, &currentMessage)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil
			}
			return fmt.Errorf("query sub node status: %w", err)
		}

		msgStr := ""
		if currentMessage.Valid {
			msgStr = currentMessage.String
		}
		if currentConnected == isConnected && currentConnecting == isConnecting && msgStr == message {
			return nil
		}

		_, err = db.Exec(`
			UPDATE sub_nodes
			SET is_connected = ?,
			    is_connecting = ?,
			    last_status_message = ?,
			    last_status_change = CURRENT_TIMESTAMP,
			    updated_at = CURRENT_TIMESTAMP
			WHERE name = ?
		`, isConnected, isConnecting, message, nodeName)
		if err != nil {
			return fmt.Errorf("update sub node status: %w", err)
		}
		return nil
	})
	if err != nil {
		sm.cfg.Logger.Warn("Failed to update subscription node status", "node", nodeName, "error", err)
	}
}
