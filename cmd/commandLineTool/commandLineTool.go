package main

import (
    "github.com/mortenterhart/trivial-tickets/logic"
)

func main() {
    logic.MainLoop()
}

// Here's how it should behave:
// On startup, show the available commands (fetch, submit and exit)
// On fetch just print out all available emails in the command line.
// ON submit ask in order for email, subject and message.
// After both fetch and submit, return to startup message.
// On exit, close the application.
