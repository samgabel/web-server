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

	// register a fileserver handler
	mux.Handle("/app/*", http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))
	// register a custom handler for readiness endpoint
	mux.HandleFunc("/healthz", handlerReadiness)

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
