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
# CLI Start script

program_name="CLI"

function info() {
    printf "[%s] INFO %s\n" "${program_name}" "$*"
}

function error() {
    printf "[%s] ERROR %s\n" "${program_name}" "$*"
}

if [ -d "cmd/commandLineTool" ]; then
    info "Checking for missing Go dependencies"
    go get -t -v ./...

    if [ "$?" -eq 0 ]; then

        info "Starting Ticketsystem webserver"
        go run cmd/commandLineTool/commandLineTool.go

    fi
else
    error "cannot find the main executable cmd/commandLineTool/commandLineTool.go"
    echo "You might be in the wrong working directory, execute"
    echo '  cd "${0%/*}"'
    echo "to change to the correct directory."

    exit 1
fi

exit 0
