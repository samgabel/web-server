package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

func respondWithJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	// marshal the payload struct into JSON
	data, err := json.Marshal(payload)
	// if there is an error in the marshalling process we want to respond with an internal server error
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// if marshalling succeeds write the status and the marshalled payload
	w.WriteHeader(status)
	w.Write(data)
}

func respondWithError(w http.ResponseWriter, status int, msg string) {
	// log the internal server error message
	if status > 499 {
		log.Printf("Responding with 5XX error: %s", msg)
	}
	// define the json go struct mapping for the error response
	type errorResponse struct {
		Error string `json:"error"`
	}
	// create JSON msg to be written of the msg string error
	respondWithJSON(w, status, errorResponse{
		Error: msg,
	})
}

func cleanBadWords(badWords map[string]struct{}, body string) string {
	// split the string into a slice of strings by " " seperator
	words := strings.Split(body, " ")
	for i, word := range words {
		lowerWord := strings.ToLower(word)
		// check if word is in the badWords map
		if _, ok := badWords[lowerWord]; ok {
			words[i] = "****"
		}
	}
	// return the joined slice of cleaned words
	return strings.Join(words, " ")
}
