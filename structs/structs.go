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

// Session is a struct that holds session variables for a certain user
type Session struct {
	IsLoggedIn bool
	User       User
	Time       time.Time
}

// User is the model for a user that works on tickets
type User struct {
	Id          int32  `json:"Id"`
	Name        string `json:"Name"`
	Mail        string `json:"Mail"`
	Hash        string `json:"Hash"`
	IsOnHoliday bool   `json:"IsOnHoliday"`
}

// Ticket represents a ticket
type Ticket struct {
	Id       int32   `json:"Id"`
	Subject  string  `json:"Subject"`
	Status   state   `json:"Status"`
	User     User    `json:"User"`
	Customer string  `json:"Customer"`
	Entries  []Entry `json:"Entries"`
}

// Entry describes a single reply within a ticket
type Entry struct {
	Date time.Time
	User string
	Text string
}

// State is an enum to represent the current status of a ticket
type state int

const (
	OPEN state = iota
	PROCESSING
	CLOSED
)

type Status interface {
	Status() state
}
