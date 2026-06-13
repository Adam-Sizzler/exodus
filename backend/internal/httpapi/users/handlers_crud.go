package users

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"exodus/internal/config"
	dbmanager "exodus/internal/db/manager"
	"exodus/internal/httpapi/shared"
	monitor "exodus/internal/nodes"
	"exodus/internal/notifications"
)

func handleGetUsers(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) {
	query, err := parseUsersTableQuery(r)
	if err != nil {
		shared.SendError(w, http.StatusBadRequest, "invalid table query", err, cfg)
		return
	}

	records, err := getAllUserRecords(r.Context(), manager)
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to fetch users", err, cfg)
		return
	}

	response, err := buildUserResponses(r.Context(), manager, records, resolveUsersSubscriptionBase(r.Context(), manager, r, cfg))
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to build users response", err, cfg)
		return
	}

	response = filterUsersTableResponse(response, query.Filters, query.FilterModes)
	sortUsersTableResponse(response, query.Sorting)
	total := len(response)
	response = paginateUsersTableResponse(response, query.Start, query.Size)

	shared.WriteJSON(w, http.StatusOK, map[string]any{
		"response": map[string]any{
			"users": response,
			"total": total,
		},
	})
}

func handleGetUser(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig, userUUID string) {
	record, err := getUserRecordByUUID(r.Context(), manager, userUUID)
	if err != nil {
		if errors.Is(err, errUserNotFound) {
			shared.SendError(w, http.StatusNotFound, "user not found", nil, cfg)
			return
		}
		shared.SendError(w, http.StatusInternalServerError, "failed to fetch user", err, cfg)
		return
	}

	response, err := buildUserResponses(r.Context(), manager, []userRecord{record}, resolveUsersSubscriptionBase(r.Context(), manager, r, cfg))
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to build user response", err, cfg)
		return
	}

	shared.WriteJSON(w, http.StatusOK, map[string]any{"response": response[0]})
}

