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

// Package jsontools provides useful tools for encoding JSON
// from a map.
package jsontools

import (
	"encoding/json"

	"github.com/mortenterhart/trivial-tickets/structs"
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
 * Package jsontools
 * Useful tools for encoding JSON
 */

// MapToJSON converts a given JSONMap with defined properties
// into a valid JSON string with four spaces of indentation.
func MapToJSON(properties structs.JSONMap) []byte {
	// The decoding error of the encoding into JSON is ignored
	// here because the json.MarshalIndent() function only returns
	// an error if an inconvertible type is given. The given json
	// map can be converted by the function and thus it never
	// returns an error.
	jsonString, _ := json.MarshalIndent(properties, "", "    ")
	return jsonString
}
