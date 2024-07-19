package database

import (
	"encoding/json"
	"errors"
	"os"
)

func (db *DB) WipeDB() error {
	err := os.Remove(db.path)
	if err != nil {
		return err
	}
	return db.ensureDB()
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
