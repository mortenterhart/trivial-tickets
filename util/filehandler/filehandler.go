package filehandler

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"os"

	"github.com/mortenterhart/trivial-tickets/structs"
)

/*
*
* Matrikelnummern
* 3040018
* 6694964
* 3478222
 */

// ReadUserFile takes a string as parameter for the location
// of the users.json file, reads the content and stores it inside
// of the hashmap for the users
func ReadUserFile(src string, users *map[string]structs.User) {

	// Read contents of users.json
	fileContent, errReadFile := ioutil.ReadFile(src)

	if errReadFile != nil {
		log.Print(errReadFile)
	}

	// Unmarshal into users hashmap
	errUnmarshal := json.Unmarshal(fileContent, users)

	if errUnmarshal != nil {
		log.Print(errUnmarshal)
	}
}

// WriteUserFile writes the contents of the users map to the
// file system to persist any changes
func WriteUserFile(dest string, users *map[string]structs.User) error {

	// Create json from the hashmap
	usersMarshal, _ := json.MarshalIndent(users, "", "   ")

	// Write json to file
	return ioutil.WriteFile(dest, usersMarshal, 0644)
}

// WriteTicketFile writes a given ticket to a given path in the json format.
// If the ticket already exists, it overwrites its contents
func WriteTicketFile(path string, ticket *structs.Ticket) error {

	// Check if path exists. If not, create it
	if _, errExists := os.Stat(path); os.IsNotExist(errExists) {

		errCreateFolders := CreateFolders(path)
		if errCreateFolders != nil {
			log.Print(errCreateFolders)
		}
	}

	// Encode the struct with json
	marshalTicket, errMarshalTicket := json.MarshalIndent(ticket, "", "   ")

	if errMarshalTicket != nil {
		log.Print(errMarshalTicket)
	}

	// Create the final output path
	finalPath := path + "/" + ticket.Id + ".json"

	// Write the file to the given path
	return ioutil.WriteFile(finalPath, marshalTicket, 0644)
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
		log.Print(err)
		return errors.New("Unable to load the ticket files directory")
	}

	// Iterate over each file
	for _, f := range files {

		// Read contents of each ticket file
		fileContent, errReadFile := ioutil.ReadFile(path + "/" + f.Name())

		if errReadFile != nil {
			log.Print(errReadFile)
		} else {
			// Create a ticket struct to hold the file contents
			ticket := structs.Ticket{}

			// Unmarshal into a ticket struct
			errUnmarshal := json.Unmarshal(fileContent, &ticket)

			if errUnmarshal != nil {
				log.Print(errUnmarshal)
			} else {
				// Store the ticket in the tickets hashmap
				(*tickets)[ticket.Id] = ticket
			}
		}
	}

	return nil
}
