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

// Package server implements the web server including
// shutdown routines and the associated handlers for
// web requests.
package server

import (
	"net/http"
	"net/http/httptest"
	"os"
	"syscall"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"

	"github.com/mortenterhart/trivial-tickets/globals"
	"github.com/mortenterhart/trivial-tickets/log/testlog"
	"github.com/mortenterhart/trivial-tickets/structs"
	"github.com/mortenterhart/trivial-tickets/structs/defaults"
	"github.com/mortenterhart/trivial-tickets/util/filehandler"
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
 * Package server [tests]
 * Server starting and handler registration
 */

// getRegisteredHandler lookups a registered handler inside the
// HTTP serving multiplexer recognizing it by its URL and returns
// the handler and its corresponding pattern. The tests assume
// that the returned pattern is equal to the URL.
func getRegisteredHandler(mux *http.ServeMux, url string) (h http.Handler, pattern string) {
	request, errRequest := http.NewRequest("GET", url, nil)
	if errRequest != nil {
		testlog.Debugf("error creating request for %s: %v", url, errRequest)
		return
	}

	return mux.Handler(request)
}

// mockConfig creates a mock configuration
// for the server.
func mockConfig() structs.ServerConfig {
	return structs.ServerConfig{
		Port:    defaults.TestPort,
		Tickets: defaults.ServerTicketsTrimmed,
		Users:   defaults.ServerUsersTrimmed,
		Mails:   defaults.ServerMailsTrimmed,
		Cert:    defaults.TestCertificateTrimmed,
		Key:     defaults.TestKeyTrimmed,
		Web:     defaults.TestWebTrimmed,
	}
}

// mockLogConfig creates a mock configuration
// for the logger.
func mockLogConfig() structs.LogConfig {
	return structs.LogConfig{
		LogLevel:  structs.AsLogLevel(defaults.LogLevelString),
		Verbose:   defaults.LogVerbose,
		FullPaths: defaults.LogFullPaths,
	}
}

// TestGetTemplates makes sure the application is
// able to correctly find the templates with the
// given standard values.
func TestGetTemplates(t *testing.T) {
	testlog.BeginTest()
	defer testlog.EndTest()

	tmpl := getTemplates("../www")
	tmplNil := getTemplates("/www")

	assert.NotNil(t, tmpl, "getTemplates() returned no found templates")
	assert.Nil(t, tmplNil, "getTemplates() found templates where it was not supposed to be")
}

// TestRedirectToTLS tests the redirect to https,
// if a request with only http is made.
func TestRedirectToTLS(t *testing.T) {
	testlog.BeginTest()
	defer testlog.EndTest()

	config := mockConfig()
	globals.ServerConfig = &config

	req, _ := http.NewRequest("GET", "localhost", nil)
	w := httptest.NewRecorder()
	handler := http.HandlerFunc(redirectToTLS)

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusTemporaryRedirect, w.Result().StatusCode, "The HTTP status code was incorrect")
}

// TestRedirectToTLSWithParams tests the redirect
// to https, if a request with parameters and
// only http is made.
func TestRedirectToTLSWithParams(t *testing.T) {
	testlog.BeginTest()
	defer testlog.EndTest()

	req, _ := http.NewRequest("GET", "localhost?id=123", nil)
	w := httptest.NewRecorder()
	handler := http.HandlerFunc(redirectToTLS)

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusTemporaryRedirect, w.Result().StatusCode, "The HTTP status code was incorrect")
}

// TestRegisterHandlers ensures that the handlers
// will be registered properly if a valid path to
// the web server's root directory is given.
func TestRegisterHandlers(t *testing.T) {
	testlog.BeginTest()
	defer testlog.EndTest()

	mux, err := registerHandlers("../www")

	assert.NoError(t, err, "Registering handlers with valid path should not error")
	assert.NotNil(t, mux, "returned mux should not be nil")

	testHandlerRegistered(t, mux, "/")
	testHandlerRegistered(t, mux, "/login")
	testHandlerRegistered(t, mux, "/logout")
	testHandlerRegistered(t, mux, "/createTicket")
	testHandlerRegistered(t, mux, "/holiday")
	testHandlerRegistered(t, mux, "/ticket")
	testHandlerRegistered(t, mux, "/updateTicket")
	testHandlerRegistered(t, mux, "/unassignTicket")
	testHandlerRegistered(t, mux, "/assignTicket")
	testHandlerRegistered(t, mux, "/api/receive")
	testHandlerRegistered(t, mux, "/api/fetchMails")
	testHandlerRegistered(t, mux, "/api/verifyMail")
}

