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

// Package logger provides a logging interface to the server
// supporting different log levels and options.
package logger

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

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
 * Package logger
 * Logging interface to the server supporting different
 * log levels and options
 */

// stdout is the default logger which writes its output
// to standard output. The standard flags are set for the
// logger so that each message is prepended with date and
// time.
var stdout = log.New(os.Stdout, "", log.LstdFlags)

// timeFormat is the representation of the time format the
// default logger in the `log` package uses with standard
// flags. It is used to precede the log messages written by
// the server error log with date and time.
const timeFormat = "2006/01/02 15:04:05"

// Info writes a standardized info message represented through
// the parameter v to the standard log with the INFO log level.
// The message is preceded by the log level string and appended
// by the function location (package, filename and line number)
// where the logging took place. The message is only written to
// the log if the log level is INFO.
func Info(v ...interface{}) {
	if canLog(structs.LevelInfo) {
		logln(structs.LevelInfo, v...)
	}
}

// Infof writes a standardized info message represented through
// a format string and the corresponding arguments to the standard
// log with the INFO log level. The message is preceded by the
// log level string and appended by the function location (package,
// filename and line number) where the logging took place. The
// message is only written to the log if the log level is INFO.
func Infof(format string, v ...interface{}) {
	if canLog(structs.LevelInfo) {
		logf(structs.LevelInfo, format, v...)
	}
}

// Warn writes a standardized warning message represented through
// the parameter v to the standard log with the WARNING log level.
// The message is preceded by the log level string and appended
// by the function location (package, filename and line number)
// where the logging took place. The message is only written to
// the log if the log level is WARNING or higher.
func Warn(v ...interface{}) {
	if canLog(structs.LevelWarning) {
		logln(structs.LevelWarning, v...)
	}
}

// Warnf writes a standardized warning message represented through
// a format string and the corresponding arguments to the standard
// log with the WARNING log level. The message is preceded by the
// log level string and appended by the function location (package,
// filename and line number) where the logging took place. The
// message is only written to the log if the log level is WARNING
// or higher.
func Warnf(format string, v ...interface{}) {
	if canLog(structs.LevelWarning) {
		logf(structs.LevelWarning, format, v...)
	}
}

// Error writes a standardized error message represented through
// the parameter v to the standard log with the ERROR log level.
// The message is preceded by the log level string and appended
// by the function location (package, filename and line number)
// where the logging took place. The message is only written to
// the log if the log level is ERROR or higher.
func Error(v ...interface{}) {
	if canLog(structs.LevelError) {
		logln(structs.LevelError, v...)
	}
}

// Errorf writes a standardized error message represented through
// a format string and the corresponding arguments to the standard
// log with the ERROR log level. The message is preceded by the
// log level string and appended by the function location (package,
// filename and line number) where the logging took place. The
// message is only written to the log if the log level is ERROR
// or higher.
func Errorf(format string, v ...interface{}) {
	if canLog(structs.LevelError) {
		logf(structs.LevelError, format, v...)
	}
}

func logf(level structs.LogLevel, format string, v ...interface{}) {
	stdout.Println(level.String(), fmt.Sprintf(format, v...), getLoggingLocationSuffix(4))
}

func logln(level structs.LogLevel, v ...interface{}) {
	stdout.Println(level.String(), joinWithoutNewline(v...), getLoggingLocationSuffix(4))
}

func joinWithoutNewline(v ...interface{}) string {
	str := fmt.Sprintln(v...)
	return str[:len(str)-1]
}

// FatalFunc is a representation of the function type of the
// log.Fatal and log.Fatalln functions.
type FatalFunc func(v ...interface{})

// fatalln defines the function used to print fatal errors. This
// is necessary to test the Fatal() function below since the
// log.Fatal and log.Fatalln function exit the program and thus
// tests are not possible.
var fatalln FatalFunc = stdout.Fatalln

// Fatal writes a standardized fatal error message represented
// through the parameter v to the standard log with the FATAL
// log level. The fatal error causes the program to stop execution
// and to terminate with an error. The log message is preceded
// by the log level string and appended by the function location
// (package, filename and line number) where the logging took place.
// Fatal error messages are always written to the log.
func Fatal(v ...interface{}) {
	fatalln(structs.LevelFatal.String(), joinWithoutNewline(v...), getLoggingLocationSuffix(3))
}

