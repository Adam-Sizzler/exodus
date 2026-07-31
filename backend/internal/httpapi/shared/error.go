package shared

import (
	"net/http"

	"exodus/internal/config"
)

func SendError(w http.ResponseWriter, code int, msg string, err error, cfg *config.BackendConfig) {
	if err != nil && cfg != nil && cfg.Logger != nil {
		if code >= 500 {
			cfg.Logger.Error(msg, "status", code, "error", err)
		} else if code == http.StatusUnauthorized || code == http.StatusForbidden {
			cfg.Logger.Debug(msg, "status", code, "error", err)
		} else {
			cfg.Logger.Warn(msg, "status", code, "error", err)
		}
	}
	WriteJSONError(w, code, msg)
}
