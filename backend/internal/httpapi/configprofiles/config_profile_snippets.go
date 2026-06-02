package configprofiles

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"exodus/internal/config"
	dbmanager "exodus/internal/db/manager"
	"exodus/internal/httpapi/shared"

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
	writeConfigProfileSnippetsResponse(w, r, manager, cfg, http.StatusOK)
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
	if err := validateConfigProfileSnippet(req.Snippet); err != nil {
		if errors.Is(err, errConfigProfileSnippetEmpty) {
			shared.SendError(w, http.StatusBadRequest, "Snippet cannot be empty", nil, cfg)
			return
		}
		if errors.Is(err, errConfigProfileSnippetContainsEmptyObjects) {
			shared.SendError(w, http.StatusBadRequest, "Snippet cannot contain empty objects", nil, cfg)
			return
		}
		shared.SendError(w, http.StatusBadRequest, "snippet must be valid JSON", nil, cfg)
		return
	}

	if err := createConfigProfileSnippet(r.Context(), manager, req); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			shared.SendError(w, http.StatusBadRequest, "Snippet name already exists", nil, cfg)
			return
		}
		shared.SendError(w, http.StatusInternalServerError, "Create snippet error", err, cfg)
		return
	}

	writeConfigProfileSnippetsResponse(w, r, manager, cfg, http.StatusCreated)
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
	if err := validateConfigProfileSnippet(req.Snippet); err != nil {
		if errors.Is(err, errConfigProfileSnippetEmpty) {
			shared.SendError(w, http.StatusBadRequest, "Snippet cannot be empty", nil, cfg)
			return
		}
		if errors.Is(err, errConfigProfileSnippetContainsEmptyObjects) {
			shared.SendError(w, http.StatusBadRequest, "Snippet cannot contain empty objects", nil, cfg)
			return
		}
		shared.SendError(w, http.StatusBadRequest, "snippet must be valid JSON", nil, cfg)
		return
	}

	if err := updateConfigProfileSnippet(r.Context(), manager, req); err != nil {
		if errors.Is(err, errConfigProfileSnippetNotFound) {
			shared.SendError(w, http.StatusNotFound, "Snippet not found", nil, cfg)
			return
		}
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			shared.SendError(w, http.StatusBadRequest, "Snippet name already exists", nil, cfg)
			return
		}
		shared.SendError(w, http.StatusInternalServerError, "Update snippet error", err, cfg)
		return
	}

	writeConfigProfileSnippetsResponse(w, r, manager, cfg, http.StatusOK)
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
			shared.SendError(w, http.StatusNotFound, "Snippet not found", nil, cfg)
			return
		}
		shared.SendError(w, http.StatusInternalServerError, "Delete snippet by name error", err, cfg)
		return
	}

	writeConfigProfileSnippetsResponse(w, r, manager, cfg, http.StatusOK)
}

func writeConfigProfileSnippetsResponse(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig, status int) {
	snippets, err := getConfigProfileSnippets(r.Context(), manager)
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "Get snippets error", err, cfg)
		return
	}

	shared.WriteJSON(w, status, map[string]any{
		"response": map[string]any{
			"total":    len(snippets),
			"snippets": snippets,
		},
	})
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

func createConfigProfileSnippet(ctx context.Context, manager *dbmanager.DatabaseManager, req configProfileSnippetRequest) error {
	return manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		_, err := db.ExecContext(ctx, `
			INSERT INTO config_profile_snippets (name, snippet)
			VALUES (?, ?::jsonb)`, req.Name, string(req.Snippet))
		return err
	})
}

func updateConfigProfileSnippet(ctx context.Context, manager *dbmanager.DatabaseManager, req configProfileSnippetRequest) error {
	return manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		res, err := db.ExecContext(ctx, `
			UPDATE config_profile_snippets
			SET snippet = ?::jsonb
			WHERE name = ?`, string(req.Snippet), req.Name)
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

var (
	errConfigProfileSnippetEmpty                = errors.New("config profile snippet is empty")
	errConfigProfileSnippetContainsEmptyObjects = errors.New("config profile snippet contains empty objects")
	errConfigProfileSnippetInvalidJSON          = errors.New("config profile snippet is invalid json")
)

func validateConfigProfileSnippet(raw json.RawMessage) error {
	if len(raw) == 0 {
		return errConfigProfileSnippetEmpty
	}
	if !json.Valid(raw) {
		return errConfigProfileSnippetInvalidJSON
	}

	var items []json.RawMessage
	if err := json.Unmarshal(raw, &items); err != nil {
		return errConfigProfileSnippetEmpty
	}
	if len(items) == 0 {
		return errConfigProfileSnippetEmpty
	}

	for _, item := range items {
		var object map[string]any
		if err := json.Unmarshal(item, &object); err != nil {
			return errConfigProfileSnippetContainsEmptyObjects
		}
		if len(object) == 0 {
			return errConfigProfileSnippetContainsEmptyObjects
		}
	}

	return nil
}
