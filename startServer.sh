#!/usr/bin/env bash

cd cmd/ticketsystem

if [ "$USER" == "root" ]; then
    go run ticketsystem.go
else
    sudo GOPATH="$HOME/Documents/GoWorkspace" go run ticketsystem.go
fi

cd ../..
