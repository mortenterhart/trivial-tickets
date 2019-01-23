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

// Package structs supplies project-wide needed data
// structures, types and constants for the server and
// the command-line tool.
package structs

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
 * Package structs
 * Project-wide needed structures for data elements
 */

// CliMessage defines different standard messages used by the
// command-line interface that are printed to the console.
type CliMessage string

const (
	// RequestCommandInput is the message for requesting
	// a command input.
	RequestCommandInput CliMessage = "\nTo fetch Mails from the server, type '0'\n" +
		"To send an email to the server type '1'\n" +
		"To exit this program type '2'\n" +
		"Command: "

	// CommandNotAccepted is the error message if a command
	// is not recognized.
	CommandNotAccepted CliMessage = "Input not accepted, error: "

	// RequestEmailAddress it the message to prompt the user
	// for a valid e-mail input.
	RequestEmailAddress CliMessage = "Please enter your email address. It has to be valid.\nEmail address: "

	// RequestSubject is the message to prompt the user for
	// a ticket subject.
	RequestSubject      CliMessage = "Please enter the subject line.\nSubject: "

	// RequestMessage is the message to prompt the user for
	// a ticket message.
	RequestMessage      CliMessage = "Please enter the body of the message.\nMessage: "

	// RequestTicketID is the message to prompt the user for
	// an optional ticket id.
	RequestTicketID     CliMessage = "If applicable please enter the ticket ID. If left empty, " +
		"a new ticket will be created.\nTicket ID (optional): "

	// To is the string for the output of the recipient of
	// an e-mail.
	To      CliMessage = "To: "

	// Subject is the string for the output of the subject
	// of a ticket.
	Subject CliMessage = "Subject: "
)

// CliErrMessage defines different error messages printed
// when the command-line tool faces an error.
type CliErrMessage string

const (
	// Various error messages for wrong user inputs
	TooManyInputs CliErrMessage = "Too many successive wrong user inputs. Aborting program execution.\n"
	NoValidOption CliErrMessage = "Not within the range of valid options"
	EmptyString   CliErrMessage = "string is empty"
	InvalidEmail  CliErrMessage = "not a valid email address"
)
