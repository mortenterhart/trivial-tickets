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
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mortenterhart/trivial-tickets/log/testlog"
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
 * Package cliutils [tests]
 * Various utilities for CLI
 */

func TestCreateSubjectLine(t *testing.T) {
	testlog.BeginTest()
	defer testlog.EndTest()

	subjectLine := createSubjectLine("abcd", "")
	assert.Equal(t, "abcd", subjectLine)

	subjectLine = createSubjectLine("abcd", "3FM")
	assert.Equal(t, "[Ticket \\\"3FM\\\"] abcd", subjectLine)
}

func TestCreateMail(t *testing.T) {
	testlog.BeginTest()
	defer testlog.EndTest()

	eMailAddr := "john.doe@example.com"
	subj := "Search field broken"
	tID := "12ab3"
	mes := "a message"
	expected := `{"from":"john.doe@example.com", "subject":"[Ticket \"12ab3\"] Search field broken", "message":"a message"}`

	actual := CreateMail(eMailAddr, subj, tID, mes)

	assert.Equal(t, expected, actual)
}

func TestCheckEmailAddress(t *testing.T) {
	testlog.BeginTest()
	defer testlog.EndTest()

	simple := "aName@address.com"
	withDots := "first.last@more.complicated.address.org"
	dashesAndUnderscores := "mike_miller@impressive-institute.de"
	failWithSpaces := "notAn Address@example.com"
	failAt := "notAnAddress.com"
	failDomain := "first.last@notADomain"

	assert.True(t, CheckEmailAddress(simple))
	assert.True(t, CheckEmailAddress(withDots))
	assert.True(t, CheckEmailAddress(dashesAndUnderscores))
	assert.False(t, CheckEmailAddress(failWithSpaces))
	assert.False(t, CheckEmailAddress(failAt))
	assert.False(t, CheckEmailAddress(failDomain))
}
