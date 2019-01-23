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

// Package client provides functions from the CLI calling
// various API endpoints from the server.
package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/mortenterhart/trivial-tickets/logger/testlogger"
	"github.com/mortenterhart/trivial-tickets/structs"
	"github.com/mortenterhart/trivial-tickets/structs/defaults"
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
 * Package client [tests]
 * Functions from the CLI calling various API endpoints from the server
 */

// testPort defines the port used for test
// CLI configurations.
const testPort = 5743

// test that makePOSTRequest is called with the correct path
func TestFetchEmails(t *testing.T) {
	testlogger.BeginTest()
	defer testlogger.EndTest()

	var inputPath string
	var outputResponse string
	var outputErr error

	get = func(path string) (response string, err error) {
		inputPath = path
		response = outputResponse
		err = outputErr
		return
	}

	testMail := structs.Mail{
		ID:      "1234abc",
		To:      "example@gmx.com",
		Subject: "this is a subject",
		Message: "The message",
	}

	testMails := make(map[string]structs.Mail)
	testMails[testMail.ID] = testMail

	jsonMail, _ := json.MarshalIndent(&testMails, "", "    ")

	outputResponse = string(jsonMail)
	outputErr = nil

	resultMails, resultErr := FetchEmails()

	assert.Equal(t, "api/fetchMails", inputPath)
	assert.Equal(t, testMails, resultMails)
	assert.NoError(t, resultErr)
}

func TestFetchEmailsConnectionError(t *testing.T) {
	testlogger.BeginTest()
	defer testlogger.EndTest()

	clientConfigured = false
	conf := structs.CLIConfig{
		Host: defaults.CliHost,
		Port: testPort,
		Cert: defaults.TestCertificate,
	}
	SetCLIConfig(conf)

	get = makeGETRequest

	mails, fetchErr := FetchEmails()

	assert.Error(t, fetchErr)
	assert.Empty(t, mails)
}

func TestFetchEmailsUnmarshalError(t *testing.T) {
	testlogger.BeginTest()
	defer testlogger.EndTest()

	clientConfigured = false
	conf := structs.CLIConfig{
		Host: defaults.CliHost,
		Port: testPort,
		Cert: defaults.TestCertificate,
	}
	SetCLIConfig(conf)

	get = func(path string) (response string, err error) {
		return "{\"invalid json\":", nil
	}

	mails, fetchErr := FetchEmails()

	assert.Error(t, fetchErr)
	assert.Empty(t, mails)
}

func TestRequests(t *testing.T) {
	testlogger.BeginTest()
	defer testlogger.EndTest()

	clientConfigured = false
	conf := structs.CLIConfig{
		Host: defaults.CliHost,
		Port: testPort,
		Cert: defaults.TestCertificate,
	}
	SetCLIConfig(conf)

	var requestURI string
	var requestPayload string
	var requestMethod string
	var responseMessage string
	var responseCode int

	go http.ListenAndServeTLS(fmt.Sprintf("%s%d", ":", conf.Port), defaults.TestCertificate, defaults.TestKey, nil)
	http.HandleFunc("/", func(responseWriter http.ResponseWriter, request *http.Request) {
		requestURI = request.RequestURI
		requestMethod = request.Method

		data, err := ioutil.ReadAll(request.Body)
		if err != nil {
			fmt.Println(err.Error())
			responseCode = http.StatusInternalServerError
		}

		requestPayload = string(data)
		responseWriter.WriteHeader(responseCode)

		_, err = responseWriter.Write([]byte(responseMessage))
		if err != nil {
			fmt.Println(err.Error())
		}
	})

	// give the server enough time to start. Makes the test more reliable
	time.Sleep(1 * time.Second)

	t.Run("TestMakeGetRequest", func(t *testing.T) {

		t.Run("verifyInputs", func(t *testing.T) {
			inputPath := "the/path"
			responseCode = http.StatusOK
			_, getRequestError := makeGETRequest(inputPath)

			assert.NoError(t, getRequestError)
			assert.Equal(t, "GET", requestMethod)
			assert.Equal(t, "", requestPayload)
			assert.Contains(t, requestURI, inputPath)
		})

		t.Run("verifyOutputs", func(t *testing.T) {
			responseCode = http.StatusOK
			responseMessage = "theResponse"
			response, getRequestError := makeGETRequest("")

			assert.NoError(t, getRequestError)
			assert.Equal(t, responseMessage, response)
		})

		t.Run("verifyServerError", func(t *testing.T) {
			responseCode = http.StatusInternalServerError
			response, getRequestError := makeGETRequest("")

			errorOccurred := getRequestError != nil
			assert.True(t, errorOccurred)
			if errorOccurred {
				assert.Contains(t, getRequestError.Error(), "received error status code:")
			}
			assert.Equal(t, "", response)
		})

		t.Run("verifyRequestError", func(t *testing.T) {
			conf.Host = "notAnIPAddress"
			SetCLIConfig(conf)
			response, getRequestError := makeGETRequest("")

			errorOccurred := getRequestError != nil
			assert.True(t, errorOccurred)
			if errorOccurred {
				assert.Contains(t, getRequestError.Error(), "error sending get request:")
			}
			assert.Equal(t, "", response)
		})

	})

	t.Run("TestMakePostRequest", func(t *testing.T) {

		conf := structs.CLIConfig{
			Host: defaults.CliHost,
			Port: testPort,
			Cert: defaults.TestCertificate,
		}

		SetCLIConfig(conf)

		t.Run("verifyInputs", func(t *testing.T) {
			requestMessage := "someString"
			requestPath := "somePath"
			responseCode = http.StatusOK
			_, sendError := makePOSTRequest(requestMessage, requestPath)

			assert.NoError(t, sendError)
			assert.Equal(t, requestMessage, requestPayload)
			assert.Contains(t, requestURI, requestPath)

		})

		t.Run("verifyOutputs", func(t *testing.T) {
			responseMessage = "theResponse"
			responseCode = http.StatusOK
			response, sendError := makePOSTRequest("", "")

			assert.Equal(t, responseMessage, response)
			assert.NoError(t, sendError)
		})

		t.Run("verifyServerError", func(t *testing.T) {
			responseCode = http.StatusNotFound
			response, sendError := makePOSTRequest("", "")

			errorOccurred := sendError != nil
			assert.True(t, errorOccurred)
			if errorOccurred {
				assert.Contains(t, sendError.Error(), "error with https request. Status code:")
			}
			assert.Equal(t, "", response)

		})

		t.Run("verifyPostError", func(t *testing.T) {
			conf.Host = "notAnIPAddress"
			SetCLIConfig(conf)
			response, sendError := makePOSTRequest("", "")
			assert.Equal(t, "", response)
			errorOccurred := sendError != nil
			assert.True(t, errorOccurred)
			if errorOccurred {
				assert.Contains(t, sendError.Error(), "error sending post request: ")
			}
		})
	})
}

