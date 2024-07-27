package main

import "os"

type apiConfig struct {
	fileserverHits int
	jwtSecret      string
	polkaKey       string
}

func newAPIConfig() apiConfig {
	return apiConfig{
		fileserverHits: 0,
		jwtSecret:      os.Getenv("JWT_SECRET"),
		polkaKey:       os.Getenv("POLKA_KEY"),
	}
}

type Chirp struct {
	ID       int    `json:"id"`
	Body     string `json:"body"`
	AuthorID int    `json:"author_id"`
}

type User struct {
	ID              int    `json:"id"`
	Email           string `json:"email"`
	ChirpyRedStatus bool   `json:"is_chirpy_red"`
}

type AuthenticatedUser struct {
	ID              int    `json:"id"`
	Email           string `json:"email"`
	Token           string `json:"token"`
	RefreshToken    string `json:"refresh_token"`
	ChirpyRedStatus bool   `json:"is_chirpy_red"`
}
