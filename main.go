package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/Michael-cmd-sys/chirpy/internal/api"
	"github.com/Michael-cmd-sys/chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	godotenv.Load()
	const FILEROOT = "./static"
	const PORT = ":8080"
	dbURL := os.Getenv("DB_URL")

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}

	dbQueries := database.New(db)

	// Server config setup
	apiCfg := api.ApiConfig{
		FileserverHits: atomic.Int32{},
		DB:             dbQueries,
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
