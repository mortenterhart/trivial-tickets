package filehandler

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "log"
    "os"
    "path"

    "github.com/mortenterhart/trivial-tickets/structs"
    "github.com/pkg/errors"
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
func ReadUserFile(src string, users *map[string]structs.User) error {

    // Read contents of users.json
    fileContent, errReadFile := ioutil.ReadFile(src)

    if errReadFile != nil {
        log.Println(errReadFile)
        return errReadFile
    }

    // Unmarshal into users hashmap
    errUnmarshal := json.Unmarshal(fileContent, users)

    if errUnmarshal != nil {
        log.Println(errUnmarshal)
        return errUnmarshal
    }

    return nil
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
            log.Println(errCreateFolders)
            return errCreateFolders
        }
    }

    // Encode the struct with json
    marshalTicket, errMarshalTicket := json.MarshalIndent(ticket, "", "   ")

    if errMarshalTicket != nil {
        log.Println(errMarshalTicket)
        return errMarshalTicket
    }

    // Create the final output path
    finalPath := path + "/" + ticket.Id + ".json"

    // Write the file to the given path
    return ioutil.WriteFile(finalPath, marshalTicket, 0644)
}

// WriteMailFile takes a mail and converts it into the JSON format to
// write it into its own file. The directory parameter is a path to
// a directory in which the new file is saved. If it does not exist yet
// it will be created.
func WriteMailFile(directory string, mail *structs.Mail) error {

    if _, existsErr := os.Stat(directory); os.IsNotExist(existsErr) {

        createFoldersErr := CreateFolders(directory)
        if createFoldersErr != nil {
            return wrapAndLogError(createFoldersErr, fmt.Sprintf("could not create directory '%s'", directory))
        }
    }

    marshaledMail, marshalErr := json.MarshalIndent(mail, "", "   ")
    if marshalErr != nil {
        return wrapAndLogError(marshalErr, "could not convert mail to JSON")
    }

    mailFilePath := path.Join(directory, mail.Id+".json")

    writeErr := ioutil.WriteFile(mailFilePath, marshaledMail, 0644)
    if writeErr != nil {
        return wrapAndLogError(writeErr, fmt.Sprintf("error while writing file '%s'", mailFilePath))
    }

    return nil
}

func ReadMailFiles(directory string) (*[]structs.Mail, error) {
    mailFiles, readErr := ioutil.ReadDir(directory)
    if readErr != nil {
        return nil, wrapAndLogError(readErr, "error while reading mail files")
    }

    mails := []structs.Mail{}
    for _, file := range mailFiles {
        jsonMail, readErr := ioutil.ReadFile(path.Join(directory, file.Name()))
        if readErr != nil {
            return nil, wrapAndLogError(readErr, "error while reading mail files")
        }

        var parsedMail structs.Mail
        if parseErr := json.Unmarshal(jsonMail, &parsedMail); parseErr != nil {
            return nil, wrapAndLogError(parseErr, "could not convert JSON mail")
        }

        mails = append(mails, parsedMail)
    }

    return &mails, nil
}

func wrapAndLogError(err error, wrapErrorMessage string) error {
    wrappedError := errors.Wrap(err, wrapErrorMessage)
    log.Println(wrappedError)
    return wrappedError
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
        log.Println(err)
        return err
    }

    // Iterate over each file
    for _, f := range files {

        // Read contents of each ticket file
        fileContent, errReadFile := ioutil.ReadFile(path + "/" + f.Name())

        if errReadFile != nil {
            log.Println(errReadFile)
            return errReadFile
        } else {
            // Create a ticket struct to hold the file contents
            ticket := structs.Ticket{}

            // Unmarshal into a ticket struct
            errUnmarshal := json.Unmarshal(fileContent, &ticket)

            if errUnmarshal != nil {
                log.Println(errUnmarshal)
                return errUnmarshal
            } else {
                // Store the ticket in the tickets hashmap
                (*tickets)[ticket.Id] = ticket
            }
        }
    }

    return nil
}
