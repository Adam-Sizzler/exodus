package metadata

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"exodus/internal/config"
	dbmanager "exodus/internal/db/manager"
	"exodus/internal/httpapi/shared"

	"github.com/google/uuid"
)

type metadataRequest struct {
	Metadata map[string]any `json:"metadata"`
}

func UserHandler(manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) http.HandlerFunc {
	return entityHandler(manager, cfg, "user")
}

func NodeHandler(manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) http.HandlerFunc {
	return entityHandler(manager, cfg, "node")
}

func entityHandler(manager *dbmanager.DatabaseManager, cfg *config.BackendConfig, entity string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		entityUUID := strings.Trim(strings.TrimPrefix(r.URL.Path, "/api/metadata/"+entity+"/"), "/")
		if _, err := uuid.Parse(entityUUID); err != nil {
			shared.SendError(w, http.StatusBadRequest, "invalid uuid format", nil, cfg)
			return
		}

		switch r.Method {
		case http.MethodGet:
			metadata, err := getMetadata(r, manager, entity, entityUUID)
			if err != nil {
				shared.SendError(w, http.StatusInternalServerError, "failed to fetch metadata", err, cfg)
				return
			}
			shared.WriteJSON(w, http.StatusOK, map[string]any{"response": map[string]any{"metadata": metadata}})
		case http.MethodPut:
			var req metadataRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				shared.SendError(w, http.StatusBadRequest, "invalid JSON", err, cfg)
				return
			}
			if req.Metadata == nil {
				req.Metadata = map[string]any{}
			}
			metadata, err := upsertMetadata(r, manager, entity, entityUUID, req.Metadata)
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					shared.SendError(w, http.StatusNotFound, entity+" not found", nil, cfg)
					return
				}
				shared.SendError(w, http.StatusInternalServerError, "failed to update metadata", err, cfg)
				return
			}
			shared.WriteJSON(w, http.StatusOK, map[string]any{"response": map[string]any{"metadata": metadata}})
		default:
			shared.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
		}
	}
}

func getMetadata(r *http.Request, manager *dbmanager.DatabaseManager, entity, entityUUID string) (map[string]any, error) {
	metadata := map[string]any{}
	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		table, column, id, err := metadataTargetID(r, db, entity, entityUUID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil
			}
			return err
		}
		var raw string
		err = db.QueryRowContext(r.Context(), "SELECT metadata::text FROM "+table+" WHERE "+column+" = ?", id).Scan(&raw)
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}
		if err != nil {
			return err
		}
		if err := json.Unmarshal([]byte(raw), &metadata); err != nil {
			return err
		}
		return nil
	})
	return metadata, err
}

func upsertMetadata(r *http.Request, manager *dbmanager.DatabaseManager, entity, entityUUID string, metadata map[string]any) (map[string]any, error) {
	raw, err := json.Marshal(metadata)
	if err != nil {
		return nil, err
	}
	err = manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		table, column, id, err := metadataTargetID(r, db, entity, entityUUID)
		if err != nil {
			return err
		}
		_, err = db.ExecContext(r.Context(), `
			INSERT INTO `+table+` (`+column+`, metadata)
			VALUES (?, ?::jsonb)
			ON CONFLICT (`+column+`) DO UPDATE
			SET metadata = EXCLUDED.metadata
		`, id, string(raw))
		return err
	})
	if err != nil {
		return nil, err
	}
	return metadata, nil
}

func metadataTable(entity string) (table string, column string) {
	if entity == "node" {
		return "node_meta", "node_id"
	}
	return "user_meta", "user_id"
}

func metadataTargetID(r *http.Request, db dbmanager.DBExecutor, entity, entityUUID string) (table string, column string, id int64, err error) {
	table, column = metadataTable(entity)
	if entity == "node" {
		err = db.QueryRowContext(r.Context(), `SELECT id FROM nodes WHERE uuid = ?`, entityUUID).Scan(&id)
		return table, column, id, err
	}
	err = db.QueryRowContext(r.Context(), `SELECT t_id FROM users WHERE uuid = ?`, entityUUID).Scan(&id)
	return table, column, id, err
}
