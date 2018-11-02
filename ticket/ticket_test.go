package ticket

import (
	"testing"

	"github.com/mortenterhart/trivial-tickets/structs"
	"github.com/stretchr/testify/assert"
)

func TestCreateTicket(t *testing.T) {

	const MAIL = "test@example.com"
	const SUBJECT = "PC geht ständig aus"
	const ENTRY = "Mein PC geht immer aus.\nIch weiß nicht was los ist."

	ticket := CreateTicket(MAIL, SUBJECT, ENTRY)

	assert.NotNil(t, ticket, "No ticket was returned")
	assert.Equal(t, ticket.Customer, MAIL, "Mail in created ticket did not match")
	assert.Equal(t, ticket.Subject, SUBJECT)
}

func TestUpdateTicket(t *testing.T) {

	const STATUS = "2"
	const TICKET_ID = "abcdef12345"
	const MAIL = "text@exmaple.com"

	ticket := UpdateTicket(STATUS, TICKET_ID, MAIL, structs.Ticket{})

	assert.NotNil(t, ticket, "No ticket was returned")
	assert.Equal(t, ticket.Status, structs.CLOSED)
}

func TestMergeTickets(t *testing.T) {

	MergeTickets("abcdef12345", "ghijkl6789")
}

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

// TestCreateTicketId makes sure the created ticket id is in line with the specification
func TestCreateTicketId(t *testing.T) {

	ticketId := createTicketId(10)

	assert.True(t, (len(ticketId) == 10), "Ticket id has the wrong length")
}
