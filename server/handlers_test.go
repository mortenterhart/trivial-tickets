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

// Package server implements the web server including
// shutdown routines and the associated handlers for
// web requests.
package server

import (
	"fmt"
	"html/template"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/mortenterhart/trivial-tickets/globals"
	"github.com/mortenterhart/trivial-tickets/logger/testlogger"
	"github.com/mortenterhart/trivial-tickets/structs"
	"github.com/mortenterhart/trivial-tickets/structs/defaults"
	"github.com/mortenterhart/trivial-tickets/util/filehandler"
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
 * Package server [tests]
 * Server handlers reacting to HTTP requests
 */

// newNonRedirectClient returns a new HTTP client that
// does not follow any redirects. If the client faces
// a 301 Moved Permanently response code with the
// Location header set the request is not automatically
// redirected to that location. This behavior is important
// for the handler tests because they return the 301
// response code for automatic redirects, but a client
// permits a maximum of 10 redirects. Since the used
// test servers use the root path '/' as handler pattern
// the request is always redirected to the same handler.
func newNonRedirectClient() *http.Client {
	return &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
}

// testServerConfig returns a test configuration for the
// server that can be used in handler tests. The path to
// the users file is also other than the default path causing
// that the directory in which the file will be written
// does not exist. If the users file is needed (such as
// in the handler HandleHoliday) the directory has to
// be created first.
func testServerConfig() structs.ServerConfig {
	return structs.ServerConfig{
		Port:    defaults.TestPort,
		Tickets: "../files/testtickets",
		Users:   "../files/testusers/users.json",
		Mails:   "../files/testmails",
		Cert:    defaults.TestCertificate,
		Key:     defaults.TestKey,
		Web:     defaults.TestWeb,
	}
}

// cleanupTestFiles removes all test tickets, mails and users
// from the paths in the given config, if they exist. It reports
// an error if a directory could not be removed.
func cleanupTestFiles(config structs.ServerConfig) {
	if filehandler.DirectoryExists(config.Tickets) {
		testlogger.Debug("Deferred: Removing test ticket directory")
		if removeErr := os.RemoveAll(config.Tickets); removeErr != nil {
			testlogger.Debug("ERROR: cannot remove test ticket directory:", removeErr)
		}
	}

	if filehandler.DirectoryExists(filepath.Dir(config.Users)) {
		testlogger.Debug("Deferred: Removing test users directory")
		if removeErr := os.RemoveAll(filepath.Dir(config.Users)); removeErr != nil {
			testlogger.Debug("ERROR: cannot remove test users directory:", removeErr)
		}
	}

	if filehandler.DirectoryExists(config.Mails) {
		testlogger.Debug("Deferred: Removing test mail directory")
		if removeErr := os.RemoveAll(config.Mails); removeErr != nil {
			testlogger.Debug("ERROR: cannot remove test mail directory:", removeErr)
		}
	}
}

// resetConfig resets the server and logging configuration
// to their initial state.
func resetConfig() {
	initializeConfig()
}

// initializeConfig assigns default values to the global
// server and logging configuration.
func initializeConfig() {
	serverConfig := testServerConfig()
	globals.ServerConfig = &serverConfig

	logConfig := mockLogConfig()
	globals.LogConfig = &logConfig
}

// TestMain is the superior test routine that is run to start the tests.
// It setups the logging configuration because the logger is used during
// the tested handlers.
func TestMain(m *testing.M) {
	initializeConfig()

	os.Exit(m.Run())
}

/*
 * To test the handlers, the ServeHTTP function was mapped to a mock struct in
 * order to call them directly via the test server.
 *
 * Various mock structs and global variables are populated to make the tests work properly.
 */

// index
// --------------------------
type indexHandler struct{}

func (h *indexHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handleIndex(w, r)
}

func TestHandleIndex(t *testing.T) {
	testlogger.BeginTest()
	defer testlogger.EndTest()

	defer resetConfig()

	config := testServerConfig()
	defer cleanupTestFiles(config)

	tmpl = getTemplates("../www")

	handler := &indexHandler{}

	globals.Tickets["abc123"] = structs.Ticket{}
	users["abc123"] = structs.User{}

	server := httptest.NewServer(handler)
	defer server.Close()

	client := newNonRedirectClient()
	resp, err := client.Get(server.URL)
	defer func() {
		if err == nil {
			resp.Body.Close()
		}
	}()

	assert.Nil(t, err, "There was an unexpected error")
	assert.Equal(t, http.StatusOK, resp.StatusCode, "The http status is not 200")
}

// login
// --------------------------

type loginHandler struct{}

func (h *loginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handleLogin(w, r)
}

func TestHandleLogin(t *testing.T) {
	testlogger.BeginTest()
	defer testlogger.EndTest()

	defer resetConfig()

	config := testServerConfig()
	defer cleanupTestFiles(config)

	handler := &loginHandler{}

	server := httptest.NewServer(handler)
	defer server.Close()

	users["testuser"] = structs.User{
		ID:          "1",
		Name:        "Test",
		Username:    "testuser",
		Mail:        "Testuser@mail.com",
		Hash:        "$2a$12$rW6Ska0DaVjTX/8sQGCp/.y7kl2RvF.9936Hmm27HyI0cJ78q1UOG",
		IsOnHoliday: false,
	}

	reader := strings.NewReader("username=testuser&password=MyPassword123!!##")

	client := newNonRedirectClient()
	resp, err := client.Post(server.URL, "application/x-www-form-urlencoded", reader)
	defer func() {
		if err == nil {
			resp.Body.Close()
		}
	}()

	assert.Equal(t, http.StatusMovedPermanently, resp.StatusCode, "Status code did not match 301")
}

// logout
// --------------------------

type logoutHandler struct{}

func (h *logoutHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handleLogout(w, r)
}

func TestHandleLogout(t *testing.T) {
	testlogger.BeginTest()
	defer testlogger.EndTest()

	defer resetConfig()

	config := testServerConfig()
	defer cleanupTestFiles(config)

	handler := &logoutHandler{}
	server := httptest.NewServer(handler)
	defer server.Close()

	client := newNonRedirectClient()
	resp, err := client.Post(server.URL, "application/x-www-form-urlencoded", nil)
	defer func() {
		if err == nil {
			resp.Body.Close()
		}
	}()

	assert.Equal(t, http.StatusMovedPermanently, resp.StatusCode, "Status code did not match 301")
}

// create ticket
// --------------------------

type createTicketHandler struct{}

func (h *createTicketHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handleCreateTicket(w, r)
}

func TestHandleCreateTicket(t *testing.T) {
	testlogger.BeginTest()
	defer testlogger.EndTest()

	defer resetConfig()

	config := testServerConfig()
	defer cleanupTestFiles(config)

	handler := &createTicketHandler{}

	server := httptest.NewServer(handler)
	defer server.Close()

	reader := strings.NewReader("mail=testuser@test.com&subject=help&text=testest")

	client := newNonRedirectClient()
	postResponse, postErr := client.Post(server.URL, "application/x-www-form-urlencoded", reader)
	defer func() {
		if postErr == nil {
			postResponse.Body.Close()
		}
	}()

	assert.Equal(t, http.StatusMovedPermanently, postResponse.StatusCode, "Status code did not match 301")

	getResponse, getErr := client.Get(server.URL)
	defer func() {
		if getErr == nil {
			getResponse.Body.Close()
		}
	}()

	assert.Equal(t, http.StatusMovedPermanently, getResponse.StatusCode, "Status code did not match 301")
}

// holiday
// --------------------------

type holidayHandler struct{}

func (h *holidayHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	cookie := http.Cookie{
		Name:  "session",
		Value: "def123",
	}
	r.AddCookie(&cookie)
	handleHoliday(w, r)
}

func TestHandleHoliday(t *testing.T) {
	testlogger.BeginTest()
	defer testlogger.EndTest()

	defer resetConfig()

	config := testServerConfig()
	defer cleanupTestFiles(config)

	handler := &holidayHandler{}
	server := httptest.NewServer(handler)
	defer server.Close()

	testUser := structs.User{
		ID:          "1",
		Name:        "Test",
		Username:    "testuser",
		Mail:        "Testuser@mail.com",
		Hash:        "$2a$12$rW6Ska0DaVjTX/8sQGCp/.y7kl2RvF.9936Hmm27HyI0cJ78q1UOG",
		IsOnHoliday: false,
	}

	users["testuser"] = testUser

	globals.Sessions["def123"] = structs.SessionManager{
		Name: "def123",
		Session: structs.Session{
			User:       testUser,
			CreateTime: time.Now(),
			IsLoggedIn: true,
			ID:         "def123",
		},
		TTL: 3600,
	}

	userDirectory := filepath.Dir(config.Users)
	filehandler.CreateFolders(userDirectory)

	client := newNonRedirectClient()
	resp, err := client.Post(server.URL, "application/x-www-form-urlencoded", nil)
	defer func() {
		if err == nil {
			resp.Body.Close()
		}
	}()

	assert.Equal(t, http.StatusMovedPermanently, resp.StatusCode, "Status code did not match 301")

	testUser.IsOnHoliday = true

	resp, err = client.Post(server.URL, "application/x-www-form-urlencoded", nil)
	defer func() {
		if err == nil {
			resp.Body.Close()
		}
	}()

	assert.Equal(t, http.StatusMovedPermanently, resp.StatusCode, "Status code did not match 301")
}

// ticket
// --------------------------

type ticketHandler struct{}

func (h *ticketHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handleTicket(w, r)
}

func TestHandleTicket(t *testing.T) {
	testlogger.BeginTest()
	defer testlogger.EndTest()

	defer resetConfig()

	config := testServerConfig()
	defer cleanupTestFiles(config)

	tmpl = getTemplates("../www")

	handler := &ticketHandler{}
	server := httptest.NewServer(handler)
	defer server.Close()

	ticket := structs.Ticket{
		ID: "abc123",
	}

	globals.Tickets["abc123"] = ticket

	client := newNonRedirectClient()
	resp, err := client.Get(server.URL + "/ticket?id=abc123")
	defer func() {
		if err == nil {
			resp.Body.Close()
		}
	}()

	assert.Nil(t, err, "There was an unexpected error")
	assert.Equal(t, http.StatusOK, resp.StatusCode, "The http response is wrong")
}

func TestHandleTicketWithMergeTo(t *testing.T) {
	testlogger.BeginTest()
	defer testlogger.EndTest()

	defer resetConfig()

	config := testServerConfig()
	defer cleanupTestFiles(config)

	tmpl = getTemplates("../www")

	handler := &ticketHandler{}
	server := httptest.NewServer(handler)
	defer server.Close()

	ticket := structs.Ticket{
		ID:      "abc123",
		MergeTo: "def123",
	}

	ticket2 := structs.Ticket{
		ID: "def123",
	}

	globals.Tickets["def123"] = ticket2
	globals.Tickets["abc123"] = ticket

	client := newNonRedirectClient()
	resp, err := client.Get(server.URL + "/ticket?id=abc123")
	defer func() {
		if err == nil {
			resp.Body.Close()
		}
	}()

	assert.Nil(t, err, "There was an unexpected error")
	assert.Equal(t, http.StatusOK, resp.StatusCode, "The http response is wrong")
}

func TestHandleTicketMissingIdParameter(t *testing.T) {
	testlogger.BeginTest()
	defer testlogger.EndTest()

	defer resetConfig()

	config := testServerConfig()
	defer cleanupTestFiles(config)

	tmpl = getTemplates("../www")

	handler := &ticketHandler{}
	server := httptest.NewServer(handler)
	defer server.Close()

	ticket := structs.Ticket{
		ID: "abc123",
	}

	globals.Tickets["abc123"] = ticket

	client := newNonRedirectClient()
	resp, err := client.Get(server.URL + "/ticket")
	defer func() {
		if err == nil {
			resp.Body.Close()
		}
	}()

	assert.Nil(t, err, "There was an unexpected error")
	assert.Equal(t, http.StatusMovedPermanently, resp.StatusCode, "The http response is wrong")
}

func TestHandleTicketRedirect(t *testing.T) {
	testlogger.BeginTest()
	defer testlogger.EndTest()

	defer resetConfig()

	config := testServerConfig()
	defer cleanupTestFiles(config)

	tmpl = getTemplates("../www")

	handler := &ticketHandler{}
	server := httptest.NewServer(handler)
	defer server.Close()

	ticket := structs.Ticket{
		ID: "abc123",
	}

	globals.Tickets["abc123"] = ticket

	client := newNonRedirectClient()
	resp, err := client.Post(server.URL+"/ticket", "application/x-www-form-urlencoded", strings.NewReader("id=abc123"))
	defer func() {
		if err == nil {
			resp.Body.Close()
		}
	}()

	assert.Nil(t, err, "There was an unexpected error")
	assert.Equal(t, http.StatusMovedPermanently, resp.StatusCode, "The http response is wrong")
}

func TestHandleTicketExecuteTemplateError(t *testing.T) {
	testlogger.BeginTest()
	defer testlogger.EndTest()

	defer resetConfig()

	config := testServerConfig()
	defer cleanupTestFiles(config)

	handler := &ticketHandler{}
	server := httptest.NewServer(handler)
	defer server.Close()

	// Define a new ticket template to replace the
	// existing template. The new template will cause
	// a template execution error.
	tmpl, _ = template.New("ticket.html").Parse("{{range $replies := .Ticket.Entries}} <p> {{$replies.Text}} </p> {{end}}")

	client := newNonRedirectClient()
	resp, err := client.Get(server.URL + "/ticket?id=abc123")
	defer func() {
		if err == nil {
			resp.Body.Close()
		}
	}()

	assert.Nil(t, err, "There was an unexpected error")
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode, "The http response is wrong")
}

