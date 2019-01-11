package api_in

import (
	"fmt"
	"github.com/mortenterhart/trivial-tickets/ticket"
	"github.com/mortenterhart/trivial-tickets/util/filehandler"
	"github.com/mortenterhart/trivial-tickets/util/jsontools"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/mortenterhart/trivial-tickets/globals"
	"github.com/mortenterhart/trivial-tickets/structs"
	"github.com/stretchr/testify/assert"
)

/*
*
* Matrikelnummern
* 3040018
* 6694964
* 3478222
 */

const jsonContentType = "application/json"

type serverSetupHandler struct {
	callUnderlying http.HandlerFunc
}

func newSetupHandler(wrappedHandler http.HandlerFunc) serverSetupHandler {
	return serverSetupHandler{wrappedHandler}
}

func (handler serverSetupHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	config := testServerConfig()
	globals.ServerConfig = &config

	handler.callUnderlying(writer, request)
}

type teardownFunc func()

func setupAndTeardown() teardownFunc {
	config := testServerConfig()
	globals.ServerConfig = &config

	return func() {
		cleanupTestFiles(config)
	}
}

func testServerConfig() structs.Config {
	return structs.Config{
		Port:    8443,
		Tickets: "../../files/testtickets",
		Users:   "../../files/users/users.json",
		Mails:   "../../files/testmails",
		Cert:    "../../ssl/server.cert",
		Key:     "../../ssl/server.key",
		Web:     "../../www",
	}
}

func cleanupTestFiles(config structs.Config) {
	if filehandler.FileExists(config.Tickets) {
		if removeErr := os.RemoveAll(globals.ServerConfig.Tickets); removeErr != nil {
			log.Println("test error: cannot remove test tickets:", removeErr)
		}
	}

	if filehandler.FileExists(config.Mails) {
		if removeErr := os.RemoveAll(globals.ServerConfig.Mails); removeErr != nil {
			log.Println("test error: cannot remove test mails:", removeErr)
		}
	}
}

func createTestServer(handler http.Handler) *httptest.Server {
	return httptest.NewServer(handler)
}

func createReader(str string) io.Reader {
	return strings.NewReader(str)
}

func buildExpectedJson(properties structs.JsonMap) []byte {
	expected, decodeErr := jsontools.MapToJson(properties)
	if decodeErr != nil {
		log.Println("error while decoding expected JSON string:", decodeErr)
		return nil
	}

	return append(expected, '\n')
}

func logResponseBody(t *testing.T, response *http.Response) {
	body, readErr := ioutil.ReadAll(response.Body)
	if readErr != nil {
		log.Println("response could not be read:", readErr)
		return
	}

	t.Log(string(body))
}

func TestReceiveMailRejectsGET(t *testing.T) {
	teardown := setupAndTeardown()
	defer teardown()

	setupHandler := newSetupHandler(ReceiveMail)

	testServer := createTestServer(setupHandler)
	defer testServer.Close()

	response, err := http.Get(testServer.URL)

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
		assert.Equal(t, buildExpectedJson(structs.JsonMap{
			"status":  http.StatusMethodNotAllowed,
			"message": "METHOD_NOT_ALLOWED (GET)",
		}), body, "response should be JSON with error message METHOD_NOT_ALLOWED")

		t.Log(string(body))
	})
}

func TestReceiveMailAcceptsPOST(t *testing.T) {
	teardown := setupAndTeardown()
	defer teardown()

	setupHandler := newSetupHandler(ReceiveMail)

	testServer := createTestServer(setupHandler)
	defer testServer.Close()

	const validJson = `{"from":"admin@example.com","subject":"Subject line","message":"Message line"}`

	response, err := http.Post(testServer.URL, jsonContentType, createReader(validJson))

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
		assert.Equal(t, buildExpectedJson(structs.JsonMap{
			"status":  http.StatusOK,
			"message": "OK",
		}), body,
			"response should be JSON with status OK")

		t.Log(string(body))
	})
}

