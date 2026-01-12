package tasks

import (
	"context"
	"encoding/json"
	"fmt"

	"v2ray-stat/backend/config"
	"v2ray-stat/backend/db"
)

// ExecuteTaskBasedOperation executes an operation using the task-based system
func ExecuteTaskBasedOperation(
	ctx context.Context,
	tm *TaskManager,
	operation string,
	payload map[string]interface{},
	targetNodes []config.NodeConfig,
	nodeClients []*db.NodeClient,
	cfg *config.BackendConfig,
) (string, error) {
	cfg.Logger.Info("Executing task-based operation", "operation", operation, "node_count", len(targetNodes))

	// Create node name list
	nodeNames := make([]string, len(targetNodes))
	for i, node := range targetNodes {
		nodeNames[i] = node.NodeName
	}

	// Create task in database
	taskID, err := tm.CreateTask(ctx, operation, payload, nodeNames)
	if err != nil {
		cfg.Logger.Error("Failed to create task", "operation", operation, "error", err)
		return "", fmt.Errorf("create task: %w", err)
	}

	cfg.Logger.Info("Task created", "task_id", taskID, "operation", operation)

	// Submit task to each node
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		cfg.Logger.Error("Failed to marshal payload", "task_id", taskID, "error", err)
		return taskID, fmt.Errorf("marshal payload: %w", err)
	}

	for _, nodeClient := range nodeClients {
		// Find corresponding TaskNode ID
		taskNodes, err := tm.GetTaskNodes(ctx, taskID)
		if err != nil {
			cfg.Logger.Error("Failed to get task nodes", "task_id", taskID, "error", err)
			continue
		}

		var nodeTaskID string
		for _, tn := range taskNodes {
			if tn.NodeName == nodeClient.NodeName {
				nodeTaskID = tn.ID
				break
			}
		}

		if nodeTaskID == "" {
			cfg.Logger.Error("Node task ID not found", "task_id", taskID, "node_name", nodeClient.NodeName)
			continue
		}

		// Submit task to node
		err = tm.SubmitTaskToNode(ctx, nodeClient, taskID, operation, payloadJSON)
		if err != nil {
			cfg.Logger.Error("Failed to submit task to node", "task_id", taskID, "node_name", nodeClient.NodeName, "error", err)
			if updateErr := tm.UpdateNodeTaskStatus(ctx, nodeTaskID, NodeTaskStatusError, err.Error()); updateErr != nil {
				cfg.Logger.Error("Failed to update node task error status", "node_task_id", nodeTaskID, "error", updateErr)
			}
			continue
		}

		// Mark task as sent
		if err := tm.UpdateNodeTaskStatus(ctx, nodeTaskID, NodeTaskStatusSent, ""); err != nil {
			cfg.Logger.Error("Failed to update node task sent status", "node_task_id", nodeTaskID, "error", err)
		}

		cfg.Logger.Debug("Task submitted to node", "task_id", taskID, "node_name", nodeClient.NodeName)
	}

	// Update parent task status to processing
	if err := tm.UpdateTaskStatus(ctx, taskID, TaskStatusProcessing); err != nil {
		cfg.Logger.Error("Failed to update task status to processing", "task_id", taskID, "error", err)
	}

	cfg.Logger.Info("Task submission completed", "task_id", taskID, "operation", operation)
	return taskID, nil
}
