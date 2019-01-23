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

// Package cliutils contains helper functions and various
// utilities for the CLI.
package cliutils

import "fmt"

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
 * Package cliutils [examples]
 * Various utilities for CLI
 */

// Example EmptyTicketID demonstrates how the command-line
// tool builds a JSON mail to be sent to the E-Mail Recipience
// API. The mail is not passed a ticket id so that a new ticket
// is created.
func ExampleCreateMail_emptyTicketID() {
	mailJSON := CreateMail("customer@example.com", "Ticket subject", "", "message")

	fmt.Println(mailJSON)
	// Output:
	// {"from":"customer@example.com", "subject":"Ticket subject", "message":"message"}
}

// Example WithTicketID shows how the subject field in the
// JSON string is manipulated if a ticket id was supplied.
// The result is a JSON string that creates a new answer
// to the ticket 42 in case that ticket existed.
func ExampleCreateMail_withTicketID() {
	mailJSON := CreateMail("customer@example.com", "Ticket subject", "42", "message")

	fmt.Println(mailJSON)
	// Output:
	// {"from":"customer@example.com", "subject":"[Ticket \"42\"] Ticket subject", "message":"message"}
}

// Example DifferentCases investigates the validity and
// invalidity of different email addresses and verifies
// or invalidates them.
func ExampleCheckEmailAddress_differentCases() {
	fmt.Println(CheckEmailAddress("customer@example.com"))
	fmt.Println(CheckEmailAddress("email.with.many.dots@host.com"))
	fmt.Println(CheckEmailAddress("email@without-domain"))
	fmt.Println(CheckEmailAddress("many@host.domain.names.com"))
	// Output:
	// true
	// true
	// false
	// true
}
