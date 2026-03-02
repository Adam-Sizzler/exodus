package users

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"v2ray-stat/backend/panel/config"
	dbmanager "v2ray-stat/backend/panel/db/manager"
	"v2ray-stat/backend/panel/httpapi/shared"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
)

// UserEntity represents a user entity from the users table for API responses and requests.
type UserEntity struct {
	TID                    int64      `json:"t_id"`
	UUID                   string     `json:"uuid"`
	ShortUUID              string     `json:"short_uuid"`
	Username               string     `json:"username"`
	Status                 string     `json:"status"`
	TrafficLimitBytes      int64      `json:"traffic_limit_bytes"`
	TrafficLimitStrategy   string     `json:"traffic_limit_strategy"`
	ExpireAt               time.Time  `json:"expire_at"`
	SubLastUserAgent       *string    `json:"sub_last_user_agent,omitempty"`
	SubLastOpenedAt        *time.Time `json:"sub_last_opened_at,omitempty"`
	LastTrafficResetAt     *time.Time `json:"last_traffic_reset_at,omitempty"`
	SubRevokedAt           *time.Time `json:"sub_revoked_at,omitempty"`
	TrojanPassword         string     `json:"trojan_password"`
	VlessUUID              string     `json:"vless_uuid"`
	SSPassword             string     `json:"ss_password"`
	Description            *string    `json:"description,omitempty"`
	Tag                    *string    `json:"tag,omitempty"`
	TelegramID             *int64     `json:"telegram_id,omitempty"`
	Email                  *string    `json:"email,omitempty"`
	HwidDeviceLimit        *int       `json:"hwid_device_limit,omitempty"`
	InternalSquadUUID      *string    `json:"internal_squad_uuid,omitempty"`
	LastTriggeredThreshold int        `json:"last_triggered_threshold"`
	CreatedAt              time.Time  `json:"created_at"`
	UpdatedAt              time.Time  `json:"updated_at"`
}

// UserCreateRequest represents a request to create a new user.
type UserCreateRequest struct {
	UUID                   *string `json:"uuid,omitempty"`            // Optional, will be generated if not provided
	Username               string  `json:"username"`                  // Required, unique
	Status                 string  `json:"status"`                    // Required: ACTIVE, DISABLED, LIMITED, EXPIRED
	TrafficLimitBytes      int64   `json:"traffic_limit_bytes"`       // Default: 0
	TrafficLimitStrategy   string  `json:"traffic_limit_strategy"`    // Default: NO_RESET
	ExpireAt               string  `json:"expire_at"`                 // Required: ISO 8601 format
	TrojanPassword         *string `json:"trojan_password,omitempty"` // Optional, will be generated if not provided
	VlessUUID              *string `json:"vless_uuid,omitempty"`      // Optional, will be generated if not provided
	SSPassword             *string `json:"ss_password,omitempty"`     // Optional, will be generated if not provided
	Description            *string `json:"description,omitempty"`
	Tag                    *string `json:"tag,omitempty"`
	TelegramID             *int64  `json:"telegram_id,omitempty"`
	Email                  *string `json:"email,omitempty"`
	HwidDeviceLimit        *int    `json:"hwid_device_limit,omitempty"`
	InternalSquadUUID      *string `json:"internal_squad_uuid,omitempty"`
	LastTriggeredThreshold int     `json:"last_triggered_threshold"` // Default: 0
}

// Validate validates the UserCreateRequest fields.
func (r *UserCreateRequest) Validate() error {
	if r.Username == "" {
		return fmt.Errorf("username is required")
	}

	// Validate username format (alphanumeric, underscore, hyphen)
	validUsername := regexp.MustCompile(`^[a-zA-Z0-9_-]+$`).MatchString
	if !validUsername(r.Username) {
		return fmt.Errorf("username can only contain letters, numbers, underscores, and hyphens")
	}

	if r.Status == "" {
		return fmt.Errorf("status is required")
	}
	validStatuses := []string{"ACTIVE", "DISABLED", "LIMITED", "EXPIRED"}
	valid := false
	for _, s := range validStatuses {
		if r.Status == s {
			valid = true
			break
		}
	}
	if !valid {
		return fmt.Errorf("status must be one of: ACTIVE, DISABLED, LIMITED, EXPIRED")
	}

	if r.ExpireAt == "" {
		return fmt.Errorf("expire_at is required")
	}
	if _, err := time.Parse(time.RFC3339, r.ExpireAt); err != nil {
		return fmt.Errorf("expire_at must be in ISO 8601 format (RFC3339)")
	}

	if r.TrafficLimitStrategy == "" {
		r.TrafficLimitStrategy = "NO_RESET"
	} else {
		validStrategies := []string{"NO_RESET", "DAY", "WEEK", "MONTH"}
		valid := false
		for _, s := range validStrategies {
			if r.TrafficLimitStrategy == s {
				valid = true
				break
			}
		}
		if !valid {
			return fmt.Errorf("traffic_limit_strategy must be one of: NO_RESET, DAY, WEEK, MONTH")
		}
	}

	if r.TrafficLimitBytes < 0 {
		return fmt.Errorf("traffic_limit_bytes must be non-negative")
	}

	if r.HwidDeviceLimit != nil && *r.HwidDeviceLimit < 0 {
		return fmt.Errorf("hwid_device_limit must be non-negative")
	}

	if r.LastTriggeredThreshold < 0 || r.LastTriggeredThreshold > 100 {
		return fmt.Errorf("last_triggered_threshold must be between 0 and 100")
	}

	if r.Email != nil && *r.Email != "" {
		if !strings.Contains(*r.Email, "@") {
			return fmt.Errorf("email must be a valid email address")
		}
	}

	return nil
}

