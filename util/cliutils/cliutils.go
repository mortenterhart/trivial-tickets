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

// Package cliutils contains helper functions and various
// utilities for the CLI.
package cliutils

import (
	"fmt"
	"regexp"
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
 * Package cliutils
 * Various utilities for CLI
 */

// createSubjectLine creates a SubjectLine based on subject
// and ticketID. If ticketID is empty, subjectLine = subject.
// The ticketID in the subjectLine is used by the API to
// assign the message to an already existing ticket.
func createSubjectLine(subject string, ticketID string) (subjectLine string) {
	if ticketID != "" {
		subjectLine = "[Ticket \\\"" + ticketID + "\\\"] "
	}
	subjectLine += subject
	return
}

// CreateMail returns a structs.Mail created with the input
// parameters. It expects the input parameters to be valid,
// no checks are being done on them. Internally it relies on
// the createSubjectLine function.
func CreateMail(emailAddress string, subject string, ticketID string, message string) (mailJSON string) {
	mailJSON = fmt.Sprintf(`{"from":"%s", "subject":"%s", "message":"%s"}`, emailAddress, createSubjectLine(subject, ticketID), message)
	return
}

// CheckEmailAddress returns true if the input string is a
// syntactically correct email address.
func CheckEmailAddress(emailAddress string) bool {
	r, _ := regexp.Compile("^([\\w\\.\\-]+)@([\\w*\\.\\-]+)(\\.\\w+)$")
	return r.MatchString(emailAddress)
}