func handleCreateUser(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) {
	var req createUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		shared.SendError(w, http.StatusBadRequest, "invalid JSON", err, cfg)
		return
	}

	if err := validateCreateUserRequest(req); err != nil {
		shared.SendError(w, http.StatusBadRequest, err.Error(), nil, cfg)
		return
	}

	userUUID := coalesceUUID(req.UUID)
	shortUUID := coalesceShortUUID(req.ShortUUID)
	trojanPassword := coalesceRandomString(req.TrojanPassword, 16)
	vlessUUID := coalesceUUID(req.VlessUUID)
	ssPassword := coalesceRandomString(req.SSPassword, 16)
	naivePassword := coalesceRandomString(req.NaivePassword, 16)
	shadowtlsPassword := coalesceRandomString(req.ShadowtlsPassword, 16)
	hysteria2Password := coalesceRandomString(req.Hysteria2Password, 16)
	anytlsPassword := coalesceRandomString(req.AnytlsPassword, 16)

	if trojanPassword == "" || ssPassword == "" || naivePassword == "" || shadowtlsPassword == "" || hysteria2Password == "" || anytlsPassword == "" {
		shared.SendError(w, http.StatusInternalServerError, "failed to generate user credentials", nil, cfg)
		return
	}

	expireAt, _ := time.Parse(time.RFC3339, req.ExpireAt)
	createdAt := time.Now().UTC()
	if req.CreatedAt != nil && strings.TrimSpace(*req.CreatedAt) != "" {
		createdAt, _ = time.Parse(time.RFC3339, strings.TrimSpace(*req.CreatedAt))
	}
	var lastTrafficResetAt any
	if req.LastTrafficResetAt != nil && strings.TrimSpace(*req.LastTrafficResetAt) != "" {
		if parsed, err := time.Parse(time.RFC3339, strings.TrimSpace(*req.LastTrafficResetAt)); err == nil {
			lastTrafficResetAt = parsed
		}
	}

	internalSquadNodeUUIDs := make([]string, 0)
	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		tx, err := db.BeginTx(r.Context(), nil)
		if err != nil {
			return err
		}

		var tID int64
		insertErr := tx.QueryRowContext(r.Context(), `
			INSERT INTO users (
					uuid, short_uuid, username, status, traffic_limit_bytes, traffic_limit_strategy,
					expire_at, last_traffic_reset_at, sub_revoked_at,
					trojan_password, vless_uuid, ss_password, naive_password, shadowtls_password, hysteria2_password, anytls_password,
					description, tag, telegram_id, email,
					hwid_device_limit, external_squad_uuid, last_triggered_threshold, created_at, updated_at
				) VALUES (?, ?, ?, ?, ?, ?, ?, ?, NULL, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, 0, ?, ?)
				RETURNING t_id
			`,
			userUUID,
			shortUUID,
			strings.TrimSpace(req.Username),
			normalizeUserStatus(req.Status),
			coalesceInt64(req.TrafficLimitBytes, 0),
			normalizeTrafficStrategy(req.TrafficLimitStrategy),
			expireAt.UTC(),
			lastTrafficResetAt,
			trojanPassword,
			vlessUUID,
			ssPassword,
			naivePassword,
			shadowtlsPassword,
			hysteria2Password,
			anytlsPassword,
			normalizeNullableString(req.Description),
			normalizeUserTag(req.Tag),
			req.TelegramID,
			normalizeNullableString(req.Email),
			req.HwidDeviceLimit,
			normalizeNullableString(req.ExternalSquadUUID),
			createdAt.UTC(),
			createdAt.UTC(),
		).Scan(&tID)
		if insertErr != nil {
			_ = tx.Rollback()
			return mapUserWriteError(insertErr)
		}

		if _, err := tx.ExecContext(r.Context(), `
			INSERT INTO user_traffic (
				t_id, used_traffic_bytes, lifetime_used_traffic_bytes, online_at,
				last_connected_node_uuid, first_connected_at
			) VALUES (?, 0, 0, NULL, NULL, NULL)
		`, tID); err != nil {
			_ = tx.Rollback()
			return err
		}

		if err := replaceUserInternalSquadsTx(r.Context(), tx, tID, req.ActiveInternalSquads); err != nil {
			_ = tx.Rollback()
			return err
		}
		requestedSquads := dedupeStrings(req.ActiveInternalSquads)
		if len(requestedSquads) > 0 {
			nodeUUIDs, nodeTargetsErr := resolveNodeUUIDsForInternalSquadsTx(r.Context(), tx, requestedSquads)
			if nodeTargetsErr != nil {
				_ = tx.Rollback()
				return nodeTargetsErr
			}
			internalSquadNodeUUIDs = nodeUUIDs
		}

		return tx.Commit()
	})
	if err != nil {
		handleUserWriteError(w, err, cfg)
		return
	}

	record, err := getUserRecordByUUID(r.Context(), manager, userUUID)
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to fetch created user", err, cfg)
		return
	}
	response, err := buildUserResponses(r.Context(), manager, []userRecord{record}, resolveUsersSubscriptionBase(r.Context(), manager, r, cfg))
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to build created user response", err, cfg)
		return
	}

	if strings.EqualFold(normalizeUserStatus(req.Status), "ACTIVE") && len(internalSquadNodeUUIDs) > 0 {
		monitor.RequestNodeDeploy(true, internalSquadNodeUUIDs...)
	}
	emitUserNotification(r.Context(), manager, cfg, notifications.EventUserCreated, record, nil)
	shared.WriteJSON(w, http.StatusCreated, map[string]any{"response": response[0]})
}

