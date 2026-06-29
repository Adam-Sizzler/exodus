package users

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"exodus/internal/config"
	dbmanager "exodus/internal/db/manager"
	"exodus/internal/httpapi/shared"
	monitor "exodus/internal/nodes"
	"exodus/internal/notifications"
)

func handleBulkDeleteUsers(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) {
	var req bulkDeleteUsersRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		shared.SendError(w, http.StatusBadRequest, "invalid JSON", err, cfg)
		return
	}
	if err := validateUUIDList(req.UUIDs); err != nil {
		shared.SendError(w, http.StatusBadRequest, err.Error(), nil, cfg)
		return
	}

	notificationRecords, err := getUserRecordsByUUIDs(r.Context(), manager, req.UUIDs)
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to fetch users", err, cfg)
		return
	}

	internalSquadNodeUUIDs := make([]string, 0)
	err = manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		tx, err := db.BeginTx(r.Context(), nil)
		if err != nil {
			return err
		}

		nodeUUIDs, nodeTargetsErr := resolveNodeUUIDsForUserUUIDsTx(r.Context(), tx, req.UUIDs)
		if nodeTargetsErr != nil {
			_ = tx.Rollback()
			return nodeTargetsErr
		}
		internalSquadNodeUUIDs = nodeUUIDs

		if _, err := tx.ExecContext(r.Context(), `DELETE FROM users WHERE uuid = ANY(?)`, req.UUIDs); err != nil {
			_ = tx.Rollback()
			return err
		}

		return tx.Commit()
	})
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to delete users", err, cfg)
		return
	}

	if len(internalSquadNodeUUIDs) > 0 {
		monitor.RequestNodeDeploy(true, internalSquadNodeUUIDs...)
	}
	emitUsersNotificationFromRecords(r.Context(), manager, cfg, notifications.EventUserDeleted, req.UUIDs, notificationRecords)
	shared.WriteJSON(w, http.StatusOK, map[string]any{"response": map[string]any{"isDeleted": true}})
}

func handleBulkResetUsersTraffic(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) {
	var req bulkDeleteUsersRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		shared.SendError(w, http.StatusBadRequest, "invalid JSON", err, cfg)
		return
	}
	if err := validateUUIDList(req.UUIDs); err != nil {
		shared.SendError(w, http.StatusBadRequest, err.Error(), nil, cfg)
		return
	}

	affectedRows, nodeUUIDs, err := resetUsersTrafficByUUIDs(r.Context(), manager, req.UUIDs)
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to reset users traffic", err, cfg)
		return
	}
	if len(nodeUUIDs) > 0 {
		monitor.RequestNodeDeploy(true, nodeUUIDs...)
	}
	emitUsersByUUIDsNotification(r.Context(), manager, cfg, notifications.EventUserTrafficReset, req.UUIDs)

	shared.WriteJSON(w, http.StatusOK, map[string]any{"response": map[string]any{"affectedRows": affectedRows}})
}

