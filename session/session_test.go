package session

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mortenterhart/trivial-tickets/globals"
	"github.com/stretchr/testify/assert"
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

// Constant mock session id for testing
const SESSION_ID = "zSkhrZiqZ1IF6nOJTxpSKFEEOGOgiZ0pn8vKxkW-S40="

// TestCreateSessionCookie makes sure a session cookie is generated
// accordingly, along with the session id itself
func TestCreateSessionCookie(t *testing.T) {

	cookie, sessionId, errCreateSessionCookie := CreateSessionCookie()

	assert.Nil(t, errCreateSessionCookie, "The cookie could not be created")
	assert.NotNil(t, sessionId, "The returned session is was nil")
	assert.NotNil(t, cookie, "The returned cookie is was nil")
	assert.Equal(t, "session", cookie.Name, "The cookie was not named session")
	assert.True(t, (len(sessionId) == 44), "The session is has the wrong length")
}

// TestGetSessionId tests that a session id is retrievable if it is set
func TestGetSessionId(t *testing.T) {

	cookie, _, _ := CreateSessionCookie()
	rr := httptest.NewRecorder()
	http.SetCookie(rr, cookie)
	request := &http.Request{Header: http.Header{"Cookie": rr.HeaderMap["Set-Cookie"]}}
	sId := GetSessionId(request)

	assert.NotNil(t, sId, "No session id was found")
	assert.True(t, (len(sId) == 44), "Session id has the wrong length")
}

// TestGetSessionIdError produces an error to make sure the function will
// return an error if the session id does not match the requirements
func TestGetSessionIdError(t *testing.T) {

	request := &http.Request{}
	sId := GetSessionId(request)

	assert.False(t, (len(sId) == 44), "No error was returned")
}

// TestDeleteSessionCookie tests the invalidation of a session cookie by overwriting its value
func TestDeleteSessionCookie(t *testing.T) {

	cookie := DeleteSessionCookie()

	assert.NotNil(t, cookie, "Cookie was not overwritten")
	assert.Equal(t, "", cookie.Value, "Value of cookie was not emptied")
}

// TestCreateSession tests the creation of a session itself with a given session id as parameter
func TestCreateSession(t *testing.T) {

	session := CreateSession(SESSION_ID)

	assert.Equal(t, SESSION_ID, session.Session.Id, "The session id was not used to create the session")
}

// TestGetSessionIfNotExist tests to get a session which does not exist. In that case a new session
// is created
func TestGetSessionIfNotExist(t *testing.T) {

	session, errGetSession := GetSession("abc123")

	assert.NotNil(t, session, "Session was not created")
	assert.NotNil(t, errGetSession, "No error was returned, although the id does not exist")
}

// TestGetSessionIfExist tests to get a session where the session does exist prior
func TestGetSessionIfExist(t *testing.T) {

	session := CreateSession(SESSION_ID)
	globals.Sessions[SESSION_ID] = session

	session2, errGetSession := GetSession(SESSION_ID)

	assert.Equal(t, SESSION_ID, session2.Id, "The session id was not used to create the session")
	assert.Nil(t, errGetSession, "An error was returned, but the session was available")
}

// TestUpdateSession makes sure the values of a session are updated correctly
func TestUpdateSession(t *testing.T) {

	session, _ := GetSession(SESSION_ID)
	session.IsLoggedIn = true

	UpdateSession(SESSION_ID, session)

	assert.True(t, session.IsLoggedIn, "Struct value was not changed")
}

// TestCheckForSession tests the function to look for a session with and without
// the pripor existence of a session
func TestCheckForSession(t *testing.T) {

	// Create a request with and without the session cookie set
	// in order to make sure it works either way

	cookie, _, _ := CreateSessionCookie()
	rr := httptest.NewRecorder()
	http.SetCookie(rr, cookie)
	request := &http.Request{Header: http.Header{"Cookie": rr.HeaderMap["Set-Cookie"]}}

	rr2 := httptest.NewRecorder()
	session, errCheckForSession := CheckForSession(rr, request)
	session1, _ := CheckForSession(httptest.NewRecorder(), &http.Request{Header: http.Header{"Cookie": rr2.HeaderMap["Set-Cookie"]}})

	assert.NotNil(t, session, "The session was not created correctly")
	assert.True(t, !session.User.IsOnHoliday, "The session was not created correctly")
	assert.Nil(t, errCheckForSession, "There was an error, although the request was valid")
	assert.NotNil(t, session1, "The session without the cookie was not created properly")
}
