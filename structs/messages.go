package structs

import (
	"fmt"
)

type Message string

const (
	REQUEST_COMMAND_INPUT Message = "To fetch Mails from the server, type '0'\n" +
		"To send an email to the server type '1'\n" +
		"To exit this program type '2'\n"
	COMMAND_NOT_ACCEPTED  = "Input not accepted, error: "
	REQUEST_EMAIL_ADDRESS = "Please enter your email address. It has to be valid.\n"
	REQUEST_SUBJECT       = "Please enter the subject line\n"
	REQUEST_MESSAGE       = "Please enter the body of the message.\n"
)

func test() {
	s := "aString"
	s = s + "to try something"
	fmt.Printf("%s", s)
}
