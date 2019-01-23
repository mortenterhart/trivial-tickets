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

// Package mail_events provides facilities to create
// standard mail messages for different actions using
// predefined templates.
package mail_events

import (
	"fmt"
	"html/template"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mortenterhart/trivial-tickets/globals"
	"github.com/mortenterhart/trivial-tickets/logger/testlogger"
	"github.com/mortenterhart/trivial-tickets/structs"
	"github.com/mortenterhart/trivial-tickets/structs/defaults"
	"github.com/mortenterhart/trivial-tickets/ticket"
	"github.com/mortenterhart/trivial-tickets/util/random"
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
 * Package mail_events [tests]
 * Mail message construction using templating
 */

func mockTicketWithEntry() structs.Ticket {
	return ticket.CreateTicket("customer@mail.com", "Something is wrong", "I don't know what")
}

func mockTicketWithUser() structs.Ticket {
	return structs.Ticket{
		ID:       random.CreateRandomID(10),
		Subject:  "Something was wrong",
		Status:   structs.CLOSED,
		User:     mockUser(),
		Customer: "customer@mail.com",
		Entries:  nil,
		MergeTo:  "",
	}
}

func mockUser() structs.User {
	return structs.User{
		ID:          "user-id",
		Name:        "Admin",
		Username:    "admin",
		Mail:        "admin@example.com",
		Hash:        "anyHashThatIsImaginable",
		IsOnHoliday: false,
	}
}

func testConfig() structs.ServerConfig {
	return structs.ServerConfig{
		Port:    defaults.TestPort,
		Tickets: defaults.TestTickets,
		Users:   defaults.TestUsers,
		Mails:   defaults.TestMails,
		Cert:    defaults.TestCertificate,
		Key:     defaults.TestKey,
		Web:     defaults.TestWeb,
	}
}

func testLogConfig() structs.LogConfig {
	return structs.LogConfig{
		LogLevel:  structs.AsLogLevel(defaults.LogLevelString),
		Verbose:   defaults.LogVerbose,
		FullPaths: defaults.LogFullPaths,
	}
}

func initializeConfig() {
	config := testConfig()
	globals.ServerConfig = &config

	logConfig := testLogConfig()
	globals.LogConfig = &logConfig
}

// Setup and teardown
func TestMain(m *testing.M) {
	initializeConfig()

	os.Exit(m.Run())
}

func TestNewMailBodyNewTicket(t *testing.T) {
	testlogger.BeginTest()
	defer testlogger.EndTest()

	t.Run("withEntrySet", func(t *testing.T) {
		testTicket := mockTicketWithEntry()

		mailBody := NewMailBody(NewTicket, testTicket)

		t.Run("containsMailEvent", func(t *testing.T) {
			assert.Contains(t, mailBody, fmt.Sprintf("Your Ticket '%s' was created successfully",
				template.HTMLEscapeString(testTicket.ID)), "mail body should contain a description of the happened event")
		})

		t.Run("containsEntry", func(t *testing.T) {
			assert.Contains(t, mailBody, template.HTMLEscapeString(testTicket.Entries[0].Text),
				"mail body should contain first written message")
		})
	})

	t.Run("withUserSet", func(t *testing.T) {
		testTicket := mockTicketWithUser()

		mailBody := NewMailBody(NewTicket, testTicket)

		t.Run("containsMailEvent", func(t *testing.T) {
			assert.Contains(t, mailBody, fmt.Sprintf("Your Ticket '%s' was created successfully",
				template.HTMLEscapeString(testTicket.ID)), "mail body should contain a description of the happened event")
		})

		t.Run("containsAssignedUser", func(t *testing.T) {
			assert.Contains(t, mailBody, fmt.Sprintf("%s (%s)", template.HTMLEscapeString(testTicket.User.Name),
				template.HTMLEscapeString(testTicket.User.Mail)), "mail body should contain the name and email of the assigned user")
		})
	})
}

func TestNewMailBodyNewAnswer(t *testing.T) {
	testlogger.BeginTest()
	defer testlogger.EndTest()

	t.Run("withEntrySet", func(t *testing.T) {
		testTicket := mockTicketWithEntry()

		mailBody := NewMailBody(NewAnswer, testTicket)

		t.Run("containsMailEvent", func(t *testing.T) {
			assert.Contains(t, mailBody, fmt.Sprintf("the user '%s' wrote a new comment to your ticket",
				template.HTMLEscapeString(testTicket.Entries[0].User)), "mail body should contain a description of the happened event")
		})

		t.Run("containsEntry", func(t *testing.T) {
			assert.Contains(t, mailBody, template.HTMLEscapeString(testTicket.Entries[0].Text),
				"mail body should contain first written message")
		})
	})

	t.Run("withUserSet", func(t *testing.T) {
		testTicket := mockTicketWithUser()

		mailBody := NewMailBody(NewAnswer, testTicket)

		t.Run("containsMailEvent", func(t *testing.T) {
			assert.Contains(t, mailBody, fmt.Sprintf("the user '%s' wrote a new comment to your ticket",
				template.HTMLEscapeString("<no user>")), "mail body should contain a description of the happened event")
		})

		t.Run("containsAssignedUser", func(t *testing.T) {
			assert.Contains(t, mailBody, fmt.Sprintf("%s (%s)", template.HTMLEscapeString(testTicket.User.Name),
				template.HTMLEscapeString(testTicket.User.Mail)), "mail body should contain the name and email of the assigned user")
		})
	})
}

func TestNewMailBodyUpdatedTicket(t *testing.T) {
	testlogger.BeginTest()
	defer testlogger.EndTest()

	t.Run("withEntrySet", func(t *testing.T) {
		testTicket := mockTicketWithEntry()

		mailBody := NewMailBody(UpdatedTicket, testTicket)

		t.Run("containsMailEvent", func(t *testing.T) {
			assert.Contains(t, mailBody, fmt.Sprintf("Your Ticket '%s' was updated with the following information",
				template.HTMLEscapeString(testTicket.ID)), "mail body should contain a description of the happened event")
		})

		t.Run("containsEntry", func(t *testing.T) {
			assert.Contains(t, mailBody, template.HTMLEscapeString(testTicket.Entries[0].Text),
				"mail body should contain first written message")
		})
	})

	t.Run("withUserSet", func(t *testing.T) {
		testTicket := mockTicketWithUser()

		mailBody := NewMailBody(UpdatedTicket, testTicket)

		t.Run("containsMailEvent", func(t *testing.T) {
			assert.Contains(t, mailBody, fmt.Sprintf("Your Ticket '%s' was updated with the following information",
				template.HTMLEscapeString(testTicket.ID)), "mail body should contain a description of the happened event")
		})

		t.Run("containsAssignedUser", func(t *testing.T) {
			assert.Contains(t, mailBody, fmt.Sprintf("%s (%s)", template.HTMLEscapeString(testTicket.User.Name),
				template.HTMLEscapeString(testTicket.User.Mail)), "mail body should contain the name and email of the assigned user")
		})
	})
}

func TestNewMailBodyAssignedTicket(t *testing.T) {
	testlogger.BeginTest()
	defer testlogger.EndTest()

	t.Run("withEntrySet", func(t *testing.T) {
		testTicket := mockTicketWithEntry()

		mailBody := NewMailBody(AssignedTicket, testTicket)

		t.Run("containsMailEvent", func(t *testing.T) {
			assert.Contains(t, mailBody, fmt.Sprintf("the editor '%s' works on Your Ticket now",
				template.HTMLEscapeString("<not assigned>")), "mail body should contain a description of the happened event")
		})

		t.Run("containsEntry", func(t *testing.T) {
			assert.Contains(t, mailBody, template.HTMLEscapeString(testTicket.Entries[0].Text),
				"mail body should contain first written message")
		})
	})

	t.Run("withUserSet", func(t *testing.T) {
		testTicket := mockTicketWithUser()

		mailBody := NewMailBody(AssignedTicket, testTicket)

		t.Run("containsMailEvent", func(t *testing.T) {
			assert.Contains(t, mailBody, fmt.Sprintf("the editor '%s' works on Your Ticket now",
				template.HTMLEscapeString(testTicket.User.Name)), "mail body should contain a description of the happened event")
		})

		t.Run("containsAssignedUser", func(t *testing.T) {
			assert.Contains(t, mailBody, fmt.Sprintf("%s (%s)", template.HTMLEscapeString(testTicket.User.Name),
				template.HTMLEscapeString(testTicket.User.Mail)), "mail body should contain the name and email of the assigned user")
		})
	})
}

func TestNewMailBodyUnassignedTicket(t *testing.T) {
	testlogger.BeginTest()
	defer testlogger.EndTest()

	t.Run("withEntrySet", func(t *testing.T) {
		testTicket := mockTicketWithEntry()

		mailBody := NewMailBody(UnassignedTicket, testTicket)

		t.Run("containsMailEvent", func(t *testing.T) {
			assert.Contains(t, mailBody, fmt.Sprintf("the editor '%s' has released Your Ticket again",
				template.HTMLEscapeString("<not assigned>")), "mail body should contain a description of the happened event")
		})

		t.Run("containsEntry", func(t *testing.T) {
			assert.Contains(t, mailBody, template.HTMLEscapeString(testTicket.Entries[0].Text),
				"mail body should contain first written message")
		})
	})

	t.Run("withUserSet", func(t *testing.T) {
		testTicket := mockTicketWithUser()

		mailBody := NewMailBody(UnassignedTicket, testTicket)

		t.Run("containsMailEvent", func(t *testing.T) {
			assert.Contains(t, mailBody, fmt.Sprintf("the editor '%s' has released Your Ticket again",
				template.HTMLEscapeString(testTicket.User.Name)), "mail body should contain a description of the happened event")
		})

		t.Run("containsAssignedUser", func(t *testing.T) {
			assert.Contains(t, mailBody, fmt.Sprintf("%s (%s)", template.HTMLEscapeString(testTicket.User.Name),
				template.HTMLEscapeString(testTicket.User.Mail)), "mail body should contain the name and email of the assigned user")
		})
	})
}

func TestEvent_String(t *testing.T) {
	testlogger.BeginTest()
	defer testlogger.EndTest()

	t.Run("newTicket", func(t *testing.T) {
		assert.Equal(t, "new ticket", NewTicket.String())
	})

	t.Run("newAnswer", func(t *testing.T) {
		assert.Equal(t, "new answer", NewAnswer.String())
	})

	t.Run("updatedTicket", func(t *testing.T) {
		assert.Equal(t, "updated ticket", UpdatedTicket.String())
	})

	t.Run("assignedTicket", func(t *testing.T) {
		assert.Equal(t, "assigned ticket", AssignedTicket.String())
	})

	t.Run("unassignedTicket", func(t *testing.T) {
		assert.Equal(t, "unassigned ticket", UnassignedTicket.String())
	})

	t.Run("undefinedEvent", func(t *testing.T) {
		assert.Equal(t, "undefined", Event(100).String())
	})
}
