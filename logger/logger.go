package logger

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"runtime"

	"github.com/mortenterhart/trivial-tickets/globals"
	"github.com/mortenterhart/trivial-tickets/structs"
)

var stdout = log.New(os.Stdout, "", log.LstdFlags)

func Info(v ...interface{}) {
	if canLog(structs.LevelInfo) {
		prependLogLevel(&v, structs.LevelInfo)
		appendFunctionLocation(&v)

		stdout.Println(v...)
	}
}

func Infof(format string, v ...interface{}) {
	if canLog(structs.LevelInfo) {
		stdout.Printf(buildFormatString(structs.LevelInfo, format), v...)
	}
}

func Warn(v ...interface{}) {
	if canLog(structs.LevelWarning) {
		prependLogLevel(&v, structs.LevelWarning)
		appendFunctionLocation(&v)

		stdout.Println(v...)
	}
}

func Warnf(format string, v ...interface{}) {
	if canLog(structs.LevelWarning) {
		stdout.Printf(buildFormatString(structs.LevelWarning, format), v...)
	}
}

func Error(v ...interface{}) {
	if canLog(structs.LevelError) {
		prependLogLevel(&v, structs.LevelError)
		appendFunctionLocation(&v)

		stdout.Println(v...)
	}
}

func Errorf(format string, v ...interface{}) {
	if canLog(structs.LevelError) {
		stdout.Printf(buildFormatString(structs.LevelError, format), v...)
	}
}

type FatalFunc func(v ...interface{})
type FatalfFunc func(format string, v ...interface{})

var fatalln FatalFunc = stdout.Fatalln
var fatalf FatalfFunc = stdout.Fatalf

func Fatal(v ...interface{}) {
	prependLogLevel(&v, structs.LevelFatal)
	appendFunctionLocation(&v)

	fatalln(v...)
}

func Fatalf(format string, v ...interface{}) {
	fatalf(buildFormatString(structs.LevelFatal, format), v...)
}

func ApiRequest(request *http.Request) {
	Infof("received API request to %s (Method = %s, Host = %s, Content-Length = %d)",
		request.RequestURI, request.Method, request.Host, request.ContentLength)
}

func prependLogLevel(v *[]interface{}, level structs.LogLevel) {
	*v = append([]interface{}{level.String()}, *v...)
}

func appendFunctionLocation(v *[]interface{}) {
	*v = append(*v, getLoggingLocationSuffix())
}

var skipFrames = 4

func getLoggingLocationSuffix() string {
	functionName, file, line := getCallerFunctionName(skipFrames)

	if !globals.LogConfig.FullPaths {
		leadingSlashes := regexp.MustCompile("^.*/")
		functionName = leadingSlashes.ReplaceAllString(functionName, "")
		file = leadingSlashes.ReplaceAllString(file, "")
	}

	if globals.LogConfig.VerboseLog {
		return fmt.Sprintf("[%s in %s:%d]", functionName, file, line)
	}

	return fmt.Sprintf("[%s:%d]", functionName, line)
}

func buildFormatString(level structs.LogLevel, formatString string) string {
	return fmt.Sprintln(level.String(), formatString, getLoggingLocationSuffix())
}

func updateLogger(writer io.Writer) {
	stdout.SetOutput(writer)
}

func canLog(level structs.LogLevel) bool {
	return globals.LogConfig.LogLevel <= level
}

func getCallerFunctionName(skipFrames int) (string, string, int) {
	frame := getFrame(skipFrames)
	return frame.Function, frame.File, frame.Line
}

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
