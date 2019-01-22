// Web API for outgoing mails to be fetched and verified to be sent
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

	"github.com/mortenterhart/trivial-tickets/globals"
	"github.com/mortenterhart/trivial-tickets/logger"
	"github.com/mortenterhart/trivial-tickets/mail_events"
	"github.com/mortenterhart/trivial-tickets/structs"
	"github.com/mortenterhart/trivial-tickets/ticket"
	"github.com/mortenterhart/trivial-tickets/util/jsontools"
	"github.com/stretchr/testify/assert"
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

const jsonContentType = "application/json"

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

func testLogConfig() structs.LogConfig {
	return structs.LogConfig{
		LogLevel:   structs.LevelInfo,
		VerboseLog: false,
		FullPaths:  false,
	}
}

func initializeConfig() {
	config := testServerConfig()
	globals.ServerConfig = &config

	logConfig := testLogConfig()
	globals.LogConfig = &logConfig
}

// Setup and teardown
func TestMain(m *testing.M) {
	initializeConfig()

	os.Exit(m.Run())
}

func mockTicket() structs.Ticket {
	return ticket.CreateTicket("customer@mail.com", "Test subject", "Test message")
}

func createTestServer(handlerFunc http.HandlerFunc) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(handlerFunc))
}

func createReader(data string) io.Reader {
	return strings.NewReader(data)
}

func buildExpectedJson(properties structs.JsonMap) []byte {
	expected, decodeErr := jsontools.MapToJson(properties)
	if decodeErr != nil {
		logger.Error("error while decoding expected JSON string:", decodeErr)
		return nil
	}

	return append(expected, '\n')
}

func cleanupMails() {
	for id := range globals.Mails {
		delete(globals.Mails, id)
	}

	os.RemoveAll(globals.ServerConfig.Mails)
}

func displayDirectoryContents(t *testing.T, dirname string, contents []os.FileInfo) {
	t.Logf("Directory contents of '%s' showing %d file(s):", dirname, len(contents))
	for index, file := range contents {
		t.Logf("%d: %s", index, file.Name())
	}
}

func TestSendMail(t *testing.T) {
	testTicket := mockTicket()

	SendMail(mail_events.NewTicket, testTicket)

	t.Run("storedInMailMap", func(t *testing.T) {
		assert.Equal(t, 1, len(globals.Mails), "sent mail should be stored in the global mail storage")
	})

	t.Run("savedInFile", func(t *testing.T) {
		dirContents, readErr := ioutil.ReadDir(globals.ServerConfig.Mails)

		assert.NoError(t, readErr, "reading mail directory should not return error")

		displayDirectoryContents(t, globals.ServerConfig.Mails, dirContents)
		assert.Equal(t, 1, len(dirContents), "mail directory should contain exactly one mail")
	})

	cleanupMails()
}

func TestFetchMails(t *testing.T) {
	testServer := createTestServer(FetchMails)
	defer testServer.Close()

	t.Run("invalidPOSTRequest", func(t *testing.T) {
		const jsonRequest = `{"id":"test-id"}`

		response, err := http.Post(testServer.URL, jsonContentType, createReader(jsonRequest))

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
				expectedJson, decodeErr := json.MarshalIndent(&globals.Mails, "", "    ")

				assert.NoError(t, decodeErr, "decoding test mail to JSON should not return an error")
				assert.Equal(t, append(expectedJson, '\n'), body, "response should contain JSON representation of mail mapped to its id")
			})
		})
	})

	cleanupMails()
}

