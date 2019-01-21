// Main package of the ticketsystem webserver
package main

import (
	"flag"
	"fmt"
	"github.com/mortenterhart/trivial-tickets/logger"
	"math"
	"os"

	"github.com/mortenterhart/trivial-tickets/server"
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
 * Package main
 * Main package of the ticketsystem webserver
 */

func main() {

	config, errConfig := initConfig()

	if errConfig != nil {
		logger.Fatal(errConfig)
	}

	exitCode, errServer := server.StartServer(&config)

	if errServer != nil {
		logger.Fatal("fatal error occurred:", errServer, "\nServer startup failed!")
	}

	logger.Info("exiting with exit code:", exitCode)
	os.Exit(exitCode)
}

// initConfig parses the command line arguments and
// populates a struct for the config parameters.
// It returns this struct
func initConfig() (structs.Config, error) {

	// Get the command line arguments
	port := flag.Int("port", 8443, "`port` on which the web server will run")
	tickets := flag.String("tickets", "../../files/tickets", "`directory` in which the tickets will be stored")
	users := flag.String("users", "../../files/users/users.json", "`directory` where the users file is stored")
	mails := flag.String("mails", "../../files/mails", "`directory` in which the mails will be cached")
	cert := flag.String("cert", "../../ssl/server.cert", "location of the ssl certificate `file`")
	key := flag.String("key", "../../ssl/server.key", "location of the ssl key `file`")
	web := flag.String("web", "../../www", "location of the www `directory`")

	// Set another usage message function
	flag.Usage = usageMessage

	// Parse all arguments, e.g. populate the variables
	flag.Parse()

	// If the port is not within boundaries, return an error
	if !isPortInBoundaries(*port) {
		return structs.Config{}, fmt.Errorf("applied port %d is not a correct port number", *port)
	}

	// Populate and return the struct
	return structs.Config{
		Port:    int16(*port),
		Tickets: *tickets,
		Users:   *users,
		Mails:   *mails,
		Cert:    *cert,
		Key:     *key,
		Web:     *web,
	}, nil
}

// isPortInBoundaries returns true if the provided port
// is within the boundaries of a 16 bit integer, false
// otherwise. Since the port numbers only go up to a 16
// bit integer
func isPortInBoundaries(port int) bool {
	return port <= math.MaxInt16
}

func usageMessage() {
	fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [options]\n\n", os.Args[0])
	fmt.Fprintf(flag.CommandLine.Output(), "options may be one of the following:\n")

	flag.PrintDefaults()
}
