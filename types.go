package main

type apiConfig struct {
	fileserverHits int
}

func newAPIConfig() apiConfig {
	return apiConfig{
		fileserverHits: 0,
	}
}