func handleUpdateUser(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) {
	var req updateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		shared.SendError(w, http.StatusBadRequest, "invalid JSON", err, cfg)
		return
	}
	if err := validateUpdateUserRequest(req); err != nil {
		shared.SendError(w, http.StatusBadRequest, err.Error(), nil, cfg)
		return
	}

	targetUUID, err := resolveUserUUIDForUpdate(r.Context(), manager, req.UUID, req.Username)
	if err != nil {
		if errors.Is(err, errUserNotFound) {
			shared.SendError(w, http.StatusNotFound, "user not found", nil, cfg)
			return
		}
		shared.SendError(w, http.StatusBadRequest, err.Error(), nil, cfg)
		return
	}

	record, err := getUserRecordByUUID(r.Context(), manager, targetUUID)
	if err != nil {
		if errors.Is(err, errUserNotFound) {
			shared.SendError(w, http.StatusNotFound, "user not found", nil, cfg)
			return
		}
		shared.SendError(w, http.StatusInternalServerError, "failed to fetch user", err, cfg)
		return
	}

	internalSquadsChanged := false
	internalSquadNodeUUIDs := make([]string, 0)
	statusNodeUUIDs := make([]string, 0)
	statusToSet, shouldSetStatus := plannedUserStatusForUpdate(record, req, time.Now().UTC())
	statusDeployRequired := shouldSetStatus && userConfigPresenceChanges(record.Status, statusToSet)
	err = manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		tx, err := db.BeginTx(r.Context(), nil)
		if err != nil {
			return err
		}

		clauses := make([]string, 0)
		args := make([]any, 0)
		add := func(column string, value any) {
			clauses = append(clauses, fmt.Sprintf("%s = ?", column))
			args = append(args, value)
		}

		if statusDeployRequired {
			nodeUUIDs, nodeTargetsErr := resolveNodeUUIDsForUserUUIDsTx(r.Context(), tx, []string{targetUUID})
			if nodeTargetsErr != nil {
				_ = tx.Rollback()
				return nodeTargetsErr
			}
			statusNodeUUIDs = nodeUUIDs
		}

		if shouldSetStatus {
			add("status", statusToSet)
		}
		if req.TrafficLimitBytes != nil {
			add("traffic_limit_bytes", *req.TrafficLimitBytes)
		}
		if req.TrafficLimitStrategy != nil {
			add("traffic_limit_strategy", strings.ToUpper(strings.TrimSpace(*req.TrafficLimitStrategy)))
		}
		if req.ExpireAt != nil {
			parsed, _ := time.Parse(time.RFC3339, strings.TrimSpace(*req.ExpireAt))
			add("expire_at", parsed.UTC())
		}
		if req.Description.Set {
			if req.Description.Value == nil || strings.TrimSpace(*req.Description.Value) == "" {
				clauses = append(clauses, "description = NULL")
			} else {
				add("description", strings.TrimSpace(*req.Description.Value))
			}
		}
		if req.Tag.Set {
			if req.Tag.Value == nil || strings.TrimSpace(*req.Tag.Value) == "" {
				clauses = append(clauses, "tag = NULL")
			} else {
				add("tag", strings.ToUpper(strings.TrimSpace(*req.Tag.Value)))
			}
		}
		if req.TelegramID.Set {
			if req.TelegramID.Value == nil {
				clauses = append(clauses, "telegram_id = NULL")
			} else {
				add("telegram_id", *req.TelegramID.Value)
			}
		}
		if req.Email.Set {
			if req.Email.Value == nil || strings.TrimSpace(*req.Email.Value) == "" {
				clauses = append(clauses, "email = NULL")
			} else {
				add("email", strings.TrimSpace(*req.Email.Value))
			}
		}
		if req.HwidDeviceLimit.Set {
			if req.HwidDeviceLimit.Value == nil {
				clauses = append(clauses, "hwid_device_limit = NULL")
			} else {
				add("hwid_device_limit", *req.HwidDeviceLimit.Value)
			}
		}
		addOptionalCredential := func(field OptionalString, column string, nullable bool) {
			if !field.Set {
				return
			}
			if field.Value == nil {
				if nullable {
					clauses = append(clauses, fmt.Sprintf("%s = NULL", column))
				}
				return
			}
			add(column, strings.TrimSpace(*field.Value))
		}
		addOptionalCredential(req.TrojanPassword, "trojan_password", false)
		addOptionalCredential(req.VlessUUID, "vless_uuid", false)
		addOptionalCredential(req.SSPassword, "ss_password", false)
		addOptionalCredential(req.NaivePassword, "naive_password", true)
		addOptionalCredential(req.ShadowtlsPassword, "shadowtls_password", true)
		addOptionalCredential(req.Hysteria2Password, "hysteria2_password", true)
		addOptionalCredential(req.AnytlsPassword, "anytls_password", true)
		if req.ExternalSquadUUID.Set {
			if req.ExternalSquadUUID.Value == nil || strings.TrimSpace(*req.ExternalSquadUUID.Value) == "" {
				clauses = append(clauses, "external_squad_uuid = NULL")
			} else {
				add("external_squad_uuid", strings.TrimSpace(*req.ExternalSquadUUID.Value))
			}
		}

		if len(clauses) > 0 {
			args = append(args, targetUUID)
			query := fmt.Sprintf("UPDATE users SET %s, updated_at = CURRENT_TIMESTAMP WHERE uuid = ?", strings.Join(clauses, ", "))
			if _, err := tx.ExecContext(r.Context(), query, args...); err != nil {
				_ = tx.Rollback()
				return mapUserWriteError(err)
			}
		}

		if req.ActiveInternalSquads != nil {
			currentSquads, loadErr := getUserInternalSquadsTx(r.Context(), tx, record.TID)
			if loadErr != nil {
				_ = tx.Rollback()
				return loadErr
			}
			requestedSquads := dedupeStrings(*req.ActiveInternalSquads)
			if internalSquadSetsDiffer(currentSquads, requestedSquads) {
				affectedSquads := dedupeStrings(append(append([]string{}, currentSquads...), requestedSquads...))
				nodeUUIDs, nodeTargetsErr := resolveNodeUUIDsForInternalSquadsTx(r.Context(), tx, affectedSquads)
				if nodeTargetsErr != nil {
					_ = tx.Rollback()
					return nodeTargetsErr
				}
				if err := replaceUserInternalSquadsTx(r.Context(), tx, record.TID, requestedSquads); err != nil {
					_ = tx.Rollback()
					return err
				}
				internalSquadNodeUUIDs = nodeUUIDs
				internalSquadsChanged = true
			}
		}

		return tx.Commit()
	})
	if err != nil {
		handleUserWriteError(w, err, cfg)
		return
	}

	updatedRecord, err := getUserRecordByUUID(r.Context(), manager, targetUUID)
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to fetch updated user", err, cfg)
		return
	}
	response, err := buildUserResponses(r.Context(), manager, []userRecord{updatedRecord}, resolveUsersSubscriptionBase(r.Context(), manager, r, cfg))
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to build updated user response", err, cfg)
		return
	}

	deployNodeUUIDs := dedupeStrings(append(statusNodeUUIDs, internalSquadNodeUUIDs...))
	if (internalSquadsChanged || statusDeployRequired) && len(deployNodeUUIDs) > 0 {
		monitor.RequestNodeDeploy(true, deployNodeUUIDs...)
	}
	emitUserNotification(r.Context(), manager, cfg, notifications.EventUserModified, updatedRecord, nil)
	if statusChanged := userStatusChangedNotification(record.Status, updatedRecord.Status); statusChanged != "" {
		emitUserNotification(r.Context(), manager, cfg, statusChanged, updatedRecord, nil)
	}
	shared.WriteJSON(w, http.StatusOK, map[string]any{"response": response[0]})
}

