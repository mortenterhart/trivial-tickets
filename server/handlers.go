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

// Package server implements the web server including
// shutdown routines and the associated handlers for
// web requests.
package server

import (
	"html/template"
	"net/http"
	"time"

	"github.com/mortenterhart/trivial-tickets/api/api_out"
	"github.com/mortenterhart/trivial-tickets/globals"
	"github.com/mortenterhart/trivial-tickets/log"
	"github.com/mortenterhart/trivial-tickets/mail_events"
	"github.com/mortenterhart/trivial-tickets/session"
	"github.com/mortenterhart/trivial-tickets/structs"
	"github.com/mortenterhart/trivial-tickets/ticket"
	"github.com/mortenterhart/trivial-tickets/util/filehandler"
	"github.com/mortenterhart/trivial-tickets/util/hashing"
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
 * Package server
 * Server handlers reacting to HTTP requests
 */

// getMethod defines the HTTP GET method string.
const getMethod string = "GET"

// postMethod defines the HTTP POST method string.
const postMethod string = "POST"

// indexURL is the base URL of the server that is
// redirected to often.
const indexURL string = "/"

// idParameter is the parameter used by some handlers
// to get the ticket id required for its action.
const idParameter string = "id"

// handleIndex handles the traffic for the index.html
func handleIndex(w http.ResponseWriter, r *http.Request) {

	userSession, errCheckForSession := session.CheckForSession(w, r)

	if errCheckForSession != nil {
		log.Error("Unable to create session:", errCheckForSession)
	}

	executeErr := tmpl.Lookup("index.html").ExecuteTemplate(w, "index",
		structs.Data{
			Session: userSession,
			Tickets: globals.Tickets,
			Users:   users,
		})
	if executeErr != nil {
		log.Error(executeErr)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// handleLogin checks the login credentials against the stored users
// and allows the user access, if their credentials are correct.
func handleLogin(w http.ResponseWriter, r *http.Request) {

	// Get session id
	sessionID := session.GetSessionID(r)

	// Only handle POST-Requests
	if r.Method == postMethod {

		// Get submitted username and password
		username := template.HTMLEscapeString(r.FormValue("username"))
		password := template.HTMLEscapeString(r.FormValue("password"))

		// Get the user with the given username from the hash map
		// Check if the given username and password are correct
		if user, errUser := users[username]; errUser {
			if username == user.Username && hashing.CheckPassword(user.Hash, password) {

				log.Infof("User '%s' (username '%s') logged in successfully", user.Name, username)

				// Create a session to update the current one
				currentSession, _ := session.GetSession(sessionID)
				currentSession.User = user
				currentSession.IsLoggedIn = true
				currentSession.CreationTime = time.Now()

				// Update the session with the one just created
				session.UpdateSession(sessionID, currentSession)
			}
		}
	}

	// Redirect the user to the index
	http.Redirect(w, r, indexURL, http.StatusMovedPermanently)
}

// handleLogout logs the user out and clears their session.
func handleLogout(w http.ResponseWriter, r *http.Request) {

	// Get session id
	sessionID := session.GetSessionID(r)

	if r.Method == postMethod {

		user := globals.Sessions[sessionID].Session.User

		// Remove the session of the user
		delete(globals.Sessions, sessionID)

		// Delete the session cookie
		http.SetCookie(w, session.DeleteSessionCookie())

		log.Infof("User '%s' (username '%s') logged out now", user.Name, user.Username)
	}

	// Redirect the user to the index
	http.Redirect(w, r, indexURL, http.StatusMovedPermanently)
}

// handleCreateTicket creates a new ticket struct and saves it.
func handleCreateTicket(w http.ResponseWriter, r *http.Request) {

	// Only react on POST request
	if r.Method == postMethod {

		// Get the form values
		mail := template.HTMLEscapeString(r.FormValue("mail"))
		subject := template.HTMLEscapeString(r.FormValue("subject"))
		text := template.HTMLEscapeString(r.FormValue("text"))

		// Create the ticket
		newTicket := ticket.CreateTicket(mail, subject, text)
		log.Infof(`Creating new ticket '%s' for customer '%s' with subject "%s"`,
			newTicket.ID, newTicket.Customer, newTicket.Subject)

		// Assign the ticket to the tickets kept in memory
		globals.Tickets[newTicket.ID] = newTicket

		// Persist the ticket to the file system
		filehandler.WriteTicketFile(globals.ServerConfig.Tickets, &newTicket)

		// Send notification mail on create ticket event
		api_out.SendMail(mail_events.NewTicket, newTicket)

		// Redirect the user to the ticket page
		http.Redirect(w, r, "/ticket?id="+newTicket.ID, http.StatusMovedPermanently)

		return
	}

	// If there is any other request, just redirect to index
	http.Redirect(w, r, indexURL, http.StatusMovedPermanently)
}

// handleHoliday activates / deactivates the holiday mode for a given user.
func handleHoliday(w http.ResponseWriter, r *http.Request) {

	// Get session id
	sessionID := session.GetSessionID(r)

	// Make sure user is logged in
	if globals.Sessions[sessionID].Session.IsLoggedIn {

		// Create a session to update the current one
		currentSession, _ := session.GetSession(sessionID)

		// Get the current user
		user := users[currentSession.User.Username]

		// Toggle IsOnHoliday
		if currentSession.User.IsOnHoliday {
			currentSession.User.IsOnHoliday, user.IsOnHoliday = false, false
		} else {
			currentSession.User.IsOnHoliday, user.IsOnHoliday = true, true
		}

		log.Infof("Updating the holiday setting for user '%s' (username '%s') to %t",
			user.Name, user.Username, user.IsOnHoliday)

		// Update the session with the one just created
		session.UpdateSession(sessionID, currentSession)

		// Update the users hash map
		users[currentSession.User.Username] = user

		// Persist the changes to the file system
		filehandler.WriteUserFile(globals.ServerConfig.Users, &users)
	}

	// Redirect the user to the index
	http.Redirect(w, r, indexURL, http.StatusMovedPermanently)
}

// handleTicket gets the requested ticket via the url GET
// parameters and serves it to the template.
func handleTicket(w http.ResponseWriter, r *http.Request) {

	// Only support GET request
	if r.Method == getMethod {

		// Extract the id url parameter
		idParam, idParamDefined := r.URL.Query()[idParameter]

		if !idParamDefined {
			log.Errorf("%s %s: missing parameter '%s'", r.Method, r.RequestURI, idParameter)
			http.Redirect(w, r, indexURL, http.StatusMovedPermanently)
			return
		}

		// Get the ticket based on the given id
		ticketID := idParam[0]
		ticket := globals.Tickets[ticketID]

		// If it is a merged ticket, redirect to the merged one
		if ticket.MergeTo != "" {
			ticket = globals.Tickets[ticket.MergeTo]
		}

		// Create or get the users session
		currentSession, errCheckForSession := session.CheckForSession(w, r)

		if errCheckForSession != nil {
			log.Error("Unable to get session")
		}

		// Serve the template to show a single ticket
		executeErr := tmpl.Lookup("ticket.html").ExecuteTemplate(w, "ticket",
			structs.DataSingleTicket{Session: currentSession, Ticket: ticket, Tickets: globals.Tickets, Users: users})
		if executeErr != nil {
			log.Error(executeErr)
			w.WriteHeader(http.StatusInternalServerError)
		}

		return
	}

	http.Redirect(w, r, indexURL, http.StatusMovedPermanently)
}

// handleUpdateTicket gets the requested ticket via the url GET
// parameters and serves it to the template.
func handleUpdateTicket(w http.ResponseWriter, r *http.Request) {

	// Only support POST request
	if r.Method == postMethod {

		// Get the session
		currentSession, errCheckForSession := session.CheckForSession(w, r)

		if errCheckForSession != nil {
			log.Error("Unable to get session:", errCheckForSession)
			return
		}

		// Get form values
		ticketID := template.HTMLEscapeString(r.FormValue("ticket"))
		status := template.HTMLEscapeString(r.FormValue("status"))
		mail := template.HTMLEscapeString(r.FormValue("mail"))
		reply := template.HTMLEscapeString(r.FormValue("reply"))
		replyType := template.HTMLEscapeString(r.FormValue("reply_type"))
		merge := template.HTMLEscapeString(r.FormValue("merge"))

		// Get the ticket which was edited
		currentTicket := globals.Tickets[ticketID]

		// Update the current ticket
		updatedTicket := ticket.UpdateTicket(status, mail, reply, replyType, currentTicket)

		if merge != "" {
			// Get the ticket to merge from the tickets map
			ticketFrom := globals.Tickets[merge]

			// Only if they have the same assigned user
			if ticketFrom.User == currentSession.User && updatedTicket.User == currentSession.User {

				// Merge structs.Ticket
				ticketMergedTo, ticketMergedFrom := ticket.MergeTickets(updatedTicket, ticketFrom)

				log.Infof("Merging ticket '%s' to ticket '%s' and saving to file system",
					ticketMergedFrom.ID, ticketMergedTo.ID)

				// Write both tickets to memory
				globals.Tickets[ticketMergedTo.ID] = ticketMergedTo
				globals.Tickets[ticketMergedFrom.ID] = ticketMergedFrom

				// Persist both tickets to file system
				filehandler.WriteTicketFile(globals.ServerConfig.Tickets, &ticketMergedTo)
				filehandler.WriteTicketFile(globals.ServerConfig.Tickets, &ticketMergedFrom)

				// Update to the merged ticket so serve to client
				updatedTicket = globals.Tickets[ticketMergedTo.ID]
			}
		} else {

			log.Infof("Updating ticket '%s' with status '%s' and %d answers", updatedTicket.ID,
				updatedTicket.Status.String(), len(updatedTicket.Entries))
			// Assign the updated ticket to the ticket map in memory
			globals.Tickets[ticketID] = updatedTicket

			// Persist the updated ticket to the file system
			filehandler.WriteTicketFile(globals.ServerConfig.Tickets, &updatedTicket)
		}

		if !currentSession.IsLoggedIn {
			replyType = "external"
		}

		// Send mail if the reply was selected for external
		if replyType == "external" {
			mailEvent := mail_events.UpdatedTicket
			if reply != "" {
				mailEvent = mail_events.NewAnswer
			}

			api_out.SendMail(mailEvent, updatedTicket)
		}

		// Redirect to the ticket again, now with updated Values
		executeErr := tmpl.Lookup("ticket.html").ExecuteTemplate(w, "ticket",
			structs.DataSingleTicket{Session: currentSession, Ticket: updatedTicket, Tickets: globals.Tickets})
		if executeErr != nil {
			log.Error(executeErr)
			w.WriteHeader(http.StatusInternalServerError)
		}

		return
	}

	http.Redirect(w, r, indexURL, http.StatusMovedPermanently)
}

// handleAssignTicket assigns a given user to a given ticket and
// returns the user name of the newly assigned user to the browser.
func handleAssignTicket(w http.ResponseWriter, r *http.Request) {

	if r.Method == getMethod {

		// Get the session
		currentSession, errCheckForSession := session.CheckForSession(w, r)

		if errCheckForSession != nil {
			log.Error("Unable to get session: ", errCheckForSession)
			return
		}

		if currentSession.IsLoggedIn {

			// Extract the GET request parameters
			params := r.URL.Query()

			ticketID := params["id"][0]
			user := params["user"][0]

			// Get the ticket based on the given id
			currentTicket := globals.Tickets[ticketID]

			// Update the ticket itself
			updatedTicket := ticket.AssignTicket(users[user], currentTicket)

			log.Infof("Assigning user '%s' (username '%s') to ticket '%s'",
				updatedTicket.User.Name, updatedTicket.User.Username, updatedTicket.ID)

			// Update the ticket in memory
			globals.Tickets[ticketID] = updatedTicket

			// Persist the change in the file system
			filehandler.WriteTicketFile(globals.ServerConfig.Tickets, &updatedTicket)

			// Return the assigned user
			response := updatedTicket.User.Username
			w.Header().Set("Content-Type", "text/html")
			w.Write([]byte(response))

			api_out.SendMail(mail_events.AssignedTicket, updatedTicket)
		}
	}
}

// handleUnassignTicket unassigns a ticket from a certain user,
// only if the actual user makes the request. Other users are
// unable to unassign a ticket from anyone apart from themselves.
func handleUnassignTicket(w http.ResponseWriter, r *http.Request) {

	if r.Method == getMethod {

		// Extract the GET request parameters
		idParam, idParamDefined := r.URL.Query()[idParameter]

		if !idParamDefined {
			log.Errorf("%s %s: missing parameter '%s'", r.Method, r.RequestURI, idParameter)
			http.Redirect(w, r, indexURL, http.StatusMovedPermanently)
			return
		}

		// Get the ticket based on the given id
		ticketID := idParam[0]
		currentTicket := globals.Tickets[ticketID]

		// Get the session
		currentSession, errCheckForSession := session.CheckForSession(w, r)

		if errCheckForSession != nil {
			log.Error("Unable to get session:", errCheckForSession)
			return
		}

		// Make sure the requesting user owns the ticket
		if currentSession.User.ID == currentTicket.User.ID {

			log.Infof("Unassigning user '%s' (username '%s') from ticket '%s'",
				currentTicket.User.Name, currentTicket.User.Username, currentTicket.ID)

			// Replace the assigned user with nobody
			updatedTicket := ticket.UnassignTicket(currentTicket)

			// Set the ticket to memory
			globals.Tickets[ticketID] = updatedTicket

			// Persist the changed ticket to the file system
			filehandler.WriteTicketFile(globals.ServerConfig.Tickets, &updatedTicket)

			// Create a response and write it to the header
			response := "The Ticket was released successfully."
			w.Header().Set("Content-Type", "text/html")
			w.Write([]byte(response))

			api_out.SendMail(mail_events.UnassignedTicket, currentTicket)
		}
	}
}
