package shared

import (
	"encoding/json"
	"net/http"
)

func WriteJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func WriteJSONError(w http.ResponseWriter, status int, message string) {
	WriteJSON(w, status, map[string]string{
		"message": message,
		"error":   message,
	})
}
