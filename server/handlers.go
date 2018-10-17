package server

import (
	"go-tickets/structs"
	"html/template"
	"net/http"
	"time"
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

	session := structs.Session{Time: time.Now()}

	// Render index.html to the browser
	tmpl.Lookup("index.html").ExecuteTemplate(w, "index", session)
}

// handleLogin checks the login credentials against the stored users
// and allows the user access, if their credentials are correct
func handleLogin(w http.ResponseWriter, r *http.Request) {

	// Only handle POST-Requests
	if r.Method == "POST" {

		// Get submitted form values
		username := template.HTMLEscapeString(r.FormValue("username"))
		password := template.HTMLEscapeString(r.FormValue("password"))

		if username == "bla" {

			// TODO: get user data from json / memory
			// hash cmp
			if password == "bb" {

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
