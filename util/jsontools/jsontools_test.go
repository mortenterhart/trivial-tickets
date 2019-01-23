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
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mortenterhart/trivial-tickets/log/testlog"
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

func TestMapToJSON(t *testing.T) {
	testlog.BeginTest()
	defer testlog.EndTest()

	testProperties := structs.JSONMap{
		"status":  200,
		"message": "OK",
	}

	expectedJSON := `{
    "message": "OK",
    "status": 200
}`

	result := MapToJSON(testProperties)

	t.Run("equalJson", func(t *testing.T) {
		assert.Equal(t, expectedJSON, string(result), "decoded JSON should be equal to the expected result")
	})
}
