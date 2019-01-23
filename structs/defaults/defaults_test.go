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

// Package defaults defines default constraints for
// the server and the test suite.
package defaults

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mortenterhart/trivial-tickets/logger/testlogger"
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
 * Package defaults [tests]
 * Default constraints for the server and the tests
 */

func TestDefaultValues(t *testing.T) {
	testlogger.BeginTest()
	defer testlogger.EndTest()

	assert.NotNil(t, ServerPort)
	assert.NotNil(t, ServerTickets)
	assert.NotNil(t, ServerUsers)
	assert.NotNil(t, ServerMails)
	assert.NotNil(t, ServerCertificate)
	assert.NotNil(t, ServerKey)
	assert.NotNil(t, ServerWeb)

	assert.NotNil(t, TestTicketsTrimmed)
	assert.NotNil(t, TestUsersTrimmed)
	assert.NotNil(t, TestMailsTrimmed)
	assert.NotNil(t, TestCertificateTrimmed)
	assert.NotNil(t, TestKeyTrimmed)
	assert.NotNil(t, TestWebTrimmed)

	assert.NotNil(t, LogVerbose)
	assert.NotNil(t, LogFullPaths)
	assert.NotNil(t, LogLevelString)

	assert.NotNil(t, CliHost)
	assert.NotNil(t, CliPort)
	assert.NotNil(t, CliCertificate)
	assert.NotNil(t, CliFetch)
	assert.NotNil(t, CliSubmit)

	assert.NotNil(t, TestPort)
	assert.NotNil(t, TestTickets)
	assert.NotNil(t, TestUsers)
	assert.NotNil(t, TestMails)
	assert.NotNil(t, TestCertificate)
	assert.NotNil(t, TestKey)
	assert.NotNil(t, TestWeb)

	assert.NotNil(t, FileModeRegular)

	assert.NotNil(t, ExitSuccessful)
	assert.NotNil(t, ExitStartError)
	assert.NotNil(t, ExitShutdownError)
}
