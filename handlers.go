package main

import (
	"encoding/json"
	"fmt"
	"net/http"
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

func handlerValidateJSON(w http.ResponseWriter, r *http.Request) {
	// define the json go struct mapping for request
	type parameters struct {
		Body string `json:"body"`
	}
	// define the json go struct mapping for the valid response
	type validResponse struct {
		CleanedBody string `json:"cleaned_body"`
	}
	// create a new decoder with the request body
	decoder := json.NewDecoder(r.Body)
	// create a new instance of the paramters struct for the decoder to dump decoded JSON
	params := parameters{}
	// decode the JSON and return an error
	err := decoder.Decode(&params)
	// return internal server error 500 if there is an error in the decoding
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}
	// return bad request if the chirp msg is longer than 140 character
	const maxChirpLength = 140
	if len(params.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}
	// define case-insensitive bad words map
	badWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}
	// if there are no decoding errors or bad requests, then we will send a reponse `valid: true` JSON msg
	respondWithJSON(w, http.StatusOK, validResponse{
		CleanedBody: cleanBadWords(badWords, params.Body),
	})
}
