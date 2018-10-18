package filehandler

import (
	"encoding/json"
	"go-tickets/structs"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

/*
*
* Matrikelnummern
* 3040018
* 6694964
* 3478222
 */

func TestReadUserFile(t *testing.T) {

	// File name for test
	const file = "testUsers.json"

	user := structs.User{
		Id:          12,
		Name:        "Max Mustermann",
		Mail:        "max.mustermann@example.com",
		Hash:        "$2a$12$n5kluCvuG3wpj18rl46bBexvTX6l0QkD7EQCkgvk1BNby5cNZPLZa",
		IsOnHoliday: false,
	}

	// Turn struct into json
	usersMarshal, _ := json.Marshal(user)

	// Write json to file
	errWriteFile := ioutil.WriteFile(file, usersMarshal, 0644)
	assert.Nil(t, errWriteFile, "Error writing file")

	// Read the file from disk and unmarshal into struct of type User
	readUsers := ReadUserFile(file)

	// Delete the test file
	errDeleteFile := os.Remove(file)
	assert.Nil(t, errDeleteFile, "Error deleting file")
	// Make sure the struct before writing to disk and after reading from disk is the same
	assert.Equal(t, user, readUsers)
}

func TestCreateFile(t *testing.T) {

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
		Id:          12,
		Name:        "Max Mustermann",
		Mail:        "max.mustermann@example.com",
		Hash:        "$2a$12$n5kluCvuG3wpj18rl46bBexvTX6l0QkD7EQCkgvk1BNby5cNZPLZa",
		IsOnHoliday: false,
	}

	ticket := structs.Ticket{
		Id:       34654522,
		Subject:  "Help",
		Status:   0,
		User:     user,
		Customer: "customer@example.com",
		Entries:  entries,
	}

	// Path to ticket files
	const usersFile = "testFiles/testTickets"

	errCreateFile := CreateFile(usersFile, &ticket)

	os.RemoveAll("testFiles/")

	assert.Nil(t, errCreateFile, "Error creating the File")
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
