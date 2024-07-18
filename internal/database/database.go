package database

import (
	"encoding/json"
	"errors"
	"os"
	"sort"
)

// CreateChirp creates a new chirp and saves it to disk
func (db *DB) CreateChirp(body string) (Chirp, error) {
	// load DB into memory
	dbStruct, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}
	// get the new ID
	newID := len(dbStruct.Chirps) + 1
	// create new Chirp
	newChirp := Chirp{
		ID:   newID,
		Body: body,
	}
	// add new Chirp to map with associated ID
	if dbStruct.Chirps == nil {
		dbStruct.Chirps = make(map[int]Chirp)
	}
	dbStruct.Chirps[newID] = newChirp
	// write chirps back to database file
	err = db.writeDB(dbStruct)
	if err != nil {
		return Chirp{}, err
	}
	return newChirp, nil
}

// GetChirps returns all chirps in the database
func (db *DB) GetChirps() ([]Chirp, error) {
	// load our DB into memory
	dbStruct, err := db.loadDB()
	if err != nil {
		return []Chirp{}, err
	}
	// grab all our Chirps and append them to a slice
	chirps := make([]Chirp, 0, len(dbStruct.Chirps))
	for _, chirp := range dbStruct.Chirps {
		chirps = append(chirps, chirp)
	}
	// sort chirps by ID
	sort.Slice(chirps, func(i, j int) bool { return chirps[i].ID < chirps[j].ID })
	return chirps, nil
}

// WipeDB removes all data from the database file, while not deleting the file itself
func (db *DB) WipeDB() error {
	db.mu.Lock()
	defer db.mu.Unlock()
	err := os.WriteFile(db.path, []byte{}, 0644)
	if err != nil {
		return err
	}
	return nil
}

// ensureDB creates a new database file if it doesn't exist
func (db *DB) ensureDB() error {
	// use our mutex to lock our DB
	// get and error from reading the file
	db.mu.RLock()
	_, err := os.ReadFile(db.path)
	db.mu.RUnlock()
	// if our error is an ErrNotExist type then we will create a new database file
	if errors.Is(err, os.ErrNotExist) {
		err := os.WriteFile(db.path, []byte{}, 0644)
		if err != nil {
			return err
		}
	}
	return nil
}

// loadDB reads the database file into memory
func (db *DB) loadDB() (DBStructure, error) {
	// use our mutex to lock our DB
	// and read from our DB
	db.mu.RLock()
	data, err := os.ReadFile(db.path)
	db.mu.RUnlock()
	if err != nil {
		return DBStructure{}, err
	}
	// unmarshal the database json file into a struct
	dbStruct := DBStructure{}
	json.Unmarshal(data, &dbStruct)
	return dbStruct, nil
}

// writeDB writes the database file to disk
func (db *DB) writeDB(dbStructure DBStructure) error {
	// marshal the input DBStructure into JSON to be stored on disk
	data, err := json.Marshal(dbStructure)
	if err != nil {
		return err
	}
	// use our mutex to lock our database
	db.mu.Lock()
	defer db.mu.Unlock()
	// write to file
	err = os.WriteFile(db.path, data, 0644)
	if err != nil {
		return err
	}
	return nil
}
