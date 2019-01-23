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

// Package log provides a logging interface to the server
// supporting different log levels and options.
package log

import (
	"bytes"
	"fmt"
	"net/http/httptest"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mortenterhart/trivial-tickets/globals"
	"github.com/mortenterhart/trivial-tickets/log/testlog"
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
 * Package log [tests]
 * Logging interface to the server supporting different
 * log levels and options
 */

//revive:disable:deep-exit

// TestMain is started to run the tests and initializes the
// configuration before running the tests. The tests' exit
// status is returned as the overall exit status.
func TestMain(m *testing.M) {
	initializeLogConfig()

	os.Exit(m.Run())
}

//revive:enable:deep-exit

// initializeLogConfig initializes the logging
// configuration with test values.
func initializeLogConfig() {
	logConfig := testLogConfig()
	globals.LogConfig = &logConfig
}

// resetLogConfig resets the logging configuration
// to the default test values, sets the log
// destination back to stdout and restores the
// fatalln function.
func resetLogConfig() {
	logConfig := testLogConfig()
	globals.LogConfig = &logConfig

	updateLogger(os.Stdout)
	fatalln = logger.Fatalln
}

// testLogConfig returns a test logging
// configuration.
func testLogConfig() structs.LogConfig {
	return structs.LogConfig{
		LogLevel:  structs.AsLogLevel(defaults.LogLevelString),
		Verbose:   defaults.LogVerbose,
		FullPaths: defaults.LogFullPaths,
	}
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

	// testLowerLogLevel denominates the subtest of writing
	// to the test log buffer with a lower log level.
	testLowerLogLevel string = "lowerLogLevel"
)

// serverPort is the port printed in the formatting functions
// of the logger.
const serverPort uint16 = 8443

func TestInfo(t *testing.T) {
	testlog.BeginTest()
	defer testlog.EndTest()

	defer resetLogConfig()

	var testLogger bytes.Buffer
	updateLogger(&testLogger)

	message := "Testing the Info function of package logger"

	Info(message)

	t.Run(testBufferNotNil, func(t *testing.T) {
		assert.NotNil(t, testLogger, "buffer should not be nil")
	})

	t.Run(testBufferNotEmpty, func(t *testing.T) {
		assert.True(t, testLogger.Len() > 0, "buffer should not be empty")
	})

	t.Run(testBufferContainsLogLevel, func(t *testing.T) {
		assert.Contains(t, testLogger.String(), structs.LevelInfo.String(), "buffer should contain INFO log level string")
	})

	t.Run(testBufferContainsMessage, func(t *testing.T) {
		assert.Contains(t, testLogger.String(), message, "buffer should contain the message that was written to it")
	})

	// Write to the logger again with a lower log level and verify that
	// no messages are written to the logger with this level
	t.Run(testLowerLogLevel, func(t *testing.T) {
		globals.LogConfig.LogLevel = structs.LevelWarning

		length := testLogger.Len()
		content := testLogger.String()

		Info(message)

		t.Run("sameLength", func(t *testing.T) {
			assert.Equal(t, length, testLogger.Len(), "length should not be affected by logging with lower log level")
		})

		t.Run("sameContent", func(t *testing.T) {
			assert.Equal(t, content, testLogger.String(), "buffer content should not be affected by logging with lower log level")
		})
	})
}

func TestInfof(t *testing.T) {
	testlog.BeginTest()
	defer testlog.EndTest()

	defer resetLogConfig()

	var testLogger bytes.Buffer
	updateLogger(&testLogger)

	format := "server listening on https://localhost:%d (PID = %d), type Ctrl-C to stop"
	arguments := []interface{}{serverPort, os.Getpid()}
	expectedMessage := fmt.Sprintf(format, arguments...)

	Infof(format, arguments...)

	t.Run(testBufferNotNil, func(t *testing.T) {
		assert.NotNil(t, testLogger, "buffer should not be nil")
	})

	t.Run(testBufferNotEmpty, func(t *testing.T) {
		assert.True(t, testLogger.Len() > 0, "buffer should not be empty")
	})

	t.Run(testBufferContainsLogLevel, func(t *testing.T) {
		assert.Contains(t, testLogger.String(), structs.LevelInfo.String(), "buffer should contain INFO log level string")
	})

	t.Run(testBufferContainsMessage, func(t *testing.T) {
		assert.Contains(t, testLogger.String(), expectedMessage, "buffer should contain the message that was written to it")
	})

	// Write to the logger again with a lower log level and verify that
	// no messages are written to the logger with this level
	t.Run(testLowerLogLevel, func(t *testing.T) {
		globals.LogConfig.LogLevel = structs.LevelWarning

		length := testLogger.Len()
		content := testLogger.String()

		Infof(format, arguments...)

		t.Run("sameLength", func(t *testing.T) {
			assert.Equal(t, length, testLogger.Len(), "length should not be affected by logging with lower log level")
		})

		t.Run("sameContent", func(t *testing.T) {
			assert.Equal(t, content, testLogger.String(), "buffer content should not be affected by logging with lower log level")
		})
	})
}

func TestWarn(t *testing.T) {
	testlog.BeginTest()
	defer testlog.EndTest()

	defer resetLogConfig()

	var testLogger bytes.Buffer
	updateLogger(&testLogger)

	message := "Testing the Warn function of package logger"

	Warn(message)

	t.Run(testBufferNotNil, func(t *testing.T) {
		assert.NotNil(t, testLogger, "buffer should not be nil")
	})

	t.Run(testBufferNotEmpty, func(t *testing.T) {
		assert.True(t, testLogger.Len() > 0, "buffer should not be empty")
	})

	t.Run(testBufferContainsLogLevel, func(t *testing.T) {
		assert.Contains(t, testLogger.String(), structs.LevelWarning.String(), "buffer should contain WARNING log level string")
	})

	t.Run(testBufferContainsMessage, func(t *testing.T) {
		assert.Contains(t, testLogger.String(), message, "buffer should contain the message that was written to it")
	})

	// Write to the logger again with a lower log level and verify that
	// no messages are written to the logger with this level
	t.Run(testLowerLogLevel, func(t *testing.T) {
		globals.LogConfig.LogLevel = structs.LevelError

		length := testLogger.Len()
		content := testLogger.String()

		Warn(message)

		t.Run("sameLength", func(t *testing.T) {
			assert.Equal(t, length, testLogger.Len(), "length should not be affected by logging with lower log level")
		})

		t.Run("sameContent", func(t *testing.T) {
			assert.Equal(t, content, testLogger.String(), "buffer content should not be affected by logging with lower log level")
		})
	})
}

func TestWarnf(t *testing.T) {
	testlog.BeginTest()
	defer testlog.EndTest()

	defer resetLogConfig()

	var testLogger bytes.Buffer
	updateLogger(&testLogger)

	format := "server listening on https://localhost:%d (PID = %d), type Ctrl-C to stop"
	arguments := []interface{}{serverPort, os.Getpid()}
	expectedMessage := fmt.Sprintf(format, arguments...)

	Warnf(format, arguments...)

	t.Run(testBufferNotNil, func(t *testing.T) {
		assert.NotNil(t, testLogger, "buffer should not be nil")
	})

	t.Run(testBufferNotEmpty, func(t *testing.T) {
		assert.True(t, testLogger.Len() > 0, "buffer should not be empty")
	})

	t.Run(testBufferContainsLogLevel, func(t *testing.T) {
		assert.Contains(t, testLogger.String(), structs.LevelWarning.String(), "buffer should contain WARNING log level string")
	})

	t.Run(testBufferContainsMessage, func(t *testing.T) {
		assert.Contains(t, testLogger.String(), expectedMessage, "buffer should contain the message that was written to it")
	})

	// Write to the logger again with a lower log level and verify that
	// no messages are written to the logger with this level
	t.Run(testLowerLogLevel, func(t *testing.T) {
		globals.LogConfig.LogLevel = structs.LevelError

		length := testLogger.Len()
		content := testLogger.String()

		Warnf(format, arguments...)

		t.Run("sameLength", func(t *testing.T) {
			assert.Equal(t, length, testLogger.Len(), "length should not be affected by logging with lower log level")
		})

		t.Run("sameContent", func(t *testing.T) {
			assert.Equal(t, content, testLogger.String(), "buffer content should not be affected by logging with lower log level")
		})
	})
}

func TestError(t *testing.T) {
	testlog.BeginTest()
	defer testlog.EndTest()

	defer resetLogConfig()

	var testLogger bytes.Buffer
	updateLogger(&testLogger)

	message := "Testing the Error function of package logger"

	Error(message)

	t.Run(testBufferNotNil, func(t *testing.T) {
		assert.NotNil(t, testLogger, "buffer should not be nil")
	})

	t.Run(testBufferNotEmpty, func(t *testing.T) {
		assert.True(t, testLogger.Len() > 0, "buffer should not be empty")
	})

	t.Run(testBufferContainsLogLevel, func(t *testing.T) {
		assert.Contains(t, testLogger.String(), structs.LevelError.String(), "buffer should contain ERROR log level string")
	})

	t.Run(testBufferContainsMessage, func(t *testing.T) {
		assert.Contains(t, testLogger.String(), message, "buffer should contain the message that was written to it")
	})

	// Write to the logger again with a lower log level and verify that
	// no messages are written to the logger with this level
	t.Run(testLowerLogLevel, func(t *testing.T) {
		globals.LogConfig.LogLevel = structs.LevelFatal

		length := testLogger.Len()
		content := testLogger.String()

		Error(message)

		t.Run("sameLength", func(t *testing.T) {
			assert.Equal(t, length, testLogger.Len(), "length should not be affected by logging with lower log level")
		})

		t.Run("sameContent", func(t *testing.T) {
			assert.Equal(t, content, testLogger.String(), "buffer content should not be affected by logging with lower log level")
		})
	})
}

func TestErrorf(t *testing.T) {
	testlog.BeginTest()
	defer testlog.EndTest()

	defer resetLogConfig()

	var testLogger bytes.Buffer
	updateLogger(&testLogger)

	format := "server listening on https://localhost:%d (PID = %d), type Ctrl-C to stop"
	arguments := []interface{}{serverPort, os.Getpid()}
	expectedMessage := fmt.Sprintf(format, arguments...)

	Errorf(format, arguments...)

	t.Run(testBufferNotNil, func(t *testing.T) {
		assert.NotNil(t, testLogger, "buffer should not be nil")
	})

	t.Run(testBufferNotEmpty, func(t *testing.T) {
		assert.True(t, testLogger.Len() > 0, "buffer should not be empty")
	})

	t.Run(testBufferContainsLogLevel, func(t *testing.T) {
		assert.Contains(t, testLogger.String(), structs.LevelError.String(), "buffer should contain ERROR log level string")
	})

	t.Run(testBufferContainsMessage, func(t *testing.T) {
		assert.Contains(t, testLogger.String(), expectedMessage, "buffer should contain the message that was written to it")
	})

	// Write to the logger again with a lower log level and verify that
	// no messages are written to the logger with this level
	t.Run(testLowerLogLevel, func(t *testing.T) {
		globals.LogConfig.LogLevel = structs.LevelFatal

		length := testLogger.Len()
		content := testLogger.String()

		Errorf(format, arguments...)

		t.Run("sameLength", func(t *testing.T) {
			assert.Equal(t, length, testLogger.Len(), "length should not be affected by logging with lower log level")
		})

		t.Run("sameContent", func(t *testing.T) {
			assert.Equal(t, content, testLogger.String(), "buffer content should not be affected by logging with lower log level")
		})
	})
}

func TestFatal(t *testing.T) {
	testlog.BeginTest()
	defer testlog.EndTest()

	defer resetLogConfig()

	fatalln = func(v ...interface{}) {
		logln(structs.LevelFatal, v...)
	}

	var testLogger bytes.Buffer
	updateLogger(&testLogger)

	message := "Testing the Fatal function of package logger"

	Fatal(message)

	t.Run(testBufferNotNil, func(t *testing.T) {
		assert.NotNil(t, testLogger, "buffer should not be nil")
	})

	t.Run(testBufferNotEmpty, func(t *testing.T) {
		assert.True(t, testLogger.Len() > 0, "buffer should not be empty")
	})

	t.Run(testBufferContainsLogLevel, func(t *testing.T) {
		assert.Contains(t, testLogger.String(), structs.LevelFatal.String(), "buffer should contain FATAL log level string")
	})

	t.Run(testBufferContainsMessage, func(t *testing.T) {
		assert.Contains(t, testLogger.String(), message, "buffer should contain the message that was written to it")
	})
}

func TestFatalf(t *testing.T) {
	testlog.BeginTest()
	defer testlog.EndTest()

	defer resetLogConfig()

	fatalln = func(v ...interface{}) {
		logln(structs.LevelFatal, v...)
	}

	var testLogger bytes.Buffer
	updateLogger(&testLogger)

	format := "server listening on https://localhost:%d (PID = %d), type Ctrl-C to stop"
	arguments := []interface{}{serverPort, os.Getpid()}
	expectedMessage := fmt.Sprintf(format, arguments...)

	Fatalf(format, arguments...)

	t.Run(testBufferNotNil, func(t *testing.T) {
		assert.NotNil(t, testLogger, "buffer should not be nil")
	})

	t.Run(testBufferNotEmpty, func(t *testing.T) {
		assert.True(t, testLogger.Len() > 0, "buffer should not be empty")
	})

	t.Run(testBufferContainsLogLevel, func(t *testing.T) {
		assert.Contains(t, testLogger.String(), structs.LevelFatal.String(), "buffer should contain FATAL log level string")
	})

	t.Run(testBufferContainsMessage, func(t *testing.T) {
		assert.Contains(t, testLogger.String(), expectedMessage, "buffer should contain the message that was written to it")
	})
}

func TestApiRequest(t *testing.T) {
	testlog.BeginTest()
	defer testlog.EndTest()

	defer resetLogConfig()

	var testLogger bytes.Buffer
	updateLogger(&testLogger)

	globals.LogConfig.LogLevel = structs.LevelInfo

	request := httptest.NewRequest("POST", "/api/receive", strings.NewReader("{}"))
	request.Host = "localhost:8443"

	APIRequest(request)

	t.Run(testBufferNotNil, func(t *testing.T) {
		assert.NotNil(t, testLogger, "buffer should not be nil")
	})

	t.Run(testBufferNotEmpty, func(t *testing.T) {
		assert.True(t, testLogger.Len() > 0, "buffer should not be empty")
	})

	t.Run(testBufferContainsLogLevel, func(t *testing.T) {
		assert.Contains(t, testLogger.String(), structs.LevelInfo.String(), "buffer should contain INFO log level string")
	})

	t.Run("bufferContainsRequestURI", func(t *testing.T) {
		expectedRequestURI := "/api/receive"

		assert.Contains(t, testLogger.String(), expectedRequestURI, "buffer should contain expected request URI in the log message")
	})

	t.Run("bufferContainsLogMessage", func(t *testing.T) {
		expectedMessage := "received API request to /api/receive"

		assert.Contains(t, testLogger.String(), expectedMessage, "buffer should contain expected message that was logged")
	})

	t.Run("bufferContainsMethod", func(t *testing.T) {
		expectedMethod := "Method = POST"

		assert.Contains(t, testLogger.String(), expectedMethod, "buffer should contain expected method in the log message")
	})

	t.Run("bufferContainsHost", func(t *testing.T) {
		expectedHost := "Host = localhost:8443"

		assert.Contains(t, testLogger.String(), expectedHost, "buffer should contain expected host in the log message")
	})

	t.Run("bufferContainsHost", func(t *testing.T) {
		expectedContentLength := "Content-Length = 2"

		assert.Contains(t, testLogger.String(), expectedContentLength, "buffer should contain expected content length in the log message")
	})
}

func TestGetLoggingLocationSuffix(t *testing.T) {
	testlog.BeginTest()
	defer testlog.EndTest()

	defer resetLogConfig()

	t.Run("fullPathsFalse", func(t *testing.T) {
		globals.LogConfig.FullPaths = false
		globals.LogConfig.Verbose = false

		logSuffix := getLoggingLocationSuffix(2)

		t.Run("suffixNotNil", func(t *testing.T) {
			assert.NotNil(t, logSuffix, "suffix should not be nil")
		})

		t.Run("suffixNotEmpty", func(t *testing.T) {
			assert.NotEmpty(t, logSuffix, "suffix should not be empty")
		})

		t.Run("containsNoSlashes", func(t *testing.T) {
			assert.NotContains(t, logSuffix, "/", "suffix should not contain slashes because full-paths option is disabled")
			assert.True(t, regexp.MustCompile("\\[log\\.TestGetLoggingLocationSuffix\\.func[0-9]+:[0-9]+\\]").MatchString(logSuffix),
				"suffix should be enclosed by squared brackets and contain no slashes")
		})
	})

	t.Run("fullPathsTrue", func(t *testing.T) {
		globals.LogConfig.FullPaths = true
		globals.LogConfig.Verbose = false

		logSuffix := getLoggingLocationSuffix(2)

		t.Run("suffixNotNil", func(t *testing.T) {
			assert.NotNil(t, logSuffix, "suffix should be not nil")
		})

		t.Run("suffixNotEmpty", func(t *testing.T) {
			assert.NotEmpty(t, logSuffix, "suffix should not be empty")
		})

		t.Run("containsSlashes", func(t *testing.T) {
			assert.Contains(t, logSuffix, "/", "suffix should contain slashes because full-paths option is enabled")
			assert.True(t, regexp.MustCompile("\\[.+/log\\.TestGetLoggingLocationSuffix\\.func[0-9]+:[0-9]+\\]").MatchString(logSuffix),
				"suffix should contain full package path with slashes")
		})
	})

	t.Run("verboseLogFalse", func(t *testing.T) {
		globals.LogConfig.FullPaths = true
		globals.LogConfig.Verbose = false

		logSuffix := getLoggingLocationSuffix(2)

		t.Run("suffixNotNil", func(t *testing.T) {
			assert.NotNil(t, logSuffix, "suffix should not be nil")
		})

		t.Run("suffixNotEmpty", func(t *testing.T) {
			assert.NotEmpty(t, logSuffix, "suffix should not be empty")
		})

		t.Run("containsNoFilename", func(t *testing.T) {
			assert.NotContains(t, logSuffix, "log_test.go", "suffix should not contain filename")
			assert.True(t, regexp.MustCompile("\\[.+/log\\.TestGetLoggingLocationSuffix\\.func[0-9]+:[0-9]+\\]").MatchString(logSuffix),
				"suffix should contain full package path and line number")
		})
	})

	t.Run("verboseLogTrue", func(t *testing.T) {
		globals.LogConfig.FullPaths = true
		globals.LogConfig.Verbose = true

		logSuffix := getLoggingLocationSuffix(2)

		t.Run("suffixNotNil", func(t *testing.T) {
			assert.NotNil(t, logSuffix, "suffix should not be nil")
		})

		t.Run("suffixNotEmpty", func(t *testing.T) {
			assert.NotEmpty(t, logSuffix, "suffix should not be empty")
		})

		t.Run("containsFilename", func(t *testing.T) {
			assert.Contains(t, logSuffix, "log_test.go", "suffix should contain filename")

			// Note: The filename in the regex is matched with '/' and '\' in the character class
			// and drive letters (e.g. 'C:\...') to be compatible to Unix and Windows systems
			assert.True(t, regexp.MustCompile("\\[.+/log\\.TestGetLoggingLocationSuffix\\.func[0-9]+ in ([/\\\\]|[A-Z]:).+[/\\\\]log_test.go:[0-9]+\\]").MatchString(logSuffix),
				"suffix should contain the full package path")
		})
	})
}

func TestErrorLogWriter(t *testing.T) {
	testlog.BeginTest()
	defer testlog.EndTest()

	defer resetLogConfig()

	t.Run("constructing", func(t *testing.T) {
		errorWriter := newErrorLogWriter(os.Stderr)

		t.Run("notNil", func(t *testing.T) {
			assert.NotNil(t, errorWriter, "constructor should return a not-nil instance of errorLogWriter")
		})

		t.Run("notEmpty", func(t *testing.T) {
			assert.NotEqual(t, errorLogWriter{}, errorWriter, "errorWriter should not be an empty writer")
		})

		t.Run("writesToStdout", func(t *testing.T) {
			assert.Equal(t, os.Stderr, errorWriter.output, "output channel of errorWriter should be stderr")
		})
	})

	t.Run("writing", func(t *testing.T) {
		var buffer bytes.Buffer
		errorWriter := newErrorLogWriter(&buffer)

		errorMessage := "TLS handshake failed"
		n, err := errorWriter.Write([]byte(errorMessage))

		t.Run("noError", func(t *testing.T) {
			assert.NoError(t, err, "writing to stdout should not cause an error")
		})

		t.Run("equalBufferLength", func(t *testing.T) {
			assert.Equal(t, buffer.Len(), n, "n should be equal to the buffer length")
		})

		t.Run("bufferContainsMessage", func(t *testing.T) {
			assert.Contains(t, buffer.String(), errorMessage, "buffer should contain the written error message alongside with prefix and suffix")
		})
	})

	t.Run("writingWithFatalLogLevel", func(t *testing.T) {
		globals.LogConfig.LogLevel = structs.LevelFatal

		var buffer bytes.Buffer
		errorWriter := newErrorLogWriter(&buffer)

		errorMessage := "error message"
		n, err := errorWriter.Write([]byte(errorMessage))

		t.Run("noError", func(t *testing.T) {
			assert.NoError(t, err, "writing to buffer should not cause an error")
		})

		t.Run("noBytesWritten", func(t *testing.T) {
			assert.Equal(t, 0, n, "there should be no bytes written to the buffer")
		})

		t.Run("emptyBuffer", func(t *testing.T) {
			assert.Equal(t, 0, buffer.Len(), "the buffer should contain nothing")
		})
	})
}

func TestNewErrorLogger(t *testing.T) {
	testlog.BeginTest()
	defer testlog.EndTest()

	errorLogger := NewErrorLogger()

	t.Run("notNil", func(t *testing.T) {
		assert.NotNil(t, errorLogger, "new logger should not be nil")
	})

	t.Run("noPrefix", func(t *testing.T) {
		assert.Empty(t, errorLogger.Prefix(), "prefix should be empty")
	})

	t.Run("noFlags", func(t *testing.T) {
		assert.Equal(t, 0, errorLogger.Flags(), "logger should not have flags set")
	})
}
