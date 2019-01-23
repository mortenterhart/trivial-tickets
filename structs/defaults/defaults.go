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
// the server and the test suite.
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

const (
	// The following constants are the default settings
	// for the productive server. Do not use those inside
	// test cases.
	ServerPort        uint16 = 8443
	ServerTickets     string = "./files/tickets"
	ServerUsers       string = "./files/users/users.json"
	ServerMails       string = "./files/mails"
	ServerCertificate string = "./ssl/server.cert"
	ServerKey         string = "./ssl/server.key"
	ServerWeb         string = "./www"

	// These is the default logging configuration of
	// the server. It can be safely used inside tests.
	LogVerbose     bool   = false
	LogFullPaths   bool   = false
	LogLevelString string = "info"

	// The following values are the default CLI configuration.
	CliHost        string = "localhost"
	CliPort        uint16 = 8443
	CliCertificate string = "./ssl/server.cert"
	CliFetch       bool   = false
	CliSubmit      bool   = false

	// These constants are testing values for the server
	// configuration. The directories for tickets and
	// mails have been changed.
	TestPort        uint16 = 8443
	TestTickets     string = "../../files/testtickets"
	TestUsers       string = "../../files/users/users.json"
	TestMails       string = "../../files/testmails"
	TestCertificate string = "../../ssl/server.cert"
	TestKey         string = "../../ssl/server.key"
	TestWeb         string = "../../www"

	// These constants are testing values as well, but
	// for packages which are only "one directory deep"
	// as seen from the project's root directory.
	TestTicketsTrimmed     string = "../files/tickets"
	TestUsersTrimmed       string = "../files/users/users.json"
	TestMailsTrimmed       string = "../files/mails"
	TestCertificateTrimmed string = "../ssl/server.cert"
	TestKeyTrimmed         string = "../ssl/server.key"
	TestWebTrimmed         string = "../www"
)

const (
	// FileModeRegular is the default file mode used
	// to write ticket and mail files as well as the
	// users file.
	FileModeRegular = os.FileMode(0644)
)

// ExitCode is a type to represent exit codes of the
// server.
type ExitCode int

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