func handleBulkDeleteUsersByStatus(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) {
	var req struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		shared.SendError(w, http.StatusBadRequest, "invalid JSON", err, cfg)
		return
	}

	status := strings.ToUpper(strings.TrimSpace(req.Status))
	if !isValidUserStatus(status) {
		shared.SendError(w, http.StatusBadRequest, "invalid status", nil, cfg)
		return
	}

	var affectedRows int64
	var internalSquadNodeUUIDs []string
	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		tx, err := db.BeginTx(r.Context(), nil)
		if err != nil {
			return err
		}

		// Collect node UUIDs before deletion for deploy notification.
		rows, queryErr := tx.QueryContext(r.Context(),
			`SELECT DISTINCT n.uuid
			   FROM users u
			   JOIN internal_squad_members ism ON ism.user_id = u.t_id
			   JOIN nodes n ON n.t_id = ism.node_id
			  WHERE u.status = ?`, status)
		if queryErr != nil {
			_ = tx.Rollback()
			return queryErr
		}
		nodeUUIDs := make([]string, 0)
		for rows.Next() {
			var nodeUUID string
			if scanErr := rows.Scan(&nodeUUID); scanErr != nil {
				_ = rows.Close()
				_ = tx.Rollback()
				return scanErr
			}
			nodeUUIDs = append(nodeUUIDs, nodeUUID)
		}
		_ = rows.Close()
		if rowsErr := rows.Err(); rowsErr != nil {
			_ = tx.Rollback()
			return rowsErr
		}
		internalSquadNodeUUIDs = nodeUUIDs

		result, execErr := tx.ExecContext(r.Context(),
			`DELETE FROM users WHERE status = ?`, status)
		if execErr != nil {
			_ = tx.Rollback()
			return execErr
		}
		n, _ := result.RowsAffected()
		affectedRows = n

		return tx.Commit()
	})
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to delete users by status", err, cfg)
		return
	}

	if len(internalSquadNodeUUIDs) > 0 {
		monitor.RequestNodeDeploy(true, internalSquadNodeUUIDs...)
	}
	emitBulkSummaryNotification(r.Context(), cfg, notifications.EventUserDeleted, affectedRows)

	shared.WriteJSON(w, http.StatusOK, map[string]any{"response": map[string]any{"affectedRows": affectedRows}})
}

func handleBulkAllResetUsersTraffic(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) {
	affectedRows, nodeUUIDs, err := resetAllUsersTraffic(r.Context(), manager)
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to reset all users traffic", err, cfg)
		return
	}
	if len(nodeUUIDs) > 0 {
		monitor.RequestNodeDeploy(true, nodeUUIDs...)
	}
	emitBulkSummaryNotification(r.Context(), cfg, notifications.EventUserTrafficReset, affectedRows)

	shared.WriteJSON(w, http.StatusOK, map[string]any{"response": map[string]any{"eventSent": affectedRows > 0}})
}

func handleBulkExtendUsersExpirationDate(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) {
	var req bulkExtendExpirationDateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		shared.SendError(w, http.StatusBadRequest, "invalid JSON", err, cfg)
		return
	}
	if err := validateUUIDList(req.UUIDs); err != nil {
		shared.SendError(w, http.StatusBadRequest, err.Error(), nil, cfg)
		return
	}
	if err := validateExtendDays(req.ExtendDays); err != nil {
		shared.SendError(w, http.StatusBadRequest, err.Error(), nil, cfg)
		return
	}

	affectedRows, nodeUUIDs, err := extendUsersExpirationByUUIDs(r.Context(), manager, req.UUIDs, req.ExtendDays)
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to extend users expiration date", err, cfg)
		return
	}
	if len(nodeUUIDs) > 0 {
		monitor.RequestNodeDeploy(true, nodeUUIDs...)
	}
	emitUsersByUUIDsNotification(r.Context(), manager, cfg, notifications.EventUserModified, req.UUIDs)

	shared.WriteJSON(w, http.StatusOK, map[string]any{"response": map[string]any{"affectedRows": affectedRows}})
}

func handleBulkAllExtendUsersExpirationDate(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) {
	var req bulkAllExtendExpirationDateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		shared.SendError(w, http.StatusBadRequest, "invalid JSON", err, cfg)
		return
	}
	if err := validateExtendDays(req.ExtendDays); err != nil {
		shared.SendError(w, http.StatusBadRequest, err.Error(), nil, cfg)
		return
	}

	affectedRows, nodeUUIDs, err := extendAllUsersExpiration(r.Context(), manager, req.ExtendDays)
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to extend all users expiration date", err, cfg)
		return
	}
	if len(nodeUUIDs) > 0 {
		monitor.RequestNodeDeploy(true, nodeUUIDs...)
	}
	emitBulkSummaryNotification(r.Context(), cfg, notifications.EventUserModified, affectedRows)

	shared.WriteJSON(w, http.StatusOK, map[string]any{"response": map[string]any{"eventSent": affectedRows > 0}})
}

