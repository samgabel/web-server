package database

import (
	"errors"
	"time"

	"github.com/samgabel/web-server/internal/auth"
)

func (db *DB) WriteRefreshToken(userID int, refreshToken string) error {
	dbStruct, err := db.loadDB()
	if err != nil {
		return err
	}
	newRefreshTokenUser := RefreshToken{
		RefreshToken:   refreshToken,
		RefreshExp: time.Now().Add(1440 * time.Hour),
	}
	if dbStruct.RefreshTokens == nil {
		dbStruct.RefreshTokens = make(map[int]RefreshToken)
	}
	dbStruct.RefreshTokens[userID] = newRefreshTokenUser
	if err := db.writeDB(dbStruct); err != nil {
		return err
	}
	return nil
}

func (db *DB) RefreshJWT(refreshToken string, jwtSecret string) (string, error) {
	dbStruct, err := db.loadDB()
	if err != nil {
		return "", err
	}
	for userID, tokenStruct := range dbStruct.RefreshTokens {
		if refreshToken == tokenStruct.RefreshToken && tokenStruct.RefreshExp.After(time.Now()) {
			return auth.NewSignedJWT(userID, jwtSecret, nil)
		}
	}
	return "", errors.New("No valid refresh token found, cannot generate new JWT: Expired or non-existent")
}

func (db *DB) DeleteRefreshToken(refreshToken string) error {
	dbStruct, err := db.loadDB()
	if err != nil {
		return err
	}
	for userID, refreshStruct := range dbStruct.RefreshTokens {
		if refreshToken == refreshStruct.RefreshToken && refreshStruct.RefreshExp.After(time.Now()) {
			dbStruct.RefreshTokens[userID] = RefreshToken{}
			if err := db.writeDB(dbStruct); err != nil {
				return err
			}
			return nil
		}
	}
	return errors.New("No valid refresh token found, cannot revoke")
}
