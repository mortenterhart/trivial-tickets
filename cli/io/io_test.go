// I/O operations for the CLI
package io

import (
	"errors"
	"github.com/mortenterhart/trivial-tickets/structs"
	"github.com/stretchr/testify/assert"
	"io"
	"strconv"
	"testing"
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

type TestWriter struct {
	output *string
}

func NewTestWriter(out *string) (w *TestWriter) {
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

func NewTestReader(data string) (r *TestReader) {
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
	reader = NewTestReader(strconv.Itoa(int(structs.FETCH)))
	command, err := readCommand()
	assert.Equal(t, structs.FETCH, command, "does not return the correct number defined in structs.FETCH.")
	assert.Equal(t, nil, err, "should run without error")
	r, ok := reader.(ITestReader)
	if ok {
		r.setData(strconv.Itoa(int(structs.SUBMIT)))
		command, _ = readCommand()
		assert.Equal(t, structs.SUBMIT, command, "does not return the correct number defined in structs.SUBMIT.")
		r.setData(strconv.Itoa(int(structs.EXIT)))
		command, _ = readCommand()
		assert.Equal(t, structs.EXIT, command, "does not return the correct number defined in structs.EXIT.")
	}
}

func TestReadCommandError(t *testing.T) {
	reader = ITestReader(NewTestReader("-1"))
	_, err := readCommand()
	assert.Error(t, err, "-1 should not be a valid argument.")
	r, ok := reader.(ITestReader)
	if ok {
		r.setData("abcd")
		_, err = readCommand()
		assert.Error(t, err, "abcd should not be a valid argument.")
		r.setData("")
		_, err = readCommand()
		assert.Error(t, err, "'' should not be a valid argument.")
	}
}

func TestNextCommand(t *testing.T) {
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
	assert.EqualError(t, outputError, string(structs.AbortExecutionDueToManyWrongUserInputs))
}

func TestOutputMessageToCommandLine(t *testing.T) {
	var output string
	writer = NewTestWriter(&output)
	testMessage := structs.RequestTicketID
	OutputMessageToCommandLine(testMessage)
	assert.Equal(t, string(testMessage), output, "string output failed.")
}

func TestGetEmailAddress(t *testing.T) {
	addressIsCorrect := true
	verifyEmailAddress = func(emailAddress string) bool {
		return addressIsCorrect
	}
	emailAddress := "john.doe@example.com"
	reader = NewTestReader(emailAddress)
	outputEmailAddress, err := getEmailAddress()
	assert.True(t, err == nil, "unexpected error with correct email address.")
	assert.Equal(t, emailAddress, outputEmailAddress, "Eamil address was distorted during reading.")
	emailAddress = ""
	if r, ok := reader.(ITestReader); ok {
		r.setData(emailAddress)
		_, err = getEmailAddress()
		assert.EqualError(t, err, string(structs.EmptyString))
		r.setData("notEmpty")
		addressIsCorrect = false
		_, err = getEmailAddress()
		assert.EqualError(t, err, string(structs.InvalidEmail))
	}
}

func TestGetString(t *testing.T) {
	testString := "abcd"
	reader = NewTestReader(testString)
	outputString, err := getString()
	assert.Equal(t, testString, outputString, "string was distorted during reading.")
	assert.Equal(t, nil, err, "function should not throw an error with a normal string.")
	r, ok := reader.(ITestReader)
	if ok {
		r.setData("")
		_, err = getString()
		assert.EqualError(t, err, string(structs.EmptyString), "getting empty strings is not very useful.")
	}
}

func TestGetEmail(t *testing.T) {
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

	//First test case: every field is filled, all is well

	inputStrings = make([]string, 0)
	inputStrings = append(inputStrings, "emailAddress", "ticketID", "subject", "and a message")
	inputErrors = make([]error, 4)
	outputJson, outputError := GetEmail()
	assert.NoError(t, outputError)
	//expectedMail := structs.Mail{
	//	Email:   "emailAddress",
	//	Subject: `[Ticket "ticketID"] subject`,
	//	Message: "and a message"}
	expectedJson := `{"from":"emailAddress", "subject":"[Ticket \"ticketID\"] subject", "message":"and a message"}`
	assert.Equal(t, expectedJson, outputJson)

	//Second test case: the ticketID is empty. This is still allowed.

	inputStrings = make([]string, 0)
	inputStrings = append(inputStrings, "emailAddress", "", "subject", "and a message")
	inputErrors[1] = errors.New(string(structs.EmptyString))
	index = 0
	outputJson, outputError = GetEmail()
	assert.NoError(t, outputError)
	expectedJson = `{"from":"emailAddress", "subject":"subject", "message":"and a message"}`
	assert.Equal(t, expectedJson, outputJson)

	//Third test case: an unexpected empty input. Throws an error after continuously receiving invalid user input.

	index = 0
	readString = func() (result string, err error) {
		return "", errors.New(string(structs.EmptyString))
	}
	_, outputError = GetEmail()
	assert.EqualError(t, outputError, string(structs.AbortExecutionDueToManyWrongUserInputs))

	//Fourth test case: the email is invalid. Throws an error after continuously receiving invalid user input.

	index = 0
	readEmailAddress = func() (result string, err error) {
		return "noValidAddress", errors.New(string(structs.InvalidEmail))
	}
	_, outputError = GetEmail()
	assert.EqualError(t, outputError, string(structs.AbortExecutionDueToManyWrongUserInputs))
}

func TestPrintEmail(t *testing.T) {
	var output string
	writer = NewTestWriter(&output)
	emailAddress := "emailAddress"
	fromAddress := "anotherEmailAddress"
	subject := "subjectline"
	message := "message"
	mail := structs.Mail{
		Id:      "AnID",
		From:    fromAddress,
		To:      emailAddress,
		Subject: subject,
		Message: message}
	PrintEmail(mail)
	expectedOutput := "From: " + fromAddress + "\n" + string(structs.To) + emailAddress + "\n\n" +
		string(structs.Subject) + subject + "\n\n" +
		message + "\n"
	assert.Equal(t, expectedOutput, output)
}
