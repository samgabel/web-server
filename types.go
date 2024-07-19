package main

type apiConfig struct {
	fileserverHits int
}

func newAPIConfig() apiConfig {
	return apiConfig{
		fileserverHits: 0,
	}
}

type Chirp struct {
	ID   int    `json:"id"`
	Body string `json:"body"`
}

type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
}
