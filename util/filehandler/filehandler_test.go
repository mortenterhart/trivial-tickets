package filehandler

import (
	"os"
	"testing"
	"time"

	"github.com/mortenterhart/trivial-tickets/structs"
	"github.com/mortenterhart/trivial-tickets/util/hashing"
	"github.com/stretchr/testify/assert"
)

/*
*
* Matrikelnummern
* 3040018
* 6694964
* 3478222
 */

// TestWriteReaduserFile tests both WriteUserFile and ReadUserFile back to back since the
// mock data can be used for both
func TestWriteReadUserFile(t *testing.T) {

	// File name for test
	const file = "testUsers.json"

	// Create hashmap for users
	var users = make(map[string]structs.User)

	// Hash their passwords
	a, _ := hashing.GenerateHash("thisisatestPw12!!")

	// Mock two users and add them to the map
	u := structs.User{
		Id:          "abc123",
		Name:        "Admin",
		Username:    "admin",
		Mail:        "admin@example.com",
		Hash:        a,
		IsOnHoliday: false,
	}

	u1 := structs.User{
		Id:          "def456",
		Name:        "Max Mustermann",
		Username:    "max4711",
		Mail:        "max.mustermann@example.com",
		Hash:        a,
		IsOnHoliday: true,
	}

	users[u.Username] = u
	users[u1.Username] = u1

	// Write json to file
	errWriteFile := WriteUserFile(file, &users)
	assert.Nil(t, errWriteFile, "Error writing file")

	// Create hashmap to store the read json
	var readUsers = make(map[string]structs.User)

	// Read the file from disk and unmarshal into the hashmap
	errReadUserFile := ReadUserFile(file, &readUsers)
	assert.Nil(t, errReadUserFile, "There was an error reading the file")

	// Delete the test file
	errDeleteFile := os.Remove(file)
	assert.Nil(t, errDeleteFile, "Error deleting file")

	// Make sure the struct before writing to disk and after reading from disk is the same
	assert.Equal(t, users, readUsers)

	errReadUserFile2 := ReadUserFile("bla.json", &readUsers)
	assert.NotNil(t, errReadUserFile2, "No error was returned")

}

func TestWriteTicketFile(t *testing.T) {

	e1 := structs.Entry{
		Date: time.Now(),
		User: "customer@example.com",
		Text: "bla bla",
	}

	e2 := structs.Entry{
		Date: time.Now(),
		User: "max.mustermann@example.com",
		Text: "ok ok",
	}

	entries := []structs.Entry{e1, e2}

	user := structs.User{
		Id:          "12",
		Name:        "Max Mustermann",
		Mail:        "max.mustermann@example.com",
		Hash:        "$2a$12$n5kluCvuG3wpj18rl46bBexvTX6l0QkD7EQCkgvk1BNby5cNZPLZa",
		IsOnHoliday: false,
	}

	ticket := structs.Ticket{
		Id:       "test123",
		Subject:  "Help",
		Status:   0,
		User:     user,
		Customer: "customer@example.com",
		Entries:  entries,
	}

	// Path to ticket files
	const usersFile = "testFiles/testTickets"

	errWriteTicketFile := WriteTicketFile(usersFile, &ticket)

	os.RemoveAll("testFiles/")

	assert.Nil(t, errWriteTicketFile, "Error creating the File")
}

func TestCreateFolder(t *testing.T) {

	// Test folders
	const ticketsFolder = "testFolder/tests"

	// Create the given folders
	errCreateFolder := CreateFolders(ticketsFolder)

	// Remove them
	os.RemoveAll("testFolder/")

	// Check that there was no error
	assert.Nil(t, errCreateFolder, "Error creating the folder(s)")
}

// TestReadTicketFiles checks if the ticket files a read correctly and if errors are returned when expected
func TestReadTicketFiles(t *testing.T) {

	var tickets = make(map[string]structs.Ticket)

	// Path does not exist
	errReadTicketFiles := ReadTicketFiles("abc", &tickets)
	assert.NotNil(t, errReadTicketFiles, "No error was returned, although the path does not exist")

	// Goes to html files, no json
	errReadTicketFiles2 := ReadTicketFiles("../../www/templates", &tickets)
	assert.NotNil(t, errReadTicketFiles2, "No error was returned, although the ticket files do not exist")

	// Correct path to ticket files
	errReadTicketFiles3 := ReadTicketFiles("../../files/tickets", &tickets)
	assert.Nil(t, errReadTicketFiles3, "An erorr was returned, although the path is correct")
}
