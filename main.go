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

	// register a handler for localhost:8080/app
	mux.Handle("/app/*",
		// strip the "app/" request prefix before passing on to the FileServer handler (essentially makes the "app/" prefix map to the filepathRoot)
		http.StripPrefix("/app",
			// our handler is a FileServer serving the root `./` of this directory
			http.FileServer(http.Dir(filepathRoot)),
		),
	)

	// register a custom handler for readiness endpoint localhost:8080/healthz
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		// write "Content-Type" header
		w.Header().Add("Content-Type", "text/plain; charset=utf-8")
		// write status code
		w.WriteHeader(http.StatusOK)
		// write "OK" to body
		w.Write([]byte(http.StatusText(http.StatusOK)))
	})

	// define the http srv
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
