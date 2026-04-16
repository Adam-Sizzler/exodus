package srslists

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"exodus/backend/config"
	dbmanager "exodus/backend/db/manager"
	"exodus/backend/httpapi/shared"
	monitor "exodus/backend/nodes"
	srscore "exodus/backend/srslists"

	"github.com/google/uuid"
)

type srsListAPI struct {
	UUID           string     `json:"uuid"`
	Tag            string     `json:"tag"`
	Format         string     `json:"format"`
	URL            string     `json:"url"`
	UpdateInterval string     `json:"updateInterval"`
	Path           *string    `json:"path"`
	FileName       string     `json:"fileName"`
	ShortName      string     `json:"shortName"`
	ViewPosition   int        `json:"viewPosition"`
	IsEnabled      bool       `json:"isEnabled"`
	IsAvailable    bool       `json:"isAvailable"`
	LastCheckedAt  *time.Time `json:"lastCheckedAt"`
	LastError      *string    `json:"lastError"`
	CreatedAt      time.Time  `json:"createdAt"`
	UpdatedAt      time.Time  `json:"updatedAt"`
}

type createSRSListsRequest struct {
	URL            string   `json:"url"`
	URLs           []string `json:"urls"`
	Tag            string   `json:"tag"`
	Format         string   `json:"format"`
	UpdateInterval string   `json:"updateInterval"`
	Path           string   `json:"path"`
	IsEnabled      *bool    `json:"isEnabled"`
}

type updateSRSListRequest struct {
	UUID           string  `json:"uuid"`
	URL            *string `json:"url"`
	Tag            *string `json:"tag"`
	Format         *string `json:"format"`
	UpdateInterval *string `json:"updateInterval"`
	Path           *string `json:"path"`
	IsEnabled      *bool   `json:"isEnabled"`
}

type reorderSRSListsRequest struct {
	Items []struct {
		UUID         string `json:"uuid"`
		ViewPosition int    `json:"viewPosition"`
	} `json:"items"`
}

type bulkDeleteRequest struct {
	UUIDs []string `json:"uuids"`
}

type bulkEnableRequest struct {
	UUIDs []string `json:"uuids"`
}

type bulkSetIntervalRequest struct {
	UUIDs          []string `json:"uuids"`
	UpdateInterval string   `json:"updateInterval"`
}

type checkListsRequest struct {
	UUIDs []string `json:"uuids"`
}

func SRSListsHandler(manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handleGetSRSLists(w, r, manager, cfg)
		case http.MethodPost:
			handleCreateSRSLists(w, r, manager, cfg)
		case http.MethodPatch:
			handleUpdateSRSList(w, r, manager, cfg)
		default:
			shared.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
		}
	}
}

func SRSListByUUIDHandler(manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uuidStr := strings.TrimSpace(strings.TrimPrefix(r.URL.Path, "/api/srs-lists/"))
		if uuidStr == "" {
			SRSListsHandler(manager, cfg)(w, r)
			return
		}
		if _, err := uuid.Parse(uuidStr); err != nil {
			shared.SendError(w, http.StatusBadRequest, "invalid UUID format", nil, cfg)
			return
		}

		switch r.Method {
		case http.MethodDelete:
			handleDeleteSRSList(w, r, manager, cfg, uuidStr)
		case http.MethodGet:
			handleGetSRSList(w, r, manager, cfg, uuidStr)
		default:
			shared.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
		}
	}
}

func SRSListsActionsHandler(manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			shared.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}
		path := strings.Trim(strings.TrimPrefix(r.URL.Path, "/api/srs-lists/actions/"), "/")
		switch path {
		case "reorder":
			handleReorderSRSLists(w, r, manager, cfg)
		case "check":
			handleCheckSRSLists(w, r, manager, cfg)
		case "sync":
			monitor.RequestSRSDeploy()
			shared.WriteJSON(w, http.StatusOK, map[string]any{"response": map[string]any{"queued": true}})
		default:
			http.NotFound(w, r)
		}
	}
}

