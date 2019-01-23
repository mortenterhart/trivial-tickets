#!/usr/bin/env bash
##
## Trivial Tickets Ticketsystem
## Copyright (C) 2019 The Contributors
##
## This program is free software: you can redistribute it and/or modify
## it under the terms of the GNU General Public License as published by
## the Free Software Foundation, either version 3 of the License, or
## (at your option) any later version.
##
## This program is distributed in the hope that it will be useful,
## but WITHOUT ANY WARRANTY; without even the implied warranty of
## MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
## GNU General Public License for more details.
##
## You should have received a copy of the GNU General Public License
## along with this program.  If not, see <http://www.gnu.org/licenses/>.
##
##
## Ticketsystem Trivial Tickets
##
## Matriculation numbers: 3040018, 6694964, 3478222
## Lecture:               Programmieren II, INF16B
## Lecturer:              Herr Prof. Dr. Helmut Neemann
## Institute:             Duale Hochschule Baden-Württemberg Mosbach
##
## ---------------
## Webserver Start script
##
## Get more information about this start script in
## the project README or here:
## https://github.com/mortenterhart/trivial-tickets/wiki/Build-and-Execution
##

program_name="trivial-tickets"
root_dir="${0%/*}"
main_executable="cmd/ticketsystem/ticketsystem.go"

exit_status=0

# info prints the program name and the message
# denoting an informative message (not an error)
# to stdout. The message is built by concatenating
# all arguments.
function info() {
    printf "[%s] INFO: %s\n" "${program_name}" "$*"
}

# error prints an error message indicating an
# occurred error to stdout. The message is built
# by concatenating all arguments.
function error() {
    printf "[%s] ERROR: %s\n" "${program_name}" "$*"
}

# Check if any Go environment is installed and report
# an error if Go is not found
if ! type -P "gol" > /dev/null 2>&1; then
    error "The 'go' command could not be found."
    error "Check your Go installation and make sure it appears on your \$PATH."
    exit 1
fi

# Check if the ticketsystem main executable exists
if [ -f "${root_dir}/${main_executable}" ]; then
    # Change into the repository root folder so
    # that the relative paths from the default
    # config match
    cd "${root_dir}"

    # Download the required productive and test
    # dependencies
    info "Checking for missing Go dependencies"
    go get -t -v ./...

    if [ "$?" -eq 0 ]; then

        # Execute the ticketsystem executable
        info "Starting Ticketsystem web server"
        go run "${main_executable}" "$@"
        exit_status=$?
    else
        error "Failed to install the missing Go dependencies"
        exit_status=2
    fi

    # Change back to the previous directory
    cd "${OLDPWD}"
else
    error "Cannot find the Ticketsystem main executable '${main_executable}'"
    exit_status=1
fi

exit "${exit_status}"
