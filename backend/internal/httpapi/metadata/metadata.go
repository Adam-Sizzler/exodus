package metadata

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"exodus/internal/config"
	"exodus/internal/httpapi/shared"

	"github.com/google/uuid"
)

type metadataRequest struct {
	Metadata map[string]any `json:"metadata"`
}

func UserHandler(db *sql.DB, cfg *config.BackendConfig) http.HandlerFunc {
	return entityHandler(db, cfg, "user")
}

func NodeHandler(db *sql.DB, cfg *config.BackendConfig) http.HandlerFunc {
	return entityHandler(db, cfg, "node")
}

func entityHandler(db *sql.DB, cfg *config.BackendConfig, entity string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pathParam := strings.Trim(strings.TrimPrefix(r.URL.Path, "/api/metadata/"+entity+"/"), "/")
		if entity == "node" {
			if _, err := uuid.Parse(pathParam); err != nil {
				shared.SendError(w, http.StatusBadRequest, "invalid uuid format", nil, cfg)
				return
			}
		} else if entity == "user" {
			if _, err := strconv.ParseInt(pathParam, 10, 64); err != nil {
				shared.SendError(w, http.StatusBadRequest, "invalid userId format", nil, cfg)
				return
			}
		}

		switch r.Method {
		case http.MethodGet:
			metadata, err := getMetadata(r, db, entity, pathParam)
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
			metadata, err := upsertMetadata(r, db, entity, pathParam, req.Metadata)
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

func getMetadata(r *http.Request, db *sql.DB, entity, pathParam string) (map[string]any, error) {
	metadata := map[string]any{}
	table, column, id, err := metadataTargetID(r, db, entity, pathParam)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return metadata, nil
		}
		return nil, err
	}
	var raw string
	err = db.QueryRowContext(r.Context(), "SELECT metadata::text FROM "+table+" WHERE "+column+" = $1", id).Scan(&raw)
	if errors.Is(err, sql.ErrNoRows) {
		return metadata, nil
	}
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal([]byte(raw), &metadata); err != nil {
		return nil, err
	}
	return metadata, nil
}

func upsertMetadata(r *http.Request, db *sql.DB, entity, pathParam string, metadata map[string]any) (map[string]any, error) {
	raw, err := json.Marshal(metadata)
	if err != nil {
		return nil, err
	}
	table, column, id, err := metadataTargetID(r, db, entity, pathParam)
	if err != nil {
		return nil, err
	}
	_, err = db.ExecContext(r.Context(), `
		INSERT INTO `+table+` (`+column+`, metadata)
		VALUES ($1, $2::jsonb)
		ON CONFLICT (`+column+`) DO UPDATE
		SET metadata = EXCLUDED.metadata
	`, id, string(raw))
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

func metadataTargetID(r *http.Request, db *sql.DB, entity, pathParam string) (table string, column string, id int64, err error) {
	table, column = metadataTable(entity)
	if entity == "node" {
		err = db.QueryRowContext(r.Context(), `SELECT id FROM nodes WHERE uuid = $1`, pathParam).Scan(&id)
		return table, column, id, err
	}
	id, err = strconv.ParseInt(pathParam, 10, 64)
	return table, column, id, err
}