func SRSListsBulkHandler(manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			shared.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}
		path := strings.Trim(strings.TrimPrefix(r.URL.Path, "/api/srs-lists/bulk/"), "/")
		switch path {
		case "delete":
			handleBulkDeleteSRSLists(w, r, manager, cfg)
		case "enable":
			handleBulkEnableSRSLists(w, r, manager, cfg, true)
		case "disable":
			handleBulkEnableSRSLists(w, r, manager, cfg, false)
		case "set-interval":
			handleBulkSetIntervalSRSLists(w, r, manager, cfg)
		default:
			http.NotFound(w, r)
		}
	}
}

func handleGetSRSLists(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) {
	items, err := srscore.LoadAll(r.Context(), manager)
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to fetch srs lists", err, cfg)
		return
	}
	apiItems := make([]srsListAPI, 0, len(items))
	for _, item := range items {
		tag := strings.TrimSpace(item.Tag)
		if tag == "" {
			tag = srscore.DeriveTagFromFileName(item.FileName)
		}
		apiItems = append(apiItems, srsListAPI{
			UUID:           item.UUID,
			Tag:            tag,
			Format:         item.Format,
			URL:            item.URL,
			UpdateInterval: item.UpdateInterval,
			Path:           item.Path,
			FileName:       item.FileName,
			ShortName:      item.FileName,
			ViewPosition:   item.ViewPosition,
			IsEnabled:      item.IsEnabled,
			IsAvailable:    item.IsAvailable,
			LastCheckedAt:  item.LastCheckedAt,
			LastError:      item.LastError,
			CreatedAt:      item.CreatedAt,
			UpdatedAt:      item.UpdatedAt,
		})
	}
	shared.WriteJSON(w, http.StatusOK, map[string]any{"response": map[string]any{"srsLists": apiItems}})
}

func handleGetSRSList(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig, listUUID string) {
	items, err := srscore.LoadAll(r.Context(), manager)
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to fetch srs list", err, cfg)
		return
	}
	for _, item := range items {
		if item.UUID == listUUID {
			tag := strings.TrimSpace(item.Tag)
			if tag == "" {
				tag = srscore.DeriveTagFromFileName(item.FileName)
			}
			shared.WriteJSON(w, http.StatusOK, map[string]any{"response": srsListAPI{
				UUID:           item.UUID,
				Tag:            tag,
				Format:         item.Format,
				URL:            item.URL,
				UpdateInterval: item.UpdateInterval,
				Path:           item.Path,
				FileName:       item.FileName,
				ShortName:      item.FileName,
				ViewPosition:   item.ViewPosition,
				IsEnabled:      item.IsEnabled,
				IsAvailable:    item.IsAvailable,
				LastCheckedAt:  item.LastCheckedAt,
				LastError:      item.LastError,
				CreatedAt:      item.CreatedAt,
				UpdatedAt:      item.UpdatedAt,
			}})
			return
		}
	}
	shared.WriteJSONError(w, http.StatusNotFound, "srs list not found")
}

