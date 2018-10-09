package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitConfig(t *testing.T) {

	arguments := []string{
		"program.exe",
		"-port",
		"443",
		"-ticketFolder",
		"files/tickets"}

	config := initConfig(arguments)

	assert.NotNil(t, config, "Config struct is nil.")
	assert.Equal(t, 443, config.port, "Config.port is not set to 443")
	assert.Equal(t, "files/tickets", config.ticketFolder, "Config.ticketFolder is not set to \"files/tickets\"")
}
