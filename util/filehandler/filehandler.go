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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"github.com/pkg/errors"

	"github.com/mortenterhart/trivial-tickets/logger"
	"github.com/mortenterhart/trivial-tickets/structs"
	"github.com/mortenterhart/trivial-tickets/structs/defaults"
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
		return wrapAndLogError(errReadFile, "unable to read users file")
	}

	// Unmarshal into users hash map
	errUnmarshal := json.Unmarshal(fileContent, users)

	if errUnmarshal != nil {
		return wrapAndLogErrorf(errUnmarshal, "unable to decode JSON in users file '%s'", src)
	}

	return nil
}

// WriteUserFile writes the contents of the users map to the
// file system to persist any changes.
func WriteUserFile(dest string, users *map[string]structs.User) error {

	// Create json from the hash map
	usersMarshal, _ := json.MarshalIndent(users, "", "    ")

	// Write json to file
	return ioutil.WriteFile(dest, usersMarshal, defaults.FileModeRegular)
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
		if hasJSONExtension(f.Name()) {

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
			(*tickets)[ticket.ID] = ticket
		}
	}

	return nil
}

// WriteTicketFile writes a given ticket to a given directory
// in the json format. If the ticket already exists, it
// overwrites its contents.
func WriteTicketFile(directory string, ticket *structs.Ticket) error {

	// Check if directory exists. If not, create it
	if !DirectoryExists(directory) {

		logger.Info("Creating missing ticket directory", directory, "for new ticket file")
		errCreateFolders := CreateFolders(directory)
		if errCreateFolders != nil {
			logger.Error(errCreateFolders)
			return errCreateFolders
		}
	}

	// Encode the struct with json
	marshalTicket, errMarshalTicket := json.MarshalIndent(ticket, "", "    ")

	if errMarshalTicket != nil {
		return wrapAndLogError(errMarshalTicket, "could not encode ticket to JSON")
	}

	// Create the final output path
	finalPath := directory + "/" + ticket.ID + ".json"

	// Write the file to the given path
	logger.Info("Writing ticket file", finalPath, "to file system (Permission = 0644 [rw-r--r--])")
	return ioutil.WriteFile(finalPath, marshalTicket, defaults.FileModeRegular)
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
		if hasJSONExtension(file.Name()) {

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
			(*mails)[parsedMail.ID] = parsedMail
		}
	}

	return nil
}

// WriteMailFile takes a mail and converts it into the JSON format to
// write it into its own file. The directory parameter is a path to
// a directory in which the new file is saved. If it does not exist yet
// it will be created.
func WriteMailFile(directory string, mail *structs.Mail) error {

	// If the directory does not exist yet, create it
	if !DirectoryExists(directory) {

		logger.Info("Creating missing mail directory", directory, "for new mail file")
		createFoldersErr := CreateFolders(directory)
		if createFoldersErr != nil {
			return wrapAndLogError(createFoldersErr, fmt.Sprintf("could not create directory '%s'", directory))
		}
	}

	// Encode the mail into JSON
	marshaledMail, marshalErr := json.MarshalIndent(mail, "", "    ")
	if marshalErr != nil {
		return wrapAndLogError(marshalErr, "could not convert mail to JSON")
	}

	// Build the final file path
	mailFilePath := path.Join(directory, mail.ID+".json")

	// Write the JSON mail into the file
	logger.Info("Writing mail file", mailFilePath, "to file system (Permission = 0644 [rw-r--r--])")
	writeErr := ioutil.WriteFile(mailFilePath, marshaledMail, defaults.FileModeRegular)
	if writeErr != nil {
		return wrapAndLogError(writeErr, fmt.Sprintf("error while writing file '%s'", mailFilePath))
	}

	return nil
}

// RemoveMailFile attempts to remove a mail with a given id in a given directory.
// If the file does not exist, it returns an non-nil error.
func RemoveMailFile(directory string, mailID string) error {
	mailPath := path.Join(directory, mailID) + ".json"
	if removeErr := os.Remove(mailPath); removeErr != nil {
		returnErr := fmt.Errorf("could not delete mail file with id '%s'", mailID)
		logger.Errorf("%v: %v", returnErr, removeErr)
		return returnErr
	}

	return nil
}

// CreateFolders creates the folders specified in the parameter.
func CreateFolders(path string) error {
	return os.MkdirAll(path, os.ModePerm)
}

// FileExists examines if a given path exists and is a regular file or not.
// Inspired by https://gist.github.com/mattes/d13e273314c3b3ade33f
func FileExists(path string) bool {
	if file, notExistsErr := os.Stat(path); !os.IsNotExist(notExistsErr) {
		return file.Mode().IsRegular()
	}

	return false
}

// DirectoryExists reports whether the given path exists and points to
// a directory or not.
// Inspired by https://stackoverflow.com/a/49697453
func DirectoryExists(path string) bool {
	if file, notExistsErr := os.Stat(path); !os.IsNotExist(notExistsErr) {
		return file.Mode().IsDir()
	}

	return false
}

// hasJSONExtension checks if the given filename argument conforms to
// the glob pattern "*.json", i.e. if the the filename has a .json file
// extension.
func hasJSONExtension(filename string) bool {
	// Note: The error can never be returned because the glob pattern
	// above is valid and filepath.Match() does not check whether
	// the file with the given filename exists
	jsonFileMatch, _ := filepath.Match("*.json", filename)
	return jsonFileMatch
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
