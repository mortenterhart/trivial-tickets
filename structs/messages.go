// Project-wide needed structures for data elements
package structs

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
 * Package structs
 * Project-wide needed structures for data elements
 */

type CliMessage string

const (
	RequestCommandInput CliMessage = "\nTo fetch Mails from the server, type '0'\n" +
		"To send an email to the server type '1'\n" +
		"To exit this program type '2'\n"
	CommandNotAccepted  CliMessage = "Input not accepted, error: "
	RequestEmailAddress CliMessage = "Please enter your email address. It has to be valid.\n"
	RequestSubject      CliMessage = "Please enter the subject line\n"
	RequestMessage      CliMessage = "Please enter the body of the message.\n"
	RequestTicketID     CliMessage = "If applicable please enter the ticket ID. If left empty, a new ticket will be created.\n"
	To                  CliMessage = "To: "
	Subject             CliMessage = "Subject: "
)

type CliErrMessage string

const (
	AbortExecutionDueToManyWrongUserInputs CliErrMessage = "Too many successive wrong user inputs. Aborting program execution.\n"
	NoValidOption                          CliErrMessage = "Not within the range of valid options"
	EmptyString                            CliErrMessage = "string is empty"
	InvalidEmail                           CliErrMessage = "not a valid email address"
)
