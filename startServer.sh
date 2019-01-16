#!/usr/bin/env bash
#
# Ticketsystem Trivial Tickets
#
# Matriculation numbers: 3040018, 6694964, 3478222
# Lecture:               Programmieren II, INF16B
# Lecturer:              Herr Prof. Dr. Helmut Neemann
# Institute:             Duale Hochschule Baden-Württemberg Mosbach
#
# ---------------
# Webserver Start script

server_name="trivial-tickets"

function info() {
    printf "[%s] INFO %s\n" "${server_name}" "$*"
}

function error() {
    printf "[%s] ERROR %s\n" "${server_name}" "$*"
}

if [ -d "cmd/ticketsystem" ]; then
    info "Checking for missing Go dependencies"
    go get -t -v ./...

    if [ "$?" -eq 0 ]; then
        cd cmd/ticketsystem

        info "Starting Ticketsystem webserver"
        go run ticketsystem.go "$@"
        exit_status=$?

        cd ../..
    fi
else
    error "cannot find the main executable cmd/ticketsystem/ticketsystem.go"
    echo "You might be in the wrong working directory, execute"
    echo "  cd \"${0%/*}\""
    echo "to change to the correct directory."

    exit 1
fi

exit "${exit_status}"
