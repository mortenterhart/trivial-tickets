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

// Package testlogger defines another logger exclusively used
// by tests.
package testlogger

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
 * Package testlogger
 * Logging interface for tests
 */

var testLog = log.New(os.Stdout, "", log.LstdFlags)

const testDebugLogLevel = "[TEST DEBUG]"

// Debug writes a standardized test debug message represented
// through the parameter v to the test log with the TEST DEBUG
// log level. The message is preceded by the log level string and
// appended by the function location (package, filename and line
// number) where the logging took place. This logging routine is
// intended exclusively for use in test cases to output relevant
// debugging information to the test log.
func Debug(v ...interface{}) {
	testLog.Println(testDebugLogLevel, joinWithoutNewline(v...), getFunctionCallLocation(4))
}

// Debugf writes a standardized test debug message represented
// through a format string and the corresponding arguments to the
// test log with the TEST DEBUG log level. The message is preceded
// by the log level string and appended by the function location
// (package, filename and line number) where the logging took place.
// This logging routine is intended exclusively for use in test cases
// to output relevant debugging information to the test log.
func Debugf(format string, v ...interface{}) {
	testLog.Println(testDebugLogLevel, fmt.Sprintf(format, v...), getFunctionCallLocation(4))
}

// BeginTest writes a begin header text into the test log signalling
// that a new test function is starting. It should be called as the
// first statement of a test. The log message contains the package
// and function name of the test as well as the logging invocation
// location consisting of filename and line number. The function
// should be called only from tests.
func BeginTest() {
	function, file, line := getCallerShortened(3)
	testLog.Printf("%s === BEGIN TEST %s === [in %s:%d]\n",
		testDebugLogLevel, function, file, line)
}

// EndTest writes an end footer text into the test log signalling
// that the test function has finished. It should be called as the
// last statement of the test or as deferred call. The log message
// contains the package and function name of the test as well as
// the logging invocation location consisting of filename and line
// number. The function should be called only from tests.
func EndTest() {
	function, file, line := getCallerShortened(3)
	testLog.Printf("%s === END TEST   %s === [in %s:%d]\n",
		testDebugLogLevel, function, file, line)
}

func joinWithoutNewline(v ...interface{}) string {
	str := fmt.Sprintln(v...)
	return str[:len(str)-1]
}

func getFunctionCallLocation(calldepth int) string {
	function, _, line := getCallerShortened(calldepth)
	return fmt.Sprintf("[%s:%d]", function, line)
}

// updateLogger sets a new output destination for the logger
// (only used within tests).
func updateLogger(writer io.Writer) {
	testLog.SetOutput(writer)
}

func getCallerShortened(calldepth int) (function string, file string, line int) {
	function, file, line = getCallerFunction(calldepth)
	return filepath.Base(function), filepath.Base(file), line
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
