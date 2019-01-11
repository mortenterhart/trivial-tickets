package structs

import (
	"time"
)

/*
*
* Matrikelnummern
* 3040018
* 6694964
* 3478222
 */

// Config is a struct to hold the config parameters provided on startup
type Config struct {
	Port    int16
	Tickets string
	Users   string
	Cert    string
	Key     string
	Web     string
}

// CLIConfig is a struct to hold the CLIs config parameters provided on startup
type CLIConfig struct {
	IPAddr string
	Port   uint16
	Cert   string
}

// Session is a struct that holds session variables for a certain user
type Session struct {
	User       User
	CreateTime time.Time
	IsLoggedIn bool
	Id         string
}

// SessionManager holds a session and operates on it
type SessionManager struct {
	Name    string
	Session Session
	TTL     int64
}

// User is the model for a user that works on tickets
type User struct {
	Id          string `json:"Id"`
	Name        string `json:"Name"`
	Username    string `json:"Username"`
	Mail        string `json:"Email"`
	Hash        string `json:"Hash"`
	IsOnHoliday bool   `json:"IsOnHoliday"`
}

// Data holds session and ticket data to parse to the web templates
type Data struct {
	Session Session
	Tickets map[string]Ticket
	Users   map[string]User
}

// DataSingleTicket holds the session and ticket data for a call to a single ticket
type DataSingleTicket struct {
	Session Session
	Ticket  Ticket
	Tickets map[string]Ticket
	Users   map[string]User
}

// Ticket represents a ticket
type Ticket struct {
	Id       string  `json:"Id"`
	Subject  string  `json:"Subject"`
	Status   State   `json:"Status"`
	User     User    `json:"User"`
	Customer string  `json:"Customer"`
	Entries  []Entry `json:"Entries"`
	MergeTo  string  `json:"MergeTo"`
}

// Entry describes a single reply within a ticket
type Entry struct {
	Date          time.Time
	FormattedDate string
	User          string
	Text          string
	Reply_Type    string
}

// State is an enum to represent the current status of a ticket
type State int

const (
	OPEN State = iota
	PROCESSING
	CLOSED
)

// Email struct holds the information for a received email in order
// to create new tickets or answers
type Mail struct {
	Id      string `json:"id"`
	From    string `json:"from"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	Message string `json:"message"`
}

type Command int

const (
	FETCH Command = iota
	SUBMIT
	EXIT
)
