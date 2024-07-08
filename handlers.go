package main

import "net/http"

func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	// write "Content-Type" header
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	// write status code
	w.WriteHeader(http.StatusOK)
	// write "OK" to body
	w.Write([]byte(http.StatusText(http.StatusOK)))
}
