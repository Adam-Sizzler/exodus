package tasks

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"v2ray-stat/backend/config"
	"v2ray-stat/backend/db"
	"v2ray-stat/backend/db/manager"
	"v2ray-stat/common"
	"v2ray-stat/proto"

	"github.com/google/uuid"
)

// TaskStatus represents task status
type TaskStatus string

const (
	TaskStatusPending   TaskStatus = "pending"
	TaskStatusSuccess   TaskStatus = "success"
	TaskStatusFailed    TaskStatus = "failed"
	TaskStatusTimeout   TaskStatus = "timeout"
	TaskStatusProcessing TaskStatus = "processing"
)

// NodeTaskStatus represents node task status
type NodeTaskStatus string

const (
	NodeTaskStatusPending   NodeTaskStatus = "pending"
	NodeTaskStatusSent      NodeTaskStatus = "sent"
	NodeTaskStatusSuccess   NodeTaskStatus = "success"
	NodeTaskStatusError     NodeTaskStatus = "error"
	NodeTaskStatusPolling   NodeTaskStatus = "polling"
)

// Task represents a parent task
type Task struct {
	ID          string
	Operation   string
	Payload     string
	Status      TaskStatus
	CreatedAt   time.Time
	CompletedAt *time.Time
	TimeoutAt   time.Time
}

// TaskNode represents a child task for a specific node
type TaskNode struct {
	ID           string
	TaskID       string
	NodeName     string
	Status       NodeTaskStatus
	ErrorMessage string
	SentAt       *time.Time
	CompletedAt  *time.Time
}

// TaskManager manages tasks in the database
type TaskManager struct {
	dbManager *manager.DatabaseManager
	cfg       *config.BackendConfig
}

// NewTaskManager creates a new task manager
func NewTaskManager(dbManager *manager.DatabaseManager, cfg *config.BackendConfig) *TaskManager {
	return &TaskManager{
		dbManager: dbManager,
		cfg:       cfg,
	}
}

// CreateTask creates a new task with child tasks for each node
func (tm *TaskManager) CreateTask(ctx context.Context, operation string, payload map[string]interface{}, nodeNames []string) (string, error) {
	taskID := uuid.New().String()
	now := time.Now().In(common.TimeLocation)
	timeoutAt := now.Add(5 * time.Minute)

	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		tm.cfg.Logger.Error("Failed to marshal task payload", "error", err)
		return "", fmt.Errorf("marshal payload: %w", err)
	}

	tm.cfg.Logger.Info("Creating task", "task_id", taskID, "operation", operation, "node_count", len(nodeNames))

	err = tm.dbManager.ExecuteHighPriority(func(db *sql.DB) error {
		tx, err := db.BeginTx(ctx, nil)
		if err != nil {
			return fmt.Errorf("begin transaction: %w", err)
		}
		defer tx.Rollback()

		// Create parent task
		_, err = tx.Exec(`
			INSERT INTO tasks (id, operation, payload, status, created_at, timeout_at)
			VALUES (?, ?, ?, ?, ?, ?)`,
			taskID, operation, string(payloadJSON), string(TaskStatusPending), now.Unix(), timeoutAt.Unix())
		if err != nil {
			return fmt.Errorf("insert task: %w", err)
		}

		// Create child tasks for each node
		for _, nodeName := range nodeNames {
			nodeTaskID := uuid.New().String()
			_, err = tx.Exec(`
				INSERT INTO task_nodes (id, task_id, node_name, status)
				VALUES (?, ?, ?, ?)`,
				nodeTaskID, taskID, nodeName, string(NodeTaskStatusPending))
			if err != nil {
				return fmt.Errorf("insert task_node: %w", err)
			}
		}

		if err := tx.Commit(); err != nil {
			return fmt.Errorf("commit transaction: %w", err)
		}

		tm.cfg.Logger.Debug("Task created successfully", "task_id", taskID, "child_tasks", len(nodeNames))
		return nil
	})

	if err != nil {
		tm.cfg.Logger.Error("Failed to create task", "task_id", taskID, "error", err)
		return "", err
	}

	return taskID, nil
}

// GetTask retrieves a task by ID
func (tm *TaskManager) GetTask(ctx context.Context, taskID string) (*Task, error) {
	var task Task
	err := tm.dbManager.ExecuteHighPriority(func(db *sql.DB) error {
		row := db.QueryRowContext(ctx, `
			SELECT id, operation, payload, status, created_at, completed_at, timeout_at
			FROM tasks WHERE id = ?`, taskID)

		var completedAt sql.NullInt64
		err := row.Scan(&task.ID, &task.Operation, &task.Payload, &task.Status, &task.CreatedAt, &completedAt, &task.TimeoutAt)
		if err != nil {
			return err
		}

		if completedAt.Valid {
			t := time.Unix(completedAt.Int64, 0).In(common.TimeLocation)
			task.CompletedAt = &t
		}

		return nil
	})
	return &task, err
}

