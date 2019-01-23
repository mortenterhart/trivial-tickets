// Interactions with files, writing and reading files and persisting
// changes to the file system
package filehandler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"github.com/pkg/errors"

	"github.com/mortenterhart/trivial-tickets/logger"
	"github.com/mortenterhart/trivial-tickets/structs"
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
 * Package filehandler
 * Interactions with files, writing and reading files and persisting
 * changes to the file system
 */

// ReadUserFile takes a string as parameter for the location
// of the users.json file, reads the content and stores it inside
// of the hash map for the users
func ReadUserFile(src string, users *map[string]structs.User) error {

	// Read contents of users.json
	fileContent, errReadFile := ioutil.ReadFile(src)

	if errReadFile != nil {
		logger.Error(errReadFile)
		return errReadFile
	}

	// Unmarshal into users hashmap
	errUnmarshal := json.Unmarshal(fileContent, users)

	if errUnmarshal != nil {
		logger.Error(errUnmarshal)
		return errUnmarshal
	}

	return nil
}

// WriteUserFile writes the contents of the users map to the
// file system to persist any changes
func WriteUserFile(dest string, users *map[string]structs.User) error {

	// Create json from the hash map
	usersMarshal, _ := json.MarshalIndent(users, "", "   ")

	// Write json to file
	return ioutil.WriteFile(dest, usersMarshal, 0644)
}

// WriteTicketFile writes a given ticket to a given path in the json format.
// If the ticket already exists, it overwrites its contents
func WriteTicketFile(path string, ticket *structs.Ticket) error {

	// Check if path exists. If not, create it
	if !DirectoryExists(path) {

		errCreateFolders := CreateFolders(path)
		if errCreateFolders != nil {
			logger.Error(errCreateFolders)
			return errCreateFolders
		}
	}

	// Encode the struct with json
	marshalTicket, errMarshalTicket := json.MarshalIndent(ticket, "", "   ")

	if errMarshalTicket != nil {
		logger.Error(errMarshalTicket)
		return errMarshalTicket
	}

	// Create the final output path
	finalPath := path + "/" + ticket.Id + ".json"

	// Write the file to the given path
	return ioutil.WriteFile(finalPath, marshalTicket, 0644)
}

// FileExists examines if a given path exists and is a regular file or not.
// Inspired by https://gist.github.com/mattes/d13e273314c3b3ade33f
func FileExists(path string) bool {
	if file, notExistsErr := os.Stat(path); notExistsErr == nil && !os.IsNotExist(notExistsErr) {
		return file.Mode().IsRegular()
	}

	return false
}

// DirectoryExists reports whether the given path exists and points to
// a directory or not.
// Inspired by https://stackoverflow.com/a/49697453
func DirectoryExists(path string) bool {
	if file, notExistsErr := os.Stat(path); notExistsErr == nil && !os.IsNotExist(notExistsErr) {
		return file.Mode().IsDir()
	}

	return false
}

// hasJsonExtension checks if the given filename argument conforms to
// the glob pattern "*.json", i.e. if the the filename has a .json file
// extension.
func hasJsonExtension(filename string) bool {
	jsonFileMatch, globErr := filepath.Match("*.json", filename)
	if globErr != nil {
		// This error can never be returned because the glob pattern above is valid
		// and it is not checked whether the file with the given filename exists
		logger.Error("file name matching failed due to syntax error in glob pattern: \"*.json\"")
		return false
	}

	return jsonFileMatch
}

// WriteMailFile takes a mail and converts it into the JSON format to
// write it into its own file. The directory parameter is a path to
// a directory in which the new file is saved. If it does not exist yet
// it will be created.
func WriteMailFile(directory string, mail *structs.Mail) error {

	// If the directory does not exist yet, create it
	if !DirectoryExists(directory) {

		createFoldersErr := CreateFolders(directory)
		if createFoldersErr != nil {
			return wrapAndLogError(createFoldersErr, fmt.Sprintf("could not create directory '%s'", directory))
		}
	}

	// Encode the mail into JSON
	marshaledMail, marshalErr := json.MarshalIndent(mail, "", "   ")
	if marshalErr != nil {
		return wrapAndLogError(marshalErr, "could not convert mail to JSON")
	}

	// Build the final file path
	mailFilePath := path.Join(directory, mail.Id+".json")

	// Write the JSON mail into the file
	writeErr := ioutil.WriteFile(mailFilePath, marshaledMail, 0644)
	if writeErr != nil {
		return wrapAndLogError(writeErr, fmt.Sprintf("error while writing file '%s'", mailFilePath))
	}

	return nil
}

