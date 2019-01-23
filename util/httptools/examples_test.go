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

// Package httptools provides useful tools for building HTTP
// responses.
package httptools

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
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
 * Package httptools [examples]
 * Useful tools for HTTP handlers
 */

// Example BadRequest shows how error messages
// in plain text are built to be used as HTTP
// response of a handler or API.
func ExampleStatusCodeError_badRequest() {
	// Create a new HTTP response recorder to
	// capture the response
	recorder := httptest.NewRecorder()

	// Write the error message as plain text to
	// the recorder
	StatusCodeError(recorder, "invalid JSON on line 4: invalid character '}' after object key", http.StatusBadRequest)

	// Get the HTTP response
	response := recorder.Result()

	// Read the response body
	body, _ := ioutil.ReadAll(response.Body)

	// Print the HTTP response status
	fmt.Println(response.Status)

	// Print the whole response body
	fmt.Println(string(body))
	// Output:
	// 400 Bad Request
	// 400 Bad Request: invalid JSON on line 4: invalid character '}' after object key
}

// Example StatusVerified shows how a JSON response
// for a successful result of a web request is
// constructed.
func ExampleJSONResponse_statusVerified() {
	// Create a new HTTP response recorder to
	// capture the response
	recorder := httptest.NewRecorder()

	// Write the verification information as
	// JSON response to the recorder
	JSONResponse(recorder, map[string]interface{}{
		"verified": true,
		"message":  "Mail was sent successfully",
		"status":   http.StatusOK,
	})

	// Get the HTTP response
	response := recorder.Result()

	// Read the response body
	body, _ := ioutil.ReadAll(response.Body)

	// Print the HTTP response status
	fmt.Println(response.Status)

	// Print the whole response body
	fmt.Println(string(body))
	// Output:
	// 200 OK
	// {
	//     "message": "Mail was sent successfully",
	//     "status": 200,
	//     "verified": true
	// }
}

// Example MethodNotAllowed explains the construction
// of a JSON error response with a custom HTTP response
// status.
func ExampleJSONError_methodNotAllowed() {
	// Create a new HTTP response recorder to
	// capture the response
	recorder := httptest.NewRecorder()

	// Write the JSON error message with a
	// response status to the recorder
	JSONError(recorder, map[string]interface{}{
		"status":  http.StatusMethodNotAllowed,
		"message": "METHOD_NOT_ALLOWED (GET)",
	}, http.StatusMethodNotAllowed)

	// Get the HTTP response
	response := recorder.Result()

	// Read the response body
	body, _ := ioutil.ReadAll(response.Body)

	// Print the HTTP response status
	fmt.Println(response.Status)

	// Print the whole response body
	fmt.Println(string(body))
	// Output:
	// 405 Method Not Allowed
	// {
	//     "message": "METHOD_NOT_ALLOWED (GET)",
	//     "status": 405
	// }
}
