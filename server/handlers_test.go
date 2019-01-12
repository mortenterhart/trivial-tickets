package server

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/mortenterhart/trivial-tickets/globals"
	"github.com/mortenterhart/trivial-tickets/structs"
	"github.com/stretchr/testify/assert"
)

/*
 * Ticketsystem Trivial Tickets
 *
 * Matriculation numbers: 3040018, 3040018, 3478222
 * Lecture:               Programmieren II, INF16B
 * Lecturer:              Herr Prof. Dr. Helmut Neemann
 * Institute:             Duale Hochschule Baden-Württemberg Mosbach
 *
 * ---------------
 *
 * Package server [tests]
 * Server handlers reacting to HTTP requests
 */

/*
 * To test the handlers, the ServeHTTP function was mapped to a mock struct in
 * order to call them directly via the test server.
 *
 * Various mock structs and global variables are populated to make the tests work properly
 */

// index
// --------------------------
type indexHandler struct{}

func (h *indexHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handleIndex(w, r)
}

func TestHandleIndex(t *testing.T) {

	tmpl = GetTemplates("../www")

	handler := &indexHandler{}

	globals.Tickets["abc123"] = structs.Ticket{}
	users["abc123"] = structs.User{}

	server := httptest.NewServer(handler)
	defer server.Close()

	resp, err := http.Get(server.URL)

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

	handler := &loginHandler{}

	server := httptest.NewServer(handler)
	defer server.Close()

	users["testuser"] = structs.User{
		Id:          "1",
		Name:        "Test",
		Username:    "testuser",
		Mail:        "Testuser@mail.com",
		Hash:        "$2a$12$rW6Ska0DaVjTX/8sQGCp/.y7kl2RvF.9936Hmm27HyI0cJ78q1UOG",
		IsOnHoliday: false,
	}

	reader := strings.NewReader("username=testuser&password=MyPassword123!!##")
	resp, _ := http.Post(server.URL, "application/x-www-form-urlencoded", reader)

	assert.Equal(t, http.StatusFound, resp.StatusCode, "Status code did not match 302")
}

// logout
// --------------------------

type logoutHandler struct{}

func (h *logoutHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handleLogout(w, r)
}

func TestHandleLogout(t *testing.T) {

	handler := &logoutHandler{}
	server := httptest.NewServer(handler)
	defer server.Close()

	resp, _ := http.Post(server.URL, "application/x-www-form-urlencoded", nil)

	assert.Equal(t, http.StatusFound, resp.StatusCode, "Status code did not match 302")
}

// create ticket
// --------------------------

type createTicketHandler struct{}

func (h *createTicketHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handleCreateTicket(w, r)
}

func TestHandleCreateTicket(t *testing.T) {

	handler := &createTicketHandler{}

	server := httptest.NewServer(handler)
	defer server.Close()

	config := structs.Config{Tickets: "../files/testtickets"}

	globals.ServerConfig = &config

	reader := strings.NewReader("mail=testuser@test.com&subject=help&text=testest")
	resp, _ := http.Post(server.URL, "application/x-www-form-urlencoded", reader)

	assert.Equal(t, http.StatusFound, resp.StatusCode, "Status code did not match 302")

	// Delete created ticket file
	os.RemoveAll("../files/testtickets/")
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

func TestHandleHandleHoliday(t *testing.T) {

	handler := &holidayHandler{}
	server := httptest.NewServer(handler)
	defer server.Close()

	users["testuser"] = structs.User{
		Id:          "1",
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
			Id:         "def123",
		},
		TTL: 3600,
	}

	resp, _ := http.Post(server.URL, "application/x-www-form-urlencoded", nil)

	assert.Equal(t, http.StatusFound, resp.StatusCode, "Status code did not match 302")
}

// ticket
// --------------------------

type ticketHandler struct{}

func (h *ticketHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handleTicket(w, r)
}

func TestHandleHandleTicket(t *testing.T) {

	tmpl = GetTemplates("../www")

	handler := &ticketHandler{}
	server := httptest.NewServer(handler)
	defer server.Close()

	ticket := structs.Ticket{
		Id: "abc123",
	}

	globals.Tickets["abc123"] = ticket

	resp, err := http.Get(server.URL + "/ticket?id=abc123")

	assert.Nil(t, err, "There was an unexpected error")
	assert.Equal(t, http.StatusOK, resp.StatusCode, "The http response is wrong")
}

