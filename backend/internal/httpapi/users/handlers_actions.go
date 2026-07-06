package users

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"exodus/internal/config"
	dbmanager "exodus/internal/db/manager"
	"exodus/internal/httpapi/shared"
	monitor "exodus/internal/nodes"
	"exodus/internal/notifications"
)

type revokeUserSubscriptionRequest struct {
	RevokeOnlyPasswords bool `json:"revokeOnlyPasswords"`
}

func handleEnableUser(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig, userUUID string) {
	nodeUUIDs, err := updateUserStatus(r.Context(), manager, userUUID, "ACTIVE")
	if err != nil {
		handleUserActionError(w, err, cfg, "failed to enable user")
		return
	}
	if len(nodeUUIDs) > 0 {
		monitor.RequestNodeDeploy(true, nodeUUIDs...)
	}
	if record, loadErr := getUserRecordByUUID(r.Context(), manager, userUUID); loadErr == nil {
		emitUserNotification(r.Context(), manager, cfg, notifications.EventUserEnabled, record, nil)
	}
	sendUpdatedUserResponse(w, r, manager, cfg, userUUID)
}

func handleDisableUser(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig, userUUID string) {
	nodeUUIDs, err := updateUserStatus(r.Context(), manager, userUUID, "DISABLED")
	if err != nil {
		handleUserActionError(w, err, cfg, "failed to disable user")
		return
	}
	if len(nodeUUIDs) > 0 {
		monitor.RequestNodeDeploy(true, nodeUUIDs...)
	}
	if record, loadErr := getUserRecordByUUID(r.Context(), manager, userUUID); loadErr == nil {
		emitUserNotification(r.Context(), manager, cfg, notifications.EventUserDisabled, record, nil)
	}
	sendUpdatedUserResponse(w, r, manager, cfg, userUUID)
}

func handleResetUserTraffic(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig, userUUID string) {
	nodeUUIDs := make([]string, 0)
	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		tx, err := db.BeginTx(r.Context(), nil)
		if err != nil {
			return err
		}

		var (
			tID    int64
			status string
		)
		if err := tx.QueryRowContext(r.Context(), `SELECT t_id, status FROM users WHERE uuid = ?`, userUUID).Scan(&tID, &status); err != nil {
			_ = tx.Rollback()
			if errors.Is(err, sql.ErrNoRows) {
				return errUserNotFound
			}
			return err
		}

		if strings.EqualFold(status, "LIMITED") {
			targets, nodeTargetsErr := resolveNodeUUIDsForUserUUIDsTx(r.Context(), tx, []string{userUUID})
			if nodeTargetsErr != nil {
				_ = tx.Rollback()
				return nodeTargetsErr
			}
			nodeUUIDs = targets
		}

		if _, err := tx.ExecContext(r.Context(), `
			UPDATE users
			SET last_traffic_reset_at = CURRENT_TIMESTAMP,
			    last_triggered_threshold = 0,
			    status = CASE WHEN status = 'LIMITED' THEN 'ACTIVE' ELSE status END,
			    updated_at = CURRENT_TIMESTAMP
			WHERE uuid = ?
		`, userUUID); err != nil {
			_ = tx.Rollback()
			return err
		}

		if _, err := tx.ExecContext(r.Context(), `
			INSERT INTO user_traffic (t_id, used_traffic_bytes, lifetime_used_traffic_bytes)
			VALUES (?, 0, 0)
			ON CONFLICT (t_id)
			DO UPDATE SET used_traffic_bytes = 0
		`, tID); err != nil {
			_ = tx.Rollback()
			return err
		}

		return tx.Commit()
	})
	if err != nil {
		handleUserActionError(w, err, cfg, "failed to reset user traffic")
		return
	}
	if len(nodeUUIDs) > 0 {
		monitor.RequestNodeDeploy(true, nodeUUIDs...)
	}
	if record, loadErr := getUserRecordByUUID(r.Context(), manager, userUUID); loadErr == nil {
		emitUserNotification(r.Context(), manager, cfg, notifications.EventUserTrafficReset, record, nil)
		if strings.EqualFold(record.Status, "ACTIVE") {
			emitUserNotification(r.Context(), manager, cfg, notifications.EventUserEnabled, record, map[string]any{"reason": "traffic_reset"})
		}
	}
	sendUpdatedUserResponse(w, r, manager, cfg, userUUID)
}