func TestReceiveMailInvalidJSONSyntax(t *testing.T) {
	teardown := setupAndTeardown()
	defer teardown()

	setupHandler := newSetupHandler(ReceiveMail)

	testServer := createTestServer(setupHandler)
	defer testServer.Close()

	// JSON string is invalid because terminating '}' is missing
	const invalidJson = `{"from":"admin@example.com","subject":"Subject line","message":"Message line"`

	response, err := http.Post(testServer.URL, jsonContentType, createReader(invalidJson))

	t.Run("POSTError", func(t *testing.T) {
		assert.NoError(t, err, "POST request should be successful")
	})

	t.Run("statusCode", func(t *testing.T) {
		assert.Equal(t, http.StatusBadRequest, response.StatusCode,
			"status code should be 400 Bad Request because of JSON parse error")

		logResponseBody(t, response)
	})
}

func TestReceiveMailMissingRequiredProperties(t *testing.T) {
	teardown := setupAndTeardown()
	defer teardown()

	setupHandler := newSetupHandler(ReceiveMail)

	testServer := createTestServer(setupHandler)
	defer testServer.Close()

	// Property "to" is not defined by the API (must be "from")
	const missingProperties = `{"to":"admin@example.com","subject":"Subject Line","message":"Message line"}`

	response, err := http.Post(testServer.URL, jsonContentType, createReader(missingProperties))

	t.Run("POSTError", func(t *testing.T) {
		assert.NoError(t, err, "POST request should be successful")
	})

	t.Run("statusCode", func(t *testing.T) {
		assert.Equal(t, http.StatusBadRequest, response.StatusCode,
			"status code should be 400 Bad Request because of missing property 'from'")

		logResponseBody(t, response)
	})
}

func TestReceiveMailAdditionalProperties(t *testing.T) {
	teardown := setupAndTeardown()
	defer teardown()

	setupHandler := newSetupHandler(ReceiveMail)

	testServer := createTestServer(setupHandler)
	defer testServer.Close()

	// JSON contains the additional property "to" which is not permitted
	const additionalProperties = `{"from":"admin@example.com","to":"no-reply@trivial-tickets.com","subject":"Subject Line","message":"Message Line"}`

	response, err := http.Post(testServer.URL, jsonContentType, createReader(additionalProperties))

	t.Run("POSTError", func(t *testing.T) {
		assert.NoError(t, err, "POST request should be successful")
	})

	t.Run("statusCode", func(t *testing.T) {
		assert.Equal(t, http.StatusBadRequest, response.StatusCode,
			"status code should be 400 Bad Request because of additional property 'to'")

		logResponseBody(t, response)
	})
}

func TestReceiveMailInvalidTypes(t *testing.T) {
	teardown := setupAndTeardown()
	defer teardown()

	setupHandler := newSetupHandler(ReceiveMail)

	testServer := createTestServer(setupHandler)
	defer testServer.Close()

	// This is a valid JSON string, but the type of "from" is invalid
	const invalidTypes = `{"from":42,"subject":"Subject Line","message":"Message Line"}`

	response, err := http.Post(testServer.URL, jsonContentType, createReader(invalidTypes))

	t.Run("POSTError", func(t *testing.T) {
		assert.NoError(t, err, "POST request should be successful")
	})

	t.Run("statusCode", func(t *testing.T) {
		assert.Equal(t, http.StatusBadRequest, response.StatusCode,
			"status code should be 400 Bad Request because of invalid type of 'from'")

		logResponseBody(t, response)
	})
}

func TestReceiveMailInvalidEmailAddress(t *testing.T) {
	teardown := setupAndTeardown()
	defer teardown()

	setupHandler := newSetupHandler(ReceiveMail)

	testServer := createTestServer(setupHandler)
	defer testServer.Close()

	// Email address in "from" is invalid because it contains no top-level domain
	const invalidEmailAddress = `{"from":"invalid@email","subject":"Subject Line","message":"Message Line"}`

	response, err := http.Post(testServer.URL, jsonContentType, createReader(invalidEmailAddress))

	t.Run("POSTError", func(t *testing.T) {
		assert.NoError(t, err, "POST request should be successful")
	})

	t.Run("statusCode", func(t *testing.T) {
		assert.Equal(t, http.StatusBadRequest, response.StatusCode,
			"status code should be 400 Bad Request because of invalid email address")

		logResponseBody(t, response)
	})
}

