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

// Command command_line_tool is the command-line interface to
// send and receive e-mails to/from the Ticketsystem server.
package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/mortenterhart/trivial-tickets/api/api_in"
	"github.com/mortenterhart/trivial-tickets/api/api_out"
	"github.com/mortenterhart/trivial-tickets/cli/client"
	localIo "github.com/mortenterhart/trivial-tickets/cli/io"
	"github.com/mortenterhart/trivial-tickets/globals"
	"github.com/mortenterhart/trivial-tickets/log"
	"github.com/mortenterhart/trivial-tickets/log/testlog"
	"github.com/mortenterhart/trivial-tickets/structs"
	"github.com/mortenterhart/trivial-tickets/structs/defaults"
	"github.com/mortenterhart/trivial-tickets/util/filehandler"
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
 * Main package of the command line utility
 */

// ICommandReader defines required methods for
// a reader being able to provide an input access
// for commands and strings that are needed at
// user input prompts.
type ICommandReader interface {
	// Embed the Reader interface to allow using
	// an ICommandReader as readable source for
	// a buffered reader.
	io.Reader

	// appendCommand appends a CLI command to
	// the reader's input buffer followed by
	// a newline to finish the line that is read.
	appendCommand(command structs.Command)

	// appendString appends a custom string to
	// the reader's input buffer followed by
	// a newline to finish the line that is read.
	appendString(data string)
}

// CommandReader is an instance of the ICommandReader
// interface. It can write commands and custom
// strings to its input buffer which can be read
// by the command-line tool.
type CommandReader struct {
	// input is the input buffer.
	input []byte
}

// newCommandReader returns a new command reader
// ready to use for input prompts.
func newCommandReader() (r *CommandReader) {
	return &CommandReader{}
}

// setInput can replace the input buffer with input.
func (r *CommandReader) setInput(input string) {
	r.input = []byte(input)
}

// appendCommand appends a command to the reader's input
// buffer followed by a newline to submit the input.
func (r *CommandReader) appendCommand(command structs.Command) {
	r.appendString(command.String())
}

// appendString appends a custom string to the reader's
// input buffer followed by a newline to submit the input.
func (r *CommandReader) appendString(data string) {
	r.input = append(r.input, []byte(data+"\n")...)
}

// readByte reads a single byte from the input buffer and
// cuts it off the input buffer.
func (r *CommandReader) readByte() byte {
	// this function assumes that eof() check was done before
	b := r.input[0]
	r.input = r.input[1:]
	return b
}

// eof checks if the reader has reached the end of
// the input buffer.
func (r *CommandReader) eof() (eof bool) {
	return len(r.input) == 0
}

// Read reads bytes of the input buffer to the slice p
// until the capacity of p exceeds. It returns the number
// of bytes read and an io.EOF error if the end of the
// buffer has reached previously.
func (r *CommandReader) Read(p []byte) (n int, err error) {
	if r.eof() {
		err = io.EOF
		return
	}

	if c := cap(p); c > 0 {
		for n < c {
			p[n] = r.readByte()
			n++
			if r.eof() {
				break
			}
		}
	}

	return
}

// dispatchInputs sets the command reader to the internal
// input reader used by the command-line tool. This operation
// is necessary so that the tool can read inputs other than
// from os.Stdin.
func (r *CommandReader) dispatchInputs() {
	localIo.Reader = r
}

// String returns the input buffer used by this command reader
// as string.
func (r *CommandReader) String() string {
	return string(r.input)
}

// defaultCliConfig returns the default CLI
// configuration.
//
// Note: The port has been changed to 8444 to
// avoid conflicts with the productive server
// that runs on port 8443 by default.
func defaultCliConfig() structs.CLIConfig {
	return structs.CLIConfig{
		Host: defaults.CliHost,
		Port: defaults.CliPort + 2,
		Cert: defaults.CliCertificate,
	}
}

