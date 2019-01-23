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

// Command ticketsystem starts the Trivial Tickets Ticketsystem
// web server to serve as support ticket platform.
package main

import (
	"bytes"
	"flag"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/mortenterhart/trivial-tickets/globals"
	"github.com/mortenterhart/trivial-tickets/log/testlog"
	"github.com/mortenterhart/trivial-tickets/server"
	"github.com/mortenterhart/trivial-tickets/structs"
	"github.com/mortenterhart/trivial-tickets/structs/defaults"
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

// testConfigs returns the test server and logging
// configuration for the tests.
func testConfigs() (structs.ServerConfig, structs.LogConfig) {
	return structs.ServerConfig{
		Port:    defaults.TestPort,
		Tickets: defaults.TestTickets,
		Users:   defaults.TestUsers,
		Mails:   defaults.TestMails,
		Cert:    defaults.TestCertificate,
		Key:     defaults.TestKey,
		Web:     defaults.TestWeb,
	}, structs.LogConfig{
		LogLevel:  structs.AsLogLevel(defaults.LogLevelString),
		Verbose:   defaults.LogVerbose,
		FullPaths: defaults.LogFullPaths,
	}
}

// productiveServerConfig creates a new server
// configuration used by the productive server.
// Note that this configuration should not be
// used to start a test server.
func productiveServerConfig() structs.ServerConfig {
	return structs.ServerConfig{
		Port:    defaults.ServerPort,
		Tickets: defaults.ServerTickets,
		Users:   defaults.ServerUsers,
		Mails:   defaults.ServerMails,
		Cert:    defaults.ServerCertificate,
		Key:     defaults.ServerKey,
		Web:     defaults.ServerWeb,
	}
}

// resetConfig resets the server and logging configuration
// to its default test values and calls resetFlags() to
// reset the command-line flags too.
func resetConfig() {
	config, logConfig := testConfigs()

	globals.ServerConfig = &config
	globals.LogConfig = &logConfig

	resetFlags(config, logConfig)
}

// resetFlags resets the global command-line flags to its
// default test values given in the applied config and
// logConfig.
func resetFlags(config structs.ServerConfig, logConfig structs.LogConfig) {
	// Reset Server configuration
	*port = uint(config.Port)
	*tickets = config.Tickets
	*users = config.Users
	*mails = config.Mails
	*cert = config.Cert
	*key = config.Key
	*web = config.Web

	// Reset Logging configuration
	*verbose = logConfig.Verbose
	*fullPaths = logConfig.FullPaths
	*logLevelString = defaults.LogLevelString
}

// TestInitConfigDefault tests the parsing of command line arguments
// and makes sure the default is applied, if there are no flags provided.
func TestInitConfigDefault(t *testing.T) {
	testlog.BeginTest()
	defer testlog.EndTest()

	defer resetConfig()

	_, logConfig := testConfigs()
	serverConfig := productiveServerConfig()

	config, err := initConfig()

	assert.NotNil(t, config, "ServerConfig struct is nil.")
	assert.Nil(t, err, "err is not nil")
	assert.Equalf(t, serverConfig.Port, config.Port, "ServerConfig.Port is not set to %d", serverConfig.Port)
	assert.Equalf(t, serverConfig.Tickets, config.Tickets, "ServerConfig.Tickets is not set to \"%s\"", serverConfig.Tickets)
	assert.Equalf(t, serverConfig.Users, config.Users, "ServerConfig.Users is not set to \"%s\"", serverConfig.Users)
	assert.Equalf(t, serverConfig.Mails, config.Mails, "ServerConfig.Mails is not set to \"%s\"", serverConfig.Mails)
	assert.Equalf(t, serverConfig.Cert, config.Cert, "ServerConfig.Cert is not set to \"%s\"", serverConfig.Cert)
	assert.Equalf(t, serverConfig.Key, config.Key, "ServerConfig.Key is not set to \"%s\"", serverConfig.Key)
	assert.Equalf(t, serverConfig.Web, config.Web, "ServerConfig.Web is not set to \"%s\"", serverConfig.Web)

	assert.NotNil(t, globals.LogConfig, "globals.LogConfig is nil")
	assert.Equal(t, logConfig.Verbose, globals.LogConfig.Verbose, "LogConfig.Verbose is not set to false")
	assert.Equal(t, logConfig.FullPaths, globals.LogConfig.FullPaths, "LogConfig.FullPaths is not set to false")
	assert.Equal(t, logConfig.LogLevel, globals.LogConfig.LogLevel, "LogConfig.LogLevel is not set to LevelInfo")
}

// TestInitConfigInvalidPort tests the check for an invalid port
// and expects an error
func TestInitConfigInvalidPort(t *testing.T) {
	testlog.BeginTest()
	defer testlog.EndTest()

	defer resetConfig()

	*port = math.MaxUint16 + 1

	config, err := initConfig()

	assert.Error(t, err, "specified port is out of valid bounds")
	assert.Empty(t, config, "config should be empty, so all values should be default values")
}

// TestInitConfigInvalidLogLevelString checks if an invalid log level
// passed as an command line argument invokes an error
func TestInitConfigInvalidLogLevelString(t *testing.T) {
	testlog.BeginTest()
	defer testlog.EndTest()

	defer resetConfig()

	*logLevelString = "invalid"

	config, err := initConfig()

	assert.Error(t, err, "invalid log level string should produce an error")
	assert.Empty(t, config, "config should be empty, so all values should be default values")
}

// TestIsPortInBoundaries checks if the provided port is within the boundaries of a 16 bit unsigned integer
func TestIsPortInBoundaries(t *testing.T) {
	testlog.BeginTest()
	defer testlog.EndTest()

	defer resetConfig()

	portInBoundaries := 80
	portOutsideBoundaries := 67534

	is80InBoundaries := isPortInBoundaries(uint(portInBoundaries))
	is67534InBoundaries := isPortInBoundaries(uint(portOutsideBoundaries))

	assert.Equal(t, true, is80InBoundaries, "Port 80 is not accepted, but it should be")
	assert.Equal(t, false, is67534InBoundaries, "Port 67534 is accepted. Should not happen.")
}

// TestUsageMessage checks that the usage message and all options
// are written to stderr or to the provided buffer
func TestUsageMessage(t *testing.T) {
	testlog.BeginTest()
	defer testlog.EndTest()

	var testBuffer bytes.Buffer
	flag.CommandLine.SetOutput(&testBuffer)

	usageMessage()

	t.Run("bufferBeginsWithUsage", func(t *testing.T) {
		assert.True(t, strings.HasPrefix(testBuffer.String(), fmt.Sprintf("Usage: %s [options]", filepath.Base(os.Args[0]))))
	})

	t.Run("bufferContainsPortOption", func(t *testing.T) {
		assert.Contains(t, testBuffer.String(), "-port")
	})
}

// TestConvertLogLevel checks that all provided strings for log levels
// in the command line arguments are mapped to the correct log level
func TestConvertLogLevel(t *testing.T) {
	testlog.BeginTest()
	defer testlog.EndTest()

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

	t.Run("undefinedLevel", func(t *testing.T) {
		level, err := convertLogLevel("undefined")

		assert.Error(t, err)
		assert.Equal(t, structs.LogLevel(-1), level)
	})
}

func TestMainFunctionStartServer(t *testing.T) {
	testlog.BeginTest()
	defer testlog.EndTest()

	defer resetConfig()

	*port = uint(defaults.TestPort) + 2

	exit = func(code int) {
		t.Run("exitCode", func(t *testing.T) {
			assert.Equal(t, int(defaults.ExitSuccessful), code, "exit code of main() should be 0")
		})
	}

	go func() {
		main()
	}()

	time.Sleep(500 * time.Millisecond)
	server.ShutdownServer()
}

func TestMainFunctionConfigError(t *testing.T) {
	testlog.BeginTest()
	defer testlog.EndTest()

	defer resetConfig()

	*logLevelString = "undefined"

	fatal = func(v ...interface{}) {
		t.Run("configError", func(t *testing.T) {
			assert.Equal(t, fmt.Sprint(v...), fmt.Sprintf("log level '%s' not defined", *logLevelString),
				"the invalid log level should cause an error")
		})
	}

	done := make(chan bool)

	go func() {
		main()
		done <- true
		close(done)
	}()

	<-done
}

func TestMainFunctionServerError(t *testing.T) {
	testlog.BeginTest()
	defer testlog.EndTest()

	defer resetConfig()

	*port = uint(defaults.TestPort) + 2
	*cert = "not/existing/server.cert"

	fatal = func(v ...interface{}) {
		t.Run("errorNotExistingCertificate", func(t *testing.T) {
			assert.Contains(t, fmt.Sprintln(v...), "error while starting server")
		})
	}

	done := make(chan bool)

	go func() {
		main()
		done <- true
		close(done)
	}()

	<-done
}