// update ticket
// --------------------------

type updateTicketHandler struct{}

func (h *updateTicketHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handleUpdateTicket(w, r)
}

func TestHandleUpdateTicketWrongMethod(t *testing.T) {
	testlogger.BeginTest()
	defer testlogger.EndTest()

	defer resetConfig()

	config := testServerConfig()
	defer cleanupTestFiles(config)

	handler := &updateTicketHandler{}
	server := httptest.NewServer(handler)
	defer server.Close()

	client := newNonRedirectClient()
	resp, err := client.Get(server.URL)
	fmt.Println(err)
	defer func() {
		if err == nil {
			resp.Body.Close()
		}
	}()

	assert.Equal(t, http.StatusMovedPermanently, resp.StatusCode, "HTTP status code did not match 301")
}

func TestHandleUpdateTicketMerge(t *testing.T) {
	testlogger.BeginTest()
	defer testlogger.EndTest()

	defer resetConfig()

	config := testServerConfig()
	defer cleanupTestFiles(config)

	tmpl = getTemplates("../www")

	handler := &updateTicketHandler{}
	server := httptest.NewServer(handler)
	defer server.Close()

	reader := strings.NewReader("ticketId=1&status=0&mail=bla@example.com&reply=hallo&reply_type=intern&merge=2")

	client := newNonRedirectClient()
	resp, err := client.Post(server.URL, "application/x-www-form-urlencoded", reader)
	defer func() {
		if err == nil {
			resp.Body.Close()
		}
	}()

	assert.Nil(t, err, "there was an error in POST request")
	assert.Equal(t, http.StatusOK, resp.StatusCode, "status code is not 200 OK")
}