// UserUpdateRequest represents a partial update request for a user.
// Only provided fields will be updated (PATCH semantics).
type UserUpdateRequest struct {
	Status                 *string `json:"status,omitempty"`
	TrafficLimitBytes      *int64  `json:"traffic_limit_bytes,omitempty"`
	TrafficLimitStrategy   *string `json:"traffic_limit_strategy,omitempty"`
	ExpireAt               *string `json:"expire_at,omitempty"` // ISO 8601 format
	Description            *string `json:"description,omitempty"`
	Tag                    *string `json:"tag,omitempty"`
	TelegramID             *int64  `json:"telegram_id,omitempty"`
	Email                  *string `json:"email,omitempty"`
	HwidDeviceLimit        *int    `json:"hwid_device_limit,omitempty"`
	InternalSquadUUID      *string `json:"internal_squad_uuid,omitempty"`
	LastTriggeredThreshold *int    `json:"last_triggered_threshold,omitempty"`
}

var (
	errUsernameAlreadyExists = errors.New("username already exists")
	errShortUUIDCollision    = errors.New("short_uuid collision")
)

// generateRandomString generates a random hex string of specified length.
func generateRandomString(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes)[:length], nil
}

// generateSubscriptionShortUUID generates a URL-safe short id for subscription links.
// 12 random bytes produce 16 chars in base64 RawURL encoding, e.g. "xefMxk3ScgY5xhoe".
func generateSubscriptionShortUUID() (string, error) {
	bytes := make([]byte, 12)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(bytes), nil
}

