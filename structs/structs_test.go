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

// Package structs supplies project-wide needed data
// structures, types and constants for the server and
// the command-line tool.
package structs

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
 * Package structs [tests]
 * Project-wide needed structures for data elements
 */

func TestLogLevel_String(t *testing.T) {
	testlog.BeginTest()
	defer testlog.EndTest()

	t.Run("infoString", func(t *testing.T) {
		assert.Equal(t, "[INFO]", LevelInfo.String())
	})

	t.Run("warningString", func(t *testing.T) {
		assert.Equal(t, "[WARNING]", LevelWarning.String())
	})

	t.Run("errorString", func(t *testing.T) {
		assert.Equal(t, "[ERROR]", LevelError.String())
	})

	t.Run("fatalErrorString", func(t *testing.T) {
		assert.Equal(t, "[FATAL ERROR]", LevelFatal.String())
	})

	t.Run("testDebugString", func(t *testing.T) {
		assert.Equal(t, "[TEST DEBUG]", LevelTestDebug.String())
	})

	t.Run("undefinedString", func(t *testing.T) {
		assert.Equal(t, "undefined", LogLevel(7).String())
	})
}

func TestAsLogLevel(t *testing.T) {
	t.Run("infoString", func(t *testing.T) {
		assert.Equal(t, LevelInfo, AsLogLevel("info"))
	})

	t.Run("warningString", func(t *testing.T) {
		assert.Equal(t, LevelWarning, AsLogLevel("warning"))
	})

	t.Run("errorString", func(t *testing.T) {
		assert.Equal(t, LevelError, AsLogLevel("error"))
	})

	t.Run("fatalString", func(t *testing.T) {
		assert.Equal(t, LevelFatal, AsLogLevel("fatal"))
	})

	t.Run("undefinedString", func(t *testing.T) {
		assert.Equal(t, LogLevel(-1), AsLogLevel("undefined"))
	})
}

func TestStatus_String(t *testing.T) {
	testlog.BeginTest()
	defer testlog.EndTest()

	t.Run("openString", func(t *testing.T) {
		assert.Equal(t, "Open", StatusOpen.String())
	})

	t.Run("processingString", func(t *testing.T) {
		assert.Equal(t, "In Progress", StatusInProgress.String())
	})

	t.Run("closedString", func(t *testing.T) {
		assert.Equal(t, "Closed", StatusClosed.String())
	})

	t.Run("undefinedStatusString", func(t *testing.T) {
		assert.Equal(t, "undefined status", Status(5).String())
	})
}

func TestCommand_String(t *testing.T) {
	testlog.BeginTest()
	defer testlog.EndTest()

	t.Run("fetchNumber", func(t *testing.T) {
		assert.Equal(t, "0", CommandFetch.String())
	})

	t.Run("submitNumber", func(t *testing.T) {
		assert.Equal(t, "1", CommandSubmit.String())
	})

	t.Run("exitNumber", func(t *testing.T) {
		assert.Equal(t, "2", CommandExit.String())
	})
}
