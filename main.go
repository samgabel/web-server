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
	// define the request multiplexer
	mux := http.NewServeMux()

	// define api config instance
	cfg := apiConfig{
		fileserverHits: 0,
	}

	// register a STATEFUL fileserver handler
	// implement the middleware to wrap the fileserver handler
	mux.Handle("/app/*", cfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))))
	// register a custom handler for readiness endpoint
	// limit HTTP method access to 'GET' only
	mux.HandleFunc("GET /healthz", handlerReadiness)
	// register a custom handler for metrics endpoint
	// limit HTTP method access to 'GET' only
	mux.HandleFunc("GET /metrics", cfg.handlerMetrics)
	// register a custom handler for metrics endpoint
	mux.HandleFunc("/reset", cfg.handlerResetMetrics)

	// define the http server
	srv := &http.Server{
		Addr:    "localhost:" + port,
		Handler: mux,
	}

	// print report
	log.Printf("Serving on port: %s\n", port)
	// listens on address and calls serve
	// log.Fatal will not close the program (os.Exit(1)) unless the server is shutdown or closed
	log.Fatal(srv.ListenAndServe())
}
