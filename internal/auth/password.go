package auth

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

var ErrPasswordTooLong = errors.New("password must not be longer than 72 bytes")

func HashPassword(password string) (string, error) {
	if len([]byte(password)) > 72 {
		return "", ErrPasswordTooLong
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}

func CheckPasswordHash(password string, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
