package main

import (
	"log"
	"net/http"
)

func main() {
	const port = "8080"

	// define the request multiplexer
	requestMultiplexer := http.NewServeMux()

	// define the http srv
	srv := &http.Server{
		Addr: "localhost:" + port,
		Handler: requestMultiplexer,
	}

	// print report
	log.Printf("Serving on port: %s\n", port)
	// listens on address and calls serve
	// log.Fatal will not close the program (os.Exit(1)) unless the server is shutdown or closed
	log.Fatal(srv.ListenAndServe())
}
