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

// Package ticket contains operations for the administration
// of ticket actions and updates.
package ticket

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/mortenterhart/trivial-tickets/logger/testlogger"
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
 * Package ticket [tests]
 * Administration of ticket actions
 */

func TestCreateTicket(t *testing.T) {
	testlogger.BeginTest()
	defer testlogger.EndTest()

	const mail = "test@example.com"
	const subject = "PC is always turning off"
	const entry = "My computer always turns off.\nI don't know what's the case there"

	ticket := CreateTicket(mail, subject, entry)

	assert.NotNil(t, ticket, "No ticket was returned")
	assert.Equal(t, mail, ticket.Customer, "Mail in created ticket did not match")
	assert.Equal(t, subject, ticket.Subject, "Subject does not match")
}

func TestUpdateTicket(t *testing.T) {
	testlogger.BeginTest()
	defer testlogger.EndTest()

	const status = "2"
	const ticketID = "abcdef12345"
	const mail = "text@exmaple.com"
	const replyType = ""

	ticket := UpdateTicket(status, ticketID, mail, replyType, structs.Ticket{})

	assert.NotNil(t, ticket, "No ticket was returned")
	assert.Equal(t, structs.CLOSED, ticket.Status, "Status does not match")
}

// TestMergeTickets makes sure that the entries of merged tickets are combined and that the
// ticket, where it is matched from has the id of the merged to ticket
func TestMergeTickets(t *testing.T) {
	testlogger.BeginTest()
	defer testlogger.EndTest()

	// Create Mock entries
	entry := structs.Entry{
		Date:          time.Now(),
		FormattedDate: time.Now().Format(time.ANSIC),
		User:          "example@example.com",
		Text:          "hello I am fine and you?\n Kind regards",
	}

	var entries []structs.Entry

	for i := 0; i < 3; i++ {
		entries = append(entries, entry)
	}

	// Create mock tickets
	ticketMergeTo := structs.Ticket{ID: "abcdef123"}
	ticketMergeFrom := structs.Ticket{}
	ticketMergeTo.Entries = entries
	ticketMergeFrom.Entries = entries

	// Merge the tickets
	ticketMergeToAfterMerge, ticketMergeFromAfterMerge := MergeTickets(ticketMergeTo, ticketMergeFrom)

	assert.NotNil(t, ticketMergeFromAfterMerge, "No ticket was returned")
	assert.NotNil(t, ticketMergeToAfterMerge, "No ticket was returned")
	assert.True(t, len(ticketMergeToAfterMerge.Entries) == 6, "The entries have not been added to the ticket")
	assert.Equal(t, "abcdef123", ticketMergeFromAfterMerge.MergeTo, "Merge to id does not match")
}

// TestAssignAndUnassignTicket tests that assign and unassigning a ticket works properly
func TestAssignAndUnassignTicket(t *testing.T) {
	testlogger.BeginTest()
	defer testlogger.EndTest()

	// Test assigning the ticket
	user := structs.User{Username: "abcdef"}
	ticket := structs.Ticket{}

	updatedTicket := AssignTicket(user, ticket)

	assert.NotNil(t, updatedTicket, "No ticket was returned")
	assert.Equal(t, "abcdef", updatedTicket.User.Username, "The assigned username does not match")
	assert.Equal(t, structs.PROCESSING, updatedTicket.Status, "The updated ticket has the wrong status")

	// Test unassigning the ticket
	updatedTicket2 := UnassignTicket(updatedTicket)

	assert.Equal(t, structs.OPEN, updatedTicket2.Status, "Status of unassigned ticket is not OPEN")
}