// testHandlerRegistered checks that the handler following
// the given URL is registered in the given serving
// multiplexer and verifies that the URL equals the actual
// handler pattern.
func testHandlerRegistered(t *testing.T, mux *http.ServeMux, url string) bool {
	handler, pattern := getRegisteredHandler(mux, url)

	assert.NotNil(t, handler, "returned handler for URL should not be nil")
	return assert.Equal(t, url, pattern, "URL and pattern should match because there should be a "+
		"registered handler for the URL")
}

// TestRegisterHandlersNoPath is used to produce an
// error to make sure the function works properly.
func TestRegisterHandlersNoPath(t *testing.T) {
	testlog.BeginTest()
	defer testlog.EndTest()

	handler, err := registerHandlers("")

	assert.NotNil(t, err, "No error occurred, although the path was incorrect")
	assert.Nil(t, handler, "returned handler should be nil")
}

// TestRegisterServerNoUsersPath produces an error to
// make sure the server will not start without a
// path to the users.json file.
func TestStartServerNoUsersPath(t *testing.T) {
	testlog.BeginTest()
	defer testlog.EndTest()

	config := mockConfig()
	config.Users = ""

	exitCode, err := StartServer(&config)

	assert.NotNil(t, err, "No error was returned, although no users path was specified")
	assert.Equal(t, defaults.ExitStartError, exitCode, "exit code should be 1 due to expected error")
}

// TestStartServerNoMailsPath produces an error to make
// sure the server will not start without a path to the
// mail folder.
func TestStartServerNoMailsPath(t *testing.T) {
	testlog.BeginTest()
	defer testlog.EndTest()

	config := mockConfig()
	config.Mails = ""

	exitCode, err := StartServer(&config)

	assert.NotNil(t, err, "No error was returned, although no mails path was specified")
	assert.Equal(t, defaults.ExitStartError, exitCode, "exit code should be 1 due to expected error")
}

// TestStartServerNoWebPath produces an error to make
// sure the server will not start without a path to the
// web files.
func TestStartServerNoWebPath(t *testing.T) {
	testlog.BeginTest()
	defer testlog.EndTest()

	config := mockConfig()
	config.Web = ""

	exitCode, err := StartServer(&config)

	assert.NotNil(t, err, "No error was returned, although no web path was specified")
	assert.Equal(t, defaults.ExitStartError, exitCode, "exit code should be 1 due to expected error")
}

// TestStartServerNoTicketsPath produces an error to make
// sure the server will not start without a path to the
// ticket folder.
func TestStartServerNoTicketsPath(t *testing.T) {
	testlog.BeginTest()
	defer testlog.EndTest()

	config := mockConfig()
	config.Tickets = ""

	exitCode, err := StartServer(&config)

	assert.NotNil(t, err, "No error was returned, although no tickets path was specified")
	assert.Equal(t, defaults.ExitStartError, exitCode, "exit code should be 1 due to expected error")
}

// TestStartServerAllConfigsSet attempts to start the productive server completely
// and tests whether the server startup works correctly with all configuration options
// set properly. The server is directly stopped by a custom start error provoked by
// an invalid port -10. This is the single solution that is working reliably. It is
// not possible to get the server instance from the StartServer() function into this
// test in order to call
//
//     server.Shutdown(context)
//
// and alternatively the blocking channels in handleServerShutdown() are also local to
// the function. Here it was also tried to send an interrupt (SIGINT), kill (SIGKILL)
// or terminating (SIGTERM) signal to the server in order to close it with the interrupt
// routine, but on the one hand this caused strange blocking of tests or even a kill
// of the `go test` process. IDE test interfaces were also terminated by the signals.
// On the other hand signal handling on Windows is not fully supported (e.g sending
// interrupt signals is unsupported). This caused compile errors because library
// functions for killing a process were not defined or caused the whole process to
// finish.
//
// Therefore it was proposed to simply provoke a start error which causes the server
// to reject startup. This can be done by setting the port either to an invalid value
// such as -10 (always provokes error) or by setting it to a value < 1024 where the
// permission to bind the port is denied if the user does not have root privileges.
// Note that the second variant only causes an error if an user other than root is
// logged in to the system.
func TestStartServerAllConfigsSet(t *testing.T) {
	testlog.BeginTest()
	defer testlog.EndTest()

	config := mockConfig()
	config.Port = defaults.TestPort + 3
	config.Cert = "not/existing/server.cert"

	shutdown := make(chan bool)

	go func() {
		exitCode, err := StartServer(&config)

		t.Run("startErrorNotNil", func(t *testing.T) {
			assert.Error(t, err, "returned error should be not-nil because the server was started with an invalid port")
		})

		t.Run("exitCode", func(t *testing.T) {
			assert.Equal(t, defaults.ExitStartError, exitCode, "exit code should be 1 because the server was not able to startup")
		})

		testlog.Debug("Server shutdown completed: Releasing channel")
		shutdown <- true
		close(shutdown)
	}()

	testlog.Debug("Waiting for server shutdown")
	<-shutdown
	testlog.Debug("Waiting finished: Test completed")
}

