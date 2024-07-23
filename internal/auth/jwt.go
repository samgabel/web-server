package auth

import (
	"errors"
	"strconv"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
)

func NewSignedJWT(userID int, jwtSecret string, durationSeconds *int) (string, error) {
	// handle JWT expiration
	var expiry int
	if durationSeconds == nil || *durationSeconds > 86400 {
		expiry = 86400
	} else {
		expiry = *durationSeconds
	}
	// create a new JWT with additional "claims"
	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "Chirpy",
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expiry) * time.Second)),
		Subject:   strconv.Itoa(userID),
	})
	// sign the new JWT (HMAC, which is compatible with the HS256 signing method that we use, requires that we sign with a []byte)
	signingKey := []byte(jwtSecret)
	return newToken.SignedString(signingKey)
}

func VerifySignedJWT(requestJWT string, jwtSecret string) (int, error) {
	// parse out token from the requestJWT
	requestToken, err := jwt.ParseWithClaims(
		requestJWT,
		&jwt.RegisteredClaims{},
		func(token *jwt.Token) (interface{}, error) { return []byte(jwtSecret), nil },
	)
	if err != nil {
		return 0, err
	}
	// check issuer
	issuer, err := requestToken.Claims.GetIssuer()
	if err != nil {
		return 0, err
	}
	if issuer != "Chirpy" {
		return 0, errors.New("Invalid issuer")
	}
	// grab id from the token
	stringID, err := requestToken.Claims.GetSubject()
	if err != nil {
		return 0, err
	}
	// convert string ID to int ID
	id, err := strconv.Atoi(stringID)
	if err != nil {
		return 0, err
	}
	return id, nil
}
