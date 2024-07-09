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
	// decouple presentation logic from api logic by providing an /api prefix path
	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("GET /api/metrics", cfg.handlerMetrics)
	mux.HandleFunc("GET /api/reset", cfg.handlerResetMetrics)

	srv := &http.Server{
		Addr:    "localhost:" + port,
		Handler: mux,
	}

	log.Printf("Serving on port: %s\n", port)
	log.Fatal(srv.ListenAndServe())
}
