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

// Command ticketsystem starts the Trivial Tickets Ticketsystem
// web server to serve as support ticket platform.
package main

import (
	"flag"
	"fmt"
	"math"
	"os"
	"path/filepath"

	"github.com/mortenterhart/trivial-tickets/globals"
	"github.com/mortenterhart/trivial-tickets/logger"
	"github.com/mortenterhart/trivial-tickets/server"
	"github.com/mortenterhart/trivial-tickets/structs"
	"github.com/mortenterhart/trivial-tickets/structs/defaults"
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
 * Package main
 * Main package of the ticketsystem webserver
 */

var (
	// Command-line options

	// Server configuration
	port    = flag.Uint("port", uint(defaults.ServerPort), "`port` on which the web server will run")
	tickets = flag.String("tickets", defaults.ServerTickets, "`directory` in which the tickets will be stored")
	users   = flag.String("users", defaults.ServerUsers, "path to the users `file`")
	mails   = flag.String("mails", defaults.ServerMails, "`directory` in which the mails will be cached")
	cert    = flag.String("cert", defaults.ServerCertificate, "location of the ssl certificate `file`")
	key     = flag.String("key", defaults.ServerKey, "location of the ssl key `file`")
	web     = flag.String("web", defaults.ServerWeb, "location of the www `directory`")

	// Logging configuration
	verbose        = flag.Bool("verbose", defaults.LogVerbose, "Enable output of verbose log (package paths, file names and line numbers)")
	fullPaths      = flag.Bool("full-paths", defaults.LogFullPaths, "Log package names and filenames with full paths instead of abbreviated ones")
	logLevelString = flag.String("log-level", defaults.LogLevelString, "Specify `level` of logging (either \"info\", \"warning\", \"error\" or \"fatal\")")
)

// exit is used as replaceable function to
// quit the program with an exit code. This
// variable is used by tests so that the
// program does not terminate.
var exit = os.Exit

// fatal is used as replaceable function to
// indicate and log a fatal error. This variable
// is used by tests so that the program does not
// terminate.
var fatal = logger.Fatal

// main is the main entry point to the ticketsystem.
func main() {

	config, errConfig := initConfig()

	if errConfig != nil {
		fatal(errConfig)
		return
	}

	exitCode, errServer := server.StartServer(&config)

	if errServer != nil {
		fatal("fatal error occurred:", errServer, "\nServer startup failed!")
		return
	}

	logger.Info("exiting with exit code:", exitCode)
	exit(int(exitCode))
}

// initConfig parses the command line arguments and
// populates a struct for the config parameters.
// It returns this struct.
func initConfig() (structs.ServerConfig, error) {
	globals.LogConfig = &structs.LogConfig{}

	// Set another usage message function
	flag.Usage = usageMessage

	// Parse all arguments, e.g. populate the variables
	flag.Parse()

	// If the port is not within boundaries, return an error
	if !isPortInBoundaries(*port) {
		return structs.ServerConfig{}, fmt.Errorf("applied port %d is not a correct port number", *port)
	}

	logLevel, convertErr := convertLogLevel(*logLevelString)
	if convertErr != nil {
		return structs.ServerConfig{}, convertErr
	}

	logConfig := structs.LogConfig{
		LogLevel:  logLevel,
		Verbose:   *verbose,
		FullPaths: *fullPaths,
	}
	globals.LogConfig = &logConfig

	// Populate and return the struct
	return structs.ServerConfig{
		Port:    uint16(*port),
		Tickets: *tickets,
		Users:   *users,
		Mails:   *mails,
		Cert:    *cert,
		Key:     *key,
		Web:     *web,
	}, nil
}

// isPortInBoundaries returns true if the provided port
// is within the boundaries of a 16 bit unsigned integer,
// false otherwise. Since the port numbers only go up to a
// 16 bit unsigned integer.
func isPortInBoundaries(port uint) bool {
	return port > 0 && port <= math.MaxUint16
}

// usageMessage writes a help message with all options to
// the output buffer (stderr by default).
func usageMessage() {
	w := flag.CommandLine.Output()
	fmt.Fprintf(w, "Usage: %s [options]\n", filepath.Base(os.Args[0]))
	fmt.Fprintln(w, "Trivial Tickets Web server")
	fmt.Fprintln(w)

	fmt.Fprintln(w, "The following options can be used to alter the default")
	fmt.Fprintln(w, "configuration of the server and the logger. No option")
	fmt.Fprintln(w, "is required and the default values are portrayed.")
	fmt.Fprintln(w)

	fmt.Fprintln(w, "Help options:")
	fmt.Fprintln(w, "  -h, -help       Print this help text and exit.")
	fmt.Fprintln(w)

	fmt.Fprintln(w, "Server options:")
	fmt.Fprintln(w, "  -port <PORT>    The port on which the web server will run. PORT has to")
	fmt.Fprintf (w, "                  be a 16 bit unsigned integer (0 < PORT <= %d) and\n", math.MaxUint16)
	fmt.Fprintln(w, "                  should not be used by another process.")
	fmt.Fprintf (w, "                  (Default: %d)\n", defaults.ServerPort)
	fmt.Fprintln(w, "  -tickets <DIR>  The directory in which the tickets will be stored. DIR")
	fmt.Fprintln(w, "                  should be an existing directory with active write privileges,")
	fmt.Fprintln(w, "                  otherwise it is created on startup.")
	fmt.Fprintf (w, "                  (Default \"%s\")\n", defaults.ServerTickets)
	fmt.Fprintln(w, "  -users <FILE>   The file path to the users.json file. FILE must be an")
	fmt.Fprintln(w, "                  existing file and should contain valid JSON with users.")
	fmt.Fprintf (w, "                  (Default: \"%s\")\n", defaults.ServerUsers)
	fmt.Fprintln(w, "  -mails <DIR>    The directory where new mails will be saved. DIR should be")
	fmt.Fprintln(w, "                  an existing directory with active write privileges,")
	fmt.Fprintln(w, "                  otherwise it is created on startup.")
	fmt.Fprintf (w, "                  (Default: \"%s\")\n", defaults.ServerMails)
	fmt.Fprintln(w, "  -cert <FILE>    The path to the ssl certificate file. FILE must be an")
	fmt.Fprintln(w, "                  existing file with a valid ssl certificate.")
	fmt.Fprintf (w, "                  (Default: \"%s\")\n", defaults.ServerCertificate)
	fmt.Fprintln(w, "  -key <FILE>     The path to the ssl key file. FILE must be an existing")
	fmt.Fprintln(w, "                  file with a valid ssl public key.")
	fmt.Fprintf (w, "                  (Default: \"%s\")\n", defaults.ServerKey)
	fmt.Fprintln(w, "  -web <DIR>      The root directory of the web server. DIR must be an")
	fmt.Fprintln(w, "                  existing directory and should match the templates and")
	fmt.Fprintln(w, "                  static paths pointing to the server resources.")
	fmt.Fprintf (w, "                  (Default: \"%s\")\n", defaults.ServerWeb)
	fmt.Fprintln(w)

	fmt.Fprintln(w, "Logging options:")
	fmt.Fprintln(w, "  -log-level <LEVEL>")
	fmt.Fprintln(w, "                  Specify the level of logging. This can be one of:")
	fmt.Fprintln(w, "                    info     all log messages (default)")
	fmt.Fprintln(w, "                    warning  only warnings, errors and fatal errors")
	fmt.Fprintln(w, "                    error    only errors and fatal errors")
	fmt.Fprintln(w, "                    fatal    only fatal errors")
	fmt.Fprintln(w, "                  Fatal errors are always logged and it is not recommended setting")
	fmt.Fprintln(w, "                  the level to 'fatal'.")
	fmt.Fprintln(w, "  -verbose        Enable verbose logging (includes package and function name, filenames")
	fmt.Fprintln(w, "                  and line numbers for every message, but with abbreviated paths).")
	fmt.Fprintln(w, "  -full-paths     Log package paths and file paths with full paths instead of")
	fmt.Fprintln(w, "                  abbreviated ones. Warning: This will extend log messages a lot")
	fmt.Fprintln(w, "                  and they will not fit on every screen in one row. This option")
	fmt.Fprintln(w, "                  is compatible with -verbose.")
}

// convertLogLevel maps a given string with the `-log-level`
// flag to its enum equivalent. If the provided log level is
// not defined, it returns an error.
func convertLogLevel(logLevelString string) (level structs.LogLevel, convertErr error) {
	level = structs.AsLogLevel(logLevelString)
	if level < 0 {
		convertErr = fmt.Errorf("log level '%s' not defined", logLevelString)
	}

	return
}
