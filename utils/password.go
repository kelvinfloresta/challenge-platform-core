package utils

import (
	"golang.org/x/crypto/bcrypt"
)

const default_cost = 10

func PasswordGen(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), default_cost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func ComparePassword(password string, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