func TestHandleHandleTicketWithMergeTo(t *testing.T) {

	tmpl = GetTemplates("../www")

	handler := &ticketHandler{}
	server := httptest.NewServer(handler)
	defer server.Close()

	ticket := structs.Ticket{
		Id:      "abc123",
		MergeTo: "def123",
	}

	ticket2 := structs.Ticket{
		Id: "def123",
	}

	globals.Tickets["def123"] = ticket2
	globals.Tickets["abc123"] = ticket

	resp, err := http.Get(server.URL + "/ticket?id=abc123")

	assert.Nil(t, err, "There was an unexpected error")
	assert.Equal(t, http.StatusOK, resp.StatusCode, "The http response is wrong")
}

// update ticket
// --------------------------

type updateTicketHandler struct{}

func (h *updateTicketHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handleUpdateTicket(w, r)
}

func TestHandleHandleUpdateTicketWrongMethod(t *testing.T) {

	handler := &updateTicketHandler{}
	server := httptest.NewServer(handler)
	defer server.Close()

	resp, _ := http.Get(server.URL)

	assert.Equal(t, http.StatusFound, resp.StatusCode, "HTTP status code did not match 302")
}

func TestHandleHandleUpdateTicketMerge(t *testing.T) {

	tmpl = GetTemplates("../www")

	config := mockConfig()
	config.Tickets = "../files/testticket"
	globals.ServerConfig = &config

	handler := &updateTicketHandler{}
	server := httptest.NewServer(handler)
	defer server.Close()

	reader := strings.NewReader("ticketId=1&status=0&mail=bla@example.com&reply=hallo&reply_type=intern&merge=2")

	resp, err := http.Post(server.URL, "application/x-www-form-urlencoded", reader)

	assert.Nil(t, err, "")
	assert.Equal(t, http.StatusOK, resp.StatusCode, "")

	os.RemoveAll("../files/testticket")
}

func TestHandleHandleUpdateTicketExtern(t *testing.T) {

	tmpl = GetTemplates("../www")

	config := mockConfig()
	config.Tickets = "../files/testticket"
	globals.ServerConfig = &config

	handler := &updateTicketHandler{}
	server := httptest.NewServer(handler)
	defer server.Close()

	reader := strings.NewReader("ticketId=1&status=0&mail=bla@example.com&reply=hallo&reply_type=extern")

	resp, err := http.Post(server.URL, "application/x-www-form-urlencoded", reader)

	assert.Nil(t, err, "There was an error")
	assert.Equal(t, http.StatusOK, resp.StatusCode, "Wrong http status code")

	os.RemoveAll("../files/testticket")
}

// unassign ticket
// --------------------------

type unassignTicketHandler struct{}

func (h *unassignTicketHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handleUnassignTicket(w, r)
}

func TestHandleHandleUnassignTicket(t *testing.T) {

	config := mockConfig()
	config.Tickets = "../files/testtickets/"

	globals.ServerConfig = &config

	handler := &unassignTicketHandler{}
	server := httptest.NewServer(handler)
	defer server.Close()

	ticket := structs.Ticket{
		Id: "abc123",
	}

	globals.Tickets["abc123"] = ticket

	resp, err := http.Get(server.URL + "/unassignTicket?id=abc123")

	assert.Nil(t, err, "An unexpected error occured")
	assert.Equal(t, http.StatusOK, resp.StatusCode, "The http status is not 200")

	os.RemoveAll("../files/testtickets/")
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

func TestHandleHandleAssignTicket(t *testing.T) {

	config := mockConfig()
	config.Tickets = "../files/testtickets/"
	globals.ServerConfig = &config

	users["testuser"] = structs.User{
		Id:          "1",
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
			Id:         "def123",
		},
		TTL: 3600,
	}

	handler := &assignTicketHandler{}
	server := httptest.NewServer(handler)
	defer server.Close()

	resp, err := http.Get(server.URL + "/assignTicket?id=abc123&user=testuser")

	assert.Nil(t, err, "An unexpected error occured")
	assert.Equal(t, http.StatusOK, resp.StatusCode, "The http status is not 200")

	os.RemoveAll("../files/testtickets/")
}
