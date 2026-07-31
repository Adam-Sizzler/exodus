package hwiduserdevices

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"exodus/internal/config"
	"exodus/internal/httpapi/shared"
	"exodus/internal/notifications"

	"github.com/google/uuid"
)

var (
	hwidRegex = regexp.MustCompile(`^[a-zA-Z0-9_-]{8,64}$`)
)

var (
	errHWIDDeviceNotFound = errors.New("hwid device not found")
	errHWIDDeviceExists   = errors.New("hwid device already exists")
	errInvalidHWIDFormat  = errors.New("invalid hwid format")
)

type HWIDUserDeviceRecord struct {
	UUID        string    `json:"uuid"`
	UserTID     int64     `json:"user_t_id"`
	HWID        string    `json:"hwid"`
	DeviceName  string    `json:"device_name"`
	FirstSeenAt time.Time `json:"first_seen_at"`
	LastSeenAt  time.Time `json:"last_seen_at"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type HWIDUserDeviceAPI struct {
	UUID        string    `json:"uuid"`
	UserTID     int64     `json:"userTId"`
	HWID        string    `json:"hwid"`
	DeviceName  string    `json:"deviceName"`
	FirstSeenAt time.Time `json:"firstSeenAt"`
	LastSeenAt  time.Time `json:"lastSeenAt"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type CreateHWIDDeviceRequest struct {
	UserUUID   string `json:"userUuid"`
	HWID       string `json:"hwid"`
	DeviceName string `json:"deviceName"`
}

type CheckHWIDRequest struct {
	UserUUID string `json:"userUuid"`
	HWID     string `json:"hwid"`
}

type CheckHWIDResponse struct {
	Allowed      bool   `json:"allowed"`
	Reason       string `json:"reason,omitempty"`
	DeviceCount  int    `json:"deviceCount"`
	DeviceLimit  *int   `json:"deviceLimit,omitempty"`
	ExistingHWID string `json:"existingHwid,omitempty"`
}

func (r *CreateHWIDDeviceRequest) Validate() error {
	if _, err := uuid.Parse(r.UserUUID); err != nil {
		return fmt.Errorf("invalid userUuid format")
	}
	if strings.TrimSpace(r.HWID) == "" {
		return fmt.Errorf("hwid is required")
	}
	if !hwidRegex.MatchString(r.HWID) {
		return errInvalidHWIDFormat
	}
	if strings.TrimSpace(r.DeviceName) == "" {
		return fmt.Errorf("deviceName is required")
	}
	if len(r.DeviceName) > 50 {
		return fmt.Errorf("deviceName must be less than 50 characters")
	}
	return nil
}

func (r *CheckHWIDRequest) Validate() error {
	if _, err := uuid.Parse(r.UserUUID); err != nil {
		return fmt.Errorf("invalid userUuid format")
	}
	if strings.TrimSpace(r.HWID) == "" {
		return fmt.Errorf("hwid is required")
	}
	return nil
}

func HWIDUserDevicesHandler(db *sql.DB, cfg *config.BackendConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handleGetHWIDDevices(w, r, db, cfg)
		case http.MethodPost:
			handleCreateHWIDDevice(w, r, db, cfg)
		default:
			shared.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
		}
	}
}

func HWIDCheckHandler(db *sql.DB, cfg *config.BackendConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			shared.WriteJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}

		var req CheckHWIDRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			shared.SendError(w, http.StatusBadRequest, "invalid JSON", err, cfg)
			return
		}

		if err := req.Validate(); err != nil {
			shared.SendError(w, http.StatusBadRequest, err.Error(), nil, cfg)
			return
		}

		handleCheckHWID(w, r, db, cfg, req.UserUUID, req.HWID)
	}
}

func handleGetHWIDDevices(w http.ResponseWriter, r *http.Request, db *sql.DB, cfg *config.BackendConfig) {
	records, err := getHWIDDevices(r.Context(), db)
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to fetch hwid devices", err, cfg)
		return
	}

	result := make([]HWIDUserDeviceAPI, 0, len(records))
	for _, rec := range records {
		result = append(result, convertHWIDRecordToAPI(rec))
	}

	shared.WriteJSON(w, http.StatusOK, map[string]any{"response": result})
}

