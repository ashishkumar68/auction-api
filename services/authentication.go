package services

import "golang.org/x/crypto/bcrypt"

const (
	HashCost = 10
)

func HashPassword(password string) (string, error) {
	passBytes, err := bcrypt.GenerateFromPassword([]byte(password), HashCost)

	return string(passBytes), err
}

func CompareHashAndPass(hash string, pass string) bool {
	return nil == bcrypt.CompareHashAndPassword([]byte(hash), []byte(pass))
}

