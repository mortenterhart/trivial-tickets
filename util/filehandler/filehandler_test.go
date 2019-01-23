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

// Package filehandler takes care of interactions with files, writing
// and reading files and persisting changes to the file system.
package filehandler

import (
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"

	"github.com/mortenterhart/trivial-tickets/globals"
	"github.com/mortenterhart/trivial-tickets/log/testlog"
	"github.com/mortenterhart/trivial-tickets/structs"
	"github.com/mortenterhart/trivial-tickets/structs/defaults"
	"github.com/mortenterhart/trivial-tickets/util/hashing"
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
 * Package filehandler [tests]
 * Interactions with files, writing and reading files and persisting
 * changes to the file system
 */

// mockTicket is a helper function to create a dummy
// ticket for the tests.
func mockTicket() structs.Ticket {
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
		ID:          "12",
		Name:        "Max Mustermann",
		Mail:        "max.mustermann@example.com",
		Hash:        "$2a$12$n5kluCvuG3wpj18rl46bBexvTX6l0QkD7EQCkgvk1BNby5cNZPLZa",
		IsOnHoliday: false,
	}

	return structs.Ticket{
		ID:       "test123",
		Subject:  "Help",
		Status:   structs.StatusOpen,
		User:     user,
		Customer: "customer@example.com",
		Entries:  entries,
	}
}

// mockMail creates a new mock mail to be used
// inside tests.
func mockMail() structs.Mail {
	return structs.Mail{
		ID:      random.CreateRandomID(structs.RandomIDLength),
		From:    "no-reply@trivial-tickets.com",
		To:      "customer@mail.com",
		Subject: "[trivial-tickets] My screen is always black",
		Message: "I cannot see anything on my screen.",
	}
}

//revive:disable:deep-exit

// TestMain is started to run the tests and initializes the
// configuration before running the tests. The tests' exit
// status is returned as the overall exit status.
func TestMain(m *testing.M) {
	initializeLogConfig()

	os.Exit(m.Run())
}

//revive:enable:deep-exit

// initializeLogConfig initializes the global logging
// configuration with test values.
func initializeLogConfig() {
	logConfig := testLogConfig()
	globals.LogConfig = &logConfig
}

// testLogConfig returns a logging configuration suitable
// to be used in tests.
func testLogConfig() structs.LogConfig {
	return structs.LogConfig{
		LogLevel:  structs.AsLogLevel(defaults.LogLevelString),
		Verbose:   defaults.LogVerbose,
		FullPaths: defaults.LogFullPaths,
	}
}

// TestWriteReadUserFile tests both WriteUserFile and ReadUserFil
// back to back since the mock data can be used for both.
func TestWriteReadUserFile(t *testing.T) {
	testlog.BeginTest()
	defer testlog.EndTest()

	// File name for test
	const userFile string = "testUsers.json"

	// Create hashmap for users
	var users = make(map[string]structs.User)

	// Hash their passwords
	a, _ := hashing.GenerateHash("thisisatestPw12!!")

	// Mock two users and add them to the map
	u := structs.User{
		ID:          "abc123",
		Name:        "Admin",
		Username:    "admin",
		Mail:        "admin@example.com",
		Hash:        a,
		IsOnHoliday: false,
	}

	u1 := structs.User{
		ID:          "def456",
		Name:        "Max Mustermann",
		Username:    "max4711",
		Mail:        "max.mustermann@example.com",
		Hash:        a,
		IsOnHoliday: true,
	}

	users[u.Username] = u
	users[u1.Username] = u1

	// Write json to file
	errWriteFile := WriteUserFile(userFile, &users)
	assert.Nil(t, errWriteFile, "Error writing file")

	// Create hashmap to store the read json
	var readUsers = make(map[string]structs.User)

	// Read the file from disk and unmarshal into the hashmap
	errReadUserFile := ReadUserFile(userFile, &readUsers)
	assert.Nil(t, errReadUserFile, "There was an error reading the file")

	// Delete the test file
	errDeleteFile := os.Remove(userFile)
	assert.Nil(t, errDeleteFile, "Error deleting file")

	// Make sure the struct before writing to disk and after reading from disk is the same
	assert.Equal(t, users, readUsers, "User structs do not match")

	errReadUserFile2 := ReadUserFile("bla.json", &readUsers)
	assert.NotNil(t, errReadUserFile2, "No error was returned")
}

