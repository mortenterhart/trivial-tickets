// Administration of ticket actions
package ticket

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

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

	const MAIL = "test@example.com"
	const SUBJECT = "PC geht ständig aus"
	const ENTRY = "Mein PC geht immer aus.\nIch weiß nicht was los ist."

	ticket := CreateTicket(MAIL, SUBJECT, ENTRY)

	assert.NotNil(t, ticket, "No ticket was returned")
	assert.Equal(t, MAIL, ticket.Customer, "Mail in created ticket did not match")
	assert.Equal(t, SUBJECT, ticket.Subject, "Subject does not match")
}

func TestUpdateTicket(t *testing.T) {

	const STATUS = "2"
	const TICKET_ID = "abcdef12345"
	const MAIL = "text@exmaple.com"
	const REPLY_TYPE = ""

	ticket := UpdateTicket(STATUS, TICKET_ID, MAIL, REPLY_TYPE, structs.Ticket{})

	assert.NotNil(t, ticket, "No ticket was returned")
	assert.Equal(t, structs.CLOSED, ticket.Status, "Status does not match")
}

// TestMergeTickets makes sure that the entries of merged tickets are combined and that the
// ticket, where it is matched from has the id of the merged to ticket
func TestMergeTickets(t *testing.T) {

	// Create Mock entries
	entry := structs.Entry{
		Date:          time.Now(),
		FormattedDate: time.Now().Format(time.ANSIC),
		User:          "example@example.com",
		Text:          "hallo mir gehts gut und dir?\n mfg",
	}

	var entries []structs.Entry

	for i := 0; i < 3; i++ {
		entries = append(entries, entry)
	}

	// Create mock tickets
	ticketMergeTo := structs.Ticket{Id: "abcdef123"}
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
