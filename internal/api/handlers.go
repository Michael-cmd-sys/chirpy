package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
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
	type errorResponse struct {
		Error string `json:"error"`
	}

	type chirp struct {
		Body string `json:"body"`
	}

	type successResponse struct {
		Valid bool `json:"valid"`
	}

	payload := chirp{}
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		log.Println(err)
		lib.SendJsonResponse(w, errorResponse{"Something went wrong"}, http.StatusBadRequest)
		return
	}

	log.Printf("Received chirp: %s", payload.Body)

	if utf8.RuneCountInString(payload.Body) > 140 {
		log.Println("Error occurred, chirp too long")
		lib.SendJsonResponse(w, errorResponse{"Chirp is too long"}, http.StatusBadRequest)
		return
	}

  lib.SendJsonResponse(w, successResponse{true}, http.StatusOK)
}

