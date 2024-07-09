package main

import "net/http"

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	// we need to modify the handler being fed into this middleware with http.HandlerFunc
	// this is because it is the http.Handler itself that gets called when a specific endpoint is requested
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits++
		// this will serve an http request with a http response writer (of the http.Handler function being fed into this middleware)
		next.ServeHTTP(w, r)
	})
}
