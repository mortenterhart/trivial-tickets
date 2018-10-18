package server

import (
	"crypto/rand"
	"encoding/base64"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/mortenterhart/go-tickets/structs"
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

func DestroySession(w http.ResponseWriter, r *http.Request) error {

	return nil
}

func GetSession(sessionId string) (structs.Session, error) {

	return sessions[sessionId].Session, nil
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