func handleBulkUpdateUsers(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) {
	var req bulkUpdateUsersRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		shared.SendError(w, http.StatusBadRequest, "invalid JSON", err, cfg)
		return
	}
	if err := validateUUIDList(req.UUIDs); err != nil {
		shared.SendError(w, http.StatusBadRequest, err.Error(), nil, cfg)
		return
	}
	if err := validateBulkUpdateUsersFields(req.Fields); err != nil {
		shared.SendError(w, http.StatusBadRequest, err.Error(), nil, cfg)
		return
	}

	clauses, args := buildBulkUpdateUserClauses(req.Fields)
	if len(clauses) == 0 {
		shared.SendError(w, http.StatusBadRequest, "at least one field must be provided", nil, cfg)
		return
	}

	cleanUUIDs := dedupeStrings(req.UUIDs)
	var affectedRows int64
	nodeUUIDs := make([]string, 0)
	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		tx, err := db.BeginTx(r.Context(), nil)
		if err != nil {
			return err
		}

		targets, nodeTargetsErr := resolveNodeUUIDsForUserUUIDsTx(r.Context(), tx, cleanUUIDs)
		if nodeTargetsErr != nil {
			_ = tx.Rollback()
			return nodeTargetsErr
		}
		nodeUUIDs = targets

		queryArgs := append(args, cleanUUIDs)
		query := fmt.Sprintf("UPDATE users SET %s, updated_at = CURRENT_TIMESTAMP WHERE uuid = ANY(?)", strings.Join(clauses, ", "))
		result, execErr := tx.ExecContext(r.Context(), query, queryArgs...)
		if execErr != nil {
			_ = tx.Rollback()
			return mapUserWriteError(execErr)
		}
		rows, rowsErr := result.RowsAffected()
		if rowsErr != nil {
			_ = tx.Rollback()
			return rowsErr
		}
		affectedRows = rows

		return tx.Commit()
	})
	if err != nil {
		handleUserWriteError(w, err, cfg)
		return
	}

	if affectedRows > 0 {
		if len(nodeUUIDs) > 0 {
			monitor.RequestNodeDeploy(true, nodeUUIDs...)
		}
		emitUsersByUUIDsNotification(r.Context(), manager, cfg, notifications.EventUserModified, cleanUUIDs)
	}

	shared.WriteJSON(w, http.StatusOK, map[string]any{"response": map[string]any{"affectedRows": affectedRows}})
}

