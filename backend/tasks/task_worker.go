package tasks

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"v2ray-stat/backend/config"
	"v2ray-stat/backend/db"
	"v2ray-stat/backend/db/manager"
	"v2ray-stat/proto"
)

// TaskWorker manages background task processing
type TaskWorker struct {
	taskManager *TaskManager
	dbManager   *manager.DatabaseManager
	nodeClients []*db.NodeClient
	cfg         *config.BackendConfig
	ctx         context.Context
	cancel      context.CancelFunc
	wg          sync.WaitGroup
}

// NewTaskWorker creates a new task worker
func NewTaskWorker(taskManager *TaskManager, dbManager *manager.DatabaseManager, nodeClients []*db.NodeClient, cfg *config.BackendConfig) *TaskWorker {
	ctx, cancel := context.WithCancel(context.Background())
	return &TaskWorker{
		taskManager: taskManager,
		dbManager:   dbManager,
		nodeClients: nodeClients,
		cfg:         cfg,
		ctx:         ctx,
		cancel:      cancel,
	}
}

// Start starts all worker goroutines
func (tw *TaskWorker) Start() {
	tw.cfg.Logger.Debug("Starting task workers")

	// Start task polling worker
	tw.wg.Add(1)
	go tw.pollTasksWorker()

	// Start timeout check worker
	tw.wg.Add(1)
	go tw.timeoutWorker()

	tw.cfg.Logger.Debug("Task workers started")
}

// Stop stops all worker goroutines
func (tw *TaskWorker) Stop() {
	tw.cfg.Logger.Info("Stopping task workers")
	tw.cancel()
	tw.wg.Wait()
	tw.cfg.Logger.Info("Task workers stopped")
}

// pollTasksWorker polls nodes for task status updates
func (tw *TaskWorker) pollTasksWorker() {
	defer tw.wg.Done()

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	tw.cfg.Logger.Debug("Task polling worker started")

	for {
		select {
		case <-ticker.C:
			if err := tw.pollTasks(); err != nil {
				tw.cfg.Logger.Error("Failed to poll tasks", "error", err)
			}
		case <-tw.ctx.Done():
			tw.cfg.Logger.Debug("Task polling worker stopped")
			return
		}
	}
}

// timeoutWorker checks for timed out tasks
func (tw *TaskWorker) timeoutWorker() {
	defer tw.wg.Done()

	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	tw.cfg.Logger.Debug("Timeout check worker started")

	for {
		select {
		case <-ticker.C:
			if err := tw.taskManager.CheckTimeouts(tw.ctx); err != nil {
				tw.cfg.Logger.Error("Failed to check timeouts", "error", err)
			}
		case <-tw.ctx.Done():
			tw.cfg.Logger.Debug("Timeout check worker stopped")
			return
		}
	}
}

// pollTasks polls all pending/processing tasks
func (tw *TaskWorker) pollTasks() error {
	// Get all pending/processing tasks
	tasks, err := tw.taskManager.GetPendingTasks(tw.ctx)
	if err != nil {
		return err
	}

	if len(tasks) == 0 {
		return nil
	}

	tw.cfg.Logger.Debug("Polling tasks", "count", len(tasks))

	// For each task, poll its nodes
	for _, task := range tasks {
		if err := tw.pollTask(task); err != nil {
			tw.cfg.Logger.Error("Failed to poll task", "task_id", task.ID, "error", err)
		}
	}

	return nil
}

// pollTask polls a single task's node statuses
func (tw *TaskWorker) pollTask(task Task) error {
	// Get all task nodes
	taskNodes, err := tw.taskManager.GetTaskNodes(tw.ctx, task.ID)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	for _, taskNode := range taskNodes {
		// Only poll nodes that are in sent or polling status
		if taskNode.Status != NodeTaskStatusSent && taskNode.Status != NodeTaskStatusPolling {
			continue
		}

		// Find the node client
		var nodeClient *db.NodeClient
		for _, nc := range tw.nodeClients {
			if nc.NodeName == taskNode.NodeName {
				nodeClient = nc
				break
			}
		}

		if nodeClient == nil {
			tw.cfg.Logger.Warn("Node client not found", "node_name", taskNode.NodeName)
			continue
		}

		wg.Add(1)
		go func(tn TaskNode, nc *db.NodeClient) {
			defer wg.Done()
			tw.pollNodeTask(task.ID, tn, nc)
		}(taskNode, nodeClient)
	}

	wg.Wait()

	// Check if task is completed
	if err := tw.taskManager.CheckTaskCompletion(tw.ctx, task.ID); err != nil {
		tw.cfg.Logger.Error("Failed to check task completion", "task_id", task.ID, "error", err)
	}

	// Sync users after task completion check
	taskAfterCheck, err := tw.taskManager.GetTask(tw.ctx, task.ID)
	if err == nil && (taskAfterCheck.Status == TaskStatusSuccess || taskAfterCheck.Status == TaskStatusFailed) {
		tw.cfg.Logger.Debug("Task completed, syncing users with database", "task_id", task.ID)
		tw.syncTaskResults(task, taskNodes)
	}

	return nil
}