func TestHandleUpdateTicketExtern(t *testing.T) {
	testlogger.BeginTest()
	defer testlogger.EndTest()

	defer resetConfig()

	config := testServerConfig()
	defer cleanupTestFiles(config)

	tmpl = getTemplates("../www")

	handler := &updateTicketHandler{}
	server := httptest.NewServer(handler)
	defer server.Close()

	reader := strings.NewReader("ticketId=1&status=0&mail=bla@example.com&reply=hallo&reply_type=extern")

	client := newNonRedirectClient()
	resp, err := client.Post(server.URL, "application/x-www-form-urlencoded", reader)
	defer func() {
		if err == nil {
			resp.Body.Close()
		}
	}()

	assert.Nil(t, err, "There was an error")
	assert.Equal(t, http.StatusOK, resp.StatusCode, "Wrong http status code")
}

func TestHandleUpdateTicketExecuteTemplateError(t *testing.T) {
	testlogger.BeginTest()
	defer testlogger.EndTest()

	defer resetConfig()

	config := testServerConfig()
	defer cleanupTestFiles(config)

	handler := &updateTicketHandler{}
	server := httptest.NewServer(handler)
	defer server.Close()

	// Define a new ticket template to replace the
	// existing template. The new template will cause
	// a template execution error.
	tmpl, _ = template.New("ticket.html").Parse("{{range $replies := .Ticket.Entries}} <p> {{$replies.Text}} </p> {{end}}")

	reader := strings.NewReader("ticketId=1&status=0&mail=bla@example.com&reply=hallo&reply_type=extern")

	client := newNonRedirectClient()
	resp, err := client.Post(server.URL, "application/x-www-form-urlencoded", reader)
	defer func() {
		if err == nil {
			resp.Body.Close()
		}
	}()

	assert.Nil(t, err, "There was an unexpected error")
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode, "The http response is wrong")
}