func handleBulkUpdateUsersSquads(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) {
	var req bulkUpdateUsersSquadsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		shared.SendError(w, http.StatusBadRequest, "invalid JSON", err, cfg)
		return
	}
	if err := validateUUIDList(req.UUIDs); err != nil {
		shared.SendError(w, http.StatusBadRequest, err.Error(), nil, cfg)
		return
	}
	if err := validateUUIDListAllowEmpty(req.ActiveInternalSquads); err != nil {
		shared.SendError(w, http.StatusBadRequest, "invalid activeInternalSquads", err, cfg)
		return
	}

	cleanUserUUIDs := dedupeStrings(req.UUIDs)
	requestedSquads := dedupeStrings(req.ActiveInternalSquads)
	var affectedRows int64
	nodeUUIDs := make([]string, 0)
	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		tx, err := db.BeginTx(r.Context(), nil)
		if err != nil {
			return err
		}

		targets, nodeTargetsErr := resolveNodeUUIDsForUserUUIDsTx(r.Context(), tx, cleanUserUUIDs)
		if nodeTargetsErr != nil {
			_ = tx.Rollback()
			return nodeTargetsErr
		}
		squadTargets, squadTargetsErr := resolveNodeUUIDsForInternalSquadsTx(r.Context(), tx, requestedSquads)
		if squadTargetsErr != nil {
			_ = tx.Rollback()
			return squadTargetsErr
		}
		nodeUUIDs = dedupeStrings(append(targets, squadTargets...))

		rows, err := tx.QueryContext(r.Context(), `SELECT t_id FROM users WHERE uuid = ANY(?)`, cleanUserUUIDs)
		if err != nil {
			_ = tx.Rollback()
			return err
		}
		userIDs := make([]int64, 0, len(cleanUserUUIDs))
		for rows.Next() {
			var userID int64
			if err := rows.Scan(&userID); err != nil {
				_ = rows.Close()
				_ = tx.Rollback()
				return err
			}
			userIDs = append(userIDs, userID)
		}
		if err := rows.Err(); err != nil {
			_ = rows.Close()
			_ = tx.Rollback()
			return err
		}
		_ = rows.Close()

		for _, userID := range userIDs {
			if err := replaceUserInternalSquadsTx(r.Context(), tx, userID, requestedSquads); err != nil {
				_ = tx.Rollback()
				return err
			}
		}
		if len(userIDs) > 0 {
			if _, err := tx.ExecContext(r.Context(), `UPDATE users SET updated_at = CURRENT_TIMESTAMP WHERE t_id = ANY(?)`, userIDs); err != nil {
				_ = tx.Rollback()
				return err
			}
		}
		affectedRows = int64(len(userIDs))

		return tx.Commit()
	})
	if err != nil {
		handleUserWriteError(w, err, cfg)
		return
	}

	if affectedRows > 0 {
		if len(nodeUUIDs) > 0 {
			monitor.RequestNodeDeploy(true, nodeUUIDs...)
		}
		emitUsersByUUIDsNotification(r.Context(), manager, cfg, notifications.EventUserModified, cleanUserUUIDs)
	}

	shared.WriteJSON(w, http.StatusOK, map[string]any{"response": map[string]any{"affectedRows": affectedRows}})
}

func handleBulkAllUpdateUsers(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) {
	var req bulkAllUpdateUsersRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		shared.SendError(w, http.StatusBadRequest, "invalid JSON", err, cfg)
		return
	}
	if err := validateBulkAllUpdateUsersRequest(req); err != nil {
		shared.SendError(w, http.StatusBadRequest, err.Error(), nil, cfg)
		return
	}

	clauses, args := buildBulkUpdateUserClauses(bulkUpdateUsersFields{
		Status:               req.Status,
		TrafficLimitBytes:    req.TrafficLimitBytes,
		TrafficLimitStrategy: req.TrafficLimitStrategy,
		ExpireAt:             req.ExpireAt,
		Description:          req.Description,
		Tag:                  req.Tag,
		TelegramID:           req.TelegramID,
		Email:                req.Email,
		HwidDeviceLimit:      req.HwidDeviceLimit,
	})
	if len(clauses) == 0 {
		shared.SendError(w, http.StatusBadRequest, "at least one field must be provided", nil, cfg)
		return
	}

	var affectedRows int64
	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		query := fmt.Sprintf("UPDATE users SET %s, updated_at = CURRENT_TIMESTAMP", strings.Join(clauses, ", "))
		result, execErr := db.ExecContext(r.Context(), query, args...)
		if execErr != nil {
			return mapUserWriteError(execErr)
		}
		rows, rowsErr := result.RowsAffected()
		if rowsErr != nil {
			return rowsErr
		}
		affectedRows = rows
		return nil
	})
	if err != nil {
		handleUserWriteError(w, err, cfg)
		return
	}

	if affectedRows > 0 {
		monitor.RequestNodeDeploy(true)
		emitBulkSummaryNotification(r.Context(), cfg, notifications.EventUserModified, affectedRows)
	}
	shared.WriteJSON(w, http.StatusOK, map[string]any{"response": map[string]any{"eventSent": affectedRows > 0}})
}
