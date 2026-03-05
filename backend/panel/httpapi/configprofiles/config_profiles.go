package configprofiles

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"v2ray-stat/backend/panel/config"
	dbmanager "v2ray-stat/backend/panel/db/manager"
	"v2ray-stat/backend/panel/dbutil"
	"v2ray-stat/backend/panel/httpapi/shared"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
)

var errConfigProfileNotFound = errors.New("config profile not found")

type ConfigProfileInbound struct {
	UUID         string          `json:"uuid"`
	ProfileUUID  string          `json:"profileUuid"`
	Tag          string          `json:"tag"`
	Type         string          `json:"type"`
	Network      *string         `json:"network"`
	Security     *string         `json:"security"`
	Port         *int            `json:"port"`
	RawInbound   json.RawMessage `json:"rawInbound"`
	ActiveSquads []string        `json:"activeSquads"`
}

type ConfigProfileNode struct {
	UUID        string `json:"uuid"`
	Name        string `json:"name"`
	CountryCode string `json:"countryCode"`
}

type ConfigProfile struct {
	UUID         string                `json:"uuid"`
	ViewPosition int                   `json:"viewPosition"`
	Name         string                `json:"name"`
	Config       json.RawMessage       `json:"config"`
	Inbounds     []ConfigProfileInbound `json:"inbounds"`
	Nodes        []ConfigProfileNode   `json:"nodes"`
	CreatedAt    time.Time             `json:"createdAt"`
	UpdatedAt    time.Time             `json:"updatedAt"`
}

type configProfileRecord struct {
	UUID         string
	ViewPosition int
	Name         string
	Config       json.RawMessage
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type createConfigProfileRequest struct {
	Name   string          `json:"name"`
	Config json.RawMessage `json:"config"`
}

type updateConfigProfileRequest struct {
	UUID   string           `json:"uuid"`
	Name   *string          `json:"name,omitempty"`
	Config *json.RawMessage `json:"config,omitempty"`
}

type reorderConfigProfilesRequest struct {
	Items []struct {
		UUID         string `json:"uuid"`
		ViewPosition int    `json:"viewPosition"`
	} `json:"items"`
}

func ConfigProfilesHandler(manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handleGetConfigProfiles(w, r, manager, cfg)
		case http.MethodPost:
			handleCreateConfigProfile(w, r, manager, cfg)
		case http.MethodPatch:
			handleUpdateConfigProfile(w, r, manager, cfg)
		default:
			shared.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
		}
	}
}

func ConfigProfileByUUIDHandler(manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := trimConfigProfilesPath(r.URL.Path, "/")
		if path == "" {
			http.NotFound(w, r)
			return
		}
		parts := strings.Split(path, "/")
		profileUUID := parts[0]
		if _, err := uuid.Parse(profileUUID); err != nil {
			shared.SendError(w, http.StatusBadRequest, "invalid UUID format", nil, cfg)
			return
		}

		if len(parts) == 2 && r.Method == http.MethodGet {
			switch parts[1] {
			case "inbounds":
				handleGetConfigProfileInbounds(w, r, manager, cfg, profileUUID)
			case "computed-config":
				handleGetComputedConfigProfile(w, r, manager, cfg, profileUUID)
			default:
				http.NotFound(w, r)
			}
			return
		}

		if len(parts) != 1 {
			http.NotFound(w, r)
			return
		}

		switch r.Method {
		case http.MethodGet:
			handleGetConfigProfile(w, r, manager, cfg, profileUUID)
		case http.MethodDelete:
			handleDeleteConfigProfile(w, r, manager, cfg, profileUUID)
		default:
			shared.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
		}
	}
}

func ConfigProfilesActionsHandler(manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			shared.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}

		path := trimConfigProfilesPath(r.URL.Path, "/actions/")
		switch path {
		case "reorder":
			handleReorderConfigProfiles(w, r, manager, cfg)
		default:
			http.NotFound(w, r)
		}
	}
}

func ConfigProfilesInboundsHandler(manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			shared.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}
		handleGetAllInbounds(w, r, manager, cfg)
	}
}

