package api_in

import (
	"fmt"
	"github.com/mortenterhart/trivial-tickets/util/jsontools"
	"io"
	"io/ioutil"
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

// Test the JSON parse function with an invalid JSON string (ending } is missing)
const invalidJson = `{"email":"admin@example.com","subject":"Subject line","message":"Message line"`

const validJson = `{"email":"admin@example.com","subject":"Subject line","message":"Message line"}`

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

	cleanupTestTickets()
}

func testServerConfig() structs.Config {
	return structs.Config{
		Port:    8443,
		Tickets: "../../files/testtickets",
		Users:   "../../files/users/users.json",
		Mails:   "../../files/mails",
		Cert:    "../../ssl/server.cert",
		Key:     "../../ssl/server.key",
		Web:     "../../www",
	}
}

func cleanupTestTickets() {
	os.RemoveAll(globals.ServerConfig.Tickets)
}

func createTestServer(handler http.Handler) *httptest.Server {
	return httptest.NewServer(handler)
}

func createReader(str string) io.Reader {
	return strings.NewReader(str)
}

func TestReceiveMailRejectsGET(t *testing.T) {
	setupHandler := newSetupHandler(ReceiveMail)

	testServer := createTestServer(setupHandler)
	defer testServer.Close()

	response, err := http.Get(testServer.URL)

	assert.NoError(t, err, "GET request should be successful")
	assert.Equal(t, http.StatusMethodNotAllowed, response.StatusCode,
		"status code should be 405 Method Not Allowed")

	body, readErr := ioutil.ReadAll(response.Body)

	assert.NoError(t, readErr, "reading response body should return no error")

	expectedJSON, _ := jsontools.MapToJson(structs.JsonMap{
		"status":  http.StatusMethodNotAllowed,
		"message": "METHOD_NOT_ALLOWED (GET)",
	})
	assert.Equal(t, append(expectedJSON, '\n'), body,
		"response should be JSON with error message METHOD_NOT_ALLOWED")
}

func TestReceiveMailAcceptsPOST(t *testing.T) {
	setupHandler := newSetupHandler(ReceiveMail)

	testServer := createTestServer(setupHandler)
	defer testServer.Close()

	response, err := http.Post(testServer.URL, jsonContentType, createReader(validJson))

	assert.NoError(t, err, "POST requests should be accepted by the mail API")
	assert.Equal(t, 200, response.StatusCode,
		"status code of POST request should be 200 OK")

	body, readErr := ioutil.ReadAll(response.Body)

	assert.NoError(t, readErr, "reading response body should return no error")
	expectedJSON, _ := jsontools.MapToJson(structs.JsonMap{
		"status":  http.StatusOK,
		"message": "OK",
	})
	assert.Equal(t, append(expectedJSON, '\n'), body,
		"response should be JSON with status OK")
}

/*func TestParseJSONMailWithEmptyBody(t *testing.T) {
	var parsedMail structs.Mail
	err := parseJSONMail(createReader(""), &parsedMail)

	assert.Error(t, err, "nil body is invalid and should cause error")
	assert.Equal(t, "EOF", err.Error(), "JSON string is malformed")
	assert.Empty(t, parsedMail.Email, "email field should be empty")
	assert.Empty(t, parsedMail.Subject, "subject field should be empty")
	assert.Empty(t, parsedMail.Message, "message field should be empty")
}

func TestParseJSONMailWithInvalidJSON(t *testing.T) {
	var parsedMail structs.Mail
	err := parseJSONMail(createReader(invalidJson), &parsedMail)

	assert.Error(t, err, "JSON string is invalid and should cause error")
	assert.Equal(t, "unexpected EOF", err.Error(), "JSON string is malformed")
	assert.Empty(t, parsedMail.Email, "email field should be empty")
	assert.Empty(t, parsedMail.Subject, "subject field should be empty")
	assert.Empty(t, parsedMail.Message, "message field should be empty")
}

func TestParseJSONMailWithValidJSON(t *testing.T) {
	var parsedMail structs.Mail
	err := parseJSONMail(createReader(validJson), &parsedMail)

	assert.NoError(t, err, "valid JSON should not cause error")
	assert.Equal(t, "admin@example.com", parsedMail.Email, "email field should be equal to JSON")
	assert.Equal(t, "Subject line", parsedMail.Subject, "subject field should be equal to JSON")
	assert.Equal(t, "Message line", parsedMail.Message, "message field should be equal to JSON")
}

func TestExtractMailWithInvalidJSON(t *testing.T) {
	request := httptest.NewRequest(http.MethodPost, "/api/submit_mail", createReader(invalidJson))

	mail, err := extractMail(request)

	assert.Error(t, err, "invalid JSON should cause error while extracting")
	assert.Equal(t, "unexpected EOF", err.Error(), "JSON string is malformed")
	assert.Empty(t, mail.Email, "email field should be empty")
	assert.Empty(t, mail.Subject, "subject field should be empty")
	assert.Empty(t, mail.Message, "message field should be empty")
}

func TestExtractMailWithValidJSON(t *testing.T) {
	request := httptest.NewRequest(http.MethodPost, "/api/submit_mail", createReader(validJson))

	mail, err := extractMail(request)

	assert.NoError(t, err, "valid JSON string should be parsed without error")
	assert.Equal(t, "admin@example.com", mail.Email, "email field should be equal to JSON")
	assert.Equal(t, "Subject line", mail.Subject, "subject field should be equal to JSON")
	assert.Equal(t, "Message line", mail.Message, "message field should be equal to JSON")
}*/

func TestCheckRequiredPropertiesSetWithInvalidJSON(t *testing.T) {
	setupHandler := newSetupHandler(ReceiveMail)

	testServer := createTestServer(setupHandler)
	defer testServer.Close()

	response, _ := http.Post(testServer.URL, jsonContentType, createReader(`{"email":"invalid@address.com","subject":"Subject","message":""}`))
	body, _ := ioutil.ReadAll(response.Body)
	fmt.Println(string(body))
}
