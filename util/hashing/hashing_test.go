// Trivial Tickets Ticketsystem
// Copyright (C) 2019 The Contributors
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

// Package hashing provides a functionality to hash
// passwords securely.
package hashing

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mortenterhart/trivial-tickets/logger/testlogger"
)

/*
 * Ticketsystem Trivial Tickets
 *
 * Matriculation numbers: 3040018, 6694964, 3478222
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
	testlogger.BeginTest()
	defer testlogger.EndTest()

	isPasswordCorrect := CheckPassword(hash, password)

	assert.True(t, isPasswordCorrect, "Password was not correct")
}

// TestGenerateHash tests that a bcrypt hash is generated
// from a given password without errors.
func TestGenerateHash(t *testing.T) {
	testlogger.BeginTest()
	defer testlogger.EndTest()

	// It is not possible to test for the specific hash, since the salt will always be different
	hash, err := GenerateHash(password)

	assert.NotNil(t, hash, "Hash is nil")
	assert.Nil(t, err, "Hashing the password did not succeed")
}