// Fatalf writes a standardized fatal error message represented
// through a format string and the corresponding arguments to the
// standard log with the FATAL log level. The fatal error causes
// the program to stop execution and to terminate with an error.
// The log message is preceded by the log level string and appended
// by the function location (package, filename and line number)
// where the logging took place. Fatal error messages are always
// written to the log.
func Fatalf(format string, v ...interface{}) {
	fatalln(structs.LevelFatal.String(), fmt.Sprintf(format, v...), getLoggingLocationSuffix(3))
}

// APIRequest writes a received API request to the standard log with
// the INFO log level. The message consists of the request URI, the
// used HTTP method, the host, the content type and length.
func APIRequest(request *http.Request) {
	Infof("received API request to %s (Method = %s, Host = %s, Content-Type = \"%s\", Content-Length = %d)",
		request.RequestURI, request.Method, request.Host, request.Header.Get("Content-Type"), request.ContentLength)
}

// errorLogWriter is a writer specifically designed for error logs
// such as the error log member of the http.Server struct. It writes
// formatted error messages with the ERROR log level and function
// location to the desired output, but only if the log level is
// ERROR or higher.
type errorLogWriter struct {
	// The output buffer where log messages are written to
	output io.Writer
}

// newErrorLogWriter constructs a new errorLogWriter writing its
// messages to the specified output buffer.
func newErrorLogWriter(output io.Writer) errorLogWriter {
	return errorLogWriter{output}
}

// Write writes the supplied message to the output buffer. The message
// is expanded with date and time stamp, the ERROR log level and the
// function location (package and function name, filename and line number).
// The message is only written if the configured log level is ERROR or
// higher.
func (writer errorLogWriter) Write(message []byte) (n int, err error) {
	if canLog(structs.LevelError) {
		timeStamp := time.Now().Format(timeFormat)

		return fmt.Fprintln(writer.output, timeStamp, structs.LevelError.String(), strings.Trim(string(message), "\n"),
			getLoggingLocationSuffix(6))
	}

	return 0, nil
}

// NewErrorLogger returns a new logger instance with an errorLogWriter
// as output. The errorLogWriter wraps os.Stderr as output and formats
// each log message according to the logging defaults in the `logger`
// package. This logger can be used for example as ErrorLog in a
// http.Server instance where the server should write formatted
// messages.
func NewErrorLogger() *log.Logger {
	return log.New(newErrorLogWriter(os.Stderr), "", 0)
}

// getLoggingLocationSuffix creates a suffix for log messages containing
// the package path, the function name, the filename and line number
// depending on the logging configuration. If the FullPaths option is
// set the package and file paths are not shortened. If the Verbose
// option is set the filename will be also written to the suffix instead
// of just package path and function name.
func getLoggingLocationSuffix(calldepth int) string {
	functionName, file, line := getCallerFunction(calldepth)

	if !globals.LogConfig.FullPaths {
		functionName = filepath.Base(functionName)
		file = filepath.Base(file)
	}

	if globals.LogConfig.Verbose {
		return fmt.Sprintf("[%s in %s:%d]", functionName, file, line)
	}

	return fmt.Sprintf("[%s:%d]", functionName, line)
}

// updateLogger sets a new output destination for the logger
// (only used within tests).
func updateLogger(writer io.Writer) {
	stdout.SetOutput(writer)
}

// canLog decides whether a log message with the specified log level will
// be output. This is the case if the specified log level is equal or
// higher than the configured log level.
func canLog(level structs.LogLevel) bool {
	return level >= globals.LogConfig.LogLevel
}

// getCallerFunction returns the package path and function name, the
// filename and the line number of the function call that is skipFrames
// calls away.
func getCallerFunction(skipFrames int) (function string, file string, line int) {
	frame := getFrame(skipFrames)
	return frame.Function, frame.File, frame.Line
}

// getFrame returns the frame that was created skipFrames numbers before.
// This is done by iterating the runtime stack frames until the desired
// function call is found.
// Taken from https://stackoverflow.com/a/35213181
func getFrame(skipFrames int) runtime.Frame {
	// We need the frame at index skipFrames+2, since we never want runtime.Callers and getFrame
	targetFrameIndex := skipFrames + 2

	// Set size to targetFrameIndex+2 to ensure we have room for one more caller than we need
	programCounters := make([]uintptr, targetFrameIndex+2)
	n := runtime.Callers(0, programCounters)

	frame := runtime.Frame{Function: "unknown"}
	if n > 0 {
		frames := runtime.CallersFrames(programCounters[:n])
		for more, frameIndex := true, 0; more && frameIndex <= targetFrameIndex; frameIndex++ {
			var frameCandidate runtime.Frame
			frameCandidate, more = frames.Next()
			if frameIndex == targetFrameIndex {
				frame = frameCandidate
			}
		}
	}

	return frame
}
