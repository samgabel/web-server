package main

import (
	"log"
	"net/http"
)

func middlewareLogging(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s Method: '%s' Pattern: '%s'", r.RemoteAddr, r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	}
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits++
		next.ServeHTTP(w, r)
	}
}
