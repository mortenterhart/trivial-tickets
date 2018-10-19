package server

import (
	"html/template"
	"net/http"
	"time"

	"github.com/mortenterhart/trivial-tickets/structs"
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

	if _, err := r.Cookie("session"); err != nil {

		cookie, sessionId := createSessionCookie()
		http.SetCookie(w, cookie)
		sessions[sessionId] = CreateSession(sessionId)

		session = sessions[sessionId].Session

	} else {
		cookie, _ := r.Cookie("session")

		session = sessions[cookie.Value].Session
	}

	tmpl.Lookup("index.html").ExecuteTemplate(w, "index", session)
}

// handleLogin checks the login credentials against the stored users
// and allows the user access, if their credentials are correct
func handleLogin(w http.ResponseWriter, r *http.Request) {

	// Get session id
	userCookie, _ := r.Cookie("session")

	// Only handle POST-Requests
	if r.Method == "POST" {

		// Get submitted form values
		username := template.HTMLEscapeString(r.FormValue("username"))
		password := template.HTMLEscapeString(r.FormValue("password"))

		if username == "admin" {

			// TODO: get user data from json / memory
			// hash cmp
			if password == "admin" {

				// Update the session for the user
				u := structs.User{
					Name: "Max Mustermann",
				}

				// Create a session to update the current one
				session, _ := GetSession(userCookie.Value)
				session.User = u
				session.IsLoggedIn = true
				session.CreateTime = time.Now()

				// Update the session with the one just created
				UpdateSession(userCookie.Value, session)
			}
		}
	}

	// Redirect the user to the index
	http.Redirect(w, r, "/", 302)
}

// handleLogout logs the user out and clears their session
func handleLogout(w http.ResponseWriter, r *http.Request) {

	// Get the cookie with the session id
	userCookie, _ := r.Cookie("session")

	if r.Method == "POST" {

		// Remove the session of the user
		delete(sessions, userCookie.Value)

		// Delete the session cookie
		http.SetCookie(w, deleteSessionCookie())
	}

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
		Expires:  time.Now().Add(365 * 24 * time.Hour)}, sessionId
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
