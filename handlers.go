package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

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
		validated, err := validateChirp(params.Body)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		chirp, err := db.CreateChirp(validated)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error creating Chirp and writing to disk: %s", err))
			return
		}
		respondWithJSON(w, http.StatusCreated, Chirp{
			ID:   chirp.ID,
			Body: chirp.Body,
		})
	}
}

func handlerGetChirps(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		chirps, err := db.GetChirps()
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error getting Chirps from database: %s", err))
			return
		}
		respondWithJSON(w, http.StatusOK, chirps)
	}
}

func handlerGetChirpByID(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		chirpIDString := r.PathValue("chirpID")
		chirpID, err := strconv.Atoi(chirpIDString)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid Chirp ID")
			return
		}
		targetChirp, err := db.GetChirp(chirpID)
		if err != nil {
			respondWithError(w, http.StatusNotFound, fmt.Sprintf("Chirp not found in database: %s", err))
			return
		}
		respondWithJSON(w, http.StatusOK, Chirp{
			ID:   targetChirp.ID,
			Body: targetChirp.Body,
		})
	}
}

func handlerPostUser(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type parameters struct {
			Email string `json:"email"`
		}
		decoder := json.NewDecoder(r.Body)
		params := parameters{}
		err := decoder.Decode(&params)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
			return
		}
		user, err := db.CreateUser(params.Email)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("User could not be created: %s", err))
			return
		}
		respondWithJSON(w, http.StatusCreated, User{
			ID:    user.ID,
			Email: user.Email,
		})
	}
}