func handleCreateHWIDDevice(w http.ResponseWriter, r *http.Request, db *sql.DB, cfg *config.BackendConfig) {
	var req CreateHWIDDeviceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		shared.SendError(w, http.StatusBadRequest, "invalid JSON", err, cfg)
		return
	}

	if err := req.Validate(); err != nil {
		shared.SendError(w, http.StatusBadRequest, err.Error(), nil, cfg)
		return
	}

	userTID, hwidLimit, err := getUserHWIDLimit(r.Context(), db, req.UserUUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			shared.SendError(w, http.StatusNotFound, "user not found", nil, cfg)
			return
		}
		shared.SendError(w, http.StatusInternalServerError, "failed to fetch user", err, cfg)
		return
	}

	exists, err := hwidExistsForUser(r.Context(), db, userTID, req.HWID)
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to check hwid", err, cfg)
		return
	}
	if exists {
		shared.SendError(w, http.StatusConflict, "hwid already registered for this user", nil, cfg)
		return
	}

	if hwidLimit != nil && *hwidLimit > 0 {
		count, err := countUserHWIDDevices(r.Context(), db, userTID)
		if err != nil {
			shared.SendError(w, http.StatusInternalServerError, "failed to count devices", err, cfg)
			return
		}
		if count >= *hwidLimit {
			shared.SendError(w, http.StatusConflict, fmt.Sprintf("hwid device limit reached (%d)", *hwidLimit), nil, cfg)
			return
		}
	}

	deviceUUID := uuid.NewString()
	now := time.Now()

	_, err = db.ExecContext(r.Context(), `
		INSERT INTO hwid_user_devices (
			uuid, user_t_id, hwid, device_name,
			first_seen_at, last_seen_at, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	`, deviceUUID, userTID, req.HWID, req.DeviceName, now, now)
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to create hwid device", err, cfg)
		return
	}

	created, err := getHWIDDeviceByUUID(r.Context(), db, deviceUUID)
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to fetch created hwid device", err, cfg)
		return
	}

	emitHWIDNotification(r.Context(), cfg, notifications.EventUserHWIDDeviceAdded, hwidRecordNotificationData(created), nil)
	shared.WriteJSON(w, http.StatusCreated, map[string]any{"response": convertHWIDRecordToAPI(created)})
}

func handleCheckHWID(w http.ResponseWriter, r *http.Request, db *sql.DB, cfg *config.BackendConfig, userUUID, hwid string) {
	userTID, hwidLimit, err := getUserHWIDLimit(r.Context(), db, userUUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			shared.SendError(w, http.StatusNotFound, "user not found", nil, cfg)
			return
		}
		shared.SendError(w, http.StatusInternalServerError, "failed to fetch user", err, cfg)
		return
	}

	response := CheckHWIDResponse{
		DeviceLimit: hwidLimit,
	}

	exists, existingHWID, err := hwidExistsForUserWithDetail(r.Context(), db, userTID, hwid)
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to check hwid", err, cfg)
		return
	}
	if exists {
		response.Allowed = true
		response.Reason = "hwid already registered"
		response.ExistingHWID = existingHWID
		shared.WriteJSON(w, http.StatusOK, map[string]any{"response": response})
		return
	}

	count, err := countUserHWIDDevices(r.Context(), db, userTID)
	if err != nil {
		shared.SendError(w, http.StatusInternalServerError, "failed to count devices", err, cfg)
		return
	}
	response.DeviceCount = count

	if hwidLimit != nil && *hwidLimit > 0 {
		if count >= *hwidLimit {
			response.Allowed = false
			response.Reason = fmt.Sprintf("device limit reached (%d)", *hwidLimit)
			shared.WriteJSON(w, http.StatusOK, map[string]any{"response": response})
			return
		}
	}

	response.Allowed = true
	response.Reason = "hwid can be registered"
	shared.WriteJSON(w, http.StatusOK, map[string]any{"response": response})
}

