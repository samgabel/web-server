package main

import (
	"log"
	"net/http"

	"github.com/samgabel/web-server/internal/database"
)

const (
	port         = "8080"
	filepathRoot = "./"
)

func main() {
	mux := http.NewServeMux()

	cfg := newAPIConfig()

	// initialize new database
	db, err := database.NewDB("database.json")
	if err != nil {
		log.Fatalf("Database failed to initialize: %s", err)
	}

	mux.Handle("/app/*", cfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))))
	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("GET /admin/metrics", cfg.handlerMetrics)
	mux.HandleFunc("GET /api/reset", cfg.handlerResetMetrics)
	// change these to a curried function(*database.DB) returning a http.HandlerFunc
	mux.HandleFunc("POST /api/chirps", handlerPostChirp(db))
	mux.HandleFunc("GET /api/chirps", handlerGetChirps(db))
	mux.HandleFunc("GET /api/chirps/{chirpID}", handlerGetChirpByID(db))
	mux.HandleFunc("DELETE /api/chirps", handlerDeleteChirps(db))

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: middlewareLogging(mux),
	}

	log.Printf("Serving on %s\n", srv.Addr)
	log.Fatal(srv.ListenAndServe())
}
