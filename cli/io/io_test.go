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
	"bytes"
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mortenterhart/trivial-tickets/log/testlog"
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

// TestWriter implements the io.Writer interface
// and can be therefore used as replacement for
// a genuine writer. The output can be tested
// against an expectation by using the contents
// of the output buffer with the pointer, for
// example in this way:
//
//   var out bytes.Buffer
//   w := newTestWriter(&output)
//
//   // Call function to test or assign test writer
//   // to a writer used in productive code.
//
//   if out.String() != "Expectation" {
//     t.Fail()
//   }
type TestWriter struct {
	// out is the buffer to which all the output
	// of this writer is written to.
	out *bytes.Buffer
}

// newTestWriter creates a new test writer with the
// specified output buffer which outputs can be
// retrieved using the `out.String()` method.
func newTestWriter(out *bytes.Buffer) (w *TestWriter) {
	return &TestWriter{out}
}

// Write writes the provided bytes to the output
// buffer and returns the number of bytes written
// and a non-nil error if an error occurred.
func (w *TestWriter) Write(p []byte) (n int, err error) {
	return w.out.Write(p)
}

// ITestReader is an interface that extends the
// `io.Reader` interface and provides the facility
// to set a custom input for tests.
type ITestReader interface {
	// Embed the `io.Reader` interface to allow
	// using this test reader as normal reader.
	io.Reader

	// setInput replaces the current input buffer
	// with the specified input string.
	setInput(input string)
}

// TestReader is a test reader which can manipulate
// the bytes read into the input buffer. The reader
// can be used to emulate user inputs from the
// command line.
type TestReader struct {
	// input is the modifiable input buffer.
	input []byte
}

// newTestReader creates a new test reader with
// the specified input string as new buffer.
func newTestReader(input string) (r *TestReader) {
	return &TestReader{[]byte(input)}
}

// setInput replaces the input buffer with the
// specified string.
func (r *TestReader) setInput(input string) {
	r.input = []byte(input)
}

// readByte reads the first byte from the input
// buffer and returns it. The read byte is then
// popped out from the buffer. If there are no
// bytes remaining to be read the function panics,
// so make sure to call `eof()` before.
func (r *TestReader) readByte() byte {
	// this function assumes that eof() check was done before
	b := r.input[0]
	r.input = r.input[1:]
	return b
}

// eof checks if the reader has reached the end
// of the input buffer.
func (r *TestReader) eof() (eof bool) {
	return len(r.input) == 0
}

// Read reads up to len(p) bytes into the slice p.
// It returns the number of bytes read and an error.
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
	testlog.BeginTest()
	defer testlog.EndTest()

	r := newTestReader(structs.CommandFetch.String())
	Reader = r

	command, err := readCommand()

	assert.Equal(t, structs.CommandFetch, command, "does not return the correct number defined in structs.CommandFetch.")
	assert.NoError(t, err, "should run without error")

	r.setInput(structs.CommandSubmit.String())
	command, _ = readCommand()

	assert.Equal(t, structs.CommandSubmit, command, "does not return the correct number defined in structs.CommandSubmit.")

	r.setInput(structs.CommandExit.String())
	command, _ = readCommand()

	assert.Equal(t, structs.CommandExit, command, "does not return the correct number defined in structs.CommandExit")
}

func TestReadCommandError(t *testing.T) {
	testlog.BeginTest()
	defer testlog.EndTest()

	r := newTestReader("-1")
	Reader = r

	_, err := readCommand()

	assert.Error(t, err, "-1 should not be a valid argument.")

	r.setInput("abcd")
	_, err = readCommand()

	assert.Error(t, err, "abcd should not be a valid argument.")

	r.setInput("")
	_, err = readCommand()

	assert.Error(t, err, "'' should not be a valid argument.")
}

func TestReadCommandWithDelimiter(t *testing.T) {
	testlog.BeginTest()
	defer testlog.EndTest()

	r := newTestReader("2\r\n")
	Reader = r

	command, err := readCommand()

	assert.Equal(t, structs.CommandExit, command, "command should be CommandExit command")
	assert.NoError(t, err, "input is valid so there should be no error")
}

func TestNextCommand(t *testing.T) {
	testlog.BeginTest()
	defer testlog.EndTest()

	var inputCommand structs.Command
	var inputError error
	readCom = func() (structs.Command, error) {
		return inputCommand, inputError
	}

	inputCommand = structs.CommandSubmit
	inputError = nil

	outputCommand, outputError := NextCommand()

	assert.Equal(t, inputCommand, outputCommand)
	assert.Equal(t, inputError, outputError)

	inputError = errors.New(string(structs.NoValidOption))
	outputCommand, outputError = NextCommand()

	assert.EqualError(t, outputError, string(structs.TooManyInputs))
}

func TestOutputMessageToCommandLine(t *testing.T) {
	testlog.BeginTest()
	defer testlog.EndTest()

	var output bytes.Buffer
	Writer = newTestWriter(&output)

	testMessage := structs.RequestTicketID
	OutputMessageToCommandLine(testMessage)

	assert.Equal(t, string(testMessage), output.String(), "string output failed.")
}

func TestGetEmailAddress(t *testing.T) {
	testlog.BeginTest()
	defer testlog.EndTest()

	verifyEmailAddress = func(emailAddress string) bool {
		return true
	}
	emailAddress := "john.doe@example.com"

	r := newTestReader(emailAddress)
	Reader = r

	outputEmailAddress, err := getEmailAddress()

	assert.True(t, err == nil, "unexpected error with correct email address.")
	assert.Equal(t, emailAddress, outputEmailAddress, "Email address was distorted during reading.")

	r.setInput("")
	_, err = getEmailAddress()
	assert.EqualError(t, err, string(structs.EmptyString))

	r.setInput("notEmpty")
	verifyEmailAddress = func(emailAddress string) bool {
		return false
	}
	_, err = getEmailAddress()
	assert.EqualError(t, err, string(structs.InvalidEmail))
}

func TestGetString(t *testing.T) {
	testlog.BeginTest()
	defer testlog.EndTest()

	testString := "abcd"

	r := newTestReader(testString)
	Reader = r

	outputString, err := getString()

	assert.Equal(t, testString, outputString, "string was distorted during reading.")
	assert.NoError(t, err, "function should not throw an error with a normal string.")

	r.setInput("")
	_, err = getString()
	assert.EqualError(t, err, string(structs.EmptyString), "getting empty strings is not very useful.")
}

func TestGetStringWithDelimiter(t *testing.T) {
	testlog.BeginTest()
	defer testlog.EndTest()

	input := "email@address.com\r\n"
	Reader = newTestReader(input)
	result, err := getString()
	assert.Equal(t, strings.Trim(input, "\r\n"), result, "string that was read should be equal to input")
	assert.NoError(t, err, "function should not throw an error with a valid input")
}

func TestGetEmail(t *testing.T) {
	testlog.BeginTest()
	defer testlog.EndTest()

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
	testlog.BeginTest()
	defer testlog.EndTest()

	var output bytes.Buffer
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
	assert.Equal(t, expectedOutput, output.String())
}