func TestSubmitEmail(t *testing.T) {
	testlogger.BeginTest()
	defer testlogger.EndTest()

	clientConfigured = false

	var inputPayload string
	var inputPath string
	var outputResponse string
	var outputErr error

	post = func(payload string, path string) (response string, err error) {
		inputPayload = payload
		inputPath = path
		response = outputResponse
		err = outputErr
		return
	}

	testMail := `{"from":"example@gmx.com", "subject":"this is a subject", "message": "The message"}`
	resultErr := SubmitEmail(testMail)

	assert.NoError(t, resultErr)
	assert.Equal(t, testMail, inputPayload)
	assert.Equal(t, "api/receive", inputPath)
}

func TestSubmitEmailConnectionError(t *testing.T) {
	testlogger.BeginTest()
	defer testlogger.EndTest()

	clientConfigured = false
	conf := structs.CLIConfig{
		Host: defaults.CliHost,
		Port: testPort,
		Cert: defaults.TestCertificate,
	}
	SetCLIConfig(conf)

	post = makePOSTRequest

	testMail := `{"from":"example@gmx.com", "subject":"this is a subject", "message": "The message"}`
	submitErr := SubmitEmail(testMail)

	assert.Error(t, submitErr)
}

func TestSetServerConfig(t *testing.T) {
	testlogger.BeginTest()
	defer testlogger.EndTest()

	conf := structs.CLIConfig{
		Host: "127.0.0.1",
		Port: 433,
	}

	SetCLIConfig(conf)

	assert.Equal(t, conf, cliConfig)

	conf = structs.CLIConfig{
		Host: "10.168.0.1",
		Port: 1010,
	}

	SetCLIConfig(conf)

	assert.Equal(t, conf, cliConfig)
}

func TestInitializeClient(t *testing.T) {
	testlogger.BeginTest()
	defer testlogger.EndTest()

	cliConfig.Cert = defaults.TestCertificate
	clientConfigured = false

	initializeClient()

	assert.True(t, clientConfigured)
	assert.Equal(t, 5*time.Second, client.Timeout)
	assert.NotEqual(t, http.Transport{}, client.Transport)
}

func TestAcknowledgeEmailReception(t *testing.T) {
	testlogger.BeginTest()
	defer testlogger.EndTest()

	clientConfigured = false
	testMail := structs.Mail{
		ID:      "IdString",
		To:      "example@gmail.com",
		Subject: "example",
		Message: "An example message",
	}

	var inputPayload string
	var inputPath string
	post = func(payload string, path string) (response string, err error) {
		inputPath = path
		inputPayload = payload
		return
	}

	acknowledgementError := AcknowledgeEmailReception(testMail)

	assert.Equal(t, `{"id":"`+testMail.ID+`"}`, inputPayload)
	assert.NoError(t, acknowledgementError)
	assert.Equal(t, "api/verifyMail", inputPath)
}

func TestAcknowledgeEmailReceptionConnectionError(t *testing.T) {
	testlogger.BeginTest()
	defer testlogger.EndTest()

	clientConfigured = false
	conf := structs.CLIConfig{
		Host: defaults.CliHost,
		Port: testPort,
		Cert: defaults.TestCertificate,
	}
	SetCLIConfig(conf)

	post = makePOSTRequest

	testMail := structs.Mail{
		ID:      "IdString",
		To:      "example@gmail.com",
		Subject: "example",
		Message: "An example message",
	}
	acknowledgementError := AcknowledgeEmailReception(testMail)

	assert.Error(t, acknowledgementError)
}