// testServerConfig returns the test configuration
// of the server as struct.
//
// Note: The port has been changed to 8444 to
// avoid conflicts with the productive server
// that runs on port 8443 by default.
func testServerConfig() structs.ServerConfig {
	return structs.ServerConfig{
		Port:    defaults.TestPort + 1,
		Tickets: "testtickets",
		Users:   defaults.TestUsers,
		Mails:   "testmails",
		Cert:    defaults.TestCertificate,
		Key:     defaults.TestKey,
		Web:     defaults.TestWeb,
	}
}

// testLogConfig returns the test configuration
// for the logging instance.
func testLogConfig() structs.LogConfig {
	return structs.LogConfig{
		LogLevel:  structs.AsLogLevel(defaults.LogLevelString),
		Verbose:   defaults.LogVerbose,
		FullPaths: defaults.LogFullPaths,
	}
}

// getTestConfigs returns all standard settings
// for the server, the logger and the CLI.
func getTestConfigs() (serverConfig structs.ServerConfig, logConfig structs.LogConfig, cliConfig structs.CLIConfig) {
	return testServerConfig(), testLogConfig(), defaultCliConfig()
}

// resetAllConfigs is called as deferred call in all
// test cases to reset the configurations to their
// default values so that it can be used in following
// test cases.
func resetAllConfigs() {
	initializeAllConfigs()
}

// resetFlags resets all configuration flags to their
// default values and should be called as deferred call
// in functions which alter the flags.
func resetFlags() {
	cliConfig := defaultCliConfig()

	*host = cliConfig.Host
	*port = uint(cliConfig.Port)
	*cert = cliConfig.Cert

	*fetch = false

	*submit = false
	*email = ""
	*ticketID = ""
	*subject = ""
	*message = ""
}

//revive:disable:deep-exit

// TestMain is the main entry point to the tests and
// exits with the result of the tests. Before the tests
// are run, it initializes the configurations.
func TestMain(m *testing.M) {
	initializeAllConfigs()

	os.Exit(m.Run())
}

//revive:enable:deep-exit

// initializeAllConfigs initializes the server configuration,
// the logging and the CLI configuration with test settings.
func initializeAllConfigs() {
	config := testServerConfig()
	globals.ServerConfig = &config

	logConfig := testLogConfig()
	globals.LogConfig = &logConfig

	client.SetCLIConfig(defaultCliConfig())
}

// cleanupTestFiles removes the directories for test tickets
// and test mails, if they exist.
func cleanupTestFiles(config structs.ServerConfig) {
	if filehandler.DirectoryExists(config.Tickets) {
		testlog.Debug("Deferred: Removing test ticket directory", config.Tickets)
		if removeErr := os.RemoveAll(config.Tickets); removeErr != nil {
			testlog.Debug("test error: could not remove ticket directory:", removeErr)
		}
	}

	if filehandler.DirectoryExists(config.Mails) {
		testlog.Debug("Deferred: Removing test mail directory", config.Mails)
		if removeErr := os.RemoveAll(config.Mails); removeErr != nil {
			testlog.Debug("test error: could not remove mail directory:", removeErr)
		}
	}
}

// displayDirectoryContents prints all files contained in a directory
// to the test log. This can be useful if certain tests depend on
// directory contents and to see what is currently inside that directory.
func displayDirectoryContents(dirname string, contents []os.FileInfo) {
	testlog.Debugf("Directory contents of '%s' showing %d file(s):", dirname, len(contents))
	for index, file := range contents {
		testlog.Debugf("%d: %s", index, file.Name())
	}
}

func TestGetConfig(t *testing.T) {
	testlog.BeginTest()
	defer testlog.EndTest()

	defer resetAllConfigs()

	_, _, cliConfig := getTestConfigs()
	cliConfig.Port = defaults.CliPort

	conf, fetch, submit, mail := getConfig()
	assert.NotNil(t, conf)
	assert.Equal(t, cliConfig.Host, conf.Host)
	assert.Equal(t, cliConfig.Port, conf.Port)
	assert.Equal(t, cliConfig.Cert, conf.Cert)
	assert.False(t, fetch)
	assert.False(t, submit)
	assert.Equal(t, `{"from":"", "subject":"", "message":""}`, mail)
}

