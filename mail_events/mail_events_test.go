// Mail message construction using templating
package mail_events

import (
	"fmt"
	"github.com/mortenterhart/trivial-tickets/globals"
	"github.com/mortenterhart/trivial-tickets/util/random"
	"github.com/stretchr/testify/assert"
	"html/template"
	"os"
	"strings"
	"testing"

	"github.com/mortenterhart/trivial-tickets/structs"
	"github.com/mortenterhart/trivial-tickets/ticket"
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
		Id:       random.CreateRandomId(10),
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
		Id:          "user-id",
		Name:        "Admin",
		Username:    "admin",
		Mail:        "admin@example.com",
		Hash:        "anyHashThatIsImaginable",
		IsOnHoliday: false,
	}
}

func testConfig() structs.Config {
	return structs.Config{
		Port:    8443,
		Tickets: "../../files/testtickets",
		Users:   "../../files/users/users.json",
		Mails:   "../../files/testmails",
		Cert:    "../../ssl/server.cert",
		Key:     "../../ssl/server.key",
		Web:     "../../www",
	}
}

func testLogConfig() structs.LogConfig {
	return structs.LogConfig{
		LogLevel:   structs.LevelInfo,
		VerboseLog: false,
		FullPaths:  false,
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
	t.Run("withEntrySet", func(t *testing.T) {
		testTicket := mockTicketWithEntry()

		mailBody := NewMailBody(NewTicket, testTicket)

		t.Run("containsMailEvent", func(t *testing.T) {
			assert.True(t, strings.Contains(mailBody, fmt.Sprintf("Ihr Ticket '%s' ist erfolgreich erstellt worden",
				template.HTMLEscapeString(testTicket.Id))), "mail body should contain a description of the happened event")
		})

		t.Run("containsEntry", func(t *testing.T) {
			assert.True(t, strings.Contains(mailBody, template.HTMLEscapeString(testTicket.Entries[0].Text)),
				"mail body should contain first written message")
		})
	})

	t.Run("withUserSet", func(t *testing.T) {
		testTicket := mockTicketWithUser()

		mailBody := NewMailBody(NewTicket, testTicket)

		t.Run("containsMailEvent", func(t *testing.T) {
			assert.True(t, strings.Contains(mailBody, fmt.Sprintf("Ihr Ticket '%s' ist erfolgreich erstellt worden",
				template.HTMLEscapeString(testTicket.Id))), "mail body should contain a description of the happened event")
		})

		t.Run("containsAssignedUser", func(t *testing.T) {
			assert.True(t, strings.Contains(mailBody, fmt.Sprintf("%s (%s)", template.HTMLEscapeString(testTicket.User.Name),
				template.HTMLEscapeString(testTicket.User.Mail))), "mail body should contain the name and email of the assigned user")
		})
	})
}

func TestNewMailBodyNewAnswer(t *testing.T) {
	t.Run("withEntrySet", func(t *testing.T) {
		testTicket := mockTicketWithEntry()

		mailBody := NewMailBody(NewAnswer, testTicket)

		t.Run("containsMailEvent", func(t *testing.T) {
			assert.True(t, strings.Contains(mailBody, fmt.Sprintf("der Benutzer '%s' hat einen neuen Kommentar geschrieben",
				template.HTMLEscapeString(testTicket.Entries[0].User))), "mail body should contain a description of the happened event")
		})

		t.Run("containsEntry", func(t *testing.T) {
			assert.True(t, strings.Contains(mailBody, template.HTMLEscapeString(testTicket.Entries[0].Text)),
				"mail body should contain first written message")
		})
	})

	t.Run("withUserSet", func(t *testing.T) {
		testTicket := mockTicketWithUser()

		mailBody := NewMailBody(NewAnswer, testTicket)

		t.Run("containsMailEvent", func(t *testing.T) {
			assert.True(t, strings.Contains(mailBody, "der Benutzer '' hat einen neuen Kommentar geschrieben"),
				"mail body should contain a description of the happened event")
		})

		t.Run("containsAssignedUser", func(t *testing.T) {
			assert.True(t, strings.Contains(mailBody, fmt.Sprintf("%s (%s)", template.HTMLEscapeString(testTicket.User.Name),
				template.HTMLEscapeString(testTicket.User.Mail))), "mail body should contain the name and email of the assigned user")
		})
	})
}

func TestNewMailBodyUpdatedTicket(t *testing.T) {
	t.Run("withEntrySet", func(t *testing.T) {
		testTicket := mockTicketWithEntry()

		mailBody := NewMailBody(UpdatedTicket, testTicket)

		t.Run("containsMailEvent", func(t *testing.T) {
			assert.True(t, strings.Contains(mailBody, fmt.Sprintf("Ihr Ticket '%s' wurde mit folgenden Informationen aktualisiert",
				template.HTMLEscapeString(testTicket.Id))), "mail body should contain a description of the happened event")
		})

		t.Run("containsEntry", func(t *testing.T) {
			assert.True(t, strings.Contains(mailBody, template.HTMLEscapeString(testTicket.Entries[0].Text)),
				"mail body should contain first written message")
		})
	})

	t.Run("withUserSet", func(t *testing.T) {
		testTicket := mockTicketWithUser()

		mailBody := NewMailBody(UpdatedTicket, testTicket)

		t.Run("containsMailEvent", func(t *testing.T) {
			assert.True(t, strings.Contains(mailBody, fmt.Sprintf("Ihr Ticket '%s' wurde mit folgenden Informationen aktualisiert",
				template.HTMLEscapeString(testTicket.Id))), "mail body should contain a description of the happened event")
		})

		t.Run("containsAssignedUser", func(t *testing.T) {
			assert.True(t, strings.Contains(mailBody, fmt.Sprintf("%s (%s)", template.HTMLEscapeString(testTicket.User.Name),
				template.HTMLEscapeString(testTicket.User.Mail))), "mail body should contain the name and email of the assigned user")
		})
	})
}

func TestNewMailBodyAssignedTicket(t *testing.T) {
	t.Run("withEntrySet", func(t *testing.T) {
		testTicket := mockTicketWithEntry()

		mailBody := NewMailBody(AssignedTicket, testTicket)

		t.Run("containsMailEvent", func(t *testing.T) {
			assert.True(t, strings.Contains(mailBody, fmt.Sprintf("unser Mitarbeiter '%s' bearbeitet nun Ihr Ticket",
				template.HTMLEscapeString("<nicht zugewiesen>"))), "mail body should contain a description of the happened event")
		})

		t.Run("containsEntry", func(t *testing.T) {
			assert.True(t, strings.Contains(mailBody, template.HTMLEscapeString(testTicket.Entries[0].Text)),
				"mail body should contain first written message")
		})
	})

	t.Run("withUserSet", func(t *testing.T) {
		testTicket := mockTicketWithUser()

		mailBody := NewMailBody(AssignedTicket, testTicket)

		t.Run("containsMailEvent", func(t *testing.T) {
			assert.True(t, strings.Contains(mailBody, fmt.Sprintf("unser Mitarbeiter '%s' bearbeitet nun Ihr Ticket",
				template.HTMLEscapeString(testTicket.User.Name))), "mail body should contain a description of the happened event")
		})

		t.Run("containsAssignedUser", func(t *testing.T) {
			assert.True(t, strings.Contains(mailBody, fmt.Sprintf("%s (%s)", template.HTMLEscapeString(testTicket.User.Name),
				template.HTMLEscapeString(testTicket.User.Mail))), "mail body should contain the name and email of the assigned user")
		})
	})
}

func TestNewMailBodyUnassignedTicket(t *testing.T) {
	t.Run("withEntrySet", func(t *testing.T) {
		testTicket := mockTicketWithEntry()

		mailBody := NewMailBody(UnassignedTicket, testTicket)

		t.Run("containsMailEvent", func(t *testing.T) {
			assert.True(t, strings.Contains(mailBody, fmt.Sprintf("der Bearbeiter '%s' hat das Ticket wieder freigegeben",
				template.HTMLEscapeString("<nicht zugewiesen>"))), "mail body should contain a description of the happened event")
		})

		t.Run("containsEntry", func(t *testing.T) {
			assert.True(t, strings.Contains(mailBody, template.HTMLEscapeString(testTicket.Entries[0].Text)),
				"mail body should contain first written message")
		})
	})

	t.Run("withUserSet", func(t *testing.T) {
		testTicket := mockTicketWithUser()

		mailBody := NewMailBody(UnassignedTicket, testTicket)

		t.Run("containsMailEvent", func(t *testing.T) {
			assert.True(t, strings.Contains(mailBody, fmt.Sprintf("der Bearbeiter '%s' hat das Ticket wieder freigegeben",
				template.HTMLEscapeString(testTicket.User.Name))), "mail body should contain a description of the happened event")
		})

		t.Run("containsAssignedUser", func(t *testing.T) {
			assert.True(t, strings.Contains(mailBody, fmt.Sprintf("%s (%s)", template.HTMLEscapeString(testTicket.User.Name),
				template.HTMLEscapeString(testTicket.User.Mail))), "mail body should contain the name and email of the assigned user")
		})
	})
}