// UsersCreateHandler handles POST /api/v1/users-list requests.
func UsersCreateHandler(manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cfg.Logger.Debug("Received UsersCreate HTTP request", "method", r.Method, "path", r.URL.Path)

		if r.Method != http.MethodPost {
			cfg.Logger.Warn("Invalid HTTP method", "method", r.Method, "expected", http.MethodPost)
			http.Error(w, `{"error": "method not allowed, use POST"}`, http.StatusMethodNotAllowed)
			return
		}

		// Parse JSON request body
		var req UserCreateRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			cfg.Logger.Error("Failed to parse JSON request", "error", err)
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "failed to parse JSON request body",
			})
			return
		}

		// Validate request
		if err := req.Validate(); err != nil {
			cfg.Logger.Warn("Validation failed for user create", "error", err)
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{
				"error": err.Error(),
			})
			return
		}

		ctx := r.Context()

		// Generate UUIDs and passwords if not provided
		userUUID := req.UUID
		if userUUID == nil || *userUUID == "" {
			newUUID := uuid.New().String()
			userUUID = &newUUID
		}

		shortUUID, err := generateSubscriptionShortUUID()
		if err != nil {
			cfg.Logger.Error("Failed to generate short_uuid", "error", err)
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "failed to generate short_uuid",
			})
			return
		}

		trojanPassword := req.TrojanPassword
		if trojanPassword == nil || *trojanPassword == "" {
			tp, err := generateRandomString(16)
			if err != nil {
				cfg.Logger.Error("Failed to generate trojan password", "error", err)
				w.Header().Set("Content-Type", "application/json; charset=utf-8")
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(map[string]string{
					"error": "failed to generate trojan password",
				})
				return
			}
			trojanPassword = &tp
		}

		vlessUUID := req.VlessUUID
		if vlessUUID == nil || *vlessUUID == "" {
			vUUID := uuid.New().String()
			vlessUUID = &vUUID
		}

		ssPassword := req.SSPassword
		if ssPassword == nil || *ssPassword == "" {
			sp, err := generateRandomString(16)
			if err != nil {
				cfg.Logger.Error("Failed to generate SS password", "error", err)
				w.Header().Set("Content-Type", "application/json; charset=utf-8")
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(map[string]string{
					"error": "failed to generate SS password",
				})
				return
			}
			ssPassword = &sp
		}

		// Parse expire_at
		expireAt, err := time.Parse(time.RFC3339, req.ExpireAt)
		if err != nil {
			cfg.Logger.Error("Failed to parse expire_at", "error", err)
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "expire_at must be in ISO 8601 format",
			})
			return
		}

		// Insert user into database
		var userTID int64
		const maxShortUUIDAttempts = 5
		for attempt := 0; attempt < maxShortUUIDAttempts; attempt++ {
			err = manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
				query := `
					INSERT INTO users (
						uuid, short_uuid, username, status, traffic_limit_bytes,
						traffic_limit_strategy, expire_at, trojan_password, vless_uuid, ss_password,
						description, tag, telegram_id, email, hwid_device_limit,
						last_triggered_threshold, created_at, updated_at
					) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
				`

				err := db.QueryRowContext(ctx, query+` RETURNING t_id`,
					*userUUID, shortUUID, req.Username, req.Status, req.TrafficLimitBytes,
					req.TrafficLimitStrategy, expireAt.Format("2006-01-02 15:04:05"), *trojanPassword, *vlessUUID, *ssPassword,
					req.Description, req.Tag, req.TelegramID, req.Email, req.HwidDeviceLimit,
					req.LastTriggeredThreshold,
				).Scan(&userTID)
				if err != nil {
					msg := err.Error()
					var pgErr *pgconn.PgError
					if errors.As(err, &pgErr) && pgErr.Code == "23505" {
						switch pgErr.ConstraintName {
						case "users_username_key":
							return errUsernameAlreadyExists
						case "users_short_uuid_key":
							return errShortUUIDCollision
						}
					}
					if strings.Contains(msg, "UNIQUE constraint failed: users.username") {
						return errUsernameAlreadyExists
					}
					if strings.Contains(msg, "UNIQUE constraint failed: users.short_uuid") {
						return errShortUUIDCollision
					}
					if strings.Contains(msg, "UNIQUE constraint failed") {
						return fmt.Errorf("failed to insert user: %w", err)
					}
					return fmt.Errorf("failed to insert user: %w", err)
				}

				// Handle internal squad assignment
				if req.InternalSquadUUID != nil && *req.InternalSquadUUID != "" {
					squadQuery := `INSERT INTO internal_squad_members (internal_squad_uuid, user_id) VALUES (?, ?)`
					_, err = db.ExecContext(ctx, squadQuery, *req.InternalSquadUUID, userTID)
					if err != nil {
						return fmt.Errorf("failed to assign squad: %w", err)
					}
				}

				return nil
			})

			if err == nil {
				break
			}

			if errors.Is(err, errShortUUIDCollision) {
				shortUUID, err = generateSubscriptionShortUUID()
				if err != nil {
					cfg.Logger.Error("Failed to regenerate short_uuid", "error", err)
					w.Header().Set("Content-Type", "application/json; charset=utf-8")
					w.WriteHeader(http.StatusInternalServerError)
					json.NewEncoder(w).Encode(map[string]string{
						"error": "failed to generate short_uuid",
					})
					return
				}
				continue
			}

			break
		}

		if err != nil {
			if errors.Is(err, errUsernameAlreadyExists) {
				cfg.Logger.Warn("Username already exists", "username", req.Username)
				w.Header().Set("Content-Type", "application/json; charset=utf-8")
				w.WriteHeader(http.StatusConflict)
				json.NewEncoder(w).Encode(map[string]string{
					"error": "username already exists",
				})
				return
			}
			if errors.Is(err, errShortUUIDCollision) {
				cfg.Logger.Error("Failed to create user due to short_uuid collisions", "username", req.Username)
				w.Header().Set("Content-Type", "application/json; charset=utf-8")
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(map[string]string{
					"error": "failed to generate unique short_uuid",
				})
				return
			}

			cfg.Logger.Error("Failed to create user", "error", err)
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "failed to create user",
			})
			return
		}

		cfg.Logger.Info("User created successfully", "username", req.Username, "uuid", *userUUID)

		// Fetch created user to return
		var createdUser UserEntity
		err = manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
			query := `
				SELECT 
					u.t_id, u.uuid, u.short_uuid, u.username, u.status, u.traffic_limit_bytes,
					u.traffic_limit_strategy, u.expire_at, u.sub_last_user_agent, u.sub_last_opened_at,
					u.last_traffic_reset_at, u.sub_revoked_at, u.trojan_password, u.vless_uuid,
					u.ss_password, u.description, u.tag, u.telegram_id, u.email, u.hwid_device_limit,
					ism.internal_squad_uuid, u.last_triggered_threshold, u.created_at, u.updated_at
				FROM users u
				LEFT JOIN internal_squad_members ism ON u.t_id = ism.user_id
				WHERE u.uuid = ?
			`

			row := db.QueryRowContext(ctx, query, *userUUID)

			var subLastUserAgent, description, tag, email, internalSquadUUID sql.NullString
			var subLastOpenedAt, lastTrafficResetAt, subRevokedAt sql.NullTime
			var telegramID, hwidDeviceLimit sql.NullInt64

			err := row.Scan(
				&createdUser.TID, &createdUser.UUID, &createdUser.ShortUUID, &createdUser.Username, &createdUser.Status, &createdUser.TrafficLimitBytes,
				&createdUser.TrafficLimitStrategy, &createdUser.ExpireAt, &subLastUserAgent, &subLastOpenedAt,
				&lastTrafficResetAt, &subRevokedAt, &createdUser.TrojanPassword, &createdUser.VlessUUID,
				&createdUser.SSPassword, &description, &tag, &telegramID, &email, &hwidDeviceLimit,
				&internalSquadUUID, &createdUser.LastTriggeredThreshold, &createdUser.CreatedAt, &createdUser.UpdatedAt,
			)

			if err != nil {
				return fmt.Errorf("failed to fetch created user: %w", err)
			}

			// Handle nullable fields
			if subLastUserAgent.Valid {
				createdUser.SubLastUserAgent = &subLastUserAgent.String
			}
			if subLastOpenedAt.Valid {
				createdUser.SubLastOpenedAt = &subLastOpenedAt.Time
			}
			if lastTrafficResetAt.Valid {
				createdUser.LastTrafficResetAt = &lastTrafficResetAt.Time
			}
			if subRevokedAt.Valid {
				createdUser.SubRevokedAt = &subRevokedAt.Time
			}
			if description.Valid {
				createdUser.Description = &description.String
			}
			if tag.Valid {
				createdUser.Tag = &tag.String
			}
			if telegramID.Valid {
				tid := int64(telegramID.Int64)
				createdUser.TelegramID = &tid
			}
			if email.Valid {
				createdUser.Email = &email.String
			}
			if hwidDeviceLimit.Valid {
				hdl := int(hwidDeviceLimit.Int64)
				createdUser.HwidDeviceLimit = &hdl
			}
			if internalSquadUUID.Valid {
				createdUser.InternalSquadUUID = &internalSquadUUID.String
			}

			return nil
		})

		if err != nil {
			cfg.Logger.Error("Failed to fetch created user", "uuid", *userUUID, "error", err)
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"user":    createdUser,
			"message": "user created successfully",
		})
	}
}

