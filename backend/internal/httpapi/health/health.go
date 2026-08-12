package health

import (
	"encoding/json"
	"exodus/internal/constant"
	"net/http"
	"time"
)

// HealthHandler godoc
// @Summary      Backend health check
// @Description  Simple liveness and health probe for load balancers and container orchestrators
// @Tags         Health
// @Produce      json
// @Success      200  {object}  map[string]string
// @Router       /health [get]
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
			"service":   "exodus-panel",
			"version":   constant.Version,
			"timestamp": time.Now().UTC().Format(time.RFC3339),
		})
	}
}
