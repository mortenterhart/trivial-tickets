package hashing

import (
	"golang.org/x/crypto/bcrypt"
)

// CheckPassword compares the given password against the stored hash.
// It returns true, if the password is correct.
func CheckPassword(hash string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// GenerateHash returns a bcrypt hash for the given password
func GenerateHash(password string) (string, error) {
	hash, error := bcrypt.GenerateFromPassword([]byte(password), 12)
	return string(hash), error
}
