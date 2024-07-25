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
	// check to see if email is already registered in the db
	for _, user := range dbStruct.Users {
		if user.Email == email {
			return User{}, errors.New("Email is already registered, please try another email")
		}
	}
	// hash password
	hash, err := auth.HashPassword(password)
	if err != nil {
		return User{}, err
	}
	newUser := User{
		ID:    newID,
		Email: email,
		// save in database as a byte slice
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
	// load database into memory
	dbStruct, err := db.loadDB()
	if err != nil {
		return User{}, err
	}
	// check if user exists by finding email in db
	var targetUser User
	for _, user := range dbStruct.Users {
		if user.Email == email {
			targetUser = user
		}
	}
	if targetUser.Email == "" {
		return User{}, errors.New("No user associated with the provided email")
	}
	// check if password hash matches the hash stored under the User in the db
	if err := auth.CheckPasswordHash(targetUser.HashedPassword, password); err != nil {
		return User{}, errors.New("Password is incorrect for given email")
	}
	return targetUser, nil
}

func (db *DB) UpdateUser(userID int, email, password string) (User, error) {
	// load db into memory
	dbStruct, err := db.loadDB()
	if err != nil {
		return User{}, err
	}
	// hash the new password
	hashedPassword, err := auth.HashPassword(password)
	if err != nil {
		return User{}, err
	}
	// check for user in database
	if _, ok := dbStruct.Users[userID]; !ok {
		return User{}, errors.New("User does not exist")
	}
	// create a new User struct to replace the old one
	newUser := User{
		ID:             userID,
		Email:          email,
		HashedPassword: hashedPassword,
	}
	// replace the old User info with the new User info
	dbStruct.Users[userID] = newUser
	err = db.writeDB(dbStruct)
	if err != nil {
		return User{}, err
	}
	return newUser, nil
}
