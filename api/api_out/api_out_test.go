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

// Package api_out implements a web interface for outgoing mails
// to be fetched and verified to be sent
package api_out

import (
	"encoding/json"
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
	"github.com/mortenterhart/trivial-tickets/logger/testlogger"
	"github.com/mortenterhart/trivial-tickets/mail_events"
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
 * Package api_out [tests]
 * Web API for outgoing mails to be fetched and verified to be sent
 */

// jsonContentTypeTest is the content type used for
// the POST request to handlers in the following tests
const jsonContentTypeTest = "application/json"

// testServerConfig returns a test server configuration.
func testServerConfig() structs.ServerConfig {
	return structs.ServerConfig{
		Port:    8443,
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

// initializeConfig is run before all tests and initializes the
// global server and logging configuration.
func initializeConfig() {
	config := testServerConfig()
	globals.ServerConfig = &config

	logConfig := testLogConfig()
	globals.LogConfig = &logConfig
}

// TestMain is started to run the tests and initializes the
// configuration before running the tests. The tests' exit
// status is returned as the overall exit status.
func TestMain(m *testing.M) {
	initializeConfig()

	os.Exit(m.Run())
}

// mockTicket creates a mocked ticket struct
func mockTicket() structs.Ticket {
	return ticket.CreateTicket("customer@mail.com", "Test subject", "Test message")
}

// createTestServer creates a server for testing purposes with
// the given handler registered
func createTestServer(handlerFunc http.HandlerFunc) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(handlerFunc))
}

// createReader creates a reader from a string suitable to use as
// request body
func createReader(data string) io.Reader {
	return strings.NewReader(data)
}

// buildExpectedJSON creates a json string from a json properties
// map that can be used to compare the actual json response of the
// APIs with the expected json
func buildExpectedJSON(properties structs.JSONMap) []byte {
	expected := jsontools.MapToJSON(properties)
	return append(expected, '\n')
}

// cleanupMails is a teardown function which cleans all created mails
// from the global hash map and the file system.
func cleanupMails() {
	// Remove each entry from the mail hash map
	for id := range globals.Mails {
		testlogger.Debugf("Deferred: Removing mail '%s' from mail hash map", id)
		delete(globals.Mails, id)
	}

	// Delete the mail directory for temporary mails if it exists
	if filehandler.DirectoryExists(globals.ServerConfig.Mails) {
		testlogger.Debug("Deferred: Removing test mail directory", globals.ServerConfig.Mails)
		if removeErr := os.RemoveAll(globals.ServerConfig.Mails); removeErr != nil {
			testlogger.Debug("ERROR: could not remove test mail directory:", removeErr)
		}
	}
}

// displayDirectoryContents is a helper function of tests that require
// certain files in a directory to exist. It prints all filenames in
// a directory to the console in alphabetical order.
func displayDirectoryContents(dirname string, contents []os.FileInfo) {
	testlogger.Debugf("Directory contents of '%s' showing %d file(s):", dirname, len(contents))
	for index, file := range contents {
		testlogger.Debugf("%d: %s", index, file.Name())
	}
}

// * ------------------------------------------- *
//          Tests for SendMail()

func TestSendMail(t *testing.T) {
	testlogger.BeginTest()
	defer testlogger.EndTest()

	defer cleanupMails()

	testTicket := mockTicket()

	SendMail(mail_events.NewTicket, testTicket)

	t.Run("storedInMailMap", func(t *testing.T) {
		assert.Equal(t, 1, len(globals.Mails), "sent mail should be stored in the global mail storage")
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
	testlogger.Debug("Creating cleanup channel to wait for subtests to complete")

	t.Run("savedInFile", func(t *testing.T) {
		dirContents, readErr := ioutil.ReadDir(globals.ServerConfig.Mails)

		t.Run("readError", func(t *testing.T) {
			assert.NoError(t, readErr, "reading mail directory should not return error")
		})

		displayDirectoryContents(globals.ServerConfig.Mails, dirContents)

		t.Run("mailFile", func(t *testing.T) {
			assert.Equal(t, 1, len(dirContents), "mail directory should contain exactly one mail")
		})

		// The subtests completed at this point so that the
		// cleanup channel can be released and closed
		testlogger.Debug("Subtests finished: writing to cleanup channel and closing it")
		cleanup <- true
		close(cleanup)
	})

	testlogger.Debug("Waiting for subtests to complete")
	<-cleanup
	testlogger.Debug("Done: Cleaning test directories")
}

// * ------------------------------------------- *
//          Tests for API FetchMails()

func TestFetchMails(t *testing.T) {
	testlogger.BeginTest()
	defer testlogger.EndTest()

	defer cleanupMails()

	testServer := createTestServer(FetchMails)
	defer testServer.Close()

	t.Run("invalidPOSTRequest", func(t *testing.T) {
		const jsonRequest = `{"id":"test-id"}`

		response, err := http.Post(testServer.URL, jsonContentTypeTest, createReader(jsonRequest))
		defer func() {
			if err == nil {
				response.Body.Close()
			}
		}()

		t.Run("POSTError", func(t *testing.T) {
			assert.NoError(t, err, "POST request should execute successfully")
		})

		t.Run("statusCode", func(t *testing.T) {
			assert.Equal(t, http.StatusMethodNotAllowed, response.StatusCode, "POST request should be declined, therefore 405 Method Not Allowed")
		})
	})

	t.Run("validGETRequest", func(t *testing.T) {
		t.Run("withoutSavedMails", func(t *testing.T) {
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
				assert.Equal(t, http.StatusOK, response.StatusCode, "GET request should be accepted and responded with 200 OK")
			})

			body, readErr := ioutil.ReadAll(response.Body)

			t.Run("bodyReadError", func(t *testing.T) {
				assert.NoError(t, readErr, "reading response body should not return error")
			})

			t.Run("jsonResponse", func(t *testing.T) {
				assert.Equal(t, "{}\n", string(body), "response should contain an empty JSON object since no mails are saved")
			})
		})

		t.Run("withSavedMails", func(t *testing.T) {
			// Create test mail to be fetched by the API request
			testTicket := mockTicket()
			SendMail(mail_events.NewTicket, testTicket)

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
				assert.Equal(t, http.StatusOK, response.StatusCode, "response status should be 200 OK")
			})

			body, readErr := ioutil.ReadAll(response.Body)

			t.Run("bodyReadErr", func(t *testing.T) {
				assert.NoError(t, readErr, "reading response body should not return error")
			})

			t.Run("jsonResponse", func(t *testing.T) {
				expectedJSON, decodeErr := json.MarshalIndent(&globals.Mails, "", "    ")

				assert.NoError(t, decodeErr, "decoding test mail to JSON should not return an error")
				assert.Equal(t, append(expectedJSON, '\n'), body, "response should contain JSON representation of mail mapped to its id")
			})
		})
	})
}

