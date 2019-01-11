package httptools

import (
	"github.com/mortenterhart/trivial-tickets/structs"
	"github.com/mortenterhart/trivial-tickets/util/jsontools"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

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
