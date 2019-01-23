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
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mortenterhart/trivial-tickets/globals"
	"github.com/mortenterhart/trivial-tickets/log/testlog"
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
 * Package session [tests]
 * Session Management
 */

//revive:disable:deep-exit

// TestMain is started to run the tests and initializes the
// configuration before running the tests. The tests' exit
// status is returned as the overall exit status.
func TestMain(m *testing.M) {
	initializeLogConfig()

	os.Exit(m.Run())
}

//revive:enable:deep-exit

// initializeLogConfig initializes the global logging
// configuration with test values.
func initializeLogConfig() {
	logConfig := testLogConfig()
	globals.LogConfig = &logConfig
}

// testLogConfig returns a test logging configuration.
func testLogConfig() structs.LogConfig {
	return structs.LogConfig{
		LogLevel:  structs.LevelInfo,
		Verbose:   false,
		FullPaths: false,
	}
}

// testSessionID is a constant mock session
// id for testing.
const testSessionID string = "zSkhrZiqZ1IF6nOJTxpSKFEEOGOgiZ0pn8vKxkW-S40="

// sessionIDLength is a comparison value
// for the length of the session id.
const sessionIDLength int = 44

// setCookieHeader is the HTTP Set-Cookie Header that
// is retrieved from the responses and set to the cookies
// in the requests.
const setCookieHeader string = "Set-Cookie"

// TestCreateSessionCookie makes sure a session cookie
// is generated accordingly, along with the session id
// itself.
func TestCreateSessionCookie(t *testing.T) {
	testlog.BeginTest()
	defer testlog.EndTest()

	cookie, sessionID, errCreateSessionCookie := CreateSessionCookie()

	assert.Nil(t, errCreateSessionCookie, "The cookie could not be created")
	assert.NotNil(t, sessionID, "The returned session is was nil")
	assert.NotNil(t, cookie, "The returned cookie is was nil")
	assert.Equalf(t, CookieName, cookie.Name, "The cookie was not named \"%s\"", CookieName)
	assert.Equal(t, sessionIDLength, len(sessionID), "The session is has the wrong length")
}

// TestGetSessionId tests that a session id is retrievable
// if it is set.
func TestGetSessionId(t *testing.T) {
	testlog.BeginTest()
	defer testlog.EndTest()

	cookie, _, _ := CreateSessionCookie()
	recorder := httptest.NewRecorder()
	http.SetCookie(recorder, cookie)
	request := &http.Request{
		Header: http.Header{
			"Cookie": recorder.Result().Header[setCookieHeader],
		},
	}

	sessionID := GetSessionID(request)

	assert.NotNil(t, sessionID, "No session id was found")
	assert.Equal(t, sessionIDLength, len(sessionID), "Session id has the wrong length")
}

// TestGetSessionIdError produces an error to make sure
// the function will return an error if the session id
// does not match the requirements.
func TestGetSessionIdError(t *testing.T) {
	testlog.BeginTest()
	defer testlog.EndTest()

	request := &http.Request{}
	sessionID := GetSessionID(request)

	assert.NotEqual(t, sessionIDLength, len(sessionID), "No error was returned")
}

// TestDeleteSessionCookie tests the invalidation of a
// session cookie by overwriting its value.
func TestDeleteSessionCookie(t *testing.T) {
	testlog.BeginTest()
	defer testlog.EndTest()

	cookie := DeleteSessionCookie()

	assert.NotNil(t, cookie, "Cookie was not overwritten")
	assert.Empty(t, cookie.Value, "Value of cookie was not emptied")
}

// TestCreateSession tests the creation of a session
// itself with a given session id as parameter.
func TestCreateSession(t *testing.T) {
	testlog.BeginTest()
	defer testlog.EndTest()

	session := CreateSession(testSessionID)

	assert.Equal(t, testSessionID, session.Session.ID, "The session id was not used to create the session")
}

// TestGetSessionIfNotExist tests to get a session
// which does not exist. In that case a new session
// is created.
func TestGetSessionIfNotExist(t *testing.T) {
	testlog.BeginTest()
	defer testlog.EndTest()

	session, errGetSession := GetSession("abc123")

	assert.NotNil(t, session, "Session was not created")
	assert.NotNil(t, errGetSession, "No error was returned, although the id does not exist")
}

// TestGetSessionIfExist tests to get a session where
// the session does exist prior.
func TestGetSessionIfExist(t *testing.T) {
	testlog.BeginTest()
	defer testlog.EndTest()

	session := CreateSession(testSessionID)
	globals.Sessions[testSessionID] = session

	session2, errGetSession := GetSession(testSessionID)

	assert.Equal(t, testSessionID, session2.ID, "The session id was not used to create the session")
	assert.Nil(t, errGetSession, "An error was returned, but the session was available")
}

// TestUpdateSession makes sure the values of a session
// are updated correctly.
func TestUpdateSession(t *testing.T) {
	testlog.BeginTest()
	defer testlog.EndTest()

	session, _ := GetSession(testSessionID)
	session.IsLoggedIn = true

	UpdateSession(testSessionID, session)

	assert.True(t, session.IsLoggedIn, "Struct value was not changed")
}

// TestCheckForSession tests the function to look for a
// session with and without the prior existence of a
// session.
func TestCheckForSession(t *testing.T) {
	testlog.BeginTest()
	defer testlog.EndTest()

	// Create a request with and without the session cookie set
	// in order to make sure it works either way

	cookie, _, _ := CreateSessionCookie()
	rr := httptest.NewRecorder()
	http.SetCookie(rr, cookie)
	request := &http.Request{
		Header: http.Header{
			"Cookie": rr.Result().Header[setCookieHeader],
		},
	}

	rr2 := httptest.NewRecorder()
	session, errCheckForSession := CheckForSession(rr, request)
	session1, _ := CheckForSession(httptest.NewRecorder(), &http.Request{
		Header: http.Header{
			"Cookie": rr2.Result().Header[setCookieHeader],
		},
	})

	assert.NotNil(t, session, "The session was not created correctly")
	assert.False(t, session.User.IsOnHoliday, "The session was not created correctly")
	assert.Nil(t, errCheckForSession, "There was an error, although the request was valid")
	assert.NotNil(t, session1, "The session without the cookie was not created properly")
}