func handleCreateSRSLists(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) {
	reqStarted := time.Now()
	defer func() {
		cfg.Logger.Debug("SRS create request completed", "duration_ms", time.Since(reqStarted).Milliseconds())
	}()

	var req createSRSListsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		shared.SendError(w, http.StatusBadRequest, "invalid request payload", err, cfg)
		return
	}

	rawURLs := make([]string, 0, len(req.URLs)+1)
	if strings.TrimSpace(req.URL) != "" {
		rawURLs = append(rawURLs, req.URL)
	}
	rawURLs = append(rawURLs, req.URLs...)

	if len(rawURLs) == 0 {
		shared.WriteJSONError(w, http.StatusBadRequest, "url or urls[] is required")
		return
	}

	type createItem struct {
		Tag            string
		Format         string
		URL            string
		UpdateInterval string
		Path           *string
		FileName       string
		IsEnabled      bool
	}

	toCreate := make([]createItem, 0, len(rawURLs))
	seenURL := make(map[string]struct{}, len(rawURLs))
	for _, raw := range rawURLs {
		cleanURL := strings.TrimSpace(raw)
		if cleanURL == "" {
			continue
		}
		if _, exists := seenURL[cleanURL]; exists {
			continue
		}
		seenURL[cleanURL] = struct{}{}

		fileName, err := srscore.DeriveFileNameFromURL(cleanURL)
		if err != nil {
			shared.SendError(w, http.StatusBadRequest, fmt.Sprintf("invalid url %q", cleanURL), err, cfg)
			return
		}
		itemFormat := strings.ToLower(strings.TrimSpace(req.Format))
		if itemFormat == "" {
			itemFormat = "binary"
		}

		updateInterval := strings.TrimSpace(req.UpdateInterval)
		if updateInterval == "" {
			updateInterval = "1d"
		}

		tag := strings.TrimSpace(req.Tag)
		if len(rawURLs) > 1 || tag == "" {
			tag = srscore.DeriveTagFromFileName(fileName)
		}

		var pathValue *string
		if p := strings.TrimSpace(req.Path); p != "" {
			pathValue = &p
		}
		isEnabled := true
		if req.IsEnabled != nil {
			isEnabled = *req.IsEnabled
		}

		toCreate = append(toCreate, createItem{
			Tag:            tag,
			Format:         itemFormat,
			URL:            cleanURL,
			UpdateInterval: updateInterval,
			Path:           pathValue,
			FileName:       fileName,
			IsEnabled:      isEnabled,
		})
	}

	if len(toCreate) == 0 {
		shared.WriteJSONError(w, http.StatusBadRequest, "no valid urls provided")
		return
	}

	dbStarted := time.Now()
	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		tx, err := db.BeginTx(r.Context(), nil)
		if err != nil {
			return err
		}
		defer func() {
			_ = tx.Rollback()
		}()

		var maxPos sql.NullInt64
		if err := tx.QueryRowContext(r.Context(), `SELECT COALESCE(MAX(view_position), -1) FROM srs_lists`).Scan(&maxPos); err != nil {
			return err
		}
		currentPos := int64(-1)
		if maxPos.Valid {
			currentPos = maxPos.Int64
		}

		for _, item := range toCreate {
			currentPos++
			if _, err := tx.ExecContext(r.Context(), `
				INSERT INTO srs_lists (tag, format, url, update_interval, path, file_name, view_position, is_enabled, is_available, created_at, updated_at)
				VALUES (?, ?, ?, ?, ?, ?, ?, ?, false, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
			`, item.Tag, item.Format, item.URL, item.UpdateInterval, item.Path, item.FileName, currentPos, item.IsEnabled); err != nil {
				return err
			}
		}

		return tx.Commit()
	})
	if err != nil {
		shared.SendError(w, http.StatusBadRequest, "failed to create srs lists", err, cfg)
		return
	}
	cfg.Logger.Debug("SRS create DB write completed", "duration_ms", time.Since(dbStarted).Milliseconds(), "created", len(toCreate))

	checkStarted := time.Now()
	if _, err := srscore.CheckAndUpdateAvailability(context.Background(), manager, cfg); err != nil {
		cfg.Logger.Warn("Failed to check SRS lists right after create", "error", err)
	}
	cfg.Logger.Debug("SRS create availability check completed", "duration_ms", time.Since(checkStarted).Milliseconds())

	queueStarted := time.Now()
	cfg.Logger.Debug("SRS create sync/deploy skipped (manual sync only)", "duration_ms", time.Since(queueStarted).Milliseconds())

	handleGetSRSLists(w, r, manager, cfg)
}

