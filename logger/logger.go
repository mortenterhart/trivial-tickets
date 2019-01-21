package logger

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"runtime"
)

type LogLevel int

const (
	LevelInfo LogLevel = iota
	LevelWarning
	LevelError
	LevelFatal
)

var stdout = log.New(os.Stdout, "", log.LstdFlags)

func Info(v ...interface{}) {
	prependLogLevel(&v, LevelInfo)
	stdout.Println(v...)
}

func prependLogLevel(v *[]interface{}, level LogLevel) {
	levelPrefix := levelText(level)
	*v = append([]interface{}{levelPrefix}, *v...)
}

func prependFunctionName(v *[]interface{}) {
	functionName, file, line := getCallerFunctionName(3)
	text := fmt.Sprintf("%s (%s:%d):", functionName, file, line)
	*v = append([]interface{}{text}, *v...)
}

func levelText(level LogLevel) string {
	switch level {
	case LevelInfo:
		return "[INFO]"

	case LevelWarning:
		return "[WARNING]"

	case LevelError:
		return "[ERROR]"

	case LevelFatal:
		return "[FATAL ERROR]"
	}

	return "undefined"
}

func Infof(format string, v ...interface{}) {
	stdout.Printf(levelText(LevelInfo)+" "+format, v...)
}

func Warn(v ...interface{}) {
	prependLogLevel(&v, LevelWarning)
	stdout.Println(v...)
}

func Warnf(format string, v ...interface{}) {
	stdout.Printf(levelText(LevelWarning)+" "+format, v...)
}

func Error(v ...interface{}) {
	prependLogLevel(&v, LevelError)
	stdout.Println(v...)
}

func Errorf(format string, v ...interface{}) {
	stdout.Printf(levelText(LevelError)+" "+format, v...)
}

func Fatal(v ...interface{}) {
	prependLogLevel(&v, LevelFatal)
	stdout.Fatalln(v...)
}

func Fatalf(format string, v ...interface{}) {
	stdout.Fatalf(levelText(LevelFatal)+" "+format, v...)
}

func ApiRequest(request *http.Request) {
	Infof("received API request to %s (Method = %s, Host = %s, Content-Length = %d)",
		request.RequestURI, request.Method, request.Host, request.ContentLength)
}

func updateLogger(writer io.Writer) {
	stdout.SetOutput(writer)
}

func getCallerFunctionName(skipFrames int) (string, string, int) {
	frame := getFrame(skipFrames)
	return frame.Function, frame.File, frame.Line
}

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