func TestReceiveMailCreateTicket(t *testing.T) {
	teardown := setupAndTeardown()
	defer teardown()

	setupHandler := newSetupHandler(ReceiveMail)

	testServer := createTestServer(setupHandler)
	defer testServer.Close()

	const createTicket = `{"from":"customer@mail.com","subject":"Issue with computer","message":"My computer is broken!"}`

	response, err := http.Post(testServer.URL, jsonContentType, createReader(createTicket))

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
		assert.Equal(t, buildExpectedJson(structs.JsonMap{
			"status":  http.StatusOK,
			"message": http.StatusText(http.StatusOK),
		}), body, "response body should be JSON with status OK")

		t.Log(string(body))
	})

	t.Run("verifyTicketCreated", func(t *testing.T) {
		t.Run("ticketFile", func(t *testing.T) {
			dirContents, readErr := ioutil.ReadDir(globals.ServerConfig.Tickets)
			assert.NoError(t, readErr, fmt.Sprintf("directory '%s' should exist and be readable within test", globals.ServerConfig.Tickets))

			assert.Equal(t, 1, len(dirContents), fmt.Sprintf("directory '%s' should contain exactly one ticket", globals.ServerConfig.Tickets))
		})

		t.Run("mailFile", func(t *testing.T) {
			dirContents, readErr := ioutil.ReadDir(globals.ServerConfig.Mails)
			assert.NoError(t, readErr, fmt.Sprintf("directory '%s' should exist and be readable within test", globals.ServerConfig.Mails))

			assert.Equal(t, 1, len(dirContents), fmt.Sprintf("directory '%s' should contain exactly one mail (ticket creation)", globals.ServerConfig.Mails))
		})
	})
}

func TestReceiveMailCreateAnswer(t *testing.T) {
	teardown := setupAndTeardown()
	defer teardown()

	setupHandler := newSetupHandler(ReceiveMail)

	testServer := createTestServer(setupHandler)
	defer testServer.Close()

	// Create new test ticket in order to submit an answer to it using the API
	testTicket := ticket.CreateTicket("customer@mail.com", "Issue with Computer", "My computer is broken")
	globals.Tickets[testTicket.Id] = testTicket
	filehandler.WriteTicketFile(globals.ServerConfig.Tickets, &testTicket)

	answerSubject := fmt.Sprintf(`[Ticket \"%s\"] Issue with Computer`, testTicket.Id)

	createAnswerJson := fmt.Sprintf(`{"from":"customer@mail.com","subject":"%s","message":"My computer is broken!"}`, answerSubject)

	response, _ := http.Post(testServer.URL, jsonContentType, createReader(createAnswerJson))

	t.Run("statusCode", func(t *testing.T) {
		assert.Equal(t, http.StatusOK, response.StatusCode,
			"response status should be 200 OK because the request is valid")
	})

	body, readErr := ioutil.ReadAll(response.Body)

	t.Run("bodyReadError", func(t *testing.T) {
		assert.NoError(t, readErr, "reading response body should not return error")
	})

	t.Run("jsonResponse", func(t *testing.T) {
		assert.Equal(t, buildExpectedJson(structs.JsonMap{
			"status":  http.StatusOK,
			"message": http.StatusText(http.StatusOK),
		}), body, "response body should be JSON with status OK")

		t.Log(string(body))
	})

	t.Run("verifyAnswerCreated", func(t *testing.T) {
		t.Run("ticketFile", func(t *testing.T) {
			dirContents, readErr := ioutil.ReadDir(globals.ServerConfig.Tickets)
			assert.NoError(t, readErr, fmt.Sprintf("directory '%s' should exist and be readable within test", globals.ServerConfig.Tickets))

			assert.Equal(t, 1, len(dirContents), fmt.Sprintf("directory '%s' should contain exactly one ticket", globals.ServerConfig.Tickets))
		})

		t.Run("mailFile", func(t *testing.T) {
			dirContents, readErr := ioutil.ReadDir(globals.ServerConfig.Mails)
			assert.NoError(t, readErr, fmt.Sprintf("directory '%s' should exist and be readable within test", globals.ServerConfig.Mails))

			assert.Equal(t, 1, len(dirContents), fmt.Sprintf("directory '%s' should contain exactly one mail (answer creation)", globals.ServerConfig.Mails))
		})
	})
}

