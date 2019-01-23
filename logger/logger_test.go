package logger

import (
	"bytes"
	"fmt"
	"net/http/httptest"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/mortenterhart/trivial-tickets/globals"
	"github.com/mortenterhart/trivial-tickets/structs"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	initializeLogConfig()

	os.Exit(m.Run())
}

func initializeLogConfig() {
	logConfig := testLogConfig()
	globals.LogConfig = &logConfig
}

func resetLogConfig() {
	logConfig := testLogConfig()
	globals.LogConfig = &logConfig

	skipFrames = 4
}

func testLogConfig() structs.LogConfig {
	return structs.LogConfig{
		LogLevel:   structs.LevelInfo,
		VerboseLog: false,
		FullPaths:  false,
	}
}

func TestInfo(t *testing.T) {
	defer resetLogConfig()

	var testLogger bytes.Buffer
	updateLogger(&testLogger)

	message := "Testing the Info function of package logger"
	Info(message)

	t.Run("bufferNotNil", func(t *testing.T) {
		assert.NotNil(t, testLogger, "buffer should not be nil")
	})

	t.Run("bufferNotEmpty", func(t *testing.T) {
		assert.True(t, testLogger.Len() > 0, "buffer should not be empty")
	})

	t.Run("bufferContainsLogLevel", func(t *testing.T) {
		assert.Contains(t, testLogger.String(), structs.LevelInfo.String(), "buffer should contain INFO log level string")
	})

	t.Run("bufferContainsMessage", func(t *testing.T) {
		assert.Contains(t, testLogger.String(), message, "buffer should contain the message that was written to it")
	})

	// Write to the logger again with a lower log level and verify that
	// no messages are written to the logger with this level
	t.Run("lowerLogLevel", func(t *testing.T) {
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
	defer resetLogConfig()

	var testLogger bytes.Buffer
	updateLogger(&testLogger)

	format := "server listening on https://localhost:%d (PID = %d), type Ctrl-C to stop"
	arguments := []interface{}{8443, os.Getpid()}
	expectedMessage := fmt.Sprintf(format, arguments...)

	Infof(format, arguments...)

	t.Run("bufferNotNil", func(t *testing.T) {
		assert.NotNil(t, testLogger, "buffer should not be nil")
	})

	t.Run("bufferNotEmpty", func(t *testing.T) {
		assert.True(t, testLogger.Len() > 0, "buffer should not be empty")
	})

	t.Run("bufferContainsLogLevel", func(t *testing.T) {
		assert.Contains(t, testLogger.String(), structs.LevelInfo.String(), "buffer should contain INFO log level string")
	})

	t.Run("bufferContainsMessage", func(t *testing.T) {
		assert.Contains(t, testLogger.String(), expectedMessage, "buffer should contain the message that was written to it")
	})

	// Write to the logger again with a lower log level and verify that
	// no messages are written to the logger with this level
	t.Run("lowerLogLevel", func(t *testing.T) {
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
	defer resetLogConfig()

	var testLogger bytes.Buffer
	updateLogger(&testLogger)

	message := "Testing the Warn function of package logger"
	Warn(message)

	t.Run("bufferNotNil", func(t *testing.T) {
		assert.NotNil(t, testLogger, "buffer should not be nil")
	})

	t.Run("bufferNotEmpty", func(t *testing.T) {
		assert.True(t, testLogger.Len() > 0, "buffer should not be empty")
	})

	t.Run("bufferContainsLogLevel", func(t *testing.T) {
		assert.Contains(t, testLogger.String(), structs.LevelWarning.String(), "buffer should contain WARNING log level string")
	})

	t.Run("bufferContainsMessage", func(t *testing.T) {
		assert.Contains(t, testLogger.String(), message, "buffer should contain the message that was written to it")
	})

	// Write to the logger again with a lower log level and verify that
	// no messages are written to the logger with this level
	t.Run("lowerLogLevel", func(t *testing.T) {
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
	defer resetLogConfig()

	var testLogger bytes.Buffer
	updateLogger(&testLogger)

	format := "server listening on https://localhost:%d (PID = %d), type Ctrl-C to stop"
	arguments := []interface{}{8443, os.Getpid()}
	expectedMessage := fmt.Sprintf(format, arguments...)

	Warnf(format, arguments...)

	t.Run("bufferNotNil", func(t *testing.T) {
		assert.NotNil(t, testLogger, "buffer should not be nil")
	})

	t.Run("bufferNotEmpty", func(t *testing.T) {
		assert.True(t, testLogger.Len() > 0, "buffer should not be empty")
	})

	t.Run("bufferContainsLogLevel", func(t *testing.T) {
		assert.Contains(t, testLogger.String(), structs.LevelWarning.String(), "buffer should contain WARNING log level string")
	})

	t.Run("bufferContainsMessage", func(t *testing.T) {
		assert.Contains(t, testLogger.String(), expectedMessage, "buffer should contain the message that was written to it")
	})

	// Write to the logger again with a lower log level and verify that
	// no messages are written to the logger with this level
	t.Run("lowerLogLevel", func(t *testing.T) {
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
	defer resetLogConfig()

	var testLogger bytes.Buffer
	updateLogger(&testLogger)

	message := "Testing the Error function of package logger"
	Error(message)

	t.Run("bufferNotNil", func(t *testing.T) {
		assert.NotNil(t, testLogger, "buffer should not be nil")
	})

	t.Run("bufferNotEmpty", func(t *testing.T) {
		assert.True(t, testLogger.Len() > 0, "buffer should not be empty")
	})

	t.Run("bufferContainsLogLevel", func(t *testing.T) {
		assert.Contains(t, testLogger.String(), structs.LevelError.String(), "buffer should contain ERROR log level string")
	})

	t.Run("bufferContainsMessage", func(t *testing.T) {
		assert.Contains(t, testLogger.String(), message, "buffer should contain the message that was written to it")
	})

	// Write to the logger again with a lower log level and verify that
	// no messages are written to the logger with this level
	t.Run("lowerLogLevel", func(t *testing.T) {
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
	defer resetLogConfig()

	var testLogger bytes.Buffer
	updateLogger(&testLogger)

	format := "server listening on https://localhost:%d (PID = %d), type Ctrl-C to stop"
	arguments := []interface{}{8443, os.Getpid()}
	expectedMessage := fmt.Sprintf(format, arguments...)

	Errorf(format, arguments...)

	t.Run("bufferNotNil", func(t *testing.T) {
		assert.NotNil(t, testLogger, "buffer should not be nil")
	})

	t.Run("bufferNotEmpty", func(t *testing.T) {
		assert.True(t, testLogger.Len() > 0, "buffer should not be empty")
	})

	t.Run("bufferContainsLogLevel", func(t *testing.T) {
		assert.Contains(t, testLogger.String(), structs.LevelError.String(), "buffer should contain ERROR log level string")
	})

	t.Run("bufferContainsMessage", func(t *testing.T) {
		assert.Contains(t, testLogger.String(), expectedMessage, "buffer should contain the message that was written to it")
	})

	// Write to the logger again with a lower log level and verify that
	// no messages are written to the logger with this level
	t.Run("lowerLogLevel", func(t *testing.T) {
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
	defer resetLogConfig()

	fatalln = func(v ...interface{}) {
		prependLogLevel(&v, structs.LevelFatal)
		appendFunctionLocation(&v)

		stdout.Println(v...)
	}

	var testLogger bytes.Buffer
	updateLogger(&testLogger)

	message := "Testing the Fatal function of package logger"
	Fatal(message)

	t.Run("bufferNotNil", func(t *testing.T) {
		assert.NotNil(t, testLogger, "buffer should not be nil")
	})

	t.Run("bufferNotEmpty", func(t *testing.T) {
		assert.True(t, testLogger.Len() > 0, "buffer should not be empty")
	})

	t.Run("bufferContainsLogLevel", func(t *testing.T) {
		assert.Contains(t, testLogger.String(), structs.LevelFatal.String(), "buffer should contain FATAL log level string")
	})

	t.Run("bufferContainsMessage", func(t *testing.T) {
		assert.Contains(t, testLogger.String(), message, "buffer should contain the message that was written to it")
	})
}

func TestFatalf(t *testing.T) {
	defer resetLogConfig()

	fatalf = func(format string, v ...interface{}) {
		stdout.Printf(buildFormatString(structs.LevelFatal, format), v...)
	}

	var testLogger bytes.Buffer
	updateLogger(&testLogger)

	format := "server listening on https://localhost:%d (PID = %d), type Ctrl-C to stop"
	arguments := []interface{}{8443, os.Getpid()}
	expectedMessage := fmt.Sprintf(format, arguments...)

	Fatalf(format, arguments...)

	t.Run("bufferNotNil", func(t *testing.T) {
		assert.NotNil(t, testLogger, "buffer should not be nil")
	})

	t.Run("bufferNotEmpty", func(t *testing.T) {
		assert.True(t, testLogger.Len() > 0, "buffer should not be empty")
	})

	t.Run("bufferContainsLogLevel", func(t *testing.T) {
		assert.Contains(t, testLogger.String(), structs.LevelFatal.String(), "buffer should contain FATAL log level string")
	})

	t.Run("bufferContainsMessage", func(t *testing.T) {
		assert.Contains(t, testLogger.String(), expectedMessage, "buffer should contain the message that was written to it")
	})
}

func TestApiRequest(t *testing.T) {
	defer resetLogConfig()

	var testLogger bytes.Buffer
	updateLogger(&testLogger)

	globals.LogConfig.LogLevel = structs.LevelInfo

	request := httptest.NewRequest("POST", "/api/receive", strings.NewReader("{}"))
	request.Host = "localhost:8443"

	ApiRequest(request)

	t.Run("bufferNotNil", func(t *testing.T) {
		assert.NotNil(t, testLogger, "buffer should not be nil")
	})

	t.Run("bufferNotEmpty", func(t *testing.T) {
		assert.True(t, testLogger.Len() > 0, "buffer should not be empty")
	})

	t.Run("bufferContainsLogLevel", func(t *testing.T) {
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
	defer resetLogConfig()

	skipFrames = 2

	t.Run("fullPathsFalse", func(t *testing.T) {
		globals.LogConfig.FullPaths = false
		globals.LogConfig.VerboseLog = false

		logSuffix := getLoggingLocationSuffix()

		t.Run("suffixNotNil", func(t *testing.T) {
			assert.NotNil(t, logSuffix, "suffix should not be nil")
		})

		t.Run("suffixNotEmpty", func(t *testing.T) {
			assert.NotEmpty(t, logSuffix, "suffix should not be empty")
		})

		t.Run("containsNoSlashes", func(t *testing.T) {
			assert.NotContains(t, logSuffix, "/", "suffix should not contain slashes because full-paths option is disabled")
			assert.True(t, regexp.MustCompile("\\[logger\\.TestGetLoggingLocationSuffix\\.func[0-9]+:[0-9]+\\]").MatchString(logSuffix),
				"suffix should be enclosed by squared brackets and contain no slashes")
		})
	})

	t.Run("fullPathsTrue", func(t *testing.T) {
		globals.LogConfig.FullPaths = true
		globals.LogConfig.VerboseLog = false

		logSuffix := getLoggingLocationSuffix()

		t.Run("suffixNotNil", func(t *testing.T) {
			assert.NotNil(t, logSuffix, "suffix should be not nil")
		})

		t.Run("suffixNotEmpty", func(t *testing.T) {
			assert.NotEmpty(t, logSuffix, "suffix should not be empty")
		})

		t.Run("containsSlashes", func(t *testing.T) {
			assert.Contains(t, logSuffix, "/", "suffix should contain slashes because full-paths option is enabled")
			assert.True(t, regexp.MustCompile("\\[.+/logger\\.TestGetLoggingLocationSuffix\\.func[0-9]+:[0-9]+\\]").MatchString(logSuffix),
				"suffix should contain full package path with slashes")
		})
	})

	t.Run("verboseLogFalse", func(t *testing.T) {
		globals.LogConfig.FullPaths = true
		globals.LogConfig.VerboseLog = false

		logSuffix := getLoggingLocationSuffix()

		t.Run("suffixNotNil", func(t *testing.T) {
			assert.NotNil(t, logSuffix, "suffix should not be nil")
		})

		t.Run("suffixNotEmpty", func(t *testing.T) {
			assert.NotEmpty(t, logSuffix, "suffix should not be empty")
		})

		t.Run("containsNoFilename", func(t *testing.T) {
			assert.NotContains(t, logSuffix, "logger_test.go", "suffix should not contain filename")
			assert.True(t, regexp.MustCompile("\\[.+/logger\\.TestGetLoggingLocationSuffix\\.func[0-9]+:[0-9]+\\]").MatchString(logSuffix),
				"suffix should contain full package path and line number")
		})
	})

	t.Run("verboseLogTrue", func(t *testing.T) {
		globals.LogConfig.FullPaths = true
		globals.LogConfig.VerboseLog = true

		logSuffix := getLoggingLocationSuffix()

		t.Run("suffixNotNil", func(t *testing.T) {
			assert.NotNil(t, logSuffix, "suffix should not be nil")
		})

		t.Run("suffixNotEmpty", func(t *testing.T) {
			assert.NotEmpty(t, logSuffix, "suffix should not be empty")
		})

		t.Run("containsFilename", func(t *testing.T) {
			assert.Contains(t, logSuffix, "logger_test.go", "suffix should contain filename")

			// Note: The filename in the regex is matched with '/' and '\' in the character
			// class to be compatible to Unix and Windows systems
			assert.True(t, regexp.MustCompile("\\[.+/logger\\.TestGetLoggingLocationSuffix\\.func[0-9]+ in [/\\\\].+[/\\\\]logger_test.go:[0-9]+\\]").MatchString(logSuffix),
				"suffix should contain the full package path")
		})
	})
}
