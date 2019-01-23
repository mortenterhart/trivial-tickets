// Useful tools for HTTP handlers
package httptools

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mortenterhart/trivial-tickets/globals"
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
 * Package httptools [tests]
 * Useful tools for HTTP handlers
 */

func TestMain(m *testing.M) {
	initializeLogConfig()

	os.Exit(m.Run())
}

func initializeLogConfig() {
	logConfig := testLogConfig()
	globals.LogConfig = &logConfig
}

func testLogConfig() structs.LogConfig {
	return structs.LogConfig{
		LogLevel:   structs.LevelInfo,
		VerboseLog: false,
		FullPaths:  false,
	}
}

func TestStatusCodeError(t *testing.T) {
	recorder := httptest.NewRecorder()

	StatusCodeError(recorder, "internal error", http.StatusInternalServerError)

	response := recorder.Result()

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
	recorder := httptest.NewRecorder()

	jsonProperties := structs.JsonMap{
		"status":  http.StatusOK,
		"message": "OK",
	}

	JsonResponse(recorder, jsonProperties)

	response := recorder.Result()

	t.Run("statusCode", func(t *testing.T) {
		assert.Equal(t, http.StatusOK, response.StatusCode, "response status code should be 200 OK")
	})

	body, readErr := ioutil.ReadAll(response.Body)

	t.Run("readBody", func(t *testing.T) {
		assert.NoError(t, readErr, "reading response body should not return error")
	})

	t.Run("equalBody", func(t *testing.T) {
		expectedJson, decodeErr := jsontools.MapToJson(jsonProperties)

		assert.NoError(t, decodeErr, "decoding map to JSON should not return error")
		assert.Equal(t, append(expectedJson, '\n'), body, "response body should match decoded expected JSON")
	})
}

func TestJsonError(t *testing.T) {
	recorder := httptest.NewRecorder()

	jsonProperties := structs.JsonMap{
		"status":  http.StatusInternalServerError,
		"message": http.StatusText(http.StatusInternalServerError),
	}

	JsonError(recorder, jsonProperties, http.StatusInternalServerError)

	response := recorder.Result()

	t.Run("statusCode", func(t *testing.T) {
		assert.Equal(t, http.StatusInternalServerError, response.StatusCode, "response status code should be 500 Internal Server Error")
	})

	body, readErr := ioutil.ReadAll(response.Body)

	t.Run("readBody", func(t *testing.T) {
		assert.NoError(t, readErr, "reading response body should not return error")
	})

	t.Run("equalBody", func(t *testing.T) {
		expectedJson, decodeErr := jsontools.MapToJson(jsonProperties)

		assert.NoError(t, decodeErr, "decoding map to JSON should not return error")
		assert.Equal(t, append(expectedJson, '\n'), body, "response body should match decoded expected JSON")
	})
}
