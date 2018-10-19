package tickets

import (
	"fmt"
	"testing"

	"github.com/mortenterhart/trivial-tickets/structs"

	"github.com/stretchr/testify/assert"
)

func TestAdoptTicket(t *testing.T) {
	ticket := structs.Ticket{
		Id:       1,
		Subject:  "user should adopt a ticket",
		Status:   structs.OPEN,
		User:     structs.User{},
		Customer: "",
		Entries:  nil,
	}

	editor := structs.User{
		Id:          1,
		Name:        "Max Mustermann",
		Mail:        "max.mustermann@email.org",
		Hash:        "hashing Wert",
		IsOnHoliday: false,
	}

	AdoptTicket(&ticket, editor)

	assert.Equal(t, editor, ticket.User, "user in ticket should be the same "+
		"than the current logged in user")

	assert.Error(t, AdoptTicket(&ticket, editor), "tickets with editors should "+
		" not be assigned a new editor")
}

func TestReleaseTicket(t *testing.T) {
	ticket := structs.Ticket{
		Id:      1,
		Subject: "user can release ticket",
		Status:  structs.OPEN,
		User: structs.User{
			Id:   1,
			Name: "Max Mustermann",
			Mail: "max.mustermann@",
		},
	}

	fmt.Println(ticket)
}
