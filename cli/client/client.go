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
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/mortenterhart/trivial-tickets/structs"
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
 * Package client
 * Functions from the CLI calling various API endpoints from the server
 */

// cliConfig holds the current command-line tool
// configuration.
var cliConfig structs.CLIConfig

// client is an instance of the HTTP Client that
// makes GET and POST requests.
var client http.Client

// clientConfigured indicates whether the HTTP client
// has already been constructed, i.e. if the client
// is already initialized.
var clientConfigured bool

// get and post are replaceable functions for the
// test cases.
var get = makeGETRequest
var post = makePOSTRequest

// FetchEmails sends a GET request to the path `api/fetchMails`
// of the server specified in the cliConfig and expects the JSON
// of an array of structs.Mail in the body of the response.
// The function returns the structs.Mail it received.
func FetchEmails() (mails map[string]structs.Mail, err error) {
	response, err := get("api/fetchMails")
	if err != nil {
		err = fmt.Errorf("error occurred while making the get request: %v", err)
		return
	}

	err = json.Unmarshal([]byte(response), &mails)
	if err != nil {
		err = fmt.Errorf("error occurred while unmarshaling the JSON: %v", err)
		return
	}
	return
}

// AcknowledgeEmailReception sends a POST request with the id
// of the received email to the server.
func AcknowledgeEmailReception(mail structs.Mail) (err error) {
	jsonID := `{"id":"` + mail.ID + `"}`
	_, err = post(jsonID, "api/verifyMail")
	if err != nil {
		err = fmt.Errorf("email acknowledgment failed: %v", err)
	}
	return
}

// SubmitEmail takes a structs.Mail and sends it to the server
// as JSON per POST request.
func SubmitEmail(mail string) (err error) {
	resp, err := post(mail, "api/receive")
	fmt.Println(resp)
	return
}
// makeGETRequest a path string and sends a GET request to the
// "path" on the server specified in the CLI Config.
func makeGETRequest(path string) (response string, err error) {
	if !clientConfigured {
		initializeClient()
	}

	url := "https://" + cliConfig.Host + ":" + strconv.Itoa(int(cliConfig.Port)) + "/" + path
	resp, err := client.Get(url)
	if err != nil {
		err = fmt.Errorf("error sending get request: %v", err)
		return
	}
	defer resp.Body.Close()

	responseData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		err = fmt.Errorf("error occurred while reading httpGet response: %v", err)
		return
	}

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("received error status code: %v %s", resp.Status, string(responseData))
		return
	}

	return string(responseData), nil
}

// makePOSTRequest takes a payload and a path string and sends
// a POST request to the "path" on the server specified in
// CLIConfig with the payload as body.
func makePOSTRequest(payload string, path string) (response string, err error) {
	if !clientConfigured {
		initializeClient()
	}

	reader := strings.NewReader(payload)
	url := fmt.Sprintf("https://%s:%d/%s", cliConfig.Host, cliConfig.Port, path)

	resp, err := client.Post(url, "application/json", reader)
	if err != nil {
		return "", fmt.Errorf("error sending post request: %v", err)
	}
	defer resp.Body.Close()

	responseData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("error with https request. Status code: %s {%s}", resp.Status, string(responseData))
	}
	response = string(responseData)
	return
}

// initializeClient creates a http.Client that accepts the
// certificate of the server when using TLS.
func initializeClient() {
	configFilePath, _ := filepath.Abs(cliConfig.Cert)
	cert, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		log.Fatal(err)
	}
	certPool := x509.NewCertPool()
	certPool.AppendCertsFromPEM(cert)
	client = http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs: certPool,
			},
		},
		Timeout: 5 * time.Second,
	}
	clientConfigured = true
}

// SetCLIConfig sets the local cliConfig variable.
func SetCLIConfig(config structs.CLIConfig) {
	cliConfig = config
}
