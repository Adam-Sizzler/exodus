package users

import (
	"errors"
	"net/http"
	"strings"

	"exodus/internal/config"
	"exodus/internal/httpapi/shared"

	"github.com/jackc/pgx/v5/pgconn"
)

func handleUserWriteError(w http.ResponseWriter, err error, cfg *config.BackendConfig) {
	switch {
	case errors.Is(err, errUsernameExists):
		shared.SendError(w, http.StatusConflict, "username already exists", nil, cfg)
	case errors.Is(err, errShortUUIDExists):
		shared.SendError(w, http.StatusConflict, "short uuid already exists", nil, cfg)
	case errors.Is(err, errVLESSUUIDExists):
		shared.SendError(w, http.StatusConflict, "vless uuid already exists", nil, cfg)
	case errors.Is(err, errUserNotFound):
		shared.SendError(w, http.StatusNotFound, "user not found", nil, cfg)
	default:
		shared.SendError(w, http.StatusInternalServerError, "failed to write user", err, cfg)
	}
}

func mapUserWriteError(err error) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		switch pgErr.ConstraintName {
		case "users_username_key":
			return errUsernameExists
		case "users_short_uuid_key":
			return errShortUUIDExists
		case "users_vless_uuid_key":
			return errVLESSUUIDExists
		}
	}

	message := err.Error()
	switch {
	case strings.Contains(message, "users_username_key"), strings.Contains(message, "UNIQUE constraint failed: users.username"):
		return errUsernameExists
	case strings.Contains(message, "users_short_uuid_key"), strings.Contains(message, "UNIQUE constraint failed: users.short_uuid"):
		return errShortUUIDExists
	case strings.Contains(message, "users_vless_uuid_key"), strings.Contains(message, "UNIQUE constraint failed: users.vless_uuid"):
		return errVLESSUUIDExists
	default:
		return err
	}
}
