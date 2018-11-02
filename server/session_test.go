package server

import (
	"testing"

	"github.com/mortenterhart/trivial-tickets/structs"
	"github.com/stretchr/testify/assert"
)

/*
*
* Matrikelnummern
* 3040018
* 6694964
* 3478222
 */

// TestCreateSession makes sure a session gets created
func TestCreateSession(t *testing.T) {

	sessionManager := CreateSession("abcdef1234")

	assert.NotNil(t, sessionManager, "SessionManager has not been created in GetSession()")
	assert.Equal(t, sessionManager.Name, "abcdef1234", "Session was not created with the id specified")
}

// TestGetSession makes sure no session gets returned, if there is none
func TestGetSession(t *testing.T) {

	session, err := GetSession("abcdef1234")

	assert.Equal(t, session, structs.Session{}, "There was not empty session returned")
	assert.NotNil(t, err, "No error was given eventho the specified session id was wrong")
}

// TestCreateSessionId makes sure generated session ids have the correct
// length and a pseudo-random
func TestCreateSessionId(t *testing.T) {

	sessionId := CreateSessionId()
	sessionId2 := CreateSessionId()

	assert.True(t, (len(sessionId) == 44), "Session id has the wrong length")
	assert.True(t, (len(sessionId2) == 44), "Session is has the wrong length")
	assert.True(t, (sessionId != sessionId2), "Created the same session id twice")
}