// * ------------------------------------------- *
//          Tests for API VerifyMailSent()

func TestVerifyMailSentRejectsGet(t *testing.T) {
	testlogger.BeginTest()
	defer testlogger.EndTest()

	defer cleanupMails()

	testServer := createTestServer(VerifyMailSent)
	defer testServer.Close()

	response, err := http.Get(testServer.URL)
	defer func() {
		if err == nil {
			response.Body.Close()
		}
	}()

	t.Run("GETError", func(t *testing.T) {
		assert.NoError(t, err, "GET request should execute successfully")
	})

	t.Run("statusCode", func(t *testing.T) {
		assert.Equal(t, http.StatusMethodNotAllowed, response.StatusCode, "GET requests should be declined with 405 Method Not Allowed")
	})
}

func TestVerifyMailSentInvalidJSON(t *testing.T) {
	testlogger.BeginTest()
	defer testlogger.EndTest()

	defer cleanupMails()

	testServer := createTestServer(VerifyMailSent)
	defer testServer.Close()

	// JSON is invalid because the key "id" is not surrounded by quotes
	const invalidJSON = `{id:"mail-id"}`

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
		assert.Equal(t, http.StatusBadRequest, response.StatusCode, "request should return 400 Bad Request due to invalid JSON")
	})
}

func TestVerifyMailSentMissingRequiredProperties(t *testing.T) {
	testlogger.BeginTest()
	defer testlogger.EndTest()

	defer cleanupMails()

	testServer := createTestServer(VerifyMailSent)
	defer testServer.Close()

	// Required property "id" is missing
	const missingProperties = `{"email":"admin@example.com"}`

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
		assert.Equal(t, http.StatusBadRequest, response.StatusCode, "request should return 400 Bad Request due to missing properties")
	})
}

func TestVerifyMailSentAdditionalProperties(t *testing.T) {
	testlogger.BeginTest()
	defer testlogger.EndTest()

	defer cleanupMails()

	testServer := createTestServer(VerifyMailSent)
	defer testServer.Close()

	// Additional property "email" is invalid
	const additionalProperties = `{"id":"mail-id","email":"admin@example.com"}`

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
		assert.Equal(t, http.StatusBadRequest, response.StatusCode, "request should return 400 Bad Request due to additional properties")
	})
}

func TestVerifyMailSentInvalidPropertyTypes(t *testing.T) {
	testlogger.BeginTest()
	defer testlogger.EndTest()

	defer cleanupMails()

	testServer := createTestServer(VerifyMailSent)
	defer testServer.Close()

	// "id" property has invalid type for API
	const propertyTypes = `{"id":200}`

	response, err := http.Post(testServer.URL, jsonContentTypeTest, createReader(propertyTypes))
	defer func() {
		if err == nil {
			response.Body.Close()
		}
	}()

	t.Run("POSTError", func(t *testing.T) {
		assert.NoError(t, err, "POST request should be successful")
	})

	t.Run("statusCode", func(t *testing.T) {
		assert.Equal(t, http.StatusBadRequest, response.StatusCode, "request should return 400 Bad Request due to invalid types")
	})
}

