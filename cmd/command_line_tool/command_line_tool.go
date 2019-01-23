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

// Command commandLineTool is the command-line interface to
// send and receive e-mails to/from the Ticketsystem server.
package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"path/filepath"

	"github.com/mortenterhart/trivial-tickets/cli/client"
	"github.com/mortenterhart/trivial-tickets/cli/io"
	"github.com/mortenterhart/trivial-tickets/structs"
	"github.com/mortenterhart/trivial-tickets/structs/defaults"
	"github.com/mortenterhart/trivial-tickets/util/cliutils"
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

var (
	// Command-line options

	// Client options
	host = flag.String("host", defaults.CliHost, "IP address of the server")
	port = flag.Uint("port", uint(defaults.CliPort), "Port the server listens to")
	cert = flag.String("cert", defaults.CliCertificate, "Location of the ssl certificate")

	// Fetch options
	fetch = flag.Bool("f", defaults.CliFetch, "fetch (fetch): If set, the application will fetch all messages from the server.")

	// Submit options
	submit   = flag.Bool("s", defaults.CliSubmit, "Use to submit a message to the server. Requires -email, -subject, -message. The use of -tID is optional.")
	email    = flag.String("email", "", "The email address of the sender")
	ticketID = flag.String("tID", "", "ID of the related Ticket. If left empty, a new ticket is created")
	subject  = flag.String("subject", "", "The subject of the message")
	message  = flag.String("message", "", "The body of the message.")
)

// fatal is a function used to log a fatal error
// which causes the program to exit. It is necessary
// for the tests that the fatal function is
// substitutable so that the tests don't exit.
var fatal = log.Fatal

// main is the main entry point to the command-line tool.
func main() {

	conf, fetch, submit, mail := getConfig()
	client.SetCLIConfig(conf)

	// Based on the flags submit and fetch either start
	// the commandLoop or directly invoke FetchEmails / SubmitEmail
	switch {
	case !submit && !fetch:
		err := commandLoop()
		if err != nil {
			fatal(err)
			return
		}

	case !submit && fetch:
		mails, err := client.FetchEmails()
		if err != nil {
			fatal(err)
			return
		}
		for _, mail := range mails {
			io.PrintEmail(mail)
			acknowledgementError := client.AcknowledgeEmailReception(mail)
			if acknowledgementError != nil {
				fmt.Println(acknowledgementError.Error())
			}
		}

	case submit && !fetch:
		submitErr := client.SubmitEmail(mail)
		if submitErr != nil {
			fatal(submitErr)
			return
		}

	default:
		fatal(structs.NoValidOption)
	}
}

// Here is how it should behave:
// On startup, show the available commands (fetch, submit
// and exit). On fetch just print out all available emails
// in the command line. On submit ask in order for email,
// ticketID, subject and message. After both fetch and submit,
// return to startup message. On exit, close the application.

// When starting the Application it should be possible
// to specify some parameters. Based on these parameters
// the IP address of the server and its port should be
// selected. The Parameters should also offer a way to
// directly invoke the send / fetch mail commands.

// getConfig parses the command line flags. It returns a
// CLIConfig struct, a mail struct and the fetch and submit
// flags as boolean. The port number is checked to be valid.
// There are no checks performed on the validity of the
// created Mail struct.
func getConfig() (conf structs.CLIConfig, fetchFlag bool, submitFlag bool, mail string) {
	flag.Usage = usageMessage

	flag.Parse()

	if *port > math.MaxUint16 {
		fatal(fmt.Sprintf("Port '%d' is not a valid port number.", *port))
		return
	}

	conf = structs.CLIConfig{
		Host: *host,
		Port: uint16(*port),
		Cert: *cert,
	}

	fetchFlag = *fetch
	submitFlag = *submit
	mail = cliutils.CreateMail(*email, *subject, *ticketID, *message)
	return
}

