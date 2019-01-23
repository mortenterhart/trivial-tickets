// Trivial Tickets Ticketsystem
// Copyright (C) 2019 The Contributors
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

// Package server implements the web server including
// shutdown routines and the associated handlers for
// web requests.
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

	"github.com/pkg/errors"

	"github.com/mortenterhart/trivial-tickets/api/api_in"
	"github.com/mortenterhart/trivial-tickets/api/api_out"
	"github.com/mortenterhart/trivial-tickets/globals"
	"github.com/mortenterhart/trivial-tickets/logger"
	"github.com/mortenterhart/trivial-tickets/structs"
	"github.com/mortenterhart/trivial-tickets/structs/defaults"
	"github.com/mortenterhart/trivial-tickets/util/filehandler"
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

// interrupt is the channel which receives potential
// interrupt or kill signals in order to shutdown
// the server. This variable is needed to provide
// a functionality for the tests to confidently
// stop the server.
var interrupt chan os.Signal

// shutdownTimeout is the duration the shutdown context
// waits before closing the server.
const shutdownTimeout = 5 * time.Second

// StartServer gets the parameters for the server and starts it.
func StartServer(config *structs.ServerConfig) (defaults.ExitCode, error) {

	// Assign given config to the global variable
	logger.Info("Initializing server configuration")
	logServerConfig(config)
	globals.ServerConfig = config

	// Create the folders for tickets and mails if they do not exist yet
	if createErr := createResourceFolders(config); createErr != nil {
		logger.Error(errors.Wrap(createErr, "unable to create resource directories"))
	}

	// Read the users file
	logger.Info("Reading users file", config.Users)
	if errReadUserFile := filehandler.ReadUserFile(config.Users, &users); errReadUserFile != nil {
		return defaults.ExitStartError, errors.Wrap(errReadUserFile, "unable to load user file")
	}

	// Read the tickets
	logger.Info("Reading ticket files in", config.Tickets)
	if errReadTicketFiles := filehandler.ReadTicketFiles(config.Tickets, &globals.Tickets); errReadTicketFiles != nil {
		return defaults.ExitStartError, errors.Wrap(errReadTicketFiles, "unable to load ticket files")
	}

	// Read the mails
	logger.Info("Reading mail files in", config.Mails)
	if errReadMailFiles := filehandler.ReadMailFiles(config.Mails, &globals.Mails); errReadMailFiles != nil {
		return defaults.ExitStartError, errors.Wrap(errReadMailFiles, "unable to load mail files")
	}

	// Read the HTML templates
	logger.Info("Loading HTML templates in", config.Web)
	if tmpl = getTemplates(config.Web); tmpl == nil {
		return defaults.ExitStartError, errors.New("unable to load HTML templates")
	}

	// Register the handlers
	logger.Info("Registering handlers")
	handler, errStartHandlers := registerHandlers(config.Web)
	if errStartHandlers != nil {
		return defaults.ExitStartError, errors.Wrap(errStartHandlers, "unable to register handlers")
	}

	// Start a GoRoutine to redirect http requests to https
	logger.Info("Starting Go routine to redirect http requests to https")
	go func() {
		err := http.ListenAndServe("localhost:80", http.HandlerFunc(redirectToTLS))
		logger.Error("error starting redirect handler:", err)
	}()

	logger.Info("Server setup completed and starting server")

	interrupt = notifyOnInterruptSignal()
	startError := make(chan error)

	server := http.Server{
		Addr:     fmt.Sprintf("localhost:%d", config.Port),
		Handler:  handler,
		ErrorLog: logger.NewErrorLogger(),
	}

	go func() {
		// Log on which socket the server is listening
		logger.Infof("Server listening on https://localhost:%d (PID = %d), type Ctrl-C to stop",
			config.Port, os.Getpid())

		// Start the server according to config
		serverErr := server.ListenAndServeTLS(config.Cert, config.Key)

		// Catch potential start errors and close the channel
		// because this is the last value sent
		startError <- serverErr
		close(startError)
	}()

	return handleServerShutdown(&server, startError, interrupt)
}

// StopServer sends an artificial interrupt signal to the
// server to make it shutdown gracefully. This function is
// useful especially for tests that have to start the
// productive server.
func StopServer() {
	logger.Info("Stopping server by sending interrupt signal")
	interrupt <- os.Interrupt
	logger.Info("Sending of interrupt signal completed")
	close(interrupt)
}

// notifyOnInterruptSignal creates a channel which reports signals sent to the
// server process in order to stop it. The channel catches interrupt (SIGINT),
// kill (SIGKILL) and terminating (SIGTERM) signals.
func notifyOnInterruptSignal() chan os.Signal {
	signalListener := make(chan os.Signal)
	signal.Notify(signalListener, os.Interrupt, os.Kill, syscall.SIGTERM)
	return signalListener
}

