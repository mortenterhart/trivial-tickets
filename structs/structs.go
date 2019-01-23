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

// Package structs supplies project-wide needed data
// structures, types and constants for the server and
// the command-line tool.
package structs

import (
	"strconv"
	"time"
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
 * Package structs
 * Project-wide needed structures for data elements
 */

// ServerConfig is a struct to hold the config parameters
// provided on startup.
type ServerConfig struct {
	// Port is the port the server is listening to.
	Port uint16

	// Tickets is the directory in which
	// all tickets are stored.
	Tickets string

	// User is the path to the users.json file.
	Users string

	// Mails is the directory in which all
	// mails are stored.
	Mails string

	// Cert is the path to the SSL
	// Certificate file.
	Cert string

	// Key is the path to the SSL Key file.
	Key string

	// Web is the root directory of the server
	// where web resources are located.
	Web string
}

// CLIConfig is a struct to hold the CLI config
// parameters provided on startup.
type CLIConfig struct {
	Host string
	Port uint16
	Cert string
}

// LogConfig defines the logging configuration for
// the server provided by the startup configuration
// flags.
type LogConfig struct {
	LogLevel  LogLevel
	Verbose   bool
	FullPaths bool
}

// RandomIDLength is the length of the ticket and
// mail ids.
const RandomIDLength int = 10

// LogLevel is a type that defines the level of logging.
// A high log level such as INFO means that more
// information and actions are logged to the console.
// A low log level suppresses the output of messages of
// a higher log level.
type LogLevel int

const (
	// LevelInfo is the default log level (all messages
	// are logged).
	LevelInfo LogLevel = iota

	// LevelWarning logs warnings, errors and fatal errors.
	LevelWarning

	// LevelError only logs errors and fatal errors.
	LevelError

	// LevelFatal only logs fatal errors (not recommended).
	// Fatal errors are always logged.
	LevelFatal

	// LevelTestDebug is an additional log level used only
	// for debugging inside tests (not configurable).
	LevelTestDebug
)

// String converts a log level to its corresponding
// output string used in the log.
func (level LogLevel) String() string {
	switch level {
	case LevelInfo:
		return "[INFO]"

	case LevelWarning:
		return "[WARNING]"

	case LevelError:
		return "[ERROR]"

	case LevelFatal:
		return "[FATAL ERROR]"

	case LevelTestDebug:
		return "[TEST DEBUG]"
	}

	return "undefined"
}

// AsLogLevel converts a given log level string to
// a log level. If the log level string is not
// defined, the return value is -1.
func AsLogLevel(logLevelString string) LogLevel {
	switch logLevelString {
	case "info":
		return LevelInfo

	case "warning":
		return LevelWarning

	case "error":
		return LevelError

	case "fatal":
		return LevelFatal
	}

	return LogLevel(-1)
}

// Session is a struct that holds session variables
// for a certain user.
type Session struct {
	ID           string
	User         User
	CreationTime time.Time
	IsLoggedIn   bool
}

// SessionManager holds a session and operates on a
// it.
type SessionManager struct {
	Name    string
	Session Session
	TTL     int64
}

// User is the model for a user that works on tickets.
type User struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Username    string `json:"username"`
	Mail        string `json:"mail"`
	Hash        string `json:"hash"`
	IsOnHoliday bool   `json:"isOnHoliday"`
}

// Data holds session and ticket data to parse
// to the web templates.
type Data struct {
	Session Session
	Tickets map[string]Ticket
	Users   map[string]User
}

// DataSingleTicket holds the session and ticket
// data for a call to a single ticket.
type DataSingleTicket struct {
	Session Session
	Ticket  Ticket
	Tickets map[string]Ticket
	Users   map[string]User
}

// Ticket represents a ticket.
type Ticket struct {
	ID       string  `json:"id"`
	Subject  string  `json:"subject"`
	Status   Status  `json:"status"`
	User     User    `json:"user"`
	Customer string  `json:"customer"`
	Entries  []Entry `json:"entries"`
	MergeTo  string  `json:"mergeTo"`
}

// Entry describes a single reply within a ticket.
type Entry struct {
	Date          time.Time `json:"id"`
	FormattedDate string    `json:"formattedDate"`
	User          string    `json:"user"`
	Text          string    `json:"text"`
	ReplyType     string    `json:"replyType"`
}

// Status is an enum to represent the current
// status of a ticket.
type Status int

const (
	// StatusOpen means the ticket is opened and has
	// no assignee.
	StatusOpen Status = iota

	// StatusInProgress means the ticket is being
	// processed by an assignee.
	StatusInProgress

	// StatusClosed means the ticket has been closed.
	StatusClosed
)

// String converts a ticket status to its
// corresponding description used in the
// outgoing mails.
func (status Status) String() string {
	switch status {
	case StatusOpen:
		return "Open"

	case StatusInProgress:
		return "In Progress"

	case StatusClosed:
		return "Closed"
	}

	return "undefined status"
}

// Mail struct holds the information for a
// received email in order to create new
// tickets or answers.
type Mail struct {
	ID      string `json:"id"`
	From    string `json:"from"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	Message string `json:"message"`
}

// JSONMap is a type for mapping JSON keys to
// JSON values. It is mostly used inside the
// Mail API.
type JSONMap map[string]interface{}

// Command represents a command-line interface
// command. This can be either Fetch, Submit
// or Exit.
type Command int

const (
	// CommandFetch is the CLI command to fetch emails
	// from the server.
	CommandFetch Command = iota

	// CommandSubmit is the CLI command to submit an
	// email to the server.
	CommandSubmit

	// CommandExit is the CLI command to exit the CLI.
	CommandExit
)

// String converts a command to its string
// representation.
func (c Command) String() string {
	return strconv.Itoa(int(c))
}