func trimConfigProfilesPath(path string, suffix string) string {
	for _, prefix := range []string{"/api/config-profiles", "/api/v1/config-profiles"} {
		if strings.HasPrefix(path, prefix+suffix) {
			return strings.Trim(strings.TrimPrefix(path, prefix+suffix), "/")
		}
	}
	return strings.Trim(path, "/")
}

func handleGetConfigProfiles(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) {
	records, err := getAllConfigProfileRecords(r.Context(), manager)
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to fetch config profiles", err, cfg)
		return
	}

	response, err := buildConfigProfileResponses(r.Context(), manager, records)
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to build config profiles response", err, cfg)
		return
	}

	shared.WriteJSON(w, http.StatusOK, map[string]any{
		"response": map[string]any{
			"total":          len(response),
			"configProfiles": response,
		},
	})
}

func handleGetConfigProfile(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig, profileUUID string) {
	record, err := getConfigProfileRecordByUUID(r.Context(), manager, profileUUID)
	if err != nil {
		if errors.Is(err, errConfigProfileNotFound) {
			shared.SendError(w, http.StatusNotFound, "config profile not found", nil, cfg)
			return
		}
		shared.SendError(w, http.StatusInternalServerError, "failed to fetch config profile", err, cfg)
		return
	}

	response, err := buildConfigProfileResponses(r.Context(), manager, []configProfileRecord{record})
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to build config profile response", err, cfg)
		return
	}
	shared.WriteJSON(w, http.StatusOK, map[string]any{"response": response[0]})
}

func handleGetComputedConfigProfile(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig, profileUUID string) {
	handleGetConfigProfile(w, r, manager, cfg, profileUUID)
}

func handleGetConfigProfileInbounds(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig, profileUUID string) {
	if _, err := getConfigProfileRecordByUUID(r.Context(), manager, profileUUID); err != nil {
		if errors.Is(err, errConfigProfileNotFound) {
			shared.SendError(w, http.StatusNotFound, "config profile not found", nil, cfg)
			return
		}
		shared.SendError(w, http.StatusInternalServerError, "failed to fetch config profile", err, cfg)
		return
	}

	inbounds, err := getConfigProfileInboundsMap(r.Context(), manager, []string{profileUUID})
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to fetch config profile inbounds", err, cfg)
		return
	}

	shared.WriteJSON(w, http.StatusOK, map[string]any{
		"response": map[string]any{
			"total":    len(inbounds[profileUUID]),
			"inbounds": inbounds[profileUUID],
		},
	})
}

func handleGetAllInbounds(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) {
	records, err := getAllConfigProfileRecords(r.Context(), manager)
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to fetch config profiles", err, cfg)
		return
	}
	profileUUIDs := make([]string, 0, len(records))
	for _, record := range records {
		profileUUIDs = append(profileUUIDs, record.UUID)
	}
	inboundsMap, err := getConfigProfileInboundsMap(r.Context(), manager, profileUUIDs)
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to fetch inbounds", err, cfg)
		return
	}
	all := make([]ConfigProfileInbound, 0)
	for _, profileUUID := range profileUUIDs {
		all = append(all, inboundsMap[profileUUID]...)
	}

	shared.WriteJSON(w, http.StatusOK, map[string]any{
		"response": map[string]any{
			"total":    len(all),
			"inbounds": all,
		},
	})
}

func handleCreateConfigProfile(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) {
	var req createConfigProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		shared.SendError(w, http.StatusBadRequest, "invalid JSON", err, cfg)
		return
	}
	if err := validateCreateConfigProfileRequest(req); err != nil {
		shared.SendError(w, http.StatusBadRequest, err.Error(), nil, cfg)
		return
	}

	profileUUID := uuid.NewString()
	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		tx, err := db.BeginTx(r.Context(), nil)
		if err != nil {
			return err
		}

		if _, err := tx.ExecContext(r.Context(), `
			INSERT INTO config_profiles (uuid, name, config, created_at, updated_at)
			VALUES (?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		`, profileUUID, strings.TrimSpace(req.Name), req.Config); err != nil {
			_ = tx.Rollback()
			return mapConfigProfileWriteError(err)
		}

		if _, err := syncConfigProfileInbounds(r.Context(), tx, profileUUID, req.Config, cfg); err != nil {
			_ = tx.Rollback()
			return err
		}
		return tx.Commit()
	})
	if err != nil {
		handleConfigProfileWriteError(w, err, cfg)
		return
	}

	handleGetConfigProfile(w, r, manager, cfg, profileUUID)
}

