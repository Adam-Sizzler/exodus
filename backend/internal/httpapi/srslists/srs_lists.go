package srslists

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"exodus/internal/config"
	"exodus/internal/httpapi/shared"
	srscore "exodus/internal/srslists"
	monitor "exodus/internal/subscriptionnodes"

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

func SRSListsHandler(db *sql.DB, cfg *config.BackendConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handleGetSRSLists(w, r, db, cfg)
		case http.MethodPost:
			handleCreateSRSLists(w, r, db, cfg)
		case http.MethodPatch:
			handleUpdateSRSList(w, r, db, cfg)
		default:
			shared.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
		}
	}
}

func SRSListByUUIDHandler(db *sql.DB, cfg *config.BackendConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uuidStr := strings.TrimSpace(strings.TrimPrefix(r.URL.Path, "/api/srs-lists/"))
		if uuidStr == "" {
			SRSListsHandler(db, cfg)(w, r)
			return
		}
		if _, err := uuid.Parse(uuidStr); err != nil {
			shared.SendError(w, http.StatusBadRequest, "invalid UUID format", nil, cfg)
			return
		}

		switch r.Method {
		case http.MethodDelete:
			handleDeleteSRSList(w, r, db, cfg, uuidStr)
		case http.MethodGet:
			handleGetSRSList(w, r, db, cfg, uuidStr)
		default:
			shared.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
		}
	}
}

func SRSListsActionsHandler(db *sql.DB, cfg *config.BackendConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			shared.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}
		path := strings.Trim(strings.TrimPrefix(r.URL.Path, "/api/srs-lists/actions/"), "/")
		switch path {
		case "reorder":
			handleReorderSRSLists(w, r, db, cfg)
		case "check":
			handleCheckSRSLists(w, r, db, cfg)
		case "sync":
			monitor.RequestSubNodeSRSDeploy()
			shared.WriteJSON(w, http.StatusOK, map[string]any{"response": map[string]any{"queued": true}})
		default:
			http.NotFound(w, r)
		}
	}
}

func SRSListsBulkHandler(db *sql.DB, cfg *config.BackendConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			shared.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}
		path := strings.Trim(strings.TrimPrefix(r.URL.Path, "/api/srs-lists/bulk/"), "/")
		switch path {
		case "delete":
			handleBulkDeleteSRSLists(w, r, db, cfg)
		case "enable":
			handleBulkEnableSRSLists(w, r, db, cfg, true)
		case "disable":
			handleBulkEnableSRSLists(w, r, db, cfg, false)
		case "set-interval":
			handleBulkSetIntervalSRSLists(w, r, db, cfg)
		default:
			http.NotFound(w, r)
		}
	}
}

func handleGetSRSLists(w http.ResponseWriter, r *http.Request, db *sql.DB, cfg *config.BackendConfig) {
	items, err := srscore.LoadAll(r.Context(), db)
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

func handleGetSRSList(w http.ResponseWriter, r *http.Request, db *sql.DB, cfg *config.BackendConfig, listUUID string) {
	items, err := srscore.LoadAll(r.Context(), db)
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

func handleCreateSRSLists(w http.ResponseWriter, r *http.Request, db *sql.DB, cfg *config.BackendConfig) {
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
	tx, err := db.BeginTx(r.Context(), nil)
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "begin transaction failed", err, cfg)
		return
	}
	defer func() {
		_ = tx.Rollback()
	}()

	var maxPos sql.NullInt64
	if err := tx.QueryRowContext(r.Context(), `SELECT COALESCE(MAX(view_position), -1) FROM srs_lists`).Scan(&maxPos); err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to query max position", err, cfg)
		return
	}
	currentPos := int64(-1)
	if maxPos.Valid {
		currentPos = maxPos.Int64
	}

	for _, item := range toCreate {
		currentPos++
		if _, err := tx.ExecContext(r.Context(), `
			INSERT INTO srs_lists (tag, format, url, update_interval, path, file_name, view_position, is_enabled, is_available, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, false, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		`, item.Tag, item.Format, item.URL, item.UpdateInterval, item.Path, item.FileName, currentPos, item.IsEnabled); err != nil {
			shared.SendError(w, http.StatusBadRequest, "failed to create srs list", err, cfg)
			return
		}
	}

	if err := tx.Commit(); err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to commit transaction", err, cfg)
		return
	}

	cfg.Logger.Debug("SRS create DB write completed", "duration_ms", time.Since(dbStarted).Milliseconds(), "created", len(toCreate))

	checkStarted := time.Now()
	if _, err := srscore.CheckAndUpdateAvailability(context.Background(), db, cfg); err != nil {
		cfg.Logger.Warn("Failed to check SRS lists right after create", "error", err)
	}
	cfg.Logger.Debug("SRS create availability check completed", "duration_ms", time.Since(checkStarted).Milliseconds())

	handleGetSRSLists(w, r, db, cfg)
}

