package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"slices"
	"strconv"
	"strings"

	"github.com/samgabel/web-server/internal/database"
)

func respondWithJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(status)
	w.Write(data) //nolint:errcheck
}

func respondWithError(w http.ResponseWriter, status int, msg string) {
	if status > 499 {
		log.Printf("Responding with 5XX error: %s", msg)
	}
	type errorResponse struct {
		Error string `json:"error"`
	}
	respondWithJSON(w, status, errorResponse{
		Error: msg,
	})
}

func validateChirp(body string) (string, error) {
	const maxChirpLength = 140
	if len(body) > maxChirpLength {
		return "", errors.New("Chirp is too long")
	}
	badWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}
	cleanBody := getCleanedBody(badWords, body)
	return cleanBody, nil
}

func getCleanedBody(badWords map[string]struct{}, body string) string {
	words := strings.Split(body, " ")
	for i, word := range words {
		lowerWord := strings.ToLower(word)
		if _, ok := badWords[lowerWord]; ok {
			words[i] = "****"
		}
	}
	return strings.Join(words, " ")
}

func processQueryAuthorID(chirps []database.Chirp, query string) ([]Chirp, error) {
	if query == "" {
		selection := []Chirp{}
		for _, chirp := range chirps {
			selection = append(selection, Chirp{
				ID:       chirp.ID,
				Body:     chirp.Body,
				AuthorID: chirp.AuthorID,
			})
		}
		return selection, nil
	}
	requestedAuthorID, err := strconv.Atoi(query)
	if err != nil {
		return []Chirp{}, errors.New("Improper value given, need int")
	}
	querySelection := []Chirp{}
	for _, chirp := range chirps {
		if chirp.AuthorID == requestedAuthorID {
			querySelection = append(querySelection, Chirp{
				ID:       chirp.ID,
				Body:     chirp.Body,
				AuthorID: chirp.AuthorID,
			})
		}
	}
	if len(querySelection) == 0 {
		return []Chirp{}, errors.New("No chirps found associated with the given value")
	}
	return querySelection, nil
}

func processQuerySort(querySelection []Chirp, querySortType string) ([]Chirp, error) {
	if querySortType != "asc" && querySortType != "desc" && querySortType != "" {
		return []Chirp{}, errors.New("Improper value given, need 'asc' or 'desc'")
	}
	if querySortType == "desc" {
		slices.SortFunc(querySelection, func(a, b Chirp) int {
			if a.ID < b.ID {
				return 1
			}
			if a.ID > b.ID {
				return -1
			}
			return 0
		})
	}
	return querySelection, nil
}