// handleServerShutdown handles the shutdown routine of the given server. The given
// channels have to be created in advance and report actions causing the server to
// interrupt and shutdown finally. The startError channel indicates if the server
// experienced an error at startup which prevented it from starting. The interrupt
// channel reports if the server has captured a system signal such as an interrupt
// or a kill signal which causes the process to stop. The function returns the
// server's exit status and an error if a start error occurred or the server was
// unable to shutdown correctly.
func handleServerShutdown(server *http.Server, startError <-chan error, interrupt <-chan os.Signal) (defaults.ExitCode, error) {
	exitCode := defaults.ExitSuccessful

	select {
	case serverErr := <-startError:
		if serverErr != http.ErrServerClosed {
			returnErr := errors.Wrap(serverErr, "error while starting server")
			logger.Error(returnErr)
			return defaults.ExitStartError, returnErr
		}
	case capturedSignal := <-interrupt:
		switch capturedSignal {
		case os.Interrupt:
			logger.Infof("Captured terminating signal '%s' (SIGINT)", capturedSignal)
			exitCode = defaults.ExitSuccessful

		case os.Kill:
			logger.Infof("Captured terminating signal '%s' (SIGKILL), preferred way is SIGINT", capturedSignal)
			exitCode = defaults.ExitShutdownError

		case syscall.SIGTERM:
			logger.Infof("Captured terminating signal '%s' (SIGTERM), preferred way is SIGINT", capturedSignal)
			exitCode = defaults.ExitShutdownError
		}

		if shutdownErr := shutdownGracefully(server, shutdownTimeout); shutdownErr != nil {
			logger.Error("Server shutdown caused error:", shutdownErr)
			return defaults.ExitShutdownError, shutdownErr
		}
		logger.Info(<-startError)

		logger.Info("Server shutdown successful")
	}

	return exitCode, nil
}

// shutdownGracefully takes a server and a timeout and attempts to shutdown the
// applied server gracefully. It is wait until all connections are finished or
// the timeout has exceeded.
func shutdownGracefully(server *http.Server, timeout time.Duration) error {
	logger.Infof("Shutting down server gracefully with timeout of %s", timeout)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if shutdownErr := server.Shutdown(ctx); shutdownErr != nil {
		return errors.Wrap(shutdownErr, "server shutdown failed")
	}

	return nil
}

// getTemplates crawls through the templates folder and reads in all
// present templates.
func getTemplates(path string) *template.Template {

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

	target := "https://" + req.Host + fmt.Sprintf(":%d", globals.ServerConfig.Port)

	if len(req.URL.RawQuery) > 0 {
		target += "?" + req.URL.RawQuery
	}

	http.Redirect(w, req, target, http.StatusTemporaryRedirect)
}

// registerHandlers maps all the various handles to the url patterns and
// returns them as request handler multiplexer.
func registerHandlers(path string) (*http.ServeMux, error) {

	// Check if the path exists
	if !filehandler.DirectoryExists(path) {
		return nil, errors.New("no path given for web folders")
	}

	logger.Info("Starting handlers for incoming HTTP requests")
	mainHandler := http.NewServeMux()
	mainHandler.HandleFunc("/", handleIndex)
	mainHandler.HandleFunc("/login", handleLogin)
	mainHandler.HandleFunc("/logout", handleLogout)
	mainHandler.HandleFunc("/createTicket", handleCreateTicket)
	mainHandler.HandleFunc("/holiday", handleHoliday)
	mainHandler.HandleFunc("/ticket", handleTicket)
	mainHandler.HandleFunc("/updateTicket", handleUpdateTicket)
	mainHandler.HandleFunc("/unassignTicket", handleUnassignTicket)
	mainHandler.HandleFunc("/assignTicket", handleAssignTicket)
	mainHandler.HandleFunc("/api/receive", api_in.ReceiveMail)
	mainHandler.HandleFunc("/api/fetchMails", api_out.FetchMails)
	mainHandler.HandleFunc("/api/verifyMail", api_out.VerifyMailSent)

	// Map the css, js and img folders to the location specified
	mainHandler.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(path+"/static"))))

	return mainHandler, nil
}

// createResourceFolders checks if the required ticket and mail
// paths given inside the server config exist and creates them
// if not.
func createResourceFolders(config *structs.ServerConfig) (returnErr error) {
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

// logServerConfig outputs the server configuration to the console.
func logServerConfig(config *structs.ServerConfig) {
	logger.Info("  Port:", config.Port)
	logger.Info("  Tickets:", config.Tickets)
	logger.Info("  Users:", config.Users)
	logger.Info("  Mails:", config.Mails)
	logger.Info("  Cert:", config.Cert)
	logger.Info("  Key:", config.Key)
	logger.Info("  Web:", config.Web)
}
