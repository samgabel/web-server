package database

import (
	"encoding/json"
	"errors"
	"os"
	"sort"
)

func (db *DB) CreateChirp(body string) (Chirp, error) {
	dbStruct, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}
	newID := len(dbStruct.Chirps) + 1
	newChirp := Chirp{
		ID:   newID,
		Body: body,
	}
	if dbStruct.Chirps == nil {
		dbStruct.Chirps = make(map[int]Chirp)
	}
	dbStruct.Chirps[newID] = newChirp
	err = db.writeDB(dbStruct)
	if err != nil {
		return Chirp{}, err
	}
	return newChirp, nil
}

func (db *DB) GetChirps() ([]Chirp, error) {
	dbStruct, err := db.loadDB()
	if err != nil {
		return []Chirp{}, err
	}
	chirps := make([]Chirp, 0, len(dbStruct.Chirps))
	for _, chirp := range dbStruct.Chirps {
		chirps = append(chirps, chirp)
	}
	sort.Slice(chirps, func(i, j int) bool { return chirps[i].ID < chirps[j].ID })
	return chirps, nil
}

func (db *DB) GetChirp(id int) (Chirp, error) {
	dbStruct, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}
	targetChirp, ok := dbStruct.Chirps[id]
	if !ok {
		return Chirp{}, errors.New("Chirp ID doesn't exist")
	}
	return targetChirp, nil
}

func (db *DB) CreateUser(email string) (User, error) {
	dbStruct, err := db.loadDB()
	if err != nil {
		return User{}, err
	}
	newID := len(dbStruct.Users) + 1
	newUser := User{
		ID:    newID,
		Email: email,
	}
	if dbStruct.Users == nil {
		dbStruct.Users = make(map[int]User)
	}
	dbStruct.Users[newID] = newUser
	err = db.writeDB(dbStruct)
	if err != nil {
		return User{}, err
	}
	return newUser, nil
}

func (db *DB) WipeDB() error {
	db.mu.Lock()
	defer db.mu.Unlock()
	err := os.WriteFile(db.path, []byte{}, 0644)
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) ensureDB() error {
	db.mu.RLock()
	_, err := os.ReadFile(db.path)
	db.mu.RUnlock()
	if errors.Is(err, os.ErrNotExist) {
		err := os.WriteFile(db.path, []byte{}, 0644)
		if err != nil {
			return err
		}
	}
	return nil
}

func (db *DB) loadDB() (DBStructure, error) {
	db.mu.RLock()
	data, err := os.ReadFile(db.path)
	db.mu.RUnlock()
	if err != nil {
		return DBStructure{}, err
	}
	dbStruct := DBStructure{}
	json.Unmarshal(data, &dbStruct)
	return dbStruct, nil
}

func (db *DB) writeDB(dbStructure DBStructure) error {
	data, err := json.Marshal(dbStructure)
	if err != nil {
		return err
	}
	db.mu.Lock()
	defer db.mu.Unlock()
	err = os.WriteFile(db.path, data, 0644)
	if err != nil {
		return err
	}
	return nil
}