// GetTaskNodes retrieves all child tasks for a parent task
func (tm *TaskManager) GetTaskNodes(ctx context.Context, taskID string) ([]TaskNode, error) {
	var taskNodes []TaskNode
	err := tm.dbManager.ExecuteHighPriority(func(db *sql.DB) error {
		rows, err := db.QueryContext(ctx, `
			SELECT id, task_id, node_name, status, error_message, sent_at, completed_at
			FROM task_nodes WHERE task_id = ?`, taskID)
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			var tn TaskNode
			var errorMessage sql.NullString
			var sentAt, completedAt sql.NullInt64

			err := rows.Scan(&tn.ID, &tn.TaskID, &tn.NodeName, &tn.Status, &errorMessage, &sentAt, &completedAt)
			if err != nil {
				return err
			}

			if errorMessage.Valid {
				tn.ErrorMessage = errorMessage.String
			}
			if sentAt.Valid {
				t := time.Unix(sentAt.Int64, 0).In(common.TimeLocation)
				tn.SentAt = &t
			}
			if completedAt.Valid {
				t := time.Unix(completedAt.Int64, 0).In(common.TimeLocation)
				tn.CompletedAt = &t
			}

			taskNodes = append(taskNodes, tn)
		}
		return rows.Err()
	})
	return taskNodes, err
}

// UpdateNodeTaskStatus updates the status of a node task
func (tm *TaskManager) UpdateNodeTaskStatus(ctx context.Context, nodeTaskID string, status NodeTaskStatus, errorMessage string) error {
	now := time.Now().In(common.TimeLocation)
	tm.cfg.Logger.Debug("Updating node task status", "node_task_id", nodeTaskID, "status", status)

	return tm.dbManager.ExecuteHighPriority(func(db *sql.DB) error {
		var query string
		var args []interface{}

		if status == NodeTaskStatusSent {
			query = `UPDATE task_nodes SET status = ?, sent_at = ? WHERE id = ?`
			args = []interface{}{string(status), now.Unix(), nodeTaskID}
		} else if status == NodeTaskStatusSuccess || status == NodeTaskStatusError {
			query = `UPDATE task_nodes SET status = ?, error_message = ?, completed_at = ? WHERE id = ?`
			args = []interface{}{string(status), errorMessage, now.Unix(), nodeTaskID}
		} else {
			query = `UPDATE task_nodes SET status = ?, error_message = ? WHERE id = ?`
			args = []interface{}{string(status), errorMessage, nodeTaskID}
		}

		_, err := db.ExecContext(ctx, query, args...)
		if err != nil {
			tm.cfg.Logger.Error("Failed to update node task status", "node_task_id", nodeTaskID, "error", err)
		}
		return err
	})
}

// UpdateTaskStatus updates the status of a parent task
func (tm *TaskManager) UpdateTaskStatus(ctx context.Context, taskID string, status TaskStatus) error {
	now := time.Now().In(common.TimeLocation)
	tm.cfg.Logger.Debug("Updating task status", "task_id", taskID, "status", status)

	return tm.dbManager.ExecuteHighPriority(func(db *sql.DB) error {
		_, err := db.ExecContext(ctx, `
			UPDATE tasks SET status = ?, completed_at = ? WHERE id = ?`,
			string(status), now.Unix(), taskID)
		if err != nil {
			tm.cfg.Logger.Error("Failed to update task status", "task_id", taskID, "error", err)
		}
		return err
	})
}

// SubmitTaskToNode submits a task to a specific node
func (tm *TaskManager) SubmitTaskToNode(ctx context.Context, nodeClient *db.NodeClient, taskID, operation string, payload []byte) error {
	tm.cfg.Logger.Debug("Submitting task to node", "task_id", taskID, "node_name", nodeClient.NodeName, "operation", operation)

	grpcCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	_, err := nodeClient.Client.SubmitTask(grpcCtx, &proto.NodeTask{
		TaskId:    taskID,
		Operation: operation,
		Payload:   payload,
	})

	if err != nil {
		tm.cfg.Logger.Error("Failed to submit task to node", "task_id", taskID, "node_name", nodeClient.NodeName, "error", err)
		return fmt.Errorf("submit task to node %s: %w", nodeClient.NodeName, err)
	}

	tm.cfg.Logger.Debug("Task submitted successfully", "task_id", taskID, "node_name", nodeClient.NodeName)
	return nil
}

