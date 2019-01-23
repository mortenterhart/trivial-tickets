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
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
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
 * Package testlog
 * Logging interface for tests
 */

// testLogger is the standard logger for writing the test log.
// Each log message is preceded with the current timestamp.
var testLogger = log.New(os.Stdout, "", log.LstdFlags)

// testDebugPrefix is the log level prefix for the the log level
// in `structs.LevelTestDebug`. Since the functions of this package
// are used in package `structs` the log level instance in that
// package cannot be used due to circle imports. Therefore this
// prefix serves as log level string.
const testDebugPrefix string = "[TEST DEBUG]"

// Debug writes a standardized test debug message represented
// through the parameter v to the test log with the TEST DEBUG
// log level. The message is preceded by the log level string and
// appended by the function location (package, filename and line
// number) where the logging took place. This logging routine is
// intended exclusively for use in test cases to output relevant
// debugging information to the test log.
func Debug(v ...interface{}) {
	testLogger.Println(testDebugPrefix, joinWithoutNewline(v...), getFunctionCallLocation(3))
}

// Debugf writes a standardized test debug message represented
// through a format string and the corresponding arguments to the
// test log with the TEST DEBUG log level. The message is preceded
// by the log level string and appended by the function location
// (package, filename and line number) where the logging took place.
// This logging routine is intended exclusively for use in test cases
// to output relevant debugging information to the test log.
func Debugf(format string, v ...interface{}) {
	testLogger.Println(testDebugPrefix, fmt.Sprintf(format, v...), getFunctionCallLocation(3))
}

// BeginTest writes a begin header text into the test log signalling
// that a new test function is starting. It should be called as the
// first statement of a test. The log message contains the package
// and function name of the test as well as the logging invocation
// location consisting of filename and line number. The function
// should be called only from tests.
func BeginTest() {
	function, file, line := getCallerShortened(2)
	testLogger.Printf("%s === BEGIN TEST %s === [in %s:%d]\n",
		testDebugPrefix, function, file, line)
}

// EndTest writes an end footer text into the test log signalling
// that the test function has finished. It should be called as the
// last statement of the test or as deferred call. The log message
// contains the package and function name of the test as well as
// the logging invocation location consisting of filename and line
// number. The function should be called only from tests.
func EndTest() {
	function, file, line := getCallerShortened(2)
	testLogger.Printf("%s === END TEST   %s === [in %s:%d]\n",
		testDebugPrefix, function, file, line)
}

// joinWithoutNewline concatenates all the provided arguments and
// joins them with a single space, but strips the trailing newline.
func joinWithoutNewline(v ...interface{}) string {
	str := fmt.Sprintln(v...)
	return str[:len(str)-1]
}

// getFunctionCallLocation returns the test log suffix containing
// the function name and line number. The suffix is appended to
// each test debug message.
func getFunctionCallLocation(calldepth int) string {
	function, _, line := getCallerShortened(calldepth)
	return fmt.Sprintf("[%s:%d]", function, line)
}

// updateLogger sets a new output destination for the logger
// (only used within tests).
func updateLogger(writer io.Writer) {
	testLogger.SetOutput(writer)
}

// getCallerShortened returns the shortened package path to the
// function, the shortened file path and the line of the function
// call that is calldepth steps away. The paths are shortened in
// a manner which removes any leading directory.
func getCallerShortened(calldepth int) (function string, file string, line int) {
	frame := getFrame(calldepth)
	return filepath.Base(frame.Function), filepath.Base(frame.File), frame.Line
}

// getFrame returns the frame that was created skipFrames numbers before.
// This is done by iterating the runtime stack frames until the desired
// function call is found.
// Taken from https://stackoverflow.com/a/35213181
func getFrame(skipFrames int) runtime.Frame {
	// We need the frame at index skipFrames+2, since we never want
	// runtime.Callers and getFrame
	targetFrameIndex := skipFrames + 2

	// Set size to targetFrameIndex+2 to ensure we have room for
	// one more caller than we need
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
