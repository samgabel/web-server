package main

import (
	"encoding/json"
	"errors"
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

func (cfg *apiConfig) handlerPostChirp(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		requestJWT, ok := strings.CutPrefix(r.Header.Get("Authorization"), "Bearer ")
		if !ok {
			respondWithError(w, http.StatusBadRequest, "Malformed Authorization request header")
			return
		}
		userID, err := auth.VerifySignedJWT(requestJWT, cfg.jwtSecret)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, fmt.Sprintf("Unauthorized attempt to login using JWT: %s", err))
			return
		}
		type parameters struct {
			Body string `json:"body"`
		}
		decoder := json.NewDecoder(r.Body)
		params := parameters{}
		err = decoder.Decode(&params)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
			return
		}
		validated, err := validateChirp(params.Body)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		chirp, err := db.CreateChirp(userID, validated)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error creating Chirp and writing to disk: %s", err))
			return
		}
		respondWithJSON(w, http.StatusCreated, Chirp{
			ID:       chirp.ID,
			Body:     chirp.Body,
			AuthorID: userID,
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

func (cfg *apiConfig) handlerDeleteChirpByID(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		requestJWT, ok := strings.CutPrefix(r.Header.Get("Authorization"), "Bearer ")
		if !ok {
			respondWithError(w, http.StatusBadRequest, "Malformed Authorization request header")
			return
		}
		requestUserID, err := auth.VerifySignedJWT(requestJWT, cfg.jwtSecret)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, fmt.Sprintf("Unauthorized attempt to login using JWT: %s", err))
		}
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
		if requestUserID != targetChirp.AuthorID {
			respondWithError(w, http.StatusForbidden, "The requester ID doesn't match the author ID of the chirp")
			return
		}
		if err := db.DeleteChirp(chirpID); err != nil {
			respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed deleting the Chirp from the database: %s", err))
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

func handlerPostUser(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type parameters struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}
		decoder := json.NewDecoder(r.Body)
		params := parameters{}
		err := decoder.Decode(&params)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
			return
		}
		user, err := db.CreateUser(params.Email, params.Password)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("User could not be created: %s", err))
			return
		}
		respondWithJSON(w, http.StatusCreated, User{
			ID:              user.ID,
			Email:           user.Email,
			ChirpyRedStatus: user.ChirpyRedStatus,
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
		refreshToken, err := auth.GenerateRefreshToken()
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Refresh Token could not be created: %s", err))
			return
		}
		err = db.WriteRefreshToken(user.ID, refreshToken)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Could not update process refresh token: %s", err))
			return
		}
		signedJWT, err := auth.NewSignedJWT(user.ID, cfg.jwtSecret, params.ExpiresInSeconds)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("JWT could not be created: %s", err))
			return
		}
		respondWithJSON(w, http.StatusOK, AuthenticatedUser{
			ID:              user.ID,
			Email:           user.Email,
			RefreshToken:    refreshToken,
			Token:           signedJWT,
			ChirpyRedStatus: user.ChirpyRedStatus,
		})
	}
}

func (cfg *apiConfig) handlerUpdateUser(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		requestJWT, ok := strings.CutPrefix(r.Header.Get("Authorization"), "Bearer ")
		if !ok {
			respondWithError(w, http.StatusBadRequest, "Malformed Authorization request header")
			return
		}
		type parameters struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}
		params := parameters{}
		if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
			respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
			return
		}
		userID, err := auth.VerifySignedJWT(requestJWT, cfg.jwtSecret)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, fmt.Sprintf("Unauthorized attempt to login using JWT: %s", err))
			return
		}
		user, err := db.UpdateUser(userID, params.Email, params.Password)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Unable to update user info: %s", err))
			return
		}
		respondWithJSON(w, http.StatusOK, User{
			ID:              user.ID,
			Email:           user.Email,
			ChirpyRedStatus: user.ChirpyRedStatus,
		})
	}
}

func (cfg *apiConfig) handlerRefresh(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		requestRefreshToken, ok := strings.CutPrefix(r.Header.Get("Authorization"), "Bearer ")
		if !ok {
			respondWithError(w, http.StatusBadRequest, "Malformed Authorization request header")
			return
		}
		newJWT, err := db.RefreshJWT(requestRefreshToken, cfg.jwtSecret)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, fmt.Sprintf("Unable to hand out new JWT: %s", err))
			return
		}
		type responseShape struct {
			Token string `json:"token"`
		}
		respondWithJSON(w, http.StatusOK, responseShape{
			Token: newJWT,
		})
	}
}

func (cfg *apiConfig) handlerRevokeRefresh(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		requestRefreshToken, ok := strings.CutPrefix(r.Header.Get("Authorization"), "Bearer ")
		if !ok {
			respondWithError(w, http.StatusBadRequest, "Malformed Authorization request header")
			return
		}
		err := db.DeleteRefreshToken(requestRefreshToken)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, fmt.Sprintf("Unable to revoke refresh token: %s", err))
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

func handlerChirpyRedConfirmation(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// identify r.Body structure
		type parameters struct {
			Event string `json:"event"`
			Data  struct {
				UserID int `json:"user_id"`
			} `json:"data"`
		}
		// decode the JSON request body into a Go struct
		params := parameters{}
		if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
			respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
			return
		}
		// if the webhook event sends anything other than the "event" field being "user.upgraded" disregard the request
		// by sending a 2XX code we are telling the third-party that we processed the request successfully
		if params.Event != "user.upgraded" {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		// if all else passes then upgrade the user in the database
		if err := db.UpgradeUserToRed(params.Data.UserID); err != nil {
			// respond with 404 http status code if user could not be found in the db
			if errors.Is(err, database.ErrUserNotExist) {
				respondWithError(w, http.StatusNotFound, fmt.Sprintf("Couldn't upgrade user to Red in the database: %s", err))
				return
			}
			// ... otherwise respond with a 500 http status code because something else went wrong
			respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Couldn't upgrade user to Red in the database: %s", err))
			return
		}
		// if there are no errors with the upgrade to database process return 204 http status code
		w.WriteHeader(http.StatusNoContent)
	}
}
