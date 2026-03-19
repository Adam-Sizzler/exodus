package shared

import (
	"net/http"

	"cerberus/backend/config"
)

func SendError(w http.ResponseWriter, code int, msg string, err error, cfg *config.BackendConfig) {
	if err != nil && cfg != nil {
		cfg.Logger.Error(msg, "error", err)
	}
	WriteJSONError(w, code, msg)
}
