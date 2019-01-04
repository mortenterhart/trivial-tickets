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

// test that sendPost is called with the correct path
// we're ignoring the body of the request so far in this test
// make sure it interprets the response payload correctly
func TestFetchEmails(t *testing.T) {
	var inputPath string
	var outputResponse string
	var outputErr error
	send = func(payload string, path string) (response string, err error) {
		inputPath = path
		response = outputResponse
		err = outputErr
		return
	}
	testMail := structs.Mail{
		Email:   "example@gmx.com",
		Subject: "this is a subject",
		Message: "The message"}
	testMails := make([]structs.Mail, 0)
	testMails = append(testMails, testMail, testMail)
	jsonMail, _ := json.MarshalIndent(&testMails, "", "	")
	outputResponse = string(jsonMail)
	outputErr = nil
	resultMails, resultErr := FetchEmails()
	assert.Equal(t, "/api/fetchMails", inputPath)
	assert.Equal(t, testMails, resultMails)
	assert.NoError(t, resultErr)
}

func TestSubmitEmail(t *testing.T) {
	var inputPayload string
	var inputPath string
	var outputResponse string
	var outputErr error
	send = func(payload string, path string) (response string, err error) {
		inputPayload = payload
		inputPath = path
		response = outputResponse
		err = outputErr
		return
	}
	testMail := structs.Mail{
		Email:   "example@gmx.com",
		Subject: "this is a subject",
		Message: "The message"}
	resultErr := SubmitEmail(testMail)
	assert.NoError(t, resultErr)
	jsonMail, _ := json.Marshal(&testMail)
	assert.Equal(t, string(jsonMail), inputPayload)
	assert.Equal(t, "/api/create_ticket", inputPath)

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
	serverConfig.Cert = "../ssl/server.cert"
	clientConfigured = false
	initializeClient()
	assert.True(t, clientConfigured)
	assert.Equal(t, 4*time.Second, client.Timeout)
	assert.NotEqual(t, http.Transport{}, client.Transport)
}

func TestSendPost(t *testing.T) {
	clientConfigured = false
	send = sendPost
	conf := structs.CLIConfig{
		IPAddr: "localhost",
		Port:   443,
		Cert:   "../ssl/server.cert"}
	SetServerConfig(conf)
	var requestURI string
	var requestPayload string
	var responseMessage string
	var responseCode int
	go http.ListenAndServeTLS(fmt.Sprintf("%s%d", ":", conf.Port), "../ssl/server.cert", "../ssl/server.key", nil)
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
		send(requestMessage, requestPath)

		assert.Equal(t, requestMessage, requestPayload)
		assert.Contains(t, requestURI, requestPath)
	})

	t.Run("verifyOutputs", func(t *testing.T) {
		responseMessage = "theResponse"
		responseCode = 200
		response, sendError := send("", "")

		assert.Equal(t, responseMessage, response)
		assert.NoError(t, sendError)
	})

	t.Run("verifyServerError", func(t *testing.T) {
		responseCode = 404
		response, sendError := send("", "")

		assert.Equal(t, "", response)
		assert.EqualError(t, sendError, "error with https request. Status code: 404 Not Found")
	})

	t.Run("verifyPostError", func(t *testing.T) {
		conf.IPAddr = "notAnIPAddress"
		SetServerConfig(conf)
		response, sendError := send("", "")
		assert.Equal(t, "", response)
		errorOccurred := sendError != nil
		assert.True(t, errorOccurred)
		if errorOccurred {
			assert.Contains(t, sendError.Error(), "error sending post request: ")
		}
	})
}
