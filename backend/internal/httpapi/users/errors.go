package users

import (
	"errors"
	"net/http"

	"exodus/internal/config"
	"exodus/internal/httpapi/shared"
	"exodus/internal/util"
)

func handleUserWriteError(w http.ResponseWriter, err error, cfg *config.BackendConfig) {
	switch {
	case errors.Is(err, errUsernameExists):
		shared.SendAPIError(w, shared.ErrUsernameAlreadyExists, cfg)
	case errors.Is(err, errShortUUIDExists):
		shared.SendAPIError(w, shared.ErrShortUUIDAlreadyExists, cfg)
	case errors.Is(err, errVLESSUUIDExists):
		shared.SendAPIError(w, shared.ErrVlessUUIDAlreadyExists, cfg)
	case errors.Is(err, errUserNotFound):
		shared.SendAPIError(w, shared.ErrUserNotFound, cfg)
	default:
		shared.SendAPIError(w, shared.ErrUserWriteFailed.WithCause(err), cfg)
	}
}

func mapUserWriteError(err error) error {
	switch {
	case util.IsUniqueViolation(err, "users_username_key"):
		return errUsernameExists
	case util.IsUniqueViolation(err, "users_short_uuid_key"):
		return errShortUUIDExists
	case util.IsUniqueViolation(err, "users_vless_uuid_key"):
		return errVLESSUUIDExists
	default:
		return err
	}
}