// unassign ticket
// --------------------------

type unassignTicketHandler struct{}

func (h *unassignTicketHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handleUnassignTicket(w, r)
}

func TestHandleUnassignTicket(t *testing.T) {
	testlogger.BeginTest()
	defer testlogger.EndTest()

	defer resetConfig()

	config := testServerConfig()
	defer cleanupTestFiles(config)

	handler := &unassignTicketHandler{}
	server := httptest.NewServer(handler)
	defer server.Close()

	ticket := structs.Ticket{
		ID: "abc123",
	}

	globals.Tickets["abc123"] = ticket

	client := newNonRedirectClient()
	resp, err := client.Get(server.URL + "/unassignTicket?id=abc123")
	defer func() {
		if err == nil {
			resp.Body.Close()
		}
	}()

	assert.Nil(t, err, "An unexpected error occurred")
	assert.Equal(t, http.StatusOK, resp.StatusCode, "The http status is not 200")
}

func TestHandleUnassignTicketMissingIdParameter(t *testing.T) {
	testlogger.BeginTest()
	defer testlogger.EndTest()

	defer resetConfig()

	config := testServerConfig()
	defer cleanupTestFiles(config)

	handler := &unassignTicketHandler{}
	server := httptest.NewServer(handler)
	defer server.Close()

	ticket := structs.Ticket{
		ID: "abc123",
	}

	globals.Tickets["abc123"] = ticket

	client := newNonRedirectClient()
	resp, err := client.Get(server.URL + "/unassignTicket")
	defer func() {
		if err == nil {
			resp.Body.Close()
		}
	}()

	assert.Nil(t, err, "An unexpected error occurred")
	assert.Equal(t, http.StatusMovedPermanently, resp.StatusCode, "The http status is not 302")
}

