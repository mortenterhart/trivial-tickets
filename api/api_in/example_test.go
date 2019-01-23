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

// Package api_in implements a web interface for incoming mails
// to create new tickets or answers
package api_in

import (
	"fmt"
	"github.com/mortenterhart/trivial-tickets/globals"
	"github.com/mortenterhart/trivial-tickets/ticket"
	"github.com/mortenterhart/trivial-tickets/util/filehandler"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
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
 * Package api_in [examples]
 * Web API for incoming mails to create new tickets or answers
 */

func ExampleReceiveMail_createNewTicket() {
	// Create a test server with the ReceiveMail
	// handler registered
	server := httptest.NewServer(http.HandlerFunc(ReceiveMail))
	defer server.Close()

	// Build the JSON Request
	jsonRequest := `
{
	"from": "email@example.com",
	"subject": "New Ticket created",
	"message": "My new ticket was just created!"
}`

	// Make the POST request to the ReceiveMail API
	client := server.Client()
	response, errPost := client.Post(server.URL, "application/json", strings.NewReader(jsonRequest))
	defer func() {
		// Close the response body if there was
		// no error
		if errPost == nil {
			response.Body.Close()
		}
	}()

	// Read the response body
	body, _ := ioutil.ReadAll(response.Body)

	fmt.Println(string(body))
	// Output:
	// {
	//     "message": "OK",
	//     "status": 200
	// }
}

func ExampleReceiveMail_createNewAnswer() {
	// Create a test server with the ReceiveMail
	// handler registered
	server := httptest.NewServer(http.HandlerFunc(ReceiveMail))
	defer server.Close()

	// Retrieve the test server config
	testConfig := testServerConfig()

	// Create a new ticket and write it to the
	// file system to attach a new answer to it
	// using the ReceiveMail API.
	newTicket := ticket.CreateTicket("email@example.com", "New ticket with answer",
		"New answers can also be created using an email request")
	globals.Tickets[newTicket.ID] = newTicket
	errWrite := filehandler.WriteTicketFile(testConfig.Tickets, &newTicket)
	if errWrite != nil {
		fmt.Println(errWrite)
	}

	// Build the JSON Request with special subject
	// markup:
	//
	//   [Ticket "<ticket-id>"] <Subject>
	//
	// If the specified ticket id exists, the message
	// is attached as new answer to the existing ticket.
	jsonRequest := fmt.Sprintf(`
{
	"from": "email@example.com",
	"subject": "[Ticket \"%s\"] New Ticket with answer",
	"message": "This answer will be attached to the existing ticket."
}`, newTicket.ID)

	// Make the POST request to the ReceiveMail API
	client := server.Client()
	response, errPost := client.Post(server.URL, "application/json", strings.NewReader(jsonRequest))
	defer func() {
		// Close the response body if there was
		// no error
		if errPost == nil {
			response.Body.Close()
		}
	}()

	// Read the response body
	body, _ := ioutil.ReadAll(response.Body)

	fmt.Println(string(body))
	// Output:
	// {
	//     "message": "OK",
	//     "status": 200
	// }
}