func TestGetConfigInvalidPort(t *testing.T) {
	testlog.BeginTest()
	defer testlog.EndTest()

	defer resetAllConfigs()
	defer resetFlags()

	*port = math.MaxUint16 + 1

	var fatalErrors bytes.Buffer
	fatal = func(v ...interface{}) {
		fatalErrors.WriteString(fmt.Sprint(v...))
	}

	clearBuffer := func() {
		fatalErrors.Reset()
	}
	defer clearBuffer()

	getConfig()

	t.Run("fatalBufferNotEmpty", func(t *testing.T) {
		assert.True(t, fatalErrors.Len() > 0, "fatal error buffer should not be empty")
	})

	t.Run("fatalErrorIsInvalidPort", func(t *testing.T) {
		assert.Contains(t, fatalErrors.String(), "not a valid port number", "the error should be the invalid port")
	})
}

func TestUsageMessage(t *testing.T) {
	testlog.BeginTest()
	defer testlog.EndTest()

	defer resetAllConfigs()

	var testBuffer bytes.Buffer
	flag.CommandLine.SetOutput(&testBuffer)

	usageMessage()

	t.Run("bufferBeginsWithUsage", func(t *testing.T) {
		assert.True(t, strings.HasPrefix(testBuffer.String(), fmt.Sprintf("Usage: %s [options]", filepath.Base(os.Args[0]))),
			"the usage message should begin with the usage string")
	})

	t.Run("bufferContainsHostOption", func(t *testing.T) {
		assert.Contains(t, testBuffer.String(), "-host", "the usage message should contain the -host option")
	})
}

// startTestServer creates a new server which has the address
// and port from the test configuration. It also has the handlers
// with the URLs of the Mail APIs used in the command-line tool's
// communication to the server registered. The server is started
// on a new Go routine and the function returns another function
// to shutdown the server gracefully with a timeout of 5 seconds.
func startTestServer(t *testing.T, serverConfig structs.ServerConfig) (shutdown func()) {
	// Setup the handlers in a serving multiplexer and provide
	// the URLs for the API handlers. This is necessary because
	// the command-line tool uses fixed URLs for its HTTP requests
	// and a httptest.NewServer() would have a static URL for
	// one single handler.
	handler := http.NewServeMux()
	handler.HandleFunc("/api/receive", api_in.ReceiveMail)
	handler.HandleFunc("/api/fetchMails", api_out.FetchMails)
	handler.HandleFunc("/api/verifyMail", api_out.VerifyMailSent)

	// Create a new HTTP server with the handler multiplexer and
	// the address defined in the server config
	server := &http.Server{
		Addr:     fmt.Sprintf("%s:%d", "localhost", serverConfig.Port),
		Handler:  handler,
		ErrorLog: log.NewErrorLogger(),
	}

	// Start the server on another go routine using TLS with the
	// configured SSL certificate and key file. Catch potential
	// server start errors with the channel startError.
	startError := make(chan error)
	go func() {
		serveError := server.ListenAndServeTLS(serverConfig.Cert, serverConfig.Key)
		startError <- serveError
		close(startError)
	}()

	// Return a function to shutdown the server gracefully
	shutdown = func() {
		// Create a new context with a timeout of 5 seconds
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Shutdown the server at the end of the tests with a timeout
		// of 5 seconds and report if there is an error
		shutdownErr := server.Shutdown(ctx)
		if shutdownErr != nil {
			t.Fatal(shutdownErr)
		}

		// Log potential start errors or server closed error
		testlog.Debug(<-startError)
	}
	return
}

