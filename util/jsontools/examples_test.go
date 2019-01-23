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
 * Package jsontools [examples]
 * Useful tools for encoding JSON
 */

// Example StatusResponse shows how a HTTP status
// response, e.g. from a web service, is converted
// from a map of properties to a valid JSON string.
func ExampleMapToJSON_statusResponse() {
	json := MapToJSON(map[string]interface{}{
		"status":  200,
		"message": "OK",
	})

	fmt.Println(string(json))
	// Output:
	// {
	//     "message": "OK",
	//     "status": 200
	// }
}
