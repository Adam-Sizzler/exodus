package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"v2ray-stat/backend/config"
	"v2ray-stat/backend/db"
	"v2ray-stat/backend/tasks"
)

// DeleteUserTaskHandler handles HTTP POST requests for deleting users via task system
func DeleteUserTaskHandler(taskManager *tasks.TaskManager, nodeClients []*db.NodeClient, cfg *config.BackendConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cfg.Logger.Debug("Received DeleteUser task-based HTTP request", "method", r.Method)

		if r.Method != http.MethodPost {
			cfg.Logger.Warn("Invalid HTTP method", "method", r.Method, "expected", http.MethodPost)
			http.Error(w, `{"error": "method not allowed, use POST"}`, http.StatusMethodNotAllowed)
			return
		}

		var req DeleteUserRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			cfg.Logger.Error("Failed to parse JSON request", "error", err)
			http.Error(w, fmt.Sprintf(`{"error": "failed to parse JSON: %v"}`, err), http.StatusBadRequest)
			return
		}

		// Parse users if comma-separated string
		if len(req.Users) == 1 {
			req.Users = strings.Split(req.Users[0], ",")
			for i := range req.Users {
				req.Users[i] = strings.TrimSpace(req.Users[i])
			}
		}

		// Validate fields
		if len(req.Users) == 0 {
			cfg.Logger.Warn("Missing users field")
			http.Error(w, `{"error": "users is required"}`, http.StatusBadRequest)
			return
		}
		if req.InboundTag == "" {
			cfg.Logger.Warn("Missing inbound_tag")
			http.Error(w, `{"error": "inbound_tag is required"}`, http.StatusBadRequest)
			return
		}

		// Validate usernames
		var validUsers []string
		for _, user := range req.Users {
			if validateUsername(user) {
				validUsers = append(validUsers, user)
			} else {
				cfg.Logger.Warn("Invalid username", "username", user)
				http.Error(w, fmt.Sprintf(`{"error": "username %s contains invalid characters or is too long"}`, user), http.StatusBadRequest)
				return
			}
		}
		if len(validUsers) == 0 {
			cfg.Logger.Warn("No valid users provided")
			http.Error(w, `{"error": "no valid users provided"}`, http.StatusBadRequest)
			return
		}
		req.Users = validUsers

		// Determine target nodes
		var targetNodes []config.NodeConfig
		var targetClients []*db.NodeClient
		if len(req.Nodes) == 0 {
			targetNodes = GetNodesFromConfig(cfg)
			targetClients = nodeClients
			cfg.Logger.Debug("Using all nodes from config", "node_count", len(targetNodes))
		} else {
			for _, nodeName := range req.Nodes {
				for _, node := range cfg.Nodes {
					if node.NodeName == nodeName {
						targetNodes = append(targetNodes, node)
						for _, nc := range nodeClients {
							if nc.NodeName == nodeName {
								targetClients = append(targetClients, nc)
								break
							}
						}
						break
					}
				}
			}
			if len(targetNodes) == 0 {
				cfg.Logger.Warn("No valid nodes found", "requested_nodes", req.Nodes)
				http.Error(w, `{"error": "no valid nodes found for the provided names"}`, http.StatusBadRequest)
				return
			}
			cfg.Logger.Debug("Selected nodes", "node_count", len(targetNodes), "nodes", req.Nodes)
		}

		// Prepare task payload
		payload := map[string]interface{}{
			"usernames":   req.Users,
			"inbound_tag": req.InboundTag,
		}

		// Execute task-based operation
		taskID, err := tasks.ExecuteTaskBasedOperation(
			r.Context(),
			taskManager,
			"delete_users",
			payload,
			targetNodes,
			targetClients,
			cfg,
		)

		if err != nil {
			cfg.Logger.Error("Failed to execute task-based operation", "error", err)
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"error": fmt.Sprintf("failed to create task: %v", err),
			})
			return
		}

		cfg.Logger.Info("Delete user task created successfully", "task_id", taskID, "users", req.Users, "node_count", len(targetNodes))

		// Return task ID for client to poll status
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusAccepted)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"task_id":   taskID,
			"message":   "task created successfully, use /api/v1/task_status?task_id=" + taskID + " to check status",
			"usernames": req.Users,
			"nodes":     req.Nodes,
		})
	}
}
