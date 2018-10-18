package server

import (
	"html/template"
	"net/http"
	"time"

	"github.com/mortenterhart/go-tickets/structs"
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

	userCookie, _ := r.Cookie("session")

	// Only handle POST-Requests
	if r.Method == "POST" {

		// Get submitted form values
		username := template.HTMLEscapeString(r.FormValue("username"))
		password := template.HTMLEscapeString(r.FormValue("password"))

		if username == "bla" {

			// TODO: get user data from json / memory
			// hash cmp
			if password == "bb" {

				session, _ := GetSession(userCookie.Value)
				session.User = users[username]
				session.IsLoggedIn = true
				session.CreateTime = time.Now()

				UpdateSession(userCookie.Value, session)

				// set session
				http.Redirect(w, r, "/", 302)
			}
		}
	}
}

// handleLogout logs the user out and clears their session
func handleLogout(w http.ResponseWriter, r *http.Request) {

	// TODO: clear session
	http.Redirect(w, r, "/", 302)
}

func createSessionCookie() (*http.Cookie, string) {

	sessionId := CreateSessionId()

	return &http.Cookie{
		Name:     "session",
		Value:    sessionId,
		HttpOnly: false,
		Expires:  time.Now().Add(365 * 24 * time.Hour)}, sessionId
}