func TestReceiveMailCreateAnswerInvalidTicketId(t *testing.T) {
	teardown := setupAndTeardown()
	defer teardown()

	setupHandler := newSetupHandler(ReceiveMail)

	testServer := createTestServer(setupHandler)
	defer testServer.Close()

	// Create new test ticket in order to submit an answer to it using the API
	testTicket := ticket.CreateTicket("customer@mail.com", "Issue with Computer", "My computer is broken")
	globals.Tickets[testTicket.Id] = testTicket
	filehandler.WriteTicketFile(globals.ServerConfig.Tickets, &testTicket)

	// Manipulate the ticket id so that it gets invalid
	manipulatedId := testTicket.Id + "x"
	answerSubject := fmt.Sprintf(`[Ticket \"%s\"] Issue with Computer`, manipulatedId)

	createAnswerJson := fmt.Sprintf(`{"from":"customer@mail.com","subject":"%s","message":"My computer is broken!"}`, answerSubject)

	http.Post(testServer.URL, jsonContentType, createReader(createAnswerJson))

	// The invalid ticket id causes a new ticket to be created,
	// therefore test for two files in the tickets directory
	t.Run("verifyAnswerCreated", func(t *testing.T) {
		t.Run("ticketFile", func(t *testing.T) {
			dirContents, readErr := ioutil.ReadDir(globals.ServerConfig.Tickets)
			assert.NoError(t, readErr, fmt.Sprintf("directory '%s' should exist and be readable within test", globals.ServerConfig.Tickets))

			assert.Equal(t, 2, len(dirContents), fmt.Sprintf("directory '%s' should contain two tickets", globals.ServerConfig.Tickets))
		})

		t.Run("mailFile", func(t *testing.T) {
			dirContents, readErr := ioutil.ReadDir(globals.ServerConfig.Mails)
			assert.NoError(t, readErr, fmt.Sprintf("directory '%s' should exist and be readable within test", globals.ServerConfig.Mails))

			assert.Equal(t, 1, len(dirContents), fmt.Sprintf("directory '%s' should contain exactly one mail (answer creation)", globals.ServerConfig.Mails))
		})
	})
}

func TestMatchAnswerSubject(t *testing.T) {
	ticketId := "100"
	subject := fmt.Sprintf(`[Ticket "%s"] Ticket subject`, ticketId)

	extractedId, matching := matchAnswerSubject(subject)

	t.Run("subjectMatching", func(t *testing.T) {
		assert.True(t, matching, "test subject should match the answering schema")
	})

	t.Run("equalTicketId", func(t *testing.T) {
		assert.Equal(t, ticketId, extractedId, "extracted id should be identical to test id")
	})
}

func TestValidEmailAddress(t *testing.T) {
	email := "admin@example.com"

	assert.True(t, validEmailAddress(email), "email should be valid")
}

func TestCheckRequiredPropertiesSet(t *testing.T) {
	t.Run("propertiesSet", func(t *testing.T) {
		requiredProperties := structs.JsonMap{
			"from":    "admin@example.com",
			"subject": "Subject Line",
			"message": "Message Line",
		}

		err := checkRequiredPropertiesSet(requiredProperties)

		assert.NoError(t, err, "all required properties are set so there should be no error")
	})

	t.Run("propertiesMissing", func(t *testing.T) {
		requiredPropertiesMissing := structs.JsonMap{
			"subject": "Subject Line",
			"message": "Message Line",
		}

		err := checkRequiredPropertiesSet(requiredPropertiesMissing)

		assert.Error(t, err, "missing required properties error")
	})
}

func TestCheckAdditionalPropertiesSet(t *testing.T) {
	t.Run("noAdditionalProperties", func(t *testing.T) {
		noAdditionalProperties := structs.JsonMap{
			"from":    "admin@example.com",
			"subject": "Subject Line",
			"message": "Message Line",
		}

		err := checkAdditionalPropertiesSet(noAdditionalProperties)

		assert.NoError(t, err, "no additional properties set, therefore no error")
	})

	t.Run("additionalProperties", func(t *testing.T) {
		additionalProperties := structs.JsonMap{
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
	t.Run("correctTypes", func(t *testing.T) {
		correctTypes := structs.JsonMap{
			"from":    "admin@example.com",
			"subject": "Subject Line",
			"message": "Message Line",
		}

		err := checkCorrectPropertyTypes(correctTypes)

		assert.NoError(t, err, "all types are correct so there should be no error")
	})

	t.Run("invalidTypes", func(t *testing.T) {
		inCorrectTypes := structs.JsonMap{
			"from":    100,
			"subject": "Subject Line",
			"message": "Message Line",
		}

		err := checkCorrectPropertyTypes(inCorrectTypes)

		assert.Error(t, err, "'from' type is invalid so there should be an error")
	})
}
