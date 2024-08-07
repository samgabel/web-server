package database

import (
	"sync"
	"time"
)

type DB struct {
	path string
	mu   *sync.RWMutex
}

func NewDB(path string) (*DB, error) {
	db := &DB{
		path: path,
		mu:   &sync.RWMutex{},
	}
	if err := db.ensureDB(); err != nil {
		return nil, err
	}
	return db, nil
}

type DBStructure struct {
	Chirps        map[int]Chirp        `json:"chirps"`
	Users         map[int]User         `json:"users"`
	RefreshTokens map[int]RefreshToken `json:"refresh_tokens"`
}

type Chirp struct {
	ID       int    `json:"id"`
	Body     string `json:"body"`
	AuthorID int    `json:"author_id"`
}

type User struct {
	ID              int    `json:"id"`
	Email           string `json:"email"`
	HashedPassword  []byte `json:"hashed_password"`
	ChirpyRedStatus bool   `json:"is_chirpy_red"`
}

type RefreshToken struct {
	RefreshToken string    `json:"refresh_token"`
	RefreshExp   time.Time `json:"refresh_expiration"`
}
