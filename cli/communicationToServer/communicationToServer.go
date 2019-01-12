package communicationToServer

import (
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
	"strings"
	"time"
)

/*
 * Ticketsystem Trivial Tickets
 *
 * Matriculation numbers: 3040018, 3040018, 3478222
 * Lecture:               Programmieren II, INF16B
 * Lecturer:              Herr Prof. Dr. Helmut Neemann
 * Institute:             Duale Hochschule Baden-Württemberg Mosbach
 *
 * ---------------
 *
 * Package communicationToServer
 * Functions from the CLI calling various API endpoints from the server
 */

var serverConfig structs.CLIConfig
var client http.Client
var clientConfigured bool
var post = makePostRequest
var get = makeGetRequest

// FetchEmails sends a Get request to the path api/fetchMails of the server specified in the cliConfig and expects the JSON of an array of structs.Mail in the body of the response.
// The function returns the structs.Mail it received.
func FetchEmails() (mails map[string]structs.Mail, err error) {
	response, err := get("api/fetchMails")
	if err != nil {
		err = fmt.Errorf("error occured while making the get request: %v", err)
		return
	}
	err = json.Unmarshal([]byte(response), &mails)
	if err != nil {
		err = fmt.Errorf("error occured while unmarshaling the JSON: %v", err)
		return
	}
	return
}

// AcknowledgeEmailReception sends a post request with the id of the received EMail to the server.
func AcknowledgeEmailReception(mail structs.Mail) (err error) {
	jsonID := `{"id":"` + mail.Id + `"}`
	_, err = post(jsonID, "api/verifyMail")
	if err != nil {
		err = fmt.Errorf("email acknowledgment failed: %v", err)
	}
	return
}

// SubmitEmail takes a structs.Mail and sens it to the server as JSON per post request.
func SubmitEmail(mail string) (err error) {
	resp, err := post(mail, "api/receive")
	println(resp)
	return
}

func makeGetRequest(path string) (response string, err error) {
	if !clientConfigured {
		initializeClient()
	}
	url := "https://" + serverConfig.IPAddr + ":" + strconv.Itoa(int(serverConfig.Port)) + "/" + path
	resp, err := client.Get(url)
	if err != nil {
		err = fmt.Errorf("error sending get request: %v", err)
		return
	}
	responseData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		err = fmt.Errorf("error occured while reading httpGet response: %v", err)
		return
	}
	if resp.Status[0] != '2' {
		err = fmt.Errorf("received error status code: %v %s", resp.Status, string(responseData))
		return
	}

	return string(responseData), nil
}

// makePostRequest takes a payload and a path string and sends a post request to the "path" on the server specified in CLIConifg with the payload as body.
func makePostRequest(payload string, path string) (response string, err error) {
	if !clientConfigured {
		initializeClient()
	}
	reader := strings.NewReader(payload)
	url := "https://" + serverConfig.IPAddr + ":" + strconv.Itoa(int(serverConfig.Port)) + "/" + path
	//if url[len(url)-1] != '/' {
	//	url += "/"
	//}
	resp, err := client.Post(url, "application/json", reader)
	if err != nil {
		return "", fmt.Errorf("error sending post request: %v", err)
	}
	responseData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response: %v", err)
	}
	status := resp.Status
	if status[0] != '2' {
		return "", errors.New("error with https request. Status code: " + status + " {" + string(responseData) + "}")
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