func handleUpdateSRSList(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) {
	reqStarted := time.Now()
	defer func() {
		cfg.Logger.Debug("SRS update request completed", "duration_ms", time.Since(reqStarted).Milliseconds())
	}()

	var req updateSRSListRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		shared.SendError(w, http.StatusBadRequest, "invalid request payload", err, cfg)
		return
	}
	if strings.TrimSpace(req.UUID) == "" {
		shared.WriteJSONError(w, http.StatusBadRequest, "uuid is required")
		return
	}
	if _, err := uuid.Parse(strings.TrimSpace(req.UUID)); err != nil {
		shared.WriteJSONError(w, http.StatusBadRequest, "invalid uuid")
		return
	}

	updates := make([]string, 0, 8)
	args := make([]any, 0, 9)

	var fileName string
	if req.URL != nil {
		urlValue := strings.TrimSpace(*req.URL)
		if urlValue == "" {
			shared.WriteJSONError(w, http.StatusBadRequest, "url cannot be empty")
			return
		}
		derivedName, err := srscore.DeriveFileNameFromURL(urlValue)
		if err != nil {
			shared.SendError(w, http.StatusBadRequest, "invalid url", err, cfg)
			return
		}
		fileName = derivedName
		updates = append(updates, "url = ?", "file_name = ?", "is_available = false", "last_error = NULL")
		args = append(args, urlValue, fileName)
	}

	if req.Tag != nil {
		tagValue := strings.TrimSpace(*req.Tag)
		if tagValue == "" {
			if fileName == "" {
				fileName = "ruleset.srs"
			}
			tagValue = srscore.DeriveTagFromFileName(fileName)
		}
		updates = append(updates, "tag = ?")
		args = append(args, tagValue)
	}

	if req.Format != nil {
		formatValue := strings.ToLower(strings.TrimSpace(*req.Format))
		if formatValue == "" {
			formatValue = "binary"
		}
		updates = append(updates, "format = ?")
		args = append(args, formatValue)
	}

	if req.UpdateInterval != nil {
		val := strings.TrimSpace(*req.UpdateInterval)
		if val == "" {
			val = "1d"
		}
		updates = append(updates, "update_interval = ?")
		args = append(args, val)
	}

	if req.Path != nil {
		pathValue := strings.TrimSpace(*req.Path)
		if pathValue == "" {
			updates = append(updates, "path = NULL")
		} else {
			updates = append(updates, "path = ?")
			args = append(args, pathValue)
		}
	}
	if req.IsEnabled != nil {
		updates = append(updates, "is_enabled = ?")
		args = append(args, *req.IsEnabled)
	}

	if len(updates) == 0 {
		shared.WriteJSONError(w, http.StatusBadRequest, "nothing to update")
		return
	}

	updates = append(updates, "updated_at = CURRENT_TIMESTAMP")
	args = append(args, strings.TrimSpace(req.UUID))

	query := fmt.Sprintf("UPDATE srs_lists SET %s WHERE uuid = ?", strings.Join(updates, ", "))

	dbStarted := time.Now()
	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		res, execErr := db.ExecContext(r.Context(), query, args...)
		if execErr != nil {
			return execErr
		}
		rows, rowsErr := res.RowsAffected()
		if rowsErr != nil {
			return rowsErr
		}
		if rows == 0 {
			return fmt.Errorf("srs list not found")
		}
		return nil
	})
	if err != nil {
		shared.SendError(w, http.StatusBadRequest, "failed to update srs list", err, cfg)
		return
	}
	cfg.Logger.Debug("SRS update DB write completed", "duration_ms", time.Since(dbStarted).Milliseconds(), "uuid", strings.TrimSpace(req.UUID))

	checkStarted := time.Now()
	if _, err := srscore.CheckAndUpdateAvailability(context.Background(), manager, cfg); err != nil {
		cfg.Logger.Warn("Failed to check SRS list after update", "error", err)
	}
	cfg.Logger.Debug("SRS update availability check completed", "duration_ms", time.Since(checkStarted).Milliseconds(), "uuid", strings.TrimSpace(req.UUID))

	queueStarted := time.Now()
	cfg.Logger.Debug("SRS update sync/deploy skipped (manual sync only)", "duration_ms", time.Since(queueStarted).Milliseconds(), "uuid", strings.TrimSpace(req.UUID))
	handleGetSRSLists(w, r, manager, cfg)
}