// Validate validates the UserUpdateRequest fields.
func (r *UserUpdateRequest) Validate() error {
	if r.Status != nil {
		validStatuses := []string{"ACTIVE", "DISABLED", "LIMITED", "EXPIRED"}
		valid := false
		for _, s := range validStatuses {
			if *r.Status == s {
				valid = true
				break
			}
		}
		if !valid {
			return fmt.Errorf("status must be one of: ACTIVE, DISABLED, LIMITED, EXPIRED")
		}
	}

	if r.TrafficLimitStrategy != nil {
		validStrategies := []string{"NO_RESET", "DAY", "WEEK", "MONTH"}
		valid := false
		for _, s := range validStrategies {
			if *r.TrafficLimitStrategy == s {
				valid = true
				break
			}
		}
		if !valid {
			return fmt.Errorf("traffic_limit_strategy must be one of: NO_RESET, DAY, WEEK, MONTH")
		}
	}

	if r.ExpireAt != nil {
		if _, err := time.Parse(time.RFC3339, *r.ExpireAt); err != nil {
			return fmt.Errorf("expire_at must be in ISO 8601 format (RFC3339)")
		}
	}

	if r.TrafficLimitBytes != nil && *r.TrafficLimitBytes < 0 {
		return fmt.Errorf("traffic_limit_bytes must be non-negative")
	}

	if r.HwidDeviceLimit != nil && *r.HwidDeviceLimit < 0 {
		return fmt.Errorf("hwid_device_limit must be non-negative")
	}

	if r.LastTriggeredThreshold != nil && (*r.LastTriggeredThreshold < 0 || *r.LastTriggeredThreshold > 100) {
		return fmt.Errorf("last_triggered_threshold must be between 0 and 100")
	}

	if r.Email != nil && *r.Email != "" {
		// Simple email validation
		if !strings.Contains(*r.Email, "@") {
			return fmt.Errorf("email must be a valid email address")
		}
	}

	return nil
}

// HasUpdates checks if any field is set for update.
func (r *UserUpdateRequest) HasUpdates() bool {
	return r.Status != nil ||
		r.TrafficLimitBytes != nil ||
		r.TrafficLimitStrategy != nil ||
		r.ExpireAt != nil ||
		r.Description != nil ||
		r.Tag != nil ||
		r.TelegramID != nil ||
		r.Email != nil ||
		r.HwidDeviceLimit != nil ||
		r.InternalSquadUUID != nil ||
		r.LastTriggeredThreshold != nil
}

// UsersAPIHandler handles GET/DELETE /api/v1/users-list requests.
// Note: This is separate from the existing /api/v1/users endpoint.
func UsersAPIHandler(manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cfg.Logger.Debug("Received UsersAPI HTTP request", "method", r.Method, "path", r.URL.Path)

		switch r.Method {
		case http.MethodGet:
			// continue below
		case http.MethodDelete:
			handleDeleteUsersBulk(w, r, manager, cfg)
			return
		default:
			cfg.Logger.Warn("Invalid HTTP method", "method", r.Method, "expected", "GET or DELETE")
			http.Error(w, `{"error": "method not allowed, use GET or DELETE"}`, http.StatusMethodNotAllowed)
			return
		}

		ctx := r.Context()

		var users []UserEntity
		err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
			query := `
				SELECT 
					u.t_id, u.uuid, u.short_uuid, u.username, u.status, u.traffic_limit_bytes,
					u.traffic_limit_strategy, u.expire_at, u.sub_last_user_agent, u.sub_last_opened_at,
					u.last_traffic_reset_at, u.sub_revoked_at, u.trojan_password, u.vless_uuid,
					u.ss_password, u.description, u.tag, u.telegram_id, u.email, u.hwid_device_limit,
					ism.internal_squad_uuid, u.last_triggered_threshold, u.created_at, u.updated_at
				FROM users u
				LEFT JOIN internal_squad_members ism ON u.t_id = ism.user_id
				ORDER BY u.t_id ASC
			`

			rows, err := db.QueryContext(ctx, query)
			if err != nil {
				return fmt.Errorf("failed to query users: %w", err)
			}
			defer rows.Close()

			for rows.Next() {
				var u UserEntity
				var subLastUserAgent, description, tag, email, internalSquadUUID sql.NullString
				var subLastOpenedAt, lastTrafficResetAt, subRevokedAt sql.NullTime
				var telegramID, hwidDeviceLimit sql.NullInt64

				err := rows.Scan(
					&u.TID, &u.UUID, &u.ShortUUID, &u.Username, &u.Status, &u.TrafficLimitBytes,
					&u.TrafficLimitStrategy, &u.ExpireAt, &subLastUserAgent, &subLastOpenedAt,
					&lastTrafficResetAt, &subRevokedAt, &u.TrojanPassword, &u.VlessUUID,
					&u.SSPassword, &description, &tag, &telegramID, &email, &hwidDeviceLimit,
					&internalSquadUUID, &u.LastTriggeredThreshold, &u.CreatedAt, &u.UpdatedAt,
				)
				if err != nil {
					return fmt.Errorf("failed to scan user row: %w", err)
				}

				// Handle nullable fields
				if subLastUserAgent.Valid {
					u.SubLastUserAgent = &subLastUserAgent.String
				}
				if subLastOpenedAt.Valid {
					u.SubLastOpenedAt = &subLastOpenedAt.Time
				}
				if lastTrafficResetAt.Valid {
					u.LastTrafficResetAt = &lastTrafficResetAt.Time
				}
				if subRevokedAt.Valid {
					u.SubRevokedAt = &subRevokedAt.Time
				}
				if description.Valid {
					u.Description = &description.String
				}
				if tag.Valid {
					u.Tag = &tag.String
				}
				if telegramID.Valid {
					tid := int64(telegramID.Int64)
					u.TelegramID = &tid
				}
				if email.Valid {
					u.Email = &email.String
				}
				if hwidDeviceLimit.Valid {
					hdl := int(hwidDeviceLimit.Int64)
					u.HwidDeviceLimit = &hdl
				}
				if internalSquadUUID.Valid {
					u.InternalSquadUUID = &internalSquadUUID.String
				}

				users = append(users, u)
			}

			if err := rows.Err(); err != nil {
				return fmt.Errorf("error iterating user rows: %w", err)
			}

			return nil
		})

		if err != nil {
			cfg.Logger.Error("Failed to fetch users", "error", err)
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "failed to fetch users",
			})
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)

		response := map[string]interface{}{
			"users": users,
			"count": len(users),
		}

		if err := json.NewEncoder(w).Encode(response); err != nil {
			cfg.Logger.Error("Failed to encode JSON response", "error", err)
		}
	}
}

