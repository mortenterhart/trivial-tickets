package ticket

import (
	"testing"

	"github.com/mortenterhart/trivial-tickets/structs"
	"github.com/stretchr/testify/assert"
)

func TestCreateTicket(t *testing.T) {

	ticket := CreateTicket("test@example.com", "PC geht ständig aus", "Mein PC geht immer aus.\nIch weiß nicht was los ist.")

	assert.NotNil(t, ticket, "No ticket was returned")
}

func TestUpdateTicket(t *testing.T) {

	ticket := UpdateTicket("2", "abcdef12345", "test@example.com", structs.Ticket{})

	assert.NotNil(t, ticket, "No ticket was returned")
}

func TestMergeTickets(t *testing.T) {

	MergeTickets("abcdef12345", "ghijkl6789")
}

func TestAssignTicket(t *testing.T) {

	ticket := AssignTicket("abcdef12345", "admin")

	assert.NotNil(t, ticket, "No ticket was returned")
}

func TestUnassignTicket(t *testing.T) {

	ticket := UnassignTicket("abcdef12345")

	assert.NotNil(t, ticket, "No ticket was returned")
}

// TestCreateTicketId makes sure the created ticket id is in line with the specification
func TestCreateTicketId(t *testing.T) {

	ticketId := createTicketId(10)

	assert.True(t, (len(ticketId) == 10), "Ticket id has the wrong length")
}