func TestVerifyMailSent(t *testing.T) {
	testServer := createTestServer(VerifyMailSent)
	defer testServer.Close()

	t.Run("rejectsGET", func(t *testing.T) {
		response, err := http.Get(testServer.URL)

		t.Run("GETError", func(t *testing.T) {
			assert.NoError(t, err, "GET request should execute successfully")
		})

		t.Run("statusCode", func(t *testing.T) {
			assert.Equal(t, http.StatusMethodNotAllowed, response.StatusCode, "GET requests should be declined with 405 Method Not Allowed")
		})
	})

	t.Run("acceptsPOST", func(t *testing.T) {
		t.Run("invalidJSON", func(t *testing.T) {
			// JSON is invalid because the key "id" is not surrounded by quotes
			const invalidJson = `{id:"mail-id"}`

			response, err := http.Post(testServer.URL, jsonContentType, createReader(invalidJson))

			t.Run("POSTError", func(t *testing.T) {
				assert.NoError(t, err, "POST request should be successful")
			})

			t.Run("statusCode", func(t *testing.T) {
				assert.Equal(t, http.StatusBadRequest, response.StatusCode, "request should return 400 Bad Request due to invalid JSON")
			})
		})

		t.Run("missingRequiredProperties", func(t *testing.T) {
			// Required property "id" is missing
			const missingProperties = `{"email":"admin@example.com"}`

			response, err := http.Post(testServer.URL, jsonContentType, createReader(missingProperties))

			t.Run("POSTError", func(t *testing.T) {
				assert.NoError(t, err, "POST request should be successful")
			})

			t.Run("statusCode", func(t *testing.T) {
				assert.Equal(t, http.StatusBadRequest, response.StatusCode, "request should return 400 Bad Request due to missing properties")
			})
		})

		t.Run("additionalProperties", func(t *testing.T) {
			// Additional property "email" is invalid
			const additionalProperties = `{"id":"mail-id","email":"admin@example.com"}`

			response, err := http.Post(testServer.URL, jsonContentType, createReader(additionalProperties))

			t.Run("POSTError", func(t *testing.T) {
				assert.NoError(t, err, "POST request should be successful")
			})

			t.Run("statusCode", func(t *testing.T) {
				assert.Equal(t, http.StatusBadRequest, response.StatusCode, "request should return 400 Bad Request due to additional properties")
			})
		})

		t.Run("propertyTypes", func(t *testing.T) {
			// "id" property has invalid type for API
			const propertyTypes = `{"id":200}`

			response, err := http.Post(testServer.URL, jsonContentType, createReader(propertyTypes))

			t.Run("POSTError", func(t *testing.T) {
				assert.NoError(t, err, "POST request should be successful")
			})

			t.Run("statusCode", func(t *testing.T) {
				assert.Equal(t, http.StatusBadRequest, response.StatusCode, "request should return 400 Bad Request due to invalid types")
			})
		})

		t.Run("verifyExistingMail", func(t *testing.T) {
			SendMail(mail_events.NewTicket, mockTicket())

			// Since there is only one mail created at this point,
			// retrieve the mail id of the just created mail
			var mailId string
			for id := range globals.Mails {
				mailId = id
			}

			verifyJson := fmt.Sprintf(`{"id":"%s"}`, mailId)

			response, err := http.Post(testServer.URL, jsonContentType, createReader(verifyJson))

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
				expectedJson := buildExpectedJson(structs.JsonMap{
					"verified": true,
					"message":  fmt.Sprintf("mail '%s' was successfully sent and deleted from server cache", mailId),
				})

				assert.Equal(t, expectedJson, body, "mail sending should be verified with correct JSON response")
			})
		})

		t.Run("verifyNonExistingMail", func(t *testing.T) {
			mailId := "notExisting"

			verifyJson := fmt.Sprintf(`{"id":"%s"}`, mailId)

			response, err := http.Post(testServer.URL, jsonContentType, createReader(verifyJson))

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
				expectedJson := buildExpectedJson(structs.JsonMap{
					"verified": false,
					"message":  fmt.Sprintf("mail '%s' does not exist or has already been deleted", mailId),
				})

				assert.Equal(t, expectedJson, body, "mail sending should be verified with correct JSON response")
			})
		})
	})
}
