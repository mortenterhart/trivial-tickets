package session

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mortenterhart/trivial-tickets/globals"
	"github.com/stretchr/testify/assert"
)

const SESSION_ID = "zSkhrZiqZ1IF6nOJTxpSKFEEOGOgiZ0pn8vKxkW-S40="

func TestCreateSessionCookie(t *testing.T) {

	cookie, sessionId, errCreateSessionCookie := CreateSessionCookie()

	assert.Nil(t, errCreateSessionCookie, "The cookie could not be created")
	assert.NotNil(t, sessionId, "The returned session is was nil")
	assert.NotNil(t, cookie, "The returned cookie is was nil")
	assert.Equal(t, "session", cookie.Name, "The cookie was not named session")
	assert.True(t, (len(sessionId) == 44), "The session is has the wrong length")
}

func TestGetSessionId(t *testing.T) {

	cookie, _, _ := CreateSessionCookie()
	rr := httptest.NewRecorder()
	http.SetCookie(rr, cookie)
	request := &http.Request{Header: http.Header{"Cookie": rr.HeaderMap["Set-Cookie"]}}
	sId := GetSessionId(request)

	assert.NotNil(t, sId, "No session id was found")
	assert.True(t, (len(sId) == 44), "Session id has the wrong length")
}

func TestGetSessionIdError(t *testing.T) {

	request := &http.Request{}
	sId := GetSessionId(request)

	assert.False(t, (len(sId) == 44), "No error was returned")
}

func TestDeleteSessionCookie(t *testing.T) {

	cookie := DeleteSessionCookie()

	assert.NotNil(t, cookie, "Cookie was not overwritten")
	assert.Equal(t, "", cookie.Value, "Value of cookie was not emptied")
}

func TestCreateSession(t *testing.T) {

	session := CreateSession(SESSION_ID)

	assert.Equal(t, SESSION_ID, session.Session.Id, "The session id was not used to create the session")
}

func TestGetSessionIfNotExist(t *testing.T) {

	session, errGetSession := GetSession(SESSION_ID)

	assert.NotNil(t, session, "Session was not created")
	assert.NotNil(t, errGetSession, "No error was returned, although the id does not exist")
}

func TestGetSessionIfExist(t *testing.T) {

	session := CreateSession(SESSION_ID)
	globals.Sessions[SESSION_ID] = session

	session2, errGetSession := GetSession(SESSION_ID)

	assert.Equal(t, SESSION_ID, session2.Id, "The session id was not used to create the session")
	assert.Nil(t, errGetSession, "An error was returned, but the session was available")
}

func TestUpdateSession(t *testing.T) {

	session, _ := GetSession(SESSION_ID)
	session.IsLoggedIn = true

	UpdateSession(SESSION_ID, session)

	assert.True(t, session.IsLoggedIn, "Struct value was not changed")
}

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
