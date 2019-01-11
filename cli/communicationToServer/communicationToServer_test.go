package communicationToServer

import (
	"encoding/json"
	"fmt"
	"github.com/mortenterhart/trivial-tickets/structs"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
)

// test that makePostRequest is called with the correct path
func TestFetchEmails(t *testing.T) {
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
		Id:      "1234abc",
		To:      "example@gmx.com",
		Subject: "this is a subject",
		Message: "The message"}
	testMails := make([]structs.Mail, 0)
	testMails = append(testMails, testMail, testMail)
	jsonMail, _ := json.MarshalIndent(&testMails, "", "	")
	outputResponse = string(jsonMail)
	outputErr = nil
	resultMails, resultErr := FetchEmails()
	assert.Equal(t, "api/fetchMails", inputPath)
	assert.Equal(t, testMails, resultMails)
	assert.NoError(t, resultErr)

}

func TestMakeGetRequest(t *testing.T) {
	clientConfigured = false
	conf := structs.CLIConfig{
		IPAddr: "localhost",
		Port:   4443,
		Cert:   "../../ssl/server.cert"}
	SetServerConfig(conf)
	var requestURI string
	var requestPayload string
	var requestMethod string
	var responseMessage string
	var responseCode int
	go http.ListenAndServeTLS(fmt.Sprintf("%s%d", ":", conf.Port), "../../ssl/server.cert", "../../ssl/server.key", nil)
	http.HandleFunc("/", func(responseWriter http.ResponseWriter, request *http.Request) {
		requestURI = request.RequestURI
		requestMethod = request.Method
		data, err := ioutil.ReadAll(request.Body)
		if err != nil {
			responseCode = 500
		}
		requestPayload = string(data)
		responseWriter.WriteHeader(responseCode)
		responseWriter.Write([]byte(responseMessage))
	})

	t.Run("verifyInputs", func(t *testing.T) {
		inputPath := "the/path"
		responseCode = 200
		makeGetRequest(inputPath)

		assert.Equal(t, "GET", requestMethod)
		assert.Equal(t, "", requestPayload)
		assert.Contains(t, requestURI, inputPath)
	})

	t.Run("verifyOutputs", func(t *testing.T) {
		responseCode = 200
		responseMessage = "theResponse"
		response, getRequestError := makeGetRequest("")

		assert.NoError(t, getRequestError)
		assert.Equal(t, responseMessage, response)
	})

	t.Run("verifyServerError", func(t *testing.T) {
		responseCode = 500
		response, getRequestError := makeGetRequest("")

		errorOccured := getRequestError != nil
		assert.True(t, errorOccured)
		if errorOccured {
			assert.Contains(t, getRequestError.Error(), "received error status code:")
		}
		assert.Equal(t, "", response)
	})

	t.Run("verifyRequestError", func(t *testing.T) {
		conf.IPAddr = "notAnIPAddress"
		SetServerConfig(conf)
		response, getRequestError := makeGetRequest("")

		errorOccured := getRequestError != nil
		assert.True(t, errorOccured)
		if errorOccured {
			assert.Contains(t, getRequestError.Error(), "error sending get request:")
		}
		assert.Equal(t, "", response)
	})
}

func TestSubmitEmail(t *testing.T) {
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

func TestSetServerConfig(t *testing.T) {
	conf := structs.CLIConfig{
		IPAddr: "127.0.0.1",
		Port:   433}
	SetServerConfig(conf)
	assert.Equal(t, conf, serverConfig)
	conf = structs.CLIConfig{
		IPAddr: "10.168.0.1",
		Port:   1010}
	SetServerConfig(conf)
	assert.Equal(t, conf, serverConfig)
}

func TestInitializeClient(t *testing.T) {
	serverConfig.Cert = "../../ssl/server.cert"
	clientConfigured = false
	initializeClient()
	assert.True(t, clientConfigured)
	assert.Equal(t, 5*time.Second, client.Timeout)
	assert.NotEqual(t, http.Transport{}, client.Transport)
}

func TestMakePostRequest(t *testing.T) {
	clientConfigured = false
	conf := structs.CLIConfig{
		IPAddr: "localhost",
		Port:   4443,
		Cert:   "../../ssl/server.cert"}
	SetServerConfig(conf)
	var requestURI string
	var requestPayload string
	var responseMessage string
	var responseCode int
	go http.ListenAndServeTLS(fmt.Sprintf("%s%d", ":", conf.Port), "../../ssl/server.cert", "../../ssl/server.key", nil)
	http.HandleFunc("/", func(responseWriter http.ResponseWriter, request *http.Request) {
		requestURI = request.RequestURI
		data, err := ioutil.ReadAll(request.Body)
		if err != nil {
			responseCode = 500
		}
		requestPayload = string(data)
		responseWriter.WriteHeader(responseCode)
		responseWriter.Write([]byte(responseMessage))
	})
	t.Run("verifyInputs", func(t *testing.T) {
		requestMessage := "someString"
		requestPath := "somePath"
		responseCode = 200
		_, sendError := makePostRequest(requestMessage, requestPath)

		assert.NoError(t, sendError)
		assert.Equal(t, requestMessage, requestPayload)
		assert.Contains(t, requestURI, requestPath)
	})

	t.Run("verifyOutputs", func(t *testing.T) {
		responseMessage = "theResponse"
		responseCode = 200
		response, sendError := makePostRequest("", "")

		assert.Equal(t, responseMessage, response)
		assert.NoError(t, sendError)
	})

	t.Run("verifyServerError", func(t *testing.T) {
		responseCode = 404
		response, sendError := makePostRequest("", "")

		errorOccured := sendError != nil
		assert.True(t, errorOccured)
		if errorOccured {
			assert.Contains(t, sendError.Error(), "error with https request. Status code:")
		}
		assert.Equal(t, "", response)

	})

	t.Run("verifyPostError", func(t *testing.T) {
		conf.IPAddr = "notAnIPAddress"
		SetServerConfig(conf)
		response, sendError := makePostRequest("", "")
		assert.Equal(t, "", response)
		errorOccurred := sendError != nil
		assert.True(t, errorOccurred)
		if errorOccurred {
			assert.Contains(t, sendError.Error(), "error sending post request: ")
		}
	})
}

func TestAcknowledgeEmailReception(t *testing.T) {
	testMail := structs.Mail{
		Id:      "IdString",
		To:      "example@gmail.com",
		Subject: "example",
		Message: "An example message"}
	var inputPayload string
	var inputPath string
	post = func(payload string, path string) (response string, err error) {
		inputPath = path
		inputPayload = payload
		return
	}
	acknowledgementError := AcknowledgeEmailReception(testMail)

	assert.Equal(t, testMail.Id, inputPayload)
	assert.NoError(t, acknowledgementError)
	assert.Equal(t, "api/verifyMail", inputPath)
}
