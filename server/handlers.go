package server

import (
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/mortenterhart/trivial-tickets/globals"
	"github.com/mortenterhart/trivial-tickets/structs"
	"github.com/mortenterhart/trivial-tickets/ticket"
	"github.com/mortenterhart/trivial-tickets/util/filehandler"
	"github.com/mortenterhart/trivial-tickets/util/hashing"
)

/*
*
* Matrikelnummern
* 3040018
* 6694964
* 3478222
 */

// handleIndex handles the traffic for the index.html
func handleIndex(w http.ResponseWriter, r *http.Request) {

	session := checkForSession(w, r)

	tmpl.Lookup("index.html").ExecuteTemplate(w, "index", structs.Data{Session: session, Tickets: globals.Tickets, Users: users})
}

// handleLogin checks the login credentials against the stored users
// and allows the user access, if their credentials are correct
func handleLogin(w http.ResponseWriter, r *http.Request) {

	// Get session id
	sessionId := getSessionId(r)

	// Only handle POST-Requests
	if r.Method == "POST" {

		// Get submitted form values
		username := template.HTMLEscapeString(r.FormValue("username"))
		password := template.HTMLEscapeString(r.FormValue("password"))

		// Get the user with the given username from the hashmap
		// Check if the given username and password are correct
		if user, errUser := users[username]; errUser {
			if username == user.Username && hashing.CheckPassword(user.Hash, password) {

				// Create a session to update the current one
				session, _ := GetSession(sessionId)
				session.User = user
				session.IsLoggedIn = true
				session.CreateTime = time.Now()

				// Update the session with the one just created
				UpdateSession(sessionId, session)
			} else {
				// TODO: Provide error of wrong login credentials
			}
		} else {
			// TODO: Provide error of wrong login credentials
		}
	}

	// Redirect the user to the index
	http.Redirect(w, r, "/", 302)
}

// handleLogout logs the user out and clears their session
func handleLogout(w http.ResponseWriter, r *http.Request) {

	// Get session id
	sessionId := getSessionId(r)

	if r.Method == "POST" {

		// Remove the session of the user
		delete(sessions, sessionId)

		// Delete the session cookie
		http.SetCookie(w, deleteSessionCookie())
	}

	// Redirect the user to the index
	http.Redirect(w, r, "/", 302)
}

// handleCreateTicket creates a new ticket struct and saves it
func handleCreateTicket(w http.ResponseWriter, r *http.Request) {

	// Only react on POST request
	if r.Method == "POST" {

		// Get the form values
		mail := template.HTMLEscapeString(r.FormValue("mail"))
		subject := template.HTMLEscapeString(r.FormValue("subject"))
		text := template.HTMLEscapeString(r.FormValue("text"))

		// Create the ticket
		newTicket := ticket.CreateTicket(mail, subject, text)

		// Assign the ticket to the tickets kept in memory
		globals.Tickets[newTicket.Id] = newTicket

		// Persist the ticket to the file system
		filehandler.WriteTicketFile(globals.ServerConfig.Tickets, &newTicket)

		// Redirect the user to the ticket page
		http.Redirect(w, r, "/ticket?id="+newTicket.Id, 302)
	}

	// If there is any other request, just redirect to index
	http.Redirect(w, r, "/", 302)
}

// handleHoliday activates / deactivates the holiday mode for a given user
func handleHoliday(w http.ResponseWriter, r *http.Request) {

	// Get session id
	sessionId := getSessionId(r)

	// Make sure user is logged in
	if sessions[sessionId].Session.IsLoggedIn {

		// Create a session to update the current one
		session, _ := GetSession(sessionId)

		// Get the current user
		user := users[session.User.Username]

		// Toggle IsOnHoliday
		if session.User.IsOnHoliday {
			session.User.IsOnHoliday, user.IsOnHoliday = false, false
		} else {
			session.User.IsOnHoliday, user.IsOnHoliday = true, true
		}

		// Update the session with the one just created
		UpdateSession(sessionId, session)

		// Update the users hash map
		users[session.User.Username] = user

		// Persist the changes to the file system
		filehandler.WriteUserFile(globals.ServerConfig.Users, &users)
	}

	// Redirect the user to the index
	http.Redirect(w, r, "/", 302)
}

