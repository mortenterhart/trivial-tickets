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

// Package mail_events provides facilities to create
// standard mail messages for different actions using
// predefined templates.
package mail_events

import (
	"fmt"

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
 * Package mail_events [examples]
 * Mail message construction using templating
 */

// Example NewTicket shows how the email message for
// a new ticket is built and looks like.
func ExampleNewMailBody_newTicket() {
	editor := structs.User{
		ID:       "editor-id",
		Name:     "Example Editor",
		Username: "editor",
		Mail:     "editor@trivial-tickets.com",
	}

	ticket := structs.Ticket{
		ID:       "example-ticket",
		Subject:  "Printer not working",
		Status:   structs.StatusOpen,
		User:     editor,
		Customer: "desparate_user@example.com",
		Entries: []structs.Entry{
			{
				User:      "desparate_user@example.com",
				Text:      "I cannot make the printer print my document :(",
				ReplyType: "external",
			},
		},
	}

	mailBody := NewMailBody(NewTicket, ticket)

	fmt.Println(mailBody)
	// Output:
	// Dear Customer,
	//
	// Your Ticket 'example-ticket' was created successfully.
	// If you want to write a new comment to this ticket,
	// please use the following link: mailto:support@trivial-tickets.com?subject=%5BTicket%20%22example-ticket%22%5D%20Printer%20not%20working
	//
	// -----------------------------
	// Customer:   desparate_user@example.com
	// Ticket Key: example-ticket
	// URL:        https://localhost:8444/ticket?id=example-ticket
	// Editor:     Example Editor (editor@trivial-tickets.com)
	// Status:     Open
	//
	// Subject: Printer not working
	//
	// I cannot make the printer print my document :(
	//
	// -----------------------------
	//
	// Kind Regards,
	// Your Trivial Tickets Team
	//
	// This message was automatically generated by trivial-tickets.com.
	// Please do not reply to this e-mail.
}

// Example NewAnswer shows how the email message for
// a new answer is built and looks like.
func ExampleNewMailBody_newAnswer() {
	editor := structs.User{
		ID:       "editor-id",
		Name:     "Example Editor",
		Username: "editor",
		Mail:     "editor@trivial-tickets.com",
	}

	answer := structs.Entry{
		User:      editor.Mail,
		Text:      "Of course you can not when the printer is turned off",
		ReplyType: "external",
	}

	ticket := structs.Ticket{
		ID:       "example-ticket",
		Subject:  "Printer not working",
		Status:   structs.StatusOpen,
		User:     editor,
		Customer: "desparate_user@example.com",
		Entries: []structs.Entry{
			{
				User:      "desparate_user@example.com",
				Text:      "I cannot make the printer print my document :(",
				ReplyType: "external",
			},
			answer,
		},
	}

	mailBody := NewMailBody(NewAnswer, ticket)

	fmt.Println(mailBody)
	// Output:
	// Dear Customer,
	//
	// the user 'editor@trivial-tickets.com' wrote a new comment to your ticket:
	//
	// -----------------------------
	// Customer:   desparate_user@example.com
	// Ticket Key: example-ticket
	// URL:        https://localhost:8444/ticket?id=example-ticket
	// Editor:     Example Editor (editor@trivial-tickets.com)
	// Status:     Open
	//
	// Subject: Printer not working
	//
	// Of course you can not when the printer is turned off
	//
	// -----------------------------
	//
	// Kind Regards,
	// Your Trivial Tickets Team
	//
	// This message was automatically generated by trivial-tickets.com.
	// Please do not reply to this e-mail.
}

// Example AssignedTicket shows how the email message
// for a newly assigned ticket is built and looks like.
func ExampleNewMailBody_assignedTicket() {
	editor := structs.User{
		ID:       "editor-id-2",
		Name:     "Example Editor 2",
		Username: "editor-2",
		Mail:     "editor-2@trivial-tickets.com",
	}

	ticket := structs.Ticket{
		ID:       "example-ticket-2",
		Subject:  "I need help",
		Status:   structs.StatusOpen,
		User:     editor,
		Customer: "desparate_user@example.com",
		Entries: []structs.Entry{
			{
				User:      "desparate_user@example.com",
				Text:      "Can somebody help me with my problem?",
				ReplyType: "external",
			},
		},
	}

	mailBody := NewMailBody(AssignedTicket, ticket)

	fmt.Println(mailBody)
	// Output:
	// Dear Customer,
	//
	// the editor 'Example Editor 2' works on Your Ticket now:
	//
	// -----------------------------
	// Customer:   desparate_user@example.com
	// Ticket Key: example-ticket-2
	// URL:        https://localhost:8444/ticket?id=example-ticket-2
	// Editor:     Example Editor 2 (editor-2@trivial-tickets.com)
	// Status:     Open
	//
	// Subject: I need help
	//
	// Can somebody help me with my problem?
	//
	// -----------------------------
	//
	// Kind Regards,
	// Your Trivial Tickets Team
	//
	// This message was automatically generated by trivial-tickets.com.
	// Please do not reply to this e-mail.
}
