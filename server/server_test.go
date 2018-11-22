package server

import (
	"net/http"
	"net/http/httptest"
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

// TestGetTemplates makes sure the application is able to correctly find the templates
// with the given standard values
func TestGetTemplates(t *testing.T) {

	tmpl := GetTemplates("../www")
	tmplNil := GetTemplates("/www")

	assert.NotNil(t, tmpl, "GetTemplates() returned no found templates")
	assert.Nil(t, tmplNil, "GetTemplates() found templates where it was not supposed to be")
}

func TestRedirectToTLS(t *testing.T) {

	req, _ := http.NewRequest("GET", "localhost", nil)
	w := httptest.NewRecorder()
	handler := http.HandlerFunc(redirectToTLS)

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusTemporaryRedirect, w.Code, "The HTTP status code was incorrect")
}

func TestRedirectToTLSWithParams(t *testing.T) {

	req, _ := http.NewRequest("GET", "localhost?id=123", nil)
	w := httptest.NewRecorder()
	handler := http.HandlerFunc(redirectToTLS)

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusTemporaryRedirect, w.Code, "The HTTP status code was incorrect")
}

func TestStartHandlers(t *testing.T) {

	err := startHandlers("../www")

	assert.Nil(t, err, "An error occured, although the path was correct")
}

func TestStartServer(t *testing.T) {

}

func TestStartHandlersNoPath(t *testing.T) {

	err := startHandlers("")

	assert.NotNil(t, err, "No error occured, although the path was incorrect")
}

func TestStartServerNoUsersPath(t *testing.T) {

	config := mockConfig()
	config.Users = ""

	err := StartServer(&config)

	assert.NotNil(t, err, "No error was returned, although no users path was specified")
}

func TestStartServerNoWebPath(t *testing.T) {

	config := mockConfig()
	config.Web = ""

	err := StartServer(&config)

	assert.NotNil(t, err, "No error was returned, although no web path was specified")
}

func TestStartServerNoTicketsPath(t *testing.T) {

	config := mockConfig()
	config.Tickets = ""

	err := StartServer(&config)

	assert.NotNil(t, err, "No error was returned, although no tickets path was specified")
}

func mockConfig() structs.Config {

	return structs.Config{
		Port:    443,
		Tickets: "../files/tickets",
		Users:   "../files/users/users.json",
		Cert:    "../ssl/server.cert",
		Key:     "../ssl/server.key",
		Web:     "../www"}
}