// handleTicket gets the requested ticket via the url GET parameters and serves it to the template
func handleTicket(w http.ResponseWriter, r *http.Request) {

	// Only support GET request
	if r.Method == "GET" {

		// Extract the id url parameter
		param, errParam := r.URL.Query()["id"]

		if !errParam || len(param[0]) < 1 {
			log.Fatal(errParam)
		}

		// Get the ticket based on the given id
		id := param[0]
		ticket := globals.Tickets[id]

		// Create or get the users session
		session := checkForSession(w, r)

		// Serve the template to show a single ticket
		tmpl.Lookup("ticket.html").ExecuteTemplate(w, "ticket", structs.DataSingleTicket{Session: session, Ticket: ticket})
	}

	http.Redirect(w, r, "/", 302)
}

// handleUpdateTicket gets the requested ticket via the url GET parameters and serves it to the template
func handleUpdateTicket(w http.ResponseWriter, r *http.Request) {

	// Only support POST request
	if r.Method == "POST" {

		// Get the session
		session := checkForSession(w, r)

		// Get form values
		ticketId := template.HTMLEscapeString(r.FormValue("ticket"))
		status := template.HTMLEscapeString(r.FormValue("status"))
		mail := template.HTMLEscapeString(r.FormValue("mail"))
		reply := template.HTMLEscapeString(r.FormValue("reply"))

		// Get the ticket which was edited
		currentTicket := globals.Tickets[ticketId]

		// Update the current ticket
		updatedTicket := ticket.UpdateTicket(status, mail, reply, currentTicket)

		// Assign the updated ticket to the ticket map in memory
		globals.Tickets[ticketId] = updatedTicket

		// Persist the updated ticket to the file system
		filehandler.WriteTicketFile(globals.ServerConfig.Tickets, &updatedTicket)

		// Redirect to the ticket again, now with updated Values
		tmpl.Lookup("ticket.html").ExecuteTemplate(w, "ticket", structs.DataSingleTicket{Session: session, Ticket: updatedTicket})
	}

	http.Redirect(w, r, "/", 302)
}

// handleUnassignTicket unassigns a ticket from a certain user, only if the actual user makes the request.
// Other users are unable to unassign a ticket from anyone apart from themselves.
func handleUnassignTicket(w http.ResponseWriter, r *http.Request) {

	if r.Method == "GET" {

		// Extract the GET request parameters
		param, errParam := r.URL.Query()["id"]

		if !errParam || len(param[0]) < 1 {
			log.Fatal(errParam)
		}

		// Get the ticket based on the given id
		ticketId := param[0]
		currentTicket := globals.Tickets[ticketId]

		// Get the session
		session := checkForSession(w, r)

		// Make sure, the requesting user owns the ticket
		if session.User.Id == currentTicket.User.Id {

			// Replace the assigned user with nobody
			updatedTicket := ticket.UnassignTicket(currentTicket)

			// Set the ticket to memory
			globals.Tickets[ticketId] = updatedTicket

			// Persist the changed ticket to the file system
			filehandler.WriteTicketFile(globals.ServerConfig.Tickets, &updatedTicket)

			// Create a response and write it to the header
			response := "Das Ticket wurde erfolgreich freigegeben"
			w.Header().Set("Content-Type", "text/html")
			w.Write([]byte(response))
		}
	}
}