// UsersReorderHandler handles POST /api/v1/users-list/reorder
func UsersReorderHandler(manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
			return
		}

		var req shared.ViewPositionReorderRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "invalid JSON"})
			return
		}
		if err := req.Validate(); err != nil {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
			return shared.ApplyViewPositionReorder(r.Context(), db, "users", req.OrderedUUIDs, cfg)
		})
		if err != nil {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "failed to reorder users"})
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "users reordered",
			"count":   len(req.OrderedUUIDs),
		})
	}
}

// UserByUUIDHandler handles GET /api/v1/users-list/{uuid}, PATCH /api/v1/users-list/{uuid}, and DELETE /api/v1/users-list/{uuid} requests.
func UserByUUIDHandler(manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cfg.Logger.Debug("Received UserByUUID HTTP request", "method", r.Method, "path", r.URL.Path)

		// Extract UUID from path: /api/v1/users-list/{uuid}
		path := strings.TrimPrefix(r.URL.Path, "/api/v1/users-list/")
		userUUID := strings.TrimSpace(path)

		if userUUID == "" {
			cfg.Logger.Warn("Missing user UUID in path", "path", r.URL.Path)
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "user UUID is required in the path",
			})
			return
		}

		cfg.Logger.Debug("Processing request for user", "uuid", userUUID, "method", r.Method)

		// Validate UUID format
		if _, err := uuid.Parse(userUUID); err != nil {
			cfg.Logger.Warn("Invalid user UUID format", "uuid", userUUID, "error", err)
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "invalid UUID format",
			})
			return
		}

		ctx := r.Context()

		switch r.Method {
		case http.MethodGet:
			cfg.Logger.Debug("Handling GET request for user", "uuid", userUUID)
			handleGetUser(w, r, manager, cfg, ctx, userUUID)
		case http.MethodPatch:
			cfg.Logger.Debug("Handling PATCH request for user", "uuid", userUUID)
			handlePatchUser(w, r, manager, cfg, ctx, userUUID)
		case http.MethodDelete:
			cfg.Logger.Debug("Handling DELETE request for user", "uuid", userUUID)
			handleDeleteUser(w, r, manager, cfg, ctx, userUUID)
		default:
			cfg.Logger.Warn("Invalid HTTP method for user endpoint", "method", r.Method, "uuid", userUUID)
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusMethodNotAllowed)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "method not allowed, use GET, PATCH or DELETE",
			})
		}
	}
}

