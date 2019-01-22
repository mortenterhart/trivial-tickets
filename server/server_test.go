// Server starting and handler registration
package server

import (
	"net/http"
	"net/http/httptest"
	"os"
	"syscall"
	"testing"
	"time"

	"github.com/mortenterhart/trivial-tickets/globals"
	"github.com/mortenterhart/trivial-tickets/structs"
	"github.com/mortenterhart/trivial-tickets/util/filehandler"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
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

// TestGetTemplates makes sure the application is able to correctly find the templates
// with the given standard values
func TestGetTemplates(t *testing.T) {

	tmpl := GetTemplates("../www")
	tmplNil := GetTemplates("/www")

	assert.NotNil(t, tmpl, "GetTemplates() returned no found templates")
	assert.Nil(t, tmplNil, "GetTemplates() found templates where it was not supposed to be")
}

// TestRedirectToTLS tests the redirect to https, if a request with only http is made
func TestRedirectToTLS(t *testing.T) {

	config := mockConfig()
	globals.ServerConfig = &config

	req, _ := http.NewRequest("GET", "localhost", nil)
	w := httptest.NewRecorder()
	handler := http.HandlerFunc(redirectToTLS)

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusTemporaryRedirect, w.Code, "The HTTP status code was incorrect")
}

// TestRedirectToTLS tests the reditect to https, if a request with parameters with only http is made
func TestRedirectToTLSWithParams(t *testing.T) {

	req, _ := http.NewRequest("GET", "localhost?id=123", nil)
	w := httptest.NewRecorder()
	handler := http.HandlerFunc(redirectToTLS)

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusTemporaryRedirect, w.Code, "The HTTP status code was incorrect")
}

// TestStartHandlersNoPath is used to produce an error to make sure the function works properly
func TestStartHandlersNoPath(t *testing.T) {

	err := startHandlers("")

	assert.NotNil(t, err, "No error occurred, although the path was incorrect")
}

// TestStartServerNoUsersPath produces an error to make sure the server will not start without a path to the users.json file
func TestStartServerNoUsersPath(t *testing.T) {

	config := mockConfig()
	config.Users = ""

	exitCode, err := StartServer(&config)

	assert.NotNil(t, err, "No error was returned, although no users path was specified")
	assert.Equal(t, 1, exitCode, "exit code should be 1 due to expected error")
}

// TestStartServerNoUsersPath produces an error to make sure the server will not start without a path to the mail folder
func TestStartServerNoMailsPath(t *testing.T) {

	config := mockConfig()
	config.Mails = ""

	exitCode, err := StartServer(&config)

	assert.NotNil(t, err, "No error was returned, although no mails path was specified")
	assert.Equal(t, 1, exitCode, "exit code should be 1 due to expected error")
}

// TestStartServerNoUsersPath produces an error to make sure the server will not start without a path to the web files
func TestStartServerNoWebPath(t *testing.T) {

	config := mockConfig()
	config.Web = ""

	exitCode, err := StartServer(&config)

	assert.NotNil(t, err, "No error was returned, although no web path was specified")
	assert.Equal(t, 1, exitCode, "exit code should be 1 due to expected error")
}

// TestStartServerNoUsersPath produces an error to make sure the server will not start without a path to the ticket folder
func TestStartServerNoTicketsPath(t *testing.T) {

	config := mockConfig()
	config.Tickets = ""

	exitCode, err := StartServer(&config)

	assert.NotNil(t, err, "No error was returned, although no tickets path was specified")
	assert.Equal(t, 1, exitCode, "exit code should be 1 due to expected error")
}

func TestStartServerAllConfigsSet(t *testing.T) {

	config := mockConfig()
	shutdown := make(chan bool)

	go func() {
		exitCode, err := StartServer(&config)

		t.Run("errorNil", func(t *testing.T) {
			assert.NoError(t, err, "returned error should be nil because the server was shutdown correctly")
		})

		t.Run("exitCode", func(t *testing.T) {
			assert.Equal(t, 0, exitCode, "exit code should be 0 because the server was shutdown correctly")
		})

		shutdown <- true
	}()

	time.Sleep(2 * time.Second)

	signalErr := syscall.Kill(os.Getpid(), syscall.SIGINT)

	assert.NoError(t, signalErr, "sending interrupt signal should not cause an error")

	<-shutdown
}

func TestCreateResourceFolders(t *testing.T) {
	config := mockConfig()
	config.Tickets = "../files/testtickets"
	config.Mails = "../files/testmails"

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

	cleanupTestFiles()
}

