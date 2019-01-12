package hashing

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

/*
 * Ticketsystem Trivial Tickets
 *
 * Matriculation numbers: 3040018, 3040018, 3478222
 * Lecture:               Programmieren II, INF16B
 * Lecturer:              Herr Prof. Dr. Helmut Neemann
 * Institute:             Duale Hochschule Baden-Württemberg Mosbach
 *
 * ---------------
 *
 * Package hashing [tests]
 * Hashing functions for passwords
 */

// Password for test and bcrypt hash generated from that password
const password = "MyPassword123!!##"
const hash = "$2a$12$rW6Ska0DaVjTX/8sQGCp/.y7kl2RvF.9936Hmm27HyI0cJ78q1UOG"

// TestCheckPassword checks a given password against a precomputed hash
// and makes sure the hashing function works properly.
func TestCheckPassword(t *testing.T) {

	isPasswordCorrect := CheckPassword(hash, password)

	assert.True(t, isPasswordCorrect, "Password was not correct")
}

// TestGenerateHash tests that a bcrypt hash is generated
// from a given password without errors.
func TestGenerateHash(t *testing.T) {

	// It is not possible to test for the specific hash, since the salt will always be different
	hash, err := GenerateHash(password)

	assert.NotNil(t, hash, "Hash is nil")
	assert.Nil(t, err, "Hashing the password did not succeed")
}
