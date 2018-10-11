package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestInitConfigDefault tests the parsing of command line arguments
// and makes sure the default is applied if there are no flags provided.
func TestInitConfigDefault(t *testing.T) {

	config, err := initConfig()

	assert.NotNil(t, config, "Config struct is nil.")
	assert.Nil(t, err, "err is not nil")
	assert.Equal(t, int16(443), config.port, "Config.port is not set to 443")
	assert.Equal(t, "files/tickets", config.tickets, "Config.tickets is not set to \"files/tickets\"")
	assert.Equal(t, "files/users", config.users, "Config.users is not set to \"files/users\"")
}

// TestIsPortInBoundaries checks if the provided port is within the boundaries of a 16 bit integer
func TestIsPortInBoundaries(t *testing.T) {

	portInBoundaries := 80
	portOutsideBoundaries := 67534

	is80InBoundaries := isPortInBoundaries(&portInBoundaries)
	is67534InBoundaries := isPortInBoundaries(&portOutsideBoundaries)

	assert.Equal(t, true, is80InBoundaries, "Port 80 is not accepted, but it should be")
	assert.Equal(t, false, is67534InBoundaries, "Port 67534 is accepted. Should not happen.")
}
