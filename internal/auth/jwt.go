package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"strconv"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
)

func NewSignedJWT(userID int, jwtSecret string, durationSeconds *int) (string, error) {
	var expiry int
	if durationSeconds == nil || *durationSeconds > 3600 {
		expiry = 3600
	} else {
		expiry = *durationSeconds
	}
	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "Chirpy",
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expiry) * time.Second)),
		Subject:   strconv.Itoa(userID),
	})
	signingKey := []byte(jwtSecret)
	return newToken.SignedString(signingKey)
}

func VerifySignedJWT(requestJWT string, jwtSecret string) (int, error) {
	requestToken, err := jwt.ParseWithClaims(
		requestJWT,
		&jwt.RegisteredClaims{},
		func(token *jwt.Token) (interface{}, error) { return []byte(jwtSecret), nil },
	)
	if err != nil {
		return 0, err
	}
	issuer, err := requestToken.Claims.GetIssuer()
	if err != nil {
		return 0, err
	}
	if issuer != "Chirpy" {
		return 0, errors.New("Invalid issuer")
	}
	stringID, err := requestToken.Claims.GetSubject()
	if err != nil {
		return 0, err
	}
	id, err := strconv.Atoi(stringID)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func GenerateRefreshToken() (string, error) {
	randBytes := make([]byte, 32)
	_, err := rand.Read(randBytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(randBytes), nil
}
