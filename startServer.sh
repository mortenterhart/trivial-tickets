#!/usr/bin/env bash

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
        go run ticketsystem.go

        cd ../..
    fi
else
    error "cannot find the main executable cmd/ticketsystem/ticketsystem.go"
    echo "You might be in the wrong working directory, execute"
    echo '  cd "${0%/*}"'
    echo "to change to the correct directory."

    exit 1
fi

exit 0
