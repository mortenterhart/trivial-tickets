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

import "fmt"

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
 * Package hashing [examples]
 * Hashing functions for passwords
 */

// Example CorrectPassword shows that the bcrypt
// hashing algorithm can reliably examine whether
// a given password matches a hash and thus is
// correct or not.
func ExampleCheckPassword_correctPassword() {
	// Set a really strong password
	password := "hello123"

	// Generate the bcrypt hash for the password
	hash, err := GenerateHash(password)
	if err != nil {
		fmt.Println(err)
	}

	// Check that the correct password is recognized
	// as the correct password
	correct := CheckPassword(hash, password)

	fmt.Println(correct)
	// Output: true
}

// Example WrongPassword proves the opposite case
// passing a wrong password to the CheckPassword
// function. It is not recognized as correct
// password and its invalidity is therefore
// verified.
func ExampleCheckPassword_wrongPassword() {
	// Set a really strong password
	password := "hello123"

	// Generate the bcrypt hash for the password
	hash, err := GenerateHash(password)
	if err != nil {
		fmt.Println(err)
	}

	// Check the password belonging to the hash with
	// a wrong password and verify its invalidity
	correct := CheckPassword(hash, "wrong password")

	fmt.Println(correct)
	// Output: false
}
