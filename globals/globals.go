package globals

import "github.com/mortenterhart/trivial-tickets/structs"

/*
*
* Matrikelnummern
* 3040018
* 6694964
* 3478222
 */

// Globals holds global variables for easy access and prevention of circle imports

// Holds all the tickets
var Tickets = make(map[string]structs.Ticket)

// Holds the given config for access to the backend systems
var ServerConfig *structs.Config

// Holds all the sessions for the users
var Sessions = make(map[string]structs.SessionManager)
