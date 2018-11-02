package server

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
	"log"
	"time"

	"github.com/mortenterhart/trivial-tickets/structs"
)

/*
*
* Matrikelnummern
* 3040018
* 6694964
* 3478222
 */

func CreateSession(sessionId string) structs.SessionManager {

	session := structs.Session{
		User:       structs.User{},
		CreateTime: time.Now(),
		IsLoggedIn: false,
		Id:         sessionId,
	}

	return structs.SessionManager{
		Name:    sessionId,
		Session: session,
		TTL:     3600,
	}
}

func GetSession(sessionId string) (structs.Session, error) {

	session := sessions[sessionId].Session

	if session != (structs.Session{}) {
		return sessions[sessionId].Session, nil
	} else {
		return structs.Session{}, errors.New("Unable to find session with id: " + sessionId)
	}
}

func UpdateSession(sessionId string, session structs.Session) {

	sessions[sessionId] = structs.SessionManager{
		Name:    sessionId,
		Session: session,
		TTL:     3600,
	}
}

func CreateSessionId() string {

	const length = 32

	sessionId := make([]byte, length)

	l, errRnd := io.ReadFull(rand.Reader, sessionId)

	if errRnd != nil && l == length {
		log.Fatal(errRnd)
	}

	return base64.URLEncoding.EncodeToString(sessionId)
}