func handleUpdateConfigProfile(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) {
	var req updateConfigProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		shared.SendError(w, http.StatusBadRequest, "invalid JSON", err, cfg)
		return
	}
	if _, err := uuid.Parse(strings.TrimSpace(req.UUID)); err != nil {
		shared.SendError(w, http.StatusBadRequest, "invalid UUID format", nil, cfg)
		return
	}
	if err := validateUpdateConfigProfileRequest(req); err != nil {
		shared.SendError(w, http.StatusBadRequest, err.Error(), nil, cfg)
		return
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
		if req.Config != nil {
			add("config", *req.Config)
		}

		if len(clauses) == 0 {
			_ = tx.Rollback()
			return fmt.Errorf("no fields to update")
		}

		args = append(args, req.UUID)
		result, err := tx.ExecContext(r.Context(), fmt.Sprintf(`
			UPDATE config_profiles
			SET %s, updated_at = CURRENT_TIMESTAMP
			WHERE uuid = ?
		`, strings.Join(clauses, ", ")), args...)
		if err != nil {
			_ = tx.Rollback()
			return mapConfigProfileWriteError(err)
		}
		rows, err := result.RowsAffected()
		if err != nil {
			_ = tx.Rollback()
			return err
		}
		if rows == 0 {
			_ = tx.Rollback()
			return errConfigProfileNotFound
		}

		if req.Config != nil {
			if _, err := syncConfigProfileInbounds(r.Context(), tx, req.UUID, *req.Config, cfg); err != nil {
				_ = tx.Rollback()
				return err
			}
		}

		return tx.Commit()
	})
	if err != nil {
		handleConfigProfileWriteError(w, err, cfg)
		return
	}

	handleGetConfigProfile(w, r, manager, cfg, req.UUID)
}

func handleDeleteConfigProfile(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig, profileUUID string) {
	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		result, err := db.ExecContext(r.Context(), `DELETE FROM config_profiles WHERE uuid = ?`, profileUUID)
		if err != nil {
			return err
		}
		rows, err := result.RowsAffected()
		if err != nil {
			return err
		}
		if rows == 0 {
			return errConfigProfileNotFound
		}
		return nil
	})
	if err != nil {
		handleConfigProfileWriteError(w, err, cfg)
		return
	}
	shared.WriteJSON(w, http.StatusOK, map[string]any{"response": map[string]any{"isDeleted": true}})
}

func handleReorderConfigProfiles(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) {
	var req reorderConfigProfilesRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		shared.SendError(w, http.StatusBadRequest, "invalid JSON", err, cfg)
		return
	}
	if len(req.Items) == 0 {
		shared.SendError(w, http.StatusBadRequest, "items cannot be empty", nil, cfg)
		return
	}
	for _, item := range req.Items {
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
		for _, item := range req.Items {
			if _, err := tx.ExecContext(r.Context(), `UPDATE config_profiles SET view_position = ? WHERE uuid = ?`, item.ViewPosition, item.UUID); err != nil {
				_ = tx.Rollback()
				return err
			}
		}
		if _, err := tx.ExecContext(r.Context(), `SELECT setval('config_profiles_view_position_seq', (SELECT COALESCE(MAX(view_position), 0) FROM config_profiles) + 1)`); err != nil {
			_ = tx.Rollback()
			return err
		}
		return tx.Commit()
	})
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to reorder config profiles", err, cfg)
		return
	}

	handleGetConfigProfiles(w, r, manager, cfg)
}

