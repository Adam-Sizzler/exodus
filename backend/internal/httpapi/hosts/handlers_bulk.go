package hosts

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"exodus/internal/config"
	dbmanager "exodus/internal/db/manager"
	"exodus/internal/httpapi/shared"

	"github.com/google/uuid"
)

func handleBulkUpdateHosts(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) {
	var req HostBulkUpdateRequestAPI
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		shared.SendError(w, http.StatusBadRequest, "invalid JSON", err, cfg)
		return
	}
	if len(req.Uuids) == 0 {
		shared.SendError(w, http.StatusBadRequest, "uuids cannot be empty", nil, cfg)
		return
	}
	for _, hostUUID := range req.Uuids {
		if _, err := uuid.Parse(hostUUID); err != nil {
			shared.SendError(w, http.StatusBadRequest, "invalid UUID format", nil, cfg)
			return
		}
	}
	if err := validateUpdateRequest(req.hostUpdateFields); err != nil {
		shared.SendError(w, http.StatusBadRequest, err.Error(), nil, cfg)
		return
	}
	if req.Inbound != nil {
		if err := ensureConfigProfileInbound(r.Context(), manager, *req.Inbound.ConfigProfileUUID, *req.Inbound.ConfigProfileInboundUUID); err != nil {
			switch {
			case errors.Is(err, errConfigProfileNotFound):
				shared.SendError(w, http.StatusBadRequest, "config profile not found", nil, cfg)
				return
			case errors.Is(err, errConfigProfileInboundNotFound):
				shared.SendError(w, http.StatusBadRequest, "config profile inbound not found in specified profile", nil, cfg)
				return
			default:
				shared.SendError(w, http.StatusInternalServerError, "failed to validate config profile inbound", err, cfg)
				return
			}
		}
	}
	if req.XrayJSONTemplateUUID.Set && req.XrayJSONTemplateUUID.Value != nil && *req.XrayJSONTemplateUUID.Value != "" {
		if err := ensureXrayJSONTemplate(r.Context(), manager, *req.XrayJSONTemplateUUID.Value); err != nil {
			switch {
			case errors.Is(err, errTemplateNotFound):
				shared.SendError(w, http.StatusBadRequest, "subscription template not found", nil, cfg)
				return
			case errors.Is(err, errTemplateTypeNotAllowed):
				shared.SendError(w, http.StatusBadRequest, "template type not allowed", nil, cfg)
				return
			default:
				shared.SendError(w, http.StatusInternalServerError, "failed to validate subscription template", err, cfg)
				return
			}
		}
	}

	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		tx, err := db.BeginTx(r.Context(), nil)
		if err != nil {
			return err
		}

		clauses, args, err := buildHostUpdateClauses(req.hostUpdateFields)
		if err != nil {
			_ = tx.Rollback()
			return err
		}

		if len(clauses) > 0 {
			args = append(args, req.Uuids)
			query := fmt.Sprintf("UPDATE hosts SET %s WHERE uuid = ANY(?)", strings.Join(clauses, ", "))
			if _, err := tx.ExecContext(r.Context(), query, args...); err != nil {
				_ = tx.Rollback()
				return err
			}
		}

		if req.Nodes != nil {
			for _, hostUUID := range req.Uuids {
				if err := replaceHostNodesTx(r.Context(), tx, hostUUID, req.Nodes); err != nil {
					_ = tx.Rollback()
					return err
				}
			}
		}
		if req.ExcludedInternalSquads != nil {
			for _, hostUUID := range req.Uuids {
				if err := replaceHostExcludedSquadsTx(r.Context(), tx, hostUUID, req.ExcludedInternalSquads); err != nil {
					_ = tx.Rollback()
					return err
				}
			}
		}

		return tx.Commit()
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			shared.SendError(w, http.StatusNotFound, "host not found", nil, cfg)
			return
		}
		shared.SendError(w, http.StatusInternalServerError, "failed to update hosts", err, cfg)
		return
	}

	handleGetHosts(w, r, manager, cfg)
}

func handleReorderHosts(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) {
	var req reorderHostsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		shared.SendError(w, http.StatusBadRequest, "invalid JSON", err, cfg)
		return
	}
	if len(req.Hosts) == 0 {
		shared.SendError(w, http.StatusBadRequest, "hosts cannot be empty", nil, cfg)
		return
	}
	for _, item := range req.Hosts {
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
		for _, item := range req.Hosts {
			if _, err := tx.ExecContext(r.Context(), `UPDATE hosts SET view_position = ? WHERE uuid = ?`, item.ViewPosition, item.UUID); err != nil {
				_ = tx.Rollback()
				return err
			}
		}
		if _, err := tx.ExecContext(r.Context(), `SELECT setval('hosts_view_position_seq', (SELECT COALESCE(MAX(view_position), 0) FROM hosts) + 1)`); err != nil {
			_ = tx.Rollback()
			return err
		}
		return tx.Commit()
	})
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to reorder hosts", err, cfg)
		return
	}

	shared.WriteJSON(w, http.StatusOK, map[string]any{"response": map[string]any{"isUpdated": true}})
}