func handleDeleteSRSList(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig, listUUID string) {
	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		res, execErr := db.ExecContext(r.Context(), `DELETE FROM srs_lists WHERE uuid = ?`, listUUID)
		if execErr != nil {
			return execErr
		}
		rows, rowsErr := res.RowsAffected()
		if rowsErr != nil {
			return rowsErr
		}
		if rows == 0 {
			return fmt.Errorf("srs list not found")
		}
		return nil
	})
	if err != nil {
		shared.SendError(w, http.StatusBadRequest, "failed to delete srs list", err, cfg)
		return
	}
	shared.WriteJSON(w, http.StatusOK, map[string]any{"response": map[string]any{"deleted": true}})
}

func handleBulkDeleteSRSLists(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) {
	var req bulkDeleteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		shared.SendError(w, http.StatusBadRequest, "invalid request payload", err, cfg)
		return
	}
	cleanUUIDs, err := normalizeUUIDs(req.UUIDs)
	if len(cleanUUIDs) == 0 {
		shared.WriteJSONError(w, http.StatusBadRequest, "uuids are required")
		return
	}
	if err != nil {
		shared.WriteJSONError(w, http.StatusBadRequest, "one or more uuids are invalid")
		return
	}

	err = manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		_, execErr := db.ExecContext(r.Context(), `DELETE FROM srs_lists WHERE uuid = ANY(?)`, cleanUUIDs)
		return execErr
	})
	if err != nil {
		shared.SendError(w, http.StatusBadRequest, "failed to bulk delete srs lists", err, cfg)
		return
	}
	shared.WriteJSON(w, http.StatusOK, map[string]any{"response": map[string]any{"deleted": true}})
}

func handleBulkEnableSRSLists(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig, enabled bool) {
	var req bulkEnableRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		shared.SendError(w, http.StatusBadRequest, "invalid request payload", err, cfg)
		return
	}
	cleanUUIDs, err := normalizeUUIDs(req.UUIDs)
	if err != nil {
		shared.WriteJSONError(w, http.StatusBadRequest, "one or more uuids are invalid")
		return
	}
	if len(cleanUUIDs) == 0 {
		shared.WriteJSONError(w, http.StatusBadRequest, "uuids are required")
		return
	}

	err = manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		_, execErr := db.ExecContext(r.Context(), `
			UPDATE srs_lists
			SET is_enabled = ?, updated_at = CURRENT_TIMESTAMP
			WHERE uuid = ANY(?)
		`, enabled, cleanUUIDs)
		return execErr
	})
	if err != nil {
		shared.SendError(w, http.StatusBadRequest, "failed to update srs lists state", err, cfg)
		return
	}

	handleGetSRSLists(w, r, manager, cfg)
}

func handleBulkSetIntervalSRSLists(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) {
	var req bulkSetIntervalRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		shared.SendError(w, http.StatusBadRequest, "invalid request payload", err, cfg)
		return
	}
	cleanUUIDs, err := normalizeUUIDs(req.UUIDs)
	if err != nil {
		shared.WriteJSONError(w, http.StatusBadRequest, "one or more uuids are invalid")
		return
	}
	if len(cleanUUIDs) == 0 {
		shared.WriteJSONError(w, http.StatusBadRequest, "uuids are required")
		return
	}
	interval := strings.TrimSpace(req.UpdateInterval)
	if interval == "" {
		shared.WriteJSONError(w, http.StatusBadRequest, "updateInterval is required")
		return
	}

	err = manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		_, execErr := db.ExecContext(r.Context(), `
			UPDATE srs_lists
			SET update_interval = ?, updated_at = CURRENT_TIMESTAMP
			WHERE uuid = ANY(?)
		`, interval, cleanUUIDs)
		return execErr
	})
	if err != nil {
		shared.SendError(w, http.StatusBadRequest, "failed to update update interval", err, cfg)
		return
	}

	handleGetSRSLists(w, r, manager, cfg)
}