func validateCreateConfigProfileRequest(req createConfigProfileRequest) error {
	name := strings.TrimSpace(req.Name)
	if len(name) < 2 || len(name) > 30 {
		return fmt.Errorf("name must be between 2 and 30 characters")
	}
	if len(req.Config) == 0 {
		return fmt.Errorf("config is required")
	}
	var parsed map[string]any
	if err := json.Unmarshal(req.Config, &parsed); err != nil {
		return fmt.Errorf("config must be valid JSON")
	}
	return nil
}

func validateUpdateConfigProfileRequest(req updateConfigProfileRequest) error {
	if req.Name == nil && req.Config == nil {
		return fmt.Errorf("no fields to update")
	}
	if req.Name != nil {
		name := strings.TrimSpace(*req.Name)
		if len(name) < 2 || len(name) > 30 {
			return fmt.Errorf("name must be between 2 and 30 characters")
		}
	}
	if req.Config != nil {
		var parsed map[string]any
		if err := json.Unmarshal(*req.Config, &parsed); err != nil {
			return fmt.Errorf("config must be valid JSON")
		}
	}
	return nil
}

func getAllConfigProfileRecords(ctx context.Context, manager *dbmanager.DatabaseManager) ([]configProfileRecord, error) {
	records := make([]configProfileRecord, 0)
	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		rows, err := db.QueryContext(ctx, `
			SELECT uuid, view_position, name, config, created_at, updated_at
			FROM config_profiles
			ORDER BY view_position ASC, name ASC
		`)
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			record, scanErr := scanConfigProfileRecord(rows)
			if scanErr != nil {
				return scanErr
			}
			records = append(records, record)
		}
		return rows.Err()
	})
	return records, err
}

func getConfigProfileRecordByUUID(ctx context.Context, manager *dbmanager.DatabaseManager, profileUUID string) (configProfileRecord, error) {
	var record configProfileRecord
	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		row := db.QueryRowContext(ctx, `
			SELECT uuid, view_position, name, config, created_at, updated_at
			FROM config_profiles
			WHERE uuid = ?
		`, profileUUID)
		var scanErr error
		record, scanErr = scanConfigProfileRecord(row)
		if errors.Is(scanErr, sql.ErrNoRows) {
			return errConfigProfileNotFound
		}
		return scanErr
	})
	return record, err
}

func scanConfigProfileRecord(scanner shared.RowScanner) (configProfileRecord, error) {
	var record configProfileRecord
	var viewPosition sql.NullInt64
	var configRaw []byte
	if err := scanner.Scan(&record.UUID, &viewPosition, &record.Name, &configRaw, &record.CreatedAt, &record.UpdatedAt); err != nil {
		return record, err
	}
	if viewPosition.Valid {
		record.ViewPosition = int(viewPosition.Int64)
	}
	record.Config = json.RawMessage(configRaw)
	return record, nil
}

func buildConfigProfileResponses(ctx context.Context, manager *dbmanager.DatabaseManager, records []configProfileRecord) ([]ConfigProfile, error) {
	profileUUIDs := make([]string, 0, len(records))
	for _, record := range records {
		profileUUIDs = append(profileUUIDs, record.UUID)
	}

	inboundsMap, err := getConfigProfileInboundsMap(ctx, manager, profileUUIDs)
	if err != nil {
		return nil, err
	}
	nodesMap, err := getConfigProfileNodesMap(ctx, manager, profileUUIDs)
	if err != nil {
		return nil, err
	}

	response := make([]ConfigProfile, 0, len(records))
	for _, record := range records {
		response = append(response, ConfigProfile{
			UUID:         record.UUID,
			ViewPosition: record.ViewPosition,
			Name:         record.Name,
			Config:       record.Config,
			Inbounds:     inboundsMap[record.UUID],
			Nodes:        nodesMap[record.UUID],
			CreatedAt:    record.CreatedAt,
			UpdatedAt:    record.UpdatedAt,
		})
	}
	return response, nil
}

