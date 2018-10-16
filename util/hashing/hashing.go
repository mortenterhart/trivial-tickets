package hashing

import (
	"golang.org/x/crypto/bcrypt"
)

/*
*
* Matrikelnummern
* 3040018
*
 */

// CheckPassword compares the given password against the stored hash.
// It returns true, if the password is correct.
func CheckPassword(hash string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// GenerateHash returns a bcrypt hash for the given password
// Bcrypt is a very secure hashing algorithm for passwords, which
// also spares the developer of having to generate a salt on his/her own
func GenerateHash(password string) (string, error) {
	hash, error := bcrypt.GenerateFromPassword([]byte(password), 12)
	return string(hash), error
}
