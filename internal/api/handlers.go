package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"unicode/utf8"

	"github.com/Michael-cmd-sys/chirpy/internal/lib"
)

func (cfg *ApiConfig) MiddlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-cache")
		cfg.FileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *ApiConfig) MetricsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	html := fmt.Sprintf(
		`
		<!doctype html/>
		<html>
		  <meta charset="utf-8" />
		  <meta name="viewport" content="width=device-width,initial-scale=1.0 />
		  <body style="min-height: 100vh; width: 100vw; display: flex; flex-direction: column; align-items: center; justify-content: center;">
        <h1>Welcome, Chirpy Admin</h1>
		    <p style="flex: 1">Chirpy has been visited %d times!</p>
      </body>
    </html>
		`, cfg.FileserverHits.Load())
	_, err := fmt.Fprintf(w, "%s", html)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (cfg *ApiConfig) MetricsResetHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	cfg.FileserverHits.Store(0)
}

func HealthHander(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprintf(w, "OK")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func ValidateChirpHandler(w http.ResponseWriter, r *http.Request) {
	var parsedString string
	foundProfanity := false
	profaneWords := []string{"kerfuffle", "sharbert", "fornax"}

	type chirp struct {
		Body string `json:"body"`
	}

	payload := chirp{}
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		lib.SendJsonResponse(w, map[string]string{"error": "Something went wrong"}, http.StatusBadRequest)
		return
	}

	if len(payload.Body) == 0 {
		lib.SendJsonResponse(w, map[string]bool{"valid": false}, http.StatusBadRequest)
		return
	}

	if utf8.RuneCountInString(payload.Body) > 140 {
		lib.SendJsonResponse(w, map[string]string{"error": "Chirp is too long"}, http.StatusBadRequest)
		return
	}

	for _, word := range profaneWords {
		if strings.Contains(strings.ToLower(payload.Body), word) {
			foundProfanity = true
			parsedString = strings.ReplaceAll(strings.ToLower(payload.Body), word, "****")
		}
	}

	if foundProfanity {
		lib.SendJsonResponse(w, map[string]string{"cleaned_body": parsedString}, http.StatusOK)
		return
	}

	lib.SendJsonResponse(w, map[string]bool{"valid": true}, http.StatusOK)
}