func handleBulkEnableHosts(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) {
	bulkUpdateHostsEnabled(w, r, manager, cfg, true)
}

func handleBulkDisableHosts(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) {
	bulkUpdateHostsEnabled(w, r, manager, cfg, false)
}

func handleBulkDeleteHosts(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) {
	var req bulkUUIDsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		shared.SendError(w, http.StatusBadRequest, "invalid JSON", err, cfg)
		return
	}
	if len(req.UUIDs) == 0 {
		shared.SendError(w, http.StatusBadRequest, "uuids cannot be empty", nil, cfg)
		return
	}

	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		_, err := db.ExecContext(r.Context(), `DELETE FROM hosts WHERE uuid = ANY(?)`, req.UUIDs)
		return err
	})
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to delete hosts", err, cfg)
		return
	}

	handleGetHosts(w, r, manager, cfg)
}

func handleBulkSetInbound(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) {
	var req setInboundRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		shared.SendError(w, http.StatusBadRequest, "invalid JSON", err, cfg)
		return
	}
	if len(req.UUIDs) == 0 {
		shared.SendError(w, http.StatusBadRequest, "uuids cannot be empty", nil, cfg)
		return
	}
	if _, err := uuid.Parse(req.ConfigProfileUUID); err != nil {
		shared.SendError(w, http.StatusBadRequest, "invalid configProfileUuid", nil, cfg)
		return
	}
	if _, err := uuid.Parse(req.ConfigProfileInboundUUID); err != nil {
		shared.SendError(w, http.StatusBadRequest, "invalid configProfileInboundUuid", nil, cfg)
		return
	}
	if err := ensureConfigProfileInbound(r.Context(), manager, req.ConfigProfileUUID, req.ConfigProfileInboundUUID); err != nil {
		switch {
		case errors.Is(err, errConfigProfileNotFound):
			shared.SendError(w, http.StatusBadRequest, "config profile not found", nil, cfg)
			return
		case errors.Is(err, errConfigProfileInboundNotFound):
			shared.SendError(w, http.StatusBadRequest, "config profile inbound not found in specified profile", nil, cfg)
			return
		default:
			shared.SendError(w, http.StatusInternalServerError, "failed to validate config profile inbound", err, cfg)
			return
		}
	}

	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		_, err := db.ExecContext(r.Context(), `
            UPDATE hosts
            SET config_profile_uuid = ?, config_profile_inbound_uuid = ?
            WHERE uuid = ANY(?)
        `, req.ConfigProfileUUID, req.ConfigProfileInboundUUID, req.UUIDs)
		return err
	})
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to set inbound", err, cfg)
		return
	}

	handleGetHosts(w, r, manager, cfg)
}

func handleBulkSetPort(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) {
	var req setPortRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		shared.SendError(w, http.StatusBadRequest, "invalid JSON", err, cfg)
		return
	}
	if len(req.UUIDs) == 0 {
		shared.SendError(w, http.StatusBadRequest, "uuids cannot be empty", nil, cfg)
		return
	}
	if req.Port < 1 || req.Port > 65535 {
		shared.SendError(w, http.StatusBadRequest, "invalid port", nil, cfg)
		return
	}

	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		_, err := db.ExecContext(r.Context(), `UPDATE hosts SET port = ? WHERE uuid = ANY(?)`, req.Port, req.UUIDs)
		return err
	})
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to set port", err, cfg)
		return
	}

	handleGetHosts(w, r, manager, cfg)
}

func bulkUpdateHostsEnabled(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig, enabled bool) {
	var req bulkUUIDsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		shared.SendError(w, http.StatusBadRequest, "invalid JSON", err, cfg)
		return
	}
	if len(req.UUIDs) == 0 {
		shared.SendError(w, http.StatusBadRequest, "uuids cannot be empty", nil, cfg)
		return
	}

	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		_, err := db.ExecContext(r.Context(), `UPDATE hosts SET is_disabled = ? WHERE uuid = ANY(?)`, !enabled, req.UUIDs)
		return err
	})
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to update hosts", err, cfg)
		return
	}

	handleGetHosts(w, r, manager, cfg)
}
