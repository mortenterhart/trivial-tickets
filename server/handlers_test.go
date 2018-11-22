package server

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/mortenterhart/trivial-tickets/globals"
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

// index
// --------------------------
type indexHandler struct{}

func (h *indexHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handleIndex(w, r)
}

func TestHandleIndex(t *testing.T) {

	handler := &indexHandler{}

	globals.Tickets["abc123"] = structs.Ticket{}
	users["abc123"] = structs.User{}

	server := httptest.NewServer(handler)
	defer server.Close()

	http.Get(server.URL)
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
	handleHoliday(w, r)
}

func TestHandleHandleHoliday(t *testing.T) {

	handler := &holidayHandler{}

	server := httptest.NewServer(handler)
	defer server.Close()

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

	handler := &ticketHandler{}

	server := httptest.NewServer(handler)
	defer server.Close()

	ticket := structs.Ticket{
		Id: "abc123",
	}

	globals.Tickets["abc123"] = ticket

	http.Get(server.URL + "/ticket?id=abc123")
}

// update ticket
// --------------------------

type updateTicketHandler struct{}

func (h *updateTicketHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handleUpdateTicket(w, r)
}

func TestHandleHandleUpdateTicket(t *testing.T) {

	handler := &updateTicketHandler{}

	server := httptest.NewServer(handler)
	defer server.Close()

	http.Post(server.URL, "application/x-www-form-urlencoded", nil)
}

// unassign ticket
// --------------------------

type unassignTicketHandler struct{}

func (h *unassignTicketHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handleUnassignTicket(w, r)
}

func TestHandleHandleUnassignTicket(t *testing.T) {

	handler := &unassignTicketHandler{}

	server := httptest.NewServer(handler)
	defer server.Close()

	ticket := structs.Ticket{
		Id: "abc123",
	}

	globals.Tickets["abc123"] = ticket

	http.Get(server.URL + "/unassignTicket?id=abc123")
}

// assign ticket
// --------------------------

type assignTicketHandler struct{}

func (h *assignTicketHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handleAssignTicket(w, r)
}

func TestHandleHandleAssignTicket(t *testing.T) {

	handler := &assignTicketHandler{}

	server := httptest.NewServer(handler)
	defer server.Close()

	http.Get(server.URL + "/assignTicket")
}
