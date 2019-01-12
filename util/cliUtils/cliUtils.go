package cliUtils

import (
	"regexp"
)

// createSubjectLine creates a SubjectLine based on subject and ticketID. If ticketID is empty, subjectLine = subject.
// The ticketID in the subjectLine is used by the API to assign the message to an already existing ticket.
func createSubjectLine(subject string, ticketID string) (subjectLine string) {
	if ticketID != "" {
		subjectLine = "[Ticket \\\"" + ticketID + "\\\"] "
	}
	subjectLine += subject
	return
}

// CreateMail returns a structs.Mail created with the input parameters.
// It expects the input parameters to be valid, no checks are being done on them.
// Internally it relies on the createSubjectLine function.
func CreateMail(eMailAddress string, subject string, ticketID string, message string) (mailJson string) {
	mailJson = `{"from":"` + eMailAddress + `", "subject":"` + createSubjectLine(subject, ticketID) + `", "message":"` + message + `"}`
	return
}

// CheckEmailAddress returns true if the input string is a syntactically correct email address.
func CheckEmailAddress(emailAddress string) bool {
	r, _ := regexp.Compile("^([\\w\\.\\-]+)@([\\w*\\.\\-]+)(\\.\\w+)$")
	return r.MatchString(emailAddress)
}