// TestCommandLoop tests the loop in commandLoop() with all
// available commands.
func TestCommandLoop(t *testing.T) {
	testlog.BeginTest()
	defer testlog.EndTest()

	defer resetAllConfigs()

	// Get the default server and CLI configuration
	serverConfig, _, cliConfig := getTestConfigs()

	// Set correct path to the ssl certificate in the CLI config
	client.SetCLIConfig(structs.CLIConfig{
		Host: cliConfig.Host,
		Port: cliConfig.Port,
		Cert: defaults.TestCertificate,
	})

	shutdown := startTestServer(t, serverConfig)
	defer shutdown()
	defer cleanupTestFiles(serverConfig)

	// give the server enough time to start. Makes the test more reliable
	time.Sleep(1 * time.Second)

	// Test the submit command with valid user inputs
	t.Run("submitCommandValidInputs", func(t *testing.T) {
		// Create a new command reader that is filled
		// with the prompt inputs for a submit command
		commandReader := newCommandReader()

		// Pretend to be the user typing the commands and inputs
		// that are expected at the prompts

		// Type submit command
		commandReader.appendCommand(structs.CommandSubmit)
		// Type the email input
		commandReader.appendString("testuser@example.com")
		// Type invalid ticket id
		commandReader.appendString("invalid-id")
		// Type the subject
		commandReader.appendString("test subject")
		// Type the message
		commandReader.appendString("test message")
		// Type exit command so that the loop exits
		commandReader.appendCommand(structs.CommandExit)

		testlog.Debug("Using following user inputs:", commandReader.String())

		// Set the command reader to the internal input reader
		// where user inputs are read from
		commandReader.dispatchInputs()

		// Execute the submit command
		submitErr := commandLoop()

		t.Run("noSubmitError", func(t *testing.T) {
			assert.NoError(t, submitErr, "unexpected error: user inputs are valid and server is running")
		})

		t.Run("verifyTicketCreated", func(t *testing.T) {
			dirContents, readErr := ioutil.ReadDir(serverConfig.Tickets)

			t.Run("readError", func(t *testing.T) {
				assert.NoErrorf(t, readErr, "reading directory contents of '%s' raised an error", serverConfig.Tickets)
			})

			displayDirectoryContents(serverConfig.Tickets, dirContents)

			t.Run("oneFileCreated", func(t *testing.T) {
				assert.Equal(t, 1, len(dirContents), "ticket directory should contain exactly one ticket")
			})
		})

		t.Run("verifyMailCreated", func(t *testing.T) {
			dirContents, readErr := ioutil.ReadDir(serverConfig.Mails)

			t.Run("readError", func(t *testing.T) {
				assert.NoErrorf(t, readErr, "reading directory contents of '%s' raised an error", serverConfig.Mails)
			})

			displayDirectoryContents(serverConfig.Mails, dirContents)

			t.Run("oneFileCreated", func(t *testing.T) {
				assert.Equal(t, 1, len(dirContents), "mail directory should contain exactly one mail")
			})
		})
	})

	// Test the submit command with invalid user inputs
	t.Run("submitMailInvalidInputs", func(t *testing.T) {
		// Create a new command reader that is filled
		// with the prompt inputs for a submit command
		commandReader := newCommandReader()

		// Pretend to be the user typing the commands and inputs
		// that are expected at the prompts

		// Type submit command
		commandReader.appendCommand(structs.CommandSubmit)
		// Type an invalid email address
		commandReader.appendString("testuser")
		// Type invalid ticket id
		commandReader.appendString("invalid-id")
		// Type the subject
		commandReader.appendString("test subject")
		// Type the message
		commandReader.appendString("test message")
		// Type exit command so that the loop exits
		commandReader.appendCommand(structs.CommandExit)

		testlog.Debug("Using following user inputs:", commandReader.String())

		// Set the command reader to the internal input reader
		// where user inputs are read from
		commandReader.dispatchInputs()

		// Execute the submit command
		submitErr := commandLoop()

		t.Run("submitError", func(t *testing.T) {
			assert.Error(t, submitErr, "invalid email input should cause a submitting error")
		})
	})

	// Test the fetch command and fetch the previously
	// created mail by the create ticket action
	t.Run("fetchCommand", func(t *testing.T) {
		// Create a new command reader that is filled
		// with the command input of the fetch command
		commandReader := newCommandReader()

		// Type the fetch command as input
		commandReader.appendCommand(structs.CommandFetch)
		// Type the exit command to quit the CLI
		commandReader.appendCommand(structs.CommandExit)

		testlog.Debug("Using following user inputs:", commandReader.String())

		// Dispatch the inputs and write it to the
		// input reader
		commandReader.dispatchInputs()

		// Execute the fetch command
		fetchErr := commandLoop()

		t.Run("noFetchError", func(t *testing.T) {
			assert.NoError(t, fetchErr, "fetching emails while the server is running should not be an error")
		})

		t.Run("mailDeletedAfterVerification", func(t *testing.T) {
			dirContents, readErr := ioutil.ReadDir(serverConfig.Mails)

			t.Run("readError", func(t *testing.T) {
				assert.NoErrorf(t, readErr, "reading contents of directory '%s' should not raise an error", serverConfig.Mails)
			})

			displayDirectoryContents(serverConfig.Mails, dirContents)

			t.Run("mailDeleted", func(t *testing.T) {
				assert.Equal(t, 0, len(dirContents), "mail directory should contain no mails any more")
			})
		})
	})
}