func TestStopServerSendingInterruptSignal(t *testing.T) {
	testlog.BeginTest()
	defer testlog.EndTest()

	interrupt = make(chan os.Signal)

	go ShutdownServer()

	t.Run("interruptChannelNotNil", func(t *testing.T) {
		assert.NotNil(t, interrupt, "interrupt channel should not be nil")
	})

	receivedSignal := <-interrupt

	t.Run("caughtSignalInterrupt", func(t *testing.T) {
		assert.Equal(t, os.Interrupt, receivedSignal, "received signal should be interrupt signal")
	})
}

func TestStopServerWithRunningServer(t *testing.T) {
	testlog.BeginTest()
	defer testlog.EndTest()

	interrupt = make(chan os.Signal)
	shutdown := make(chan bool)

	config := mockConfig()
	config.Port = defaults.TestPort + 3

	go func() {
		exitCode, startError := StartServer(&config)

		t.Run("successfulExit", func(t *testing.T) {
			assert.Equal(t, defaults.ExitSuccessful, exitCode, "the exit status of the server should be 0 because it was shutdown correctly")
		})

		t.Run("startErrorNil", func(t *testing.T) {
			assert.NoError(t, startError, "there should be no server start error")
		})

		testlog.Debug("Server shutdown completed: Releasing channel")
		shutdown <- true
		close(shutdown)
	}()

	time.Sleep(500 * time.Millisecond)

	ShutdownServer()

	testlog.Debug("Waiting for server shutdown")
	<-shutdown
	testlog.Debug("Waiting finished: Test completed")
}

func TestCreateResourceFolders(t *testing.T) {
	testlog.BeginTest()
	defer testlog.EndTest()

	config := testServerConfig()
	defer cleanupTestFiles(config)

	t.Run("notExistingDirectories", func(t *testing.T) {
		assert.False(t, filehandler.DirectoryExists(config.Tickets), "testtickets directory should not exist yet")
		assert.False(t, filehandler.DirectoryExists(config.Mails), "testmails directory should not exist yet")

		createErr := createResourceFolders(&config)

		assert.NoError(t, createErr, "creating directories should not return error")
		assert.True(t, filehandler.DirectoryExists(config.Tickets), "testtickets directory should exist now")
		assert.True(t, filehandler.DirectoryExists(config.Mails), "testmails directory should exist now")
	})

	t.Run("existingDirectories", func(t *testing.T) {
		createErr := createResourceFolders(&config)

		assert.NoError(t, createErr, "no error because directories already exist")
		assert.True(t, filehandler.DirectoryExists(config.Tickets), "testtickets directory should already exist")
		assert.True(t, filehandler.DirectoryExists(config.Mails), "testmails directory should already exist")
	})

	t.Run("createTicketsError", func(t *testing.T) {
		errorConfig := mockConfig()
		errorConfig.Tickets = ""

		createErr := createResourceFolders(&errorConfig)

		assert.Error(t, createErr, "error because ticket directory with empty name cannot be created")
	})

	t.Run("createMailsError", func(t *testing.T) {
		errorConfig := mockConfig()
		errorConfig.Mails = ""

		createErr := createResourceFolders(&errorConfig)

		assert.Error(t, createErr, "error because mail directory with empty name cannot be created")
	})
}

func TestNotifyOnInterruptSignal(t *testing.T) {
	testlog.BeginTest()
	defer testlog.EndTest()

	interrupt := notifyOnInterruptSignal()

	done := make(chan bool)

	go t.Run("catchSignal", func(t *testing.T) {
		testlog.Debug("Waiting for signal")
		capturedSignal := <-interrupt

		testlog.Debug("Caught signal SIGINT in Go routine")

		t.Run("caughtSignalNotNil", func(t *testing.T) {
			assert.NotNil(t, capturedSignal, "captured signal should not be nil")
		})

		t.Run("isInterruptSignal", func(t *testing.T) {
			assert.Equal(t, os.Interrupt, capturedSignal, "captured signal should be SIGINT signal")
		})

		testlog.Debug("Tests done: Releasing channel")
		done <- true
		close(done)
	})

	interrupt <- os.Interrupt

	testlog.Debug("Waiting for signal tests to finish")
	<-done
	testlog.Debug("Waiting finished: Test completed")
}

