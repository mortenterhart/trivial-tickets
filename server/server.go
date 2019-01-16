// Server starting and handler registration
package server

import (
	"fmt"
	"github.com/mortenterhart/trivial-tickets/api/api_in"
	"github.com/mortenterhart/trivial-tickets/api/api_out"
	"github.com/mortenterhart/trivial-tickets/globals"
	"github.com/mortenterhart/trivial-tickets/structs"
	"github.com/mortenterhart/trivial-tickets/util/filehandler"
	"github.com/pkg/errors"
	"html/template"
	"log"
	"net/http"
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
 * Server starting and handler registration
 */

// Holds the parsed templates. Is defined as a global variable to only parse
// the templates once on startup, instead of on every GET - request to index.
var tmpl *template.Template

// Holds all the users
var users = make(map[string]structs.User)

// StartServer gets the parameters for the server and starts it
func StartServer(config *structs.Config) error {

	// Assign given config to the global variable
	log.Println("Initializing server configuration")
	logServerConfig(config)
	globals.ServerConfig = config

	// Create the folders for tickets, mails and users if it does not exist
	if createErr := createResourceFolders(config); createErr != nil {
		log.Println(errors.Wrap(createErr, "unable to create resource directories"))
	}

	// Read the users file
	log.Println("Reading users file", config.Users)
	if errReadUserFile := filehandler.ReadUserFile(config.Users, &users); errReadUserFile != nil {
		return errors.Wrap(errReadUserFile, "unable to load user file")
	}

	// Read the tickets
	log.Println("Reading ticket files in", config.Tickets)
	if errReadTicketFiles := filehandler.ReadTicketFiles(config.Tickets, &globals.Tickets); errReadTicketFiles != nil {
		return errors.Wrap(errReadTicketFiles, "unable to load ticket files")
	}

	// Read the mails
	log.Println("Reading mail files in", config.Mails)
	if errReadMailFiles := filehandler.ReadMailFiles(config.Mails, &globals.Mails); errReadMailFiles != nil {
		return errors.Wrap(errReadMailFiles, "unable to load mail files")
	}

	// Read the HTML templates
	log.Println("Loading HTML templates in", config.Web)
	if tmpl = GetTemplates(config.Web); tmpl == nil {
		return errors.New("unable to load HTML templates")
	}

	// Register the handlers
	log.Println("Registering handlers")
	if errStartHandlers := startHandlers(config.Web); errStartHandlers != nil {
		return errors.Wrap(errStartHandlers, "unable to register handlers")
	}

	// Start a GoRoutine to redirect http requests to https
	log.Println("Starting Go routine to redirect http requests to https")
	go http.ListenAndServe(":80", http.HandlerFunc(redirectToTLS))

	log.Println("Server startup completed and ready to use")

	// Log on which socket the server is listening
	log.Printf("server listening on https://localhost:%d, type Ctrl-C to stop", config.Port)

	// Start the server according to config
	return http.ListenAndServeTLS(fmt.Sprintf("%s%d", ":", config.Port), config.Cert, config.Key, nil)
}

// GetTemplates crawls through the templates folder and reads in all
// present templates.
func GetTemplates(path string) *template.Template {

	// Crawl via relative path, since our current work dir is in cmd/ticketsystem
	t, errTemplates := template.ParseGlob(path + "/templates/*.html")

	if errTemplates != nil {
		log.Println("unable to load the templates:", errTemplates)
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
	if !filehandler.DirectoryExists(path) {
		return errors.New("no path given for web folders")
	}

	log.Println("Starting handlers for incoming HTTP requests")
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

// createResourceFolders checks if the required ticket and mail
// paths given inside the server config exist and creates them
// if not.
func createResourceFolders(config *structs.Config) (returnErr error) {
	if !filehandler.DirectoryExists(config.Tickets) {
		log.Printf("Creating missing ticket directory '%s'", config.Tickets)
		if createErr := filehandler.CreateFolders(config.Tickets); createErr != nil {
			returnErr = errors.Wrap(createErr, "error while creating ticket directory")
			log.Println(returnErr)
		}
	}

	if !filehandler.DirectoryExists(config.Mails) {
		log.Printf("Creating missing mail directory '%s'", config.Mails)
		if createErr := filehandler.CreateFolders(config.Mails); createErr != nil {
			returnErr = errors.Wrap(createErr, "error while creating mail directory")
			log.Println(returnErr)
		}
	}

	return
}

// logServerConfig outputs the server configuration to the console
func logServerConfig(config *structs.Config) {
	log.Println("  Port:", config.Port)
	log.Println("  Tickets:", config.Tickets)
	log.Println("  Users:", config.Users)
	log.Println("  Mails:", config.Mails)
	log.Println("  Cert:", config.Cert)
	log.Println("  Key:", config.Key)
	log.Println("  Web:", config.Web)
}
