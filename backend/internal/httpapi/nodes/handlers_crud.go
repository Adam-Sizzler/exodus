package nodes

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
	"exodus/internal/security"

	"github.com/google/uuid"
)

func handleGetNodes(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) {
	nodes, err := getAllNodeRecords(r.Context(), manager)
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to fetch nodes", err, cfg)
		return
	}

	response, err := buildNodeResponses(r.Context(), manager, cfg, nodes)
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to build node response", err, cfg)
		return
	}

	shared.WriteJSON(w, http.StatusOK, map[string]any{"response": response})
}

func handleGetNode(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig, nodeUUID string) {
	node, err := getNodeByUUID(r.Context(), manager, nodeUUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			shared.SendError(w, http.StatusNotFound, "node not found", nil, cfg)
			return
		}
		shared.SendError(w, http.StatusInternalServerError, "failed to fetch node", err, cfg)
		return
	}

	response, err := buildNodeResponses(r.Context(), manager, cfg, []nodeRecord{node})
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to build node response", err, cfg)
		return
	}

	shared.WriteJSON(w, http.StatusOK, map[string]any{"response": response[0]})
}

func handleCreateNode(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) {
	var req createNodeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		shared.SendError(w, http.StatusBadRequest, "invalid JSON", err, cfg)
		return
	}
	if err := validateCreateRequest(req); err != nil {
		shared.SendError(w, http.StatusBadRequest, err.Error(), nil, cfg)
		return
	}
	if err := ensureConfigProfileInbounds(r.Context(), manager, req.ConfigProfile.ActiveConfigProfileUUID, req.ConfigProfile.ActiveInbounds); err != nil {
		handleConfigProfileValidationError(w, err, cfg)
		return
	}
	if req.ProviderUUID != nil && strings.TrimSpace(*req.ProviderUUID) != "" {
		if _, err := uuid.Parse(strings.TrimSpace(*req.ProviderUUID)); err != nil {
			shared.SendError(w, http.StatusBadRequest, "invalid providerUuid", nil, cfg)
			return
		}
	}
	grpcAuthToken := ""
	if req.GRPCAuthToken != nil {
		grpcAuthToken = *req.GRPCAuthToken
	}
	grpcAuthToken, err := security.ResolveGRPCAuthToken(grpcAuthToken)
	if err != nil {
		shared.SendError(w, http.StatusBadRequest, err.Error(), nil, cfg)
		return
	}

	nodeUUID := uuid.NewString()
	now := time.Now().UTC()
	err = manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		tx, err := db.BeginTx(r.Context(), nil)
		if err != nil {
			return err
		}

		_, err = tx.ExecContext(r.Context(), `
			INSERT INTO nodes (
				uuid, name, address, port, proxy_url, api_schema, api_path, grpc_auth_token, active_config_profile_uuid, active_plugin_uuid,
				is_connected, is_connecting, is_disabled, last_status_change, last_status_message,
				consumption_multiplier, node_consumption_multiplier,
				is_traffic_tracking_active, traffic_reset_day, traffic_limit_bytes, traffic_used_bytes,
				notify_percent, provider_uuid, country_code, tags, note, created_at, updated_at
			) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		`,
			nodeUUID,
			strings.TrimSpace(req.Name),
			strings.TrimSpace(req.Address),
			req.Port,
			normalizeNullableString(req.ProxyURL),
			normalizeAPISchema(req.APISchema),
			normalizeAPIPath(req.APIPath),
			grpcAuthToken,
			req.ConfigProfile.ActiveConfigProfileUUID,
			normalizeNullableString(req.ActivePluginUUID),
			false,
			false,
			false,
			nil,
			nil,
			toNanoMultiplier(coalesceFloat(req.ConsumptionMultiplier, 1)),
			toNanoMultiplier(coalesceFloat(req.NodeConsumptionMultiplier, 1)),
			coalesceBool(req.IsTrafficTrackingActive, false),
			coalesceInt(req.TrafficResetDay, 1),
			coalesceInt64(req.TrafficLimitBytes, 0),
			0,
			coalesceInt(req.NotifyPercent, 0),
			normalizeNullableString(req.ProviderUUID),
			normalizeCountryCode(req.CountryCode),
			normalizeTags(req.Tags),
			normalizeNullableString(req.Note),
			now,
			now,
		)
		if err != nil {
			_ = tx.Rollback()
			errStr := err.Error()
			if strings.Contains(errStr, "nodes_name_key") {
				return fmt.Errorf("node with this name already exists")
			}
			if strings.Contains(errStr, "nodes_address_key") {
				return fmt.Errorf("node with this address already exists")
			}
			return err
		}

		if err := replaceNodeInboundsTx(r.Context(), tx, nodeUUID, req.ConfigProfile.ActiveInbounds); err != nil {
			_ = tx.Rollback()
			return err
		}

		return tx.Commit()
	})
	if err != nil {
		errStr := err.Error()
		if strings.Contains(errStr, "already exists") {
			shared.SendError(w, http.StatusBadRequest, errStr, err, cfg)
		} else {
			shared.SendError(w, http.StatusInternalServerError, "failed to create node", err, cfg)
		}
		return
	}

	monitor.RequestNodeSync()
	monitor.RequestNodeDeploy(true, nodeUUID)
	node, err := getNodeByUUID(r.Context(), manager, nodeUUID)
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to fetch created node", err, cfg)
		return
	}
	response, err := buildNodeResponses(r.Context(), manager, cfg, []nodeRecord{node})
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to build node response", err, cfg)
		return
	}
	emitNodeNotification(r.Context(), cfg, notifications.EventNodeCreated, node, nil)
	shared.WriteJSON(w, http.StatusCreated, map[string]any{"response": response[0]})
}

