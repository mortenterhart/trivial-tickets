package server

import (
    "html/template"
    "log"
    "net/http"
    "time"

    "github.com/mortenterhart/trivial-tickets/api/api_out"
    "github.com/mortenterhart/trivial-tickets/globals"
    "github.com/mortenterhart/trivial-tickets/session"
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

    session, errCheckForSession := session.CheckForSession(w, r)

    if errCheckForSession != nil {
        log.Print("Unable to create session")
    }

    tmpl.Lookup("index.html").ExecuteTemplate(w, "index", structs.Data{Session: session, Tickets: globals.Tickets, Users: users})
}

// handleLogin checks the login credentials against the stored users
// and allows the user access, if their credentials are correct
func handleLogin(w http.ResponseWriter, r *http.Request) {

    // Get session id
    sessionId := session.GetSessionId(r)

    // Only handle POST-Requests
    if r.Method == "POST" {

        // Get submitted username
        username := template.HTMLEscapeString(r.FormValue("username"))

        // Get the user with the given username from the hashmap
        // Check if the given username and password are correct
        if user, errUser := users[username]; errUser {
            if username == user.Username && hashing.CheckPassword(user.Hash, template.HTMLEscapeString(r.FormValue("password"))) {

                // Create a session to update the current one
                currentSession, _ := session.GetSession(sessionId)
                currentSession.User = user
                currentSession.IsLoggedIn = true
                currentSession.CreateTime = time.Now()

                // Update the session with the one just created
                session.UpdateSession(sessionId, currentSession)
            }
        }
    }

    // Redirect the user to the index
    http.Redirect(w, r, "/", http.StatusFound)
}

// handleLogout logs the user out and clears their session
func handleLogout(w http.ResponseWriter, r *http.Request) {

    // Get session id
    sessionId := session.GetSessionId(r)

    if r.Method == "POST" {

        // Remove the session of the user
        delete(globals.Sessions, sessionId)

        // Delete the session cookie
        http.SetCookie(w, session.DeleteSessionCookie())
    }

    // Redirect the user to the index
    http.Redirect(w, r, "/", http.StatusFound)
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
        http.Redirect(w, r, "/ticket?id="+newTicket.Id, http.StatusFound)

        return
    }

    // If there is any other request, just redirect to index
    http.Redirect(w, r, "/", http.StatusFound)
}

// handleHoliday activates / deactivates the holiday mode for a given user
func handleHoliday(w http.ResponseWriter, r *http.Request) {

    // Get session id
    sessionId := session.GetSessionId(r)

    // Make sure user is logged in
    if globals.Sessions[sessionId].Session.IsLoggedIn {

        // Create a session to update the current one
        currentSession, _ := session.GetSession(sessionId)

        // Get the current user
        user := users[currentSession.User.Username]

        // Toggle IsOnHoliday
        if currentSession.User.IsOnHoliday {
            currentSession.User.IsOnHoliday, user.IsOnHoliday = false, false
        } else {
            currentSession.User.IsOnHoliday, user.IsOnHoliday = true, true
        }

        // Update the session with the one just created
        session.UpdateSession(sessionId, currentSession)

        // Update the users hash map
        users[currentSession.User.Username] = user

        // Persist the changes to the file system
        filehandler.WriteUserFile(globals.ServerConfig.Users, &users)
    }

    // Redirect the user to the index
    http.Redirect(w, r, "/", http.StatusFound)
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

        // If it is a merged ticket, redirect to the merged one
        if ticket.MergeTo != "" {
            ticket = globals.Tickets[ticket.MergeTo]
        }

        // Create or get the users session
        currentSession, errCheckForSession := session.CheckForSession(w, r)

        if errCheckForSession != nil {
            log.Println("Unable to get session")
        }

        // Serve the template to show a single ticket
        tmpl.Lookup("ticket.html").ExecuteTemplate(w, "ticket", structs.DataSingleTicket{Session: currentSession, Ticket: ticket, Tickets: globals.Tickets, Users: users})

        return
    }

    http.Redirect(w, r, "/", http.StatusFound)
}

// handleUpdateTicket gets the requested ticket via the url GET parameters and serves it to the template
func handleUpdateTicket(w http.ResponseWriter, r *http.Request) {

    // Only support POST request
    if r.Method == "POST" {

        // Get the session
        currentSession, errCheckForSession := session.CheckForSession(w, r)

        if errCheckForSession != nil {
            log.Println("Unable to get session")
        }

        // Get form values
        ticketId := template.HTMLEscapeString(r.FormValue("ticket"))
        status := template.HTMLEscapeString(r.FormValue("status"))
        mail := template.HTMLEscapeString(r.FormValue("mail"))
        reply := template.HTMLEscapeString(r.FormValue("reply"))
        replyType := template.HTMLEscapeString(r.FormValue("replyType"))
        merge := template.HTMLEscapeString(r.FormValue("merge"))

        // Get the ticket which was edited
        currentTicket := globals.Tickets[ticketId]

        // Update the current ticket
        updatedTicket := ticket.UpdateTicket(status, mail, reply, replyType, currentTicket)

        if merge != "" {
            // Get the ticket to merge from the tickets map
            ticketFrom := globals.Tickets[merge]

            // Only if they have the same assigned user
            if ticketFrom.User == currentSession.User && updatedTicket.User == currentSession.User {

                // Merge structs.Ticket
                ticketMergedTo, ticketMergedFrom := ticket.MergeTickets(updatedTicket, ticketFrom)

                // Write both tickets to memory
                globals.Tickets[ticketMergedTo.Id] = ticketMergedTo
                globals.Tickets[ticketMergedFrom.Id] = ticketMergedFrom

                // Persist both tickets to file system
                filehandler.WriteTicketFile(globals.ServerConfig.Tickets, &ticketMergedTo)
                filehandler.WriteTicketFile(globals.ServerConfig.Tickets, &ticketMergedFrom)

                // Update to the merged ticket so serve to client
                updatedTicket = globals.Tickets[ticketMergedTo.Id]
            }
        } else {

            // Assign the updated ticket to the ticket map in memory
            globals.Tickets[ticketId] = updatedTicket

            // Persist the updated ticket to the file system
            filehandler.WriteTicketFile(globals.ServerConfig.Tickets, &updatedTicket)
        }

        // Publish mail if the reply was selected for external
        if replyType == "extern" {
            api_out.SendMail(updatedTicket.Customer, updatedTicket.Subject, reply)
        }

        // Redirect to the ticket again, now with updated Values
        tmpl.Lookup("ticket.html").ExecuteTemplate(w, "ticket", structs.DataSingleTicket{Session: currentSession, Ticket: updatedTicket, Tickets: globals.Tickets})

        return
    }

    http.Redirect(w, r, "/", http.StatusFound)
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
        currentSession, errCheckForSession := session.CheckForSession(w, r)

        if errCheckForSession != nil {
            log.Println("Unable to get session")
        }

        // Make sure the requesting user owns the ticket
        if currentSession.User.Id == currentTicket.User.Id {

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
        currentSession, errCheckForSession := session.CheckForSession(w, r)

        if errCheckForSession != nil {
            log.Println("Unable to get session")
        }

        if currentSession.IsLoggedIn {

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
