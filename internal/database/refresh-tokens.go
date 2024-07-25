package database

import (
	"errors"
	"time"

	"github.com/samgabel/web-server/internal/auth"
)

func (db *DB) WriteRefreshToken(userID int, refreshToken string) error {
	// load db into memory
	dbStruct, err := db.loadDB()
	if err != nil {
		return err
	}
	// create a new RefreshToken struct
	newRefreshTokenUser := RefreshToken{
		RefreshToken:   refreshToken,
		// 60 day expiration from creation
		RefreshExp: time.Now().Add(1440 * time.Hour),
	}
	// initialize map if non has been created yet
	if dbStruct.RefreshTokens == nil {
		dbStruct.RefreshTokens = make(map[int]RefreshToken)
	}
	// write the new refreshToken to the db (the map key corresponds to the userID associated with the Refresh Token)
	dbStruct.RefreshTokens[userID] = newRefreshTokenUser
	if err := db.writeDB(dbStruct); err != nil {
		return err
	}
	// return nil if successful
	return nil
}

func (db *DB) RefreshJWT(refreshToken string, jwtSecret string) (string, error) {
	// load db into memory
	dbStruct, err := db.loadDB()
	if err != nil {
		return "", err
	}
	// range over users in db
	for userID, tokenStruct := range dbStruct.RefreshTokens {
		// if our user has a refresh token that matches our given token and it still hasn't expired, we return a new JWT
		if refreshToken == tokenStruct.RefreshToken && tokenStruct.RefreshExp.After(time.Now()) {
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
	for userID, refreshStruct := range dbStruct.RefreshTokens {
		// if our user has a refresh token that matches our given token and it still hasn't expired, we remove it
		if refreshToken == refreshStruct.RefreshToken && refreshStruct.RefreshExp.After(time.Now()) {
			// remove the refresh token and expiration date and write to database
			dbStruct.RefreshTokens[userID] = RefreshToken{}
			if err := db.writeDB(dbStruct); err != nil {
				return err
			}
			// return nil error if successful
			return nil
		}
	}
	// return an error if no valid token is found
	return errors.New("No valid refresh token found, cannot revoke")
}
