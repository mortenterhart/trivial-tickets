package communicationToServer

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/mortenterhart/trivial-tickets/structs"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"time"
)

var serverConfig structs.CLIConfig
var client http.Client
var clientConfigured bool
var send = sendPost

// FetchEmails sends a Post request to the path api/fetchMails of the server specified in the cliConfig and expects a JSON of a structs.Mail in the body of the response.
// The function returns the structs.Mail it received.
func FetchEmails() (mails []structs.Mail, err error) {
	response, err := send("", "api/fetchMails")
	if err != nil {
		return
	}
	err = json.Unmarshal([]byte(response), &mails)
	return
}

// AcknowledgeEmailReception sends a post request with the id of the received EMail to the server.
func AcknowledgeEmailReception(mail structs.Mail) (err error) {
	_, err = send(mail.Id, "api/verifyMail")
	if err != nil {
		err = fmt.Errorf("email acknowledgment failed: %v", err)
	}
	return
}

// SubmitEmail takes a structs.Mail and sens it to the server as JSON per post request.
func SubmitEmail(mail string) (err error) {
	resp, err := send(mail, "api/receive")
	println(resp)
	return
}

// sendPost takes a payload and a path string and sends a post request to the "path" on the server specified in CLIConifg with the payload as body.
func sendPost(payload string, path string) (response string, err error) {
	if !clientConfigured {
		initializeClient()
	}
	buffer := bytes.NewBufferString(payload)
	url := "https://" + serverConfig.IPAddr + ":" + strconv.Itoa(int(serverConfig.Port)) + "/" + path
	//if url[len(url)-1] != '/' {
	//	url += "/"
	//}
	resp, err := client.Post(url, "application/json", buffer)
	if err != nil {
		return "", fmt.Errorf("error sending post request: %v", err)
	}
	responseData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response: %v", err)
	}
	status := resp.Status
	if status[0] != '2' {
		return "", errors.New("error with https request. Status code: " + status)
	}
	response = string(responseData)
	return
}

// initializeClient creates a http.Client that accepts the certificate of the server when using TLS.
func initializeClient() {
	configFilePath, _ := filepath.Abs(serverConfig.Cert)
	cert, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		log.Fatal(err)
	}
	certPool := x509.NewCertPool()
	certPool.AppendCertsFromPEM(cert)
	client = http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs: certPool}},
		Timeout: (5 * time.Second)}
	clientConfigured = true
}

// Sets the local serverConfig variable
func SetServerConfig(config structs.CLIConfig) {
	serverConfig = config
}
