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
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mortenterhart/trivial-tickets/logger/testlogger"
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
 * Package io [tests]
 * I/O operations for the CLI
 */

// TestWriter
type TestWriter struct {
	output *string
}

func newTestWriter(out *string) (w *TestWriter) {
	return &TestWriter{out}
}

func (w *TestWriter) Write(p []byte) (n int, err error) {
	*w.output = string(p)
	return len(p), nil
}

type ITestReader interface {
	io.Reader
	setData(d string)
}

type TestReader struct {
	data []byte
}

// newTestReader creates a new test reader
func newTestReader(data string) (r *TestReader) {
	return &TestReader{[]byte(data)}
}

func (r *TestReader) setData(d string) {
	r.data = []byte(d)
}

func (r *TestReader) readByte() byte {
	// this function assumes that eof() check was done before
	b := r.data[0]
	r.data = r.data[1:]
	return b
}

func (r *TestReader) eof() (eof bool) {
	return len(r.data) == 0
}

func (r *TestReader) Read(p []byte) (n int, err error) {
	if r.eof() {
		err = io.EOF
		return
	}

	if c := cap(p); c > 0 {
		for n < c {
			p[n] = r.readByte()
			n++
			if r.eof() {
				break
			}
		}
	}
	return
}

func TestReadCommandSuccess(t *testing.T) {
	testlogger.BeginTest()
	defer testlogger.EndTest()

	r := newTestReader(structs.FETCH.String())
	Reader = r

	command, err := readCommand()

	assert.Equal(t, structs.FETCH, command, "does not return the correct number defined in structs.FETCH.")
	assert.NoError(t, err, "should run without error")

	r.setData(structs.SUBMIT.String())
	command, _ = readCommand()

	assert.Equal(t, structs.SUBMIT, command, "does not return the correct number defined in structs.SUBMIT.")

	r.setData(structs.EXIT.String())
	command, _ = readCommand()

	assert.Equal(t, structs.EXIT, command, "does not return the correct number defined in structs.EXIT")
}

func TestReadCommandError(t *testing.T) {
	testlogger.BeginTest()
	defer testlogger.EndTest()

	r := newTestReader("-1")
	Reader = r

	_, err := readCommand()

	assert.Error(t, err, "-1 should not be a valid argument.")

	r.setData("abcd")
	_, err = readCommand()

	assert.Error(t, err, "abcd should not be a valid argument.")

	r.setData("")
	_, err = readCommand()

	assert.Error(t, err, "'' should not be a valid argument.")
}

func TestReadCommandWithDelimiter(t *testing.T) {
	testlogger.BeginTest()
	defer testlogger.EndTest()

	r := newTestReader("2\r\n")
	Reader = r

	command, err := readCommand()

	assert.Equal(t, structs.EXIT, command, "command should be EXIT command")
	assert.NoError(t, err, "input is valid so there should be no error")
}

func TestNextCommand(t *testing.T) {
	testlogger.BeginTest()
	defer testlogger.EndTest()

	var inputCommand structs.Command
	var inputError error
	readCom = func() (structs.Command, error) {
		return inputCommand, inputError
	}
	inputCommand = structs.SUBMIT
	inputError = nil
	outputCommand, outputError := NextCommand()
	assert.Equal(t, inputCommand, outputCommand)
	assert.Equal(t, inputError, outputError)
	inputError = errors.New(string(structs.NoValidOption))
	outputCommand, outputError = NextCommand()
	assert.EqualError(t, outputError, string(structs.TooManyInputs))
}

func TestOutputMessageToCommandLine(t *testing.T) {
	testlogger.BeginTest()
	defer testlogger.EndTest()

	var output string
	Writer = newTestWriter(&output)
	testMessage := structs.RequestTicketID
	OutputMessageToCommandLine(testMessage)
	assert.Equal(t, string(testMessage), output, "string output failed.")
}

func TestGetEmailAddress(t *testing.T) {
	testlogger.BeginTest()
	defer testlogger.EndTest()

	verifyEmailAddress = func(emailAddress string) bool {
		return true
	}
	emailAddress := "john.doe@example.com"

	r := newTestReader(emailAddress)
	Reader = r

	outputEmailAddress, err := getEmailAddress()

	assert.True(t, err == nil, "unexpected error with correct email address.")
	assert.Equal(t, emailAddress, outputEmailAddress, "Email address was distorted during reading.")

	r.setData("")
	_, err = getEmailAddress()
	assert.EqualError(t, err, string(structs.EmptyString))

	r.setData("notEmpty")
	verifyEmailAddress = func(emailAddress string) bool {
		return false
	}
	_, err = getEmailAddress()
	assert.EqualError(t, err, string(structs.InvalidEmail))
}

func TestGetString(t *testing.T) {
	testlogger.BeginTest()
	defer testlogger.EndTest()

	testString := "abcd"

	r := newTestReader(testString)
	Reader = r

	outputString, err := getString()

	assert.Equal(t, testString, outputString, "string was distorted during reading.")
	assert.NoError(t, err, "function should not throw an error with a normal string.")

	r.setData("")
	_, err = getString()
	assert.EqualError(t, err, string(structs.EmptyString), "getting empty strings is not very useful.")
}