func TestHandleServerShutdown(t *testing.T) {
	testlog.BeginTest()
	defer testlog.EndTest()

	testServer := &http.Server{}

	startError := make(chan error)
	interrupt := make(chan os.Signal)
	done := make(chan bool)

	t.Run("fatalStartError", func(t *testing.T) {
		provokedError := errors.New("listen tcp :8443: address already in use")

		go t.Run("handleServerStartError", func(t *testing.T) {
			exitCode, err := handleServerShutdown(testServer, startError, nil)

			t.Run("errorNotNil", func(t *testing.T) {
				assert.Error(t, err, "error should not be nil")
			})

			t.Run("equalErrorMessage", func(t *testing.T) {
				expectedError := errors.Wrap(provokedError, "error while starting server")

				assert.Equal(t, expectedError.Error(), err.Error(),
					"the returned error should be the one that was written in the startError channel")
			})

			t.Run("exitCode", func(t *testing.T) {
				assert.Equal(t, defaults.ExitStartError, exitCode, "exit code on start error should be 1")
			})

			testlog.Debug("Shutdown routine finished")
			done <- true
		})

		startError <- provokedError

		testlog.Debug("Waiting for start error tests to finish")
		<-done
		testlog.Debug("Waiting finished: Test completed")
	})

	t.Run("serverClosedError", func(t *testing.T) {
		go t.Run("handleStartErrorIsErrServerClosed", func(t *testing.T) {
			exitCode, err := handleServerShutdown(testServer, startError, nil)

			t.Run("errorNil", func(t *testing.T) {
				assert.NoError(t, err, "returned error should be nil")
			})

			t.Run("exitCode", func(t *testing.T) {
				assert.Equal(t, defaults.ExitSuccessful, exitCode, "exit code should be 0 because the start error is ErrServerClosed")
			})

			testlog.Debug("Shutdown routine finished")
			done <- true
		})

		startError <- http.ErrServerClosed

		testlog.Debug("Waiting for start error tests with ErrServerClosed to finish")
		<-done
		testlog.Debug("Waiting finished: Test completed")
	})

	t.Run("interrupt", func(t *testing.T) {
		go t.Run("handleInterruptSignal", func(t *testing.T) {
			exitCode, err := handleServerShutdown(testServer, startError, interrupt)

			t.Run("errorNil", func(t *testing.T) {
				assert.NoError(t, err, "returned error should be nil")
			})

			t.Run("exitCode", func(t *testing.T) {
				assert.Equal(t, defaults.ExitSuccessful, exitCode, "exit code should be 0 because the server was shutdown correctly")
			})

			testlog.Debug("Shutdown routine finished")
			done <- true
		})

		interrupt <- os.Interrupt
		startError <- http.ErrServerClosed

		testlog.Debug("Waiting for interrupt signal handling to finish")
		<-done
		testlog.Debug("Waiting finished: Test completed")
	})

	t.Run("kill", func(t *testing.T) {
		go t.Run("handleKillSignal", func(t *testing.T) {
			exitCode, err := handleServerShutdown(testServer, startError, interrupt)

			t.Run("errorNil", func(t *testing.T) {
				assert.NoError(t, err, "returned error should be nil")
			})

			t.Run("exitCode", func(t *testing.T) {
				assert.Equal(t, defaults.ExitShutdownError, exitCode, "exit code should be 2 because the server was not "+
					"shutdown correctly, interrupt signal should be used")
			})

			testlog.Debug("Shutdown routine finished")
			done <- true
		})

		interrupt <- os.Kill
		startError <- http.ErrServerClosed

		testlog.Debug("Waiting for kill signal handling to finish")
		<-done
		testlog.Debug("Waiting finished: Test completed")
	})

	t.Run("terminate", func(t *testing.T) {
		go t.Run("handleTerminateSignal", func(t *testing.T) {
			exitCode, err := handleServerShutdown(testServer, startError, interrupt)

			t.Run("errorNil", func(t *testing.T) {
				assert.NoError(t, err, "returned error should be nil")
			})

			t.Run("exitCode", func(t *testing.T) {
				assert.Equal(t, defaults.ExitShutdownError, exitCode, "exit code should be 2 because the server was not "+
					"shutdown correctly, interrupt signal should be used")
			})

			testlog.Debug("Shutdown routine finished")
			done <- true
			close(done)
		})

		interrupt <- syscall.SIGTERM
		startError <- http.ErrServerClosed

		testlog.Debug("Closing interrupt and start error channels")
		close(interrupt)
		close(startError)

		testlog.Debug("Waiting for terminate signal handling to finish")
		<-done
		testlog.Debug("Waiting finished: Test completed")
	})
}
