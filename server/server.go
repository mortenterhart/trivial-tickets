// Server starting and handler registration
package server

import (
	"context"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mortenterhart/trivial-tickets/api/api_in"
	"github.com/mortenterhart/trivial-tickets/api/api_out"
	"github.com/mortenterhart/trivial-tickets/globals"
	"github.com/mortenterhart/trivial-tickets/logger"
	"github.com/mortenterhart/trivial-tickets/structs"
	"github.com/mortenterhart/trivial-tickets/util/filehandler"
	"github.com/pkg/errors"
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
func StartServer(config *structs.Config) (int, error) {

	// Assign given config to the global variable
	logger.Info("Initializing server configuration")
	logServerConfig(config)
	globals.ServerConfig = config

	// Create the folders for tickets, mails and users if it does not exist
	if createErr := createResourceFolders(config); createErr != nil {
		logger.Error(errors.Wrap(createErr, "unable to create resource directories"))
	}

	// Read the users file
	logger.Info("Reading users file", config.Users)
	if errReadUserFile := filehandler.ReadUserFile(config.Users, &users); errReadUserFile != nil {
		return 1, errors.Wrap(errReadUserFile, "unable to load user file")
	}

	// Read the tickets
	logger.Info("Reading ticket files in", config.Tickets)
	if errReadTicketFiles := filehandler.ReadTicketFiles(config.Tickets, &globals.Tickets); errReadTicketFiles != nil {
		return 1, errors.Wrap(errReadTicketFiles, "unable to load ticket files")
	}

	// Read the mails
	logger.Info("Reading mail files in", config.Mails)
	if errReadMailFiles := filehandler.ReadMailFiles(config.Mails, &globals.Mails); errReadMailFiles != nil {
		return 1, errors.Wrap(errReadMailFiles, "unable to load mail files")
	}

	// Read the HTML templates
	logger.Info("Loading HTML templates in", config.Web)
	if tmpl = GetTemplates(config.Web); tmpl == nil {
		return 1, errors.New("unable to load HTML templates")
	}

	// Register the handlers
	logger.Info("Registering handlers")
	if errStartHandlers := startHandlers(config.Web); errStartHandlers != nil {
		return 1, errors.Wrap(errStartHandlers, "unable to register handlers")
	}

	// Start a GoRoutine to redirect http requests to https
	logger.Info("Starting Go routine to redirect http requests to https")
	go http.ListenAndServe(":80", http.HandlerFunc(redirectToTLS))

	logger.Info("Server setup completed and starting server")

	interrupt := notifyOnInterruptSignal()
	startError := make(chan error)

	server := http.Server{
		Addr: fmt.Sprintf("%s%d", ":", config.Port),
	}

	go func() {
		// Log on which socket the server is listening
		logger.Infof("server listening on https://localhost:%d (PID = %d), type Ctrl-C to stop",
			config.Port, os.Getpid())

		// Start the server according to config
		serverErr := server.ListenAndServeTLS(config.Cert, config.Key)
		startError <- serverErr
	}()

	return handleServerShutdown(&server, startError, interrupt)
}

func notifyOnInterruptSignal() <-chan os.Signal {
	signalListener := make(chan os.Signal)
	signal.Notify(signalListener, os.Interrupt, os.Kill, syscall.SIGTERM)
	return signalListener
}

func handleServerShutdown(server *http.Server, startError <-chan error, interrupt <-chan os.Signal) (int, error) {
	exitCode := 0

	select {
	case serverErr := <-startError:
		if serverErr != http.ErrServerClosed {
			return 1, errors.Wrap(serverErr, "error while starting server")
		}
	case capturedSignal := <-interrupt:
		switch capturedSignal {
		case os.Interrupt:
			logger.Infof("Captured terminating signal %s (SIGINT)", capturedSignal)
			exitCode = 0

		case os.Kill:
			logger.Infof("Captured terminating signal %s (SIGKILL), preferred way is SIGINT", capturedSignal)
			exitCode = 1

		case syscall.SIGTERM:
			logger.Infof("Captured terminating signal %s (SIGTERM), preferred way is SIGINT", capturedSignal)
			exitCode = 1

		}

		timeout := 5 * time.Second
		if shutdownErr := shutdownGracefully(server, timeout); shutdownErr != nil {
			logger.Error("Server shutdown caused error:", shutdownErr)
			return 1, shutdownErr
		}

		logger.Info("Server shutdown successful")
	}

	return exitCode, nil
}

func shutdownGracefully(server *http.Server, timeout time.Duration) error {
	logger.Infof("Shutting down server gracefully with timeout of %s", timeout)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if shutdownErr := server.Shutdown(ctx); shutdownErr != nil {
		return errors.Wrap(shutdownErr, "server shutdown failed")
	}

	return nil
}

// GetTemplates crawls through the templates folder and reads in all
// present templates.
func GetTemplates(path string) *template.Template {

	// Crawl via relative path, since our current work dir is in cmd/ticketsystem
	t, errTemplates := template.ParseGlob(path + "/templates/*.html")

	if errTemplates != nil {
		logger.Error("unable to load the templates:", errTemplates)
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

	logger.Info("Starting handlers for incoming HTTP requests")
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
	http.NewServeMux()

	return nil
}

// createResourceFolders checks if the required ticket and mail
// paths given inside the server config exist and creates them
// if not.
func createResourceFolders(config *structs.Config) (returnErr error) {
	if !filehandler.DirectoryExists(config.Tickets) {
		logger.Infof("Creating missing ticket directory '%s'", config.Tickets)
		if createErr := filehandler.CreateFolders(config.Tickets); createErr != nil {
			returnErr = errors.Wrap(createErr, "error while creating ticket directory")
			logger.Error(returnErr)
		}
	}

	if !filehandler.DirectoryExists(config.Mails) {
		logger.Infof("Creating missing mail directory '%s'", config.Mails)
		if createErr := filehandler.CreateFolders(config.Mails); createErr != nil {
			returnErr = errors.Wrap(createErr, "error while creating mail directory")
			logger.Error(returnErr)
		}
	}

	return
}

// logServerConfig outputs the server configuration to the console
func logServerConfig(config *structs.Config) {
	logger.Info("  Port:", config.Port)
	logger.Info("  Tickets:", config.Tickets)
	logger.Info("  Users:", config.Users)
	logger.Info("  Mails:", config.Mails)
	logger.Info("  Cert:", config.Cert)
	logger.Info("  Key:", config.Key)
	logger.Info("  Web:", config.Web)
}
