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
	"golang.org/x/crypto/bcrypt"
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
 * Package hashing
 * Hashing functions for passwords
 */

// CheckPassword compares the given password against the stored hash.
// It returns true, if the password is correct.
func CheckPassword(hash string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// GenerateHash returns a bcrypt hash for the given password.
// Bcrypt is a very secure hashing algorithm for passwords, which
// also spares the developer from having to generate a salt on
// his/her own.
func GenerateHash(password string) (string, error) {
	hash, hashErr := bcrypt.GenerateFromPassword([]byte(password), 12)
	return string(hash), hashErr
}
