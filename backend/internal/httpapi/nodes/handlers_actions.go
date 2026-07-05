package nodes

import (
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"

	"exodus/internal/config"
	dbmanager "exodus/internal/db/manager"
	"exodus/internal/httpapi/shared"
	"exodus/internal/nodehotcache"
	monitor "exodus/internal/nodes"
	"exodus/internal/notifications"

	"github.com/google/uuid"
)

func handleEnableNode(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig, nodeUUID string) {
	node, err := getNodeByUUID(r.Context(), manager, nodeUUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			shared.SendError(w, http.StatusNotFound, "node not found", nil, cfg)
			return
		}
		shared.SendError(w, http.StatusInternalServerError, "failed to fetch node", err, cfg)
		return
	}

	inboundsMap, err := getNodeInbounds(r.Context(), manager, []string{nodeUUID})
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to fetch node inbounds", err, cfg)
		return
	}

	err = manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		if node.ActiveConfigProfileUUID == nil || len(inboundsMap[nodeUUID]) == 0 {
			_, execErr := db.ExecContext(r.Context(), `
				UPDATE nodes
				SET is_disabled = true, active_config_profile_uuid = NULL, is_connecting = false,
					is_connected = false, last_status_message = NULL, last_status_change = ?
				WHERE uuid = ?
			`, time.Now().UTC(), nodeUUID)
			return execErr
		}
		_, execErr := db.ExecContext(r.Context(), `UPDATE nodes SET is_disabled = false, updated_at = CURRENT_TIMESTAMP WHERE uuid = ?`, nodeUUID)
		return execErr
	})
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to enable node", err, cfg)
		return
	}

	monitor.RequestNodeSync()
	monitor.RequestNodeDeploy(true, nodeUUID)
	if updated, loadErr := getNodeByUUID(r.Context(), manager, nodeUUID); loadErr == nil {
		emitNodeNotification(r.Context(), cfg, notifications.EventNodeEnabled, updated, nil)
	}
	sendUpdatedNodeResponse(w, r, manager, cfg, nodeUUID)
}

func handleDisableNode(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig, nodeUUID string) {
	node, err := getNodeByUUID(r.Context(), manager, nodeUUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			shared.SendError(w, http.StatusNotFound, "node not found", nil, cfg)
			return
		}
		shared.SendError(w, http.StatusInternalServerError, "failed to fetch node", err, cfg)
		return
	}

	inboundsMap, err := getNodeInbounds(r.Context(), manager, []string{nodeUUID})
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to fetch node inbounds", err, cfg)
		return
	}

	err = manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		if node.ActiveConfigProfileUUID == nil || len(inboundsMap[nodeUUID]) == 0 {
			if _, execErr := db.ExecContext(r.Context(), `UPDATE nodes SET active_config_profile_uuid = NULL WHERE uuid = ?`, nodeUUID); execErr != nil {
				return execErr
			}
		}
		_, execErr := db.ExecContext(r.Context(), `
			UPDATE nodes
			SET is_disabled = true, is_connecting = false, is_connected = false,
				last_status_message = NULL, last_status_change = ?,
				updated_at = CURRENT_TIMESTAMP
			WHERE uuid = ?
		`, time.Now().UTC(), nodeUUID)
		return execErr
	})
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to disable node", err, cfg)
		return
	}

	monitor.RequestNodeSync()
	monitor.RequestNodeDeploy(true, nodeUUID)
	_ = nodehotcache.Default(cfg).DeleteTransient(r.Context(), nodeUUID)
	if updated, loadErr := getNodeByUUID(r.Context(), manager, nodeUUID); loadErr == nil {
		emitNodeNotification(r.Context(), cfg, notifications.EventNodeDisabled, updated, nil)
	}
	sendUpdatedNodeResponse(w, r, manager, cfg, nodeUUID)
}

func handleRestartNode(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig, nodeUUID string) {
	req, err := decodeOptionalRestartNodesRequest(r)
	if err != nil {
		shared.SendError(w, http.StatusBadRequest, "invalid JSON", err, cfg)
		return
	}

	node, err := getNodeByUUID(r.Context(), manager, nodeUUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			shared.SendError(w, http.StatusNotFound, "node not found", nil, cfg)
			return
		}
		shared.SendError(w, http.StatusInternalServerError, "failed to fetch node", err, cfg)
		return
	}
	if node.IsDisabled {
		shared.SendError(w, http.StatusBadRequest, "node is disabled", nil, cfg)
		return
	}

	monitor.RequestNodeSync()
	monitor.RequestNodeDeployWithForce(true, isForceRestartRequested(req), nodeUUID)
	shared.WriteJSON(w, http.StatusOK, map[string]any{"response": map[string]any{"eventSent": true}})
}

