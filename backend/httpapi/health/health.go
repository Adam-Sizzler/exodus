package health

import (
	"encoding/json"
	"net/http"
	"time"
	"v2ray-stat/constant"
)

// HealthHandler handles backend liveness checks.
func HealthHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.Header().Set("Allow", http.MethodGet)
			w.WriteHeader(http.StatusMethodNotAllowed)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "method not allowed"})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"status":    "ok",
			"service":   "v2rs-panel",
			"version":   constant.Version,
			"timestamp": time.Now().UTC().Format(time.RFC3339),
		})
	}
}
