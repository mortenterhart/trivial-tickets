package server

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/mortenterhart/trivial-tickets/structs"
	"github.com/mortenterhart/trivial-tickets/util/filehandler"
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

// Holds all the tickets
var tickets = make(map[string]structs.Ticket)

// Holds the given config for access to the backend systems
var serverConfig *structs.Config

// StartServer gets the parameters for the server and starts it
func StartServer(config *structs.Config) error {

	// Assign given config to the global variable
	serverConfig = config

	// Read in the users
	filehandler.ReadUserFile(serverConfig.Users, &users)

	// Read in the templates
	tmpl = GetTemplates(serverConfig.Web)

	// Register the handlers
	startHandlers(serverConfig.Web)

	// Start the server according to config
	return http.ListenAndServe(fmt.Sprintf("%s%d", ":", serverConfig.Port), nil)
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
	http.HandleFunc("/create_ticket", handleCreateTicket)
	http.HandleFunc("/holiday", handleHoliday)

	// Map the css, js and img folders to the location specified
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(path+"/static"))))
}
