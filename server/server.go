package server

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/mortenterhart/go-tickets/structs"
)

/*
*
* Matrikelnummern
* 3040018
* 6694964
* 3478222
 */

// Holds the parsed templates. Is defined as a global variable to only parse
// the templates once on startup, instead of on every GET - request to index.
var tmpl *template.Template

// Holds all the sessions for the users
var sessions = make(map[string]structs.SessionManager)

// Holds all the users
var users = make(map[string]structs.User)

// StartServer gets the parameters for the server and starts it
func StartServer(config *structs.Config) error {

	tmpl = GetTemplates(config.Web)

	startHandlers(config.Web)

	return http.ListenAndServe(fmt.Sprintf("%s%d", ":", config.Port), nil)
}

// GetTemplates crawls through the templates folder and reads in all
// present templates.
func GetTemplates(path string) *template.Template {

	// Crawl via relative path, since our current work dir is in cmd/ticketsystem
	t, errTemplates := template.ParseGlob(path + "/templates/*.html")

	if errTemplates != nil {
		log.Fatal("Unable to load the templates: ", errTemplates)
	}

	return t
}

// startHandlers maps all the various handles to the url patterns.
func startHandlers(path string) {
	http.HandleFunc("/", handleIndex)
	http.HandleFunc("/login", handleLogin)
	http.HandleFunc("/logout", handleLogout)
	// Map the css, js and img folders to the location specified
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(path+"/static"))))
}
