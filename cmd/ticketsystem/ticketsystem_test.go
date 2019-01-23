// Main package of the ticketsystem webserver
package main

import (
	"bytes"
	"flag"
	"fmt"
	"math"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mortenterhart/trivial-tickets/globals"
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
 * Package main [tests]
 * Main package of the ticketsystem webserver
 */

func defaultConfig() (structs.Config, structs.LogConfig) {
	return structs.Config{
		Port:    8443,
		Tickets: "../../files/tickets",
		Users:   "../../files/users/users.json",
		Mails:   "../../files/mails",
		Cert:    "../../ssl/server.cert",
		Key:     "../../ssl/server.key",
		Web:     "../../www",
	}, structs.LogConfig{
		LogLevel:   structs.LevelInfo,
		VerboseLog: false,
		FullPaths:  false,
	}
}

func resetConfig() {
	config, logConfig := defaultConfig()

	globals.ServerConfig = &config
	globals.LogConfig = &logConfig

	resetFlags(config, logConfig)
}

func resetFlags(config structs.Config, logConfig structs.LogConfig) {
	// Reset Server configuration
	*port = int(config.Port)
	*tickets = config.Tickets
	*users = config.Users
	*mails = config.Mails
	*cert = config.Cert
	*key = config.Key
	*web = config.Web

	// Reset Logging configuration
	*verbose = logConfig.VerboseLog
	*fullPaths = logConfig.FullPaths
	*logLevelString = "info"
}

// TestInitConfigDefault tests the parsing of command line arguments
// and makes sure the default is applied, if there are no flags provided.
func TestInitConfigDefault(t *testing.T) {
	defer resetConfig()

	config, err := initConfig()

	assert.NotNil(t, config, "Config struct is nil.")
	assert.Nil(t, err, "err is not nil")
	assert.Equal(t, uint16(8443), config.Port, "Config.port is not set to 8443")
	assert.Equal(t, "../../files/tickets", config.Tickets, "Config.tickets is not set to \"files/tickets\"")
	assert.Equal(t, "../../files/users/users.json", config.Users, "Config.users is not set to \"files/users\"")
	assert.Equal(t, "../../www", config.Web, "Config.web is not set to \"../../www\"")

	assert.NotNil(t, globals.LogConfig, "globals.LogConfig is nil")
	assert.Equal(t, false, globals.LogConfig.VerboseLog, "LogConfig.VerboseLog is not set to false")
	assert.Equal(t, false, globals.LogConfig.FullPaths, "LogConfig.FullPaths is not set to false")
	assert.Equal(t, structs.LevelInfo, globals.LogConfig.LogLevel, "LogConfig.LogLevel is not set to LevelInfo")
}

// TestInitConfigInvalidPort tests the check for an invalid port
// and expects an error
func TestInitConfigInvalidPort(t *testing.T) {
	defer resetConfig()

	*port = math.MaxUint16 + 1

	config, err := initConfig()

	assert.Error(t, err, "specified port is out of valid bounds")
	assert.Empty(t, config, "config should be empty, so all values should be default values")
}

// TestInitConfigInvalidLogLevelString checks if an invalid log level
// passed as an command line argument invokes an error
func TestInitConfigInvalidLogLevelString(t *testing.T) {
	defer resetConfig()

	*logLevelString = "invalid"

	config, err := initConfig()

	assert.Error(t, err, "invalid log level string should produce an error")
	assert.Empty(t, config, "config should be empty, so all values should be default values")
}

// TestIsPortInBoundaries checks if the provided port is within the boundaries of a 16 bit unsigned integer
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
