package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestInitConfigDefault tests the parsing of command line arguments
// and makes sure the default is applied if there are no flags provided.
func TestInitConfigDefault(t *testing.T) {

	config := initConfig()

	assert.NotNil(t, config, "Config struct is nil.")
	assert.Equal(t, 443, config.port, "Config.port is not set to 443")
	assert.Equal(t, "files/tickets", config.tickets, "Config.tickets is not set to \"files/tickets\"")
	assert.Equal(t, "files/users", config.users, "Config.users is not set to \"files/users\"")
}
