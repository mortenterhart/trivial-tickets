package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

/*
 * Ticketsystem Trivial Tickets
 *
 * Matriculation numbers: 3040018, 3040018, 3478222
 * Lecture:               Programmieren II, INF16B
 * Lecturer:              Herr Prof. Dr. Helmut Neemann
 * Institute:             Duale Hochschule Baden-Württemberg Mosbach
 *
 * ---------------
 *
 * Package main [tests]
 * Main package of the ticketsystem webserver
 */

// TestInitConfigDefault tests the parsing of command line arguments
// and makes sure the default is applied, if there are no flags provided.
func TestInitConfigDefault(t *testing.T) {

	config, err := initConfig()

	assert.NotNil(t, config, "Config struct is nil.")
	assert.Nil(t, err, "err is not nil")
	assert.Equal(t, int16(8443), config.Port, "Config.port is not set to 8443")
	assert.Equal(t, "../../files/tickets", config.Tickets, "Config.tickets is not set to \"files/tickets\"")
	assert.Equal(t, "../../files/users/users.json", config.Users, "Config.users is not set to \"files/users\"")
	assert.Equal(t, "../../www", config.Web, "Config.web is not set to \"../../www\"")
}

// TestIsPortInBoundaries checks if the provided port is within the boundaries of a 16 bit integer
func TestIsPortInBoundaries(t *testing.T) {

	portInBoundaries := 80
	portOutsideBoundaries := 67534

	is80InBoundaries := isPortInBoundaries(portInBoundaries)
	is67534InBoundaries := isPortInBoundaries(portOutsideBoundaries)

	assert.Equal(t, true, is80InBoundaries, "Port 80 is not accepted, but it should be")
	assert.Equal(t, false, is67534InBoundaries, "Port 67534 is accepted. Should not happen.")
}
