package main

import (
	//"fmt"
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-cache")
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) MetricsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	html := fmt.Sprintf(
		`<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>`, cfg.fileserverHits.Load())
	_, err := fmt.Fprintf(w, "%s", html)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (cfg *apiConfig) MetricsResetHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	cfg.fileserverHits.Store(0)
}

func healthHander(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprintf(w, "OK")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func main() {
	const FILEROOT = "./static"
	const PORT = ":8080"

	// Server config setup
	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
	}

	// Create a server instance handler
	mux := http.NewServeMux()

	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(FILEROOT)))))

	mux.HandleFunc("GET /admin/metrics", apiCfg.MetricsHandler)

	mux.HandleFunc("POST /admin/reset", apiCfg.MetricsResetHandler)

	mux.HandleFunc("GET /api/healthz", healthHander)

	srv := &http.Server{
		Addr:    PORT,
		Handler: mux,
	}

	// Start the server
	log.Printf("Server started on port http://localhost%s\n", PORT)
	log.Fatal(srv.ListenAndServe())
}
