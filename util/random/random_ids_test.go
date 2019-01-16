// Random id creation
package random

import (
	"testing"

	"github.com/stretchr/testify/assert"
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
 * Package random [tests]
 * Random id creation
 */

// TestCreateRandomId makes sure the created ticket id is in line with the specification
func TestCreateRandomId(t *testing.T) {

	ticketId := CreateRandomId(10)

	assert.True(t, len(ticketId) == 10, "Random id has the wrong length")
}