func handleGetUserSubscriptionRequestHistory(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig, userUUID string) {
	records := make([]userSubscriptionRequestHistoryRecord, 0)
	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		var exists bool
		if err := db.QueryRowContext(r.Context(), `SELECT EXISTS(SELECT 1 FROM users WHERE uuid = ?)`, userUUID).Scan(&exists); err != nil {
			return err
		}
		if !exists {
			return errUserNotFound
		}

		rows, err := db.QueryContext(r.Context(), `
			SELECT id, user_uuid, request_ip, user_agent, request_at
			FROM user_subscription_request_history
			WHERE user_uuid = ?
			ORDER BY request_at DESC
			LIMIT 24
		`, userUUID)
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			var item userSubscriptionRequestHistoryRecord
			var requestAt time.Time
			if scanErr := rows.Scan(&item.ID, &item.UserUUID, &item.RequestIP, &item.UserAgent, &requestAt); scanErr != nil {
				return scanErr
			}
			item.RequestAt = requestAt.UTC().Format("2006-01-02T15:04:05.000Z")
			records = append(records, item)
		}
		return rows.Err()
	})
	if err != nil {
		if errors.Is(err, errUserNotFound) {
			shared.SendError(w, http.StatusNotFound, "user not found", nil, cfg)
			return
		}
		shared.SendError(w, http.StatusInternalServerError, "failed to fetch user subscription request history", err, cfg)
		return
	}

	shared.WriteJSON(w, http.StatusOK, map[string]any{
		"response": map[string]any{
			"records": records,
			"total":   len(records),
		},
	})
}

