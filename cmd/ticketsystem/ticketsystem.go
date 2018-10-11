package main

import (
	"errors"
	"flag"
	"log"
	"math"
)

// Config is a struct to hold the config parameters provided on startup
type Config struct {
	port    int16
	tickets string
	users   string
}

func main() {

	config, err := initConfig()

	if err != nil {
		log.Fatal(err)
	}

	// Log output to use config variable to make it compile
	log.Println("\nProgram initialized with port: ", config.port, "\nThe ticket folder located at ", config.tickets, "\nThe user folder is located at ", config.users)
}

// initConfig parses the command line arguments and
// populates a struct for the config parameters.
// It returns this struct
func initConfig() (Config, error) {

	// Get the command line arguments
	port := flag.Int("port", 443, "Port on which the web server will run")
	tickets := flag.String("tickets", "files/tickets", "Folder in which the tickets will be stored")
	users := flag.String("users", "files/users", "Folder in which the users will be stored")

	// Parse all arguments, e.g. populate the variables
	flag.Parse()

	// If the port is not within boundaries, return an error
	if !isPortInBoundaries(port) {
		return Config{}, errors.New("Port is not a correct port number")
	}

	// Populate and return the struct
	return Config{
		port:    int16(*port),
		tickets: *tickets,
		users:   *users}, nil
}

// isPortInBoundaries returns true if the provided port
// is within the boundaries of a 16 bit integer, false
// otherwise. Since the port numbers only go up to a 16
// bit integer
func isPortInBoundaries(port *int) bool {
	return *port <= math.MaxInt16
}