// handleGetUser handles GET request for a single user.
func handleGetUser(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig, ctx context.Context, userUUID string) {
	cfg.Logger.Debug("Fetching user details", "uuid", userUUID)

	var user UserEntity

	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		query := `
			SELECT
				u.t_id, u.uuid, u.short_uuid, u.username, u.status, u.traffic_limit_bytes,
				u.traffic_limit_strategy, u.expire_at, u.sub_last_user_agent, u.sub_last_opened_at,
				u.last_traffic_reset_at, u.sub_revoked_at, u.trojan_password, u.vless_uuid,
				u.ss_password, u.description, u.tag, u.telegram_id, u.email, u.hwid_device_limit,
				ism.internal_squad_uuid, u.last_triggered_threshold, u.created_at, u.updated_at
			FROM users u
			LEFT JOIN internal_squad_members ism ON u.t_id = ism.user_id
			WHERE u.uuid = ?
		`

		row := db.QueryRowContext(ctx, query, userUUID)

		var subLastUserAgent, description, tag, email, internalSquadUUID sql.NullString
		var subLastOpenedAt, lastTrafficResetAt, subRevokedAt sql.NullTime
		var telegramID, hwidDeviceLimit sql.NullInt64

		err := row.Scan(
			&user.TID, &user.UUID, &user.ShortUUID, &user.Username, &user.Status, &user.TrafficLimitBytes,
			&user.TrafficLimitStrategy, &user.ExpireAt, &subLastUserAgent, &subLastOpenedAt,
			&lastTrafficResetAt, &subRevokedAt, &user.TrojanPassword, &user.VlessUUID,
			&user.SSPassword, &description, &tag, &telegramID, &email, &hwidDeviceLimit,
			&internalSquadUUID, &user.LastTriggeredThreshold, &user.CreatedAt, &user.UpdatedAt,
		)

		if err == sql.ErrNoRows {
			return fmt.Errorf("user not found")
		}
		if err != nil {
			return fmt.Errorf("failed to query user: %w", err)
		}

		// Handle nullable fields
		if subLastUserAgent.Valid {
			user.SubLastUserAgent = &subLastUserAgent.String
		}
		if subLastOpenedAt.Valid {
			user.SubLastOpenedAt = &subLastOpenedAt.Time
		}
		if lastTrafficResetAt.Valid {
			user.LastTrafficResetAt = &lastTrafficResetAt.Time
		}
		if subRevokedAt.Valid {
			user.SubRevokedAt = &subRevokedAt.Time
		}
		if description.Valid {
			user.Description = &description.String
		}
		if tag.Valid {
			user.Tag = &tag.String
		}
		if telegramID.Valid {
			tid := int64(telegramID.Int64)
			user.TelegramID = &tid
		}
		if email.Valid {
			user.Email = &email.String
		}
		if hwidDeviceLimit.Valid {
			hdl := int(hwidDeviceLimit.Int64)
			user.HwidDeviceLimit = &hdl
		}
		if internalSquadUUID.Valid {
			user.InternalSquadUUID = &internalSquadUUID.String
		}

		return nil
	})

	if err != nil {
		if err.Error() == "user not found" {
			cfg.Logger.Warn("User not found", "uuid", userUUID)
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "user not found",
			})
			return
		}

		cfg.Logger.Error("Failed to fetch user", "uuid", userUUID, "error", err)
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "failed to fetch user",
		})
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"user": user,
	})
}

