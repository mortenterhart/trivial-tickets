package session

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mortenterhart/trivial-tickets/globals"
	"github.com/stretchr/testify/assert"
)

// Testing the various session functions, combined for easier testing
func TestSessionFunctions(t *testing.T) {

	// Test CreateSessionCookie
	_, sessionId, errCreateSessionCookie := CreateSessionCookie()

	assert.Nil(t, errCreateSessionCookie, "There was an error creating the cookie")

	// Test CreateSession
	session1 := CreateSession(sessionId)

	assert.Equal(t, sessionId, session1.Session.Id, "The session id was not used to create the session")

	// Test GetSession
	session2, errGetSession := GetSession(sessionId)

	assert.NotNil(t, session2, "Session was not created")
	assert.NotNil(t, errGetSession, "No error was returned, although the id does not exist")

	globals.Sessions[sessionId] = session1

	session3, errGetSession3 := GetSession(sessionId)

	assert.Equal(t, sessionId, session3.Id, "The session id was not used to create the session")
	assert.Nil(t, errGetSession3, "An error was returned, but the session was available")

	// Test UpdateSession
	session3.IsLoggedIn = true

	UpdateSession(sessionId, session3)

	assert.True(t, session3.IsLoggedIn, "Struct value was not changed")

	// Test CheckForSession
	cookie, _, _ := CreateSessionCookie()
	rr := httptest.NewRecorder()
	http.SetCookie(rr, cookie)
	request := &http.Request{Header: http.Header{"Cookie": rr.HeaderMap["Set-Cookie"]}}

	rr2 := httptest.NewRecorder()
	session4, errCheckForSession := CheckForSession(rr, request)
	session5, _ := CheckForSession(httptest.NewRecorder(), &http.Request{Header: http.Header{"Cookie": rr2.HeaderMap["Set-Cookie"]}})

	assert.NotNil(t, session4, "The session was not created correctly")
	assert.True(t, !session4.User.IsOnHoliday, "The session was not created correctly")
	assert.Nil(t, errCheckForSession, "There was an error, although the request was valid")
	assert.NotNil(t, session5, "The session without the cookie was not created properly")

}

// TestCookieFunctions tests all cookie related functions including errors
func TestCookieFunctions(t *testing.T) {

	// Test creating a new session cookie
	cookie, sessionId, errCreateSessionCookie := CreateSessionCookie()

	assert.Nil(t, errCreateSessionCookie, "The cookie could not be created")
	assert.NotNil(t, sessionId, "The returned session is was nil")
	assert.NotNil(t, cookie, "The returned cookie is was nil")
	assert.Equal(t, "session", cookie.Name, "The cookie was not named session")
	assert.True(t, (len(sessionId) == 44), "The session is has the wrong length")

	// Test to retrieve the session id from a set cookie
	rr := httptest.NewRecorder()
	http.SetCookie(rr, cookie)
	request := &http.Request{Header: http.Header{"Cookie": rr.HeaderMap["Set-Cookie"]}}
	sId := GetSessionId(request)

	assert.NotNil(t, sId, "No session id was found")
	assert.True(t, (len(sId) == 44), "Session id has the wrong length")

	// Test that an error string is returned, when there is no cookie
	request2 := &http.Request{}
	sId2 := GetSessionId(request2)

	assert.False(t, (len(sId2) == 44), "No error was returned")

	// Overwrite the session cookie, so that it is deleted
	cookie2 := DeleteSessionCookie()

	assert.NotNil(t, cookie2, "Cookie was not overwritten")
	assert.Equal(t, "", cookie2.Value, "Value of cookie was not emptied")
}

// TestCheckForSession creates a request to test the creation of a session for a user
func TestCheckForSession(t *testing.T) {

}