func getConfigProfileInboundsMap(ctx context.Context, manager *dbmanager.DatabaseManager, profileUUIDs []string) (map[string][]ConfigProfileInbound, error) {
	result := make(map[string][]ConfigProfileInbound, len(profileUUIDs))
	if len(profileUUIDs) == 0 {
		return result, nil
	}

	activeSquadsByInbound := make(map[string][]string)
	if err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		rows, err := db.QueryContext(ctx, `
			SELECT isi.inbound_uuid, isi.internal_squad_uuid
			FROM internal_squad_inbounds isi
			JOIN config_profile_inbounds cpi ON cpi.uuid = isi.inbound_uuid
			WHERE cpi.profile_uuid = ANY(?)
		`, profileUUIDs)
		if err != nil {
			return err
		}
		defer rows.Close()
		for rows.Next() {
			var inboundUUID, squadUUID string
			if err := rows.Scan(&inboundUUID, &squadUUID); err != nil {
				return err
			}
			activeSquadsByInbound[inboundUUID] = append(activeSquadsByInbound[inboundUUID], squadUUID)
		}
		return rows.Err()
	}); err != nil {
		return nil, err
	}

	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		rows, err := db.QueryContext(ctx, `
			SELECT uuid, profile_uuid, tag, type, network, security, port, raw_inbound
			FROM config_profile_inbounds
			WHERE profile_uuid = ANY(?)
			ORDER BY tag ASC
		`, profileUUIDs)
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			var (
				item     ConfigProfileInbound
				network  sql.NullString
				security sql.NullString
				port     sql.NullInt64
				raw      []byte
			)
			if err := rows.Scan(&item.UUID, &item.ProfileUUID, &item.Tag, &item.Type, &network, &security, &port, &raw); err != nil {
				return err
			}
			if network.Valid {
				item.Network = &network.String
			}
			if security.Valid {
				item.Security = &security.String
			}
			if port.Valid {
				value := int(port.Int64)
				item.Port = &value
			}
			item.RawInbound = json.RawMessage(raw)
			item.ActiveSquads = dedupeStrings(activeSquadsByInbound[item.UUID])
			result[item.ProfileUUID] = append(result[item.ProfileUUID], item)
		}
		return rows.Err()
	})

	return result, err
}

func getConfigProfileNodesMap(ctx context.Context, manager *dbmanager.DatabaseManager, profileUUIDs []string) (map[string][]ConfigProfileNode, error) {
	result := make(map[string][]ConfigProfileNode, len(profileUUIDs))
	if len(profileUUIDs) == 0 {
		return result, nil
	}

	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		rows, err := db.QueryContext(ctx, `
			SELECT active_config_profile_uuid, uuid, name, country_code
			FROM nodes
			WHERE active_config_profile_uuid = ANY(?)
			ORDER BY view_position ASC, name ASC
		`, profileUUIDs)
		if err != nil {
			return err
		}
		defer rows.Close()
		for rows.Next() {
			var profileUUID string
			var node ConfigProfileNode
			if err := rows.Scan(&profileUUID, &node.UUID, &node.Name, &node.CountryCode); err != nil {
				return err
			}
			result[profileUUID] = append(result[profileUUID], node)
		}
		return rows.Err()
	})

	return result, err
}

func handleConfigProfileWriteError(w http.ResponseWriter, err error, cfg *config.BackendConfig) {
	switch {
	case errors.Is(err, errConfigProfileNotFound):
		shared.SendError(w, http.StatusNotFound, "config profile not found", nil, cfg)
	case strings.Contains(err.Error(), "no fields to update"):
		shared.SendError(w, http.StatusBadRequest, "no fields to update", nil, cfg)
	default:
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			shared.SendError(w, http.StatusConflict, "config profile name already exists or inbound tags are not unique", nil, cfg)
			return
		}
		if strings.Contains(err.Error(), "UNIQUE constraint failed") || strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			shared.SendError(w, http.StatusConflict, "config profile name already exists or inbound tags are not unique", nil, cfg)
			return
		}
		shared.SendError(w, http.StatusInternalServerError, "failed to write config profile", err, cfg)
	}
}

