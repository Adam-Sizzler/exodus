package server

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"v2ray-stat/node/config"
	"v2ray-stat/proto"
)

// TaskStatus represents the status of a task
type TaskStatus string

const (
	TaskStatusPending    TaskStatus = "pending"
	TaskStatusProcessing TaskStatus = "processing"
	TaskStatusSuccess    TaskStatus = "success"
	TaskStatusError      TaskStatus = "error"
)

// TaskResult stores the result of a task execution
type TaskResult struct {
	TaskID       string
	Status       TaskStatus
	ErrorMessage string
	Users        *proto.ListUsersResponse
	Credentials  map[string]string
	CompletedAt  time.Time
}

// TaskManager manages tasks in memory
type TaskManager struct {
	tasks map[string]*TaskResult
	mutex sync.RWMutex
	cfg   *config.NodeConfig
}

// NewTaskManager creates a new task manager
func NewTaskManager(cfg *config.NodeConfig) *TaskManager {
	return &TaskManager{
		tasks: make(map[string]*TaskResult),
		cfg:   cfg,
	}
}

// StoreTask stores a task result in memory
func (tm *TaskManager) StoreTask(taskID string, status TaskStatus, errorMessage string, users *proto.ListUsersResponse, credentials map[string]string) {
	tm.mutex.Lock()
	defer tm.mutex.Unlock()

	tm.tasks[taskID] = &TaskResult{
		TaskID:       taskID,
		Status:       status,
		ErrorMessage: errorMessage,
		Users:        users,
		Credentials:  credentials,
		CompletedAt:  time.Now(),
	}
	tm.cfg.Logger.Debug("Task result stored in memory", "task_id", taskID, "status", status)
}

// GetTask retrieves a task result from memory
func (tm *TaskManager) GetTask(taskID string) (*TaskResult, bool) {
	tm.mutex.RLock()
	defer tm.mutex.RUnlock()

	result, exists := tm.tasks[taskID]
	return result, exists
}

// DeleteTask removes a task from memory
func (tm *TaskManager) DeleteTask(taskID string) {
	tm.mutex.Lock()
	defer tm.mutex.Unlock()

	delete(tm.tasks, taskID)
	tm.cfg.Logger.Debug("Task deleted from memory", "task_id", taskID)
}

// CleanupOldTasks removes tasks older than the specified duration
func (tm *TaskManager) CleanupOldTasks(maxAge time.Duration) {
	tm.mutex.Lock()
	defer tm.mutex.Unlock()

	now := time.Now()
	for taskID, result := range tm.tasks {
		if now.Sub(result.CompletedAt) > maxAge {
			delete(tm.tasks, taskID)
			tm.cfg.Logger.Debug("Old task cleaned up", "task_id", taskID)
		}
	}
}

// ProcessTask processes a task asynchronously
func (tm *TaskManager) ProcessTask(ctx context.Context, task *proto.NodeTask, server *NodeServer) {
	tm.cfg.Logger.Info("Processing task", "task_id", task.TaskId, "operation", task.Operation)

	// Set task status to processing
	tm.StoreTask(task.TaskId, TaskStatusProcessing, "", nil, nil)

	var users *proto.ListUsersResponse
	var credentials map[string]string
	var err error

	// Unmarshal payload
	var payload map[string]interface{}
	if err := json.Unmarshal(task.Payload, &payload); err != nil {
		tm.cfg.Logger.Error("Failed to unmarshal task payload", "task_id", task.TaskId, "error", err)
		tm.StoreTask(task.TaskId, TaskStatusError, fmt.Sprintf("failed to unmarshal payload: %v", err), nil, nil)
		return
	}

	// Execute operation based on type
	switch task.Operation {
	case "add_users":
		err = tm.processAddUsers(ctx, server, payload, &users, &credentials)
	case "delete_users":
		err = tm.processDeleteUsers(ctx, server, payload, &users)
	case "set_enabled":
		err = tm.processSetEnabled(ctx, server, payload, &users)
	default:
		err = fmt.Errorf("unknown operation: %s", task.Operation)
	}

	if err != nil {
		tm.cfg.Logger.Error("Task execution failed", "task_id", task.TaskId, "operation", task.Operation, "error", err)
		tm.StoreTask(task.TaskId, TaskStatusError, err.Error(), nil, nil)
		return
	}

	tm.cfg.Logger.Info("Task completed successfully", "task_id", task.TaskId, "operation", task.Operation)
	tm.StoreTask(task.TaskId, TaskStatusSuccess, "", users, credentials)
}

func (tm *TaskManager) processAddUsers(ctx context.Context, server *NodeServer, payload map[string]interface{}, users **proto.ListUsersResponse, credentials *map[string]string) error {
	usernames, ok := payload["usernames"].([]interface{})
	if !ok {
		return fmt.Errorf("invalid usernames in payload")
	}
	inboundTag, ok := payload["inbound_tag"].(string)
	if !ok {
		return fmt.Errorf("invalid inbound_tag in payload")
	}

	usernamesStr := make([]string, len(usernames))
	for i, u := range usernames {
		usernamesStr[i] = u.(string)
	}

	req := &proto.AddUsersRequest{
		Usernames:  usernamesStr,
		InboundTag: inboundTag,
	}

	resp, err := server.AddUsers(ctx, req)
	if err != nil {
		return err
	}

	*users = resp.Users
	*credentials = resp.Credentials
	return nil
}

func (tm *TaskManager) processDeleteUsers(ctx context.Context, server *NodeServer, payload map[string]interface{}, users **proto.ListUsersResponse) error {
	usernames, ok := payload["usernames"].([]interface{})
	if !ok {
		return fmt.Errorf("invalid usernames in payload")
	}
	inboundTag, ok := payload["inbound_tag"].(string)
	if !ok {
		return fmt.Errorf("invalid inbound_tag in payload")
	}

	usernamesStr := make([]string, len(usernames))
	for i, u := range usernames {
		usernamesStr[i] = u.(string)
	}

	req := &proto.DeleteUsersRequest{
		Usernames:  usernamesStr,
		InboundTag: inboundTag,
	}

	resp, err := server.DeleteUsers(ctx, req)
	if err != nil {
		return err
	}

	*users = resp.Users
	return nil
}

func (tm *TaskManager) processSetEnabled(ctx context.Context, server *NodeServer, payload map[string]interface{}, users **proto.ListUsersResponse) error {
	usernames, ok := payload["usernames"].([]interface{})
	if !ok {
		return fmt.Errorf("invalid usernames in payload")
	}
	enabled, ok := payload["enabled"].(bool)
	if !ok {
		return fmt.Errorf("invalid enabled in payload")
	}

	usernamesStr := make([]string, len(usernames))
	for i, u := range usernames {
		usernamesStr[i] = u.(string)
	}

	req := &proto.SetUserEnabledRequest{
		Usernames: usernamesStr,
		Enabled:   enabled,
	}

	resp, err := server.SetUserEnabled(ctx, req)
	if err != nil {
		return err
	}

	*users = resp.Users
	return nil
}

// StartCleanupWorker starts a background worker to cleanup old tasks
func (tm *TaskManager) StartCleanupWorker(ctx context.Context, interval time.Duration, maxAge time.Duration) {
	tm.cfg.Logger.Info("Starting task cleanup worker", "interval", interval, "max_age", maxAge)

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			tm.CleanupOldTasks(maxAge)
		case <-ctx.Done():
			tm.cfg.Logger.Info("Task cleanup worker stopped")
			return
		}
	}
}
