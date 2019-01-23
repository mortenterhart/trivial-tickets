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

// Package globals contains global settings and resources for
// the server and logger.
package globals

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
 * Package globals [tests]
 * Global hash maps and structs
 */

// TestGlobalStructures checks that all variables in the
// global package are correctly initialized.
func TestGlobalStructures(t *testing.T) {
	testlog.BeginTest()
	defer testlog.EndTest()

	t.Run("ticketsNotNil", func(t *testing.T) {
		assert.NotNil(t, Tickets, "Tickets should not be nil")
	})

	t.Run("mailsNotNil", func(t *testing.T) {
		assert.NotNil(t, Mails, "Mails should not be nil")
	})

	t.Run("emptyServerConfig", func(t *testing.T) {
		assert.Empty(t, ServerConfig, "ServerConfig should be an empty struct")
	})

	t.Run("emptyLogConfig", func(t *testing.T) {
		assert.Empty(t, LogConfig, "LogConfig should be an empty struct")
	})

	t.Run("sessionsNotNil", func(t *testing.T) {
		assert.NotNil(t, Sessions, "Sessions should not be nil")
	})
}
