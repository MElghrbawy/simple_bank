package util

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword hashes the password using bcrypt
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("error hashing password: %w", err)
	}
	return string(bytes), err
}

// CheckPasswordHash checks if the password matches the hash
func CheckPasswordHash(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))

}