func TestCommandLoopWithoutServer(t *testing.T) {
	testlog.BeginTest()
	defer testlog.EndTest()

	defer resetAllConfigs()

	// Get the default CLI configuration
	_, _, cliConfig := getTestConfigs()

	// Set correct path to the ssl certificate in the CLI config
	client.SetCLIConfig(structs.CLIConfig{
		Host: cliConfig.Host,
		Port: cliConfig.Port,
		Cert: defaults.TestCertificate,
	})

	// Test the submit command with valid user inputs
	t.Run("submitCommand", func(t *testing.T) {
		// Create a new command reader that is filled
		// with the prompt inputs for a submit command
		commandReader := newCommandReader()

		// Pretend to be the user typing the commands and inputs
		// that are expected at the prompts

		// Type submit command
		commandReader.appendCommand(structs.CommandSubmit)
		// Type the email input
		commandReader.appendString("testuser@example.com")
		// Type invalid ticket id
		commandReader.appendString("invalid-id")
		// Type the subject
		commandReader.appendString("test subject")
		// Type the message
		commandReader.appendString("test message")
		// Type exit command so that the loop exits
		commandReader.appendCommand(structs.CommandExit)

		testlog.Debug("Using following user inputs:", commandReader.String())

		// Set the command reader to the internal input reader
		// where user inputs are read from
		commandReader.dispatchInputs()

		// Execute the submit command
		submitErr := commandLoop()

		t.Run("submitError", func(t *testing.T) {
			assert.Error(t, submitErr, "There should be an error because the server is not running")
		})
	})

	// Test the fetch command without access to the server
	t.Run("fetchCommand", func(t *testing.T) {
		// Create a new command reader that is filled
		// with the command input of the fetch command
		commandReader := newCommandReader()

		// Type the fetch command as input
		commandReader.appendCommand(structs.CommandFetch)
		// Type the exit command to quit the CLI
		commandReader.appendCommand(structs.CommandExit)

		testlog.Debug("Using following user inputs:", commandReader.String())

		// Dispatch the inputs and write it to the
		// input reader
		commandReader.dispatchInputs()

		// Execute the fetch command
		fetchErr := commandLoop()

		t.Run("fetchError", func(t *testing.T) {
			assert.Error(t, fetchErr, "fetching should cause an error because the server is not accessible")
		})
	})

	// Test the input of an invalid command
	t.Run("invalidCommandInput", func(t *testing.T) {
		// Create a new command reader that is filled
		// with the invalid command
		commandReader := newCommandReader()

		// Write the invalid command 'fetch'
		// (the input expects a number as command)
		commandReader.appendString("fetch")

		// Dispatch the invalid input amd write it
		// to the input reader
		commandReader.dispatchInputs()

		testlog.Debug("Using following user inputs:", commandReader.String())

		// Execute the input with an invalid command
		inputErr := commandLoop()

		t.Run("inputError", func(t *testing.T) {
			assert.Errorf(t, inputErr, "user input '%s' should cause an error because it is an invalid command", "fetch")
		})
	})

	// Test the input of a valid integer, but a
	// command out of range of valid options
	t.Run("commandOutOfRange", func(t *testing.T) {
		// Create a new command reader that is filled
		// with the command '5' which does not exist
		commandReader := newCommandReader()

		// Write a valid integer, but a command
		// out of range to the reader
		commandReader.appendString("5")

		testlog.Debug("Using following user inputs:", commandReader.String())

		// Dispatch the invalid command and write it
		// to the input reader
		commandReader.dispatchInputs()

		// Execute the input with the invalid command
		commandErr := commandLoop()

		t.Run("commandError", func(t *testing.T) {
			assert.Errorf(t, commandErr, "invalid command '%d' should cause an error", 5)
		})
	})
}

