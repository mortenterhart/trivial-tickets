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

// Package globals contains global settings and resources for
// the server and logger.
package globals

import (
	"github.com/mortenterhart/trivial-tickets/structs"
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
 * Package globals
 * Global hash maps and structs
 */

// Globals holds global variables for easy access
// and prevention of circle imports.

// Tickets holds all the created tickets.
var Tickets = make(map[string]structs.Ticket)

// Mails holds all currently cached mails.
var Mails = make(map[string]structs.Mail)

// ServerConfig holds the given server config
// for access to the backend systems.
var ServerConfig *structs.ServerConfig

// LogConfig contains the global logging
// configuration.
var LogConfig *structs.LogConfig

// Sessions holds all the sessions for the users.
var Sessions = make(map[string]structs.SessionManager)
