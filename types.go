package main

import "os"

type apiConfig struct {
	fileserverHits int
	jwtSecret      string
}

func newAPIConfig() apiConfig {
	return apiConfig{
		fileserverHits: 0,
		jwtSecret:      os.Getenv("JWT_SECRET"),
	}
}

type Chirp struct {
	ID   int    `json:"id"`
	Body string `json:"body"`
}

type User struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
}

type AuthenticatedUser struct {
	ID           int    `json:"id"`
	Email        string `json:"email"`
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}