func TestMainFunction(t *testing.T) {
	testlog.BeginTest()
	defer testlog.EndTest()

	defer resetAllConfigs()

	serverConfig, _, _ := getTestConfigs()

	shutdown := startTestServer(t, serverConfig)
	defer shutdown()
	defer cleanupTestFiles(serverConfig)

	// give the server enough time to start. Makes the test more reliable
	time.Sleep(1 * time.Second)

	var fatalErrors bytes.Buffer
	fatal = func(v ...interface{}) {
		fatalErrors.WriteString(fmt.Sprint(v...))
	}

	clearBuffer := func() {
		fatalErrors.Reset()
	}

	t.Run("noFlagsSet", func(t *testing.T) {
		defer clearBuffer()

		commandReader := newCommandReader()

		commandReader.appendCommand(structs.CommandExit)

		testlog.Debug("Using following user inputs:", commandReader.String())

		commandReader.dispatchInputs()

		main()

		t.Run("noFatalErrors", func(t *testing.T) {
			assert.Equal(t, 0, fatalErrors.Len(), "The fatal error buffer should be empty")
		})
	})

	t.Run("noFlagsSetError", func(t *testing.T) {
		defer resetFlags()
		defer clearBuffer()

		*cert = defaults.TestCertificate

		testlog.Debug("Using following flags:")
		testlog.Debug("  -cert", *cert)

		commandReader := newCommandReader()

		commandReader.appendCommand(structs.CommandFetch)

		testlog.Debug("Using following user inputs:", commandReader.String())

		commandReader.dispatchInputs()

		main()

		t.Run("fatalBufferNotEmpty", func(t *testing.T) {
			assert.True(t, fatalErrors.Len() > 0, "fatal error buffer should be not empty")
		})

		t.Run("fatalError", func(t *testing.T) {
			assert.Contains(t, fatalErrors.String(), structs.TooManyInputs,
				"There should be a fatal error because of too many wrong inputs")
		})
	})

	t.Run("submitFlagSet", func(t *testing.T) {
		defer resetFlags()
		defer clearBuffer()

		*submit = true
		*email = "testuser@example.com"
		*subject = "test subject"
		*message = "test message"

		*cert = defaults.TestCertificate

		testlog.Debug("Using following flags:")
		testlog.Debug("  -s", *submit)
		testlog.Debug("  -email", *email)
		testlog.Debug("  -subject", *subject)
		testlog.Debug("  -message", *message)
		testlog.Debug("  -cert", *cert)

		main()

		t.Run("fatalBufferEmpty", func(t *testing.T) {
			assert.Equal(t, 0, fatalErrors.Len(), "fatal error buffer should be empty")
		})

		t.Run("verifyTicketAndMailCreated", func(t *testing.T) {
			t.Run("ticketCreated", func(t *testing.T) {
				dirContents, readErr := ioutil.ReadDir(serverConfig.Tickets)

				t.Run("readError", func(t *testing.T) {
					assert.NoErrorf(t, readErr, "reading directory '%s' should not be an error", serverConfig.Tickets)
				})

				displayDirectoryContents(serverConfig.Tickets, dirContents)

				t.Run("oneTicketCreated", func(t *testing.T) {
					assert.Equal(t, 1, len(dirContents), "ticket directory should contain exactly one ticket")
				})
			})

			t.Run("mailCreated", func(t *testing.T) {
				dirContents, readErr := ioutil.ReadDir(serverConfig.Mails)

				t.Run("readError", func(t *testing.T) {
					assert.NoErrorf(t, readErr, "reading directory '%s' should not be an error", serverConfig.Mails)
				})

				displayDirectoryContents(serverConfig.Mails, dirContents)

				t.Run("oneMailCreated", func(t *testing.T) {
					assert.Equal(t, 1, len(dirContents), "mail directory should contain exactly one mail")
				})
			})
		})
	})

	t.Run("submitFlagError", func(t *testing.T) {
		defer resetFlags()
		defer clearBuffer()

		*host = "unknown-host"

		*submit = true
		*email = "testuser@example.com"
		*subject = "test subject"
		*message = "test message"

		*cert = defaults.TestCertificate

		testlog.Debug("Using following flags:")
		testlog.Debug("  -host", *host)
		testlog.Debug("  -s", *submit)
		testlog.Debug("  -email", *email)
		testlog.Debug("  -subject", *subject)
		testlog.Debug("  -message", *message)
		testlog.Debug("  -cert", *cert)

		main()

		t.Run("fatalBufferNotEmpty", func(t *testing.T) {
			assert.True(t, fatalErrors.Len() > 0, "fatal error buffer should be not empty")
		})

		t.Run("fatalErrorSendingPostRequest", func(t *testing.T) {
			assert.Containsf(t, fatalErrors.String(), "error sending post request",
				"fatal error expected because the host name '%s' does not exist", *host)
		})
	})

	t.Run("fetchFlagSet", func(t *testing.T) {
		defer resetFlags()
		defer clearBuffer()

		*fetch = true
		*cert = defaults.TestCertificate

		testlog.Debug("Using following flags:")
		testlog.Debug("  -f", *fetch)
		testlog.Debug("  -cert", *cert)

		main()

		t.Run("fatalBufferEmpty", func(t *testing.T) {
			assert.Equal(t, 0, fatalErrors.Len(), "fatal error buffer should be empty")
		})

		t.Run("mailFileDeletedDueToVerification", func(t *testing.T) {
			dirContents, readErr := ioutil.ReadDir(serverConfig.Mails)

			t.Run("readError", func(t *testing.T) {
				assert.NoErrorf(t, readErr, "reading directory '%s' should not be an error", serverConfig.Mails)
			})

			displayDirectoryContents(serverConfig.Mails, dirContents)

			t.Run("mailDeleted", func(t *testing.T) {
				assert.Equal(t, 0, len(dirContents), "mail directory should not contain any mail files more")
			})
		})
	})

	t.Run("fetchFlagError", func(t *testing.T) {
		defer resetFlags()
		defer clearBuffer()

		*host = "unknown-host"

		*fetch = true
		*cert = defaults.TestCertificate

		testlog.Debug("Using following flags:")
		testlog.Debug("  -host", *host)
		testlog.Debug("  -f", *fetch)
		testlog.Debug("  -cert", *cert)

		main()

		t.Run("fatalBufferNotEmpty", func(t *testing.T) {
			assert.True(t, fatalErrors.Len() > 0, "fatal error buffer should be not empty")
		})

		t.Run("fatalErrorSendingGetRequest", func(t *testing.T) {
			assert.Containsf(t, fatalErrors.String(), "error sending get request",
				"fatal error expected because the host name '%s' does not exist", *host)
		})
	})

	t.Run("submitAndFetchFlagSet", func(t *testing.T) {
		defer resetFlags()
		defer clearBuffer()

		*submit = true
		*fetch = true

		testlog.Debug("Using following flags:")
		testlog.Debug("  -s", *submit)
		testlog.Debug("  -f", *submit)

		main()

		t.Run("fatalBufferNotEmpty", func(t *testing.T) {
			assert.True(t, fatalErrors.Len() > 0)
		})

		t.Run("fatalErrorMessageIsOutsideRangeOptions", func(t *testing.T) {
			assert.Contains(t, fatalErrors.String(), structs.NoValidOption,
				"There should be a fatal error because both the submit and fetch flag were specified")
		})
	})
}
