package auth

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) ([]byte, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return []byte{}, err
	}
	return hash, nil
}

func CheckPasswordHash(hash []byte, password string) error{
	return bcrypt.CompareHashAndPassword(hash, []byte(password))
}
