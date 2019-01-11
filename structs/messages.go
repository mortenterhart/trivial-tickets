package structs

type Message string

const (
	RequestCommandInput Message = "\nTo fetch Mails from the server, type '0'\n" +
		"To send an email to the server type '1'\n" +
		"To exit this program type '2'\n"
	CommandNotAccepted  Message = "Input not accepted, error: "
	RequestEmailAddress Message = "Please enter your email address. It has to be valid.\n"
	RequestSubject      Message = "Please enter the subject line\n"
	RequestMessage      Message = "Please enter the body of the message.\n"
	RequestTicketID     Message = "If applicable please enter the ticket ID. If left empty, a new ticket will be created.\n"
	Receiver            Message = "Receiver: "
	Subject             Message = "Subject: "
)

type ErrMessage string

const (
	AbortExecutionDueToManyWrongUserInputs ErrMessage = "Too many successive wrong user inputs. Aborting program execution.\n"
	NoValidOption                          ErrMessage = "Not within the range of valid options"
	EmptyString                            ErrMessage = "string is empty"
	InvalidEmail                           ErrMessage = "not a valid email address"
)
