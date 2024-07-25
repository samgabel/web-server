package database

import (
	"errors"
	"sort"
)

func (db *DB) CreateChirp(authorID int, body string) (Chirp, error) {
	dbStruct, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}
	newID := len(dbStruct.Chirps) + 1
	newChirp := Chirp{
		ID:       newID,
		Body:     body,
		AuthorID: authorID,
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

func (db *DB) DeleteChirp(id int) error {
	dbStruct, err := db.loadDB()
	if err != nil {
		return err
	}
	if _, ok := dbStruct.Chirps[id]; !ok {
		return errors.New("Chirp ID doesn't exist")
	}
	dbStruct.Chirps[id] = Chirp{}
	return db.writeDB(dbStruct)
}
