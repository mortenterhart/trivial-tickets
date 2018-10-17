package filehandler

import (
	"encoding/json"
	"errors"
	"go-tickets/structs"
	"io/ioutil"
	"log"
	"os"
	"strconv"
)

/*
*
* Matrikelnummern
* 3040018
* 6694964
* 3478222
 */

// ReadUserFile takes a string as parameter for the location
// of the users.json file, reads the content and returns a struct
// of type Users
func ReadUserFile(src string) structs.User {

	// Struct to be populated
	var users structs.User

	// Read contents of users.json
	fileContent, errReadFile := ioutil.ReadFile(src)

	if errReadFile != nil {
		log.Fatal(errReadFile)
	}

	// Unmarshal into users struct
	errUnmarshal := json.Unmarshal(fileContent, &users)

	if errUnmarshal != nil {
		log.Fatal(errUnmarshal)
	}

	return users
}

// CreateFile writes a given ticket to a given path in the json format
func CreateFile(path string, ticket *structs.Ticket) error {

	// TODO: If already exists, overwrite contents

	// Check if path exists. If not, create it
	if _, errExists := os.Stat(path); os.IsNotExist(errExists) {

		errCreateFolders := CreateFolders(path)
		if errCreateFolders != nil {
			log.Fatal(errCreateFolders)
		}
	}

	// Encode the struct with json
	marshalTicket, errMarshalTicket := json.Marshal(ticket)

	if errMarshalTicket != nil {
		log.Fatal(errMarshalTicket)
	}

	// Create the final output path
	finalPath := path + "/" + strconv.Itoa(int(ticket.Id)) + ".json"

	// If the file already exists, return an error. Otherwise write the file
	if _, errExistsFile := os.Stat(finalPath); os.IsNotExist(errExistsFile) {
		return ioutil.WriteFile(finalPath, marshalTicket, 0644)
	} else {
		return errors.New("File already exists on disk.\nPath: " + finalPath)
	}
}

// CreateFolders creates the folders specified in the parameter
func CreateFolders(path string) error {
	return os.MkdirAll(path, os.ModePerm)
}
