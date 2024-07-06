package main

import (
	"log"
	"net/http"
)

func main() {
	const port = "8080"
	const filepathRoot = "./"

	// define the request multiplexer
	mux := http.NewServeMux()

	// add handler for root path
	// our handler is a FileServer serving the root `.`
	mux.Handle("/", http.FileServer(http.Dir(filepathRoot)))

	// define the http srv
	srv := &http.Server{
		Addr: "localhost:" + port,
		Handler: mux,
	}

	// print report
	log.Printf("Serving on port: %s\n", port)
	// listens on address and calls serve
	// log.Fatal will not close the program (os.Exit(1)) unless the server is shutdown or closed
	log.Fatal(srv.ListenAndServe())
}
