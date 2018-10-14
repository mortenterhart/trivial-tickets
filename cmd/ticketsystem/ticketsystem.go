package main

import (
	"errors"
	"flag"
	"go-tickets/server"
	"go-tickets/structs"
	"log"
	"math"
)

/*
*
* Matrikelnummern
* 3040018
*
 */

func main() {

	config, errConfig := initConfig()

	if errConfig != nil {
		log.Fatal(errConfig)
	}

	errServer := server.StartServer(&config)

	if errServer != nil {
		log.Fatal(errServer)
	}
}

// initConfig parses the command line arguments and
// populates a struct for the config parameters.
// It returns this struct
func initConfig() (structs.Config, error) {

	// Get the command line arguments
	port := flag.Int("port", 443, "Port on which the web server will run")
	tickets := flag.String("tickets", "files/tickets", "Folder in which the tickets will be stored")
	users := flag.String("users", "files/users", "Folder in which the users will be stored")
	cert := flag.String("cert", "server.cert", "Folder in which the ssl certificate is located")
	key := flag.String("key", "server.key", "Folder in which the ssl key file is located")

	// Parse all arguments, e.g. populate the variables
	flag.Parse()

	// If the port is not within boundaries, return an error
	if !isPortInBoundaries(port) {
		return structs.Config{}, errors.New("Port is not a correct port number")
	}

	// Populate and return the struct
	return structs.Config{
		Port:    int16(*port),
		Tickets: *tickets,
		Users:   *users,
		Cert:    *cert,
		Key:     *key}, nil
}

// isPortInBoundaries returns true if the provided port
// is within the boundaries of a 16 bit integer, false
// otherwise. Since the port numbers only go up to a 16
// bit integer
func isPortInBoundaries(port *int) bool {
	return *port <= math.MaxInt16
}
