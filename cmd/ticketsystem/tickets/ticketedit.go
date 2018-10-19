package tickets

import (
	"errors"

	"github.com/mortenterhart/trivial-tickets/structs"
)

func AdoptTicket(ticket *structs.Ticket, editor structs.User) error {
	if ticket.User != (structs.User{}) {
		return errors.New("ticket already has a user assigned")
	}

	ticket.User = editor
	return nil
}

func ReleaseTicket(ticket *structs.Ticket) error {
	if ticket.User == (structs.User{}) {
		return errors.New("ticket has no user assigned, unable to release user")
	}

	ticket.User = (structs.User{})
	return nil
}