func handleRevokeUserSubscription(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig, userUUID string) {
	req := revokeUserSubscriptionRequest{}
	if r.Body != nil {
		body, _ := io.ReadAll(io.LimitReader(r.Body, 1<<20))
		if len(strings.TrimSpace(string(body))) > 0 {
			if err := json.Unmarshal(body, &req); err != nil {
				shared.SendError(w, http.StatusBadRequest, "invalid JSON", err, cfg)
				return
			}
		}
	}

	shortUUID := generateSubscriptionShortUUID()
	credentials, err := newUserProtocolCredentials(nil, nil, nil, nil, nil, nil, nil)
	if shortUUID == "" || err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to generate user credentials", err, cfg)
		return
	}

	nodeUUIDs := make([]string, 0)
	err = manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		tx, err := db.BeginTx(r.Context(), nil)
		if err != nil {
			return err
		}

		targets, nodeTargetsErr := resolveNodeUUIDsForUserUUIDsTx(r.Context(), tx, []string{userUUID})
		if nodeTargetsErr != nil {
			_ = tx.Rollback()
			return nodeTargetsErr
		}
		nodeUUIDs = targets

		query := `
			UPDATE users
			SET trojan_password = ?,
			    vless_uuid = ?,
			    ss_password = ?,
			    naive_password = ?,
			    shadowtls_password = ?,
			    hysteria2_password = ?,
			    anytls_password = ?,
			    sub_revoked_at = CURRENT_TIMESTAMP,
			    updated_at = CURRENT_TIMESTAMP`
		args := []any{
			credentials.TrojanPassword,
			credentials.VlessUUID,
			credentials.SSPassword,
			credentials.NaivePassword,
			credentials.ShadowtlsPassword,
			credentials.Hysteria2Password,
			credentials.AnytlsPassword,
		}
		if !req.RevokeOnlyPasswords {
			query += `, short_uuid = ?`
			args = append(args, shortUUID)
		}
		query += `
			WHERE uuid = ?
		`
		args = append(args, userUUID)

		result, err := tx.ExecContext(r.Context(), query, args...)
		if err != nil {
			_ = tx.Rollback()
			return mapUserWriteError(err)
		}
		rows, err := result.RowsAffected()
		if err != nil {
			_ = tx.Rollback()
			return err
		}
		if rows == 0 {
			_ = tx.Rollback()
			return errUserNotFound
		}
		return tx.Commit()
	})
	if err != nil {
		handleUserWriteError(w, err, cfg)
		return
	}
	if len(nodeUUIDs) > 0 {
		monitor.RequestNodeDeploy(true, nodeUUIDs...)
	}
	if record, loadErr := getUserRecordByUUID(r.Context(), manager, userUUID); loadErr == nil {
		emitUserNotification(r.Context(), manager, cfg, notifications.EventUserRevoked, record, nil)
	}
	sendUpdatedUserResponse(w, r, manager, cfg, userUUID)
}

func sendUpdatedUserResponse(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig, userUUID string) {
	record, err := getUserRecordByUUID(r.Context(), manager, userUUID)
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to fetch updated user", err, cfg)
		return
	}
	response, err := buildUserResponses(r.Context(), manager, []userRecord{record}, resolveUsersSubscriptionBase(r.Context(), manager, r, cfg))
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to build updated user response", err, cfg)
		return
	}
	shared.WriteJSON(w, http.StatusOK, map[string]any{"response": response[0]})
}

func handleUserActionError(w http.ResponseWriter, err error, cfg *config.BackendConfig, message string) {
	if errors.Is(err, errUserNotFound) {
		shared.SendError(w, http.StatusNotFound, "user not found", nil, cfg)
		return
	}
	shared.SendError(w, http.StatusInternalServerError, message, err, cfg)
}

func updateUserStatus(ctx context.Context, manager *dbmanager.DatabaseManager, userUUID string, status string) ([]string, error) {
	nodeUUIDs := make([]string, 0)
	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		tx, err := db.BeginTx(ctx, nil)
		if err != nil {
			return err
		}

		targets, nodeTargetsErr := resolveNodeUUIDsForUserUUIDsTx(ctx, tx, []string{userUUID})
		if nodeTargetsErr != nil {
			_ = tx.Rollback()
			return nodeTargetsErr
		}
		nodeUUIDs = targets

		result, err := tx.ExecContext(ctx, `
			UPDATE users
			SET status = ?, updated_at = CURRENT_TIMESTAMP
			WHERE uuid = ?
		`, status, userUUID)
		if err != nil {
			_ = tx.Rollback()
			return err
		}
		rows, err := result.RowsAffected()
		if err != nil {
			_ = tx.Rollback()
			return err
		}
		if rows == 0 {
			_ = tx.Rollback()
			return errUserNotFound
		}
		return tx.Commit()
	})
	return nodeUUIDs, err
}

func plannedUserStatusForUpdate(record userRecord, req updateUserRequest, now time.Time) (string, bool) {
	if req.Status != nil {
		return strings.ToUpper(strings.TrimSpace(*req.Status)), true
	}

	if req.TrafficLimitBytes != nil && strings.EqualFold(record.Status, "LIMITED") {
		if *req.TrafficLimitBytes == 0 || *req.TrafficLimitBytes > record.TrafficLimitBytes {
			return "ACTIVE", true
		}
	}

	if req.ExpireAt != nil && strings.EqualFold(record.Status, "EXPIRED") {
		parsed, err := time.Parse(time.RFC3339, strings.TrimSpace(*req.ExpireAt))
		if err == nil {
			newExpireAt := parsed.UTC()
			if !newExpireAt.Equal(record.ExpireAt.UTC()) && newExpireAt.After(now.UTC()) {
				return "ACTIVE", true
			}
		}
	}

	return "", false
}

func userConfigPresenceChanges(previousStatus string, nextStatus string) bool {
	previousActive := strings.EqualFold(previousStatus, "ACTIVE")
	nextActive := strings.EqualFold(nextStatus, "ACTIVE")
	return previousActive != nextActive
}

func validateExtendDays(days int) error {
	if days < 1 || days > 9999 {
		return fmt.Errorf("extendDays must be between 1 and 9999")
	}
	return nil
}
