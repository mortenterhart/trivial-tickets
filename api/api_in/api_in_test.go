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

// Package api_in implements a web interface for incoming mails
// to create new tickets or answers
package api_in

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mortenterhart/trivial-tickets/globals"
	"github.com/mortenterhart/trivial-tickets/log"
	"github.com/mortenterhart/trivial-tickets/log/testlog"
	"github.com/mortenterhart/trivial-tickets/structs"
	"github.com/mortenterhart/trivial-tickets/structs/defaults"
	"github.com/mortenterhart/trivial-tickets/ticket"
	"github.com/mortenterhart/trivial-tickets/util/filehandler"
	"github.com/mortenterhart/trivial-tickets/util/jsontools"
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
 * Package api_in [tests]
 * Web API for incoming mails to create new tickets or answers
 */

// jsonContentTypeTest is the content type for the JSON format
// used in POST requests.
const jsonContentTypeTest string = "application/json"

// serverSetupHandler is a handler wrapper for initializing the
// server configuration and setting up the server.
type serverSetupHandler struct {
	// callUnderlying is the underlying handler that
	// is going to be tested
	callUnderlying http.HandlerFunc
}

// newSetupHandler creates a new setup wrapper with the
// specified handler wrapped inside and called after
// initialization.
func newSetupHandler(wrappedHandler http.HandlerFunc) serverSetupHandler {
	return serverSetupHandler{wrappedHandler}
}

// ServeHTTP is the handler function of the setup handler
// which initializes the server config with a test configuration
// and then calls the underlying handler.
func (handler serverSetupHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	config := testServerConfig()
	globals.ServerConfig = &config

	handler.callUnderlying(writer, request)
}

// teardownFunc is the type for a returned function that does
// cleanups and actions after the tests ran.
type teardownFunc func()

// setupAndCleanup setups the the test prerequisites such as
// server and logging configuration and returns a teardown function
// which cleanups test resources and files created by the tests.
// The teardown function is advised to be called as deferred call.
func setupAndCleanup() teardownFunc {
	testlog.Debug("Initializing server and logging configuration with default values")

	config := testServerConfig()
	globals.ServerConfig = &config

	logConfig := testLogConfig()
	globals.LogConfig = &logConfig

	return func() {
		cleanupTestFiles(config)
	}
}

// testServerConfig returns a test server configuration.
func testServerConfig() structs.ServerConfig {
	return structs.ServerConfig{
		Port:    defaults.TestPort,
		Tickets: "testtickets",
		Users:   defaults.TestUsers,
		Mails:   "testmails",
		Cert:    defaults.TestCertificate,
		Key:     defaults.TestKey,
		Web:     defaults.TestWeb,
	}
}

// testLogConfig returns a test logging configuration.
func testLogConfig() structs.LogConfig {
	return structs.LogConfig{
		LogLevel:  structs.AsLogLevel(defaults.LogLevelString),
		Verbose:   defaults.LogVerbose,
		FullPaths: defaults.LogFullPaths,
	}
}

// cleanupTestFiles removes all created temporary tickets
// and mails, if any.
func cleanupTestFiles(config structs.ServerConfig) {
	if filehandler.DirectoryExists(config.Tickets) {
		testlog.Debug("Deferred: Removing test ticket directory", config.Tickets)
		if removeErr := os.RemoveAll(config.Tickets); removeErr != nil {
			testlog.Debug("ERROR: cannot remove test ticket directory:", removeErr)
		}
	}

	if filehandler.DirectoryExists(config.Mails) {
		testlog.Debug("Deferred: Removing test mail directory", config.Tickets)
		if removeErr := os.RemoveAll(config.Mails); removeErr != nil {
			log.Error("ERROR: cannot remove test mail directory:", removeErr)
		}
	}
}

// createTestServer creates a test server with the given
// handler registered.
func createTestServer(handler http.Handler) *httptest.Server {
	return httptest.NewServer(handler)
}

// createReader creates a reader from a string suitable to
// use as a request body.
func createReader(data string) io.Reader {
	return strings.NewReader(data)
}

