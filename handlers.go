package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/samgabel/web-server/internal/auth"
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
			// add password parameter in the API request
			Password string `json:"password"`
		}
		decoder := json.NewDecoder(r.Body)
		params := parameters{}
		err := decoder.Decode(&params)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
			return
		}
		// CreateUser will now hash the API request password parameter and store it int the database
		user, err := db.CreateUser(params.Email, params.Password)
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

func (cfg *apiConfig) handlerLogin(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type parameters struct {
			Email            string `json:"email"`
			Password         string `json:"password"`
			ExpiresInSeconds *int   `json:"expires_in_seconds"`
		}
		decoder := json.NewDecoder(r.Body)
		params := parameters{}
		err := decoder.Decode(&params)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
			return
		}
		user, err := db.AuthenticateUser(params.Email, params.Password)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, fmt.Sprintf("User could not be authenticated: %s", err))
			return
		}
		// create our signed JWT string
		signedJWT, err := auth.NewSignedJWT(user.ID, cfg.jwtSecret, params.ExpiresInSeconds)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("JWT could not be created: %s", err))
			return
		}
		respondWithJSON(w, http.StatusOK, AuthenticatedUser{
			ID:    user.ID,
			Email: user.Email,
			// add our JWT to our response body
			Token: signedJWT,
		})
	}
}

func (cfg *apiConfig) handlerUpdateUser(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// get authorization request header (the JWT sent in the request)
		requestJWT, ok := strings.CutPrefix(r.Header.Get("Authorization"), "Bearer ")
		if !ok {
			respondWithError(w, http.StatusBadRequest, "Malformed Authorization request header")
			return
		}
		// request body parameters
		type parameters struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}
		// decode request body
		params := parameters{}
		if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
			respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
			return
		}
		// make sure the request JWT matches and is signed
		userID, err := auth.VerifySignedJWT(requestJWT, cfg.jwtSecret)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, fmt.Sprintf("Unauthorized attempt to login using JWT: %s", err))
			return
		}
		// update the user in the database once the JWT validation has completed
		user, err := db.UpdateUser(userID, params.Email, params.Password)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Unable to update user info: %s", err))
			return
		}
		// respond with the updated user info
		respondWithJSON(w, http.StatusOK, User{
			ID:    user.ID,
			Email: user.Email,
		})
	}
}
