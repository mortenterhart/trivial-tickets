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
	"sort"
	"strconv"
	"time"

	"github.com/mortenterhart/trivial-tickets/structs"
	"github.com/mortenterhart/trivial-tickets/util/random"
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
 * Package ticket
 * Administration of ticket actions
 */

// CreateTicket takes the arguments from either web or
// the Mail API and returns a populated ticket.
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
		ID:       random.CreateRandomID(structs.RandomIDLength),
		Subject:  subject,
		Status:   structs.StatusOpen,
		User:     structs.User{},
		Customer: mail,
		Entries:  entries,
		MergeTo:  "",
	}
}

// UpdateTicket gets update parameters as well as the
// ticket to be updated and returns it with the values
// overwritten.
func UpdateTicket(status, mail, reply, replyType string, currentTicket structs.Ticket) structs.Ticket {

	// Set the status to the one provided by the form
	statusValue, _ := strconv.Atoi(status)
	currentTicket.Status = structs.Status(statusValue)

	// If there has been a reply, attach it to the entries slice of the ticket
	if reply != "" {

		newEntry := structs.Entry{
			Date:          time.Now(),
			FormattedDate: time.Now().Format(time.ANSIC),
			User:          mail,
			Text:          reply,
			ReplyType:     replyType,
		}

		entries := currentTicket.Entries
		entries = append(entries, newEntry)

		currentTicket.Entries = entries
	}

	return currentTicket
}

// MergeTickets merges two tickets if they share the
// same customer. The entries of the mergeFromTicket
// are appended to the mergeToTicket and all entries
// are sorted by their creation time.
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
		mergeFromTicket.MergeTo = mergeToTicket.ID
		mergeFromTicket.Status = structs.StatusClosed
		mergeFromTicket.User = structs.User{}
	}

	return mergeToTicket, mergeFromTicket
}

// AssignTicket adds a user to a ticket.
func AssignTicket(user structs.User, currentTicket structs.Ticket) structs.Ticket {

	// Assign the user to the specified ticket
	// and change the Status
	currentTicket.User = user
	currentTicket.Status = structs.StatusInProgress

	return currentTicket
}

// UnassignTicket removes a user from a ticket.
func UnassignTicket(currentTicket structs.Ticket) structs.Ticket {

	// Replace the assigned user with an empty struct
	// and set the status to open
	currentTicket.User = structs.User{}
	currentTicket.Status = structs.StatusOpen

	return currentTicket
}
