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
 * Package trivial_tickets
 * Root Package of the Trivial Tickets Ticketsystem
 */

// Package trivial_tickets is a basic implementation of a support ticket
// system in Go. Customers of a company can easily create tickets with
// the form on the home page and display the created tickets over a
// static permalink. The author of a ticket is informed about every
// action the ticket experiences by an email written to his email
// address.
//
// Available Operations
//
// Registered users are the assignees of the system and work on the
// customer's tickets. They can log in and off the system and assign
// a ticket to themselves or to other users. An assignee can only
// edit information on those ticket that he is assigned to. He can
// release his tickets, add comments to it or change the status of
// the ticket. It is also possible to merge two tickets to one, but
// only if the customer and the assignee of both tickets match.
// Additionally, an assignee may indicate that he is on holiday. In
// this case, tickets cannot be assigned to him.
//
// The E-Mail Recipience API
//
// The ticket system offers an E-Mail Recipience and Dispatch API for
// a mailing service to interact with the system. A request with the
// properties `from` (the sender's e-mail address, i.e the customer),
// `subject` (the ticket's subject) and `message` (the actual message)
// set can be delivered to the server in order to create new tickets by
// mail. If the subject contains the ticket id of an already existing
// ticket in a special markup the e-mail creates a new answer to this
// ticket instead of a new ticket.
//
// The E-Mail Dispatch API
//
// Likewise, the mailing service can fetch e-mails created by the server
// on special events (such as creation of a new ticket or update of a
// ticket). The server then sends back all e-mails remaining to be sent
// to the customer. After the service has sent the received e-mails it
// can confirm the successful sending of the concerned e-mails by sending
// a verification request to the server. This is done for each e-mail by
// applying the id of the certain mail to the request. If the server can
// ensure that the respective e-mail exists the mail can be safely deleted
// and the verification result is returned to the service. More information
// on the Mail API and its functionality can be found in the Mail API Reference
// on https://github.com/mortenterhart/trivial-tickets/wiki/Mail-API-Reference.
//
// The Command-line Tool
//
// The repository contains a command-line tool which serves as simple
// example for an external mailing service. Although it cannot send or
// receive genuine e-mails it can communicate with the mail API of the
// server. The user can write really simple e-mails with e-mail address,
// subject and message to be sent to the server. Alternatively, he may
// retrieve all server-side created e-mails that are print to the console
// then. All fetched e-mails are verified against the server who can delete
// the mails then safely.
//
// Further Information
//
// Consult the User Manual for more information (currently only available
// in German) or alternatively checkout our Wiki Pages on
// https://github.com/mortenterhart/trivial-tickets/wiki.
package trivial_tickets
