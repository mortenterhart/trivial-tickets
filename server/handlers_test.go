package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

/*
*
* Matrikelnummern
* 3040018
* 6694964
* 3478222
 */

// TestCreateTicketId makes sure the created ticket id is in line with the specification
func TestCreateTicketId(t *testing.T) {

	ticketId := createTicketId(10)

	assert.True(t, (len(ticketId) == 10), "Ticket id has the wrong length")
}

// TestCookieFunctions tests all cookie related functions including errors
func TestCookieFunctions(t *testing.T) {

	// Test creating a new session cookie
	cookie, sessionId := createSessionCookie()

	assert.NotNil(t, sessionId, "The returned session is was nil")
	assert.NotNil(t, cookie, "The returned cookie is was nil")
	assert.Equal(t, cookie.Name, "session", "The cookie was not named session")
	assert.True(t, (len(sessionId) == 44), "The session is has the wrong length")

	// Test to retrieve the session id from a set cookie
	rr := httptest.NewRecorder()
	http.SetCookie(rr, cookie)
	request := &http.Request{Header: http.Header{"Cookie": rr.HeaderMap["Set-Cookie"]}}
	sId := getSessionId(request)

	assert.NotNil(t, sId, "No session id was found")
	assert.True(t, (len(sId) == 44), "Session id has the wrong length")

	// Test that an error string is returned, when there is no cookie
	request2 := &http.Request{}
	sId2 := getSessionId(request2)

	assert.False(t, (len(sId2) == 44), "No error was returned")

	// Overwrite the session cookie, so that it is deleted
	cookie2 := deleteSessionCookie()

	assert.NotNil(t, cookie2, "Cookie was not overwritten")
	assert.Equal(t, cookie2.Value, "", "Value of cookie was not emptied")
}

// TestCheckForSession creates a request to test the creation of a session for a user
func TestCheckForSession(t *testing.T) {

	// Create a mock request
	cookie, _ := createSessionCookie()
	rr := httptest.NewRecorder()
	http.SetCookie(rr, cookie)
	request := &http.Request{Header: http.Header{"Cookie": rr.HeaderMap["Set-Cookie"]}}

	// Get session with mock request
	session := checkForSession(rr, request)

	session2 := checkForSession(httptest.NewRecorder(), &http.Request{Header: http.Header{"Cookie": rr.HeaderMap["Set-Cookie"]}})

	assert.NotNil(t, session, "The session was not created correctly")
	assert.True(t, !session.User.IsOnHoliday, "The session was not created correctly")
	assert.NotNil(t, session2, "The session without the cookie was not created properly")
}
