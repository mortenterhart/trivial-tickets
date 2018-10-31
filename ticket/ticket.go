package ticket

import (
	"math/rand"
	"strconv"
	"time"

	"github.com/mortenterhart/trivial-tickets/structs"
)

// CreateTicket takes the arguments from either web or the mail api and returns a populated ticket
func CreateTicket(mail, subject, text string) structs.Ticket {

	// Create a new entry for the ticket
	entry := structs.Entry{
		Date:          time.Now(),
		FormattedDate: time.Now().Format(time.ANSIC),
		User:          mail,
		Text:          text,
	}

	var entries []structs.Entry
	entries = append(entries, entry)

	// Construct the ticket
	return structs.Ticket{
		Id:       createTicketId(10),
		Subject:  subject,
		Status:   structs.OPEN,
		User:     structs.User{},
		Customer: mail,
		Entries:  entries,
	}
}

func UpdateTicket(status, mail, reply string, currentTicket structs.Ticket) structs.Ticket {

	// Set the status to the one provided by the form
	statusValue, _ := strconv.Atoi(status)
	currentTicket.Status = structs.State(statusValue)

	// If there has been a reply, attach it to the entries slice of the ticket
	if reply != "" {

		newEntry := structs.Entry{
			Date:          time.Now(),
			FormattedDate: time.Now().Format(time.ANSIC),
			User:          mail,
			Text:          reply,
		}

		entries := currentTicket.Entries
		entries = append(entries, newEntry)
		currentTicket.Entries = entries
	}

	return currentTicket
}

func MergeTickets(mergeToTicketId, mergeFromTicketId string) {

}

func AssignTicket(ticketId, userName string) structs.Ticket {

	return structs.Ticket{}
}

func UnassignTicket(ticketId string) structs.Ticket {

	return structs.Ticket{}
}

// letters are the valid characters for the ticket id
var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

// createTicketId generates a pseudo random id for the tickets
// Tweaked example from https://stackoverflow.com/a/22892986
func createTicketId(n int) string {

	// Seed the random function to make it more random
	rand.Seed(time.Now().UnixNano())

	// Create a slice, big enough to hold the id
	b := make([]rune, n)

	// Randomly choose a letter from the letters rune
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}

	return string(b)
}
