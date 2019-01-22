// Main package of the ticketsystem webserver
package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/mortenterhart/trivial-tickets/structs"
	"os"
	"strings"
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

func TestUsageMessage(t *testing.T) {
	var testBuffer bytes.Buffer
	flag.CommandLine.SetOutput(&testBuffer)

	usageMessage()

	t.Run("bufferBeginsWithUsage", func(t *testing.T) {
		assert.True(t, strings.HasPrefix(testBuffer.String(), fmt.Sprintf("Usage: %s [options]", os.Args[0])))
	})

	t.Run("bufferContainsOptions", func(t *testing.T) {
		assert.Contains(t, testBuffer.String(), "options may be one of the following")
	})
}

func TestConvertLogLevel(t *testing.T) {
	t.Run("levelInfo", func(t *testing.T) {
		level, err := convertLogLevel("info")

		assert.NoError(t, err)
		assert.Equal(t, structs.LevelInfo, level)
	})

	t.Run("levelWarning", func(t *testing.T) {
		level, err := convertLogLevel("warning")

		assert.NoError(t, err)
		assert.Equal(t, structs.LevelWarning, level)
	})

	t.Run("levelError", func(t *testing.T) {
		level, err := convertLogLevel("error")

		assert.NoError(t, err)
		assert.Equal(t, structs.LevelError, level)
	})

	t.Run("levelFatal", func(t *testing.T) {
		level, err := convertLogLevel("fatal")

		assert.NoError(t, err)
		assert.Equal(t, structs.LevelFatal, level)
	})

	t.Run("levelInfo", func(t *testing.T) {
		level, err := convertLogLevel("info")

		assert.NoError(t, err)
		assert.Equal(t, structs.LevelInfo, level)
	})

	t.Run("undefinedLevel", func(t *testing.T) {
		level, err := convertLogLevel("undefined")

		assert.Error(t, err)
		assert.Empty(t, level)
	})
}
