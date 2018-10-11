package main

import (
	"log"
	"os"
	"strconv"
)

// Config is a struct to hold the config parameters provided on startup
type Config struct {
	port         int
	ticketFolder string
}

func main() {

	config := initConfig(os.Args)
	log.Println("Program initialized with port: ", config.port, " and ticket folder located at ", config.ticketFolder) // Log output to use config variable to make it compile
}

// initConfig takes the command line arguments parsed via
// and populates a struct for the config parameters.
// It returns this struct
func initConfig(arguments []string) Config {

	// Init the struct
	config := Config{}

	// Loop through the command line arguments
	for i := 1; i < len(arguments); i++ {

		// Check for the argument flags and set the struct fields
		// to the corresponding value being passed
		switch arguments[i] {

		case "-port":
			config.port, _ = strconv.Atoi(arguments[i+1])
		case "-ticketFolder":
			config.ticketFolder = arguments[i+1]
		}
	}

	return config
}

/* Refactoring for the above code, not implemented until a unit test for this is created
return Config{
	port:         *flag.Int("port", 443, "Port on which the web server will run"),
	ticketFolder: *flag.String("ticketFolder", "files/tickets", "Folder in which the tickets will be stored")}
*/
