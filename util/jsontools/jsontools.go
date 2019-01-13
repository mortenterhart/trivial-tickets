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

func MapToJson(properties structs.JsonMap) ([]byte, error) {
	return json.MarshalIndent(properties, "", "    ")
}
