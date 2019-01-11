package session

import (
    "crypto/rand"
    "encoding/base64"
    "errors"
    "io"
    "log"
    "net/http"
    "time"

    "github.com/mortenterhart/trivial-tickets/globals"
    "github.com/mortenterhart/trivial-tickets/structs"
)

/*
*
* Matrikelnummern
* 3040018
* 6694964
* 3478222
 */

// CreateSession returns a new session manager struct with an empty user
// and the given session id
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

// GetSession retrieves a session from the global session map
// with a given session id
func GetSession(sessionId string) (structs.Session, error) {

    session := globals.Sessions[sessionId].Session

    if session != (structs.Session{}) {
        return globals.Sessions[sessionId].Session, nil
    } else {
        return structs.Session{}, errors.New("Unable to find session with id: " + sessionId)
    }
}

// UpdateSession updates a session manager struct with the given session id
// with a given session struct
func UpdateSession(sessionId string, session structs.Session) {

    globals.Sessions[sessionId] = structs.SessionManager{
        Name:    sessionId,
        Session: session,
        TTL:     3600,
    }
}

// CreateSessionId  generates a pseudo random id for session
// with a length of 32 characters and returns it as a base64
// encoded string
func CreateSessionId() (string, error) {

    const length = 32

    sessionId := make([]byte, length)

    l, errRnd := io.ReadFull(rand.Reader, sessionId)

    if errRnd != nil && l == length {
        return "", errors.New("Unable to create session id")
    }

    return base64.URLEncoding.EncodeToString(sessionId), nil
}

// createSessionCookie returns a http cookie to hold the session
// id for the user
func CreateSessionCookie() (*http.Cookie, string, error) {

    sessionId, errSessionId := CreateSessionId()

    if errSessionId != nil {
        return nil, "", errors.New("Unable to create session id")
    }

    return &http.Cookie{
        Name:     "session",
        Value:    sessionId,
        HttpOnly: false,
        Expires:  time.Now().Add(2 * time.Hour)}, sessionId, nil
}

// deleteSessionCookie returns a http cookie which will overwrite the
// existing session cookie in order to nulify it
func DeleteSessionCookie() *http.Cookie {

    return &http.Cookie{
        Name:     "session",
        Value:    "",
        HttpOnly: false,
        Expires:  time.Now().Add(-100 * time.Hour)}
}

// getSessionId retrieves the session id from the cookie of the user
func GetSessionId(r *http.Request) string {

    // Get the cookie with the session id
    userCookie, errUserCookie := r.Cookie("session")

    if errUserCookie != nil {
        log.Print(errUserCookie)
        return errUserCookie.Error()
    }

    return userCookie.Value
}

// checkForSession either returns a new session or the existing session of a user.
func CheckForSession(w http.ResponseWriter, r *http.Request) (structs.Session, error) {

    var newSession structs.Session

    // Check if the user already has a session
    // If not, create one
    // Otherwise read the session id and load the index with his session
    if _, err := r.Cookie("session"); err != nil {

        cookie, sessionId, errCreateSessionCookie := CreateSessionCookie()

        if errCreateSessionCookie != nil {
            return structs.Session{}, errors.New("Unable to create a session cookie")
        }

        http.SetCookie(w, cookie)
        globals.Sessions[sessionId] = CreateSession(sessionId)

        newSession = globals.Sessions[sessionId].Session

    } else {
        sessionId := GetSessionId(r)

        newSession = globals.Sessions[sessionId].Session
    }

    return newSession, nil
}
