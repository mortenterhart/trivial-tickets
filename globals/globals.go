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

// Globals holds global variables for easy access and prevention of circle imports

// Holds all the tickets
var Tickets = make(map[string]structs.Ticket)

// Holds all currently cached mails
var Mails = make(map[string]structs.Mail)

// Holds the given config for access to the backend systems
var ServerConfig *structs.Config

// Holds all the sessions for the users
var Sessions = make(map[string]structs.SessionManager)
