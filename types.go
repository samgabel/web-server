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
	Body string `json:"body"`
}
