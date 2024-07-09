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

	cfg := apiConfig{
		fileserverHits: 0,
	}

	mux.Handle("/app/*", cfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))))
	mux.HandleFunc("GET /healthz", handlerReadiness)
	mux.HandleFunc("GET /metrics", cfg.handlerMetrics)
	mux.HandleFunc("/reset", cfg.handlerResetMetrics)

	srv := &http.Server{
		Addr:    "localhost:" + port,
		Handler: mux,
	}

	log.Printf("Serving on port: %s\n", port)
	log.Fatal(srv.ListenAndServe())
}
