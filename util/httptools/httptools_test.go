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
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mortenterhart/trivial-tickets/globals"
	"github.com/mortenterhart/trivial-tickets/logger/testlogger"
	"github.com/mortenterhart/trivial-tickets/structs"
	"github.com/mortenterhart/trivial-tickets/structs/defaults"
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
 * Package httptools [tests]
 * Useful tools for HTTP handlers
 */

func TestMain(m *testing.M) {
	initializeLogConfig()

	os.Exit(m.Run())
}

// initializeLogConfig initializes the global logging
// configuration with test values.
func initializeLogConfig() {
	logConfig := testLogConfig()
	globals.LogConfig = &logConfig
}

// testLogConfig returns a logging configuration suitable
// to be used in tests.
func testLogConfig() structs.LogConfig {
	return structs.LogConfig{
		LogLevel:  structs.AsLogLevel(defaults.LogLevelString),
		Verbose:   defaults.LogVerbose,
		FullPaths: defaults.LogFullPaths,
	}
}

func TestStatusCodeError(t *testing.T) {
	testlogger.BeginTest()
	defer testlogger.EndTest()

	recorder := httptest.NewRecorder()

	StatusCodeError(recorder, "internal error", http.StatusInternalServerError)

	response := recorder.Result()
	defer response.Body.Close()

	t.Run("statusCode", func(t *testing.T) {
		assert.Equal(t, http.StatusInternalServerError, response.StatusCode, "response status code should be 500 Internal Server Error")
	})

	body, readErr := ioutil.ReadAll(response.Body)

	t.Run("readBody", func(t *testing.T) {
		assert.NoError(t, readErr, "reading response body should not return error")
	})

	t.Run("equalErrorMessage", func(t *testing.T) {
		expectedErrorMessage := "500 Internal Server Error: internal error\n"

		assert.Equal(t, expectedErrorMessage, string(body), "error message written to recorder should be equal to expected error message")
	})
}

func TestJsonResponse(t *testing.T) {
	testlogger.BeginTest()
	defer testlogger.EndTest()

	recorder := httptest.NewRecorder()

	jsonProperties := structs.JSONMap{
		"status":  http.StatusOK,
		"message": "OK",
	}

	JSONResponse(recorder, jsonProperties)

	response := recorder.Result()
	defer response.Body.Close()

	t.Run("statusCode", func(t *testing.T) {
		assert.Equal(t, http.StatusOK, response.StatusCode, "response status code should be 200 OK")
	})

	t.Run("jsonContentType", func(t *testing.T) {
		assert.Equal(t, jsonContentType, response.Header.Get("Content-Type"), "response should have the JSON content type")
	})

	body, readErr := ioutil.ReadAll(response.Body)

	t.Run("readBody", func(t *testing.T) {
		assert.NoError(t, readErr, "reading response body should not return error")
	})

	t.Run("equalBody", func(t *testing.T) {
		expectedJSON := jsontools.MapToJSON(jsonProperties)

		assert.Equal(t, append(expectedJSON, '\n'), body, "response body should match decoded expected JSON")
	})
}

func TestJsonError(t *testing.T) {
	testlogger.BeginTest()
	defer testlogger.EndTest()

	recorder := httptest.NewRecorder()

	jsonProperties := structs.JSONMap{
		"status":  http.StatusInternalServerError,
		"message": http.StatusText(http.StatusInternalServerError),
	}

	JSONError(recorder, jsonProperties, http.StatusInternalServerError)

	response := recorder.Result()
	defer response.Body.Close()

	t.Run("statusCode", func(t *testing.T) {
		assert.Equal(t, http.StatusInternalServerError, response.StatusCode, "response status code should be 500 Internal Server Error")
	})

	t.Run("jsonContentType", func(t *testing.T) {
		assert.Equal(t, jsonContentType, response.Header.Get("Content-Type"), "response should have JSON content type")
	})

	t.Run("noSniffOption", func(t *testing.T) {
		assert.Equal(t, contentTypeOptions, response.Header.Get("X-Content-Type-Options"), "response should have the 'nosniff' option")
	})

	body, readErr := ioutil.ReadAll(response.Body)

	t.Run("readBody", func(t *testing.T) {
		assert.NoError(t, readErr, "reading response body should not return error")
	})

	t.Run("equalBody", func(t *testing.T) {
		expectedJSON := jsontools.MapToJSON(jsonProperties)

		assert.Equal(t, append(expectedJSON, '\n'), body, "response body should match decoded expected JSON")
	})
}
