// Useful tools for encoding JSON
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

// MapToJson converts a given JsonMap with defined properties
// into a valid JSON string with four spaces of indentation.
func MapToJson(properties structs.JsonMap) ([]byte, error) {
	return json.MarshalIndent(properties, "", "    ")
}
