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
	"net/http"

	"github.com/mortenterhart/trivial-tickets/logger"
	"github.com/mortenterhart/trivial-tickets/structs"
	"github.com/mortenterhart/trivial-tickets/util/jsontools"
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
 * Package httptools
 * Useful tools for HTTP handlers
 */

// jsonContentType is used as a constant content type for
// json responses
const jsonContentType = "application/json; charset=utf-8"

// contentTypeOptions is the default value for the
// "X-Content-Type-Options" header used in error responses
const contentTypeOptions = "nosniff"

// StatusCodeError writes an error message given by the cause parameter alongside
// with a HTTP status code and its description to the given response writer. The
// status code should be an erroneous one as the function implies it.
func StatusCodeError(writer http.ResponseWriter, cause string, statusCode int) {
	errorMessage := fmt.Sprintf("%d %s: %s", statusCode, http.StatusText(statusCode), cause)
	http.Error(writer, errorMessage, statusCode)
	logger.Error(errorMessage)
}

// JSONResponse writes a json response to the supplied response writer by converting
// the map containing json properties to a valid json string.
func JSONResponse(writer http.ResponseWriter, jsonProperties structs.JSONMap) {
	writer.Header().Set("Content-Type", jsonContentType)
	fmt.Fprintln(writer, responseToJSON(jsonProperties))
}

// JSONError writes a json error response with the properties contained in jsonProperties
// and sets the response status code which should be an erroneous one.
func JSONError(writer http.ResponseWriter, jsonProperties structs.JSONMap, statusCode int) {
	writer.Header().Set("Content-Type", jsonContentType)
	writer.Header().Set("X-Content-Type-Options", contentTypeOptions)
	writer.WriteHeader(statusCode)
	fmt.Fprintln(writer, responseToJSON(jsonProperties))
}

// responseToJSON converts the json properties map to a json string and returns it.
func responseToJSON(jsonProperties structs.JSONMap) string {
	return string(jsontools.MapToJSON(jsonProperties))
}
