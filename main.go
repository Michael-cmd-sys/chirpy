package main 

import (
	//"fmt"
	"net/http"
	"log"
)

func main() {
	// Create a server instance handler
	mux := http.NewServeMux()

	// Static file server instance
	fileServer := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fileServer))

	// Serve static html file to users
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/index.html")
	})

	srv := &http.Server{
		Addr   : ":8080",
		Handler: mux,
	}

	// Start the server
	log.Println("Server started on port http://localhost:8080...")
	log.Fatal(srv.ListenAndServe())
}