// assign ticket
// --------------------------

type assignTicketHandler struct{}

func (h *assignTicketHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	cookie := http.Cookie{
		Name:  "session",
		Value: "def123",
	}
	r.AddCookie(&cookie)

	handleAssignTicket(w, r)
}

func TestHandleAssignTicket(t *testing.T) {
	testlogger.BeginTest()
	defer testlogger.EndTest()

	defer resetConfig()

	config := testServerConfig()
	defer cleanupTestFiles(config)

	users["testuser"] = structs.User{
		ID:          "1",
		Name:        "Test",
		Username:    "testuser",
		Mail:        "Testuser@mail.com",
		Hash:        "$2a$12$rW6Ska0DaVjTX/8sQGCp/.y7kl2RvF.9936Hmm27HyI0cJ78q1UOG",
		IsOnHoliday: false,
	}

	globals.Sessions["def123"] = structs.SessionManager{
		Name: "def123",
		Session: structs.Session{
			User:       users["testuser"],
			CreateTime: time.Now(),
			IsLoggedIn: true,
			ID:         "def123",
		},
		TTL: 3600,
	}

	handler := &assignTicketHandler{}
	server := httptest.NewServer(handler)
	defer server.Close()

	client := newNonRedirectClient()
	resp, err := client.Get(server.URL + "/assignTicket?id=abc123&user=testuser")
	defer func() {
		if err == nil {
			resp.Body.Close()
		}
	}()

	assert.Nil(t, err, "An unexpected error occurred")
	assert.Equal(t, http.StatusOK, resp.StatusCode, "The http status is not 200")
}
