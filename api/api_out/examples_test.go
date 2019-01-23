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

// Package api_out implements a web interface for outgoing mails
// to be fetched and verified to be sent
package api_out

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/mortenterhart/trivial-tickets/globals"
	"github.com/mortenterhart/trivial-tickets/structs"
	"github.com/mortenterhart/trivial-tickets/util/filehandler"
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
 * Package api_out [examples]
 * Web API for outgoing mails to be fetched and verified to be sent
 */

// Example NoCachedMails requests the created mails from the server
// while there is no mail cached at runtime. The result is an empty
// JSON object.
func ExampleFetchMails_noCachedMails() {
	// Create a test server to handle the API
	// request tp the FetchMails API
	server := httptest.NewServer(http.HandlerFunc(FetchMails))
	defer server.Close()

	// Create a client to make the request
	client := server.Client()

	// Make the GET request to the API while
	// the server does not cache any mail
	response, errGet := client.Get(server.URL)
	defer func() {
		// Close the response body if there
		// was no error
		if errGet == nil {
			response.Body.Close()
		}
	}()

	// Read the response body
	body, _ := ioutil.ReadAll(response.Body)

	fmt.Println(string(body))
	// Output: {}
}

// Example OneCachedMail performs a GET request to the FetchMails API
// when one mail is cached by the server. The mail is created explicitly
// for this example which would be usually resumed by the server. The
// caller of this API has therefore nothing to do with the mail creation.
// Finally, the response contains the server-created mail as JSON.
func ExampleFetchMails_oneCachedMail() {
	// Create a test server to handle the API
	// request to the FetchMails handler
	server := httptest.NewServer(http.HandlerFunc(FetchMails))
	defer server.Close()

	// Create a client to make the request
	client := server.Client()

	// Get the server config
	config := testServerConfig()

	// --- Begin mail creation (irrelevant for client) ---

	// Simulate a created mail by generating
	// a test mail. Usually the server adopts
	// this task and the caller has nothing to
	// do with it.
	newMail := structs.Mail{
		ID:      "mail-id",
		From:    "customer@example.com",
		To:      "editor@trivial-tickets.com",
		Subject: "New ticket created",
		Message: "This is a notification about a newly created ticket",
	}

	// Store the mail in the mail map
	globals.Mails[newMail.ID] = newMail

	// Save the mail to a file
	errWrite := filehandler.WriteMailFile(config.Mails, &newMail)
	if errWrite != nil {
		fmt.Println(errWrite)
	}

	// --- End mail creation ---

	// Get the just created mail with a call to the
	// FetchMails API
	response, errGet := client.Get(server.URL)
	defer func() {
		// Close the response body if there was
		// no error
		if errGet == nil {
			response.Body.Close()
		}
	}()

	// Read the response body
	body, _ := ioutil.ReadAll(response.Body)

	fmt.Println(string(body))
	// Output:
	// {
	//     "mail-id": {
	//         "id": "mail-id",
	//         "from": "customer@example.com",
	//         "to": "editor@trivial-tickets.com",
	//         "subject": "New ticket created",
	//         "message": "This is a notification about a newly created ticket"
	//     }
	// }
}

// Example SendingVerified illustrates a successful sending verification
// of a mail fetched from the server. The mail with the corresponding
// mail id was found at the server and could be deleted successfully
// from cache. The JSON response denotes a successful verification with
// the `verified` field.
func ExampleVerifyMailSent_sendingVerified() {
	// Create a test server to handle the API
	// request to the VerifyMailSent handler
	server := httptest.NewServer(http.HandlerFunc(VerifyMailSent))
	defer server.Close()

	// Create a client to do the request with
	client := server.Client()

	// Get the server config
	config := testServerConfig()

	// --- Begin mail creation (irrelevant for client) ---

	// Simulate a created mail by generating
	// a test mail. Usually the server adopts
	// this task and the caller has nothing to
	// do with it.
	newMail := structs.Mail{
		ID:      "mail-id",
		From:    "customer@example.com",
		To:      "editor@trivial-tickets.com",
		Subject: "New ticket created",
		Message: "This is a notification about a newly created ticket",
	}

	// Store the mail in the mail map
	globals.Mails[newMail.ID] = newMail

	// Save the mail to a file
	errWrite := filehandler.WriteMailFile(config.Mails, &newMail)
	if errWrite != nil {
		fmt.Println(errWrite)
	}

	// --- End mail creation ---

	// Write the JSON request with the mail id to
	// be deleted
	jsonRequest := `{ "id": "mail-id" }`

	// Do the actual POST request with the JSON body
	response, errPost := client.Post(server.URL, "application/json", strings.NewReader(jsonRequest))
	defer func() {
		// Close the response body if there was
		// no error during request
		if errPost == nil {
			response.Body.Close()
		}
	}()

	// Read the response body
	body, _ := ioutil.ReadAll(response.Body)

	fmt.Println(string(body))
	// Output:
	// {
	//     "message": "mail 'mail-id' was successfully sent and deleted from server cache",
	//     "status": 200,
	//     "verified": true
	// }
}

// Example SendingNotVerified shows the API response in
// case the applied mail id does not belong to an existing
// mail. In this example, the mail id 'another-mail-id' in
// the JSON request does not exist. The server can not find
// any associated mail and sends an unverified response.
func ExampleVerifyMailSent_sendingNotVerified() {
	// Create a test server to handle the API
	// request to the VerifyMailSent handler
	server := httptest.NewServer(http.HandlerFunc(VerifyMailSent))
	defer server.Close()

	// Create a client to do the request with
	client := server.Client()

	// Write the JSON request with the mail id of a not
	// existing mail at the server
	jsonRequest := `{ "id": "another-mail-id" }`

	// Do the actual POST request with the JSON body
	response, errPost := client.Post(server.URL, "application/json", strings.NewReader(jsonRequest))
	defer func() {
		// Close the response body if there was
		// no error during request
		if errPost == nil {
			response.Body.Close()
		}
	}()

	// Read the response body
	body, _ := ioutil.ReadAll(response.Body)

	fmt.Println(string(body))
	// Output:
	// {
	//     "message": "mail 'another-mail-id' does not exist or has already been deleted",
	//     "status": 200,
	//     "verified": false
	// }
}