// ReadMailFiles lookups the files in the given directory, reads them and decodes
// JSON files into mail structures. Those structures are added to a mail hash map
// with its id as key.
func ReadMailFiles(directory string, mails *map[string]structs.Mail) error {

	// Read directory contents
	mailFiles, readErr := ioutil.ReadDir(directory)
	if readErr != nil {
		return wrapAndLogError(readErr, "error while reading mail files")
	}

	// Iterate over each file in the directory
	for _, file := range mailFiles {

		// Only read files with .json file extension and
		// ignore all other files in that directory. Otherwise
		// this would cause errors when unmarshaling the file
		// contents which is expected to be JSON.
		if hasJsonExtension(file.Name()) {

			// Read the mail from the .json file
			jsonMail, readErr := ioutil.ReadFile(path.Join(directory, file.Name()))
			if readErr != nil {
				return wrapAndLogErrorf(readErr, "error while reading mail file '%s/%s'", directory, file.Name())
			}

			// Parse the read JSON into a mail struct
			var parsedMail structs.Mail
			if parseErr := json.Unmarshal(jsonMail, &parsedMail); parseErr != nil {
				return wrapAndLogErrorf(parseErr, "could not decode JSON in mail file '%s/%s'", directory, file.Name())
			}

			// Add the parsed mail to the mail hash map
			(*mails)[parsedMail.Id] = parsedMail
		}
	}

	return nil
}

// RemoveMailFile attempts to remove a mail with a given id in a given directory.
// If the file does not exist, it returns an non-nil error.
func RemoveMailFile(directory string, mailId string) error {
	mailPath := path.Join(directory, mailId) + ".json"
	if removeErr := os.Remove(mailPath); removeErr != nil {
		returnErr := fmt.Errorf("could not delete mail file with id '%s'", mailId)
		logger.Errorf("%v: %v", returnErr, removeErr)
		return returnErr
	}

	return nil
}

// wrapAndLogError creates a new error by wrapping an error message around an
// existing error. Before it returns the error, the function logs its message
// to the console.
func wrapAndLogError(err error, wrapErrorMessage string) error {
	wrappedError := errors.Wrap(err, wrapErrorMessage)
	logger.Error(wrappedError)
	return wrappedError
}

// wrapAndLogErrorf creates a new error by wrapping an error message around an
// existing error. The error message is constructed using a printf like format
// string and the corresponding arguments to the placeholders. Before it returns
// the error, the function logs its message to the console.
func wrapAndLogErrorf(err error, errorFormat string, arguments ...interface{}) error {
	return wrapAndLogError(err, fmt.Sprintf(errorFormat, arguments...))
}

// CreateFolders creates the folders specified in the parameter
func CreateFolders(path string) error {
	return os.MkdirAll(path, os.ModePerm)
}

// ReadTicketFiles reads all the tickets into memory at the server start
func ReadTicketFiles(path string, tickets *map[string]structs.Ticket) error {

	// Get all the files in given directory
	files, err := ioutil.ReadDir(path)
	if err != nil {
		logger.Error(err)
		return err
	}

	// Iterate over each file
	for _, f := range files {

		// Only read files with .json file extension and
		// ignore all other files in that directory. Otherwise
		// this would cause errors when unmarshaling the file
		// contents which is expected to be JSON.
		if hasJsonExtension(f.Name()) {

			// Read contents of each ticket file
			fileContent, errReadFile := ioutil.ReadFile(path + "/" + f.Name())

			if errReadFile != nil {
				return wrapAndLogErrorf(errReadFile, "error while reading ticket file '%s/%s'", path, f.Name())
			}

			// Create a ticket struct to hold the file contents
			ticket := structs.Ticket{}

			// Unmarshal into a ticket struct
			errUnmarshal := json.Unmarshal(fileContent, &ticket)

			if errUnmarshal != nil {
				return wrapAndLogErrorf(errUnmarshal, "could not decode JSON in ticket file '%s/%s'", path, f.Name())
			}

			// Store the ticket in the tickets hash map
			(*tickets)[ticket.Id] = ticket
		}
	}

	return nil
}