func TestReadUserFileInvalidJson(t *testing.T) {
	testlog.BeginTest()
	defer testlog.EndTest()

	const userDirectory string = "../../files/testusers"
	const userFile string = userDirectory + "/users.json"

	createErr := CreateFolders(userDirectory)
	t.Run("createFoldersError", func(t *testing.T) {
		assert.NoError(t, createErr, "creating testusers directory should not error")
	})

	writeErr := ioutil.WriteFile(userFile, []byte("{"), defaults.FileModeRegular)
	t.Run("writeUserFileError", func(t *testing.T) {
		assert.NoError(t, writeErr, "writing invalid user json file should not error")
	})

	users := make(map[string]structs.User)
	readErr := ReadUserFile(userFile, &users)

	t.Run("readError", func(t *testing.T) {
		assert.Error(t, readErr, "reading users.json with invalid json should be an error")
	})

	t.Run("emptyUsersMap", func(t *testing.T) {
		assert.Equal(t, 0, len(users), "users map should be empty due to read error")
	})

	removeErr := os.RemoveAll(userDirectory)
	t.Run("removeError", func(t *testing.T) {
		assert.NoError(t, removeErr, "removing users directory should not error")
	})
}

func TestWriteTicketFile(t *testing.T) {
	testlog.BeginTest()
	defer testlog.EndTest()

	// Path to ticket files
	const ticketFile string = "testFiles/testTickets"

	ticket := mockTicket()

	errWriteTicketFile := WriteTicketFile(ticketFile, &ticket)

	errRemove := os.RemoveAll("testFiles/")

	assert.Nil(t, errRemove, "Unexpected error while removing ticket directory")

	assert.Nil(t, errWriteTicketFile, "Error creating the File")
}

// TestWriteTicketFileError produces an error on creating
// a directory.
func TestWriteTicketFileError(t *testing.T) {
	testlog.BeginTest()
	defer testlog.EndTest()

	// Invalid Path
	const ticketFile string = ""

	ticket := mockTicket()

	errWriteTicketFile := WriteTicketFile(ticketFile, &ticket)

	assert.NotNil(t, errWriteTicketFile, "Error creating the File")
}

// TestReadTicketFiles checks if the ticket files are read
// correctly and if errors are returned when expected.
func TestReadTicketFiles(t *testing.T) {
	testlog.BeginTest()
	defer testlog.EndTest()

	var tickets = make(map[string]structs.Ticket)

	// Path does not exist
	errReadTicketFiles := ReadTicketFiles("abc", &tickets)
	assert.NotNil(t, errReadTicketFiles, "No error was returned, although the path does not exist")

	// Create folder for temporary test tickets
	const testTicketPath string = defaults.TestTickets
	errCreate := CreateFolders(testTicketPath)
	assert.NoError(t, errCreate, "Unexpected error while creating folders")

	// Write invalid JSON into file with .json extension
	const invalidJSONFile string = testTicketPath + "/invalid.json"
	errWrite := ioutil.WriteFile(invalidJSONFile, []byte("{"), defaults.FileModeRegular)
	assert.NoError(t, errWrite, "Unexpected error while writing ticket file")

	// Read invalid json file from test directory
	errReadTicketFiles2 := ReadTicketFiles(testTicketPath, &tickets)
	assert.Error(t, errReadTicketFiles2, "No error was returned, although the ticket file contains invalid json")

	// Remove invalid json file for next tests
	errRemove := os.Remove(invalidJSONFile)
	assert.NoError(t, errRemove, "Unexpected error while removing ticket file")

	ticket := mockTicket()
	errWrite = WriteTicketFile(testTicketPath, &ticket)
	assert.NoError(t, errWrite, "Unexpected error while writing ticket file")

	// Correct path to ticket files
	errReadTicketFiles3 := ReadTicketFiles(testTicketPath, &tickets)
	assert.Nil(t, errReadTicketFiles3, "An error was returned, although the path is correct")

	errRemove = os.RemoveAll(testTicketPath + "/")
	assert.NoError(t, errRemove, "Unexpected error while removing test directory")
}

func TestReadTicketFilesReadErrorPermissionDenied(t *testing.T) {
	testlog.BeginTest()
	defer testlog.EndTest()

	const testTicketDirectory string = defaults.TestTickets

	createErr := CreateFolders(testTicketDirectory)
	assert.NoError(t, createErr, "Creating test ticket directory should not return an error")

	const readPermissionDenied os.FileMode = 0200

	writeErr := ioutil.WriteFile(testTicketDirectory+"/noReadPermission.json", []byte("{}"), readPermissionDenied)
	assert.NoError(t, writeErr, "Writing the file with no read permission should not be an error")

	tickets := make(map[string]structs.Ticket)

	readErr := ReadTicketFiles(testTicketDirectory, &tickets)
	testlog.Debug(readErr)

	removeErr := os.RemoveAll(testTicketDirectory)
	assert.NoError(t, removeErr, "Removing test ticket directory should not be an error")
}