// pollNodeTask polls a single node task
func (tw *TaskWorker) pollNodeTask(taskID string, taskNode TaskNode, nodeClient *db.NodeClient) {
	// Update status to polling if it was sent
	if taskNode.Status == NodeTaskStatusSent {
		if err := tw.taskManager.UpdateNodeTaskStatus(tw.ctx, taskNode.ID, NodeTaskStatusPolling, ""); err != nil {
			tw.cfg.Logger.Error("Failed to update node task status to polling", "node_task_id", taskNode.ID, "error", err)
			return
		}
	}

	// Poll the node for task status
	resp, err := tw.taskManager.PollTaskStatus(tw.ctx, nodeClient, taskID)
	if err != nil {
		tw.cfg.Logger.Error("Failed to poll task status", "task_id", taskID, "node_name", nodeClient.NodeName, "error", err)
		if err := tw.taskManager.UpdateNodeTaskStatus(tw.ctx, taskNode.ID, NodeTaskStatusError, err.Error()); err != nil {
			tw.cfg.Logger.Error("Failed to update node task error status", "node_task_id", taskNode.ID, "error", err)
		}
		return
	}

	// Update node task status based on response
	switch resp.Status {
	case "success":
		tw.cfg.Logger.Info("Node task completed successfully", "task_id", taskID, "node_name", nodeClient.NodeName)
		if err := tw.taskManager.UpdateNodeTaskStatus(tw.ctx, taskNode.ID, NodeTaskStatusSuccess, ""); err != nil {
			tw.cfg.Logger.Error("Failed to update node task success status", "node_task_id", taskNode.ID, "error", err)
		}

		// Store the result in database (update user_traffic and user_ids)
		if resp.Users != nil {
			tw.syncUsersFromTaskResponse(taskNode.NodeName, resp)
		}

	case "error":
		tw.cfg.Logger.Warn("Node task failed", "task_id", taskID, "node_name", nodeClient.NodeName, "error", resp.ErrorMessage)
		if err := tw.taskManager.UpdateNodeTaskStatus(tw.ctx, taskNode.ID, NodeTaskStatusError, resp.ErrorMessage); err != nil {
			tw.cfg.Logger.Error("Failed to update node task error status", "node_task_id", taskNode.ID, "error", err)
		}

	case "pending", "processing":
		tw.cfg.Logger.Debug("Node task still processing", "task_id", taskID, "node_name", nodeClient.NodeName, "status", resp.Status)
		// Keep polling

	default:
		tw.cfg.Logger.Warn("Unknown node task status", "task_id", taskID, "node_name", nodeClient.NodeName, "status", resp.Status)
	}
}

// syncUsersFromTaskResponse syncs users from task response to database
func (tw *TaskWorker) syncUsersFromTaskResponse(nodeName string, resp *proto.TaskStatusResponse) {
	if resp.Users == nil || len(resp.Users.Users) == 0 {
		tw.cfg.Logger.Debug("No users to sync from task response", "node_name", nodeName)
		return
	}
	tw.cfg.Logger.Debug("Syncing users from task response", "node_name", nodeName, "user_count", len(resp.Users.Users))

	err := tw.dbManager.ExecuteHighPriority(func(db *sql.DB) error {
		tx, err := db.BeginTx(tw.ctx, nil)
		if err != nil {
			return err
		}
		defer tx.Rollback()

		// 1. Upsert user_traffic — ТОЛЬКО enabled обновляем, created/rate НЕ трогаем
		stmtUpsertUser, err := tx.Prepare(`
            INSERT INTO user_traffic (node_name, user, rate, created, enabled)
            VALUES (?, ?, 0, ?, ?)
            ON CONFLICT(node_name, user) DO UPDATE SET
                enabled = excluded.enabled`)
		if err != nil {
			return err
		}
		defer stmtUpsertUser.Close()

		// 2. DELETE всех старых id для этого user (чтобы заменить UUID)
		stmtDeleteIDs, err := tx.Prepare(`
            DELETE FROM user_ids WHERE node_name = ? AND user = ?`)
		if err != nil {
			return err
		}
		defer stmtDeleteIDs.Close()

		// 3. INSERT новых id (обычный INSERT — после DELETE конфликтов не будет)
		stmtInsertID, err := tx.Prepare(`
            INSERT INTO user_ids (node_name, user, id, inbound_tag)
            VALUES (?, ?, ?, ?)`)
		if err != nil {
			return err
		}
		defer stmtInsertID.Close()

		currentTime := time.Now().Unix()
		for _, user := range resp.Users.Users {
			enabledStr := "false"
			if user.Enabled {
				enabledStr = "true"
			}

			// Upsert user_traffic — created НЕ обновляется на conflict
			_, err := stmtUpsertUser.Exec(nodeName, user.Username, currentTime, enabledStr)
			if err != nil {
				tw.cfg.Logger.Error("Failed to upsert user_traffic (only enabled)", "node_name", nodeName, "user", user.Username, "error", err)
				continue
			}

			// 1) Удаляем ВСЕ старые UUID для этого пользователя
			_, err = stmtDeleteIDs.Exec(nodeName, user.Username)
			if err != nil {
				tw.cfg.Logger.Error("Failed to delete old user_ids", "node_name", nodeName, "user", user.Username, "error", err)
				// Продолжаем — возможно, просто нет старых записей
			}

			// 2) Вставляем НОВЫЕ UUID из ответа ноды
			for _, ui := range user.IdInbounds {
				_, err = stmtInsertID.Exec(nodeName, user.Username, ui.Id, ui.InboundTag)
				if err != nil {
					tw.cfg.Logger.Error("Failed to insert new user_id", "node_name", nodeName, "user", user.Username, "new_id", ui.Id, "error", err)
				} else {
					tw.cfg.Logger.Info("UUID replaced", "node_name", nodeName, "user", user.Username, "new_uuid", ui.Id, "inbound_tag", ui.InboundTag)
				}
			}
		}

		return tx.Commit()
	})

	if err != nil {
		tw.cfg.Logger.Error("Failed to sync users (UUID replace)", "node_name", nodeName, "error", err)
	} else {
		tw.cfg.Logger.Info("UUID replace completed successfully", "node_name", nodeName)
	}
}