func handleUpdateNode(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) {
	var req updateNodeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		shared.SendError(w, http.StatusBadRequest, "invalid JSON", err, cfg)
		return
	}
	if _, err := uuid.Parse(req.UUID); err != nil {
		shared.SendError(w, http.StatusBadRequest, "invalid UUID format", nil, cfg)
		return
	}
	if err := validateUpdateRequest(req); err != nil {
		shared.SendError(w, http.StatusBadRequest, err.Error(), nil, cfg)
		return
	}
	if req.ConfigProfile != nil {
		if err := ensureConfigProfileInbounds(r.Context(), manager, req.ConfigProfile.ActiveConfigProfileUUID, req.ConfigProfile.ActiveInbounds); err != nil {
			handleConfigProfileValidationError(w, err, cfg)
			return
		}
	}
	if req.ProviderUUID.Set && req.ProviderUUID.Value != nil && strings.TrimSpace(*req.ProviderUUID.Value) != "" {
		if _, err := uuid.Parse(strings.TrimSpace(*req.ProviderUUID.Value)); err != nil {
			shared.SendError(w, http.StatusBadRequest, "invalid providerUuid", nil, cfg)
			return
		}
	}
	var grpcAuthToken *string
	if req.GRPCAuthToken != nil {
		token, err := security.ResolveGRPCAuthToken(*req.GRPCAuthToken)
		if err != nil {
			shared.SendError(w, http.StatusBadRequest, err.Error(), nil, cfg)
			return
		}
		grpcAuthToken = &token
	}

	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
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

		if req.Name != nil {
			add("name", strings.TrimSpace(*req.Name))
		}
		if req.Address != nil {
			add("address", strings.TrimSpace(*req.Address))
		}
		if req.Port != nil {
			add("port", *req.Port)
		}
		if req.ProxyURL.Set {
			if req.ProxyURL.Value == nil || strings.TrimSpace(*req.ProxyURL.Value) == "" {
				clauses = append(clauses, "proxy_url = NULL")
			} else {
				add("proxy_url", strings.TrimSpace(*req.ProxyURL.Value))
			}
		}
		if req.APISchema != nil {
			add("api_schema", normalizeAPISchema(req.APISchema))
		}
		if req.APIPath != nil {
			add("api_path", normalizeAPIPath(req.APIPath))
		}
		if grpcAuthToken != nil {
			add("grpc_auth_token", *grpcAuthToken)
		}
		if req.IsTrafficTrackingActive != nil {
			add("is_traffic_tracking_active", *req.IsTrafficTrackingActive)
		}
		if req.TrafficLimitBytes != nil {
			add("traffic_limit_bytes", *req.TrafficLimitBytes)
		}
		if req.NotifyPercent != nil {
			add("notify_percent", *req.NotifyPercent)
		}
		if req.TrafficResetDay != nil {
			add("traffic_reset_day", *req.TrafficResetDay)
		}
		if req.CountryCode != nil {
			add("country_code", strings.ToUpper(strings.TrimSpace(*req.CountryCode)))
		}
		if req.ConsumptionMultiplier != nil {
			add("consumption_multiplier", toNanoMultiplier(*req.ConsumptionMultiplier))
		}
		if req.NodeConsumptionMultiplier != nil {
			add("node_consumption_multiplier", toNanoMultiplier(*req.NodeConsumptionMultiplier))
		}
		if req.Tags != nil {
			add("tags", normalizeTags(*req.Tags))
		}
		if req.Note.Set {
			if req.Note.Value == nil || strings.TrimSpace(*req.Note.Value) == "" {
				clauses = append(clauses, "note = NULL")
			} else {
				add("note", strings.TrimSpace(*req.Note.Value))
			}
		}
		if req.ProviderUUID.Set {
			if req.ProviderUUID.Value == nil || strings.TrimSpace(*req.ProviderUUID.Value) == "" {
				clauses = append(clauses, "provider_uuid = NULL")
			} else {
				add("provider_uuid", strings.TrimSpace(*req.ProviderUUID.Value))
			}
		}
		if req.ActivePluginUUID.Set {
			if req.ActivePluginUUID.Value == nil || strings.TrimSpace(*req.ActivePluginUUID.Value) == "" {
				clauses = append(clauses, "active_plugin_uuid = NULL")
			} else {
				add("active_plugin_uuid", strings.TrimSpace(*req.ActivePluginUUID.Value))
			}
		}
		if req.ConfigProfile != nil {
			add("active_config_profile_uuid", req.ConfigProfile.ActiveConfigProfileUUID)
		}

		if len(clauses) > 0 {
			args = append(args, req.UUID)
			query := fmt.Sprintf("UPDATE nodes SET %s, updated_at = CURRENT_TIMESTAMP WHERE uuid = ?", strings.Join(clauses, ", "))
			result, err := tx.ExecContext(r.Context(), query, args...)
			if err != nil {
				_ = tx.Rollback()
				errStr := err.Error()
				if strings.Contains(errStr, "nodes_name_key") {
					return fmt.Errorf("node with this name already exists")
				}
				if strings.Contains(errStr, "nodes_address_key") {
					return fmt.Errorf("node with this address already exists")
				}
				return err
			}
			rows, err := result.RowsAffected()
			if err != nil {
				_ = tx.Rollback()
				return err
			}
			if rows == 0 {
				_ = tx.Rollback()
				return sql.ErrNoRows
			}
		}

		if req.ConfigProfile != nil {
			if err := replaceNodeInboundsTx(r.Context(), tx, req.UUID, req.ConfigProfile.ActiveInbounds); err != nil {
				_ = tx.Rollback()
				return err
			}
		}

		return tx.Commit()
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			shared.SendError(w, http.StatusNotFound, "node not found", nil, cfg)
			return
		}
		errStr := err.Error()
		if strings.Contains(errStr, "already exists") {
			shared.SendError(w, http.StatusBadRequest, errStr, err, cfg)
		} else {
			shared.SendError(w, http.StatusInternalServerError, "failed to update node", err, cfg)
		}
		return
	}

	monitor.RequestNodeSync()
	monitor.RequestNodeDeploy(true, req.UUID)
	node, err := getNodeByUUID(r.Context(), manager, req.UUID)
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to fetch updated node", err, cfg)
		return
	}
	response, err := buildNodeResponses(r.Context(), manager, cfg, []nodeRecord{node})
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to build node response", err, cfg)
		return
	}
	emitNodeNotification(r.Context(), cfg, notifications.EventNodeModified, node, nil)
	shared.WriteJSON(w, http.StatusOK, map[string]any{"response": response[0]})
}

func handleDeleteNode(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig, nodeUUID string) {
	node, nodeErr := getNodeByUUID(r.Context(), manager, nodeUUID)
	if nodeErr != nil && !errors.Is(nodeErr, sql.ErrNoRows) {
		shared.SendError(w, http.StatusInternalServerError, "failed to fetch node", nodeErr, cfg)
		return
	}
	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		result, err := db.ExecContext(r.Context(), `DELETE FROM nodes WHERE uuid = ?`, nodeUUID)
		if err != nil {
			return err
		}
		rows, err := result.RowsAffected()
		if err != nil {
			return err
		}
		if rows == 0 {
			return sql.ErrNoRows
		}
		return nil
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			shared.SendError(w, http.StatusNotFound, "node not found", nil, cfg)
			return
		}
		shared.SendError(w, http.StatusInternalServerError, "failed to delete node", err, cfg)
		return
	}

	monitor.RequestNodeSync()
	monitor.RequestNodeDeploy(true, nodeUUID)
	if nodeErr == nil {
		emitNodeNotification(r.Context(), cfg, notifications.EventNodeDeleted, node, nil)
	}
	shared.WriteJSON(w, http.StatusOK, map[string]any{"response": map[string]any{"isDeleted": true}})
}