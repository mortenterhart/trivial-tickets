package server

import (
	"go-tickets/structs"
)

/*
*
* Matrikelnummern
* 3040018
* 6694964
* 3478222
 */

func CreateSession() structs.Session {

	return structs.Session{}
}

func DestroySession(session *structs.Session) bool {

	return true
}

func IsSessionValid(session *structs.Session) bool {
	return true
}
