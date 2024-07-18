package main

import (
	"flag"
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

	// create a --debug flag for the binary to wipe the database before startup
	dbg := flag.Bool("debug", false, "Enable debug mode")
	flag.Parse()
	if *dbg {
		if err := db.WipeDB(); err != nil {
			log.Printf("Unable to wipe databse on --debug flag: %s", err)
		}
	}

	mux.Handle("/app/*", cfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))))
	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("GET /admin/metrics", cfg.handlerMetrics)
	mux.HandleFunc("GET /api/reset", cfg.handlerResetMetrics)
	// change these to a curried function(*database.DB) returning a http.HandlerFunc
	mux.HandleFunc("POST /api/chirps", handlerPostChirp(db))
	mux.HandleFunc("GET /api/chirps", handlerGetChirps(db))
	mux.HandleFunc("GET /api/chirps/{chirpID}", handlerGetChirpByID(db))
	mux.HandleFunc("POST /api/users", handlerPostUser(db))

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: middlewareLogging(mux),
	}

	log.Printf("Serving on %s\n", srv.Addr)
	log.Fatal(srv.ListenAndServe())
}