func TestVerifyMailSentExistingMail(t *testing.T) {
	testlogger.BeginTest()
	defer testlogger.EndTest()

	defer cleanupMails()

	testServer := createTestServer(VerifyMailSent)
	defer testServer.Close()

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
	testlogger.Debug("Creating cleanup channel to wait for subtests to complete")

	SendMail(mail_events.NewTicket, mockTicket())

	// Since there is only one mail created at this point,
	// retrieve the mail id of the just created mail
	var mailID string
	for id := range globals.Mails {
		mailID = id
	}

	verifyJSON := fmt.Sprintf(`{"id":"%s"}`, mailID)

	response, err := http.Post(testServer.URL, jsonContentTypeTest, createReader(verifyJSON))
	defer func() {
		if err == nil {
			response.Body.Close()
		}
	}()

	t.Run("POSTError", func(t *testing.T) {
		assert.NoError(t, err, "POST request should be successful")
	})

	t.Run("statusCode", func(t *testing.T) {
		assert.Equal(t, http.StatusOK, response.StatusCode, "POST request with valid JSON properties and valid mail id should return 200 OK")
	})

	body, readErr := ioutil.ReadAll(response.Body)

	t.Run("bodyReadError", func(t *testing.T) {
		assert.NoError(t, readErr, "reading response body should not return an error")
	})

	t.Run("jsonResponse", func(t *testing.T) {
		expectedJSON := buildExpectedJSON(structs.JSONMap{
			"verified": true,
			"status":   http.StatusOK,
			"message":  fmt.Sprintf("mail '%s' was successfully sent and deleted from server cache", mailID),
		})

		assert.Equal(t, expectedJSON, body, "mail sending should be verified with correct JSON response")

		// The subtests completed at this point so that the
		// cleanup channel can be released and closed
		testlogger.Debug("Subtests finished: writing to cleanup channel and closing it")
		cleanup <- true
		close(cleanup)
	})

	testlogger.Debug("Waiting for subtests to complete")
	<-cleanup
	testlogger.Debug("Done: Cleaning test directories")
}

func TestVerifyMailSentNotExistingMail(t *testing.T) {
	testlogger.BeginTest()
	defer testlogger.EndTest()

	defer cleanupMails()

	testServer := createTestServer(VerifyMailSent)
	defer testServer.Close()

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
	testlogger.Debug("Creating cleanup channel to wait for subtests to complete")

	mailID := "notExisting"

	verifyJSON := fmt.Sprintf(`{"id":"%s"}`, mailID)

	response, err := http.Post(testServer.URL, jsonContentTypeTest, createReader(verifyJSON))
	defer func() {
		if err == nil {
			response.Body.Close()
		}
	}()

	t.Run("POSTError", func(t *testing.T) {
		assert.NoError(t, err, "POST request should be successful")
	})

	t.Run("statusCode", func(t *testing.T) {
		assert.Equal(t, http.StatusOK, response.StatusCode, "POST request with valid JSON properties and invalid mail id should return 200 OK")
	})

	body, readErr := ioutil.ReadAll(response.Body)

	t.Run("bodyReadError", func(t *testing.T) {
		assert.NoError(t, readErr, "reading response body should not return an error")
	})

	t.Run("jsonResponse", func(t *testing.T) {
		expectedJSON := buildExpectedJSON(structs.JSONMap{
			"verified": false,
			"status":   http.StatusOK,
			"message":  fmt.Sprintf("mail '%s' does not exist or has already been deleted", mailID),
		})

		assert.Equal(t, expectedJSON, body, "mail sending should be verified with correct JSON response")

		// The subtests completed at this point so that the
		// cleanup channel can be released and closed
		testlogger.Debug("Subtests finished: writing to cleanup channel and closing it")
		cleanup <- true
		close(cleanup)
	})

	testlogger.Debug("Waiting for subtests to complete")
	<-cleanup
	testlogger.Debug("Done: Cleaning test directories")
}

func TestVerifyMailSentExistingMailInCacheAndNotExistingMailFile(t *testing.T) {
	testlogger.BeginTest()
	defer testlogger.EndTest()

	defer cleanupMails()

	testServer := createTestServer(VerifyMailSent)
	defer testServer.Close()

	// Add a new mail to the mail hash map, but
	// do not write the corresponding mail file
	mailID := "mail-id"
	globals.Mails[mailID] = structs.Mail{
		ID: mailID,
	}

	verifyJSON := fmt.Sprintf(`{"id":"%s"}`, mailID)

	response, err := http.Post(testServer.URL, jsonContentTypeTest, createReader(verifyJSON))
	defer func() {
		if err == nil {
			response.Body.Close()
		}
	}()

	t.Run("POSTError", func(t *testing.T) {
		assert.NoError(t, err, "POST request should be successful")
	})

	t.Run("statusCode", func(t *testing.T) {
		assert.Equal(t, http.StatusInternalServerError, response.StatusCode, "response status should be "+
			"500 Internal Server Error because the mail file could not be deleted")
	})
}
