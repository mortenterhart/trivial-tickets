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

// Package random implements a pseudo-random algorithm to
// create unique ticket and mail ids.
package random

import (
	"math/rand"
	"time"
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
 * Package random
 * Random id creation
 */

// letters are the valid characters for the ids
var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

// CreateRandomID generates a pseudo random id for
// tickets and mails.
// Tweaked example from https://stackoverflow.com/a/22892986
func CreateRandomID(n int) string {

	// Seed the random function to make it more random
	rand.Seed(time.Now().UnixNano())

	// Create a slice, big enough to hold the id
	b := make([]rune, n)

	// Randomly choose a letter from the letters rune
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}

	return string(b)
}