// handlePatchUser handles PATCH request for partial user update.
func handlePatchUser(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig, ctx context.Context, userUUID string) {
	cfg.Logger.Debug("Starting user update", "uuid", userUUID, "method", r.Method, "path", r.URL.Path)

	// Parse JSON request body
	var req UserUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		cfg.Logger.Error("Failed to parse JSON request", "uuid", userUUID, "error", err)
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "failed to parse JSON request body",
		})
		return
	}

	cfg.Logger.Debug("Parsed update request", "uuid", userUUID, "internal_squad_uuid", req.InternalSquadUUID)

	// Check if any fields are provided
	if !req.HasUpdates() {
		cfg.Logger.Warn("No update fields provided in PATCH request", "uuid", userUUID)
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "no update fields provided",
		})
		return
	}

	// Validate request
	if err := req.Validate(); err != nil {
		cfg.Logger.Warn("Validation failed for user update", "uuid", userUUID, "error", err)
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	// Build dynamic UPDATE query based on provided fields
	var updateClauses []string
	var args []interface{}
	argIndex := 1

	// Helper function to add a field to the UPDATE clause
	addField := func(column string, value interface{}) {
		updateClauses = append(updateClauses, fmt.Sprintf("%s = $%d", column, argIndex))
		args = append(args, value)
		argIndex++
	}

	// Helper function to add a nullable field (empty string = NULL)
	addNullableField := func(column string, value *string) {
		if *value == "" {
			updateClauses = append(updateClauses, fmt.Sprintf("%s = NULL", column))
		} else {
			addField(column, *value)
		}
	}

	// Process each field using a structured approach
	type fieldUpdate struct {
		condition bool
		action    func()
	}

	fieldUpdates := []fieldUpdate{
		{req.Status != nil, func() { addField("status", *req.Status) }},
		{req.TrafficLimitBytes != nil, func() { addField("traffic_limit_bytes", *req.TrafficLimitBytes) }},
		{req.TrafficLimitStrategy != nil, func() { addField("traffic_limit_strategy", *req.TrafficLimitStrategy) }},
		{req.ExpireAt != nil, func() {
			// Parse ISO 8601 date and convert to SQLite format
			if t, err := time.Parse(time.RFC3339, *req.ExpireAt); err == nil {
				addField("expire_at", t.Format("2006-01-02 15:04:05"))
			}
		}},
		{req.Description != nil, func() { addNullableField("description", req.Description) }},
		{req.Tag != nil, func() { addNullableField("tag", req.Tag) }},
		{req.TelegramID != nil, func() { addField("telegram_id", *req.TelegramID) }},
		{req.Email != nil, func() { addNullableField("email", req.Email) }},
		{req.HwidDeviceLimit != nil, func() { addField("hwid_device_limit", *req.HwidDeviceLimit) }},
		{req.LastTriggeredThreshold != nil, func() { addField("last_triggered_threshold", *req.LastTriggeredThreshold) }},
	}

	// Execute field updates
	for _, update := range fieldUpdates {
		if update.condition {
			update.action()
		}
	}

	// Add WHERE clause argument
	args = append(args, userUUID)

	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		// 1. Update users table if there are changes
		if len(updateClauses) > 0 {
			query := fmt.Sprintf(`
				UPDATE users 
				SET %s, updated_at = CURRENT_TIMESTAMP
				WHERE uuid = $%d
			`, strings.Join(updateClauses, ", "), argIndex)

			result, err := db.ExecContext(ctx, query, args...)
			if err != nil {
				return fmt.Errorf("failed to update user: %w", err)
			}

			rowsAffected, err := result.RowsAffected()
			if err != nil {
				return fmt.Errorf("failed to get rows affected: %w", err)
			}

			if rowsAffected == 0 {
				return fmt.Errorf("user not found")
			}
		}

		// 2. Handle internal squad assignment if provided
		if req.InternalSquadUUID != nil {
			cfg.Logger.Debug("Processing squad assignment", "uuid", userUUID, "squad_uuid", *req.InternalSquadUUID)

			// Get t_id first
			var tID int64
			err := db.QueryRowContext(ctx, "SELECT t_id FROM users WHERE uuid = ?", userUUID).Scan(&tID)
			if err != nil {
				cfg.Logger.Error("Failed to get user t_id for squad assignment", "uuid", userUUID, "error", err)
				return fmt.Errorf("failed to get user t_id: %w", err)
			}
			cfg.Logger.Debug("Got user t_id", "uuid", userUUID, "t_id", tID)

			// Delete existing squad assignments
			_, err = db.ExecContext(ctx, "DELETE FROM internal_squad_members WHERE user_id = ?", tID)
			if err != nil {
				cfg.Logger.Error("Failed to clear existing squad assignments", "uuid", userUUID, "t_id", tID, "error", err)
				return fmt.Errorf("failed to clear existing squad assignments: %w", err)
			}
			cfg.Logger.Debug("Cleared existing squad assignments", "uuid", userUUID, "t_id", tID)

			// Add new assignment if provided and not empty
			if *req.InternalSquadUUID != "" {
				// First verify the squad exists
				var squadExists bool
				err = db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM internal_squads WHERE uuid = ?)", *req.InternalSquadUUID).Scan(&squadExists)
				if err != nil {
					cfg.Logger.Error("Failed to check squad existence", "uuid", userUUID, "squad_uuid", *req.InternalSquadUUID, "error", err)
					return fmt.Errorf("failed to verify squad exists: %w", err)
				}
				cfg.Logger.Debug("Squad existence check", "uuid", userUUID, "squad_uuid", *req.InternalSquadUUID, "exists", squadExists)

				if !squadExists {
					cfg.Logger.Error("Squad does not exist", "uuid", userUUID, "squad_uuid", *req.InternalSquadUUID)
					return fmt.Errorf("squad does not exist: %s", *req.InternalSquadUUID)
				}

				_, err = db.ExecContext(ctx, "INSERT INTO internal_squad_members (internal_squad_uuid, user_id) VALUES (?, ?)", *req.InternalSquadUUID, tID)
				if err != nil {
					cfg.Logger.Error("Failed to assign new squad", "uuid", userUUID, "squad_uuid", *req.InternalSquadUUID, "t_id", tID, "error", err)
					return fmt.Errorf("failed to assign new squad: %w", err)
				}
				cfg.Logger.Info("Squad assigned successfully", "uuid", userUUID, "squad_uuid", *req.InternalSquadUUID)
			}
		}

		return nil
	})

	if err != nil {
		if err.Error() == "user not found" {
			cfg.Logger.Warn("User not found for update", "uuid", userUUID)
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "user not found",
			})
			return
		}

		// Log detailed error information
		cfg.Logger.Error("Failed to update user", "uuid", userUUID, "error", err, "error_type", fmt.Sprintf("%T", err))
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "failed to update user",
		})
		return
	}

	cfg.Logger.Info("User updated successfully", "uuid", userUUID)

	// Fetch updated user to return
	var updatedUser UserEntity
	err = manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		query := `
			SELECT
				u.t_id, u.uuid, u.short_uuid, u.username, u.status, u.traffic_limit_bytes,
				u.traffic_limit_strategy, u.expire_at, u.sub_last_user_agent, u.sub_last_opened_at,
				u.last_traffic_reset_at, u.sub_revoked_at, u.trojan_password, u.vless_uuid,
				u.ss_password, u.description, u.tag, u.telegram_id, u.email, u.hwid_device_limit,
				ism.internal_squad_uuid, u.last_triggered_threshold, u.created_at, u.updated_at
			FROM users u
			LEFT JOIN internal_squad_members ism ON u.t_id = ism.user_id
			WHERE u.uuid = ?
		`

		row := db.QueryRowContext(ctx, query, userUUID)

		var subLastUserAgent, description, tag, email, internalSquadUUID sql.NullString
		var subLastOpenedAt, lastTrafficResetAt, subRevokedAt sql.NullTime
		var telegramID, hwidDeviceLimit sql.NullInt64

		err := row.Scan(
			&updatedUser.TID, &updatedUser.UUID, &updatedUser.ShortUUID, &updatedUser.Username, &updatedUser.Status, &updatedUser.TrafficLimitBytes,
			&updatedUser.TrafficLimitStrategy, &updatedUser.ExpireAt, &subLastUserAgent, &subLastOpenedAt,
			&lastTrafficResetAt, &subRevokedAt, &updatedUser.TrojanPassword, &updatedUser.VlessUUID,
			&updatedUser.SSPassword, &description, &tag, &telegramID, &email, &hwidDeviceLimit,
			&internalSquadUUID, &updatedUser.LastTriggeredThreshold, &updatedUser.CreatedAt, &updatedUser.UpdatedAt,
		)

		if err != nil {
			return fmt.Errorf("failed to fetch updated user: %w", err)
		}

		// Handle nullable fields (same as in handleGetUser)
		if subLastUserAgent.Valid {
			updatedUser.SubLastUserAgent = &subLastUserAgent.String
		}
		if subLastOpenedAt.Valid {
			updatedUser.SubLastOpenedAt = &subLastOpenedAt.Time
		}
		if lastTrafficResetAt.Valid {
			updatedUser.LastTrafficResetAt = &lastTrafficResetAt.Time
		}
		if subRevokedAt.Valid {
			updatedUser.SubRevokedAt = &subRevokedAt.Time
		}
		if description.Valid {
			updatedUser.Description = &description.String
		}
		if tag.Valid {
			updatedUser.Tag = &tag.String
		}
		if telegramID.Valid {
			tid := int64(telegramID.Int64)
			updatedUser.TelegramID = &tid
		}
		if email.Valid {
			updatedUser.Email = &email.String
		}
		if hwidDeviceLimit.Valid {
			hdl := int(hwidDeviceLimit.Int64)
			updatedUser.HwidDeviceLimit = &hdl
		}
		if internalSquadUUID.Valid {
			updatedUser.InternalSquadUUID = &internalSquadUUID.String
		}

		return nil
	})

	if err != nil {
		cfg.Logger.Error("Failed to fetch updated user", "uuid", userUUID, "error", err)
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"user":    updatedUser,
		"message": "user updated successfully",
	})
}

