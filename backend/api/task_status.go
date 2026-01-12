package api

import (
	"encoding/json"
	"net/http"

	"v2ray-stat/backend/config"
	"v2ray-stat/backend/tasks"
)

// TaskStatusHandler handles requests for task status
func TaskStatusHandler(taskManager *tasks.TaskManager, cfg *config.BackendConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cfg.Logger.Debug("Received task status request", "method", r.Method)

		if r.Method != http.MethodGet {
			cfg.Logger.Warn("Invalid HTTP method", "method", r.Method, "expected", http.MethodGet)
			http.Error(w, `{"error": "method not allowed, use GET"}`, http.StatusMethodNotAllowed)
			return
		}

		taskID := r.URL.Query().Get("task_id")
		if taskID == "" {
			cfg.Logger.Warn("Missing task_id parameter")
			http.Error(w, `{"error": "task_id parameter is required"}`, http.StatusBadRequest)
			return
		}

		// Get task from database
		task, err := taskManager.GetTask(r.Context(), taskID)
		if err != nil {
			cfg.Logger.Error("Failed to get task", "task_id", taskID, "error", err)
			http.Error(w, `{"error": "task not found"}`, http.StatusNotFound)
			return
		}

		// Get task nodes
		taskNodes, err := taskManager.GetTaskNodes(r.Context(), taskID)
		if err != nil {
			cfg.Logger.Error("Failed to get task nodes", "task_id", taskID, "error", err)
			http.Error(w, `{"error": "failed to get task nodes"}`, http.StatusInternalServerError)
			return
		}

		// Build response
		nodeStatuses := make(map[string]interface{})
		for _, tn := range taskNodes {
			nodeStatuses[tn.NodeName] = map[string]interface{}{
				"status":        tn.Status,
				"error_message": tn.ErrorMessage,
				"sent_at":       tn.SentAt,
				"completed_at":  tn.CompletedAt,
			}
		}

		response := map[string]interface{}{
			"task_id":      task.ID,
			"operation":    task.Operation,
			"status":       task.Status,
			"created_at":   task.CreatedAt,
			"completed_at": task.CompletedAt,
			"timeout_at":   task.TimeoutAt,
			"nodes":        nodeStatuses,
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			cfg.Logger.Error("Failed to encode JSON response", "error", err)
		}
	}
}