func TestReadMailFilesInvalidDirectory(t *testing.T) {
	testlog.BeginTest()
	defer testlog.EndTest()

	// Try to read mails from not existing directory
	const notExistingMailDirectory string = defaults.TestMails

	mails := make(map[string]structs.Mail)

	readErr := ReadMailFiles(notExistingMailDirectory, &mails)

	t.Run("readError", func(t *testing.T) {
		assert.Error(t, readErr, "reading from not existent directory should be an error")
	})

	t.Run("emptyMailMap", func(t *testing.T) {
		assert.Equal(t, 0, len(mails), "mail map should not contain any entries")
	})
}

func TestReadMailFilesInvalidJson(t *testing.T) {
	testlog.BeginTest()
	defer testlog.EndTest()

	const mailDirectory string = defaults.TestMails

	createErr := CreateFolders(mailDirectory)
	t.Run("noCreateError", func(t *testing.T) {
		assert.NoError(t, createErr, "creating mail directory should not return an error")
	})

	const invalidJSONFile string = mailDirectory + "/invalid.json"

	writeErr := ioutil.WriteFile(invalidJSONFile, []byte("{"), defaults.FileModeRegular)
	t.Run("noWriteError", func(t *testing.T) {
		assert.NoError(t, writeErr, "writing json file should not be an error")
	})

	mails := make(map[string]structs.Mail)

	readErr := ReadMailFiles(mailDirectory, &mails)
	t.Run("readError", func(t *testing.T) {
		assert.Error(t, readErr, "reading invalid.json with invalid json content should be a decoding error")
	})

	removeErr := os.RemoveAll(mailDirectory)
	t.Run("noRemoveError", func(t *testing.T) {
		assert.NoError(t, removeErr, "removing existing mail directory should not be an error")
	})
}

func TestReadMailFilesReadErrorPermissionDenied(t *testing.T) {
	testlog.BeginTest()
	defer testlog.EndTest()

	const testMailDirectory string = defaults.TestMails

	createErr := CreateFolders(testMailDirectory)
	assert.NoError(t, createErr, "Creating test mail directory should not return an error")

	const readPermissionDenied os.FileMode = 0200

	writeErr := ioutil.WriteFile(testMailDirectory+"/noReadPermission.json", []byte("{}"), readPermissionDenied)
	assert.NoError(t, writeErr, "Writing the file with no read permission should not be an error")

	mails := make(map[string]structs.Mail)

	readErr := ReadMailFiles(testMailDirectory, &mails)
	testlog.Debug(readErr)

	removeErr := os.RemoveAll(testMailDirectory)
	assert.NoError(t, removeErr, "Removing test mail directory should not be an error")
}

func TestWriteReadMailFile(t *testing.T) {
	testlog.BeginTest()
	defer testlog.EndTest()

	const mailDirectory string = defaults.TestMails

	testMail := mockMail()

	t.Run("writeMailFile", func(t *testing.T) {
		writeErr := WriteMailFile(mailDirectory, &testMail)

		t.Run("writeError", func(t *testing.T) {
			assert.NoError(t, writeErr, "writing mail file should not return error")
		})

		t.Run("mailDirectoryExists", func(t *testing.T) {
			assert.True(t, DirectoryExists(mailDirectory), "mailDirectory should exist because function creates missing folders")
		})

		t.Run("mailWritten", func(t *testing.T) {
			dirContents, readErr := ioutil.ReadDir(mailDirectory)

			assert.NoError(t, readErr, "reading contents of mail directory should not return an error")
			assert.Equal(t, 1, len(dirContents), "mail directory should contain exactly one mail")
		})
	})

	t.Run("readMailFiles", func(t *testing.T) {
		mails := make(map[string]structs.Mail)

		readErr := ReadMailFiles(mailDirectory, &mails)

		t.Run("readError", func(t *testing.T) {
			assert.NoError(t, readErr, "reading mail files should not return error since a mail exists")
		})

		t.Run("numberOfReadMails", func(t *testing.T) {
			assert.Equal(t, 1, len(mails), "there should be one mail read because one mail was previously written")
		})

		t.Run("identicalMailId", func(t *testing.T) {
			readMail, mailIDDefined := mails[testMail.ID]
			assert.True(t, mailIDDefined, "test mail id should be defined in read mail map")
			assert.NotNil(t, readMail, "read mail should be non-nil")
		})
	})

	errRemove := os.RemoveAll(mailDirectory)
	assert.NoError(t, errRemove, "error while removing mail directory")
}

func TestWriteMailFileInvalidDirectory(t *testing.T) {
	testlog.BeginTest()
	defer testlog.EndTest()

	const mailDirectory string = ""

	testMail := mockMail()

	writeErr := WriteMailFile(mailDirectory, &testMail)

	assert.Error(t, writeErr, "writing a mail in an empty directory name should return an error")
}