// handleDeleteUser handles DELETE request for deleting a user.
func handleDeleteUser(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig, ctx context.Context, userUUID string) {
	cfg.Logger.Debug("Starting user deletion", "uuid", userUUID)

	// Get user info for logging before deletion
	var username string
	var userTID int64
	err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		query := "SELECT t_id, username FROM users WHERE uuid = ?"
		return db.QueryRowContext(ctx, query, userUUID).Scan(&userTID, &username)
	})

	if err != nil {
		if err == sql.ErrNoRows {
			cfg.Logger.Warn("User not found for deletion", "uuid", userUUID)
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "user not found",
			})
			return
		}

		cfg.Logger.Error("Failed to fetch user for deletion", "uuid", userUUID, "error", err)
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "failed to fetch user",
		})
		return
	}

	// Delete user (CASCADE will handle related records)
	err = manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		query := "DELETE FROM users WHERE uuid = ?"
		result, err := db.ExecContext(ctx, query, userUUID)
		if err != nil {
			return fmt.Errorf("failed to delete user: %w", err)
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			return fmt.Errorf("failed to get rows affected: %w", err)
		}

		if rowsAffected == 0 {
			return fmt.Errorf("user not found")
		}

		return nil
	})

	if err != nil {
		cfg.Logger.Error("Failed to delete user", "uuid", userUUID, "error", err)
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "failed to delete user",
		})
		return
	}

	cfg.Logger.Info("User deleted successfully", "uuid", userUUID, "username", username, "t_id", userTID)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":  "user deleted successfully",
		"uuid":     userUUID,
		"username": username,
		"t_id":     userTID,
	})
}

func handleDeleteUsersBulk(w http.ResponseWriter, r *http.Request, manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) {
	rawUUIDs := r.URL.Query().Get("uuids")
	uuids, err := shared.ParseUUIDCSV(rawUUIDs)
	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	var deleted int64
	err = manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
		tx, err := db.BeginTx(r.Context(), nil)
		if err != nil {
			return err
		}

		for _, id := range uuids {
			res, execErr := tx.ExecContext(r.Context(), "DELETE FROM users WHERE uuid = ?", id)
			if execErr != nil {
				_ = tx.Rollback()
				return execErr
			}
			rows, rowsErr := res.RowsAffected()
			if rowsErr != nil {
				_ = tx.Rollback()
				return rowsErr
			}
			deleted += rows
		}

		return tx.Commit()
	})
	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "failed to delete selected users"})
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "selected users deleted",
		"count":   deleted,
	})
}

// ==================== USERS LIST SUMMARY (FOR UI SELECTION) ====================

// UserSummary represents a brief user info for selection UI.
type UserSummary struct {
	TID      int64  `json:"t_id"`
	UUID     string `json:"uuid"`
	Username string `json:"username"`
	Email    string `json:"email,omitempty"`
	Tag      string `json:"tag,omitempty"`
	Status   string `json:"status"`
}

// UsersListSummaryHandler handles GET /api/v1/users-list/summary
func UsersListSummaryHandler(manager *dbmanager.DatabaseManager, cfg *config.BackendConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, `{"error": "method not allowed"}`, http.StatusMethodNotAllowed)
			return
		}

		ctx := r.Context()
		var users []UserSummary

		err := manager.ExecuteHighPriority(func(db dbmanager.DBExecutor) error {
			query := `
				SELECT t_id, uuid, username, email, tag, status
				FROM users
				ORDER BY t_id ASC`

			rows, err := db.QueryContext(ctx, query)
			if err != nil {
				return err
			}
			defer rows.Close()

			for rows.Next() {
				var u UserSummary
				var email, tag sql.NullString

				if err := rows.Scan(&u.TID, &u.UUID, &u.Username, &email, &tag, &u.Status); err != nil {
					return err
				}

				if email.Valid {
					u.Email = email.String
				}
				if tag.Valid {
					u.Tag = tag.String
				}

				users = append(users, u)
			}

			return rows.Err()
		})

		if err != nil {
			shared.SendError(w, http.StatusInternalServerError, "failed to fetch users summary", err, cfg)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"users": users,
			"count": len(users),
		})
	}
}
