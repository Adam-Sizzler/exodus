package nodes

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"exodus/internal/config"
	dbmanager "exodus/internal/db/manager"
	"exodus/internal/httpapi/shared"
	monitor "exodus/internal/nodes"
	"exodus/internal/notifications"
)

func handleBulkProfileModification(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) {
	var req bulkProfileModificationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		shared.SendError(w, http.StatusBadRequest, "invalid JSON", err, cfg)
		return
	}
	if err := validateUUIDs(req.UUIDs); err != nil {
		shared.SendError(w, http.StatusBadRequest, err.Error(), nil, cfg)
		return
	}
	if err := ensureConfigProfileInbounds(r.Context(), manager, req.ConfigProfile.ActiveConfigProfileUUID, req.ConfigProfile.ActiveInbounds); err != nil {
		handleConfigProfileValidationError(w, err, cfg)
		return
	}

	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		tx, err := db.BeginTx(r.Context(), nil)
		if err != nil {
			return err
		}
		for _, nodeUUID := range req.UUIDs {
			if _, err := tx.ExecContext(r.Context(), `UPDATE nodes SET active_config_profile_uuid = ?, updated_at = CURRENT_TIMESTAMP WHERE uuid = ?`, req.ConfigProfile.ActiveConfigProfileUUID, nodeUUID); err != nil {
				_ = tx.Rollback()
				return err
			}
			if err := replaceNodeInboundsTx(r.Context(), tx, nodeUUID, req.ConfigProfile.ActiveInbounds); err != nil {
				_ = tx.Rollback()
				return err
			}
		}
		return tx.Commit()
	})
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to modify nodes profile", err, cfg)
		return
	}

	monitor.RequestNodeSync()
	monitor.RequestNodeDeploy(true, req.UUIDs...)
	shared.WriteJSON(w, http.StatusOK, map[string]any{"response": map[string]any{"eventSent": true}})
}

func handleBulkNodesActions(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) {
	var req bulkNodesActionsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		shared.SendError(w, http.StatusBadRequest, "invalid JSON", err, cfg)
		return
	}
	if err := validateUUIDs(req.UUIDs); err != nil {
		shared.SendError(w, http.StatusBadRequest, err.Error(), nil, cfg)
		return
	}

	for _, nodeUUID := range req.UUIDs {
		switch req.Action {
		case "ENABLE":
			if err := performEnableAction(r.Context(), manager, nodeUUID); err != nil {
				shared.SendError(w, http.StatusInternalServerError, "failed to enable nodes", err, cfg)
				return
			}
		case "DISABLE":
			if err := performDisableAction(r.Context(), manager, nodeUUID); err != nil {
				shared.SendError(w, http.StatusInternalServerError, "failed to disable nodes", err, cfg)
				return
			}
		case "RESTART":
			if _, err := getNodeByUUID(r.Context(), manager, nodeUUID); err != nil {
				shared.SendError(w, http.StatusInternalServerError, "failed to restart nodes", err, cfg)
				return
			}
		case "RESET_TRAFFIC":
			if err := performResetTrafficAction(r.Context(), manager, nodeUUID); err != nil {
				shared.SendError(w, http.StatusInternalServerError, "failed to reset nodes traffic", err, cfg)
				return
			}
		default:
			shared.SendError(w, http.StatusBadRequest, "invalid bulk action", nil, cfg)
			return
		}
	}

	monitor.RequestNodeSync()
	monitor.RequestNodeDeploy(true, req.UUIDs...)
	switch req.Action {
	case "ENABLE":
		emitNodesByUUIDsNotification(r.Context(), manager, cfg, notifications.EventNodeEnabled, req.UUIDs)
	case "DISABLE":
		emitNodesByUUIDsNotification(r.Context(), manager, cfg, notifications.EventNodeDisabled, req.UUIDs)
	case "RESET_TRAFFIC":
		emitNodesByUUIDsNotification(r.Context(), manager, cfg, notifications.EventNodeModified, req.UUIDs)
	}
	shared.WriteJSON(w, http.StatusOK, map[string]any{"response": map[string]any{"eventSent": true}})
}

func performEnableAction(ctx context.Context, manager *dbmanager.DatabaseManager, nodeUUID string) error {
	node, err := getNodeByUUID(ctx, manager, nodeUUID)
	if err != nil {
		return err
	}
	inboundsMap, err := getNodeInbounds(ctx, manager, []string{nodeUUID})
	if err != nil {
		return err
	}
	return manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		if node.ActiveConfigProfileUUID == nil || len(inboundsMap[nodeUUID]) == 0 {
			_, execErr := db.ExecContext(ctx, `
				UPDATE nodes
				SET is_disabled = true, active_config_profile_uuid = NULL, is_connecting = false,
					is_connected = false, last_status_message = NULL, last_status_change = ?, users_online = 0
				WHERE uuid = ?
			`, time.Now().UTC(), nodeUUID)
			return execErr
		}
		_, execErr := db.ExecContext(ctx, `UPDATE nodes SET is_disabled = false, updated_at = CURRENT_TIMESTAMP WHERE uuid = ?`, nodeUUID)
		return execErr
	})
}

func performDisableAction(ctx context.Context, manager *dbmanager.DatabaseManager, nodeUUID string) error {
	node, err := getNodeByUUID(ctx, manager, nodeUUID)
	if err != nil {
		return err
	}
	inboundsMap, err := getNodeInbounds(ctx, manager, []string{nodeUUID})
	if err != nil {
		return err
	}
	return manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		if node.ActiveConfigProfileUUID == nil || len(inboundsMap[nodeUUID]) == 0 {
			if _, execErr := db.ExecContext(ctx, `UPDATE nodes SET active_config_profile_uuid = NULL WHERE uuid = ?`, nodeUUID); execErr != nil {
				return execErr
			}
		}
		_, execErr := db.ExecContext(ctx, `
			UPDATE nodes
			SET is_disabled = true, is_connecting = false, is_connected = false,
				last_status_message = NULL, last_status_change = ?, users_online = 0,
				updated_at = CURRENT_TIMESTAMP
			WHERE uuid = ?
		`, time.Now().UTC(), nodeUUID)
		return execErr
	})
}

func performResetTrafficAction(ctx context.Context, manager *dbmanager.DatabaseManager, nodeUUID string) error {
	node, err := getNodeByUUID(ctx, manager, nodeUUID)
	if err != nil {
		return err
	}
	return manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		tx, err := db.BeginTx(ctx, nil)
		if err != nil {
			return err
		}
		_, err = tx.ExecContext(ctx, `
			INSERT INTO nodes_traffic_usage_history (node_uuid, traffic_bytes, reset_at)
			VALUES (?, ?, ?)
		`, nodeUUID, coalesceInt64Ptr(node.TrafficUsedBytes), time.Now().UTC())
		if err != nil {
			_ = tx.Rollback()
			return err
		}
		if _, err := tx.ExecContext(ctx, `UPDATE nodes SET traffic_used_bytes = 0, updated_at = CURRENT_TIMESTAMP WHERE uuid = ?`, nodeUUID); err != nil {
			_ = tx.Rollback()
			return err
		}
		return tx.Commit()
	})
}
