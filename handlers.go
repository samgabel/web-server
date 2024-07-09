package main

import (
	"fmt"
	"net/http"
)

func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	// write "Content-Type" header
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	// write status code
	w.WriteHeader(http.StatusOK)
	// write "OK" to body
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	// write the number of api hits
	w.Write([]byte(fmt.Sprintf("Hits: %d", cfg.fileserverHits)))
}

func (cfg *apiConfig) handlerResetMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	// write the success of the reset
	w.Write([]byte("Reset fileserverHits to 0"))
	// reset the apiConfig struct fields
	cfg.fileserverHits = 0
}
