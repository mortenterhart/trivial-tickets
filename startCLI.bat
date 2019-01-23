@echo off

REM  Skip the comments to speed up execution
GOTO startCLI

REM  Trivial Tickets Ticketsystem
REM  Copyright (C) 2019 The Contributors
REM
REM  This program is free software: you can redistribute it and/or modify
REM  it under the terms of the GNU General Public License as published by
REM  the Free Software Foundation, either version 3 of the License, or
REM  (at your option) any later version.
REM
REM  This program is distributed in the hope that it will be useful,
REM  but WITHOUT ANY WARRANTY; without even the implied warranty of
REM  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
REM  GNU General Public License for more details.
REM
REM  You should have received a copy of the GNU General Public License
REM  along with this program.  If not, see <http://www.gnu.org/licenses/>.
REM
REM
REM  Ticketsystem Trivial Tickets
REM
REM  Matriculation numbers: 3040018, 6694964, 3478222
REM  Lecture:               Programmieren II, INF16B
REM  Lecturer:              Herr Prof. Dr. Helmut Neemann
REM  Institute:             Duale Hochschule Baden-Württemberg Mosbach
REM
REM  ---------------
REM  CLI Start script
REM
REM  Get more information about this start script in
REM  the project README or here:
REM  https://github.com/mortenterhart/trivial-tickets/wiki/Build-and-Execution
REM

REM  Here the actual script starts
:startCLI

REM  Start localization of environment changes from this
REM  batch script and enable command extensions if
REM  available on that system.
SETLOCAL ENABLEEXTENSIONS
IF ERRORLEVEL 1 (
    call :error "Failed to enable command extensions: not available on your system"
    exit /B 1
)

REM  Script constants
set program_name=trivial-tickets CLI
set root_dir=%~dp0
set main_executable=cmd\command_line_tool\command_line_tool.go

REM  The exit status of this script
set exit_status=0

REM  Check if any Go environment is installed and
REM  report an error if Go is not found
WHERE /Q go
IF ERRORLEVEL 1 (
    call :error "The 'go' command could not be found."
    call :error "Check your Go installation and make sure it appears on your %%%%PATH%%%%."
    exit /B 1
)

REM  Check if the main executable of the command-line
REM  tool exists
IF EXIST "%root_dir%\%main_executable%" (
    REM  Save the current working directory to
    REM  return here later
    set OLDPWD=%CD%

    REM  Change into the repository root folder so
    REM  that the relative paths from the default
    REM  config match
    cd %root_dir%

    REM  Download the required productive and test
    REM  dependencies
    call :info "Checking for missing Go dependencies"
    go get -t -v ./...

    IF %ERRORLEVEL% EQU 0 (

        REM  Execute the command-line tool executable
        call :info "Starting command line interface"
        go run %main_executable% %*
        set exit_status=%ERRORLEVEL%
    ) ELSE (
        call :error "Failed to install the missing Go dependencies"
        set exit_status=2
    )

    REM  Change back to the previous directory
    cd %OLDPWD%
) ELSE (
    call :error "Cannot find the Ticketsystem main executable %main_executable%"
    set exit_status=1
)

REM  End local environment
ENDLOCAL

REM  Exit the script with the exit status
exit /B %exit_status%

REM  info prints the program name and the message
REM  denoting an informative message (not an error)
REM  to stdout. The message is passed as first
REM  argument $1 and should be quoted with double
REM  quotes.
:info
    echo [%program_name%] INFO: %~1
exit /B 0

REM  error prints an error message indicating an
REM  occurred error to stdout. The message is passed
REM  as first argument $1 and should be surrounded
REM  by double quotes.
:error
    IF "%program_name%" == "" set program_name=trivial-tickets CLI
    echo [%program_name%] ERROR: %~1 >&2
exit /B 0
