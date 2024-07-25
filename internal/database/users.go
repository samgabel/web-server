package database

import (
	"errors"
	"time"

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
	user, ok := dbStruct.Users[userID]
	if !ok {
		return User{}, errors.New("User does not exist")
	}
	// create a new User struct to replace the old one
	newUser := User{
		ID:             userID,
		Email:          email,
		HashedPassword: hashedPassword,
		// making sure to keep refresh token info in place
		RefreshToken:   user.RefreshToken,
		RefreshExp:     user.RefreshExp,
	}
	// replace the old User info with the new User info
	dbStruct.Users[userID] = newUser
	err = db.writeDB(dbStruct)
	if err != nil {
		return User{}, err
	}
	return newUser, nil
}

func (db *DB) WriteRefreshToken(userID int, refreshToken string) error {
	// load db into memory
	dbStruct, err := db.loadDB()
	if err != nil {
		return err
	}
	// Find User
	user, ok := dbStruct.Users[userID]
	if !ok {
		return errors.New("Could not find user")
	}
	// create a new User struct to replace the old one
	newRefreshTokenUser := User{
		ID:             userID,
		Email:          user.Email,
		HashedPassword: user.HashedPassword,
		RefreshToken:   refreshToken,
		// 60 day expiration from creation
		RefreshExp: time.Now().Add(1440 * time.Hour),
	}
	// write the new refreshToken to the db
	dbStruct.Users[userID] = newRefreshTokenUser
	err = db.writeDB(dbStruct)
	if err != nil {
		return err
	}
	// return our new Refresh Token
	return nil
}

func (db *DB) RefreshJWT(refreshToken string, jwtSecret string) (string, error) {
	// load db into memory
	dbStruct, err := db.loadDB()
	if err != nil {
		return "", err
	}
	// range over users in db
	for userID, user := range dbStruct.Users {
		// if our user has a refresh token that matches our given token and it still hasn't expired, we return a new JWT
		if refreshToken == user.RefreshToken && user.RefreshExp.After(time.Now()) {
			return auth.NewSignedJWT(userID, jwtSecret, nil)
		}
	}
	// return an error if no valid token is found
	return "", errors.New("No valid refresh token found, cannot generate new JWT: Expired or non-existent")
}

func (db *DB) DeleteRefreshToken(refreshToken string) error {
	// load db into memory
	dbStruct, err := db.loadDB()
	if err != nil {
		return err
	}
	// range over users in db
	for userID, user := range dbStruct.Users {
		// if our user has a refresh token that matches our given token and it still hasn't expired, we remove it
		if refreshToken == user.RefreshToken && user.RefreshExp.After(time.Now()) {
			// remove the refresh token and expiration date
			dbStruct.Users[userID] = User{
				ID:             user.ID,
				Email:          user.Email,
				HashedPassword: user.HashedPassword,
			}
			// write the edited dbStruct to the db
			if err := db.writeDB(dbStruct); err != nil {
				return err
			}
			return nil
		}
	}
	// return an error if no valid token is found
	return errors.New("No valid refresh token found, cannot revoke")
}