// buildExpectedJSON creates a json string from a json properties
// map that can be used to compare the actual API json response
// to the expected json.
func buildExpectedJSON(properties structs.JSONMap) []byte {
	expected := jsontools.MapToJSON(properties)
	return append(expected, '\n')
}

// logResponseBody reads the body of a response and logs it to the
// test log. Warning: After a call to this function the response body
// cannot be read again because it is a one-time reader. If it is
// required to read the body after this function, consider making a
// copy of the response and pass it to this function.
func logResponseBody(response *http.Response) {
	body, readErr := ioutil.ReadAll(response.Body)
	if readErr != nil {
		log.Error("response could not be read:", readErr)
		return
	}

	testlog.Debug(string(body))
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

//revive:disable:deep-exit

// TestMain is started to run the tests and initializes the
// configuration before running the tests. The tests' exit
// status is returned as the overall exit status.
func TestMain(m *testing.M) {
	os.Exit(runTests(m))
}

//revive:enable:deep-exit

// runTests setups the server and logging configuration and executes
// all the tests that were triggered by a call to 'go test'. It ensures
// that deferred calls such as the cleanup function are executed before
// the test binary is quit with `os.Exit()`.
func runTests(m *testing.M) (code int) {
	cleanupFiles := setupAndCleanup()
	defer cleanupFiles()

	return m.Run()
}

// * ------------------------------------------- *
//          Tests for API ReceiveMail()

func TestReceiveMailRejectsGET(t *testing.T) {
	testlog.BeginTest()
	defer testlog.EndTest()

	cleanupFiles := setupAndCleanup()
	defer cleanupFiles()

	setupHandler := newSetupHandler(ReceiveMail)

	testServer := createTestServer(setupHandler)
	defer testServer.Close()

	response, err := http.Get(testServer.URL)
	defer func() {
		if err == nil {
			response.Body.Close()
		}
	}()

	t.Run("GETError", func(t *testing.T) {
		assert.NoError(t, err, "GET request should be successful")
	})

	t.Run("statusCode", func(t *testing.T) {
		assert.Equal(t, http.StatusMethodNotAllowed, response.StatusCode,
			"status code should be 405 Method Not Allowed")
	})

	body, readErr := ioutil.ReadAll(response.Body)

	t.Run("bodyReadError", func(t *testing.T) {
		assert.NoError(t, readErr, "reading response body should return no error")
	})

	t.Run("jsonResponse", func(t *testing.T) {
		assert.Equal(t, buildExpectedJSON(structs.JSONMap{
			"status":  http.StatusMethodNotAllowed,
			"message": "METHOD_NOT_ALLOWED (GET)",
		}), body, "response should be JSON with error message METHOD_NOT_ALLOWED")

		testlog.Debug(string(body))
	})
}

func TestReceiveMailAcceptsPOST(t *testing.T) {
	testlog.BeginTest()
	defer testlog.EndTest()

	cleanupFiles := setupAndCleanup()
	defer cleanupFiles()

	setupHandler := newSetupHandler(ReceiveMail)

	testServer := createTestServer(setupHandler)
	defer testServer.Close()

	const validJSON string = `{"from":"admin@example.com","subject":"Subject line","message":"Message line"}`

	response, err := http.Post(testServer.URL, jsonContentTypeTest, createReader(validJSON))
	defer func() {
		if err == nil {
			response.Body.Close()
		}
	}()

	t.Run("POSTError", func(t *testing.T) {
		assert.NoError(t, err, "POST requests should be successful")
	})

	t.Run("statusCode", func(t *testing.T) {
		assert.Equal(t, http.StatusOK, response.StatusCode,
			"status code of POST request should be 200 OK")
	})

	body, readErr := ioutil.ReadAll(response.Body)

	t.Run("bodyReadError", func(t *testing.T) {
		assert.NoError(t, readErr, "reading response body should return no error")
	})

	t.Run("jsonResponse", func(t *testing.T) {
		assert.Equal(t, buildExpectedJSON(structs.JSONMap{
			"status":  http.StatusOK,
			"message": "OK",
		}), body,
			"response should be JSON with status OK")

		testlog.Debug(string(body))
	})
}

func TestReceiveMailInvalidJsonSyntax(t *testing.T) {
	testlog.BeginTest()
	defer testlog.EndTest()

	cleanupFiles := setupAndCleanup()
	defer cleanupFiles()

	setupHandler := newSetupHandler(ReceiveMail)

	testServer := createTestServer(setupHandler)
	defer testServer.Close()

	// JSON string is invalid because terminating '}' is missing
	const invalidJSON string = `{"from":"admin@example.com","subject":"Subject line","message":"Message line"`

	response, err := http.Post(testServer.URL, jsonContentTypeTest, createReader(invalidJSON))
	defer func() {
		if err == nil {
			response.Body.Close()
		}
	}()

	t.Run("POSTError", func(t *testing.T) {
		assert.NoError(t, err, "POST request should be successful")
	})

	t.Run("statusCode", func(t *testing.T) {
		assert.Equal(t, http.StatusBadRequest, response.StatusCode,
			"status code should be 400 Bad Request because of JSON parse error")

		logResponseBody(response)
	})
}

func TestReceiveMailMissingRequiredProperties(t *testing.T) {
	testlog.BeginTest()
	defer testlog.EndTest()

	cleanupFiles := setupAndCleanup()
	defer cleanupFiles()

	setupHandler := newSetupHandler(ReceiveMail)

	testServer := createTestServer(setupHandler)
	defer testServer.Close()

	// Property "to" is not defined by the API (must be "from")
	const missingProperties string = `{"to":"admin@example.com","subject":"Subject Line","message":"Message line"}`

	response, err := http.Post(testServer.URL, jsonContentTypeTest, createReader(missingProperties))
	defer func() {
		if err == nil {
			response.Body.Close()
		}
	}()

	t.Run("POSTError", func(t *testing.T) {
		assert.NoError(t, err, "POST request should be successful")
	})

	t.Run("statusCode", func(t *testing.T) {
		assert.Equal(t, http.StatusBadRequest, response.StatusCode,
			"status code should be 400 Bad Request because of missing property 'from'")

		logResponseBody(response)
	})
}

func TestReceiveMailAdditionalProperties(t *testing.T) {
	testlog.BeginTest()
	defer testlog.EndTest()

	cleanupFiles := setupAndCleanup()
	defer cleanupFiles()

	setupHandler := newSetupHandler(ReceiveMail)

	testServer := createTestServer(setupHandler)
	defer testServer.Close()

	// JSON contains the additional property "to" which is not permitted
	const additionalProperties string = `{"from":"admin@example.com","to":"no-reply@trivial-tickets.com","subject":"Subject Line","message":"Message Line"}`

	response, err := http.Post(testServer.URL, jsonContentTypeTest, createReader(additionalProperties))
	defer func() {
		if err == nil {
			response.Body.Close()
		}
	}()

	t.Run("POSTError", func(t *testing.T) {
		assert.NoError(t, err, "POST request should be successful")
	})

	t.Run("statusCode", func(t *testing.T) {
		assert.Equal(t, http.StatusBadRequest, response.StatusCode,
			"status code should be 400 Bad Request because of additional property 'to'")

		logResponseBody(response)
	})
}

func TestReceiveMailInvalidTypes(t *testing.T) {
	testlog.BeginTest()
	defer testlog.EndTest()

	cleanupFiles := setupAndCleanup()
	defer cleanupFiles()

	setupHandler := newSetupHandler(ReceiveMail)

	testServer := createTestServer(setupHandler)
	defer testServer.Close()

	// This is a valid JSON string, but the type of "from" is invalid
	const invalidTypes string = `{"from":42,"subject":"Subject Line","message":"Message Line"}`

	response, err := http.Post(testServer.URL, jsonContentTypeTest, createReader(invalidTypes))
	defer func() {
		if err == nil {
			response.Body.Close()
		}
	}()

	t.Run("POSTError", func(t *testing.T) {
		assert.NoError(t, err, "POST request should be successful")
	})

	t.Run("statusCode", func(t *testing.T) {
		assert.Equal(t, http.StatusBadRequest, response.StatusCode,
			"status code should be 400 Bad Request because of invalid type of 'from'")

		logResponseBody(response)
	})
}

func TestReceiveMailInvalidEmailAddress(t *testing.T) {
	testlog.BeginTest()
	defer testlog.EndTest()

	cleanupFiles := setupAndCleanup()
	defer cleanupFiles()

	setupHandler := newSetupHandler(ReceiveMail)

	testServer := createTestServer(setupHandler)
	defer testServer.Close()

	// Email address in "from" is invalid because it contains no top-level domain
	const invalidEmailAddress string = `{"from":"invalid@email","subject":"Subject Line","message":"Message Line"}`

	response, err := http.Post(testServer.URL, jsonContentTypeTest, createReader(invalidEmailAddress))
	defer func() {
		if err == nil {
			response.Body.Close()
		}
	}()

	t.Run("POSTError", func(t *testing.T) {
		assert.NoError(t, err, "POST request should be successful")
	})

	t.Run("statusCode", func(t *testing.T) {
		assert.Equal(t, http.StatusBadRequest, response.StatusCode,
			"status code should be 400 Bad Request because of invalid email address")

		logResponseBody(response)
	})
}

func TestReceiveMailCreateTicket(t *testing.T) {
	testlog.BeginTest()
	defer testlog.EndTest()

	cleanupFiles := setupAndCleanup()
	defer cleanupFiles()

	setupHandler := newSetupHandler(ReceiveMail)

	testServer := createTestServer(setupHandler)
	defer testServer.Close()

	const createTicket string = `{"from":"customer@mail.com","subject":"Issue with computer","message":"My computer is broken!"}`

	response, err := http.Post(testServer.URL, jsonContentTypeTest, createReader(createTicket))
	defer func() {
		if err == nil {
			response.Body.Close()
		}
	}()

	t.Run("POSTError", func(t *testing.T) {
		assert.NoError(t, err, "POST request should be successful")
	})

	t.Run("statusCode", func(t *testing.T) {
		assert.Equal(t, http.StatusOK, response.StatusCode,
			"response status should be 200 OK because the request is valid")
	})

	body, readErr := ioutil.ReadAll(response.Body)

	t.Run("bodyReadError", func(t *testing.T) {
		assert.NoError(t, readErr, "reading response body should not return error")
	})

	t.Run("jsonResponse", func(t *testing.T) {
		assert.Equal(t, buildExpectedJSON(structs.JSONMap{
			"status":  http.StatusOK,
			"message": http.StatusText(http.StatusOK),
		}), body, "response body should be JSON with status OK")

		testlog.Debug(string(body))
	})

	// Usually, the t.Run() function should block until the
	// subtest finished. However, there were issues that the
	// deferred cleanup function was executed before the
	// actual subtest and this made the tests fail. This might
	// be a bug as pointed out in this issue:
	//   https://github.com/golang/go/issues/17791
	//
	// Solution:
	// Create a new buffered cleanup channel with a buffer
	// size of 1 to hold a bool value. This is used to wait
	// for the following subtests that require reading the
	// ticket or mail directories with the test files to
	// complete before the deferred function cleanupMails()
	// removes the test directories. If this happens the tests
	// fail with unpredictable errors because the directory
	// cannot be read because it does not exist or there
	// are more mails in the directory than it should be.
	cleanup := make(chan bool, 1)
	testlog.Debug("Creating cleanup channel to wait for subtests to complete")

	t.Run("verifyTicketCreated", func(t *testing.T) {
		t.Run("ticketFile", func(t *testing.T) {
			dirContents, readErr := ioutil.ReadDir(globals.ServerConfig.Tickets)
			assert.NoError(t, readErr, fmt.Sprintf("directory '%s' should exist and be readable within test", globals.ServerConfig.Tickets))

			displayDirectoryContents(globals.ServerConfig.Tickets, dirContents)
			assert.Equal(t, 1, len(dirContents), fmt.Sprintf("directory '%s' should contain exactly one ticket", globals.ServerConfig.Tickets))
		})

		t.Run("mailFile", func(t *testing.T) {
			dirContents, readErr := ioutil.ReadDir(globals.ServerConfig.Mails)
			assert.NoError(t, readErr, fmt.Sprintf("directory '%s' should exist and be readable within test", globals.ServerConfig.Mails))

			displayDirectoryContents(globals.ServerConfig.Mails, dirContents)
			assert.Equal(t, 1, len(dirContents), fmt.Sprintf("directory '%s' should contain exactly one mail (ticket creation)", globals.ServerConfig.Mails))
		})

		// The subtests completed at this point so that the
		// cleanup channel can be released and closed
		testlog.Debug("Subtests finished: writing to cleanup channel and closing it")
		cleanup <- true
		close(cleanup)
	})

	testlog.Debug("Waiting for subtests to complete")
	<-cleanup
	testlog.Debug("Done: Cleaning test directories")
}

func TestReceiveMailCreateAnswer(t *testing.T) {
	testlog.BeginTest()
	defer testlog.EndTest()

	cleanupFiles := setupAndCleanup()
	defer cleanupFiles()

	setupHandler := newSetupHandler(ReceiveMail)

	testServer := createTestServer(setupHandler)
	defer testServer.Close()

	// Create new test ticket in order to submit an answer to it using the API
	testTicket := ticket.CreateTicket("customer@mail.com", "Issue with Computer", "My computer is broken")
	globals.Tickets[testTicket.ID] = testTicket
	writeErr := filehandler.WriteTicketFile(globals.ServerConfig.Tickets, &testTicket)
	if writeErr != nil {
		testlog.Debug("ERROR:", writeErr)
	}

	t.Run("writeError", func(t *testing.T) {
		assert.NoError(t, writeErr, "writing the ticket file caused an error")
	})

	answerSubject := fmt.Sprintf(`[Ticket \"%s\"] Issue with Computer`, testTicket.ID)

	createAnswerJSON := fmt.Sprintf(`{"from":"customer@mail.com","subject":"%s","message":"My computer is broken!"}`, answerSubject)

	response, err := http.Post(testServer.URL, jsonContentTypeTest, createReader(createAnswerJSON))
	defer func() {
		if err == nil {
			response.Body.Close()
		}
	}()

	t.Run("statusCode", func(t *testing.T) {
		assert.Equal(t, http.StatusOK, response.StatusCode,
			"response status should be 200 OK because the request is valid")
	})

	body, readErr := ioutil.ReadAll(response.Body)

	t.Run("bodyReadError", func(t *testing.T) {
		assert.NoError(t, readErr, "reading response body should not return error")
	})

	t.Run("jsonResponse", func(t *testing.T) {
		assert.Equal(t, buildExpectedJSON(structs.JSONMap{
			"status":  http.StatusOK,
			"message": http.StatusText(http.StatusOK),
		}), body, "response body should be JSON with status OK")

		testlog.Debug(string(body))
	})

	// Usually, the t.Run() function should block until the
	// subtest finished. However, there were issues that the
	// deferred cleanup function was executed before the
	// actual subtest and this made the tests fail. This might
	// be a bug as pointed out in this issue:
	//   https://github.com/golang/go/issues/17791
	//
	// Solution:
	// Create a new buffered cleanup channel with a buffer
	// size of 1 to hold a bool value. This is used to wait
	// for the following subtests that require reading the
	// ticket or mail directories with the test files to
	// complete before the deferred function cleanupMails()
	// removes the test directories. If this happens the tests
	// fail with unpredictable errors because the directory
	// cannot be read because it does not exist or there
	// are more mails in the directory than it should be.
	cleanup := make(chan bool, 1)
	testlog.Debug("Creating cleanup channel to wait for subtests to complete")

	t.Run("verifyAnswerCreated", func(t *testing.T) {
		t.Run("ticketFile", func(t *testing.T) {
			dirContents, readErr := ioutil.ReadDir(globals.ServerConfig.Tickets)
			assert.NoError(t, readErr, fmt.Sprintf("directory '%s' should exist and be readable within test", globals.ServerConfig.Tickets))

			displayDirectoryContents(globals.ServerConfig.Mails, dirContents)
			assert.Equal(t, 1, len(dirContents), fmt.Sprintf("directory '%s' should contain exactly one ticket", globals.ServerConfig.Tickets))
		})

		t.Run("mailFile", func(t *testing.T) {
			dirContents, readErr := ioutil.ReadDir(globals.ServerConfig.Mails)
			assert.NoError(t, readErr, fmt.Sprintf("directory '%s' should exist and be readable within test", globals.ServerConfig.Mails))

			displayDirectoryContents(globals.ServerConfig.Mails, dirContents)
			assert.Equal(t, 1, len(dirContents), fmt.Sprintf("directory '%s' should contain exactly one mail (answer creation)", globals.ServerConfig.Mails))
		})

		// The subtests completed at this point so that the
		// cleanup channel can be released and closed
		testlog.Debug("Subtests finished: writing to cleanup channel and closing it")
		cleanup <- true
		close(cleanup)
	})

	testlog.Debug("Waiting for subtests to complete")
	<-cleanup
	testlog.Debug("Done: Cleaning test directories")
}

func TestReceiveMailCreateAnswerInvalidTicketId(t *testing.T) {
	testlog.BeginTest()
	defer testlog.EndTest()

	cleanupFiles := setupAndCleanup()
	defer cleanupFiles()

	setupHandler := newSetupHandler(ReceiveMail)

	testServer := createTestServer(setupHandler)
	defer testServer.Close()

	// Create new test ticket in order to submit an answer to it using the API
	testTicket := ticket.CreateTicket("customer@mail.com", "Issue with Computer", "My computer is broken")
	globals.Tickets[testTicket.ID] = testTicket
	writeErr := filehandler.WriteTicketFile(globals.ServerConfig.Tickets, &testTicket)
	if writeErr != nil {
		testlog.Debug("ERROR:", writeErr)
	}

	t.Run("writeError", func(t *testing.T) {
		assert.NoError(t, writeErr, "writing the ticket file caused an error")
	})

	// Manipulate the ticket id so that it gets invalid
	manipulatedID := testTicket.ID + "x"
	answerSubject := fmt.Sprintf(`[Ticket \"%s\"] Issue with Computer`, manipulatedID)

	createAnswerJSON := fmt.Sprintf(`{"from":"customer@mail.com","subject":"%s","message":"My computer is broken!"}`, answerSubject)

	response, err := http.Post(testServer.URL, jsonContentTypeTest, createReader(createAnswerJSON))
	defer func() {
		if err == nil {
			response.Body.Close()
		}
	}()

	t.Run("POSTError", func(t *testing.T) {
		assert.NoError(t, err, "POST request should be successful")
	})

	t.Run("statusCode", func(t *testing.T) {
		assert.Equal(t, http.StatusOK, response.StatusCode, "response status code should be 200 OK")
	})

	// Usually, the t.Run() function should block until the
	// subtest finished. However, there were issues that the
	// deferred cleanup function was executed before the
	// actual subtest and this made the tests fail. This might
	// be a bug as pointed out in this issue:
	//   https://github.com/golang/go/issues/17791
	//
	// Solution:
	// Create a new buffered cleanup channel with a buffer
	// size of 1 to hold a bool value. This is used to wait
	// for the following subtests that require reading the
	// ticket or mail directories with the test files to
	// complete before the deferred function cleanupMails()
	// removes the test directories. If this happens the tests
	// fail with unpredictable errors because the directory
	// cannot be read because it does not exist or there
	// are more mails in the directory than it should be.
	cleanup := make(chan bool, 1)
	testlog.Debug("Creating cleanup channel to wait for subtests to complete")

	// The invalid ticket id causes a new ticket to be created,
	// therefore test for two files in the tickets directory
	t.Run("verifyAnswerCreated", func(t *testing.T) {
		t.Run("ticketFile", func(t *testing.T) {
			dirContents, readErr := ioutil.ReadDir(globals.ServerConfig.Tickets)
			assert.NoError(t, readErr, fmt.Sprintf("directory '%s' should exist and be readable within test", globals.ServerConfig.Tickets))

			displayDirectoryContents(globals.ServerConfig.Tickets, dirContents)
			assert.Equal(t, 2, len(dirContents), fmt.Sprintf("directory '%s' should contain two tickets", globals.ServerConfig.Tickets))
		})

		t.Run("mailFile", func(t *testing.T) {
			dirContents, readErr := ioutil.ReadDir(globals.ServerConfig.Mails)
			assert.NoError(t, readErr, fmt.Sprintf("directory '%s' should exist and be readable within test", globals.ServerConfig.Mails))

			displayDirectoryContents(globals.ServerConfig.Mails, dirContents)
			assert.Equal(t, 1, len(dirContents), fmt.Sprintf("directory '%s' should contain exactly one mail (answer creation)", globals.ServerConfig.Mails))
		})

		// The subtests completed at this point so that the
		// cleanup channel can be released and closed
		testlog.Debug("Subtests finished: writing to cleanup channel and closing it")
		cleanup <- true
		close(cleanup)
	})

	testlog.Debug("Waiting for subtests to complete")
	<-cleanup
	testlog.Debug("Done: Cleaning test directories")
}

func TestReceiveMailCreateAnswerClosedTicket(t *testing.T) {
	testlog.BeginTest()
	defer testlog.EndTest()

	cleanupFiles := setupAndCleanup()
	defer cleanupFiles()

	setupHandler := newSetupHandler(ReceiveMail)

	testServer := createTestServer(setupHandler)
	defer testServer.Close()

	// Create new test ticket in order to submit an answer to it using the API
	testTicket := ticket.CreateTicket("customer@mail.com", "Issue with Computer", "My computer is broken")

	// Set the ticket status to closed
	testTicket.Status = structs.StatusClosed

	globals.Tickets[testTicket.ID] = testTicket
	writeErr := filehandler.WriteTicketFile(globals.ServerConfig.Tickets, &testTicket)
	if writeErr != nil {
		testlog.Debug("ERROR:", writeErr)
	}

	t.Run("writeError", func(t *testing.T) {
		assert.NoError(t, writeErr, "writing the ticket file caused an error")
	})

	answerSubject := fmt.Sprintf(`[Ticket \"%s\"] Issue with Computer`, testTicket.ID)

	answerJSON := fmt.Sprintf(`{"from":"%s","subject":"%s","message":"Answer to an existing ticket"}`,
		testTicket.Customer, answerSubject)

	response, err := http.Post(testServer.URL, jsonContentTypeTest, createReader(answerJSON))
	defer func() {
		if err == nil {
			response.Body.Close()
		}
	}()

	t.Run("POSTError", func(t *testing.T) {
		assert.NoError(t, err, "POST request should be successful")
	})

	t.Run("statusCode", func(t *testing.T) {
		assert.Equal(t, http.StatusOK, response.StatusCode, "response status code should be 200 OK")
	})

	updatedTicket := globals.Tickets[testTicket.ID]

	t.Run("statusChangedToOpen", func(t *testing.T) {
		assert.Equal(t, structs.StatusOpen, updatedTicket.Status, "ticket status should be reset to StatusOpen")
	})
}

// * ------------------------------------------- *
//          Tests for helper functions
//               of ReceiveMail()

func TestMatchAnswerSubject(t *testing.T) {
	testlog.BeginTest()
	defer testlog.EndTest()

	ticketID := "100"
	subject := fmt.Sprintf(`[Ticket "%s"] Ticket subject`, ticketID)

	extractedID, matching := matchAnswerSubject(subject)

	t.Run("subjectMatching", func(t *testing.T) {
		assert.True(t, matching, "test subject should match the answering schema")
	})

	t.Run("equalTicketId", func(t *testing.T) {
		assert.Equal(t, ticketID, extractedID, "extracted id should be identical to test id")
	})
}

func TestValidEmailAddress(t *testing.T) {
	testlog.BeginTest()
	defer testlog.EndTest()

	email := "admin@example.com"

	assert.True(t, validEmailAddress(email), "email should be valid")
}

func TestCheckRequiredPropertiesSet(t *testing.T) {
	testlog.BeginTest()
	defer testlog.EndTest()

	t.Run("propertiesSet", func(t *testing.T) {
		requiredProperties := structs.JSONMap{
			"from":    "admin@example.com",
			"subject": "Subject Line",
			"message": "Message Line",
		}

		err := checkRequiredPropertiesSet(requiredProperties)

		assert.NoError(t, err, "all required properties are set so there should be no error")
	})

	t.Run("propertiesMissing", func(t *testing.T) {
		requiredPropertiesMissing := structs.JSONMap{
			"subject": "Subject Line",
			"message": "Message Line",
		}

		err := checkRequiredPropertiesSet(requiredPropertiesMissing)

		assert.Error(t, err, "missing required properties error")
	})
}

func TestCheckAdditionalPropertiesSet(t *testing.T) {
	testlog.BeginTest()
	defer testlog.EndTest()

	t.Run("noAdditionalProperties", func(t *testing.T) {
		noAdditionalProperties := structs.JSONMap{
			"from":    "admin@example.com",
			"subject": "Subject Line",
			"message": "Message Line",
		}

		err := checkAdditionalPropertiesSet(noAdditionalProperties)

		assert.NoError(t, err, "no additional properties set, therefore no error")
	})

	t.Run("additionalProperties", func(t *testing.T) {
		additionalProperties := structs.JSONMap{
			"from":    "admin@example.com",
			"to":      "no-reply@trivial-tickets.com",
			"subject": "Subject Line",
			"message": "Message Line",
		}

		err := checkAdditionalPropertiesSet(additionalProperties)

		assert.Error(t, err, "error because of additional properties set")
	})
}

func TestCheckCorrectPropertyTypes(t *testing.T) {
	testlog.BeginTest()
	defer testlog.EndTest()

	t.Run("correctTypes", func(t *testing.T) {
		correctTypes := structs.JSONMap{
			"from":    "admin@example.com",
			"subject": "Subject Line",
			"message": "Message Line",
		}

		err := checkCorrectPropertyTypes(correctTypes)

		assert.NoError(t, err, "all types are correct so there should be no error")
	})

	t.Run("invalidTypes", func(t *testing.T) {
		inCorrectTypes := structs.JSONMap{
			"from":    100,
			"subject": "Subject Line",
			"message": "Message Line",
		}

		err := checkCorrectPropertyTypes(inCorrectTypes)

		assert.Error(t, err, "'from' type is invalid so there should be an error")
	})
}

func TestCheckCorrectPropertyTypesPropertyNotGiven(t *testing.T) {
	testlog.BeginTest()
	defer testlog.EndTest()

	missingPropertyMap := structs.JSONMap{
		"from":    "admin@example.com",
		"subject": "message not given here",
	}

	typeErr := checkCorrectPropertyTypes(missingPropertyMap)

	t.Run("nonNilError", func(t *testing.T) {
		assert.Error(t, typeErr, "function should also check if required properties are not "+
			"given and return an error in this case")
	})

	t.Run("isPropertyNotDefinedError", func(t *testing.T) {
		assert.IsType(t, propertyNotDefinedError{}, typeErr, "error should be of type propertyNotDefinedError")
	})

	t.Run("errorMessageContainsMissingPropertyName", func(t *testing.T) {
		assert.Contains(t, typeErr.Error(), "message", "message property is missing and should "+
			"be contained in error message")
	})
}

func TestWriteJsonProperty(t *testing.T) {
	testlog.BeginTest()
	defer testlog.EndTest()

	t.Run("intValue", func(t *testing.T) {
		expectedStatus := `"status":200`
		actualStatus := writeJSONProperty("status", 200)

		assert.Equal(t, expectedStatus, actualStatus, "actual json status should match expected json")
	})

	t.Run("stringValue", func(t *testing.T) {
		expectedEmail := `"email":"admin@example.com"`
		actualEmail := writeJSONProperty("email", "admin@example.com")

		assert.Equal(t, expectedEmail, actualEmail, "actual json email should match expected json")
	})
}
