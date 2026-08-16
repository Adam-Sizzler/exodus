package shared

import (
	"net/http"

	"exodus/internal/config"
)

func SendError(w http.ResponseWriter, code int, msg string, err error, cfg *config.BackendConfig) {
	if err != nil && cfg != nil && cfg.Logger != nil {
		if code >= 500 {
			cfg.Logger.Error(msg, "error", err)
		} else {
			cfg.Logger.Debug(msg, "error", err)
		}
	}
	WriteJSONError(w, code, msg)
}