func getHWIDDevices(ctx context.Context, db *sql.DB) ([]HWIDUserDeviceRecord, error) {
	rows, err := db.QueryContext(ctx, `
		SELECT uuid, user_t_id, hwid, device_name,
			first_seen_at, last_seen_at, created_at, updated_at
		FROM hwid_user_devices
		ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	records := make([]HWIDUserDeviceRecord, 0)
	for rows.Next() {
		rec, err := scanHWIDDevice(rows)
		if err != nil {
			return nil, err
		}
		records = append(records, rec)
	}

	return records, rows.Err()
}

func getHWIDDeviceByUUID(ctx context.Context, db *sql.DB, deviceUUID string) (HWIDUserDeviceRecord, error) {
	row := db.QueryRowContext(ctx, `
		SELECT uuid, user_t_id, hwid, device_name,
			first_seen_at, last_seen_at, created_at, updated_at
		FROM hwid_user_devices
		WHERE uuid = $1
	`, deviceUUID)

	return scanHWIDDevice(row)
}

func getUserHWIDLimit(ctx context.Context, db *sql.DB, userUUID string) (int64, *int, error) {
	var userTID int64
	var hwidDeviceLimit sql.NullInt64

	row := db.QueryRowContext(ctx, `
		SELECT t_id, hwid_device_limit
		FROM users
		WHERE uuid = $1
	`, userUUID)

	if err := row.Scan(&userTID, &hwidDeviceLimit); err != nil {
		return 0, nil, err
	}

	var limit *int
	if hwidDeviceLimit.Valid && hwidDeviceLimit.Int64 > 0 {
		l := int(hwidDeviceLimit.Int64)
		limit = &l
	}

	return userTID, limit, nil
}

func hwidExistsForUser(ctx context.Context, db *sql.DB, userTID int64, hwid string) (bool, error) {
	var exists bool

	row := db.QueryRowContext(ctx, `
		SELECT EXISTS(
			SELECT 1 FROM hwid_user_devices
			WHERE user_t_id = $1 AND hwid = $2
		)
	`, userTID, hwid)

	err := row.Scan(&exists)
	return exists, err
}

func hwidExistsForUserWithDetail(ctx context.Context, db *sql.DB, userTID int64, hwid string) (bool, string, error) {
	var hwidValue string

	row := db.QueryRowContext(ctx, `
		SELECT hwid FROM hwid_user_devices
		WHERE user_t_id = $1 AND hwid = $2
		LIMIT 1
	`, userTID, hwid)

	err := row.Scan(&hwidValue)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, "", nil
		}
		return false, "", err
	}

	return true, hwidValue, nil
}

func countUserHWIDDevices(ctx context.Context, db *sql.DB, userTID int64) (int, error) {
	var count int

	row := db.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM hwid_user_devices
		WHERE user_t_id = $1
	`, userTID)

	err := row.Scan(&count)
	return count, err
}

func scanHWIDDevice(scanner shared.RowScanner) (HWIDUserDeviceRecord, error) {
	var rec HWIDUserDeviceRecord

	err := scanner.Scan(
		&rec.UUID,
		&rec.UserTID,
		&rec.HWID,
		&rec.DeviceName,
		&rec.FirstSeenAt,
		&rec.LastSeenAt,
		&rec.CreatedAt,
		&rec.UpdatedAt,
	)
	if err != nil {
		return rec, err
	}

	return rec, nil
}

func convertHWIDRecordToAPI(rec HWIDUserDeviceRecord) HWIDUserDeviceAPI {
	return HWIDUserDeviceAPI{
		UUID:        rec.UUID,
		UserTID:     rec.UserTID,
		HWID:        rec.HWID,
		DeviceName:  rec.DeviceName,
		FirstSeenAt: rec.FirstSeenAt,
		LastSeenAt:  rec.LastSeenAt,
		CreatedAt:   rec.CreatedAt,
		UpdatedAt:   rec.UpdatedAt,
	}
}
