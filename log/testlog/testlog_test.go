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

// Package testlog defines another logger exclusively used
// by tests.
package testlog

import (
	"bytes"
	"fmt"
	"os"
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
 * Package testlog [tests]
 * Logging interface for tests
 */

// resetTestLogger resets the logging destination
// of the test logger to `os.Stdout`.
func resetTestLogger() {
	updateLogger(os.Stdout)
}

const (
	// testBufferNotNil denominates the subtest for checking
	// that the test log buffer is not nil.
	testBufferNotNil string = "bufferNotNil"

	// testBufferNotEmpty denominates the subtest of checking
	// that the test log buffer is not empty.
	testBufferNotEmpty string = "bufferNotEmpty"

	// testBufferContainsLogLevel denominates the subtest of
	// checking that the test log buffer contains the mentioned
	// log level.
	testBufferContainsLogLevel string = "bufferContainsLogLevel"

	// testBufferContainsMessage denominates the subtest of
	// checking that the test log buffer contains the logged
	// message.
	testBufferContainsMessage string = "bufferContainsMessage"
)

func TestDebug(t *testing.T) {
	BeginTest()
	defer EndTest()

	defer resetTestLogger()

	var testLogger bytes.Buffer
	updateLogger(&testLogger)

	message := "Calling the Debug function of package testlogger"

	Debug(message)

	t.Run(testBufferNotNil, func(t *testing.T) {
		assert.NotNil(t, testLogger, "test buffer should not be nil")
	})

	t.Run(testBufferNotEmpty, func(t *testing.T) {
		assert.True(t, testLogger.Len() > 0, "test buffer should not be empty")
	})

	t.Run(testBufferContainsLogLevel, func(t *testing.T) {
		assert.Contains(t, testLogger.String(), testDebugPrefix, "test buffer should contain TEST DEBUG log level")
	})

	t.Run(testBufferContainsMessage, func(t *testing.T) {
		assert.Contains(t, testLogger.String(), message, "test buffer should contain the message")
	})
}

func TestDebugf(t *testing.T) {
	BeginTest()
	defer EndTest()

	defer resetTestLogger()

	var testLogger bytes.Buffer
	updateLogger(&testLogger)

	const contentLength int = 30
	format := "Received API request to %s (Method = %s, Host = %s, Content-Type = %s, Content-Length = %d)"
	arguments := []interface{}{"/api/receive", "POST", "localhost:8443", "application/json", contentLength}
	expectedMessage := fmt.Sprintf(format, arguments...)

	Debugf(format, arguments...)

	t.Run(testBufferNotNil, func(t *testing.T) {
		assert.NotNil(t, testLogger, "test buffer should not be nil")
	})

	t.Run(testBufferNotEmpty, func(t *testing.T) {
		assert.True(t, testLogger.Len() > 0, "test buffer should not be empty")
	})

	t.Run(testBufferContainsLogLevel, func(t *testing.T) {
		assert.Contains(t, testLogger.String(), testDebugPrefix, "test buffer should contain TEST DEBUG log level")
	})

	t.Run(testBufferContainsMessage, func(t *testing.T) {
		assert.Contains(t, testLogger.String(), expectedMessage, "test buffer should contain the formatted message")
	})
}

func TestBeginTest(t *testing.T) {
	BeginTest()
	defer EndTest()

	defer resetTestLogger()

	var testLogger bytes.Buffer
	updateLogger(&testLogger)

	BeginTest()

	t.Run(testBufferNotNil, func(t *testing.T) {
		assert.NotNil(t, testLogger, "test buffer should not be nil")
	})

	t.Run(testBufferNotEmpty, func(t *testing.T) {
		assert.True(t, testLogger.Len() > 0, "test buffer should not be empty")
	})

	t.Run(testBufferContainsLogLevel, func(t *testing.T) {
		assert.Contains(t, testLogger.String(), testDebugPrefix, "test buffer should contain TEST DEBUG log level")
	})

	t.Run("bufferContainsBeginTest", func(t *testing.T) {
		assert.Contains(t, testLogger.String(), "BEGIN TEST", "the BEGIN TEST string should be contained in the test buffer")
	})

	t.Run("bufferContainsFunctionName", func(t *testing.T) {
		assert.Contains(t, testLogger.String(), "testlog.TestBeginTest", "the test buffer should contain the package and function name")
	})

	t.Run("bufferContainsTestFilename", func(t *testing.T) {
		assert.Contains(t, testLogger.String(), "testlog_test.go")
	})
}

func TestEndTest(t *testing.T) {
	BeginTest()
	defer EndTest()

	defer resetTestLogger()

	var testLogger bytes.Buffer
	updateLogger(&testLogger)

	EndTest()

	t.Run(testBufferNotNil, func(t *testing.T) {
		assert.NotNil(t, testLogger, "test buffer should not be nil")
	})

	t.Run(testBufferNotEmpty, func(t *testing.T) {
		assert.True(t, testLogger.Len() > 0, "test buffer should not be empty")
	})

	t.Run(testBufferContainsLogLevel, func(t *testing.T) {
		assert.Contains(t, testLogger.String(), testDebugPrefix, "test buffer should contain TEST DEBUG log level")
	})

	t.Run("bufferContainsBeginTest", func(t *testing.T) {
		assert.Contains(t, testLogger.String(), "END TEST", "the END TEST string should be contained in the test buffer")
	})

	t.Run("bufferContainsFunctionName", func(t *testing.T) {
		assert.Contains(t, testLogger.String(), "testlog.TestEndTest", "the test buffer should contain the package and function name")
	})

	t.Run("bufferContainsTestFilename", func(t *testing.T) {
		assert.Contains(t, testLogger.String(), "testlog_test.go")
	})
}