func mapConfigProfileWriteError(err error) error {
	return err
}

func parseConfigInbounds(profileUUID string, configJSON json.RawMessage) ([]ConfigProfileInbound, error) {
	var configData map[string]any
	if err := json.Unmarshal(configJSON, &configData); err != nil {
		return nil, fmt.Errorf("failed to parse config JSON: %w", err)
	}

	inboundsRaw, ok := configData["inbounds"]
	if !ok {
		return []ConfigProfileInbound{}, nil
	}
	inboundsArray, ok := inboundsRaw.([]any)
	if !ok {
		return nil, fmt.Errorf("inbounds must be an array")
	}

	seenTags := make(map[string]struct{})
	result := make([]ConfigProfileInbound, 0, len(inboundsArray))

	for _, inboundRaw := range inboundsArray {
		inboundMap, ok := inboundRaw.(map[string]any)
		if !ok {
			continue
		}

		tag, _ := inboundMap["tag"].(string)
		if strings.TrimSpace(tag) == "" {
			continue
		}
		if _, ok := seenTags[tag]; ok {
			continue
		}
		seenTags[tag] = struct{}{}

		item := ConfigProfileInbound{
			UUID:        uuid.NewString(),
			ProfileUUID: profileUUID,
			Tag:         tag,
			ActiveSquads: []string{},
		}
		if typeValue, ok := inboundMap["type"].(string); ok {
			item.Type = typeValue
		} else if protocolValue, ok := inboundMap["protocol"].(string); ok {
			item.Type = protocolValue
		}
		if networkValue, ok := inboundMap["network"].(string); ok && networkValue != "" {
			item.Network = &networkValue
		}
		if securityValue, ok := inboundMap["security"].(string); ok && securityValue != "" {
			item.Security = &securityValue
		}
		if portValue, ok := inboundMap["listen_port"].(float64); ok {
			p := int(portValue)
			item.Port = &p
		} else if portValue, ok := inboundMap["port"].(float64); ok {
			p := int(portValue)
			item.Port = &p
		}
		rawInbound, err := json.Marshal(inboundMap)
		if err != nil {
			continue
		}
		item.RawInbound = rawInbound
		result = append(result, item)
	}

	return result, nil
}

func syncConfigProfileInbounds(ctx context.Context, db dbmanager.TxExecutor, profileUUID string, configJSON json.RawMessage, cfg *config.BackendConfig) (int, error) {
	inbounds, err := parseConfigInbounds(profileUUID, configJSON)
	if err != nil {
		return 0, err
	}

	if _, err := db.ExecContext(ctx, `DELETE FROM config_profile_inbounds WHERE profile_uuid = ?`, profileUUID); err != nil {
		return 0, err
	}

	for _, inbound := range inbounds {
		var networkVal, securityVal, portVal any
		if inbound.Network != nil {
			networkVal = *inbound.Network
		}
		if inbound.Security != nil {
			securityVal = *inbound.Security
		}
		if inbound.Port != nil {
			portVal = *inbound.Port
		}

		if _, err := db.ExecContext(ctx, `
			INSERT INTO config_profile_inbounds (
				uuid, profile_uuid, tag, type, network, security, port, raw_inbound
			) VALUES (?, ?, ?, ?, ?, ?, ?, ?)
		`, inbound.UUID, inbound.ProfileUUID, inbound.Tag, inbound.Type, networkVal, securityVal, portVal, inbound.RawInbound); err != nil {
			return 0, err
		}
	}

	if cfg != nil && cfg.Logger != nil {
		cfg.Logger.Debug("Synced config profile inbounds", "profile_uuid", profileUUID, "inbounds_count", len(inbounds))
	}
	return len(inbounds), nil
}

func dedupeStrings(values []string) []string {
	seen := make(map[string]struct{}, len(values))
	result := make([]string, 0, len(values))
	for _, value := range values {
		if _, ok := seen[value]; ok {
			continue
		}
		seen[value] = struct{}{}
		result = append(result, value)
	}
	return result
}

func _unused(_ context.Context, _ dbutil.StringArray) {}
