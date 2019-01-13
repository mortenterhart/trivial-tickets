package jsontools

import (
	"github.com/stretchr/testify/assert"
	"testing"

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
 * Package jsontools [tests]
 * Useful tools for encoding JSON
 */

func TestMapToJson(t *testing.T) {
	testProperties := structs.JsonMap{
		"status":  200,
		"message": "OK",
	}

	expectedJson := `{
    "message": "OK",
    "status": 200
}`

	result, decodeErr := MapToJson(testProperties)

	t.Run("noDecodeError", func(t *testing.T) {
		assert.NoError(t, decodeErr, "map type should be valid to be decoded to JSON")
	})

	t.Run("equalJson", func(t *testing.T) {
		assert.Equal(t, expectedJson, string(result), "decoded JSON should be equal to the expected result")
	})
}