func handleResetNodeTraffic(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig, nodeUUID string) {
	node, err := getNodeByUUID(r.Context(), manager, nodeUUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			shared.SendError(w, http.StatusNotFound, "node not found", nil, cfg)
			return
		}
		shared.SendError(w, http.StatusInternalServerError, "failed to fetch node", err, cfg)
		return
	}

	err = manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		tx, err := db.BeginTx(r.Context(), nil)
		if err != nil {
			return err
		}
		_, err = tx.ExecContext(r.Context(), `
			INSERT INTO nodes_traffic_usage_history (node_uuid, traffic_bytes, reset_at)
			VALUES (?, ?, ?)
		`, nodeUUID, coalesceInt64Ptr(node.TrafficUsedBytes), time.Now().UTC())
		if err != nil {
			_ = tx.Rollback()
			return err
		}
		_, err = tx.ExecContext(r.Context(), `
			UPDATE nodes SET traffic_used_bytes = 0, updated_at = CURRENT_TIMESTAMP WHERE uuid = ?
		`, nodeUUID)
		if err != nil {
			_ = tx.Rollback()
			return err
		}
		return tx.Commit()
	})
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to reset node traffic", err, cfg)
		return
	}

	shared.WriteJSON(w, http.StatusOK, map[string]any{"response": map[string]any{"eventSent": true}})
}

func handleRestartAllNodes(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) {
	req, err := decodeOptionalRestartNodesRequest(r)
	if err != nil {
		shared.SendError(w, http.StatusBadRequest, "invalid JSON", err, cfg)
		return
	}

	var enabledCount int
	err = manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		return db.QueryRowContext(r.Context(), `SELECT COUNT(*) FROM nodes WHERE is_disabled = false`).Scan(&enabledCount)
	})
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to inspect nodes", err, cfg)
		return
	}
	if enabledCount == 0 {
		shared.SendError(w, http.StatusBadRequest, errNoEnabledNodes.Error(), nil, cfg)
		return
	}

	monitor.RequestNodeSync()
	monitor.RequestNodeDeployWithForce(true, isForceRestartRequested(req))
	shared.WriteJSON(w, http.StatusOK, map[string]any{"response": map[string]any{"eventSent": true}})
}

func decodeOptionalRestartNodesRequest(r *http.Request) (restartAllNodesRequest, error) {
	var req restartAllNodesRequest
	if r.Body == nil || r.ContentLength == 0 {
		return req, nil
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		if errors.Is(err, io.EOF) {
			return req, nil
		}
		return req, err
	}
	return req, nil
}

func isForceRestartRequested(req restartAllNodesRequest) bool {
	return req.ForceRestart != nil && *req.ForceRestart
}

func handleReorderNodes(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) {
	var req reorderNodesRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		shared.SendError(w, http.StatusBadRequest, "invalid JSON", err, cfg)
		return
	}
	if len(req.Nodes) == 0 {
		shared.SendError(w, http.StatusBadRequest, "nodes cannot be empty", nil, cfg)
		return
	}
	for _, item := range req.Nodes {
		if _, err := uuid.Parse(item.UUID); err != nil {
			shared.SendError(w, http.StatusBadRequest, "invalid UUID format", nil, cfg)
			return
		}
	}

	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		tx, err := db.BeginTx(r.Context(), nil)
		if err != nil {
			return err
		}
		for _, item := range req.Nodes {
			if _, err := tx.ExecContext(r.Context(), `UPDATE nodes SET view_position = ? WHERE uuid = ?`, item.ViewPosition, item.UUID); err != nil {
				_ = tx.Rollback()
				return err
			}
		}
		if _, err := tx.ExecContext(r.Context(), `SELECT setval('nodes_view_position_seq', (SELECT COALESCE(MAX(view_position), 0) FROM nodes) + 1)`); err != nil {
			_ = tx.Rollback()
			return err
		}
		return tx.Commit()
	})
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to reorder nodes", err, cfg)
		return
	}

	handleGetNodes(w, r, manager, cfg)
}

func sendUpdatedNodeResponse(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig, nodeUUID string) {
	node, err := getNodeByUUID(r.Context(), manager, nodeUUID)
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to fetch updated node", err, cfg)
		return
	}
	response, err := buildNodeResponses(r.Context(), manager, cfg, []nodeRecord{node})
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to build node response", err, cfg)
		return
	}
	shared.WriteJSON(w, http.StatusOK, map[string]any{"response": response[0]})
}