// syncTaskResults syncs task results to database after completion
func (tw *TaskWorker) syncTaskResults(task Task, taskNodes []TaskNode) {
	tw.cfg.Logger.Debug("Syncing task results", "task_id", task.ID)

	// Общая логика для всех операций
	for _, tn := range taskNodes {
		if tn.Status == NodeTaskStatusSuccess {
			tw.cfg.Logger.Debug("Task node completed successfully", "task_id", task.ID, "node_name", tn.NodeName)
		} else if tn.Status == NodeTaskStatusError {
			tw.cfg.Logger.Warn("Task node failed", "task_id", task.ID, "node_name", tn.NodeName, "error", tn.ErrorMessage)
		}
	}

	// Специфическая логика для delete_users: удаление из БД
	if task.Operation == "delete_users" {
		// Парсим payload, чтобы получить usernames
		var payload map[string]interface{}
		if err := json.Unmarshal([]byte(task.Payload), &payload); err != nil {
			tw.cfg.Logger.Error("Failed to unmarshal task payload for delete_users", "task_id", task.ID, "error", err)
			return
		}
		usernamesRaw, ok := payload["usernames"]
		if !ok {
			tw.cfg.Logger.Error("Payload missing 'usernames' for delete_users", "task_id", task.ID)
			return
		}
		usernames, ok := usernamesRaw.([]interface{})
		if !ok {
			tw.cfg.Logger.Error("Invalid 'usernames' format in payload for delete_users", "task_id", task.ID)
			return
		}
		var userList []string
		for _, u := range usernames {
			userStr, ok := u.(string)
			if ok {
				userList = append(userList, userStr)
			}
		}
		if len(userList) == 0 {
			tw.cfg.Logger.Warn("No usernames to delete from DB", "task_id", task.ID)
			return
		}

		// Удаляем только для success-нод
		err := tw.dbManager.ExecuteHighPriority(func(db *sql.DB) error {
			tx, err := db.BeginTx(tw.ctx, nil)
			if err != nil {
				return fmt.Errorf("begin transaction for delete_users: %w", err)
			}
			defer tx.Rollback()

			stmtDeleteUser, err := tx.Prepare("DELETE FROM user_traffic WHERE node_name = ? AND user = ?")
			if err != nil {
				return fmt.Errorf("prepare delete user_traffic stmt: %w", err)
			}
			defer stmtDeleteUser.Close()

			stmtDeleteID, err := tx.Prepare("DELETE FROM user_ids WHERE node_name = ? AND user = ?")
			if err != nil {
				return fmt.Errorf("prepare delete user_ids stmt: %w", err)
			}
			defer stmtDeleteID.Close()

			for _, tn := range taskNodes {
				if tn.Status != NodeTaskStatusSuccess {
					continue // Пропускаем failed ноды
				}
				for _, user := range userList {
					_, err := stmtDeleteUser.Exec(tn.NodeName, user)
					if err != nil {
						tw.cfg.Logger.Error("Failed to delete from user_traffic", "node_name", tn.NodeName, "user", user, "error", err)
						continue
					}
					_, err = stmtDeleteID.Exec(tn.NodeName, user)
					if err != nil {
						tw.cfg.Logger.Error("Failed to delete from user_ids", "node_name", tn.NodeName, "user", user, "error", err)
						continue
					}
					tw.cfg.Logger.Info("Deleted user from DB", "node_name", tn.NodeName, "user", user, "task_id", task.ID)
				}
			}

			return tx.Commit()
		})
		if err != nil {
			tw.cfg.Logger.Error("Failed to delete users from DB", "task_id", task.ID, "error", err)
		} else {
			tw.cfg.Logger.Info("Users deleted from DB successfully", "task_id", task.ID, "users", userList)
		}
	}
}