// PollTaskStatus polls the status of a task from a node
func (tm *TaskManager) PollTaskStatus(ctx context.Context, nodeClient *db.NodeClient, taskID string) (*proto.TaskStatusResponse, error) {
	tm.cfg.Logger.Debug("Polling task status from node", "task_id", taskID, "node_name", nodeClient.NodeName)

	grpcCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	resp, err := nodeClient.Client.GetTaskStatus(grpcCtx, &proto.TaskStatusRequest{
		TaskId: taskID,
	})

	if err != nil {
		tm.cfg.Logger.Error("Failed to poll task status", "task_id", taskID, "node_name", nodeClient.NodeName, "error", err)
		return nil, fmt.Errorf("poll task status from node %s: %w", nodeClient.NodeName, err)
	}

	tm.cfg.Logger.Debug("Task status polled", "task_id", taskID, "node_name", nodeClient.NodeName, "status", resp.Status)
	return resp, nil
}

// CheckTaskCompletion checks if all child tasks are completed and updates parent task status
func (tm *TaskManager) CheckTaskCompletion(ctx context.Context, taskID string) error {
	tm.cfg.Logger.Debug("Checking task completion", "task_id", taskID)

	taskNodes, err := tm.GetTaskNodes(ctx, taskID)
	if err != nil {
		return fmt.Errorf("get task nodes: %w", err)
	}

	allSuccess := true
	anyError := false
	allCompleted := true

	for _, tn := range taskNodes {
		if tn.Status != NodeTaskStatusSuccess && tn.Status != NodeTaskStatusError {
			allCompleted = false
		}
		if tn.Status == NodeTaskStatusError {
			anyError = true
			allSuccess = false
		}
		if tn.Status != NodeTaskStatusSuccess {
			allSuccess = false
		}
	}

	if allCompleted {
		if allSuccess {
			tm.cfg.Logger.Info("Task completed successfully", "task_id", taskID)
			return tm.UpdateTaskStatus(ctx, taskID, TaskStatusSuccess)
		} else if anyError {
			tm.cfg.Logger.Warn("Task completed with errors", "task_id", taskID)
			return tm.UpdateTaskStatus(ctx, taskID, TaskStatusFailed)
		}
	}

	tm.cfg.Logger.Debug("Task not yet completed", "task_id", taskID, "all_completed", allCompleted)
	return nil
}

// GetPendingTasks retrieves tasks that are pending or processing
func (tm *TaskManager) GetPendingTasks(ctx context.Context) ([]Task, error) {
	var tasks []Task
	err := tm.dbManager.ExecuteHighPriority(func(db *sql.DB) error {
		rows, err := db.QueryContext(ctx, `
			SELECT id, operation, payload, status, created_at, completed_at, timeout_at
			FROM tasks WHERE status IN (?, ?)`,
			string(TaskStatusPending), string(TaskStatusProcessing))
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			var task Task
			var completedAt sql.NullInt64
			var createdAtUnix, timeoutAtUnix int64

			err := rows.Scan(&task.ID, &task.Operation, &task.Payload, &task.Status, &createdAtUnix, &completedAt, &timeoutAtUnix)
			if err != nil {
				return err
			}

			task.CreatedAt = time.Unix(createdAtUnix, 0).In(common.TimeLocation)
			task.TimeoutAt = time.Unix(timeoutAtUnix, 0).In(common.TimeLocation)

			if completedAt.Valid {
				t := time.Unix(completedAt.Int64, 0).In(common.TimeLocation)
				task.CompletedAt = &t
			}

			tasks = append(tasks, task)
		}
		return rows.Err()
	})
	return tasks, err
}

// CheckTimeouts checks for timed out tasks and marks them as timeout
func (tm *TaskManager) CheckTimeouts(ctx context.Context) error {
	now := time.Now().In(common.TimeLocation)
	tm.cfg.Logger.Debug("Checking for timed out tasks")

	return tm.dbManager.ExecuteHighPriority(func(db *sql.DB) error {
		result, err := db.ExecContext(ctx, `
			UPDATE tasks SET status = ?, completed_at = ?
			WHERE status IN (?, ?) AND timeout_at < ?`,
			string(TaskStatusTimeout), now.Unix(), string(TaskStatusPending), string(TaskStatusProcessing), now.Unix())
		if err != nil {
			tm.cfg.Logger.Error("Failed to check timeouts", "error", err)
			return err
		}

		rowsAffected, _ := result.RowsAffected()
		if rowsAffected > 0 {
			tm.cfg.Logger.Warn("Timed out tasks found", "count", rowsAffected)
		}
		return nil
	})
}
