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

// Package io contains I/O operations for the CLI
package io

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/mortenterhart/trivial-tickets/structs"
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
 * Package io
 * I/O operations for the CLI
 */

// Reader is the reader used to read inputs from
// the standard input.
var Reader = io.Reader(os.Stdin)

// Writer is the writer that is used to write the
// e-mails, messages and prompts to the standard
// output.
var Writer = io.Writer(os.Stdout)

// maxInputAttempts is the maximum amount of permitted
// successive user inputs before the prompt is quit.
const maxInputAttempts int = 10

// These variables are used within tests to redefine
// functions that are expensive to test.
var readCom = readCommand
var output = OutputMessageToCommandLine
var verifyEmailAddress = cliutils.CheckEmailAddress
var readString = getString
var readEmailAddress = getEmailAddress

// readCommand reads and validates a command
// entered by the user.
func readCommand() (structs.Command, error) {
	bufReader := bufio.NewReader(Reader)
	var ret structs.Command
	var asInt int
	input, err := bufReader.ReadString('\n')
	Reader = bufio.NewReader(bufReader)
	// gets rid of the delimiter if there was no error
	if err == nil {
		input = input[:(len(input) - 1)]
		// if it's a windows machine remove carriage return
		if len(input) > 0 && input[len(input)-1] == '\r' {
			input = input[:(len(input) - 1)]
		}
	}
	asInt, err = strconv.Atoi(input)
	switch structs.Command(asInt) {
	case structs.CommandFetch,
		structs.CommandSubmit,
		structs.CommandExit:
		ret = structs.Command(asInt)
	default:
		err = errors.New("not within range of valid options")

	}
	return ret, err
}

// OutputMessageToCommandLine writes the message of
// an email to the output.
func OutputMessageToCommandLine(output structs.CliMessage) {
	fmt.Fprint(Writer, string(output))
}

// PrintEmail outputs a received e-mail to
// the console.
func PrintEmail(mail structs.Mail) {
	fmt.Fprintf(Writer, "From: %s\n"+
		"To: %s\n\n"+
		"Subject: %s\n\n"+
		"%s\n", mail.From, mail.To, mail.Subject, mail.Message)
}

// getEmailAddress reads and validates an
// email input.
func getEmailAddress() (addr string, err error) {
	addr, err = readString()
	if !verifyEmailAddress(addr) {
		err = errors.New(string(structs.InvalidEmail))
	}
	return
}

// getString makes an user input and returns the
// result and an error if something went wrong.
func getString() (result string, err error) {
	bufReader := bufio.NewReader(Reader)
	result, err = bufReader.ReadString('\n')
	Reader = bufio.NewReader(bufReader)
	// gets rid of the delimiter if there was no error
	if err == nil {
		result = result[:(len(result) - 1)]
		// if it's a windows machine remove carriage return
		if len(result) > 0 && result[len(result)-1] == '\r' {
			result = result[:(len(result) - 1)]
		}
	} else if err == io.EOF {
		err = nil
	}

	if result == "" {
		err = errors.New(string(structs.EmptyString))
	}
	return
}

// NextCommand prompts for a new command. If the
// command was entered invalidly 10 times in a row
// the function exits with an error.
func NextCommand() (com structs.Command, err error) {
	counter := 0
	for {
		output(structs.RequestCommandInput)
		if counter > maxInputAttempts {
			return 0, errors.New(string(structs.TooManyInputs))
		}
		command, err := readCom()
		if err == nil {
			return command, err
		}
		output(structs.CommandNotAccepted + structs.CliMessage(err.Error()))
		counter++
	}
}

// GetEmail prompts the user to enter the required
// information for a new e-mail to be sent to the
// server. These are the e-mail address, the optional
// ticket id, the subject and the message.
func GetEmail() (jsonMail string, err error) {
	counter := 0
	output(structs.RequestEmailAddress)
	emailAddress, err := readEmailAddress()
	for err != nil {
		if counter > maxInputAttempts {
			return "", errors.New(string(structs.TooManyInputs))
		}
		output(structs.CommandNotAccepted + structs.CliMessage(err.Error()) + "\n" + structs.RequestEmailAddress)
		emailAddress, err = readEmailAddress()
		counter++
	}
	counter = 0
	output(structs.RequestTicketID)
	ticketID, err := readString()
	for err != nil && err.Error() != string(structs.EmptyString) {
		if counter > maxInputAttempts {
			return "", errors.New(string(structs.TooManyInputs))
		}
		output(structs.CommandNotAccepted + structs.CliMessage(err.Error()) + "\n" + structs.RequestTicketID)
		ticketID, err = readString()
		counter++
	}
	counter = 0
	output(structs.RequestSubject)
	subject, err := readString()
	for err != nil {
		if counter > maxInputAttempts {
			return "", errors.New(string(structs.TooManyInputs))
		}
		output(structs.CommandNotAccepted + structs.CliMessage(err.Error()) + "\n" + structs.RequestSubject)
		subject, err = readString()
		counter++
	}
	counter = 0
	output(structs.RequestMessage)
	message, err := readString()
	for err != nil {
		if counter > maxInputAttempts {
			return "", errors.New(string(structs.TooManyInputs))
		}
		output(structs.CommandNotAccepted + structs.CliMessage(err.Error()) + "\n" + structs.RequestMessage)
		message, err = readString()
		counter++
	}
	return cliutils.CreateMail(emailAddress, subject, ticketID, message), nil
}