// handleAssignTicket assigns a given user to a given ticket and returns the user name
// of the newly assigned user to the browser
func handleAssignTicket(w http.ResponseWriter, r *http.Request) {

	if r.Method == "GET" {

		// Get the session
		session := checkForSession(w, r)

		if session.IsLoggedIn {

			// Extract the GET request parameters
			params := r.URL.Query()

			ticketId := params["id"][0]
			user := params["user"][0]

			// Get the ticket based on the given id
			currentTicket := globals.Tickets[ticketId]

			// Update the ticket itself
			updatedTicket := ticket.AssignTicket(users[user], currentTicket)

			// Update the ticket in memory
			globals.Tickets[ticketId] = updatedTicket

			// Persist the change in the file system
			filehandler.WriteTicketFile(globals.ServerConfig.Tickets, &updatedTicket)

			// Return the assigned user
			response := updatedTicket.User.Username
			w.Header().Set("Content-Type", "text/html")
			w.Write([]byte(response))
		}
	}
}

// handleMergeTickets merges two tickets, writes them to memory and persist them.
// Then, the newly merged ticket is returned to the view
func handleMergeTickets(w http.ResponseWriter, r *http.Request) {

	// Only support POST
	if r.Method == "POST" {

		// Create or get the users session
		session := checkForSession(w, r)

		// Make sure user is logged in
		if session.IsLoggedIn {

			// Get both ticket ids
			ticketIdMergeTo := globals.Tickets[template.HTMLEscapeString(r.FormValue("merge_to"))]
			ticketIdMergeFrom := globals.Tickets[template.HTMLEscapeString(r.FormValue("merge_from"))]

			// Only if they have the same assigned user
			if ticketIdMergeFrom.User == session.User && ticketIdMergeTo.User == session.User {

				// Merge structs.Ticket
				ticketMergedTo, ticketMergedFrom := ticket.MergeTickets(ticketIdMergeTo, ticketIdMergeFrom)

				// Write both tickets to memory
				globals.Tickets[ticketMergedTo.Id] = ticketMergedTo
				globals.Tickets[ticketMergedFrom.Id] = ticketMergedFrom

				// Persist both tickets to file system
				filehandler.WriteTicketFile(globals.ServerConfig.Tickets, &ticketMergedTo)
				filehandler.WriteTicketFile(globals.ServerConfig.Tickets, &ticketMergedFrom)

				// Serve the template to show a single ticket
				tmpl.Lookup("ticket.html").ExecuteTemplate(w, "ticket", structs.DataSingleTicket{Session: session, Ticket: ticketMergedTo})
			}
		}
	}

	http.Redirect(w, r, "/", 302)
}

// createSessionCookie returns a http cookie to hold the session
// id for the user
func createSessionCookie() (*http.Cookie, string) {

	sessionId := CreateSessionId()

	return &http.Cookie{
		Name:     "session",
		Value:    sessionId,
		HttpOnly: false,
		Expires:  time.Now().Add(2 * time.Hour)}, sessionId
}

// deleteSessionCookie returns a http cookie which will overwrite the
// existing session cookie in order to nulify it
func deleteSessionCookie() *http.Cookie {

	return &http.Cookie{
		Name:     "session",
		Value:    "",
		HttpOnly: false,
		Expires:  time.Now().Add(-100 * time.Hour)}
}

// getSessionId retrieves the session id from the cookie of the user
func getSessionId(r *http.Request) string {

	// Get the cookie with the session id
	userCookie, errUserCookie := r.Cookie("session")

	if errUserCookie != nil {
		log.Print(errUserCookie)
		return errUserCookie.Error()
	}

	return userCookie.Value
}

// checkForSession either returns a new session or the existing session of a user.
func checkForSession(w http.ResponseWriter, r *http.Request) structs.Session {

	var session structs.Session

	// Check if the user already has a session
	// If not, create one
	// Otherwise read the session id and load the index with his session
	if _, err := r.Cookie("session"); err != nil {

		cookie, sessionId := createSessionCookie()
		http.SetCookie(w, cookie)
		sessions[sessionId] = CreateSession(sessionId)

		session = sessions[sessionId].Session

	} else {
		sessionId := getSessionId(r)

		session = sessions[sessionId].Session
	}

	return session
}
