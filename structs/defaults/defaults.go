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

// Package defaults defines default constraints for
// the server and the test suite. Do not modify these
// settings because the server and the tests rely on
// them.
package defaults

import "os"

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
 * Package defaults
 * Default constraints for the server and the tests
 */

// Global default constraints for the server and the
// included test cases. Do not modify these settings
// because the server and the tests rely on them.
const (
	// The following constants are the default settings
	// for the productive server. Do not modify these or
	// use them in test cases.
	ServerPort        uint16 = 8443                       // The default server port
	ServerTickets     string = "./files/tickets"          // The default ticket directory path
	ServerUsers       string = "./files/users/users.json" // The default user file path
	ServerMails       string = "./files/mails"            // The default mail directory path
	ServerCertificate string = "./ssl/server.cert"        // The default SSL certificate file
	ServerKey         string = "./ssl/server.key"         // The default SSL private key file
	ServerWeb         string = "./www"                    // The default web directory

	// The following values are an addition to the default
	// server configuration. They can be used in packages
	// and tests which are located in a subdirectory of
	// the project's root directory.
	ServerTicketsTrimmed string = "../files/tickets"          // The trimmed default ticket directory path
	ServerUsersTrimmed   string = "../files/users/users.json" // The trimmed default user file path
	ServerMailsTrimmed   string = "../files/mails"            // The trimmed default mail directory path

	// These constants form the default logging configuration
	// of the server. It can be safely used inside tests.
	LogVerbose     bool   = false  // The default value for the verbose logging option
	LogFullPaths   bool   = false  // The default value for the full paths logging option
	LogLevelString string = "info" // The default value for the log level option

	// The following values are the default CLI configuration.
	CliHost        string = "localhost"         // The default CLI hostname
	CliPort        uint16 = 8443                // The default CLI port
	CliCertificate string = "./ssl/server.cert" // The default SSL certificate file
	CliFetch       bool   = false               // The default value for the fetch option
	CliSubmit      bool   = false               // The default value for the submit option

	// These constants are testing values for the server
	// configuration. The directories for tickets and
	// mails have been changed.
	TestPort        uint16 = 8444                           // The default test port
	TestTickets     string = "../../files/testtickets"      // The default path to the test ticket directory
	TestUsers       string = "../../files/users/users.json" // The default path to the users file
	TestMails       string = "../../files/testmails"        // The default path to the test mail directory
	TestCertificate string = "../../ssl/server.cert"        // The default file path to the SSL certificate
	TestKey         string = "../../ssl/server.key"         // The default file path to the SSL private key
	TestWeb         string = "../../www"                    // The default path to the web directory

	// These constants are testing values as well, but
	// for packages which are only "one directory deep"
	// as seen from the project's root directory.
	TestTicketsTrimmed     string = "../files/testtickets"          // The trimmed default path to the test ticket directory
	TestUsersTrimmed       string = "../files/testusers/users.json" // The trimmed default path to the test users file
	TestMailsTrimmed       string = "../files/testmails"            // The trimmed default path to the test mail directory
	TestCertificateTrimmed string = "../ssl/server.cert"            // The trimmed default file path to the SSL certificate
	TestKeyTrimmed         string = "../ssl/server.key"             // The trimmed default file path to the SSL private key
	TestWebTrimmed         string = "../www"                        // The trimmed default path to the web directory
)

// Standard file modes for writing of ticket
// and mail files.
const (
	// FileModeRegular is the default file mode used
	// to write ticket and mail files as well as the
	// users file.
	FileModeRegular os.FileMode = 0644
)

// ExitCode is a type to represent exit codes of the
// server.
type ExitCode int

// The exit codes defined by the server.
const (
	// ExitSuccessful is the exit code for a successful
	// server startup and shutdown.
	ExitSuccessful ExitCode = iota

	// ExitStartError is an exit code that denotes a
	// server startup error.
	ExitStartError

	// ExitShutdownError is an exit code that denotes a
	// server shutdown error.
	ExitShutdownError
)