func TestGetStringWithDelimiter(t *testing.T) {
	testlogger.BeginTest()
	defer testlogger.EndTest()

	input := "email@address.com\r\n"
	Reader = newTestReader(input)
	result, err := getString()
	assert.Equal(t, strings.Trim(input, "\r\n"), result, "string that was read should be equal to input")
	assert.NoError(t, err, "function should not throw an error with a valid input")
}

func TestGetEmail(t *testing.T) {
	testlogger.BeginTest()
	defer testlogger.EndTest()

	var inputStrings []string
	var inputErrors []error
	var index int
	readEmailAddress = func() (emailAddr string, err error) {
		emailAddr, err = inputStrings[index], inputErrors[index]
		index++
		return
	}
	readString = func() (result string, err error) {
		result, err = inputStrings[index], inputErrors[index]
		index++
		return
	}

	// First test case: every field is filled, all is well
	t.Run("everyFieldFilled", func(t *testing.T) {
		inputStrings = []string{"emailAddress", "ticketID", "subject", "and a message"}
		inputErrors = make([]error, 4)

		outputJSON, outputError := GetEmail()

		assert.NoError(t, outputError)

		// expectedMail := structs.Mail{
		// 	Email:   "emailAddress",
		// 	Subject: `[Ticket "ticketID"] subject`,
		// 	Message: "and a message",
		// }
		expectedJSON := `{"from":"emailAddress", "subject":"[Ticket \"ticketID\"] subject", "message":"and a message"}`
		assert.Equal(t, expectedJSON, outputJSON)
	})

	// Second test case: the ticketID is empty. This is still allowed.
	t.Run("everyFieldExceptTicketIDFilled", func(t *testing.T) {
		inputStrings = []string{"emailAddress", "", "subject", "and a message"}
		inputErrors[1] = errors.New(string(structs.EmptyString))
		index = 0

		outputJSON, outputError := GetEmail()

		assert.NoError(t, outputError)

		expectedJSON := `{"from":"emailAddress", "subject":"subject", "message":"and a message"}`
		assert.Equal(t, expectedJSON, outputJSON)
	})

	// Third test case: an unexpected error while reading the ticket
	// id occurs. Throws an error after continuously receiving invalid
	// user input.
	t.Run("ticketIdError", func(t *testing.T) {
		index = 0

		readString = func() (result string, err error) {
			return "", errors.New("invalid ticket id")
		}

		readEmailAddress = func() (result string, err error) {
			return "validAddress", nil
		}

		_, outputError := GetEmail()
		assert.EqualError(t, outputError, string(structs.TooManyInputs))
	})

	// Fourth test case: an unexpected empty input. Throws an error
	// after continuously receiving invalid user input.
	t.Run("emptyInput", func(t *testing.T) {
		index = 0

		readString = func() (result string, err error) {
			return "", errors.New(string(structs.EmptyString))
		}

		_, outputError := GetEmail()
		assert.EqualError(t, outputError, string(structs.TooManyInputs))
	})

	// Fifth test case: the email is invalid. Throws an error
	// after continuously receiving invalid user input.
	t.Run("invalidEmail", func(t *testing.T) {
		index = 0

		readEmailAddress = func() (result string, err error) {
			return "noValidAddress", errors.New(string(structs.InvalidEmail))
		}

		_, outputError := GetEmail()
		assert.EqualError(t, outputError, string(structs.TooManyInputs))
	})

	// Sixth test case: the message is empty. Throws an error
	// after continuously receiving invalid user input.
	t.Run("emptyMessage", func(t *testing.T) {
		index = 0

		inputStrings = []string{"emailAddress", "ticketID", "subject", ""}
		inputErrors = []error{nil, nil, errors.New(string(structs.EmptyString))}

		readString = func() (result string, err error) {
			result, err = inputStrings[index], inputErrors[index]
			if index < 2 {
				index++
			}
			return
		}

		readEmailAddress = func() (result string, err error) {
			return "validAddress", nil
		}

		_, outputError := GetEmail()
		assert.EqualError(t, outputError, string(structs.TooManyInputs))
	})
}

func TestPrintEmail(t *testing.T) {
	testlogger.BeginTest()
	defer testlogger.EndTest()

	var output string
	Writer = newTestWriter(&output)
	emailAddress := "emailAddress"
	fromAddress := "anotherEmailAddress"
	subject := "subjectLine"
	message := "message"
	mail := structs.Mail{
		ID:      "AnID",
		From:    fromAddress,
		To:      emailAddress,
		Subject: subject,
		Message: message,
	}

	PrintEmail(mail)

	expectedOutput := "From: " + fromAddress + "\n" + string(structs.To) + emailAddress + "\n\n" +
		string(structs.Subject) + subject + "\n\n" +
		message + "\n"
	assert.Equal(t, expectedOutput, output)
}
