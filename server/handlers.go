package server

import (
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/mortenterhart/trivial-tickets/structs"
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

	tmpl.Lookup("index.html").ExecuteTemplate(w, "index", session)
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

}

// handleHoliday activates / deactivates the holiday mode for a given user
func handleHoliday(w http.ResponseWriter, r *http.Request) {

	// Get session id
	sessionId := getSessionId(r)

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
	filehandler.WriteUserFile(serverConfig.Users, &users)

	// Redirect the user to the index
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
		log.Fatal(errUserCookie)
	}

	return userCookie.Value
}
