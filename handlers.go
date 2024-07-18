package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/samgabel/web-server/internal/database"
)

func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`
<html>
<body>
	<h1>Welcome, Chirpy Admin</h1>
	<p>Chirpy has been visited %d times!</p>
</body>
</html>
	`, cfg.fileserverHits)))
}

func (cfg *apiConfig) handlerResetMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Reset fileserverHits to 0"))
	cfg.fileserverHits = 0
}

// We use currying in order to allow us to return a handler without modifying it's function signature,
// instead we pass in a pointer to a database and inject that into our handler
func handlerPostChirp(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type parameters struct {
			Body string `json:"body"`
		}
		decoder := json.NewDecoder(r.Body)
		params := parameters{}
		err := decoder.Decode(&params)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
			return
		}
		// refactor and abstract away validation
		validated, err := validateChirp(params.Body)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		// create chirp and post to DB
		chirp, err := db.CreateChirp(validated)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error creating Chirp and writing to disk: %s", err))
			return
		}
		respondWithJSON(w, http.StatusCreated, chirp)
	}
}

// Create a handler that will grab the contents of the database.json file using the db.GetChirps() method.
// We also want to use currying in order to inject a database pointer into the http.HandlerFunc without changing the function signature.
func handlerGetChirps(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// load all the chirps into memory
		chirps, err := db.GetChirps()
		if err != nil {
			respondWithError(w, http.StatusOK, fmt.Sprintf("Error getting Chirps from database: %s", err))
			return
		}
		// respond with an array of Chirp messages
		respondWithJSON(w, http.StatusOK, chirps)
	}
}

// Not part of the assignment, but I wanted to create a way to quickly wipe the Database with an API call
func handlerDeleteChirps(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// use the WipeDB method in order to load the database.json file with an empty slice of bytes []byte{}
		err := db.WipeDB()
		if err != nil {
			respondWithError(w, http.StatusOK, fmt.Sprintf("Error removing Chirps from database: %s", err))
			return
		}
		// respond with "status": "deleted" message
		respondWithJSON(w, http.StatusOK, struct {
			Status string `json:"status"`
		}{
			Status: "deleted",
		})
	}
}