func TestNotifyOnInterruptSignal(t *testing.T) {
	interrupt := notifyOnInterruptSignal()

	done := make(chan bool)

	go t.Run("catchSignal", func(t *testing.T) {
		t.Log("Waiting for signal")
		capturedSignal := <-interrupt

		t.Log("Caught signal SIGINT in Go routine")

		t.Run("caughtSignalNotNil", func(t *testing.T) {
			assert.NotNil(t, capturedSignal, "captured signal should not be nil")
		})

		t.Run("isInterruptSignal", func(t *testing.T) {
			assert.Equal(t, syscall.SIGINT, capturedSignal, "captured signal should be SIGINT signal")
		})

		done <- true
	})

	signalErr := syscall.Kill(os.Getpid(), syscall.SIGINT)

	t.Run("signalErrNil", func(t *testing.T) {
		assert.NoError(t, signalErr, "kill error should be nil")
	})

	t.Log("Waiting for signal tests done")
	<-done
}

func TestHandleServerShutdown(t *testing.T) {
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
				assert.Equal(t, 1, exitCode, "exit code on start error should be 1")
			})

			done <- true
		})

		startError <- provokedError

		t.Log("Waiting for start error tests to finish")
		<-done
	})

	t.Run("serverClosedError", func(t *testing.T) {
		go t.Run("handleStartErrorIsErrServerClosed", func(t *testing.T) {
			exitCode, err := handleServerShutdown(testServer, startError, nil)

			t.Run("errorNil", func(t *testing.T) {
				assert.NoError(t, err, "returned error should be nil")
			})

			t.Run("exitCode", func(t *testing.T) {
				assert.Equal(t, 0, exitCode, "exit code should be 0 because the start error is ErrServerClosed")
			})

			done <- true
		})

		startError <- http.ErrServerClosed

		t.Log("Waiting for start error tests with ErrServerClosed to finish")
		<-done
	})

	t.Run("interrupt", func(t *testing.T) {
		go t.Run("handleInterruptSignal", func(t *testing.T) {
			exitCode, err := handleServerShutdown(testServer, nil, interrupt)

			t.Run("errorNil", func(t *testing.T) {
				assert.NoError(t, err, "returned error should be nil")
			})

			t.Run("exitCode", func(t *testing.T) {
				assert.Equal(t, 0, exitCode, "exit code should be 0 because the server was shutdown correctly")
			})

			done <- true
		})

		interrupt <- os.Interrupt

		t.Log("Waiting for interrupt signal handling to finish")
		<-done
	})

	t.Run("kill", func(t *testing.T) {
		go t.Run("handleKillSignal", func(t *testing.T) {
			exitCode, err := handleServerShutdown(testServer, nil, interrupt)

			t.Run("errorNil", func(t *testing.T) {
				assert.NoError(t, err, "returned error should be nil")
			})

			t.Run("exitCode", func(t *testing.T) {
				assert.Equal(t, 1, exitCode, "exit code should be 1 because the server was not "+
					"shutdown correctly, interrupt signal should be used")
			})

			done <- true
		})

		interrupt <- os.Kill

		t.Log("Waiting for kill signal handling to finish")
		<-done
	})

	t.Run("terminate", func(t *testing.T) {
		go t.Run("handleTerminateSignal", func(t *testing.T) {
			exitCode, err := handleServerShutdown(testServer, nil, interrupt)

			t.Run("errorNil", func(t *testing.T) {
				assert.NoError(t, err, "returned error should be nil")
			})

			t.Run("exitCode", func(t *testing.T) {
				assert.Equal(t, 1, exitCode, "exit code should be 1 because the server was not "+
					"shutdown correctly, interrupt signal should be used")
			})

			done <- true
		})

		interrupt <- syscall.SIGTERM

		t.Log("Waiting for terminate signal handling to finish")
		<-done
	})
}

// Utility function to create a mock configuration for the server
func mockConfig() structs.Config {

	return structs.Config{
		Port:    8443,
		Tickets: "../files/tickets",
		Mails:   "../files/mails",
		Users:   "../files/users/users.json",
		Cert:    "../ssl/server.cert",
		Key:     "../ssl/server.key",
		Web:     "../www",
	}
}

func mockLogConfig() structs.LogConfig {
	return structs.LogConfig{
		LogLevel:   structs.LevelInfo,
		VerboseLog: false,
		FullPaths:  false,
	}
}