func handleGetUserAccessibleNodes(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig, userUUID string) {
	activeNodes := make([]userAccessibleNode, 0)
	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		var userID int64
		if err := db.QueryRowContext(r.Context(), `SELECT t_id FROM users WHERE uuid = ?`, userUUID).Scan(&userID); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return errUserNotFound
			}
			return err
		}

		rows, err := db.QueryContext(r.Context(), `
			SELECT
				n.uuid,
				n.name,
				n.country_code,
				cp.uuid,
				cp.name,
				sq.uuid,
				sq.name,
				cpi.tag
			FROM nodes n
			INNER JOIN config_profiles cp ON cp.uuid = n.active_config_profile_uuid
			INNER JOIN config_profile_inbounds cpi ON cpi.profile_uuid = cp.uuid
			INNER JOIN config_profile_inbounds_to_nodes cpin
				ON cpin.config_profile_inbound_uuid = cpi.uuid
				AND cpin.node_uuid = n.uuid
			INNER JOIN internal_squad_inbounds isi ON isi.inbound_uuid = cpi.uuid
			INNER JOIN internal_squads sq ON sq.uuid = isi.internal_squad_uuid
			INNER JOIN internal_squad_members ism
				ON ism.internal_squad_uuid = sq.uuid
				AND ism.user_id = ?
			ORDER BY n.view_position ASC, sq.view_position ASC, cpi.tag ASC
		`, userID)
		if err != nil {
			return err
		}
		defer rows.Close()

		nodeIndexes := make(map[string]int)
		squadIndexesByNode := make(map[string]map[string]int)
		for rows.Next() {
			var nodeUUID, nodeName, countryCode, profileUUID, profileName string
			var squadUUID, squadName, inboundTag string
			if scanErr := rows.Scan(&nodeUUID, &nodeName, &countryCode, &profileUUID, &profileName, &squadUUID, &squadName, &inboundTag); scanErr != nil {
				return scanErr
			}

			nodeIndex, ok := nodeIndexes[nodeUUID]
			if !ok {
				activeNodes = append(activeNodes, userAccessibleNode{
					UUID:              nodeUUID,
					NodeName:          nodeName,
					CountryCode:       countryCode,
					ConfigProfileUUID: profileUUID,
					ConfigProfileName: profileName,
					ActiveSquads:      make([]userAccessibleSquad, 0),
				})
				nodeIndex = len(activeNodes) - 1
				nodeIndexes[nodeUUID] = nodeIndex
				squadIndexesByNode[nodeUUID] = make(map[string]int)
			}

			squadIndexes := squadIndexesByNode[nodeUUID]
			squadIndex, ok := squadIndexes[squadUUID]
			if !ok {
				activeNodes[nodeIndex].ActiveSquads = append(activeNodes[nodeIndex].ActiveSquads, userAccessibleSquad{
					SquadName:      squadName,
					ActiveInbounds: make([]string, 0),
				})
				squadIndex = len(activeNodes[nodeIndex].ActiveSquads) - 1
				squadIndexes[squadUUID] = squadIndex
			}

			activeNodes[nodeIndex].ActiveSquads[squadIndex].ActiveInbounds = append(
				activeNodes[nodeIndex].ActiveSquads[squadIndex].ActiveInbounds,
				inboundTag,
			)
		}
		return rows.Err()
	})
	if err != nil {
		if errors.Is(err, errUserNotFound) {
			shared.SendError(w, http.StatusNotFound, "user not found", nil, cfg)
			return
		}
		shared.SendError(w, http.StatusInternalServerError, "failed to fetch user accessible nodes", err, cfg)
		return
	}

	shared.WriteJSON(w, http.StatusOK, map[string]any{
		"response": map[string]any{
			"userUuid":    userUUID,
			"activeNodes": activeNodes,
		},
	})
}

func handleDeleteUser(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig, userUUID string) {
	record, recordErr := getUserRecordByUUID(r.Context(), manager, userUUID)
	if recordErr != nil && !errors.Is(recordErr, errUserNotFound) {
		shared.SendError(w, http.StatusInternalServerError, "failed to fetch user", recordErr, cfg)
		return
	}
	internalSquadNodeUUIDs := make([]string, 0)
	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		tx, err := db.BeginTx(r.Context(), nil)
		if err != nil {
			return err
		}

		var tID int64
		if err := tx.QueryRowContext(r.Context(), `SELECT t_id FROM users WHERE uuid = ?`, userUUID).Scan(&tID); err != nil {
			_ = tx.Rollback()
			if errors.Is(err, sql.ErrNoRows) {
				return errUserNotFound
			}
			return err
		}

		currentSquads, loadErr := getUserInternalSquadsTx(r.Context(), tx, tID)
		if loadErr != nil {
			_ = tx.Rollback()
			return loadErr
		}

		nodeUUIDs, nodeTargetsErr := resolveNodeUUIDsForInternalSquadsTx(r.Context(), tx, currentSquads)
		if nodeTargetsErr != nil {
			_ = tx.Rollback()
			return nodeTargetsErr
		}
		internalSquadNodeUUIDs = nodeUUIDs

		result, err := tx.ExecContext(r.Context(), `DELETE FROM users WHERE uuid = ?`, userUUID)
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
	if err != nil {
		if errors.Is(err, errUserNotFound) {
			shared.SendError(w, http.StatusNotFound, "user not found", nil, cfg)
			return
		}
		shared.SendError(w, http.StatusInternalServerError, "failed to delete user", err, cfg)
		return
	}

	if len(internalSquadNodeUUIDs) > 0 {
		monitor.RequestNodeDeploy(true, internalSquadNodeUUIDs...)
	}
	if recordErr == nil {
		emitUserNotification(r.Context(), manager, cfg, notifications.EventUserDeleted, record, nil)
	}
	shared.WriteJSON(w, http.StatusOK, map[string]any{"response": map[string]any{"isDeleted": true}})
}
