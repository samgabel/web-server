package database

import (
	"sync"
)

type DB struct {
	path string
	mu   *sync.RWMutex
}

// NewDB creates a new database connection
// and creates the database file if it doesn't exist
func NewDB(path string) (*DB, error) {
	// create new DB struct instance
	db := &DB{
		path: path,
		mu:   &sync.RWMutex{},
	}
	// create a database file if one doesn't exist
	if err := db.ensureDB(); err != nil {
		return nil, err
	}
	return db, nil
}

type DBStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
	Users  map[int]User  `json:"users"`
}

type Chirp struct {
	ID   int    `json:"id"`
	Body string `json:"body"`
}

type User struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
}