func handleUpdateSRSList(w http.ResponseWriter, r *http.Request, db *sql.DB, cfg *config.BackendConfig) {
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
	idx := 1

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
		updates = append(updates, fmt.Sprintf("url = $%d", idx), fmt.Sprintf("file_name = $%d", idx+1), "is_available = false", "last_error = NULL")
		args = append(args, urlValue, fileName)
		idx += 2
	}

	if req.Tag != nil {
		tagValue := strings.TrimSpace(*req.Tag)
		if tagValue == "" {
			if fileName == "" {
				fileName = "ruleset.srs"
			}
			tagValue = srscore.DeriveTagFromFileName(fileName)
		}
		updates = append(updates, fmt.Sprintf("tag = $%d", idx))
		args = append(args, tagValue)
		idx++
	}

	if req.Format != nil {
		formatValue := strings.ToLower(strings.TrimSpace(*req.Format))
		if formatValue == "" {
			formatValue = "binary"
		}
		updates = append(updates, fmt.Sprintf("format = $%d", idx))
		args = append(args, formatValue)
		idx++
	}

	if req.UpdateInterval != nil {
		val := strings.TrimSpace(*req.UpdateInterval)
		if val == "" {
			val = "1d"
		}
		updates = append(updates, fmt.Sprintf("update_interval = $%d", idx))
		args = append(args, val)
		idx++
	}

	if req.Path != nil {
		pathValue := strings.TrimSpace(*req.Path)
		if pathValue == "" {
			updates = append(updates, "path = NULL")
		} else {
			updates = append(updates, fmt.Sprintf("path = $%d", idx))
			args = append(args, pathValue)
			idx++
		}
	}
	if req.IsEnabled != nil {
		updates = append(updates, fmt.Sprintf("is_enabled = $%d", idx))
		args = append(args, *req.IsEnabled)
		idx++
	}

	if len(updates) == 0 {
		shared.WriteJSONError(w, http.StatusBadRequest, "nothing to update")
		return
	}

	updates = append(updates, "updated_at = CURRENT_TIMESTAMP")
	args = append(args, strings.TrimSpace(req.UUID))

	query := fmt.Sprintf("UPDATE srs_lists SET %s WHERE uuid = $%d", strings.Join(updates, ", "), idx)

	dbStarted := time.Now()
	res, execErr := db.ExecContext(r.Context(), query, args...)
	if execErr != nil {
		shared.SendError(w, http.StatusBadRequest, "failed to update srs list", execErr, cfg)
		return
	}
	rows, rowsErr := res.RowsAffected()
	if rowsErr != nil {
		shared.SendError(w, http.StatusBadRequest, "failed to read rows affected", rowsErr, cfg)
		return
	}
	if rows == 0 {
		shared.WriteJSONError(w, http.StatusNotFound, "srs list not found")
		return
	}

	cfg.Logger.Debug("SRS update DB write completed", "duration_ms", time.Since(dbStarted).Milliseconds(), "uuid", strings.TrimSpace(req.UUID))

	checkStarted := time.Now()
	if _, err := srscore.CheckAndUpdateAvailability(context.Background(), db, cfg); err != nil {
		cfg.Logger.Warn("Failed to check SRS list after update", "error", err)
	}
	cfg.Logger.Debug("SRS update availability check completed", "duration_ms", time.Since(checkStarted).Milliseconds(), "uuid", strings.TrimSpace(req.UUID))

	handleGetSRSLists(w, r, db, cfg)
}

func handleDeleteSRSList(w http.ResponseWriter, r *http.Request, db *sql.DB, cfg *config.BackendConfig, listUUID string) {
	res, execErr := db.ExecContext(r.Context(), `DELETE FROM srs_lists WHERE uuid = $1`, listUUID)
	if execErr != nil {
		shared.SendError(w, http.StatusBadRequest, "failed to delete srs list", execErr, cfg)
		return
	}
	rows, rowsErr := res.RowsAffected()
	if rowsErr != nil {
		shared.SendError(w, http.StatusBadRequest, "failed to read rows affected", rowsErr, cfg)
		return
	}
	if rows == 0 {
		shared.WriteJSONError(w, http.StatusNotFound, "srs list not found")
		return
	}
	shared.WriteJSON(w, http.StatusOK, map[string]any{"response": map[string]any{"deleted": true}})
}

func handleBulkDeleteSRSLists(w http.ResponseWriter, r *http.Request, db *sql.DB, cfg *config.BackendConfig) {
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

	_, execErr := db.ExecContext(r.Context(), `DELETE FROM srs_lists WHERE uuid = ANY($1)`, cleanUUIDs)
	if execErr != nil {
		shared.SendError(w, http.StatusBadRequest, "failed to bulk delete srs lists", execErr, cfg)
		return
	}
	shared.WriteJSON(w, http.StatusOK, map[string]any{"response": map[string]any{"deleted": true}})
}

