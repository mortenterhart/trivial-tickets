package main

import (
	"flag"
	"log"
)

// Config is a struct to hold the config parameters provided on startup
type Config struct {
	port    int
	tickets string
	users   string
}

func main() {

	config := initConfig()
	log.Println("\nProgram initialized with port: ", config.port, "\nThe ticket folder located at ", config.tickets, "\nThe user folder is located at ", config.users) // Log output to use config variable to make it compile
}

// initConfig parses the command line arguments and
// populates a struct for the config parameters.
// It returns this struct
func initConfig() Config {

	// Get the command line arguments
	port := flag.Int("port", 443, "Port on which the web server will run")
	tickets := flag.String("tickets", "files/tickets", "Folder in which the tickets will be stored")
	users := flag.String("users", "files/users", "Folder in which the users will be stored")

	// Parse all arguments, e.g. populate the variables
	flag.Parse()

	// Populate and return the struct
	return Config{
		port:    *port,
		tickets: *tickets,
		users:   *users}
}
