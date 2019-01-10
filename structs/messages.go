package structs

type Message string

const (
	RequestCommandInput Message = "To fetch Mails from the server, type '0'\n" +
		"To send an email to the server type '1'\n" +
		"To exit this program type '2'\n"
	CommandNotAccepted  = "Input not accepted, error: "
	RequestEmailAddress = "Please enter your email address. It has to be valid.\n"
	RequestSubject      = "Please enter the subject line\n"
	RequestMessage      = "Please enter the body of the message.\n"
)
