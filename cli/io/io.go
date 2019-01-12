package io

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/mortenterhart/trivial-tickets/structs"
	"github.com/mortenterhart/trivial-tickets/util/cliUtils"
	"io"
	"os"
	"strconv"
)

/*
 * Ticketsystem Trivial Tickets
 *
 * Matriculation numbers: 3040018, 3040018, 3478222
 * Lecture:               Programmieren II, INF16B
 * Lecturer:              Herr Prof. Dr. Helmut Neemann
 * Institute:             Duale Hochschule Baden-Württemberg Mosbach
 *
 * ---------------
 *
 * Package io
 * I/O operations for the CLI
 */

var reader = io.Reader(os.Stdin)
var writer = io.Writer(os.Stdout)
var readCom = readCommand
var output = OutputMessageToCommandLine
var verifyEmailAddress = cliUtils.CheckEmailAddress
var readString = getString
var readEmailAddress = getEmailAddress

func readCommand() (structs.Command, error) {
	bufReader := bufio.NewReader(reader)
	var ret structs.Command
	var asInt int
	input, err := bufReader.ReadString('\n')
	// gets rid of the delimiter if there was no error
	if err == nil {
		input = input[:(len(input) - 1)]
	}
	asInt, err = strconv.Atoi(input)
	switch structs.Command(asInt) {
	case structs.FETCH,
		structs.SUBMIT,
		structs.EXIT:
		ret = structs.Command(asInt)
	default:
		err = errors.New("not within range of valid options")

	}
	return ret, err
}

func OutputMessageToCommandLine(output structs.CliMessage) {
	fmt.Fprintf(writer, "%s", string(output))
}

func PrintEmail(mail structs.Mail) {
	fmt.Fprintf(writer, "Receiver: %s\n\n"+
		"Subject: %s\n\n"+
		"%s", mail.To, mail.Subject, mail.Message)
}

func getEmailAddress() (addr string, err error) {
	addr, err = readString()
	if !verifyEmailAddress(addr) {
		err = errors.New(string(structs.InvalidEmail))
	}
	return
}

func getString() (result string, err error) {
	bufReader := bufio.NewReader(reader)
	result, err = bufReader.ReadString('\n')
	// gets rid of the delimiter if there was no error
	if err == nil {
		result = result[:(len(result) - 1)]
	} else if err == io.EOF {
		err = nil
	} else {
		return
	}
	if result == "" {
		err = errors.New(string(structs.EmptyString))
	}
	return
}

func NextCommand() (com structs.Command, err error) {
	counter := 0
	for {
		output(structs.RequestCommandInput)
		if counter > 10 {
			return 0, errors.New(string(structs.AbortExecutionDueToManyWrongUserInputs))
		}
		command, err := readCom()
		if err == nil {
			return command, err
		}
		output(structs.CommandNotAccepted + structs.CliMessage(err.Error()))
		counter++
	}
}

func GetEmail() (jsonMail string, err error) {
	counter := 0
	output(structs.RequestEmailAddress)
	emailAddress, err := readEmailAddress()
	for err != nil {
		if counter > 10 {
			return "", errors.New(string(structs.AbortExecutionDueToManyWrongUserInputs))
		}
		output(structs.CommandNotAccepted + structs.CliMessage(err.Error()) + "\n" + structs.RequestEmailAddress)
		emailAddress, err = readEmailAddress()
		counter++
	}
	counter = 0
	output(structs.RequestTicketID)
	ticketID, err := readString()
	for err != nil && err.Error() != string(structs.EmptyString) {
		if counter > 10 {
			return "", errors.New(string(structs.AbortExecutionDueToManyWrongUserInputs))
		}
		output(structs.CommandNotAccepted + structs.CliMessage(err.Error()) + "\n" + structs.RequestTicketID)
		ticketID, err = readString()
		counter++
	}
	counter = 0
	output(structs.RequestSubject)
	subject, err := readString()
	for err != nil {
		if counter > 10 {
			return "", errors.New(string(structs.AbortExecutionDueToManyWrongUserInputs))
		}
		output(structs.CommandNotAccepted + structs.CliMessage(err.Error()) + "\n" + structs.RequestSubject)
		subject, err = readString()
		counter++
	}
	counter = 0
	output(structs.RequestMessage)
	message, err := readString()
	for err != nil {
		if counter > 10 {
			return "", errors.New(string(structs.AbortExecutionDueToManyWrongUserInputs))
		}
		output(structs.CommandNotAccepted + structs.CliMessage(err.Error()) + "\n" + structs.RequestMessage)
		message, err = readString()
		counter++
	}
	return cliUtils.CreateMail(emailAddress, subject, ticketID, message), nil
}