// commandLoop is called when neither the fetch nor the submit flag
// are set. It guides the User through the operation of the CLI.
// commandLoop returns when either the user selects an exit command,
// when the user fails to input a valid command N times in a row, or
// when FetchEmails / SubmitEmail returns an error.
func commandLoop() error {
	ok := true
	for ok {
		com, err := io.NextCommand()
		if err != nil {
			return err
		}

		switch com {
		case structs.FETCH:
			mails, err := client.FetchEmails()
			if err != nil {
				return err
			}
			for _, mail := range mails {
				io.PrintEmail(mail)
				acknowledgementError := client.AcknowledgeEmailReception(mail)
				if acknowledgementError != nil {
					fmt.Println(acknowledgementError.Error())
				}
			}
		case structs.SUBMIT:
			mailJSON, err := io.GetEmail()
			if err != nil {
				return err
			}
			err = client.SubmitEmail(mailJSON)
			if err != nil {
				return err
			}
		case structs.EXIT:
			ok = false
		}
	}

	return nil
}

// usageMessage prints a complete help text about the usage
// of the command-line tool to the output writer.
func usageMessage() {
	w := flag.CommandLine.Output()
	fmt.Fprintf (w, "Usage: %s [options]\n", filepath.Base(os.Args[0]))
	fmt.Fprintln(w, "Trivial Tickets Command-line Tool")
	fmt.Fprintln(w)

	fmt.Fprintln(w, "The command-line tool can submit emails to the web server")
	fmt.Fprintln(w, "and create new tickets or answers this way. It offers an")
	fmt.Fprintln(w, "interactive menu to create these emails when called without")
	fmt.Fprintln(w, "options. The tool can also fetch created emails by the server")
	fmt.Fprintln(w, "and sends acknowledgements for each received email.")
	fmt.Fprintln(w)

	fmt.Fprintln(w, "When called either with the -f or the -s flag, the interactive")
	fmt.Fprintln(w, "menu is skipped and the request is done directly. The following")
	fmt.Fprintln(w, "options are available.")
	fmt.Fprintln(w)

	fmt.Fprintln(w, "Help options:")
	fmt.Fprintln(w, "  -h, -help       Print this help text and exit.")
	fmt.Fprintln(w)

	fmt.Fprintln(w, "Client options:")
	fmt.Fprintln(w, "  -host <HOST>   Change the host name of the server the client is")
	fmt.Fprintln(w, "                 connecting to. HOST must be an existing host name")
	fmt.Fprintln(w, "                 or IP address.")
	fmt.Fprintf (w, "                 (Default: \"%s\")\n", defaults.CliHost)
	fmt.Fprintln(w, "  -port <PORT>   The port the server listens to. PORT has to be a")
	fmt.Fprintf (w, "                 16 bit unsigned integer (0 < PORT <= %d) and", math.MaxUint16)
	fmt.Fprintln(w, "                 must be the port used by the server process.")
	fmt.Fprintf (w, "                 (Default: %d)\n", defaults.CliPort)
	fmt.Fprintln(w, "  -cert <FILE>   The path to the ssl certificate file. FILE has to")
	fmt.Fprintln(w, "                 be an existing file with a valid ssl certificate.")
	fmt.Fprintf (w, "                 (Default: \"%s\")\n", defaults.CliCertificate)
	fmt.Fprintln(w)

	fmt.Fprintln(w, "Fetch options:")
	fmt.Fprintln(w, "  -f             Fetch emails from the server and skip the interactive")
	fmt.Fprintln(w, "                 menu. Each mail is verified against the server to be")
	fmt.Fprintln(w, "                 sent.")
	fmt.Fprintln(w)

	fmt.Fprintln(w, "Submit options:")
	fmt.Fprintln(w, "  -s             Use this flag to directly submit an email to the server.")
	fmt.Fprintln(w, "                 The flags -email, -subject and -message are required to")
	fmt.Fprintln(w, "                 set the email's properties. The -tID flag can be specified")
	fmt.Fprintln(w, "                 to provide an optional ticket id in order to create an")
	fmt.Fprintln(w, "                 rather than a new ticket.")
	fmt.Fprintln(w, "  -email <EMAIL> Provide the email address of the sender. EMAIL should be")
	fmt.Fprintln(w, "                 a valid email address.")
	fmt.Fprintln(w, "  -subject <STRING>")
	fmt.Fprintln(w, "                 Provide the subject of the email.")
	fmt.Fprintln(w, "  -message <STRING>")
	fmt.Fprintln(w, "                 Provide the content of the email.")
	fmt.Fprintln(w, "  -tID <ID>      (optional) Provide the ticket id of an existing ticket")
	fmt.Fprintln(w, "                 to create a new answer to this ticket rather than a new")
	fmt.Fprintln(w, "                 ticket.")
}
