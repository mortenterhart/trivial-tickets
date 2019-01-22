// Interactions with files, writing and reading files and persisting
// changes to the file system
package filehandler

import (
	"github.com/mortenterhart/trivial-tickets/globals"
	"github.com/pkg/errors"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/mortenterhart/trivial-tickets/structs"
	"github.com/mortenterhart/trivial-tickets/util/hashing"
	"github.com/mortenterhart/trivial-tickets/util/random"
	"github.com/stretchr/testify/assert"
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

func TestMain(m *testing.M) {
	initializeLogConfig()

	os.Exit(m.Run())
}

func initializeLogConfig() {
	logConfig := testLogConfig()
	globals.LogConfig = &logConfig
}

func testLogConfig() structs.LogConfig {
	return structs.LogConfig{
		LogLevel:   structs.LevelInfo,
		VerboseLog: false,
		FullPaths:  false,
	}
}

// TestWriteReadUserFile tests both WriteUserFile and ReadUserFile back to back since the
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
	assert.Equal(t, users, readUsers, "User structs do not match")

	errReadUserFile2 := ReadUserFile("bla.json", &readUsers)
	assert.NotNil(t, errReadUserFile2, "No error was returned")
}

func TestWriteTicketFile(t *testing.T) {

	// Path to ticket files
	const usersFile = "testFiles/testTickets"

	ticket := mockTicket()

	errWriteTicketFile := WriteTicketFile(usersFile, &ticket)

	os.RemoveAll("testFiles/")

	assert.Nil(t, errWriteTicketFile, "Error creating the File")
}

// TestWriteTicketFileError produces an error on creating a directory
func TestWriteTicketFileError(t *testing.T) {

	// Invalid Path
	const usersFile = ""

	ticket := mockTicket()

	errWriteTicketFile := WriteTicketFile(usersFile, &ticket)

	assert.NotNil(t, errWriteTicketFile, "Error creating the File")
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

	// Create folder for temporary test tickets
	const testTicketPath = "../../files/testtickets"
	CreateFolders(testTicketPath)

	// Write invalid JSON into file with .json extension
	const invalidJsonFile = testTicketPath + "/invalid.json"
	ioutil.WriteFile(invalidJsonFile, []byte("{"), 0644)

	// Read invalid json file from test directory
	errReadTicketFiles2 := ReadTicketFiles(testTicketPath, &tickets)
	assert.NotNil(t, errReadTicketFiles2, "No error was returned, although the ticket files do not exist")

	// Remove invalid json file for next tests
	os.Remove(invalidJsonFile)

	ticket := mockTicket()
	WriteTicketFile(testTicketPath, &ticket)

	// Correct path to ticket files
	errReadTicketFiles3 := ReadTicketFiles(testTicketPath, &tickets)
	assert.Nil(t, errReadTicketFiles3, "An error was returned, although the path is correct")

	os.RemoveAll(testTicketPath + "/")
}

// mockTicket is a helper function to create a dummy ticket for the tests
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
		Id:          "12",
		Name:        "Max Mustermann",
		Mail:        "max.mustermann@example.com",
		Hash:        "$2a$12$n5kluCvuG3wpj18rl46bBexvTX6l0QkD7EQCkgvk1BNby5cNZPLZa",
		IsOnHoliday: false,
	}

	return structs.Ticket{
		Id:       "test123",
		Subject:  "Help",
		Status:   0,
		User:     user,
		Customer: "customer@example.com",
		Entries:  entries,
	}
}

func TestFileExists(t *testing.T) {
	t.Run("existingFile", func(t *testing.T) {
		const existingFile = "../../files/users/users.json"

		assert.True(t, FileExists(existingFile), "users.json file should always exist")
	})

	t.Run("notExistingFile", func(t *testing.T) {
		const notExistingFile = "../../files/users/passwords.json"

		assert.False(t, FileExists(notExistingFile), "nobody would store passwords in a JSON-file, so why should it exist?")
	})
}

func TestDirectoryExists(t *testing.T) {
	t.Run("existingDirectory", func(t *testing.T) {
		const existingDirectory = "../../files/tickets"

		assert.True(t, DirectoryExists(existingDirectory), "tickets directory should exist")
	})

	t.Run("notExistingDirectory", func(t *testing.T) {
		const notExistingDirectory = "../../files/secret_keys"

		assert.False(t, DirectoryExists(notExistingDirectory), "again, secret keys should not be stored here")
	})
}

func TestHasJsonExtension(t *testing.T) {
	t.Run("jsonExtension", func(t *testing.T) {
		assert.True(t, hasJsonExtension("ticket.json"), "should be true because file has a .json file extension")
	})

	t.Run("noJsonExtension", func(t *testing.T) {
		assert.False(t, hasJsonExtension("ticket.xml"), "should be false because file has .xml file extension instead of .json")
	})
}

func mockMail() structs.Mail {
	return structs.Mail{
		Id:      random.CreateRandomId(10),
		From:    "no-reply@trivial-tickets.com",
		To:      "customer@mail.com",
		Subject: "[trivial-tickets] My screen is always black",
		Message: "I cannot see anything on my screen.",
	}
}

func TestWriteReadMailFile(t *testing.T) {
	const mailDirectory = "../../files/testmails"

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
			readMail, mailIdDefined := mails[testMail.Id]
			assert.True(t, mailIdDefined, "test mail id should be defined in read mail map")
			assert.NotNil(t, readMail, "read mail should be non-nil")
		})
	})

	os.RemoveAll(mailDirectory)
}

func TestRemoveMailFile(t *testing.T) {
	const mailDirectory = "../../files/testmails"

	t.Run("existingMail", func(t *testing.T) {
		// Create test ticket to be removed
		testMail := mockMail()
		WriteMailFile(mailDirectory, &testMail)

		removeErr := RemoveMailFile(mailDirectory, testMail.Id)

		assert.NoError(t, removeErr, "removing mail file should not return error since the file exists")
	})

	t.Run("notExistingMail", func(t *testing.T) {
		notExistingMailId := "mail-id"

		removeErr := RemoveMailFile(mailDirectory, notExistingMailId)

		assert.Error(t, removeErr, "remove error should be non-nil because mail file does not exist")
	})

	os.RemoveAll(mailDirectory)
}

func TestWrapAndLogError(t *testing.T) {
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
