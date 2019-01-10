package server

import (
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/mortenterhart/trivial-tickets/api/api_in"
	"github.com/mortenterhart/trivial-tickets/api/api_out"
	"github.com/mortenterhart/trivial-tickets/globals"
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

// Holds all the users
var users = make(map[string]structs.User)

// StartServer gets the parameters for the server and starts it
func StartServer(config *structs.Config) error {

	// Assign given config to the global variable
	globals.ServerConfig = config

	// Read in the users
	errReadUserFile := filehandler.ReadUserFile(globals.ServerConfig.Users, &users)

	if errReadUserFile == nil {
		// Read in the tickets
		errReadTicketFiles := filehandler.ReadTicketFiles(globals.ServerConfig.Tickets, &globals.Tickets)

		if errReadTicketFiles == nil {
			// Read in the templates
			tmpl = GetTemplates(globals.ServerConfig.Web)

			// Register the handlers
			errStartHandlers := startHandlers(globals.ServerConfig.Web)

			if tmpl != nil && errStartHandlers == nil {

				// Start a GoRoutine to redirect http requests to https
				go http.ListenAndServe(":80", http.HandlerFunc(redirectToTLS))

				// Start the server according to config
				return http.ListenAndServeTLS(fmt.Sprintf("%s%d", ":", globals.ServerConfig.Port), globals.ServerConfig.Cert, globals.ServerConfig.Key, nil)
			} else {
				return errors.New("unable to load templates / register handlers")
			}
		} else {
			return errors.New("unable to load ticket files")
		}
	} else {
		return errors.New("unable to load user file")
	}
}

// GetTemplates crawls through the templates folder and reads in all
// present templates.
func GetTemplates(path string) *template.Template {

	// Crawl via relative path, since our current work dir is in cmd/ticketsystem
	t, errTemplates := template.ParseGlob(path + "/templates/*.html")

	if errTemplates != nil {
		log.Println("Unable to load the templates: ", errTemplates)
		return nil
	}

	return t
}

// redirectToTLS is invoked as soon as someone tries to reach the ticket system
// via http, the request is then redirected to https.
// Taken from https://gist.github.com/d-schmidt/587ceec34ce1334a5e60
func redirectToTLS(w http.ResponseWriter, req *http.Request) {

	target := "https://" + req.Host + fmt.Sprintf("%s%d", ":", globals.ServerConfig.Port)

	if len(req.URL.RawQuery) > 0 {
		target += "?" + req.URL.RawQuery
	}

	http.Redirect(w, req, target, http.StatusTemporaryRedirect)
}

// startHandlers maps all the various handles to the url patterns.
func startHandlers(path string) error {

	// Check if the path exists
	// Taken from https://gist.github.com/mattes/d13e273314c3b3ade33f
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return errors.New("no path given for web folders")
	}

	http.HandleFunc("/", handleIndex)
	http.HandleFunc("/login", handleLogin)
	http.HandleFunc("/logout", handleLogout)
	http.HandleFunc("/create_ticket", handleCreateTicket)
	http.HandleFunc("/holiday", handleHoliday)
	http.HandleFunc("/ticket", handleTicket)
	http.HandleFunc("/updateTicket", handleUpdateTicket)
	http.HandleFunc("/unassignTicket", handleUnassignTicket)
	http.HandleFunc("/assignTicket", handleAssignTicket)
	http.HandleFunc("/api/receive", api_in.ReceiveMail)
	http.HandleFunc("/api/fetchMails", api_out.FetchMails)
	http.HandleFunc("/api/verifyMail", api_out.VerifyMailSent)

	// Map the css, js and img folders to the location specified
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(path+"/static"))))

	return nil
}