func handleBulkEnableSRSLists(w http.ResponseWriter, r *http.Request, db *sql.DB, cfg *config.BackendConfig, enabled bool) {
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

	_, execErr := db.ExecContext(r.Context(), `
		UPDATE srs_lists
		SET is_enabled = $1, updated_at = CURRENT_TIMESTAMP
		WHERE uuid = ANY($2)
	`, enabled, cleanUUIDs)
	if execErr != nil {
		shared.SendError(w, http.StatusBadRequest, "failed to update srs lists state", execErr, cfg)
		return
	}

	handleGetSRSLists(w, r, db, cfg)
}

func handleBulkSetIntervalSRSLists(w http.ResponseWriter, r *http.Request, db *sql.DB, cfg *config.BackendConfig) {
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

	_, execErr := db.ExecContext(r.Context(), `
		UPDATE srs_lists
		SET update_interval = $1, updated_at = CURRENT_TIMESTAMP
		WHERE uuid = ANY($2)
	`, interval, cleanUUIDs)
	if execErr != nil {
		shared.SendError(w, http.StatusBadRequest, "failed to update update interval", execErr, cfg)
		return
	}

	handleGetSRSLists(w, r, db, cfg)
}

func handleReorderSRSLists(w http.ResponseWriter, r *http.Request, db *sql.DB, cfg *config.BackendConfig) {
	var req reorderSRSListsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		shared.SendError(w, http.StatusBadRequest, "invalid request payload", err, cfg)
		return
	}
	if len(req.Items) == 0 {
		shared.WriteJSONError(w, http.StatusBadRequest, "items are required")
		return
	}

	tx, err := db.BeginTx(r.Context(), nil)
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "begin transaction failed", err, cfg)
		return
	}
	defer func() {
		_ = tx.Rollback()
	}()

	for _, item := range req.Items {
		if _, err := uuid.Parse(strings.TrimSpace(item.UUID)); err != nil {
			shared.SendError(w, http.StatusBadRequest, fmt.Sprintf("invalid uuid: %s", item.UUID), nil, cfg)
			return
		}
		if _, err := tx.ExecContext(r.Context(), `
			UPDATE srs_lists SET view_position = $1, updated_at = CURRENT_TIMESTAMP WHERE uuid = $2
		`, item.ViewPosition, item.UUID); err != nil {
			shared.SendError(w, http.StatusBadRequest, "failed to update item position", err, cfg)
			return
		}
	}

	if err := tx.Commit(); err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to commit transaction", err, cfg)
		return
	}

	handleGetSRSLists(w, r, db, cfg)
}

func handleCheckSRSLists(w http.ResponseWriter, r *http.Request, db *sql.DB, cfg *config.BackendConfig) {
	started := time.Now()
	defer func() {
		cfg.Logger.Debug("SRS check request completed", "duration_ms", time.Since(started).Milliseconds())
	}()

	var req checkListsRequest
	_ = json.NewDecoder(r.Body).Decode(&req)

	if len(req.UUIDs) > 0 {
		if err := checkSelectedLists(r.Context(), db, cfg, req.UUIDs); err != nil {
			shared.SendError(w, http.StatusBadRequest, "failed to check selected srs lists", err, cfg)
			return
		}
	} else {
		if _, err := srscore.CheckAndUpdateAvailability(r.Context(), db, cfg); err != nil {
			shared.SendError(w, http.StatusBadRequest, "failed to check srs lists", err, cfg)
			return
		}
	}
	handleGetSRSLists(w, r, db, cfg)
}

func checkSelectedLists(ctx context.Context, db *sql.DB, cfg *config.BackendConfig, uuids []string) error {
	clean, err := normalizeUUIDs(uuids)
	if err != nil {
		return err
	}
	if len(clean) == 0 {
		return nil
	}

	itemsByUUID := make(map[string]string, len(clean))
	rows, err := db.QueryContext(ctx, `SELECT uuid, url FROM srs_lists WHERE uuid = ANY($1)`, clean)
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
	if err := rows.Err(); err != nil {
		return err
	}

	for id, rawURL := range itemsByUUID {
		err := srscore.CheckOneURL(ctx, rawURL)
		isAvailable := err == nil
		var errText any
		if err != nil {
			errText = err.Error()
		}
		if _, writeErr := db.ExecContext(ctx, `
			UPDATE srs_lists
			SET is_available = $1,
				last_checked_at = CURRENT_TIMESTAMP,
				last_error = $2,
				updated_at = CURRENT_TIMESTAMP
			WHERE uuid = $3
		`, isAvailable, errText, id); writeErr != nil {
			if cfg != nil && cfg.Logger != nil {
				cfg.Logger.Warn("Failed to write selected srs check result", "uuid", id, "error", writeErr)
			}
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
