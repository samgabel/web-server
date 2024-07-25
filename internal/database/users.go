package database

import (
	"errors"

	"github.com/samgabel/web-server/internal/auth"
)

func (db *DB) CreateUser(email, password string) (User, error) {
	dbStruct, err := db.loadDB()
	if err != nil {
		return User{}, err
	}
	newID := len(dbStruct.Users) + 1
	for _, user := range dbStruct.Users {
		if user.Email == email {
			return User{}, errors.New("Email is already registered, please try another email")
		}
	}
	hash, err := auth.HashPassword(password)
	if err != nil {
		return User{}, err
	}
	newUser := User{
		ID:    newID,
		Email: email,
		HashedPassword: hash,
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

func (db *DB) AuthenticateUser(email, password string) (User, error) {
	dbStruct, err := db.loadDB()
	if err != nil {
		return User{}, err
	}
	var targetUser User
	for _, user := range dbStruct.Users {
		if user.Email == email {
			targetUser = user
		}
	}
	if targetUser.Email == "" {
		return User{}, errors.New("No user associated with the provided email")
	}
	if err := auth.CheckPasswordHash(targetUser.HashedPassword, password); err != nil {
		return User{}, errors.New("Password is incorrect for given email")
	}
	return targetUser, nil
}

func (db *DB) UpdateUser(userID int, email, password string) (User, error) {
	dbStruct, err := db.loadDB()
	if err != nil {
		return User{}, err
	}
	hashedPassword, err := auth.HashPassword(password)
	if err != nil {
		return User{}, err
	}
	if _, ok := dbStruct.Users[userID]; !ok {
		return User{}, errors.New("User does not exist")
	}
	newUser := User{
		ID:             userID,
		Email:          email,
		HashedPassword: hashedPassword,
	}
	dbStruct.Users[userID] = newUser
	err = db.writeDB(dbStruct)
	if err != nil {
		return User{}, err
	}
	return newUser, nil
}
