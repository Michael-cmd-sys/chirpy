package lib

import (
	"encoding/json"
	"net/http"
)

func SendJsonResponse(w http.ResponseWriter, payload any, statusCode int) {
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(statusCode)

	err := json.NewEncoder(w).Encode(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "Service unavailable",
		})
	}
}
