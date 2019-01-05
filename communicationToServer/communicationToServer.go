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

func FetchEmails() (mails []structs.Mail, err error) {
	response, err := send("", "/api/fetchMails")
	if err != nil {
		return
	}
	err = json.Unmarshal([]byte(response), &mails)
	// TODO implement call to api acknowledging every mail successfully received. Doesn't have to be implemented in this function.
	return
}

func SubmitEmail(mail structs.Mail) (err error) {
	jsonMail, err := json.Marshal(&mail)
	if err != nil {
		return
	}
	resp, err := send(string(jsonMail), "/api/create_ticket")
	println(resp)
	return
}

func sendPost(payload string, path string) (response string, err error) {
	if !clientConfigured {
		initializeClient()
	}
	buffer := bytes.NewBufferString(payload)
	url := "https://" + serverConfig.IPAddr + ":" + strconv.Itoa(int(serverConfig.Port)) + "/" + path
	if url[len(url)-1] != '/' {
		url += "/"
	}
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

func SetServerConfig(config structs.CLIConfig) {
	serverConfig = config
}
