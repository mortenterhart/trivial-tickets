// Main package of the command line utility
package main

import (
	"flag"
	"fmt"
	"github.com/mortenterhart/trivial-tickets/cli/communicationToServer"
	"github.com/mortenterhart/trivial-tickets/cli/io"
	"github.com/mortenterhart/trivial-tickets/structs"
	"github.com/mortenterhart/trivial-tickets/util/cliUtils"
	"log"
	"math"
	"os"
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
 * Main package of the command line utility
 */

func main() {

	conf, fetch, submit, mail := getConfig()
	communicationToServer.SetServerConfig(conf)

	// based on the flags submit and fetch either start the commandLoop or directly invoke FetchEmail / SubmitEmail
	switch {
	case !submit && !fetch:
		commandLoop()
	case !submit && fetch:
		mails, err := communicationToServer.FetchEmails()
		if err != nil {
			log.Fatal(err)
		}
		for _, mail := range mails {
			io.PrintEmail(mail)
			acknowledgementError := communicationToServer.AcknowledgeEmailReception(mail)
			if acknowledgementError != nil {
				println(acknowledgementError.Error())
			}
		}
	case submit && !fetch:
		communicationToServer.SubmitEmail(mail)
	default:
		log.Fatal(structs.NoValidOption)
	}

}

// Here's how it should behave:
// On startup, show the available commands (fetch, submit and exit)
// On fetch just print out all available emails in the command line.
// ON submit ask in order for email, ticketID, subject and message.
// After both fetch and submit, return to startup message.
// On exit, close the application.

// When starting the Application it should be possible to specify some parameters.
// Based on these parameters the IP address of the server and its port should be selected.
// The Parameters should also offer a way to directly invoke the send / fetch mail commands.

// getConfig parses the command line flags. It returns a CLIConfig struct, a mail struct and the fetch and submit flags as boolean.
// The port number is checked to be valid. There are no checks performed on the validity of the created Mail struct.
func getConfig() (conf structs.CLIConfig, fetch bool, submit bool, mail string) {
	IPAddr := flag.String("ip", "localhost", "IP address of the server")
	port := flag.Uint("port", 8443, "Port the server listens to")
	cert := flag.String("cert", "./ssl/server.cert", "Location of the ssl certificate")
	f := flag.Bool("f", false, "fetch (fetch): If set, the application will fetch all messages from the server.")
	s := flag.Bool("s", false, "Use to submit a message to the server. Requires -email, -subject, -message. The use of -tID is optional.")
	email := flag.String("email", "", "The email address of the sender")
	ticketID := flag.String("tID", "", "ID of the related Ticket. If left empty, a new ticket is created")
	subject := flag.String("subject", "", "The subject of the message")
	message := flag.String("message", "", "The body of the message.")

	flag.Usage = usageMessage

	flag.Parse()

	if *port > math.MaxUint16 {
		log.Fatal("Port is not a valid port number.")
	}

	conf = structs.CLIConfig{
		IPAddr: *IPAddr,
		Port:   uint16(*port),
		Cert:   *cert,
	}

	fetch = *f
	submit = *s
	mail = cliUtils.CreateMail(*email, *subject, *ticketID, *message)
	return
}

// commandLoop is called when neither the fetch nor the submit flag are set. It guides the User through the operation of the CLI.
// commandLoop returns when either the user selects an exit command, when the user fails to input a valid command N times in a row, or when FetchEmails / SubmitEmail returns an error
func commandLoop() {
	ok := true
	for ok {
		com, err := io.NextCommand()
		if err != nil {
			log.Fatal(err)
		}
		switch com {
		case structs.FETCH:
			mails, err := communicationToServer.FetchEmails()
			if err != nil {
				log.Fatal(err)
			}
			for _, mail := range mails {
				io.PrintEmail(mail)
				acknowledgementError := communicationToServer.AcknowledgeEmailReception(mail)
				if acknowledgementError != nil {
					println(acknowledgementError.Error())
				}
			}
		case structs.SUBMIT:
			mailJson, err := io.GetEmail()
			if err != nil {
				log.Fatal(err)
			}
			err = communicationToServer.SubmitEmail(mailJson)
			if err != nil {
				log.Fatal(err)
			}
		case structs.EXIT:
			ok = false
		default:
			log.Fatal(com, err)
		}
	}
}

func usageMessage() {
	fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [options]\n\n", os.Args[0])
	fmt.Fprintf(flag.CommandLine.Output(), "options may be one of the following:\n")

	flag.PrintDefaults()
}
