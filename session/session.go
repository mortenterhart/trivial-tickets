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

// Package session utilizes a session management for
// registered users.
package session

import (
	"crypto/rand"
	"encoding/base64"
	"io"
	"net/http"
	"time"

	"github.com/pkg/errors"

	"github.com/mortenterhart/trivial-tickets/globals"
	"github.com/mortenterhart/trivial-tickets/log"
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
 * Package session
 * Session Management
 */

// CookieName is the used name for the
// session cookie.
const CookieName string = "session"

// CookieTTL is the used time to live for
// the session cookie.
const CookieTTL int64 = 3600

// CreateSession returns a new session manager struct with an empty user
// and the given session id.
func CreateSession(sessionID string) structs.SessionManager {

	session := structs.Session{
		User:         structs.User{},
		CreationTime: time.Now(),
		IsLoggedIn:   false,
		ID:           sessionID,
	}

	return structs.SessionManager{
		Name:    sessionID,
		Session: session,
		TTL:     CookieTTL,
	}
}

// GetSession retrieves a session from the global session map
// with a given session id.
func GetSession(sessionID string) (structs.Session, error) {

	session := globals.Sessions[sessionID].Session

	if session != (structs.Session{}) {
		return globals.Sessions[sessionID].Session, nil
	}

	return structs.Session{}, errors.New("unable to find session with id: " + sessionID)
}

// UpdateSession updates a session manager struct with the
// given session id with a given session struct.
func UpdateSession(sessionID string, session structs.Session) {

	globals.Sessions[sessionID] = structs.SessionManager{
		Name:    sessionID,
		Session: session,
		TTL:     CookieTTL,
	}
}

// CreateSessionID generates a pseudo random id for session
// with a length of 32 characters and returns it as a base64
// encoded string.
func CreateSessionID() (string, error) {

	const length int = 32

	sessionID := make([]byte, length)

	l, errRnd := io.ReadFull(rand.Reader, sessionID)

	if errRnd != nil && l == length {
		return "", errors.New("unable to create session id")
	}

	return base64.URLEncoding.EncodeToString(sessionID), nil
}

// CreateSessionCookie returns a http cookie to hold the session
// id for the user.
func CreateSessionCookie() (*http.Cookie, string, error) {

	sessionID, errSessionID := CreateSessionID()

	if errSessionID != nil {
		return nil, "", errors.New("unable to create session id")
	}

	return &http.Cookie{
		Name:     CookieName,
		Value:    sessionID,
		HttpOnly: false,
		Expires:  time.Now().Add(2 * time.Hour),
	}, sessionID, nil
}

// DeleteSessionCookie returns a http cookie which will overwrite the
// existing session cookie in order to nullify it.
func DeleteSessionCookie() *http.Cookie {

	return &http.Cookie{
		Name:     CookieName,
		Value:    "",
		HttpOnly: false,
		Expires:  time.Now().Add(-100 * time.Hour)}
}

// GetSessionID retrieves the session id from the cookie of the user.
func GetSessionID(r *http.Request) string {

	// Get the cookie with the session id
	userCookie, errUserCookie := r.Cookie(CookieName)

	if errUserCookie != nil {
		log.Error(errUserCookie)
		return errUserCookie.Error()
	}

	return userCookie.Value
}

// CheckForSession either returns a new session or the existing
// session of a user.
func CheckForSession(w http.ResponseWriter, r *http.Request) (structs.Session, error) {

	var newSession structs.Session

	// Check if the user already has a session
	// If not, create one
	// Otherwise read the session id and load the index with his session
	if _, err := r.Cookie(CookieName); err != nil {

		cookie, sessionID, errCreateSessionCookie := CreateSessionCookie()

		if errCreateSessionCookie != nil {
			return structs.Session{}, errors.Wrap(errCreateSessionCookie, "unable to create a session cookie")
		}

		http.SetCookie(w, cookie)
		globals.Sessions[sessionID] = CreateSession(sessionID)

		newSession = globals.Sessions[sessionID].Session

	} else {
		sessionID := GetSessionID(r)

		newSession = globals.Sessions[sessionID].Session
	}

	return newSession, nil
}
