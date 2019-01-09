package ticket

import (
	"github.com/mortenterhart/trivial-tickets/util/random"
	"sort"
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
		Id:       random.CreateRandomId(10),
		Subject:  subject,
		Status:   structs.OPEN,
		User:     structs.User{},
		Customer: mail,
		Entries:  entries,
		MergeTo:  "",
	}
}

// UpdateTicket gets update parameters as well as the ticket to be updated
// and returns it with the values overwritten
func UpdateTicket(status, mail, reply, replyType string, currentTicket structs.Ticket) structs.Ticket {

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
			Reply_Type:    replyType,
		}

		entries := currentTicket.Entries
		entries = append(entries, newEntry)

		currentTicket.Entries = entries
	}

	return currentTicket
}

// MergeTickets merges two tickets if they share the same customer
func MergeTickets(mergeToTicket, mergeFromTicket structs.Ticket) (structs.Ticket, structs.Ticket) {

	if mergeToTicket.Customer == mergeFromTicket.Customer {
		// Get and merge the entries
		entriesMerged := append(mergeFromTicket.Entries, mergeToTicket.Entries...)

		// Sort the merged entries by date from earliest to latest
		sort.Slice(entriesMerged, func(i, j int) bool {
			return entriesMerged[i].Date.Before(entriesMerged[j].Date)
		})

		// Assign the merged entries
		mergeToTicket.Entries = entriesMerged

		// Point to the newly merged ticket
		mergeFromTicket.MergeTo = mergeToTicket.Id
		mergeFromTicket.Status = structs.CLOSED
		mergeFromTicket.User = structs.User{}
	}

	return mergeToTicket, mergeFromTicket
}

// AssignTicket adds a user to a ticket
func AssignTicket(user structs.User, currentTicket structs.Ticket) structs.Ticket {

	// Assign the user to the specified ticket
	// and change the Status
	currentTicket.User = user
	currentTicket.Status = structs.PROCESSING

	return currentTicket
}

// UnassignTicket removes a user from a ticket
func UnassignTicket(currentTicket structs.Ticket) structs.Ticket {

	// Replace the assigned user with an empty struct
	// and set the status to open
	currentTicket.User = structs.User{}
	currentTicket.Status = structs.OPEN

	return currentTicket
}