func handleReorderSRSLists(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) {
	var req reorderSRSListsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		shared.SendError(w, http.StatusBadRequest, "invalid request payload", err, cfg)
		return
	}
	if len(req.Items) == 0 {
		shared.WriteJSONError(w, http.StatusBadRequest, "items are required")
		return
	}

	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		tx, err := db.BeginTx(r.Context(), nil)
		if err != nil {
			return err
		}
		defer func() {
			_ = tx.Rollback()
		}()

		for _, item := range req.Items {
			if _, err := uuid.Parse(strings.TrimSpace(item.UUID)); err != nil {
				return fmt.Errorf("invalid uuid: %s", item.UUID)
			}
			if _, err := tx.ExecContext(r.Context(), `
				UPDATE srs_lists SET view_position = ?, updated_at = CURRENT_TIMESTAMP WHERE uuid = ?
			`, item.ViewPosition, item.UUID); err != nil {
				return err
			}
		}

		return tx.Commit()
	})
	if err != nil {
		shared.SendError(w, http.StatusBadRequest, "failed to reorder srs lists", err, cfg)
		return
	}
	handleGetSRSLists(w, r, manager, cfg)
}

func handleCheckSRSLists(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) {
	started := time.Now()
	defer func() {
		cfg.Logger.Debug("SRS check request completed", "duration_ms", time.Since(started).Milliseconds())
	}()

	var req checkListsRequest
	_ = json.NewDecoder(r.Body).Decode(&req)

	if len(req.UUIDs) > 0 {
		if err := checkSelectedLists(r.Context(), manager, cfg, req.UUIDs); err != nil {
			shared.SendError(w, http.StatusBadRequest, "failed to check selected srs lists", err, cfg)
			return
		}
	} else {
		if _, err := srscore.CheckAndUpdateAvailability(r.Context(), manager, cfg); err != nil {
			shared.SendError(w, http.StatusBadRequest, "failed to check srs lists", err, cfg)
			return
		}
	}
	handleGetSRSLists(w, r, manager, cfg)
}

func checkSelectedLists(ctx context.Context, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig, uuids []string) error {
	clean, err := normalizeUUIDs(uuids)
	if err != nil {
		return err
	}
	if len(clean) == 0 {
		return nil
	}

	itemsByUUID := make(map[string]string, len(clean))
	if err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		rows, err := db.QueryContext(ctx, `SELECT uuid, url FROM srs_lists WHERE uuid = ANY(?)`, clean)
		if err != nil {
			return err
		}
		defer rows.Close()
		for rows.Next() {
			var id, rawURL string
			if err := rows.Scan(&id, &rawURL); err != nil {
				return err
			}
			itemsByUUID[id] = rawURL
		}
		return rows.Err()
	}); err != nil {
		return err
	}

	for id, rawURL := range itemsByUUID {
		err := srscore.CheckOneURL(ctx, rawURL)
		isAvailable := err == nil
		var errText any
		if err != nil {
			errText = err.Error()
		}
		if writeErr := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
			_, execErr := db.ExecContext(ctx, `
				UPDATE srs_lists
				SET is_available = ?,
					last_checked_at = CURRENT_TIMESTAMP,
					last_error = ?,
					updated_at = CURRENT_TIMESTAMP
				WHERE uuid = ?
			`, isAvailable, errText, id)
			return execErr
		}); writeErr != nil {
			cfg.Logger.Warn("Failed to write selected srs check result", "uuid", id, "error", writeErr)
		}
	}

	return nil
}

func normalizeUUIDs(values []string) ([]string, error) {
	clean := make([]string, 0, len(values))
	seen := make(map[string]struct{}, len(values))
	for _, raw := range values {
		u := strings.TrimSpace(raw)
		if u == "" {
			continue
		}
		if _, err := uuid.Parse(u); err != nil {
			return nil, fmt.Errorf("invalid uuid: %s", u)
		}
		if _, ok := seen[u]; ok {
			continue
		}
		seen[u] = struct{}{}
		clean = append(clean, u)
	}
	return clean, nil
}
