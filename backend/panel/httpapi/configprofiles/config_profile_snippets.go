package configprofiles

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"v2ray-stat/backend/panel/config"
	dbmanager "v2ray-stat/backend/panel/db/manager"
	"v2ray-stat/backend/panel/httpapi/shared"

	"github.com/jackc/pgx/v5/pgconn"
)

type ConfigProfileSnippet struct {
	Name      string          `json:"name"`
	Snippet   json.RawMessage `json:"snippet"`
	CreatedAt time.Time       `json:"createdAt"`
}

type configProfileSnippetRequest struct {
	Name    string          `json:"name"`
	Snippet json.RawMessage `json:"snippet"`
}

var errConfigProfileSnippetNotFound = errors.New("config profile snippet not found")

func ConfigProfileSnippetsHandler(manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handleGetConfigProfileSnippets(w, r, manager, cfg)
		case http.MethodPost:
			handleCreateConfigProfileSnippet(w, r, manager, cfg)
		case http.MethodPatch:
			handleUpdateConfigProfileSnippet(w, r, manager, cfg)
		case http.MethodDelete:
			handleDeleteConfigProfileSnippet(w, r, manager, cfg)
		default:
			shared.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
		}
	}
}

func handleGetConfigProfileSnippets(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) {
	snippets, err := getConfigProfileSnippets(r.Context(), manager)
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to fetch config profile snippets", err, cfg)
		return
	}

	shared.WriteJSON(w, http.StatusOK, map[string]any{
		"response": map[string]any{
			"total":    len(snippets),
			"snippets": snippets,
		},
	})
}

func handleCreateConfigProfileSnippet(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) {
	var req configProfileSnippetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		shared.SendError(w, http.StatusBadRequest, "invalid request body", err, cfg)
		return
	}

	req.Name = strings.TrimSpace(req.Name)
	if req.Name == "" {
		shared.SendError(w, http.StatusBadRequest, "name is required", nil, cfg)
		return
	}
	if len(req.Name) > 255 {
		shared.SendError(w, http.StatusBadRequest, "name must be 255 characters or less", nil, cfg)
		return
	}
	if len(req.Snippet) == 0 || !json.Valid(req.Snippet) {
		shared.SendError(w, http.StatusBadRequest, "snippet must be valid JSON", nil, cfg)
		return
	}

	snippet, err := createConfigProfileSnippet(r.Context(), manager, req)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			shared.SendError(w, http.StatusConflict, "snippet with this name already exists", nil, cfg)
			return
		}
		shared.SendError(w, http.StatusInternalServerError, "failed to create config profile snippet", err, cfg)
		return
	}

	shared.WriteJSON(w, http.StatusCreated, map[string]any{"response": snippet})
}

func handleUpdateConfigProfileSnippet(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) {
	var req configProfileSnippetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		shared.SendError(w, http.StatusBadRequest, "invalid request body", err, cfg)
		return
	}

	req.Name = strings.TrimSpace(req.Name)
	if req.Name == "" {
		shared.SendError(w, http.StatusBadRequest, "name is required", nil, cfg)
		return
	}
	if len(req.Snippet) == 0 || !json.Valid(req.Snippet) {
		shared.SendError(w, http.StatusBadRequest, "snippet must be valid JSON", nil, cfg)
		return
	}

	snippet, err := updateConfigProfileSnippet(r.Context(), manager, req)
	if err != nil {
		if errors.Is(err, errConfigProfileSnippetNotFound) {
			shared.SendError(w, http.StatusNotFound, "snippet not found", nil, cfg)
			return
		}
		shared.SendError(w, http.StatusInternalServerError, "failed to update config profile snippet", err, cfg)
		return
	}

	shared.WriteJSON(w, http.StatusOK, map[string]any{"response": snippet})
}

func handleDeleteConfigProfileSnippet(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) {
	var req struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		shared.SendError(w, http.StatusBadRequest, "invalid request body", err, cfg)
		return
	}

	req.Name = strings.TrimSpace(req.Name)
	if req.Name == "" {
		shared.SendError(w, http.StatusBadRequest, "name is required", nil, cfg)
		return
	}

	if err := deleteConfigProfileSnippet(r.Context(), manager, req.Name); err != nil {
		if errors.Is(err, errConfigProfileSnippetNotFound) {
			shared.SendError(w, http.StatusNotFound, "snippet not found", nil, cfg)
			return
		}
		shared.SendError(w, http.StatusInternalServerError, "failed to delete config profile snippet", err, cfg)
		return
	}

	shared.WriteJSON(w, http.StatusOK, map[string]any{"response": map[string]any{"isDeleted": true}})
}

func getConfigProfileSnippets(ctx context.Context, manager *dbmanager.DatabaseManager) ([]ConfigProfileSnippet, error) {
	var snippets []ConfigProfileSnippet
	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		rows, err := db.QueryContext(ctx, `
			SELECT name, snippet, created_at
			FROM config_profile_snippets
			ORDER BY name ASC`)
		if err != nil {
			return err
		}
		defer rows.Close()

		result := make([]ConfigProfileSnippet, 0)
		for rows.Next() {
			var snippet ConfigProfileSnippet
			if scanErr := rows.Scan(&snippet.Name, &snippet.Snippet, &snippet.CreatedAt); scanErr != nil {
				return scanErr
			}
			result = append(result, snippet)
		}
		if err := rows.Err(); err != nil {
			return err
		}
		snippets = result
		return nil
	})
	if err != nil {
		return nil, err
	}
	return snippets, nil
}

func createConfigProfileSnippet(ctx context.Context, manager *dbmanager.DatabaseManager, req configProfileSnippetRequest) (ConfigProfileSnippet, error) {
	var snippet ConfigProfileSnippet
	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		row := db.QueryRowContext(ctx, `
			INSERT INTO config_profile_snippets (name, snippet)
			VALUES (?, ?::jsonb)
			RETURNING name, snippet, created_at`, req.Name, string(req.Snippet))
		return row.Scan(&snippet.Name, &snippet.Snippet, &snippet.CreatedAt)
	})
	if err != nil {
		return ConfigProfileSnippet{}, err
	}
	return snippet, nil
}

func updateConfigProfileSnippet(ctx context.Context, manager *dbmanager.DatabaseManager, req configProfileSnippetRequest) (ConfigProfileSnippet, error) {
	var snippet ConfigProfileSnippet
	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		row := db.QueryRowContext(ctx, `
			UPDATE config_profile_snippets
			SET snippet = ?::jsonb
			WHERE name = ?
			RETURNING name, snippet, created_at`, string(req.Snippet), req.Name)
		if err := row.Scan(&snippet.Name, &snippet.Snippet, &snippet.CreatedAt); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return errConfigProfileSnippetNotFound
			}
			return err
		}
		return nil
	})
	if err != nil {
		return ConfigProfileSnippet{}, err
	}
	return snippet, nil
}

func deleteConfigProfileSnippet(ctx context.Context, manager *dbmanager.DatabaseManager, name string) error {
	return manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		res, err := db.ExecContext(ctx, `DELETE FROM config_profile_snippets WHERE name = ?`, name)
		if err != nil {
			return err
		}
		affected, err := res.RowsAffected()
		if err != nil {
			return err
		}
		if affected == 0 {
			return errConfigProfileSnippetNotFound
		}
		return nil
	})
}
