package main

import (
	"log"
	"net/http"
	"sync/atomic"

	"github.com/Michael-cmd-sys/chirpy/internal/api"
)

func main() {
	const FILEROOT = "./static"
	const PORT = ":8080"

	// Server config setup
	apiCfg := api.ApiConfig{
		FileserverHits: atomic.Int32{},
	}

	// Create a server instance handler
	mux := http.NewServeMux()

	mux.Handle("/app/", apiCfg.MiddlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(FILEROOT)))))

	mux.HandleFunc("GET /admin/metrics", apiCfg.MetricsHandler)

	mux.HandleFunc("POST /admin/reset", apiCfg.MetricsResetHandler)

	mux.HandleFunc("GET /api/healthz", api.HealthHander)

	mux.HandleFunc("POST /api/validate_chirp", api.ValidateChirpHandler)

	srv := &http.Server{
		Addr:    PORT,
		Handler: mux,
	}

	// Start the server
	log.Printf("Server started on port http://localhost%s\n", PORT)
	log.Fatal(srv.ListenAndServe())
}
