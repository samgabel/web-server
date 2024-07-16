package main

import (
	"log"
	"net/http"
)

const (
	port         = "8080"
	filepathRoot = "./"
)

func main() {
	mux := http.NewServeMux()

	cfg := newAPIConfig()

	mux.Handle("/app/*", cfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))))
	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("GET /admin/metrics", cfg.handlerMetrics)
	mux.HandleFunc("GET /api/reset", cfg.handlerResetMetrics)
	mux.HandleFunc("POST /api/validate_chirp", handlerValidateJSON)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: middlewareLogging(mux),
	}

	log.Printf("Serving on %s\n", srv.Addr)
	log.Fatal(srv.ListenAndServe())
}