func TestRemoveMailFile(t *testing.T) {
	testlog.BeginTest()
	defer testlog.EndTest()

	const mailDirectory string = defaults.TestMails

	t.Run("existingMail", func(t *testing.T) {
		// Create test ticket to be removed
		testMail := mockMail()
		errWrite := WriteMailFile(mailDirectory, &testMail)

		t.Run("noWriteError", func(t *testing.T) {
			assert.NoError(t, errWrite, "writing mail file should not return an error")
		})

		removeErr := RemoveMailFile(mailDirectory, testMail.ID)

		t.Run("noRemoveError", func(t *testing.T) {
			assert.NoError(t, removeErr, "removing mail file should not return error since the file exists")
		})
	})

	t.Run("notExistingMail", func(t *testing.T) {
		notExistingMailID := "mail-id"

		removeErr := RemoveMailFile(mailDirectory, notExistingMailID)

		assert.Error(t, removeErr, "remove error should be non-nil because mail file does not exist")
	})

	removeErr := os.RemoveAll(mailDirectory)
	assert.NoError(t, removeErr, "removing mail directory should not return an error because the directory exists")
}

func TestFileExists(t *testing.T) {
	testlog.BeginTest()
	defer testlog.EndTest()

	t.Run("existingFile", func(t *testing.T) {
		const existingFile string = "../../files/users/users.json"

		assert.True(t, FileExists(existingFile), "users.json file should always exist")
	})

	t.Run("notExistingFile", func(t *testing.T) {
		const notExistingFile string = "../../files/users/passwords.json"

		assert.False(t, FileExists(notExistingFile), "nobody would store passwords in a JSON-file, so why should it exist?")
	})
}

func TestDirectoryExists(t *testing.T) {
	testlog.BeginTest()
	defer testlog.EndTest()

	t.Run("existingDirectory", func(t *testing.T) {
		const existingDirectory string = "../../files/tickets"

		assert.True(t, DirectoryExists(existingDirectory), "tickets directory should exist")
	})

	t.Run("notExistingDirectory", func(t *testing.T) {
		const notExistingDirectory string = "../../files/secret_keys"

		assert.False(t, DirectoryExists(notExistingDirectory), "again, secret keys should not be stored here")
	})
}

func TestHasJsonExtension(t *testing.T) {
	testlog.BeginTest()
	defer testlog.EndTest()

	t.Run("jsonExtension", func(t *testing.T) {
		assert.True(t, hasJSONExtension("ticket.json"), "should be true because file has a .json file extension")
	})

	t.Run("noJsonExtension", func(t *testing.T) {
		assert.False(t, hasJSONExtension("ticket.xml"), "should be false because file has .xml file extension instead of .json")
	})
}

func TestCreateFolders(t *testing.T) {
	testlog.BeginTest()
	defer testlog.EndTest()

	// Test folders
	const ticketsFolder string = "testFolder/tests"

	// Create the given folders
	errCreateFolder := CreateFolders(ticketsFolder)

	// Check that there was no error
	assert.Nil(t, errCreateFolder, "Error creating the folder(s)")

	// Check that the created directory exists
	assert.DirExists(t, ticketsFolder, "ticketsFolder/tests should exist now")

	// Remove them
	errRemove := os.RemoveAll("testFolder/")

	assert.Nil(t, errRemove, "Unexpected error while removing test directory")
}

func TestWrapAndLogError(t *testing.T) {
	testlog.BeginTest()
	defer testlog.EndTest()

	readErr := errors.New("open ../../files/testmails/mail.json: no such file or directory")
	expectedErr := errors.Wrap(readErr, "could not read mail file")

	wrapErr := wrapAndLogError(readErr, "could not read mail file")

	t.Run("notNil", func(t *testing.T) {
		assert.NotNil(t, wrapErr, "wrap error should not be nil")
	})

	t.Run("equalWrappedError", func(t *testing.T) {
		assert.Equal(t, expectedErr.Error(), wrapErr.Error(), "expected and wrapped error should be identical")
	})
}

func TestWrapAndLogErrorf(t *testing.T) {
	testlog.BeginTest()
	defer testlog.EndTest()

	err := errors.New("mkdir: '': no such file or directory")
	expectedErr := errors.Wrap(err, "could not create mail directory '../../files/testmails'")

	wrapErr := wrapAndLogErrorf(err, "could not create mail directory '%s'", "../../files/testmails")

	t.Run("notNil", func(t *testing.T) {
		assert.NotNil(t, wrapErr, "wrap error should not be nil")
	})

	t.Run("equalWrappedError", func(t *testing.T) {
		assert.Equal(t, expectedErr.Error(), wrapErr.Error(), "expected and wrapped error should be identical")
	})
}
